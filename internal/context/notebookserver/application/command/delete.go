package command

import (
	"context"

	"github.com/Bio-OS/bioos/internal/context/notebookserver/domain"
)

type DeleteCommand struct {
	ID          string
	WorkspaceID string
}

type DeleteHandler interface {
	Handle(context.Context, *DeleteCommand) error
}

func NewDeleteHandler(svc domain.Service) DeleteHandler {
	return &deleteHandler{
		service: svc,
	}
}

type deleteHandler struct {
	service domain.Service
}

func (h *deleteHandler) Handle(ctx context.Context, cmd *DeleteCommand) error {
	return h.service.Delete(ctx, cmd.ID)
}
