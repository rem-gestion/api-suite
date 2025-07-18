package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ---- enums ----
type PropertyType string

const (
	PropertyApartment           PropertyType = "APARTMENT"
	PropertyHouse               PropertyType = "HOUSE"
	PropertyCommercialSpace     PropertyType = "COMMERCIAL_SPACE"
	PropertyOffice              PropertyType = "OFFICE"
	PropertyLand                PropertyType = "LAND"
	PropertyIndustrialWarehouse PropertyType = "INDUSTRIAL_WAREHOUSE"
)

// ---- core tables ----
type Property struct {
	ID             uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	OwnerPersonID  uuid.UUID    `gorm:"type:uuid;not null;index"`
	AddressID      uuid.UUID    `gorm:"type:uuid;not null;unique"`
	PropertyType   PropertyType `gorm:"type:property_type;not null;index"`
	InternalCode   *string      `gorm:"size:50;unique"`
	YearBuilt      *int
	Bedrooms       *int
	Bathrooms      *float32 `gorm:"type:decimal(3,1)"`
	TotalAreaSqm   *float64 `gorm:"type:decimal(10,2)"`
	CoveredAreaSqm *float64 `gorm:"type:decimal(10,2)"`
	Description    *string  `gorm:"type:text"`

	CreatedAt time.Time
	UpdatedAt *time.Time
	UpdatedBy *uuid.UUID
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Relations
	PropertyManagement []PropertyManagement `gorm:"foreignKey:PropertyID"`
	PropertyAmenities  []PropertyAmenity    `gorm:"foreignKey:PropertyID"`
}

type PropertyManagement struct {
	ID                uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	PropertyID        uuid.UUID `gorm:"type:uuid;not null;index"`
	OrganizationID    uuid.UUID `gorm:"type:uuid;not null;index"`
	ManagedSince      time.Time `gorm:"not null;default:now()"`
	ManagedUntil      *time.Time
	IsActive          bool     `gorm:"not null;default:true;index"`
	CommissionPercent *float64 `gorm:"type:decimal(5,2)"`
	Notes             *string  `gorm:"type:text"`

	CreatedAt time.Time
	UpdatedAt *time.Time
	UpdatedBy *uuid.UUID
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Relations
	Property Property `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Amenity struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string    `gorm:"size:100;not null;unique"`
	Description *string   `gorm:"type:text"`
	Icon        *string   `gorm:"size:50"`
	Category    *string   `gorm:"size:50"`

	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Relations
	PropertyAmenities []PropertyAmenity `gorm:"foreignKey:AmenityID"`
}

type PropertyAmenity struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	PropertyID uuid.UUID `gorm:"type:uuid;not null;index"`
	AmenityID  uuid.UUID `gorm:"type:uuid;not null;index"`

	CreatedAt time.Time

	// Relations
	Property Property `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Amenity  Amenity  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// Table names
func (Property) TableName() string           { return "property" }
func (PropertyManagement) TableName() string { return "property_management" }
func (Amenity) TableName() string            { return "amenity" }
func (PropertyAmenity) TableName() string    { return "property_amenities" }
