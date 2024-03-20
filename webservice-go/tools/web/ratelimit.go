package web

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// ClientRateLimiterMiddleware
type ClientRateLimiterMiddleware struct {
	sync.Mutex
	limitersByClients  map[string]*clientLimiter
	cleanInterval      time.Duration
	inactivityDuration time.Duration
	rateLimit          rate.Limit
	burst              int
}

type clientLimiter struct {
	rate.Limiter
	lastSeen time.Time
}

// NewClientRateLimiterMiddleware
func NewClientRateLimiterMiddleware(rateLimit float64, burst int) Middleware {
	m := &ClientRateLimiterMiddleware{
		limitersByClients:  make(map[string]*clientLimiter),
		cleanInterval:      1 * time.Minute,
		inactivityDuration: 10 * time.Minute,
		rateLimit:          rate.Limit(rateLimit),
		burst:              burst,
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

func (m *ClientRateLimiterMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, err := GetRemoteAddr(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		clientID := ip.String()
		m.Lock()
		if _, found := m.limitersByClients[clientID]; !found {
			m.limitersByClients[clientID] = &clientLimiter{
				Limiter: *rate.NewLimiter(m.rateLimit, m.burst),
			}
		}
		m.limitersByClients[clientID].lastSeen = time.Now()
		if !m.limitersByClients[clientID].Allow() {
			m.Unlock()
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		m.Unlock()
		next.ServeHTTP(w, r)
	})
}
