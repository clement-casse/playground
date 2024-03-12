package rest

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func TestWithLogger(t *testing.T) {
	l := &slog.Logger{}
	s := NewAPIHandler(WithLogger(l))
	assert.Same(t, l, s.logger)
}

func TestWithMeter(t *testing.T) {
	mp := sdkmetric.NewMeterProvider()
	meter := mp.Meter("some meter")
	s := NewAPIHandler(WithMeter(meter))
	assert.Same(t, meter, s.otelMeter)
}

func TestWithTracer(t *testing.T) {
	tp := sdktrace.NewTracerProvider()
	tracer := tp.Tracer("some tracer")
	s := NewAPIHandler(WithTracer(tracer))
	assert.Same(t, tracer, s.otelTracer)
}

var (
	errTesting     = fmt.Errorf("Testing Error with a private Text")
	workingHandler = func(w http.ResponseWriter, _ *http.Request) error {
		fmt.Fprint(w, `{"test": "ok"}`)
		return nil
	}

	apiErrorHandlerNotFound = func(_ http.ResponseWriter, _ *http.Request) error {
		return newAPIError(errTesting, http.StatusNotFound, "A public message")
	}

	wildErrorHandler = func(_ http.ResponseWriter, _ *http.Request) error {
		return errTesting
	}
)

func TestHandleErrors(t *testing.T) {
	for _, tt := range []struct {
		name         string
		innerHandler handlerFuncWithError
		expectStatus int
		expectBody   string
		expectLog    string
	}{
		{
			name:         "a normal handler should not be affected",
			innerHandler: workingHandler,
			expectStatus: http.StatusOK,
			expectBody:   `{"test": "ok"}`,
			expectLog:    "",
		}, {
			name:         "status code should be extracted from apiError structure",
			innerHandler: apiErrorHandlerNotFound,
			expectStatus: http.StatusNotFound,
			expectBody:   `{"status":404,"message":"A public message"}` + "\n",
			expectLog:    errTesting.Error(),
		}, {
			name:         "an unwrapped error should be reported as 500 internal server error with no body",
			innerHandler: wildErrorHandler,
			expectStatus: http.StatusInternalServerError,
			expectBody:   "\n",
			expectLog:    errTesting.Error(),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var recorder bytes.Buffer
			apiHandler := &APIHandler{logger: slog.New(slog.NewTextHandler(&recorder, &slog.HandlerOptions{Level: slog.LevelError}))}

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "http:///", nil)

			apiHandler.handleErrors(tt.innerHandler)(w, req)

			resp := w.Result()
			assert.Equal(t, tt.expectStatus, resp.StatusCode, "unexpected status")
			body, err := io.ReadAll(resp.Body)
			assert.Nil(t, err)
			assert.Equal(t, tt.expectBody, string(body))

			logLines := recorder.String()
			if tt.expectLog == "" {
				assert.Equal(t, tt.expectLog, logLines, "unexpected log line")
			} else {
				assert.True(t, strings.Contains(logLines, tt.expectLog), "cannot find error field in log")
			}
		})
	}
}

func TestRegisterRoute(t *testing.T) {
	var recorder bytes.Buffer
	testLogger := slog.New(slog.NewTextHandler(&recorder, &slog.HandlerOptions{Level: slog.LevelError}))

	for _, tt := range []struct {
		name       string
		apiHandler *APIHandler

		expectedErrorLogLines int
	}{
		{
			name:       "a normal APIHandler",
			apiHandler: NewAPIHandler(),
		}, {
			name:                  "an APIHandler with a logger",
			apiHandler:            NewAPIHandler(WithLogger(testLogger)),
			expectedErrorLogLines: 2,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			simpleAPIHandler := tt.apiHandler

			simpleAPIHandler.registerRoute("GET /", workingHandler)
			simpleAPIHandler.registerRoute("GET /notfound", apiErrorHandlerNotFound)
			simpleAPIHandler.registerRoute("GET /wilderror", wildErrorHandler)

			testServer := httptest.NewServer(simpleAPIHandler.mux)
			defer testServer.Close()

			resOK, err := http.Get(testServer.URL + "/")
			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resOK.StatusCode)
			bodyOK, err := io.ReadAll(resOK.Body)
			assert.Nil(t, err)
			assert.Equal(t, `{"test": "ok"}`, string(bodyOK))
			assert.True(t, resOK.Header.Get("Content-Type") == "application/json")

			resNotFound, err := http.Get(testServer.URL + "/notfound")
			assert.Nil(t, err)
			assert.Equal(t, http.StatusNotFound, resNotFound.StatusCode)

			resInternalError, err := http.Get(testServer.URL + "/wilderror")
			assert.Nil(t, err)
			assert.Equal(t, http.StatusInternalServerError, resInternalError.StatusCode)

			if tt.expectedErrorLogLines != 0 {
				logScanner := bufio.NewScanner(&recorder)
				lineCounter := 0
				for logScanner.Scan() {
					lineCounter++
				}
				assert.NoError(t, logScanner.Err())
				assert.Equal(t, tt.expectedErrorLogLines, lineCounter)
			}

			recorder.Reset()
		})
	}
}
