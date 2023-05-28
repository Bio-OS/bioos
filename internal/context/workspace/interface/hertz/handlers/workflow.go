package handlers

import (
	"context"
	"encoding/json"

	"github.com/cloudwego/hertz/pkg/app"

	command "github.com/Bio-OS/bioos/internal/context/workspace/application/command/workflow"
	query "github.com/Bio-OS/bioos/internal/context/workspace/application/query/workflow"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	applog "github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/utils"
)

// CreateWorkflow create or update workflow
//
//	@Summary		use to create or update workflow
//	@Description	create workflow, add workflow version if id is given
//	@Tags			workflow
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/workspace/{workspace-id}/workflow [post]
//	@Security		basicAuth
//	@Param			workspace-id	path		string					true	"workspace id"
//	@Param			request			body		createWorkflowRequest	true	"create workflow request"
//	@Success		200				{object}	createWorkflowResponse
//	@Failure		400				{object}	apperrors.AppError	"invalid param"
//	@Failure		401				{object}	apperrors.AppError	"unauthorized"
//	@Failure		403				{object}	apperrors.AppError	"forbidden"
//	@Failure		500				{object}	apperrors.AppError	"internal system error"
func CreateWorkflow(ctx context.Context, c *app.RequestContext, handler command.CreateHandler) {
	var req createWorkflowRequest
	if err := c.Bind(&req); err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	reqStr, _ := json.Marshal(req)
	applog.Infow("CreateWorkflow", "reqStr", reqStr, "req", req)
	id, err := handler.Handle(ctx, req.toDTO())
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}

	resp := &createWorkflowResponse{ID: id}
	utils.WriteHertzCreatedResponse(c, resp)

}

// GetWorkflow get workflow by id
//
//	@Summary		get workflow
//	@Description	get workflow
//	@Tags			workflow
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/workspace/{workspace-id}/workflow/{id} [get]
//	@Security		basicAuth
//	@Param			id				path		string	true	"workflow id"
//	@Param			workspace-id	path		string	true	"workspace id"
//	@Success		200				{object}	getWorkflowResponse
//	@Failure		400				{object}	apperrors.AppError	"invalid param"
//	@Failure		401				{object}	apperrors.AppError	"unauthorized"
//	@Failure		403				{object}	apperrors.AppError	"forbidden"
//	@Failure		500				{object}	apperrors.AppError	"internal system error"
func GetWorkflow(ctx context.Context, c *app.RequestContext, handler query.GetHandler) {
	var req getWorkflowRequest
	if err := c.Bind(&req); err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	workflow, err := handler.Handle(ctx, req.toDTO())
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}
	resp := &getWorkflowResponse{Workflow: WorkflowDTOtoWorkflowItemVO(workflow)}
	utils.WriteHertzOKResponse(c, resp)
}

// ListWorkflows list workflows
//
//	@Summary		use to list workflows of workspace
//	@Description	list workflow of workspace
//	@Tags			workflow
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/workspace/{workspace-id}/workflow [get]
//	@Security		basicAuth
//	@Param			workspace-id	path		string	true	"workspace id"
//	@Param			page			query		int		false	"page number"
//	@Param			size			query		int		false	"page size"
//	@Param			orderBy			query		string	false	"support order field: name/createdAt, support order: asc/desc, seperated by comma, eg: createdAt:desc,name:asc"
//	@Param			searchWord		query		string	false	"workflow name"
//	@Param			exact			query		bool	false	"exact"
//	@Param			ids				query		string	false	"workspace ids seperated by comma"
//	@Success		200				{object}	listWorkflowsResponse
//	@Failure		400				{object}	apperrors.AppError	"invalid param"
//	@Failure		401				{object}	apperrors.AppError	"unauthorized"
//	@Failure		403				{object}	apperrors.AppError	"forbidden"
//	@Failure		500				{object}	apperrors.AppError	"internal system error"
func ListWorkflows(ctx context.Context, c *app.RequestContext, handler query.ListHandler) {
	var req listWorkflowsRequest
	if err := c.Bind(&req); err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	reqDTO, err := req.toDTO()
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}
	workflows, total, err := handler.Handle(ctx, reqDTO)
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}
	items := make([]*WorkflowItem, 0, len(workflows))
	for _, wf := range workflows {
		items = append(items, WorkflowDTOtoWorkflowItemVO(wf))
	}
	resp := &listWorkflowsResponse{
		Page:  req.Page,
		Size:  req.Size,
		Total: total,
		Items: items,
	}
	utils.WriteHertzOKResponse(c, resp)
}

