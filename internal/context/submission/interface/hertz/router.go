package hertz

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/route"

	"github.com/Bio-OS/bioos/pkg/server"

	"github.com/Bio-OS/bioos/internal/context/submission/application"
	"github.com/Bio-OS/bioos/internal/context/submission/interface/hertz/handlers"
	apphertz "github.com/Bio-OS/bioos/pkg/middlewares/hertz"
)

type register struct {
	svc *application.SubmissionService
}

func NewRouteRegister(submissionService *application.SubmissionService) server.RouteRegister {
	return &register{
		svc: submissionService,
	}
}

func (r *register) AddRoute(h route.IRouter) {
	submission := h.Group("/workspace/:workspace_id/submission")
	submission.Use(apphertz.Authn())
	submission.POST("", apphertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		return fmt.Sprintf("Workspace-%s:CreateSubmission", c.Param("workspace_id"))
	}), func(c context.Context, ctx *app.RequestContext) {
		handlers.CreateSubmission(c, ctx, r.svc.SubmissionCommands.CreateSubmission)
	})

	submission.GET("/:name", apphertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		return fmt.Sprintf("Workspace-%s:CheckSubmission", c.Param("workspace_id"))
	}), func(c context.Context, ctx *app.RequestContext) {
		handlers.CheckSubmission(c, ctx, r.svc.SubmissionQueries.Check)
	})

	submission.GET("", apphertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		return fmt.Sprintf("Workspace-%s:ListSubmissions", c.Param("workspace_id"))
	}), func(c context.Context, ctx *app.RequestContext) {
		handlers.ListSubmissions(c, ctx, r.svc.SubmissionQueries.List, r.svc.RunQueries.CountRunsResult)
	})

	submission.POST("/:id/cancel", apphertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		return fmt.Sprintf("Workspace-%s:CancelSubmission", c.Param("workspace_id"))
	}), func(c context.Context, ctx *app.RequestContext) {
		handlers.CancelSubmission(c, ctx, r.svc.SubmissionCommands.CancelSubmission)
	})

	submission.DELETE("/:id", apphertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		return fmt.Sprintf("Workspace-%s:DeleteSubmission", c.Param("workspace_id"))
	}), func(c context.Context, ctx *app.RequestContext) {
		handlers.DeleteSubmission(c, ctx, r.svc.SubmissionCommands.DeleteSubmission)
	})

	addNotebookRoute(submission, r.svc)
	return
}

func addNotebookRoute(group *route.RouterGroup, submissionService *application.SubmissionService) {
	group.POST("/:submission_id/run/:id/cancel", apphertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		return fmt.Sprintf("Workspace-%s:CancelRun", c.Param("workspace_id"))
	}), func(c context.Context, ctx *app.RequestContext) {
		handlers.CancelRun(c, ctx, submissionService.RunCommands.CancelRun)
	})

	group.GET("/:submission_id/run", apphertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		return fmt.Sprintf("Workspace-%s:ListRuns", c.Param("workspace_id"))
	}), func(c context.Context, ctx *app.RequestContext) {
		handlers.ListRuns(c, ctx, submissionService.RunQueries.ListRuns)
	})

	group.GET("/:submission_id/run/:run_id/task", apphertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		return fmt.Sprintf("Workspace-%s:ListTasks", c.Param("workspace_id"))
	}), func(c context.Context, ctx *app.RequestContext) {
		handlers.ListTasks(c, ctx, submissionService.RunQueries.ListTasks)
	})
}
