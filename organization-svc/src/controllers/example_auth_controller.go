package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rem-gestion/rem-common/middleware"
)

// ExampleAuthenticatedController muestra cómo usar el contexto de autenticación
// proporcionado por el JWT Auth Bridge
func ExampleAuthenticatedController(c *gin.Context) {
	// El JWT Auth Bridge ya validó el token y estableció el contexto del usuario
	// Ahora podemos acceder a la información del usuario autenticado

	// Obtener información del usuario desde el contexto
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get user ID from context",
			"code":  "CONTEXT_ERROR",
		})
		return
	}

	userEmail, err := middleware.GetUserEmail(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get user email from context",
			"code":  "CONTEXT_ERROR",
		})
		return
	}

	userRole, err := middleware.GetUserRole(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get user role from context",
			"code":  "CONTEXT_ERROR",
		})
		return
	}

	// También podemos acceder directamente al contexto si necesitamos
	authenticated, _ := c.Get("authenticated")
	tokenExpiresAt, _ := c.Get("token_expires_at")

	// Ejemplo de lógica de negocio que usa la información del usuario
	response := gin.H{
		"message": "Authentication bridge working successfully!",
		"authenticated_user": gin.H{
			"user_id":          userID,
			"email":            userEmail,
			"role":             userRole,
			"authenticated":    authenticated,
			"token_expires_at": tokenExpiresAt,
		},
		"bridge_info": gin.H{
			"description": "This response proves the JWT Auth Bridge is working",
			"flow": []string{
				"1. Client sent JWT token to organization-svc",
				"2. JWT Auth Bridge forwarded token to auth-identity-svc",
				"3. auth-identity-svc validated token and returned user info",
				"4. Bridge set user context in organization-svc",
				"5. Controller accessed authenticated user information",
			},
		},
	}

	c.JSON(http.StatusOK, response)
}

// ExampleOrganizationListController muestra cómo filtrar organizaciones según el usuario autenticado
func ExampleOrganizationListController(c *gin.Context) {
	// Obtener información del usuario autenticado
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not authenticated"})
		return
	}

	// Aquí iría la lógica para:
	// 1. Consultar la BD de organization-svc
	// 2. Filtrar organizaciones donde el usuario es miembro/owner
	// 3. Aplicar permisos según el rol del usuario

	// Por ahora, ejemplo sin BD real
	response := gin.H{
		"message": "This would return organizations for authenticated user",
		"user_id": userID,
		"organizations": []gin.H{
			{
				"id":   "org-1",
				"name": "Example Organization 1",
				"role": "admin",
			},
			{
				"id":   "org-2",
				"name": "Example Organization 2",
				"role": "member",
			},
		},
		"bridge_success": true,
	}

	c.JSON(http.StatusOK, response)
}
