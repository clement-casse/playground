package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func NewAPIRouter() http.Handler {
	apiMux := mux.NewRouter()

	apiMux.HandleFunc("/health", healthHandler).Methods(http.MethodGet)

	apiMux.Use(
		mux.CORSMethodMiddleware(apiMux),
	)
	return apiMux
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
