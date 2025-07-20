// middleware/user_auth.go
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	jwtutils "github.com/rem-gestion/rem-common/utils" // tu helper de JWT
)

// UserClaims represents JWT user claims structure
type UserClaims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	PersonID  string `json:"person_id"`
	AccountID string `json:"account_id"`
}

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
		c.Set("userClaims", claims)
		c.Next()
	}
}

// GetUserClaims extracts user claims from gin context
func GetUserClaims(c *gin.Context) (*UserClaims, bool) {
	if userID, exists := c.Get("userID"); exists {
		claims := &UserClaims{
			UserID: userID.(string),
		}

		// Try to get enriched user data
		if accountID, exists := c.Get("user_account_id"); exists {
			claims.AccountID = accountID.(string)
		}
		if personID, exists := c.Get("user_person_id"); exists {
			claims.PersonID = personID.(string)
		}
		if email, exists := c.Get("account_email"); exists {
			claims.Email = email.(string)
		}

		return claims, true
	}
	return nil, false
}
