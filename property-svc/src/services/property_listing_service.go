package services

import (
	"time"

	"github.com/google/uuid"
	"github.com/rem-gestion/api-suite/property/src/dto"
	"github.com/rem-gestion/api-suite/property/src/grpc/clients"
	"github.com/rem-gestion/api-suite/property/src/models"
	"github.com/rem-gestion/api-suite/property/src/repository"
	rerrors "github.com/rem-gestion/rem-common/errors"
	"go.uber.org/zap"
)

type PropertyListingService struct {
	propertyListingRepo repository.PropertyListingRepository
	propertyRepo        repository.PropertyRepository
	clientManager       *clients.ClientManager
	lg                  *zap.Logger
}

func NewPropertyListingService(
	propertyListingRepo repository.PropertyListingRepository,
	propertyRepo repository.PropertyRepository,
	clientManager *clients.ClientManager,
	lg *zap.Logger,
) *PropertyListingService {
	return &PropertyListingService{
		propertyListingRepo: propertyListingRepo,
		propertyRepo:        propertyRepo,
		clientManager:       clientManager,
		lg:                  lg.Named("property-listing-service"),
	}
}

func (s *PropertyListingService) Create(in dto.PropertyListingCreate, createdBy *uuid.UUID) (*dto.PropertyListingResponse, error) {
	s.lg.Debug("create property listing request", zap.Any("payload", in))

	// Validate property exists
	property, err := s.propertyRepo.GetByID(in.PropertyID)
	if err != nil {
		s.lg.Warn("property not found", zap.String("property_id", in.PropertyID.String()))
		return nil, &rerrors.ValidationError{Msg: "Property not found"}
	}

	// Check if there's already an active listing for this property and operation type
	// Note: This is a business logic check - we'll get all listings and check manually
	existingListings, err := s.propertyListingRepo.GetByPropertyID(in.PropertyID)
	if err == nil {
		for _, existing := range existingListings {
			if existing.ListingStatus == models.ListingStatusActive && existing.OperationType == in.OperationType {
				s.lg.Warn("active listing already exists",
					zap.String("property_id", in.PropertyID.String()),
					zap.String("operation_type", string(in.OperationType)))
				return nil, &rerrors.ValidationError{Msg: "An active listing already exists for this property and operation type"}
			}
		}
	}

	// Create listing model
	listing := &models.PropertyListing{
		ID:               uuid.New(),
		PropertyID:       in.PropertyID,
		OperationType:    in.OperationType,
		ListingStatus:    models.ListingStatusDraft, // Default to draft
		Price:            in.Price,
		Currency:         "USD", // Default currency
		PricePerSqm:      in.PricePerSqm,
		IncludesExpenses: in.IncludesExpenses,
		MonthlyExpenses:  in.MonthlyExpenses,
		DepositAmount:    in.DepositAmount,
		RentalTerms:      in.RentalTerms,
		SaleTerms:        in.SaleTerms,
		Title:            in.Title,
		Description:      in.Description,
		Highlighted:      in.Highlighted,
		PublicationStart: in.PublicationStart,
		PublicationEnd:   in.PublicationEnd,
		AutoRenew:        in.AutoRenew,
		PublishOnWeb:     in.PublishOnWeb,
		PublishOnPortals: in.PublishOnPortals,
		PortalConfig:     in.PortalConfig,
		ViewsCount:       0,
		InquiriesCount:   0,
		FavoritesCount:   0,
		CreatedAt:        time.Now(),
		CreatedBy:        createdBy,
	}

	if in.Currency != nil {
		listing.Currency = *in.Currency
	}

	// Create listing
	createdListing, err := s.propertyListingRepo.Create(listing)
	if err != nil {
		s.lg.Error("failed to create property listing", zap.Error(err))
		return nil, err
	}

	// Convert to response DTO
	response := s.convertToDTO(createdListing, property)

	s.lg.Info("property listing created successfully", zap.String("id", createdListing.ID.String()))
	return response, nil
}

func (s *PropertyListingService) GetByID(id uuid.UUID) (*dto.PropertyListingResponse, error) {
	listing, err := s.propertyListingRepo.GetByID(id)
	if err != nil {
		s.lg.Warn("property listing not found", zap.String("id", id.String()))
		return nil, err
	}

	// Get property details
	property, err := s.propertyRepo.GetByID(listing.PropertyID)
	if err != nil {
		s.lg.Error("failed to get property for listing", zap.String("property_id", listing.PropertyID.String()))
		// Continue without property details
		property = nil
	}

	response := s.convertToDTO(listing, property)
	return response, nil
}

func (s *PropertyListingService) GetByPropertyID(propertyID uuid.UUID) ([]dto.PropertyListingResponse, error) {
	listings, err := s.propertyListingRepo.GetByPropertyID(propertyID)
	if err != nil {
		s.lg.Error("failed to get listings by property ID", zap.String("property_id", propertyID.String()))
		return nil, err
	}

	// Get property details once
	property, err := s.propertyRepo.GetByID(propertyID)
	if err != nil {
		s.lg.Error("failed to get property for listings", zap.String("property_id", propertyID.String()))
		property = nil
	}

	responses := make([]dto.PropertyListingResponse, len(listings))
	for i, listing := range listings {
		responses[i] = *s.convertToDTO(&listing, property)
	}

	return responses, nil
}

