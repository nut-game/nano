package slog

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"

	"github.com/fatih/color"
)

type PrettyHandler struct {
	slog.Handler
	l *log.Logger
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {

	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString("[D]")
	case slog.LevelInfo:
		level = color.BlueString("[I]")
	case slog.LevelWarn:
		level = color.YellowString("[W]")
	case slog.LevelError:
		level = color.RedString("[E]")
	}
	// 组合键值对数据
	attrs := ""
	r.Attrs(func(a slog.Attr) bool {
		attrs += fmt.Sprintf("%s=%v ", a.Key, a.Value.Any())
		return true
	})

	// b, err := json.MarshalIndent(fields, "", "  ")
	// if err != nil {
	// 	return err
	// }

	h.l.Println(r.Time.Format("2006/01/02 15:04:05 -0700"),
		level,
		color.WhiteString(r.Message),
		color.CyanString(attrs))
	return nil

}

func NewPrettyHandler(out io.Writer, level slog.Level) *PrettyHandler {
	h := &PrettyHandler{
		Handler: slog.NewTextHandler(out, &slog.HandlerOptions{Level: level}),
		l:       log.New(out, "", 0),
	}
	return h
}
