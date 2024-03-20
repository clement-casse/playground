package web

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggingMiddleware(t *testing.T) {
	var recorder bytes.Buffer
	recordedLogger := slog.New(slog.NewTextHandler(&recorder, nil))
	lm := NewAccessLoggingMiddleware(recordedLogger)
	testServer := httptest.NewServer(lm.Handle(testingHandler))
	req, err := http.NewRequest(http.MethodGet, testServer.URL, nil)
	assert.NoError(t, err)

	_, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	testServer.Close()

	logLines := recorder.String()
	assert.Contains(t, logLines, fmt.Sprintf("method=%s", req.Method), "cannot find method field in log")
	assert.Contains(t, logLines, fmt.Sprintf("path=%s", req.URL.Path), "cannot find path field in log")

}
