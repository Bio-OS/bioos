package workflow

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/Bio-OS/bioos/internal/context/workspace/application/query/workflow"
	"github.com/Bio-OS/bioos/internal/context/workspace/domain/workspace"
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/workflowparser"
	"github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	applog "github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/schema"
	"github.com/Bio-OS/bioos/pkg/utils/git"
	"github.com/Bio-OS/bioos/pkg/validator"
)

type WorkflowVersionAddedHandler struct {
	repo Repository
}

func NewWorkflowVersionAddedHandler(repo Repository) *WorkflowVersionAddedHandler {
	return &WorkflowVersionAddedHandler{
		repo: repo,
	}
}

func (h *WorkflowVersionAddedHandler) Handle(ctx context.Context, event *WorkflowVersionAddedEvent) (err error) {
	// get workflow
	workflow, err := h.repo.Get(ctx, event.WorkspaceID, event.WorkflowID)
	if err != nil {
		return err
	}
	version, exist := workflow.Versions[event.WorkflowVersionID]
	if !exist {
		return proto.ErrorWorkflowVersionNotFound("workspace:%s workflow:%s version:%s not found", event.WorkspaceID, event.WorkflowID, event.WorkflowVersionID)
	}
	defer func() {
		if err == nil {
			version.Status = WorkflowVersionSuccessStatus
			version.Message = "success"
		} else {
			version.Status = WorkflowVersionFailedStatus
			version.Message = err.Error()
		}
		applog.Infow("start to save workflow", "workflowID", workflow.ID, "workflowVersionID", version.ID, "status", version.Status, "err", err)
		// save workflow
		if err := h.repo.Save(ctx, workflow); err != nil {
			applog.Errorw("fail to save workflow version", "workflowVersion", version.ID, "err", err)
		}
	}()

	switch version.Status {
	case WorkflowVersionSuccessStatus:
		return nil
	case WorkflowVersionFailedStatus, WorkflowVersionPendingStatus:
		// TODO when and how to handle fail status when retrying
		if err := h.handle(ctx, workflow.ID, version, event); err != nil {
			return err
		}
	}

	return nil
}
func (h *WorkflowVersionAddedHandler) handle(ctx context.Context, workflowID string, version *WorkflowVersion, event *WorkflowVersionAddedEvent) error {
	var dir string
	var err error
	workflowParser, err := workflowparser.GetFactory().GetParser(version.Language)
	if err != nil {
		applog.Errorf("wrong workflow parser: %s", err)
		return err
	}

	if event.FilesBaseDir != "" {
		dir = event.FilesBaseDir
	} else {
		applog.Infow("start to git clone", "metadata", version.Metadata, "source", version.Source)

		dir, err = os.MkdirTemp("", fmt.Sprintf("workflow-%s-%s-", workflowID, version.ID))
		if err != nil {
			applog.Errorf("fail to make tmp dir: %s", err)
			return err
		}
		// clean workflow files if source is git
		defer os.RemoveAll(dir)

		// step1: clone workflow
		if err := git.Clone(dir, event.GitRepo, event.GitToken, event.GitTag); err != nil {
			applog.Errorw("fail to clone", "err", err)
			return err
		}
	}

	// step2: validate main workflow path exist
	mainWorkflowPath := path.Join(dir, version.MainWorkflowPath)
	if _, err = os.Stat(mainWorkflowPath); err != nil {
		if os.IsNotExist(err) {
			return proto.ErrorWorkflowMainWorkflowFileNotExist("main workflow file not exist")
		}
		return apperrors.NewInternalError(err)
	}
	// parse workfile version
	languageVersion, err := workflowParser.ParseWorkflowVersion(ctx, mainWorkflowPath)
	if err != nil {
		return apperrors.NewInternalError(err)
	}
	version.LanguageVersion = languageVersion

	// step3: validate and save workflow files
	fileParamsStr, err := workflowParser.ValidateWorkflowFiles(ctx, dir, version.MainWorkflowPath)
	if err != nil {
		return err
	}
	fileParams, err := fileParamPOToDO(fileParamsStr)
	if err != nil {
		return err
	}
	for _, value := range fileParams {
		valueCopy := value
		workflowFile, err := version.AddFile(&valueCopy)
		if err != nil {
			return err
		}
		applog.Infow("success add workflow file", "workflowVersionID", version.ID, "fileID", workflowFile.ID, "path", workflowFile.Path)
	}

	// step4: get workflow inputs
	inputsStr, err := workflowParser.GetWorkflowInputs(ctx, mainWorkflowPath)
	if err != nil {
		return err
	}
	inputs, err := workflowParamPOToDO(inputsStr)
	if err != nil {
		return err
	}
	version.Inputs = inputs

	// step5: get workflow outputs
	outputsStr, err := workflowParser.GetWorkflowOutputs(ctx, mainWorkflowPath)
	if err != nil {
		return err
	}
	outputs, err := workflowParamPOToDO(outputsStr)
	if err != nil {
		return err
	}
	version.Outputs = outputs

	// step6: get workflow graph
	graph, err := workflowParser.GetWorkflowGraph(ctx, mainWorkflowPath)
	if err != nil {
		return err
	}
	version.Graph = graph
	return nil
}

