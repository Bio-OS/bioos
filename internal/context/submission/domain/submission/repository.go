package submission

import "context"

// Repository allows to get/save events from/to event store.
type Repository interface {
	Save(ctx context.Context, s *Submission) error
	Get(ctx context.Context, id string) (*Submission, error)
	Delete(ctx context.Context, s *Submission) error
	SoftDelete(ctx context.Context, s *Submission) error
}
