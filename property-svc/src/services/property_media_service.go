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

type PropertyMediaService struct {
	propertyMediaRepo repository.PropertyMediaRepository
	propertyRepo      repository.PropertyRepository
	clientManager     *clients.ClientManager
	lg                *zap.Logger
}

func NewPropertyMediaService(
	propertyMediaRepo repository.PropertyMediaRepository,
	propertyRepo repository.PropertyRepository,
	clientManager *clients.ClientManager,
	lg *zap.Logger,
) *PropertyMediaService {
	return &PropertyMediaService{
		propertyMediaRepo: propertyMediaRepo,
		propertyRepo:      propertyRepo,
		clientManager:     clientManager,
		lg:                lg.Named("property-media-service"),
	}
}

func (s *PropertyMediaService) Create(in dto.PropertyMediaCreate, createdBy *uuid.UUID) (*dto.PropertyMediaResponse, error) {
	s.lg.Debug("create property media request", zap.Any("payload", in))

	// Validate property exists
	_, err := s.propertyRepo.GetByID(in.PropertyID)
	if err != nil {
		s.lg.Warn("property not found", zap.String("property_id", in.PropertyID.String()))
		return nil, &rerrors.ValidationError{Msg: "Property not found"}
	}

	// If this is set as main, ensure no other media is main for this property
	if in.IsMain != nil && *in.IsMain {
		// This will be handled by the repository's SetAsMain method
		// For now, we'll create the media and then update it
	}

	// Create media model
	media := &models.PropertyMedia{
		ID:            uuid.New(),
		PropertyID:    in.PropertyID,
		MediaType:     in.MediaType,
		FileName:      in.FileName,
		FileURL:       in.FileURL,
		FileSizeBytes: in.FileSizeBytes,
		MimeType:      in.MimeType,
		Title:         in.Title,
		Description:   in.Description,
		AltText:       in.AltText,
		DisplayOrder:  1, // Default display order
		IsMain:        in.IsMain,
		IsPublic:      in.IsPublic,
		WidthPx:       in.WidthPx,
		HeightPx:      in.HeightPx,
		CreatedAt:     time.Now(),
		CreatedBy:     createdBy,
	}

	// Set display order if provided
	if in.DisplayOrder != nil {
		media.DisplayOrder = *in.DisplayOrder
	} else {
		// Get next display order for this property
		existingMedia, err := s.propertyMediaRepo.GetByPropertyID(in.PropertyID)
		if err == nil && len(existingMedia) > 0 {
			maxOrder := int32(0)
			for _, existing := range existingMedia {
				if existing.DisplayOrder > maxOrder {
					maxOrder = existing.DisplayOrder
				}
			}
			media.DisplayOrder = maxOrder + 1
		}
	}

	// Create media
	createdMedia, err := s.propertyMediaRepo.Create(media)
	if err != nil {
		s.lg.Error("failed to create property media", zap.Error(err))
		return nil, err
	}

	// If this should be the main media, update accordingly
	if in.IsMain != nil && *in.IsMain {
		err = s.propertyMediaRepo.SetAsMain(createdMedia.ID, in.PropertyID, *createdBy)
		if err != nil {
			s.lg.Error("failed to set media as main", zap.String("id", createdMedia.ID.String()), zap.Error(err))
			// Continue anyway, the media was created successfully
		}
	}

	// Convert to response DTO
	response := s.convertToDTO(createdMedia)

	s.lg.Info("property media created successfully", zap.String("id", createdMedia.ID.String()))
	return response, nil
}

func (s *PropertyMediaService) GetByID(id uuid.UUID) (*dto.PropertyMediaResponse, error) {
	media, err := s.propertyMediaRepo.GetByID(id)
	if err != nil {
		s.lg.Warn("property media not found", zap.String("id", id.String()))
		return nil, err
	}

	response := s.convertToDTO(media)
	return response, nil
}

func (s *PropertyMediaService) GetByPropertyID(propertyID uuid.UUID) ([]dto.PropertyMediaResponse, error) {
	mediaList, err := s.propertyMediaRepo.GetByPropertyID(propertyID)
	if err != nil {
		s.lg.Error("failed to get media by property ID", zap.String("property_id", propertyID.String()))
		return nil, err
	}

	responses := make([]dto.PropertyMediaResponse, len(mediaList))
	for i, media := range mediaList {
		responses[i] = *s.convertToDTO(&media)
	}

	return responses, nil
}

func (s *PropertyMediaService) GetByPropertyIDAndType(propertyID uuid.UUID, mediaType models.MediaType) ([]dto.PropertyMediaResponse, error) {
	mediaList, err := s.propertyMediaRepo.GetByPropertyIDAndType(propertyID, mediaType)
	if err != nil {
		s.lg.Error("failed to get media by property ID and type",
			zap.String("property_id", propertyID.String()),
			zap.String("media_type", string(mediaType)))
		return nil, err
	}

	responses := make([]dto.PropertyMediaResponse, len(mediaList))
	for i, media := range mediaList {
		responses[i] = *s.convertToDTO(&media)
	}

	return responses, nil
}

