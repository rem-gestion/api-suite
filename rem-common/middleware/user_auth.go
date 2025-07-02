// middleware/user_auth.go
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	jwtutils "github.com/rem-gestion/rem-common/utils" // tu helper de JWT
)

// UserAuth exige un JWT válido y pone "userID" en c.Keys
func UserAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		parts := strings.SplitN(h, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token inválido"})
			return
		}
		claims, err := jwtutils.ParseToken(parts[1], secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
