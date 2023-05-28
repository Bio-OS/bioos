package handlers

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	runquery "github.com/Bio-OS/bioos/internal/context/submission/application/query/run"

	command "github.com/Bio-OS/bioos/internal/context/submission/application/command/submission"
	query "github.com/Bio-OS/bioos/internal/context/submission/application/query/submission"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	applog "github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/utils"
)

// CreateSubmission create submission
//
//	@Summary		use to create submission
//	@Description	create submission
//	@Tags			submission
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/workspace/{workspace_id}/submission [post]
//	@Security		basicAuth
//	@Param			workspace_id	path		string					true	"workspace id"
//	@Param			request			body		CreateSubmissionRequest	true	"create submission request"
//	@Success		201				{object}	CreateSubmissionResponse
//	@Failure		400				{object}	apperrors.AppError	"invalid param"
//	@Failure		401				{object}	apperrors.AppError	"unauthorized"
//	@Failure		403				{object}	apperrors.AppError	"forbidden"
//	@Failure		500				{object}	apperrors.AppError	"internal system error"
func CreateSubmission(ctx context.Context, c *app.RequestContext, handler command.CreateSubmissionHandler) {
	var req CreateSubmissionRequest
	err := c.Bind(&req)
	if err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	cmd := createSubmissionVoToDto(req)
	id, err := handler.Handle(ctx, cmd)
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}

	resp := &CreateSubmissionResponse{ID: id}
	utils.WriteHertzCreatedResponse(c, resp)
}

// CancelSubmission cancel submission
//
//	@Summary		use to cancel submission
//	@Description	cancel submission
//	@Tags			submission
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/workspace/{workspace_id}/submission/{id}/cancel [post]
//	@Security		basicAuth
//	@Param			workspace_id	path	string	true	"workspace id"
//	@Param			id				path	string	true	"submission id"
//	@Success		202
//	@Failure		400	{object}	apperrors.AppError	"invalid param"
//	@Failure		401	{object}	apperrors.AppError	"unauthorized"
//	@Failure		403	{object}	apperrors.AppError	"forbidden"
//	@Failure		500	{object}	apperrors.AppError	"internal system error"
func CancelSubmission(ctx context.Context, c *app.RequestContext, handler command.CancelSubmissionHandler) {
	var req CancelSubmissionRequest
	err := c.Bind(&req)
	if err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	cmd := cancelSubmissionVoToDto(req)
	err = handler.Handle(ctx, cmd)
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}

	utils.WriteHertzAcceptedResponse(c)
}

// DeleteSubmission delete submission
//
//	@Summary		use to delete submission
//	@Description	delete submission
//	@Tags			submission
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/workspace/{workspace_id}/submission/{id} [delete]
//	@Security		basicAuth
//	@Param			workspace_id	path	string	true	"workspace id"
//	@Param			id				path	string	true	"submission id"
//	@Success		202
//	@Failure		400	{object}	apperrors.AppError	"invalid param"
//	@Failure		401	{object}	apperrors.AppError	"unauthorized"
//	@Failure		403	{object}	apperrors.AppError	"forbidden"
//	@Failure		500	{object}	apperrors.AppError	"internal system error"
func DeleteSubmission(ctx context.Context, c *app.RequestContext, handler command.DeleteSubmissionHandler) {
	var req DeleteSubmissionRequest
	err := c.Bind(&req)
	if err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	cmd := deleteSubmissionVoToDto(req)
	err = handler.Handle(ctx, cmd)
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}

	utils.WriteHertzAcceptedResponse(c)
}

// CheckSubmission check submission name
//
//	@Summary		use to check submission name unique
//	@Description	check submission name unique
//	@Tags			submission
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/workspace/{workspace_id}/submission/{name} [get]
//	@Security		basicAuth
//	@Param			workspace_id	path	string	true	"workspace id"
//	@Param			name			path	string	true	"submission name"
//	@Success		202
//	@Failure		400	{object}	apperrors.AppError	"invalid param"
//	@Failure		401	{object}	apperrors.AppError	"unauthorized"
//	@Failure		403	{object}	apperrors.AppError	"forbidden"
//	@Failure		500	{object}	apperrors.AppError	"internal system error"
func CheckSubmission(ctx context.Context, c *app.RequestContext, handler query.CheckHandler) {
	var req CheckSubmissionRequest
	err := c.Bind(&req)
	if err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	cmd := checkSubmissionVoToDto(req)
	flag, err := handler.Handle(ctx, cmd)
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}

	resp := &CheckSubmissionResponse{IsNameExist: flag}
	utils.WriteHertzOKResponse(c, resp)
}

// ListSubmissions list submissions
//
//	@Summary		use to list submissions
//	@Description	list submissions
//	@Tags			submission
//	@Accept			application/json
//	@Produce		application/json
//	@Router			/workspace/{workspace_id}/submission [get]
//	@Security		basicAuth
//	@Param			workspace_id	path		string		true	"workspace id"
//	@Param			page			query		int			false	"query page"
//	@Param			size			query		int			false	"query size"
//	@Param			orderBy			query		string		false	"query order, just like field1,field2:desc"
//	@Param			searchWord		query		string		false	"query searchWord"
//	@Param			exact			query		bool		false	"query exact"
//	@Param			ids				query		[]string	false	"query ids"
//	@Param			workflowID		query		string		false	"workflow id"
//	@Param			status			query		[]string	false	"query status"
//	@Success		200				{object}	ListSubmissionsResponse
//	@Failure		400				{object}	apperrors.AppError	"invalid param"
//	@Failure		401				{object}	apperrors.AppError	"unauthorized"
//	@Failure		403				{object}	apperrors.AppError	"forbidden"
//	@Failure		500				{object}	apperrors.AppError	"internal system error"
func ListSubmissions(ctx context.Context, c *app.RequestContext, subHandler query.ListHandler, runHandler runquery.CountRunsResultHandler) {
	var req ListSubmissionsRequest
	err := c.Bind(&req)
	if err != nil {
		applog.Errorw("hertz bind error", "err", err)
		utils.WriteHertzErrorResponse(c, apperrors.NewHertzBindError(err))
		return
	}

	query, err := listSubmissionsVoToDto(req)
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}
	submissions, total, err := subHandler.Handle(ctx, query)
	if err != nil {
		utils.WriteHertzErrorResponse(c, err)
		return
	}
	for _, sub := range submissions {
		runStatus, err := runHandler.Handle(ctx, &runquery.CountRunsResultQuery{SubmissionID: sub.ID})
		if err != nil {
			utils.WriteHertzErrorResponse(c, err)
			return
		}
		sub.RunStatus = runStatusDtoToSubmissionRunStatus(runStatus)
	}
	items := make([]SubmissionItem, 0, len(submissions))
	for _, submission := range submissions {
		items = append(items, submissionItemDtoToVo(submission))
	}
	resp := &ListSubmissionsResponse{
		query.Pg.Page,
		query.Pg.Size,
		total,
		items,
	}
	utils.WriteHertzOKResponse(c, resp)
}
