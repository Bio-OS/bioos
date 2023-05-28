package submission

import (
	"context"
	"encoding/json"

	"github.com/Bio-OS/bioos/internal/context/submission/domain/submission"
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
	"github.com/Bio-OS/bioos/pkg/consts"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	"github.com/Bio-OS/bioos/pkg/validator"
)

type CreateSubmissionHandler interface {
	Handle(ctx context.Context, cmd *CreateSubmissionCommand) (string, error)
}

type createSubmissionHandler struct {
	service           submission.Service
	submissionFactory *submission.Factory
	eventBus          eventbus.EventBus
}

var _ CreateSubmissionHandler = &createSubmissionHandler{}

func NewCreateSubmissionHandler(service submission.Service, submissionFactory *submission.Factory, eventBus eventbus.EventBus) CreateSubmissionHandler {
	return &createSubmissionHandler{
		service:           service,
		submissionFactory: submissionFactory,
		eventBus:          eventBus,
	}
}

func (c *createSubmissionHandler) Handle(ctx context.Context, cmd *CreateSubmissionCommand) (string, error) {
	if err := validator.Validate(cmd); err != nil {
		return "", err
	}

	if err := c.service.CheckWorkspaceExist(ctx, cmd.WorkspaceID); err != nil {
		return "", err
	}

	if err := c.service.CheckSubmissionExist(ctx, cmd.WorkspaceID, cmd.Name); err != nil {
		return "", err
	}

	param := submission.CreateSubmissionParam{
		Name:        cmd.Name,
		Description: cmd.Description,
		WorkflowID:  cmd.WorkflowID,
		WorkspaceID: cmd.WorkspaceID,
		Type:        cmd.Type,
		ExposedOptions: submission.ExposedOptions{
			ReadFromCache: cmd.ExposedOptions.ReadFromCache,
		},
		Inputs:  make(map[string]interface{}),
		Outputs: make(map[string]interface{}),
	}

	switch param.Type {
	case consts.DataModelTypeSubmission:
		if len(cmd.Entity.DataModelRowIDs) == 0 || len(cmd.Entity.DataModelID) == 0 {
			return "", apperrors.NewInvalidError("data model id & row ids should not empty")
		}
		param.DataModelID = &cmd.Entity.DataModelID
		param.DataModelRowIDs = cmd.Entity.DataModelRowIDs
		var inputs map[string]interface{}
		if err := json.Unmarshal([]byte(cmd.Entity.InputsTemplate), &inputs); err != nil {
			return "", apperrors.NewInvalidError(err.Error())
		}
		param.Inputs = inputs
		if cmd.Entity.OutputsTemplate != "" {
			var output map[string]interface{}
			if err := json.Unmarshal([]byte(cmd.Entity.OutputsTemplate), &output); err != nil {
				return "", apperrors.NewInvalidError(err.Error())
			}
			param.Outputs = output
		}
	case consts.FilePathTypeSubmission:
		if cmd.InOutMaterial.OutputsMaterial != "" {
			var output map[string]interface{}
			if err := json.Unmarshal([]byte(cmd.InOutMaterial.OutputsMaterial), &output); err != nil {
				return "", apperrors.NewInvalidError(err.Error())
			}
			param.Outputs = output
		}
		var inputs map[string]interface{}
		err := json.Unmarshal([]byte(cmd.InOutMaterial.InputsMaterial), &inputs)
		if err != nil {
			return "", apperrors.NewInvalidError(err.Error())
		}
		param.Inputs = inputs
	default:
		return "", apperrors.NewInvalidError("unsupported submission type: %s", param.Type)
	}
	sub, err := c.submissionFactory.CreateWithSubmissionParam(param)
	if err != nil {
		return "", err
	}
	if err = c.service.Create(ctx, sub); err != nil {
		return "", err
	}
	return sub.ID, nil
}
