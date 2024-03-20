package web

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRemoteAddr(t *testing.T) {
	for _, tt := range []struct {
		name        string
		reqHeaders  map[string]string
		remoteAddr  string
		expectIP    net.IP
		expectError bool
	}{
		{
			name:       "no request headers",
			reqHeaders: map[string]string{},
			remoteAddr: "8.8.4.4:16781",
			expectIP:   net.ParseIP("8.8.4.4"),
		}, {
			name:       "X-Real-IP request header",
			reqHeaders: map[string]string{"X-Real-Ip": "1.1.1.1"},
			remoteAddr: "8.8.4.4:16781",
			expectIP:   net.ParseIP("1.1.1.1"),
		}, {
			name:       "X-Forwarded-For request header",
			reqHeaders: map[string]string{"X-Forwarded-For": "1.1.1.1"},
			remoteAddr: "8.8.4.4:16781",
			expectIP:   net.ParseIP("1.1.1.1"),
		}, {
			name:       "X-Forwarded-For request header with multiple values",
			reqHeaders: map[string]string{"X-Forwarded-For": "1.1.1.1, 2.2.2.2, 4.4.4.4, 8.8.8.8"},
			remoteAddr: "8.8.4.4:16781",
			expectIP:   net.ParseIP("8.8.8.8"),
		}, {
			name:       "X-Forwarded-For request header with multiple values and a private address",
			reqHeaders: map[string]string{"X-Forwarded-For": "1.1.1.1, 2.2.2.2, 4.4.4.4, 8.8.8.8, 10.10.10.10"},
			remoteAddr: "8.8.4.4:16781",
			expectIP:   net.ParseIP("8.8.8.8"),
		}, {
			name: "Both X-Real-Ip and X-Forwarded-For request headers",
			reqHeaders: map[string]string{
				"X-Real-Ip":       "1.1.1.1",
				"X-Forwarded-For": "1.1.1.1, 2.2.2.2, 4.4.4.4, 8.8.8.8",
			},
			remoteAddr: "8.8.4.4:16781",
			expectIP:   net.ParseIP("8.8.8.8"),
		}, {
			name: "Both X-Real-Ip and X-Forwarded-For request headers with privates xff",
			reqHeaders: map[string]string{
				"X-Real-Ip":       "1.1.1.1",
				"X-Forwarded-For": "172.18.19.20, 192.168.104.2",
			},
			remoteAddr: "8.8.4.4:16781",
			expectIP:   net.ParseIP("1.1.1.1"),
		}, {
			name:       "X-Forwarded-For request headers with privates xff",
			reqHeaders: map[string]string{"X-Forwarded-For": "172.18.19.20, 192.168.104.2"},
			remoteAddr: "8.8.4.4:16781",
			expectIP:   net.ParseIP("8.8.4.4"),
		}, {
			name:       "no request headers ipv6",
			reqHeaders: map[string]string{},
			remoteAddr: "[2a01:cb19:85a0:4d00:7016:6357:9140:f300]:16781",
			expectIP:   net.ParseIP("2a01:cb19:85a0:4d00:7016:6357:9140:f300"),
		}, {
			name:       "X-Real-IP request header ipv6",
			reqHeaders: map[string]string{"X-Real-Ip": "2a01:cb19:85a0:4d00:7016:6357:9140:f300"},
			remoteAddr: "8.8.4.4:16781",
			expectIP:   net.ParseIP("2a01:cb19:85a0:4d00:7016:6357:9140:f300"),
		}, {
			name:       "X-Forwarded-For request header ipv6",
			reqHeaders: map[string]string{"X-Forwarded-For": "2a01:cb19:85a0:4d00:7016:6357:9140:f300"},
			remoteAddr: "8.8.4.4:16781",
			expectIP:   net.ParseIP("2a01:cb19:85a0:4d00:7016:6357:9140:f300"),
		}, {
			name:       "X-Forwarded-For request header multiple values ipv6",
			reqHeaders: map[string]string{"X-Forwarded-For": "2a01:cb19:85a0:4d00:7016:6357:9140:f300, fc00::abcd"},
			remoteAddr: "8.8.4.4:16781",
			expectIP:   net.ParseIP("2a01:cb19:85a0:4d00:7016:6357:9140:f300"),
		}, {
			name:       "X-Forwarded-For with bad invalid value falls back to remoteAddr",
			reqHeaders: map[string]string{"X-Forwarded-For": "2a01::abcdef"},
			remoteAddr: "8.8.4.4:16781",
			expectIP:   net.ParseIP("8.8.4.4"),
		}, {
			name:        "invalid remote address with no headers",
			remoteAddr:  "11.12.13.256:21345",
			expectError: true,
		}, {
			name:        "invalid remote address with no headers",
			remoteAddr:  "[2a01::abcd];21345",
			expectError: true,
		}, {
			name:        "No remote address with no headers",
			expectError: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.RemoteAddr = tt.remoteAddr
			for rhKey, rhValue := range tt.reqHeaders {
				req.Header.Set(rhKey, rhValue)
			}

			ip, err := GetRemoteAddr(req)
			assert.Equal(t, tt.expectError, err != nil)
			assert.Equal(t, tt.expectIP, ip)
		})
	}
}
