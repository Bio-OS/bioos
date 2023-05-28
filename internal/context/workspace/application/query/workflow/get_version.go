package workflow

import (
	"context"

	"github.com/Bio-OS/bioos/internal/context/workspace/application/query/workspace"
	"github.com/Bio-OS/bioos/pkg/validator"
)

type GetVersionQuery struct {
	ID          string `validate:"required"`
	WorkspaceID string `validate:"required"`
	WorkflowID  string `validate:"required"`
}

type GetVersionHandler interface {
	Handle(context.Context, *GetVersionQuery) (*WorkflowVersion, error)
}

type getVersionHandler struct {
	readModel          ReadModel
	workspaceReadModel workspace.WorkspaceReadModel
}

var _ GetVersionHandler = &getVersionHandler{}

func NewGetVersionHandler(readModel ReadModel, workspaceReadModel workspace.WorkspaceReadModel) GetVersionHandler {
	return &getVersionHandler{
		readModel:          readModel,
		workspaceReadModel: workspaceReadModel,
	}
}

func (h *getVersionHandler) Handle(ctx context.Context, query *GetVersionQuery) (*WorkflowVersion, error) {
	if err := validator.Validate(query); err != nil {
		return nil, err
	}
	// check workspace exist
	if _, err := h.workspaceReadModel.GetWorkspaceById(ctx, query.WorkspaceID); err != nil {
		return nil, err
	}
	// check workflow exist
	if _, err := h.readModel.GetById(ctx, query.WorkspaceID, query.WorkflowID); err != nil {
		return nil, err
	}
	// get workflow by id
	res, err := h.readModel.GetVersion(ctx, query.ID)
	if err != nil {
		return nil, err
	}
	return res, nil
}
