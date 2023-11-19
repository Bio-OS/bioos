package submission

import (
	"context"
	"errors"

	"github.com/Bio-OS/bioos/internal/context/submission/application/query/run"
	submissionquery "github.com/Bio-OS/bioos/internal/context/submission/application/query/submission"
	"github.com/Bio-OS/bioos/internal/context/workspace/domain/workflow"
	"github.com/Bio-OS/bioos/internal/context/workspace/domain/workspace"
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
	workspaceproto "github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	"github.com/Bio-OS/bioos/pkg/consts"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	"github.com/Bio-OS/bioos/pkg/log"
	applog "github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/utils"
	"github.com/Bio-OS/bioos/pkg/utils/grpc"
)

type Service interface {
	Get(context.Context, string) (*Submission, error)
	Upsert(context.Context, *Submission) error
	Create(context.Context, *Submission) error
	Update(context.Context, *Submission) error
	Delete(ctx context.Context, id string) error
	SoftDelete(ctx context.Context, id string) error
	Cancel(ctx context.Context, id string) error
	CheckWorkspaceExist(ctx context.Context, workspaceID string) error
	CheckSubmissionExist(ctx context.Context, workspaceID, submissionName string) error
}

func NewService(grpcFactory grpc.Factory, repo Repository, eventbus eventbus.EventBus, submissionReadModel submissionquery.ReadModel, runReadModel run.ReadModel) Service {
	dataModelClient, err := grpcFactory.DataModelClient()
	if err != nil {
		log.Fatalf(err.Error())
	}
	workspaceClient, err := grpcFactory.WorkspaceClient()
	if err != nil {
		log.Fatalf(err.Error())
	}
	workflowClient, err := grpcFactory.WorkflowClient()
	if err != nil {
		log.Fatalf(err.Error())
	}
	svc := &service{
		repository:          repo,
		eventbus:            eventbus,
		submissionReadModel: submissionReadModel,
		runReadModel:        runReadModel,
		dataModelClient:     dataModelClient,
		workspaceClient:     workspaceClient,
		workflowClient:      workflowClient,
	}
	svc.subscribeEvents()
	return svc
}

type service struct {
	repository          Repository
	eventbus            eventbus.EventBus
	submissionReadModel submissionquery.ReadModel
	runReadModel        run.ReadModel
	dataModelClient     grpc.DataModelClient
	workflowClient      grpc.WorkflowClient
	workspaceClient     grpc.WorkspaceClient
}

func (s *service) Get(ctx context.Context, id string) (*Submission, error) {
	return s.repository.Get(ctx, id)
}

func (s *service) Upsert(ctx context.Context, submission *Submission) error {
	return s.repository.Save(ctx, submission)
}

func (s *service) Create(ctx context.Context, submission *Submission) error {
	if stored, err := s.repository.Get(ctx, submission.ID); err != nil {
		var apperror apperrors.Error
		if !(errors.As(err, &apperror) && (apperror.GetCode() == apperrors.NotFoundCode)) {
			return err
		}
	} else if stored != nil {
		return apperrors.NewAlreadyExistError("submission", submission.Name)
	}
	event := NewCreateEvent(submission.WorkspaceID, submission.ID, submission.WorkflowID, submission.WorkflowVersionID, submission.Language, submission.DataModelID, submission.DataModelRowIDs)
	if err := s.eventbus.Publish(ctx, event); err != nil {
		return apperrors.NewInternalError(err)
	}
	return s.Upsert(ctx, submission)
}

func (s *service) Update(ctx context.Context, submission *Submission) error {
	if stored, err := s.repository.Get(ctx, submission.ID); err != nil {
		return err
	} else if stored == nil {
		return apperrors.NewNotFoundError("submission", submission.Name)
	}
	return s.Upsert(ctx, submission)
}

func (s service) Delete(ctx context.Context, id string) error {
	submission, err := s.repository.Get(ctx, id)
	if err != nil {
		return err
	} else if submission == nil {
		return apperrors.NewNotFoundError("submission", submission.Name)
	}
	if err = s.repository.Delete(ctx, submission); err != nil {
		return err
	}
	return nil
}

func (s service) SoftDelete(ctx context.Context, id string) error {
	submission, err := s.repository.Get(ctx, id)
	if err != nil {
		return err
	} else if submission == nil {
		return apperrors.NewNotFoundError("submission", submission.Name)
	}
	event := NewDeleteSubmissionEvent(submission.ID, 0)
	if err = s.eventbus.Publish(ctx, event); err != nil {
		return apperrors.NewInternalError(err)
	}
	return s.repository.SoftDelete(ctx, submission)
}

