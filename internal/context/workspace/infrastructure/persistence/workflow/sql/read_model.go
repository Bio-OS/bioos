package sql

import (
	"context"
	"errors"

	"gorm.io/gorm"

	query "github.com/Bio-OS/bioos/internal/context/workspace/application/query/workflow"
	"github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	applog "github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/utils"
)

func NewReadModel(_ context.Context, db *gorm.DB) (query.ReadModel, error) {
	return &readModel{db: db}, nil
}

type readModel struct {
	db *gorm.DB
}

var _ query.ReadModel = &readModel{}

func (r *readModel) GetById(ctx context.Context, workspaceID, workflowID string) (*query.Workflow, error) {
	wf := &workflow{
		ID:          workflowID,
		WorkspaceID: workspaceID,
	}
	if err := r.db.WithContext(ctx).First(wf).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, proto.ErrorWorkflowNotFound("workflow:%s in workspace:%s not found", workflowID, workspaceID)
		}
		return nil, apperrors.NewInternalError(err)
	}

	wfDTO := wf.toDTO()

	if wf.LatestVersion != "" {
		wvDTO, err := r.GetVersion(ctx, wf.LatestVersion)
		if err != nil {
			return nil, err
		}
		wfDTO.LatestVersion = wvDTO
	}
	return wfDTO, nil
}

func (r *readModel) GetByName(ctx context.Context, workspaceID, workflowName string) (*query.Workflow, error) {
	wf := &workflow{
		Name:        workflowName,
		WorkspaceID: workspaceID,
	}
	if err := r.db.WithContext(ctx).Where("name = ?", wf.Name).Where("workspace_id", wf.WorkspaceID).First(wf).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, proto.ErrorWorkflowNotFound("workflow:%s in workspace:%s not found", workflowName, workspaceID)
		}
		return nil, apperrors.NewInternalError(err)
	}

	wfDTO := wf.toDTO()

	if wf.LatestVersion != "" {
		wvDTO, err := r.GetVersion(ctx, wf.LatestVersion)
		if err != nil {
			return nil, err
		}
		wfDTO.LatestVersion = wvDTO
	}
	return wfDTO, nil
}

func (r *readModel) List(ctx context.Context, workspaceID string, pg *utils.Pagination, filter *query.ListWorkflowsFilter) ([]*query.Workflow, int, error) {
	var wfs []*workflow
	db := r.db.Model(&workflow{}).WithContext(ctx).Where("workspace_id = ?", workspaceID)
	if filter != nil {
		if len(filter.IDs) > 0 {
			db = db.Where("id in ?", filter.IDs)
		}
		if filter.SearchWord != "" {
			db = utils.SearchWordFilter(db, filter.SearchWord, []string{"name"}, filter.Exact)
		}
	}
	var cnt int64
	if err := db.Count(&cnt).Error; err != nil {
		applog.Errorw("failed to count workflows", "err", err)
		return nil, 0, apperrors.NewInternalError(err)
	}
	if pg != nil {
		db = db.Limit(pg.GetLimit()).
			Offset(pg.GetOffset()).
			Order(utils.DBOrder(pg.Orders, map[string]string{
				query.OrderByName:       "name",
				query.OrderByCreateTime: "created_at",
			}))
	}
	if err := db.Find(&wfs).Error; err != nil {
		return nil, 0, err
	}
	ret := make([]*query.Workflow, len(wfs))
	for index, wfPO := range wfs {
		wfDTO := wfPO.toDTO()

		if wfPO.LatestVersion != "" {
			wvDTO, err := r.GetVersion(ctx, wfPO.LatestVersion)
			if err != nil {
				return nil, 0, err
			}
			wfDTO.LatestVersion = wvDTO
		}

		ret[index] = wfDTO
	}

	return ret, int(cnt), nil
}

