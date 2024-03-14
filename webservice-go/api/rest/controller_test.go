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
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/clement-casse/playground/webservice-go/tools/web"
)

func TestWithLogger(t *testing.T) {
	l := &slog.Logger{}
	s := NewAPIController(WithLogger(l))
	assert.Same(t, l, s.logger)
}

func TestWithMeter(t *testing.T) {
	mp := sdkmetric.NewMeterProvider()
	meter := mp.Meter("some meter")
	s := NewAPIController(WithMeter(meter))
	assert.Same(t, meter, s.otelMeter)
}

func TestWithTracer(t *testing.T) {
	tp := sdktrace.NewTracerProvider()
	tracer := tp.Tracer("some tracer")
	s := NewAPIController(WithTracer(tracer))
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
			expectBody:   "{}\n",
			expectLog:    errTesting.Error(),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var recorder bytes.Buffer
			ctrl := &APIController{logger: slog.New(slog.NewTextHandler(&recorder, &slog.HandlerOptions{Level: slog.LevelError}))}

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "http:///", nil)

			ctrl.handleErrors(tt.innerHandler)(w, req)

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
		name          string
		apiController *APIController

		expectedErrorLogLines int
	}{
		{
			name:          "a normal APIHandler",
			apiController: NewAPIController(),
		}, {
			name:                  "an APIHandler with a logger",
			apiController:         NewAPIController(WithLogger(testLogger)),
			expectedErrorLogLines: 2,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := tt.apiController
			ctrl.registerRoute("GET /", workingHandler)
			ctrl.registerRoute("GET /notfound", apiErrorHandlerNotFound)
			ctrl.registerRoute("GET /wilderror", wildErrorHandler)

			testServer := httptest.NewServer(ctrl.mux)
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

func TestRegisterRouteWithMiddlewares(t *testing.T) {
	secretKey := []byte("privatekeyformytest")

	ctrl := NewAPIController()
	ctrl.registerRoute("GET /", workingHandler, web.NewJWTAuthMiddleware(secretKey))

	testServer := httptest.NewServer(ctrl.mux)
	defer testServer.Close()

	resUnauthorized, err := http.Get(testServer.URL + "/")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, resUnauthorized.StatusCode)

	// Testing a flow with JWT Signature
	validClaims := &jwt.RegisteredClaims{
		Issuer:    "test",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
	}
	validToken := jwt.NewWithClaims(jwt.SigningMethodHS256, validClaims)
	signedToken, err := validToken.SignedString(secretKey)
	assert.Nil(t, err)

	req, err := http.NewRequest(http.MethodGet, testServer.URL, nil)
	assert.Nil(t, err)
	req.Header.Set("Authorization", "Bearer "+signedToken)
	res, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}
