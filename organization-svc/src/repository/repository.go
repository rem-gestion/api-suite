package repository

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Repositories is a struct that holds all repository instances
type Repositories struct {
	Organization        OrganizationRepository
	OrganizationSetting OrganizationSettingRepository
	Employee            EmployeeRepository
	EmployeeRole        EmployeeRoleRepository
	Branch              BranchRepository
	Role                RoleRepository
	Owner               OwnerRepository
	Invitation          InvitationRepository
	InvitationLog       InvitationLogRepository
}

// NewRepositories creates all repository instances with the given database connection
func NewRepositories(db *gorm.DB, logger *zap.Logger) *Repositories {
	return &Repositories{
		Organization:        NewOrganizationRepository(db, logger),
		OrganizationSetting: NewOrganizationSettingRepository(db, logger),
		Employee:            NewEmployeeRepository(db, logger),
		EmployeeRole:        NewEmployeeRoleRepository(db, logger),
		Branch:              NewBranchRepository(db, logger),
		Role:                NewRoleRepository(db, logger),
		Owner:               NewOwnerRepository(db, logger),
		Invitation:          NewInvitationRepository(db, logger),
		InvitationLog:       NewInvitationLogRepository(db, logger),
	}
}

// Repository interface aggregator for dependency injection
type Repository interface {
	// Organization operations
	GetOrganizationRepository() OrganizationRepository
	GetOrganizationSettingRepository() OrganizationSettingRepository

	// Employee operations
	GetEmployeeRepository() EmployeeRepository
	GetEmployeeRoleRepository() EmployeeRoleRepository

	// Branch operations
	GetBranchRepository() BranchRepository

	// Role operations
	GetRoleRepository() RoleRepository

	// Owner operations
	GetOwnerRepository() OwnerRepository

	// Invitation operations
	GetInvitationRepository() InvitationRepository
	GetInvitationLogRepository() InvitationLogRepository
}

// repositoryManager implements Repository interface
type repositoryManager struct {
	repos *Repositories
}

// NewRepositoryManager creates a new repository manager
func NewRepositoryManager(db *gorm.DB, logger *zap.Logger) Repository {
	return &repositoryManager{
		repos: NewRepositories(db, logger),
	}
}

// Organization repository getters
func (rm *repositoryManager) GetOrganizationRepository() OrganizationRepository {
	return rm.repos.Organization
}

func (rm *repositoryManager) GetOrganizationSettingRepository() OrganizationSettingRepository {
	return rm.repos.OrganizationSetting
}

// Employee repository getters
func (rm *repositoryManager) GetEmployeeRepository() EmployeeRepository {
	return rm.repos.Employee
}

func (rm *repositoryManager) GetEmployeeRoleRepository() EmployeeRoleRepository {
	return rm.repos.EmployeeRole
}

// Branch repository getters
func (rm *repositoryManager) GetBranchRepository() BranchRepository {
	return rm.repos.Branch
}

// Role repository getters
func (rm *repositoryManager) GetRoleRepository() RoleRepository {
	return rm.repos.Role
}

// Owner repository getters
func (rm *repositoryManager) GetOwnerRepository() OwnerRepository {
	return rm.repos.Owner
}

// Invitation repository getters
func (rm *repositoryManager) GetInvitationRepository() InvitationRepository {
	return rm.repos.Invitation
}

func (rm *repositoryManager) GetInvitationLogRepository() InvitationLogRepository {
	return rm.repos.InvitationLog
}
