package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	ginlimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

var rateLimitStore = memory.NewStore()

var (
	AuthRateLimit   = limiter.Rate{Period: 1 * time.Minute, Limit: 10}
	CreateRateLimit = limiter.Rate{Period: 1 * time.Minute, Limit: 10}
)

func RateLimit(rate limiter.Rate) gin.HandlerFunc {
	return ginlimiter.NewMiddleware(
		limiter.New(rateLimitStore, rate),
		ginlimiter.WithLimitReachedHandler(func(c *gin.Context) {
			c.JSON(429, gin.H{"error": "Too many requests"})
			c.Abort()
		}),
	)
}
