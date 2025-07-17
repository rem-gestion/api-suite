package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rem-gestion/rem-common/errors"
	"github.com/rem-gestion/rem-common/utils"
)

// JWTAuthLocal creates a middleware for local JWT validation using shared secret
// This approach is more efficient than calling auth-identity-svc via REST/gRPC
// and follows microservices best practices for stateless authentication
func JWTAuthLocal(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract Bearer token from Authorization header
		token := extractBearerToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Missing or invalid authorization header",
				"code":  "MISSING_TOKEN",
			})
			c.Abort()
			return
		}

		// Parse and validate JWT using shared secret (LOCAL VALIDATION)
		claims, err := utils.ParseToken(token, jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid or expired token",
				"code":    "INVALID_TOKEN",
				"details": err.Error(),
			})
			c.Abort()
			return
		}

		// Set user context for downstream middlewares/controllers
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Set("authenticated", true)

		// Continue to next middleware/controller
		c.Next()
	}
}

// extractBearerToken extracts JWT token from Authorization header
func extractBearerToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	// Check for Bearer prefix
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return ""
	}

	// Extract token part
	token := strings.TrimPrefix(authHeader, bearerPrefix)
	token = strings.TrimSpace(token)

	if token == "" {
		return ""
	}

	return token
}

// RequireAuth ensures the request is authenticated (use after JWTAuthBridge)
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authenticated, exists := c.Get("authenticated")
		if !exists || !authenticated.(bool) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
				"code":  "AUTHENTICATION_REQUIRED",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// GetUserID extracts user ID from context (helper for controllers)
func GetUserID(c *gin.Context) (string, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", &errors.BadRequestError{Msg: "User ID not found in context"}
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return "", &errors.InternalServerError{Msg: "Invalid user ID in context"}
	}

	return userIDStr, nil
}

// GetUserEmail extracts user email from context (helper for controllers)
// Note: Email is not stored in JWT claims by default, use gRPC call to auth-identity-svc if needed
func GetUserEmail(c *gin.Context) (string, error) {
	// Email is not included in JWT claims for security/performance reasons
	// If email is needed, make a gRPC call to auth-identity-svc with user_id
	return "", &errors.BadRequestError{Msg: "Email not available in JWT claims, use gRPC call to auth-identity-svc"}
}

// GetUserRole extracts user role from context (helper for controllers)
func GetUserRole(c *gin.Context) (string, error) {
	role, exists := c.Get("user_role")
	if !exists {
		return "", &errors.BadRequestError{Msg: "User role not found in context"}
	}

	roleStr, ok := role.(string)
	if !ok {
		return "", &errors.InternalServerError{Msg: "Invalid user role in context"}
	}

	return roleStr, nil
}
