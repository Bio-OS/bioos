package workspace

import (
	"context"

	"github.com/Bio-OS/bioos/pkg/validator"
)

type GetWorkspaceByIDQueryHandler interface {
	Handle(ctx context.Context, command *GetWorkspaceByIDQuery) (*WorkspaceItem, error)
}

type getWorkspaceByIDHandler struct {
	workspaceReadModel WorkspaceReadModel
}

func NewGetWorkspaceByIDHandler(workspaceReadModel WorkspaceReadModel) GetWorkspaceByIDQueryHandler {
	return &getWorkspaceByIDHandler{workspaceReadModel: workspaceReadModel}
}

func (q *getWorkspaceByIDHandler) Handle(ctx context.Context, query *GetWorkspaceByIDQuery) (*WorkspaceItem, error) {
	if err := validator.Validate(query); err != nil {
		return nil, err
	}

	ws, err := q.workspaceReadModel.GetWorkspaceById(ctx, query.ID)
	if err != nil {
		return nil, err
	}
	return ws, nil
}
