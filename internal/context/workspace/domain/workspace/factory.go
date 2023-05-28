package workspace

import (
	"context"
	"time"

	"github.com/Bio-OS/bioos/pkg/utils"
)

// CreateWorkspaceParam use to create Workspace
type CreateWorkspaceParam struct {
	ID          string
	Name        string
	Description string
	Storage     Storage
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (p CreateWorkspaceParam) validate() error {
	return nil
}

// Factory workspace factory.
type Factory struct{}

// NewWorkspaceFactory return a workspace factory.
func NewWorkspaceFactory(_ context.Context) *Factory {
	return &Factory{}
}

// CreateWithWorkspaceParam ...
func (fac *Factory) CreateWithWorkspaceParam(param CreateWorkspaceParam) (*Workspace, error) {
	if err := param.validate(); err != nil {
		return nil, err
	}
	if len(param.ID) == 0 {
		param.ID = utils.GenWorkspaceID()
	}
	if param.CreatedAt.IsZero() {
		param.CreatedAt = time.Now()
		param.UpdatedAt = time.Now()
	}
	if param.UpdatedAt.IsZero() {
		param.UpdatedAt = time.Now()
	}

	return &Workspace{
		ID:          param.ID,
		Name:        param.Name,
		Description: param.Description,
		CreatedAt:   param.CreatedAt,
		UpdatedAt:   param.UpdatedAt,
		Storage:     param.Storage,
	}, nil
}
