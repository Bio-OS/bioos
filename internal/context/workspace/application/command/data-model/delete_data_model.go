package datamodel

import (
	"context"

	datamodelquery "github.com/Bio-OS/bioos/internal/context/workspace/application/query/data-model"
	workspacequery "github.com/Bio-OS/bioos/internal/context/workspace/application/query/workspace"
	datamodel "github.com/Bio-OS/bioos/internal/context/workspace/domain/data-model"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	"github.com/Bio-OS/bioos/pkg/utils"
	"github.com/Bio-OS/bioos/pkg/validator"
)

type DeleteDataModelHandler interface {
	Handle(ctx context.Context, cmd *DeleteDataModelCommand) error
}

type deleteDataModelHandler struct {
	svc                datamodel.Service
	workspaceReadModel workspacequery.WorkspaceReadModel
	dataModelReadModel datamodelquery.DataModelReadModel
}

var _ DeleteDataModelHandler = &deleteDataModelHandler{}

func NewDeleteDataModelHandler(svc datamodel.Service, workspaceReadModel workspacequery.WorkspaceReadModel, dataModelReadModel datamodelquery.DataModelReadModel) DeleteDataModelHandler {
	return &deleteDataModelHandler{
		svc,
		workspaceReadModel,
		dataModelReadModel,
	}
}

func (d *deleteDataModelHandler) Handle(ctx context.Context, cmd *DeleteDataModelCommand) error {
	if err := validator.Validate(cmd); err != nil {
		return err
	}

	if err := workspacequery.CheckWorkspaceExist(ctx, d.workspaceReadModel, cmd.WorkspaceID); err != nil {
		return err
	}

	dm, err := d.svc.Get(ctx, cmd.ID)
	if err != nil {
		return err
	}
	if dm.WorkspaceID != cmd.WorkspaceID {
		return apperrors.NewInvalidError("data model[%s] is not belong to workspace[%s]", dm.Name, cmd.WorkspaceID)
	}
	if len(cmd.Headers) == 0 && len(cmd.RowIDs) == 0 {
		// delete the whole data model
		return d.svc.Delete(ctx, dm)
	}
	headers, err := d.dataModelReadModel.ListDataModelHeaders(ctx, dm.ID, dm.Name, dm.Type)
	if err != nil {
		return err
	}
	newHeaders := utils.DeleteStrSliceElms(headers, cmd.Headers...)

	// check the new rowIDs
	rowIDs, err := d.dataModelReadModel.ListAllDataModelRowIDs(ctx, cmd.ID, dm.Type)
	if err != nil {
		return err
	}
	newRowIDs := utils.DeleteStrSliceElms(rowIDs, cmd.RowIDs...)
	if len(newRowIDs) != 0 && utils.In(utils.GenDataModelHeaderOfID(dm.Name), newHeaders) {
		dm.Headers = newHeaders
		dm.RowIDs = newRowIDs
	}
	return d.svc.Delete(ctx, dm)
}
