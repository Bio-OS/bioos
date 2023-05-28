//
// Copyright 2023 Beijing Volcano Engine Technology Ltd.
// Copyright 2023 Guangzhou Laboratory
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoPaginate struct {
	limit int64
	page  int64
}

func NewMongoPaginate(limit, page int) *MongoPaginate {
	return &MongoPaginate{
		limit: int64(limit),
		page:  int64(page),
	}
}

func (mp *MongoPaginate) GetPaginatedOpts() *options.FindOptions {
	l := mp.limit
	skip := mp.page*mp.limit - mp.limit
	fOpt := options.FindOptions{Limit: &l, Skip: &skip}

	return &fOpt
}

func EnsureIndex(ctx context.Context, collection *mongo.Collection, indexName string, unique bool, index primitive.D) (err error) {
	var cursor *mongo.Cursor
	if cursor, err = collection.Indexes().List(ctx); err != nil {
		return
	}

	var exist bool
	for cursor.Next(ctx) {
		var index = bson.M{}
		if err = cursor.Decode(&index); err != nil {
			return
		}
		if index["name"].(string) == indexName {
			exist = true
		}
	}

	if !exist {
		var indexOptions = &options.IndexOptions{}
		indexOptions.SetUnique(unique).SetName(indexName)
		if _, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys:    index,
			Options: indexOptions,
		}); err != nil {
			return
		}
	}
	return
}
