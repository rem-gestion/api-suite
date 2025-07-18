package services

import (
	"github.com/rem-gestion/api-suite/organization/src/clients"
	"github.com/rem-gestion/api-suite/organization/src/repository"
	"go.uber.org/zap"
)

// Services contiene todos los servicios de dominio del sistema de organizaciones
type Services struct {
	Organization *OrganizationService
	Employee     *EmployeeService
	Branch       *BranchService
	Role         *RoleService
	Owner        *OwnerService
	Invite       *InvitationService
}

// NewServices crea una nueva instancia con todos los servicios configurados
func NewServices(
	repos *repository.Repositories,
	clients *clients.ClientManager,
	cache CacheService,
	events EventService,
	logger *zap.Logger,
) *Services {
	// Crear servicios base
	orgService := NewOrganizationService(repos, clients, cache, events, logger)
	employeeService := NewEmployeeService(repos, clients, cache, events, logger)
	branchService := NewBranchService(repos, clients, cache, events, logger)
	roleService := NewRoleService(repos, clients, cache, events, logger)
	ownerService := NewOwnerService(repos, clients, cache, events, logger)

	// El invitation service necesita referencia al employee service
	inviteService := NewInvitationService(repos, clients, cache, events, employeeService, logger)

	return &Services{
		Organization: orgService,
		Employee:     employeeService,
		Branch:       branchService,
		Role:         roleService,
		Owner:        ownerService,
		Invite:       inviteService,
	}
}
