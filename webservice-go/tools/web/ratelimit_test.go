package web

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	maxNumberOfRequests                 = 10
	numberOfRequestRegainedEverySeconds = 5.0
)

func TestRateLimitMiddleware(t *testing.T) {
	crlm := NewClientRateLimiterMiddleware(numberOfRequestRegainedEverySeconds, maxNumberOfRequests)
	testServer := httptest.NewServer(crlm.Handle(testingHandler))
	defer testServer.Close()

	// Consume a full burst of 10 requests
	for i := range maxNumberOfRequests {
		req, err := http.NewRequest(http.MethodGet, testServer.URL, nil)
		assert.NoError(t, err)
		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "10", res.Header.Get(HeaderRateLimitLimit))
		assert.Equal(t, fmt.Sprintf("%d", maxNumberOfRequests-i-1), res.Header.Get(HeaderRateLimitRemaining))
	}
	// Then the following requests should be discarded with status 429
	for range 3 {
		req, err := http.NewRequest(http.MethodGet, testServer.URL, nil)
		assert.NoError(t, err)
		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusTooManyRequests, res.StatusCode)
		assert.Equal(t, "10", res.Header.Get(HeaderRateLimitLimit))
		assert.Equal(t, "0", res.Header.Get(HeaderRateLimitRemaining))
	}
	// making a new client should not limit the requests
	for i := range 3 {
		req, err := http.NewRequest(http.MethodGet, testServer.URL, nil)
		assert.NoError(t, err)
		req.Header.Set("X-Forwarded-For", "8.8.4.4")
		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "10", res.Header.Get(HeaderRateLimitLimit))
		assert.Equal(t, fmt.Sprintf("%d", maxNumberOfRequests-i-1), res.Header.Get(HeaderRateLimitRemaining))
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
