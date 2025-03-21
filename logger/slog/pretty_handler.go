package slog

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/fatih/color"
)

type PrettyHandler struct {
	slog.Handler
	o io.Writer
}

func (h *PrettyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
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

	io.WriteString(h.o, fmt.Sprintf("%s %s %s %s\n", r.Time.Format("2006/01/02 15:04:05 -0700"),
		level,
		color.WhiteString(r.Message),
		color.CyanString(attrs)))

	// _, err := os.Stdout.WriteString(fmt.Sprintf("%s %s %s %s\n",
	// 	r.Time.Format("2006/01/02 15:04:05 -0700"),
	// 	level,
	// 	color.WhiteString(r.Message),
	// 	color.CyanString(attrs),
	// ))

	return nil

}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// for k, v := range attrs {
	// 	fmt.Println(k, v)
	// }
	return h

}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	// fmt.Println(name)
	return h
}

func NewPrettyHandler(out io.Writer, level slog.Level) *PrettyHandler {
	return &PrettyHandler{
		Handler: slog.NewTextHandler(out,
			&slog.HandlerOptions{
				Level: level,
			}),
		o: out,
	}
}