type WorkspaceDeletedHandler struct {
	repo Repository
}

func NewWorkspaceDeletedHandler(repo Repository) *WorkspaceDeletedHandler {
	return &WorkspaceDeletedHandler{
		repo: repo,
	}
}

func (h *WorkspaceDeletedHandler) Handle(ctx context.Context, event *workspace.WorkspaceEvent) (err error) {
	if event == nil {
		return nil
	}
	applog.Infow("start to gc workflow", "workspace", event.WorkspaceID)
	workflowIDs, err := h.repo.List(ctx, event.WorkspaceID)
	if err != nil {
		applog.Errorw("fail to list workflows", "err", err, "workspace", event.WorkspaceID)
		return apperrors.NewInternalError(err)
	}
	for _, workflowID := range workflowIDs {
		applog.Infow("start to delete workflow", "workspace", event.WorkspaceID, "workflow", workflowID)
		if err := h.repo.Delete(ctx, &Workflow{ID: workflowID}); err != nil {
			applog.Errorw("fail to delete workflow", "err", err, "workspace", event.WorkspaceID, "workflow", workflowID)
			return apperrors.NewInternalError(err)
		}
	}
	return nil
}

type ImportWorkflowsHandler struct {
	readModel workflow.ReadModel
	repo      Repository
	eventbus  eventbus.EventBus
	factory   *Factory
}

func NewImportWorkflowsHandler(repo Repository, readModel workflow.ReadModel, bus eventbus.EventBus, factory *Factory) *ImportWorkflowsHandler {
	return &ImportWorkflowsHandler{
		readModel: readModel,
		repo:      repo,
		eventbus:  bus,
		factory:   factory,
	}
}

func (h *ImportWorkflowsHandler) Handle(ctx context.Context, event *ImportWorkflowsEvent) error {
	workflowSet := sets.New[string]()
	workflowVersionSet := sets.New[string]()
	for _, workflow := range event.Schemas {
		//workflowPath := strings.ReplaceAll(workflow.Path, " ", "")
		if workflowSet.Has(workflow.Name) {
			return fmt.Errorf("workflow name[%s] is not unique ", workflow.Name)
		}
		workflowSet.Insert(workflow.Name)
		if err := validateWorkflow(workflow); err != nil {
			return err
		}
		newWorkflow, err := h.factory.NewWorkflow(event.WorkspaceID,
			&WorkflowOption{
				Name:        workflow.Name,
				Description: workflow.Description,
			},
		)
		if err != nil {
			return err
		}
		versionOption := &VersionOption{
			Language:         workflow.Language,
			MainWorkflowPath: workflow.MainWorkflowPath,
			Source:           WorkflowSourceFile,
			Url:              fmt.Sprintf("%s://%s", workflow.Metadata.Scheme, workflow.Metadata.Repo),
			Tag:              workflow.Metadata.Tag,
		}
		if workflow.Metadata.Token != nil {
			versionOption.Token = *workflow.Metadata.Token
		}

		version, err := newWorkflow.AddVersion(versionOption)
		if err != nil {
			return err
		}
		workflowVersionSet.Insert(version.ID)

		versionAddedEvent := NewWorkflowVersionAddedEvent(event.WorkspaceID, newWorkflow.ID, version.ID, versionOption.Url, versionOption.Tag, versionOption.Token, path.Join(event.ImportFileBaseDir, workflow.Path))

		if err = h.repo.Save(ctx, newWorkflow); err != nil {
			return err
		}

		if err = h.eventbus.Publish(ctx, versionAddedEvent); err != nil {
			return err
		}
	}
	// check if all version has been imported
	deadline := time.Now().Add(10 * time.Minute)
	for {
		time.Sleep(1 * time.Second)
		if time.Now().After(deadline) {
			return fmt.Errorf("importing workflows timeout")
		}
		versionList := workflowVersionSet.UnsortedList()
		for _, id := range versionList {
			ver, err := h.readModel.GetVersion(ctx, id)
			//not return error because we don't want to use retry mechanism in eventbus
			if err != nil {
				applog.Errorw("fail to get workflow version", "err", err, "version", id)
				continue
			}
			if ver != nil && ver.Status == WorkflowVersionSuccessStatus {
				workflowVersionSet.Delete(id)
			}
		}

		if workflowVersionSet.Len() == 0 {
			break
		}
	}

	for _, workflow := range event.Schemas {
		err := os.Remove(path.Join(event.ImportFileBaseDir, workflow.Path))
		if err != nil {
			//remove file error should not lead to import fail
			applog.Errorf("remove file failed: %w", err)
		}
	}

	//importedEvent := NewWorkflowsImportedEvent(event.WorkspaceID, event.ImportFileBaseDir)
	//err := h.eventbus.Publish(ctx, importedEvent)
	//if err != nil {
	//	return err
	//}
	return nil
}

func validateWorkflow(workflow schema.WorkflowTypedSchema) error {
	if !validator.ValidateResNameInString(workflow.Name) {
		return fmt.Errorf("workflow name[%s] not passed the validation ", workflow.Name)
	}
	if workflow.Language != "WDL" {
		return fmt.Errorf("workflow language [%s] not passed the validation ", workflow.Language)
	}
	//TODO Validation will be consistent with that of commercial version in the future
	return nil
}
