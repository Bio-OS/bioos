package hertz

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/Bio-OS/bioos/internal/context/notebookserver/application/command"
	"github.com/Bio-OS/bioos/internal/context/notebookserver/application/query"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	"github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/utils"
)

// CreateNotebookServer create notebook server
//
//	@Summary		use to create notebook server
//	@Description	create notebook server
//	@Tags			notebook server
//	@Produce		application/json
//	@Router			/workspace/{workspace-id}/notebookserver [post]
//	@Security		basicAuth
//	@Param			workspace-id	path		string			true	"workspace id "
//	@Param			request			body		createRequest	true	"notebook server settings"
//	@Success		201				{object}	createResponse
//	@Failure		400				{object}	apperrors.AppError	"invalid param"
//	@Failure		401				{object}	apperrors.AppError	"unauthorized"
//	@Failure		403				{object}	apperrors.AppError	"forbidden"
//	@Failure		500				{object}	apperrors.AppError	"internal system error"
func CreateNotebookServer(ctx context.Context, c *app.RequestContext, handler command.CreateHandler) {
	var req createRequest
	if err := c.Bind(&req); err != nil {
		log.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	id, err := handler.Handle(ctx, req.toDTO())
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}
	utils.WriteHertzCreatedResponse(c, createResponse{
		ID: id,
	})
}

// UpdateNotebookServerSettings update notebook server settings
//
//	@Summary		use to update notebook server settings
//	@Description	update notebook server settings
//	@Tags			notebook server
//	@Produce		application/json
//	@Router			/workspace/{workspace-id}/notebookserver/{id} [put]
//	@Security		basicAuth
//	@Param			workspace-id	path	string					true	"workspace id "
//	@Param			id				path	string					true	"notebook server id"
//	@Param			request			body	updateSettingsRequest	true	"notebook server settings"
//	@Success		202
//	@Failure		400	{object}	apperrors.AppError	"invalid param"
//	@Failure		401	{object}	apperrors.AppError	"unauthorized"
//	@Failure		403	{object}	apperrors.AppError	"forbidden"
//	@Failure		500	{object}	apperrors.AppError	"internal system error"
func UpdateNotebookServerSettings(ctx context.Context, c *app.RequestContext, handler command.UpdateHandler) {
	var req updateSettingsRequest
	if err := c.Bind(&req); err != nil {
		log.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	if err := handler.Handle(ctx, req.toDTO()); err != nil {
		utils.WriteHertzErrorResponse(c, err)
	} else {
		utils.WriteHertzAcceptedResponse(c)
	}
}

// SwitchNotebookServer switch notebook server
//
//	@Summary		use to turn notebook server on or off
//	@Description	turn notebook server on or off
//	@Tags			notebook server
//	@Produce		application/json
//	@Router			/workspace/{workspace-id}/notebookserver/{id} [post]
//	@Security		basicAuth
//	@Param			workspace-id	path	string	true	"workspace id "
//	@Param			id				path	string	true	"notebook server id"
//	@Param			on				query	boolean	false	"turn on notebook server"
//	@Param			off				query	boolean	false	"turn off notebook server"
//	@Success		202
//	@Failure		400	{object}	apperrors.AppError	"invalid param"
//	@Failure		401	{object}	apperrors.AppError	"unauthorized"
//	@Failure		403	{object}	apperrors.AppError	"forbidden"
//	@Failure		500	{object}	apperrors.AppError	"internal system error"
func SwitchNotebookServer(ctx context.Context, c *app.RequestContext, handler command.SwitchHandler) {
	var req switchRequest
	if err := c.Bind(&req); err != nil {
		log.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}
	on := c.Request.URI().QueryArgs().Has("on")
	off := c.Request.URI().QueryArgs().Has("off")
	if on && off {
		utils.WriteHertzErrorResponse(c, apperrors.NewInvalidError("on", "off", "can not set both"))
		return
	}
	if !on && !off {
		utils.WriteHertzErrorResponse(c, apperrors.NewInvalidError("on", "off", "must set one"))
		return
	}

	if err := handler.Handle(ctx, req.toDTO(on)); err != nil {
		utils.WriteHertzErrorResponse(c, err)
	} else {
		utils.WriteHertzAcceptedResponse(c)
	}
}

// DeleteNotebookServer delete notebook server
//
//	@Summary		use to delete notebook server
//	@Description	delete notebook server
//	@Tags			notebook server
//	@Produce		application/json
//	@Router			/workspace/{workspace-id}/notebookserver/{id} [delete]
//	@Security		basicAuth
//	@Param			workspace-id	path	string	true	"workspace id "
//	@Param			id				path	string	true	"notebook server id"
//	@Success		202
//	@Failure		400	{object}	apperrors.AppError	"invalid param"
//	@Failure		401	{object}	apperrors.AppError	"unauthorized"
//	@Failure		403	{object}	apperrors.AppError	"forbidden"
//	@Failure		500	{object}	apperrors.AppError	"internal system error"
func DeleteNotebookServer(ctx context.Context, c *app.RequestContext, handler command.DeleteHandler) {
	var req deleteRequest
	if err := c.Bind(&req); err != nil {
		log.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	if err := handler.Handle(ctx, req.toDTO()); err != nil {
		utils.WriteHertzErrorResponse(c, err)
	} else {
		utils.WriteHertzAcceptedResponse(c)
	}
}

// ListNotebookServers list notebook server of workspace
//
//	@Summary		use to list notebook server
//	@Description	list notebook server
//	@Tags			notebook server
//	@Produce		application/json
//	@Router			/workspace/{workspace-id}/notebookserver [get]
//	@Security		basicAuth
//	@Param			workspace-id	path		string	true	"workspace id "
//	@Success		200				{object}	[]listResponseItem
//	@Failure		400				{object}	apperrors.AppError	"invalid param"
//	@Failure		401				{object}	apperrors.AppError	"unauthorized"
//	@Failure		403				{object}	apperrors.AppError	"forbidden"
//	@Failure		500				{object}	apperrors.AppError	"internal system error"
func ListNotebookServers(ctx context.Context, c *app.RequestContext, handler query.ListHandler) {
	var req listRequest
	if err := c.Bind(&req); err != nil {
		log.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	list, err := handler.Handle(ctx, req.toDTO())
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}
	res := make([]*listResponseItem, len(list))
	for i := range list {
		res[i] = newListResponseItem(&list[i])
	}
	utils.WriteHertzCreatedResponse(c, res)
}

// GetNotebookServer get notebook server of workspace
//
//	@Summary		use to get notebook server
//	@Description	get notebook server
//	@Tags			notebook server
//	@Produce		application/json
//	@Router			/workspace/{workspace-id}/notebookserver/{id} [get]
//	@Security		basicAuth
//	@Param			workspace-id	path		string	true	"workspace id "
//	@Param			id				path		string	true	"notebook server id"
//	@Param			notebook		query		string	false	"notebook object to edit"
//	@Success		200				{object}	getResponse
//	@Failure		400				{object}	apperrors.AppError	"invalid param"
//	@Failure		401				{object}	apperrors.AppError	"unauthorized"
//	@Failure		403				{object}	apperrors.AppError	"forbidden"
//	@Failure		500				{object}	apperrors.AppError	"internal system error"
func GetNotebookServer(ctx context.Context, c *app.RequestContext, handler query.GetHandler) {
	var req getRequest
	if err := c.Bind(&req); err != nil {
		log.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	get, err := handler.Handle(ctx, req.toDTO())
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}
	utils.WriteHertzCreatedResponse(c, newGetResponse(get))
}
