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
	Category    *string `json:"category" binding:"omitempty,oneof=residential commercial industrial"`
	IsActive    *bool   `json:"is_active"`
}

type PropertyTypeUpdate struct {
	Code        *string `json:"code,omitempty" binding:"omitempty,max=32"`
	Name        *string `json:"name,omitempty" binding:"omitempty,max=100"`
	Description *string `json:"description,omitempty"`
	Category    *string `json:"category,omitempty" binding:"omitempty,oneof=residential commercial industrial"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

type PropertyTypeResponse struct {
	ID          int32      `json:"id"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	Category    *string    `json:"category"`
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
	IsActive *bool                   `json:"is_active"`
}

type AmenityUpdate struct {
	Name     *string                 `json:"name,omitempty" binding:"omitempty,max=100"`
	Category *models.AmenityCategory `json:"category,omitempty"`
	IconURL  *string                 `json:"icon_url,omitempty"`
	IsActive *bool                   `json:"is_active,omitempty"`
}

type AmenityResponse struct {
	ID        int32                   `json:"id"`
	Name      string                  `json:"name"`
	Category  *models.AmenityCategory `json:"category"`
	IconURL   *string                 `json:"icon_url"`
	IsActive  bool                    `json:"is_active"`
	CreatedAt time.Time               `json:"created_at"`
	UpdatedAt *time.Time              `json:"updated_at"`
}

// Property DTOs
type PropertyCreate struct {
	OwnerPersonID  uuid.UUID `json:"owner_person_id" binding:"required"`
	AddressID      uuid.UUID `json:"address_id" binding:"required"`
	PropertyTypeID int32     `json:"property_type_id" binding:"required"`
	InternalCode   *string   `json:"internal_code,omitempty" binding:"omitempty,max=50"`

	// Physical characteristics
	YearBuilt *int32   `json:"year_built,omitempty" binding:"omitempty,min=1800,max=2100"`
	Bedrooms  *int32   `json:"bedrooms,omitempty" binding:"omitempty,min=0"`
	Bathrooms *float32 `json:"bathrooms,omitempty" binding:"omitempty,min=0"`
	Toilets   *int32   `json:"toilets,omitempty" binding:"omitempty,min=0"`
	Rooms     *int32   `json:"rooms,omitempty" binding:"omitempty,min=0"`
	Levels    *int32   `json:"levels,omitempty" binding:"omitempty,min=1"`

	// Areas
	TotalAreaSqm       *float64 `json:"total_area_sqm,omitempty" binding:"omitempty,min=0"`
	CoveredAreaSqm     *float64 `json:"covered_area_sqm,omitempty" binding:"omitempty,min=0"`
	SemicoveredAreaSqm *float64 `json:"semicovered_area_sqm,omitempty" binding:"omitempty,min=0"`
	UncoveredAreaSqm   *float64 `json:"uncovered_area_sqm,omitempty" binding:"omitempty,min=0"`

	// Features
	IncludesGarage *bool  `json:"includes_garage,omitempty"`
	GarageSpaces   *int32 `json:"garage_spaces,omitempty" binding:"omitempty,min=0"`
	Furnished      *bool  `json:"furnished,omitempty"`

	// Condition
	ConditionRating    *int32 `json:"condition_rating,omitempty" binding:"omitempty,min=1,max=5"`
	LastRenovationYear *int32 `json:"last_renovation_year,omitempty" binding:"omitempty,min=1900,max=2100"`

	Description *string `json:"description,omitempty"`
}

