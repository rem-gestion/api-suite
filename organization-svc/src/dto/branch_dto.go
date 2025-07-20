package dto

import (
	"time"

	"github.com/google/uuid"
)

// ===============================
// BRANCH REQUEST DTOs
// ===============================

// CreateBranchRequest represents request to create a new branch
type CreateBranchRequest struct {
	DisplayName string     `json:"display_name" validate:"required,min=2,max=120" example:"Downtown Office"`
	AddressID   *uuid.UUID `json:"address_id,omitempty" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	Phone       *string    `json:"phone,omitempty" validate:"omitempty,phone" example:"+1-555-0123"`
	Email       *string    `json:"email,omitempty" validate:"omitempty,email" example:"downtown@tech-solutions.com"`
	IsMain      bool       `json:"is_main" example:"false"`
	MetadataRequest
}

// UpdateBranchRequest represents request to update branch data
type UpdateBranchRequest struct {
	DisplayName *string    `json:"display_name,omitempty" validate:"omitempty,min=2,max=120" example:"Downtown Office"`
	AddressID   *uuid.UUID `json:"address_id,omitempty" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	Phone       *string    `json:"phone,omitempty" validate:"omitempty,phone" example:"+1-555-0123"`
	Email       *string    `json:"email,omitempty" validate:"omitempty,email" example:"downtown@tech-solutions.com"`
	MetadataRequest
}

// BranchFiltersRequest represents filters for listing branches
type BranchFiltersRequest struct {
	FilterRequest
	OrganizationID uuid.UUID `json:"organization_id" form:"organization_id" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	IsMain         *bool     `json:"is_main" form:"is_main" example:"true"`
	HasAddress     *bool     `json:"has_address" form:"has_address" example:"true"`
	PaginationRequest
	IncludeOptions
}

// ===============================
// BRANCH RESPONSE DTOs
// ===============================

// OrganizationBranchResponse represents branch data in responses
type OrganizationBranchResponse struct {
	BaseResponse
	OrganizationID uuid.UUID  `json:"organization_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	DisplayName    string     `json:"display_name" example:"Downtown Office"`
	AddressID      *uuid.UUID `json:"address_id,omitempty" example:"123e4567-e89b-12d3-a456-426614174000"`
	Phone          *string    `json:"phone,omitempty" example:"+1-555-0123"`
	Email          *string    `json:"email,omitempty" example:"downtown@tech-solutions.com"`
	IsMain         bool       `json:"is_main" example:"false"`

	// Embedded optional data
	Organization *OrganizationBasicResponse `json:"organization,omitempty"`

	// External service data (populated via gRPC)
	Address *AddressDetails `json:"address,omitempty"`

	// Statistics (included when include_stats=true)
	Stats *BranchStats `json:"stats,omitempty"`
}

// BranchListResponse represents paginated branch list
type BranchListResponse struct {
	Data       []OrganizationBranchResponse `json:"data"`
	Pagination PaginationResponse           `json:"pagination"`
}

// BranchBasicResponse represents minimal branch data
type BranchBasicResponse struct {
	ID             uuid.UUID  `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	OrganizationID uuid.UUID  `json:"organization_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	DisplayName    string     `json:"display_name" example:"Downtown Office"`
	AddressID      *uuid.UUID `json:"address_id,omitempty" example:"123e4567-e89b-12d3-a456-426614174000"`
	IsMain         bool       `json:"is_main" example:"false"`
}

// ===============================
// BRANCH STATISTICS DTOs
// ===============================

// BranchStats represents branch statistics
type BranchStats struct {
	TotalEmployees   int       `json:"total_employees" example:"15"`
	ActiveEmployees  int       `json:"active_employees" example:"14"`
	TotalProperties  int       `json:"total_properties" example:"50"`
	ActiveProperties int       `json:"active_properties" example:"45"`
	LastActivityAt   time.Time `json:"last_activity_at" example:"2024-01-15T10:30:00Z"`
}

// ===============================
// ROLE REQUEST DTOs
// ===============================

// CreateRoleRequest represents request to create a new role
type CreateRoleRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=64" example:"Senior Agent"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500" example:"Senior real estate agent with advanced permissions"`
	IsDefault   bool    `json:"is_default" example:"false"`
	MetadataRequest
}

// UpdateRoleRequest represents request to update role data
type UpdateRoleRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=2,max=64" example:"Senior Agent"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500" example:"Senior real estate agent with advanced permissions"`
	MetadataRequest
}

// RoleFiltersRequest represents filters for listing roles
type RoleFiltersRequest struct {
	FilterRequest
	OrganizationID    uuid.UUID `json:"organization_id" form:"organization_id" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	IsDefault         *bool     `json:"is_default" form:"is_default" example:"true"`
	IncludeUsageStats *bool     `json:"include_usage_stats" form:"include_usage_stats" example:"true"`
	PaginationRequest
	IncludeOptions
}

// DeleteRoleRequest represents request to delete a role
type DeleteRoleRequest struct {
	ReplacementRoleID *uuid.UUID `json:"replacement_role_id,omitempty" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174001"`
	Reason            string     `json:"reason,omitempty" validate:"omitempty,max=500" example:"Role consolidation"`
}

// ===============================
// ROLE RESPONSE DTOs
// ===============================

// OrganizationRoleResponse represents role data in responses
type OrganizationRoleResponse struct {
	BaseResponse
	OrganizationID uuid.UUID `json:"organization_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name           string    `json:"name" example:"Senior Agent"`
	Description    *string   `json:"description,omitempty" example:"Senior real estate agent with advanced permissions"`
	IsDefault      bool      `json:"is_default" example:"false"`

	// Embedded optional data
	Organization  *OrganizationBasicResponse `json:"organization,omitempty"`
	EmployeeRoles []EmployeeRoleResponse     `json:"employee_roles,omitempty"`

	// Statistics (included when include_stats=true)
	Stats *RoleStats `json:"stats,omitempty"`
}

// RoleListResponse represents paginated role list
type RoleListResponse struct {
	Data       []OrganizationRoleResponse `json:"data"`
	Pagination PaginationResponse         `json:"pagination"`
}

// RoleBasicResponse represents minimal role data
type RoleBasicResponse struct {
	ID             uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	OrganizationID uuid.UUID `json:"organization_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name           string    `json:"name" example:"Senior Agent"`
	Description    *string   `json:"description,omitempty" example:"Senior real estate agent with advanced permissions"`
	IsDefault      bool      `json:"is_default" example:"false"`
}

// ===============================
// ROLE STATISTICS DTOs
// ===============================

// RoleStats represents role statistics
type RoleStats struct {
	TotalEmployees   int       `json:"total_employees" example:"8"`
	ActiveEmployees  int       `json:"active_employees" example:"7"`
	PrimaryRoleCount int       `json:"primary_role_count" example:"5"`
	UsagePercentage  float64   `json:"usage_percentage" example:"35.5"`
	LastAssignedAt   time.Time `json:"last_assigned_at" example:"2024-01-15T10:30:00Z"`
}
