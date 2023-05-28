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
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var defaultLogger Logger

type Logger interface {
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Panicf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})

	Debugw(msg string, keysAndValues ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Panicw(msg string, keysAndValues ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})

	Sync()
}

func Debugf(template string, args ...interface{}) {
	defaultLogger.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	defaultLogger.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	defaultLogger.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	defaultLogger.Errorf(template, args...)
}

func Panicf(template string, args ...interface{}) {
	defaultLogger.Panicf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	defaultLogger.Fatalf(template, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Debugw(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	defaultLogger.Infow(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Warnw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Errorw(msg, keysAndValues...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Panicw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Fatalw(msg, keysAndValues...)
}

// Sync ...
func Sync() {
	defaultLogger.Sync()
}

var registerLogger sync.Once

// RegisterLogger register global logger.
func RegisterLogger(options *Options) {
	registerLogger.Do(func() {
		defaultLogger = NewLogger(options)
	})
}

// NewLogger new a logger.
func NewLogger(opts *Options) Logger {
	if opts == nil {
		opts = NewOptions()
	}

	cores := getZapCores(opts)

	zapTee := zapcore.NewTee(cores...)
	logger := zap.New(zapTee, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &zapLogger{logger}
}

func getZapCores(opts *Options) []zapcore.Core {
	zapLevel, err := zapcore.ParseLevel(opts.Level)
	if err != nil {
		zapLevel = zapcore.InfoLevel
	}
	encoderConfig := getEncoderConfig(opts)
	zapCores := []zapcore.Core{zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(os.Stdout), zapLevel)}
	if opts.OutputPath != "" {
		writer := &lumberjack.Logger{
			Filename:   opts.OutputPath, // make sure logfile parent dir exist
			MaxSize:    opts.MaxSize,
			MaxBackups: opts.MaxBackups,
			MaxAge:     opts.MaxAge,
			Compress:   opts.Compress,
		}
		zapCores = append(zapCores, zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(writer), zapLevel))
	}
	return zapCores
}

func getEncoderConfig(opts *Options) zapcore.EncoderConfig {
	encoderConfig := zap.NewProductionEncoderConfig()
	if opts.MessageKey != "" {
		encoderConfig.MessageKey = opts.MessageKey
	}
	if opts.LevelKey != "" {
		encoderConfig.LevelKey = opts.LevelKey
	}
	if opts.CallerKey != "" {
		encoderConfig.CallerKey = opts.CallerKey
	}
	if opts.TimeKey != "" {
		encoderConfig.TimeKey = opts.TimeKey
	}
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	return encoderConfig
}

type zapLogger struct {
	*zap.Logger
}

// Debugf ...
func (z *zapLogger) Debugf(template string, args ...interface{}) {
	z.Sugar().Debugf(template, args...)
}

// Infof ...
func (z *zapLogger) Infof(template string, args ...interface{}) {
	z.Sugar().Infof(template, args...)
}

// Warnf ...
func (z *zapLogger) Warnf(template string, args ...interface{}) {
	z.Sugar().Warnf(template, args...)
}

// Errorf ...
func (z *zapLogger) Errorf(template string, args ...interface{}) {
	z.Sugar().Errorf(template, args...)
}

// Panicf ...
func (z *zapLogger) Panicf(template string, args ...interface{}) {
	z.Sugar().Panicf(template, args...)
}

// Fatalf ...
func (z *zapLogger) Fatalf(template string, args ...interface{}) {
	z.Sugar().Fatalf(template, args...)
}

// Debug ...
func (z *zapLogger) Debugw(msg string, keysAndValues ...interface{}) {
	z.Sugar().Debugw(msg, keysAndValues...)
}

// Info ...
func (z *zapLogger) Infow(msg string, keysAndValues ...interface{}) {
	z.Sugar().Infow(msg, keysAndValues...)
}

// Warn ...
func (z *zapLogger) Warnw(msg string, keysAndValues ...interface{}) {
	z.Sugar().Warnw(msg, keysAndValues...)
}

// Error ...
func (z *zapLogger) Errorw(msg string, keysAndValues ...interface{}) {
	z.Sugar().Errorw(msg, keysAndValues...)
}

// Panic ...
func (z *zapLogger) Panicw(msg string, keysAndValues ...interface{}) {
	z.Sugar().Panicw(msg, keysAndValues...)
}

// Fatal ...
func (z *zapLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	z.Sugar().Fatalw(msg, keysAndValues...)
}

// Sync ...
func (z *zapLogger) Sync() {
	_ = z.Sugar().Sync()
}

var _ Logger = &zapLogger{}
