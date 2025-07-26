package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rem-gestion/api-suite/property/src/dto"
	"github.com/rem-gestion/api-suite/property/src/grpc/clients"
	"github.com/rem-gestion/api-suite/property/src/models"
	"github.com/rem-gestion/api-suite/property/src/repository"
	"go.uber.org/zap"
)

type PropertyService struct {
	propertyRepo           repository.PropertyRepository
	propertyTypeRepo       repository.PropertyTypeRepository
	propertyManagementRepo repository.PropertyManagementRepository
	propertyAmenityRepo    repository.PropertyAmenityRepository
	clientManager          *clients.ClientManager
	lg                     *zap.Logger
}

func NewPropertyService(
	propertyRepo repository.PropertyRepository,
	propertyTypeRepo repository.PropertyTypeRepository,
	propertyManagementRepo repository.PropertyManagementRepository,
	propertyAmenityRepo repository.PropertyAmenityRepository,
	clientManager *clients.ClientManager,
	lg *zap.Logger,
) *PropertyService {
	return &PropertyService{
		propertyRepo:           propertyRepo,
		propertyTypeRepo:       propertyTypeRepo,
		propertyManagementRepo: propertyManagementRepo,
		propertyAmenityRepo:    propertyAmenityRepo,
		clientManager:          clientManager,
		lg:                     lg.Named("property-service"),
	}
}

func (s *PropertyService) Create(in dto.PropertyCreate, createdBy *uuid.UUID) (*dto.PropertyResponse, error) {
	s.lg.Debug("create property request", zap.Any("payload", in))

	ctx := context.Background()

	// Validate property type exists
	_, err := s.propertyTypeRepo.GetByID(in.PropertyTypeID)
	if err != nil {
		s.lg.Warn("property type not found", zap.Int32("property_type_id", in.PropertyTypeID))
		return nil, err
	}

	// Validate address exists (via gRPC)
	if err := s.clientManager.ValidateAddressExists(ctx, in.AddressID.String()); err != nil {
		s.lg.Warn("address validation failed",
			zap.String("address_id", in.AddressID.String()),
			zap.Error(err))
		return nil, err
	}

	// Validate owner person exists (via gRPC)
	if err := s.clientManager.ValidatePersonExists(ctx, in.OwnerPersonID.String()); err != nil {
		s.lg.Warn("owner person validation failed",
			zap.String("owner_person_id", in.OwnerPersonID.String()),
			zap.Error(err))
		return nil, err
	}

	property := &models.Property{
		ID:             uuid.New(),
		OwnerPersonID:  in.OwnerPersonID,
		AddressID:      in.AddressID,
		PropertyTypeID: in.PropertyTypeID,
		InternalCode:   in.InternalCode,
		YearBuilt:      in.YearBuilt,
		Bedrooms:       in.Bedrooms,
		Bathrooms:      in.Bathrooms,
		TotalAreaSqm:   in.TotalAreaSqm,
		CoveredAreaSqm: in.CoveredAreaSqm,
		Description:    in.Description,
		CreatedAt:      time.Now(),
	}

	if createdBy != nil {
		property.UpdatedBy = createdBy
	}

	created, err := s.propertyRepo.Create(property)
	if err != nil {
		s.lg.Error("create property failed", zap.Error(err))
		return nil, err
	}

	s.lg.Info("property created successfully", zap.String("id", created.ID.String()))
	return s.toPropertyResponse(created), nil
}

func (s *PropertyService) GetByID(id uuid.UUID, includeRelations bool) (*dto.PropertyResponse, error) {
	property, err := s.propertyRepo.GetByID(id)
	if err != nil {
		s.lg.Warn("get property failed", zap.String("id", id.String()), zap.Error(err))
		return nil, err
	}

	response := s.toPropertyResponse(property)

	if includeRelations {
		// Load additional relations if needed
		s.loadPropertyRelations(response, property)
	}

	return response, nil
}