func (s *PropertyListingService) List(filters dto.PropertyListingFilter, page, limit int) ([]dto.PropertyListingResponse, int64, error) {
	offset := (page - 1) * limit

	repoFilters := repository.PropertyListingFilters{
		PropertyID:       nil, // Can be set from filters if needed
		OperationType:    filters.OperationType,
		ListingStatus:    filters.ListingStatus,
		MinPrice:         filters.MinPrice,
		MaxPrice:         filters.MaxPrice,
		Currency:         filters.Currency,
		PublishOnWeb:     filters.PublishOnWeb,
		PublishOnPortals: filters.PublishOnPortals,
		Highlighted:      filters.Highlighted,
		IncludeDeleted:   false,
	}

	listings, total, err := s.propertyListingRepo.List(repoFilters, limit, offset)
	if err != nil {
		s.lg.Error("failed to list property listings", zap.Error(err))
		return nil, 0, err
	}

	responses := make([]dto.PropertyListingResponse, len(listings))
	for i, listing := range listings {
		// For list operations, we might not want to load full property details
		responses[i] = *s.convertToDTO(&listing, nil)
	}

	return responses, total, nil
}

func (s *PropertyListingService) Update(id uuid.UUID, in dto.PropertyListingUpdate, updatedBy *uuid.UUID) (*dto.PropertyListingResponse, error) {
	s.lg.Debug("update property listing request", zap.String("id", id.String()), zap.Any("payload", in))

	// Get existing listing
	listing, err := s.propertyListingRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields
	now := time.Now()
	if in.OperationType != nil {
		listing.OperationType = *in.OperationType
	}
	if in.ListingStatus != nil {
		listing.ListingStatus = *in.ListingStatus
	}
	if in.Price != nil {
		listing.Price = in.Price
	}
	if in.Currency != nil {
		listing.Currency = *in.Currency
	}
	if in.PricePerSqm != nil {
		listing.PricePerSqm = in.PricePerSqm
	}
	if in.IncludesExpenses != nil {
		listing.IncludesExpenses = in.IncludesExpenses
	}
	if in.MonthlyExpenses != nil {
		listing.MonthlyExpenses = in.MonthlyExpenses
	}
	if in.DepositAmount != nil {
		listing.DepositAmount = in.DepositAmount
	}
	if in.RentalTerms != nil {
		listing.RentalTerms = in.RentalTerms
	}
	if in.SaleTerms != nil {
		listing.SaleTerms = in.SaleTerms
	}
	if in.Title != nil {
		listing.Title = *in.Title
	}
	if in.Description != nil {
		listing.Description = in.Description
	}
	if in.Highlighted != nil {
		listing.Highlighted = in.Highlighted
	}
	if in.PublicationStart != nil {
		listing.PublicationStart = in.PublicationStart
	}
	if in.PublicationEnd != nil {
		listing.PublicationEnd = in.PublicationEnd
	}
	if in.AutoRenew != nil {
		listing.AutoRenew = in.AutoRenew
	}
	if in.PublishOnWeb != nil {
		listing.PublishOnWeb = in.PublishOnWeb
	}
	if in.PublishOnPortals != nil {
		listing.PublishOnPortals = in.PublishOnPortals
	}
	if in.PortalConfig != nil {
		listing.PortalConfig = in.PortalConfig
	}

	listing.UpdatedAt = &now
	listing.UpdatedBy = updatedBy

	// Update listing
	err = s.propertyListingRepo.Update(listing)
	if err != nil {
		s.lg.Error("failed to update property listing", zap.String("id", id.String()), zap.Error(err))
		return nil, err
	}

	// Get property details
	property, err := s.propertyRepo.GetByID(listing.PropertyID)
	if err != nil {
		property = nil
	}

	response := s.convertToDTO(listing, property)

	s.lg.Info("property listing updated successfully", zap.String("id", id.String()))
	return response, nil
}

func (s *PropertyListingService) Delete(id uuid.UUID, deletedBy *uuid.UUID) error {
	s.lg.Debug("delete property listing request", zap.String("id", id.String()))

	// Check if listing exists
	_, err := s.propertyListingRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Soft delete listing
	if deletedBy != nil {
		err = s.propertyListingRepo.SoftDelete(id, *deletedBy)
	} else {
		err = s.propertyListingRepo.Delete(id)
	}

	if err != nil {
		s.lg.Error("failed to delete property listing", zap.String("id", id.String()), zap.Error(err))
		return err
	}

	s.lg.Info("property listing deleted successfully", zap.String("id", id.String()))
	return nil
}

func (s *PropertyListingService) UpdateStatus(id uuid.UUID, status models.ListingStatus, updatedBy uuid.UUID) error {
	s.lg.Debug("update property listing status", zap.String("id", id.String()), zap.String("status", string(status)))

	err := s.propertyListingRepo.UpdateStatus(id, status, updatedBy)
	if err != nil {
		s.lg.Error("failed to update property listing status", zap.String("id", id.String()), zap.Error(err))
		return err
	}

	s.lg.Info("property listing status updated successfully", zap.String("id", id.String()), zap.String("status", string(status)))
	return nil
}

