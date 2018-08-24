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
	"runtime"
	"sync"
	"encoding/json"
)

// X is a instance of Logger, and it will be initialized when logx is imported.
var X *Logger

// field store every entries in different goroutines.
var field sync.Map

// Config is struct about configuration file.
type Config struct {
	// Level is the minimum log level
	// -1 is debug level,
	// 0 is info level,
	// 1 is warn level,
	// 2 is error level
	Level int8 `yaml:"level"`

	// keys used for each log entry. If any key is empty, that portion
	// of the entry is omitted.
	MessageKey string `yaml:"message_key"` // key of message
	LevelKey   string `yaml:"level_key"`   // key of level
	TimeKey    string `yaml:"time_key"`    // key of time
	CallerKey  string `yaml:"caller_key"`  // key of caller

	// encoding of log, just is json or console
	Encoding string `yaml:"encoding"`

	// Filename is the file to write logs to.  Backup log files will be retained
	// in the same directory.
	Filename string `yaml:"file_name"`

	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	MaxSize int `yaml:"max_size"`

	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxAge int `yaml:"max_age"`

	// MaxBackups is the maximum number of old log files to retain. The default
	// is to retain all old log files (though MaxAge may still cause them to get
	// deleted.)
	MaxBackups int `yaml:"max_backups"`

	// LocalTime determines if the time used for formatting the timestamps in
	// backup files is the computer's local time.  The default is to use UTC
	// time.
	LocalTime bool `yaml:"local_time"`

	// Compress determines if the rotated log files should be compressed
	// using gzip.
	Compress bool `yaml:"compress"`
}

// Log is a logger interface. It contains all API about logx.
type Log interface {
	Flush()
	Clean()
	Show(interface{})
	With(...interface{}) *Logger
	Add(string, interface{})
	Debug(interface{})
	Debugf(string, ...interface{})
	Info(interface{})
	Infof(string, ...interface{})
	Warn(interface{})
	Warnf(string, ...interface{})
	Error(interface{})
	Errorf(string, ...interface{})
}

// Logger is the implement about Log.
type Logger struct {
	zapLogger *zap.Logger
	sugar     *zap.SugaredLogger
}

// Flush calls the underlying Core's Sync method, flushing any buffered log
// entries. Applications should take care to call Sync before exiting.
func (l *Logger) Flush() {
	l.zapLogger.Sync()
}

// Add adds a entry in current goroutine. Entry will be print when logging.
func (l *Logger) Add(key string, value interface{}) {
	goid := runtime.Goid()
	actual, _ := field.LoadOrStore(goid, make([]interface{}, 0))
	arr := append(actual.([]interface{}), key, value)
	field.Store(goid, arr)
}

// Clean all entries in current goroutine.
func (l *Logger) Clean() {
	goid := runtime.Goid()
	field.Delete(goid)
}

// With adds entries and constructs a new Logger.
// Note that the keys in key-value pairs should be strings.
func (l *Logger) With(keysAndValues ...interface{}) (log *Logger){
	log = new(Logger)
	log.sugar = l.sugar.With(keysAndValues...)
	return
}

// Show prints value by JSON format on stdout.
func (l *Logger) Show(value interface{}){
	b, _ := json.MarshalIndent(value, "", "  ")
	msg := fmt.Sprintf("%s", string(b))
	fmt.Println(msg)
}

// Debug uses fmt.Sprint to construct and logs a message.
// If value is struct, it will be converted to JSON.
func (l *Logger) Debug(v interface{}) {
	generate(l.sugar.Debugw, v)
}

// Debugf uses fmt.Sprintf to log a templated message.
func (l *Logger) Debugf(format string, params ...interface{}) {
	generate(l.sugar.Debugw, format, params...)
}

// Info uses fmt.Sprint to construct and log a message.
// If value is struct, it will be converted to JSON.
func (l *Logger) Info(v interface{}) {
	generate(l.sugar.Infow, v)
}

// Infof uses fmt.Sprintf to log a templated message.
func (l *Logger) Infof(format string, params ...interface{}) {
	generate(l.sugar.Infow, format, params...)
}

// Warn uses fmt.Sprint to construct and log a message.
// If value is struct, it will be converted to JSON.
func (l *Logger) Warn(v interface{}) {
	generate(l.sugar.Warnw, v)
}

// Warnf uses fmt.Sprintf to log a templated message.
func (l *Logger) Warnf(format string, params ...interface{}) {
	generate(l.sugar.Warnw, format, params...)
}

// Error uses fmt.Sprint to construct and log a message.
// If value is struct, it will be converted to JSON.
func (l *Logger) Error(v interface{}) {
	generate(l.sugar.Errorw, v)
}

// Errorf uses fmt.Sprintf to log a templated message.
func (l *Logger) Errorf(format string, params ...interface{}) {
	generate(l.sugar.Errorw, format, params...)
}

func generate(fun func(string, ...interface{}), format interface{}, params ...interface{}) {
	var msg string
	if len(params) > 0 {
		msg = fmt.Sprintf(format.(string), params...)
	} else {
		msg = fmt.Sprintf("%+v", format)
	}
	fun(msg, getField()...)
	return
}

func getField() []interface{} {
	goid := runtime.Goid()
	value, _ := field.LoadOrStore(goid, make([]interface{}, 0))
	return value.([]interface{})
}
