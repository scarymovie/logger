package slogx

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"log/slog"
)

var (
	String   = slog.String
	Int      = slog.Int
	Int64    = slog.Int64
	Uint64   = slog.Uint64
	Bool     = slog.Bool
	Float64  = slog.Float64
	Time     = slog.Time
	Duration = slog.Duration
	Any      = slog.Any
	Group    = slog.Group
)

type Attr = slog.Attr

type redactSet map[string]struct{}

func toRedactSet(keys []string) redactSet {
	s := make(redactSet, len(keys))
	for _, k := range keys {
		if k == "" {
			continue
		}
		s[k] = struct{}{}
	}
	return s
}

var defaultLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

func L() *slog.Logger { return slog.Default() }

func Configure(c Config) (*slog.Logger, error) {
	cfg := c.WithDefaults()
	opts := slog.HandlerOptions{
		AddSource: cfg.AddSource,
		Level:     cfg.Level,
	}

	redactions := toRedactSet(cfg.RedactKeys)

	replace := func(groups []string, a slog.Attr) slog.Attr {
		if redactions != nil {
			if _, found := redactions[a.Key]; found {
				return slog.String(a.Key, "***")
			}
		}
		if a.Key == slog.TimeKey {
			t := a.Value.Time()
			if cfg.UseUTC {
				t = t.UTC()
			}
			return slog.String(slog.TimeKey, t.Format(cfg.TimeFormat))
		}
		if a.Key == slog.SourceKey {
			const skip = 9
			pc, file, line, ok := runtime.Caller(skip)
			if ok {
				src := slog.Source{
					Function: runtime.FuncForPC(pc).Name(),
					File:     file,
					Line:     line,
				}
				if cfg.ShortSource {
					src.File = filepath.Base(src.File)
				}
				return slog.Any(slog.SourceKey, src)
			}
		}
		if cfg.ReplaceAttr != nil {
			return cfg.ReplaceAttr(groups, a)
		}
		return a
	}
	opts.ReplaceAttr = replace

	var w io.Writer
	if len(cfg.Writers) > 0 {
		w = io.MultiWriter(cfg.Writers...)
	} else {
		w = cfg.Writer
	}

	var h slog.Handler
	switch strings.ToLower(cfg.Format) {
	case "json":
		h = slog.NewJSONHandler(w, &opts)
	case "text":
		h = slog.NewTextHandler(w, &opts)
	default:
		return nil, fmt.Errorf("unknown format %q (expected 'json' or 'text')", cfg.Format)
	}

	logger := slog.New(h)
	if len(cfg.DefaultAttrs) > 0 {
		logger = logger.With(attrsToAny(cfg.DefaultAttrs)...)
	}
	if cfg.UseAsDefault {
		slog.SetDefault(logger)
		defaultLogger = logger
	}
	return logger, nil
}

func MustConfigure(c Config) *slog.Logger {
	l, err := Configure(c)
	if err != nil {
		panic(err)
	}
	return l
}

func With(attrs ...Attr) *slog.Logger { return L().With(attrsToAny(attrs)...) }

// --- simplified logging ---

func Debug(ctx context.Context, msg string, args ...any) { FromContext(ctx).Debug(msg, args...) }
func Info(ctx context.Context, msg string, args ...any)  { FromContext(ctx).Info(msg, args...) }
func Warn(ctx context.Context, msg string, args ...any)  { FromContext(ctx).Warn(msg, args...) }
func Error(ctx context.Context, msg string, args ...any) { FromContext(ctx).Error(msg, args...) }

func Debugf(ctx context.Context, format string, args ...any) {
	FromContext(ctx).Debug(fmt.Sprintf(format, args...))
}
func Infof(ctx context.Context, format string, args ...any) {
	FromContext(ctx).Info(fmt.Sprintf(format, args...))
}
func Warnf(ctx context.Context, format string, args ...any) {
	FromContext(ctx).Warn(fmt.Sprintf(format, args...))
}
func Errorf(ctx context.Context, format string, args ...any) {
	FromContext(ctx).Error(fmt.Sprintf(format, args...))
}

// --- utils ---

func Background() context.Context { return context.Background() }

var Now = time.Now

func attrsToAny(attrs []slog.Attr) []any {
	switch len(attrs) {
	case 0:
		return nil
	case 1:
		return []any{attrs[0]}
	default:
		out := make([]any, len(attrs))
		for i := range attrs {
			out[i] = attrs[i]
		}
		return out
	}
}
