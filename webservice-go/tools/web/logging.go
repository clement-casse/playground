package web

import (
	"log/slog"
	"net/http"
	"time"
)

// AccessLoggingMiddleware handles requests and log with the logger requests that pass through it
type AccessLoggingMiddleware struct {
	logger *slog.Logger
}

// NewAccessLoggingMiddleware creates a middleware that logs requests that pass through it
func NewAccessLoggingMiddleware(logger *slog.Logger) *AccessLoggingMiddleware {
	return &AccessLoggingMiddleware{logger}
}

func (lm *AccessLoggingMiddleware) Chain(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		rww := newRespWriterWrapper(w)

		next.ServeHTTP(rww, r)

		lm.logger.InfoContext(r.Context(), "Incoming request",
			slog.Group("request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Duration("duration", time.Since(now)),
			),
			slog.Group("response",
				slog.Int("code", rww.status),
				slog.Int("size", int(rww.written.Load())),
			),
		)
	})
}
