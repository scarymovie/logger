package logger

import (
	"context"
	"errors"
)

type errorWithCtx struct {
	err error
	ctx logCtx
}

func (e *errorWithCtx) Error() string {
	return e.err.Error()
}

func (e *errorWithCtx) Unwrap() error {
	return e.err
}

func WrapError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}
	c, _ := ctx.Value(ctxKey).(logCtx)
	return &errorWithCtx{err: err, ctx: c}
}

func ErrorCtx(ctx context.Context, err error) context.Context {
	var e *errorWithCtx
	if errors.As(err, &e) {
		return context.WithValue(ctx, ctxKey, e.ctx)
	}
	return ctx
}
