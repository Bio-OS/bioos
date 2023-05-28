package datamodel

import (
	"context"

	"github.com/Bio-OS/bioos/pkg/utils"
)

type DataModelReadModel interface {
	ListDataModels(ctx context.Context, workspaceID string, filter *ListDataModelsFilter) ([]*DataModel, error)

	CountDataModel(ctx context.Context, workspaceID string, filter *ListDataModelsFilter) (int64, error)

	GetDataModelName(ctx context.Context, workspaceID, id string) (string, error)

	GetDataModelWithID(ctx context.Context, id string) (*DataModel, error)

	GetDataModelWithName(ctx context.Context, workspaceID, name string) (*DataModel, error)

	ListDataModelHeaders(ctx context.Context, id, name, _type string) ([]string, error)
	ListEntityDataModelHeaders(ctx context.Context, id string) ([]string, error)

	ListEntityDataModelColumnsWithRowIDs(ctx context.Context, id string, headers []string, rowIDs []string) (map[string][]string, error)

	ListDataModelRows(ctx context.Context, id, _type string, pagination *utils.Pagination, order *utils.Order, filter *ListDataModelRowsFilter) ([][]string, int64, error)

	ListAllDataModelRowIDs(ctx context.Context, id, _type string) ([]string, error)
	CountDataModelRows(ctx context.Context, id, _type string, filter *ListDataModelRowsFilter) (int64, error)
}
