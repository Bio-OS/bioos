package workflow

import (
	"context"

	"github.com/Bio-OS/bioos/internal/context/workspace/application/query/workspace"
	"github.com/Bio-OS/bioos/pkg/utils"
	"github.com/Bio-OS/bioos/pkg/validator"
)

const (
	OrderByName       = "name"
	OrderByCreateTime = "createdAt"
)

type ListQuery struct {
	Pg          *utils.Pagination
	Filter      *ListWorkflowsFilter
	WorkspaceID string
}

type ListWorkflowsFilter struct {
	SearchWord string
	IDs        []string
	Exact      bool
}

type ListHandler interface {
	Handle(context.Context, *ListQuery) ([]*Workflow, int, error)
}

type listHandler struct {
	readModel          ReadModel
	workspaceReadModel workspace.WorkspaceReadModel
}

var _ ListHandler = &listHandler{}

func NewListHandler(readModel ReadModel, workspaceReadModel workspace.WorkspaceReadModel) ListHandler {
	return &listHandler{
		readModel:          readModel,
		workspaceReadModel: workspaceReadModel,
	}
}

func (h *listHandler) Handle(ctx context.Context, query *ListQuery) ([]*Workflow, int, error) {
	if err := validator.Validate(query); err != nil {
		return nil, 0, err
	}
	// check if workspace exist
	if _, err := h.workspaceReadModel.GetWorkspaceById(ctx, query.WorkspaceID); err != nil {
		return nil, 0, err
	}

	// list workflows
	return h.readModel.List(ctx, query.WorkspaceID, query.Pg, query.Filter)
}
