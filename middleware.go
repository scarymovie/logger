package logger

import (
	"context"
	"log/slog"
)

type middlewareHandler struct {
	next slog.Handler
}

func NewMiddleware(next slog.Handler) slog.Handler {
	return &middlewareHandler{next: next}
}

func (h *middlewareHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h *middlewareHandler) Handle(ctx context.Context, rec slog.Record) error {
	if c, ok := ctx.Value(ctxKey).(logCtx); ok {
		rec.AddAttrs(c.Attrs()...)
	}
	return h.next.Handle(ctx, rec)
}

func (h *middlewareHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &middlewareHandler{next: h.next.WithAttrs(attrs)}
}

func (h *middlewareHandler) WithGroup(name string) slog.Handler {
	return &middlewareHandler{next: h.next.WithGroup(name)}
}
