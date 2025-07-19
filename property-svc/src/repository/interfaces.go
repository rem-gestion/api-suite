package repository

import (
	"github.com/google/uuid"
	"github.com/rem-gestion/api-suite/property/src/models"
)

// PropertyTypeRepository interface
type PropertyTypeRepository interface {
	Create(propertyType *models.PropertyType) (*models.PropertyType, error)
	GetByID(id int32) (*models.PropertyType, error)
	GetByCode(code string) (*models.PropertyType, error)
	List(isActive *bool, limit, offset int) ([]models.PropertyType, int64, error)
	Update(propertyType *models.PropertyType) error
	Delete(id int32) error
}

// ManagerTypeRepository interface
type ManagerTypeRepository interface {
	Create(managerType *models.ManagerType) (*models.ManagerType, error)
	GetByID(id int32) (*models.ManagerType, error)
	GetByCode(code string) (*models.ManagerType, error)
	List(isActive *bool, limit, offset int) ([]models.ManagerType, int64, error)
	Update(managerType *models.ManagerType) error
	Delete(id int32) error
}

// AmenityRepository interface
type AmenityRepository interface {
	Create(amenity *models.Amenity) (*models.Amenity, error)
	GetByID(id int32) (*models.Amenity, error)
	GetByName(name string) (*models.Amenity, error)
	List(category *models.AmenityCategory, limit, offset int) ([]models.Amenity, int64, error)
	Update(amenity *models.Amenity) error
	Delete(id int32) error
}

// PropertyRepository interface
type PropertyRepository interface {
	Create(property *models.Property) (*models.Property, error)
	GetByID(id uuid.UUID) (*models.Property, error)
	GetByInternalCode(code string) (*models.Property, error)
	List(filters PropertyFilters, limit, offset int) ([]models.Property, int64, error)
	Update(property *models.Property) error
	SoftDelete(id uuid.UUID, deletedBy uuid.UUID) error
	Delete(id uuid.UUID) error
	Search(query string, filters PropertyFilters, limit, offset int) ([]models.Property, int64, error)
	BulkCreate(properties []models.Property) ([]models.Property, []BulkError, error)
}

// PropertyManagementRepository interface
type PropertyManagementRepository interface {
	Create(management *models.PropertyManagement) (*models.PropertyManagement, error)
	GetByID(id uuid.UUID) (*models.PropertyManagement, error)
	GetByPropertyID(propertyID uuid.UUID, active bool) ([]models.PropertyManagement, error)
	GetByManagerID(managerID uuid.UUID, active bool) ([]models.PropertyManagement, error)
	List(filters PropertyManagementFilters, limit, offset int) ([]models.PropertyManagement, int64, error)
	Update(management *models.PropertyManagement) error
	SoftDelete(id uuid.UUID, deletedBy uuid.UUID) error
	Delete(id uuid.UUID) error
}

// PropertyAmenityRepository interface
type PropertyAmenityRepository interface {
	Create(amenity *models.PropertyAmenity) (*models.PropertyAmenity, error)
	GetByPropertyAndAmenity(propertyID uuid.UUID, amenityID string) (*models.PropertyAmenity, error)
	GetByPropertyID(propertyID uuid.UUID) ([]models.PropertyAmenity, error)
	Update(amenity *models.PropertyAmenity) error
	Delete(propertyID uuid.UUID, amenityID string) error
	BulkCreateForProperty(propertyID uuid.UUID, amenities []models.PropertyAmenity) error
	BulkDeleteForProperty(propertyID uuid.UUID) error
}

// Filter structs
type PropertyFilters struct {
	OwnerPersonID  *uuid.UUID
	PropertyTypeID *int32
	InternalCode   *string
	YearBuilt      *int32
	Bedrooms       *int32
	Bathrooms      *float32
	MinTotalArea   *float64
	MaxTotalArea   *float64
	City           *string
	State          *string
	IncludeDeleted bool
}

type PropertyManagementFilters struct {
	PropertyID     *uuid.UUID
	ManagerID      *uuid.UUID
	ManagerTypeID  *int32
	ActiveOnly     bool
	IncludeDeleted bool
}

// BulkError for bulk operations
type BulkError struct {
	Index int
	Error string
}
