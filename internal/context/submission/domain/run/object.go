package run

import (
	"time"

	"github.com/jinzhu/copier"

	"github.com/Bio-OS/bioos/pkg/consts"
)

// Run ...
type Run struct {
	ID           string
	Name         string
	SubmissionID string
	Inputs       map[string]interface{}  `gorm:"serializer:json"`
	Outputs      *map[string]interface{} `gorm:"serializer:json"`
	EngineRunID  string
	Status       string
	Log          *string
	Message      *string
	StartTime    time.Time
	FinishTime   *time.Time
	Tasks        []*Task
}

// Task ...
type Task struct {
	Name       string
	RunID      string
	Status     string
	Stdout     string
	Stderr     string
	StartTime  time.Time
	FinishTime *time.Time
}

func (run *Run) Copy() *Run {
	copy := &Run{}
	copier.Copy(copy, run)
	return copy
}

func (run *Run) IsCancelling() bool {
	return run.Status == consts.RunCancelling
}

func (run *Run) IsFinished() bool {
	return run.Status == consts.RunFailed || run.Status == consts.RunSucceeded || run.Status == consts.RunCancelled
}
