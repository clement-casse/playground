package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var (
	verison   = "0.0.1"
	startTime = time.Now()
)

func NewAPIRouter() http.Handler {
	apiMux := mux.NewRouter()

	// List API routes managed by the API
	apiMux.Methods(http.MethodGet).Path("/health").HandlerFunc(healthHandler)
	apiMux.Methods(http.MethodGet).Path("/status").HandlerFunc(statusHandler)

	// Set up the HTTP middlewares specifically for the API router
	apiMux.Use(mux.CORSMethodMiddleware(apiMux))
	return apiMux
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(&statusPayload{
		Version: verison,
		Uptime:  time.Since(startTime),
	})
}

type statusPayload struct {
	Version string        `json:"version"`
	Uptime  time.Duration `json:"uptime"`
}
