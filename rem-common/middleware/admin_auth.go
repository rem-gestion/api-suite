package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminAuth asume que UserAuth ya puso claims en c.Keys
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role") // pusiste role en tu helper de jwt
		if role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "solo admins"})
			return
		}
		c.Next()
	}
}
