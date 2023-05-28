package submission

import (
	"context"

	"github.com/Bio-OS/bioos/internal/context/submission/domain/submission"
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
	"github.com/Bio-OS/bioos/pkg/validator"
)

type CancelSubmissionHandler interface {
	Handle(ctx context.Context, cmd *CancelSubmissionCommand) error
}

type cancelSubmissionHandler struct {
	service  submission.Service
	eventBus eventbus.EventBus
}

var _ CancelSubmissionHandler = &cancelSubmissionHandler{}

func NewCancelSubmissionHandler(service submission.Service, eventBus eventbus.EventBus) CancelSubmissionHandler {
	return &cancelSubmissionHandler{
		service:  service,
		eventBus: eventBus,
	}
}

func (c *cancelSubmissionHandler) Handle(ctx context.Context, cmd *CancelSubmissionCommand) error {
	if err := validator.Validate(cmd); err != nil {
		return err
	}

	if err := c.service.CheckWorkspaceExist(ctx, cmd.WorkspaceID); err != nil {
		return err
	}
	return c.service.Cancel(ctx, cmd.ID)
}
