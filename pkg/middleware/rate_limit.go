package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hafiztri123/internal/core/response"
)

type visitor struct {
	count int
	lastSeen time.Time
}

type RateLimiter struct {
	visitors map[string]*visitor
	mu sync.Mutex
	rate int
	per time.Duration
}

func NewRateLimiter(rate int, per time.Duration) *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]*visitor),
		rate: rate,
		per: per,
	}
}

func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	threshold := time.Now().Add(-rl.per)
	for ip, v := range rl.visitors {
		if v.lastSeen.Before(threshold) {
			delete(rl.visitors, ip)
		}
	}
}

func RateLimit(rate int, per time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(rate, per)

	go func(){
		ticker := time.NewTicker(per)
		defer ticker.Stop()
		
		for range ticker.C {
			limiter.cleanup()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter.mu.Lock()
		defer limiter.mu.Unlock()

		v, exists := limiter.visitors[ip]
		now := time.Now()

		if !exists {
			limiter.visitors[ip] = &visitor{
				count: 1,
				lastSeen: now,
			}
			c.Next()
			return
		}

		if now.Sub(v.lastSeen) > limiter.per {
			v.count = 0
			v.lastSeen = now
		}

		if v.count >= limiter.rate {
			appError := response.NewAppError(
				429,
				"Rate limit exceeded. Please try again later",
			)
			c.AbortWithStatusJSON(appError.Code, appError)
			return 
		}

		v.count++
		v.lastSeen = now
		c.Next()
	}
}