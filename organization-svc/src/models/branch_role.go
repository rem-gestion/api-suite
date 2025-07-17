package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OrganizationBranch represents branches/offices of an organization
type OrganizationBranch struct {
	ID             uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrganizationID uuid.UUID  `json:"organization_id" gorm:"type:uuid;not null;index"`
	DisplayName    string     `json:"display_name" gorm:"type:varchar(120);not null" validate:"required,min=2,max=120"`
	AddressID      *uuid.UUID `json:"address_id" gorm:"type:uuid"` // FK to address-svc
	Phone          *string    `json:"phone" gorm:"type:varchar(32)" validate:"omitempty,phone"`
	Email          *string    `json:"email" gorm:"type:varchar(160)" validate:"omitempty,email"`
	IsMain         bool       `json:"is_main" gorm:"not null;default:false"`

	// Audit fields
	CreatedAt time.Time       `json:"created_at" gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`
	CreatedBy uuid.UUID       `json:"created_by" gorm:"type:uuid;not null"` // FK to auth-identity-svc
	UpdatedAt time.Time       `json:"updated_at" gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`
	UpdatedBy *uuid.UUID      `json:"updated_by" gorm:"type:uuid"` // FK to auth-identity-svc
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	Organization Organization `json:"organization,omitempty" gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
}

// OrganizationRole represents internal roles per organization
type OrganizationRole struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrganizationID uuid.UUID `json:"organization_id" gorm:"type:uuid;not null;index"`
	Name           string    `json:"name" gorm:"type:varchar(64);not null" validate:"required,min=2,max=64"`
	Description    *string   `json:"description" gorm:"type:text"`
	IsDefault      bool      `json:"is_default" gorm:"not null;default:false"`

	// Audit fields
	CreatedAt time.Time       `json:"created_at" gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`
	CreatedBy uuid.UUID       `json:"created_by" gorm:"type:uuid;not null"` // FK to auth-identity-svc
	UpdatedAt time.Time       `json:"updated_at" gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`
	UpdatedBy *uuid.UUID      `json:"updated_by" gorm:"type:uuid"` // FK to auth-identity-svc
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	Organization  Organization   `json:"organization,omitempty" gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
	EmployeeRoles []EmployeeRole `json:"employee_roles,omitempty" gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE"`
}

// TableName specifies the table name for OrganizationBranch
func (OrganizationBranch) TableName() string {
	return "organization_branch"
}

// TableName specifies the table name for OrganizationRole
func (OrganizationRole) TableName() string {
	return "organization_role"
}
