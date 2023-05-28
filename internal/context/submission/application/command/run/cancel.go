package run

import (
	"context"

	"github.com/Bio-OS/bioos/internal/context/submission/application/query/submission"
	"github.com/Bio-OS/bioos/internal/context/submission/domain/run"
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
	"github.com/Bio-OS/bioos/pkg/validator"
)

type CancelRunHandler interface {
	Handle(ctx context.Context, cmd *CancelRunCommand) error
}

type cancelRunHandler struct {
	service             run.Service
	eventBus            eventbus.EventBus
	submissionReadModel submission.ReadModel
}

var _ CancelRunHandler = &cancelRunHandler{}

func NewCancelRunHandler(service run.Service, eventBus eventbus.EventBus, submissionReadModel submission.ReadModel) CancelRunHandler {
	return &cancelRunHandler{
		service:             service,
		eventBus:            eventBus,
		submissionReadModel: submissionReadModel,
	}
}

func (c *cancelRunHandler) Handle(ctx context.Context, cmd *CancelRunCommand) error {
	if err := validator.Validate(cmd); err != nil {
		return err
	}

	if err := c.service.CheckWorkspaceExist(ctx, cmd.WorkspaceID); err != nil {
		return err
	}
	if err := submission.CheckSubmissionExist(ctx, c.submissionReadModel, cmd.WorkspaceID, cmd.SubmissionID); err != nil {
		return err
	}
	return c.service.Cancel(ctx, cmd.ID)
}
