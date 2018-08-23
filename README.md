## Logx

powerful and simple logger for golang.

## Installation

`go get -u github.com/souhup/logx`

Note that logs need a method to get goroutine id.

My suggestion is modifying source code. 

You can add the following method to runtime/proc.go, and then compile source code.

```go
func Goid() int64 {
    _g_ := getg()
    return _g_.goid
}
```

## Quick Start

### log 

You can use it directly.
```go
import "github.com/souhup/logx"

func main() {
	logx.X.Debug("hi")
	logx.X.Debugf("%s", "hello")
}
```

```json
{"level":"DEBUG","time":"2018-08-23 15:49:35","caller":"example/main.go:6","msg":"hi"}
{"level":"DEBUG","time":"2018-08-23 15:49:35","caller":"example/main.go:7","msg":"hello"}
```

Key-value pairs is Ok.

```go
func main() {
	logx.X.With("Years", 1, "Name", "souhup").Debug("hi")
}
```

```json
{"level":"DEBUG","time":"2018-08-23 15:53:12","caller":"example/main.go:9","msg":"hi","Years":1,"Name":"souhup"}
```

More useful function is that, it can cache entry in every Goroutine.
```go
func TestLogger_Add(t *testing.T) {
	group := sync.WaitGroup{}
	group.Add(1)
	go func() {
		X.Add("foo1", 1)
		X.Debug("test Add 1")
		group.Done()
	}()
	group.Add(1)
	go func() {
		X.Add("foo2", 2)
		X.Debug("test Add 2")
		group.Done()
	}()
	group.Wait()
}
```

```json
{"level":"DEBUG","time":"2018-08-23 16:03:25","caller":"example/main.go:19","msg":"test Add 2","foo2":2}
{"level":"DEBUG","time":"2018-08-23 16:03:25","caller":"example/main.go:13","msg":"test Add 1","foo1":1}
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



