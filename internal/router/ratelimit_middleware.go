package router

import (
	"net/http"
	"strings"
	"sync"
	"time"

	ratelimiter "github.com/mikandro/rate-limiter"
)

// perIPRateLimiter manages rate limiters for each IP address
type perIPRateLimiter struct {
	limiters map[string]*ratelimiter.RateLimiter
	mu       sync.RWMutex
	capacity int
	rate     time.Duration
}

// newPerIPRateLimiter creates a new per-IP rate limiter
func newPerIPRateLimiter(capacity int, rate time.Duration) *perIPRateLimiter {
	return &perIPRateLimiter{
		limiters: make(map[string]*ratelimiter.RateLimiter),
		capacity: capacity,
		rate:     rate,
	}
}

// getLimiter returns or creates a rate limiter for the given IP
func (p *perIPRateLimiter) getLimiter(ip string) *ratelimiter.RateLimiter {
	p.mu.RLock()
	limiter, exists := p.limiters[ip]
	p.mu.RUnlock()

	if exists {
		return limiter
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	// Check again in case another goroutine created it
	if limiter, exists = p.limiters[ip]; exists {
		return limiter
	}

	limiter = ratelimiter.NewRateLimiter(p.capacity, p.rate)
	p.limiters[ip] = limiter
	return limiter
}

// extractIPAddress extracts the client's IP address from the request
func extractIPAddress(r *http.Request) string {
	// Check X-Forwarded-For header first
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	// RemoteAddr includes port, remove it
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}

// RateLimitMiddleware creates a Chi middleware that enforces per-IP rate limiting
func RateLimitMiddleware(capacity int, rate time.Duration) func(next http.Handler) http.Handler {
	limiter := newPerIPRateLimiter(capacity, rate)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := extractIPAddress(r)
			ipLimiter := limiter.getLimiter(ip)

			if !ipLimiter.Allow() {
				http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
