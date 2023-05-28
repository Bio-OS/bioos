package run

import (
	"context"
	"time"

	"github.com/Bio-OS/bioos/internal/context/submission/domain/submission"
	"github.com/Bio-OS/bioos/internal/context/submission/infrastructure/client/wes"
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
	"github.com/Bio-OS/bioos/pkg/consts"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	applog "github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/utils"
)

type CancelHandler struct {
	wes        wes.Client
	repository Repository
	eventbus   eventbus.EventBus
}

func NewEventHandlerCancelRun(wes wes.Client, repository Repository, eventbus eventbus.EventBus) *CancelHandler {
	return &CancelHandler{
		wes:        wes,
		repository: repository,
		eventbus:   eventbus,
	}
}

func (h *CancelHandler) Handle(ctx context.Context, event *submission.EventRun) (err error) {
	if event == nil {
		return nil
	}
	run, err := h.repository.Get(ctx, event.RunID)
	if err != nil {
		return err
	}
	if len(run.EngineRunID) == 0 {
		run.Status = consts.RunCancelled
		run.FinishTime = utils.PointTime(time.Now())
		eventSyncSubmission := submission.NewSyncSubmissionEvent(run.SubmissionID)
		if err := h.eventbus.Publish(ctx, eventSyncSubmission); err != nil {
			return apperrors.NewInternalError(err)
		}
		return h.repository.Save(ctx, run)
	}
	if _, err = h.wes.CancelRun(ctx, &wes.CancelRunRequest{RunID: run.EngineRunID}); err != nil {
		if !wes.IsNotFound(err) {
			return apperrors.NewInternalError(err)
		}
		applog.Warnf("engine run %s not found", run.EngineRunID)
	}
	syncEvent := submission.NewEventSyncRun(event.RunID, 0)
	if err = h.eventbus.Publish(ctx, syncEvent); err != nil {
		return apperrors.NewInternalError(err)
	}
	if run.Status != consts.RunCancelling {
		run.Status = consts.RunCancelling
		return h.repository.Save(ctx, run)
	}
	return nil
}
