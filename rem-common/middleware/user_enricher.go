package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/rem-gestion/rem-common/errors"
	authpb "github.com/rem-gestion/rem-common/protos/auth-identity/v1"
)

// UserInfoEnricher enriches user context with additional information via gRPC
// Use this middleware AFTER JWTAuthLocal when you need user details beyond JWT claims
func UserInfoEnricher(authGRPCAddress string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip if user is not authenticated
		authenticated, exists := c.Get("authenticated")
		if !exists || !authenticated.(bool) {
			c.Next()
			return
		}

		// Get user ID from JWT context
		userID, err := GetUserID(c)
		if err != nil {
			c.Next() // Continue without enrichment if no user ID
			return
		}

		// Connect to auth-identity-svc via gRPC
		conn, err := grpc.NewClient(authGRPCAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			// Log error but don't fail the request
			c.Next()
			return
		}
		defer conn.Close()

		// Create gRPC client
		client := authpb.NewAuthIdentityServiceClient(conn)

		// Create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Call GetUserById via gRPC
		resp, err := client.GetUserById(ctx, &authpb.GetUserByIdRequest{
			UserId: userID,
		})
		if err != nil {
			// Log error but don't fail the request
			c.Next()
			return
		}

		// Enrich context with additional user information
		if resp != nil {
			c.Set("user_id", resp.Id)
			c.Set("user_account_id", resp.AccountId)
			c.Set("user_person_id", resp.PersonId)
			c.Set("user_onboard_status", resp.OnboardStatus.String())
			c.Set("user_last_login", resp.LastLogin)
			c.Set("user_created_at", resp.CreatedAt)
			c.Set("user_enriched", true)

			// Also add account info if available
			if resp.Account != nil {
				c.Set("account_email", resp.Account.Email)
				c.Set("account_status", resp.Account.Status.String())
			}
		}

		c.Next()
	}
}

// GetUserEmailViaGRPC gets user email via gRPC call (use sparingly, prefer JWT claims)
func GetUserEmailViaGRPC(c *gin.Context) (string, error) {
	// Check if already enriched
	if email, exists := c.Get("account_email"); exists {
		if emailStr, ok := email.(string); ok {
			return emailStr, nil
		}
	}

	return "", &errors.BadRequestError{Msg: "User email not available, use UserInfoEnricher middleware"}
}

// GetUserAccountID gets user account ID from enriched context
func GetUserAccountID(c *gin.Context) (string, error) {
	accountID, exists := c.Get("user_account_id")
	if !exists {
		return "", &errors.BadRequestError{Msg: "User account ID not available, use UserInfoEnricher middleware"}
	}

	accountIDStr, ok := accountID.(string)
	if !ok {
		return "", &errors.InternalServerError{Msg: "Invalid user account ID in context"}
	}

	return accountIDStr, nil
}

// GetUserPersonID gets user person ID from enriched context
func GetUserPersonID(c *gin.Context) (string, error) {
	personID, exists := c.Get("user_person_id")
	if !exists {
		return "", &errors.BadRequestError{Msg: "User person ID not available, use UserInfoEnricher middleware"}
	}

	personIDStr, ok := personID.(string)
	if !ok {
		return "", &errors.InternalServerError{Msg: "Invalid user person ID in context"}
	}

	return personIDStr, nil
}

// IsUserEnriched checks if user context has been enriched with gRPC data
func IsUserEnriched(c *gin.Context) bool {
	enriched, exists := c.Get("user_enriched")
	return exists && enriched.(bool)
}
