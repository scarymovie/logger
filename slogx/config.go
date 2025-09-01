package slogx

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"log/slog"
)

type Config struct {
	// Format selects the handler format: "json" or "text". Default: "json".
	Format string
	// Level sets the minimum level for the handler. Default: slog.LevelInfo.
	Level slog.Level
	// AddSource toggles source location reporting. Default: false.
	AddSource bool
	// ShortSource trims file paths to base names if true. Default: true when AddSource is true.
	ShortSource bool
	// TimeFormat allows customizing the timestamp format. If empty, RFC3339Nano is used.
	TimeFormat string
	// UseUTC forces timestamps to be written in UTC. Default: false.
	UseUTC bool
	// Writer is the destination for logs. Default: os.Stdout.
	Writer io.Writer
	// Writers allows multiple outputs (e.g., stdout + file). Overrides Writer if set.
	Writers []io.Writer
	// RedactKeys is a set of attribute keys whose values should be redacted.
	RedactKeys []string
	// DefaultAttrs are attributes added to every log entry.
	DefaultAttrs []slog.Attr
	// ReplaceAttr allows callers to customize attributes after builtin transforms.
	ReplaceAttr func(groups []string, a slog.Attr) slog.Attr
	// UseAsDefault sets both slog's global default and slogx global if true. Default: true.
	UseAsDefault bool
}

func (c Config) Clone() Config { return c }

func (c Config) WithDefaults() Config {
	cfg := c.Clone()
	if cfg.Format == "" {
		cfg.Format = "json"
	}
	if cfg.Writer == nil && len(cfg.Writers) == 0 {
		cfg.Writer = os.Stdout
	}
	if cfg.Level == 0 {
		cfg.Level = slog.LevelInfo
	}
	if cfg.TimeFormat == "" {
		cfg.TimeFormat = time.RFC3339Nano
	}
	if c.AddSource && !c.ShortSource {
		cfg.ShortSource = true
	}
	if !c.UseAsDefault {
		cfg.UseAsDefault = true
	}
	return cfg
}

func ParseLevel(level string) (slog.Level, error) {
	s := strings.TrimSpace(strings.ToLower(level))
	switch s {
	case "debug":
		return slog.LevelDebug, nil
	case "info", "":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error", "err":
		return slog.LevelError, nil
	case "fatal":
		return slog.LevelError, nil
	default:
		if n, err := parseInt(level); err == nil {
			return slog.Level(n), nil
		}
		return 0, fmt.Errorf("unknown level: %q", level)
	}
}

func parseInt(s string) (int, error) {
	var sign = 1
	if len(s) > 0 && (s[0] == '+' || s[0] == '-') {
		if s[0] == '-' {
			sign = -1
		}
		s = s[1:]
	}
	var n int
	if s == "" {
		return 0, fmt.Errorf("empty")
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("invalid digit")
		}
		n = n*10 + int(c-'0')
	}
	return sign * n, nil
}
