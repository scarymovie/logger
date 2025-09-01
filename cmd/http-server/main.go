package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"log/slog"

	"github.com/google/uuid"
	"github.com/scarymovie/logger/slogx"
)

func main() {
	slogx.MustConfigure(slogx.Config{
		Format:       "json",
		Level:        slog.LevelDebug,
		AddSource:    true,
		RedactKeys:   []string{"password"},
		DefaultAttrs: []slog.Attr{slogx.String("service", "http-demo")},
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: requestIDMiddleware(mux),
	}

	slogx.Info(context.Background(), "HTTP server starting", slogx.String("addr", srv.Addr))
	if err := srv.ListenAndServe(); err != nil {
		slogx.Error(context.Background(), "server failed", slogx.String("error", err.Error()))
	}
}

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.New().String()
		ctx := slogx.WithContext(r.Context(), slogx.String("request_id", reqID))
		start := time.Now()

		slogx.Info(ctx, "incoming request",
			slogx.String("method", r.Method),
			slogx.String("path", r.URL.Path),
		)

		next.ServeHTTP(w, r.WithContext(ctx))

		slogx.Info(ctx, "request completed",
			slogx.Duration("duration", time.Since(start)),
		)
	})
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slogx.Infof(ctx, "handling request to %s", r.URL.Path)
	fmt.Fprintf(w, "Hello, World!")
}
