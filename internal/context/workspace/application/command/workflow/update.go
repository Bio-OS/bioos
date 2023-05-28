package workflow

import (
	"context"
	stderrs "errors"

	"github.com/go-kratos/kratos/v2/errors"

	workflowquery "github.com/Bio-OS/bioos/internal/context/workspace/application/query/workflow"
	workspacequery "github.com/Bio-OS/bioos/internal/context/workspace/application/query/workspace"
	"github.com/Bio-OS/bioos/internal/context/workspace/domain/workflow"
	"github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	"github.com/Bio-OS/bioos/pkg/validator"
)

type UpdateCommand struct {
	WorkspaceID      string
	ID               string
	Name             *string
	Description      *string
	Language         *string
	Source           *string
	URL              *string
	Tag              *string
	Token            *string
	MainWorkflowPath *string
}

type UpdateHandler interface {
	Handle(ctx context.Context, cmd *UpdateCommand) error
}

func NewUpdateHandler(workflowService workflow.Service, workflowReadModel workflowquery.ReadModel, workspaceReadModel workspacequery.WorkspaceReadModel) UpdateHandler {
	return &updateHandler{
		workflowService:    workflowService,
		workflowReadModel:  workflowReadModel,
		workspaceReadModel: workspaceReadModel,
	}
}

type updateHandler struct {
	workflowService    workflow.Service
	workflowReadModel  workflowquery.ReadModel
	workspaceReadModel workspacequery.WorkspaceReadModel
}

func (h updateHandler) Handle(ctx context.Context, cmd *UpdateCommand) error {
	if err := validator.Validate(cmd); err != nil {
		return err
	}

	// check workspace exist
	_, err := h.workspaceReadModel.GetWorkspaceById(ctx, cmd.WorkspaceID)
	if err != nil {
		return err
	}

	wf, err := h.workflowReadModel.GetById(ctx, cmd.WorkspaceID, cmd.ID)
	if err != nil {
		return err
	}

	// check name duplicate
	if cmd.Name != nil && wf.Name != *cmd.Name {
		_, err = h.workflowReadModel.GetByName(ctx, cmd.WorkspaceID, *cmd.Name)
		var e *errors.Error
		if stderrs.As(err, &e) && e.Code == 404 {
			// not exist same name, do nothing
		} else if err != nil {
			return err
		} else if err == nil {
			return proto.ErrorWorkflowNameDuplicated("workflow update failed name:%s already exists", *cmd.Name)
		}
	}

	// update workflow
	if (cmd.Name != nil && *cmd.Name != wf.Name) || (cmd.Description != nil && *cmd.Description != wf.Description) {
		workflowParam := &workflow.WorkflowOption{}
		if cmd.Name != nil {
			workflowParam.Name = *cmd.Name
		}
		if cmd.Description != nil {
			workflowParam.Description = cmd.Description
		}
		if err := h.workflowService.Update(ctx, cmd.WorkspaceID, cmd.ID, workflowParam); err != nil {
			return err
		}
	}

	// update workflow version if repo info updated
	if cmd.Language != nil || cmd.MainWorkflowPath != nil ||
		cmd.Source != nil || cmd.URL != nil || cmd.Tag != nil || cmd.Token != nil {
		versionUpdateOpt := &workflow.VersionUpdateOption{
			Language:         cmd.Language,
			MainWorkflowPath: cmd.MainWorkflowPath,
			Source:           cmd.Source,
			URL:              cmd.URL,
			Tag:              cmd.Tag,
		}

		if cmd.Token != nil {
			versionUpdateOpt.Token = cmd.Token
		}

		if err := h.workflowService.UpdateVersion(ctx, cmd.WorkspaceID,
			&workflow.WorkflowOption{
				ID: cmd.ID,
			}, versionUpdateOpt); err != nil {
			return err
		}
	}

	return nil
}

var _ UpdateHandler = &updateHandler{}
