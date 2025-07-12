package router

import (
	"github.com/gin-gonic/gin"
	"github.com/rem-gestion/api-suite/auth-identity/src/controllers"
	"github.com/rem-gestion/rem-common/middleware"
)

// SetupRoutes configura todas las rutas del microservicio
func SetupRoutes(
	router *gin.Engine,
	authController *controllers.AuthController,
	userController *controllers.UserController,
	jwtSecret string,
) {

	// Middleware globales ya aplicados en main.go
	// router.Use(gin.Logger())
	// router.Use(gin.Recovery())
	// router.Use(middleware.RequestID())
	// router.Use(middleware.ErrorHandler())

	// Health check endpoint ya configurado en main.go
	// router.GET("/health", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"status":  "healthy",
	// 		"service": "auth-identity-svc",
	// 	})
	// })

	// Rutas de autenticación (públicas)
	// API Gateway enrutará /api/auth/ aquí, por lo que usamos "/" como base
	auth := router.Group("/")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
		auth.POST("/refresh", authController.RefreshToken)
		auth.POST("/forgot-password", authController.ForgotPassword)
		auth.POST("/reset-password", authController.ResetPassword)
		auth.GET("/validate", authController.ValidateToken) // Para validar tokens externamente
	}

	// Rutas protegidas de autenticación
	authProtected := router.Group("/")
	authProtected.Use(middleware.UserAuth(jwtSecret))
	{
		authProtected.GET("/me", authController.GetProfile)
		authProtected.POST("/logout", authController.Logout)
	}

	// Rutas de usuarios protegidas
	// API Gateway enrutará /api/users/ a estas rutas
	users := router.Group("/users")
	users.Use(middleware.UserAuth(jwtSecret))
	{
		// Rutas para el usuario actual
		users.PUT("/profile", userController.UpdateProfile)

		// Rutas administrativas (requieren permisos de admin)
		// TODO: Agregar middleware de admin cuando esté disponible
		users.GET("/", userController.ListUsers)                     // GET /users
		users.GET("/email/:email", userController.GetUserByEmail)    // GET /users/email/{email}
		users.POST("/:id/deactivate", userController.DeactivateUser) // POST /users/{id}/deactivate
	}
}

// SetupPublicRoutes configura solo las rutas públicas (útil para testing)
func SetupPublicRoutes(authController *controllers.AuthController) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.RequestID())
	router.Use(middleware.ErrorHandler())

	// Para testing, mantenemos las rutas originales con /api/v1
	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
			auth.POST("/refresh", authController.RefreshToken)
			auth.POST("/forgot-password", authController.ForgotPassword)
			auth.POST("/reset-password", authController.ResetPassword)
			auth.GET("/validate", authController.ValidateToken)
		}
	}

	return router
}
