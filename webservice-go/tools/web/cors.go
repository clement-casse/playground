package web

import (
	"net/http"

	"github.com/rs/cors"
)

// CORSMiddleware wraps github.com/rs/cors with no additional logic.
type CORSMiddleware struct {
	handler http.Handler
	cors    *cors.Cors
}

// NewCORSMiddleware
func NewCORSMiddleware(allowsOrigins ...string) *CORSMiddleware {
	if len(allowsOrigins) == 0 {
		allowsOrigins = []string{"*"}
	}
	return &CORSMiddleware{cors: cors.New(cors.Options{AllowedOrigins: allowsOrigins})}
}

func (cm *CORSMiddleware) Chain(handler http.Handler) http.Handler {
	cm.handler = handler
	return cm
}

func (cm *CORSMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cm.cors.HandlerFunc(w, r)
	if cm.cors.OriginAllowed(r) {
		cm.handler.ServeHTTP(w, r)
	} else {
		http.Error(w, "cors", http.StatusNoContent)
	}
}
