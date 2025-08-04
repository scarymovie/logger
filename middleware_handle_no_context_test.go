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

func TestMiddlewareHandleWithoutContext(t *testing.T) {
	var buf bytes.Buffer
	baseHandler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	middleware := logger.NewMiddleware(baseHandler)
	log := slog.New(middleware)

	ctx := context.Background()
	log.InfoContext(ctx, "test message", "key1", "value1", "key2", 42)

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

	if key2, ok := logEntry["key2"].(float64); !ok || key2 != 42 {
		t.Errorf("Expected key2 42, got %v", key2)
	}

	if _, hasMessage := logEntry["message"]; hasMessage {
		t.Error("Expected no 'message' attribute from context, but found one")
	}
}
