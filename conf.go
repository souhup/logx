// Copyright (c) 2018 souhup
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package logx

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
	"errors"
)

// GetLoggerByConf constructs a new Logger by Config.
func GetLoggerByConf(config *Config) (logger *Logger, err error) {
	proConf := zapcore.EncoderConfig{
		MessageKey:     config.MessageKey,
		LevelKey:       config.LevelKey,
		TimeKey:        config.TimeKey,
		CallerKey:      config.CallerKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// choose the type of encoding.
	var encoder zapcore.Encoder
	if config.Encoding == "json" {
		encoder = zapcore.NewJSONEncoder(proConf)
	} else if config.Encoding == "console" {
		encoder = zapcore.NewConsoleEncoder(proConf)
	} else {
		err = errors.New("encoding must be one of the json or console")
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	// writer logs to rolling files
	zapWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   config.Filename,
		MaxSize:    config.MaxSize,
		MaxAge:     config.MaxAge,
		MaxBackups: config.MaxBackups,
		LocalTime:  config.LocalTime,
	})
	output := zapcore.NewMultiWriteSyncer(zapWriter)
	if config.Filename == "" {
		output = os.Stdout
	}


	newCore := zapcore.NewCore(encoder, output, zap.NewAtomicLevelAt(zapcore.Level(config.Level)))
	opts := []zap.Option{zap.ErrorOutput(zapWriter)}
	opts = append(opts, zap.AddCaller(), zap.AddCallerSkip(3))

	logger = new(Logger)
	logger.zapLogger = zap.New(newCore, opts...)
	logger.sugar = logger.zapLogger.Sugar()
	return
}

// timeEncoder sets the time format
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}
