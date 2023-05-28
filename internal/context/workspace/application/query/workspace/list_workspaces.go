package workspace

import (
	"context"

	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	"github.com/Bio-OS/bioos/pkg/validator"
)

type ListWorkspacesHandler interface {
	Handle(ctx context.Context, query *ListWorkspacesQuery) ([]*WorkspaceItem, int, error)
}

type listWorkspacesHandler struct {
	workspaceReadModel WorkspaceReadModel
}

func NewListWorkspacesHandler(workspaceReadModel WorkspaceReadModel) ListWorkspacesHandler {
	return &listWorkspacesHandler{workspaceReadModel: workspaceReadModel}
}

func (q *listWorkspacesHandler) Handle(ctx context.Context, query *ListWorkspacesQuery) ([]*WorkspaceItem, int, error) {
	if err := validator.Validate(query); err != nil {
		return nil, 0, err
	}

	for _, order := range query.Pg.Orders {
		if order.Field != OrderByName && order.Field != OrderByCreateTime {
			return nil, 0, apperrors.NewInvalidError("orderField")
		}
	}

	ws, err := q.workspaceReadModel.ListWorkspaces(ctx, query.Pg, query.Filter)
	if err != nil {
		return nil, 0, err
	}
	cnt, err := q.workspaceReadModel.CountWorkspaces(ctx, query.Filter)
	if err != nil {
		return nil, 0, err
	}
	return ws, cnt, nil
}
