package run

import (
	"context"

	"github.com/Bio-OS/bioos/pkg/utils"
)

type ReadModel interface {
	ListAllRunIDs(ctx context.Context, submissionID string) ([]string, error)
	ListRuns(ctx context.Context, submissionID string, pg *utils.Pagination, filter *ListRunsFilter) ([]*RunItem, error)
	CountRuns(ctx context.Context, submissionID string, filter *ListRunsFilter) (int, error)

	ListTasks(ctx context.Context, runID string, pg *utils.Pagination) ([]*TaskItem, error)
	CountTasks(ctx context.Context, runID string) (int, error)
	CountRunsResult(ctx context.Context, submissionID string) ([]*StatusCount, error)
	CountTasksResult(ctx context.Context, runID string) ([]*StatusCount, error)
}
