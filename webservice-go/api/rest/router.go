package rest

import (
	"net/http"
)

// Router returns the http handler for serving REST API Server routes
func (s *APIController) Router() http.Handler {
	s.registerRoute("GET /api/somepath/...", nil)
	return s.mux
}
