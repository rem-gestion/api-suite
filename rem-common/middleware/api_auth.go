package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIKeyAuth devuelve un middleware que obliga a incluir el header
// headerName con el valor apiKey. Usa comparación constante para evitar timing attacks.
func APIKeyAuth(headerName, apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		incoming := c.GetHeader(headerName)
		// comparación en tiempo constante
		if subtle.ConstantTimeCompare([]byte(incoming), []byte(apiKey)) != 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "api key inválida",
			})
			return
		}
		c.Next()
	}
}
