package workflow

import "context"

// Repository repository for workflow
type Repository interface {
	Save(context.Context, *Workflow) error
	Get(context.Context, string, string) (*Workflow, error)
	Delete(context.Context, *Workflow) error
	List(ctx context.Context, workspaceID string) ([]string, error)
}
