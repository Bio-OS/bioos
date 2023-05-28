package domain

import (
	"context"
	"fmt"
	"reflect"

	"github.com/Bio-OS/bioos/pkg/errors"
	"github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/notebook"
)

type Service interface {
	Create(context.Context, *NotebookServer) error
	Update(context.Context, *NotebookServer) error
	Start(ctx context.Context, id string) error
	Stop(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
}

type service struct {
	repository Repository
	runtime    Runtime // TODO new it by cluster
}

func NewService(repo Repository, runtime Runtime) Service {
	return &service{
		repository: repo,
		runtime:    runtime,
	}
}

func (s *service) Create(ctx context.Context, srv *NotebookServer) error {
	if err := s.repository.Save(ctx, srv); err != nil {
		return errors.NewInternalError(fmt.Errorf("write settings fail: %w", err))
	}
	if err := s.runtime.Create(ctx, srv); err != nil {
		if err2 := s.repository.Delete(ctx, srv); err2 != nil {
			log.Errorf("revert notebookserver %s settings creating fail: %s", srv.ID, err2)
		}
		return errors.NewInternalError(fmt.Errorf("runtime create fail: %w", err))
	}
	return nil
}

func (s *service) Update(ctx context.Context, srv *NotebookServer) error {
	stored, err := s.repository.Get(ctx, srv.ID)
	if err != nil {
		return errors.NewInternalError(fmt.Errorf("check notebook exist fail: %w", err))
	} else if stored == nil {
		return errors.NewNotFoundError("notebookserver", srv.ID)
	}

	if srv.Settings.DockerImage != "" {
		stored.Settings.DockerImage = srv.Settings.DockerImage
	}
	if !reflect.DeepEqual(srv.Settings.ResourceSize, notebook.ResourceSize{}) &&
		!reflect.DeepEqual(srv.Settings.ResourceSize, stored.Settings.ResourceSize) {
		// TODO check if size in options
		stored.Settings.ResourceSize = srv.Settings.ResourceSize
	}

	if err := s.repository.Save(ctx, stored); err != nil {
		return errors.NewInternalError(err)
	}
	return nil
}

func (s *service) Start(ctx context.Context, id string) error {
	srv, err := s.repository.Get(ctx, id)
	if err != nil {
		return errors.NewInternalError(fmt.Errorf("check notebook exist fail: %w", err))
	} else if srv == nil {
		return errors.NewNotFoundError("notebookserver", id)
	}
	return s.runtime.Start(ctx, srv)
}

func (s *service) Stop(ctx context.Context, id string) error {
	srv, err := s.repository.Get(ctx, id)
	if err != nil {
		return errors.NewInternalError(fmt.Errorf("check notebook exist fail: %w", err))
	} else if srv == nil {
		return errors.NewNotFoundError("notebookserver", id)
	}
	return s.runtime.Stop(ctx, srv)
}

func (s *service) Delete(ctx context.Context, id string) error {
	srv, err := s.repository.Get(ctx, id)
	if err != nil {
		return errors.NewInternalError(fmt.Errorf("check notebook exist fail: %w", err))
	} else if srv == nil {
		return errors.NewNotFoundError("notebookserver", id)
	}
	if err = s.repository.Delete(ctx, srv); err != nil {
		return errors.NewInternalError(fmt.Errorf("delete settings fail: %w", err))
	}
	if err = s.runtime.Delete(ctx, srv); err != nil {
		if err2 := s.repository.Save(ctx, srv); err2 != nil {
			log.Errorf("revert notebookserver %s settings deleting fail: %s", id, err2)
		}
		return errors.NewInternalError(fmt.Errorf("delete resource fail: %w", err))
	}
	return nil
}