func (s *PropertyMediaService) GetMainImage(propertyID uuid.UUID) (*dto.PropertyMediaResponse, error) {
	media, err := s.propertyMediaRepo.GetMainImage(propertyID)
	if err != nil {
		s.lg.Warn("main image not found for property", zap.String("property_id", propertyID.String()))
		return nil, err
	}

	response := s.convertToDTO(media)
	return response, nil
}

func (s *PropertyMediaService) List(filters dto.PropertyMediaFilter, page, limit int) ([]dto.PropertyMediaResponse, int64, error) {
	offset := (page - 1) * limit

	repoFilters := repository.PropertyMediaFilters{
		PropertyID:     filters.PropertyID,
		MediaType:      filters.MediaType,
		IsMain:         filters.IsMain,
		IsPublic:       filters.IsPublic,
		IncludeDeleted: false,
	}

	mediaList, total, err := s.propertyMediaRepo.List(repoFilters, limit, offset)
	if err != nil {
		s.lg.Error("failed to list property media", zap.Error(err))
		return nil, 0, err
	}

	responses := make([]dto.PropertyMediaResponse, len(mediaList))
	for i, media := range mediaList {
		responses[i] = *s.convertToDTO(&media)
	}

	return responses, total, nil
}

func (s *PropertyMediaService) Update(id uuid.UUID, in dto.PropertyMediaUpdate, updatedBy *uuid.UUID) (*dto.PropertyMediaResponse, error) {
	s.lg.Debug("update property media request", zap.String("id", id.String()), zap.Any("payload", in))

	// Get existing media
	media, err := s.propertyMediaRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields
	now := time.Now()
	if in.MediaType != nil {
		media.MediaType = *in.MediaType
	}
	if in.FileName != nil {
		media.FileName = *in.FileName
	}
	if in.FileURL != nil {
		media.FileURL = *in.FileURL
	}
	if in.FileSizeBytes != nil {
		media.FileSizeBytes = in.FileSizeBytes
	}
	if in.MimeType != nil {
		media.MimeType = in.MimeType
	}
	if in.Title != nil {
		media.Title = in.Title
	}
	if in.Description != nil {
		media.Description = in.Description
	}
	if in.AltText != nil {
		media.AltText = in.AltText
	}
	if in.DisplayOrder != nil {
		media.DisplayOrder = *in.DisplayOrder
	}
	if in.IsMain != nil {
		media.IsMain = in.IsMain
	}
	if in.IsPublic != nil {
		media.IsPublic = in.IsPublic
	}
	if in.WidthPx != nil {
		media.WidthPx = in.WidthPx
	}
	if in.HeightPx != nil {
		media.HeightPx = in.HeightPx
	}

	media.UpdatedAt = &now
	media.UpdatedBy = updatedBy

	// Update media
	err = s.propertyMediaRepo.Update(media)
	if err != nil {
		s.lg.Error("failed to update property media", zap.String("id", id.String()), zap.Error(err))
		return nil, err
	}

	// If this should be the main media, update accordingly
	if in.IsMain != nil && *in.IsMain && updatedBy != nil {
		err = s.propertyMediaRepo.SetAsMain(media.ID, media.PropertyID, *updatedBy)
		if err != nil {
			s.lg.Error("failed to set media as main", zap.String("id", media.ID.String()), zap.Error(err))
			// Continue anyway, the media was updated successfully
		}
	}

	response := s.convertToDTO(media)

	s.lg.Info("property media updated successfully", zap.String("id", id.String()))
	return response, nil
}

func (s *PropertyMediaService) Delete(id uuid.UUID, deletedBy *uuid.UUID) error {
	s.lg.Debug("delete property media request", zap.String("id", id.String()))

	// Check if media exists
	_, err := s.propertyMediaRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Soft delete media
	if deletedBy != nil {
		err = s.propertyMediaRepo.SoftDelete(id, *deletedBy)
	} else {
		err = s.propertyMediaRepo.Delete(id)
	}

	if err != nil {
		s.lg.Error("failed to delete property media", zap.String("id", id.String()), zap.Error(err))
		return err
	}

	s.lg.Info("property media deleted successfully", zap.String("id", id.String()))
	return nil
}

func (s *PropertyMediaService) SetAsMain(id uuid.UUID, propertyID uuid.UUID, updatedBy uuid.UUID) error {
	s.lg.Debug("set media as main", zap.String("id", id.String()), zap.String("property_id", propertyID.String()))

	err := s.propertyMediaRepo.SetAsMain(id, propertyID, updatedBy)
	if err != nil {
		s.lg.Error("failed to set media as main", zap.String("id", id.String()), zap.Error(err))
		return err
	}

	s.lg.Info("media set as main successfully", zap.String("id", id.String()))
	return nil
}

