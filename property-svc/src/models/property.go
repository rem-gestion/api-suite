package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/rem-gestion/rem-common/db"
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

// ListingStatus enum type for property listing status
type ListingStatus string

const (
	ListingStatusDraft     ListingStatus = "draft"
	ListingStatusActive    ListingStatus = "active"
	ListingStatusPaused    ListingStatus = "paused"
	ListingStatusReserved  ListingStatus = "reserved"
	ListingStatusInProcess ListingStatus = "in_process"
	ListingStatusCompleted ListingStatus = "completed"
	ListingStatusExpired   ListingStatus = "expired"
	ListingStatusCancelled ListingStatus = "cancelled"
)

// OperationType enum type for operation types
type OperationType string

const (
	OperationTypeSale       OperationType = "sale"
	OperationTypeRental     OperationType = "rental"
	OperationTypeTemporary  OperationType = "temporary"
	OperationTypeCommercial OperationType = "commercial"
	OperationTypeInvestment OperationType = "investment"
)

// MediaType enum type for property media
type MediaType string

const (
	MediaTypePhoto       MediaType = "photo"
	MediaTypeVideo       MediaType = "video"
	MediaTypeFloorPlan   MediaType = "floor_plan"
	MediaTypeDocument    MediaType = "document"
	MediaTypeVirtualTour MediaType = "virtual_tour"
	MediaTypeOther       MediaType = "other"
)

// PropertyType model - catalog of property types
type PropertyType struct {
	ID          int32      `gorm:"primaryKey;autoIncrement" json:"id"`
	Code        string     `gorm:"size:32;unique;not null" json:"code"`
	Name        string     `gorm:"size:100;not null" json:"name"`
	Description *string    `gorm:"type:text" json:"description"`
	Category    *string    `gorm:"size:50" json:"category"` // residential, commercial, industrial
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`

	// Relations
	Properties []Property `gorm:"foreignKey:PropertyTypeID" json:"properties,omitempty"`
}

func (PropertyType) TableName() string { return db.GetPropertyTableName("property_type") }

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

func (ManagerType) TableName() string { return db.GetPropertyTableName("manager_type") }

// Amenity model - catalog of amenities
type Amenity struct {
	ID        int32            `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string           `gorm:"size:100;unique;not null" json:"name"`
	Category  *AmenityCategory `gorm:"type:amenity_category" json:"category"`
	IconURL   *string          `gorm:"type:text" json:"icon_url"`
	IsActive  bool             `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time        `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt *time.Time       `json:"updated_at"`

	// Relations
	PropertyAmenities []PropertyAmenity `gorm:"foreignKey:AmenityID" json:"property_amenities,omitempty"`
}

func (Amenity) TableName() string { return db.GetPropertyTableName("amenity") }

