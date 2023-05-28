package notebook

import (
	"context"

	"github.com/Bio-OS/bioos/internal/context/workspace/application/query/workspace"
	"github.com/Bio-OS/bioos/internal/context/workspace/domain/notebook"
)

// CreateCommand ...
type CreateCommand struct {
	Name        string
	WorkspaceID string
	Content     []byte
}

// CreateHandler ...
type CreateHandler interface {
	Handle(context.Context, *CreateCommand) error
}

// NewCreateHandler ...
func NewCreateHandler(svc notebook.Service, workspaceReadModel workspace.WorkspaceReadModel) CreateHandler {
	return &createHandler{
		factory:            notebook.NewFactory(),
		service:            svc,
		workspaceReadModel: workspaceReadModel,
	}
}

type createHandler struct {
	factory            *notebook.Factory
	service            notebook.Service
	workspaceReadModel workspace.WorkspaceReadModel
}

func (h *createHandler) Handle(ctx context.Context, cmd *CreateCommand) error {
	if err := workspace.CheckWorkspaceExist(ctx, h.workspaceReadModel, cmd.WorkspaceID); err != nil {
		return err
	}
	nb, err := h.factory.New(&notebook.CreateParam{
		Name:        cmd.Name,
		WorkspaceID: cmd.WorkspaceID,
		Content:     cmd.Content,
	})
	if err != nil {
		return err
	}
	if err := h.service.Upsert(ctx, nb); err != nil {
		return err
	}
	return nil
}