func (r *readModel) GetVersion(ctx context.Context, id string) (*query.WorkflowVersion, error) {
	wv := &workflowVersion{
		ID: id,
	}
	if err := r.db.WithContext(ctx).First(wv).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, proto.ErrorWorkflowVersionNotFound("workflow version:%s not found", id)
		}
		return nil, apperrors.NewInternalError(err)
	}
	workflowFilesDTO, _, err := r.ListFiles(ctx, id, nil, nil)
	if err != nil {
		return nil, err
	}
	fileInfos := make([]*query.WorkflowFileInfo, len(workflowFilesDTO))
	for index, file := range workflowFilesDTO {
		fileInfos[index] = file.ToWorkflowFileInfo()
	}
	wvDTO, err := wv.toDTO()
	if err != nil {
		applog.Errorw("fail to convert workflow version PO to DTO", "err", err)
		return nil, apperrors.NewInternalError(err)
	}
	wvDTO.Files = fileInfos

	return wvDTO, nil
}

func (r *readModel) ListVersions(ctx context.Context, workflowID string, pg *utils.Pagination, filter *query.ListWorkflowVersionsFilter) ([]*query.WorkflowVersion, int, error) {
	var wvs []*workflowVersion
	db := r.db.WithContext(ctx).Where("workflow_id = ?", workflowID)
	if pg != nil {
		db = db.Limit(pg.GetLimit()).
			Offset(pg.GetOffset()).
			Order(utils.DBOrder(pg.Orders, map[string]string{
				query.VersionOrderByStatus:   "status",
				query.VersionOrderByLanguage: "language",
				query.VersionOrderBySource:   "source",
			}))
	}
	if filter != nil {
		if len(filter.IDs) > 0 {
			db = db.Where("id IN ?", filter.IDs)
		}
	}
	if err := db.Find(&wvs).Error; err != nil {
		return nil, 0, err
	}
	var cnt int64

	if err := db.Count(&cnt).Error; err != nil {
		applog.Errorw("failed to count workflow versions", "err", err)
		return nil, 0, apperrors.NewInternalError(err)
	}

	ret := make([]*query.WorkflowVersion, len(wvs))
	for index, workflowVersionPO := range wvs {
		workflowVersionDTO, err := workflowVersionPO.toDTO()
		if err != nil {
			applog.Errorw("fail to convert workflow version PO to DTO", "err", err)
			return nil, 0, apperrors.NewInternalError(err)
		}
		// list workflow version files
		workflowFilesDTO, _, err := r.ListFiles(ctx, workflowVersionPO.ID, nil, nil)
		if err != nil {
			return nil, 0, err
		}
		fileInfos := make([]*query.WorkflowFileInfo, len(workflowFilesDTO))
		for i, file := range workflowFilesDTO {
			fileInfos[i] = file.ToWorkflowFileInfo()
		}
		workflowVersionDTO.Files = fileInfos
		ret[index] = workflowVersionDTO
	}

	return ret, int(cnt), nil
}

func (r *readModel) GetFile(ctx context.Context, id string) (*query.WorkflowFile, error) {
	wf := &workflowFile{
		ID: id,
	}
	if err := r.db.WithContext(ctx).First(wf).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, proto.ErrorWorkflowFileNotFound("workflow file:%s not found", id)
		}
		return nil, err
	}

	return wf.toDTO(), nil
}

func (r *readModel) ListFiles(ctx context.Context, workflowVersionID string, pg *utils.Pagination, filter *query.ListWorkflowFilesFilter) ([]*query.WorkflowFile, int, error) {
	var wf []*workflowFile
	db := r.db.WithContext(ctx).Where("workflow_version_id = ?", workflowVersionID)
	if pg != nil {
		db = db.Limit(pg.GetLimit()).
			Offset(pg.GetOffset()).
			Order(utils.DBOrder(pg.Orders, map[string]string{
				query.FileOrderByPath:    "path",
				query.FileOrderByVersion: "version",
			}))
	}

	if filter != nil {
		if len(filter.IDs) > 0 {
			db = db.Where("id IN ?", filter.IDs)
		}
	}

	if err := db.Find(&wf).Error; err != nil {
		return nil, 0, err
	}
	var cnt int64

	if err := db.Count(&cnt).Error; err != nil {
		applog.Errorw("failed to count workflow files", "err", err)
		return nil, 0, apperrors.NewInternalError(err)
	}

	ret := make([]*query.WorkflowFile, len(wf))
	for index, po := range wf {
		ret[index] = po.toDTO()
	}

	return ret, int(cnt), nil
}
