package sql

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"

	applog "github.com/Bio-OS/bioos/pkg/log"

	query "github.com/Bio-OS/bioos/internal/context/submission/application/query/run"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	"github.com/Bio-OS/bioos/pkg/utils"
)

type runReadModel struct {
	db *gorm.DB
}

// NewRunReadModel ...
func NewRunReadModel(ctx context.Context, db *gorm.DB) (query.ReadModel, error) {
	if err := db.WithContext(ctx).AutoMigrate(&Run{}, &Task{}); err != nil {
		return nil, apperrors.NewInternalError(err)
	}

	return &runReadModel{db: db}, nil
}

func (r *runReadModel) ListAllRunIDs(ctx context.Context, submissionID string) ([]string, error) {
	var ids []string
	if err := r.db.WithContext(ctx).Model(&Run{}).Select("id").Where("submission_id = ?", submissionID).Find(&ids).Error; err != nil {
		applog.Errorw("failed to list all run ids", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	return ids, nil
}

func (r *runReadModel) ListRuns(ctx context.Context, submissionID string, pg *utils.Pagination, filter *query.ListRunsFilter) ([]*query.RunItem, error) {
	dbChain := r.db.WithContext(ctx).Model(&Run{}).Where("submission_id = ?", submissionID).Limit(pg.GetLimit()).Offset(pg.GetOffset()).Order(ordersToOrderDB(pg.Orders))
	dbChain = listRunsFilter(dbChain, filter)
	var runs []*Run
	if err := dbChain.Find(&runs).Error; err != nil {
		applog.Errorw("failed to list runs", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	ret := make([]*query.RunItem, len(runs))
	for index, po := range runs {
		item, err := RunPOToRunDTO(ctx, po)
		if err != nil {
			applog.Errorw("failed to convert run po to dto", "err", err)
			return nil, apperrors.NewInternalError(err)
		}
		ret[index] = item
	}
	return ret, nil
}

func (r *runReadModel) CountRuns(ctx context.Context, submissionID string, filter *query.ListRunsFilter) (int, error) {
	dbChain := r.db.WithContext(ctx).Model(&Run{}).Where("submission_id = ?", submissionID)
	dbChain = listRunsFilter(dbChain, filter)
	var count int64
	if err := dbChain.Count(&count).Error; err != nil {
		applog.Errorw("failed to count runs", "err", err)
		return 0, apperrors.NewInternalError(err)
	}
	return int(count), nil
}

func (r *runReadModel) ListTasks(ctx context.Context, runID string, pg *utils.Pagination) ([]*query.TaskItem, error) {
	dbChain := r.db.WithContext(ctx).Model(&Task{}).Where("run_id = ?", runID).Limit(pg.GetLimit()).Offset(pg.GetOffset()).Order(ordersToOrderDB(pg.Orders))
	var tasks []*Task
	if err := dbChain.Find(&tasks).Error; err != nil {
		applog.Errorw("failed to list tasks", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	ret := make([]*query.TaskItem, len(tasks))
	for index, po := range tasks {
		item := TaskPOToTaskDTO(ctx, po)
		ret[index] = item
	}
	return ret, nil
}

func (r *runReadModel) CountTasks(ctx context.Context, runID string) (int, error) {
	dbChain := r.db.WithContext(ctx).Model(&Task{}).Where("run_id = ?", runID)
	var count int64
	if err := dbChain.Count(&count).Error; err != nil {
		applog.Errorw("failed to count tasks", "err", err)
		return 0, apperrors.NewInternalError(err)
	}
	return int(count), nil
}

func listRunsFilter(db *gorm.DB, filter *query.ListRunsFilter) *gorm.DB {
	if filter == nil {
		return db
	}
	if len(filter.SearchWord) != 0 {
		db = utils.SearchWordFilter(db, filter.SearchWord, []string{"name"}, filter.Exact)
	}
	if len(filter.IDs) != 0 {
		db = db.Where("id IN ?", filter.IDs)
	}
	if len(filter.Status) != 0 {
		db = db.Where("status IN ?", filter.Status)
	}

	return db
}

func (r *runReadModel) CountTasksResult(ctx context.Context, runID string) ([]*query.StatusCount, error) {
	dbChain := r.db.WithContext(ctx).Table("task").Where("run_id = ?", runID)
	statusCounts, err := countByStatus(dbChain)
	if err != nil {
		applog.Errorw("failed to count tasks result", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	ret := make([]*query.StatusCount, len(statusCounts))
	for index, statusCount := range statusCounts {
		ret[index] = StatusCountPOToStatusCountDTO(statusCount)
	}
	return ret, nil
}

func (r *runReadModel) CountRunsResult(ctx context.Context, submissionID string) ([]*query.StatusCount, error) {
	dbChain := r.db.WithContext(ctx).Table("run").Where("submission_id = ?", submissionID)
	statusCounts, err := countByStatus(dbChain)
	if err != nil {
		applog.Errorw("failed to count runs result", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	ret := make([]*query.StatusCount, len(statusCounts))
	for index, statusCount := range statusCounts {
		ret[index] = StatusCountPOToStatusCountDTO(statusCount)
	}
	return ret, nil
}

func countByStatus(db *gorm.DB) ([]*StatusCount, error) {
	var counts []*StatusCount
	if err := db.
		Select("count(*) as count, status").
		Group("status").Find(&counts).Error; err != nil {
		return nil, fmt.Errorf("failed to count group by status: %w", err)
	}
	return counts, nil
}

func ordersToOrderDB(orders []utils.Order) string {
	orderStrs := make([]string, 0, len(orders))
	for _, order := range orders {
		var orderStr string
		switch order.Field {
		case query.OrderByName:
			orderStr = "name"
		case query.OrderByStartTime:
			orderStr = "start_time"
		default:
			continue
		}
		if order.Ascending {
			orderStr += " ASC"
		} else {
			orderStr += " DESC"
		}
		orderStrs = append(orderStrs, orderStr)
	}
	return strings.Join(orderStrs, ", ")
}
