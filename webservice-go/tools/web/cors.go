package web

import (
	"net/http"

	"github.com/rs/cors"
)

// CORSMiddleware wraps github.com/rs/cors, it returns Forbidden when the headers do not match the same resource policy.
type CORSMiddleware struct {
	*cors.Cors
}

// NewCORSMiddleware
func NewCORSMiddleware(allowsOrigins ...string) Middleware {
	if len(allowsOrigins) == 0 {
		allowsOrigins = []string{"*"}
	}
	return &CORSMiddleware{
		Cors: cors.New(cors.Options{AllowedOrigins: allowsOrigins}),
	}
}

func (m *CORSMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.HandlerFunc(w, r)
		if m.OriginAllowed(r) {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "cors", http.StatusForbidden)
		}
	})
}
