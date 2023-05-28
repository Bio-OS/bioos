package sql

import (
	"context"
	"errors"
	"reflect"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	applog "github.com/Bio-OS/bioos/pkg/log"

	datamodel "github.com/Bio-OS/bioos/internal/context/workspace/domain/data-model"
	"github.com/Bio-OS/bioos/pkg/consts"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	"github.com/Bio-OS/bioos/pkg/utils"
)

type dataModelRepository struct {
	db *gorm.DB
}

// NewDataModelRepository ...
func NewDataModelRepository(ctx context.Context, db *gorm.DB) (datamodel.Repository, error) {
	return &dataModelRepository{db: db}, nil
}

func (d *dataModelRepository) Save(ctx context.Context, dm *datamodel.DataModel) error {
	dataModel := DataModelDOtoDataModelPO(ctx, dm)
	switch dm.Type {
	case consts.DataModelTypeEntity:
		return d.saveEntityTypeDataModel(ctx, dataModel, dm)
	case consts.DataModelTypeEntitySet:
		return d.saveEntitySetTypeDataModel(ctx, dataModel, dm)
	case consts.DataModelTypeWorkspace:
		return d.saveWorkspaceTypeDataModel(ctx, dataModel, dm)
	default:
		return apperrors.NewInvalidError("unsupported data model type")
	}
}

func (d *dataModelRepository) Get(ctx context.Context, id string) (*datamodel.DataModel, error) {
	var dm DataModel
	if err := d.db.WithContext(ctx).Where("id = ?", id).First(&dm).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError("data model", err.Error())
		}
		applog.Errorw("failed to get data model", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	return DataModelPOToDataModelDO(ctx, &dm), nil
}

func (d *dataModelRepository) Delete(ctx context.Context, dm *datamodel.DataModel) error {
	dataModel := DataModelDOtoDataModelPO(ctx, dm)
	switch dm.Type {
	case consts.DataModelTypeEntity:
		var eh []*EntityHeader
		if err := d.db.WithContext(ctx).Where("name IN ? AND data_model_id = ?", dm.Headers, dm.ID).Order(ordersToOrderDB([]utils.Order{{
			Field:     "column_index",
			Ascending: true,
		}})).Find(&eh).Error; err != nil {
			applog.Errorw("failed to list entity data model headers", "err", err)
			return apperrors.NewInternalError(err)
		}
		index := make([]int, 0)
		for _, header := range eh {
			index = append(index, header.ColumnIndex)
		}
		if err := d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if dm.Headers == nil && dm.RowIDs == nil {
				if err := d.db.WithContext(ctx).Delete(&dataModel).Error; err != nil {
					applog.Errorw("failed to delete entity data model", "err", err)
					return apperrors.NewInternalError(err)
				}
				if err := d.db.WithContext(ctx).Where("data_model_id = ?", dm.ID).Delete(&EntityHeader{}).Error; err != nil {
					applog.Errorw("failed to delete entity data model headers", "err", err)
					return apperrors.NewInternalError(err)
				}
				if err := d.db.WithContext(ctx).Where("data_model_id = ?", dm.ID).Delete(&EntityGrid{}).Error; err != nil {
					applog.Errorw("failed to delete entity data model grids", "err", err)
					return apperrors.NewInternalError(err)
				}
				return d.deleteEntitySetWhenEntityDeleted(ctx, dm)
			}
			if len(dm.Headers) != 0 {
				if err := d.db.WithContext(ctx).Where("data_model_id = ?", dm.ID).Delete(&EntityHeader{}, "name NOT IN ?", dm.Headers).Error; err != nil {
					applog.Errorw("failed to delete entity data model headers", "err", err)
					return apperrors.NewInternalError(err)
				}
				if err := d.db.WithContext(ctx).Where("data_model_id = ?", dm.ID).Delete(&EntityGrid{}, "column_index NOT IN ?", index).Error; err != nil {
					applog.Errorw("failed to delete entity data model grids", "err", err)
					return apperrors.NewInternalError(err)
				}
			}
			if len(dm.RowIDs) != 0 {
				if err := d.db.WithContext(ctx).Where("data_model_id = ?", dm.ID).Delete(&EntityGrid{}, "row_id NOT IN ?", dm.RowIDs).Error; err != nil {
					applog.Errorw("failed to delete entity data model grids", "err", err)
					return apperrors.NewInternalError(err)
				}
				if err := d.updateEntitySetWhenEntityRowDeleted(ctx, dm); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return err
		}
	case consts.DataModelTypeEntitySet:
		if err := d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if dm.RowIDs == nil {
				if err := d.db.WithContext(ctx).Delete(&dataModel).Error; err != nil {
					applog.Errorw("failed to delete entity_set data model", "err", err)
					return apperrors.NewInternalError(err)
				}
				if err := d.db.WithContext(ctx).Where("data_model_id = ?", dm.ID).Delete(&EntitySetRow{}).Error; err != nil {
					applog.Errorw("failed to delete entity_set data model rows", "err", err)
					return apperrors.NewInternalError(err)
				}
				return d.deleteEntitySetWhenEntityDeleted(ctx, dm)
			} else if len(dm.RowIDs) != 0 {
				if err := d.db.WithContext(ctx).Where("data_model_id = ?", dm.ID).Delete(&EntitySetRow{}, "row_id NOT IN ?", dm.RowIDs).Error; err != nil {
					applog.Errorw("failed to delete entity_set data model rows", "err", err)
					return apperrors.NewInternalError(err)
				}
				if err := d.updateEntitySetWhenEntityRowDeleted(ctx, dm); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return err
		}
	case consts.DataModelTypeWorkspace:
		if err := d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if dm.RowIDs == nil {
				if err := d.db.WithContext(ctx).Delete(&dataModel).Error; err != nil {
					applog.Errorw("failed to delete workspace data model", "err", err)
					return apperrors.NewInternalError(err)
				}
				if err := d.db.WithContext(ctx).Where("data_model_id = ?", dm.ID).Delete(&WorkspaceRow{}).Error; err != nil {
					applog.Errorw("failed to delete workspace data model rows", "err", err)
					return apperrors.NewInternalError(err)
				}
			} else if len(dm.RowIDs) != 0 {
				if err := d.db.WithContext(ctx).Where("data_model_id = ?", dm.ID).Delete(&WorkspaceRow{}, "`key` NOT IN ?", dm.RowIDs).Error; err != nil {
					applog.Errorw("failed to delete workspace data model rows", "err", err)
					return apperrors.NewInternalError(err)
				}
			}
			return nil
		}); err != nil {
			return err
		}
	default:
		return apperrors.NewInvalidError("unsupported data model type")
	}
	return nil
}

