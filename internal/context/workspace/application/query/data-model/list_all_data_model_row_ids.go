package datamodel

import (
	"context"

	workspacequery "github.com/Bio-OS/bioos/internal/context/workspace/application/query/workspace"
	"github.com/Bio-OS/bioos/pkg/utils"
	"github.com/Bio-OS/bioos/pkg/validator"
)

type ListAllDataModelRowIDsHandler interface {
	Handle(ctx context.Context, query *ListAllDataModelRowIDsQuery) ([]string, error)
}

type listAllDataModelRowIDsHandler struct {
	workspaceReadModel workspacequery.WorkspaceReadModel
	dataModelReadModel DataModelReadModel
}

var _ ListAllDataModelRowIDsHandler = &listAllDataModelRowIDsHandler{}

func NewListAllDataModelRowIDsHandler(workspaceReadModel workspacequery.WorkspaceReadModel, dataModelReadModel DataModelReadModel) ListAllDataModelRowIDsHandler {
	return &listAllDataModelRowIDsHandler{
		workspaceReadModel,
		dataModelReadModel,
	}
}

func (l *listAllDataModelRowIDsHandler) Handle(ctx context.Context, query *ListAllDataModelRowIDsQuery) ([]string, error) {
	if err := validator.Validate(query); err != nil {
		return nil, err
	}

	if err := workspacequery.CheckWorkspaceExist(ctx, l.workspaceReadModel, query.WorkspaceID); err != nil {
		return nil, err
	}
	dataModelName, err := l.dataModelReadModel.GetDataModelName(ctx, query.WorkspaceID, query.ID)
	if err != nil {
		return nil, err
	}
	typ := utils.GetDataModelType(dataModelName)
	return l.dataModelReadModel.ListAllDataModelRowIDs(ctx, query.ID, typ)
}
