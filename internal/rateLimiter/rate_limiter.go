package rateLimiter

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.Mutex
	rate     int           // max requests
	window   time.Duration // time window
}

type Visitor struct {
	lastSeen time.Time
	tokens   int
}

// NewLimiter creates a new IP rate limiter
func NewLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
		rate:     rate,
		window:   window,
	}
	go rl.cleanupVisitors()
	return rl
}

// Middleware function
func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		rl.mu.Lock()
		visitor, exists := rl.visitors[ip]
		now := time.Now()
		if !exists || now.Sub(visitor.lastSeen) > rl.window {
			rl.visitors[ip] = &Visitor{
				lastSeen: now,
				tokens:   rl.rate - 1,
			}
			rl.mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		if visitor.tokens <= 0 {
			rl.mu.Unlock()
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		visitor.tokens--
		visitor.lastSeen = now
		rl.mu.Unlock()
		next.ServeHTTP(w, r)
	})
}

// cleanupVisitors prunes old IPs
func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > rl.window {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}