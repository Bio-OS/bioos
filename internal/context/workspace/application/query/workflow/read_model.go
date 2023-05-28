package workflow

import (
	"context"

	"github.com/Bio-OS/bioos/pkg/utils"
)

// ReadModel workflow read model
type ReadModel interface {
	GetById(ctx context.Context, workspaceID, id string) (*Workflow, error)
	GetByName(ctx context.Context, workspaceID, name string) (*Workflow, error)
	List(ctx context.Context, workspaceID string, pg *utils.Pagination, filter *ListWorkflowsFilter) ([]*Workflow, int, error)
	GetVersion(ctx context.Context, id string) (*WorkflowVersion, error)
	ListVersions(ctx context.Context, workflowID string, pg *utils.Pagination, filter *ListWorkflowVersionsFilter) ([]*WorkflowVersion, int, error)
	GetFile(ctx context.Context, id string) (*WorkflowFile, error)
	ListFiles(ctx context.Context, workflowVersionID string, pg *utils.Pagination, filter *ListWorkflowFilesFilter) ([]*WorkflowFile, int, error)
}
