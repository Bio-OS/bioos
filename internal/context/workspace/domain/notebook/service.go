package notebook

import (
	"context"
	"fmt"

	"github.com/Bio-OS/bioos/pkg/errors"
)

type Service interface {
	Upsert(context.Context, *Notebook) error
	Create(context.Context, *Notebook) error
	Update(context.Context, *Notebook) error
	Delete(ctx context.Context, path string) error
}

func NewService(repo Repository) Service {
	return &service{
		repository: repo,
	}
}

type service struct {
	repository Repository
}

func (s *service) Upsert(ctx context.Context, nb *Notebook) error {
	if err := s.repository.Save(ctx, nb); err != nil {
		return errors.NewInternalError(err)
	}
	return nil
}

func (s *service) Create(ctx context.Context, nb *Notebook) error {
	if stored, err := s.repository.Get(ctx, nb.Path()); err != nil {
		return errors.NewInternalError(fmt.Errorf("check notebook exist fail: %w", err))
	} else if stored != nil {
		return errors.NewAlreadyExistError("notebook", nb.Name)
	}
	return s.Upsert(ctx, nb)
}

func (s *service) Update(ctx context.Context, nb *Notebook) error {
	if stored, err := s.repository.Get(ctx, nb.Path()); err != nil {
		return errors.NewInternalError(fmt.Errorf("check notebook exist fail: %w", err))
	} else if stored == nil {
		return errors.NewNotFoundError("notebook", nb.Name)
	}
	return s.Upsert(ctx, nb)
}

func (s *service) Delete(ctx context.Context, path string) error {
	nb, err := s.repository.Get(ctx, path)
	if err != nil {
		return errors.NewInternalError(fmt.Errorf("check notebook exist fail: %w", err))
	} else if nb == nil {
		return errors.NewNotFoundError("notebook", path)
	}
	if err = s.repository.Delete(ctx, nb); err != nil {
		return errors.NewInternalError(err)
	}
	return nil
}
