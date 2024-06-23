package web

import (
	"fmt"
	"log/slog"
	"net/http"

	"go.opentelemetry.io/otel"
	metricapi "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
)

type recoveryMiddleware struct {
	logger *slog.Logger

	errorCounter metricapi.Int64Counter
}

// verify Middleware interface compliance
var _ Middleware = (*recoveryMiddleware)(nil)

// NewRecoveryMiddleware creates a middleware that tries to recover from panics that
// happen when they reach the it and returns a 500 instead
func NewRecoveryMiddleware(logger *slog.Logger, meter metricapi.Meter) Middleware {
	m := &recoveryMiddleware{logger: logger}
	if meter == nil {
		meter = noop.NewMeterProvider().Meter("noop-meter")
	}
	var err error
	m.errorCounter, err = meter.Int64Counter("application.panics.recovered",
		metricapi.WithDescription("counts the number of panics recovered from the recovery middleware"),
	)
	if err != nil {
		otel.Handle(err)
	}
	return m
}

func (m *recoveryMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if q := recover(); q != nil {
				w.WriteHeader(http.StatusInternalServerError)
				m.errorCounter.Add(r.Context(), 1)
				m.logger.ErrorContext(r.Context(), fmt.Sprintf("recovering from a panic: %+v", q))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
