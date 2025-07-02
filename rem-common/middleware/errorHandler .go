package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rem-gestion/rem-common/errors"
)

// ErrorHandler atrapa c.Errors y responde con JSON y código apropiado.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		switch e := err.(type) {
		case *errors.ValidationError:
			// 422 Unprocessable Entity
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": e.Error(), "fields": e.Fields})
		case *errors.BadRequestError:
			// 400 Bad Request
			c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
		case *errors.UnauthorizedError:
			// 401 Unauthorized
			c.JSON(http.StatusUnauthorized, gin.H{"error": e.Error()})
		case *errors.ForbiddenError:
			// 403 Forbidden
			c.JSON(http.StatusForbidden, gin.H{"error": e.Error()})
		case *errors.NotFoundError:
			// 404 Not Found
			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
		case *errors.ConflictError:
			// 409 Conflict
			c.JSON(http.StatusConflict, gin.H{"error": e.Error()})
		case *errors.TooManyRequestsError:
			// 429 Too Many Requests
			c.JSON(http.StatusTooManyRequests, gin.H{"error": e.Error()})
		case *errors.InternalServerError:
			// 500 Internal Server Error
			c.JSON(http.StatusInternalServerError, gin.H{"error": e.Error()})
		default:
			// Cualquier otro error no tipificado → 500
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}
}
