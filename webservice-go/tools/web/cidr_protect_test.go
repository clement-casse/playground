package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCIDRProtectMiddleware(t *testing.T) {
	for _, tt := range []struct {
		name            string
		allowedNetworks []string
		remoteAddr      string
		expectStatus    int
	}{
		{
			name:            "public address should not reach inner middleware",
			allowedNetworks: []string{},
			remoteAddr:      "8.8.4.4:16781",
			expectStatus:    http.StatusUnauthorized,
		}, {
			name:            "loopback should be allowed enventhough allowedNetworks is empty",
			allowedNetworks: []string{},
			remoteAddr:      "127.0.0.1:21345",
			expectStatus:    http.StatusOK,
		}, {
			name:            "allowing tailscale network for example",
			allowedNetworks: []string{"100.64.0.0/10"},
			remoteAddr:      "100.108.113.29:21345",
			expectStatus:    http.StatusOK,
		}, {
			name:            "allowing tailscale network should still allow localhost",
			allowedNetworks: []string{"100.64.0.0/10"},
			remoteAddr:      "127.0.0.1:21345",
			expectStatus:    http.StatusOK,
		}, {
			name:            "allowing private ranges first then tailscale network should still allow the request",
			allowedNetworks: []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "100.64.0.0/10"},
			remoteAddr:      "100.108.113.29:21345",
			expectStatus:    http.StatusOK,
		}, {
			name:            "disallowing public ip when only allowingprivate ranges",
			allowedNetworks: []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "100.64.0.0/10"},
			remoteAddr:      "8.8.4.4:16781",
			expectStatus:    http.StatusUnauthorized,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			cpm := NewCIDRProtectMiddleware(tt.allowedNetworks...)
			chainedHandler := cpm.Chain(testingHandler)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.RemoteAddr = tt.remoteAddr

			chainedHandler.ServeHTTP(w, req)

			res := w.Result()
			assert.Equal(t, tt.expectStatus, res.StatusCode)
		})
	}
}
