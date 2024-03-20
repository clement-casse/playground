package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

var (
	secretKey = []byte("wouiiiiiiiiiiiiiii")

	validClaims = &jwt.RegisteredClaims{
		Issuer:    "test",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
	}
	validToken = jwt.NewWithClaims(jwt.SigningMethodHS256, validClaims)

	expiredClaims = &jwt.RegisteredClaims{
		Issuer:    "test",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-10 * time.Minute)),
	}
	expiredToken = jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
)

func TestJWTVerification(t *testing.T) {
	validSignedString, err := validToken.SignedString(secretKey)
	assert.NoError(t, err)

	invalidSiggnedString, err := validToken.SignedString([]byte("invalid-secret"))
	assert.NoError(t, err)

	expiredSignedString, err := expiredToken.SignedString(secretKey)
	assert.NoError(t, err)

	for _, tt := range []struct {
		name         string
		reqHeaders   map[string]string
		expectStatus int
	}{
		{
			name:         "No Authorization header",
			reqHeaders:   map[string]string{},
			expectStatus: http.StatusUnauthorized,
		}, {
			name:         "Authorization header malformed",
			reqHeaders:   map[string]string{"Authorization": "BadlyFormattedHeader"},
			expectStatus: http.StatusBadRequest,
		}, {
			name:         "Authorization header with malformed bearer",
			reqHeaders:   map[string]string{"Authorization": "Bearer not-a-jwt"},
			expectStatus: http.StatusBadRequest,
		}, {
			name:         "Authorization header with a correct token",
			reqHeaders:   map[string]string{"Authorization": "Bearer " + validSignedString},
			expectStatus: http.StatusOK,
		}, {
			name:         "Authorization header with an expired token",
			reqHeaders:   map[string]string{"Authorization": "Bearer " + expiredSignedString},
			expectStatus: http.StatusUnauthorized,
		}, {
			name:         "Authorization header with an invalid token",
			reqHeaders:   map[string]string{"Authorization": "Bearer " + invalidSiggnedString},
			expectStatus: http.StatusUnauthorized,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			jm := NewJWTAuthMiddleware(secretKey)
			testServer := httptest.NewServer(jm.Handle(testingHandler))
			defer testServer.Close()
			req, err := http.NewRequest(http.MethodGet, testServer.URL, nil)
			assert.NoError(t, err)
			for rhKey, rhValue := range tt.reqHeaders {
				req.Header.Set(rhKey, rhValue)
			}

			res, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectStatus, res.StatusCode)
		})
	}

}
