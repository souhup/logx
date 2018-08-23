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
	"runtime"
	"sync/atomic"
	"testing"
	"time"
	"sync"
)

func BenchmarkParallelLogger_Debug(b *testing.B) {
	return
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			X.Debug(1)
		}
	})
	b.StopTimer()
}

func TestPerformance(t *testing.T) {
	fmt.Println("test performance")
	Init("./config/logs.yml")
	var ops int64 = 0
	for i := 0; i < 100000; i++ {
		go func() {
			for {
				X.Add("goid", runtime.Goid())
				X.Info(0)
				X.Clean()
				atomic.AddInt64(&ops, 1)
			}
		}()
	}
	fmt.Println("create goroutine done.")
	time.Sleep(1 * time.Second)
	score := atomic.LoadInt64(&ops)
	fmt.Printf("number of execution in 1s: %d\n", score)
	return
}

func TestInit(t *testing.T) {
	// test error
	Init("./config/err_test.yml")
	// test error
	Init("./foo/bar")

	Init("./config/logs.yml")
	X.Debug("test Init")
}

func TestGetLogger(t *testing.T) {
	logger, err := GetLogger("./config/logs.yml")
	if err != nil {
		panic(err)
	}
	logger.Debug("test GetLogger")
}

func TestGetLoggerByConf(t *testing.T) {
	conf := Config{
		MessageKey: "msg",
		LevelKey:   "level",
		TimeKey:    "time",
		Encoding:   "console",
		CallerKey:  "caller",
		Level:      -1,
	}
	logger, err := GetLoggerByConf(&conf)
	if err != nil {
		t.Fatal(err)
	}
	logger.Debug("test GetLogger")

	// test error
	conf.Encoding = "foo"
	_, err = GetLoggerByConf(&conf)
	if err != nil {
		logger.Debug("test GetLogger")
	}
}

func TestLogger_Show(t *testing.T) {
	type Book struct {
		ID   int
		Name string
	}
	book := &Book{
		ID:   1,
		Name: "test Show",
	}
	X.Show(book)
}

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

func TestLogger_With(t *testing.T) {
	X.With("foo", "bar", "foo2", 0).With("foo3", 1).Debug("test With")
	X.Debug("test With")
}

func TestLogger_Clean(t *testing.T) {
	X.Add("foo", "bar")
	X.Debug("test Clean")
	X.Clean()
	X.Debug("test Clean")
}

func TestLogger_Flush(t *testing.T) {
	defer X.Flush()
	X.Debug("test Flush")
}

func TestLogger_Debug(t *testing.T) {
	X.Debug("test Debug")
}

func TestLogger_Debugf(t *testing.T) {
	X.Debugf("test %s", "Debugf")
}

func TestLogger_Info(t *testing.T) {
	X.Info("test Info")
}

func TestLogger_Infof(t *testing.T) {
	X.Infof("test %s", "Infof")
}

func TestLogger_Warn(t *testing.T) {
	X.Warn("test Warn")
}

func TestLogger_Warnf(t *testing.T) {
	X.Warnf("test %s", "Warnf")
}

func TestLogger_Error(t *testing.T) {
	X.Error("test Error")
}

func TestLogger_Errorf(t *testing.T) {
	X.Errorf("test %s", "Errorf")
}
