package notebook

import (
	"context"

	"github.com/Bio-OS/bioos/internal/context/workspace/application/query/workspace"
	"github.com/Bio-OS/bioos/internal/context/workspace/domain/notebook"
)

type DeleteCommand struct {
	Name        string
	WorkspaceID string
}

// DeleteHandler ...
type DeleteHandler interface {
	Handle(context.Context, *DeleteCommand) error
}

// NewDeleteHandler ...
func NewDeleteHandler(svc notebook.Service, workspaceReadModel workspace.WorkspaceReadModel) DeleteHandler {
	return &deleteHandler{
		factory:            notebook.NewFactory(),
		service:            svc,
		workspaceReadModel: workspaceReadModel,
	}
}

type deleteHandler struct {
	factory            *notebook.Factory
	service            notebook.Service
	workspaceReadModel workspace.WorkspaceReadModel
}

func (h *deleteHandler) Handle(ctx context.Context, cmd *DeleteCommand) error {
	if err := workspace.CheckWorkspaceExist(ctx, h.workspaceReadModel, cmd.WorkspaceID); err != nil {
		return err
	}
	if err := h.service.Delete(ctx, notebook.Path(cmd.WorkspaceID, cmd.Name)); err != nil {
		return err
	}
	return nil
}
