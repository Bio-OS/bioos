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
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	MySQLDefaultPort            = "3306"
	MySQLDefaultDatabase        = "bioos"
	MySQLDefaultMaxIdleConns    = 10
	MySQLDefaultMaxOpenConns    = 100
	MySQLDefaultCreateBatchSize = 1000
	MySQLDefaultConnMaxLifeTime = time.Hour
	MySQLDefaultConnMaxIdleTime = 30 * time.Second
)

type MySQLOptions struct {
	Username        string        `json:"username" mapstructure:"username"`
	Password        string        `json:"password" mapstructure:"password,omitempty"`
	Host            string        `json:"host" mapstructure:"host"`
	Port            string        `json:"port" mapstructure:"port"`
	Database        string        `json:"database" mapstructure:"database"`
	MaxIdleConns    int           `json:"maxIdleConns" mapstructure:"maxIdleConns"`
	MaxOpenConns    int           `json:"maxOpenConns" mapstructure:"maxOpenConns"`
	CreateBatchSize int           `json:"createBatchSize" mapstructure:"createBatchSize,omitempty"`
	ConnMaxLifetime time.Duration `json:"connMaxLifetime" mapstructure:"connMaxLifetime"`
	ConnMaxIdleTime time.Duration `json:"connMaxIdleTime" mapstructure:"connMaxIdletime"`
}

// NewMySQLOptions new a mysql option.
func NewMySQLOptions() *MySQLOptions {
	return &MySQLOptions{}
}

// Validate validate log options is valid.
func (o *MySQLOptions) Validate() error {
	return nil
}

func (o *MySQLOptions) Enabled() bool {
	return len(o.Host) > 0
}

func (o *MySQLOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Username, "mysql-username", "", "mysql db username")
	fs.StringVar(&o.Password, "mysql-password", "", "mysql db password")
	fs.StringVar(&o.Host, "mysql-host", "", "mysql db host")
	fs.StringVar(&o.Port, "mysql-port", MySQLDefaultPort, "mysql db port")
	fs.StringVar(&o.Database, "mysql-database", MySQLDefaultDatabase, "mysql database name")
	fs.IntVar(&o.MaxIdleConns, "mysql-max-idle-conns", MySQLDefaultMaxIdleConns, "mysql max idle conns")
	fs.IntVar(&o.MaxOpenConns, "mysql-max-open-conns", MySQLDefaultMaxOpenConns, "mysql max open conns")
	fs.IntVar(&o.CreateBatchSize, "mysql-create-batch-size", MySQLDefaultCreateBatchSize, "mysql create batch size")
	fs.DurationVar(&o.ConnMaxIdleTime, "mysql-conn-max-idle-time", MySQLDefaultConnMaxIdleTime, "mysql conn max idle time")
	fs.DurationVar(&o.ConnMaxLifetime, "mysql-conn-max-life-time", MySQLDefaultConnMaxLifeTime, "mysql conn max life time")
}

// GetGORMInstance ...
func (o *MySQLOptions) GetGORMInstance(ctx context.Context) (*gorm.DB, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		o.Username,
		o.Password,
		o.Host,
		o.Port,
		o.Database,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:          logger.Default.LogMode(logger.Info),
		CreateBatchSize: o.CreateBatchSize,
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	// set connection pool
	sqlDB.SetMaxIdleConns(o.MaxIdleConns)
	sqlDB.SetMaxOpenConns(o.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(o.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(o.ConnMaxIdleTime)

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}
