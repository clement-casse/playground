package web

import "net/http"

// MiddlewareChainer represents HTTP Middleware that can be chained to compose HTTP Handlers for the web server.
type MiddlewareChainer interface {
	// Chain allows chaining middlewares by having a function(handler)->handler
	Chain(http.Handler) http.Handler
}
