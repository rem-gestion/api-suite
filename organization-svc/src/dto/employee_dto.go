package dto

import (
	"time"

	"github.com/google/uuid"
)

// ===============================
// EMPLOYEE REQUEST DTOs
// ===============================

// AddEmployeeRequest represents request to add a new employee to organization
type AddEmployeeRequest struct {
	UserID        uuid.UUID   `json:"user_id" validate:"required,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	PersonID      *uuid.UUID  `json:"person_id,omitempty" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174001"`
	PrimaryRoleID *uuid.UUID  `json:"primary_role_id,omitempty" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174002"`
	RoleIDs       []uuid.UUID `json:"role_ids,omitempty" validate:"omitempty,dive,uuid" example:"[\"123e4567-e89b-12d3-a456-426614174002\"]"`
	HiredAt       *time.Time  `json:"hired_at,omitempty" example:"2024-01-15"`
	Status        string      `json:"status" validate:"omitempty,oneof=active inactive suspended terminated" example:"active"`
	MetadataRequest
}

// UpdateEmployeeRequest represents request to update employee data
type UpdateEmployeeRequest struct {
	PersonID      *uuid.UUID `json:"person_id,omitempty" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174001"`
	PrimaryRoleID *uuid.UUID `json:"primary_role_id,omitempty" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174002"`
	HiredAt       *time.Time `json:"hired_at,omitempty" example:"2024-01-15"`
	FiredAt       *time.Time `json:"fired_at,omitempty" example:"2024-12-15"`
	MetadataRequest
}

// EmployeeFiltersRequest represents filters for listing employees
type EmployeeFiltersRequest struct {
	FilterRequest
	OrganizationID uuid.UUID  `json:"organization_id" form:"organization_id" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	RoleID         *uuid.UUID `json:"role_id" form:"role_id" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174002"`
	BranchID       *uuid.UUID `json:"branch_id" form:"branch_id" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174003"`
	UserID         *uuid.UUID `json:"user_id" form:"user_id" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	HiredAfter     *time.Time `json:"hired_after" form:"hired_after" example:"2024-01-01"`
	HiredBefore    *time.Time `json:"hired_before" form:"hired_before" example:"2024-12-31"`
	PaginationRequest
	IncludeOptions
}

// AssignRoleRequest represents request to assign role to employee
type AssignRoleRequest struct {
	RoleID    uuid.UUID `json:"role_id" validate:"required,uuid" example:"123e4567-e89b-12d3-a456-426614174002"`
	IsPrimary bool      `json:"is_primary" example:"false"`
}

// ===============================
// EMPLOYEE RESPONSE DTOs
// ===============================

// EmployeeResponse represents employee data in responses
type EmployeeResponse struct {
	BaseResponse
	OrganizationID uuid.UUID  `json:"organization_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	UserID         uuid.UUID  `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	PersonID       *uuid.UUID `json:"person_id,omitempty" example:"123e4567-e89b-12d3-a456-426614174001"`
	PrimaryRoleID  *uuid.UUID `json:"primary_role_id,omitempty" example:"123e4567-e89b-12d3-a456-426614174002"`
	Status         string     `json:"status" example:"active"`
	HiredAt        *time.Time `json:"hired_at,omitempty" example:"2024-01-15"`
	FiredAt        *time.Time `json:"fired_at,omitempty" example:"2024-12-15"`

	// Embedded optional data
	Organization  *OrganizationBasicResponse `json:"organization,omitempty"`
	PrimaryRole   *OrganizationRoleResponse  `json:"primary_role,omitempty"`
	EmployeeRoles []EmployeeRoleResponse     `json:"employee_roles,omitempty"`

	// External service data (populated via gRPC)
	UserDetails   *UserDetails   `json:"user_details,omitempty"`
	PersonDetails *PersonDetails `json:"person_details,omitempty"`

	// Statistics (included when include_stats=true)
	Stats *EmployeeStats `json:"stats,omitempty"`
}

// EmployeeListResponse represents paginated employee list
type EmployeeListResponse struct {
	Data       []EmployeeResponse `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}

// EmployeeBasicResponse represents minimal employee data
type EmployeeBasicResponse struct {
	ID             uuid.UUID  `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	OrganizationID uuid.UUID  `json:"organization_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	UserID         uuid.UUID  `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	PersonID       *uuid.UUID `json:"person_id,omitempty" example:"123e4567-e89b-12d3-a456-426614174001"`
	Status         string     `json:"status" example:"active"`
	HiredAt        *time.Time `json:"hired_at,omitempty" example:"2024-01-15"`

	// External service data (basic)
	UserDetails   *UserBasicDetails   `json:"user_details,omitempty"`
	PersonDetails *PersonBasicDetails `json:"person_details,omitempty"`
}

// ===============================
// EMPLOYEE ROLES DTOs
// ===============================

