package run

import (
	"context"
	"math/rand"
	"time"

	"github.com/Bio-OS/bioos/internal/context/submission/domain/submission"
	"github.com/Bio-OS/bioos/internal/context/submission/infrastructure/client/wes"
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
	"github.com/Bio-OS/bioos/pkg/consts"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	applog "github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/utils"
)

// EventHandlerSyncRun submit run -> list run/task
type EventHandlerSyncRun struct {
	wes        wes.Client
	runRepo    Repository
	runFactory Factory
	eventBus   eventbus.EventBus
}

func NewEventHandlerSyncRun(wesClient wes.Client, runRepo Repository, runFactory Factory, eventBus eventbus.EventBus) *EventHandlerSyncRun {
	return &EventHandlerSyncRun{
		wes:        wesClient,
		runRepo:    runRepo,
		runFactory: runFactory,
		eventBus:   eventBus,
	}
}

func (e *EventHandlerSyncRun) Handle(ctx context.Context, event *submission.EventRun) error {
	curRun, err := e.runRepo.Get(ctx, event.RunID)
	if err != nil {
		return err
	}

	if curRun.IsFinished() {
		return nil
	}

	if curRun.EngineRunID == "" {
		// mark failed
		return e.markRunFailed(ctx, curRun, "nil engineRunID while sync run")
	}

	resp, err := e.wes.GetRunLog(ctx, &wes.GetRunLogRequest{RunID: curRun.EngineRunID})
	if err != nil {
		if wes.IsNotFound(err) {
			return e.markRunFailed(ctx, curRun, "not found ")
		}

		applog.Errorw("failed to get run log", "err", err)
		return e.republicCurrentEvent(ctx, event.RunID)
	}
	updatedRun := e.syncRunStatus(curRun, resp)
	taskList, err := e.genTasks(event.RunID, updatedRun.Status, resp)
	if err != nil {
		return err
	}
	updatedRun.Tasks = taskList
	if err := e.runRepo.Save(ctx, updatedRun); err != nil {
		return err
	}

	// public sync submission event to update output row datamodel
	eventSyncSubmission := submission.NewSyncSubmissionEvent(updatedRun.SubmissionID)
	if err := e.eventBus.Publish(ctx, eventSyncSubmission); err != nil {
		return apperrors.NewInternalError(err)
	}

	if updatedRun.IsFinished() {
		return nil
	}
	return e.republicCurrentEvent(ctx, event.RunID)
}

func (e *EventHandlerSyncRun) republicCurrentEvent(ctx context.Context, runID string) error {
	newEventSyncRun := submission.NewEventSyncRun(runID, genReSyncDelayTime())
	if err := e.eventBus.Publish(ctx, newEventSyncRun); err != nil {
		return apperrors.NewInternalError(err)
	}
	return nil
}

func (e *EventHandlerSyncRun) markRunFailed(ctx context.Context, run *Run, message string) error {
	tempRun := run.Copy()
	tempRun.Message = utils.PointString(message)
	tempRun.FinishTime = utils.PointTime(time.Now())
	tempRun.Status = consts.RunFailed
	if err := e.runRepo.Save(ctx, tempRun); err != nil {
		return apperrors.NewInternalError(err)
	}
	return nil
}

func (e *EventHandlerSyncRun) syncRunStatus(run *Run, resp *wes.GetRunLogResponse) *Run {
	runStatus := convertRunStatus(run, resp.State)
	tempRun := run.Copy()
	tempRun.Status = runStatus
	tempRun.EngineRunID = resp.RunID
	if tempRun.IsFinished() {
		tempRun.FinishTime = utils.PointTime(time.Now())
		if len(resp.RunLog.Log) != 0 {
			tempRun.Log = utils.PointString(resp.RunLog.Log)
		}
	}
	if tempRun.Status == consts.RunFailed {
		tempRun.Message = utils.PointString(resp.RunLog.Stderr)
	}

	if len(resp.Outputs) > 0 {
		tempRun.Outputs = &resp.Outputs
	}
	return tempRun
}

func (e *EventHandlerSyncRun) genTasks(runID string, runState string, resp *wes.GetRunLogResponse) ([]*Task, error) {
	if len(resp.TaskLogs) == 0 {
		return nil, nil
	}
	taskList := make([]*Task, 0)
	for _, log := range resp.TaskLogs {
		taskParam := CreateTaskParam{
			Name:   log.Name,
			RunID:  runID,
			Status: convertTaskStatus(runState, log.ExitCode),
			Stdout: log.Stdout,
			Stderr: log.Stderr,
		}
		if log.StartTime != nil {
			taskParam.StartTime = log.StartTime.Time()
		} else {
			taskParam.StartTime = resp.RunLog.StartTime.Time()
		}
		taskParam.FinishTime = log.EndTime.PointTime()
		task, err := e.runFactory.CreateWithTaskParam(taskParam)
		if err != nil {
			return nil, apperrors.NewInternalError(err)
		}
		taskList = append(taskList, task)
	}
	return taskList, nil
}

func convertRunStatus(run *Run, s wes.RunState) string {
	if run.IsCancelling() {
		return convertCancellingRunStatus(s)
	}
	return convertNormalRunStatus(s)
}

// If user click cancel, run is Cancelling, but the run finished before the abort request
// actually be sent to cromwell, then the run should be Succeeded or Failed, not Cancelled.
func convertCancellingRunStatus(s wes.RunState) string {
	switch s {
	case wes.RunStateComplete:
		return consts.RunSucceeded
	case wes.RunStateExecutorError, wes.RunStateSystemError:
		return consts.RunFailed
	case wes.RunStateCanceled:
		return consts.RunCancelled
	default:
		return consts.RunCancelling
	}
}

// If user doesn't click cancel, but send abort request to cromwell directly, the run should
// not be Cancelling or Cancelled, but Running or Failed.
func convertNormalRunStatus(s wes.RunState) string {
	switch s {
	case wes.RunStateComplete:
		return consts.RunSucceeded
	case wes.RunStateExecutorError, wes.RunStateSystemError, wes.RunStateCanceled:
		return consts.RunFailed
	default:
		return consts.RunRunning
	}
}

func convertTaskStatus(runState string, exitCode *int32) string {
	if exitCode == nil {
		switch runState {
		case consts.RunCancelled:
			return consts.TaskCancelled
		case consts.RunFailed:
			return consts.TaskFailed
		case consts.RunSucceeded:
			return consts.TaskSucceeded
		default:
			return consts.TaskRunning
		}
	}
	if *exitCode == 0 {
		return consts.TaskSucceeded
	}
	return consts.TaskFailed
}

// get 5s ~ 10s random delay time
func genReSyncDelayTime() time.Duration {
	rand.Seed(time.Now().UnixNano())
	randomFloat := rand.Float64()*5 + 5
	return time.Duration(randomFloat * float64(time.Second))
}
