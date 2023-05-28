package command

import (
	"context"

	"github.com/Bio-OS/bioos/internal/context/notebookserver/domain"
	"github.com/Bio-OS/bioos/pkg/notebook"
)

type UpdateCommand struct {
	ID           string
	WorkspaceID  string
	Image        *string
	ResourceSize *notebook.ResourceSize
}

type UpdateHandler interface {
	Handle(context.Context, *UpdateCommand) error
}

func NewUpdateHandler(svc domain.Service) UpdateHandler {
	return &updateHandler{
		service: svc,
	}
}

type updateHandler struct {
	service domain.Service
}

func (h *updateHandler) Handle(ctx context.Context, cmd *UpdateCommand) error {
	if cmd.Image == nil && cmd.ResourceSize == nil {
		return nil // nothing to update
	}
	// TODO check workspace exist
	do := domain.NotebookServer{
		ID:          cmd.ID,
		WorkspaceID: cmd.WorkspaceID,
	}
	if cmd.Image != nil {
		do.Settings.DockerImage = *cmd.Image
	}
	if cmd.ResourceSize != nil {
		do.Settings.ResourceSize = *cmd.ResourceSize
	}
	return h.service.Update(ctx, &do)
}
