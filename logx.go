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
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

// init sets default format when logx is imported.
//
// The default format of log is "<time> <level> <caller> <message>".
func init() {
	conf := Config{
		MessageKey: "msg",
		LevelKey:   "level",
		TimeKey:    "time",
		Encoding:   "simple",
		CallerKey:  "caller",
		Level:      -1,
	}
	X, _ = GetLoggerByConf(&conf)
}

// Init is a high-level wrapper that takes a URL, open specified configuration file,
// and initializes X.
func Init(path string) (err error) {
	logger, err := GetLogger(path)
	if err != nil {
		return
	}
	X = logger
	return
}

// Init is a high-level wrapper that takes a URL, open specified configuration file,
// and generate a Logger.
func GetLogger(path string) (logger *Logger, err error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("read log configuration %v, error: %v", path, err)
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		err = fmt.Errorf("unmarshal log configuration %v, error: %v", path, err)
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	return GetLoggerByConf(&config)
}
