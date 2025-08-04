package logger_test

import (
	"bytes"
	"context"
	"github.com/scarymovie/logger"
	"log/slog"
	"testing"
)

func TestMiddlewareEnabled(t *testing.T) {
	tests := []struct {
		name          string
		handlerLevel  slog.Level
		testLevel     slog.Level
		shouldEnabled bool
	}{
		{
			name:          "Debug handler should enable debug level",
			handlerLevel:  slog.LevelDebug,
			testLevel:     slog.LevelDebug,
			shouldEnabled: true,
		},
		{
			name:          "Info handler should not enable debug level",
			handlerLevel:  slog.LevelInfo,
			testLevel:     slog.LevelDebug,
			shouldEnabled: false,
		},
		{
			name:          "Info handler should enable info level",
			handlerLevel:  slog.LevelInfo,
			testLevel:     slog.LevelInfo,
			shouldEnabled: true,
		},
		{
			name:          "Error handler should not enable warn level",
			handlerLevel:  slog.LevelError,
			testLevel:     slog.LevelWarn,
			shouldEnabled: false,
		},
		{
			name:          "Error handler should enable error level",
			handlerLevel:  slog.LevelError,
			testLevel:     slog.LevelError,
			shouldEnabled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			baseHandler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: tt.handlerLevel})
			middleware := logger.NewMiddleware(baseHandler)

			ctx := context.Background()
			enabled := middleware.Enabled(ctx, tt.testLevel)

			if enabled != tt.shouldEnabled {
				t.Errorf("Expected Enabled to return %v, got %v", tt.shouldEnabled, enabled)
			}
		})
	}
}
