package submission

import (
	"context"
	"time"

	"github.com/Bio-OS/bioos/pkg/consts"
	"github.com/Bio-OS/bioos/pkg/utils"
)

// CreateSubmissionParam use to create Submission
type CreateSubmissionParam struct {
	Name              string
	Description       *string
	WorkflowID        string
	WorkflowVersionID string
	WorkspaceID       string
	DataModelID       *string
	DataModelRowIDs   []string
	Type              string
	Language          string
	Inputs            map[string]interface{}
	Outputs           map[string]interface{}
	ExposedOptions    ExposedOptions
}

func (p CreateSubmissionParam) validate() error {
	return nil
}

// Factory workspace factory.
type Factory struct{}

// NewSubmissionFactory return a workspace factory.
func NewSubmissionFactory(_ context.Context) *Factory {
	return &Factory{}
}

// CreateWithSubmissionParam ...
func (fac *Factory) CreateWithSubmissionParam(param CreateSubmissionParam) (*Submission, error) {
	if err := param.validate(); err != nil {
		return nil, err
	}

	return &Submission{
		ID:                utils.GenSubmissionID(),
		Name:              param.Name,
		Description:       param.Description,
		WorkflowID:        param.WorkflowID,
		WorkflowVersionID: param.WorkflowVersionID,
		DataModelID:       param.DataModelID,
		DataModelRowIDs:   param.DataModelRowIDs,
		WorkspaceID:       param.WorkspaceID,
		Type:              param.Type,
		Language:          param.Language,
		Inputs:            param.Inputs,
		Outputs:           param.Outputs,
		ExposedOptions:    param.ExposedOptions,
		Status:            consts.SubmissionPending,
		StartTime:         time.Now(),
	}, nil
}
