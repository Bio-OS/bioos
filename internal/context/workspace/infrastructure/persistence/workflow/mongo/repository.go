package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	domain "github.com/Bio-OS/bioos/internal/context/workspace/domain/workflow"
	"github.com/Bio-OS/bioos/internal/context/workspace/interface/grpc/proto"
	apperrors "github.com/Bio-OS/bioos/pkg/errors"
	applog "github.com/Bio-OS/bioos/pkg/log"
	"github.com/Bio-OS/bioos/pkg/utils"
)

func NewRepository(ctx context.Context, mongoDB *mongo.Database, mongoClient *mongo.Client) (domain.Repository, error) {
	r := &repository{
		client:                    mongoClient,
		workflowCollection:        mongoDB.Collection(workflowCollection),
		workflowVersionCollection: mongoDB.Collection(workflowVersionCollection),
		workflowFileCollection:    mongoDB.Collection(workflowFileCollection),
	}
	err := utils.EnsureIndex(ctx, r.workflowCollection, "idx_name_ws", true, primitive.D{
		{Key: "name", Value: 1},
		{Key: "workspaceID", Value: 1},
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

type repository struct {
	client                    *mongo.Client
	workflowCollection        *mongo.Collection
	workflowVersionCollection *mongo.Collection
	workflowFileCollection    *mongo.Collection
}

var _ domain.Repository = &repository{}

func (r *repository) Save(ctx context.Context, wf *domain.Workflow) error {
	session, err := r.client.StartSession()
	if err != nil {
		return apperrors.NewInternalError(err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// save workflow
		workflowPO := workflowDOToPO(wf)
		filter := bson.M{"id": workflowPO.ID}
		update := bson.M{"$set": workflowPO}
		opts := options.Update().SetUpsert(true)

		if _, err := r.workflowCollection.UpdateOne(ctx, filter, update, opts); err != nil {
			return nil, apperrors.NewInternalError(err)
		}

		// save workflow versions
		for _, workflowVersionDO := range wf.Versions {
			workflowVersionPO, err := workflowVersionDOtoPO(wf.ID, workflowVersionDO)
			if err != nil {
				return nil, apperrors.NewInternalError(err)
			}
			versionFilter := bson.M{"id": workflowVersionPO.ID}
			versionUpdate := bson.M{"$set": workflowVersionPO}

			if _, err := r.workflowVersionCollection.UpdateOne(ctx, versionFilter, versionUpdate, opts); err != nil {
				return nil, apperrors.NewInternalError(err)
			}

			// save workflow files
			for _, workflowFileDO := range workflowVersionDO.Files {
				workflowFilePO := workflowFileDOToPO(workflowVersionDO.ID, workflowFileDO)

				fileFilter := bson.M{"id": workflowFilePO.ID}
				fileUpdate := bson.M{"$set": workflowFilePO}

				if _, err := r.workflowFileCollection.UpdateOne(ctx, fileFilter, fileUpdate, opts); err != nil {
					return nil, apperrors.NewInternalError(err)
				}
			}
		}

		return nil, nil
	})

	if err != nil {
		if e := session.AbortTransaction(ctx); e != nil {
			return apperrors.NewInternalError(e)
		}
		return apperrors.NewInternalError(err)
	}
	return nil
}

func (r *repository) Get(ctx context.Context, workspaceID string, workflowID string) (*domain.Workflow, error) {
	filter := bson.M{
		"id":          workflowID,
		"workspaceID": workspaceID,
	}
	var workflowPO workflow
	err := r.workflowCollection.FindOne(ctx, filter).Decode(&workflowPO)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, proto.ErrorWorkflowNotFound("workflow:%s in workspace:%s not found", workflowID, workspaceID)
		}
		applog.Errorw("failed to get workflow", "err", err)
		return nil, apperrors.NewInternalError(err)
	}

	// list workflow versions
	versionFilter := bson.M{
		"workflowID": workflowID,
	}
	cursor, err := r.workflowVersionCollection.Find(ctx, versionFilter)
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	workflowDO := workflowPO.toDO()
	for cursor.Next(ctx) {
		var workflowVersionPO workflowVersion
		if err := cursor.Decode(&workflowVersionPO); err != nil {
			return nil, apperrors.NewInternalError(err)
		}

		// list workflow files
		fileFilter := bson.M{
			"workflowID": workflowID,
		}
		fileCursor, err := r.workflowFileCollection.Find(ctx, fileFilter)
		if err != nil {
			return nil, apperrors.NewInternalError(err)
		}
		workflowVersionDO, err := workflowVersionPO.toDO()
		if err != nil {
			applog.Errorw("failed to convert workflow version PO to DO", "err", err)

			return nil, apperrors.NewInternalError(err)
		}
		for fileCursor.Next(ctx) {
			var workflowFilePO workflowFile
			if err := fileCursor.Decode(&workflowFilePO); err != nil {
				return nil, apperrors.NewInternalError(err)
			}
			workflowVersionDO.Files[workflowFilePO.ID] = workflowFilePO.toDO()
		}
		workflowDO.Versions[workflowVersionDO.ID] = workflowVersionDO
	}

	return workflowDO, nil
}

func (r *repository) Delete(ctx context.Context, wf *domain.Workflow) error {
	session, err := r.client.StartSession()
	if err != nil {
		return apperrors.NewInternalError(err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// delete workflow files
		var filter = bson.M{
			"workflowID": wf.ID,
		}
		cursor, err := r.workflowVersionCollection.Find(ctx, filter)
		if err != nil {
			return nil, apperrors.NewInternalError(err)
		}
		for cursor.Next(ctx) {
			var wv workflowVersion
			if err := cursor.Decode(&wv); err != nil {
				return nil, apperrors.NewInternalError(err)
			}
			_, err = r.workflowFileCollection.DeleteMany(sessCtx, bson.M{"workflowVersionID": wv.ID})
			if err != nil {
				return nil, apperrors.NewInternalError(err)
			}
		}

		// delete workflow versions
		_, err = r.workflowVersionCollection.DeleteMany(sessCtx, bson.M{"workflowID": wf.ID})
		if err != nil {
			return nil, apperrors.NewInternalError(err)
		}

		// delete workflow
		_, err = r.workflowCollection.DeleteOne(sessCtx, bson.M{"id": wf.ID})
		if err != nil {
			return nil, apperrors.NewInternalError(err)
		}
		return nil, nil
	})
	if err != nil {
		if e := session.AbortTransaction(ctx); e != nil {
			return apperrors.NewInternalError(e)
		}
		return apperrors.NewInternalError(err)
	}
	return nil
}

func (r *repository) List(ctx context.Context, workspaceID string) ([]string, error) {
	filter := bson.M{
		"workspaceID": workspaceID,
	}
	cursor, err := r.workflowCollection.Find(ctx, filter)
	if err != nil {
		return nil, apperrors.NewInternalError(err)
	}
	workflowIDs := make([]string, 0)
	for cursor.Next(ctx) {
		var workflowPO workflow
		if err := cursor.Decode(&workflowPO); err != nil {
			return nil, apperrors.NewInternalError(err)
		}
		workflowIDs = append(workflowIDs, workflowPO.ID)
	}

	return workflowIDs, nil
}
