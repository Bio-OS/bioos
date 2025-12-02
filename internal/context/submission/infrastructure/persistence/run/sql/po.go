package sql

import (
	"time"
)

// Run ...
type Run struct {
	ID                string
	Name              string                  `gorm:"type:varchar(200);not null;uniqueIndex:sub_run"`
	SubmissionID      string                  `gorm:"type:varchar(32);not null;uniqueIndex:sub_run"`
	WorkflowVersionID string                  `gorm:"type:varchar(32)"`
	Cache             bool                    `gorm:"type:bool"`
	Inputs            map[string]interface{}  `gorm:"serializer:json"`
	Outputs           *map[string]interface{} `gorm:"serializer:json"`
	EngineRunID       string                  `gorm:"type:varchar(128);not null"`
	Status            string                  `gorm:"type:varchar(32);not null"`
	Log               *string                 `gorm:"type:longtext"`
	Message           *string                 `gorm:"type:longtext"`
	StartTime         time.Time
	FinishTime        *time.Time
}

func (r *Run) TableName() string {
	return "run"
}

// Task ...
type Task struct {
	Name       string `gorm:"type:varchar(267);primary_key"`
	RunID      string `gorm:"type:varchar(32);primary_key"`
	Status     string `gorm:"type:varchar(32);not null"`
	Stdout     string `gorm:"type:longtext;not null"`
	Stderr     string `gorm:"type:longtext;not null"`
	StartTime  time.Time
	FinishTime *time.Time
}

func (t *Task) TableName() string {
	return "task"
}

// StatusCount ...
type StatusCount struct {
	Count  int64
	Status string
}
