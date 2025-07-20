package dto

import (
	"time"

	"github.com/google/uuid"
)

// ===============================
// COMMON DTOs - Shared structures
// ===============================

// PaginationRequest represents pagination parameters
type PaginationRequest struct {
	Page    int `json:"page" form:"page" validate:"omitempty,min=1" example:"1"`
	PerPage int `json:"per_page" form:"per_page" validate:"omitempty,min=1,max=100" example:"10"`
}

// PaginationResponse represents pagination metadata in responses
type PaginationResponse struct {
	Page       int   `json:"page" example:"1"`
	PerPage    int   `json:"per_page" example:"10"`
	Total      int64 `json:"total" example:"150"`
	TotalPages int   `json:"total_pages" example:"15"`
	HasNext    bool  `json:"has_next" example:"true"`
	HasPrev    bool  `json:"has_prev" example:"false"`
}

// FilterRequest represents common filter parameters
type FilterRequest struct {
	Search    string     `json:"search" form:"search" validate:"omitempty,max=100" example:"tech company"`
	Status    string     `json:"status" form:"status" validate:"omitempty,oneof=active inactive suspended pending terminated expired cancelled" example:"active"`
	CreatedAt *time.Time `json:"created_at" form:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt *time.Time `json:"updated_at" form:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// BaseResponse represents common response fields
type BaseResponse struct {
	ID        uuid.UUID  `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	CreatedAt time.Time  `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time  `json:"updated_at" example:"2024-01-01T00:00:00Z"`
	CreatedBy uuid.UUID  `json:"created_by" example:"123e4567-e89b-12d3-a456-426614174001"`
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty" example:"123e4567-e89b-12d3-a456-426614174002"`
}

// StatusUpdateRequest represents status change requests
type StatusUpdateRequest struct {
	Status string `json:"status" validate:"required,oneof=active inactive suspended pending terminated expired cancelled" example:"active"`
	Reason string `json:"reason,omitempty" validate:"omitempty,max=500" example:"Administrative update"`
}

// MetadataRequest represents optional metadata in requests
type MetadataRequest struct {
	Metadata map[string]interface{} `json:"metadata,omitempty" example:"{\"source\":\"admin_panel\"}"`
}

// ErrorResponse represents error responses
type ErrorResponse struct {
	Error   string                 `json:"error" example:"Validation failed"`
	Message string                 `json:"message" example:"Invalid input data"`
	Code    string                 `json:"code,omitempty" example:"VALIDATION_ERROR"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// SuccessResponse represents successful operation responses
type SuccessResponse struct {
	Message string                 `json:"message" example:"Operation completed successfully"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// IDResponse represents responses that only return an ID
type IDResponse struct {
	ID uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// BulkRequest represents bulk operation requests
type BulkRequest struct {
	IDs    []uuid.UUID            `json:"ids" validate:"required,min=1,max=100" example:"[\"123e4567-e89b-12d3-a456-426614174000\"]"`
	Action string                 `json:"action" validate:"required,oneof=activate deactivate delete update" example:"activate"`
	Data   map[string]interface{} `json:"data,omitempty" example:"{\"status\":\"active\"}"`
}

// BulkResponse represents bulk operation responses
type BulkResponse struct {
	Successful []uuid.UUID          `json:"successful" example:"[\"123e4567-e89b-12d3-a456-426614174000\"]"`
	Failed     []BulkFailureDetails `json:"failed,omitempty"`
	Total      int                  `json:"total" example:"10"`
	Success    int                  `json:"success" example:"8"`
	Errors     int                  `json:"errors" example:"2"`
}

// BulkFailureDetails represents failed operations in bulk requests
type BulkFailureDetails struct {
	ID     uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Error  string    `json:"error" example:"Record not found"`
	Reason string    `json:"reason" example:"The specified record does not exist"`
}

// IncludeOptions represents include options for nested data
type IncludeOptions struct {
	IncludeDetails    bool `json:"include_details" form:"include_details" example:"true"`
	IncludeStats      bool `json:"include_stats" form:"include_stats" example:"false"`
	IncludeHistorical bool `json:"include_historical" form:"include_historical" example:"false"`
	IncludeAddress    bool `json:"include_address" form:"include_address" example:"true"`
	IncludePerson     bool `json:"include_person" form:"include_person" example:"true"`
	IncludeRoles      bool `json:"include_roles" form:"include_roles" example:"true"`
}

// UpdateOrganizationStatusRequest represents request to update organization status
type UpdateOrganizationStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=active inactive suspended pending terminated expired cancelled" example:"active"`
	Reason string `json:"reason,omitempty" validate:"omitempty,max=500" example:"Administrative update"`
}

// MessageResponse represents simple message responses
type MessageResponse struct {
	Message string `json:"message" example:"Operation completed successfully"`
}

// UserClaims represents JWT user claims structure
type UserClaims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	PersonID  string `json:"person_id"`
	AccountID string `json:"account_id"`
}