func (s service) Cancel(ctx context.Context, id string) error {
	submission, err := s.repository.Get(ctx, id)
	if err != nil {
		return err
	} else if submission == nil {
		return apperrors.NewNotFoundError("submission", submission.Name)
	}
	if !utils.In(submission.Status, consts.AllowCancelSubmissionStatuses) {
		return apperrors.NewInvalidError("cannot cancel finished submission")
	}

	submission.Status = consts.SubmissionCancelling
	event := NewCancelSubmissionEvent(submission.ID, 0)
	if err = s.eventbus.Publish(ctx, event); err != nil {
		return apperrors.NewInternalError(err)
	}
	return s.repository.Save(ctx, submission)
}

func (s *service) CheckWorkspaceExist(ctx context.Context, workspaceID string) error {
	if _, err := s.workspaceClient.GetWorkspace(ctx, &workspaceproto.GetWorkspaceRequest{Id: workspaceID}); err != nil {
		return apperrors.NewInternalError(err)
	}
	return nil
}

func (s *service) CheckSubmissionExist(ctx context.Context, workspaceID, submissionName string) error {
	count, err := s.submissionReadModel.CountSubmissions(ctx, workspaceID, &submissionquery.ListSubmissionsFilter{Name: submissionName})
	if err != nil {
		return apperrors.NewInternalError(err)
	}
	if count > 0 {
		return apperrors.NewAlreadyExistError("submission", submissionName)
	}
	return nil
}

func (s *service) subscribeEvents() {
	s.eventbus.Subscribe(CreateSubmission, eventbus.EventHandlerFunc(func(ctx context.Context, payload string) error {
		applog.Infow("start to consume submission create event", "payload", payload)
		event, err := NewCreateEventFromPayload([]byte(payload))
		if err != nil {
			return err
		}

		handler := NewCreateHandler(s.repository, s.eventbus, s.workflowClient)
		return handler.Handle(ctx, event)
	}))

	s.eventbus.Subscribe(CancelSubmission, eventbus.EventHandlerFunc(func(ctx context.Context, payload string) error {
		applog.Infow("start to consume submission cancel event", "payload", payload)
		event, err := NewEventFromPayload([]byte(payload))
		if err != nil {
			return err
		}

		handler := NewCancelHandler(s.repository, s.eventbus, s.runReadModel)
		return handler.Handle(ctx, event)
	}))

	s.eventbus.Subscribe(DeleteSubmission, eventbus.EventHandlerFunc(func(ctx context.Context, payload string) error {
		applog.Infow("start to consume submission deleted event", "payload", payload)
		event, err := NewEventFromPayload([]byte(payload))
		if err != nil {
			return err
		}

		handler := NewDeleteHandler(s.repository, s.eventbus, s.runReadModel)
		return handler.Handle(ctx, event)
	}))

	s.eventbus.Subscribe(SyncSubmission, eventbus.EventHandlerFunc(func(ctx context.Context, payload string) error {
		applog.Infow("start to consume submission sync event", "payload", payload)
		event, err := NewEventFromPayload([]byte(payload))
		if err != nil {
			return err
		}

		handler := NewSyncHandler(s.repository, s.runReadModel, s.dataModelClient)
		return handler.Handle(ctx, event)
	}))

	s.eventbus.Subscribe(CascadeDeleteSubmission, eventbus.EventHandlerFunc(func(ctx context.Context, payload string) (err error) {
		applog.Infow("start to consume submission cascade deleted event", "payload", payload)

		event, err := NewEventCascadeDeleteSubmissionFromPayload([]byte(payload))
		if err != nil {
			return err
		}

		handler := NewCascadeDeleteHandler(s.repository, s.eventbus, s.submissionReadModel)
		return handler.Handle(ctx, event)
	}))

	s.eventbus.Subscribe(workflow.WorkflowDeleted, eventbus.EventHandlerFunc(func(ctx context.Context, payload string) error {
		applog.Infow("start to consume workflow deleted event", "payload", payload)
		event, err := workflow.NewWorkflowEventFromPayload([]byte(payload))
		if err != nil {
			return err
		}
		deleteSubmissionEvent := NewEventCascadeDeleteSubmission(event.WorkspaceID, utils.PointString(event.WorkflowID))

		if err = s.eventbus.Publish(ctx, deleteSubmissionEvent); err != nil {
			return apperrors.NewInternalError(err)
		}
		return nil
	}))

	s.eventbus.Subscribe(workspace.WorkspaceDeleted, eventbus.EventHandlerFunc(func(ctx context.Context, payload string) (err error) {
		applog.Infow("start to consume workspace deleted event", "payload", payload)

		event, err := workspace.NewWorkspaceEventFromPayload([]byte(payload))
		if err != nil {
			return err
		}

		deleteSubmissionEvent := NewEventCascadeDeleteSubmission(event.WorkspaceID, nil)

		if err = s.eventbus.Publish(ctx, deleteSubmissionEvent); err != nil {
			return apperrors.NewInternalError(err)
		}
		return nil
	}))
}
