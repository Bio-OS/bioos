package run

import (
	"context"
	"errors"

	"github.com/Bio-OS/bioos/internal/context/submission/domain/submission"
	"github.com/Bio-OS/bioos/internal/context/submission/infrastructure/client/wes"
	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
	workspaceproto "github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	"github.com/Bio-OS/bioos/pkg/consts"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	"github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/utils"
	"github.com/Bio-OS/bioos/pkg/utils/grpc"
)

type Service interface {
	Upsert(context.Context, *Run, []*Task) error
	Create(context.Context, *Run, []*Task) error
	Update(context.Context, *Run, []*Task) error
	Delete(ctx context.Context, id string) error
	Cancel(ctx context.Context, id string) error
	CheckWorkspaceExist(ctx context.Context, workspaceID string) error
}

func NewService(grpcFactory grpc.Factory, repo Repository, eventbus eventbus.EventBus, wesClient wes.Client) Service {
	dataModelClient, err := grpcFactory.DataModelClient()
	if err != nil {
		log.Fatalf(err.Error())
	}
	workspaceClient, err := grpcFactory.WorkspaceClient()
	if err != nil {
		log.Fatalf(err.Error())
	}
	svc := &service{
		repository:      repo,
		eventbus:        eventbus,
		factory:         Factory{},
		dataModelClient: dataModelClient,
		workspaceClient: workspaceClient,
		wesClient:       wesClient,
	}
	svc.subscribeEvents()
	return svc
}

type service struct {
	repository      Repository
	eventbus        eventbus.EventBus
	factory         Factory
	dataModelClient grpc.DataModelClient
	workspaceClient grpc.WorkspaceClient
	wesClient       wes.Client
}

func (s *service) Upsert(ctx context.Context, run *Run, tasks []*Task) error {
	return s.repository.Save(ctx, run)
}

func (s *service) Create(ctx context.Context, run *Run, tasks []*Task) error {
	if stored, err := s.repository.Get(ctx, run.ID); err != nil {
		var apperror apperrors.Error
		if !(errors.As(err, &apperror) && (apperror.GetCode() == apperrors.NotFoundCode)) {
			return err
		}
	} else if stored != nil {
		return apperrors.NewAlreadyExistError("run", run.Name)
	}
	return s.Upsert(ctx, run, tasks)
}

func (s *service) Update(ctx context.Context, run *Run, tasks []*Task) error {
	if stored, err := s.repository.Get(ctx, run.ID); err != nil {
		return err
	} else if stored == nil {
		return apperrors.NewNotFoundError("run", run.Name)
	}
	return s.Upsert(ctx, run, tasks)
}

func (s *service) Delete(ctx context.Context, id string) error {
	run, err := s.repository.Get(ctx, id)
	if err != nil {
		return err
	} else if run == nil {
		return apperrors.NewNotFoundError("run", run.Name)
	}
	return s.repository.Delete(ctx, run)
}

func (s *service) Cancel(ctx context.Context, id string) error {
	run, err := s.repository.Get(ctx, id)
	if err != nil {
		return err
	} else if run == nil {
		return apperrors.NewNotFoundError("run", run.Name)
	}
	if !utils.In(run.Status, consts.AllowCancelRunStatuses) {
		return apperrors.NewInvalidError("cannot cancel finished run")
	}

	event := submission.NewEventCancelRun(id)
	if err = s.eventbus.Publish(ctx, event); err != nil {
		return apperrors.NewInternalError(err)
	}

	run.Status = consts.RunCancelling
	return s.repository.Save(ctx, run)
}

func (s *service) CheckWorkspaceExist(ctx context.Context, workspaceID string) error {
	if _, err := s.workspaceClient.GetWorkspace(ctx, &workspaceproto.GetWorkspaceRequest{Id: workspaceID}); err != nil {
		return apperrors.NewInternalError(err)
	}
	return nil
}

func (s *service) subscribeEvents() {
	s.eventbus.Subscribe(submission.CreateRuns, eventbus.EventHandlerFunc(func(ctx context.Context, payload string) error {
		event, err := submission.NewEventCreateRunFromPayload([]byte(payload))
		if err != nil {
			return err
		}

		handler := NewEventHandlerCreateRuns(s.repository, s.dataModelClient, s.eventbus, s.factory)
		return handler.Handle(ctx, event)
	}))

	s.eventbus.Subscribe(submission.SubmitRun, eventbus.EventHandlerFunc(func(ctx context.Context, payload string) error {
		event, err := submission.NewEventSubmitRunFromPayload([]byte(payload))
		if err != nil {
			return err
		}

		handler := NewEventHandlerSubmitRun(s.wesClient, s.eventbus, s.repository)
		return handler.Handle(ctx, event)
	}))

	s.eventbus.Subscribe(submission.SyncRun, eventbus.EventHandlerFunc(func(ctx context.Context, payload string) error {
		event, err := submission.NewEventRunFromPayload([]byte(payload))
		if err != nil {
			return err
		}

		handler := NewEventHandlerSyncRun(s.wesClient, s.repository, s.factory, s.eventbus)
		return handler.Handle(ctx, event)
	}))
	s.eventbus.Subscribe(submission.DeleteRun, eventbus.EventHandlerFunc(func(ctx context.Context, payload string) error {
		event, err := submission.NewEventRunFromPayload([]byte(payload))
		if err != nil {
			return err
		}

		handler := NewEventHandlerDeleteRun(s.repository, s.eventbus)
		return handler.Handle(ctx, event)
	}))
	s.eventbus.Subscribe(submission.CancelRun, eventbus.EventHandlerFunc(func(ctx context.Context, payload string) error {
		event, err := submission.NewEventRunFromPayload([]byte(payload))
		if err != nil {
			return err
		}

		handler := NewEventHandlerCancelRun(s.wesClient, s.repository, s.eventbus)
		return handler.Handle(ctx, event)
	}))
}
