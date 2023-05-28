package workflow

import (
	"context"
	stderrs "errors"

	errors "github.com/go-kratos/kratos/v2/errors"

	workflowquery "github.com/Bio-OS/bioos/internal/context/workspace/application/query/workflow"
	workspacequery "github.com/Bio-OS/bioos/internal/context/workspace/application/query/workspace"
	"github.com/Bio-OS/bioos/internal/context/workspace/domain/workflow"
	"github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	"github.com/Bio-OS/bioos/pkg/validator"
)

type CreateCommand struct {
	ID               string
	Name             string
	Description      *string
	WorkspaceID      string
	Language         string
	Source           string
	URL              string
	Tag              string
	Token            string
	MainWorkflowPath string
}

type CreateHandler interface {
	Handle(context.Context, *CreateCommand) (string, error)
}

type createHandler struct {
	service            workflow.Service
	workspaceReadModel workspacequery.WorkspaceReadModel
	workflowReadModel  workflowquery.ReadModel
}

func NewCreateHandler(service workflow.Service, workflowReadModel workflowquery.ReadModel, workspaceReadModel workspacequery.WorkspaceReadModel) CreateHandler {
	return &createHandler{
		service:            service,
		workspaceReadModel: workspaceReadModel,
		workflowReadModel:  workflowReadModel,
	}
}

func (h *createHandler) Handle(ctx context.Context, cmd *CreateCommand) (string, error) {
	if err := validator.Validate(cmd); err != nil {
		return "", err
	}
	// check workspace exist
	_, err := h.workspaceReadModel.GetWorkspaceById(ctx, cmd.WorkspaceID)
	if err != nil {
		return "", err
	}

	// check name duplicate
	_, err = h.workflowReadModel.GetByName(ctx, cmd.WorkspaceID, cmd.Name)
	var e *errors.Error
	if stderrs.As(err, &e) && e.Code == 404 {
		// not exist same name, do nothing
	} else if err != nil {
		return "", err
	} else if err == nil {
		return "", proto.ErrorWorkflowNameDuplicated("workflow update failed name:%s already exists", cmd.Name)
	}

	workflowID, _, err := h.service.AddVersion(ctx, cmd.WorkspaceID,
		&workflow.WorkflowOption{
			ID:          cmd.ID,
			Name:        cmd.Name,
			Description: cmd.Description,
		},
		&workflow.VersionOption{
			Language:         cmd.Language,
			MainWorkflowPath: cmd.MainWorkflowPath,
			Source:           cmd.Source,
			Url:              cmd.URL,
			Tag:              cmd.Tag,
			Token:            cmd.Token,
		})
	if err != nil {
		return "", err
	}

	return workflowID, nil
}
