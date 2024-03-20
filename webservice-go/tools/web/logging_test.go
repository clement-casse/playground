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

func TestLoggingMiddleware(t *testing.T) {
	var recorder bytes.Buffer
	recordedLogger := slog.New(slog.NewTextHandler(&recorder, nil))
	lm := NewAccessLoggingMiddleware(recordedLogger)
	testServer := httptest.NewServer(lm.Handle(testingHandler))
	req, err := http.NewRequest(http.MethodGet, testServer.URL, nil)
	assert.NilError(t, err)

	_, err = http.DefaultClient.Do(req)
	assert.NilError(t, err)
	testServer.Close()

	logLines := recorder.String()
	assert.Assert(t, strings.Contains(logLines, fmt.Sprintf("method=%s", req.Method)), "cannot find method field in log")
	assert.Assert(t, strings.Contains(logLines, fmt.Sprintf("path=%s", req.URL.Path)), "cannot find path field in log")

}
