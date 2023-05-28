package workflow

import (
	"context"
	"time"

	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	"github.com/Bio-OS/bioos/pkg/utils"
)

type Factory struct{}

func NewFactory(_ context.Context) *Factory {
	return &Factory{}
}

type WorkflowOption struct {
	ID          string
	Name        string
	Description *string
}

func (p *WorkflowOption) Validate() error {
	return nil
}

type Param struct {
	ID            string
	Name          string
	WorkspaceID   string
	Description   string
	LatestVersion string
	CreateTime    time.Time
	UpdateTime    time.Time
	DeletedAt     time.Time
}

// nolint
func (p Param) validate() error {
	return nil
}

func (f *Factory) NewWorkflow(workspaceID string, param *WorkflowOption) (*Workflow, error) {
	if err := param.Validate(); err != nil {
		return nil, err
	}
	if len(param.ID) == 0 {
		param.ID = utils.GenWorkflowID()
	}

	// validate param
	workflow := &Workflow{
		ID:          param.ID,
		Name:        param.Name,
		WorkspaceID: workspaceID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if param.Description != nil {
		workflow.Description = *param.Description
	}
	return workflow, nil
}

type VersionOption struct {
	Language         string
	MainWorkflowPath string
	Source           string
	Url              string
	Tag              string
	Token            string
}

func (p VersionOption) validate() error {
	if p.Source != WorkflowSourceGit && p.Source != WorkflowSourceFile {
		return apperrors.NewInvalidError("source")
	}
	if p.Source == WorkflowSourceGit {
		if p.Url == "" {
			return apperrors.NewInvalidError("url")
		}
		if p.Tag == "" {
			return apperrors.NewInvalidError("tag")
		}
	}
	return nil
}

type VersionUpdateOption struct {
	Language         *string
	MainWorkflowPath *string
	Source           *string
	URL              *string
	Tag              *string
	Token            *string
}

type FileParam struct {
	Path    string
	Content string
}

func (p FileParam) validate() error {
	if p.Content == "" {
		return apperrors.NewInvalidError("content")
	}
	if p.Path == "" {
		return apperrors.NewInvalidError("path")
	}
	return nil
}
