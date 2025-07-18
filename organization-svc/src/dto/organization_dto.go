package dto

import (
	"time"

	"github.com/google/uuid"
)

// ===============================
// ORGANIZATION REQUEST DTOs
// ===============================

// CreateOrganizationRequest represents request to create a new organization
type CreateOrganizationRequest struct {
	Name            string     `json:"name" validate:"required,min=2,max=200" example:"Tech Solutions Inc"`
	DisplayName     string     `json:"display_name" validate:"required,min=2,max=300" example:"Tech Solutions - Real Estate Division"`
	Slug            string     `json:"slug" validate:"required,min=2,max=100,slug" example:"tech-solutions-re"`
	Description     *string    `json:"description,omitempty" validate:"omitempty,max=1000" example:"Leading technology solutions for real estate"`
	Type            string     `json:"type" validate:"required,oneof=real_estate property_management construction architecture appraisal mortgage_broker insurance legal consulting investment" example:"real_estate"`
	LegalName       *string    `json:"legal_name,omitempty" validate:"omitempty,max=300" example:"Tech Solutions Incorporated"`
	TaxID           *string    `json:"tax_id,omitempty" validate:"omitempty,max=50" example:"12-3456789"`
	Website         *string    `json:"website,omitempty" validate:"omitempty,url" example:"https://tech-solutions.com"`
	Phone           *string    `json:"phone,omitempty" validate:"omitempty,phone" example:"+1-555-0123"`
	Email           *string    `json:"email,omitempty" validate:"omitempty,email" example:"contact@tech-solutions.com"`
	LogoURL         *string    `json:"logo_url,omitempty" validate:"omitempty,url" example:"https://cdn.tech-solutions.com/logo.png"`
	TimezoneID      string     `json:"timezone_id" validate:"required,timezone" example:"America/New_York"`
	FiscalAddressID *uuid.UUID `json:"fiscal_address_id,omitempty" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	MetadataRequest
}

// UpdateOrganizationRequest represents request to update organization data
type UpdateOrganizationRequest struct {
	Name            *string    `json:"name,omitempty" validate:"omitempty,min=2,max=200" example:"Tech Solutions Inc"`
	DisplayName     *string    `json:"display_name,omitempty" validate:"omitempty,min=2,max=300" example:"Tech Solutions - Real Estate Division"`
	Description     *string    `json:"description,omitempty" validate:"omitempty,max=1000" example:"Leading technology solutions for real estate"`
	LegalName       *string    `json:"legal_name,omitempty" validate:"omitempty,max=300" example:"Tech Solutions Incorporated"`
	TaxID           *string    `json:"tax_id,omitempty" validate:"omitempty,max=50" example:"12-3456789"`
	Website         *string    `json:"website,omitempty" validate:"omitempty,url" example:"https://tech-solutions.com"`
	Phone           *string    `json:"phone,omitempty" validate:"omitempty,phone" example:"+1-555-0123"`
	Email           *string    `json:"email,omitempty" validate:"omitempty,email" example:"contact@tech-solutions.com"`
	LogoURL         *string    `json:"logo_url,omitempty" validate:"omitempty,url" example:"https://cdn.tech-solutions.com/logo.png"`
	TimezoneID      *string    `json:"timezone_id,omitempty" validate:"omitempty,timezone" example:"America/New_York"`
	FiscalAddressID *uuid.UUID `json:"fiscal_address_id,omitempty" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	MetadataRequest
}

