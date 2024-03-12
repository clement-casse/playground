package rest

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	metricapi "go.opentelemetry.io/otel/metric"
	traceapi "go.opentelemetry.io/otel/trace"

	"github.com/clement-casse/playground/webservice-go/tools/web"
)

// APIHandler is a server that registers endpoints for a REST API
type APIHandler struct {
	mux *http.ServeMux

	otelMeter  metricapi.Meter
	otelTracer traceapi.Tracer
	logger     *slog.Logger
}

// APIHandlerOpt in an interface for applying APIHandler options.
type APIHandlerOpt interface {
	applyOpt(*APIHandler) *APIHandler
}

type apiHandlerOptFunc func(*APIHandler) *APIHandler

func (fn apiHandlerOptFunc) applyOpt(s *APIHandler) *APIHandler {
	return fn(s)
}

// NewAPIHandler creates an API Handler for REST API
func NewAPIHandler(opts ...APIHandlerOpt) *APIHandler {
	apiHandler := &APIHandler{
		mux:        http.NewServeMux(),
		otelMeter:  nil,
		otelTracer: nil,
		logger:     slog.Default(),
	}

	for _, opt := range opts {
		apiHandler = opt.applyOpt(apiHandler)
	}
	return apiHandler
}

// WithLogger applies a custom logger for the APIHandler
func WithLogger(l *slog.Logger) APIHandlerOpt {
	return apiHandlerOptFunc(func(s *APIHandler) *APIHandler {
		s.logger = l
		return s
	})
}

// WithMeter applies a custom OpenTelemetry Meter for the APIHandler (if not set no metrics are collected)
func WithMeter(m metricapi.Meter) APIHandlerOpt {
	return apiHandlerOptFunc(func(s *APIHandler) *APIHandler {
		s.otelMeter = m
		return s
	})
}

// WithTracer applies a custom OpenTelemetry Tracer for the APIHandler (if not set no traces are collected)
func WithTracer(t traceapi.Tracer) APIHandlerOpt {
	return apiHandlerOptFunc(func(s *APIHandler) *APIHandler {
		s.otelTracer = t
		return s
	})
}

func setJSONHeader(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	})
}

type handlerFuncWithError func(http.ResponseWriter, *http.Request) error

func (s *APIHandler) registerRoute(pattern string, handlerFunc handlerFuncWithError) {
	handler := setJSONHeader(s.handleErrors(handlerFunc))
	if s.otelMeter != nil {
		handler = web.NewMetricsMiddleware(s.otelMeter, pattern).Chain(handler)
	}
	if s.otelTracer != nil {
		handler = otelhttp.NewHandler(handler, pattern)
	}
	s.mux.Handle(pattern, handler)
}

func (s *APIHandler) handleErrors(hwe handlerFuncWithError) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := hwe(w, r)
		if err != nil {
			var apiErr apiError
			if errors.As(err, &apiErr) {
				s.logger.ErrorContext(r.Context(), apiErr.Error())
				resp, err2 := json.Marshal(apiErr)
				if err2 != nil {
					s.logger.ErrorContext(r.Context(), "cannot marshal APIError: "+err2.Error(), "error", err.Error())
				}
				http.Error(w, string(resp), apiErr.status)
			} else {
				s.logger.ErrorContext(r.Context(), err.Error())
				http.Error(w, `{}`, http.StatusInternalServerError)
			}
		}
	}
}
