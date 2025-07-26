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
	List(category *string, isActive *bool, limit, offset int) ([]models.PropertyType, int64, error)
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
	List(category *models.AmenityCategory, isActive *bool, limit, offset int) ([]models.Amenity, int64, error)
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
	GetByPropertyAndAmenity(propertyID uuid.UUID, amenityID int32) (*models.PropertyAmenity, error)
	GetByPropertyID(propertyID uuid.UUID) ([]models.PropertyAmenity, error)
	ListByPropertyID(propertyID uuid.UUID) ([]*models.Amenity, error)
	AddAmenityToProperty(propertyID uuid.UUID, amenityID int32, note string) error
	RemoveAmenityFromProperty(propertyID uuid.UUID, amenityID int32) error
	Update(amenity *models.PropertyAmenity) error
	Delete(propertyID uuid.UUID, amenityID int32) error
	BulkCreateForProperty(propertyID uuid.UUID, amenities []models.PropertyAmenity) error
	BulkDeleteForProperty(propertyID uuid.UUID) error
}

// PropertyListingRepository interface
type PropertyListingRepository interface {
	Create(listing *models.PropertyListing) (*models.PropertyListing, error)
	GetByID(id uuid.UUID) (*models.PropertyListing, error)
	GetByPropertyID(propertyID uuid.UUID) ([]models.PropertyListing, error)
	List(filters PropertyListingFilters, limit, offset int) ([]models.PropertyListing, int64, error)
	Update(listing *models.PropertyListing) error
	UpdateStatus(id uuid.UUID, status models.ListingStatus, updatedBy uuid.UUID) error
	SoftDelete(id uuid.UUID, deletedBy uuid.UUID) error
	Delete(id uuid.UUID) error
	GetActiveListingsByType(operationType models.OperationType, limit, offset int) ([]models.PropertyListing, int64, error)
	Search(query string, filters PropertyListingFilters, limit, offset int) ([]models.PropertyListing, int64, error)
	IncrementViews(id uuid.UUID) error
	IncrementInquiries(id uuid.UUID) error
	IncrementFavorites(id uuid.UUID) error
}

// PropertyMediaRepository interface
type PropertyMediaRepository interface {
	Create(media *models.PropertyMedia) (*models.PropertyMedia, error)
	GetByID(id uuid.UUID) (*models.PropertyMedia, error)
	GetByPropertyID(propertyID uuid.UUID) ([]models.PropertyMedia, error)
	GetByPropertyIDAndType(propertyID uuid.UUID, mediaType models.MediaType) ([]models.PropertyMedia, error)
	GetMainImage(propertyID uuid.UUID) (*models.PropertyMedia, error)
	List(filters PropertyMediaFilters, limit, offset int) ([]models.PropertyMedia, int64, error)
	Update(media *models.PropertyMedia) error
	SetAsMain(propertyID uuid.UUID, mediaID uuid.UUID, updatedBy uuid.UUID) error
	UpdateDisplayOrder(propertyID uuid.UUID, mediaOrder []uuid.UUID, updatedBy uuid.UUID) error
	SoftDelete(id uuid.UUID, deletedBy uuid.UUID) error
	Delete(id uuid.UUID) error
	BulkCreateForProperty(propertyID uuid.UUID, media []models.PropertyMedia) error
	BulkDeleteForProperty(propertyID uuid.UUID) error
}

// PropertyValuationRepository interface
type PropertyValuationRepository interface {
	Create(valuation *models.PropertyValuation) (*models.PropertyValuation, error)
	GetByID(id uuid.UUID) (*models.PropertyValuation, error)
	GetByPropertyID(propertyID uuid.UUID) ([]models.PropertyValuation, error)
	GetLatestByPropertyID(propertyID uuid.UUID) (*models.PropertyValuation, error)
	GetByPropertyIDAndType(propertyID uuid.UUID, valuationType string) ([]models.PropertyValuation, error)
	List(filters PropertyValuationFilters, limit, offset int) ([]models.PropertyValuation, int64, error)
	Update(valuation *models.PropertyValuation) error
	Delete(id uuid.UUID) error
}

// Filter structs
type PropertyFilters struct {
	OwnerPersonID   *uuid.UUID
	PropertyTypeID  *int32
	InternalCode    *string
	YearBuilt       *int32
	Bedrooms        *int32
	Bathrooms       *float32
	MinTotalArea    *float64
	MaxTotalArea    *float64
	MinPrice        *float64
	MaxPrice        *float64
	City            *string
	State           *string
	ConditionRating *int32
	IncludeDeleted  bool
}

type PropertyManagementFilters struct {
	PropertyID     *uuid.UUID
	ManagerID      *uuid.UUID
	ManagerTypeID  *int32
	OrganizationID *uuid.UUID
	ActiveOnly     bool
	IncludeDeleted bool
}

type PropertyListingFilters struct {
	PropertyID       *uuid.UUID
	OperationType    *models.OperationType
	ListingStatus    *models.ListingStatus
	MinPrice         *float64
	MaxPrice         *float64
	Currency         *string
	PublishOnWeb     *bool
	PublishOnPortals *bool
	Highlighted      *bool
	CreatedBy        *uuid.UUID
	IncludeDeleted   bool
}

type PropertyMediaFilters struct {
	PropertyID     *uuid.UUID
	MediaType      *models.MediaType
	IsMain         *bool
	IsPublic       *bool
	CreatedBy      *uuid.UUID
	IncludeDeleted bool
}

type PropertyValuationFilters struct {
	PropertyID      *uuid.UUID
	ValuationType   *string
	ValuationMethod *string
	MinValue        *float64
	MaxValue        *float64
	Currency        *string
	ValidOnly       bool // Only non-expired valuations
	CreatedBy       *uuid.UUID
}

// BulkError for bulk operations
type BulkError struct {
	Index int
	Error string
}
