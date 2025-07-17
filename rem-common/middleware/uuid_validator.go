package middleware

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UUIDParamValidator valida que los parámetros de ruta especificados sean UUIDs válidos
// Útil para prevenir errores 500 por parámetros malformados
func UUIDParamValidator(paramNames ...string) gin.HandlerFunc {
	// Pre-compilar regex para mejor rendimiento
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

	return func(c *gin.Context) {
		for _, paramName := range paramNames {
			paramValue := c.Param(paramName)
			if paramValue == "" {
				continue // Parámetro opcional no presente
			}

			// Validar formato UUID usando regex (más rápido que uuid.Parse)
			if !uuidRegex.MatchString(paramValue) {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error":           "Invalid UUID format",
					"parameter":       paramName,
					"value":           paramValue,
					"expected_format": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
				})
				return
			}

			// Validación adicional con uuid.Parse para casos edge
			if _, err := uuid.Parse(paramValue); err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error":     "Invalid UUID",
					"parameter": paramName,
					"value":     paramValue,
					"details":   err.Error(),
				})
				return
			}
		}

		c.Next()
	}
}
