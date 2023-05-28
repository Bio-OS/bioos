package workspace

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	"gopkg.in/yaml.v3"

	notebookserver "github.com/Bio-OS/bioos/internal/context/notebookserver/domain"
	datamodel "github.com/Bio-OS/bioos/internal/context/workspace/domain/data-model"
	"github.com/Bio-OS/bioos/internal/context/workspace/domain/notebook"
	"github.com/Bio-OS/bioos/internal/context/workspace/domain/workflow"
	"github.com/Bio-OS/bioos/internal/context/workspace/domain/workspace"
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
	"github.com/Bio-OS/bioos/pkg/consts"
	applog "github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/schema"
)

func addEventHandle(eb eventbus.EventBus,
	workspaceRepo workspace.Repository,
	eventRepo eventbus.EventRepository,
	workspaceFactory *workspace.Factory,

) {
	eb.Subscribe(workspace.ImportWorkspace, &importWorkspaceHandler{
		repo:     workspaceRepo,
		eventbus: eb,
		factory:  workspaceFactory,
	})
	eb.Subscribe(workspace.WorkspaceImported, &workspaceImportedHandler{
		repo:      workspaceRepo,
		eventbus:  eb,
		eventRepo: eventRepo,
	})
}

type importWorkspaceHandler struct {
	repo     workspace.Repository
	eventbus eventbus.EventBus
	factory  *workspace.Factory
}

func (h *importWorkspaceHandler) Handle(ctx context.Context, payload string) (err error) {
	applog.Infow("start to consume import workspace event", "payload", payload)

	event, err := workspace.NewImportWorkspaceEventFromPayload([]byte(payload))
	if err != nil {
		return err
	}
	baseDir := fmt.Sprintf(path.Join(event.Storage.NFS.MountPath, event.WorkspaceID))
	bytes, err := os.ReadFile(path.Join(baseDir, consts.WorkspaceYAMLName))
	if err != nil {
		return err
	}
	schema := &schema.WorkspaceTypedSchema{}
	err = yaml.Unmarshal(bytes, &schema)
	if err != nil {
		return err
	}
	err = h.createWorkspace(ctx, schema, event)
	if err != nil {
		return err
	}
	err = h.eventbus.Publish(ctx, workspace.NewWorkspaceImportedEvent(event.WorkspaceID, baseDir))
	if err != nil {
		return err
	}
	err = h.publishNotebooksEvent(ctx, event.WorkspaceID, baseDir, schema.Notebooks)
	if err != nil {
		return err
	}
	err = h.publishWorkflowsEvent(ctx, event.WorkspaceID, baseDir, schema.Workflows)
	if err != nil {
		return err
	}
	err = h.publishDataModelsEvent(ctx, event.WorkspaceID, baseDir, schema.DataModels)
	if err != nil {
		return err
	}
	err = h.publishNotebookServersEvent(ctx, event.WorkspaceID, baseDir, schema.Notebooks)
	if err != nil {
		return err
	}
	return nil
}

func (h *importWorkspaceHandler) createWorkspace(ctx context.Context, schema *schema.WorkspaceTypedSchema, event *workspace.ImportWorkspaceEvent) error {
	param := workspace.CreateWorkspaceParam{
		ID:          event.WorkspaceID,
		Name:        schema.Name,
		Description: schema.Description,
	}
	if event.Storage.NFS != nil {
		param.Storage.NFS = &workspace.NFSStorage{MountPath: event.Storage.NFS.MountPath}
	}
	ws, err := h.factory.CreateWithWorkspaceParam(param)
	if err != nil {
		return err
	}
	if err := h.repo.Save(ctx, ws); err != nil {
		return err
	}

	//Be consistent with that of directly create workspace
	applog.Infow("publish workspace created event", "ID", event.WorkspaceID)
	return nil
}

func (h *importWorkspaceHandler) publishWorkflowsEvent(ctx context.Context, workspaceID, baseDir string, schemas []schema.WorkflowTypedSchema) error {
	importWorkflowsEvent := workflow.NewImportWorkflowsEvent(workspaceID, baseDir, schemas)
	if err := h.eventbus.Publish(ctx, importWorkflowsEvent); err != nil {
		return err
	}
	return nil
}

