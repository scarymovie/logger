package slogx

import (
	"context"

	"log/slog"
)

type ctxKey struct{}

type ctxAttrs struct{ list []slog.Attr }

func WithContext(ctx context.Context, attrs ...slog.Attr) context.Context {
	if len(attrs) == 0 {
		return ctx
	}
	existing := contextAttrs(ctx)
	combined := make([]slog.Attr, 0, len(existing)+len(attrs))
	combined = append(combined, existing...)
	combined = append(combined, attrs...)
	return context.WithValue(ctx, ctxKey{}, ctxAttrs{list: combined})
}

func contextAttrs(ctx context.Context) []slog.Attr {
	if v := ctx.Value(ctxKey{}); v != nil {
		if ca, ok := v.(ctxAttrs); ok {
			return ca.list
		}
	}
	return nil
}

func FromContext(ctx context.Context) *slog.Logger {
	attrs := contextAttrs(ctx)
	if len(attrs) == 0 {
		return L()
	}
	return L().With(attrsToAny(attrs)...)
}
