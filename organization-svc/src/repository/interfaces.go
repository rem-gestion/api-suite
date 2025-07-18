package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rem-gestion/api-suite/organization/src/dto"
	"github.com/rem-gestion/api-suite/organization/src/models"
)

// ===============================
// ORGANIZATION REPOSITORY INTERFACES
// ===============================

// OrganizationRepository defines the interface for organization data operations
type OrganizationRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, org *models.Organization) (*models.Organization, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Organization, error)
	GetBySlug(ctx context.Context, slug string) (*models.Organization, error)
	Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.Organization, error)
	Delete(ctx context.Context, id uuid.UUID) error
	SoftDelete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error

	// Listing and filtering
	List(ctx context.Context, filters dto.OrganizationFiltersRequest) ([]models.Organization, int64, error)
	ListByOwner(ctx context.Context, userID uuid.UUID, filters dto.OrganizationFiltersRequest) ([]models.Organization, int64, error)
	ListByEmployee(ctx context.Context, userID uuid.UUID, filters dto.OrganizationFiltersRequest) ([]models.Organization, int64, error)

	// Status management
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.OrganizationStatus, updatedBy uuid.UUID) error
	GetByStatus(ctx context.Context, status models.OrganizationStatus, limit int) ([]models.Organization, error)

	// Verification
	MarkAsVerified(ctx context.Context, id uuid.UUID, verifiedBy uuid.UUID) error
	GetUnverifiedOrganizations(ctx context.Context, limit int) ([]models.Organization, error)

	// Relationships
	GetWithSettings(ctx context.Context, id uuid.UUID) (*models.Organization, error)
	GetWithOwner(ctx context.Context, id uuid.UUID) (*models.Organization, error)
	GetWithEmployees(ctx context.Context, id uuid.UUID) (*models.Organization, error)
	GetWithBranches(ctx context.Context, id uuid.UUID) (*models.Organization, error)
	GetWithRoles(ctx context.Context, id uuid.UUID) (*models.Organization, error)

	// Business logic support
	ExistsBySlug(ctx context.Context, slug string) (bool, error)
	ExistsByName(ctx context.Context, name string, excludeID *uuid.UUID) (bool, error)
	CountByUser(ctx context.Context, userID uuid.UUID) (int64, error)
	GetActiveCount(ctx context.Context) (int64, error)
}

// OrganizationSettingRepository defines the interface for organization settings operations
type OrganizationSettingRepository interface {
	// Settings CRUD
	Create(ctx context.Context, setting *models.OrganizationSetting) (*models.OrganizationSetting, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.OrganizationSetting, error)
	GetByKey(ctx context.Context, orgID uuid.UUID, key string) (*models.OrganizationSetting, error)
	Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.OrganizationSetting, error)
	Delete(ctx context.Context, id uuid.UUID) error

	// Bulk operations
	GetAllByOrganization(ctx context.Context, orgID uuid.UUID) ([]models.OrganizationSetting, error)
	SetMultiple(ctx context.Context, orgID uuid.UUID, settings map[string]interface{}, createdBy uuid.UUID) error
	DeleteByKey(ctx context.Context, orgID uuid.UUID, key string) error

	// Configuration management
	GetEditableSettings(ctx context.Context, orgID uuid.UUID) ([]models.OrganizationSetting, error)
	GetDefaultSettings(ctx context.Context) (map[string]interface{}, error)
	ResetToDefaults(ctx context.Context, orgID uuid.UUID, keys []string, updatedBy uuid.UUID) error

	// Business logic support
	SettingExists(ctx context.Context, orgID uuid.UUID, key string) (bool, error)
	GetSettingsAsMap(ctx context.Context, orgID uuid.UUID) (map[string]interface{}, error)
}

// ===============================
// EMPLOYEE REPOSITORY INTERFACES
// ===============================

