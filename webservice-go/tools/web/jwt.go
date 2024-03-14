package web

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var authHeaderRegex = regexp.MustCompile(`[^\s]+\s[^\s]`)

type JWTAuthMiddleware struct {
	handler http.Handler

	parser    *jwt.Parser
	secretKey []byte
}

func NewJWTAuthMiddleware(secretKey []byte) *JWTAuthMiddleware {
	return &JWTAuthMiddleware{
		handler:   nil,
		parser:    jwt.NewParser(),
		secretKey: secretKey,
	}
}

func (jm *JWTAuthMiddleware) Chain(handler http.Handler) http.Handler {
	jm.handler = handler
	return jm
}

func (jm *JWTAuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bearerToken := r.Header.Get("Authorization")
	if bearerToken == "" {
		http.Error(w, "authentication required", http.StatusUnauthorized)
		return
	} else if !authHeaderRegex.Match([]byte(bearerToken)) {
		http.Error(w, "header authorization not containing a bearer token", http.StatusBadRequest)
		return
	}
	reqToken := strings.Split(bearerToken, " ")[1]
	claims := &jwt.RegisteredClaims{}
	token, err := jm.parser.ParseWithClaims(reqToken, claims, func(_ *jwt.Token) (any, error) {
		return jm.secretKey, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) ||
			errors.Is(err, jwt.ErrTokenExpired) ||
			errors.Is(err, jwt.ErrTokenNotValidYet) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !token.Valid {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}
	jm.handler.ServeHTTP(w, r)
}
