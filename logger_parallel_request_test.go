package logger_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/scarymovie/logger"
)

type logEntry struct {
	GoroutineID int
	UniqueID    string
	Level       string
}

type safeWriter struct {
	writer io.Writer
	mu     *sync.Mutex
}

func (sw *safeWriter) Write(p []byte) (n int, err error) {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	return sw.writer.Write(p)
}

func TestConcurrentLoggingWithUniqueIDs(t *testing.T) {
	logger.ResetGlobal()

	var buf bytes.Buffer
	var mu sync.Mutex

	safeWriter := &safeWriter{writer: &buf, mu: &mu}

	jsonHandler := slog.NewJSONHandler(safeWriter, &slog.HandlerOptions{Level: slog.LevelDebug})
	handler := logger.NewMiddleware(jsonHandler)
	log := slog.New(handler)

	const numGoroutines = 5
	const logsPerGoroutine = 4

	var wg sync.WaitGroup
	logEntries := make(chan logEntry, numGoroutines*logsPerGoroutine)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			uniqueID := fmt.Sprintf("goroutine-%d-%d", goroutineID, time.Now().UnixNano())

			ctx := logger.WithLogMessage(context.Background(), uniqueID)

			log.DebugContext(ctx, "Debug message from goroutine",
				"goroutine_id", goroutineID, "operation", "debug_op")

			log.InfoContext(ctx, "Info message from goroutine",
				"goroutine_id", goroutineID, "operation", "info_op")

			log.WarnContext(ctx, "Warn message from goroutine",
				"goroutine_id", goroutineID, "operation", "warn_op")

			log.ErrorContext(ctx, "Error message from goroutine",
				"goroutine_id", goroutineID, "operation", "error_op")

			for _, level := range []string{"DEBUG", "INFO", "WARN", "ERROR"} {
				logEntries <- logEntry{
					GoroutineID: goroutineID,
					UniqueID:    uniqueID,
					Level:       level,
				}
			}
		}(i)
	}

	wg.Wait()
	close(logEntries)

	expectedEntries := make([]logEntry, 0, numGoroutines*logsPerGoroutine)
	for entry := range logEntries {
		expectedEntries = append(expectedEntries, entry)
	}

	mu.Lock()
	output := buf.String()
	mu.Unlock()

	lines := strings.Split(strings.TrimSpace(output), "\n")
	actualLogs := make([]map[string]interface{}, 0)

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		var logEntry map[string]interface{}
		if err := json.Unmarshal([]byte(line), &logEntry); err != nil {
			t.Fatalf("Failed to parse log entry: %s, error: %v", line, err)
		}
		actualLogs = append(actualLogs, logEntry)
	}

	expectedCount := numGoroutines * logsPerGoroutine
	if len(actualLogs) != expectedCount {
		t.Fatalf("Expected %d log entries, got %d", expectedCount, len(actualLogs))
	}

	uniqueIDToGoroutineID := make(map[string]int)
	goroutineIDCounts := make(map[int]int)
	levelCounts := make(map[string]int)

	for _, logEntry := range actualLogs {
		uniqueID, hasUniqueID := logEntry["message"].(string)
		if !hasUniqueID {
			t.Error("Log entry missing unique ID in message field")
			continue
		}

		goroutineIDFloat, hasGoroutineID := logEntry["goroutine_id"].(float64)
		if !hasGoroutineID {
			t.Error("Log entry missing goroutine_id field")
			continue
		}
		goroutineID := int(goroutineIDFloat)

		level, hasLevel := logEntry["level"].(string)
		if !hasLevel {
			t.Error("Log entry missing level field")
			continue
		}

		expectedPrefix := fmt.Sprintf("goroutine-%d-", goroutineID)
		if !strings.HasPrefix(uniqueID, expectedPrefix) {
			t.Errorf("Unique ID %s doesn't match expected format for goroutine %d",
				uniqueID, goroutineID)
		}

		if existingGoroutineID, exists := uniqueIDToGoroutineID[uniqueID]; exists {
			if existingGoroutineID != goroutineID {
				t.Errorf("Unique ID %s used by multiple goroutines: %d and %d",
					uniqueID, existingGoroutineID, goroutineID)
			}
		} else {
			uniqueIDToGoroutineID[uniqueID] = goroutineID
		}

		goroutineIDCounts[goroutineID]++
		levelCounts[level]++
	}

	for i := 0; i < numGoroutines; i++ {
		if count := goroutineIDCounts[i]; count != logsPerGoroutine {
			t.Errorf("Goroutine %d produced %d logs, expected %d", i, count, logsPerGoroutine)
		}
	}

	expectedLevelCount := numGoroutines
	for _, expectedLevel := range []string{"DEBUG", "INFO", "WARN", "ERROR"} {
		if count := levelCounts[expectedLevel]; count != expectedLevelCount {
			t.Errorf("Level %s has %d entries, expected %d", expectedLevel, count, expectedLevelCount)
		}
	}

	if len(uniqueIDToGoroutineID) != numGoroutines {
		t.Errorf("Expected %d unique IDs, got %d", numGoroutines, len(uniqueIDToGoroutineID))
	}

	t.Logf("Successfully verified %d log entries from %d goroutines with unique IDs",
		len(actualLogs), numGoroutines)
}
