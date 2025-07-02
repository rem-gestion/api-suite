package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

// RateLimit crea un middleware que permite rps tokens por segundo y burst tamaño
func RateLimit(rps float64, burst int64) gin.HandlerFunc {
	bucket := ratelimit.NewBucketWithRate(rps, burst)
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit excedido"})
			return
		}
		c.Next()
	}
}