// OrganizationFiltersRequest represents filters for listing organizations
type OrganizationFiltersRequest struct {
	FilterRequest
	Type       string  `json:"type" form:"type" validate:"omitempty,oneof=real_estate property_management construction architecture appraisal mortgage_broker insurance legal consulting investment" example:"real_estate"`
	IsVerified *bool   `json:"is_verified" form:"is_verified" example:"true"`
	OwnerID    *string `json:"owner_id" form:"owner_id" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	PaginationRequest
	IncludeOptions
}

// ===============================
// ORGANIZATION RESPONSE DTOs
// ===============================

// OrganizationResponse represents organization data in responses
type OrganizationResponse struct {
	BaseResponse
	Name            string     `json:"name" example:"Tech Solutions Inc"`
	DisplayName     string     `json:"display_name" example:"Tech Solutions - Real Estate Division"`
	Slug            string     `json:"slug" example:"tech-solutions-re"`
	Description     *string    `json:"description,omitempty" example:"Leading technology solutions for real estate"`
	Status          string     `json:"status" example:"active"`
	Type            string     `json:"type" example:"real_estate"`
	LegalName       *string    `json:"legal_name,omitempty" example:"Tech Solutions Incorporated"`
	TaxID           *string    `json:"tax_id,omitempty" example:"12-3456789"`
	Website         *string    `json:"website,omitempty" example:"https://tech-solutions.com"`
	Phone           *string    `json:"phone,omitempty" example:"+1-555-0123"`
	Email           *string    `json:"email,omitempty" example:"contact@tech-solutions.com"`
	LogoURL         *string    `json:"logo_url,omitempty" example:"https://cdn.tech-solutions.com/logo.png"`
	TimezoneID      string     `json:"timezone_id" example:"America/New_York"`
	FiscalAddressID *uuid.UUID `json:"fiscal_address_id,omitempty" example:"123e4567-e89b-12d3-a456-426614174000"`
	IsVerified      bool       `json:"is_verified" example:"true"`
	VerifiedAt      *time.Time `json:"verified_at,omitempty" example:"2024-01-15T10:30:00Z"`
	SubscriptionID  *uuid.UUID `json:"subscription_id,omitempty" example:"123e4567-e89b-12d3-a456-426614174000"`
	PlanID          *uuid.UUID `json:"plan_id,omitempty" example:"123e4567-e89b-12d3-a456-426614174000"`

	// Embedded optional data
	Settings     []OrganizationSettingResponse `json:"settings,omitempty"`
	Owner        *OrganizationOwnerResponse    `json:"owner,omitempty"`
	Branches     []OrganizationBranchResponse  `json:"branches,omitempty"`
	Roles        []OrganizationRoleResponse    `json:"roles,omitempty"`
	Employees    []EmployeeResponse            `json:"employees,omitempty"`
	Invitations  []OrganizationInviteResponse  `json:"invitations,omitempty"`
	Subscription *OrganizationSubscription     `json:"subscription,omitempty"`

	// External service data (populated via gRPC)
	FiscalAddress *AddressDetails `json:"fiscal_address,omitempty"`

	// Statistics (included when include_stats=true)
	Stats *OrganizationStats `json:"stats,omitempty"`
}

// OrganizationListResponse represents paginated organization list
type OrganizationListResponse struct {
	Data       []OrganizationResponse `json:"data"`
	Pagination PaginationResponse     `json:"pagination"`
}

// OrganizationBasicResponse represents minimal organization data
type OrganizationBasicResponse struct {
	ID          uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name        string    `json:"name" example:"Tech Solutions Inc"`
	DisplayName string    `json:"display_name" example:"Tech Solutions - Real Estate Division"`
	Slug        string    `json:"slug" example:"tech-solutions-re"`
	Status      string    `json:"status" example:"active"`
	Type        string    `json:"type" example:"real_estate"`
	LogoURL     *string   `json:"logo_url,omitempty" example:"https://cdn.tech-solutions.com/logo.png"`
	IsVerified  bool      `json:"is_verified" example:"true"`
}

// ===============================
// ORGANIZATION SETTINGS DTOs
// ===============================

// OrganizationSettingResponse represents organization setting data
type OrganizationSettingResponse struct {
	BaseResponse
	OrganizationID uuid.UUID   `json:"organization_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	SettingKey     string      `json:"setting_key" example:"currency"`
	SettingValue   interface{} `json:"setting_value" example:"USD"`
	Description    *string     `json:"description,omitempty" example:"Default currency for transactions"`
	IsEditable     bool        `json:"is_editable" example:"true"`
}

// UpdateOrganizationSettingsRequest represents batch settings update
type UpdateOrganizationSettingsRequest struct {
	Settings map[string]interface{} `json:"settings" validate:"required,min=1" example:"{\"currency\":\"USD\",\"timezone\":\"America/New_York\"}"`
}

// OrganizationSettingsResponse represents all organization settings
type OrganizationSettingsResponse struct {
	OrganizationID uuid.UUID              `json:"organization_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Settings       map[string]interface{} `json:"settings" example:"{\"currency\":\"USD\",\"timezone\":\"America/New_York\"}"`
	UpdatedAt      time.Time              `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// ===============================
// ORGANIZATION SUBSCRIPTION DTOs
// ===============================

// OrganizationSubscription represents subscription information
type OrganizationSubscription struct {
	ID              uuid.UUID   `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	OrganizationID  uuid.UUID   `json:"organization_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	PlanName        string      `json:"plan_name" example:"Professional"`
	Status          string      `json:"status" example:"active"`
	StartDate       time.Time   `json:"start_date" example:"2024-01-01"`
	EndDate         *time.Time  `json:"end_date,omitempty" example:"2024-12-31"`
	TrialEndDate    *time.Time  `json:"trial_end_date,omitempty" example:"2024-02-01"`
	BillingCycle    string      `json:"billing_cycle" example:"monthly"`
	Amount          float64     `json:"amount" example:"99.99"`
	Currency        string      `json:"currency" example:"USD"`
	Features        []string    `json:"features" example:"[\"unlimited_properties\",\"advanced_analytics\"]"`
	Limits          interface{} `json:"limits" example:"{\"max_users\":50,\"max_properties\":1000}"`
	PaymentMethodID *string     `json:"payment_method_id,omitempty" example:"pm_1234567890"`
	CreatedAt       time.Time   `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt       time.Time   `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// ===============================
// ORGANIZATION STATISTICS DTOs
// ===============================

// OrganizationStats represents organization statistics
type OrganizationStats struct {
	TotalEmployees   int       `json:"total_employees" example:"25"`
	ActiveEmployees  int       `json:"active_employees" example:"23"`
	TotalBranches    int       `json:"total_branches" example:"3"`
	TotalInvitations int       `json:"total_invitations" example:"5"`
	PendingInvites   int       `json:"pending_invites" example:"2"`
	CompletionRate   float64   `json:"completion_rate" example:"92.5"`
	LastActivityAt   time.Time `json:"last_activity_at" example:"2024-01-15T10:30:00Z"`
}

// ===============================
// EXTERNAL SERVICE DTOs (gRPC)
// ===============================

// AddressDetails represents address data from address-svc
type AddressDetails struct {
	ID      string  `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Floor   *string `json:"floor,omitempty" example:"2"`
	Unit    *string `json:"unit,omitempty" example:"A"`
	Street  string  `json:"street" example:"Main Street"`
	Number  int32   `json:"number" example:"123"`
	City    string  `json:"city" example:"New York"`
	State   string  `json:"state" example:"NY"`
	Zip     string  `json:"zip" example:"10001"`
	Country string  `json:"country" example:"US"`
}
