package hertz

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/route"

	"github.com/Bio-OS/bioos/internal/context/notebookserver/application"
	"github.com/Bio-OS/bioos/pkg/middlewares/hertz"
	"github.com/Bio-OS/bioos/pkg/server"
)

type register struct {
	svc *application.Service
}

func NewRouteRegister(service *application.Service) server.RouteRegister {
	return &register{
		svc: service,
	}
}

func (r *register) AddRoute(h route.IRouter) {
	workspace := h.Group("/workspace")
	workspace.Use(hertz.Authn())

	workspace.POST("/:workspace-id/notebookserver", hertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		workspaceID := c.Param("workspace-id")
		return fmt.Sprintf("Workspace-%s:CreateNotebookServer", workspaceID)
	}), func(c context.Context, ctx *app.RequestContext) {
		CreateNotebookServer(c, ctx, r.svc.Commands.Create)
	})

	workspace.PUT("/:workspace-id/notebookserver/:id", hertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		workspaceID := c.Param("workspace-id")
		return fmt.Sprintf("Workspace-%s:UpdateNotebookServerSettings", workspaceID)
	}), func(c context.Context, ctx *app.RequestContext) {
		UpdateNotebookServerSettings(c, ctx, r.svc.Commands.Update)
	})

	workspace.POST("/:workspace-id/notebookserver/:id", hertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		workspaceID := c.Param("workspace-id")
		return fmt.Sprintf("Workspace-%s:SwitchNotebookServer", workspaceID)
	}), func(c context.Context, ctx *app.RequestContext) {
		SwitchNotebookServer(c, ctx, r.svc.Commands.Switch)
	})

	workspace.DELETE("/:workspace-id/notebookserver/:id", hertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		workspaceID := c.Param("workspace-id")
		return fmt.Sprintf("Workspace-%s:DeleteNotebookServer", workspaceID)
	}), func(c context.Context, ctx *app.RequestContext) {
		DeleteNotebookServer(c, ctx, r.svc.Commands.Delete)
	})

	workspace.GET("/:workspace-id/notebookserver", hertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		workspaceID := c.Param("workspace-id")
		return fmt.Sprintf("Workspace-%s:ListNotebookServers", workspaceID)
	}), func(c context.Context, ctx *app.RequestContext) {
		ListNotebookServers(c, ctx, r.svc.Queries.List)
	})

	workspace.GET("/:workspace-id/notebookserver/:id", hertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		workspaceID := c.Param("workspace-id")
		return fmt.Sprintf("Workspace-%s:GetNotebookServer", workspaceID)
	}), func(c context.Context, ctx *app.RequestContext) {
		GetNotebookServer(c, ctx, r.svc.Queries.Get)
	})
}
