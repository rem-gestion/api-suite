package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rem-gestion/api-suite/property/src/dto"
	"github.com/rem-gestion/api-suite/property/src/grpc/clients"
	"github.com/rem-gestion/api-suite/property/src/models"
	"github.com/rem-gestion/api-suite/property/src/repository"
	rerrors "github.com/rem-gestion/rem-common/errors"
	"go.uber.org/zap"
)

type PropertyValuationService struct {
	propertyValuationRepo repository.PropertyValuationRepository
	propertyRepo          repository.PropertyRepository
	clientManager         *clients.ClientManager
	lg                    *zap.Logger
}

func NewPropertyValuationService(
	propertyValuationRepo repository.PropertyValuationRepository,
	propertyRepo repository.PropertyRepository,
	clientManager *clients.ClientManager,
	lg *zap.Logger,
) *PropertyValuationService {
	return &PropertyValuationService{
		propertyValuationRepo: propertyValuationRepo,
		propertyRepo:          propertyRepo,
		clientManager:         clientManager,
		lg:                    lg.Named("property-valuation-service"),
	}
}

func (s *PropertyValuationService) Create(in dto.PropertyValuationCreate, createdBy *uuid.UUID) (*dto.PropertyValuationResponse, error) {
	s.lg.Debug("create property valuation request", zap.Any("payload", in))

	// Validate property exists
	_, err := s.propertyRepo.GetByID(in.PropertyID)
	if err != nil {
		s.lg.Warn("property not found", zap.String("property_id", in.PropertyID.String()))
		return nil, &rerrors.ValidationError{Msg: "Property not found"}
	}

	// Validate valuation date is not in the future
	if in.ValuationDate.After(time.Now()) {
		return nil, &rerrors.ValidationError{Msg: "Valuation date cannot be in the future"}
	}

	// Create valuation model
	valuation := &models.PropertyValuation{
		ID:                    uuid.New(),
		PropertyID:            in.PropertyID,
		ValuationDate:         in.ValuationDate,
		ValuationType:         in.ValuationType,
		AppraisedValue:        in.AppraisedValue,
		Currency:              "USD", // Default currency
		ValuePerSqm:           in.ValuePerSqm,
		AppraiserName:         in.AppraiserName,
		AppraiserLicense:      in.AppraiserLicense,
		AppraiserOrganization: in.AppraiserOrganization,
		ValuationMethod:       in.ValuationMethod,
		MarketConditions:      in.MarketConditions,
		AdjustmentsApplied:    in.AdjustmentsApplied,
		ComparableProperties:  in.ComparableProperties,
		ValuationReportURL:    in.ValuationReportURL,
		PhotosURLs:            in.PhotosURLs,
		ValidUntil:            in.ValidUntil,
		Purpose:               in.Purpose,
		CreatedAt:             time.Now(),
		CreatedBy:             createdBy,
	}

	if in.Currency != nil {
		valuation.Currency = *in.Currency
	}

	// Create valuation
	createdValuation, err := s.propertyValuationRepo.Create(valuation)
	if err != nil {
		s.lg.Error("failed to create property valuation", zap.Error(err))
		return nil, err
	}

	// Convert to response DTO
	response := s.convertToDTO(createdValuation)

	s.lg.Info("property valuation created successfully", zap.String("id", createdValuation.ID.String()))
	return response, nil
}

func (s *PropertyValuationService) GetByID(id uuid.UUID) (*dto.PropertyValuationResponse, error) {
	valuation, err := s.propertyValuationRepo.GetByID(id)
	if err != nil {
		s.lg.Warn("property valuation not found", zap.String("id", id.String()))
		return nil, err
	}

	response := s.convertToDTO(valuation)
	return response, nil
}

func (s *PropertyValuationService) GetByPropertyID(propertyID uuid.UUID) ([]dto.PropertyValuationResponse, error) {
	valuations, err := s.propertyValuationRepo.GetByPropertyID(propertyID)
	if err != nil {
		s.lg.Error("failed to get valuations by property ID", zap.String("property_id", propertyID.String()))
		return nil, err
	}

	responses := make([]dto.PropertyValuationResponse, len(valuations))
	for i, valuation := range valuations {
		responses[i] = *s.convertToDTO(&valuation)
	}

	return responses, nil
}