// Property model - main property entity
type Property struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	OwnerPersonID  uuid.UUID `gorm:"type:uuid;not null" json:"owner_person_id"`
	AddressID      uuid.UUID `gorm:"type:uuid;unique;not null" json:"address_id"`
	PropertyTypeID int32     `gorm:"not null" json:"property_type_id"`
	InternalCode   *string   `gorm:"size:50;unique" json:"internal_code"`

	// Physical characteristics
	YearBuilt *int32   `json:"year_built"`
	Bedrooms  *int32   `json:"bedrooms"`
	Bathrooms *float32 `gorm:"type:decimal(3,1)" json:"bathrooms"`
	Toilets   *int32   `json:"toilets"`
	Rooms     *int32   `json:"rooms"`
	Levels    *int32   `json:"levels"`

	// Areas
	TotalAreaSqm       *float64 `gorm:"type:decimal(10,2)" json:"total_area_sqm"`
	CoveredAreaSqm     *float64 `gorm:"type:decimal(10,2)" json:"covered_area_sqm"`
	SemicoveredAreaSqm *float64 `gorm:"type:decimal(10,2)" json:"semicovered_area_sqm"`
	UncoveredAreaSqm   *float64 `gorm:"type:decimal(10,2)" json:"uncovered_area_sqm"`

	// Additional features
	IncludesGarage *bool  `gorm:"default:false" json:"includes_garage"`
	GarageSpaces   *int32 `gorm:"default:0" json:"garage_spaces"`
	Furnished      *bool  `gorm:"default:false" json:"furnished"`

	// Condition
	ConditionRating    *int32 `json:"condition_rating"` // 1-5 scale
	LastRenovationYear *int32 `json:"last_renovation_year"`

	// General info
	Description *string `gorm:"type:text" json:"description"`

	// Audit fields
	CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt *time.Time     `json:"updated_at"`
	UpdatedBy *uuid.UUID     `gorm:"type:uuid" json:"updated_by"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// Relations
	PropertyType        PropertyType         `gorm:"foreignKey:PropertyTypeID" json:"property_type,omitempty"`
	PropertyManagements []PropertyManagement `gorm:"foreignKey:PropertyID" json:"property_managements,omitempty"`
	PropertyAmenities   []PropertyAmenity    `gorm:"foreignKey:PropertyID" json:"property_amenities,omitempty"`
	PropertyListings    []PropertyListing    `gorm:"foreignKey:PropertyID" json:"property_listings,omitempty"`
	PropertyMedia       []PropertyMedia      `gorm:"foreignKey:PropertyID" json:"property_media,omitempty"`
	PropertyValuations  []PropertyValuation  `gorm:"foreignKey:PropertyID" json:"property_valuations,omitempty"`
}

func (Property) TableName() string { return db.GetPropertyTableName("property") }

// PropertyManagement model - manages who handles the property
type PropertyManagement struct {
	ID                   uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	PropertyID           uuid.UUID      `gorm:"type:uuid;not null" json:"property_id"`
	ManagerID            uuid.UUID      `gorm:"type:uuid;not null" json:"manager_id"`
	ManagerTypeID        int32          `gorm:"not null" json:"manager_type_id"`
	OrganizationID       *uuid.UUID     `gorm:"type:uuid" json:"organization_id"`
	StartDate            time.Time      `gorm:"type:date;not null" json:"start_date"`
	EndDate              *time.Time     `gorm:"type:date" json:"end_date"`
	IsActive             *bool          `gorm:"default:true" json:"is_active"`
	CommissionPercentage *float64       `gorm:"type:decimal(5,2)" json:"commission_percentage"`
	FixedFee             *float64       `gorm:"type:decimal(10,2)" json:"fixed_fee"`
	ExclusiveManagement  *bool          `gorm:"default:false" json:"exclusive_management"`
	ServicesIncluded     *string        `gorm:"type:jsonb" json:"services_included"`
	CreatedAt            time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt            *time.Time     `json:"updated_at"`
	UpdatedBy            *uuid.UUID     `gorm:"type:uuid" json:"updated_by"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// Relations
	Property    Property    `gorm:"foreignKey:PropertyID" json:"property,omitempty"`
	ManagerType ManagerType `gorm:"foreignKey:ManagerTypeID" json:"manager_type,omitempty"`
}

func (PropertyManagement) TableName() string { return db.GetPropertyTableName("property_management") }

// PropertyAmenity model - junction table for property amenities
type PropertyAmenity struct {
	PropertyID uuid.UUID `gorm:"type:uuid;not null;primaryKey" json:"property_id"`
	AmenityID  int32     `gorm:"not null;primaryKey" json:"amenity_id"`
	Note       *string   `gorm:"type:text" json:"note"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relations
	Property Property `gorm:"foreignKey:PropertyID" json:"property,omitempty"`
	Amenity  Amenity  `gorm:"foreignKey:AmenityID" json:"amenity,omitempty"`
}

func (PropertyAmenity) TableName() string { return db.GetPropertyTableName("property_amenities") }

// PropertyListing model - property operations/listings
type PropertyListing struct {
	ID            uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	PropertyID    uuid.UUID     `gorm:"type:uuid;not null" json:"property_id"`
	OperationType OperationType `gorm:"type:operation_type;not null" json:"operation_type"`
	ListingStatus ListingStatus `gorm:"type:listing_status;default:'draft'" json:"listing_status"`

	// Pricing and conditions
	Price            *float64 `gorm:"type:decimal(15,2)" json:"price"`
	Currency         string   `gorm:"size:3;default:'ARS'" json:"currency"`
	PricePerSqm      *float64 `gorm:"type:decimal(10,2)" json:"price_per_sqm"`
	IncludesExpenses *bool    `gorm:"default:false" json:"includes_expenses"`
	MonthlyExpenses  *float64 `gorm:"type:decimal(10,2)" json:"monthly_expenses"`
	DepositAmount    *float64 `gorm:"type:decimal(10,2)" json:"deposit_amount"`

	// Terms
	RentalTerms *string `gorm:"type:jsonb" json:"rental_terms"`
	SaleTerms   *string `gorm:"type:jsonb" json:"sale_terms"`

	// Marketing
	Title            string     `gorm:"size:200;not null" json:"title"`
	Description      *string    `gorm:"type:text" json:"description"`
	Highlighted      *bool      `gorm:"default:false" json:"highlighted"`
	PublicationStart *time.Time `gorm:"type:date" json:"publication_start"`
	PublicationEnd   *time.Time `gorm:"type:date" json:"publication_end"`
	AutoRenew        *bool      `gorm:"default:false" json:"auto_renew"`

	// Publication config
	PublishOnWeb     *bool   `gorm:"default:true" json:"publish_on_web"`
	PublishOnPortals *bool   `gorm:"default:false" json:"publish_on_portals"`
	PortalConfig     *string `gorm:"type:jsonb" json:"portal_config"`

	// Metrics
	ViewsCount     int32 `gorm:"default:0" json:"views_count"`
	InquiriesCount int32 `gorm:"default:0" json:"inquiries_count"`
	FavoritesCount int32 `gorm:"default:0" json:"favorites_count"`

	// Audit
	CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	CreatedBy *uuid.UUID     `gorm:"type:uuid" json:"created_by"`
	UpdatedAt *time.Time     `json:"updated_at"`
	UpdatedBy *uuid.UUID     `gorm:"type:uuid" json:"updated_by"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// Relations
	Property Property `gorm:"foreignKey:PropertyID" json:"property,omitempty"`
}

