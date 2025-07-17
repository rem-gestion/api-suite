package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PaginationDefaults define los valores por defecto para paginación
type PaginationDefaults struct {
	DefaultPage    int `json:"default_page"`
	DefaultPerPage int `json:"default_per_page"`
	MaxPerPage     int `json:"max_per_page"`
}

// DefaultPaginationConfig configuración estándar de paginación
var DefaultPaginationConfig = PaginationDefaults{
	DefaultPage:    1,
	DefaultPerPage: 20,
	MaxPerPage:     100,
}

// PaginationValidator valida y normaliza parámetros de paginación
// Establece defaults sensatos y previene valores excesivos
func PaginationValidator(config ...PaginationDefaults) gin.HandlerFunc {
	// Usar configuración por defecto si no se especifica
	cfg := DefaultPaginationConfig
	if len(config) > 0 {
		cfg = config[0]
	}

	return func(c *gin.Context) {
		// Solo aplicar a endpoints GET que probablemente usen paginación
		if c.Request.Method != "GET" {
			c.Next()
			return
		}

		// Validar y normalizar 'page'
		page := cfg.DefaultPage
		if pageStr := c.Query("page"); pageStr != "" {
			if parsedPage, err := strconv.Atoi(pageStr); err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error":     "Invalid page parameter",
					"parameter": "page",
					"value":     pageStr,
					"expected":  "positive integer",
				})
				return
			} else if parsedPage < 1 {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error":     "Page must be greater than 0",
					"parameter": "page",
					"value":     parsedPage,
					"minimum":   1,
				})
				return
			} else {
				page = parsedPage
			}
		}

		// Validar y normalizar 'per_page'
		perPage := cfg.DefaultPerPage
		if perPageStr := c.Query("per_page"); perPageStr != "" {
			if parsedPerPage, err := strconv.Atoi(perPageStr); err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error":     "Invalid per_page parameter",
					"parameter": "per_page",
					"value":     perPageStr,
					"expected":  "positive integer",
				})
				return
			} else if parsedPerPage < 1 {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error":     "per_page must be greater than 0",
					"parameter": "per_page",
					"value":     parsedPerPage,
					"minimum":   1,
				})
				return
			} else if parsedPerPage > cfg.MaxPerPage {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error":     "per_page exceeds maximum allowed",
					"parameter": "per_page",
					"value":     parsedPerPage,
					"maximum":   cfg.MaxPerPage,
				})
				return
			} else {
				perPage = parsedPerPage
			}
		}

		// Calcular offset y establecer en contexto para uso en handlers
		offset := (page - 1) * perPage

		c.Set("pagination_page", page)
		c.Set("pagination_per_page", perPage)
		c.Set("pagination_offset", offset)

		c.Next()
	}
}
