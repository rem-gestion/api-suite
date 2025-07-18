package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rem-gestion/api-suite/property/src/dto"
	"github.com/rem-gestion/api-suite/property/src/models"
	"github.com/rem-gestion/api-suite/property/src/repository"
	rerrors "github.com/rem-gestion/rem-common/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

/* ───────────────────── struct & ctor ─────────────────────── */

type PropertyService struct {
	repo    *repository.PropertyRepo
	lg      *zap.Logger
	timeout time.Duration
	mockSvc *MockExternalServices // Para simulaciones hasta que tengamos gRPC
}

func New(repo *repository.PropertyRepo, lg *zap.Logger) *PropertyService {
	return &PropertyService{
		repo:    repo,
		lg:      lg.Named("service"),
		timeout: 3 * time.Second,
		mockSvc: &MockExternalServices{}, // Inicializar mock service
	}
}

/* ─────────────────────── helpers ─────────────────────────── */

func (s *PropertyService) validateCreate(in dto.CreatePropertyDTO) error {
	// Validar datos básicos
	if in.OwnerPersonID == uuid.Nil {
		return &rerrors.ValidationError{Msg: "owner_person_id es requerido"}
	}
	if in.AddressID == uuid.Nil {
		return &rerrors.ValidationError{Msg: "address_id es requerido"}
	}

	// Validar que el owner_person_id existe (usando mock por ahora)
	if !s.mockSvc.ValidatePersonExists(in.OwnerPersonID) {
		return &rerrors.NotFoundError{Msg: "owner_person_id no existe"}
	}

	// Validar que el address_id existe (usando mock por ahora)
	if !s.mockSvc.ValidateAddressExists(in.AddressID) {
		return &rerrors.NotFoundError{Msg: "address_id no existe"}
	}

	// Validar property_type
	validTypes := []string{"APARTMENT", "HOUSE", "COMMERCIAL_SPACE", "OFFICE", "LAND", "INDUSTRIAL_WAREHOUSE"}
	isValid := false
	for _, t := range validTypes {
		if in.PropertyType == t {
			isValid = true
			break
		}
	}
	if !isValid {
		return &rerrors.ValidationError{Msg: "property_type debe ser uno de: " + fmt.Sprintf("%v", validTypes)}
	}

	// Si se especifica management, validar organization_id
	if in.Management != nil && !s.mockSvc.ValidateOrganizationExists(in.Management.OrganizationID) {
		return &rerrors.NotFoundError{Msg: "organization_id en management no existe"}
	}

	return nil
}

func (s *PropertyService) validateUpdate(in dto.UpdatePropertyDTO) error {
	// Si se especifica property_type, validarlo
	if in.PropertyType != nil {
		validTypes := []string{"APARTMENT", "HOUSE", "COMMERCIAL_SPACE", "OFFICE", "LAND", "INDUSTRIAL_WAREHOUSE"}
		isValid := false
		for _, t := range validTypes {
			if *in.PropertyType == t {
				isValid = true
				break
			}
		}
		if !isValid {
			return &rerrors.ValidationError{Msg: "property_type debe ser uno de: " + fmt.Sprintf("%v", validTypes)}
		}
	}

	return nil
}

/* ─────────────────────── Properties CRUD ─────────────────────── */