// EmployeeRepository defines the interface for employee data operations
type EmployeeRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, employee *models.Employee) (*models.Employee, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Employee, error)
	GetByUserAndOrg(ctx context.Context, userID, orgID uuid.UUID) (*models.Employee, error)
	Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.Employee, error)
	Delete(ctx context.Context, id uuid.UUID) error
	SoftDelete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error

	// Listing and filtering
	List(ctx context.Context, filters dto.EmployeeFiltersRequest) ([]models.Employee, int64, error)
	ListByOrganization(ctx context.Context, orgID uuid.UUID, filters dto.EmployeeFiltersRequest) ([]models.Employee, int64, error)
	ListByRole(ctx context.Context, roleID uuid.UUID, filters dto.EmployeeFiltersRequest) ([]models.Employee, int64, error)
	ListByBranch(ctx context.Context, branchID uuid.UUID, filters dto.EmployeeFiltersRequest) ([]models.Employee, int64, error)

	// Status management
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.EmployeeStatus, updatedBy uuid.UUID) error
	GetByStatus(ctx context.Context, orgID uuid.UUID, status models.EmployeeStatus) ([]models.Employee, error)
	TerminateEmployee(ctx context.Context, id uuid.UUID, terminatedBy uuid.UUID, reason string) error

	// Relationships
	GetWithRoles(ctx context.Context, id uuid.UUID) (*models.Employee, error)
	GetWithOrganization(ctx context.Context, id uuid.UUID) (*models.Employee, error)

	// Business logic support
	ExistsByUserAndOrg(ctx context.Context, userID, orgID uuid.UUID) (bool, error)
	CountByOrganization(ctx context.Context, orgID uuid.UUID) (int64, error)
	CountActiveByOrganization(ctx context.Context, orgID uuid.UUID) (int64, error)
	GetActiveEmployeesByOrg(ctx context.Context, orgID uuid.UUID) ([]models.Employee, error)
	IsUserEmployeeOfOrg(ctx context.Context, userID, orgID uuid.UUID) (bool, error)
}

// EmployeeRoleRepository defines the interface for employee-role relationship operations
type EmployeeRoleRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, empRole *models.EmployeeRole) (*models.EmployeeRole, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.EmployeeRole, error)
	GetByEmployeeAndRole(ctx context.Context, employeeID, roleID uuid.UUID) (*models.EmployeeRole, error)
	Delete(ctx context.Context, id uuid.UUID) error
	SoftDelete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error

	// Role assignments
	AssignRole(ctx context.Context, employeeID, roleID uuid.UUID, isPrimary bool, createdBy uuid.UUID) (*models.EmployeeRole, error)
	RemoveRole(ctx context.Context, employeeID, roleID uuid.UUID) error
	SetPrimaryRole(ctx context.Context, employeeID, roleID uuid.UUID) error
	ClearPrimaryRole(ctx context.Context, employeeID uuid.UUID) error

	// Listing operations
	GetByEmployee(ctx context.Context, employeeID uuid.UUID) ([]models.EmployeeRole, error)
	GetByRole(ctx context.Context, roleID uuid.UUID) ([]models.EmployeeRole, error)
	GetPrimaryRoleByEmployee(ctx context.Context, employeeID uuid.UUID) (*models.EmployeeRole, error)

	// Business logic support
	HasRole(ctx context.Context, employeeID, roleID uuid.UUID) (bool, error)
	CountByRole(ctx context.Context, roleID uuid.UUID) (int64, error)
	GetEmployeeRoleCount(ctx context.Context, employeeID uuid.UUID) (int64, error)
	RemoveAllRoles(ctx context.Context, employeeID uuid.UUID) error
}

// ===============================
// BRANCH REPOSITORY INTERFACES
// ===============================

// BranchRepository defines the interface for organization branch operations
type BranchRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, branch *models.OrganizationBranch) (*models.OrganizationBranch, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.OrganizationBranch, error)
	Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.OrganizationBranch, error)
	Delete(ctx context.Context, id uuid.UUID) error
	SoftDelete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error

	// Listing and filtering
	List(ctx context.Context, filters dto.BranchFiltersRequest) ([]models.OrganizationBranch, int64, error)
	ListByOrganization(ctx context.Context, orgID uuid.UUID, filters dto.BranchFiltersRequest) ([]models.OrganizationBranch, int64, error)

	// Main branch management
	GetMainBranch(ctx context.Context, orgID uuid.UUID) (*models.OrganizationBranch, error)
	SetAsMain(ctx context.Context, id uuid.UUID, updatedBy uuid.UUID) error
	ClearMainStatus(ctx context.Context, orgID uuid.UUID) error

	// Relationships
	GetWithOrganization(ctx context.Context, id uuid.UUID) (*models.OrganizationBranch, error)

	// Business logic support
	CountByOrganization(ctx context.Context, orgID uuid.UUID) (int64, error)
	ExistsByName(ctx context.Context, orgID uuid.UUID, name string, excludeID *uuid.UUID) (bool, error)
	HasMainBranch(ctx context.Context, orgID uuid.UUID) (bool, error)
	CanDelete(ctx context.Context, id uuid.UUID) (bool, error) // Check if branch can be deleted (not last one, no active employees, etc.)
}

