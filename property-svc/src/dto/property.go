package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/rem-gestion/api-suite/property/src/models"
)

/* ─────────────────────────  PROPERTY  ───────────────────────── */

type CreatePropertyDTO struct {
	OwnerPersonID  uuid.UUID `json:"owner_person_id"   validate:"required"`
	AddressID      uuid.UUID `json:"address_id"        validate:"required"`
	PropertyType   string    `json:"property_type"     validate:"required,oneof=APARTMENT HOUSE COMMERCIAL_SPACE OFFICE LAND INDUSTRIAL_WAREHOUSE"`
	InternalCode   *string   `json:"internal_code,omitempty"`
	YearBuilt      *int      `json:"year_built,omitempty"        validate:"omitempty,min=1800,max=2100"`
	Bedrooms       *int      `json:"bedrooms,omitempty"          validate:"omitempty,min=0,max=50"`
	Bathrooms      *float32  `json:"bathrooms,omitempty"         validate:"omitempty,min=0,max=20"`
	TotalAreaSqm   *float64  `json:"total_area_sqm,omitempty"    validate:"omitempty,min=0"`
	CoveredAreaSqm *float64  `json:"covered_area_sqm,omitempty"  validate:"omitempty,min=0"`
	Description    *string   `json:"description,omitempty"`

	// Management opcional al crear
	Management *CreatePropertyManagementDTO `json:"management,omitempty"`

	// Amenities opcionales al crear
	AmenityIDs []uuid.UUID `json:"amenity_ids,omitempty"`
}

type UpdatePropertyDTO struct {
	PropertyType   *string  `json:"property_type,omitempty"     validate:"omitempty,oneof=APARTMENT HOUSE COMMERCIAL_SPACE OFFICE LAND INDUSTRIAL_WAREHOUSE"`
	InternalCode   *string  `json:"internal_code,omitempty"`
	YearBuilt      *int     `json:"year_built,omitempty"        validate:"omitempty,min=1800,max=2100"`
	Bedrooms       *int     `json:"bedrooms,omitempty"          validate:"omitempty,min=0,max=50"`
	Bathrooms      *float32 `json:"bathrooms,omitempty"         validate:"omitempty,min=0,max=20"`
	TotalAreaSqm   *float64 `json:"total_area_sqm,omitempty"    validate:"omitempty,min=0"`
	CoveredAreaSqm *float64 `json:"covered_area_sqm,omitempty"  validate:"omitempty,min=0"`
	Description    *string  `json:"description,omitempty"`
}

type PropertyResponseDTO struct {
	ID             uuid.UUID `json:"id"`
	OwnerPersonID  uuid.UUID `json:"owner_person_id"`
	AddressID      uuid.UUID `json:"address_id"`
	PropertyType   string    `json:"property_type"`
	InternalCode   *string   `json:"internal_code,omitempty"`
	YearBuilt      *int      `json:"year_built,omitempty"`
	Bedrooms       *int      `json:"bedrooms,omitempty"`
	Bathrooms      *float32  `json:"bathrooms,omitempty"`
	TotalAreaSqm   *float64  `json:"total_area_sqm,omitempty"`
	CoveredAreaSqm *float64  `json:"covered_area_sqm,omitempty"`
	Description    *string   `json:"description,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`

	// Relations
	Management []PropertyManagementResponseDTO `json:"management,omitempty"`
	Amenities  []AmenityResponseDTO            `json:"amenities,omitempty"`
}

/* ─────────────────────────  PROPERTY MANAGEMENT  ───────────────────────── */

type CreatePropertyManagementDTO struct {
	PropertyID        uuid.UUID  `json:"property_id"        validate:"required"`
	OrganizationID    uuid.UUID  `json:"organization_id"    validate:"required"`
	ManagedSince      *time.Time `json:"managed_since,omitempty"`
	ManagedUntil      *time.Time `json:"managed_until,omitempty"`
	CommissionPercent *float64   `json:"commission_percent,omitempty" validate:"omitempty,min=0,max=100"`
	Notes             *string    `json:"notes,omitempty"`
}

