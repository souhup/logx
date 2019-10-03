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
	. "gopkg.in/check.v1"
	"reflect"
	"runtime"
	"testing"
)

type MySuite struct {
}

var _ = Suite(&MySuite{})

func Test(t *testing.T) { TestingT(t) }

func (it *MySuite) SetUpSuite(c *C) {
}

func (it *MySuite) TearDownSuite(c *C) {
}

func (it *MySuite) SetUpTest(c *C) {
}

func (it *MySuite) TearDownTest(c *C) {
}

func (it *MySuite) TestGetLogger(c *C) {
	logger, err := GetLogger("./config/logs.yml")
	c.Assert(err, IsNil)
	logger.Info("test GetLogger")
}

func (it *MySuite) TestGetLoggerByConf(c *C) {
	conf := Config{
		MessageKey: "msg",
		LevelKey:   "level",
		TimeKey:    "time",
		Encoding:   "json",
		CallerKey:  "caller",
		Level:      -1,
	}
	logger, err := GetLoggerByConf(&conf)
	c.Assert(err, IsNil)
	logger.Info("test GetLogger")
}

func (it *MySuite) TestWith(c *C) {
	X.With("a", "1", "b", 2).
		With("c", 3).
		With("d", 5).
		Debug("test With")
}

func (it *MySuite) TestWithc(c *C) {
	ctx := context.TODO()
	ctx = X.Withc(ctx, "a", 1, "b", 2)
	ctx = X.Withc(ctx, "c", 3)
	X.Debugc(ctx, "test Withc")
}

func (it *MySuite) TestWithcf(c *C) {
	ctx := context.TODO()
	ctx = X.Withc(ctx, "a", 1, "b", 2)
	ctx = X.Withcf(ctx, "tag", "key%v", 13)
	X.Debugc(ctx, "test Withcf")
}

func (it *MySuite) TestFlush(c *C) {
	defer X.Flush()
	X.Debug("test Flush")
}

func (it *MySuite) TestPrint(c *C) {
	funArr := []func(...interface{}){X.Debug, X.Info, X.Warn, X.Error}
	for _, fun := range funArr {
		fun("testing", getFunctionName(fun))
	}
}

func (it *MySuite) TestPrintf(c *C) {
	funArr := []func(string, ...interface{}){X.Debugf, X.Infof, X.Warnf, X.Errorf}
	for _, fun := range funArr {
		fun("testing %v", getFunctionName(fun))
	}
}

func (it *MySuite) TestPrintc(c *C) {
	ctx := X.Withc(nil, "a", 1, "b", 2)
	ctx = X.Withc(ctx, "c", 3)
	funArr := []func(context.Context, ...interface{}){X.Debugc, X.Infoc, X.Warnc, X.Errorc}
	for _, fun := range funArr {
		fun(ctx, "testing", getFunctionName(fun))
	}
}

func (it *MySuite) TestPrintcf(c *C) {
	ctx := X.Withc(nil, "a", 1, "b", 2)
	ctx = X.Withc(ctx, "c", 3)
	funArr := []func(context.Context, string, ...interface{}){X.Debugcf, X.Infocf, X.Warncf, X.Errorcf}
	for _, fun := range funArr {
		fun(ctx, "testing %v", getFunctionName(fun))
	}
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
