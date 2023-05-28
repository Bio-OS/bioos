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
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	MongoDefaultPort                   = "27017"
	MongoDefaultDatabase               = "bioos"
	MongoDefaultConnectTimeout         = 30 * time.Second
	MongoDefaultServerSelectionTimeout = 30 * time.Second
	MongoDefaultSocketTimeout          = 30 * time.Second
	MongoDefaultHeartbeatInterval      = 5 * time.Second
	MongoDefaultMaxConnIdleTime        = 30 * time.Second
	MongoDefaultMaxPoolSize            = 100
	MongoDefaultMinPoolSize            = 1
)

type MongoOptions struct {
	Host                   string        `json:"host" mapstructure:"host"`
	Port                   string        `json:"port" mapstructure:"port"`
	Username               string        `json:"username" mapstructure:"username"`
	Password               string        `json:"password" mapstructure:"password,omitempty"`
	Database               string        `json:"database" mapstructure:"database"`
	ConnectTimeout         time.Duration `json:"connectTimeout" mapstructure:"connectTimeout"`
	ServerSelectionTimeout time.Duration `json:"serverSelectionTimeout" mapstructure:"serverSelectionTimeout"`
	SocketTimeout          time.Duration `json:"socketTimeout" mapstructure:"socketTimeout"`
	HeartbeatInterval      time.Duration `json:"heartbeatInterval" mapstructure:"heartbeatInterval"`
	MaxConnIdleTime        time.Duration `json:"maxConnIdleTime" mapstructure:"maxConnIdleTime"`
	MaxPoolSize            uint64        `json:"maxPoolSize" mapstructure:"maxPoolSize"`
	MinPoolSize            uint64        `json:"minPoolSize" mapstructure:"minPoolSize"`
}

// NewMongoOptions new a mongo option.
func NewMongoOptions() *MongoOptions {
	return &MongoOptions{}
}

// Validate validate mongo options is valid.
func (o *MongoOptions) Validate() error {
	return nil
}

func (o *MongoOptions) Enabled() bool {
	return len(o.Host) > 0
}

// AddFlags ...
func (o *MongoOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Host, "mongo-host", "", "mongo db host")
	fs.StringVar(&o.Port, "mongo-port", MongoDefaultPort, "mongo db port")
	fs.StringVar(&o.Username, "mongo-username", "", "mongo db username")
	fs.StringVar(&o.Password, "mongo-password", "", "mongo db password")
	fs.StringVar(&o.Database, "mongo-database", MongoDefaultDatabase, "mongo database name")
	fs.DurationVar(&o.ConnectTimeout, "mongo-connect-timeout", MongoDefaultConnectTimeout, "mongo connect timeout")
	fs.DurationVar(&o.ServerSelectionTimeout, "mongo-server-selection-timeout", MongoDefaultServerSelectionTimeout, "mongo server selection timeout")
	fs.DurationVar(&o.SocketTimeout, "mongo-socket-timeout", MongoDefaultSocketTimeout, "mongo socket timeout")
	fs.DurationVar(&o.HeartbeatInterval, "mongo-heartbeat-interval", MongoDefaultHeartbeatInterval, "mongo heartbeat interval")
	fs.DurationVar(&o.MaxConnIdleTime, "mongo-conn-idle-timeout", MongoDefaultMaxConnIdleTime, "mongo connect idle timeout")
	fs.Uint64Var(&o.MaxPoolSize, "mongo-max-pool-size", MongoDefaultMaxPoolSize, "mongo max pool size")
	fs.Uint64Var(&o.MinPoolSize, "mongo-min-pool-size", MongoDefaultMinPoolSize, "mongo min pool size")
}

// GetDBInstance ...
func (o *MongoOptions) GetDBInstance(ctx context.Context) (*mongo.Client, *mongo.Database, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	mongoURL := fmt.Sprintf("mongodb://%s:%s/?authSource=admin",
		viper.GetString(o.Host),
		viper.GetString(o.Port),
	)
	if o.Username != "" && o.Password != "" {
		mongoURL = fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin",
			viper.GetString(o.Username),
			viper.GetString(o.Password),
			viper.GetString(o.Host),
			viper.GetString(o.Port),
		)
	}

	// set connect pool
	clientOpts := options.Client().ApplyURI(mongoURL)
	clientOpts.SetConnectTimeout(o.ConnectTimeout)
	clientOpts.SetServerSelectionTimeout(o.ServerSelectionTimeout)
	clientOpts.SetSocketTimeout(o.SocketTimeout)
	clientOpts.SetMaxPoolSize(o.MaxPoolSize)
	clientOpts.SetMinPoolSize(o.MinPoolSize)
	clientOpts.SetHeartbeatInterval(o.HeartbeatInterval)
	clientOpts.SetMaxConnIdleTime(o.MaxConnIdleTime)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, nil, fmt.Errorf("mongo open fail: %w", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, nil, err
	}

	return client, client.Database(viper.GetString(o.Database)), nil
}
