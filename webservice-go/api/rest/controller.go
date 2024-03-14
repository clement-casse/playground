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

// APIController is a server that registers endpoints for a REST API
type APIController struct {
	mux *http.ServeMux

	otelMeter  metricapi.Meter
	otelTracer traceapi.Tracer
	logger     *slog.Logger
}

// APIControllerOpt in an interface for applying APIHandler options.
type APIControllerOpt interface {
	applyOpt(*APIController) *APIController
}

type apiControllerOptFunc func(*APIController) *APIController

func (fn apiControllerOptFunc) applyOpt(s *APIController) *APIController {
	return fn(s)
}

// NewAPIController creates an API Handler for REST API
func NewAPIController(opts ...APIControllerOpt) *APIController {
	apiController := &APIController{
		mux:        http.NewServeMux(),
		otelMeter:  nil,
		otelTracer: nil,
		logger:     slog.Default(),
	}

	for _, opt := range opts {
		apiController = opt.applyOpt(apiController)
	}
	return apiController
}

// WithLogger applies a custom logger for the APIHandler
func WithLogger(l *slog.Logger) APIControllerOpt {
	return apiControllerOptFunc(func(s *APIController) *APIController {
		s.logger = l
		return s
	})
}

// WithMeter applies a custom OpenTelemetry Meter for the APIHandler (if not set no metrics are collected)
func WithMeter(m metricapi.Meter) APIControllerOpt {
	return apiControllerOptFunc(func(s *APIController) *APIController {
		s.otelMeter = m
		return s
	})
}

// WithTracer applies a custom OpenTelemetry Tracer for the APIHandler (if not set no traces are collected)
func WithTracer(t traceapi.Tracer) APIControllerOpt {
	return apiControllerOptFunc(func(s *APIController) *APIController {
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

func (c *APIController) registerRoute(pattern string, handlerFunc handlerFuncWithError, middlewares ...web.MiddlewareChainer) {
	handler := setJSONHeader(c.handleErrors(handlerFunc))
	if c.otelMeter != nil {
		handler = web.NewMetricsMiddleware(c.otelMeter, pattern).Chain(handler)
	}
	if c.otelTracer != nil {
		handler = otelhttp.NewHandler(handler, pattern)
	}
	for _, mw := range middlewares {
		handler = mw.Chain(handler)
	}
	c.mux.Handle(pattern, handler)
}

func (c *APIController) handleErrors(hwe handlerFuncWithError) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := hwe(w, r)
		if err != nil {
			var apiErr apiError
			if errors.As(err, &apiErr) {
				c.logger.ErrorContext(r.Context(), apiErr.Error())
				resp, err2 := json.Marshal(apiErr)
				if err2 != nil {
					c.logger.ErrorContext(r.Context(), "cannot marshal APIError: "+err2.Error(), "error", err.Error())
				}
				http.Error(w, string(resp), apiErr.status)
			} else {
				c.logger.ErrorContext(r.Context(), err.Error())
				http.Error(w, `{}`, http.StatusInternalServerError)
			}
		}
	}
}
