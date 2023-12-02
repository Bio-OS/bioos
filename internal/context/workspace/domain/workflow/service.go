package workflow

import (
	"context"

	"github.com/Bio-OS/bioos/internal/context/workspace/application/query/workflow"
	"github.com/Bio-OS/bioos/internal/context/workspace/domain/workspace"
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
	"github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	applog "github.com/Bio-OS/bioos/pkg/log"
)

type Service interface {
	Create(ctx context.Context, workspaceID string, workflowOption *WorkflowOption) (string, error)
	Delete(ctx context.Context, workspaceID string, workflowID string) error
	Update(ctx context.Context, workspaceID string, workflowID string, workflowOption *WorkflowOption) error
	AddVersion(ctx context.Context, workspaceID string, workflowOption *WorkflowOption, versionOption *VersionOption) (string, *WorkflowVersion, error)
	UpdateVersion(ctx context.Context, workspaceID string, workflowOption *WorkflowOption, versionOption *VersionUpdateOption) error
}

type service struct {
	readModel  workflow.ReadModel
	repository Repository
	eventbus   eventbus.EventBus
	factory    *Factory
}

func NewService(repo Repository, readModel workflow.ReadModel, bus eventbus.EventBus, factory *Factory, womtoolPath string) Service {
	svc := &service{
		readModel:  readModel,
		repository: repo,
		eventbus:   bus,
		factory:    factory,
	}
	svc.subscribeEvents(womtoolPath)
	return svc
}

func (s *service) Update(ctx context.Context, workspaceID, workflowID string, workflowOption *WorkflowOption) error {
	if workflowOption == nil {
		return nil
	}
	workflow, err := s.repository.Get(ctx, workspaceID, workflowID)
	if err != nil {
		if proto.IsWorkflowNotFound(err) {
			return proto.ErrorWorkflowNotFound("workflow:%s in workspace:%s not found", workflowID, workspaceID)
		}
		return err
	}
	if workflowOption.Name != "" {
		workflow.UpdateName(workflowOption.Name)
	}
	if workflowOption.Description != nil {
		workflow.UpdateDescription(*workflowOption.Description)
	}

	return s.repository.Save(ctx, workflow)
}

// Create new a workflow and return created workflow id.
func (s *service) Create(ctx context.Context, workspaceID string, workflowOption *WorkflowOption) (workflowID string, err error) {
	workflow, err := s.factory.NewWorkflow(workspaceID, workflowOption)
	if err != nil {
		return "", err
	}

	event := NewWorkflowCreatedEvent(workflow.WorkspaceID, workflow.ID)

	if err := s.repository.Save(ctx, workflow); err != nil {
		return "", err
	}

	return workflow.ID, s.eventbus.Publish(ctx, event)
}

// Delete delete a workflow by id
func (s *service) Delete(ctx context.Context, workspaceID, workflowID string) error {
	workflow, err := s.repository.Get(ctx, workspaceID, workflowID)
	if err != nil {
		if proto.IsWorkflowNotFound(err) {
			return nil
		}
		return err
	}

	event := NewWorkflowDeletedEvent(workflow.WorkspaceID, workflow.ID)

	if err := s.repository.Delete(ctx, workflow); err != nil {
		return err
	}

	return s.eventbus.Publish(ctx, event)
}

// AddVersion add workflow version
func (s *service) AddVersion(ctx context.Context, workspaceID string, workflowOption *WorkflowOption, versionOption *VersionOption) (workflowID string, version *WorkflowVersion, err error) {
	// check params
	if workflowOption == nil || versionOption == nil {
		return "", nil, apperrors.NewInvalidError("workflow or version param is nil")
	}
	var workflow *Workflow
	if workflowOption.ID == "" {
		// need to create workflow
		workflow, err = s.factory.NewWorkflow(workspaceID, workflowOption)
		if err != nil {
			return "", nil, err
		}
	} else {
		workflow, err = s.repository.Get(ctx, workspaceID, workflowOption.ID)
		if err != nil {
			return "", nil, err
		}
	}

	version, err = workflow.AddVersion(versionOption)
	if err != nil {
		return "", nil, err
	}

	//TODO support create by file through http/grpc in the future
	event := NewWorkflowVersionAddedEvent(workspaceID, workflow.ID, version.ID, versionOption.Url, versionOption.Tag, versionOption.Token, "")

	if err := s.repository.Save(ctx, workflow); err != nil {
		return "", nil, err
	}

	//TODO if save success but publish failed?
	return workflow.ID, version, s.eventbus.Publish(ctx, event)
}

