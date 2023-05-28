package sql

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	applog "github.com/Bio-OS/bioos/pkg/log"

	"github.com/Bio-OS/bioos/internal/context/submission/domain/submission"
)

type submissionRepository struct {
	db *gorm.DB
}

// NewSubmissionRepository ...
func NewSubmissionRepository(ctx context.Context, db *gorm.DB) (submission.Repository, error) {
	return &submissionRepository{db: db}, nil
}

func (s *submissionRepository) Save(ctx context.Context, sub *submission.Submission) error {
	po, err := SubmissionDOToSubmissionPO(ctx, sub)
	if err != nil {
		applog.Errorw("failed to convert submission do to po", "err", err)
		return apperrors.NewInternalError(err)
	}
	if err = s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(po).Error; err != nil {
		applog.Errorw("failed to create submission", "err", err)
		return apperrors.NewInternalError(err)
	}
	return nil
}

func (s *submissionRepository) Get(ctx context.Context, id string) (*submission.Submission, error) {
	var sub *Submission
	if err := s.db.WithContext(ctx).Where("id = ?", id).First(&sub).Error; err != nil {
		// applog.Errorw("failed to get submission", "err", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError("submission", id)
		}
		return nil, apperrors.NewInternalError(err)
	}
	return SubmissionPOToSubmissionDO(ctx, sub)
}

func (s *submissionRepository) Delete(ctx context.Context, sub *submission.Submission) error {
	if err := s.db.WithContext(ctx).Model(&Submission{}).Unscoped().Where("id = ?", sub.ID).Delete(sub).Error; err != nil {
		applog.Errorw("failed to delete submission", "err", err)
		return apperrors.NewInternalError(err)
	}
	return nil
}

func (s *submissionRepository) SoftDelete(ctx context.Context, sub *submission.Submission) error {
	if err := s.db.WithContext(ctx).Model(&SubmissionModel{}).Where("id = ?", sub.ID).Delete(&SubmissionModel{}).Error; err != nil {
		applog.Errorw("failed to soft delete submission", "err", err)
		return apperrors.NewInternalError(err)
	}
	return nil
}
