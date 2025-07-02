package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RecoveryWithZap atrapa panics, los loggea con Zap y devuelve 500
func RecoveryWithZap(lg *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				lg.Error("panic recovered",
					zap.Any("error", rec),
					zap.String("path", c.Request.URL.Path),
				)
				c.AbortWithStatusJSON(500, gin.H{"error": "internal server error"})
			}
		}()
		c.Next()
	}
}
