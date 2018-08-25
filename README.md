## Logx

powerful and simple logger for golang, and help you identify different goroutine automatically.

## Installation

`go get -u github.com/souhup/logx`

## Quick Start

### log 

You can use it directly.
```go
import "github.com/souhup/logx"

func main() {
	logx.X.Debug("hello")
	logx.X.Debugf("%s", "world")
}
```


```
018-08-25 08:47:33	DEBUG	igo/t5.go:5	hello
2018-08-25 08:47:33	DEBUG	igo/t5.go:6	world
```

Key-value pairs is Ok.

Note that the default mode is **simple**, it will ignore the keys. if you want 
show them, you can use **console** or **json** in configuration file.

```go
func main() {
	logx.X.With("Method", "Get", "DeviceID", "aptx4896").Debug("authenticate")
}
```

```
2018-08-25 08:48:40	DEBUG	Get	aptx4896	igo/t5.go:5	authenticate
```

More useful function is that, it can cache entry in every Goroutine.

```go
func main() {
	rand.Seed(time.Now().UnixNano())
	go business()
	go business()
	time.Sleep(time.Second)
}

func business() {
	// clean goroutine cache when you destroy a goroutine.
	defer logx.X.Clean()
	logx.X.Add("trace_id", rand.Int())
	logx.X.Debug("authenticate")
	logx.X.Debug("done")
}
```

```json
func main() {
	rand.Seed(time.Now().UnixNano())
	go business()
	go business()
	time.Sleep(time.Second)
}

func business() {
	// clean goroutine cache when you destroy a goroutine.
	defer logx.X.Clean()
	logx.X.Add("trace_id", rand.Int())
	logx.X.Debug("authenticate")
	logx.X.Debug("done")
}
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



