package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCORSMiddleware(t *testing.T) {
	for _, tt := range []struct {
		name           string
		allowedOrigins []string
		reqHeaders     map[string]string
		resHeaders     map[string]string
		expectStatus   int
	}{
		{
			name:           "No allowed origins, no request headers",
			allowedOrigins: []string{},
			reqHeaders:     map[string]string{},
			resHeaders:     map[string]string{"Vary": "Origin"},
			expectStatus:   http.StatusOK,
		}, {
			name:           "No allowed origins should allow all origins",
			allowedOrigins: []string{}, // should then be "*"
			reqHeaders:     map[string]string{"Origin": "http://example.com"},
			resHeaders:     map[string]string{"Vary": "Origin", "Access-Control-Allow-Origin": "*"},
			expectStatus:   http.StatusOK,
		}, {
			name:           "with allowed origins should allow this origin",
			allowedOrigins: []string{"http://example.com"},
			reqHeaders:     map[string]string{"Origin": "http://example.com"},
			resHeaders:     map[string]string{"Vary": "Origin", "Access-Control-Allow-Origin": "http://example.com"},
			expectStatus:   http.StatusOK,
		}, {
			name:           "with allowed origins should disallow other origins",
			allowedOrigins: []string{"http://example.com"},
			reqHeaders:     map[string]string{"Origin": "http://not-allowed.com"},
			resHeaders:     map[string]string{"Vary": "Origin"},
			expectStatus:   http.StatusForbidden,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			cm := NewCORSMiddleware(tt.allowedOrigins...)
			testServer := httptest.NewServer(cm.Handle(testingHandler))
			defer testServer.Close()
			req, err := http.NewRequest(http.MethodGet, testServer.URL, nil)
			assert.NoError(t, err)
			for rhKey, rhValue := range tt.reqHeaders {
				req.Header.Set(rhKey, rhValue)
			}
			res, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			for rhKey, rhValue := range tt.resHeaders {
				assert.Equal(t, rhValue, res.Header.Get(rhKey))
			}
			assert.Equal(t, tt.expectStatus, res.StatusCode)
		})
	}
}
