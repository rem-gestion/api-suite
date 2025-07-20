package controllers

import (
	"github.com/rem-gestion/api-suite/organization/src/services"
	"go.uber.org/zap"
)

// Controllers agrupa todos los controladores del servicio de organización
type Controllers struct {
	Organization *OrganizationController
	Employee     *EmployeeController
	Branch       *BranchController
	Role         *RoleController
	Owner        *OwnerController
	Invite       *InviteController
	Integration  *IntegrationController
	logger       *zap.Logger
}

// NewControllers crea una nueva instancia de todos los controladores usando la estructura Services
func NewControllers(
	services *services.Services,
	logger *zap.Logger,
) *Controllers {
	return &Controllers{
		Organization: NewOrganizationController(services.Organization, logger),
		Employee:     NewEmployeeController(services.Employee, logger),
		Branch:       NewBranchController(services.Branch, logger),
		Role:         NewRoleController(services.Role, logger),
		Owner:        NewOwnerController(services.Owner, logger),
		Invite:       NewInviteController(services.Invite, logger),
		Integration:  NewIntegrationController(logger),
		logger:       logger,
	}
}
