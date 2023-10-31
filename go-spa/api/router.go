package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// NewAPIRouter creates a new http.Handler that will handle API requests
func NewAPIRouter() http.Handler {
	apiMux := mux.NewRouter()

	// List API routes managed by the API
	apiMux.Methods(http.MethodGet).Path("/health").HandlerFunc(healthHandler)

	// Set up the HTTP middlewares specifically for the API router
	apiMux.Use(mux.CORSMethodMiddleware(apiMux))
	return apiMux
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
