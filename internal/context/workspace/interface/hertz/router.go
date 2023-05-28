package hertz

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/route"

	"github.com/Bio-OS/bioos/internal/context/workspace/application"
	"github.com/Bio-OS/bioos/internal/context/workspace/interface/hertz/handlers"
	apphertz "github.com/Bio-OS/bioos/pkg/middlewares/hertz"
	"github.com/Bio-OS/bioos/pkg/server"
)

type register struct {
	svc *application.WorkspaceService
}

func NewRouteRegister(workspaceService *application.WorkspaceService) server.RouteRegister {
	return &register{
		svc: workspaceService,
	}
}

func (r *register) AddRoute(h route.IRouter) {
	workspace := h.Group("/workspace")
	workspace.Use(apphertz.Authn())

	addWorkspaceRoute(workspace, r.svc)
	addWorkflowRoute(workspace, r.svc)
	addNotebookRoute(workspace, r.svc)
	addDataModelRouter(workspace, r.svc)
	return
}

func addWorkspaceRoute(group *route.RouterGroup, workspaceService *application.WorkspaceService) {
	group.GET("/:id", apphertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		workspaceId := c.Param("id")
		return fmt.Sprintf("Workspace-%s:Get", workspaceId)
	}), func(c context.Context, ctx *app.RequestContext) {
		handlers.GetWorkspaceById(c, ctx, workspaceService.WorkspaceQueries.GetWorkspaceByID)
	})

	group.DELETE("/:id", apphertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		workspaceId := c.Param("id")
		return fmt.Sprintf("Workspace-%s:Delete", workspaceId)
	}), func(c context.Context, ctx *app.RequestContext) {
		handlers.DeleteWorkspace(c, ctx, workspaceService.WorkspaceCommands.DeleteWorkspace)
	})

	group.PATCH("/:id", apphertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		workspaceId := c.Param("id")
		return fmt.Sprintf("Workspace-%s:Patch", workspaceId)
	}), func(c context.Context, ctx *app.RequestContext) {
		handlers.UpdateWorkspace(c, ctx, workspaceService.WorkspaceCommands.UpdateWorkspace)
	})

	group.POST("", apphertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		return "Workspace:Create"
	}), func(c context.Context, ctx *app.RequestContext) {
		handlers.CreateWorkspace(c, ctx, workspaceService.WorkspaceCommands.CreateWorkspace)
	})

	group.GET("", apphertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		return "Workspace:List"
	}), func(c context.Context, ctx *app.RequestContext) {
		handlers.ListWorkspaces(c, ctx, workspaceService.WorkspaceQueries.ListWorkspaces)
	})

	group.PUT("", apphertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		return "Workspace:Import"
	}), func(c context.Context, ctx *app.RequestContext) {
		handlers.ImportWorkspace(c, ctx, workspaceService.WorkspaceCommands.ImportWorkspace)
	})

	return
}

func addWorkflowRoute(group *route.RouterGroup, workspaceService *application.WorkspaceService) {
	group.POST("/:workspace-id/workflow", apphertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		return fmt.Sprintf("Workspace-%s:CreateWorkflow", c.Param("workspace-id"))
	}), func(c context.Context, ctx *app.RequestContext) {
		handlers.CreateWorkflow(c, ctx, workspaceService.WorkflowCommands.Create)
	})
	group.GET("/:workspace-id/workflow/:id", apphertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		return fmt.Sprintf("Workspace-%s:GetWorkflow", c.Param("workspace-id"))
	}), func(c context.Context, ctx *app.RequestContext) {
		handlers.GetWorkflow(c, ctx, workspaceService.WorkflowQueries.GetByID)
	})
	group.GET("/:workspace-id/workflow", apphertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		return fmt.Sprintf("Workspace-%s:ListWorkflow", c.Param("workspace-id"))
	}), func(c context.Context, ctx *app.RequestContext) {
		handlers.ListWorkflows(c, ctx, workspaceService.WorkflowQueries.ListWorkflows)
	})
	group.PATCH("/:workspace-id/workflow/:id", apphertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		return fmt.Sprintf("Workspace-%s:UpdateWorkflow", c.Param("workspace-id"))
	}), func(c context.Context, ctx *app.RequestContext) {
		handlers.UpdateWorkflow(c, ctx, workspaceService.WorkflowCommands.Update)
	})
	group.DELETE("/:workspace-id/workflow/:id", apphertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		return fmt.Sprintf("Workspace-%s:DeleteWorkflow", c.Param("workspace-id"))
	}), func(c context.Context, ctx *app.RequestContext) {
		handlers.DeleteWorkflow(c, ctx, workspaceService.WorkflowCommands.Delete)
	})
	group.GET("/:workspace-id/workflow/:workflow-id/version", apphertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		return fmt.Sprintf("Workspace-%s:ListWorkflowVersions", c.Param("workspace-id"))
	}), func(c context.Context, ctx *app.RequestContext) {
		handlers.ListWorkflowVersions(c, ctx, workspaceService.WorkflowQueries.ListVersions)
	})
	group.GET("/:workspace-id/workflow/:workflow-id/version/:id", apphertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		return fmt.Sprintf("Workspace-%s:GetWorkflowVersions", c.Param("workspace-id"))
	}), func(c context.Context, ctx *app.RequestContext) {
		handlers.GetWorkflowVersion(c, ctx, workspaceService.WorkflowQueries.GetVersion)
	})
	group.GET("/:workspace-id/workflow/:workflow-id/file", apphertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		return fmt.Sprintf("Workspace-%s:ListWorkflowFiles", c.Param("workspace-id"))
	}), func(c context.Context, ctx *app.RequestContext) {
		handlers.ListWorkflowFiles(c, ctx, workspaceService.WorkflowQueries.ListFiles)
	})
	group.GET("/:workspace-id/workflow/:workflow-id/file/:id", apphertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		return fmt.Sprintf("Workspace-%s:GetWorkflowFile", c.Param("workspace-id"))
	}), func(c context.Context, ctx *app.RequestContext) {
		handlers.GetWorkflowFile(c, ctx, workspaceService.WorkflowQueries.GetFile)
	})
}

