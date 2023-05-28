package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	query "github.com/Bio-OS/bioos/internal/context/workspace/application/query/workflow"
	"github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	applog "github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/utils"
)

const (
	workflowCollection        = "workflow"
	workflowVersionCollection = "workflowVersion"
	workflowFileCollection    = "workflowFile"
)

func NewReadModel(_ context.Context, mongoDB *mongo.Database) (query.ReadModel, error) {
	return &readModel{
		workflowCollection:        mongoDB.Collection(workflowCollection),
		workflowVersionCollection: mongoDB.Collection(workflowVersionCollection),
		workflowFileCollection:    mongoDB.Collection(workflowFileCollection),
	}, nil
}

type readModel struct {
	workflowCollection        *mongo.Collection
	workflowVersionCollection *mongo.Collection
	workflowFileCollection    *mongo.Collection
}

var _ query.ReadModel = &readModel{}

func (r *readModel) GetById(ctx context.Context, workspaceID, workflowID string) (*query.Workflow, error) {
	filter := bson.M{
		"id":          workflowID,
		"workspaceID": workspaceID,
	}
	var wf workflow
	err := r.workflowCollection.FindOne(ctx, filter).Decode(&wf)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, proto.ErrorWorkflowNotFound("workflow:%s in workspace:%s not found", workflowID, workspaceID)
		}
		return nil, err
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
	filter := bson.M{
		"name":        workflowName,
		"workspaceID": workspaceID,
	}
	var wf workflow
	err := r.workflowCollection.FindOne(ctx, filter).Decode(&wf)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, proto.ErrorWorkflowNotFound("workflow:%s in workspace:%s not found", workflowName, workspaceID)
		}
		return nil, err
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

func (r *readModel) List(ctx context.Context, workspaceID string, pg *utils.Pagination, listFilter *query.ListWorkflowsFilter) ([]*query.Workflow, int, error) {
	var filter = bson.M{
		"workspaceID": workspaceID,
	}
	if listFilter != nil {
		if len(listFilter.SearchWord) > 0 {
			if listFilter.Exact {
				filter["name"] = listFilter.SearchWord
			} else {
				filter["name"] = bson.M{"$regex": listFilter.SearchWord, "$options": ""}
			}
		}
		if len(listFilter.IDs) > 0 {
			filter["id"] = bson.M{"$in": listFilter.IDs}
		}
	}

	var cursor *mongo.Cursor
	var err error
	if pg != nil {
		cursor, err = r.workflowCollection.Find(ctx, filter, utils.NewMongoPaginate(pg.GetLimit(), pg.GetPage()).GetPaginatedOpts())
	} else {
		cursor, err = r.workflowCollection.Find(ctx, filter)
	}
	if err != nil {
		return nil, 0, err
	}

	var result []*query.Workflow
	for cursor.Next(ctx) {
		var wf workflow
		if err := cursor.Decode(&wf); err != nil {
			return nil, 0, err
		}
		wfDTO := wf.toDTO()
		wvDTO, err := r.GetVersion(ctx, wf.LatestVersion)
		if err != nil {
			return nil, 0, err
		}
		wfDTO.LatestVersion = wvDTO
		result = append(result, wfDTO)
	}
	count, err := r.workflowCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return result, int(count), nil
}

func (r *readModel) GetVersion(ctx context.Context, id string) (*query.WorkflowVersion, error) {
	filter := bson.M{
		"id": id,
	}
	var wv workflowVersion
	err := r.workflowVersionCollection.FindOne(ctx, filter).Decode(&wv)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, proto.ErrorWorkflowVersionNotFound("workflow file:%s  not found", id)
		}
		return nil, err
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

func (r *readModel) ListVersions(ctx context.Context, workflowID string, pg *utils.Pagination, versionFilter *query.ListWorkflowVersionsFilter) ([]*query.WorkflowVersion, int, error) {
	var filter = bson.M{
		"workflowID": workflowID,
	}
	if versionFilter != nil {
		if len(versionFilter.IDs) > 0 {
			filter["id"] = bson.M{"$in": versionFilter.IDs}
		}
	}

	var cursor *mongo.Cursor
	var err error
	if pg != nil {
		cursor, err = r.workflowVersionCollection.Find(ctx, filter, utils.NewMongoPaginate(pg.GetLimit(), pg.GetPage()).GetPaginatedOpts())
	} else {
		cursor, err = r.workflowVersionCollection.Find(ctx, filter)
	}
	if err != nil {
		return nil, 0, apperrors.NewInternalError(err)
	}

	var result []*query.WorkflowVersion
	for cursor.Next(ctx) {
		var wv workflowVersion
		if err := cursor.Decode(&wv); err != nil {
			return nil, 0, err
		}
		workflowVersionDTO, err := wv.toDTO()
		if err != nil {
			applog.Errorw("fail to convert workflow version PO to DTO", "err", err)
			return nil, 0, apperrors.NewInternalError(err)
		}
		// list workflow version files
		workflowFilesDTO, _, err := r.ListFiles(ctx, wv.ID, nil, nil)
		if err != nil {
			return nil, 0, err
		}
		fileInfos := make([]*query.WorkflowFileInfo, len(workflowFilesDTO))
		for i, file := range workflowFilesDTO {
			fileInfos[i] = file.ToWorkflowFileInfo()
		}
		workflowVersionDTO.Files = fileInfos

		result = append(result, workflowVersionDTO)
	}
	count, err := r.workflowVersionCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return result, int(count), nil
}

func (r *readModel) GetFile(ctx context.Context, id string) (*query.WorkflowFile, error) {
	filter := bson.M{
		"id": id,
	}
	var wf workflowFile
	err := r.workflowFileCollection.FindOne(ctx, filter).Decode(&wf)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, proto.ErrorWorkflowFileNotFound("workflow file:%s  not found", id)
		}
		return nil, err
	}

	return wf.toDTO(), nil
}

func (r *readModel) ListFiles(ctx context.Context, workflowVersionID string, pg *utils.Pagination, fileFilter *query.ListWorkflowFilesFilter) ([]*query.WorkflowFile, int, error) {
	var filter = bson.M{
		"workflowVersionID": workflowVersionID,
	}
	if fileFilter != nil {
		if len(fileFilter.IDs) > 0 {
			filter["id"] = bson.M{"$in": fileFilter.IDs}
		}
	}

	var cursor *mongo.Cursor
	var err error
	if pg != nil {
		cursor, err = r.workflowFileCollection.Find(ctx, filter, utils.NewMongoPaginate(pg.GetLimit(), pg.GetPage()).GetPaginatedOpts())
	} else {
		cursor, err = r.workflowFileCollection.Find(ctx, filter)
	}

	if err != nil {
		return nil, 0, err
	}

	var result []*query.WorkflowFile
	for cursor.Next(ctx) {
		var wf workflowFile
		if err := cursor.Decode(&wf); err != nil {
			return nil, 0, err
		}

		result = append(result, wf.toDTO())
	}
	count, err := r.workflowFileCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return result, int(count), nil
}
