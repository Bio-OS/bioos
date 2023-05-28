package submission

import (
	"context"

	"github.com/Bio-OS/bioos/internal/context/submission/domain/submission"
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
	"github.com/Bio-OS/bioos/pkg/validator"
)

type DeleteSubmissionHandler interface {
	Handle(ctx context.Context, cmd *DeleteSubmissionCommand) error
}

type deleteSubmissionHandler struct {
	service  submission.Service
	eventBus eventbus.EventBus
}

var _ DeleteSubmissionHandler = &deleteSubmissionHandler{}

func NewDeleteSubmissionHandler(service submission.Service, eventBus eventbus.EventBus) DeleteSubmissionHandler {
	return &deleteSubmissionHandler{
		service:  service,
		eventBus: eventBus,
	}
}

func (c *deleteSubmissionHandler) Handle(ctx context.Context, cmd *DeleteSubmissionCommand) error {
	if err := validator.Validate(cmd); err != nil {
		return err
	}

	if err := c.service.CheckWorkspaceExist(ctx, cmd.WorkspaceID); err != nil {
		return err
	}
	_, err := c.service.Get(ctx, cmd.ID)
	if err != nil {
		return err
	}

	return c.service.SoftDelete(ctx, cmd.ID)
}