type PropertyUpdate struct {
	OwnerPersonID  *uuid.UUID `json:"owner_person_id,omitempty"`
	AddressID      *uuid.UUID `json:"address_id,omitempty"`
	PropertyTypeID *int32     `json:"property_type_id,omitempty"`
	InternalCode   *string    `json:"internal_code,omitempty" binding:"omitempty,max=50"`

	// Physical characteristics
	YearBuilt *int32   `json:"year_built,omitempty" binding:"omitempty,min=1800,max=2100"`
	Bedrooms  *int32   `json:"bedrooms,omitempty" binding:"omitempty,min=0"`
	Bathrooms *float32 `json:"bathrooms,omitempty" binding:"omitempty,min=0"`
	Toilets   *int32   `json:"toilets,omitempty" binding:"omitempty,min=0"`
	Rooms     *int32   `json:"rooms,omitempty" binding:"omitempty,min=0"`
	Levels    *int32   `json:"levels,omitempty" binding:"omitempty,min=1"`

	// Areas
	TotalAreaSqm       *float64 `json:"total_area_sqm,omitempty" binding:"omitempty,min=0"`
	CoveredAreaSqm     *float64 `json:"covered_area_sqm,omitempty" binding:"omitempty,min=0"`
	SemicoveredAreaSqm *float64 `json:"semicovered_area_sqm,omitempty" binding:"omitempty,min=0"`
	UncoveredAreaSqm   *float64 `json:"uncovered_area_sqm,omitempty" binding:"omitempty,min=0"`

	// Features
	IncludesGarage *bool  `json:"includes_garage,omitempty"`
	GarageSpaces   *int32 `json:"garage_spaces,omitempty" binding:"omitempty,min=0"`
	Furnished      *bool  `json:"furnished,omitempty"`

	// Condition
	ConditionRating    *int32 `json:"condition_rating,omitempty" binding:"omitempty,min=1,max=5"`
	LastRenovationYear *int32 `json:"last_renovation_year,omitempty" binding:"omitempty,min=1900,max=2100"`

	Description *string `json:"description,omitempty"`
}

