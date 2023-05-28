package domain

import (
	"context"
	"fmt"
)

type Runtime interface {
	Create(context.Context, *NotebookServer) error
	Start(context.Context, *NotebookServer) error
	Stop(context.Context, *NotebookServer) error
	Delete(context.Context, *NotebookServer) error
	GetStatus(context.Context, *NotebookServer) (*Status, error)
}

type UnimplementedRuntime struct{}

func (UnimplementedRuntime) Create(context.Context, *NotebookServer) error {
	return fmt.Errorf("notebookserver runtime unimplement")
}

func (UnimplementedRuntime) Start(context.Context, *NotebookServer) error {
	return fmt.Errorf("notebookserver runtime unimplement")
}

func (UnimplementedRuntime) Stop(context.Context, *NotebookServer) error {
	return fmt.Errorf("notebookserver runtime unimplement")
}

func (UnimplementedRuntime) Delete(context.Context, *NotebookServer) error {
	return fmt.Errorf("notebookserver runtime unimplement")
}

func (UnimplementedRuntime) GetStatus(context.Context, *NotebookServer) (*Status, error) {
	return nil, fmt.Errorf("notebookserver runtime unimplement")
}
