## Logx

powerful and simple logger for golang

## Installation

`go get -u github.com/go/logx`

## Quick Start

### log 

You can use it directly.
```go
package main

import "github.com/go/logx"

func main() {
	logx.X.Debug("hi")
	logx.X.Debug("hello", "world", 42)
	logx.X.Debugf("%s", "world")
}
```
```
{"level":"DEBUG","time":"2019-09-21 16:47:19","caller":"test/main.go:6","msg":"hi"}
{"level":"DEBUG","time":"2019-09-21 16:47:19","caller":"test/main.go:7","msg":"hello world 42"}
{"level":"DEBUG","time":"2019-09-21 16:47:19","caller":"test/main.go:8","msg":"world"}
```

Key-value pairs is Ok.
```go
func main() {
	logxX.With("a", "1", "b", 2).
		With("c", 3).
		With("d", 5).
		Debug("test With")
}
```

```
{"level":"DEBUG","time":"2019-09-21 16:49:11","caller":"test/main.go:9","msg":"test With","a":"1","b":2,"c":3,"d":5}
```

More useful function is that, it can store entry in context 
```go
func main() {
	ctx := logx.X.Withc(nil, "a", 1, "b", 2)
	ctx = logx.X.Withc(ctx, "c", 3)
	logx.X.Debugc(ctx, "test Withc")
}
```

```
{"level":"DEBUG","time":"2019-09-21 16:51:17","caller":"test/main.go:10","msg":"test Withc","a":1,"b":2,"c":3}
```

### Init

of course, just printing on the console does not meet our needs. We can write logs to files.

```go
func main() {
	logx.Init("./config/logs.yml")
	logx.X.Debug("hi")
}
```

the configuration file as the following.

```yaml
# Level is the minimum log level
# -1 is debug level,
# 0 is info level,
# 1 is warn level,
# is error level
level: -1

# keys used for each log entry. If any key is empty, that portion
# of the entry is omitted.
message_key: msg
level_key: level
time_key: time
caller_key: caller

# encoding of log, just is json or console
encoding: json
# Filename is the file to write logs to.  Backup log files will be retained
# in the same directory.
file_name: "logs/test.log"
# MaxSize is the maximum size in megabytes of the log file before it gets
# rotated. It defaults to 100 megabytes.
max_size: 1
# MaxAge is the maximum number of days to retain old log files based on the
# timestamp encoded in their filename.  Note that a day is defined as 24
# hours and may not exactly correspond to calendar days due to daylight
# savings, leap seconds, etc. The default is not to remove old log files
# based on age.
max_age: 7
# MaxBackups is the maximum number of old log files to retain. The default
# is to retain all old log files (though MaxAge may still cause them to get
# deleted.)
max_backups: 50
# LocalTime determines if the time used for formatting the timestamps in
# backup files is the computer's local time.  The default is to use UTC
# time.
local_time: true
# Compress determines if the rotated log files should be compressed
# using gzip.
compress: false
```
