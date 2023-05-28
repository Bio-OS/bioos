package workflow

import (
	"context"

	workspacequery "github.com/Bio-OS/bioos/internal/context/workspace/application/query/workspace"
	"github.com/Bio-OS/bioos/internal/context/workspace/domain/workflow"
	"github.com/Bio-OS/bioos/pkg/validator"
)

type DeleteCommand struct {
	ID          string `validate:"required"`
	WorkspaceID string `validate:"required"`
}

type DeleteHandler interface {
	Handle(ctx context.Context, cmd *DeleteCommand) error
}

func NewDeleteHandler(workflowService workflow.Service, workspaceReadModel workspacequery.WorkspaceReadModel) DeleteHandler {
	return &deleteHandler{
		workflowService:    workflowService,
		workspaceReadModel: workspaceReadModel,
	}
}

type deleteHandler struct {
	workflowService    workflow.Service
	workspaceReadModel workspacequery.WorkspaceReadModel
}

func (h *deleteHandler) Handle(ctx context.Context, cmd *DeleteCommand) error {
	if err := validator.Validate(cmd); err != nil {
		return err
	}
	// check workspace exist
	_, err := h.workspaceReadModel.GetWorkspaceById(ctx, cmd.WorkspaceID)
	if err != nil {
		return err
	}

	return h.workflowService.Delete(ctx, cmd.WorkspaceID, cmd.ID)
}

var _ DeleteHandler = &deleteHandler{}
