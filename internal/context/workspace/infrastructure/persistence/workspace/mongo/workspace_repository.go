package mongo

import (
	"context"

	"github.com/vinllen/mgo/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Bio-OS/bioos/internal/context/workspace/domain/workspace"
	"github.com/Bio-OS/bioos/pkg/utils"
)

const WorkspaceCollection = "workspace"

type workspaceRepository struct {
	collection *mongo.Collection
}

var _ workspace.Repository = &workspaceRepository{}

func (r *workspaceRepository) Delete(ctx context.Context, w *workspace.Workspace) error {
	filter := bson.M{
		"id": w.GetID(),
	}

	if _, err := r.collection.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}

// NewWorkspaceRepository ...
func NewWorkspaceRepository(ctx context.Context, mongoDB *mongo.Database) (workspace.Repository, error) {
	collection := mongoDB.Collection(WorkspaceCollection)
	if _, err := collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.M{"id": 1}, Options: options.Index().SetUnique(true)},
	}); err != nil {
		return nil, err
	}
	err := utils.EnsureIndex(ctx, collection, "idx_name", true, primitive.D{{Key: "name", Value: 1}})
	if err != nil {
		return nil, err
	}

	return &workspaceRepository{
		collection: collection,
	}, nil
}

func (r *workspaceRepository) Get(ctx context.Context, id string) (*workspace.Workspace, error) {
	filter := bson.M{
		"id": id,
	}
	var result workspacePO
	if err := r.collection.FindOne(ctx, filter).Decode(&result); err != nil {
		return nil, err
	}

	return workspacePOToWorkspaceDO(ctx, &result)
}

func (r *workspaceRepository) Save(ctx context.Context, w *workspace.Workspace) error {
	ws, err := workspaceDOtoWorkspacePO(ctx, w)
	if err != nil {
		return err
	}
	filter := bson.M{"id": ws.ID}
	result, err := r.collection.ReplaceOne(ctx, filter, ws)
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		_, err = r.collection.InsertOne(ctx, ws)
	}
	return err
}
