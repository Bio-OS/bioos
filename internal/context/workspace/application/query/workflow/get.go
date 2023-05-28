package workflow

import (
	"context"

	"github.com/Bio-OS/bioos/internal/context/workspace/application/query/workspace"
	"github.com/Bio-OS/bioos/pkg/validator"
)

type GetQuery struct {
	ID          string `validate:"required"`
	WorkspaceID string `validate:"required"`
}

type GetHandler interface {
	Handle(context.Context, *GetQuery) (*Workflow, error)
}

type getHandler struct {
	readModel          ReadModel
	workspaceReadModel workspace.WorkspaceReadModel
}

var _ GetHandler = &getHandler{}

func NewGetHandler(readModel ReadModel, workspaceReadModel workspace.WorkspaceReadModel) GetHandler {
	return &getHandler{
		readModel:          readModel,
		workspaceReadModel: workspaceReadModel,
	}
}

func (h *getHandler) Handle(ctx context.Context, query *GetQuery) (*Workflow, error) {
	if err := validator.Validate(query); err != nil {
		return nil, err
	}
	// check workspace exist
	if _, err := h.workspaceReadModel.GetWorkspaceById(ctx, query.WorkspaceID); err != nil {
		return nil, err
	}
	// get workflow by id
	res, err := h.readModel.GetById(ctx, query.WorkspaceID, query.ID)
	if err != nil {
		return nil, err
	}
	return res, nil
}
