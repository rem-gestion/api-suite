package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/rem-gestion/api-suite/property/src/models"
)

// PropertyType DTOs
type PropertyTypeCreate struct {
	Code        string  `json:"code" binding:"required,max=32"`
	Name        string  `json:"name" binding:"required,max=100"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}

type PropertyTypeUpdate struct {
	Code        *string `json:"code,omitempty" binding:"omitempty,max=32"`
	Name        *string `json:"name,omitempty" binding:"omitempty,max=100"`
	Description *string `json:"description,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

type PropertyTypeResponse struct {
	ID          int32      `json:"id"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

// ManagerType DTOs
type ManagerTypeCreate struct {
	Code        string  `json:"code" binding:"required,max=32"`
	Name        string  `json:"name" binding:"required,max=100"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}

type ManagerTypeUpdate struct {
	Code        *string `json:"code,omitempty" binding:"omitempty,max=32"`
	Name        *string `json:"name,omitempty" binding:"omitempty,max=100"`
	Description *string `json:"description,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

type ManagerTypeResponse struct {
	ID          int32      `json:"id"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

// Amenity DTOs
type AmenityCreate struct {
	Name     string                  `json:"name" binding:"required,max=100"`
	Category *models.AmenityCategory `json:"category"`
	IconURL  *string                 `json:"icon_url"`
}

type AmenityUpdate struct {
	Name     *string                 `json:"name,omitempty" binding:"omitempty,max=100"`
	Category *models.AmenityCategory `json:"category,omitempty"`
	IconURL  *string                 `json:"icon_url,omitempty"`
}

type AmenityResponse struct {
	ID       int32                   `json:"id"`
	Name     string                  `json:"name"`
	Category *models.AmenityCategory `json:"category"`
	IconURL  *string                 `json:"icon_url"`
}

// Property DTOs
type PropertyCreate struct {
	OwnerPersonID  uuid.UUID `json:"owner_person_id" binding:"required"`
	AddressID      uuid.UUID `json:"address_id" binding:"required"`
	PropertyTypeID int32     `json:"property_type_id" binding:"required"`
	InternalCode   *string   `json:"internal_code,omitempty" binding:"omitempty,max=50"`
	YearBuilt      *int32    `json:"year_built,omitempty" binding:"omitempty,min=1800,max=2100"`
	Bedrooms       *int32    `json:"bedrooms,omitempty" binding:"omitempty,min=0"`
	Bathrooms      *float32  `json:"bathrooms,omitempty" binding:"omitempty,min=0"`
	TotalAreaSqm   *float64  `json:"total_area_sqm,omitempty" binding:"omitempty,min=0"`
	CoveredAreaSqm *float64  `json:"covered_area_sqm,omitempty" binding:"omitempty,min=0"`
	Description    *string   `json:"description,omitempty"`
}

type PropertyUpdate struct {
	OwnerPersonID  *uuid.UUID `json:"owner_person_id,omitempty"`
	AddressID      *uuid.UUID `json:"address_id,omitempty"`
	PropertyTypeID *int32     `json:"property_type_id,omitempty"`
	InternalCode   *string    `json:"internal_code,omitempty" binding:"omitempty,max=50"`
	YearBuilt      *int32     `json:"year_built,omitempty" binding:"omitempty,min=1800,max=2100"`
	Bedrooms       *int32     `json:"bedrooms,omitempty" binding:"omitempty,min=0"`
	Bathrooms      *float32   `json:"bathrooms,omitempty" binding:"omitempty,min=0"`
	TotalAreaSqm   *float64   `json:"total_area_sqm,omitempty" binding:"omitempty,min=0"`
	CoveredAreaSqm *float64   `json:"covered_area_sqm,omitempty" binding:"omitempty,min=0"`
	Description    *string    `json:"description,omitempty"`
}

type PropertyResponse struct {
	ID             uuid.UUID  `json:"id"`
	OwnerPersonID  uuid.UUID  `json:"owner_person_id"`
	AddressID      uuid.UUID  `json:"address_id"`
	PropertyTypeID int32      `json:"property_type_id"`
	InternalCode   *string    `json:"internal_code"`
	YearBuilt      *int32     `json:"year_built"`
	Bedrooms       *int32     `json:"bedrooms"`
	Bathrooms      *float32   `json:"bathrooms"`
	TotalAreaSqm   *float64   `json:"total_area_sqm"`
	CoveredAreaSqm *float64   `json:"covered_area_sqm"`
	Description    *string    `json:"description"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
	UpdatedBy      *uuid.UUID `json:"updated_by"`

	// Embedded relations (optional)
	PropertyType        *PropertyTypeResponse        `json:"property_type,omitempty"`
	PropertyManagements []PropertyManagementResponse `json:"property_managements,omitempty"`
	PropertyAmenities   []PropertyAmenityResponse    `json:"property_amenities,omitempty"`
}

// PropertyManagement DTOs
type PropertyManagementCreate struct {
	PropertyID           uuid.UUID  `json:"property_id" binding:"required"`
	ManagerID            uuid.UUID  `json:"manager_id" binding:"required"`
	ManagerTypeID        int32      `json:"manager_type_id" binding:"required"`
	StartDate            time.Time  `json:"start_date" binding:"required"`
	EndDate              *time.Time `json:"end_date,omitempty"`
	CommissionPercentage *float64   `json:"commission_percentage,omitempty" binding:"omitempty,min=0,max=100"`
}

type PropertyManagementUpdate struct {
	PropertyID           *uuid.UUID `json:"property_id,omitempty"`
	ManagerID            *uuid.UUID `json:"manager_id,omitempty"`
	ManagerTypeID        *int32     `json:"manager_type_id,omitempty"`
	StartDate            *time.Time `json:"start_date,omitempty"`
	EndDate              *time.Time `json:"end_date,omitempty"`
	CommissionPercentage *float64   `json:"commission_percentage,omitempty" binding:"omitempty,min=0,max=100"`
}

type PropertyManagementResponse struct {
	ID                   uuid.UUID  `json:"id"`
	PropertyID           uuid.UUID  `json:"property_id"`
	ManagerID            uuid.UUID  `json:"manager_id"`
	ManagerTypeID        int32      `json:"manager_type_id"`
	StartDate            time.Time  `json:"start_date"`
	EndDate              *time.Time `json:"end_date"`
	CommissionPercentage *float64   `json:"commission_percentage"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            *time.Time `json:"updated_at"`
	UpdatedBy            *uuid.UUID `json:"updated_by"`

	// Embedded relations (optional)
	Property    *PropertyResponse    `json:"property,omitempty"`
	ManagerType *ManagerTypeResponse `json:"manager_type,omitempty"`
}

// PropertyAmenity DTOs
type PropertyAmenityCreate struct {
	PropertyID uuid.UUID `json:"property_id" binding:"required"`
	AmenityID  string    `json:"amenity_id" binding:"required"` // JSONB string
	Note       *string   `json:"note,omitempty"`
}

type PropertyAmenityUpdate struct {
	Note *string `json:"note,omitempty"`
}

type PropertyAmenityResponse struct {
	PropertyID uuid.UUID `json:"property_id"`
	AmenityID  string    `json:"amenity_id"`
	Note       *string   `json:"note"`

	// Embedded relations (optional)
	Property *PropertyResponse `json:"property,omitempty"`
}

// List/Pagination DTOs
type PropertyListRequest struct {
	OwnerPersonID  *uuid.UUID `form:"owner_person_id"`
	PropertyTypeID *int32     `form:"property_type_id"`
	InternalCode   *string    `form:"internal_code"`
	YearBuilt      *int32     `form:"year_built"`
	Bedrooms       *int32     `form:"bedrooms"`
	Bathrooms      *float32   `form:"bathrooms"`
	MinTotalArea   *float64   `form:"min_total_area"`
	MaxTotalArea   *float64   `form:"max_total_area"`
	City           *string    `form:"city"`
	State          *string    `form:"state"`
	Page           int        `form:"page,default=1" binding:"min=1"`
	Limit          int        `form:"limit,default=20" binding:"min=1,max=100"`
	IncludeDeleted bool       `form:"include_deleted,default=false"`
}

type ListResponse struct {
	Data  interface{} `json:"data"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
	Total int64       `json:"total"`
}

// Search DTOs
type PropertySearchRequest struct {
	Query          string     `form:"query" binding:"omitempty,min=3"`
	OwnerPersonID  *uuid.UUID `form:"owner_person_id"`
	PropertyTypeID *int32     `form:"property_type_id"`
	Page           int        `form:"page,default=1" binding:"min=1"`
	Limit          int        `form:"limit,default=20" binding:"min=1,max=100"`
}

// Bulk operations DTOs
type BulkPropertyCreateRequest struct {
	Properties []PropertyCreate `json:"properties" binding:"required,min=1,max=100"`
}

type BulkPropertyCreateResponse struct {
	Created []PropertyResponse `json:"created"`
	Errors  []BulkError        `json:"errors,omitempty"`
}

type BulkError struct {
	Index int    `json:"index"`
	Error string `json:"error"`
}

// PropertyManagement DTOs
type PropertyManagementCreateRequest struct {
	PropertyID           uuid.UUID `json:"property_id" binding:"required"`
	ManagerTypeID        int32     `json:"manager_type_id" binding:"required"`
	ManagerPersonID      uuid.UUID `json:"manager_person_id" binding:"required"`
	StartDate            time.Time `json:"start_date" binding:"required"`
	EndDate              *time.Time `json:"end_date,omitempty"`
	CommissionPercentage *float64  `json:"commission_percentage,omitempty" binding:"omitempty,min=0,max=100"`
	ContractReference    *string   `json:"contract_reference,omitempty"`
	Notes                *string   `json:"notes,omitempty"`
	IsActive             *bool     `json:"is_active,omitempty"`
}

type PropertyManagementUpdateRequest struct {
	ManagerTypeID        *int32     `json:"manager_type_id,omitempty"`
	ManagerPersonID      *uuid.UUID `json:"manager_person_id,omitempty"`
	StartDate            *time.Time `json:"start_date,omitempty"`
	EndDate              *time.Time `json:"end_date,omitempty"`
	CommissionPercentage *float64   `json:"commission_percentage,omitempty" binding:"omitempty,min=0,max=100"`
	ContractReference    *string    `json:"contract_reference,omitempty"`
	Notes                *string    `json:"notes,omitempty"`
	IsActive             *bool      `json:"is_active,omitempty"`
}

type PropertyManagementListRequest struct {
	Page       int        `form:"page" binding:"omitempty,min=1"`
	Limit      int        `form:"limit" binding:"omitempty,min=1,max=100"`
	PropertyID *uuid.UUID `form:"property_id"`
	IsActive   *bool      `form:"is_active"`
	StartDate  *time.Time `form:"start_date"`
	EndDate    *time.Time `form:"end_date"`
}

type PropertyManagementListResponse struct {
	Data  []PropertyManagementResponse `json:"data"`
	Page  int                          `json:"page"`
	Limit int                          `json:"limit"`
	Total int64                        `json:"total"`
}
