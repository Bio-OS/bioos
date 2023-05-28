package handlers

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	datamodelcommand "github.com/Bio-OS/bioos/internal/context/workspace/application/command/data-model"
	datamodelquery "github.com/Bio-OS/bioos/internal/context/workspace/application/query/data-model"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	applog "github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/utils"
)

// PatchDataModel patch data model
//
//	@Summary		use to patch data model
//	@Description	patch data model
//	@Tags			datamodel
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/workspace/{workspace_id}/data_model  [patch]
//	@Security		basicAuth
//	@Param			request	body		PatchDataModelRequest	true	"patch data model request"
//	@Success		201		{object}	PatchDataModelResponse
//	@Failure		400		{object}	apperrors.AppError	"invalid param"
//	@Failure		401		{object}	apperrors.AppError	"unauthorized"
//	@Failure		403		{object}	apperrors.AppError	"forbidden"
//	@Failure		500		{object}	apperrors.AppError	"internal system error"
func PatchDataModel(ctx context.Context, c *app.RequestContext, handler datamodelcommand.PatchDataModelHandler) {
	var req PatchDataModelRequest
	err := c.Bind(&req)
	if err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	cmd := patchDataModelVoToDto(req)
	id, err := handler.Handle(ctx, cmd)
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}

	resp := &PatchDataModelResponse{ID: id}
	utils.WriteHertzCreatedResponse(c, resp)
}

// DeleteDataModel delete data model
//
//	@Summary		use to delete data model,support delete with data model name/row ids/headers
//	@Description	delete data model
//	@Tags			datamodel
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/workspace/{workspace_id}/data_model/{id} [delete]
//	@Security		basicAuth
//	@Param			workspace_id	path	string		true	"get workspace id"
//	@Param			id				path	string		true	"get data model id"
//	@Param			headers			query	[]string	false	"the data model headers should delete"
//	@Param			rowIDs			query	[]string	false	"the data model row ids should delete"
//	@Success		202
//	@Failure		400	{object}	apperrors.AppError	"invalid param"
//	@Failure		401	{object}	apperrors.AppError	"unauthorized"
//	@Failure		403	{object}	apperrors.AppError	"forbidden"
//	@Failure		500	{object}	apperrors.AppError	"internal system error"
func DeleteDataModel(ctx context.Context, c *app.RequestContext, handler datamodelcommand.DeleteDataModelHandler) {
	var req DeleteDataModelRequest
	err := c.Bind(&req)
	if err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	cmd := deleteDataModelVoToDto(req)
	err = handler.Handle(ctx, cmd)
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}

	utils.WriteHertzAcceptedResponse(c)
}

// ListDataModels list data models
//
//	@Summary		use to list data models
//	@Description	list data models
//	@Tags			datamodel
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/workspace/{workspace_id}/data_model [get]
//	@Security		basicAuth
//	@Param			workspace_id	path		string		true	"get workspace id"
//	@Param			types			query		[]string	false	"data model types"
//	@Param			searchWord		query		string		false	"query searchWord"
//	@Param			ids				query		[]string	false	"data model ids"
//	@Success		200				{object}	ListDataModelsResponse
//	@Failure		400				{object}	apperrors.AppError	"invalid param"
//	@Failure		401				{object}	apperrors.AppError	"unauthorized"
//	@Failure		403				{object}	apperrors.AppError	"forbidden"
//	@Failure		500				{object}	apperrors.AppError	"internal system error"
func ListDataModels(ctx context.Context, c *app.RequestContext, handler datamodelquery.ListDataModelsHandler) {
	var req ListDataModelsRequest
	err := c.Bind(&req)
	if err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	query := listDataModelsVoToDto(req)
	dataModels, err := handler.Handle(ctx, query)
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}

	items := make([]*DataModel, 0, len(dataModels))
	for _, dm := range dataModels {
		items = append(items, dataModelDtoToVo(dm))
	}
	resp := &ListDataModelsResponse{
		Items: items,
	}
	utils.WriteHertzOKResponse(c, resp)
}

