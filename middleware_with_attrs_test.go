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

func TestMiddlewareWithAttrs(t *testing.T) {
	var buf bytes.Buffer
	baseHandler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	middleware := logger.NewMiddleware(baseHandler)

	middlewareWithAttrs := middleware.WithAttrs([]slog.Attr{
		slog.String("service", "test-service"),
		slog.Int("version", 1),
	})

	log := slog.New(middlewareWithAttrs)

	ctx := logger.WithLogMessage(context.Background(), "attrs_test")
	log.InfoContext(ctx, "test message", "extra", "value")

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

	if message, ok := logEntry["message"].(string); !ok || message != "attrs_test" {
		t.Errorf("Expected context message 'attrs_test', got %v", message)
	}

	if service, ok := logEntry["service"].(string); !ok || service != "test-service" {
		t.Errorf("Expected service 'test-service', got %v", service)
	}

	if version, ok := logEntry["version"].(float64); !ok || version != 1 {
		t.Errorf("Expected version 1, got %v", version)
	}

	if extra, ok := logEntry["extra"].(string); !ok || extra != "value" {
		t.Errorf("Expected extra 'value', got %v", extra)
	}
}
