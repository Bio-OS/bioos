package datamodel

import (
	"time"

	"github.com/Bio-OS/bioos/pkg/utils"
)

// Factory workspace factory.
type Factory struct{}

// NewDataModelFactory return a workspace factory.
func NewDataModelFactory() *Factory {
	return &Factory{}
}

type CreateParam struct {
	WorkspaceID string
	Name        string
	Type        string
	Headers     []string
	Rows        [][]string
}

func (f *Factory) New(param *CreateParam) *DataModel {
	return &DataModel{
		WorkspaceID: param.WorkspaceID,
		Name:        param.Name,
		ID:          utils.GenDataModelID(),
		Type:        param.Type,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Headers:     param.Headers,
		Rows:        param.Rows,
	}
}
