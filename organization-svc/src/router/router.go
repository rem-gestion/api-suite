package router

import (
	"github.com/gin-gonic/gin"
	"github.com/rem-gestion/api-suite/organization/src/controllers"
	"github.com/rem-gestion/rem-common/middleware"
)

// SetupBasicOrganizationRoutes configura rutas básicas solo para organizaciones
func SetupBasicOrganizationRoutes(
	router *gin.Engine,
	orgController *controllers.OrganizationController,
	jwtSecret string,
) {
	// Rutas protegidas - todas las operaciones requieren autenticación
	protected := router.Group("/")
	protected.Use(middleware.UserAuth(jwtSecret))
	{
		// === RUTAS DE ORGANIZACIONES ===
		// Rutas principales de organizaciones
		protected.POST("/organizations", orgController.CreateOrganization)
		protected.GET("/organizations", orgController.ListOrganizations)
		protected.GET("/organizations/:org_id", orgController.GetOrganization)
		protected.PUT("/organizations/:org_id", orgController.UpdateOrganization)
		protected.PUT("/organizations/:org_id/status", orgController.UpdateOrganizationStatus)

		// Configuración de organizaciones
		protected.GET("/organizations/:org_id/settings", orgController.GetOrganizationSettings)
		protected.PUT("/organizations/:org_id/settings", orgController.UpdateOrganizationSettings)
	}
}

// SetupRoutes configura todas las rutas del microservicio
func SetupRoutes(
	router *gin.Engine,
	controllers *controllers.Controllers,
	jwtSecret string,
) {

	// Middleware globales ya aplicados en main.go
	// router.Use(gin.Logger())
	// router.Use(gin.Recovery())
	// router.Use(middleware.RequestID())
	// router.Use(middleware.ErrorHandler())

	// Health check endpoint ya configurado en main.go

	// Rutas protegidas - todas las operaciones requieren autenticación
	protected := router.Group("/")
	protected.Use(middleware.UserAuth(jwtSecret))
	{
		// === RUTAS DE ORGANIZACIONES ===
		// Rutas principales de organizaciones
		protected.POST("/organizations", controllers.Organization.CreateOrganization)
		protected.GET("/organizations", controllers.Organization.ListOrganizations)
		protected.GET("/organizations/:org_id", controllers.Organization.GetOrganization)
		protected.PUT("/organizations/:org_id", controllers.Organization.UpdateOrganization)
		protected.PUT("/organizations/:org_id/status", controllers.Organization.UpdateOrganizationStatus)

		// Configuración de organizaciones
		protected.GET("/organizations/:org_id/settings", controllers.Organization.GetOrganizationSettings)
		protected.PUT("/organizations/:org_id/settings", controllers.Organization.UpdateOrganizationSettings)

		// === RUTAS DE EMPLEADOS ===
		// Gestión de empleados
		protected.POST("/organizations/:org_id/employees", controllers.Employee.AddEmployee)
		protected.GET("/organizations/:org_id/employees", controllers.Employee.ListEmployees)
		protected.GET("/organizations/:org_id/employees/:id", controllers.Employee.GetEmployee)
		protected.PUT("/organizations/:org_id/employees/:id", controllers.Employee.UpdateEmployee)
		protected.PUT("/organizations/:org_id/employees/:id/status", controllers.Employee.UpdateEmployeeStatus)

		// Gestión de roles de empleados
		protected.GET("/organizations/:org_id/employees/:id/roles", controllers.Employee.GetEmployeeRoles)
		protected.POST("/organizations/:org_id/employees/:id/roles", controllers.Employee.AssignRole)
		protected.DELETE("/organizations/:org_id/employees/:id/roles/:role_id", controllers.Employee.RemoveRole)

		// === RUTAS DE ROLES ===
		// Gestión de roles organizacionales
		protected.POST("/organizations/:org_id/roles", controllers.Role.CreateRole)
		protected.GET("/organizations/:org_id/roles", controllers.Role.ListRoles)
		protected.GET("/organizations/:org_id/roles/:role_id", controllers.Role.GetRole)
		protected.PUT("/organizations/:org_id/roles/:role_id", controllers.Role.UpdateRole)
		protected.DELETE("/organizations/:org_id/roles/:role_id", controllers.Role.DeleteRole)

		// === RUTAS DE SUCURSALES ===
		// Gestión de sucursales
		protected.POST("/organizations/:org_id/branches", controllers.Branch.CreateBranch)
		protected.GET("/organizations/:org_id/branches", controllers.Branch.ListBranches)
		protected.GET("/organizations/:org_id/branches/:branch_id", controllers.Branch.GetBranch)
		protected.PUT("/organizations/:org_id/branches/:branch_id", controllers.Branch.UpdateBranch)
		protected.DELETE("/organizations/:org_id/branches/:branch_id", controllers.Branch.DeleteBranch)

		// === RUTAS DE PROPIETARIOS ===
		// Gestión de propietarios
		protected.GET("/organizations/:org_id/owners", controllers.Owner.ListOwners)
		protected.POST("/organizations/:org_id/owners", controllers.Owner.AddOwner)
		protected.GET("/organizations/:org_id/owners/:owner_id", controllers.Owner.GetOwner)
		protected.DELETE("/organizations/:org_id/owners/:owner_id", controllers.Owner.RemoveOwner)

		// === RUTAS DE INVITACIONES ===
		// Gestión de invitaciones
		protected.POST("/organizations/:org_id/invitations", controllers.Invite.SendInvitation)
		protected.GET("/organizations/:org_id/invitations", controllers.Invite.ListInvitations)
		protected.GET("/organizations/:org_id/invitations/:id", controllers.Invite.GetInvitation)
		protected.POST("/organizations/:org_id/invitations/:id/resend", controllers.Invite.ResendInvitation)
		protected.PUT("/organizations/:org_id/invitations/:id/cancel", controllers.Invite.CancelInvitation)

		// Validación y procesamiento de invitaciones (públicas)
		protected.GET("/invitations/:token/validate", controllers.Invite.ValidateInvitationToken)
		protected.POST("/invitations/:token/accept", controllers.Invite.AcceptInvitation)
		protected.POST("/invitations/:token/decline", controllers.Invite.DeclineInvitation)
	}
}
