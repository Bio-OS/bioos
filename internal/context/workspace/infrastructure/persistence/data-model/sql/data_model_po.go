package sql

import "time"

type DataModel struct {
	ID          string
	WorkspaceID string `gorm:"type:varchar(32);not null;uniqueIndex:data_ws"`
	Name        string `gorm:"type:varchar(50) CHARACTER SET gbk COLLATE gbk_bin;not null;uniqueIndex:data_ws"`
	Type        string `gorm:"type:varchar(32);not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (d *DataModel) TableName() string {
	return "data_model"
}

type EntityHeader struct {
	ColumnIndex int    `gorm:"primaryKey;not null"`
	DataModelID string `gorm:"primaryKey;type:varchar(32);not null;index"`
	Name        string `gorm:"type:varchar(100) CHARACTER SET gbk COLLATE gbk_bin;not null"`
	Type        string `gorm:"type:varchar(32);not null"`
}

func (d *EntityHeader) TableName() string {
	return "data_model_entity_header"
}

type EntityGrid struct {
	RowID       string `gorm:"primaryKey;type:varchar(100) CHARACTER SET gbk COLLATE gbk_bin;not null"`
	ColumnIndex int    `gorm:"primaryKey;not null"`
	DataModelID string `gorm:"primaryKey;type:varchar(32);not null;index"`
	Value       string `gorm:"type:longtext CHARACTER SET gbk COLLATE gbk_bin;not null"`
}

func (d *EntityGrid) TableName() string {
	return "data_model_entity_grid"
}

type EntitySetRow struct {
	ID          uint
	RowID       string `gorm:"type:varchar(100) CHARACTER SET gbk COLLATE gbk_bin;not null"`
	DataModelID string `gorm:"type:varchar(32);not null;index"`
	RefRowID    string `gorm:"type:varchar(100) CHARACTER SET gbk COLLATE gbk_bin;not null"`
}

func (d *EntitySetRow) TableName() string {
	return "data_model_entity_set_row"
}

type WorkspaceRow struct {
	Key         string `gorm:"primaryKey;type:varchar(100) CHARACTER SET gbk COLLATE gbk_bin;not null"`
	DataModelID string `gorm:"primaryKey;type:varchar(32);not null;index"`
	Value       string `gorm:"type:longtext CHARACTER SET gbk COLLATE gbk_bin;not null"`
	Type        string `gorm:"type:varchar(32);not null"`
}

func (d *WorkspaceRow) TableName() string {
	return "data_model_workspace_row"
}
