package dto

import (
	"time"

	"github.com/google/uuid"
)

// ===============================
// OWNER REQUEST DTOs
// ===============================

// AddOwnerRequest represents request to add a new owner to organization
type AddOwnerRequest struct {
	PersonID     uuid.UUID `json:"person_id" validate:"required,uuid" example:"123e4567-e89b-12d3-a456-426614174001"`
	UserID       uuid.UUID `json:"user_id" validate:"required,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	IsFounder    bool      `json:"is_founder" example:"false"`
	OwnershipPct *float64  `json:"ownership_pct,omitempty" validate:"omitempty,min=0,max=100" example:"25.5"`
	MetadataRequest
}

// UpdateOwnerRequest represents request to update owner data
type UpdateOwnerRequest struct {
	IsFounder    *bool    `json:"is_founder,omitempty" example:"false"`
	OwnershipPct *float64 `json:"ownership_pct,omitempty" validate:"omitempty,min=0,max=100" example:"25.5"`
	MetadataRequest
}

// OwnerFiltersRequest represents filters for listing owners
type OwnerFiltersRequest struct {
	FilterRequest
	OrganizationID uuid.UUID  `json:"organization_id" form:"organization_id" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	IsFounder      *bool      `json:"is_founder" form:"is_founder" example:"true"`
	MinOwnership   *float64   `json:"min_ownership" form:"min_ownership" validate:"omitempty,min=0,max=100" example:"10.0"`
	MaxOwnership   *float64   `json:"max_ownership" form:"max_ownership" validate:"omitempty,min=0,max=100" example:"50.0"`
	PersonID       *uuid.UUID `json:"person_id" form:"person_id" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174001"`
	UserID         *uuid.UUID `json:"user_id" form:"user_id" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	PaginationRequest
	IncludeOptions
}

// RemoveOwnerRequest represents request to remove an owner
type RemoveOwnerRequest struct {
	Reason                string     `json:"reason" validate:"required,max=500" example:"Ownership transfer completed"`
	RedistributeOwnership bool       `json:"redistribute_ownership" example:"true"`
	TransferToOwnerID     *uuid.UUID `json:"transfer_to_owner_id,omitempty" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174002"`
}

// ===============================
// OWNER RESPONSE DTOs
// ===============================

// OrganizationOwnerResponse represents owner data in responses
type OrganizationOwnerResponse struct {
	BaseResponse
	OrganizationID uuid.UUID `json:"organization_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	PersonID       uuid.UUID `json:"person_id" example:"123e4567-e89b-12d3-a456-426614174001"`
	UserID         uuid.UUID `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	IsFounder      bool      `json:"is_founder" example:"false"`
	OwnershipPct   *float64  `json:"ownership_pct,omitempty" example:"25.5"`

	// Embedded optional data
	Organization *OrganizationBasicResponse `json:"organization,omitempty"`

	// External service data (populated via gRPC)
	PersonDetails *PersonDetails    `json:"person_details,omitempty"`
	UserDetails   *UserBasicDetails `json:"user_details,omitempty"`

	// Statistics (included when include_stats=true)
	Stats *OwnerStats `json:"stats,omitempty"`
}

// OwnerListResponse represents paginated owner list
type OwnerListResponse struct {
	Data       []OrganizationOwnerResponse `json:"data"`
	Pagination PaginationResponse          `json:"pagination"`
	Summary    *OwnershipSummary           `json:"summary,omitempty"`
}

// OwnerBasicResponse represents minimal owner data
type OwnerBasicResponse struct {
	ID             uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	OrganizationID uuid.UUID `json:"organization_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	PersonID       uuid.UUID `json:"person_id" example:"123e4567-e89b-12d3-a456-426614174001"`
	UserID         uuid.UUID `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	IsFounder      bool      `json:"is_founder" example:"false"`
	OwnershipPct   *float64  `json:"ownership_pct,omitempty" example:"25.5"`

	// External service data (basic)
	PersonDetails *PersonBasicDetails `json:"person_details,omitempty"`
	UserDetails   *UserBasicDetails   `json:"user_details,omitempty"`
}

// ===============================
// OWNERSHIP SUMMARY DTOs
// ===============================

