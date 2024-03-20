package web

import "net/http"

// RateLimiterMiddleware
type RateLimiterMiddleware struct {
	handler http.Handler
}

// NewRateLimiterMiddleware
func NewRateLimiterMiddleware() *RateLimiterMiddleware {
	rlm := &RateLimiterMiddleware{
		handler: nil,
	}

	return rlm
}

func (rlm *RateLimiterMiddleware) Chain(handler http.Handler) http.Handler {
	rlm.handler = handler
	return rlm
}

func (rlm *RateLimiterMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ip, err := GetRemoteAddr(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	_, _ = ctx, ip
}
