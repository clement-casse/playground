package web

import (
	"fmt"
	"log/slog"
	"net/http"
)

// RecoveryMiddleware handles inner handlers that panic and log reason as an error
type RecoveryMiddleware struct {
	handler http.Handler
	logger  *slog.Logger
}

// NewRecoveryMiddleware creates a middleware that tries to recover from panics that happen when they reach the it and returns a 500 instead
func NewRecoveryMiddleware(logger *slog.Logger) *RecoveryMiddleware {
	return &RecoveryMiddleware{nil, logger}
}

func (rm *RecoveryMiddleware) Chain(handler http.Handler) http.Handler {
	rm.handler = handler
	return rm
}

func (rm *RecoveryMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if q := recover(); q != nil {
			w.WriteHeader(http.StatusInternalServerError)
			rm.logger.ErrorContext(r.Context(), fmt.Sprintf("recovering from a panic: %+v", q))
		}
	}()

	rm.handler.ServeHTTP(w, r)
}
