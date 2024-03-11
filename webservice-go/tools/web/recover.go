package web

import (
	"fmt"
	"log/slog"
	"net/http"

	"go.opentelemetry.io/otel"
	metricapi "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
)

// RecoveryMiddleware handles inner handlers that panic and log reason as an error
type RecoveryMiddleware struct {
	handler http.Handler
	logger  *slog.Logger

	errorCounter metricapi.Int64Counter
}

// NewRecoveryMiddleware creates a middleware that tries to recover from panics that happen when they reach the it and returns a 500 instead
func NewRecoveryMiddleware(logger *slog.Logger, meter metricapi.Meter) *RecoveryMiddleware {
	rm := &RecoveryMiddleware{handler: nil, logger: logger}
	if meter == nil {
		meter = noop.NewMeterProvider().Meter("noop-meter")
	}
	var err error
	rm.errorCounter, err = meter.Int64Counter("application.panics.recovered",
		metricapi.WithDescription("counts the number of panics recovered from the recovery middleware"),
	)
	if err != nil {
		otel.Handle(err)
	}
	return rm
}

func (rm *RecoveryMiddleware) Chain(handler http.Handler) http.Handler {
	rm.handler = handler
	return rm
}

func (rm *RecoveryMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if q := recover(); q != nil {
			w.WriteHeader(http.StatusInternalServerError)
			rm.errorCounter.Add(r.Context(), 1)
			rm.logger.ErrorContext(r.Context(), fmt.Sprintf("recovering from a panic: %+v", q))
		}
	}()

	rm.handler.ServeHTTP(w, r)
}
