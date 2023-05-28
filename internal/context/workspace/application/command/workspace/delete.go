package workspace

import (
	"context"

	"github.com/Bio-OS/bioos/internal/context/workspace/domain/workspace"
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
	"github.com/Bio-OS/bioos/pkg/validator"
)

type DeleteWorkspaceHandler interface {
	Handle(ctx context.Context, cmd *DeleteWorkspaceCommand) error
}

func NewDeleteWorkspaceHandler(workspaceRepo workspace.Repository, eventBus eventbus.EventBus) DeleteWorkspaceHandler {
	return &deleteWorkspaceHandler{
		workspaceRepo: workspaceRepo,
		eventBus:      eventBus,
	}
}

type deleteWorkspaceHandler struct {
	workspaceRepo workspace.Repository
	eventBus      eventbus.EventBus
}

func (d *deleteWorkspaceHandler) Handle(ctx context.Context, cmd *DeleteWorkspaceCommand) error {
	if err := validator.Validate(cmd); err != nil {
		return err
	}

	ws, err := d.workspaceRepo.Get(ctx, cmd.ID)
	if err != nil {
		return err
	}
	err = d.workspaceRepo.Delete(ctx, ws)

	event := workspace.NewWorkspaceDeletedEvent(ws.ID)

	if eventErr := d.eventBus.Publish(ctx, event); eventErr != nil {
		return eventErr
	}
	return err
}

var _ DeleteWorkspaceHandler = &deleteWorkspaceHandler{}
