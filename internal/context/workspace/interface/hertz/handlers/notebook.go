package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"

	command "github.com/Bio-OS/bioos/internal/context/workspace/application/command/notebook"
	query "github.com/Bio-OS/bioos/internal/context/workspace/application/query/notebook"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	applog "github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/utils"
)

// CreateNotebook create or update notebook
//
//	@Summary		use to create or update notebook
//	@Description	create notebook, update if name exist, set ipynb content in http body
//	@Tags			notebook
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/workspace/{workspace-id}/notebook/{name} [put]
//	@Security		basicAuth
//	@Param			request	body		notebook.IPythonNotebook	true	"ipynb content"
//	@Failure		400		{object}	apperrors.AppError			"invalid param"
//	@Failure		401		{object}	apperrors.AppError			"unauthorized"
//	@Failure		403		{object}	apperrors.AppError			"forbidden"
//	@Failure		500		{object}	apperrors.AppError			"internal system error"
func CreateNotebook(ctx context.Context, c *app.RequestContext, handler command.CreateHandler) {
	var req createNotebookRequest
	if err := c.Bind(&req); err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	if err := handler.Handle(ctx, req.toDTO()); err != nil {
		utils.WriteHertzErrorResponse(c, err)
	} else {
		utils.WriteHertzCreatedResponse(c, nil)
	}
}

// GetNotebook get notebook content
//
//	@Summary		get notebook content
//	@Description	get notebook content
//	@Tags			notebook
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/workspace/{workspace-id}/notebook/{name} [get]
//	@Security		basicAuth
//	@Failure		400	{object}	apperrors.AppError	"invalid param"
//	@Failure		401	{object}	apperrors.AppError	"unauthorized"
//	@Failure		403	{object}	apperrors.AppError	"forbidden"
//	@Failure		500	{object}	apperrors.AppError	"internal system error"
func GetNotebook(ctx context.Context, c *app.RequestContext, handler query.GetHandler) {
	var req getNotebookRequest
	if err := c.Bind(&req); err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	dto, err := handler.Handle(ctx, req.toDTO())
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}
	if dto == nil {
		utils.WriteHertzErrorResponse(c, apperrors.NewNotFoundError("notebook", req.Name))
		return
	}
	if len(dto.Content) == 0 {
		utils.WriteHertzErrorResponse(c, fmt.Errorf("ipynb content is empty"))
		return
	}
	var data interface{}
	if err = json.Unmarshal(dto.Content, &data); err != nil {
		applog.Errorf("notebook %s/%s is not a valid json: %s", req.WorkspaceID, req.Name, err)
		utils.WriteHertzErrorResponse(c, fmt.Errorf("ipynb content is not a valid json"))
		return
	}
	utils.WriteHertzOKResponse(c, data)
}

// ListNotebooks list notebook
//
//	@Summary		use to list notebook of workspace
//	@Description	list notebook of workspace
//	@Tags			notebook
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/workspace/{workspace-id}/notebook [get]
//	@Security		basicAuth
//	@Success		200	{object}	listNotebooksResponse
//	@Failure		400	{object}	apperrors.AppError	"invalid param"
//	@Failure		401	{object}	apperrors.AppError	"unauthorized"
//	@Failure		403	{object}	apperrors.AppError	"forbidden"
//	@Failure		500	{object}	apperrors.AppError	"internal system error"
func ListNotebooks(ctx context.Context, c *app.RequestContext, handler query.ListHandler) {
	var req listNotebookRequest
	if err := c.Bind(&req); err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	list, err := handler.Handle(ctx, req.toDTO())
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
	} else {
		utils.WriteHertzCreatedResponse(c, newListNotebooksResponse(list))
	}
}

// DeleteNotebook delete notebook
//
//	@Summary		use to delete notebook
//	@Description	delete notebook
//	@Tags			notebook
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/workspace/{workspace-id}/notebook/{name} [delete]
//	@Security		basicAuth
//	@Success		200
//	@Failure		400	{object}	apperrors.AppError	"invalid param"
//	@Failure		401	{object}	apperrors.AppError	"unauthorized"
//	@Failure		403	{object}	apperrors.AppError	"forbidden"
//	@Failure		500	{object}	apperrors.AppError	"internal system error"
func DeleteNotebook(ctx context.Context, c *app.RequestContext, handler command.DeleteHandler) {
	var req deleteNotebookRequest
	if err := c.Bind(&req); err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	if err := handler.Handle(ctx, req.toDTO()); err != nil {
		utils.WriteHertzErrorResponse(c, err)
	} else {
		utils.WriteHertzAcceptedResponse(c)
	}
}
