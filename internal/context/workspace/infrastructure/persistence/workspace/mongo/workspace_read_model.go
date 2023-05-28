package mongo

import (
	"context"
	"fmt"

	"github.com/vinllen/mgo/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	query "github.com/Bio-OS/bioos/internal/context/workspace/application/query/workspace"
	"github.com/Bio-OS/bioos/pkg/utils"
)

type workspaceReadModel struct {
	collection *mongo.Collection
}

func NewWorkspaceReadModel(ctx context.Context, mongoDB *mongo.Database) (query.WorkspaceReadModel, error) {
	collection := mongoDB.Collection(WorkspaceCollection)

	return &workspaceReadModel{
		collection: collection,
	}, nil
}

type mongoPaginate struct {
	limit int64
	page  int64
}

func newMongoPaginate(limit, page int) *mongoPaginate {
	return &mongoPaginate{
		limit: int64(limit),
		page:  int64(page),
	}
}

func (mp *mongoPaginate) getPaginatedOpts() *options.FindOptions {
	l := mp.limit
	skip := mp.page*mp.limit - mp.limit
	fOpt := options.FindOptions{Limit: &l, Skip: &skip}

	return &fOpt
}

func toMongoFindOptions(pg *utils.Pagination) *options.FindOptions {
	var l int64 = int64(pg.GetLimit())
	var skip int64 = int64(pg.GetPage()*pg.GetLimit() - pg.GetLimit())
	var sort interface{}
	if len(pg.Orders) > 0 {
		d := make(bson.D, 0, len(pg.Orders))
		for _, order := range pg.Orders {
			elem := bson.DocElem{}
			switch order.Field {
			case query.OrderByName:
				elem.Name = "name"
			case query.OrderByCreateTime:
				elem.Name = "createTime"
			default:
				continue
			}
			if order.Ascending {
				elem.Value = 1
			} else {
				elem.Value = -1
			}
			d = append(d, elem)
		}
		sort = d
	}
	return &options.FindOptions{Limit: &l, Skip: &skip, Sort: sort}
}

func (w workspaceReadModel) SearchWorkspaceByName(ctx context.Context, name string, pg *utils.Pagination) ([]*query.WorkspaceItem, error) {
	w.collection.Find(ctx, bson.M{
		"name": fmt.Sprintf("/%s/", name),
	})

	cursor, err := w.collection.Find(ctx, bson.M{
		"name": fmt.Sprintf("/%s/", name),
	}, toMongoFindOptions(pg))
	if err != nil {
		return nil, err
	}

	var result []*query.WorkspaceItem
	for cursor.Next(ctx) {
		var wsPO workspacePO
		if err := cursor.Decode(&wsPO); err != nil {
			return nil, err
		}
		ws, err := workspacePOToQueryItem(ctx, &wsPO)
		if err != nil {
			return nil, err
		}

		result = append(result, ws)
	}

	return result, nil
}

func (w workspaceReadModel) ListWorkspaces(ctx context.Context, pg utils.Pagination, filter *query.ListWorkspacesFilter) ([]*query.WorkspaceItem, error) {
	cursor, err := w.collection.Find(ctx, getFilter(filter), newMongoPaginate(pg.GetLimit(), pg.GetPage()).getPaginatedOpts())
	if err != nil {
		return nil, err
	}

	var result []*query.WorkspaceItem
	for cursor.Next(ctx) {
		var wsPO workspacePO
		if err := cursor.Decode(&wsPO); err != nil {
			return nil, err
		}
		ws, err := workspacePOToQueryItem(ctx, &wsPO)
		if err != nil {
			return nil, err
		}

		result = append(result, ws)
	}

	return result, nil
}

func getFilter(filter *query.ListWorkspacesFilter) bson.M {
	res := bson.M{}
	if filter != nil {
		if len(filter.SearchWord) > 0 {
			if filter.Exact {
				res["name"] = filter.SearchWord
			} else {
				res["name"] = bson.M{"$regex": filter.SearchWord, "$options": ""}
			}
		}
		if len(filter.IDs) > 0 {
			res["id"] = bson.M{"$in": filter.IDs}
		}
	}
	return res
}

func (w workspaceReadModel) GetWorkspaceById(ctx context.Context, id string) (*query.WorkspaceItem, error) {
	filter := bson.M{
		"id": id,
	}
	var result workspacePO
	if err := w.collection.FindOne(ctx, filter).Decode(&result); err != nil {
		return nil, err
	}

	return workspacePOToQueryItem(ctx, &result)
}

func (w workspaceReadModel) CountWorkspaces(ctx context.Context, filter *query.ListWorkspacesFilter) (int, error) {
	opts := options.Count().SetHint("_id_")
	count, err := w.collection.CountDocuments(ctx, getFilter(filter), opts)
	return int(count), err
}