func (s *PropertyMediaService) UpdateDisplayOrder(propertyID uuid.UUID, mediaOrder []uuid.UUID, updatedBy uuid.UUID) error {
	s.lg.Debug("update media display order",
		zap.String("property_id", propertyID.String()),
		zap.Int("count", len(mediaOrder)))

	err := s.propertyMediaRepo.UpdateDisplayOrder(propertyID, mediaOrder, updatedBy)
	if err != nil {
		s.lg.Error("failed to update media display order",
			zap.String("property_id", propertyID.String()),
			zap.Error(err))
		return err
	}

	s.lg.Info("media display order updated successfully",
		zap.String("property_id", propertyID.String()),
		zap.Int("count", len(mediaOrder)))
	return nil
}

func (s *PropertyMediaService) BulkCreateForProperty(propertyID uuid.UUID, mediaList []dto.PropertyMediaCreate, createdBy *uuid.UUID) ([]dto.PropertyMediaResponse, error) {
	s.lg.Debug("bulk create media for property",
		zap.String("property_id", propertyID.String()),
		zap.Int("count", len(mediaList)))

	// Validate property exists
	_, err := s.propertyRepo.GetByID(propertyID)
	if err != nil {
		s.lg.Warn("property not found", zap.String("property_id", propertyID.String()))
		return nil, &rerrors.ValidationError{Msg: "Property not found"}
	}

	// Convert DTOs to models
	mediaModels := make([]models.PropertyMedia, len(mediaList))
	for i, mediaDTO := range mediaList {
		mediaModels[i] = models.PropertyMedia{
			ID:            uuid.New(),
			PropertyID:    propertyID,
			MediaType:     mediaDTO.MediaType,
			FileName:      mediaDTO.FileName,
			FileURL:       mediaDTO.FileURL,
			FileSizeBytes: mediaDTO.FileSizeBytes,
			MimeType:      mediaDTO.MimeType,
			Title:         mediaDTO.Title,
			Description:   mediaDTO.Description,
			AltText:       mediaDTO.AltText,
			DisplayOrder:  int32(i + 1), // Sequential display order
			IsMain:        mediaDTO.IsMain,
			IsPublic:      mediaDTO.IsPublic,
			WidthPx:       mediaDTO.WidthPx,
			HeightPx:      mediaDTO.HeightPx,
			CreatedAt:     time.Now(),
			CreatedBy:     createdBy,
		}

		if mediaDTO.DisplayOrder != nil {
			mediaModels[i].DisplayOrder = *mediaDTO.DisplayOrder
		}
	}

	// Bulk create
	err = s.propertyMediaRepo.BulkCreateForProperty(propertyID, mediaModels)
	if err != nil {
		s.lg.Error("failed to bulk create property media",
			zap.String("property_id", propertyID.String()),
			zap.Error(err))
		return nil, err
	}

	// Convert to response DTOs
	responses := make([]dto.PropertyMediaResponse, len(mediaModels))
	for i, media := range mediaModels {
		responses[i] = *s.convertToDTO(&media)
	}

	s.lg.Info("property media bulk created successfully",
		zap.String("property_id", propertyID.String()),
		zap.Int("count", len(mediaModels)))
	return responses, nil
}

func (s *PropertyMediaService) BulkDeleteForProperty(propertyID uuid.UUID) error {
	s.lg.Debug("bulk delete media for property", zap.String("property_id", propertyID.String()))

	err := s.propertyMediaRepo.BulkDeleteForProperty(propertyID)
	if err != nil {
		s.lg.Error("failed to bulk delete property media",
			zap.String("property_id", propertyID.String()),
			zap.Error(err))
		return err
	}

	s.lg.Info("property media bulk deleted successfully", zap.String("property_id", propertyID.String()))
	return nil
}

func (s *PropertyMediaService) convertToDTO(media *models.PropertyMedia) *dto.PropertyMediaResponse {
	return &dto.PropertyMediaResponse{
		ID:            media.ID,
		PropertyID:    media.PropertyID,
		MediaType:     media.MediaType,
		FileName:      media.FileName,
		FileURL:       media.FileURL,
		FileSizeBytes: media.FileSizeBytes,
		MimeType:      media.MimeType,
		Title:         media.Title,
		Description:   media.Description,
		AltText:       media.AltText,
		DisplayOrder:  media.DisplayOrder,
		IsMain:        media.IsMain,
		IsPublic:      media.IsPublic,
		WidthPx:       media.WidthPx,
		HeightPx:      media.HeightPx,
		CreatedAt:     media.CreatedAt,
		CreatedBy:     media.CreatedBy,
		UpdatedAt:     media.UpdatedAt,
		UpdatedBy:     media.UpdatedBy,
	}
}
