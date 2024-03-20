package web

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const (
	// HeaderRateLimitLimit and HeaderRateLimitRemaining are the recommended return header
	// values from IETF on rate limiting.
	HeaderRateLimitLimit     = "X-RateLimit-Limit"
	HeaderRateLimitRemaining = "X-RateLimit-Remaining"

	// HeaderRetryAfter is the header used to indicate when a client should retry
	// requests (when the rate limit expires), in UTC time.
	HeaderRetryAfter = "Retry-After"
)

type clientRateLimiterMiddleware struct {
	sync.Mutex
	limitersByClients  map[string]*clientLimiter
	cleanInterval      time.Duration
	inactivityDuration time.Duration
	rateLimit          float64
	burst              int
	limitHeader        string
}

type clientLimiter struct {
	rate.Limiter
	lastSeen time.Time
}

// NewClientRateLimiterMiddleware creates a middleware that will allow clients to make
// `rateLimitPerSeconds` requests per seconds. For each clients, the middleware creates
// a "token bucket" limiter of size `burst` which is implemented in "golang.org/x/time/rate".
// The middleware will cleanup the list of its clients every `cleanInterval` and remove
// clients inactive for longer than `inactivityDuration`.
func NewClientRateLimiterMiddleware(rateLimitPerSeconds float64, burst int) Middleware {
	m := &clientRateLimiterMiddleware{
		limitersByClients:  make(map[string]*clientLimiter),
		cleanInterval:      1 * time.Minute,
		inactivityDuration: 10 * time.Minute,
		rateLimit:          rateLimitPerSeconds,
		burst:              burst,

		// cache the computation of the limit header as a string as the value won't change
		limitHeader: strconv.FormatInt(int64(burst), 10),
	}
	go func() { // periodic cleaning of the limitersbyclients every cleanInterval
		for {
			time.Sleep(m.cleanInterval)
			m.Lock()
			for ip, cl := range m.limitersByClients {
				if time.Since(cl.lastSeen) > m.inactivityDuration {
					delete(m.limitersByClients, ip)
				}
			}
			m.Unlock()
		}
	}()
	return m
}

func (m *clientRateLimiterMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Use the remote IP to identify the client
		ip, err := GetRemoteAddr(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		clientID := ip.String()

		// Get the correct rate.Limiter for this client
		m.Lock()
		if _, found := m.limitersByClients[clientID]; !found {
			m.limitersByClients[clientID] = &clientLimiter{
				Limiter: *rate.NewLimiter(rate.Limit(m.rateLimit), m.burst),
			}
		}
		m.limitersByClients[clientID].lastSeen = time.Now()

		if !m.limitersByClients[clientID].Allow() {
			m.Unlock()

			w.Header().Set(HeaderRateLimitLimit, m.limitHeader)
			w.Header().Set(HeaderRateLimitRemaining, "0")
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		remaining := m.limitersByClients[clientID].Tokens()
		m.Unlock()

		w.Header().Set(HeaderRateLimitLimit, m.limitHeader)
		w.Header().Set(HeaderRateLimitRemaining, strconv.FormatFloat(remaining, 'f', 0, 64))

		next.ServeHTTP(w, r)
	})
}
