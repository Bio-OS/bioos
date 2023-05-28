package datamodel

import (
	"context"
)

// Repository allows to get/save events from/to event store.
type Repository interface {
	Save(ctx context.Context, dm *DataModel) error
	Get(ctx context.Context, id string) (*DataModel, error)
	Delete(ctx context.Context, dm *DataModel) error
}
