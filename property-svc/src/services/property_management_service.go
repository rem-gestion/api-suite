package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/rem-gestion/api-suite/property/src/dto"
	"github.com/rem-gestion/api-suite/property/src/models"
	"github.com/rem-gestion/api-suite/property/src/repository"
	"go.uber.org/zap"
)

type PropertyManagementService interface {
	Create(ctx context.Context, req *dto.PropertyManagementCreateRequest) (*dto.PropertyManagementResponse, error)
	List(ctx context.Context, req *dto.PropertyManagementListRequest) (*dto.PropertyManagementListResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dto.PropertyManagementResponse, error)
	Update(ctx context.Context, id uuid.UUID, req *dto.PropertyManagementUpdateRequest) (*dto.PropertyManagementResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type propertyManagementService struct {
	repo   repository.PropertyManagementRepository
	logger *zap.Logger
}

func NewPropertyManagementService(repo repository.PropertyManagementRepository, logger *zap.Logger) PropertyManagementService {
	return &propertyManagementService{
		repo:   repo,
		logger: logger.Named("property-management-service"),
	}
}

func (s *propertyManagementService) Create(ctx context.Context, req *dto.PropertyManagementCreateRequest) (*dto.PropertyManagementResponse, error) {
	s.logger.Info("creating property management", zap.String("property_id", req.PropertyID.String()))

	propertyManagement := &models.PropertyManagement{
		ID:                   uuid.New(),
		PropertyID:           req.PropertyID,
		ManagerID:            req.ManagerPersonID,
		ManagerTypeID:        req.ManagerTypeID,
		StartDate:            req.StartDate,
		EndDate:              req.EndDate,
		CommissionPercentage: req.CommissionPercentage,
	}

	createdManagement, err := s.repo.Create(propertyManagement)
	if err != nil {
		s.logger.Error("failed to create property management", zap.Error(err))
		return nil, err
	}

	return s.modelToResponse(createdManagement), nil
}

func (s *propertyManagementService) List(ctx context.Context, req *dto.PropertyManagementListRequest) (*dto.PropertyManagementListResponse, error) {
	s.logger.Info("listing property managements")

	// Set defaults
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 20
	}

	// Convert DTO filters to repository filters
	filters := repository.PropertyManagementFilters{
		PropertyID: req.PropertyID,
		ActiveOnly: req.IsActive != nil && *req.IsActive,
	}

	offset := (req.Page - 1) * req.Limit
	managements, total, err := s.repo.List(filters, req.Limit, offset)
	if err != nil {
		s.logger.Error("failed to list property managements", zap.Error(err))
		return nil, err
	}

	responses := make([]dto.PropertyManagementResponse, len(managements))
	for i, management := range managements {
		responses[i] = *s.modelToResponse(&management)
	}

	return &dto.PropertyManagementListResponse{
		Data:  responses,
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}, nil
}

func (s *propertyManagementService) GetByID(ctx context.Context, id uuid.UUID) (*dto.PropertyManagementResponse, error) {
	s.logger.Info("getting property management by ID", zap.String("id", id.String()))

	management, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("failed to get property management", zap.String("id", id.String()), zap.Error(err))
		return nil, err
	}

	return s.modelToResponse(management), nil
}

func (s *propertyManagementService) Update(ctx context.Context, id uuid.UUID, req *dto.PropertyManagementUpdateRequest) (*dto.PropertyManagementResponse, error) {
	s.logger.Info("updating property management", zap.String("id", id.String()))

	// Get existing management
	existing, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("failed to get existing property management", zap.String("id", id.String()), zap.Error(err))
		return nil, err
	}

	// Update fields if provided
	if req.ManagerTypeID != nil {
		existing.ManagerTypeID = *req.ManagerTypeID
	}
	if req.ManagerPersonID != nil {
		existing.ManagerID = *req.ManagerPersonID
	}
	if req.StartDate != nil {
		existing.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		existing.EndDate = req.EndDate
	}
	if req.CommissionPercentage != nil {
		existing.CommissionPercentage = req.CommissionPercentage
	}

	err = s.repo.Update(existing)
	if err != nil {
		s.logger.Error("failed to update property management", zap.String("id", id.String()), zap.Error(err))
		return nil, err
	}

	// Get updated record
	updatedManagement, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("failed to get updated property management", zap.String("id", id.String()), zap.Error(err))
		return nil, err
	}

	return s.modelToResponse(updatedManagement), nil
}

func (s *propertyManagementService) Delete(ctx context.Context, id uuid.UUID) error {
	s.logger.Info("deleting property management", zap.String("id", id.String()))

	err := s.repo.Delete(id)
	if err != nil {
		s.logger.Error("failed to delete property management", zap.String("id", id.String()), zap.Error(err))
		return err
	}

	return nil
}

func (s *propertyManagementService) modelToResponse(management *models.PropertyManagement) *dto.PropertyManagementResponse {
	return &dto.PropertyManagementResponse{
		ID:                   management.ID,
		PropertyID:           management.PropertyID,
		ManagerID:            management.ManagerID,
		ManagerTypeID:        management.ManagerTypeID,
		StartDate:            management.StartDate,
		EndDate:              management.EndDate,
		CommissionPercentage: management.CommissionPercentage,
		CreatedAt:            management.CreatedAt,
		UpdatedAt:            management.UpdatedAt,
		UpdatedBy:            management.UpdatedBy,
	}
}