func (h *importWorkspaceHandler) publishNotebooksEvent(ctx context.Context, workspaceID, baseDir string, schemas schema.NotebookTypedSchema) error {
	importNotebooksEvent := notebook.NewImportNotebooksEvent(workspaceID, baseDir, schemas)
	if err := h.eventbus.Publish(ctx, importNotebooksEvent); err != nil {
		return err
	}
	return nil
}

func (h *importWorkspaceHandler) publishDataModelsEvent(ctx context.Context, workspaceID, baseDir string, schemas []schema.DataModelTypedSchema) error {
	importDataModelsEvent := datamodel.NewImportDataModelsEvent(workspaceID, baseDir, schemas)
	if err := h.eventbus.Publish(ctx, importDataModelsEvent); err != nil {
		return err
	}
	return nil
}

func (h *importWorkspaceHandler) publishNotebookServersEvent(ctx context.Context, workspaceID, baseDir string, schemas schema.NotebookTypedSchema) error {
	importNotebooksEvent := notebookserver.NewImportNotebookServersEvent(workspaceID, baseDir, schemas)
	if err := h.eventbus.Publish(ctx, importNotebooksEvent); err != nil {
		return err
	}
	return nil
}

func (h *workspaceImportedHandler) checkImportProcess(ctx context.Context, types []string, workspaceID string) (completed bool, hasFailed bool, err error) {
	events, err := h.eventRepo.Search(ctx, &eventbus.Filter{
		Type:    types,
		Status:  []string{eventbus.EventStatusCompleted, eventbus.EventStatusFailed},
		Payload: workspaceID,
	})
	if err != nil {
		return false, false, err
	}
	completedNum := 0
	//workspaceID is unique in import event among all events, thus we will only get one corresponding event each types
	for _, event := range events {
		//TODO consider the situation that the event may mark failed in the future
		// 1) short task and outdated
		// 2) long time task
		if event.Status == eventbus.EventStatusFailed {
			return false, true, nil
		}
		if event.Status == eventbus.EventStatusCompleted {
			completedNum += 1
			if completedNum == len(types) {
				return true, false, nil
			}
		}
	}
	return false, false, nil
}

type workspaceImportedHandler struct {
	repo      workspace.Repository
	eventRepo eventbus.EventRepository
	eventbus  eventbus.EventBus
}

func (h *workspaceImportedHandler) Handle(ctx context.Context, payload string) error {
	applog.Infow("start to consume workspace imported event", "payload", payload)

	event, err := workspace.NewWorkspaceImportedEventFromPayload([]byte(payload))
	if err != nil {
		return err
	}
	defer os.RemoveAll(event.ImportBaseDir)

	deadline := time.Now().Add(30 * time.Minute)
	for {
		time.Sleep(1 * time.Second)
		// prevent the dead cycle
		if time.Now().After(deadline) {
			return fmt.Errorf("importing workspace timeout")
		}
		completed, failed, err := h.checkImportProcess(ctx, []string{
			notebook.ImportNotebooks,
			workflow.ImportWorkflows,
			notebookserver.ImportNotebookServers,
			datamodel.ImportDataModels,
		}, event.WorkspaceID)
		//not return error because we don't want to use retry mechanism in eventbus
		if err != nil {
			applog.Errorw("fail to check event status", "err", err)
			continue
		}

		if completed {
			applog.Infof("%s importing completed", event.WorkspaceID)
			break
		} else if failed {
			// clean workspace if import process failed
			ws, err := h.repo.Get(ctx, event.WorkspaceID)
			if err != nil {
				applog.Errorw("fail to get workspace while deleting workspace", "err", err, "workspace", event.WorkspaceID)
			}
			err = h.repo.Delete(ctx, ws)
			if err != nil {
				applog.Errorw("fail to delete workspace", "err", err, "workspace", event.WorkspaceID)
			}
			deletedEvent := workspace.NewWorkspaceDeletedEvent(ws.ID)

			if eventErr := h.eventbus.Publish(ctx, deletedEvent); eventErr != nil {
				return eventErr
			}
			break

		}
	}
	return nil
}
