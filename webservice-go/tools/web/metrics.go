package web

import (
	"cmp"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	api "go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

type metricsMiddleware struct {
	pattern string

	requestDurationHist  api.Int64Histogram
	requestBytesCounter  api.Int64Counter
	responseBytesCounter api.Int64Counter
}

// NewMetricsMiddleware creates a new metric monitoring middleware
func NewMetricsMiddleware(otelMeter api.Meter, pattern string) Middleware {
	m := &metricsMiddleware{pattern: pattern}
	var err error
	m.requestDurationHist, err = otelMeter.Int64Histogram(
		"http.server.request.duration",
		metric.WithUnit("ms"),
		api.WithDescription("Measures the duration of inbound HTTP requests."),
	)
	if err != nil {
		otel.Handle(err)
	}
	m.requestBytesCounter, err = otelMeter.Int64Counter(
		string(semconv.HTTPRequestBodySizeKey),
		metric.WithUnit("By"),
		api.WithDescription("Measures the size of HTTP request messages."),
	)
	if err != nil {
		otel.Handle(err)
	}
	m.responseBytesCounter, err = otelMeter.Int64Counter(
		string(semconv.HTTPResponseBodySizeKey),
		metric.WithUnit("By"),
		api.WithDescription("Measures the size of HTTP response messages."),
	)
	if err != nil {
		otel.Handle(err)
	}

	return m
}

func (m *metricsMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		now := time.Now()
		rww := newRespWriterWrapper(w)

		var bw bodyWrapper
		// if request body is nil or NoBody, we don't want to mutate the body as it
		// will affect the identity of it in an unforeseeable way because we assert
		// ReadCloser fulfills a certain interface and it is indeed nil or NoBody.
		if r.Body != nil && r.Body != http.NoBody {
			bw.ReadCloser = r.Body
			r.Body = &bw
		}

		next.ServeHTTP(rww, r)

		httpRouteKey := cmp.Or(m.pattern, r.URL.Path) // does a coalesce operation since go 1.22: i.e. if m.pattern == "" then r.URL.Path is used
		o := metric.WithAttributes(
			semconv.HTTPRequestMethodKey.String(r.Method),
			semconv.HTTPResponseStatusCode(rww.status),
			semconv.HTTPRouteKey.String(httpRouteKey),
		)

		m.requestDurationHist.Record(ctx, time.Since(now).Milliseconds(), o)
		m.requestBytesCounter.Add(ctx, bw.read.Load(), o)
		m.responseBytesCounter.Add(ctx, rww.written.Load(), o)
	})
}
