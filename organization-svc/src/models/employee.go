package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EmployeeStatus represents the status enum for employees
type EmployeeStatus string

const (
	EmpStatusActive     EmployeeStatus = "active"
	EmpStatusInactive   EmployeeStatus = "inactive"
	EmpStatusSuspended  EmployeeStatus = "suspended"
	EmpStatusTerminated EmployeeStatus = "terminated"
)

// Employee represents the link between user and organization
type Employee struct {
	ID             uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrganizationID uuid.UUID      `json:"organization_id" gorm:"type:uuid;not null;index"`
	UserID         uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"` // FK to auth-identity-svc
	PersonID       *uuid.UUID     `json:"person_id" gorm:"type:uuid"`              // FK to person-svc
	PrimaryRoleID  *uuid.UUID     `json:"primary_role_id" gorm:"type:uuid"`        // FK to organization_role
	Status         EmployeeStatus `json:"status" gorm:"type:employee_status_enum;not null;default:'active'" validate:"required,oneof=active inactive suspended terminated"`
	HiredAt        *time.Time     `json:"hired_at" gorm:"type:date"`
	FiredAt        *time.Time     `json:"fired_at" gorm:"type:date"`

	// Audit fields
	CreatedAt time.Time       `json:"created_at" gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`
	CreatedBy uuid.UUID       `json:"created_by" gorm:"type:uuid;not null"` // FK to auth-identity-svc
	UpdatedAt time.Time       `json:"updated_at" gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`
	UpdatedBy *uuid.UUID      `json:"updated_by" gorm:"type:uuid"` // FK to auth-identity-svc
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	Organization  Organization      `json:"organization,omitempty" gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
	PrimaryRole   *OrganizationRole `json:"primary_role,omitempty" gorm:"foreignKey:PrimaryRoleID;constraint:OnDelete:SET NULL"`
	EmployeeRoles []EmployeeRole    `json:"employee_roles,omitempty" gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE"`
}

// EmployeeRole represents the N-to-N relationship between employees and roles
type EmployeeRole struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	EmployeeID uuid.UUID `json:"employee_id" gorm:"type:uuid;not null;index"`
	RoleID     uuid.UUID `json:"role_id" gorm:"type:uuid;not null;index"`
	IsPrimary  bool      `json:"is_primary" gorm:"not null;default:false"`

	// Audit fields
	CreatedAt time.Time       `json:"created_at" gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`
	CreatedBy uuid.UUID       `json:"created_by" gorm:"type:uuid;not null"` // FK to auth-identity-svc
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	Employee Employee         `json:"employee,omitempty" gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE"`
	Role     OrganizationRole `json:"role,omitempty" gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE"`
}

// TableName specifies the table name for Employee
func (Employee) TableName() string {
	return "employees"
}

// TableName specifies the table name for EmployeeRole
func (EmployeeRole) TableName() string {
	return "employee_roles"
}