func (s *PropertyListingService) IncrementViews(id uuid.UUID) error {
	return s.propertyListingRepo.IncrementViews(id)
}

func (s *PropertyListingService) IncrementInquiries(id uuid.UUID) error {
	return s.propertyListingRepo.IncrementInquiries(id)
}

func (s *PropertyListingService) IncrementFavorites(id uuid.UUID) error {
	return s.propertyListingRepo.IncrementFavorites(id)
}

func (s *PropertyListingService) DecrementFavorites(id uuid.UUID) error {
	// The repository method signature is different from what we expected
	// Let's use a manual update approach
	listing, err := s.propertyListingRepo.GetByID(id)
	if err != nil {
		return err
	}

	if listing.FavoritesCount > 0 {
		listing.FavoritesCount--
		now := time.Now()
		listing.UpdatedAt = &now
		return s.propertyListingRepo.Update(listing)
	}

	return nil
}

func (s *PropertyListingService) Search(query string, filters dto.PropertyListingFilter, page, limit int) ([]dto.PropertyListingResponse, int64, error) {
	offset := (page - 1) * limit

	repoFilters := repository.PropertyListingFilters{
		OperationType:    filters.OperationType,
		ListingStatus:    filters.ListingStatus,
		MinPrice:         filters.MinPrice,
		MaxPrice:         filters.MaxPrice,
		Currency:         filters.Currency,
		PublishOnWeb:     filters.PublishOnWeb,
		PublishOnPortals: filters.PublishOnPortals,
		Highlighted:      filters.Highlighted,
		IncludeDeleted:   false,
	}

	listings, total, err := s.propertyListingRepo.Search(query, repoFilters, limit, offset)
	if err != nil {
		s.lg.Error("failed to search property listings", zap.String("query", query), zap.Error(err))
		return nil, 0, err
	}

	responses := make([]dto.PropertyListingResponse, len(listings))
	for i, listing := range listings {
		responses[i] = *s.convertToDTO(&listing, nil)
	}

	return responses, total, nil
}

func (s *PropertyListingService) convertToDTO(listing *models.PropertyListing, property *models.Property) *dto.PropertyListingResponse {
	response := &dto.PropertyListingResponse{
		ID:               listing.ID,
		PropertyID:       listing.PropertyID,
		OperationType:    listing.OperationType,
		ListingStatus:    listing.ListingStatus,
		Price:            listing.Price,
		Currency:         listing.Currency,
		PricePerSqm:      listing.PricePerSqm,
		IncludesExpenses: listing.IncludesExpenses,
		MonthlyExpenses:  listing.MonthlyExpenses,
		DepositAmount:    listing.DepositAmount,
		RentalTerms:      listing.RentalTerms,
		SaleTerms:        listing.SaleTerms,
		Title:            listing.Title,
		Description:      listing.Description,
		Highlighted:      listing.Highlighted,
		PublicationStart: listing.PublicationStart,
		PublicationEnd:   listing.PublicationEnd,
		AutoRenew:        listing.AutoRenew,
		PublishOnWeb:     listing.PublishOnWeb,
		PublishOnPortals: listing.PublishOnPortals,
		PortalConfig:     listing.PortalConfig,
		ViewsCount:       listing.ViewsCount,
		InquiriesCount:   listing.InquiriesCount,
		FavoritesCount:   listing.FavoritesCount,
		CreatedAt:        listing.CreatedAt,
		CreatedBy:        listing.CreatedBy,
		UpdatedAt:        listing.UpdatedAt,
		UpdatedBy:        listing.UpdatedBy,
	}

	// Add property details if available
	if property != nil {
		response.Property = &dto.PropertyResponse{
			ID:                 property.ID,
			OwnerPersonID:      property.OwnerPersonID,
			PropertyTypeID:     property.PropertyTypeID,
			AddressID:          property.AddressID,
			InternalCode:       property.InternalCode,
			YearBuilt:          property.YearBuilt,
			Bedrooms:           property.Bedrooms,
			Bathrooms:          property.Bathrooms,
			Toilets:            property.Toilets,
			Rooms:              property.Rooms,
			Levels:             property.Levels,
			TotalAreaSqm:       property.TotalAreaSqm,
			CoveredAreaSqm:     property.CoveredAreaSqm,
			SemicoveredAreaSqm: property.SemicoveredAreaSqm,
			UncoveredAreaSqm:   property.UncoveredAreaSqm,
			IncludesGarage:     property.IncludesGarage,
			GarageSpaces:       property.GarageSpaces,
			Furnished:          property.Furnished,
			ConditionRating:    property.ConditionRating,
			LastRenovationYear: property.LastRenovationYear,
			Description:        property.Description,
			CreatedAt:          property.CreatedAt,
			UpdatedAt:          property.UpdatedAt,
			UpdatedBy:          property.UpdatedBy,
		}
	}

	return response
}
