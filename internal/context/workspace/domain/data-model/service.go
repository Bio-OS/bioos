package datamodel

import (
	"context"
	"errors"
	"fmt"

	"github.com/Bio-OS/bioos/internal/context/workspace/infrastructure/eventbus"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	applog "github.com/Bio-OS/bioos/pkg/log"
)

type Service interface {
	Get(context.Context, string) (*DataModel, error)
	Upsert(context.Context, *DataModel) error
	Create(context.Context, *DataModel) error
	Update(context.Context, *DataModel) error
	Delete(context.Context, *DataModel) error
}

func NewService(repo Repository, eventbus eventbus.EventBus, factory *Factory) Service {
	svc := &service{
		repository: repo,
		eventbus:   eventbus,
		factory:    factory,
	}
	svc.subscribeEvents()
	return svc
}

type service struct {
	repository Repository
	eventbus   eventbus.EventBus
	factory    *Factory
}

func (s *service) Get(ctx context.Context, id string) (*DataModel, error) {
	return s.repository.Get(ctx, id)
}

func (s *service) Upsert(ctx context.Context, dm *DataModel) error {
	return s.repository.Save(ctx, dm)
}

func (s *service) Create(ctx context.Context, dm *DataModel) error {
	if stored, err := s.repository.Get(ctx, dm.ID); err != nil {
		var apperror apperrors.Error
		if !(errors.As(err, &apperror) && (apperror.GetCode() == apperrors.NotFoundCode)) {
			return err
		}
	} else if stored != nil {
		return apperrors.NewAlreadyExistError("data model", dm.Name)
	}
	return s.Upsert(ctx, dm)
}

func (s *service) Update(ctx context.Context, dm *DataModel) error {
	if stored, err := s.repository.Get(ctx, dm.ID); err != nil {
		return apperrors.NewInternalError(fmt.Errorf("check data model exist fail: %w", err))
	} else if stored == nil {
		return apperrors.NewNotFoundError("data model", dm.Name)
	}
	return s.Upsert(ctx, dm)
}

func (s *service) Delete(ctx context.Context, dm *DataModel) error {
	return s.repository.Delete(ctx, dm)
}

func (s *service) subscribeEvents() {
	s.eventbus.Subscribe(ImportDataModels, eventbus.EventHandlerFunc(func(ctx context.Context, payload string) (err error) {
		applog.Infow("start to consume import data-models event", "payload", payload)

		event, err := NewImportDataModelsEventFromPayload([]byte(payload))
		if err != nil {
			return err
		}

		handler := NewImportDataModelsHandler(s.repository, s.eventbus, s.factory)
		return handler.Handle(ctx, event)
	}))

}
