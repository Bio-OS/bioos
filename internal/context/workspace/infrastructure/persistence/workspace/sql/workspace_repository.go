package mysql

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/Bio-OS/bioos/internal/context/workspace/domain/workspace"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	applog "github.com/Bio-OS/bioos/pkg/log"
)

type workspaceRepository struct {
	db *gorm.DB
}

// NewWorkspaceRepository ...
func NewWorkspaceRepository(ctx context.Context, db *gorm.DB) (workspace.Repository, error) {
	if err := db.WithContext(ctx).AutoMigrate(&Workspace{}); err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	return &workspaceRepository{db: db}, nil
}

func (r *workspaceRepository) Get(ctx context.Context, id string) (*workspace.Workspace, error) {
	var ws Workspace
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&ws).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError("workspace", "id")
		}
		applog.Errorw("failed to get workspace", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	return WorkspacePOToWorkspaceDO(ctx, &ws)
}

func (r *workspaceRepository) Save(ctx context.Context, w *workspace.Workspace) error {
	ws := &Workspace{}
	if err := r.db.WithContext(ctx).Where("id = ?", w.ID).First(&ws).Error; err != nil {
		// Create if workspace with same id not exist
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ws = WorkspaceDOtoWorkspacePO(w)
			if err := r.db.WithContext(ctx).Create(ws).Error; err != nil {
				applog.Errorw("failed to save workspace", "err", err)
				return apperrors.NewInternalError(err)
			}
			return nil
		}
		applog.Errorw("failed to get workspace", "err", err)
		return apperrors.NewInternalError(err)
	}
	// Update if workspace with same id exist
	ws = WorkspaceDOtoWorkspacePO(w)
	if err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(ws).Error; err != nil {
		applog.Errorw("failed to save workspace", "err", err)
		return apperrors.NewInternalError(err)
	}
	return nil
}

func (r *workspaceRepository) Delete(ctx context.Context, w *workspace.Workspace) error {
	ws := WorkspaceDOtoWorkspacePO(w)
	if err := r.db.WithContext(ctx).Delete(ws).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return apperrors.NewInternalError(err)
	}
	return nil
}
