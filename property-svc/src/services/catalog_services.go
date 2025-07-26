package services

import (
	"time"

	"github.com/rem-gestion/api-suite/property/src/dto"
	"github.com/rem-gestion/api-suite/property/src/models"
	"github.com/rem-gestion/api-suite/property/src/repository"
	"go.uber.org/zap"
)

// PropertyTypeService handles property type operations
type PropertyTypeService struct {
	repo repository.PropertyTypeRepository
	lg   *zap.Logger
}

func NewPropertyTypeService(repo repository.PropertyTypeRepository, lg *zap.Logger) *PropertyTypeService {
	return &PropertyTypeService{
		repo: repo,
		lg:   lg.Named("property-type-service"),
	}
}

func (s *PropertyTypeService) Create(in dto.PropertyTypeCreate) (*dto.PropertyTypeResponse, error) {
	s.lg.Debug("create property type request", zap.Any("payload", in))

	propertyType := &models.PropertyType{
		Code:        in.Code,
		Name:        in.Name,
		Description: in.Description,
		IsActive:    true,
		CreatedAt:   time.Now(),
	}

	if in.IsActive != nil {
		propertyType.IsActive = *in.IsActive
	}

	created, err := s.repo.Create(propertyType)
	if err != nil {
		s.lg.Error("create property type failed", zap.Error(err))
		return nil, err
	}

	s.lg.Info("property type created successfully", zap.String("code", created.Code))
	return s.toPropertyTypeResponse(created), nil
}

func (s *PropertyTypeService) GetByID(id int32) (*dto.PropertyTypeResponse, error) {
	propertyType, err := s.repo.GetByID(id)
	if err != nil {
		s.lg.Warn("get property type failed", zap.Int32("id", id), zap.Error(err))
		return nil, err
	}

	return s.toPropertyTypeResponse(propertyType), nil
}

func (s *PropertyTypeService) GetByCode(code string) (*dto.PropertyTypeResponse, error) {
	propertyType, err := s.repo.GetByCode(code)
	if err != nil {
		s.lg.Warn("get property type by code failed", zap.String("code", code), zap.Error(err))
		return nil, err
	}

	return s.toPropertyTypeResponse(propertyType), nil
}

func (s *PropertyTypeService) List(isActive *bool, page, limit int) (*dto.ListResponse, error) {
	offset := (page - 1) * limit
	propertyTypes, total, err := s.repo.List(nil, isActive, limit, offset)
	if err != nil {
		s.lg.Warn("list property types failed", zap.Error(err))
		return nil, err
	}

	var responses []dto.PropertyTypeResponse
	for _, propertyType := range propertyTypes {
		responses = append(responses, *s.toPropertyTypeResponse(&propertyType))
	}

	return &dto.ListResponse{
		Data:  responses,
		Page:  page,
		Limit: limit,
		Total: total,
	}, nil
}

func (s *PropertyTypeService) Update(id int32, in dto.PropertyTypeUpdate) (*dto.PropertyTypeResponse, error) {
	s.lg.Debug("update property type request", zap.Int32("id", id), zap.Any("payload", in))

	propertyType, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if in.Code != nil {
		propertyType.Code = *in.Code
	}
	if in.Name != nil {
		propertyType.Name = *in.Name
	}
	if in.Description != nil {
		propertyType.Description = in.Description
	}
	if in.IsActive != nil {
		propertyType.IsActive = *in.IsActive
	}

	now := time.Now()
	propertyType.UpdatedAt = &now

	if err := s.repo.Update(propertyType); err != nil {
		s.lg.Error("update property type failed", zap.Int32("id", id), zap.Error(err))
		return nil, err
	}

	s.lg.Info("property type updated successfully", zap.Int32("id", id))
	return s.toPropertyTypeResponse(propertyType), nil
}

func (s *PropertyTypeService) Delete(id int32) error {
	// Verify property type exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		s.lg.Error("delete property type failed", zap.Int32("id", id), zap.Error(err))
		return err
	}

	s.lg.Info("property type deleted successfully", zap.Int32("id", id))
	return nil
}

func (s *PropertyTypeService) toPropertyTypeResponse(propertyType *models.PropertyType) *dto.PropertyTypeResponse {
	return &dto.PropertyTypeResponse{
		ID:          propertyType.ID,
		Code:        propertyType.Code,
		Name:        propertyType.Name,
		Description: propertyType.Description,
		IsActive:    propertyType.IsActive,
		CreatedAt:   propertyType.CreatedAt,
		UpdatedAt:   propertyType.UpdatedAt,
	}
}

// ManagerTypeService handles manager type operations
type ManagerTypeService struct {
	repo repository.ManagerTypeRepository
	lg   *zap.Logger
}

func NewManagerTypeService(repo repository.ManagerTypeRepository, lg *zap.Logger) *ManagerTypeService {
	return &ManagerTypeService{
		repo: repo,
		lg:   lg.Named("manager-type-service"),
	}
}

func (s *ManagerTypeService) Create(in dto.ManagerTypeCreate) (*dto.ManagerTypeResponse, error) {
	s.lg.Debug("create manager type request", zap.Any("payload", in))

	managerType := &models.ManagerType{
		Code:        in.Code,
		Name:        in.Name,
		Description: in.Description,
		IsActive:    true,
		CreatedAt:   time.Now(),
	}

	if in.IsActive != nil {
		managerType.IsActive = *in.IsActive
	}

	created, err := s.repo.Create(managerType)
	if err != nil {
		s.lg.Error("create manager type failed", zap.Error(err))
		return nil, err
	}

	s.lg.Info("manager type created successfully", zap.String("code", created.Code))
	return s.toManagerTypeResponse(created), nil
}

