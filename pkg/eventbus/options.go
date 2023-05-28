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

package eventbus

import (
	"time"

	"github.com/spf13/pflag"
)

const (
	DefaultWorkers        = 5
	DefaultBatchSize      = 100
	DefaultSyncPeriod     = time.Second * 30
	DefaultMaxRetries     = 10
	DefaultDequeueTimeout = time.Minute * 5
	DefaultRunningTimeout = time.Hour * 24 * 365 // 1year
)

type Options struct {
	MaxRetries     int           `json:"maxRetries" mapstructure:"maxRetries"`
	SyncPeriod     time.Duration `json:"syncPeriod" mapstructure:"syncPeriod"`
	BatchSize      int           `json:"batchSize" mapstructure:"batchSize"`
	Workers        int           `json:"workers" mapstructure:"workers"`
	DequeueTimeout time.Duration `json:"dequeueTimeout" mapstructure:"dequeueTimeout"`
	RunningTimeout time.Duration `json:"runningTimeout" mapstructure:"runningTimeout"`
}

// NewOptions new an event bus option.
func NewOptions() *Options {
	return &Options{}
}

// Validate validate log options is valid.
func (o *Options) Validate() error {
	return nil
}

// Enabled check if event bus is enabled
func (o *Options) Enabled() bool {
	return true
}

// AddFlags add event bus flags
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.IntVar(&o.MaxRetries, "event-bus-max-retries", DefaultMaxRetries, "event max try times")
	fs.DurationVar(&o.SyncPeriod, "event-bus-sync-period", DefaultSyncPeriod, "sync period in seconds")
	fs.IntVar(&o.BatchSize, "event-bus-batch-size", DefaultBatchSize, "batch size to get events")
	fs.IntVar(&o.Workers, "event-bus-workers", DefaultWorkers, "concurrent workers")
	fs.DurationVar(&o.DequeueTimeout, "event-bus-dequeue-timeout", DefaultDequeueTimeout, "dequeue timeout")
	fs.DurationVar(&o.RunningTimeout, "event-bus-running-timeout", DefaultRunningTimeout, "running timeout")
}
