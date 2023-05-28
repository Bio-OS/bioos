package notebook

import (
	"context"

	"github.com/Bio-OS/bioos/internal/context/workspace/application/query/workspace"
	"github.com/Bio-OS/bioos/internal/context/workspace/domain/notebook"
)

type UpdateCommand struct {
	Name        string
	WorkspaceID string
	Content     []byte
}

// UpdateHandler ...
type UpdateHandler interface {
	Handle(context.Context, *UpdateCommand) error
}

// NewUpdateHandler ...
func NewUpdateHandler(svc notebook.Service, workspaceReadModel workspace.WorkspaceReadModel) UpdateHandler {
	return &updateHandler{
		factory:            notebook.NewFactory(),
		service:            svc,
		workspaceReadModel: workspaceReadModel,
	}
}

type updateHandler struct {
	factory            *notebook.Factory
	service            notebook.Service
	workspaceReadModel workspace.WorkspaceReadModel
}

func (h *updateHandler) Handle(ctx context.Context, cmd *UpdateCommand) error {
	if err := workspace.CheckWorkspaceExist(ctx, h.workspaceReadModel, cmd.WorkspaceID); err != nil {
		return err
	}
	// TODO how to get and merge ?
	nb, err := h.factory.New(&notebook.CreateParam{
		Name:        cmd.Name,
		WorkspaceID: cmd.WorkspaceID,
		Content:     cmd.Content,
	})
	if err != nil {
		return err
	}
	if err := h.service.Update(ctx, nb); err != nil {
		return err
	}
	return nil
}
