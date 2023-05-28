package handlers

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/cloudwego/hertz/pkg/app"

	command "github.com/Bio-OS/bioos/internal/context/workspace/application/command/workspace"
	query "github.com/Bio-OS/bioos/internal/context/workspace/application/query/workspace"
	"github.com/Bio-OS/bioos/pkg/consts"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	applog "github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/utils"
)

// CreateWorkspace create workspace
//
//	@Summary		use to create workspace
//	@Description	create workspace
//	@Tags			workspace
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/workspace [post]
//	@Security		basicAuth
//	@Param			request	body		CreateWorkspaceRequest	true	"create workspace request"
//	@Success		201		{object}	CreateWorkspaceResponse
//	@Failure		400		{object}	apperrors.AppError	"invalid param"
//	@Failure		401		{object}	apperrors.AppError	"unauthorized"
//	@Failure		403		{object}	apperrors.AppError	"forbidden"
//	@Failure		500		{object}	apperrors.AppError	"internal system error"
func CreateWorkspace(ctx context.Context, c *app.RequestContext, handler command.CreateWorkspaceHandler) {
	var req CreateWorkspaceRequest
	err := c.Bind(&req)
	if err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	cmd := createWorkspaceVoToDto(req)
	id, err := handler.Handle(ctx, cmd)
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}

	resp := &CreateWorkspaceResponse{Id: id}
	utils.WriteHertzCreatedResponse(c, resp)
}

// DeleteWorkspace delete workspace
//
//	@Summary		use to delete workspace
//	@Description	delete workspace
//	@Tags			workspace
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/workspace/{id} [delete]
//	@Security		basicAuth
//	@Param			id	path	string	true	"delete workspace id"
//	@Success		202
//	@Failure		400	{object}	apperrors.AppError	"invalid param"
//	@Failure		401	{object}	apperrors.AppError	"unauthorized"
//	@Failure		403	{object}	apperrors.AppError	"forbidden"
//	@Failure		404	{object}	apperrors.AppError	"not found"
//	@Failure		500	{object}	apperrors.AppError	"internal system error"
func DeleteWorkspace(ctx context.Context, c *app.RequestContext, handler command.DeleteWorkspaceHandler) {
	var req DeleteWorkspaceRequest
	err := c.Bind(&req)
	if err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	cmd := &command.DeleteWorkspaceCommand{ID: req.Id}
	err = handler.Handle(ctx, cmd)
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}

	utils.WriteHertzAcceptedResponse(c)
}

// UpdateWorkspace update workspace
//
//	@Summary		use to update workspace
//	@Description	update workspace
//	@Tags			workspace
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/workspace/{id} [patch]
//	@Security		basicAuth
//	@Param			id		path	string					true	"update workspace id"
//	@Param			request	body	UpdateWorkspaceRequest	true	"update workspace request"
//	@Success		200
//	@Failure		400	{object}	apperrors.AppError	"invalid param"
//	@Failure		401	{object}	apperrors.AppError	"unauthorized"
//	@Failure		403	{object}	apperrors.AppError	"forbidden"
//	@Failure		404	{object}	apperrors.AppError	"not found"
//	@Failure		500	{object}	apperrors.AppError	"internal system error"
func UpdateWorkspace(ctx context.Context, c *app.RequestContext, handler command.UpdateWorkspaceHandler) {
	var req UpdateWorkspaceRequest
	err := c.Bind(&req)
	if err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	cmd := updateWorkspaceVoToDto(req)
	err = handler.Handle(ctx, cmd)
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}

	utils.WriteHertzOKResponse(c, nil)
}

// GetWorkspaceById get workspace
//
//	@Summary		use to get workspace by id
//	@Description	get workspace
//	@Tags			workspace
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/workspace/{id} [get]
//	@Security		basicAuth
//	@Param			id	path		string	true	"get workspace id"
//	@Success		200	{object}	GetWorkspaceByIdResponse
//	@Failure		400	{object}	apperrors.AppError	"invalid param"
//	@Failure		401	{object}	apperrors.AppError	"unauthorized"
//	@Failure		403	{object}	apperrors.AppError	"forbidden"
//	@Failure		404	{object}	apperrors.AppError	"not found"
//	@Failure		500	{object}	apperrors.AppError	"internal system error"
func GetWorkspaceById(ctx context.Context, c *app.RequestContext, handler query.GetWorkspaceByIDQueryHandler) {
	var req GetWorkspaceByIdRequest
	err := c.Bind(&req)
	if err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	qry := &query.GetWorkspaceByIDQuery{ID: req.Id}
	ws, err := handler.Handle(ctx, qry)
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}

	resp := &GetWorkspaceByIdResponse{WorkspaceItem: workspaceItemDtoToVo(ws)}
	utils.WriteHertzOKResponse(c, resp)
}

