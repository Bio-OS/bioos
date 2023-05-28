package run

import (
	"context"

	"github.com/Bio-OS/bioos/internal/context/submission/domain/submission"
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
	"github.com/Bio-OS/bioos/pkg/consts"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	"github.com/Bio-OS/bioos/pkg/utils"
)

type DeleteHandler struct {
	repository Repository
	eventbus   eventbus.EventBus
}

func NewEventHandlerDeleteRun(repository Repository, eventbus eventbus.EventBus) *DeleteHandler {
	return &DeleteHandler{
		repository: repository,
		eventbus:   eventbus,
	}
}

func (h *DeleteHandler) Handle(ctx context.Context, event *submission.EventRun) (err error) {
	if event == nil {
		return nil
	}
	run, err := h.repository.Get(ctx, event.RunID)
	if err != nil {
		return err
	}
	if utils.In(run.Status, consts.NonFinishedRunStatuses) {
		cancelEvent := submission.NewEventCancelRun(event.RunID)
		if err = h.eventbus.Publish(ctx, cancelEvent); err != nil {
			return apperrors.NewInternalError(err)
		}
		deleteEvent := submission.NewEventDeleteRun(event.RunID)
		if err = h.eventbus.Publish(ctx, deleteEvent); err != nil {
			return apperrors.NewInternalError(err)
		}
	}
	return h.repository.Delete(ctx, run)
}
