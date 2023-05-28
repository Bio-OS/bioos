package sql

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	"github.com/Bio-OS/bioos/pkg/consts"
	applog "github.com/Bio-OS/bioos/pkg/log"

	query "github.com/Bio-OS/bioos/internal/context/workspace/application/query/data-model"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	"github.com/Bio-OS/bioos/pkg/utils"
)

type dataModelReadModel struct {
	db *gorm.DB
}

// NewDataModelReadModel ...
func NewDataModelReadModel(ctx context.Context, db *gorm.DB) (query.DataModelReadModel, error) {
	if err := db.WithContext(ctx).AutoMigrate(&DataModel{}, &EntityHeader{}, &EntityGrid{}, &EntitySetRow{}, &WorkspaceRow{}); err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	return &dataModelReadModel{db: db}, nil
}

func (d *dataModelReadModel) ListDataModels(ctx context.Context, workspaceID string, filter *query.ListDataModelsFilter) ([]*query.DataModel, error) {
	db := d.db.WithContext(ctx).Where("workspace_id = ?", workspaceID).Order(ordersToOrderDB([]utils.Order{{
		Field:     "name",
		Ascending: true,
	}}))
	db = listDataModelsFilter(db, filter)
	// ignore _${data_model}, which is the data model name stored for submission
	// db = db.Where("name NOT REGEXP '^[_*]'")
	var ws []*DataModel
	if err := db.Find(&ws).Error; err != nil {
		applog.Errorw("failed to list data model", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	ret := make([]*query.DataModel, len(ws))
	for index, po := range ws {
		count, err := d.CountDataModelRows(ctx, po.ID, po.Type, &query.ListDataModelRowsFilter{})
		if err != nil {
			applog.Errorw("failed to count data model row", "err", err)
			return nil, apperrors.NewInternalError(err)
		}
		ret[index] = DataModelPOToDataModelDTO(ctx, po, count)
	}
	return ret, nil
}

func (d *dataModelReadModel) CountDataModel(ctx context.Context, workspaceID string, filter *query.ListDataModelsFilter) (int64, error) {
	db := d.db.WithContext(ctx).Where("workspace_id = ?", workspaceID).Order(ordersToOrderDB([]utils.Order{{
		Field:     "name",
		Ascending: true,
	}}))
	db = listDataModelsFilter(db, filter)
	var count int64
	if err := db.Count(&count).Error; err != nil {
		applog.Errorw("failed to count data model", "err", err)
		return 0, apperrors.NewInternalError(err)
	}
	return count, nil
}

func (d *dataModelReadModel) GetDataModelName(ctx context.Context, workspaceID, id string) (string, error) {
	var name string
	if err := d.db.WithContext(ctx).Model(&DataModel{}).Where("workspace_id = ? AND id = ?", workspaceID, id).Select("name").Find(&name).Error; err != nil {
		applog.Errorw("failed to get data model name", "err", err)
		return "", apperrors.NewInternalError(err)
	}
	if len(name) == 0 {
		return "", apperrors.NewNotFoundError("data model", id)
	}
	return name, nil
}

func (d *dataModelReadModel) GetDataModelWithID(ctx context.Context, id string) (*query.DataModel, error) {
	var dm DataModel
	if err := d.db.WithContext(ctx).Where("id = ?", id).First(&dm).Error; err != nil {
		applog.Errorw("failed to get data model by id", "err", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError("data model", id)
		}
		return nil, apperrors.NewInternalError(err)
	}
	count, err := d.CountDataModelRows(ctx, dm.ID, dm.Type, &query.ListDataModelRowsFilter{})
	if err != nil {
		applog.Errorw("failed to count data model rows", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	return DataModelPOToDataModelDTO(ctx, &dm, count), nil
}

func (d *dataModelReadModel) GetDataModelWithName(ctx context.Context, workspaceID, name string) (*query.DataModel, error) {
	var dm DataModel
	if err := d.db.WithContext(ctx).Where("workspace_id = ? AND name = ?", workspaceID, name).First(&dm).Error; err != nil {
		applog.Errorw("failed to get data model by name", "err", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError("data model", name)
		}
		return nil, apperrors.NewInternalError(err)
	}
	count, err := d.CountDataModelRows(ctx, dm.ID, dm.Type, &query.ListDataModelRowsFilter{})
	if err != nil {
		applog.Errorw("failed to count data model rows", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	return DataModelPOToDataModelDTO(ctx, &dm, count), nil
}

func (d *dataModelReadModel) ListDataModelHeaders(ctx context.Context, id, name, _type string) ([]string, error) {
	switch _type {
	case consts.DataModelTypeEntitySet:
		return []string{utils.GenDataModelHeaderOfID(name), strings.TrimSuffix(name, consts.DataModelEntitySetNameSuffix)}, nil
	case consts.DataModelTypeWorkspace:
		return []string{consts.WorkspaceTypeDataModelHeaderKey, consts.WorkspaceTypeDataModelHeaderValue}, nil
	case consts.DataModelTypeEntity:
		return d.ListEntityDataModelHeaders(ctx, id)
	default:
		return nil, apperrors.NewInvalidError("unsupport data model type")
	}
}

func (d *dataModelReadModel) ListEntityDataModelHeaders(ctx context.Context, id string) ([]string, error) {
	ws, err := d.listEntityDataModelHeaders(ctx, id)
	if err != nil {
		return nil, err
	}
	ret := make([]string, len(ws))
	for index, po := range ws {
		ret[index] = EntityHeadersPOToHeadersDTO(ctx, po)
	}
	return ret, nil
}

func (d *dataModelReadModel) listEntityDataModelHeaders(ctx context.Context, id string) ([]*EntityHeader, error) {
	db := d.db.WithContext(ctx).Where("data_model_id = ?", id).Order(ordersToOrderDB([]utils.Order{{
		Field:     "column_index",
		Ascending: true,
	}}))
	var eh []*EntityHeader
	if err := db.Find(&eh).Error; err != nil {
		applog.Errorw("failed to list data model entity headers", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	return eh, nil
}

func (d *dataModelReadModel) ListEntityDataModelColumnsWithRowIDs(ctx context.Context, id string, headers []string, rowIDs []string) (map[string][]string, error) {
	db := d.db.WithContext(ctx).Where("data_model_id = ?", id).Order(ordersToOrderDB([]utils.Order{{
		Field:     "column_index",
		Ascending: true,
	}}))
	if len(headers) > 0 {
		db = db.Where("name IN ?", headers)
	}
	var eh []*EntityHeader
	if err := db.Find(&eh).Error; err != nil {
		applog.Errorw("failed to list data model entity headers", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	res := make(map[string][]string)
	for _, entityHeader := range eh {
		db := d.db.WithContext(ctx).Where("data_model_id = ? AND column_index = ?", id, entityHeader.ColumnIndex).Order(ordersToOrderDB([]utils.Order{{
			Field:     "row_id",
			Ascending: true,
		}}))
		if len(rowIDs) > 0 {
			db = db.Where("row_id IN ?", rowIDs)
		}
		var eg []*EntityGrid
		if err := db.Find(&eg).Error; err != nil {
			applog.Errorw("failed to list data model entity grids", "err", err)
			return nil, apperrors.NewInternalError(err)
		}
		res[entityHeader.Name] = EntityGridsPOToColumnDTO(ctx, eg, rowIDs)
	}
	return res, nil
}

func (d *dataModelReadModel) ListDataModelRows(ctx context.Context, id, _type string, pagination *utils.Pagination, order *utils.Order, filter *query.ListDataModelRowsFilter) ([][]string, int64, error) {
	switch _type {
	case consts.DataModelTypeEntity:
		rows, err := d.listEntityDataModelRows(ctx, id, pagination, order, filter)
		if err != nil {
			return nil, 0, err
		}
		count, err := d.countEntityDataModelRows(ctx, id, filter)
		if err != nil {
			return nil, 0, err
		}
		return rows, count, nil
	case consts.DataModelTypeEntitySet:
		rows, err := d.listEntitySetDataModelRows(ctx, id, pagination, order, filter)
		if err != nil {
			return nil, 0, err
		}
		count, err := d.countEntitySetDataModelRows(ctx, id, filter)
		if err != nil {
			return nil, 0, err
		}
		return rows, count, nil
	case consts.DataModelTypeWorkspace:
		rows, err := d.listWorkspaceDataModelRows(ctx, id, pagination, filter)
		if err != nil {
			return nil, 0, err
		}
		count, err := d.countWorkspaceDataModelRows(ctx, id, filter)
		if err != nil {
			return nil, 0, err
		}
		return rows, count, nil
	default:
		return nil, 0, apperrors.NewInvalidError("unsupported data model type")
	}
}

func (d *dataModelReadModel) getEntityDataModelRowIDsWithFilter(ctx context.Context, id string, filter *query.ListDataModelRowsFilter) ([]string, error) {
	db := d.db.WithContext(ctx).Where("data_model_id = ?", id)
	db = listEntityDataModelRowsFilter(db, filter)
	var egs []*EntityGrid
	if err := db.Select("row_id").Find(&egs).Error; err != nil {
		applog.Errorw("failed to list entity data model headers", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	return EntityGridsPOToRowIDsDTO(ctx, egs), nil
}

func (d *dataModelReadModel) getEntityDataModelRowIDsWithPagination(ctx context.Context, id string, rowIdsWithFilter []string, pagination *utils.Pagination, order *utils.Order) ([]string, error) {
	db := d.db.WithContext(ctx).Where("data_model_id = ? AND column_index = ? AND row_id IN ?", id, order.Field, rowIdsWithFilter).Limit(pagination.GetLimit()).Offset(pagination.GetOffset()).Order(ordersToOrderDB([]utils.Order{{
		Field:     "`value`",
		Ascending: order.Ascending,
	}, {
		Field:     "row_id",
		Ascending: true,
	}}))
	var egs []*EntityGrid
	if err := db.Select("row_id").Find(&egs).Error; err != nil {
		applog.Errorw("failed to list entity data model row ids", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	// this row ids([]string)'s order should be equal to []*EntityGrid
	var ids []string
	for _, eg := range egs {
		ids = append(ids, eg.RowID)
	}
	return ids, nil
}

func (d *dataModelReadModel) listEntityDataModelRows(ctx context.Context, id string, pagination *utils.Pagination, order *utils.Order, filter *query.ListDataModelRowsFilter) ([][]string, error) {
	rowIdsWithFilter, err := d.getEntityDataModelRowIDsWithFilter(ctx, id, filter)
	if err != nil {
		applog.Errorw("failed to list entity data model ids", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	if len(rowIdsWithFilter) == 0 {
		return nil, nil
	}
	rowIDs, err := d.getEntityDataModelRowIDsWithPagination(ctx, id, rowIdsWithFilter, pagination, order)
	if err != nil {
		applog.Errorw("failed to list entity data model ids", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	rows := make([][]string, 0)
	var eg []*EntityGrid
	if err := d.db.WithContext(ctx).Where("data_model_id = ? AND row_id IN  ?", id, rowIDs).Order(ordersToOrderDB([]utils.Order{{
		Field:     "row_id",
		Ascending: true,
	}, {
		Field:     "column_index",
		Ascending: true,
	}})).Find(&eg).Error; err != nil {
		applog.Errorw("failed to list entity data model grids", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	mappedRows := make(map[string][]string)
	for _, grid := range eg {
		row, ok := mappedRows[grid.RowID]
		if !ok {
			mappedRows[grid.RowID] = make([]string, 0)
		}
		mappedRows[grid.RowID] = append(row, grid.Value)
	}
	for _, rowId := range rowIDs {
		rows = append(rows, mappedRows[rowId])
	}
	return rows, nil
}

func (d *dataModelReadModel) getEntitySetDataModelRowIDsWithFilter(ctx context.Context, id string, filter *query.ListDataModelRowsFilter) ([]string, error) {
	db := d.db.WithContext(ctx).Where("data_model_id = ?", id)
	db = listEntitySetDataModelRowsFilter(db, filter)
	var eg []*EntitySetRow
	if err := db.Find(&eg).Error; err != nil {
		applog.Errorw("failed to list entity_set data model rows", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	return EntitySetRowsPOToEntitySetRowIDsDTO(ctx, eg), nil
}

func (d *dataModelReadModel) getEntitySetDataModelRowIDsWithPagination(ctx context.Context, id string, rowIdsWithFilter []string, pagination *utils.Pagination, order *utils.Order) ([]string, error) {
	db := d.db.WithContext(ctx).Where("data_model_id = ? AND row_id IN ?", id, rowIdsWithFilter).Limit(pagination.GetLimit()).Offset(pagination.GetOffset()).Order(ordersToOrderDB([]utils.Order{*order}))
	var egs []*EntitySetRow
	if err := db.Distinct("row_id").Find(&egs).Error; err != nil {
		applog.Errorw("failed to list entity data model row ids", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	// this row ids([]string)'s order should be equal to []*EntityGrid
	var ids []string
	for _, eg := range egs {
		ids = append(ids, eg.RowID)
	}
	return ids, nil
}

func (d *dataModelReadModel) listEntitySetDataModelRows(ctx context.Context, id string, pagination *utils.Pagination, order *utils.Order, filter *query.ListDataModelRowsFilter) ([][]string, error) {
	rowIdsWithFilter, err := d.getEntitySetDataModelRowIDsWithFilter(ctx, id, filter)
	if err != nil {
		applog.Errorw("failed to list entity_set data model row ids", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	if len(rowIdsWithFilter) == 0 {
		return nil, nil
	}
	rowIDs, err := d.getEntitySetDataModelRowIDsWithPagination(ctx, id, rowIdsWithFilter, pagination, order)
	if err != nil {
		applog.Errorw("failed to list entity data model ids", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	db := d.db.WithContext(ctx).Where("data_model_id = ? AND row_id IN  ?", id, rowIDs)
	var eg []*EntitySetRow
	if err := db.Find(&eg).Error; err != nil {
		applog.Errorw("failed to list entity_set data model rows", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	return EntitySetRowsPOToEntitySetRowsDTO(ctx, eg, rowIDs), nil
}

func (d *dataModelReadModel) listWorkspaceDataModelRows(ctx context.Context, id string, pagination *utils.Pagination, filter *query.ListDataModelRowsFilter) ([][]string, error) {
	db := d.db.WithContext(ctx).Where("data_model_id = ?", id).Limit(pagination.GetLimit()).Offset(pagination.GetOffset()).Order(ordersToOrderDB(pagination.Orders))
	db = listWorkspaceDataModelRowsFilter(db, filter)
	var eg []*WorkspaceRow
	if err := db.Find(&eg).Error; err != nil {
		applog.Errorw("failed to list workspace data model rows", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	return WorkspaceRowsPOToWorkspaceRowsDTO(ctx, eg), nil
}

func (d *dataModelReadModel) CountDataModelRows(ctx context.Context, id, _type string, filter *query.ListDataModelRowsFilter) (int64, error) {
	switch _type {
	case consts.DataModelTypeEntity:
		return d.countEntityDataModelRows(ctx, id, filter)
	case consts.DataModelTypeEntitySet:
		return d.countEntitySetDataModelRows(ctx, id, filter)
	case consts.DataModelTypeWorkspace:
		return d.countWorkspaceDataModelRows(ctx, id, filter)
	default:
		return 0, apperrors.NewInvalidError("unsupported data model type")
	}
}

func (d *dataModelReadModel) ListAllDataModelRowIDs(ctx context.Context, id, _type string) ([]string, error) {
	switch _type {
	case consts.DataModelTypeEntity:
		return d.listAllEntityDataModelRowIDs(ctx, id)
	case consts.DataModelTypeEntitySet:
		return d.listAllEntitySetDataModelRowIDs(ctx, id)
	case consts.DataModelTypeWorkspace:
		return d.listAllWorkspaceDataModelRowIDs(ctx, id)
	default:
		return nil, apperrors.NewInvalidError("unsupported data model type ")
	}
}

func (d *dataModelReadModel) listAllEntityDataModelRowIDs(ctx context.Context, id string) ([]string, error) {
	db := d.db.WithContext(ctx).Where("data_model_id = ?", id).Distinct("row_id")
	var eg []*EntityGrid
	if err := db.Find(&eg).Error; err != nil {
		applog.Errorw("failed to list all entity data model row ids", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	return EntityGridsPOToRowIDsDTO(ctx, eg), nil
}

func (d *dataModelReadModel) listAllEntitySetDataModelRowIDs(ctx context.Context, id string) ([]string, error) {
	db := d.db.WithContext(ctx).Where("data_model_id = ?", id).Distinct("row_id")
	var eg []*EntitySetRow
	if err := db.Find(&eg).Error; err != nil {
		applog.Errorw("failed to list all entity_set data model row ids", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	return EntitySetRowsPOToRowIDsDTO(ctx, eg), nil
}

func (d *dataModelReadModel) listAllWorkspaceDataModelRowIDs(ctx context.Context, id string) ([]string, error) {
	db := d.db.WithContext(ctx).Where("data_model_id = ?", id).Distinct("key")
	var eg []*WorkspaceRow
	if err := db.Find(&eg).Error; err != nil {
		applog.Errorw("failed to list all workspace data model row ids", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	return WorkspaceRowsPOToRowIDsDTO(ctx, eg), nil
}

func (d *dataModelReadModel) countEntityDataModelRows(ctx context.Context, id string, filter *query.ListDataModelRowsFilter) (int64, error) {
	db := d.db.WithContext(ctx).Model(EntityGrid{}).Where("data_model_id = ?", id).Distinct("row_id")
	db = listEntityDataModelRowsFilter(db, filter)
	var count int64
	if err := db.Count(&count).Error; err != nil {
		applog.Errorw("failed to count entity data model rows", "err", err)
		return 0, apperrors.NewInternalError(err)
	}
	return count, nil
}

func (d *dataModelReadModel) countEntitySetDataModelRows(ctx context.Context, id string, filter *query.ListDataModelRowsFilter) (int64, error) {
	db := d.db.WithContext(ctx).Model(EntitySetRow{}).Where("data_model_id = ?", id).Distinct("row_id")
	db = listEntitySetDataModelRowsFilter(db, filter)
	var count int64
	if err := db.Count(&count).Error; err != nil {
		applog.Errorw("failed to count entity_set data model rows", "err", err)
		return count, apperrors.NewInternalError(err)
	}
	return count, nil
}

func (d *dataModelReadModel) countWorkspaceDataModelRows(ctx context.Context, id string, filter *query.ListDataModelRowsFilter) (int64, error) {
	db := d.db.WithContext(ctx).Model(WorkspaceRow{}).Where("data_model_id = ?", id).Distinct("key")
	db = listWorkspaceDataModelRowsFilter(db, filter)
	var count int64
	if err := db.Count(&count).Error; err != nil {
		applog.Errorw("failed to count workspace data model rows", "err", err)
		return count, apperrors.NewInternalError(err)
	}
	return count, nil
}

func ordersToOrderDB(orders []utils.Order) string {
	orderStrs := make([]string, 0, len(orders))
	for _, order := range orders {
		orderStr := order.Field
		if order.Ascending {
			orderStr += " ASC"
		} else {
			orderStr += " DESC"
		}
		orderStrs = append(orderStrs, orderStr)
	}
	return strings.Join(orderStrs, ", ")
}

func listDataModelsFilter(db *gorm.DB, filter *query.ListDataModelsFilter) *gorm.DB {
	if filter == nil {
		return db
	}
	if len(filter.SearchWord) > 0 {
		db = utils.SearchWordFilter(db, filter.SearchWord, []string{"name"}, false)
	}
	if len(filter.IDs) > 0 {
		db = db.Where("id IN ?", filter.IDs)
	}
	if len(filter.Types) > 0 {
		db = db.Where("type IN ?", filter.Types)
	}
	return db
}

func listEntityDataModelRowsFilter(db *gorm.DB, filter *query.ListDataModelRowsFilter) *gorm.DB {
	if filter == nil {
		return db
	}
	if len(filter.SearchWord) > 0 {
		db = utils.SearchWordFilter(db, filter.SearchWord, []string{"value"}, false)
	}
	if len(filter.RowIDs) > 0 {
		db = db.Where("row_id IN ?", filter.RowIDs)
	}
	return db
}

func listEntitySetDataModelRowsFilter(db *gorm.DB, filter *query.ListDataModelRowsFilter) *gorm.DB {
	if filter == nil {
		return db
	}
	if len(filter.SearchWord) > 0 {
		db = utils.SearchWordFilter(db, filter.SearchWord, []string{"ref_row_id", "row_id"}, false)
	}
	if len(filter.InSetIDs) > 0 {
		db = db.Where("ref_row_id IN ?", filter.InSetIDs)
	}
	if len(filter.RowIDs) > 0 {
		db = db.Where("row_id IN ?", filter.RowIDs)
	}
	return db
}

func listWorkspaceDataModelRowsFilter(db *gorm.DB, filter *query.ListDataModelRowsFilter) *gorm.DB {
	if filter == nil {
		return db
	}
	if len(filter.SearchWord) > 0 {
		db = utils.SearchWordFilter(db, filter.SearchWord, []string{"`key`", "`value`"}, false)
	}
	if len(filter.RowIDs) > 0 {
		db = db.Where("`key` IN ?", filter.RowIDs)
	}
	return db
}
