package logger

import (
	"log/slog"
	"os"
	"sync"
)

var (
	once   sync.Once
	global *slog.Logger
)

type GroupedAttrs struct {
	Group string
	Attrs []slog.Attr
}

type Config struct {
	Level        slog.Level
	JSONFormat   bool
	DefaultAttrs []slog.Attr
	GroupedAttrs []GroupedAttrs
}

func NewLogger(cfg Config) {
	once.Do(func() {
		var handler slog.Handler

		opts := &slog.HandlerOptions{Level: cfg.Level}
		if cfg.JSONFormat {
			handler = slog.NewJSONHandler(os.Stdout, opts)
		} else {
			handler = slog.NewTextHandler(os.Stdout, opts)
		}

		for _, group := range cfg.GroupedAttrs {
			if group.Group != "" && len(group.Attrs) > 0 {
				handler = handler.WithGroup(group.Group).WithAttrs(group.Attrs)
			}
		}

		if len(cfg.DefaultAttrs) > 0 {
			handler = handler.WithAttrs(cfg.DefaultAttrs)
		}

		handler = NewMiddleware(handler)
		global = slog.New(handler)
		slog.SetDefault(global)
	})
}

func GetLogger() *slog.Logger {
	return global
}
