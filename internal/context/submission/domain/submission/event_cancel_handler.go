package submission

import (
	"context"
	"time"

	"github.com/Bio-OS/bioos/internal/context/submission/application/query/run"
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
	"github.com/Bio-OS/bioos/pkg/consts"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	"github.com/Bio-OS/bioos/pkg/utils"
)

type CancelHandler struct {
	repository   Repository
	runReadModel run.ReadModel
	eventbus     eventbus.EventBus
}

func NewCancelHandler(repository Repository, eventbus eventbus.EventBus, runReadModel run.ReadModel) *CancelHandler {
	return &CancelHandler{
		repository:   repository,
		eventbus:     eventbus,
		runReadModel: runReadModel,
	}
}

func (h *CancelHandler) Handle(ctx context.Context, event *EventSubmission) (err error) {
	if event == nil {
		return nil
	}
	sub, err := h.repository.Get(ctx, event.SubmissionID)
	if err != nil {
		return err
	}
	if sub.Status == consts.SubmissionCancelled {
		return nil
	}
	runCount := 0
	switch sub.Type {
	case consts.FilePathTypeSubmission:
		runCount = len(sub.Inputs)
	case consts.DataModelTypeSubmission:
		runCount = len(sub.DataModelRowIDs)
	}
	runIDs, err := h.runReadModel.ListAllRunIDs(ctx, event.SubmissionID)
	if err != nil {
		return err
	}
	for _, runID := range runIDs {
		newEvent := NewEventCancelRun(runID)
		if err = h.eventbus.Publish(ctx, newEvent); err != nil {
			return apperrors.NewInternalError(err)
		}
	}
	if len(runIDs) != runCount {
		if time.Since(sub.StartTime) < 5*time.Minute {
			newEvent := NewCancelSubmissionEvent(event.SubmissionID, 10)
			if err = h.eventbus.Publish(ctx, newEvent); err != nil {
				return apperrors.NewInternalError(err)
			}
		} else {
			sub.Status = consts.SubmissionCancelled
			if sub.FinishTime == nil {
				sub.FinishTime = utils.PointTime(time.Now())
			}
			if err = h.repository.Save(ctx, sub); err != nil {
				return err
			}
			return
		}
	}
	if sub.Status != consts.SubmissionCancelling {
		sub.Status = consts.SubmissionCancelling
		if err = h.repository.Save(ctx, sub); err != nil {
			return err
		}
	}
	return nil
}