type UpdatePropertyManagementDTO struct {
	ManagedUntil      *time.Time `json:"managed_until,omitempty"`
	IsActive          *bool      `json:"is_active,omitempty"`
	CommissionPercent *float64   `json:"commission_percent,omitempty" validate:"omitempty,min=0,max=100"`
	Notes             *string    `json:"notes,omitempty"`
}

type PropertyManagementResponseDTO struct {
	ID                uuid.UUID  `json:"id"`
	PropertyID        uuid.UUID  `json:"property_id"`
	OrganizationID    uuid.UUID  `json:"organization_id"`
	ManagedSince      time.Time  `json:"managed_since"`
	ManagedUntil      *time.Time `json:"managed_until,omitempty"`
	IsActive          bool       `json:"is_active"`
	CommissionPercent *float64   `json:"commission_percent,omitempty"`
	Notes             *string    `json:"notes,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

/* ─────────────────────────  AMENITY  ───────────────────────── */

type CreateAmenityDTO struct {
	Name        string  `json:"name"        validate:"required,min=1,max=100"`
	Description *string `json:"description,omitempty"`
	Icon        *string `json:"icon,omitempty"        validate:"omitempty,max=50"`
	Category    *string `json:"category,omitempty"    validate:"omitempty,max=50"`
}

type UpdateAmenityDTO struct {
	Name        *string `json:"name,omitempty"        validate:"omitempty,min=1,max=100"`
	Description *string `json:"description,omitempty"`
	Icon        *string `json:"icon,omitempty"        validate:"omitempty,max=50"`
	Category    *string `json:"category,omitempty"    validate:"omitempty,max=50"`
}

type AmenityResponseDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Icon        *string   `json:"icon,omitempty"`
	Category    *string   `json:"category,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

/* ─────────────────────────  CONVERTERS  ───────────────────────── */

func ToPropertyResponse(property *models.Property) PropertyResponseDTO {
	response := PropertyResponseDTO{
		ID:             property.ID,
		OwnerPersonID:  property.OwnerPersonID,
		AddressID:      property.AddressID,
		PropertyType:   string(property.PropertyType),
		InternalCode:   property.InternalCode,
		YearBuilt:      property.YearBuilt,
		Bedrooms:       property.Bedrooms,
		Bathrooms:      property.Bathrooms,
		TotalAreaSqm:   property.TotalAreaSqm,
		CoveredAreaSqm: property.CoveredAreaSqm,
		Description:    property.Description,
		CreatedAt:      property.CreatedAt,
		UpdatedAt:      property.UpdatedAt,
	}

	// Convert management
	for _, mgmt := range property.PropertyManagement {
		response.Management = append(response.Management, PropertyManagementResponseDTO{
			ID:                mgmt.ID,
			PropertyID:        mgmt.PropertyID,
			OrganizationID:    mgmt.OrganizationID,
			ManagedSince:      mgmt.ManagedSince,
			ManagedUntil:      mgmt.ManagedUntil,
			IsActive:          mgmt.IsActive,
			CommissionPercent: mgmt.CommissionPercent,
			Notes:             mgmt.Notes,
			CreatedAt:         mgmt.CreatedAt,
			UpdatedAt:         mgmt.UpdatedAt,
		})
	}

	// Convert amenities
	for _, pa := range property.PropertyAmenities {
		response.Amenities = append(response.Amenities, AmenityResponseDTO{
			ID:          pa.Amenity.ID,
			Name:        pa.Amenity.Name,
			Description: pa.Amenity.Description,
			Icon:        pa.Amenity.Icon,
			Category:    pa.Amenity.Category,
			CreatedAt:   pa.Amenity.CreatedAt,
			UpdatedAt:   pa.Amenity.UpdatedAt,
		})
	}

	return response
}

func ToAmenityResponse(amenity *models.Amenity) AmenityResponseDTO {
	return AmenityResponseDTO{
		ID:          amenity.ID,
		Name:        amenity.Name,
		Description: amenity.Description,
		Icon:        amenity.Icon,
		Category:    amenity.Category,
		CreatedAt:   amenity.CreatedAt,
		UpdatedAt:   amenity.UpdatedAt,
	}
}
