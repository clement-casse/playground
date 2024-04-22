package web

import (
	"log/slog"
	"net/http"
	"time"
)

type accessLoggingMiddleware struct {
	logger *slog.Logger
}

// verify Middleware interface compliance
var _ Middleware = (*accessLoggingMiddleware)(nil)

// NewAccessLoggingMiddleware creates a middleware that logs requests that pass through it.
func NewAccessLoggingMiddleware(logger *slog.Logger) Middleware {
	return &accessLoggingMiddleware{logger}
}

func (m *accessLoggingMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		rww := newRespWriterWrapper(w)

		next.ServeHTTP(rww, r)

		m.logger.InfoContext(r.Context(), "Incoming request",
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
