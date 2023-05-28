package workspace

import (
	"context"
)

// Repository allows to get/save events from/to event store.
type Repository interface {
	Save(ctx context.Context, w *Workspace) error
	Get(ctx context.Context, id string) (*Workspace, error)
	Delete(ctx context.Context, w *Workspace) error
}
