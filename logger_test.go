package logger_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/scarymovie/logger"
)

func TestLogging(t *testing.T) {
	logger.NewLogger(slog.LevelDebug, false)

	ctx := context.Background()
	ctx = logger.WithLogUserID(ctx, 42)
	ctx = logger.WithLogPhone(ctx, "+79991234567")
	ctx = logger.WithLogMessage(ctx, "hello")

	slog.InfoContext(ctx, "Test log")
}
