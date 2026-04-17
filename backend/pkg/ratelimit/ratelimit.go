package ratelimit

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// IPRateLimiter implements per-IP rate limiting using token bucket.
type IPRateLimiter struct {
	limiters sync.Map
	rate     rate.Limit
	burst    int
}

// NewIPRateLimiter creates a rate limiter that allows `r` requests per second
// with a burst of `b` for each unique client IP.
// Example: NewIPRateLimiter(rate.Every(12*time.Second), 5) allows 5 requests per minute.
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		rate:  r,
		burst: b,
	}
}

func (l *IPRateLimiter) getLimiter(ip string) *rate.Limiter {
	if v, ok := l.limiters.Load(ip); ok {
		return v.(*rate.Limiter)
	}
	limiter := rate.NewLimiter(l.rate, l.burst)
	l.limiters.Store(ip, limiter)
	return limiter
}

// Middleware returns a Gin middleware that rate-limits by client IP.
func (l *IPRateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := l.getLimiter(ip)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "too many requests, please try again later",
			})
			return
		}

		c.Next()
	}
}
