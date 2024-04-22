package web

import (
	"net/http"

	"github.com/rs/cors"
)

type corsMiddleware struct {
	*cors.Cors
}

// verify Middleware interface compliance
var _ Middleware = (*corsMiddleware)(nil)

// NewCORSMiddleware creates a CORS Middleware that returns Forbidden when the headers
// do not match the same resource policy. It wraps github.com/rs/cors.
func NewCORSMiddleware(allowsOrigins ...string) Middleware {
	if len(allowsOrigins) == 0 {
		allowsOrigins = []string{"*"}
	}
	return &corsMiddleware{
		Cors: cors.New(cors.Options{AllowedOrigins: allowsOrigins}),
	}
}

func (m *corsMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.HandlerFunc(w, r)
		if m.OriginAllowed(r) {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "cors", http.StatusForbidden)
		}
	})
}
