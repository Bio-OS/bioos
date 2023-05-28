package notebook

import (
	"context"
	"fmt"
	"os"
	"path"

	query "github.com/Bio-OS/bioos/internal/context/workspace/application/query/notebook"
	"github.com/Bio-OS/bioos/internal/context/workspace/domain/notebook"
	"github.com/Bio-OS/bioos/internal/context/workspace/domain/workspace"
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
	"github.com/Bio-OS/bioos/pkg/log"
)

func addEventHandle(eb eventbus.EventBus, readModel query.ReadModel, svc notebook.Service, factory *notebook.Factory) {
	eb.Subscribe(workspace.WorkspaceDeleted, &workspaceDeleteHandler{
		service:   svc,
		readModel: readModel,
	})
	eb.Subscribe(notebook.ImportNotebooks, &importNotebooksHandler{
		service:  svc,
		eventbus: eb,
		factory:  factory,
	})
}

type workspaceDeleteHandler struct {
	service   notebook.Service
	readModel query.ReadModel
}

func (h *workspaceDeleteHandler) Handle(ctx context.Context, payload string) error {
	log.Infow("start to consume workspace deleted event", "payload", payload)
	event, err := workspace.NewWorkspaceEventFromPayload([]byte(payload))
	if err != nil {
		return fmt.Errorf("decode event payload fail: %w", err)
	}
	list, err := h.readModel.ListByWorkspace(ctx, event.WorkspaceID)
	if err != nil {
		return fmt.Errorf("list workspace %s notebook fail: %w", event.WorkspaceID, err)
	}
	for _, n := range list {
		if err = h.service.Delete(ctx, notebook.Path(n.WorkspaceID, n.Name)); err != nil {
			return fmt.Errorf("delete notebook %s/%s fail: %w", n.WorkspaceID, n.Name, err)
		}
	}
	return nil
}

type importNotebooksHandler struct {
	service  notebook.Service
	eventbus eventbus.EventBus
	factory  *notebook.Factory
}

func (h *importNotebooksHandler) Handle(ctx context.Context, payload string) error {
	log.Infow("start to consume import notebooks event", "payload", payload)
	event, err := notebook.NewImportNotebooksEventFromPayload([]byte(payload))
	if err != nil {
		return fmt.Errorf("decode event payload fail: %w", err)
	}

	for _, nb := range event.Schema.Artifacts {
		bytes, err := os.ReadFile(path.Join(event.ImportFileBaseDir, nb.Path))
		if err != nil {
			return err
		}
		newNotebook, err := h.factory.New(&notebook.CreateParam{
			Name:        nb.Name,
			WorkspaceID: event.WorkspaceID,
			Content:     bytes,
		})

		if err != nil {
			return err
		}
		if err := h.service.Upsert(ctx, newNotebook); err != nil {
			return err
		}
	}

	// clean files
	for _, notebook := range event.Schema.Artifacts {
		err = os.Remove(path.Join(event.ImportFileBaseDir, notebook.Path))
		if err != nil {
			//remove file error should not lead to import fail
			log.Errorf("remove file failed: %s", err.Error())
		}
	}
	return nil
}
