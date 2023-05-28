package domain

import "context"

// Repository ...
type Repository interface {
	Save(context.Context, *NotebookServer) error
	Get(context.Context, string) (*NotebookServer, error)
	Delete(context.Context, *NotebookServer) error
}