// UpdateWorkflow update workflow
//
//	@Summary		use to update workflow
//	@Description	update workspace
//	@Tags			workflow
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/workspace/{workspace-id}/workflow/{id} [patch]
//	@Security		basicAuth
//	@Param			workspace-id	path	string					true	"update workspace id"
//	@Param			id				path	string					true	"update workflow id"
//	@Param			request			body	updateWorkflowRequest	true	"update workflow request"
//	@Success		200
//	@Failure		400	{object}	apperrors.AppError	"invalid param"
//	@Failure		401	{object}	apperrors.AppError	"unauthorized"
//	@Failure		403	{object}	apperrors.AppError	"forbidden"
//	@Failure		404	{object}	apperrors.AppError	"not found"
//	@Failure		500	{object}	apperrors.AppError	"internal system error"
func UpdateWorkflow(ctx context.Context, c *app.RequestContext, handler command.UpdateHandler) {
	var req updateWorkflowRequest
	err := c.Bind(&req)
	if err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	err = handler.Handle(ctx, req.toDTO())
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}

	utils.WriteHertzOKResponse(c, updateWorkflowResponse{})
}

// DeleteWorkflow delete workflow
//
//	@Summary		use to delete workflow
//	@Description	delete workflow
//	@Tags			workflow
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/workspace/{workspace-id}/workflow/{id} [delete]
//	@Security		basicAuth
//	@Param			workspace-id	path	string	true	"delete workspace id"
//	@Param			id				path	string	true	"delete workflow id"
//	@Success		202
//	@Failure		400	{object}	apperrors.AppError	"invalid param"
//	@Failure		401	{object}	apperrors.AppError	"unauthorized"
//	@Failure		403	{object}	apperrors.AppError	"forbidden"
//	@Failure		404	{object}	apperrors.AppError	"not found"
//	@Failure		500	{object}	apperrors.AppError	"internal system error"
func DeleteWorkflow(ctx context.Context, c *app.RequestContext, handler command.DeleteHandler) {
	var req deleteWorkflowRequest
	err := c.Bind(&req)
	if err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	err = handler.Handle(ctx, req.toDTO())
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}

	utils.WriteHertzOKResponse(c, deleteWorkflowResponse{})
}

// ListWorkflowFiles list workflow files
//
//	@Summary		use to list workflow files
//	@Description	list workflow of workspace
//	@Tags			workflow
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/workspace/{workspace-id}/workflow/{workflow-id}/file [get]
//	@Security		basicAuth
//	@Param			workspace-id		path		string						true	"workspace id"
//	@Param			workflow-id			path		string						true	"workflow id"
//	@Param			request				body		listWorkflowFilesRequest	true	"list workflow files"
//	@Param			page				query		int							false	"page number"
//	@Param			size				query		int							false	"page size"
//	@Param			orderBy				query		string						false	"support order field: version/path, support order: asc/desc, seperated by comma, eg: version:desc,path:asc"
//	@Param			searchWord			query		string						false	"workflow name"
//	@Param			ids					query		string						false	"workspace file ids seperated by comma"
//	@Param			workflowVersionID	query		string						false	"workspace version id"
//	@Success		200					{object}	listWorkflowFilesResponse
//	@Failure		400					{object}	apperrors.AppError	"invalid param"
//	@Failure		401					{object}	apperrors.AppError	"unauthorized"
//	@Failure		403					{object}	apperrors.AppError	"forbidden"
//	@Failure		500					{object}	apperrors.AppError	"internal system error"
func ListWorkflowFiles(ctx context.Context, c *app.RequestContext, handler query.ListFilesHandler) {
	var req listWorkflowFilesRequest
	if err := c.Bind(&req); err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	reqDTO, err := req.toDTO()
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}
	workflowFiles, total, err := handler.Handle(ctx, reqDTO)
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}
	items := make([]*WorkflowFile, 0, len(workflowFiles))
	for _, wf := range workflowFiles {
		items = append(items, WorkflowFileDTOtoWorkflowFileVO(wf))
	}
	resp := &listWorkflowFilesResponse{
		Page:        req.Page,
		Size:        req.Size,
		Total:       total,
		WorkspaceID: req.WorkspaceID,
		WorkflowID:  req.WorkflowID,
		Items:       items,
	}
	utils.WriteHertzOKResponse(c, resp)
}

