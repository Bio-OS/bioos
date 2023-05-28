package workflow

import (
	"context"

	"github.com/Bio-OS/bioos/internal/context/workspace/application/query/workspace"
	"github.com/Bio-OS/bioos/pkg/validator"
)

type GetFileQuery struct {
	ID          string `validate:"required"`
	WorkspaceID string `validate:"required"`
	WorkflowID  string `validate:"required"`
}

type GetFileHandler interface {
	Handle(context.Context, *GetFileQuery) (*WorkflowFile, error)
}

type getFileHandler struct {
	readModel          ReadModel
	workspaceReadModel workspace.WorkspaceReadModel
}

var _ GetFileHandler = &getFileHandler{}

func NewGetFileHandler(readModel ReadModel, workspaceReadModel workspace.WorkspaceReadModel) GetFileHandler {
	return &getFileHandler{
		readModel:          readModel,
		workspaceReadModel: workspaceReadModel,
	}
}

func (h *getFileHandler) Handle(ctx context.Context, query *GetFileQuery) (*WorkflowFile, error) {
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
	res, err := h.readModel.GetFile(ctx, query.ID)
	if err != nil {
		return nil, err
	}
	return res, nil
}
