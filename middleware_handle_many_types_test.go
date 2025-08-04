package logger_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/scarymovie/logger"
	"log/slog"
	"strings"
	"testing"
)

func TestMiddlewareWithDifferentHandlerTypes(t *testing.T) {
	tests := []struct {
		name          string
		createHandler func(*bytes.Buffer) slog.Handler
		isJSON        bool
	}{
		{
			name: "JSONHandler",
			createHandler: func(buf *bytes.Buffer) slog.Handler {
				return slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})
			},
			isJSON: true,
		},
		{
			name: "TextHandler",
			createHandler: func(buf *bytes.Buffer) slog.Handler {
				return slog.NewTextHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})
			},
			isJSON: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			baseHandler := tt.createHandler(&buf)
			middleware := logger.NewMiddleware(baseHandler)
			log := slog.New(middleware)

			ctx := logger.WithLogMessage(context.Background(), "handler_test")
			log.InfoContext(ctx, "test message", "key", "value")

			output := buf.String()
			if output == "" {
				t.Fatal("Expected log output, got empty string")
			}

			if tt.isJSON {
				var logEntry map[string]interface{}
				if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &logEntry); err != nil {
					t.Fatalf("Failed to parse JSON log: %v", err)
				}

				if message, ok := logEntry["message"].(string); !ok || message != "handler_test" {
					t.Errorf("Expected context message 'handler_test', got %v", message)
				}
			} else {
				if !strings.Contains(output, "handler_test") {
					t.Error("Expected context message 'handler_test' to appear in text output")
				}
				if !strings.Contains(output, "test message") {
					t.Error("Expected log message 'test message' to appear in text output")
				}
			}
		})
	}
}
