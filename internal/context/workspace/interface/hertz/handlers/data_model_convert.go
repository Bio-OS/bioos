package handlers

import (
	datamodelcommand "github.com/Bio-OS/bioos/internal/context/workspace/application/command/data-model"
	datamodelquery "github.com/Bio-OS/bioos/internal/context/workspace/application/query/data-model"
	"github.com/Bio-OS/bioos/pkg/utils"
)

func patchDataModelVoToDto(req PatchDataModelRequest) *datamodelcommand.PatchDataModelCommand {
	return &datamodelcommand.PatchDataModelCommand{
		WorkspaceID: req.WorkspaceID,
		Name:        req.Name,
		Async:       req.Async,
		Headers:     req.Headers,
		Rows:        req.Rows,
	}
}

func deleteDataModelVoToDto(req DeleteDataModelRequest) *datamodelcommand.DeleteDataModelCommand {
	return &datamodelcommand.DeleteDataModelCommand{
		ID:          req.ID,
		WorkspaceID: req.WorkspaceID,
		Headers:     req.Headers,
		RowIDs:      req.RowIDs,
	}
}

func listDataModelsVoToDto(req ListDataModelsRequest) *datamodelquery.ListDataModelsQuery {
	return &datamodelquery.ListDataModelsQuery{
		WorkspaceID: req.WorkspaceID,
		Filter: &datamodelquery.ListDataModelsFilter{
			Types:      req.Types,
			SearchWord: req.SearchWord,
			IDs:        req.IDs,
		},
	}
}

func getDataModelVoToDto(req GetDataModelRequest) *datamodelquery.GetDataModelQuery {
	return &datamodelquery.GetDataModelQuery{
		WorkspaceID: req.WorkspaceID,
		ID:          req.ID,
	}
}

func listDataModelRowsVoToDto(req ListDataModelRowsRequest) (*datamodelquery.ListDataModelRowsQuery, error) {
	pg := utils.NewPagination(int(req.Size), int(req.Page))
	if err := pg.SetOrderBy(req.OrderBy); err != nil {
		return nil, err
	}
	return &datamodelquery.ListDataModelRowsQuery{
		WorkspaceID: req.WorkspaceID,
		ID:          req.ID,
		Pagination:  pg,
		Filter: &datamodelquery.ListDataModelRowsFilter{
			SearchWord: req.SearchWord,
			InSetIDs:   req.InSetIDs,
			RowIDs:     req.RowIDs,
		},
	}, nil
}

func listAllDataModelRowIDsVoToDto(req ListAllDataModelRowIDsRequest) *datamodelquery.ListAllDataModelRowIDsQuery {
	return &datamodelquery.ListAllDataModelRowIDsQuery{
		WorkspaceID: req.WorkspaceID,
		ID:          req.ID,
	}
}

func dataModelDtoToVo(dataModel *datamodelquery.DataModel) *DataModel {
	return &DataModel{
		dataModel.ID,
		dataModel.Name,
		dataModel.RowCount,
		dataModel.Type,
	}
}

/*
func rowDtoToVo(row *datamodelquery.Row) *Row {
	grids := make([]*Grid, 0, len(row.Grids))
	for _, gird := range row.Grids {
		grids = append(grids, gridDtoToVo(gird))
	}
	return &Row{
		Grid: grids,
	}
}

func gridDtoToVo(grid *datamodelquery.Grid) *Grid {
	return &Grid{
		grid.Value,
		grid.Type,
	}
}
*/
