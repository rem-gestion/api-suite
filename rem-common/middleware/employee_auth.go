package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type EmployeeChecker interface {
	IsMember(userID, orgID, role string) bool
}

func EmployeeAuth(ec EmployeeChecker, requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.GetString("userID")
		org := c.Param("orgID")
		if !ec.IsMember(uid, org, requiredRole) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "acceso denegado"})
			return
		}
		c.Next()
	}
}
