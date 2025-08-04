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

func TestMiddlewareHandleWithEmptyContext(t *testing.T) {
	var buf bytes.Buffer
	baseHandler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	middleware := logger.NewMiddleware(baseHandler)
	log := slog.New(middleware)

	ctx := logger.WithLogMessage(context.Background(), "")
	log.InfoContext(ctx, "test message")

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

	if _, hasMessage := logEntry["message"]; hasMessage {
		t.Error("Expected no 'message' attribute for empty context message, but found one")
	}
}
