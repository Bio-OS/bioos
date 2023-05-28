package submission

import (
	"context"

	"github.com/Bio-OS/bioos/pkg/utils"
)

type ReadModel interface {
	ListSubmissions(ctx context.Context, workspaceID string, pg *utils.Pagination, filter *ListSubmissionsFilter) ([]*SubmissionItem, error)

	CountSubmissions(ctx context.Context, workspaceID string, filter *ListSubmissionsFilter) (int, error)
}
