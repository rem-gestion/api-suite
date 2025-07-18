package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/rem-gestion/api-suite/property/src/models"
)

/* ─────────────────────────  PROPERTY  ───────────────────────── */

type CreatePropertyDTO struct {
	OwnerPersonID  uuid.UUID          `json:"owner_person_id"   validate:"required"`
	AddressID      *uuid.UUID         `json:"address_id,omitempty"`      // escenario 1: usar dirección existente
	Address        *AddressPayloadDTO `json:"address_payload,omitempty"` // escenario 2: crear nueva dirección
	PropertyTypeID int                `json:"property_type_id"  validate:"required,min=1"`
	InternalCode   *string            `json:"internal_code,omitempty"`
	YearBuilt      *int               `json:"year_built,omitempty"        validate:"omitempty,min=1800,max=2100"`
	Bedrooms       *int               `json:"bedrooms,omitempty"          validate:"omitempty,min=0,max=50"`
	Bathrooms      *float32           `json:"bathrooms,omitempty"         validate:"omitempty,min=0,max=20"`
	TotalAreaSqm   *float64           `json:"total_area_sqm,omitempty"    validate:"omitempty,min=0"`
	CoveredAreaSqm *float64           `json:"covered_area_sqm,omitempty"  validate:"omitempty,min=0"`
	Description    *string            `json:"description,omitempty"`

	// Management opcional al crear
	Management *CreatePropertyManagementDTO `json:"management,omitempty"`

	// Amenities opcionales al crear
	AmenityIDs []uuid.UUID `json:"amenity_ids,omitempty"`
}

// AddressPayloadDTO contiene la información necesaria para crear una dirección
// (copiado del person service para mantener consistencia)
type AddressPayloadDTO struct {
	Floor   *string `json:"floor,omitempty"`
	Unit    *string `json:"unit,omitempty"`
	Street  string  `json:"street"  validate:"required"`
	Number  int     `json:"number"  validate:"required"`
	City    string  `json:"city"    validate:"required"`
	State   string  `json:"state,omitempty"`
	Zip     string  `json:"zip,omitempty"`
	Country string  `json:"country" validate:"required,len=2"`
}

type UpdatePropertyDTO struct {
	PropertyTypeID *int     `json:"property_type_id,omitempty"  validate:"omitempty,min=1"`
	InternalCode   *string  `json:"internal_code,omitempty"`
	YearBuilt      *int     `json:"year_built,omitempty"        validate:"omitempty,min=1800,max=2100"`
	Bedrooms       *int     `json:"bedrooms,omitempty"          validate:"omitempty,min=0,max=50"`
	Bathrooms      *float32 `json:"bathrooms,omitempty"         validate:"omitempty,min=0,max=20"`
	TotalAreaSqm   *float64 `json:"total_area_sqm,omitempty"    validate:"omitempty,min=0"`
	CoveredAreaSqm *float64 `json:"covered_area_sqm,omitempty"  validate:"omitempty,min=0"`
	Description    *string  `json:"description,omitempty"`
}

type PropertyResponseDTO struct {
	ID               uuid.UUID `json:"id"`
	OwnerPersonID    uuid.UUID `json:"owner_person_id"`
	AddressID        uuid.UUID `json:"address_id"`
	PropertyTypeID   int       `json:"property_type_id"`
	PropertyTypeName *string   `json:"property_type_name,omitempty"` // For convenience in API responses
	InternalCode     *string   `json:"internal_code,omitempty"`
	YearBuilt        *int      `json:"year_built,omitempty"`
	Bedrooms         *int      `json:"bedrooms,omitempty"`
	Bathrooms        *float32  `json:"bathrooms,omitempty"`
	TotalAreaSqm     *float64  `json:"total_area_sqm,omitempty"`
	CoveredAreaSqm   *float64  `json:"covered_area_sqm,omitempty"`
	Description      *string   `json:"description,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`

	// Relations
	Management []PropertyManagementResponseDTO `json:"management,omitempty"`
	Amenities  []AmenityResponseDTO            `json:"amenities,omitempty"`
}

