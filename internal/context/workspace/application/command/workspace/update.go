package workspace

import (
	"context"

	"github.com/Bio-OS/bioos/internal/context/workspace/domain/workspace"
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
	"github.com/Bio-OS/bioos/pkg/validator"
)

type UpdateWorkspaceHandler interface {
	Handle(ctx context.Context, cmd *UpdateWorkspaceCommand) error
}

func NewUpdateWorkspaceHandler(workspaceRepo workspace.Repository, eventBus eventbus.EventBus) UpdateWorkspaceHandler {
	return &updateWorkspaceHandler{
		workspaceRepo: workspaceRepo,
		eventBus:      eventBus,
	}
}

type updateWorkspaceHandler struct {
	workspaceRepo workspace.Repository
	eventBus      eventbus.EventBus
}

func (d updateWorkspaceHandler) Handle(ctx context.Context, cmd *UpdateWorkspaceCommand) error {
	if err := validator.Validate(cmd); err != nil {
		return err
	}

	ws, err := d.workspaceRepo.Get(ctx, cmd.ID)
	if err != nil {
		return err
	}
	if cmd.Name != nil {
		ws.UpdateName(*cmd.Name)
	}
	if cmd.Description != nil {
		ws.UpdateDescription(*cmd.Description)
	}

	return d.workspaceRepo.Save(ctx, ws)
}

var _ UpdateWorkspaceHandler = &updateWorkspaceHandler{}
