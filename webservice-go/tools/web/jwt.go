package web

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var authHeaderRegex = regexp.MustCompile(`[^\s]+\s[^\s]`)

type jwtAuthMiddleware struct {
	parser    *jwt.Parser
	secretKey []byte
}

// verify http.Handler interface compliance
var _ http.Handler = (*jwtAuthMiddleware)(nil)

// NewJWTAuthMiddleware creates a JWTAuthMiddleware with the given secret key used to check
// requests signature. The middleware verifies that the requests addressed to the inner handler
// are signed with JWT without checking users.
func NewJWTAuthMiddleware(secretKey []byte) Middleware {
	return &jwtAuthMiddleware{
		parser:    jwt.NewParser(),
		secretKey: secretKey,
	}
}

func (m *jwtAuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		token, err := m.parser.ParseWithClaims(reqToken, claims, func(_ *jwt.Token) (any, error) {
			return m.secretKey, nil
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
		next.ServeHTTP(w, r)
	})
}
