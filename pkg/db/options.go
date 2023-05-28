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

package db

import (
	"fmt"

	"github.com/spf13/pflag"

	apperrors "github.com/Bio-OS/bioos/pkg/errors"
)

type Options struct {
	SQLite3 *SQLite3Options `json:"sqlite3,omitempty" mapstructure:"sqlite3,omitempty"`
	MySQL   *MySQLOptions   `json:"mysql,omitempty" mapstructure:"mysql,omitempty"`
	Mongo   *MongoOptions   `json:"mongo,omitempty" mapstructure:"mongo,omitempty"`
}

// NewOptions new a db option.
func NewOptions() *Options {
	return &Options{
		SQLite3: NewSQLite3Options(),
		MySQL:   NewMySQLOptions(),
		Mongo:   NewMongoOptions(),
	}
}

// Validate check db options is valid.
func (o *Options) Validate() error {
	if o.SQLite3 != nil && (*o.SQLite3 != SQLite3Options{}) {
		if err := o.SQLite3.Validate(); err != nil {
			return apperrors.NewInternalError(err)
		}
	}
	if o.MySQL != nil && (*o.MySQL != MySQLOptions{}) {
		fmt.Println("mysql")
		if err := o.MySQL.Validate(); err != nil {
			return apperrors.NewInternalError(err)
		}
	}
	if o.Mongo != nil && (*o.Mongo != MongoOptions{}) {
		fmt.Println("mongo")
		if err := o.Mongo.Validate(); err != nil {
			return apperrors.NewInternalError(err)
		}
	}
	return nil
}

// AddFlags ...
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	o.SQLite3.AddFlags(fs)
	o.MySQL.AddFlags(fs)
	o.Mongo.AddFlags(fs)
}
