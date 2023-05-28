package notebook

import "context"

type ReadModel interface {
	ListByWorkspace(ctx context.Context, workspaceID string) ([]*Notebook, error)
	Get(ctx context.Context, workspaceID, name string) (*Notebook, error)
}
