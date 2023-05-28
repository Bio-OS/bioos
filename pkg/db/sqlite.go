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
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/pflag"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const inmemory = ":memory:"

// SQLite3Options stands for sqlite options.
type SQLite3Options struct {
	File string `json:"file" mapstructure:"file"`
}

// NewSQLite3Options return a sqlite option.
func NewSQLite3Options() *SQLite3Options {
	return &SQLite3Options{}
}

func (o *SQLite3Options) Validate() error {
	if len(o.File) == 0 {
		return errors.New("options file no set")
	}

	if strings.ToLower(o.File) == inmemory {
		return nil
	}

	_, err := os.Stat(o.File)

	if os.IsNotExist(err) {
		return fmt.Errorf("options file not exist: %w", err)
	}

	if err != nil {
		return fmt.Errorf("os stat err:%w", err)
	}

	return nil
}

func (o *SQLite3Options) Enabled() bool {
	return len(o.File) > 0
}

func (o *SQLite3Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.File, "sqlite-file", "", "sqlite file name e.g. test.db")
}

// GetGORMInstance ...
func (o *SQLite3Options) GetGORMInstance(ctx context.Context) (*gorm.DB, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	orm, err := gorm.Open(sqlite.Open(o.File), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("sqlite open fail: %w", err)
	}
	instance, err := orm.DB()
	if err != nil {
		return nil, fmt.Errorf("sqlite get db instance fail: %w", err)
	}
	instance.SetMaxOpenConns(1)
	if err := instance.PingContext(ctx); err != nil {
		return nil, err
	}
	return orm, nil
}
