package rest

import (
	"net/http"

	"github.com/clement-casse/playground/webservice-go/tools/web"
)

// Router returns the http handler for serving REST API Server routes
func (c *APIController) Router() http.Handler {
	jwtMw := web.NewJWTAuthMiddleware(c.secret)

	c.registerRoute("GET /api/somepath/...", nil, jwtMw)
	return c.mux
}