// ListWorkspaces list workspaces
//
//	@Summary		use to list workspaces
//	@Description	list workspaces
//	@Tags			workspace
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/workspace [get]
//	@Security		basicAuth
//	@Param			page		query		int		false	"query page"
//	@Param			size		query		int		false	"query size"
//	@Param			orderBy		query		string	false	"query order, just like field1,field2:desc"
//	@Param			searchWord	query		string	false	"query searchWord"
//	@Param			exact		query		bool	false	"query exact"
//	@Param			ids			query		string	false	"query ids, split by comma"
//	@Success		200			{object}	ListWorkspacesResponse
//	@Failure		400			{object}	apperrors.AppError	"invalid param"
//	@Failure		401			{object}	apperrors.AppError	"unauthorized"
//	@Failure		403			{object}	apperrors.AppError	"forbidden"
//	@Failure		500			{object}	apperrors.AppError	"internal system error"
func ListWorkspaces(ctx context.Context, c *app.RequestContext, handler query.ListWorkspacesHandler) {
	var req ListWorkspacesRequest
	err := c.Bind(&req)
	if err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	qry, err := listWorkspacesVoToDto(req)
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}
	workspaces, total, err := handler.Handle(ctx, qry)
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}

	items := make([]WorkspaceItem, 0, len(workspaces))
	for _, ws := range workspaces {
		items = append(items, workspaceItemDtoToVo(ws))
	}
	resp := &ListWorkspacesResponse{
		Page:  req.Page,
		Size:  req.Size,
		Total: total,
		Items: items,
	}
	utils.WriteHertzOKResponse(c, resp)
}

// ImportWorkspace import workspace
//
//	@Summary		use to import workspace
//	@Description	import workspace
//	@Tags			workspace
//	@Accept			multipart/form-data
//	@Produce		application/json
//	@Router			/workspace [put]
//	@Security		basicAuth
//	@Param			file		formData	file	true	"file"
//	@Param			mountType	query		string	true	"workspace mount path"
//	@Param			mountPath	query		string	true	"workspace mount type, only support nfs"
//	@Success		200			{object}	ImportWorkspaceResponse
//	@Failure		400			{object}	apperrors.AppError	"invalid param"
//	@Failure		401			{object}	apperrors.AppError	"unauthorized"
//	@Failure		403			{object}	apperrors.AppError	"forbidden"
//	@Failure		500			{object}	apperrors.AppError	"internal system error"
func ImportWorkspace(ctx context.Context, c *app.RequestContext, handler command.ImportWorkspaceHandler) {
	var req ImportWorkspaceRequest
	err := c.Bind(&req)
	if err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	//In order to avoid the cost of wrong file writing, we validate mountType here
	if req.MountType != "nfs" {
		err = fmt.Errorf("unsupported mount type: %w", err)
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	file, err := c.Request.FormFile("file")
	if err != nil {
		applog.Errorw("hertz get form file error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzFormFileError(err))
		return
	}
	req.FileName = file.Filename

	cmd := importWorkspaceVoToDto(req)

	//In order to avoid the cost of wrong file writing, we validate file type here
	if path.Ext(cmd.FileName) != consts.ImportWorkspaceFileTypeExt {
		err = fmt.Errorf("unsupported file type: %s", path.Ext(cmd.FileName))
		applog.Errorw("hertz get form file error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzFormFileError(err))
		return
	}

	dst := path.Join(cmd.Storage.NFS.MountPath, cmd.ID, cmd.FileName)
	if err = os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		applog.Errorw("create dir error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzFormFileError(err))
		return
	}
	// Upload the file to specific dst
	err = c.SaveUploadedFile(file, path.Join(cmd.Storage.NFS.MountPath, cmd.ID, cmd.FileName))
	if err != nil {
		applog.Errorw("hertz get form file error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzFormFileError(err))
		return
	}

	err = handler.Handle(ctx, cmd)
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}
	resp := &ImportWorkspaceResponse{
		Id: cmd.ID,
	}

	utils.WriteHertzOKResponse(c, resp)
}
