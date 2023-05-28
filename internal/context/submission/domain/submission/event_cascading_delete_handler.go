package submission

import (
	"context"

	submissionquery "github.com/Bio-OS/bioos/internal/context/submission/application/query/submission"
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	"github.com/Bio-OS/bioos/pkg/utils"
)

type CascadeDeleteHandler struct {
	repository          Repository
	submissionReadModel submissionquery.ReadModel
	eventbus            eventbus.EventBus
}

func NewCascadeDeleteHandler(repository Repository, eventbus eventbus.EventBus, submissionReadModel submissionquery.ReadModel) *CascadeDeleteHandler {
	return &CascadeDeleteHandler{
		repository:          repository,
		eventbus:            eventbus,
		submissionReadModel: submissionReadModel,
	}
}

func (h *CascadeDeleteHandler) Handle(ctx context.Context, event *CascadeDeleteSubmissionEvent) (err error) {
	filter := &submissionquery.ListSubmissionsFilter{}
	if event.Workflow != nil {
		filter.WorkflowID = *event.Workflow
	}
	count, err := h.submissionReadModel.CountSubmissions(ctx, event.WorkspaceID, filter)
	if err != nil {
		return err
	}

	submissions, err := h.submissionReadModel.ListSubmissions(ctx, event.WorkspaceID, utils.NewPagination(count, 1), filter)
	if err != nil {
		return err
	}

	for _, submission := range submissions {
		deleteEvent := NewDeleteSubmissionEvent(submission.ID, 0)
		if err = h.eventbus.Publish(ctx, deleteEvent); err != nil {
			return apperrors.NewInternalError(err)
		}
	}

	return nil
}