func (s *ManagerTypeService) GetByID(id int32) (*dto.ManagerTypeResponse, error) {
	managerType, err := s.repo.GetByID(id)
	if err != nil {
		s.lg.Warn("get manager type failed", zap.Int32("id", id), zap.Error(err))
		return nil, err
	}

	return s.toManagerTypeResponse(managerType), nil
}

func (s *ManagerTypeService) List(isActive *bool, page, limit int) (*dto.ListResponse, error) {
	offset := (page - 1) * limit
	managerTypes, total, err := s.repo.List(isActive, limit, offset)
	if err != nil {
		s.lg.Warn("list manager types failed", zap.Error(err))
		return nil, err
	}

	var responses []dto.ManagerTypeResponse
	for _, managerType := range managerTypes {
		responses = append(responses, *s.toManagerTypeResponse(&managerType))
	}

	return &dto.ListResponse{
		Data:  responses,
		Page:  page,
		Limit: limit,
		Total: total,
	}, nil
}

func (s *ManagerTypeService) Update(id int32, in dto.ManagerTypeUpdate) (*dto.ManagerTypeResponse, error) {
	s.lg.Debug("update manager type request", zap.Int32("id", id), zap.Any("payload", in))

	managerType, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if in.Code != nil {
		managerType.Code = *in.Code
	}
	if in.Name != nil {
		managerType.Name = *in.Name
	}
	if in.Description != nil {
		managerType.Description = in.Description
	}
	if in.IsActive != nil {
		managerType.IsActive = *in.IsActive
	}

	now := time.Now()
	managerType.UpdatedAt = &now

	if err := s.repo.Update(managerType); err != nil {
		s.lg.Error("update manager type failed", zap.Int32("id", id), zap.Error(err))
		return nil, err
	}

	s.lg.Info("manager type updated successfully", zap.Int32("id", id))
	return s.toManagerTypeResponse(managerType), nil
}

func (s *ManagerTypeService) Delete(id int32) error {
	// Verify manager type exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		s.lg.Error("delete manager type failed", zap.Int32("id", id), zap.Error(err))
		return err
	}

	s.lg.Info("manager type deleted successfully", zap.Int32("id", id))
	return nil
}

func (s *ManagerTypeService) toManagerTypeResponse(managerType *models.ManagerType) *dto.ManagerTypeResponse {
	return &dto.ManagerTypeResponse{
		ID:          managerType.ID,
		Code:        managerType.Code,
		Name:        managerType.Name,
		Description: managerType.Description,
		IsActive:    managerType.IsActive,
		CreatedAt:   managerType.CreatedAt,
		UpdatedAt:   managerType.UpdatedAt,
	}
}

// AmenityService handles amenity operations
type AmenityService struct {
	repo repository.AmenityRepository
	lg   *zap.Logger
}

func NewAmenityService(repo repository.AmenityRepository, lg *zap.Logger) *AmenityService {
	return &AmenityService{
		repo: repo,
		lg:   lg.Named("amenity-service"),
	}
}

func (s *AmenityService) Create(in dto.AmenityCreate) (*dto.AmenityResponse, error) {
	s.lg.Debug("create amenity request", zap.Any("payload", in))

	amenity := &models.Amenity{
		Name:     in.Name,
		Category: in.Category,
		IconURL:  in.IconURL,
	}

	created, err := s.repo.Create(amenity)
	if err != nil {
		s.lg.Error("create amenity failed", zap.Error(err))
		return nil, err
	}

	s.lg.Info("amenity created successfully", zap.String("name", created.Name))
	return s.toAmenityResponse(created), nil
}

func (s *AmenityService) GetByID(id int32) (*dto.AmenityResponse, error) {
	amenity, err := s.repo.GetByID(id)
	if err != nil {
		s.lg.Warn("get amenity failed", zap.Int32("id", id), zap.Error(err))
		return nil, err
	}

	return s.toAmenityResponse(amenity), nil
}

func (s *AmenityService) List(category *models.AmenityCategory, page, limit int) (*dto.ListResponse, error) {
	offset := (page - 1) * limit
	amenities, total, err := s.repo.List(category, nil, limit, offset)
	if err != nil {
		s.lg.Warn("list amenities failed", zap.Error(err))
		return nil, err
	}

	var responses []dto.AmenityResponse
	for _, amenity := range amenities {
		responses = append(responses, *s.toAmenityResponse(&amenity))
	}

	return &dto.ListResponse{
		Data:  responses,
		Page:  page,
		Limit: limit,
		Total: total,
	}, nil
}

func (s *AmenityService) Update(id int32, in dto.AmenityUpdate) (*dto.AmenityResponse, error) {
	s.lg.Debug("update amenity request", zap.Int32("id", id), zap.Any("payload", in))

	amenity, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if in.Name != nil {
		amenity.Name = *in.Name
	}
	if in.Category != nil {
		amenity.Category = in.Category
	}
	if in.IconURL != nil {
		amenity.IconURL = in.IconURL
	}

	if err := s.repo.Update(amenity); err != nil {
		s.lg.Error("update amenity failed", zap.Int32("id", id), zap.Error(err))
		return nil, err
	}

	s.lg.Info("amenity updated successfully", zap.Int32("id", id))
	return s.toAmenityResponse(amenity), nil
}

func (s *AmenityService) Delete(id int32) error {
	// Verify amenity exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		s.lg.Error("delete amenity failed", zap.Int32("id", id), zap.Error(err))
		return err
	}

	s.lg.Info("amenity deleted successfully", zap.Int32("id", id))
	return nil
}

func (s *AmenityService) toAmenityResponse(amenity *models.Amenity) *dto.AmenityResponse {
	return &dto.AmenityResponse{
		ID:       amenity.ID,
		Name:     amenity.Name,
		Category: amenity.Category,
		IconURL:  amenity.IconURL,
	}
}