func (s *PropertyService) CreateProperty(ctx context.Context, in dto.CreatePropertyDTO) (*dto.PropertyResponseDTO, error) {
	s.lg.Debug("Creating property", zap.Any("input", in))

	if err := s.validateCreate(in); err != nil {
		return nil, err
	}

	// TODO: Aquí deberíamos validar que owner_person_id y address_id existen en sus microservicios
	// Por ahora, creamos datos mock/dummy

	// Mapear DTO a modelo
	property := &models.Property{
		OwnerPersonID:  in.OwnerPersonID,
		AddressID:      in.AddressID,
		PropertyType:   models.PropertyType(in.PropertyType),
		InternalCode:   in.InternalCode,
		YearBuilt:      in.YearBuilt,
		Bedrooms:       in.Bedrooms,
		Bathrooms:      in.Bathrooms,
		TotalAreaSqm:   in.TotalAreaSqm,
		CoveredAreaSqm: in.CoveredAreaSqm,
		Description:    in.Description,
	}

	// Crear en BD
	created, err := s.repo.Create(property)
	if err != nil {
		s.lg.Error("Failed to create property", zap.Error(err))
		return nil, &rerrors.InternalServerError{Msg: "failed to create property"}
	}

	// Crear management si se especificó
	if in.Management != nil {
		management := &models.PropertyManagement{
			PropertyID:        created.ID,
			OrganizationID:    in.Management.OrganizationID,
			ManagedSince:      time.Now(),
			CommissionPercent: in.Management.CommissionPercent,
			Notes:             in.Management.Notes,
		}
		if in.Management.ManagedSince != nil {
			management.ManagedSince = *in.Management.ManagedSince
		}
		if in.Management.ManagedUntil != nil {
			management.ManagedUntil = in.Management.ManagedUntil
		}

		_, err := s.repo.CreateManagement(management)
		if err != nil {
			s.lg.Warn("Failed to create property management", zap.Error(err))
		}
	}

	// Agregar amenities si se especificaron
	for _, amenityID := range in.AmenityIDs {
		err := s.repo.AddAmenityToProperty(created.ID, amenityID)
		if err != nil {
			s.lg.Warn("Failed to add amenity to property", zap.String("amenity_id", amenityID.String()), zap.Error(err))
		}
	}

	// Recargar con relaciones
	created, err = s.repo.Get(created.ID.String())
	if err != nil {
		s.lg.Error("Failed to reload created property", zap.Error(err))
		return nil, &rerrors.InternalServerError{Msg: "failed to reload property"}
	}

	response := dto.ToPropertyResponse(created)
	return &response, nil
}

func (s *PropertyService) GetProperty(ctx context.Context, id string) (*dto.PropertyResponseDTO, error) {
	s.lg.Debug("Getting property", zap.String("id", id))

	property, err := s.repo.Get(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &rerrors.NotFoundError{Msg: "property not found"}
		}
		s.lg.Error("Failed to get property", zap.Error(err))
		return nil, &rerrors.InternalServerError{Msg: "failed to get property"}
	}

	response := dto.ToPropertyResponse(property)
	return &response, nil
}

func (s *PropertyService) ListProperties(ctx context.Context, search string, propertyType *string, ownerPersonID *string, limit, offset int) ([]dto.PropertyResponseDTO, int64, error) {
	s.lg.Debug("Listing properties", zap.String("search", search), zap.Any("type", propertyType))

	var pType *models.PropertyType
	if propertyType != nil && *propertyType != "" {
		pt := models.PropertyType(*propertyType)
		pType = &pt
	}

	var ownerUUID *uuid.UUID
	if ownerPersonID != nil && *ownerPersonID != "" {
		parsed, err := uuid.Parse(*ownerPersonID)
		if err != nil {
			return nil, 0, &rerrors.BadRequestError{Msg: "invalid owner_person_id format"}
		}
		ownerUUID = &parsed
	}

	properties, total, err := s.repo.List(search, pType, ownerUUID, limit, offset)
	if err != nil {
		s.lg.Error("Failed to list properties", zap.Error(err))
		return nil, 0, &rerrors.InternalServerError{Msg: "failed to list properties"}
	}

	var responses []dto.PropertyResponseDTO
	for _, property := range properties {
		responses = append(responses, dto.ToPropertyResponse(&property))
	}

	return responses, total, nil
}