// ===============================
// ROLE REPOSITORY INTERFACES
// ===============================

// RoleRepository defines the interface for organization role operations
type RoleRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, role *models.OrganizationRole) (*models.OrganizationRole, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.OrganizationRole, error)
	Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.OrganizationRole, error)
	Delete(ctx context.Context, id uuid.UUID) error
	SoftDelete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error

	// Listing and filtering
	List(ctx context.Context, filters dto.RoleFiltersRequest) ([]models.OrganizationRole, int64, error)
	ListByOrganization(ctx context.Context, orgID uuid.UUID, filters dto.RoleFiltersRequest) ([]models.OrganizationRole, int64, error)

	// Default roles management
	GetDefaultRoles(ctx context.Context, orgID uuid.UUID) ([]models.OrganizationRole, error)
	GetNonDefaultRoles(ctx context.Context, orgID uuid.UUID) ([]models.OrganizationRole, error)
	CreateDefaultRoles(ctx context.Context, orgID uuid.UUID, createdBy uuid.UUID) ([]models.OrganizationRole, error)

	// Relationships
	GetWithEmployees(ctx context.Context, id uuid.UUID) (*models.OrganizationRole, error)
	GetWithOrganization(ctx context.Context, id uuid.UUID) (*models.OrganizationRole, error)

	// Business logic support
	ExistsByName(ctx context.Context, orgID uuid.UUID, name string, excludeID *uuid.UUID) (bool, error)
	CountByOrganization(ctx context.Context, orgID uuid.UUID) (int64, error)
	GetRoleUsageStats(ctx context.Context, roleID uuid.UUID) (map[string]interface{}, error)
	CanDelete(ctx context.Context, id uuid.UUID) (bool, error) // Check if role can be deleted (not default, no employees assigned, etc.)
}

// ===============================
// OWNER REPOSITORY INTERFACES
// ===============================

// OwnerRepository defines the interface for organization owner operations
type OwnerRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, owner *models.OrganizationOwner) (*models.OrganizationOwner, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.OrganizationOwner, error)
	GetByOrganization(ctx context.Context, orgID uuid.UUID) (*models.OrganizationOwner, error)
	GetByPersonAndOrg(ctx context.Context, personID, orgID uuid.UUID) (*models.OrganizationOwner, error)
	Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.OrganizationOwner, error)
	Delete(ctx context.Context, id uuid.UUID) error

	// Listing and filtering
	List(ctx context.Context, filters dto.OwnerFiltersRequest) ([]models.OrganizationOwner, int64, error)
	ListByOrganization(ctx context.Context, orgID uuid.UUID, filters dto.OwnerFiltersRequest) ([]models.OrganizationOwner, int64, error)

	// Ownership management
	GetFounders(ctx context.Context, orgID uuid.UUID) ([]models.OrganizationOwner, error)
	GetMajorityOwner(ctx context.Context, orgID uuid.UUID) (*models.OrganizationOwner, error)
	UpdateOwnershipPercentage(ctx context.Context, id uuid.UUID, percentage float64, updatedBy uuid.UUID) error
	TransferOwnership(ctx context.Context, fromOwnerID, toOwnerID uuid.UUID, percentage float64, updatedBy uuid.UUID) error

	// Validation and statistics
	GetTotalOwnership(ctx context.Context, orgID uuid.UUID) (float64, error)
	ValidateOwnershipPercentages(ctx context.Context, orgID uuid.UUID) (bool, error)
	GetOwnershipSummary(ctx context.Context, orgID uuid.UUID) (*dto.OwnershipSummary, error)

	// Relationships
	GetWithOrganization(ctx context.Context, id uuid.UUID) (*models.OrganizationOwner, error)

	// Business logic support
	ExistsByPersonAndOrg(ctx context.Context, personID, orgID uuid.UUID) (bool, error)
	CountByOrganization(ctx context.Context, orgID uuid.UUID) (int64, error)
	IsUserOwnerOfOrg(ctx context.Context, userID, orgID uuid.UUID) (bool, error)
	CanRemoveOwner(ctx context.Context, id uuid.UUID) (bool, error) // Check business rules for owner removal
}

