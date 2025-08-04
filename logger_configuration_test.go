package logger_test

import (
	"github.com/scarymovie/logger"
	"log/slog"
	"testing"
)

func TestLoggerConfiguration(t *testing.T) {
	tests := []struct {
		name   string
		config logger.Config
	}{
		{
			name: "JSON format with Debug level",
			config: logger.Config{
				Level:      slog.LevelDebug,
				JSONFormat: true,
			},
		},
		{
			name: "Text format with Info level",
			config: logger.Config{
				Level:      slog.LevelInfo,
				JSONFormat: false,
			},
		},
		{
			name: "JSON format with Error level",
			config: logger.Config{
				Level:      slog.LevelError,
				JSONFormat: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.ResetGlobal()

			logger.NewLogger(tt.config)

			log := logger.GetLogger()
			if log == nil {
				t.Error("Expected logger to be initialized, got nil")
			}

			firstLogger := logger.GetLogger()
			logger.NewLogger(tt.config)
			secondLogger := logger.GetLogger()

			if firstLogger != secondLogger {
				t.Error("Expected same logger instance after multiple NewLogger calls")
			}
		})
	}
}
