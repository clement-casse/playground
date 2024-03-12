package web

import (
	"net/http"

	"github.com/rs/cors"
)

// CORSMiddleware wraps github.com/rs/cors with no additional logic.
type CORSMiddleware struct {
	*cors.Cors
}

// NewCORSMiddleware
func NewCORSMiddleware(allowsOrigins ...string) *CORSMiddleware {
	if len(allowsOrigins) == 0 {
		allowsOrigins = []string{"*"}
	}
	return &CORSMiddleware{cors.New(cors.Options{AllowedOrigins: allowsOrigins})}
}

func (cm *CORSMiddleware) Chain(handler http.Handler) http.Handler {
	return cm.Handler(handler)
}
