package run

import (
	"context"
	"fmt"
	"time"

	"github.com/Bio-OS/bioos/internal/context/submission/domain/submission"
	"github.com/Bio-OS/bioos/internal/context/submission/infrastructure/client/wes"
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	applog "github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/utils"
)

const (
	BioosRunIDKey = "bioos-run-id"
)

type EventHandlerSubmitRun struct {
	wes      wes.Client
	runRepo  Repository
	eventBus eventbus.EventBus
}

func NewEventHandlerSubmitRun(wesClient wes.Client, eventBus eventbus.EventBus, runRepo Repository) *EventHandlerSubmitRun {
	return &EventHandlerSubmitRun{
		wes:      wesClient,
		runRepo:  runRepo,
		eventBus: eventBus,
	}
}

func (e *EventHandlerSubmitRun) Handle(ctx context.Context, event *submission.EventSubmitRun) error {
	run, err := e.runRepo.Get(ctx, event.RunID)
	if err != nil {
		return apperrors.NewInternalError(err)
	}

	if run == nil {
		applog.Warnf("can not find run with ID:%s", event.RunID)
		return nil
	}

	if run.IsCancelling() || run.IsFinished() {
		return nil
	}

	if run.EngineRunID != "" {
		// todo check delay
		eventSync := submission.NewEventSyncRun(event.RunID, 0)
		if err := e.eventBus.Publish(ctx, eventSync); err != nil {
			return apperrors.NewInternalError(err)
		}
		return nil
	}

	// not submit before
	resp, err := e.wes.RunWorkflow(ctx, &wes.RunWorkflowRequest{
		RunRequest: wes.RunRequest{
			WorkflowParams:      run.Inputs,
			WorkflowType:        event.RunConfig.Language,
			WorkflowTypeVersion: event.RunConfig.Version,
			Tags: map[string]interface{}{
				BioosRunIDKey: run.ID,
			},
			WorkflowEngineParameters: event.RunConfig.WorkflowEngineParameters,
		},
		WorkflowAttachment: event.RunConfig.WorkflowContents,
	})
	if err != nil {
		if wes.IsBadRequest(err) {
			// mark failed
			applog.Errorw("bad request to submit run", "err", err)
			return e.markRunFailed(ctx, run, fmt.Sprintf("bad request to submit run: %s", err.Error()))
		}
		// republish
		return e.republicCurrentEvent(ctx, event)
	}

	return e.markRunRunningAndPublicEventSync(ctx, run, resp.RunID, event)
}

func (e *EventHandlerSubmitRun) markRunFailed(ctx context.Context, run *Run, message string) error {
	tempRun := run.Copy()
	tempRun.Message = utils.PointString(message)
	tempRun.FinishTime = utils.PointTime(time.Now())
	if err := e.runRepo.Save(ctx, run); err != nil {
		return apperrors.NewInternalError(err)
	}
	return nil
}

func (e *EventHandlerSubmitRun) markRunRunningAndPublicEventSync(ctx context.Context, run *Run, engineRunID string, event *submission.EventSubmitRun) error {
	tempRun := run.Copy()
	tempRun.EngineRunID = engineRunID
	if err := e.runRepo.Save(ctx, tempRun); err != nil {
		return apperrors.NewInternalError(err)
	}
	eventSync := submission.NewEventSyncRun(event.RunID, 0)
	if err := e.eventBus.Publish(ctx, eventSync); err != nil {
		return apperrors.NewInternalError(err)
	}
	return nil
}

func (e *EventHandlerSubmitRun) republicCurrentEvent(ctx context.Context, event *submission.EventSubmitRun) error {
	newEventSubmitRun := submission.NewEventSubmitRun(event.RunID, event.RunConfig)
	if err := e.eventBus.Publish(ctx, newEventSubmitRun); err != nil {
		return apperrors.NewInternalError(err)
	}
	return nil
}
