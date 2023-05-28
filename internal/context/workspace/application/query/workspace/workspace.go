package workspace

import (
	"context"

	"github.com/Bio-OS/bioos/pkg/utils"
)

type WorkspaceReadModel interface {
	ListWorkspaces(ctx context.Context, pg utils.Pagination, filter *ListWorkspacesFilter) ([]*WorkspaceItem, error)
	CountWorkspaces(ctx context.Context, filter *ListWorkspacesFilter) (int, error)
	GetWorkspaceById(ctx context.Context, id string) (*WorkspaceItem, error)
}
