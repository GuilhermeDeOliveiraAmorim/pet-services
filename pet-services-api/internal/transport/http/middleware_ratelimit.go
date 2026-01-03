package http

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter holds per-IP rate limiters
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
}

// NewRateLimiter creates a new rate limiter manager
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
	}
}

// getLimiter returns the rate limiter for a given IP, creating one if needed
func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if limiter, exists := rl.limiters[ip]; exists {
		return limiter
	}

	// 100 requests per minute per IP (100/60 ≈ 1.67 req/sec)
	limiter := rate.NewLimiter(rate.Limit(100.0/60.0), 10)
	rl.limiters[ip] = limiter
	return limiter
}

// RateLimitMiddleware returns a Gin middleware that applies rate limiting per IP
func (rl *RateLimiter) RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client IP (respects X-Forwarded-For behind proxies)
		ip := c.ClientIP()

		// Get or create limiter for this IP
		limiter := rl.getLimiter(ip)

		// Check if we can process this request
		if !limiter.Allow() {
			c.AbortWithStatusJSON(
				http.StatusTooManyRequests,
				errorResponse("rate_limit_exceeded", "too many requests, try again later"),
			)
			return
		}

		c.Next()
	}
}

// CriticalEndpointRateLimiter creates a stricter limiter for sensitive endpoints
// (e.g., auth, password reset) - 20 requests per minute per IP
func CriticalEndpointRateLimiter() gin.HandlerFunc {
	limiters := make(map[string]*rate.Limiter)
	mu := sync.RWMutex{}

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.RLock()
		limiter, exists := limiters[ip]
		mu.RUnlock()

		if !exists {
			mu.Lock()
			// 20 requests per minute per IP (stricter for auth endpoints)
			limiter = rate.NewLimiter(rate.Limit(20.0/60.0), 5)
			limiters[ip] = limiter
			mu.Unlock()
		}

		if !limiter.Allow() {
			c.AbortWithStatusJSON(
				http.StatusTooManyRequests,
				errorResponse("rate_limit_exceeded", "too many attempts, try again later"),
			)
			return
		}

		c.Next()
	}
}
