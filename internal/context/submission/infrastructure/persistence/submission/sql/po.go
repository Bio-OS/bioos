package sql

import (
	"time"

	"gorm.io/gorm"
)

type SubmissionModel struct {
	Submission
	// '软删除标识, 空表示未删除，非空值为删除时间',
	DeletedAt gorm.DeletedAt
}

type Submission struct {
	ID                string
	WorkspaceID       string                 `gorm:"type:varchar(32);not null;uniqueIndex:sub_ws"`
	Name              string                 `gorm:"type:varchar(410) CHARACTER SET gbk COLLATE gbk_bin;not null;uniqueIndex:sub_ws"`
	Description       *string                `gorm:"type:text"`
	WorkflowID        string                 `gorm:"type:varchar(32);not null"`
	WorkflowVersionID string                 `gorm:"type:varchar(32);not null"`
	DataModelID       *string                `gorm:"type:varchar(32)"`
	DataModelRowIDs   *string                `gorm:"type:text"`
	Type              string                 `gorm:"type:varchar(32);not null"`
	Inputs            map[string]interface{} `gorm:"serializer:json"`
	Outputs           map[string]interface{} `gorm:"serializer:json"`
	ExposedOptions    ExposedOptions         `gorm:"serializer:json"`
	Status            string                 `gorm:"type:varchar(32);not null"`
	StartTime         time.Time              `gorm:"not null"`
	FinishTime        *time.Time
	UserID            *int64
}

type ExposedOptions struct {
	ReadFromCache bool `json:"readFromCache"`
}

func (s *SubmissionModel) TableName() string {
	return "submission"
}

func (s *Submission) TableName() string {
	return "submission"
}
