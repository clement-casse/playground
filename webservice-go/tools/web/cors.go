package web

import (
	"net/http"

	"github.com/rs/cors"
)

// CORSMiddleware
type CORSMiddleware struct {
	handler http.Handler
	cors    *cors.Cors
}

// NewCORSMiddleware
func NewCORSMiddleware(allowsOrigins ...string) *CORSMiddleware {
	if len(allowsOrigins) == 0 {
		allowsOrigins = []string{"*"}
	}
	return &CORSMiddleware{nil, cors.New(cors.Options{AllowedOrigins: allowsOrigins})}
}

func (cm *CORSMiddleware) Chain(handler http.Handler) http.Handler {
	return cm.cors.Handler(handler)
}
