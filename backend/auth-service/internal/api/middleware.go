package api

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"korun.io/auth-service/internal/infrastructure/metrics"
)

type tokenBucket struct {
	lastRefillTime time.Time
	tokens         float64
	capacity       float64
	fillRate       float64 // tokens per second
	mu             sync.Mutex
}

func newTokenBucket(capacity, fillRate float64) *tokenBucket {
	return &tokenBucket{
		lastRefillTime: time.Now(),
		tokens:         capacity,
		capacity:       capacity,
		fillRate:       fillRate,
	}
}

func (tb *tokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefillTime).Seconds()
	tb.lastRefillTime = now

	tb.tokens = tb.tokens + elapsed*tb.fillRate
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}

	if tb.tokens >= 1.0 {
		tb.tokens = tb.tokens - 1.0
		return true
	}

	return false
}

// Simple in-memory store for rate limiting.
// For a distributed system, a shared store like Redis would be more appropriate.
var (
	clients = make(map[string]*tokenBucket)
	mu      sync.Mutex
)

func getVisitorLimiter(ip string) *tokenBucket {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := clients[ip]
	if !exists {
		// Allow 5 requests per second with a burst of 10.
		// For simplicity, capacity is 10, fillRate is 5 tokens/sec.
		limiter = newTokenBucket(10, 5)
		clients[ip] = limiter
	}

	return limiter
}

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		limiter := getVisitorLimiter(c.ClientIP())
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			return
		}
		c.Next()
	}
}

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath() // use template path

		c.Next() // process request

		duration := time.Since(start)
		statusCode := c.Writer.Status()

		// Record metrics
		metrics.HttpRequestDuration.WithLabelValues(path, c.Request.Method).Observe(duration.Seconds())
		metrics.HttpRequestsTotal.WithLabelValues(path, c.Request.Method, strconv.Itoa(statusCode)).Inc()
	}
}
