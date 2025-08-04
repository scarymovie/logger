package logger

import (
	"context"
	"log/slog"
	"strings"
)

type logCtx struct {
	UserID  int
	Phone   string
	Gate    string
	Message string
}

func (l logCtx) Attrs() []slog.Attr {
	var attrs []slog.Attr
	if l.UserID != 0 {
		attrs = append(attrs, slog.Int("user_id", l.UserID))
	}
	if l.Phone != "" {
		attrs = append(attrs, slog.String("phone", l.Phone))
	}
	if l.Gate != "" {
		attrs = append(attrs, slog.String("sms_gate", l.Gate))
	}
	if l.Message != "" {
		attrs = append(attrs, slog.String("message", l.Message))
	}
	return attrs
}

type ctxKeyType int

const ctxKey ctxKeyType = 0

func mergeCtx(a, b logCtx) logCtx {
	if b.UserID != 0 {
		a.UserID = b.UserID
	}
	if b.Phone != "" {
		a.Phone = b.Phone
	}
	if b.Gate != "" {
		a.Gate = b.Gate
	}
	if b.Message != "" {
		a.Message = b.Message
	}
	return a
}

func WithLogUserID(ctx context.Context, userID int) context.Context {
	c, _ := ctx.Value(ctxKey).(logCtx)
	c = mergeCtx(c, logCtx{UserID: userID})
	return context.WithValue(ctx, ctxKey, c)
}

func WithLogPhone(ctx context.Context, phone string) context.Context {
	if len(phone) > 4 {
		phone = strings.Repeat("*", len(phone)-4) + phone[len(phone)-4:]
	}
	c, _ := ctx.Value(ctxKey).(logCtx)
	c = mergeCtx(c, logCtx{Phone: phone})
	return context.WithValue(ctx, ctxKey, c)
}

func WithLogGate(ctx context.Context, gate string) context.Context {
	c, _ := ctx.Value(ctxKey).(logCtx)
	c = mergeCtx(c, logCtx{Gate: gate})
	return context.WithValue(ctx, ctxKey, c)
}

func WithLogMessage(ctx context.Context, message string) context.Context {
	c, _ := ctx.Value(ctxKey).(logCtx)
	c = mergeCtx(c, logCtx{Message: message})
	return context.WithValue(ctx, ctxKey, c)
}
