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
func (l *slogImpl) Debugln(args ...any) {
	l.Logger.Debug(fmt.Sprintln(args...))
}

// Errorln implements interfaces.Logger.
func (l *slogImpl) Errorln(args ...any) {
	l.Logger.Error(fmt.Sprintln(args...))
}

// Fatalln implements interfaces.Logger.
func (l *slogImpl) Fatalln(args ...any) {
	l.Fatal(fmt.Sprintln(args...))
}

// Infoln implements interfaces.Logger.
func (l *slogImpl) Infoln(args ...any) {
	l.Logger.Info(fmt.Sprintln(args...))
}

// Warnln implements interfaces.Logger.
func (l *slogImpl) Warnln(args ...any) {
	l.Logger.Warn(fmt.Sprintln(args...))
}

func New() interfaces.Logger {
	level := slog.LevelInfo
	// todo 目前没有处理With相关的，加了也没有用。
	// slog.Default().With(slog.String("key", "value"))
	logger := slog.New(NewPrettyHandler(os.Stdout, level))
	return &slogImpl{Logger: logger}
}

// 定义 Infof 方法
func (l *slogImpl) Info(args ...any) {
	l.Logger.Info(fmt.Sprint(args...))
}
func (l *slogImpl) Infof(format string, args ...any) {
	l.Logger.Info(fmt.Sprintf(format, args...))
}

// 定义 Warnf 方法
func (l *slogImpl) Warn(args ...any) {
	l.Logger.Warn(fmt.Sprint(args...))
}
func (l *slogImpl) Warnf(format string, args ...any) {
	l.Logger.Warn(fmt.Sprintf(format, args...))
}

// 定义 Errorf 方法
func (l *slogImpl) Error(args ...any) {
	l.Logger.Error(fmt.Sprint(args...))
}
func (l *slogImpl) Errorf(format string, args ...any) {
	l.Logger.Error(fmt.Sprintf(format, args...))
}

// 其他 slog 的原始方法仍可直接使用
func (l *slogImpl) Debug(args ...any) {
	l.Logger.Debug(fmt.Sprint(args...))
}
func (l *slogImpl) Debugf(format string, args ...any) {
	l.Logger.Debug(fmt.Sprintf(format, args...))
}

// 其他 slog 的原始方法仍可直接使用
func (l *slogImpl) Fatal(args ...any) {
	l.Logger.Error(fmt.Sprint(args...))
}
func (l *slogImpl) Fatalf(format string, args ...any) {
	l.Logger.Error(fmt.Sprintf(format, args...))
}
