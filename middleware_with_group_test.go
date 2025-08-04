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

func TestMiddlewareWithGroup(t *testing.T) {
	var buf bytes.Buffer
	baseHandler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	middleware := logger.NewMiddleware(baseHandler)

	middlewareWithGroup := middleware.WithGroup("app")
	log := slog.New(middlewareWithGroup)

	ctx := logger.WithLogMessage(context.Background(), "group_test")
	log.InfoContext(ctx, "test message", "nested", "value")

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

	if app, ok := logEntry["app"].(map[string]interface{}); ok {
		if message, ok := app["message"].(string); !ok || message != "group_test" {
			t.Errorf("Expected app.message 'group_test', got %v", message)
		}

		if nested, ok := app["nested"].(string); !ok || nested != "value" {
			t.Errorf("Expected app.nested 'value', got %v", nested)
		}
	} else {
		t.Errorf("Expected 'app' group, got %v", logEntry["app"])
	}

	if _, hasTopLevelMessage := logEntry["message"]; hasTopLevelMessage {
		t.Error("Expected no top-level 'message' attribute when using WithGroup")
	}
}