func (s *PropertyService) GetByInternalCode(code string) (*dto.PropertyResponse, error) {
	property, err := s.propertyRepo.GetByInternalCode(code)
	if err != nil {
		s.lg.Warn("get property by code failed", zap.String("code", code), zap.Error(err))
		return nil, err
	}

	return s.toPropertyResponse(property), nil
}

func (s *PropertyService) List(req dto.PropertyListRequest) (*dto.ListResponse, error) {
	filters := repository.PropertyFilters{
		OwnerPersonID:  req.OwnerPersonID,
		PropertyTypeID: req.PropertyTypeID,
		InternalCode:   req.InternalCode,
		YearBuilt:      req.YearBuilt,
		Bedrooms:       req.Bedrooms,
		Bathrooms:      req.Bathrooms,
		MinTotalArea:   req.MinTotalArea,
		MaxTotalArea:   req.MaxTotalArea,
		City:           req.City,
		State:          req.State,
		IncludeDeleted: req.IncludeDeleted,
	}

	offset := (req.Page - 1) * req.Limit
	properties, total, err := s.propertyRepo.List(filters, req.Limit, offset)
	if err != nil {
		s.lg.Warn("list properties failed", zap.Error(err))
		return nil, err
	}

	var responses []dto.PropertyResponse
	for _, property := range properties {
		responses = append(responses, *s.toPropertyResponse(&property))
	}

	s.lg.Debug("property list result",
		zap.Int("responses_count", len(responses)),
		zap.Int64("total_count", total))

	return &dto.ListResponse{
		Data:  responses,
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}, nil
}

func (s *PropertyService) Update(id uuid.UUID, in dto.PropertyUpdate, updatedBy *uuid.UUID) (*dto.PropertyResponse, error) {
	s.lg.Debug("update property request", zap.String("id", id.String()), zap.Any("payload", in))

	property, err := s.propertyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Validate property type if being updated
	if in.PropertyTypeID != nil {
		_, err := s.propertyTypeRepo.GetByID(*in.PropertyTypeID)
		if err != nil {
			s.lg.Warn("property type not found", zap.Int32("property_type_id", *in.PropertyTypeID))
			return nil, err
		}
		property.PropertyTypeID = *in.PropertyTypeID
	}

	// Update fields
	s.updatePropertyFields(property, in)

	now := time.Now()
	property.UpdatedAt = &now
	if updatedBy != nil {
		property.UpdatedBy = updatedBy
	}

	if err := s.propertyRepo.Update(property); err != nil {
		s.lg.Error("update property failed", zap.String("id", id.String()), zap.Error(err))
		return nil, err
	}

	s.lg.Info("property updated successfully", zap.String("id", id.String()))
	return s.toPropertyResponse(property), nil
}

func (s *PropertyService) SoftDelete(id uuid.UUID, deletedBy uuid.UUID) error {
	// Verify property exists
	_, err := s.propertyRepo.GetByID(id)
	if err != nil {
		return err
	}

	if err := s.propertyRepo.SoftDelete(id, deletedBy); err != nil {
		s.lg.Error("soft delete property failed", zap.String("id", id.String()), zap.Error(err))
		return err
	}

	s.lg.Info("property soft deleted successfully", zap.String("id", id.String()))
	return nil
}

func (s *PropertyService) Delete(id uuid.UUID) error {
	// Verify property exists
	_, err := s.propertyRepo.GetByID(id)
	if err != nil {
		return err
	}

	if err := s.propertyRepo.Delete(id); err != nil {
		s.lg.Error("delete property failed", zap.String("id", id.String()), zap.Error(err))
		return err
	}

	s.lg.Info("property deleted successfully", zap.String("id", id.String()))
	return nil
}