func (s *PropertyValuationService) GetLatestByPropertyID(propertyID uuid.UUID) (*dto.PropertyValuationResponse, error) {
	valuation, err := s.propertyValuationRepo.GetLatestByPropertyID(propertyID)
	if err != nil {
		s.lg.Warn("latest valuation not found for property", zap.String("property_id", propertyID.String()))
		return nil, err
	}

	response := s.convertToDTO(valuation)
	return response, nil
}

func (s *PropertyValuationService) GetByPropertyIDAndType(propertyID uuid.UUID, valuationType string) ([]dto.PropertyValuationResponse, error) {
	valuations, err := s.propertyValuationRepo.GetByPropertyIDAndType(propertyID, valuationType)
	if err != nil {
		s.lg.Error("failed to get valuations by property ID and type",
			zap.String("property_id", propertyID.String()),
			zap.String("valuation_type", valuationType))
		return nil, err
	}

	responses := make([]dto.PropertyValuationResponse, len(valuations))
	for i, valuation := range valuations {
		responses[i] = *s.convertToDTO(&valuation)
	}

	return responses, nil
}

func (s *PropertyValuationService) List(filters dto.PropertyValuationFilter, page, limit int) ([]dto.PropertyValuationResponse, int64, error) {
	offset := (page - 1) * limit

	repoFilters := repository.PropertyValuationFilters{
		PropertyID:      filters.PropertyID,
		ValuationType:   filters.ValuationType,
		ValuationMethod: filters.ValuationMethod,
		MinValue:        filters.MinValue,
		MaxValue:        filters.MaxValue,
		ValidOnly:       false, // Can be set based on filters if needed
	}

	valuations, total, err := s.propertyValuationRepo.List(repoFilters, limit, offset)
	if err != nil {
		s.lg.Error("failed to list property valuations", zap.Error(err))
		return nil, 0, err
	}

	responses := make([]dto.PropertyValuationResponse, len(valuations))
	for i, valuation := range valuations {
		responses[i] = *s.convertToDTO(&valuation)
	}

	return responses, total, nil
}

func (s *PropertyValuationService) Update(id uuid.UUID, in dto.PropertyValuationUpdate, updatedBy *uuid.UUID) (*dto.PropertyValuationResponse, error) {
	s.lg.Debug("update property valuation request", zap.String("id", id.String()), zap.Any("payload", in))

	// Get existing valuation
	valuation, err := s.propertyValuationRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields
	now := time.Now()
	if in.ValuationDate != nil {
		// Validate valuation date is not in the future
		if in.ValuationDate.After(time.Now()) {
			return nil, &rerrors.ValidationError{Msg: "Valuation date cannot be in the future"}
		}
		valuation.ValuationDate = *in.ValuationDate
	}
	if in.ValuationType != nil {
		valuation.ValuationType = *in.ValuationType
	}
	if in.AppraisedValue != nil {
		valuation.AppraisedValue = *in.AppraisedValue
	}
	if in.Currency != nil {
		valuation.Currency = *in.Currency
	}
	if in.ValuePerSqm != nil {
		valuation.ValuePerSqm = in.ValuePerSqm
	}
	if in.AppraiserName != nil {
		valuation.AppraiserName = in.AppraiserName
	}
	if in.AppraiserLicense != nil {
		valuation.AppraiserLicense = in.AppraiserLicense
	}
	if in.AppraiserOrganization != nil {
		valuation.AppraiserOrganization = in.AppraiserOrganization
	}
	if in.ValuationMethod != nil {
		valuation.ValuationMethod = in.ValuationMethod
	}
	if in.MarketConditions != nil {
		valuation.MarketConditions = in.MarketConditions
	}
	if in.AdjustmentsApplied != nil {
		valuation.AdjustmentsApplied = in.AdjustmentsApplied
	}
	if in.ComparableProperties != nil {
		valuation.ComparableProperties = in.ComparableProperties
	}
	if in.ValuationReportURL != nil {
		valuation.ValuationReportURL = in.ValuationReportURL
	}
	if in.PhotosURLs != nil {
		valuation.PhotosURLs = in.PhotosURLs
	}
	if in.ValidUntil != nil {
		valuation.ValidUntil = in.ValidUntil
	}
	if in.Purpose != nil {
		valuation.Purpose = in.Purpose
	}

	valuation.UpdatedAt = &now
	valuation.UpdatedBy = updatedBy

	// Update valuation
	err = s.propertyValuationRepo.Update(valuation)
	if err != nil {
		s.lg.Error("failed to update property valuation", zap.String("id", id.String()), zap.Error(err))
		return nil, err
	}

	response := s.convertToDTO(valuation)

	s.lg.Info("property valuation updated successfully", zap.String("id", id.String()))
	return response, nil
}

