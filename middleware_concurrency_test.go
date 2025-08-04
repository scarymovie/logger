package logger_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/scarymovie/logger"
	"log/slog"
	"strings"
	"sync"
	"testing"
)

func TestMiddlewareConcurrency(t *testing.T) {
	var buf bytes.Buffer
	var mu sync.Mutex

	safeWriter := &struct {
		*bytes.Buffer
		*sync.Mutex
	}{&buf, &mu}

	baseHandler := slog.NewJSONHandler(safeWriter, &slog.HandlerOptions{Level: slog.LevelDebug})
	middleware := logger.NewMiddleware(baseHandler)
	log := slog.New(middleware)

	const numGoroutines = 10
	const logsPerGoroutine = 5

	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < logsPerGoroutine; j++ {
				ctx := logger.WithLogMessage(context.Background(),
					fmt.Sprintf("goroutine-%d-log-%d", goroutineID, j))
				log.InfoContext(ctx, "concurrent test", "goroutine", goroutineID, "log", j)
			}
		}(i)
	}

	wg.Wait()

	mu.Lock()
	output := buf.String()
	mu.Unlock()

	lines := strings.Split(strings.TrimSpace(output), "\n")
	actualLogs := 0
	contextMessages := make(map[string]bool)

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		actualLogs++
		var logEntry map[string]interface{}
		if err := json.Unmarshal([]byte(line), &logEntry); err != nil {
			t.Errorf("Failed to parse JSON log: %v", err)
			continue
		}

		if message, ok := logEntry["message"].(string); ok {
			if contextMessages[message] {
				t.Errorf("Duplicate context message found: %s", message)
			}
			contextMessages[message] = true
		} else {
			t.Error("Expected context message in log entry")
		}
	}

	expectedLogs := numGoroutines * logsPerGoroutine
	if actualLogs != expectedLogs {
		t.Errorf("Expected %d log entries, got %d", expectedLogs, actualLogs)
	}

	if len(contextMessages) != expectedLogs {
		t.Errorf("Expected %d unique context messages, got %d", expectedLogs, len(contextMessages))
	}
}
