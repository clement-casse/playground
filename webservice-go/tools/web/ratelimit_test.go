package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRateLimitMiddleware(t *testing.T) {
	crlm := NewClientRateLimiterMiddleware(5.0, 10)
	testServer := httptest.NewServer(crlm.Handle(testingHandler))
	defer testServer.Close()

	// Consume a full burst of 10 requests
	for range 10 {
		req, err := http.NewRequest(http.MethodGet, testServer.URL, nil)
		assert.NoError(t, err)
		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	}
	// Then the following requests should be discarded with status 429
	for range 3 {
		req, err := http.NewRequest(http.MethodGet, testServer.URL, nil)
		assert.NoError(t, err)
		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusTooManyRequests, res.StatusCode)
	}
	// Cooldown for regaining 5 requests (rateLimitPerSeconds)
	time.Sleep(1 * time.Second)
	for range 5 {
		req, err := http.NewRequest(http.MethodGet, testServer.URL, nil)
		assert.NoError(t, err)
		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	}
	// All requests are consumed again then remaining requests should be 429
	for range 3 {
		req, err := http.NewRequest(http.MethodGet, testServer.URL, nil)
		assert.NoError(t, err)
		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusTooManyRequests, res.StatusCode)
	}

}
