package workspace

import (
	"context"

	"github.com/Bio-OS/bioos/internal/context/workspace/domain/workspace"
	"github.com/Bio-OS/bioos/pkg/validator"
)

type ImportWorkspaceHandler interface {
	Handle(ctx context.Context, cmd *ImportWorkspaceCommand) error
}

type importWorkspaceHandlerImpl struct {
	service          workspace.Service
	workspaceFactory *workspace.Factory
}

var _ ImportWorkspaceHandler = &importWorkspaceHandlerImpl{}

func NewImportWorkspaceHandler(service workspace.Service, workspaceFactory *workspace.Factory) ImportWorkspaceHandler {
	return &importWorkspaceHandlerImpl{
		service:          service,
		workspaceFactory: workspaceFactory,
	}
}

func (h *importWorkspaceHandlerImpl) Handle(ctx context.Context, cmd *ImportWorkspaceCommand) error {
	if err := validator.Validate(cmd); err != nil {
		return err
	}
	if err := h.service.Import(ctx, cmd.ID, cmd.FileName, workspace.Storage{
		NFS: &workspace.NFSStorage{
			MountPath: cmd.Storage.NFS.MountPath,
		},
	}); err != nil {
		return err
	}
	return nil
}