func (PropertyListing) TableName() string { return db.GetPropertyTableName("listings") }

// PropertyMedia model - multimedia content for properties
type PropertyMedia struct {
	ID            uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	PropertyID    uuid.UUID      `gorm:"type:uuid;not null" json:"property_id"`
	MediaType     MediaType      `gorm:"type:media_type;not null" json:"media_type"`
	FileName      string         `gorm:"size:255;not null" json:"file_name"`
	FileURL       string         `gorm:"type:text;not null" json:"file_url"`
	FileSizeBytes *int64         `json:"file_size_bytes"`
	MimeType      *string        `gorm:"size:100" json:"mime_type"`
	Title         *string        `gorm:"size:200" json:"title"`
	Description   *string        `gorm:"type:text" json:"description"`
	AltText       *string        `gorm:"size:255" json:"alt_text"`
	DisplayOrder  int32          `gorm:"default:0" json:"display_order"`
	IsMain        *bool          `gorm:"default:false" json:"is_main"`
	IsPublic      *bool          `gorm:"default:true" json:"is_public"`
	WidthPx       *int32         `json:"width_px"`
	HeightPx      *int32         `json:"height_px"`
	CreatedAt     time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	CreatedBy     *uuid.UUID     `gorm:"type:uuid" json:"created_by"`
	UpdatedAt     *time.Time     `json:"updated_at"`
	UpdatedBy     *uuid.UUID     `gorm:"type:uuid" json:"updated_by"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// Relations
	Property Property `gorm:"foreignKey:PropertyID" json:"property,omitempty"`
}

func (PropertyMedia) TableName() string { return db.GetPropertyTableName("media") }

// PropertyValuation model - property valuations and appraisals
type PropertyValuation struct {
	ID                    uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	PropertyID            uuid.UUID  `gorm:"type:uuid;not null" json:"property_id"`
	ValuationDate         time.Time  `gorm:"type:date;not null" json:"valuation_date"`
	ValuationType         string     `gorm:"size:50;not null" json:"valuation_type"` // market, fiscal, insurance, mortgage
	AppraisedValue        float64    `gorm:"type:decimal(15,2);not null" json:"appraised_value"`
	Currency              string     `gorm:"size:3;default:'ARS'" json:"currency"`
	ValuePerSqm           *float64   `gorm:"type:decimal(10,2)" json:"value_per_sqm"`
	AppraiserName         *string    `gorm:"size:100" json:"appraiser_name"`
	AppraiserLicense      *string    `gorm:"size:50" json:"appraiser_license"`
	AppraiserOrganization *string    `gorm:"size:100" json:"appraiser_organization"`
	ValuationMethod       *string    `gorm:"size:50" json:"valuation_method"` // comparative, income, cost
	MarketConditions      *string    `gorm:"type:text" json:"market_conditions"`
	AdjustmentsApplied    *string    `gorm:"type:jsonb" json:"adjustments_applied"`
	ComparableProperties  *string    `gorm:"type:jsonb" json:"comparable_properties"`
	ValuationReportURL    *string    `gorm:"type:text" json:"valuation_report_url"`
	PhotosURLs            *string    `gorm:"type:jsonb" json:"photos_urls"`
	ValidUntil            *time.Time `gorm:"type:date" json:"valid_until"`
	Purpose               *string    `gorm:"size:100" json:"purpose"`
	CreatedAt             time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	CreatedBy             *uuid.UUID `gorm:"type:uuid" json:"created_by"`
	UpdatedAt             *time.Time `json:"updated_at"`
	UpdatedBy             *uuid.UUID `gorm:"type:uuid" json:"updated_by"`

	// Relations
	Property Property `gorm:"foreignKey:PropertyID" json:"property,omitempty"`
}

func (PropertyValuation) TableName() string { return db.GetPropertyTableName("valuations") }
