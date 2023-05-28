package datamodel

import (
	"context"
	"strconv"

	workspacequery "github.com/Bio-OS/bioos/internal/context/workspace/application/query/workspace"
	"github.com/Bio-OS/bioos/pkg/consts"
	"github.com/Bio-OS/bioos/pkg/utils"
	"github.com/Bio-OS/bioos/pkg/validator"
)

type ListDataModelRowsHandler interface {
	Handle(ctx context.Context, query *ListDataModelRowsQuery) ([]string, [][]string, int64, error)
}

type listDataModelRowsHandler struct {
	workspaceReadModel workspacequery.WorkspaceReadModel
	dataModelReadModel DataModelReadModel
}

var _ ListDataModelRowsHandler = &listDataModelRowsHandler{}

func NewListDataModelRowsHandler(workspaceReadModel workspacequery.WorkspaceReadModel, dataModelReadModel DataModelReadModel) ListDataModelRowsHandler {
	return &listDataModelRowsHandler{
		workspaceReadModel,
		dataModelReadModel,
	}
}

func (l *listDataModelRowsHandler) Handle(ctx context.Context, query *ListDataModelRowsQuery) ([]string, [][]string, int64, error) {
	if err := validator.Validate(query); err != nil {
		return nil, nil, 0, err
	}

	if err := workspacequery.CheckWorkspaceExist(ctx, l.workspaceReadModel, query.WorkspaceID); err != nil {
		return nil, nil, 0, err
	}
	dataModelName, err := l.dataModelReadModel.GetDataModelName(ctx, query.WorkspaceID, query.ID)
	if err != nil {
		return nil, nil, 0, err
	}
	typ := utils.GetDataModelType(dataModelName)
	headers, err := l.dataModelReadModel.ListDataModelHeaders(ctx, query.ID, dataModelName, typ)
	if err != nil {
		return nil, nil, 0, err
	}
	var order *utils.Order
	switch typ {
	case consts.DataModelTypeEntity:
		order = &utils.Order{
			Field:     utils.GenDataModelHeaderOfID(dataModelName),
			Ascending: true,
		}
		if len(query.Pagination.Orders) != 0 {
			order = &query.Pagination.Orders[0]
		}
		for index, header := range headers {
			if order.Field == header {
				order.Field = strconv.Itoa(index)
			}
		}
	case consts.DataModelTypeEntitySet:
		order = &utils.Order{
			Field:     "row_id",
			Ascending: true,
		}
		for _, pgOrder := range query.Pagination.Orders {
			if pgOrder.Field == utils.GenDataModelHeaderOfID(dataModelName) {
				order.Ascending = pgOrder.Ascending
			}
		}
	}
	rows, count, err := l.dataModelReadModel.ListDataModelRows(ctx, query.ID, typ, query.Pagination, order, query.Filter)
	if err != nil {
		return nil, nil, 0, err
	}
	return headers, rows, count, nil
}
