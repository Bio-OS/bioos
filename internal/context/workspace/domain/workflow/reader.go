package workflow

import (
	"context"
)

// Reader workflow reader
type Reader interface {
	ParseWorkflowVersion(ctx context.Context, mainWorkflowPath string) (string, error)
	ValidateWorkflowFiles(ctx context.Context, version *WorkflowVersion, baseDir, mainWorkflowPath string) error
	GetWorkflowInputs(ctx context.Context, WorkflowFilePath string) ([]WorkflowParam, error)
	GetWorkflowOutputs(ctx context.Context, WorkflowFilePath string) ([]WorkflowParam, error)
	GetWorkflowGraph(ctx context.Context, WorkflowFilePath string) (string, error)
}

var (
	_ Reader = &WDLReader{}
	_ Reader = &NextflowReader{}
)
