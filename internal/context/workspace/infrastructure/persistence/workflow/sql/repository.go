package sql

import (
	"context"
	"errors"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	domain "github.com/Bio-OS/bioos/internal/context/workspace/domain/workflow"
	"github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	applog "github.com/Bio-OS/bioos/pkg/log"
)

const MySQLDuplicatedCode = 1062

func NewRepository(ctx context.Context, db *gorm.DB) (domain.Repository, error) {
	if err := db.WithContext(ctx).AutoMigrate(&workflow{}, &workflowVersion{}, &workflowFile{}); err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	return &repository{db: db}, nil
}

type repository struct {
	db *gorm.DB
}

var _ domain.Repository = &repository{}

func (r *repository) Save(ctx context.Context, wf *domain.Workflow) (err error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else if tx.Error != nil {
			applog.Errorw("transaction error, start to rollback", "err", tx.Error)
			tx.Rollback()
		} else if err != nil {
			applog.Errorw("save workflow error, start to rollback", "err", err)
			tx.Rollback()
		} else {
			applog.Infow("start to commit workflow", "workflow", wf)
			tx.Commit()
		}
	}()

	// save workflow
	workflowPO := workflowDOToPO(wf)
	createTX := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(workflowPO)

	if createTX.Error != nil {
		applog.Errorw("failed to save workflow", "err", createTX.Error)
		var mysqlErr *mysql.MySQLError
		// Error 1062 (23000): Duplicate entry
		if errors.As(createTX.Error, &mysqlErr) && mysqlErr.Number == MySQLDuplicatedCode {
			return proto.ErrorWorkflowNameDuplicated("workflow update failed name:%s already exists", workflowPO.Name)
		}

		return apperrors.NewInternalError(tx.Error)
	}

	// save workflow versions
	for _, workflowVersionDO := range wf.Versions {
		workflowVersionPO, err := workflowVersionDOToPO(wf.ID, workflowVersionDO)
		if err != nil {
			return apperrors.NewInternalError(err)
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(workflowVersionPO).Error; err != nil {
			applog.Errorw("failed to save workflow version", "err", err)
			return apperrors.NewInternalError(err)
		}

		// save workflow files
		for _, workflowFileDO := range workflowVersionDO.Files {
			workflowFilePO := workflowFileDOToPO(workflowVersionDO.ID, workflowFileDO)
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				UpdateAll: true,
			}).Create(workflowFilePO).Error; err != nil {
				applog.Errorw("failed to save workflow file", "err", err)
				return apperrors.NewInternalError(err)
			}
		}
	}

	return nil
}

func (r *repository) Get(ctx context.Context, workspaceID string, workflowID string) (*domain.Workflow, error) {
	workflowPO := workflow{
		ID:          workflowID,
		WorkspaceID: workspaceID,
	}
	if err := r.db.WithContext(ctx).First(&workflowPO).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, proto.ErrorWorkflowNotFound("workflow:%s in workspace:%s not found", workflowID, workspaceID)
		}
		applog.Errorw("failed to get workflow", "err", err)
		return nil, apperrors.NewInternalError(err)
	}

	// list workflow versions
	var workflowVersionPOs []*workflowVersion
	if err := r.db.WithContext(ctx).Where("workflow_id = ?", workflowID).Find(&workflowVersionPOs).Error; err != nil {
		applog.Errorw("failed to list workflow versions", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	workflowDO := workflowPO.toDO()
	for _, workflowVersionPO := range workflowVersionPOs {
		// list workflow files
		var workflowFilePOs []*workflowFile
		if err := r.db.WithContext(ctx).Where("workflow_version_id = ?", workflowVersionPO.ID).Find(&workflowFilePOs).Error; err != nil {
			applog.Errorw("failed to list workflow versions", "err", err)
			return nil, apperrors.NewInternalError(err)
		}
		workflowVersionDO, err := workflowVersionPO.toDO()
		if err != nil {
			applog.Errorw("failed to convert workflow version PO to DO", "err", err)
			return nil, apperrors.NewInternalError(err)
		}
		for _, workflowFilePO := range workflowFilePOs {
			workflowVersionDO.Files[workflowFilePO.ID] = workflowFilePO.toDO()
		}
		workflowDO.Versions[workflowVersionDO.ID] = workflowVersionDO
	}

	return workflowDO, nil
}

func (r *repository) Delete(ctx context.Context, wf *domain.Workflow) (err error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else if tx.Error != nil {
			applog.Errorw("delete workflow transaction error, start to rollback", "err", tx.Error)
			tx.Rollback()
		} else if err != nil {
			applog.Errorw("delete workflow error, start to rollback", "err", err)
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// delete WorkflowFiles
	var workflowVersionDOs []*workflowVersion
	if err := r.db.WithContext(ctx).Where("workflow_id = ?", wf.ID).Find(&workflowVersionDOs).Error; err != nil {
		applog.Errorw("failed to list workflow versions", "err", err)
		return apperrors.NewInternalError(err)
	}
	for _, workflowVersionDO := range workflowVersionDOs {
		if err := tx.Where("workflow_version_id = ?", workflowVersionDO.ID).Delete(&workflowFile{}).Error; err != nil {
			tx.Rollback()
			return apperrors.NewInternalError(err)
		}
	}

	// delete WorkflowVersion
	if err := tx.Where("workflow_id = ?", wf.ID).Delete(&workflowVersion{}).Error; err != nil {
		tx.Rollback()
		return apperrors.NewInternalError(err)
	}

	// delete Workflow
	if err := tx.Where("id = ?", wf.ID).Delete(&workflow{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return apperrors.NewInternalError(err)
	}

	return nil
}

func (r *repository) List(ctx context.Context, workspaceID string) ([]string, error) {
	var ids []struct {
		ID string
	}
	if err := r.db.WithContext(ctx).Table("workflow").Select("id").Where("workspace_id = ?", workspaceID).Find(&ids).Error; err != nil {
		applog.Errorw("failed to list workflow", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	workflowIDs := make([]string, len(ids))
	for i, id := range ids {
		workflowIDs[i] = id.ID
	}
	return workflowIDs, nil
}
