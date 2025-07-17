package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BodySizeLimit limita el tamaño del body de las peticiones HTTP
// Previene ataques DoS por carga excesiva de datos
func BodySizeLimit(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Solo aplicar a métodos que pueden tener body
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			if c.Request.ContentLength > maxSize {
				c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
					"error":          "Request body too large",
					"max_size_bytes": maxSize,
					"received_bytes": c.Request.ContentLength,
				})
				return
			}

			// Limitar el reader para prevenir bypass del Content-Length
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		}

		c.Next()
	}
}
