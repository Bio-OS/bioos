package sql

import (
	"context"
	"encoding/json"

	query "github.com/Bio-OS/bioos/internal/context/workspace/application/query/data-model"
	datamodel "github.com/Bio-OS/bioos/internal/context/workspace/domain/data-model"
	"github.com/Bio-OS/bioos/pkg/log"
)

func DataModelPOToDataModelDTO(ctx context.Context, d *DataModel, count int64) *query.DataModel {
	item := &query.DataModel{
		ID:          d.ID,
		Name:        d.Name,
		Type:        d.Type,
		WorkspaceID: d.WorkspaceID,
		RowCount:    count,
	}
	return item
}

func EntityHeadersPOToHeadersDTO(ctx context.Context, h *EntityHeader) string {
	return h.Name
}

func EntityGridPOToEntityGridDTO(ctx context.Context, e *EntityGrid) string {
	return e.Value
}

func EntityGridsPOToColumnDTO(ctx context.Context, e []*EntityGrid, rowIDs []string) []string {
	res := make([]string, len(rowIDs), len(rowIDs))
	mappedRowIDs := make(map[string]int, 0)
	for index, rowID := range rowIDs {
		mappedRowIDs[rowID] = index
	}
	for _, grid := range e {
		index := mappedRowIDs[grid.RowID]
		res[index] = EntityGridPOToEntityGridDTO(ctx, grid)
	}
	return res
}

func EntityGridPOToRowIDDTO(ctx context.Context, e *EntityGrid) string {
	return e.RowID
}

func EntityGridsPOToRowIDsDTO(ctx context.Context, e []*EntityGrid) []string {
	rowIDMapped := make(map[string]struct{})
	for _, grid := range e {
		rowIDMapped[EntityGridPOToRowIDDTO(ctx, grid)] = struct{}{}
	}
	res := make([]string, 0, len(e))
	for key, _ := range rowIDMapped {
		res = append(res, key)
	}
	return res
}

func EntitySetRowPOToRowIDDTO(ctx context.Context, e *EntitySetRow) string {
	return e.RowID
}

func EntitySetRowsPOToRowIDsDTO(ctx context.Context, e []*EntitySetRow) []string {
	res := make([]string, 0, len(e))
	for _, grid := range e {
		res = append(res, EntitySetRowPOToRowIDDTO(ctx, grid))
	}
	return res
}

func WorkspaceRowPOToRowIDDTO(ctx context.Context, w *WorkspaceRow) string {
	return w.Key
}

func WorkspaceRowsPOToRowIDsDTO(ctx context.Context, e []*WorkspaceRow) []string {
	res := make([]string, 0, len(e))
	for _, grid := range e {
		res = append(res, WorkspaceRowPOToRowIDDTO(ctx, grid))
	}
	return res
}

func EntitySetRowPOToEntitySetRowDTO(ctx context.Context, e *EntitySetRow, m map[string][]string) map[string][]string {
	row, ok := m[e.RowID]
	if !ok {
		m[e.RowID] = []string{e.RefRowID}
	} else {
		m[e.RowID] = append(row, e.RefRowID)
	}
	return m
}

func EntitySetRowsPOToEntitySetRowIDsDTO(ctx context.Context, e []*EntitySetRow) []string {
	rowIDMapped := make(map[string]struct{})
	for _, grid := range e {
		rowIDMapped[grid.RowID] = struct{}{}
	}
	res := make([]string, 0, len(e))
	for key, _ := range rowIDMapped {
		res = append(res, key)
	}
	return res
}

func EntitySetRowsPOToEntitySetRowsDTO(ctx context.Context, e []*EntitySetRow, rowIDs []string) [][]string {
	EntitySetRowMap := make(map[string][]string)
	res := make([][]string, 0, len(e))
	for _, row := range e {
		EntitySetRowMap = EntitySetRowPOToEntitySetRowDTO(ctx, row, EntitySetRowMap)
	}
	for _, index := range rowIDs {
		rowInBytes, err := json.Marshal(EntitySetRowMap[index])
		if err != nil {
			log.Fatalf(err.Error())
		}
		res = append(res, []string{index, string(rowInBytes)})
	}
	return res
}

func WorkspaceRowPOToWorkspaceRowDTO(ctx context.Context, e *WorkspaceRow) []string {
	res := make([]string, 0, 2)
	res = append(res, e.Key, e.Value)
	return res
}

func WorkspaceRowsPOToWorkspaceRowsDTO(ctx context.Context, e []*WorkspaceRow) [][]string {
	res := make([][]string, 0, len(e))
	for _, row := range e {
		res = append(res, WorkspaceRowPOToWorkspaceRowDTO(ctx, row))
	}
	return res
}

func DataModelPOToDataModelDO(ctx context.Context, model *DataModel) *datamodel.DataModel {
	return &datamodel.DataModel{
		WorkspaceID: model.WorkspaceID,
		ID:          model.ID,
		Name:        model.Name,
		Type:        model.Type,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

func DataModelDOtoDataModelPO(ctx context.Context, model *datamodel.DataModel) *DataModel {
	return &DataModel{
		WorkspaceID: model.WorkspaceID,
		ID:          model.ID,
		Name:        model.Name,
		Type:        model.Type,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

func DataModelDOtoEntitySetRowsPO(ctx context.Context, model *datamodel.DataModel) ([]*EntitySetRow, error) {
	workspaceRows := make([]*EntitySetRow, 0)
	for _, row := range model.Rows {
		var rowInList []string
		err := json.Unmarshal([]byte(row[1]), &rowInList)
		if err != nil {
			return nil, err
		}
		for _, grid := range rowInList {
			workspaceRow := &EntitySetRow{
				RowID:       row[0],
				DataModelID: model.ID,
				RefRowID:    grid,
			}
			workspaceRows = append(workspaceRows, workspaceRow)
		}
	}
	return workspaceRows, nil
}

func DataModelDOtoWorkspaceRowsPO(ctx context.Context, model *datamodel.DataModel) []*WorkspaceRow {
	workspaceRows := make([]*WorkspaceRow, 0)
	for _, row := range model.Rows {
		workspaceRow := &WorkspaceRow{
			Key:         row[0],
			DataModelID: model.ID,
			Value:       row[1],
			Type:        "string",
		}
		workspaceRows = append(workspaceRows, workspaceRow)
	}
	return workspaceRows
}

func DataModelDOtoEntityHeadersPO(ctx context.Context, model *datamodel.DataModel) []*EntityHeader {
	entityHeaders := make([]*EntityHeader, 0, len(model.Headers))
	for index, header := range model.Headers {
		entityHeaders = append(entityHeaders, &EntityHeader{
			ColumnIndex: index,
			DataModelID: model.ID,
			Name:        header,
			Type:        "string",
		})
	}
	return entityHeaders
}

func DataModelDOtoEntityGridsPO(ctx context.Context, model *datamodel.DataModel) []*EntityGrid {
	entityGrid := make([]*EntityGrid, 0)
	for _, row := range model.Rows {
		for index, grid := range row {
			rowID := row[0]
			entityGrid = append(entityGrid, &EntityGrid{
				RowID:       rowID,
				ColumnIndex: index,
				DataModelID: model.ID,
				Value:       grid,
			})
		}
	}
	return entityGrid
}
