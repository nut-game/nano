package slog

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/nut-game/nano/logger/interfaces"
)

type slogImpl struct {
	slog.Logger
}

func New() interfaces.Logger {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	return &slogImpl{Logger: logger}
}

// 定义 Infof 方法
func (l *slogImpl) Info(args ...interface{}) {
	l.Info(fmt.Sprint(args...))
}
func (l *slogImpl) Infof(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

// 定义 Warnf 方法
func (l *slogImpl) Warn(args ...interface{}) {
	l.Warn(fmt.Sprint(args...))
}
func (l *slogImpl) Warnf(format string, args ...interface{}) {
	l.Warn(fmt.Sprintf(format, args...))
}

// 定义 Errorf 方法
func (l *slogImpl) Error(args ...interface{}) {
	l.Error(fmt.Sprint(args...))
}
func (l *slogImpl) Errorf(format string, args ...interface{}) {
	l.Error(fmt.Sprintf(format, args...))
}

// 其他 slog 的原始方法仍可直接使用
func (l *slogImpl) Debug(args ...interface{}) {
	l.Debug(fmt.Sprint(args...))
}
func (l *slogImpl) Debugf(format string, args ...interface{}) {
	l.Debug(fmt.Sprintf(format, args...))
}

// 其他 slog 的原始方法仍可直接使用
func (l *slogImpl) Fatal(args ...interface{}) {
	l.Error(fmt.Sprint(args...))
}
func (l *slogImpl) Fatalf(format string, args ...interface{}) {
	l.Error(fmt.Sprintf(format, args...))
}
