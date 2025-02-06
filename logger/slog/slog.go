package slog

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/nut-game/nano/logger/interfaces"
)

type slogImpl struct {
	*slog.Logger
}

// Debugln implements interfaces.Logger.
func (l *slogImpl) Debugln(args ...interface{}) {
	l.Logger.Debug(fmt.Sprintln(args...))
}

// Errorln implements interfaces.Logger.
func (l *slogImpl) Errorln(args ...interface{}) {
	l.Logger.Error(fmt.Sprintln(args...))
}

// Fatalln implements interfaces.Logger.
func (l *slogImpl) Fatalln(args ...interface{}) {
	l.Fatal(fmt.Sprintln(args...))
}

// Infoln implements interfaces.Logger.
func (l *slogImpl) Infoln(args ...interface{}) {
	l.Logger.Info(fmt.Sprintln(args...))
}

// Warnln implements interfaces.Logger.
func (l *slogImpl) Warnln(args ...interface{}) {
	l.Logger.Warn(fmt.Sprintln(args...))
}

func New() interfaces.Logger {
	level := slog.LevelInfo
	// todo 目前没有处理With相关的，加了也没有用。
	logger := slog.New(NewPrettyHandler(os.Stdout, level)).With("a", "b")
	return &slogImpl{Logger: logger}
}

// 定义 Infof 方法
func (l *slogImpl) Info(args ...interface{}) {
	l.Logger.Info(fmt.Sprint(args...))
}
func (l *slogImpl) Infof(format string, args ...interface{}) {
	l.Logger.Info(fmt.Sprintf(format, args...))
}

// 定义 Warnf 方法
func (l *slogImpl) Warn(args ...interface{}) {
	l.Logger.Warn(fmt.Sprint(args...))
}
func (l *slogImpl) Warnf(format string, args ...interface{}) {
	l.Logger.Warn(fmt.Sprintf(format, args...))
}

// 定义 Errorf 方法
func (l *slogImpl) Error(args ...interface{}) {
	l.Logger.Error(fmt.Sprint(args...))
}
func (l *slogImpl) Errorf(format string, args ...interface{}) {
	l.Logger.Error(fmt.Sprintf(format, args...))
}

// 其他 slog 的原始方法仍可直接使用
func (l *slogImpl) Debug(args ...interface{}) {
	l.Logger.Debug(fmt.Sprint(args...))
}
func (l *slogImpl) Debugf(format string, args ...interface{}) {
	l.Logger.Debug(fmt.Sprintf(format, args...))
}

// 其他 slog 的原始方法仍可直接使用
func (l *slogImpl) Fatal(args ...interface{}) {
	l.Logger.Error(fmt.Sprint(args...))
}
func (l *slogImpl) Fatalf(format string, args ...interface{}) {
	l.Logger.Error(fmt.Sprintf(format, args...))
}
