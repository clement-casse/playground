package web

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"gotest.tools/assert"
)

var (
	testListeningAddr = "127.0.0.1:8080"
	testHTTPHandler   = http.DefaultServeMux
)

func TestServerConfig(t *testing.T) {
	// After creating a new server instance none of the pointer should be nil to avoid runtime NPE
	newInstance1 := NewServer(testListeningAddr, testHTTPHandler)
	assert.Assert(t, newInstance1.logger != nil && newInstance1.server != nil)

	notTheDefaultLogger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	newInstance2 := NewServer(testListeningAddr, testHTTPHandler, WithLogger(notTheDefaultLogger))
	assert.Equal(t, notTheDefaultLogger, newInstance2.logger, "the logger is not applied properly")

	// Setting the logger to nil should instantiate a default logger instead
	newInstance3 := NewServer(testListeningAddr, testHTTPHandler, WithLogger(nil))
	assert.Assert(t, newInstance3.logger != nil, "it is possible to set the instance logger to nil")
	assert.Equal(t, slog.Default(), newInstance3.logger, "the nil logger does not fallback to slog default logger")
}

func TestWithReadTimeout(t *testing.T) {
	timeout, err := time.ParseDuration("42s")
	assert.NilError(t, err)
	rc := NewServer("", testHTTPHandler, WithReadTimeout(timeout))
	assert.Equal(t, timeout, rc.server.ReadTimeout)
}

func TestWithWriteTimeout(t *testing.T) {
	timeout, err := time.ParseDuration("42s")
	assert.NilError(t, err)
	rc := NewServer("", testHTTPHandler, WithWriteTimeout(timeout))
	assert.Equal(t, timeout, rc.server.WriteTimeout)
}

func TestWithIdleTimeout(t *testing.T) {
	timeout, err := time.ParseDuration("42s")
	assert.NilError(t, err)
	rc := NewServer("", testHTTPHandler, WithIdleTimeout(timeout))
	assert.Equal(t, timeout, rc.server.IdleTimeout)
}

func TestBaseHandler(t *testing.T) {
	s := &Server{}
	handler := s.makeHandler(http.DefaultServeMux)
	for _, tt := range []struct {
		name           string
		httpMethod     string
		httpURL        string
		expectedStatus int
		expectedBody   string // Response body should CONTAIN expectedBody
	}{
		{
			name:           "health endpoint",
			httpMethod:     "GET",
			httpURL:        "http://example.com/health",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"ok":true}`,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.httpMethod, tt.httpURL, nil)

			handler.ServeHTTP(w, req)

			resp := w.Result()
			body, err := io.ReadAll(resp.Body)
			assert.NilError(t, err)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			assert.Assert(t, strings.Contains(string(body), tt.expectedBody))
		})
	}
}

type testMiddleware struct {
	middlewareFunc func(http.Handler) http.Handler
}

func (tm *testMiddleware) Handle(next http.Handler) http.Handler {
	return tm.middlewareFunc(next)
}

var (
	// testMiddlewareFunc1 set the Header X-Testing to `Middleware#1` if not set before and prepend body with `Middleware 1: `
	testMiddlewareFunc1 = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("X-Testing", "Middleware#1")
			_, _ = w.Write([]byte("Middleware 1: "))
			h.ServeHTTP(w, r)
		})
	}

	// testMiddlewareFunc2 set the Header X-Testing to `Middleware#2` if not set before
	testMiddlewareFunc2 = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("X-Testing", "Middleware#2")
			h.ServeHTTP(w, r)
		})
	}
)

func TestWithMiddleware(t *testing.T) {
	testMux := testHTTPHandler
	testMux.HandleFunc("GET /testing", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Add("X-Testing", "No middlewares")
		_, err := w.Write([]byte("Test Body"))
		assert.NilError(t, err)
	})

	testMiddleware1 := &testMiddleware{testMiddlewareFunc1}
	testMiddleware2 := &testMiddleware{testMiddlewareFunc2}

	for _, tt := range []struct {
		name             string
		middlewares      []Middleware
		expectTestHeader string
		expectBody       string
	}{
		{
			name:             "no middlewares",
			middlewares:      []Middleware{},
			expectTestHeader: "No middlewares",
			expectBody:       "Test Body",
		}, {
			name:             "with one middleware",
			middlewares:      []Middleware{testMiddleware1},
			expectTestHeader: "Middleware#1",
			expectBody:       "Middleware 1: Test Body",
		}, {
			name:             "with two middlewares the last takes precedence over the previous ones",
			middlewares:      []Middleware{testMiddleware1, testMiddleware2},
			expectTestHeader: "Middleware#2",
			expectBody:       "Middleware 1: Test Body",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			testServer := NewServer(testListeningAddr, testMux, WithMiddlewares(tt.middlewares...))

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "http://example.com/testing", nil)

			testServer.mainHandler.ServeHTTP(w, req)

			resp := w.Result()
			body, err := io.ReadAll(resp.Body)
			assert.NilError(t, err)

			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.Equal(t, tt.expectBody, string(body))
			assert.Equal(t, tt.expectTestHeader, resp.Header.Get("X-Testing"))
		})
	}
}
