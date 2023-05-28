package mongo

import (
	"context"

	"github.com/vinllen/mgo/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Bio-OS/bioos/internal/context/notebookserver/application/query"
)

type readModel struct {
	collection *mongo.Collection
}

// NewReadModel ...
func NewReadModel(ctx context.Context, mongoDB *mongo.Database) (query.ReadModel, error) {
	collection := mongoDB.Collection(notebookServerCollection)

	return &readModel{
		collection: collection,
	}, nil
}

func (r *readModel) ListSettingsByWorkspace(ctx context.Context, workspaceID string) ([]*query.NotebookSettings, error) {
	filter := bson.M{
		"workspaceID": workspaceID,
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var res []*query.NotebookSettings
	for cursor.Next(ctx) {
		var po notebookServer
		if err = cursor.Decode(&po); err != nil {
			return nil, err
		}
		res = append(res, po.toDTO())
	}
	return res, nil
}

func (r *readModel) GetSettingsByID(ctx context.Context, workspaceID, id string) (*query.NotebookSettings, error) {
	filter := bson.M{
		"workspaceID": workspaceID,
		"id":          id,
	}
	var result notebookServer
	if err := r.collection.FindOne(ctx, filter).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return result.toDTO(), nil
}
