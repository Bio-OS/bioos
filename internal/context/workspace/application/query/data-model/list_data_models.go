package datamodel

import (
	"context"

	workspacequery "github.com/Bio-OS/bioos/internal/context/workspace/application/query/workspace"
	"github.com/Bio-OS/bioos/pkg/validator"
)

type ListDataModelsHandler interface {
	Handle(ctx context.Context, query *ListDataModelsQuery) ([]*DataModel, error)
}

type listDataModelsHandler struct {
	workspaceReadModel workspacequery.WorkspaceReadModel
	dataModelReadModel DataModelReadModel
}

var _ ListDataModelsHandler = &listDataModelsHandler{}

func NewListDataModelsHandler(workspaceReadModel workspacequery.WorkspaceReadModel, dataModelReadModel DataModelReadModel) ListDataModelsHandler {
	return &listDataModelsHandler{
		workspaceReadModel,
		dataModelReadModel,
	}
}

func (l *listDataModelsHandler) Handle(ctx context.Context, query *ListDataModelsQuery) ([]*DataModel, error) {
	if err := validator.Validate(query); err != nil {
		return nil, err
	}

	if err := workspacequery.CheckWorkspaceExist(ctx, l.workspaceReadModel, query.WorkspaceID); err != nil {
		return nil, err
	}
	return l.dataModelReadModel.ListDataModels(ctx, query.WorkspaceID, query.Filter)
}
