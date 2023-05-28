package submission

import (
	"context"

	"github.com/Bio-OS/bioos/internal/context/submission/application/query/run"
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
	"github.com/Bio-OS/bioos/pkg/consts"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	"github.com/Bio-OS/bioos/pkg/utils"
)

type DeleteHandler struct {
	repository   Repository
	eventbus     eventbus.EventBus
	runReadModel run.ReadModel
}

func NewDeleteHandler(repository Repository, eventbus eventbus.EventBus, runReadModel run.ReadModel) *DeleteHandler {
	return &DeleteHandler{
		repository:   repository,
		eventbus:     eventbus,
		runReadModel: runReadModel,
	}
}

func (h *DeleteHandler) Handle(ctx context.Context, event *EventSubmission) (err error) {
	if event == nil {
		return nil
	}
	sub, err := h.repository.Get(ctx, event.SubmissionID)
	if err != nil {
		return err
	}
	if utils.In(sub.Status, consts.AllowCancelSubmissionStatuses) {
		cancelEvent := NewCancelSubmissionEvent(event.SubmissionID, 0)
		if err = h.eventbus.Publish(ctx, cancelEvent); err != nil {
			return apperrors.NewInternalError(err)
		}
		sub.Status = consts.SubmissionCancelling
		if err = h.repository.Save(ctx, sub); err != nil {
			return err
		}
	}
	if utils.In(sub.Status, consts.NonFinishedSubmissionStatuses) {
		deleteEvent := NewDeleteSubmissionEvent(event.SubmissionID, 100)
		return h.eventbus.Publish(ctx, deleteEvent)
	}
	if utils.In(sub.Status, consts.FinishedSubmissionStatuses) {
		runIDs, err := h.runReadModel.ListAllRunIDs(ctx, event.SubmissionID)
		if err != nil {
			return err
		}
		for _, runID := range runIDs {
			deleteEvent := NewEventDeleteRun(runID)
			if err = h.eventbus.Publish(ctx, deleteEvent); err != nil {
				return apperrors.NewInternalError(err)
			}
		}
		return h.repository.Delete(ctx, sub)
	}
	return nil
}