// OwnershipSummary represents ownership distribution summary
type OwnershipSummary struct {
	TotalOwners          int                  `json:"total_owners" example:"4"`
	TotalFounders        int                  `json:"total_founders" example:"1"`
	TotalOwnership       float64              `json:"total_ownership" example:"100.0"`
	AllocatedOwnership   float64              `json:"allocated_ownership" example:"85.5"`
	UnallocatedOwnership float64              `json:"unallocated_ownership" example:"14.5"`
	MajorityOwner        *OwnerBasicResponse  `json:"majority_owner,omitempty"`
	OwnershipBreakdown   []OwnershipBreakdown `json:"ownership_breakdown"`
	LastUpdated          time.Time            `json:"last_updated" example:"2024-01-15T10:30:00Z"`
}

// OwnershipBreakdown represents ownership percentage ranges
type OwnershipBreakdown struct {
	Range      string  `json:"range" example:"0-25%"`
	Count      int     `json:"count" example:"2"`
	Percentage float64 `json:"percentage" example:"35.0"`
}

// ===============================
// OWNER STATISTICS DTOs
// ===============================

// OwnerStats represents owner statistics
type OwnerStats struct {
	YearsAsOwner      float64    `json:"years_as_owner" example:"2.5"`
	DividendsReceived *float64   `json:"dividends_received,omitempty" example:"15000.00"`
	LastDividendDate  *time.Time `json:"last_dividend_date,omitempty" example:"2024-01-01T00:00:00Z"`
	VotingPower       *float64   `json:"voting_power,omitempty" example:"25.5"`
	LastActivityAt    time.Time  `json:"last_activity_at" example:"2024-01-15T10:30:00Z"`
}

// ===============================
// OWNERSHIP VALIDATION DTOs
// ===============================

// OwnershipValidationResponse represents ownership validation results
type OwnershipValidationResponse struct {
	IsValid         bool                       `json:"is_valid" example:"true"`
	TotalPercentage float64                    `json:"total_percentage" example:"100.0"`
	Issues          []OwnershipValidationIssue `json:"issues,omitempty"`
	Recommendations []string                   `json:"recommendations,omitempty" example:"[\"Consider allocating remaining 14.5% ownership\"]"`
	LastValidated   time.Time                  `json:"last_validated" example:"2024-01-15T10:30:00Z"`
}

// OwnershipValidationIssue represents ownership validation issues
type OwnershipValidationIssue struct {
	Type        string     `json:"type" example:"OVER_ALLOCATION"`
	Severity    string     `json:"severity" example:"ERROR"`
	Description string     `json:"description" example:"Total ownership exceeds 100%"`
	OwnerID     *uuid.UUID `json:"owner_id,omitempty" example:"123e4567-e89b-12d3-a456-426614174000"`
	Suggestion  string     `json:"suggestion" example:"Reduce ownership percentages to total 100%"`
}

// ===============================
// OWNER TRANSFER DTOs
// ===============================

// TransferOwnershipRequest represents request to transfer ownership
type TransferOwnershipRequest struct {
	FromOwnerID      uuid.UUID  `json:"from_owner_id" validate:"required,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	ToOwnerID        uuid.UUID  `json:"to_owner_id" validate:"required,uuid" example:"123e4567-e89b-12d3-a456-426614174001"`
	PercentageAmount float64    `json:"percentage_amount" validate:"required,min=0.01,max=100" example:"10.5"`
	TransferReason   string     `json:"transfer_reason" validate:"required,max=500" example:"Partial sale of ownership stake"`
	EffectiveDate    *time.Time `json:"effective_date,omitempty" example:"2024-02-01T00:00:00Z"`
}

// TransferOwnershipResponse represents ownership transfer response
type TransferOwnershipResponse struct {
	TransferID        uuid.UUID           `json:"transfer_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Success           bool                `json:"success" example:"true"`
	FromOwner         *OwnerBasicResponse `json:"from_owner"`
	ToOwner           *OwnerBasicResponse `json:"to_owner"`
	TransferredAmount float64             `json:"transferred_amount" example:"10.5"`
	NewFromPercentage float64             `json:"new_from_percentage" example:"14.5"`
	NewToPercentage   float64             `json:"new_to_percentage" example:"35.5"`
	TransferReason    string              `json:"transfer_reason" example:"Partial sale of ownership stake"`
	ProcessedAt       time.Time           `json:"processed_at" example:"2024-01-15T10:30:00Z"`
	EffectiveDate     time.Time           `json:"effective_date" example:"2024-02-01T00:00:00Z"`
}
