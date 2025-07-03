package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GinLogger emite una línea por request.
// 2xx/3xx → INFO · 4xx → WARN · 5xx → ERROR
func GinLogger(lg *zap.Logger) gin.HandlerFunc {
	nolog := lg.WithOptions(zap.WithCaller(false)) // evita caller en esta capa

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		code := c.Writer.Status()
		elap := float64(time.Since(start).Microseconds()) / 1000 // ms con decimales

		level := zap.InfoLevel
		switch {
		case code >= 500:
			level = zap.ErrorLevel
		case code >= 400:
			level = zap.WarnLevel
		}

		if ce := nolog.Check(level, "http request"); ce != nil {
			ce.Write(
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.Int("status", code),
				zap.Float64("latency_ms", elap),
			)
		}
	}
}
