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
			name:            "well formated cidr range should not panic",
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
		reqHeaders      map[string]string
		remoteAddr      string
		expectStatus    int
	}{
		{
			name:            "public address should not reach inner middleware",
			allowedNetworks: []string{},
			remoteAddr:      "8.8.4.4:16781",
			expectStatus:    http.StatusUnauthorized,
		}, {
			name:            "loopback should be allowed eventhough allowedNetworks is empty",
			allowedNetworks: []string{},
			remoteAddr:      "127.0.0.1:21345",
			expectStatus:    http.StatusOK,
		}, {
			name:            "X-Real-Ip should take prcedence over request remoteAddr",
			allowedNetworks: []string{},
			reqHeaders:      map[string]string{"X-Real-Ip": "8.8.4.4"},
			remoteAddr:      "127.0.0.1:21345",
			expectStatus:    http.StatusUnauthorized,
		}, {
			name:            "X-Forwarded-For should take prcedence over request remoteAddr",
			allowedNetworks: []string{},
			reqHeaders:      map[string]string{"X-Forwarded-For": "8.8.4.4"},
			remoteAddr:      "127.0.0.1:21345",
			expectStatus:    http.StatusUnauthorized,
		}, {
			name:            "X-Forwarded-For should discard private ips",
			allowedNetworks: []string{"192.168.0.0/16"},
			reqHeaders:      map[string]string{"X-Forwarded-For": "8.8.4.4, 192.168.104.2"},
			remoteAddr:      "127.0.0.1:21345",
			expectStatus:    http.StatusUnauthorized,
		}, {
			name:            "X-Forwarded-For containing only private IP should fall back to remoteAddr",
			allowedNetworks: []string{"10.0.0.0/8"},
			reqHeaders:      map[string]string{"X-Forwarded-For": "172.18.19.20, 192.168.104.2"},
			remoteAddr:      "10.0.1.254:21345",
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
		}, {
			name:            "public IPv6 address should not reach inner middleware",
			allowedNetworks: []string{},
			remoteAddr:      "[2a01:cb19:85a0:4d00:7016:6357:9140:f300]:16781",
			expectStatus:    http.StatusUnauthorized,
		}, {
			name:            "loopback IPv6 should be allowed eventhough allowedNetworks is empty",
			allowedNetworks: []string{},
			remoteAddr:      "[::1]:21345",
			expectStatus:    http.StatusOK,
		}, {
			name:            "allowing IPv6 Unique Local Addressing should allow a local ipv6 address",
			allowedNetworks: []string{"fc00::/7"},
			remoteAddr:      "[fc00::abcd]:21345",
			expectStatus:    http.StatusOK,
		}, {
			name:            "allowing IPv6 Unique Local Addressing should disallow not local ipv6 addresses",
			allowedNetworks: []string{"fc00::/7"},
			remoteAddr:      "[2a01::abcd]:21345",
			expectStatus:    http.StatusUnauthorized,
		}, {
			name:            "malformed RemoteAddress IP should trigger a 500",
			allowedNetworks: []string{"fc00::/7"},
			remoteAddr:      "[2a01::abcdefgh]:21345",
			expectStatus:    http.StatusInternalServerError,
		}, {
			name:            "malformed RemoteAddress field should trigger a 500",
			allowedNetworks: []string{"fc00::/7"},
			remoteAddr:      "[2a01::abcd];21345",
			expectStatus:    http.StatusInternalServerError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			cpm := NewCIDRProtectMiddleware(tt.allowedNetworks...)
			protectededHandler := cpm.Chain(testingHandler)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.RemoteAddr = tt.remoteAddr
			for rhKey, rhValue := range tt.reqHeaders {
				req.Header.Set(rhKey, rhValue)
			}

			protectededHandler.ServeHTTP(w, req)

			res := w.Result()
			assert.Equal(t, tt.expectStatus, res.StatusCode)
		})
	}
}