func (s *PropertyService) Search(req dto.PropertySearchRequest) (*dto.ListResponse, error) {
	filters := repository.PropertyFilters{
		OwnerPersonID:  req.OwnerPersonID,
		PropertyTypeID: req.PropertyTypeID,
	}

	offset := (req.Page - 1) * req.Limit
	properties, total, err := s.propertyRepo.Search(req.Query, filters, req.Limit, offset)
	if err != nil {
		s.lg.Warn("search properties failed", zap.String("query", req.Query), zap.Error(err))
		return nil, err
	}

	var responses []dto.PropertyResponse
	for _, property := range properties {
		responses = append(responses, *s.toPropertyResponse(&property))
	}

	return &dto.ListResponse{
		Data:  responses,
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}, nil
}

func (s *PropertyService) BulkCreate(req dto.BulkPropertyCreateRequest, createdBy *uuid.UUID) (*dto.BulkPropertyCreateResponse, error) {
	s.lg.Debug("bulk create properties request", zap.Int("count", len(req.Properties)))

	var properties []models.Property
	for _, in := range req.Properties {
		// Validate property type exists
		_, err := s.propertyTypeRepo.GetByID(in.PropertyTypeID)
		if err != nil {
			s.lg.Warn("property type not found in bulk create", zap.Int32("property_type_id", in.PropertyTypeID))
			return nil, err
		}

		property := models.Property{
			ID:             uuid.New(),
			OwnerPersonID:  in.OwnerPersonID,
			AddressID:      in.AddressID,
			PropertyTypeID: in.PropertyTypeID,
			InternalCode:   in.InternalCode,
			YearBuilt:      in.YearBuilt,
			Bedrooms:       in.Bedrooms,
			Bathrooms:      in.Bathrooms,
			TotalAreaSqm:   in.TotalAreaSqm,
			CoveredAreaSqm: in.CoveredAreaSqm,
			Description:    in.Description,
			CreatedAt:      time.Now(),
		}

		if createdBy != nil {
			property.UpdatedBy = createdBy
		}

		properties = append(properties, property)
	}

	created, bulkErrors, err := s.propertyRepo.BulkCreate(properties)
	if err != nil {
		s.lg.Error("bulk create properties failed", zap.Error(err))
		return nil, err
	}

	var responses []dto.PropertyResponse
	for _, property := range created {
		responses = append(responses, *s.toPropertyResponse(&property))
	}

	var dtoErrors []dto.BulkError
	for _, bulkErr := range bulkErrors {
		dtoErrors = append(dtoErrors, dto.BulkError{
			Index: bulkErr.Index,
			Error: bulkErr.Error,
		})
	}

	s.lg.Info("bulk create properties completed",
		zap.Int("total", len(req.Properties)),
		zap.Int("created", len(created)),
		zap.Int("errors", len(dtoErrors)))

	return &dto.BulkPropertyCreateResponse{
		Created: responses,
		Errors:  dtoErrors,
	}, nil
}

// Helper methods
func (s *PropertyService) updatePropertyFields(property *models.Property, in dto.PropertyUpdate) {
	if in.OwnerPersonID != nil {
		property.OwnerPersonID = *in.OwnerPersonID
	}
	if in.AddressID != nil {
		property.AddressID = *in.AddressID
	}
	if in.InternalCode != nil {
		property.InternalCode = in.InternalCode
	}
	if in.YearBuilt != nil {
		property.YearBuilt = in.YearBuilt
	}
	if in.Bedrooms != nil {
		property.Bedrooms = in.Bedrooms
	}
	if in.Bathrooms != nil {
		property.Bathrooms = in.Bathrooms
	}
	if in.TotalAreaSqm != nil {
		property.TotalAreaSqm = in.TotalAreaSqm
	}
	if in.CoveredAreaSqm != nil {
		property.CoveredAreaSqm = in.CoveredAreaSqm
	}
	if in.Description != nil {
		property.Description = in.Description
	}
}

