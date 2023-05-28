package datamodel

import (
	"context"

	workspacequery "github.com/Bio-OS/bioos/internal/context/workspace/application/query/workspace"
	"github.com/Bio-OS/bioos/pkg/validator"
)

type GetDataModelHandler interface {
	Handle(ctx context.Context, query *GetDataModelQuery) (*DataModel, []string, error)
}

type getDataModelHandler struct {
	workspaceReadModel workspacequery.WorkspaceReadModel
	dataModelReadModel DataModelReadModel
}

var _ GetDataModelHandler = &getDataModelHandler{}

func NewGetDataModelHandler(workspaceReadModel workspacequery.WorkspaceReadModel, dataModelReadModel DataModelReadModel) GetDataModelHandler {
	return &getDataModelHandler{
		workspaceReadModel,
		dataModelReadModel,
	}
}

func (g *getDataModelHandler) Handle(ctx context.Context, query *GetDataModelQuery) (*DataModel, []string, error) {
	if err := validator.Validate(query); err != nil {
		return nil, nil, err
	}

	if err := workspacequery.CheckWorkspaceExist(ctx, g.workspaceReadModel, query.WorkspaceID); err != nil {
		return nil, nil, err
	}
	model, err := g.dataModelReadModel.GetDataModelWithID(ctx, query.ID)
	if err != nil {
		return nil, nil, err
	}
	headers, err := g.dataModelReadModel.ListDataModelHeaders(ctx, model.ID, model.Name, model.Type)
	if err != nil {
		return nil, nil, err
	}
	return model, headers, nil
}
