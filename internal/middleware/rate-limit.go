package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type clientLimit struct {
	requests int
	resetAt  time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	limit   int
	window  time.Duration
	clients map[string]*clientLimit
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		limit:   limit,
		window:  window,
		clients: make(map[string]*clientLimit),
	}
}

func (limiter *RateLimiter) RateLimiterMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientIP := ctx.ClientIP()
		now := time.Now()

		limiter.mu.Lock()

		client, exists := limiter.clients[clientIP]

		if !exists || now.After(client.resetAt) {
			client = &clientLimit{
				requests: 0,
				resetAt:  now.Add(limiter.window),
			}

			limiter.clients[clientIP] = client
		}

		if client.requests >= limiter.limit {
			retryAfter := int(time.Until(client.resetAt).Seconds())

			if retryAfter < 1 {
				retryAfter = 1
			}

			limiter.mu.Unlock()

			ctx.Header("X-Limit-Remaining", "0")
			ctx.Header("Retry-After", strconv.Itoa(retryAfter))

			ctx.JSON(http.StatusTooManyRequests, gin.H{
				"error": "слишком много запросов, попробуйте позже",
			})

			ctx.Abort()
			return
		}

		client.requests++

		remaining := limiter.limit - client.requests

		limiter.mu.Unlock()

		ctx.Header("X-Limit-Remaining", strconv.Itoa(remaining))

		ctx.Next()
	}
}