// GetWorkflowFile get workflow file by id
//
//	@Summary		get workflow file
//	@Description	get workflow file
//	@Tags			workflow
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/workspace/{workspace-id}/workflow/{workflow-id}/file/{id} [get]
//	@Security		basicAuth
//	@Param			id				path		string	true	"workflow file id"
//	@Param			workspace-id	path		string	true	"workspace id"
//	@Param			workflow-id		path		string	true	"workflow id"
//	@Success		200				{object}	getWorkflowFileResponse
//	@Failure		400				{object}	apperrors.AppError	"invalid param"
//	@Failure		401				{object}	apperrors.AppError	"unauthorized"
//	@Failure		403				{object}	apperrors.AppError	"forbidden"
//	@Failure		500				{object}	apperrors.AppError	"internal system error"
func GetWorkflowFile(ctx context.Context, c *app.RequestContext, handler query.GetFileHandler) {
	var req getWorkflowFileRequest
	if err := c.Bind(&req); err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	workflowFile, err := handler.Handle(ctx, req.toDTO())
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}
	resp := &getWorkflowFileResponse{File: WorkflowFileDTOtoWorkflowFileVO(workflowFile)}
	utils.WriteHertzOKResponse(c, resp)
}

// ListWorkflowVersions list workflow version
//
//	@Summary		use to list workflow versions
//	@Description	list workflow verions of workspace
//	@Tags			workflow
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/workspace/{workspace-id}/workflow/{workflow-id}/versions [get]
//	@Security		basicAuth
//	@Param			workspace-id	path		string	true	"workspace id"
//	@Param			workflow-id		path		string	true	"workflow id"
//	@Param			page			query		int		false	"page number"
//	@Param			size			query		int		false	"page size"
//	@Param			orderBy			query		string	false	"support order field: source/language/status, support order: asc/desc, seperated by comma, eg: status:desc,language:asc"
//	@Param			ids				query		string	false	"workspace version ids seperated by comma"
//	@Success		200				{object}	listWorkflowVersionsResponse
//	@Failure		400				{object}	apperrors.AppError	"invalid param"
//	@Failure		401				{object}	apperrors.AppError	"unauthorized"
//	@Failure		403				{object}	apperrors.AppError	"forbidden"
//	@Failure		500				{object}	apperrors.AppError	"internal system error"
func ListWorkflowVersions(ctx context.Context, c *app.RequestContext, handler query.ListVersionsHandler) {
	var req listWorkflowVersionsRequest
	if err := c.Bind(&req); err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	reqDTO, err := req.toDTO()
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}
	workflowVersions, total, err := handler.Handle(ctx, reqDTO)
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}
	items := make([]*WorkflowVersion, 0, len(workflowVersions))
	for _, wf := range workflowVersions {
		items = append(items, WorkflowVersionDTOtoVO(wf))
	}
	resp := &listWorkflowVersionsResponse{
		Page:        req.Page,
		Size:        req.Size,
		Total:       total,
		WorkspaceID: req.WorkspaceID,
		WorkflowID:  req.WorkflowID,
		Items:       items,
	}
	utils.WriteHertzOKResponse(c, resp)
}

// GetWorkflowVersion get workflow version by id
//
//	@Summary		get workflow version
//	@Description	get workflow version
//	@Tags			workflow
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/workspace/{workspace-id}/workflow/{workflow-id}/version/{id} [get]
//	@Security		basicAuth
//	@Param			id				path		string	true	"workflow version id"
//	@Param			workspace-id	path		string	true	"workspace id"
//	@Param			workflow-id		path		string	true	"workflow id"
//	@Success		200				{object}	getWorkflowVersionResponse
//	@Failure		400				{object}	apperrors.AppError	"invalid param"
//	@Failure		401				{object}	apperrors.AppError	"unauthorized"
//	@Failure		403				{object}	apperrors.AppError	"forbidden"
//	@Failure		500				{object}	apperrors.AppError	"internal system error"
func GetWorkflowVersion(ctx context.Context, c *app.RequestContext, handler query.GetVersionHandler) {
	var req getWorkflowVersionRequest
	if err := c.Bind(&req); err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	workflowVersion, err := handler.Handle(ctx, req.toDTO())
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}
	resp := &getWorkflowVersionResponse{Version: WorkflowVersionDTOtoVO(workflowVersion)}
	utils.WriteHertzOKResponse(c, resp)
}
