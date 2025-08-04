package logger

import (
	"context"
	"log/slog"
)

type logCtx struct {
	Message   string
	RequestID string
}

func (l logCtx) Attrs() []slog.Attr {
	var attrs []slog.Attr
	if l.Message != "" {
		attrs = append(attrs, slog.String("message", l.Message))
	}
	return attrs
}

type ctxKeyType int

const ctxKey ctxKeyType = 0

func mergeCtx(a, b logCtx) logCtx {
	if b.Message != "" {
		a.Message = b.Message
	}
	if b.RequestID != "" {
		a.RequestID = b.RequestID
	}
	return a
}

func WithLogMessage(ctx context.Context, message string) context.Context {
	c, _ := ctx.Value(ctxKey).(logCtx)
	c = mergeCtx(c, logCtx{Message: message})
	return context.WithValue(ctx, ctxKey, c)
}

func WithRequestID(ctx context.Context, id string) context.Context {
	c, _ := ctx.Value(ctxKey).(logCtx)
	c = mergeCtx(c, logCtx{RequestID: id})
	return context.WithValue(ctx, ctxKey, c)
}
