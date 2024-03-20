package web

import "net/http"

// RateLimiterMiddleware
type RateLimiterMiddleware struct {
}

// NewRateLimiterMiddleware
func NewRateLimiterMiddleware() *RateLimiterMiddleware {
	rlm := &RateLimiterMiddleware{}

	return rlm
}

func (rlm *RateLimiterMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ip, err := GetRemoteAddr(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		_, _, _ = ctx, ip, next
	})
}
