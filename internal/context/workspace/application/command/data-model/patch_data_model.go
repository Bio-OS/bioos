package datamodel

import (
	"context"
	"errors"

	datamodelquery "github.com/Bio-OS/bioos/internal/context/workspace/application/query/data-model"
	workspacequery "github.com/Bio-OS/bioos/internal/context/workspace/application/query/workspace"
	datamodel "github.com/Bio-OS/bioos/internal/context/workspace/domain/data-model"
	"github.com/Bio-OS/bioos/pkg/consts"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	"github.com/Bio-OS/bioos/pkg/utils"
	"github.com/Bio-OS/bioos/pkg/validator"
)

type PatchDataModelHandler interface {
	Handle(ctx context.Context, cmd *PatchDataModelCommand) (string, error)
}

type patchDataModelHandler struct {
	svc                datamodel.Service
	workspaceReadModel workspacequery.WorkspaceReadModel
	dataModelReadModel datamodelquery.DataModelReadModel
	factory            *datamodel.Factory
}

var _ PatchDataModelHandler = &patchDataModelHandler{}

func NewPatchDataModelHandler(svc datamodel.Service, workspaceReadModel workspacequery.WorkspaceReadModel, dataModelReadModel datamodelquery.DataModelReadModel) PatchDataModelHandler {
	return &patchDataModelHandler{
		svc,
		workspaceReadModel,
		dataModelReadModel,
		datamodel.NewDataModelFactory(),
	}
}

func (p *patchDataModelHandler) Handle(ctx context.Context, cmd *PatchDataModelCommand) (string, error) {
	if err := validator.Validate(cmd); err != nil {
		return "", err
	}

	if err := workspacequery.CheckWorkspaceExist(ctx, p.workspaceReadModel, cmd.WorkspaceID); err != nil {
		return "", err
	}

	dataModelType := utils.GetDataModelType(cmd.Name)

	dm, err := p.dataModelReadModel.GetDataModelWithName(ctx, cmd.WorkspaceID, cmd.Name)
	if err != nil {
		var apperror apperrors.Error
		if errors.As(err, &apperror) && (apperror.GetCode() == apperrors.NotFoundCode) {
			newDataModel := p.factory.New(&datamodel.CreateParam{
				WorkspaceID: cmd.WorkspaceID,
				Name:        cmd.Name,
				Type:        dataModelType,
				Headers:     cmd.Headers,
				Rows:        cmd.Rows,
			})
			if err = p.svc.Create(ctx, newDataModel); err != nil {
				return "", err
			}
			return newDataModel.ID, nil
		}
		return "", err
	}
	model, err := p.svc.Get(ctx, dm.ID)
	if err != nil {
		return "", nil
	}

	if dm.WorkspaceID != cmd.WorkspaceID {
		return "", apperrors.NewInvalidError("data model[%s] is not belong to workspace[%s]", dm.Name, cmd.WorkspaceID)
	}
	headers := cmd.Headers
	rows := cmd.Rows
	if dataModelType == consts.DataModelTypeEntity {
		rowIDs := getRowIDs(cmd.Rows)

		dbHeaders, err := p.dataModelReadModel.ListEntityDataModelHeaders(ctx, model.ID)
		if err != nil {
			return "", err
		}

		dbColumns, err := p.dataModelReadModel.ListEntityDataModelColumnsWithRowIDs(ctx, model.ID, dbHeaders, rowIDs)
		if err != nil {
			return "", err
		}

		headers, rows = genNewHeadersAndRows(dbHeaders, cmd.Headers, dbColumns, cmd.Rows)
	}
	model.Headers = headers
	model.Rows = rows
	if err = p.svc.Upsert(ctx, model); err != nil {
		return "", err
	}
	return model.ID, nil
}

func rows2Columns(headers []string, rows [][]string) (columns map[string][]string) {
	columns = make(map[string][]string, 0)
	for index, header := range headers {
		column := make([]string, 0, len(rows))
		for _, row := range rows {
			column = append(column, row[index])
		}
		columns[header] = column
	}
	return columns
}

func columns2Rows(headers []string, columns map[string][]string) (rows [][]string) {
	for index, header := range headers {
		column := columns[header]
		if index == 0 {
			for _, _ = range column {
				row := make([]string, 0, len(headers))
				rows = append(rows, row)
			}
		}
		for columnIndex, grid := range column {
			rows[columnIndex] = append(rows[columnIndex], grid)
		}
	}
	return rows
}

func getRowIDs(rows [][]string) (rowIDs []string) {
	for _, row := range rows {
		rowIDs = append(rowIDs, row[0])
	}
	return rowIDs
}

func genNewHeadersAndRows(dbHeaders, fileHeaders []string, dbColumns map[string][]string, fileRows [][]string) (headers []string, rows [][]string) {
	fileColumns := rows2Columns(fileHeaders, fileRows)

	headers = append(headers, dbHeaders...)

	for _, header := range fileHeaders {
		if _, ok := dbColumns[header]; !ok {
			headers = append(headers, header)
		}
	}
	newMappedColumns := make(map[string][]string, 0)
	for header, column := range fileColumns {
		_, ok := dbColumns[header]
		if ok {
			delete(dbColumns, header)
		}
		newMappedColumns[header] = column
	}
	for header, column := range dbColumns {
		newMappedColumns[header] = column
	}
	rows = columns2Rows(headers, newMappedColumns)

	return headers, rows
}
