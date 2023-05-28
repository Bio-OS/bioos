package sql

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	query "github.com/Bio-OS/bioos/internal/context/submission/application/query/run"
	"github.com/Bio-OS/bioos/internal/context/submission/domain/run"
	"github.com/Bio-OS/bioos/pkg/utils"
)

func RunPOToRunDTO(ctx context.Context, run *Run) (*query.RunItem, error) {
	item := &query.RunItem{
		ID:          run.ID,
		Name:        run.Name,
		Status:      run.Status,
		StartTime:   run.StartTime.Unix(),
		EngineRunID: run.EngineRunID,
		Log:         run.Log,
		Message:     run.Message,
	}
	var inputs, outputs []byte
	var err error
	inputs, err = json.Marshal(run.Inputs)
	if err != nil {
		return nil, fmt.Errorf("failed to transform inputs: %w", err)
	}
	if run.Outputs != nil {
		outputs, err = json.Marshal(*run.Outputs)
		if err != nil {
			return nil, fmt.Errorf("failed to transform outputs: %w", err)
		}
	}
	if run.FinishTime != nil {
		item.FinishTime = utils.PointInt64(run.FinishTime.Unix())
		item.Duration = run.FinishTime.Unix() - run.StartTime.Unix()
	} else {
		item.Duration = time.Now().Unix() - run.StartTime.Unix()
	}
	item.Inputs = string(inputs)
	item.Outputs = string(outputs)
	return item, nil
}

func TaskPOToTaskDTO(ctx context.Context, task *Task) *query.TaskItem {
	item := &query.TaskItem{
		Name:      task.Name,
		RunID:     task.RunID,
		Status:    task.Status,
		StartTime: task.StartTime.Unix(),
		Stdout:    task.Stdout,
		Stderr:    task.Stderr,
	}
	if task.FinishTime != nil {
		item.FinishTime = utils.PointInt64(task.FinishTime.Unix())
		item.Duration = task.FinishTime.Unix() - task.StartTime.Unix()
	} else {
		item.Duration = time.Now().Unix() - task.StartTime.Unix()
	}
	return item
}

func StatusCountPOToStatusCountDTO(count *StatusCount) *query.StatusCount {
	return &query.StatusCount{
		Count:  count.Count,
		Status: count.Status,
	}
}

func RunPOToRunDO(runPO *Run) *run.Run {
	return &run.Run{
		ID:           runPO.ID,
		Name:         runPO.Name,
		SubmissionID: runPO.SubmissionID,
		Inputs:       runPO.Inputs,
		Outputs:      runPO.Outputs,
		EngineRunID:  runPO.EngineRunID,
		Status:       runPO.Status,
		Log:          runPO.Log,
		Message:      runPO.Message,
		StartTime:    runPO.StartTime,
		FinishTime:   runPO.FinishTime,
	}
}

func RunDOToTaskPOList(runDO *run.Run) []*Task {
	taskPOList := make([]*Task, 0)
	for _, curTask := range runDO.Tasks {
		taskPO := &Task{
			Name:       curTask.Name,
			RunID:      curTask.RunID,
			Status:     curTask.Status,
			Stdout:     curTask.Stdout,
			Stderr:     curTask.Stderr,
			StartTime:  curTask.StartTime,
			FinishTime: curTask.FinishTime,
		}
		taskPOList = append(taskPOList, taskPO)
	}
	return taskPOList
}

func RunDOToRunPO(runDO *run.Run) *Run {
	return &Run{
		ID:           runDO.ID,
		Name:         runDO.Name,
		SubmissionID: runDO.SubmissionID,
		Inputs:       runDO.Inputs,
		Outputs:      runDO.Outputs,
		EngineRunID:  runDO.EngineRunID,
		Status:       runDO.Status,
		Log:          runDO.Log,
		Message:      runDO.Message,
		StartTime:    runDO.StartTime,
		FinishTime:   runDO.FinishTime,
	}
}
