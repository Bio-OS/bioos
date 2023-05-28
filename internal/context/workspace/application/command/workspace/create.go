package workspace

import (
	"context"

	"github.com/Bio-OS/bioos/internal/context/workspace/domain/workspace"
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
	applog "github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/validator"
)

type CreateWorkspaceHandler interface {
	Handle(ctx context.Context, cmd *CreateWorkspaceCommand) (string, error)
}

type createWorkspaceHandler struct {
	workspaceRepo    workspace.Repository
	workspaceFactory *workspace.Factory
	eventBus         eventbus.EventBus
}

var _ CreateWorkspaceHandler = &createWorkspaceHandler{}

func NewCreateWorkspaceHandler(workspaceRepo workspace.Repository, workspaceFactory *workspace.Factory, eventBus eventbus.EventBus) CreateWorkspaceHandler {
	return &createWorkspaceHandler{
		workspaceRepo:    workspaceRepo,
		workspaceFactory: workspaceFactory,
		eventBus:         eventBus,
	}
}

func (h *createWorkspaceHandler) Handle(ctx context.Context, cmd *CreateWorkspaceCommand) (string, error) {
	if err := validator.Validate(cmd); err != nil {
		return "", err
	}

	param := workspace.CreateWorkspaceParam{
		Name:        cmd.Name,
		Description: cmd.Description,
	}
	if cmd.Storage.NFS != nil {
		param.Storage.NFS = &workspace.NFSStorage{MountPath: cmd.Storage.NFS.MountPath}
	}
	ws, err := h.workspaceFactory.CreateWithWorkspaceParam(param)
	if err != nil {
		return "", err
	}
	if err := h.workspaceRepo.Save(ctx, ws); err != nil {
		return "", err
	}

	applog.Infow("publish workspace created event", "ID", ws.ID)
	return ws.ID, nil
}