func addNotebookRoute(group *route.RouterGroup, workspaceService *application.WorkspaceService) {
	group.PUT("/:workspace-id/notebook/:name", apphertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		return fmt.Sprintf("Workspace-%s:CreateNotebook", c.Param("workspace-id"))
	}), func(c context.Context, ctx *app.RequestContext) {
		handlers.CreateNotebook(c, ctx, workspaceService.NotebookCommands.Create)
	})
	group.GET("/:workspace-id/notebook/:name", apphertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		return fmt.Sprintf("Workspace-%s:GetNotebook", c.Param("workspace-id"))
	}), func(c context.Context, ctx *app.RequestContext) {
		handlers.GetNotebook(c, ctx, workspaceService.NotebookQueries.Get)
	})
	group.GET("/:workspace-id/notebook", apphertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		return fmt.Sprintf("Workspace-%s:ListNotebooks", c.Param("workspace-id"))
	}), func(c context.Context, ctx *app.RequestContext) {
		handlers.ListNotebooks(c, ctx, workspaceService.NotebookQueries.List)
	})
	group.DELETE("/:workspace-id/notebook/:name", apphertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		return fmt.Sprintf("Workspace-%s:DeleteNotebook", c.Param("workspace-id"))
	}), func(c context.Context, ctx *app.RequestContext) {
		handlers.DeleteNotebook(c, ctx, workspaceService.NotebookCommands.Delete)
	})
}

func addDataModelRouter(group *route.RouterGroup, service *application.WorkspaceService) {
	group.DELETE("/:workspace_id/data_model/:id", apphertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		return fmt.Sprintf("Workspace-%s:DeleteDataModel", c.Param("workspace_id"))
	}), func(c context.Context, ctx *app.RequestContext) {
		handlers.DeleteDataModel(c, ctx, service.DataModelCommands.DeleteDataModel)
	})

	group.PATCH("/:workspace_id/data_model", apphertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		return fmt.Sprintf("Workspace-%s:PatchDataModel", c.Param("workspace_id"))
	}), func(c context.Context, ctx *app.RequestContext) {
		handlers.PatchDataModel(c, ctx, service.DataModelCommands.PatchDataModel)
	})

	group.GET("/:workspace_id/data_model", apphertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		return fmt.Sprintf("Workspace-%s:ListDataModels", c.Param("workspace_id"))
	}), func(c context.Context, ctx *app.RequestContext) {
		handlers.ListDataModels(c, ctx, service.DataModelQueries.ListDataModels)
	})

	group.GET("/:workspace_id/data_model/:id", apphertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		return fmt.Sprintf("Workspace-%s:GetDataModel", c.Param("workspace_id"))
	}), func(c context.Context, ctx *app.RequestContext) {
		handlers.GetDataModel(c, ctx, service.DataModelQueries.GetDataModel)
	})

	group.GET("/:workspace_id/data_model/:id/rows", apphertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		return fmt.Sprintf("Workspace-%s:ListDataModelRows", c.Param("workspace_id"))
	}), func(c context.Context, ctx *app.RequestContext) {
		handlers.ListDataModelRows(c, ctx, service.DataModelQueries.ListDataModelRows)
	})

	group.GET("/:workspace_id/data_model/:id/rows/ids", apphertz.Authz(func(ctx context.Context, c *app.RequestContext) string {
		return fmt.Sprintf("Workspace-%s:ListAllDataModelRowIDs", c.Param("workspace-id"))
	}), func(c context.Context, ctx *app.RequestContext) {
		handlers.ListAllDataModelRowIDs(c, ctx, service.DataModelQueries.ListAllDataModelRowIDs)
	})
}
