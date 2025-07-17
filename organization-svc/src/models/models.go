// Package models contains all the data models for organization-svc
// Based on PostgreSQL migrations 0001-0008 with comprehensive entity modeling
package models

// Re-export all model types for convenience
type (
	// Organization Core Models
	OrganizationModel             = Organization
	OrganizationSettingModel      = OrganizationSetting
	OrganizationOwnerModel        = OrganizationOwner
	OrganizationSubscriptionModel = OrganizationSubscription

	// Structure Models
	OrganizationBranchModel = OrganizationBranch
	OrganizationRoleModel   = OrganizationRole

	// Employee Models
	EmployeeModel     = Employee
	EmployeeRoleModel = EmployeeRole

	// Invitation Models
	OrganizationInviteModel    = OrganizationInvite
	OrganizationInviteLogModel = OrganizationInviteLog

	// Integration Models
	IntegrationTypeModel              = IntegrationType
	OrganizationIntegrationModel      = OrganizationIntegration
	OrganizationIntegrationEventModel = OrganizationIntegrationEvent

	// Domain Models
	OrganizationDomainModel                = OrganizationDomain
	OrganizationDomainDNSModel             = OrganizationDomainDNS
	OrganizationDomainVerificationLogModel = OrganizationDomainVerificationLog
)

// AllModels returns a slice of all model types for GORM migrations
func AllModels() []interface{} {
	return []interface{}{
		// Core organization models
		&Organization{},
		&OrganizationSetting{},
		&OrganizationOwner{},
		&OrganizationSubscription{},

		// Structure models
		&OrganizationBranch{},
		&OrganizationRole{},

		// Employee models
		&Employee{},
		&EmployeeRole{},

		// Invitation models
		&OrganizationInvite{},
		&OrganizationInviteLog{},

		// Integration models
		&IntegrationType{},
		&OrganizationIntegration{},
		&OrganizationIntegrationEvent{},

		// Domain models
		&OrganizationDomain{},
		&OrganizationDomainDNS{},
		&OrganizationDomainVerificationLog{},
	}
}

// CoreModels returns the essential models needed for basic organization functionality
func CoreModels() []interface{} {
	return []interface{}{
		&Organization{},
		&OrganizationSetting{},
		&OrganizationOwner{},
		&OrganizationBranch{},
		&OrganizationRole{},
		&Employee{},
		&EmployeeRole{},
	}
}

// ExtendedModels returns models for advanced features (invitations, integrations, domains)
func ExtendedModels() []interface{} {
	return []interface{}{
		&OrganizationInvite{},
		&OrganizationInviteLog{},
		&IntegrationType{},
		&OrganizationIntegration{},
		&OrganizationIntegrationEvent{},
		&OrganizationDomain{},
		&OrganizationDomainDNS{},
		&OrganizationDomainVerificationLog{},
		&OrganizationSubscription{},
	}
}

// LogModels returns only the audit/log models for partitioned tables
func LogModels() []interface{} {
	return []interface{}{
		&OrganizationInviteLog{},
		&OrganizationIntegrationEvent{},
		&OrganizationDomainVerificationLog{},
	}
}