// UpdateVersion update workflow version
func (s *service) UpdateVersion(ctx context.Context, workspaceID string, workflowOpt *WorkflowOption, updateOpt *VersionUpdateOption) error {
	// need to create workflow
	applog.Infow("UpdateVersion", "workflowOpt", workflowOpt, "updateOpt", updateOpt)
	workflow, err := s.repository.Get(ctx, workspaceID, workflowOpt.ID)
	if err != nil {
		return err
	}

	workflowVersion, exist := workflow.Versions[workflow.LatestVersion]
	if !exist {
		return proto.ErrorWorkflowVersionNotFound("workflow:%s version:%s not found", workflow.ID, workflow.LatestVersion)
	}

	if updateOpt == nil {
		return nil
	}

	// check if only update token
	addVersionFlag := false
	versionOption := &VersionOption{
		Language:         workflowVersion.Language,
		MainWorkflowPath: workflowVersion.MainWorkflowPath,
		Source:           workflowVersion.Source,
		Url:              workflowVersion.Metadata[WorkflowGitURL],
		Tag:              workflowVersion.Metadata[WorkflowGitTag],
		Token:            "",
	}
	if updateOpt.Language != nil && *updateOpt.Language != workflowVersion.Language {
		versionOption.Language = *updateOpt.Language
		addVersionFlag = true
	}
	if updateOpt.MainWorkflowPath != nil && *updateOpt.MainWorkflowPath != workflowVersion.MainWorkflowPath {
		versionOption.MainWorkflowPath = *updateOpt.MainWorkflowPath
		addVersionFlag = true
	}
	if updateOpt.Source != nil && *updateOpt.Source != workflowVersion.Source {
		versionOption.Source = *updateOpt.Source
		addVersionFlag = true
	}
	if updateOpt.URL != nil && *updateOpt.URL != workflowVersion.Metadata[WorkflowGitURL] {
		versionOption.Url = *updateOpt.URL
		addVersionFlag = true
	}
	if updateOpt.Tag != nil && *updateOpt.Tag != workflowVersion.Metadata[WorkflowGitTag] {
		versionOption.Tag = *updateOpt.Tag
		addVersionFlag = true
	}
	if updateOpt.Token != nil {
		versionOption.Token = *updateOpt.Token
		addVersionFlag = true
	}
	applog.Infow("need to add version", "addVersionFlag", addVersionFlag)
	// need to add new version
	if addVersionFlag {
		version, err := workflow.AddVersion(versionOption)
		if err != nil {
			return err
		}
		//TODO support create by file through http/grpc in the future
		event := NewWorkflowVersionAddedEvent(workspaceID, workflow.ID, version.ID, versionOption.Url, versionOption.Tag, versionOption.Token, "")

		if err := s.repository.Save(ctx, workflow); err != nil {
			return err
		}
		return s.eventbus.Publish(ctx, event)
	}

	return nil
}

func (s *service) subscribeEvents(womtoolPath string) {
	s.eventbus.Subscribe(WorkflowVersionAdded, eventbus.EventHandlerFunc(func(ctx context.Context, payload string) (err error) {
		applog.Infow("start to consume workflow version added event", "payload", payload)

		event, err := NewWorkflowVersionAddedEventFromPayload([]byte(payload))
		if err != nil {
			return err
		}

		handler := NewWorkflowVersionAddedHandler(s.repository, &ReaderOptions{WomtoolPath: womtoolPath})
		return handler.Handle(ctx, event)
	}))

	s.eventbus.Subscribe(workspace.WorkspaceDeleted, eventbus.EventHandlerFunc(func(ctx context.Context, payload string) (err error) {
		applog.Infow("start to consume workspace deleted event", "payload", payload)

		event, err := workspace.NewWorkspaceEventFromPayload([]byte(payload))
		if err != nil {
			return err
		}

		handler := NewWorkspaceDeletedHandler(s.repository)
		return handler.Handle(ctx, event)
	}))

	s.eventbus.Subscribe(ImportWorkflows, eventbus.EventHandlerFunc(func(ctx context.Context, payload string) (err error) {
		applog.Infow("start to consume import workflows event", "payload", payload)

		event, err := NewImportWorkflowsEventFromPayload([]byte(payload))
		if err != nil {
			return err
		}

		handler := NewImportWorkflowsHandler(s.repository, s.readModel, s.eventbus, s.factory)
		return handler.Handle(ctx, event)
	}))
}
