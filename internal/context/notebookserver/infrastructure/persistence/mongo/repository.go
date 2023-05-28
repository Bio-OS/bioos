package mongo

import (
	"context"

	"github.com/vinllen/mgo/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Bio-OS/bioos/internal/context/notebookserver/domain"
)

const notebookServerCollection = "notebookserver"

type repository struct {
	collection *mongo.Collection
}

// NewRepository ...
func NewRepository(ctx context.Context, mongoDB *mongo.Database) (domain.Repository, error) {
	collection := mongoDB.Collection(notebookServerCollection)
	if _, err := collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.M{"id": 1}, Options: options.Index().SetUnique(true)},
	}); err != nil {
		return nil, err
	}

	return &repository{
		collection: collection,
	}, nil
}

func (r *repository) Get(ctx context.Context, id string) (*domain.NotebookServer, error) {
	filter := bson.M{
		"id": id,
	}
	var result notebookServer
	if err := r.collection.FindOne(ctx, filter).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return result.toDO(), nil
}

func (r *repository) Save(ctx context.Context, do *domain.NotebookServer) error {
	po := newNotebookServer(do)
	filter := bson.M{"id": po.ID}
	result, err := r.collection.ReplaceOne(ctx, filter, po)
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		_, err = r.collection.InsertOne(ctx, po)
	}
	return err
}

func (r *repository) Delete(ctx context.Context, do *domain.NotebookServer) error {
	filter := bson.M{
		"id": do.ID,
	}
	if _, err := r.collection.DeleteOne(ctx, filter); err != nil {
		return err
	}
	return nil
}
