package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCIDRProtectMiddleware(t *testing.T) {
	for _, tt := range []struct {
		name            string
		allowedNetworks []string
		expectPanic     bool
	}{
		{
			name:            "an empty allowed network should not panic",
			allowedNetworks: []string{},
			expectPanic:     false,
		}, {
			name:            "come well formated cidr range should not panic",
			allowedNetworks: []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "100.64.0.0/10"},
			expectPanic:     false,
		}, {
			name:            "a not parsable CIDR range should panic",
			allowedNetworks: []string{"not a parsable cidr range"},
			expectPanic:     true,
		}, {
			name:            "one allowed network not being a parsable CIDR range should panic",
			allowedNetworks: []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.O.O/16", "100.64.0.0/10"}, // some O instead of 0 in 192.168/16
			expectPanic:     true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectPanic {
				assert.Panics(t, func() { _ = NewCIDRProtectMiddleware(tt.allowedNetworks...) })
			} else {
				assert.NotPanics(t, func() { _ = NewCIDRProtectMiddleware(tt.allowedNetworks...) })
			}
		})
	}
}

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
