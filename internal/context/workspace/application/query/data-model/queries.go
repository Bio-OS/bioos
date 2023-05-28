package datamodel

import (
	workspacequery "github.com/Bio-OS/bioos/internal/context/workspace/application/query/workspace"
	"github.com/Bio-OS/bioos/pkg/utils"
)

type GetDataModelQuery struct {
	WorkspaceID string `validate:"required"`
	ID          string `validate:"required"`
}

type ListDataModelsQuery struct {
	WorkspaceID string `validate:"required"`
	Filter      *ListDataModelsFilter
}

type ListDataModelRowsQuery struct {
	WorkspaceID string `validate:"required"`
	ID          string `validate:"required"`
	Pagination  *utils.Pagination
	Filter      *ListDataModelRowsFilter
}

type ListAllDataModelRowIDsQuery struct {
	WorkspaceID string `validate:"required"`
	ID          string `validate:"required"`
}

type ListDataModelsFilter struct {
	Types      []string
	SearchWord string
	Exact      bool
	IDs        []string
}

type ListDataModelRowsFilter struct {
	SearchWord string
	InSetIDs   []string
	RowIDs     []string
}

type DataModel struct {
	ID          string
	Name        string
	RowCount    int64
	Type        string
	WorkspaceID string
}

type Queries struct {
	GetDataModel           GetDataModelHandler
	ListDataModels         ListDataModelsHandler
	ListDataModelRows      ListDataModelRowsHandler
	ListAllDataModelRowIDs ListAllDataModelRowIDsHandler
}

func NewQueries(workspaceReadModel workspacequery.WorkspaceReadModel, dataModelReadModel DataModelReadModel) *Queries {
	return &Queries{
		GetDataModel:           NewGetDataModelHandler(workspaceReadModel, dataModelReadModel),
		ListDataModels:         NewListDataModelsHandler(workspaceReadModel, dataModelReadModel),
		ListDataModelRows:      NewListDataModelRowsHandler(workspaceReadModel, dataModelReadModel),
		ListAllDataModelRowIDs: NewListAllDataModelRowIDsHandler(workspaceReadModel, dataModelReadModel),
	}
}
