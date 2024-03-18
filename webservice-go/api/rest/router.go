package rest

import (
	"net/http"
)

// Router returns the http handler for serving REST API Server routes
func (c *APIController) Router() http.Handler {
	c.registerRoute("GET /api/somepath/...", nil)
	return c.mux
}
