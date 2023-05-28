package run

import "context"

// Repository allows to get/save events from/to event store.
type Repository interface {
	Save(ctx context.Context, r *Run) error
	Get(ctx context.Context, id string) (*Run, error)
	Delete(ctx context.Context, r *Run) error
}
