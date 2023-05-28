package sql

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	applog "github.com/Bio-OS/bioos/pkg/log"

	"github.com/Bio-OS/bioos/internal/context/submission/domain/run"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
)

type runRepository struct {
	db *gorm.DB
}

// NewRunRepository ...
func NewRunRepository(ctx context.Context, db *gorm.DB) (run.Repository, error) {
	return &runRepository{db: db}, nil
}

func (r *runRepository) Save(ctx context.Context, w *run.Run) error {
	runPO := RunDOToRunPO(w)
	taskPOList := RunDOToTaskPOList(w)
	// only update run
	if len(taskPOList) == 0 {
		if err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}, {Name: "submission_id"}},
			UpdateAll: true,
		}).Create(&runPO).Error; err != nil {
			applog.Errorw("failed to create/update run", "err", err)
			return apperrors.NewInternalError(err)
		}
		return nil
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(&runPO).Error; err != nil {
			applog.Errorw("failed to create/update run", "err", err)
			return apperrors.NewInternalError(err)
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}, {Name: "run_id"}},
			UpdateAll: true,
		}).Create(&taskPOList).Error; err != nil {
			applog.Errorw("failed to create/update task", "err", err)
			return apperrors.NewInternalError(err)
		}
		return nil
	})
}

func (r *runRepository) Get(ctx context.Context, id string) (*run.Run, error) {
	var runPo Run
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&runPo).Error; err != nil {
		applog.Errorw("failed to get task", "err", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError("Run", id)
		}
		return nil, apperrors.NewInternalError(err)
	}
	return RunPOToRunDO(&runPo), nil
}

func (r *runRepository) Delete(ctx context.Context, w *run.Run) error {
	runPO := RunDOToRunPO(w)
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := r.db.WithContext(ctx).Delete(&runPO).Error; err != nil {
			applog.Errorw("failed to delete run", "err", err)
			return apperrors.NewInternalError(err)
		}
		if err := r.db.WithContext(ctx).Where("run_id = ?", runPO.ID).Delete(&Task{}).Error; err != nil {
			applog.Errorw("failed to delete task", "err", err)
			return apperrors.NewInternalError(err)
		}
		return nil
	})
}
