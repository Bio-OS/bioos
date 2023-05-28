package handlers

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	command "github.com/Bio-OS/bioos/internal/context/submission/application/command/run"
	query "github.com/Bio-OS/bioos/internal/context/submission/application/query/run"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	applog "github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/utils"
)

// CancelRun cancel run
//
//	@Summary		use to cancel run
//	@Description	cancel run
//	@Tags			submission
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/workspace/{workspace_id}/submission/{submission_id}/run/{id}/cancel [post]
//	@Security		basicAuth
//	@Param			workspace_id	path	string	true	"workspace id"
//	@Param			submission_id	path	string	true	"submission id"
//	@Param			id				path	string	true	"run id"
//	@Success		202
//	@Failure		400	{object}	apperrors.AppError	"invalid param"
//	@Failure		401	{object}	apperrors.AppError	"unauthorized"
//	@Failure		403	{object}	apperrors.AppError	"forbidden"
//	@Failure		500	{object}	apperrors.AppError	"internal system error"
func CancelRun(ctx context.Context, c *app.RequestContext, handler command.CancelRunHandler) {
	var req CancelRunRequest
	err := c.Bind(&req)
	if err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	cmd := cancelRunVoToDto(req)
	err = handler.Handle(ctx, cmd)
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}

	utils.WriteHertzAcceptedResponse(c)
}

// ListRuns list runs
//
//	@Summary		use to list runs
//	@Description	list runs
//	@Tags			submission
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/workspace/{workspace_id}/submission/{submission_id}/run [get]
//	@Security		basicAuth
//	@Param			workspace_id	path		string		true	"workspace id"
//	@Param			submission_id	path		string		true	"submission id"
//	@Param			page			query		int			false	"query page"
//	@Param			size			query		int			false	"query size"
//	@Param			orderBy			query		string		false	"query order, just like field1,field2:desc"
//	@Param			searchWord		query		string		false	"query searchWord"
//	@Param			ids				query		[]string	false	"query ids"
//	@Param			status			query		[]string	false	"query status"
//	@Success		200				{object}	ListRunsResponse
//	@Failure		400				{object}	apperrors.AppError	"invalid param"
//	@Failure		401				{object}	apperrors.AppError	"unauthorized"
//	@Failure		403				{object}	apperrors.AppError	"forbidden"
//	@Failure		500				{object}	apperrors.AppError	"internal system error"
func ListRuns(ctx context.Context, c *app.RequestContext, handler query.ListRunsHandler) {
	var req ListRunsRequest
	err := c.Bind(&req)
	if err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	query, err := listRunsVoToDto(req)
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}
	runs, total, err := handler.Handle(ctx, query)
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}
	items := make([]RunItem, 0, len(runs))
	for _, run := range runs {
		items = append(items, runItemDtoToVo(run))
	}
	resp := &ListRunsResponse{
		query.Pg.Page,
		query.Pg.Size,
		total,
		items,
	}
	utils.WriteHertzOKResponse(c, resp)
}

// ListTasks list tasks
//
//	@Summary		use to list tasks
//	@Description	list tasks
//	@Tags			submission
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/workspace/{workspace_id}/submission/{submission_id}/run/{run_id}/task [get]
//	@Security		basicAuth
//	@Param			workspace_id	path		string	true	"workspace id"
//	@Param			submission_id	path		string	true	"submission id"
//	@Param			run_id			path		string	true	"run id"
//	@Param			page			query		int		false	"query page"
//	@Param			size			query		int		false	"query size"
//	@Param			orderBy			query		string	false	"query order, just like field1,field2:desc"
//	@Success		200				{object}	ListTasksResponse
//	@Failure		400				{object}	apperrors.AppError	"invalid param"
//	@Failure		401				{object}	apperrors.AppError	"unauthorized"
//	@Failure		403				{object}	apperrors.AppError	"forbidden"
//	@Failure		500				{object}	apperrors.AppError	"internal system error"
func ListTasks(ctx context.Context, c *app.RequestContext, handler query.ListTasksHandler) {
	var req ListTasksRequest
	err := c.Bind(&req)
	if err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	query, err := listTasksVoToDto(req)
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}
	tasks, total, err := handler.Handle(ctx, query)
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}
	items := make([]TaskItem, 0, len(tasks))
	for _, task := range tasks {
		items = append(items, taskItemDtoToVo(task))
	}
	resp := &ListTasksResponse{
		query.Pg.Page,
		query.Pg.Size,
		total,
		items,
	}
	utils.WriteHertzOKResponse(c, resp)
}
