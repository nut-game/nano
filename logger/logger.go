// Copyright (c) TFG Co. All Rights Reserved.
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

package logger

import (
	"github.com/nut-game/nano/logger/interfaces"
	"github.com/nut-game/nano/logger/slog"
)

var (
	Fatal   func(v ...interface{})
	Fatalf  func(format string, v ...interface{})
	Fatalln func(v ...interface{})

	Debug   func(v ...interface{})
	Debugf  func(format string, v ...interface{})
	Debugln func(v ...interface{})

	Error   func(v ...interface{})
	Errorf  func(format string, v ...interface{})
	Errorln func(v ...interface{})

	Info   func(v ...interface{})
	Infof  func(format string, v ...interface{})
	Infoln func(v ...interface{})

	Warn   func(v ...interface{})
	Warnf  func(format string, v ...interface{})
	Warnln func(v ...interface{})
)

// Log is the default logger
var Log interfaces.Logger

func init() {
	logger := slog.New()

	SetLogger(logger)
}

// SetLogger rewrites the default logger
func SetLogger(logger interfaces.Logger) {
	if logger == nil {
		return
	}

	Log = logger

	Fatal = logger.Fatal
	Fatalf = logger.Fatalf
	Fatalln = logger.Fatal

	Debug = logger.Debug
	Debugf = logger.Debugf
	Debugln = logger.Debug

	Error = logger.Error
	Errorf = logger.Errorf
	Errorln = logger.Error

	Info = logger.Info
	Infof = logger.Infof
	Infoln = logger.Info

	Warn = logger.Warn
	Warnf = logger.Warnf
	Warnln = logger.Warn

}