// GetDataModel get data model
//
//	@Summary		use to get data model
//	@Description	get data model
//	@Tags			datamodel
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/workspace/{workspace_id}/data_model/{id} [get]
//	@Security		basicAuth
//	@Param			workspace_id	path		string	true	"get workspace id"
//	@Param			id				path		string	true	"get data model id"
//	@Success		200				{object}	GetDataModelResponse
//	@Failure		400				{object}	apperrors.AppError	"invalid param"
//	@Failure		401				{object}	apperrors.AppError	"unauthorized"
//	@Failure		403				{object}	apperrors.AppError	"forbidden"
//	@Failure		500				{object}	apperrors.AppError	"internal system error"
func GetDataModel(ctx context.Context, c *app.RequestContext, handler datamodelquery.GetDataModelHandler) {
	var req GetDataModelRequest
	err := c.Bind(&req)
	if err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	query := getDataModelVoToDto(req)
	dataModel, headers, err := handler.Handle(ctx, query)
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}

	resp := &GetDataModelResponse{
		DataModel: dataModelDtoToVo(dataModel),
		Headers:   headers,
	}
	utils.WriteHertzOKResponse(c, resp)
}

// ListDataModelRows list data model rows
//
//	@Summary		use to list data model rows
//	@Description	list data model rows
//	@Tags			datamodel
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/workspace/{workspace_id}/data_model/{id}/rows [get]
//	@Security		basicAuth
//	@Param			workspace_id	path		string		true	"get workspace id"
//	@Param			id				path		string		true	"get data model id"
//	@Param			page			query		int			false	"query page"
//	@Param			size			query		int			false	"query size"
//	@Param			orderBy			query		string		false	"query order, just like field1,field2:desc"
//	@Param			inSetIDs		query		[]string	false	"data model entity set reffed entity row ids"
//	@Param			searchWord		query		string		false	"query searchWord"
//	@Param			rowIDs			query		[]string	false	"data model row ids"
//	@Success		200				{object}	ListDataModelRowsResponse
//	@Failure		400				{object}	apperrors.AppError	"invalid param"
//	@Failure		401				{object}	apperrors.AppError	"unauthorized"
//	@Failure		403				{object}	apperrors.AppError	"forbidden"
//	@Failure		500				{object}	apperrors.AppError	"internal system error"
func ListDataModelRows(ctx context.Context, c *app.RequestContext, handler datamodelquery.ListDataModelRowsHandler) {
	var req ListDataModelRowsRequest
	err := c.Bind(&req)
	if err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	query, err := listDataModelRowsVoToDto(req)
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}
	headers, rows, total, err := handler.Handle(ctx, query)
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}

	resp := &ListDataModelRowsResponse{
		Headers: headers,
		Rows:    rows,
		Size:    int32(query.Pagination.Size),
		Page:    int32(query.Pagination.Page),
		Total:   total,
	}
	utils.WriteHertzOKResponse(c, resp)
}

// ListAllDataModelRowIDs list all data model row ids
//
//	@Summary		use to list all data model row ids
//	@Description	list all data model row ids
//	@Tags			datamodel
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/workspace/{workspace_id}/data_model/{id}/rows/ids [get]
//	@Security		basicAuth
//	@Param			workspace_id	path		string	true	"get workspace id"
//	@Param			id				path		string	true	"get data model id"
//	@Success		200				{object}	ListAllDataModelRowIDsResponse
//	@Failure		400				{object}	apperrors.AppError	"invalid param"
//	@Failure		401				{object}	apperrors.AppError	"unauthorized"
//	@Failure		403				{object}	apperrors.AppError	"forbidden"
//	@Failure		500				{object}	apperrors.AppError	"internal system error"
func ListAllDataModelRowIDs(ctx context.Context, c *app.RequestContext, handler datamodelquery.ListAllDataModelRowIDsHandler) {
	var req ListAllDataModelRowIDsRequest
	err := c.Bind(&req)
	if err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	query := listAllDataModelRowIDsVoToDto(req)
	ids, err := handler.Handle(ctx, query)
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}

	resp := &ListAllDataModelRowIDsResponse{
		RowIDs: ids,
	}
	utils.WriteHertzOKResponse(c, resp)
}
