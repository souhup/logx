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
	"context"
	"encoding/json"
	"fmt"
)

// X is a instance of Logger, and it will be initialized when logx is imported.
var X *Logger

const contextLogKey = "_logx"

type method uint8

const (
	Debug method = iota
	Info
	Warn
	Error
	Fatal
	Panic
)

// Show prints value by JSON format on stdout.
func (it *Logger) Show(value interface{}) {
	b, _ := json.MarshalIndent(value, "", "  ")
	msg := fmt.Sprintf("%s", string(b))
	fmt.Println(msg)
}

// Flush calls the underlying Core's Sync method, flushing any buffered log
// entries. Applications should take care to call Sync before exiting.
func (it *Logger) Flush() {
	it.zapLogger.Sync()
}

// With adds entries and constructs a new Logger.
// Note that the keys in key-value pairs should be strings.
func (it *Logger) With(keysAndValues ...interface{}) (log *Logger) {
	log = new(Logger)
	log.sugar = it.sugar.With(keysAndValues...)
	return
}

// Withc adds entries and constructs a new Logger, and uses fmt.Sprintf to store a templated message.
func (it *Logger) Withf(key string, format string, params ...interface{}) (log *Logger) {
	log = new(Logger)
	log.sugar = it.sugar.With(key, fmt.Sprintf(format, params...))
	return
}

//  same as With, but store in context
func (it *Logger) Withc(ctx context.Context, keysAndValues ...interface{}) context.Context {
	if ctx == nil {
		ctx = context.TODO()
	}
	value := ctx.Value(contextLogKey)
	var log *Logger
	switch value.(type) {
	case *Logger:
		log = value.(*Logger)
		log.sugar = log.sugar.With(keysAndValues...)
	default:
		log = new(Logger)
		log.sugar = it.sugar.With(keysAndValues...)
	}
	return context.WithValue(ctx, contextLogKey, log)
}

//  same as Withf, but store in context
func (it *Logger) Withcf(ctx context.Context, key string, format string, params ...interface{}) context.Context {
	if ctx == nil {
		ctx = context.TODO()
	}
	value := ctx.Value(contextLogKey)
	var log *Logger
	switch value.(type) {
	case *Logger:
		log = value.(*Logger)
		log.sugar = log.sugar.With(key, fmt.Sprintf(format, params...))
	default:
		log = new(Logger)
		log.sugar = it.sugar.With(key, fmt.Sprintf(format, params...))
	}
	return context.WithValue(ctx, contextLogKey, log)
}

// Debug uses fmt.Sprint to construct and logs a message.
// If value is struct, it will be converted to JSON.
func (it *Logger) Debug(v ...interface{}) {
	generate(nil, it, Debug, "", v...)
}

// Debugf uses fmt.Sprintf to log a templated message.
func (it *Logger) Debugf(format string, params ...interface{}) {
	generate(nil, it, Debug, format, params...)
}

// same as Debug, but print with content in context
func (it *Logger) Debugc(ctx context.Context, v ...interface{}) {
	generate(ctx, it, Debug, "", v...)
}

// same as Debugf, but print with content in context
func (it *Logger) Debugcf(ctx context.Context, format string, params ...interface{}) {
	generate(ctx, it, Debug, format, params...)
}

// Info uses fmt.Sprint to construct and log a message.
// If value is struct, it will be converted to JSON.
func (it *Logger) Info(v ...interface{}) {
	generate(nil, it, Info, "", v...)
}

// Infof uses fmt.Sprintf to log a templated message.
func (it *Logger) Infof(format string, params ...interface{}) {
	generate(nil, it, Info, format, params...)
}

// same as Info, but print with content in context
func (it *Logger) Infoc(ctx context.Context, v ...interface{}) {
	generate(ctx, it, Info, "", v...)
}

// same as Infof, but print with content in context
func (it *Logger) Infocf(ctx context.Context, format string, params ...interface{}) {
	generate(ctx, it, Info, format, params...)
}

