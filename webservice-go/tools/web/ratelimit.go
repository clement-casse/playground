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

	defaultCleanInterval      = 1 * time.Minute
	defaultInactivityDuration = 5 * time.Minute
)

type clientRateLimiterMiddleware struct {
	sync.Mutex
	limitersByClients map[string]*clientLimiter

	rateLimit          float64
	burst              int
	cleanInterval      time.Duration
	inactivityDuration time.Duration

	limitHeader string
}

// verify http.Handler interface compliance
var _ http.Handler = (*clientRateLimiterMiddleware)(nil)

type clientLimiter struct {
	rate.Limiter
	lastSeen time.Time
}

// RateLimiterOpt in an interface for applying RateLimiterMiddleware options.
type RateLimiterOpt interface {
	applyOpt(*clientRateLimiterMiddleware) *clientRateLimiterMiddleware
}

type rateLimiterOptFunc func(*clientRateLimiterMiddleware) *clientRateLimiterMiddleware

func (fn rateLimiterOptFunc) applyOpt(s *clientRateLimiterMiddleware) *clientRateLimiterMiddleware {
	return fn(s)
}

// NewClientRateLimiterMiddleware creates a middleware that will allow clients to make
// `rateLimitPerSeconds` requests per seconds. For each clients, the middleware creates
// a "token bucket" limiter of size `burst` which is implemented in "golang.org/x/time/rate".
// The middleware will cleanup the list of its clients every `cleanInterval` and remove
// clients inactive for longer than `inactivityDuration`.
func NewClientRateLimiterMiddleware(rateLimitPerSeconds float64, burst int, opts ...RateLimiterOpt) Middleware {
	m := &clientRateLimiterMiddleware{
		limitersByClients:  make(map[string]*clientLimiter),
		cleanInterval:      defaultCleanInterval,
		inactivityDuration: defaultInactivityDuration,
		rateLimit:          rateLimitPerSeconds,
		burst:              burst,

		// cache the computation of the limit header as a string as the value won't change
		limitHeader: strconv.FormatInt(int64(burst), 10),
	}
	for _, opt := range opts {
		m = opt.applyOpt(m)
	}
	go func() { // periodic cleaning of the limitersByClients every cleanInterval
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

// WithCleanInterval configure a NewRateLimiterMiddleware by setting the cleaning
// interval to the specified value (by default to 1 minute).
func WithCleanInterval(d time.Duration) RateLimiterOpt {
	return rateLimiterOptFunc(func(m *clientRateLimiterMiddleware) *clientRateLimiterMiddleware {
		m.cleanInterval = d
		return m
	})
}

// WithInactivityDuration configures a NewRateLimiterMiddleware by setting the
// client inactivity duration to the desired value (by default to 5 minutes).
func WithInactivityDuration(d time.Duration) RateLimiterOpt {
	return rateLimiterOptFunc(func(m *clientRateLimiterMiddleware) *clientRateLimiterMiddleware {
		m.inactivityDuration = d
		return m
	})
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
