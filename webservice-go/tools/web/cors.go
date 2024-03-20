package web

import (
	"net/http"

	"github.com/rs/cors"
)

// CORSMiddleware wraps github.com/rs/cors, it returns Forbidden when the headers do not match the same resource policy.
type CORSMiddleware struct {
	cors *cors.Cors
}

// NewCORSMiddleware
func NewCORSMiddleware(allowsOrigins ...string) *CORSMiddleware {
	if len(allowsOrigins) == 0 {
		allowsOrigins = []string{"*"}
	}
	return &CORSMiddleware{cors: cors.New(cors.Options{AllowedOrigins: allowsOrigins})}
}

func (cm *CORSMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cm.cors.HandlerFunc(w, r)
		if cm.cors.OriginAllowed(r) {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "cors", http.StatusForbidden)
		}
	})
}