// ===============================
// INVITATION REPOSITORY INTERFACES
// ===============================

// InvitationRepository defines the interface for organization invitation operations
type InvitationRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, invitation *models.OrganizationInvite) (*models.OrganizationInvite, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.OrganizationInvite, error)
	GetByToken(ctx context.Context, token string) (*models.OrganizationInvite, error)
	Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.OrganizationInvite, error)
	Delete(ctx context.Context, id uuid.UUID) error

	// Listing and filtering
	List(ctx context.Context, filters dto.InvitationFiltersRequest) ([]models.OrganizationInvite, int64, error)
	ListByOrganization(ctx context.Context, orgID uuid.UUID, filters dto.InvitationFiltersRequest) ([]models.OrganizationInvite, int64, error)
	ListByEmail(ctx context.Context, email string) ([]models.OrganizationInvite, error)

	// Status management
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.InvitationStatus, updatedBy *uuid.UUID) error
	AcceptInvitation(ctx context.Context, id uuid.UUID, acceptedBy *uuid.UUID) error
	RejectInvitation(ctx context.Context, id uuid.UUID, rejectedBy *uuid.UUID, reason *string) error
	CancelInvitation(ctx context.Context, id uuid.UUID, cancelledBy uuid.UUID, reason string) error
	ExpireInvitation(ctx context.Context, id uuid.UUID) error

	// Token management
	RegenerateToken(ctx context.Context, id uuid.UUID) (string, error)
	IsTokenValid(ctx context.Context, token string) (bool, error)

	// Expiration management
	GetExpiredInvitations(ctx context.Context, limit int) ([]models.OrganizationInvite, error)
	MarkExpiredInvitations(ctx context.Context) (int64, error)

	// Relationships
	GetWithOrganization(ctx context.Context, id uuid.UUID) (*models.OrganizationInvite, error)
	GetWithRole(ctx context.Context, id uuid.UUID) (*models.OrganizationInvite, error)
	GetWithBranch(ctx context.Context, id uuid.UUID) (*models.OrganizationInvite, error)
	GetWithLogs(ctx context.Context, id uuid.UUID) (*models.OrganizationInvite, error)

	// Business logic support
	ExistsByEmailAndOrg(ctx context.Context, email string, orgID uuid.UUID, excludeStatuses []models.InvitationStatus) (bool, error)
	CountPendingByOrganization(ctx context.Context, orgID uuid.UUID) (int64, error)
	GetInvitationStats(ctx context.Context, orgID uuid.UUID) (map[string]interface{}, error)
}

// InvitationLogRepository defines the interface for invitation audit log operations
type InvitationLogRepository interface {
	// Basic operations
	Create(ctx context.Context, log *models.OrganizationInviteLog) (*models.OrganizationInviteLog, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.OrganizationInviteLog, error)

	// Listing operations
	ListByInvitation(ctx context.Context, inviteID uuid.UUID) ([]models.OrganizationInviteLog, error)
	ListByOrganization(ctx context.Context, orgID uuid.UUID, limit int) ([]models.OrganizationInviteLog, error)
	ListByAction(ctx context.Context, action models.InvitationLogAction, limit int) ([]models.OrganizationInviteLog, error)

	// Convenience methods for common log entries
	LogInvitationCreated(ctx context.Context, inviteID uuid.UUID, actorPersonID *uuid.UUID, clientIP, userAgent *string) error
	LogInvitationAccepted(ctx context.Context, inviteID uuid.UUID, actorPersonID *uuid.UUID, clientIP, userAgent *string) error
	LogInvitationRejected(ctx context.Context, inviteID uuid.UUID, actorPersonID *uuid.UUID, reason *string, clientIP, userAgent *string) error
	LogInvitationCancelled(ctx context.Context, inviteID uuid.UUID, actorPersonID *uuid.UUID, reason *string, clientIP, userAgent *string) error
	LogInvitationExpired(ctx context.Context, inviteID uuid.UUID) error
	LogInvitationResent(ctx context.Context, inviteID uuid.UUID, actorPersonID *uuid.UUID, clientIP, userAgent *string) error

	// Cleanup operations
	DeleteOldLogs(ctx context.Context, olderThan int) (int64, error) // Delete logs older than X days
}
