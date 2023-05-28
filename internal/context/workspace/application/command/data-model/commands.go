package datamodel

import (
	datamodelquery "github.com/Bio-OS/bioos/internal/context/workspace/application/query/data-model"
	workspacequery "github.com/Bio-OS/bioos/internal/context/workspace/application/query/workspace"
	datamodel "github.com/Bio-OS/bioos/internal/context/workspace/domain/data-model"
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
)

type PatchDataModelCommand struct {
	WorkspaceID string `validate:"required"`
	Name        string `validate:"required,dataModelName"`
	Async       bool
	Headers     []string   `validate:"required,dataModelHeaders"`
	Rows        [][]string `validate:"required,dataModelRows"`
}

type DeleteDataModelCommand struct {
	ID          string   `validate:"required"`
	WorkspaceID string   `validate:"required"`
	Headers     []string `validate:"deleteDataModelHeaders"`
	RowIDs      []string
}

type Commands struct {
	PatchDataModel  PatchDataModelHandler
	DeleteDataModel DeleteDataModelHandler
}

func NewCommands(dataModelRepo datamodel.Repository, workspaceReadModel workspacequery.WorkspaceReadModel, dataModelFactory *datamodel.Factory, dataModelReadModel datamodelquery.DataModelReadModel, eventBus eventbus.EventBus) *Commands {
	svc := datamodel.NewService(dataModelRepo, eventBus, dataModelFactory)
	return &Commands{
		PatchDataModel:  NewPatchDataModelHandler(svc, workspaceReadModel, dataModelReadModel),
		DeleteDataModel: NewDeleteDataModelHandler(svc, workspaceReadModel, dataModelReadModel),
	}
}
