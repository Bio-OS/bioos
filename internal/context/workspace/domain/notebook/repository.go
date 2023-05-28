package notebook

import "context"

// Repository allows to get/save events from/to event store.
type Repository interface {
	Save(context.Context, *Notebook) error
	Get(context.Context, string) (*Notebook, error)
	Delete(context.Context, *Notebook) error
}
