package sql

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/Bio-OS/bioos/internal/context/notebookserver/application/query"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
)

type readModel struct {
	db *gorm.DB
}

func NewReadModel(ctx context.Context, db *gorm.DB) (query.ReadModel, error) {
	return &readModel{db: db}, nil
}

func (r *readModel) ListSettingsByWorkspace(ctx context.Context, workspaceID string) ([]*query.NotebookSettings, error) {
	var po []notebookServer
	if err := r.db.WithContext(ctx).Where("workspace_id = ?", workspaceID).Find(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, apperrors.NewInternalError(err)
	}
	var res []*query.NotebookSettings
	for i := range po {
		res = append(res, po[i].toDTO())
	}
	return res, nil
}

func (r *readModel) GetSettingsByID(ctx context.Context, workspaceID, id string) (*query.NotebookSettings, error) {
	var po notebookServer
	if err := r.db.WithContext(ctx).Where("id = ?", id).Where("workspace_id = ?", workspaceID).First(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, apperrors.NewInternalError(err)
	}
	return po.toDTO(), nil
}
