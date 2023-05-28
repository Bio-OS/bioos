package workflow

import (
	"context"

	"github.com/Bio-OS/bioos/internal/context/workspace/application/query/workspace"
	"github.com/Bio-OS/bioos/pkg/utils"
	"github.com/Bio-OS/bioos/pkg/validator"
)

const (
	VersionOrderByStatus   = "status"
	VersionOrderByLanguage = "language"
	VersionOrderBySource   = "source"
)

type ListVersionsHandler interface {
	Handle(context.Context, *ListVersionsQuery) ([]*WorkflowVersion, int, error)
}

type ListVersionsQuery struct {
	Pg          *utils.Pagination           `json:"pg"`
	WorkflowID  string                      `json:"workflowID"`
	WorkspaceID string                      `json:"workspaceID"`
	Filter      *ListWorkflowVersionsFilter `json:"filter"`
}

type ListWorkflowVersionsFilter struct {
	IDs []string
}

type listVersionsHandler struct {
	workspaceReadModel workspace.WorkspaceReadModel
	readModel          ReadModel
}

var _ ListVersionsHandler = &listVersionsHandler{}

func (l listVersionsHandler) Handle(ctx context.Context, query *ListVersionsQuery) ([]*WorkflowVersion, int, error) {
	if err := validator.Validate(query); err != nil {
		return nil, 0, err
	}
	// check if workspace exist
	if _, err := l.workspaceReadModel.GetWorkspaceById(ctx, query.WorkspaceID); err != nil {
		return nil, 0, err
	}
	// check if workflow exist
	if _, err := l.readModel.GetById(ctx, query.WorkspaceID, query.WorkflowID); err != nil {
		return nil, 0, err
	}
	return l.readModel.ListVersions(ctx, query.WorkflowID, query.Pg, query.Filter)
}

// NewListVersionsHandler new a ListVersionsHandler
func NewListVersionsHandler(readModel ReadModel, workspaceReadModel workspace.WorkspaceReadModel) ListVersionsHandler {
	return &listVersionsHandler{
		workspaceReadModel: workspaceReadModel,
		readModel:          readModel,
	}
}