func (s *PropertyService) toPropertyResponse(property *models.Property) *dto.PropertyResponse {
	response := &dto.PropertyResponse{
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
		UpdatedBy:      property.UpdatedBy,
	}

	// Include embedded relations if loaded
	if property.PropertyType.ID != 0 {
		response.PropertyType = &dto.PropertyTypeResponse{
			ID:          property.PropertyType.ID,
			Code:        property.PropertyType.Code,
			Name:        property.PropertyType.Name,
			Description: property.PropertyType.Description,
			IsActive:    property.PropertyType.IsActive,
			CreatedAt:   property.PropertyType.CreatedAt,
			UpdatedAt:   property.PropertyType.UpdatedAt,
		}
	}

	return response
}

// ListPropertyAmenities returns all amenities associated with a property
func (s *PropertyService) ListPropertyAmenities(propertyID uuid.UUID) ([]*dto.AmenityResponse, error) {
	s.lg.Debug("listing property amenities", zap.String("property_id", propertyID.String()))

	// First verify the property exists
	_, err := s.propertyRepo.GetByID(propertyID)
	if err != nil {
		s.lg.Warn("property not found", zap.String("property_id", propertyID.String()), zap.Error(err))
		return nil, err
	}

	// Get amenities for this property
	amenities, err := s.propertyAmenityRepo.ListByPropertyID(propertyID)
	if err != nil {
		s.lg.Error("failed to list property amenities", zap.String("property_id", propertyID.String()), zap.Error(err))
		return nil, err
	}

	var response []*dto.AmenityResponse
	for _, amenity := range amenities {
		response = append(response, &dto.AmenityResponse{
			ID:       amenity.ID,
			Name:     amenity.Name,
			Category: amenity.Category,
			IconURL:  amenity.IconURL,
		})
	}

	s.lg.Debug("property amenities listed successfully", zap.String("property_id", propertyID.String()), zap.Int("count", len(response)))
	return response, nil
}

// AddAmenityToProperty adds an amenity to a property
func (s *PropertyService) AddAmenityToProperty(propertyID uuid.UUID, amenityID int32, note string) error {
	s.lg.Debug("adding amenity to property",
		zap.String("property_id", propertyID.String()),
		zap.Int32("amenity_id", amenityID),
		zap.String("note", note))

	// Verify property exists
	_, err := s.propertyRepo.GetByID(propertyID)
	if err != nil {
		s.lg.Warn("property not found", zap.String("property_id", propertyID.String()), zap.Error(err))
		return err
	}

	// Add amenity to property
	err = s.propertyAmenityRepo.AddAmenityToProperty(propertyID, amenityID, note)
	if err != nil {
		s.lg.Error("failed to add amenity to property",
			zap.String("property_id", propertyID.String()),
			zap.Int32("amenity_id", amenityID),
			zap.Error(err))
		return err
	}

	s.lg.Debug("amenity added to property successfully",
		zap.String("property_id", propertyID.String()),
		zap.Int32("amenity_id", amenityID))
	return nil
}

// RemoveAmenityFromProperty removes an amenity from a property
func (s *PropertyService) RemoveAmenityFromProperty(propertyID uuid.UUID, amenityID int32) error {
	s.lg.Debug("removing amenity from property",
		zap.String("property_id", propertyID.String()),
		zap.Int32("amenity_id", amenityID))

	// Verify property exists
	_, err := s.propertyRepo.GetByID(propertyID)
	if err != nil {
		s.lg.Warn("property not found", zap.String("property_id", propertyID.String()), zap.Error(err))
		return err
	}

	// Remove amenity from property
	err = s.propertyAmenityRepo.RemoveAmenityFromProperty(propertyID, amenityID)
	if err != nil {
		s.lg.Error("failed to remove amenity from property",
			zap.String("property_id", propertyID.String()),
			zap.Int32("amenity_id", amenityID),
			zap.Error(err))
		return err
	}

	s.lg.Debug("amenity removed from property successfully",
		zap.String("property_id", propertyID.String()),
		zap.Int32("amenity_id", amenityID))
	return nil
}

func (s *PropertyService) loadPropertyRelations(response *dto.PropertyResponse, property *models.Property) {
	// Load additional relations like managements and amenities if needed
	// This can be implemented based on specific requirements
}