/* ─────────────────────────  PROPERTY MANAGEMENT  ───────────────────────── */

type CreatePropertyManagementDTO struct {
	PropertyID           uuid.UUID  `json:"property_id"         validate:"required"`
	ManagerID            uuid.UUID  `json:"manager_id"          validate:"required"`
	ManagerTypeID        int        `json:"manager_type_id"     validate:"required,min=1"`
	StartDate            *time.Time `json:"start_date,omitempty"`
	EndDate              *time.Time `json:"end_date,omitempty"`
	CommissionPercentage *float64   `json:"commission_percentage,omitempty" validate:"omitempty,min=0,max=100"`
}

type UpdatePropertyManagementDTO struct {
	EndDate              *time.Time `json:"end_date,omitempty"`
	CommissionPercentage *float64   `json:"commission_percentage,omitempty" validate:"omitempty,min=0,max=100"`
}

type PropertyManagementResponseDTO struct {
	ID                   uuid.UUID  `json:"id"`
	PropertyID           uuid.UUID  `json:"property_id"`
	ManagerID            uuid.UUID  `json:"manager_id"`
	ManagerTypeID        int        `json:"manager_type_id"`
	ManagerTypeName      *string    `json:"manager_type_name,omitempty"` // For convenience
	StartDate            time.Time  `json:"start_date"`
	EndDate              *time.Time `json:"end_date,omitempty"`
	CommissionPercentage *float64   `json:"commission_percentage,omitempty"`

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

/* ─────────────────────────  PROPERTY TYPE  ───────────────────────── */

type PropertyTypeResponseDTO struct {
	ID          int        `json:"id"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

/* ─────────────────────────  MANAGER TYPE  ───────────────────────── */

type ManagerTypeResponseDTO struct {
	ID          int        `json:"id"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

/* ─────────────────────────  CONVERTERS  ───────────────────────── */

func ToPropertyResponse(property *models.Property) PropertyResponseDTO {
	response := PropertyResponseDTO{
		ID:             property.ID,
		OwnerPersonID:  property.OwnerPersonID,
		AddressID:      property.AddressID,
		PropertyTypeID: property.PropertyTypeID,
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

	// Add property type name if available
	if property.PropertyType.Name != "" {
		response.PropertyTypeName = &property.PropertyType.Name
	}

	// Convert management
	for _, mgmt := range property.PropertyManagement {
		mgmtResponse := PropertyManagementResponseDTO{
			ID:                   mgmt.ID,
			PropertyID:           mgmt.PropertyID,
			ManagerID:            mgmt.ManagerID,
			ManagerTypeID:        mgmt.ManagerTypeID,
			StartDate:            mgmt.StartDate,
			EndDate:              mgmt.EndDate,
			CommissionPercentage: mgmt.CommissionPercentage,
			CreatedAt:            mgmt.CreatedAt,
			UpdatedAt:            mgmt.UpdatedAt,
		}

		// Add manager type name if available
		if mgmt.ManagerType.Name != "" {
			mgmtResponse.ManagerTypeName = &mgmt.ManagerType.Name
		}

		response.Management = append(response.Management, mgmtResponse)
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

func ToPropertyTypeResponse(propertyType *models.PropertyType) PropertyTypeResponseDTO {
	return PropertyTypeResponseDTO{
		ID:          propertyType.ID,
		Code:        propertyType.Code,
		Name:        propertyType.Name,
		Description: propertyType.Description,
		IsActive:    propertyType.IsActive,
		CreatedAt:   propertyType.CreatedAt,
		UpdatedAt:   propertyType.UpdatedAt,
	}
}

func ToManagerTypeResponse(managerType *models.ManagerType) ManagerTypeResponseDTO {
	return ManagerTypeResponseDTO{
		ID:          managerType.ID,
		Code:        managerType.Code,
		Name:        managerType.Name,
		Description: managerType.Description,
		IsActive:    managerType.IsActive,
		CreatedAt:   managerType.CreatedAt,
		UpdatedAt:   managerType.UpdatedAt,
	}
}