func (s *PropertyService) UpdateProperty(ctx context.Context, id string, in dto.UpdatePropertyDTO) (*dto.PropertyResponseDTO, error) {
	s.lg.Debug("Updating property", zap.String("id", id), zap.Any("input", in))

	if err := s.validateUpdate(in); err != nil {
		return nil, err
	}

	// Construir map de updates
	updates := make(map[string]interface{})
	if in.PropertyType != nil {
		updates["property_type"] = models.PropertyType(*in.PropertyType)
	}
	if in.InternalCode != nil {
		updates["internal_code"] = *in.InternalCode
	}
	if in.YearBuilt != nil {
		updates["year_built"] = *in.YearBuilt
	}
	if in.Bedrooms != nil {
		updates["bedrooms"] = *in.Bedrooms
	}
	if in.Bathrooms != nil {
		updates["bathrooms"] = *in.Bathrooms
	}
	if in.TotalAreaSqm != nil {
		updates["total_area_sqm"] = *in.TotalAreaSqm
	}
	if in.CoveredAreaSqm != nil {
		updates["covered_area_sqm"] = *in.CoveredAreaSqm
	}
	if in.Description != nil {
		updates["description"] = *in.Description
	}

	if len(updates) == 0 {
		return nil, &rerrors.BadRequestError{Msg: "no fields to update"}
	}

	updates["updated_at"] = time.Now()

	updated, err := s.repo.Update(id, updates)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &rerrors.NotFoundError{Msg: "property not found"}
		}
		s.lg.Error("Failed to update property", zap.Error(err))
		return nil, &rerrors.InternalServerError{Msg: "failed to update property"}
	}

	response := dto.ToPropertyResponse(updated)
	return &response, nil
}

func (s *PropertyService) DeleteProperty(ctx context.Context, id string) error {
	s.lg.Debug("Deleting property", zap.String("id", id))

	err := s.repo.Delete(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &rerrors.NotFoundError{Msg: "property not found"}
		}
		s.lg.Error("Failed to delete property", zap.Error(err))
		return &rerrors.InternalServerError{Msg: "failed to delete property"}
	}

	return nil
}

/* ─────────────────────── Amenities CRUD ─────────────────────── */

func (s *PropertyService) CreateAmenity(ctx context.Context, in dto.CreateAmenityDTO) (*dto.AmenityResponseDTO, error) {
	s.lg.Debug("Creating amenity", zap.Any("input", in))

	// Verificar que no exista ya una amenity con el mismo nombre
	_, err := s.repo.GetAmenityByName(in.Name)
	if err == nil {
		return nil, &rerrors.ConflictError{Msg: "amenity with this name already exists"}
	}

	amenity := &models.Amenity{
		Name:        in.Name,
		Description: in.Description,
		Icon:        in.Icon,
		Category:    in.Category,
	}

	created, err := s.repo.CreateAmenity(amenity)
	if err != nil {
		s.lg.Error("Failed to create amenity", zap.Error(err))
		return nil, &rerrors.InternalServerError{Msg: "failed to create amenity"}
	}

	response := dto.ToAmenityResponse(created)
	return &response, nil
}

func (s *PropertyService) GetAmenity(ctx context.Context, id string) (*dto.AmenityResponseDTO, error) {
	s.lg.Debug("Getting amenity", zap.String("id", id))

	amenity, err := s.repo.GetAmenity(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &rerrors.NotFoundError{Msg: "amenity not found"}
		}
		s.lg.Error("Failed to get amenity", zap.Error(err))
		return nil, &rerrors.InternalServerError{Msg: "failed to get amenity"}
	}

	response := dto.ToAmenityResponse(amenity)
	return &response, nil
}

func (s *PropertyService) ListAmenities(ctx context.Context, search string, category *string, limit, offset int) ([]dto.AmenityResponseDTO, int64, error) {
	s.lg.Debug("Listing amenities", zap.String("search", search), zap.Any("category", category))

	amenities, total, err := s.repo.ListAmenities(search, category, limit, offset)
	if err != nil {
		s.lg.Error("Failed to list amenities", zap.Error(err))
		return nil, 0, &rerrors.InternalServerError{Msg: "failed to list amenities"}
	}

	var responses []dto.AmenityResponseDTO
	for _, amenity := range amenities {
		responses = append(responses, dto.ToAmenityResponse(&amenity))
	}

	return responses, total, nil
}