func (d *dataModelRepository) updateEntitySetWhenEntityRowDeleted(ctx context.Context, dm *datamodel.DataModel) error {
	var entitySet DataModel
	if err := d.db.WithContext(ctx).Where("workspace_id = ? AND name = ?", dm.WorkspaceID, utils.GenDataModelEntitySetName(dm.Name)).First(&entitySet).Error; err != nil {
		applog.Errorw("failed to get data model by name", "err", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return apperrors.NewInternalError(err)
	}
	if err := d.db.WithContext(ctx).Where("data_model_id = ?", entitySet.ID).Delete(&EntitySetRow{}, "ref_row_id NOT IN ?", dm.RowIDs).Error; err != nil {
		applog.Errorw("failed to delete entity set data model row", "err", err)
		return apperrors.NewInternalError(err)
	}
	return nil
}

func (d *dataModelRepository) deleteEntitySetWhenEntityDeleted(ctx context.Context, dm *datamodel.DataModel) error {
	var entitySet DataModel
	if err := d.db.WithContext(ctx).Where("workspace_id = ? AND name = ?", dm.WorkspaceID, utils.GenDataModelEntitySetName(dm.Name)).First(&entitySet).Error; err != nil {
		applog.Errorw("failed to get data model by name", "err", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return apperrors.NewInternalError(err)
	}
	entitySetDO := DataModelPOToDataModelDO(ctx, &entitySet)
	if err := d.Delete(ctx, entitySetDO); err != nil {
		return err
	}
	return d.deleteEntitySetWhenEntityDeleted(ctx, entitySetDO)
}

func (d *dataModelRepository) saveEntityTypeDataModel(ctx context.Context, dataModel *DataModel, dm *datamodel.DataModel) error {
	entityHeaders := DataModelDOtoEntityHeadersPO(ctx, dm)
	grids := DataModelDOtoEntityGridsPO(ctx, dm)
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}, {Name: "workspace_id"}, {Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"updated_at"}),
		}).Create(dataModel).Error; err != nil {
			applog.Errorw("failed to create entity data model", "err", err)
			return apperrors.NewInternalError(err)
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}, {Name: "data_model_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "type"}),
		}).Create(&entityHeaders).Error; err != nil {
			applog.Errorw("failed to create entity data model headers", "err", err)
			return apperrors.NewInternalError(err)
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "data_model_id"}, {Name: "row_id"}, {Name: "column_index"}},
			DoUpdates: clause.AssignmentColumns([]string{"value"}),
		}).Create(&grids).Error; err != nil {
			applog.Errorw("failed to create entity data model grids", "err", err)
			return apperrors.NewInternalError(err)
		}
		return nil
	})
}

func (d *dataModelRepository) saveEntitySetTypeDataModel(ctx context.Context, dataModel *DataModel, dm *datamodel.DataModel) error {
	rows, err := DataModelDOtoEntitySetRowsPO(ctx, dm)
	if err != nil {
		return apperrors.NewInvalidError(err.Error())
	}
	if !reflect.DeepEqual(dm.Headers, []string{utils.GenDataModelHeaderOfID(dm.Name), strings.TrimSuffix(dm.Name, consts.DataModelEntitySetNameSuffix)}) {
		return apperrors.NewInvalidError("should not update entity set type data model's headers")
	}
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err = tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}, {Name: "workspace_id"}, {Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"updated_at"}),
		}).Create(dataModel).Error; err != nil {
			applog.Errorw("failed to create entity_set data model", "err", err)
			return apperrors.NewInternalError(err)
		}
		if err = tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}, {Name: "data_model_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"ref_row_id"}),
		}).Create(&rows).Error; err != nil {
			applog.Errorw("failed to create entity_set data model rows", "err", err)
			return apperrors.NewInternalError(err)
		}
		return nil
	})
}

func (d *dataModelRepository) saveWorkspaceTypeDataModel(ctx context.Context, dataModel *DataModel, dm *datamodel.DataModel) error {
	rows := DataModelDOtoWorkspaceRowsPO(ctx, dm)
	if !reflect.DeepEqual(dm.Headers, []string{consts.WorkspaceTypeDataModelHeaderKey, consts.WorkspaceTypeDataModelHeaderValue}) {
		return apperrors.NewInvalidError("should not update workspace type data model's headers")
	}
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}, {Name: "workspace_id"}, {Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"updated_at"}),
		}).Create(dataModel).Error; err != nil {
			applog.Errorw("failed to create workspace data model", "err", err)
			return apperrors.NewInternalError(err)
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "data_model_id"}, {Name: "key"}},
			DoUpdates: clause.AssignmentColumns([]string{"value"}),
		}).Create(&rows).Error; err != nil {
			applog.Errorw("failed to create workspace data model rows", "err", err)
			return apperrors.NewInternalError(err)
		}
		return nil
	})
}
