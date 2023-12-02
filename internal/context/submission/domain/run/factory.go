package run

import (
	"context"
	"time"

	"github.com/Bio-OS/bioos/pkg/utils"
)

type CreateRunParam struct {
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
	WorkflowType string
}

type CreateTaskParam struct {
	Name       string
	RunID      string
	Status     string
	Stdout     string
	Stderr     string
	StartTime  time.Time
	FinishTime *time.Time
}

func (p CreateRunParam) validate() error {
	return nil
}

func (p CreateTaskParam) validate() error {
	return nil
}

// Factory workspace factory.
type Factory struct{}

// NewRunFactory return a workspace factory.
func NewRunFactory(_ context.Context) *Factory {
	return &Factory{}
}

// CreateWithRunParam ...
func (fac *Factory) CreateWithRunParam(param CreateRunParam) (*Run, error) {
	if err := param.validate(); err != nil {
		return nil, err
	}
	if len(param.ID) == 0 {
		param.ID = utils.GenRunID()
	}
	if param.StartTime.IsZero() {
		param.StartTime = time.Now()
	}

	return &Run{
		ID:           param.ID,
		Name:         param.Name,
		SubmissionID: param.SubmissionID,
		Inputs:       param.Inputs,
		Outputs:      param.Outputs,
		EngineRunID:  param.EngineRunID,
		Status:       param.Status,
		Log:          param.Log,
		Message:      param.Message,
		StartTime:    param.StartTime,
		FinishTime:   param.FinishTime,
		WorkflowType: param.WorkflowType,
	}, nil
}

// CreateWithTaskParam ...
func (fac *Factory) CreateWithTaskParam(param CreateTaskParam) (*Task, error) {
	if err := param.validate(); err != nil {
		return nil, err
	}
	if param.StartTime.IsZero() {
		param.StartTime = time.Now()
	}

	return &Task{
		Name:       param.Name,
		RunID:      param.RunID,
		Status:     param.Status,
		Stdout:     param.Stdout,
		Stderr:     param.Stderr,
		StartTime:  param.StartTime,
		FinishTime: param.FinishTime,
	}, nil
}