type PropertyResponse struct {
	ID             uuid.UUID `json:"id"`
	OwnerPersonID  uuid.UUID `json:"owner_person_id"`
	AddressID      uuid.UUID `json:"address_id"`
	PropertyTypeID int32     `json:"property_type_id"`
	InternalCode   *string   `json:"internal_code"`

	// Physical characteristics
	YearBuilt *int32   `json:"year_built"`
	Bedrooms  *int32   `json:"bedrooms"`
	Bathrooms *float32 `json:"bathrooms"`
	Toilets   *int32   `json:"toilets"`
	Rooms     *int32   `json:"rooms"`
	Levels    *int32   `json:"levels"`

	// Areas
	TotalAreaSqm       *float64 `json:"total_area_sqm"`
	CoveredAreaSqm     *float64 `json:"covered_area_sqm"`
	SemicoveredAreaSqm *float64 `json:"semicovered_area_sqm"`
	UncoveredAreaSqm   *float64 `json:"uncovered_area_sqm"`

	// Features
	IncludesGarage *bool  `json:"includes_garage"`
	GarageSpaces   *int32 `json:"garage_spaces"`
	Furnished      *bool  `json:"furnished"`

	// Condition
	ConditionRating    *int32 `json:"condition_rating"`
	LastRenovationYear *int32 `json:"last_renovation_year"`

	Description *string    `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	UpdatedBy   *uuid.UUID `json:"updated_by"`

	// Embedded relations (optional)
	PropertyType        *PropertyTypeResponse        `json:"property_type,omitempty"`
	PropertyManagements []PropertyManagementResponse `json:"property_managements,omitempty"`
	PropertyAmenities   []PropertyAmenityResponse    `json:"property_amenities,omitempty"`
	PropertyListings    []PropertyListingResponse    `json:"property_listings,omitempty"`
	PropertyMedia       []PropertyMediaResponse      `json:"property_media,omitempty"`
	PropertyValuations  []PropertyValuationResponse  `json:"property_valuations,omitempty"`
}

// PropertyManagement DTOs
type PropertyManagementCreate struct {
	PropertyID           uuid.UUID  `json:"property_id" binding:"required"`
	ManagerID            uuid.UUID  `json:"manager_id" binding:"required"`
	ManagerTypeID        int32      `json:"manager_type_id" binding:"required"`
	OrganizationID       *uuid.UUID `json:"organization_id,omitempty"`
	StartDate            time.Time  `json:"start_date" binding:"required"`
	EndDate              *time.Time `json:"end_date,omitempty"`
	CommissionPercentage *float64   `json:"commission_percentage,omitempty" binding:"omitempty,min=0,max=100"`
	FixedFee             *float64   `json:"fixed_fee,omitempty" binding:"omitempty,min=0"`
	ExclusiveManagement  *bool      `json:"exclusive_management,omitempty"`
	ServicesIncluded     *string    `json:"services_included,omitempty"`
}

type PropertyManagementUpdate struct {
	PropertyID           *uuid.UUID `json:"property_id,omitempty"`
	ManagerID            *uuid.UUID `json:"manager_id,omitempty"`
	ManagerTypeID        *int32     `json:"manager_type_id,omitempty"`
	OrganizationID       *uuid.UUID `json:"organization_id,omitempty"`
	StartDate            *time.Time `json:"start_date,omitempty"`
	EndDate              *time.Time `json:"end_date,omitempty"`
	CommissionPercentage *float64   `json:"commission_percentage,omitempty" binding:"omitempty,min=0,max=100"`
	FixedFee             *float64   `json:"fixed_fee,omitempty" binding:"omitempty,min=0"`
	ExclusiveManagement  *bool      `json:"exclusive_management,omitempty"`
	ServicesIncluded     *string    `json:"services_included,omitempty"`
}

type PropertyManagementResponse struct {
	ID                   uuid.UUID  `json:"id"`
	PropertyID           uuid.UUID  `json:"property_id"`
	ManagerID            uuid.UUID  `json:"manager_id"`
	ManagerTypeID        int32      `json:"manager_type_id"`
	OrganizationID       *uuid.UUID `json:"organization_id"`
	StartDate            time.Time  `json:"start_date"`
	EndDate              *time.Time `json:"end_date"`
	IsActive             *bool      `json:"is_active"`
	CommissionPercentage *float64   `json:"commission_percentage"`
	FixedFee             *float64   `json:"fixed_fee"`
	ExclusiveManagement  *bool      `json:"exclusive_management"`
	ServicesIncluded     *string    `json:"services_included"`
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
	AmenityID  int32     `json:"amenity_id" binding:"required"`
	Note       *string   `json:"note,omitempty"`
}

type PropertyAmenityUpdate struct {
	Note *string `json:"note,omitempty"`
}

type PropertyAmenityResponse struct {
	PropertyID uuid.UUID `json:"property_id"`
	AmenityID  int32     `json:"amenity_id"`
	Note       *string   `json:"note"`
	CreatedAt  time.Time `json:"created_at"`

	// Embedded relations (optional)
	Property *PropertyResponse `json:"property,omitempty"`
	Amenity  *AmenityResponse  `json:"amenity,omitempty"`
}

// PropertyListing DTOs
type PropertyListingCreate struct {
	PropertyID       uuid.UUID            `json:"property_id" binding:"required"`
	OperationType    models.OperationType `json:"operation_type" binding:"required,oneof=sale rental temporary commercial investment"`
	Price            *float64             `json:"price,omitempty" binding:"omitempty,min=0"`
	Currency         *string              `json:"currency,omitempty" binding:"omitempty,len=3"`
	PricePerSqm      *float64             `json:"price_per_sqm,omitempty" binding:"omitempty,min=0"`
	IncludesExpenses *bool                `json:"includes_expenses,omitempty"`
	MonthlyExpenses  *float64             `json:"monthly_expenses,omitempty" binding:"omitempty,min=0"`
	DepositAmount    *float64             `json:"deposit_amount,omitempty" binding:"omitempty,min=0"`
	RentalTerms      *string              `json:"rental_terms,omitempty"`
	SaleTerms        *string              `json:"sale_terms,omitempty"`
	Title            string               `json:"title" binding:"required,max=200"`
	Description      *string              `json:"description,omitempty"`
	Highlighted      *bool                `json:"highlighted,omitempty"`
	PublicationStart *time.Time           `json:"publication_start,omitempty"`
	PublicationEnd   *time.Time           `json:"publication_end,omitempty"`
	AutoRenew        *bool                `json:"auto_renew,omitempty"`
	PublishOnWeb     *bool                `json:"publish_on_web,omitempty"`
	PublishOnPortals *bool                `json:"publish_on_portals,omitempty"`
	PortalConfig     *string              `json:"portal_config,omitempty"`
}

type PropertyListingUpdate struct {
	OperationType    *models.OperationType `json:"operation_type,omitempty"`
	ListingStatus    *models.ListingStatus `json:"listing_status,omitempty"`
	Price            *float64              `json:"price,omitempty" binding:"omitempty,min=0"`
	Currency         *string               `json:"currency,omitempty" binding:"omitempty,len=3"`
	PricePerSqm      *float64              `json:"price_per_sqm,omitempty" binding:"omitempty,min=0"`
	IncludesExpenses *bool                 `json:"includes_expenses,omitempty"`
	MonthlyExpenses  *float64              `json:"monthly_expenses,omitempty" binding:"omitempty,min=0"`
	DepositAmount    *float64              `json:"deposit_amount,omitempty" binding:"omitempty,min=0"`
	RentalTerms      *string               `json:"rental_terms,omitempty"`
	SaleTerms        *string               `json:"sale_terms,omitempty"`
	Title            *string               `json:"title,omitempty" binding:"omitempty,max=200"`
	Description      *string               `json:"description,omitempty"`
	Highlighted      *bool                 `json:"highlighted,omitempty"`
	PublicationStart *time.Time            `json:"publication_start,omitempty"`
	PublicationEnd   *time.Time            `json:"publication_end,omitempty"`
	AutoRenew        *bool                 `json:"auto_renew,omitempty"`
	PublishOnWeb     *bool                 `json:"publish_on_web,omitempty"`
	PublishOnPortals *bool                 `json:"publish_on_portals,omitempty"`
	PortalConfig     *string               `json:"portal_config,omitempty"`
}

type PropertyListingResponse struct {
	ID               uuid.UUID            `json:"id"`
	PropertyID       uuid.UUID            `json:"property_id"`
	OperationType    models.OperationType `json:"operation_type"`
	ListingStatus    models.ListingStatus `json:"listing_status"`
	Price            *float64             `json:"price"`
	Currency         string               `json:"currency"`
	PricePerSqm      *float64             `json:"price_per_sqm"`
	IncludesExpenses *bool                `json:"includes_expenses"`
	MonthlyExpenses  *float64             `json:"monthly_expenses"`
	DepositAmount    *float64             `json:"deposit_amount"`
	RentalTerms      *string              `json:"rental_terms"`
	SaleTerms        *string              `json:"sale_terms"`
	Title            string               `json:"title"`
	Description      *string              `json:"description"`
	Highlighted      *bool                `json:"highlighted"`
	PublicationStart *time.Time           `json:"publication_start"`
	PublicationEnd   *time.Time           `json:"publication_end"`
	AutoRenew        *bool                `json:"auto_renew"`
	PublishOnWeb     *bool                `json:"publish_on_web"`
	PublishOnPortals *bool                `json:"publish_on_portals"`
	PortalConfig     *string              `json:"portal_config"`
	ViewsCount       int32                `json:"views_count"`
	InquiriesCount   int32                `json:"inquiries_count"`
	FavoritesCount   int32                `json:"favorites_count"`
	CreatedAt        time.Time            `json:"created_at"`
	CreatedBy        *uuid.UUID           `json:"created_by"`
	UpdatedAt        *time.Time           `json:"updated_at"`
	UpdatedBy        *uuid.UUID           `json:"updated_by"`

	// Embedded relations (optional)
	Property *PropertyResponse `json:"property,omitempty"`
}

// PropertyMedia DTOs
type PropertyMediaCreate struct {
	PropertyID    uuid.UUID        `json:"property_id" binding:"required"`
	MediaType     models.MediaType `json:"media_type" binding:"required,oneof=photo video floor_plan document virtual_tour other"`
	FileName      string           `json:"file_name" binding:"required,max=255"`
	FileURL       string           `json:"file_url" binding:"required"`
	FileSizeBytes *int64           `json:"file_size_bytes,omitempty" binding:"omitempty,min=0"`
	MimeType      *string          `json:"mime_type,omitempty" binding:"omitempty,max=100"`
	Title         *string          `json:"title,omitempty" binding:"omitempty,max=200"`
	Description   *string          `json:"description,omitempty"`
	AltText       *string          `json:"alt_text,omitempty" binding:"omitempty,max=255"`
	DisplayOrder  *int32           `json:"display_order,omitempty" binding:"omitempty,min=0"`
	IsMain        *bool            `json:"is_main,omitempty"`
	IsPublic      *bool            `json:"is_public,omitempty"`
	WidthPx       *int32           `json:"width_px,omitempty" binding:"omitempty,min=0"`
	HeightPx      *int32           `json:"height_px,omitempty" binding:"omitempty,min=0"`
}

type PropertyMediaUpdate struct {
	MediaType     *models.MediaType `json:"media_type,omitempty"`
	FileName      *string           `json:"file_name,omitempty" binding:"omitempty,max=255"`
	FileURL       *string           `json:"file_url,omitempty"`
	FileSizeBytes *int64            `json:"file_size_bytes,omitempty" binding:"omitempty,min=0"`
	MimeType      *string           `json:"mime_type,omitempty" binding:"omitempty,max=100"`
	Title         *string           `json:"title,omitempty" binding:"omitempty,max=200"`
	Description   *string           `json:"description,omitempty"`
	AltText       *string           `json:"alt_text,omitempty" binding:"omitempty,max=255"`
	DisplayOrder  *int32            `json:"display_order,omitempty" binding:"omitempty,min=0"`
	IsMain        *bool             `json:"is_main,omitempty"`
	IsPublic      *bool             `json:"is_public,omitempty"`
	WidthPx       *int32            `json:"width_px,omitempty" binding:"omitempty,min=0"`
	HeightPx      *int32            `json:"height_px,omitempty" binding:"omitempty,min=0"`
}

type PropertyMediaResponse struct {
	ID            uuid.UUID        `json:"id"`
	PropertyID    uuid.UUID        `json:"property_id"`
	MediaType     models.MediaType `json:"media_type"`
	FileName      string           `json:"file_name"`
	FileURL       string           `json:"file_url"`
	FileSizeBytes *int64           `json:"file_size_bytes"`
	MimeType      *string          `json:"mime_type"`
	Title         *string          `json:"title"`
	Description   *string          `json:"description"`
	AltText       *string          `json:"alt_text"`
	DisplayOrder  int32            `json:"display_order"`
	IsMain        *bool            `json:"is_main"`
	IsPublic      *bool            `json:"is_public"`
	WidthPx       *int32           `json:"width_px"`
	HeightPx      *int32           `json:"height_px"`
	CreatedAt     time.Time        `json:"created_at"`
	CreatedBy     *uuid.UUID       `json:"created_by"`
	UpdatedAt     *time.Time       `json:"updated_at"`
	UpdatedBy     *uuid.UUID       `json:"updated_by"`

	// Embedded relations (optional)
	Property *PropertyResponse `json:"property,omitempty"`
}

// PropertyValuation DTOs
type PropertyValuationCreate struct {
	PropertyID            uuid.UUID  `json:"property_id" binding:"required"`
	ValuationDate         time.Time  `json:"valuation_date" binding:"required"`
	ValuationType         string     `json:"valuation_type" binding:"required,oneof=market fiscal insurance mortgage"`
	AppraisedValue        float64    `json:"appraised_value" binding:"required,min=0"`
	Currency              *string    `json:"currency,omitempty" binding:"omitempty,len=3"`
	ValuePerSqm           *float64   `json:"value_per_sqm,omitempty" binding:"omitempty,min=0"`
	AppraiserName         *string    `json:"appraiser_name,omitempty" binding:"omitempty,max=100"`
	AppraiserLicense      *string    `json:"appraiser_license,omitempty" binding:"omitempty,max=50"`
	AppraiserOrganization *string    `json:"appraiser_organization,omitempty" binding:"omitempty,max=100"`
	ValuationMethod       *string    `json:"valuation_method,omitempty" binding:"omitempty,oneof=comparative income cost"`
	MarketConditions      *string    `json:"market_conditions,omitempty"`
	AdjustmentsApplied    *string    `json:"adjustments_applied,omitempty"`
	ComparableProperties  *string    `json:"comparable_properties,omitempty"`
	ValuationReportURL    *string    `json:"valuation_report_url,omitempty"`
	PhotosURLs            *string    `json:"photos_urls,omitempty"`
	ValidUntil            *time.Time `json:"valid_until,omitempty"`
	Purpose               *string    `json:"purpose,omitempty" binding:"omitempty,max=100"`
}

type PropertyValuationUpdate struct {
	ValuationDate         *time.Time `json:"valuation_date,omitempty"`
	ValuationType         *string    `json:"valuation_type,omitempty" binding:"omitempty,oneof=market fiscal insurance mortgage"`
	AppraisedValue        *float64   `json:"appraised_value,omitempty" binding:"omitempty,min=0"`
	Currency              *string    `json:"currency,omitempty" binding:"omitempty,len=3"`
	ValuePerSqm           *float64   `json:"value_per_sqm,omitempty" binding:"omitempty,min=0"`
	AppraiserName         *string    `json:"appraiser_name,omitempty" binding:"omitempty,max=100"`
	AppraiserLicense      *string    `json:"appraiser_license,omitempty" binding:"omitempty,max=50"`
	AppraiserOrganization *string    `json:"appraiser_organization,omitempty" binding:"omitempty,max=100"`
	ValuationMethod       *string    `json:"valuation_method,omitempty" binding:"omitempty,oneof=comparative income cost"`
	MarketConditions      *string    `json:"market_conditions,omitempty"`
	AdjustmentsApplied    *string    `json:"adjustments_applied,omitempty"`
	ComparableProperties  *string    `json:"comparable_properties,omitempty"`
	ValuationReportURL    *string    `json:"valuation_report_url,omitempty"`
	PhotosURLs            *string    `json:"photos_urls,omitempty"`
	ValidUntil            *time.Time `json:"valid_until,omitempty"`
	Purpose               *string    `json:"purpose,omitempty" binding:"omitempty,max=100"`
}

type PropertyValuationResponse struct {
	ID                    uuid.UUID  `json:"id"`
	PropertyID            uuid.UUID  `json:"property_id"`
	ValuationDate         time.Time  `json:"valuation_date"`
	ValuationType         string     `json:"valuation_type"`
	AppraisedValue        float64    `json:"appraised_value"`
	Currency              string     `json:"currency"`
	ValuePerSqm           *float64   `json:"value_per_sqm"`
	AppraiserName         *string    `json:"appraiser_name"`
	AppraiserLicense      *string    `json:"appraiser_license"`
	AppraiserOrganization *string    `json:"appraiser_organization"`
	ValuationMethod       *string    `json:"valuation_method"`
	MarketConditions      *string    `json:"market_conditions"`
	AdjustmentsApplied    *string    `json:"adjustments_applied"`
	ComparableProperties  *string    `json:"comparable_properties"`
	ValuationReportURL    *string    `json:"valuation_report_url"`
	PhotosURLs            *string    `json:"photos_urls"`
	ValidUntil            *time.Time `json:"valid_until"`
	Purpose               *string    `json:"purpose"`
	CreatedAt             time.Time  `json:"created_at"`
	CreatedBy             *uuid.UUID `json:"created_by"`
	UpdatedAt             *time.Time `json:"updated_at"`
	UpdatedBy             *uuid.UUID `json:"updated_by"`

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
	PropertyID           uuid.UUID  `json:"property_id" binding:"required"`
	ManagerTypeID        int32      `json:"manager_type_id" binding:"required"`
	ManagerPersonID      uuid.UUID  `json:"manager_person_id" binding:"required"`
	StartDate            time.Time  `json:"start_date" binding:"required"`
	EndDate              *time.Time `json:"end_date,omitempty"`
	CommissionPercentage *float64   `json:"commission_percentage,omitempty" binding:"omitempty,min=0,max=100"`
	ContractReference    *string    `json:"contract_reference,omitempty"`
	Notes                *string    `json:"notes,omitempty"`
	IsActive             *bool      `json:"is_active,omitempty"`
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

// Filter DTOs for pagination and searching
type PropertyListingFilter struct {
	OperationType    *models.OperationType `json:"operation_type,omitempty"`
	ListingStatus    *models.ListingStatus `json:"listing_status,omitempty"`
	MinPrice         *float64              `json:"min_price,omitempty"`
	MaxPrice         *float64              `json:"max_price,omitempty"`
	Currency         *string               `json:"currency,omitempty"`
	Highlighted      *bool                 `json:"highlighted,omitempty"`
	PublishOnWeb     *bool                 `json:"publish_on_web,omitempty"`
	PublishOnPortals *bool                 `json:"publish_on_portals,omitempty"`
}

type PropertyMediaFilter struct {
	PropertyID *uuid.UUID        `json:"property_id,omitempty"`
	MediaType  *models.MediaType `json:"media_type,omitempty"`
	IsMain     *bool             `json:"is_main,omitempty"`
	IsPublic   *bool             `json:"is_public,omitempty"`
}

type PropertyValuationFilter struct {
	PropertyID      *uuid.UUID `json:"property_id,omitempty"`
	ValuationType   *string    `json:"valuation_type,omitempty"`
	MinValue        *float64   `json:"min_value,omitempty"`
	MaxValue        *float64   `json:"max_value,omitempty"`
	AppraiserName   *string    `json:"appraiser_name,omitempty"`
	ValuationMethod *string    `json:"valuation_method,omitempty"`
	DateFrom        *time.Time `json:"date_from,omitempty"`
	DateTo          *time.Time `json:"date_to,omitempty"`
}