// EmployeeRoleResponse represents employee-role relationship data
type EmployeeRoleResponse struct {
	BaseResponse
	EmployeeID uuid.UUID `json:"employee_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	RoleID     uuid.UUID `json:"role_id" example:"123e4567-e89b-12d3-a456-426614174002"`
	IsPrimary  bool      `json:"is_primary" example:"false"`

	// Embedded optional data
	Employee *EmployeeBasicResponse    `json:"employee,omitempty"`
	Role     *OrganizationRoleResponse `json:"role,omitempty"`
}

// EmployeeRolesResponse represents all roles for an employee
type EmployeeRolesResponse struct {
	EmployeeID  uuid.UUID                 `json:"employee_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Roles       []EmployeeRoleResponse    `json:"roles"`
	PrimaryRole *OrganizationRoleResponse `json:"primary_role,omitempty"`
	UpdatedAt   time.Time                 `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// ===============================
// EMPLOYEE STATISTICS DTOs
// ===============================

// EmployeeStats represents employee statistics
type EmployeeStats struct {
	TotalRoles        int       `json:"total_roles" example:"3"`
	DaysEmployed      int       `json:"days_employed" example:"365"`
	LastActivityAt    time.Time `json:"last_activity_at" example:"2024-01-15T10:30:00Z"`
	PerformanceRating *float64  `json:"performance_rating,omitempty" example:"4.5"`
}

// ===============================
// EXTERNAL SERVICE DTOs (gRPC)
// ===============================

// UserDetails represents user data from auth-identity-svc
type UserDetails struct {
	ID            string    `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	AccountID     string    `json:"account_id" example:"123e4567-e89b-12d3-a456-426614174001"`
	PersonID      string    `json:"person_id" example:"123e4567-e89b-12d3-a456-426614174002"`
	OnboardStatus string    `json:"onboard_status" example:"DONE"`
	LastLogin     time.Time `json:"last_login" example:"2024-01-15T10:30:00Z"`
	CreatedAt     time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt     time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
	UpdatedBy     string    `json:"updated_by" example:"123e4567-e89b-12d3-a456-426614174000"`

	// Account details
	Account *AccountDetails `json:"account,omitempty"`
}

// UserBasicDetails represents minimal user data from auth-identity-svc
type UserBasicDetails struct {
	ID            string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	AccountID     string `json:"account_id" example:"123e4567-e89b-12d3-a456-426614174001"`
	OnboardStatus string `json:"onboard_status" example:"DONE"`

	// Account details
	Account *AccountBasicDetails `json:"account,omitempty"`
}

// AccountDetails represents account data from auth-identity-svc
type AccountDetails struct {
	ID        string    `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Provider  string    `json:"provider" example:"EMAIL"`
	Email     string    `json:"email" example:"user@example.com"`
	Status    string    `json:"status" example:"ACTIVE"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
	UpdatedBy string    `json:"updated_by" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// AccountBasicDetails represents minimal account data from auth-identity-svc
type AccountBasicDetails struct {
	ID       string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Provider string `json:"provider" example:"EMAIL"`
	Email    string `json:"email" example:"user@example.com"`
	Status   string `json:"status" example:"ACTIVE"`
}

// PersonDetails represents person data from person-svc
type PersonDetails struct {
	ID        string    `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Type      string    `json:"type" example:"INDIVIDUAL"`
	AddressID string    `json:"address_id" example:"123e4567-e89b-12d3-a456-426614174001"`
	AvatarURL string    `json:"avatar_url" example:"https://example.com/avatar.jpg"`
	Sexo      string    `json:"sexo" example:"MASCULINO"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`

	// Person type specific details
	Individual *IndividualDetails `json:"individual,omitempty"`
	Company    *CompanyDetails    `json:"company,omitempty"`

	// Related data
	Contacts []ContactDetails `json:"contacts,omitempty"`
	Address  *AddressDetails  `json:"address,omitempty"`
}

// PersonBasicDetails represents minimal person data from person-svc
type PersonBasicDetails struct {
	ID        string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Type      string `json:"type" example:"INDIVIDUAL"`
	AvatarURL string `json:"avatar_url" example:"https://example.com/avatar.jpg"`

	// Person type specific details
	Individual *IndividualBasicDetails `json:"individual,omitempty"`
	Company    *CompanyBasicDetails    `json:"company,omitempty"`
}

// IndividualDetails represents individual person data from person-svc
type IndividualDetails struct {
	FirstName string `json:"first_name" example:"John"`
	LastName  string `json:"last_name" example:"Doe"`
	DNI       string `json:"dni" example:"12345678"`
}

// IndividualBasicDetails represents minimal individual person data
type IndividualBasicDetails struct {
	FirstName string `json:"first_name" example:"John"`
	LastName  string `json:"last_name" example:"Doe"`
}

// CompanyDetails represents company person data from person-svc
type CompanyDetails struct {
	LegalName   string `json:"legal_name" example:"Acme Corporation"`
	CUIT        string `json:"cuit" example:"20-12345678-9"`
	SocietyType string `json:"society_type" example:"SA"`
}

// CompanyBasicDetails represents minimal company person data
type CompanyBasicDetails struct {
	LegalName string `json:"legal_name" example:"Acme Corporation"`
}

// ContactDetails represents contact data from person-svc
type ContactDetails struct {
	ID        string    `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Tipo      string    `json:"tipo" example:"EMAIL"`
	Dato      string    `json:"dato" example:"contact@example.com"`
	IsPrimary bool      `json:"is_primary" example:"true"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
}
