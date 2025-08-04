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

func NewLogger(level slog.Level, jsonFormat bool) {
	once.Do(func() {
		var handler slog.Handler

		opts := &slog.HandlerOptions{Level: level}
		if jsonFormat {
			handler = slog.NewJSONHandler(os.Stdout, opts)
		} else {
			handler = slog.NewTextHandler(os.Stdout, opts)
		}

		handler = NewMiddleware(handler)
		global = slog.New(handler)
		slog.SetDefault(global)
	})
}
