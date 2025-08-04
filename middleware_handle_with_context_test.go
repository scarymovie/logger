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

func TestMiddlewareHandleWithContext(t *testing.T) {
	var buf bytes.Buffer
	baseHandler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	middleware := logger.NewMiddleware(baseHandler)
	log := slog.New(middleware)

	ctx := logger.WithLogMessage(context.Background(), "context_message")
	log.InfoContext(ctx, "test message", "key1", "value1")

	output := buf.String()
	if output == "" {
		t.Fatal("Expected log output, got empty string")
	}

	var logEntry map[string]interface{}
	if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &logEntry); err != nil {
		t.Fatalf("Failed to parse JSON log: %v", err)
	}

	if msg, ok := logEntry["msg"].(string); !ok || msg != "test message" {
		t.Errorf("Expected msg 'test message', got %v", msg)
	}

	if key1, ok := logEntry["key1"].(string); !ok || key1 != "value1" {
		t.Errorf("Expected key1 'value1', got %v", key1)
	}

	if message, ok := logEntry["message"].(string); !ok || message != "context_message" {
		t.Errorf("Expected context message 'context_message', got %v", message)
	}
}
