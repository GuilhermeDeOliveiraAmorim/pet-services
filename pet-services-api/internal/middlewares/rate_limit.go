package middlewares

import (
	"net/http"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiterConfig struct {
	RequestsPerMinute int
	BurstSize         int
}

type ClientInfo struct {
	tokens         int
	lastRefillTime time.Time
	mu             sync.Mutex
}

type RateLimiter struct {
	clients    map[string]*ClientInfo
	mu         sync.RWMutex
	maxTokens  int
	refillRate time.Duration
	burstSize  int
}

func NewRateLimiter(config RateLimiterConfig) *RateLimiter {
	if config.RequestsPerMinute <= 0 {
		config.RequestsPerMinute = 60
	}
	if config.BurstSize <= 0 {
		config.BurstSize = 10
	}

	rl := &RateLimiter{
		clients:    make(map[string]*ClientInfo),
		maxTokens:  config.RequestsPerMinute,
		refillRate: time.Minute / time.Duration(config.RequestsPerMinute),
		burstSize:  config.BurstSize,
	}

	go rl.cleanupInactiveClients()

	return rl
}

func (rl *RateLimiter) Allow(clientID string) bool {
	rl.mu.Lock()
	client, exists := rl.clients[clientID]
	if !exists {
		client = &ClientInfo{
			tokens:         rl.burstSize,
			lastRefillTime: time.Now(),
		}
		rl.clients[clientID] = client
	}
	rl.mu.Unlock()

	client.mu.Lock()
	defer client.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(client.lastRefillTime)
	tokensToAdd := int(elapsed / rl.refillRate)

	if tokensToAdd > 0 {
		client.tokens += tokensToAdd
		if client.tokens > rl.burstSize {
			client.tokens = rl.burstSize
		}
		client.lastRefillTime = now
	}

	if client.tokens > 0 {
		client.tokens--
		return true
	}

	return false
}

func (rl *RateLimiter) cleanupInactiveClients() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for clientID, client := range rl.clients {
			client.mu.Lock()
			if now.Sub(client.lastRefillTime) > 30*time.Minute {
				delete(rl.clients, clientID)
			}
			client.mu.Unlock()
		}
		rl.mu.Unlock()
	}
}

func RateLimitMiddleware(config RateLimiterConfig, logger logging.LoggerInterface) gin.HandlerFunc {
	limiter := NewRateLimiter(config)

	return func(c *gin.Context) {
		ctx := c.Request.Context()

		clientID := c.ClientIP()

		if !limiter.Allow(clientID) {
			problem := exceptions.NewProblemDetails(exceptions.TooManyRequests, exceptions.ErrorMessage{
				Title:  "Taxa de requisições excedida",
				Detail: "Você excedeu o limite de requisições permitidas. Tente novamente em alguns instantes.",
			})

			logger.LogWarning(ctx, "RateLimitMiddleware", "Rate limit excedido para IP: "+clientID, nil)

			c.Header("X-RateLimit-Limit", string(rune(config.RequestsPerMinute)))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("Retry-After", "60")

			c.AbortWithStatusJSON(http.StatusTooManyRequests, problem)
			return
		}

		c.Next()
	}
}

func DefaultRateLimitMiddleware(logger logging.LoggerInterface) gin.HandlerFunc {
	return RateLimitMiddleware(RateLimiterConfig{
		RequestsPerMinute: 60,
		BurstSize:         10,
	}, logger)
}

func StrictRateLimitMiddleware(logger logging.LoggerInterface) gin.HandlerFunc {
	return RateLimitMiddleware(RateLimiterConfig{
		RequestsPerMinute: 10,
		BurstSize:         3,
	}, logger)
}