func (s *PropertyService) UpdateAmenity(ctx context.Context, id string, in dto.UpdateAmenityDTO) (*dto.AmenityResponseDTO, error) {
	s.lg.Debug("Updating amenity", zap.String("id", id), zap.Any("input", in))

	// Si se está cambiando el nombre, verificar que no exista otro con ese nombre
	if in.Name != nil {
		existing, err := s.repo.GetAmenityByName(*in.Name)
		if err == nil && existing.ID.String() != id {
			return nil, &rerrors.ConflictError{Msg: "amenity with this name already exists"}
		}
	}

	// Construir map de updates
	updates := make(map[string]interface{})
	if in.Name != nil {
		updates["name"] = *in.Name
	}
	if in.Description != nil {
		updates["description"] = *in.Description
	}
	if in.Icon != nil {
		updates["icon"] = *in.Icon
	}
	if in.Category != nil {
		updates["category"] = *in.Category
	}

	if len(updates) == 0 {
		return nil, &rerrors.BadRequestError{Msg: "no fields to update"}
	}

	updates["updated_at"] = time.Now()

	updated, err := s.repo.UpdateAmenity(id, updates)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &rerrors.NotFoundError{Msg: "amenity not found"}
		}
		s.lg.Error("Failed to update amenity", zap.Error(err))
		return nil, &rerrors.InternalServerError{Msg: "failed to update amenity"}
	}

	response := dto.ToAmenityResponse(updated)
	return &response, nil
}

func (s *PropertyService) DeleteAmenity(ctx context.Context, id string) error {
	s.lg.Debug("Deleting amenity", zap.String("id", id))

	err := s.repo.DeleteAmenity(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &rerrors.NotFoundError{Msg: "amenity not found"}
		}
		s.lg.Error("Failed to delete amenity", zap.Error(err))
		return &rerrors.InternalServerError{Msg: "failed to delete amenity"}
	}

	return nil
}

/* ─────────────────────── Property Amenities Management ─────────────────────── */

func (s *PropertyService) AddAmenityToProperty(ctx context.Context, propertyID, amenityID string) error {
	s.lg.Debug("Adding amenity to property", zap.String("property_id", propertyID), zap.String("amenity_id", amenityID))

	propUUID, err := uuid.Parse(propertyID)
	if err != nil {
		return &rerrors.BadRequestError{Msg: "invalid property_id format"}
	}

	amenUUID, err := uuid.Parse(amenityID)
	if err != nil {
		return &rerrors.BadRequestError{Msg: "invalid amenity_id format"}
	}

	// Verificar que existan
	_, err = s.repo.Get(propertyID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &rerrors.NotFoundError{Msg: "property not found"}
		}
		return &rerrors.InternalServerError{Msg: "failed to verify property"}
	}

	_, err = s.repo.GetAmenity(amenityID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &rerrors.NotFoundError{Msg: "amenity not found"}
		}
		return &rerrors.InternalServerError{Msg: "failed to verify amenity"}
	}

	err = s.repo.AddAmenityToProperty(propUUID, amenUUID)
	if err != nil {
		s.lg.Error("Failed to add amenity to property", zap.Error(err))
		return &rerrors.InternalServerError{Msg: "failed to add amenity to property"}
	}

	return nil
}

func (s *PropertyService) RemoveAmenityFromProperty(ctx context.Context, propertyID, amenityID string) error {
	s.lg.Debug("Removing amenity from property", zap.String("property_id", propertyID), zap.String("amenity_id", amenityID))

	propUUID, err := uuid.Parse(propertyID)
	if err != nil {
		return &rerrors.BadRequestError{Msg: "invalid property_id format"}
	}

	amenUUID, err := uuid.Parse(amenityID)
	if err != nil {
		return &rerrors.BadRequestError{Msg: "invalid amenity_id format"}
	}

	err = s.repo.RemoveAmenityFromProperty(propUUID, amenUUID)
	if err != nil {
		s.lg.Error("Failed to remove amenity from property", zap.Error(err))
		return &rerrors.InternalServerError{Msg: "failed to remove amenity from property"}
	}

	return nil
}
