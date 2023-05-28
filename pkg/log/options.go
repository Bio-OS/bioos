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

package log

import (
	"github.com/spf13/pflag"
)

type Options struct {
	Level      string `json:"level" mapstructure:"level"`
	OutputPath string `json:"output-path" mapstructure:"output-path,omitempty"`
	MaxSize    int    `json:"max-size" mapstructure:"max-size,omitempty"`
	MaxBackups int    `json:"max-backups" mapstructure:"max-backups,omitempty"`
	MaxAge     int    `json:"max-age" mapstructure:"max-age,omitempty"`
	Compress   bool   `json:"compress" mapstructure:"compress,omitempty"`
	MessageKey string `json:"message-key" mapstructure:"message-key,omitempty"`
	LevelKey   string `json:"level-key" mapstructure:"level-key,omitempty"`
	CallerKey  string `json:"caller-key" mapstructure:"caller-key,omitempty"`
	TimeKey    string `json:"time-key" mapstructure:"time-key,omitempty"`
}

// NewOptions new a log option.
func NewOptions() *Options {
	return &Options{
		Level:      "info",
		OutputPath: "app.log",
		MaxSize:    100,
		MaxBackups: 5,
		MaxAge:     1,
		Compress:   true,
	}
}

// Validate validate log options is valid.
func (o *Options) Validate() error {
	// if o.OutputPath == "" {
	//	return errors.New("log output empty")
	//}
	return nil
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&o.Level, "log-level", "", "info", "log level")
	fs.StringVarP(&o.OutputPath, "log-path", "", "", "log path")
	fs.IntVarP(&o.MaxSize, "log-max-size", "", 100, "the maximum size in megabytes of the log file")
	fs.IntVarP(&o.MaxBackups, "log-max-backups", "", 5, "the maximum number of old log files to retain")
	fs.IntVarP(&o.MaxAge, "log-max-age", "", 1, "log maximum number of days")
	fs.BoolVarP(&o.Compress, "log-compress", "", true, "log compress")
	fs.StringVarP(&o.MessageKey, "log-message-key", "", "msg", "log message key")
	fs.StringVarP(&o.LevelKey, "log-level-key", "", "level", "log level key")
	fs.StringVarP(&o.CallerKey, "log-caller-key", "", "caller", "log caller key")
	fs.StringVarP(&o.TimeKey, "log-time-key", "", "", "log time key")
}
