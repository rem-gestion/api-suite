package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AmenityCategory enum type
type AmenityCategory string

const (
	AmenityCategoryGeneral      AmenityCategory = "general"
	AmenityCategoryServices     AmenityCategory = "services"
	AmenityCategoryEnvironments AmenityCategory = "environments"
	AmenityCategorySecurity     AmenityCategory = "security"
	AmenityCategoryComfort      AmenityCategory = "comfort"
)

// PropertyType model - catalog of property types
type PropertyType struct {
	ID          int32      `gorm:"primaryKey;autoIncrement" json:"id"`
	Code        string     `gorm:"size:32;unique;not null" json:"code"`
	Name        string     `gorm:"size:100;not null" json:"name"`
	Description *string    `gorm:"type:text" json:"description"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`

	// Relations
	Properties []Property `gorm:"foreignKey:PropertyTypeID" json:"properties,omitempty"`
}

func (PropertyType) TableName() string { return "property_property_type" }

// ManagerType model - catalog of manager types
type ManagerType struct {
	ID          int32      `gorm:"primaryKey;autoIncrement" json:"id"`
	Code        string     `gorm:"size:32;unique;not null" json:"code"`
	Name        string     `gorm:"size:100;not null" json:"name"`
	Description *string    `gorm:"type:text" json:"description"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`

	// Relations
	PropertyManagements []PropertyManagement `gorm:"foreignKey:ManagerTypeID" json:"property_managements,omitempty"`
}

func (ManagerType) TableName() string { return "property_manager_type" }

// Amenity model - catalog of amenities
type Amenity struct {
	ID       int32            `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string           `gorm:"size:100;unique;not null" json:"name"`
	Category *AmenityCategory `gorm:"type:amenity_category" json:"category"`
	IconURL  *string          `gorm:"type:text" json:"icon_url"`

	// Relations
	PropertyAmenities []PropertyAmenity `gorm:"foreignKey:AmenityID" json:"property_amenities,omitempty"`
}

func (Amenity) TableName() string { return "property_amenity" }

// Property model - main property entity
type Property struct {
	ID             uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	OwnerPersonID  uuid.UUID      `gorm:"type:uuid;not null" json:"owner_person_id"`
	AddressID      uuid.UUID      `gorm:"type:uuid;unique;not null" json:"address_id"`
	PropertyTypeID int32          `gorm:"not null" json:"property_type_id"`
	InternalCode   *string        `gorm:"size:50" json:"internal_code"`
	YearBuilt      *int32         `json:"year_built"`
	Bedrooms       *int32         `json:"bedrooms"`
	Bathrooms      *float32       `gorm:"type:decimal(3,1)" json:"bathrooms"`
	TotalAreaSqm   *float64       `gorm:"type:decimal(10,2)" json:"total_area_sqm"`
	CoveredAreaSqm *float64       `gorm:"type:decimal(10,2)" json:"covered_area_sqm"`
	Description    *string        `gorm:"type:text" json:"description"`
	CreatedAt      time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      *time.Time     `json:"updated_at"`
	UpdatedBy      *uuid.UUID     `gorm:"type:uuid" json:"updated_by"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// Relations
	PropertyType        PropertyType         `gorm:"foreignKey:PropertyTypeID" json:"property_type,omitempty"`
	PropertyManagements []PropertyManagement `gorm:"foreignKey:PropertyID" json:"property_managements,omitempty"`
	PropertyAmenities   []PropertyAmenity    `gorm:"foreignKey:PropertyID" json:"property_amenities,omitempty"`
}

func (Property) TableName() string { return "property_property" }

// PropertyManagement model - manages who handles the property
type PropertyManagement struct {
	ID                   uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	PropertyID           uuid.UUID      `gorm:"type:uuid;not null" json:"property_id"`
	ManagerID            uuid.UUID      `gorm:"type:uuid;not null" json:"manager_id"`
	ManagerTypeID        int32          `gorm:"not null" json:"manager_type_id"`
	StartDate            time.Time      `gorm:"type:date;not null" json:"start_date"`
	EndDate              *time.Time     `gorm:"type:date" json:"end_date"`
	CommissionPercentage *float64       `gorm:"type:decimal(5,2)" json:"commission_percentage"`
	CreatedAt            time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt            *time.Time     `json:"updated_at"`
	UpdatedBy            *uuid.UUID     `gorm:"type:uuid" json:"updated_by"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// Relations
	Property    Property    `gorm:"foreignKey:PropertyID" json:"property,omitempty"`
	ManagerType ManagerType `gorm:"foreignKey:ManagerTypeID" json:"manager_type,omitempty"`
}

func (PropertyManagement) TableName() string { return "property_property_management" }

// PropertyAmenity model - junction table for property amenities
type PropertyAmenity struct {
	PropertyID uuid.UUID `gorm:"type:uuid;not null;primaryKey" json:"property_id"`
	AmenityID  int32     `gorm:"not null;primaryKey" json:"amenity_id"`
	Note       *string   `gorm:"type:text" json:"note"`

	// Relations
	Property Property `gorm:"foreignKey:PropertyID" json:"property,omitempty"`
	Amenity  Amenity  `gorm:"foreignKey:AmenityID" json:"amenity,omitempty"`
}

func (PropertyAmenity) TableName() string { return "property_property_amenities" }

// Custom composite primary key for PropertyAmenity
func (PropertyAmenity) BeforeCreate(tx *gorm.DB) error {
	// Any validation logic can go here
	return nil
}