// Warn uses fmt.Sprint to construct and log a message.
// If value is struct, it will be converted to JSON.
func (it *Logger) Warn(v ...interface{}) {
	generate(nil, it, Warn, "", v...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func (it *Logger) Warnf(format string, params ...interface{}) {
	generate(nil, it, Warn, format, params...)
}

// same as Warn, but print with content in context
func (it *Logger) Warnc(ctx context.Context, v ...interface{}) {
	generate(ctx, it, Warn, "", v...)
}

// same as Warnf, but print with content in context
func (it *Logger) Warncf(ctx context.Context, format string, params ...interface{}) {
	generate(ctx, it, Warn, format, params...)
}

// Error uses fmt.Sprint to construct and log a message.
// If value is struct, it will be converted to JSON.
func (it *Logger) Error(v ...interface{}) {
	generate(nil, it, Error, "", v...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func (it *Logger) Errorf(format string, params ...interface{}) {
	generate(nil, it, Error, format, params...)
}

// same as Error, but print with content in context
func (it *Logger) Errorc(ctx context.Context, v ...interface{}) {
	generate(ctx, it, Error, "", v...)
}

// same as Errorf, but print with content in context
func (it *Logger) Errorcf(ctx context.Context, format string, params ...interface{}) {
	generate(ctx, it, Error, format, params)
}

// Fatal uses fmt.Sprint to construct and log a message.
// If value is struct, it will be converted to JSON.
func (it *Logger) Fatal(v ...interface{}) {
	generate(nil, it, Fatal, "", v...)
}

// Fatalf uses fmt.Sprintf to log a templated message.
func (it *Logger) Fatalf(format string, params ...interface{}) {
	generate(nil, it, Fatal, format, params...)
}

// same as Fatal, but print with content in context
func (it *Logger) Fatalc(ctx context.Context, v ...interface{}) {
	generate(ctx, it, Fatal, "", v...)
}

// same as Fatalf, but print with content in context
func (it *Logger) Fatalcf(ctx context.Context, format string, params ...interface{}) {
	generate(ctx, it, Fatal, format, params...)
}

// Panic uses fmt.Sprint to construct and log a message.
// If value is struct, it will be converted to JSON.
func (it *Logger) Panic(v ...interface{}) {
	generate(nil, it, Panic, "", v...)
}

// Panicf uses fmt.Sprintf to log a templated message.
func (it *Logger) Panicf(format string, params ...interface{}) {
	generate(nil, it, Panic, format, params...)
}

// same as Panic, but print with content in context
func (it *Logger) Panicc(ctx context.Context, v ...interface{}) {
	generate(ctx, it, Panic, "", v...)
}

// same as Panicf, but print with content in context
func (it *Logger) Paniccf(ctx context.Context, format string, params ...interface{}) {
	generate(ctx, it, Panic, format, params...)
}

func generate(ctx context.Context, self *Logger, fun method, format interface{}, params ...interface{}) {

	var msg string
	if len(format.(string)) > 0 {
		msg = fmt.Sprintf(format.(string), params...)
	} else {
		for i, param := range params {
			if i == 0 {
				msg += fmt.Sprintf("%+v", param)
			} else {
				msg += fmt.Sprintf(" %+v", param)
			}
		}
	}
	var basicLog *Logger
	if ctx != nil {
		value := ctx.Value(contextLogKey)
		switch value.(type) {
		case *Logger:
			basicLog = value.(*Logger)
		default:
			basicLog = self
		}
	} else {
		basicLog = self
	}

	switch fun {
	case Debug:
		basicLog.sugar.Debug(msg)
	case Info:
		basicLog.sugar.Info(msg)
	case Warn:
		basicLog.sugar.Warn(msg)
	case Error:
		basicLog.sugar.Error(msg)
	case Fatal:
		basicLog.sugar.Fatal(msg)
	case Panic:
		basicLog.sugar.Panic(msg)
	}
	return
}
