package notebook

import (
	"context"
	"fmt"

	"github.com/Bio-OS/bioos/internal/context/workspace/application/query/workspace"
	"github.com/Bio-OS/bioos/pkg/errors"
	"github.com/Bio-OS/bioos/pkg/validator"
)

type ListQuery struct {
	WorkspaceID string `validate:"required"`
}

type ListHandler interface {
	Handle(context.Context, *ListQuery) ([]*Notebook, error)
}

type listHandler struct {
	readModel          ReadModel
	workspaceReadModel workspace.WorkspaceReadModel
}

func NewListHandler(readModel ReadModel, workspaceReadModel workspace.WorkspaceReadModel) ListHandler {
	return &listHandler{
		readModel:          readModel,
		workspaceReadModel: workspaceReadModel,
	}
}

func (h *listHandler) Handle(ctx context.Context, query *ListQuery) ([]*Notebook, error) {
	if err := validator.Validate(query); err != nil {
		return nil, err
	}
	if err := workspace.CheckWorkspaceExist(ctx, h.workspaceReadModel, query.WorkspaceID); err != nil {
		return nil, err
	}
	res, err := h.readModel.ListByWorkspace(ctx, query.WorkspaceID)
	if err != nil {
		return nil, errors.NewInternalError(fmt.Errorf("find notebook fail: %w", err))
	}
	return res, nil
}
