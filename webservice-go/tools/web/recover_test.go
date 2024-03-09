package web

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gotest.tools/assert"
)

var (
	testingHandler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "test ok")
	})

	panicReason     = "something weird"
	panicingHandler = http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		panic(panicReason)
	})
)

func TestRecoveryMiddleware(t *testing.T) {
	for _, tt := range []struct {
		name         string
		handlerFunc  http.HandlerFunc
		expectStatus int
	}{
		{
			name:         "a normal handler should not be affected",
			handlerFunc:  testingHandler,
			expectStatus: http.StatusOK,
		}, {
			name:         "a panic in the handler should be caught",
			handlerFunc:  panicingHandler,
			expectStatus: http.StatusInternalServerError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			rm := NewRecoveryMiddleware(slog.Default())
			testServer := httptest.NewServer(rm.Chain(tt.handlerFunc))
			defer testServer.Close()

			res, err := http.Get(testServer.URL)
			assert.NilError(t, err)
			assert.Equal(t, tt.expectStatus, res.StatusCode, "unexpected status")
		})
	}
}

func TestRecoveryMiddlewareLogsPanicReason(t *testing.T) {
	var recorder bytes.Buffer
	strLogger := slog.New(slog.NewTextHandler(&recorder, nil))

	rm := NewRecoveryMiddleware(strLogger)
	testServer := httptest.NewServer(rm.Chain(panicingHandler))

	_, err := http.Get(testServer.URL)
	assert.NilError(t, err)
	testServer.Close()

	logLines := recorder.String()
	assert.Assert(t, strings.Contains(logLines, panicReason), "recovery middleware does not print the inner panic reason")
}