func (s *PropertyValuationService) Delete(id uuid.UUID, deletedBy *uuid.UUID) error {
	s.lg.Debug("delete property valuation request", zap.String("id", id.String()))

	// Check if valuation exists
	_, err := s.propertyValuationRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Delete valuation (hard delete since soft delete is not available)
	err = s.propertyValuationRepo.Delete(id)

	if err != nil {
		s.lg.Error("failed to delete property valuation", zap.String("id", id.String()), zap.Error(err))
		return err
	}

	s.lg.Info("property valuation deleted successfully", zap.String("id", id.String()))
	return nil
}

// GetValuationsByType returns valuations by property ID and type
func (s *PropertyValuationService) GetValuationsByType(propertyID uuid.UUID, valuationType string) ([]dto.PropertyValuationResponse, error) {
	s.lg.Debug("get valuations by type",
		zap.String("property_id", propertyID.String()),
		zap.String("valuation_type", valuationType))

	// Verify property exists (removed organization check for now)
	_, err := s.propertyRepo.GetByID(propertyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get property: %w", err)
	}

	valuations, err := s.propertyValuationRepo.GetByPropertyIDAndType(propertyID, valuationType)
	if err != nil {
		s.lg.Error("failed to get valuations by type",
			zap.String("property_id", propertyID.String()),
			zap.String("valuation_type", valuationType),
			zap.Error(err))
		return nil, err
	}

	responses := make([]dto.PropertyValuationResponse, len(valuations))
	for i, valuation := range valuations {
		responses[i] = *s.convertToDTO(&valuation)
	}

	return responses, nil
}

// Note: Bulk operations are not available in the current repository interface
// These methods would need to be implemented in the repository layer first

// convertToDTO converts a PropertyValuation model to DTO
func (s *PropertyValuationService) convertToDTO(valuation *models.PropertyValuation) *dto.PropertyValuationResponse {
	return &dto.PropertyValuationResponse{
		ID:                    valuation.ID,
		PropertyID:            valuation.PropertyID,
		ValuationDate:         valuation.ValuationDate,
		ValuationType:         valuation.ValuationType,
		AppraisedValue:        valuation.AppraisedValue,
		Currency:              valuation.Currency,
		ValuePerSqm:           valuation.ValuePerSqm,
		AppraiserName:         valuation.AppraiserName,
		AppraiserLicense:      valuation.AppraiserLicense,
		AppraiserOrganization: valuation.AppraiserOrganization,
		ValuationMethod:       valuation.ValuationMethod,
		MarketConditions:      valuation.MarketConditions,
		AdjustmentsApplied:    valuation.AdjustmentsApplied,
		ComparableProperties:  valuation.ComparableProperties,
		ValuationReportURL:    valuation.ValuationReportURL,
		PhotosURLs:            valuation.PhotosURLs,
		ValidUntil:            valuation.ValidUntil,
		Purpose:               valuation.Purpose,
		CreatedAt:             valuation.CreatedAt,
		CreatedBy:             valuation.CreatedBy,
		UpdatedAt:             valuation.UpdatedAt,
		UpdatedBy:             valuation.UpdatedBy,
	}
}
