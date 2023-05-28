package sql

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/Bio-OS/bioos/internal/context/notebookserver/domain"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(ctx context.Context, db *gorm.DB) (domain.Repository, error) {
	if err := db.WithContext(ctx).AutoMigrate(&notebookServer{}); err != nil {
		return nil, fmt.Errorf("notebookserver sql migrate fail: %w", err)
	}
	return &repository{db: db}, nil
}

func (r *repository) Save(ctx context.Context, do *domain.NotebookServer) error {
	po := newNotebookServer(do)
	if err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(po).Error; err != nil {
		return apperrors.NewInternalError(err)
	}
	return nil
}

func (r *repository) Get(ctx context.Context, id string) (*domain.NotebookServer, error) {
	var po notebookServer
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, apperrors.NewInternalError(err)
	}
	return po.toDO(), nil
}

func (r *repository) Delete(ctx context.Context, do *domain.NotebookServer) error {
	po := newNotebookServer(do)
	if err := r.db.WithContext(ctx).Delete(po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return apperrors.NewInternalError(err)
	}
	return nil
}
