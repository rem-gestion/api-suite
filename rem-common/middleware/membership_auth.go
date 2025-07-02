package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MembershipChecker es tu interfaz para chequear suscripciones
type MembershipChecker interface {
	HasActive(userID string) bool
}

func MembershipAuth(mc MembershipChecker) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.GetString("userID")
		if !mc.HasActive(uid) {
			c.AbortWithStatusJSON(http.StatusPaymentRequired, gin.H{"error": "plan inactivo"})
			return
		}
		c.Next()
	}
}
