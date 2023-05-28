package mysql

import (
	"context"
	"errors"

	"gorm.io/gorm"

	query "github.com/Bio-OS/bioos/internal/context/workspace/application/query/workspace"
	"github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	applog "github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/utils"
)

type workspaceReadModel struct {
	db *gorm.DB
}

// NewWorkspaceReadModel ...
func NewWorkspaceReadModel(_ context.Context, db *gorm.DB) (query.WorkspaceReadModel, error) {
	return &workspaceReadModel{db: db}, nil
}

func (w *workspaceReadModel) ListWorkspaces(ctx context.Context, pg utils.Pagination, filter *query.ListWorkspacesFilter) ([]*query.WorkspaceItem, error) {
	db := w.db.WithContext(ctx).
		Limit(pg.GetLimit()).
		Offset(pg.GetOffset()).
		Order(utils.DBOrder(pg.Orders, map[string]string{
			query.OrderByName:       "name",
			query.OrderByCreateTime: "create_time",
		}))
	db = listWorkspacesFilter(db, filter)
	var ws []*Workspace
	if err := db.Find(&ws).Error; err != nil {
		applog.Errorw("failed to list workspaces", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	ret := make([]*query.WorkspaceItem, len(ws))
	for index, po := range ws {
		ret[index] = WorkspacePOToWorkspaceDTO(ctx, po)
	}
	return ret, nil
}

func (w *workspaceReadModel) CountWorkspaces(ctx context.Context, filter *query.ListWorkspacesFilter) (int, error) {
	db := w.db.WithContext(ctx).Model(&Workspace{})
	db = listWorkspacesFilter(db, filter)
	var cnt int64
	if err := db.Count(&cnt).Error; err != nil {
		applog.Errorw("failed to count workspaces", "err", err)
		return 0, apperrors.NewInternalError(err)
	}
	return int(cnt), nil
}

func (w *workspaceReadModel) GetWorkspaceById(ctx context.Context, id string) (*query.WorkspaceItem, error) {
	var ws Workspace
	if err := w.db.WithContext(ctx).Where("id = ?", id).First(&ws).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, proto.ErrorWorkspaceNotFound("workspace: %s not found", id)
		}
		applog.Errorw("failed to get workspace by id", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	return WorkspacePOToWorkspaceDTO(ctx, &ws), nil
}

func listWorkspacesFilter(db *gorm.DB, filter *query.ListWorkspacesFilter) *gorm.DB {
	if filter == nil {
		return db
	}
	if len(filter.SearchWord) > 0 {
		db = utils.SearchWordFilter(db, filter.SearchWord, []string{"name", "description"}, filter.Exact)
	}
	if len(filter.IDs) > 0 {
		db = db.Where("id IN ?", filter.IDs)
	}
	return db
}
