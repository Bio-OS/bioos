package sql

import (
	"context"
	"strings"

	"gorm.io/gorm"

	applog "github.com/Bio-OS/bioos/pkg/log"

	query "github.com/Bio-OS/bioos/internal/context/submission/application/query/submission"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	"github.com/Bio-OS/bioos/pkg/utils"
)

type submissionReadModel struct {
	db *gorm.DB
}

// NewSubmissionReadModel ...
func NewSubmissionReadModel(ctx context.Context, db *gorm.DB) (query.ReadModel, error) {
	if err := db.WithContext(ctx).AutoMigrate(&SubmissionModel{}); err != nil {
		return nil, apperrors.NewInternalError(err)
	}

	return &submissionReadModel{db: db}, nil
}

func (s *submissionReadModel) ListSubmissions(ctx context.Context, workspaceID string, pg *utils.Pagination, filter *query.ListSubmissionsFilter) ([]*query.SubmissionItem, error) {
	dbChain := s.db.WithContext(ctx).Model(&SubmissionModel{}).Where("workspace_id = ?", workspaceID).Limit(pg.GetLimit()).Offset(pg.GetOffset()).Order(ordersToOrderDB(pg.Orders))
	dbChain = listSubmissionsFilter(dbChain, filter)
	var sbs []*Submission
	if err := dbChain.Find(&sbs).Error; err != nil {
		applog.Errorw("failed to list submissions", "err", err)
		return nil, apperrors.NewInternalError(err)

	}
	ret := make([]*query.SubmissionItem, len(sbs))
	for index, po := range sbs {
		item, err := SubmissionPOToSubmissionDTO(ctx, po)
		if err != nil {
			applog.Errorw("failed to convert submission po to DTO", "err", err)
			return nil, apperrors.NewInternalError(err)
		}
		ret[index] = item
	}
	return ret, nil
}

func (s *submissionReadModel) CountSubmissions(ctx context.Context, workspaceID string, filter *query.ListSubmissionsFilter) (int, error) {
	dbChain := s.db.WithContext(ctx).Model(&SubmissionModel{}).Where("workspace_id = ?", workspaceID)
	dbChain = listSubmissionsFilter(dbChain, filter)
	var count int64
	if err := dbChain.Count(&count).Error; err != nil {
		applog.Errorw("failed to count submissions", "err", err)
		return 0, apperrors.NewInternalError(err)
	}
	return int(count), nil
}

func listSubmissionsFilter(db *gorm.DB, filter *query.ListSubmissionsFilter) *gorm.DB {
	if filter == nil {
		return db
	}
	if filter.IDs != nil && len(filter.IDs) > 0 {
		db = db.Where("id IN ?", filter.IDs)
	}
	if len(filter.Name) > 0 {
		db = db.Where("name = ?", filter.Name)
	}
	if len(filter.Status) > 0 {
		db = db.Where("status IN ?", filter.Status)
	}
	if len(filter.SearchWord) > 0 {
		db = utils.SearchWordFilter(db, filter.SearchWord, []string{"name"}, filter.Exact)
	}
	if len(filter.WorkflowID) != 0 {
		db = db.Where("workflow_id = ?", filter.WorkflowID)
	}

	return db
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
