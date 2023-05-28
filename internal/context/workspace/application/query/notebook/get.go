package notebook

import (
	"context"

	"github.com/Bio-OS/bioos/internal/context/workspace/application/query/workspace"
	"github.com/Bio-OS/bioos/pkg/errors"
	"github.com/Bio-OS/bioos/pkg/validator"
)

type GetQuery struct {
	Name        string `validate:"required"`
	WorkspaceID string `validate:"required"`
}

type GetHandler interface {
	Handle(context.Context, *GetQuery) (*Notebook, error)
}

type getHandler struct {
	readModel          ReadModel
	workspaceReadModel workspace.WorkspaceReadModel
}

func NewGetHandler(readModel ReadModel, workspaceReadModel workspace.WorkspaceReadModel) GetHandler {
	return &getHandler{
		readModel:          readModel,
		workspaceReadModel: workspaceReadModel,
	}
}

func (h *getHandler) Handle(ctx context.Context, query *GetQuery) (*Notebook, error) {
	if err := validator.Validate(query); err != nil {
		return nil, err
	}
	if err := workspace.CheckWorkspaceExist(ctx, h.workspaceReadModel, query.WorkspaceID); err != nil {
		return nil, err
	}
	nb, err := h.readModel.Get(ctx, query.WorkspaceID, query.Name)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	if nb == nil {
		return nil, errors.NewNotFoundError("notebook", query.Name)
	}
	return nb, nil
}
