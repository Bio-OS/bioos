package workflow

import (
	"context"

	"github.com/Bio-OS/bioos/internal/context/workspace/application/query/workspace"
	"github.com/Bio-OS/bioos/pkg/utils"
	"github.com/Bio-OS/bioos/pkg/validator"
)

const (
	FileOrderByVersion = "version"
	FileOrderByPath    = "path"
)

type ListFilesQuery struct {
	Pg                *utils.Pagination        `json:"pg"`
	Filter            *ListWorkflowFilesFilter `json:"filter"`
	WorkspaceID       string                   `json:"workspaceID"`
	WorkflowID        string                   `json:"workflowID"`
	WorkflowVersionID string                   `json:"workflowVersionID,omitempty"`
}

type ListWorkflowFilesFilter struct {
	IDs []string
}

type ListFilesHandler interface {
	Handle(context.Context, *ListFilesQuery) ([]*WorkflowFile, int, error)
}

func NewListFilesHandler(readModel ReadModel, workspaceReadModel workspace.WorkspaceReadModel) ListFilesHandler {
	return &listFilesHandler{
		readModel:          readModel,
		workspaceReadModel: workspaceReadModel,
	}
}

type listFilesHandler struct {
	readModel          ReadModel
	workspaceReadModel workspace.WorkspaceReadModel
}

var _ ListFilesHandler = &listFilesHandler{}

func (l listFilesHandler) Handle(ctx context.Context, query *ListFilesQuery) ([]*WorkflowFile, int, error) {
	if err := validator.Validate(query); err != nil {
		return nil, 0, err
	}
	// check workspace exist
	if _, err := l.workspaceReadModel.GetWorkspaceById(ctx, query.WorkspaceID); err != nil {
		return nil, 0, err
	}
	workflow, err := l.readModel.GetById(ctx, query.WorkspaceID, query.WorkflowID)
	if err != nil {
		return nil, 0, err
	}
	// get workflow by id
	workflowVersionID := workflow.LatestVersion.ID
	if query.WorkflowVersionID != "" {
		workflowVersionID = query.WorkflowVersionID
	}
	return l.readModel.ListFiles(ctx, workflowVersionID, query.Pg, query.Filter)
}
