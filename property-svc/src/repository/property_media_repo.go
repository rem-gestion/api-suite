package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/rem-gestion/api-suite/property/src/models"
	rerrors "github.com/rem-gestion/rem-common/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PropertyMediaRepo struct {
	db *gorm.DB
	lg *zap.Logger
}

func NewPropertyMediaRepo(db *gorm.DB, lg *zap.Logger) PropertyMediaRepository {
	return &PropertyMediaRepo{db: db, lg: lg.Named("property-media-repo")}
}

func (r *PropertyMediaRepo) Create(media *models.PropertyMedia) (*models.PropertyMedia, error) {
	if err := r.db.Create(media).Error; err != nil {
		r.lg.Error("failed to create property media", zap.Error(err))
		return nil, err
	}
	r.lg.Info("property media created successfully", zap.String("id", media.ID.String()))
	return media, nil
}

func (r *PropertyMediaRepo) GetByID(id uuid.UUID) (*models.PropertyMedia, error) {
	var media models.PropertyMedia
	err := r.db.Preload("Property").
		First(&media, "id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &rerrors.NotFoundError{Msg: "property media not found"}
		}
		r.lg.Error("failed to get property media by ID", zap.String("id", id.String()), zap.Error(err))
		return nil, err
	}

	return &media, nil
}

func (r *PropertyMediaRepo) GetByPropertyID(propertyID uuid.UUID) ([]models.PropertyMedia, error) {
	var medias []models.PropertyMedia
	err := r.db.Where("property_id = ?", propertyID).
		Order("display_order ASC, created_at ASC").
		Find(&medias).Error

	if err != nil {
		r.lg.Error("failed to get property media by property ID", zap.String("property_id", propertyID.String()), zap.Error(err))
		return nil, err
	}

	return medias, nil
}

func (r *PropertyMediaRepo) GetByPropertyIDAndType(propertyID uuid.UUID, mediaType models.MediaType) ([]models.PropertyMedia, error) {
	var medias []models.PropertyMedia
	err := r.db.Where("property_id = ? AND media_type = ?", propertyID, mediaType).
		Order("display_order ASC, created_at ASC").
		Find(&medias).Error

	if err != nil {
		r.lg.Error("failed to get property media by property ID and type",
			zap.String("property_id", propertyID.String()),
			zap.String("media_type", string(mediaType)),
			zap.Error(err))
		return nil, err
	}

	return medias, nil
}

func (r *PropertyMediaRepo) GetMainByPropertyID(propertyID uuid.UUID) (*models.PropertyMedia, error) {
	var media models.PropertyMedia
	err := r.db.Where("property_id = ? AND is_main = true", propertyID).
		First(&media).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &rerrors.NotFoundError{Msg: "main property media not found"}
		}
		r.lg.Error("failed to get main property media", zap.String("property_id", propertyID.String()), zap.Error(err))
		return nil, err
	}

	return &media, nil
}

func (r *PropertyMediaRepo) GetMainImage(propertyID uuid.UUID) (*models.PropertyMedia, error) {
	var media models.PropertyMedia
	err := r.db.Where("property_id = ? AND media_type = ? AND (is_main = true OR is_main IS NULL)",
		propertyID, models.MediaTypePhoto).
		Order("is_main DESC, display_order ASC, created_at ASC").
		First(&media).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &rerrors.NotFoundError{Msg: "main property image not found"}
		}
		r.lg.Error("failed to get main property image", zap.String("property_id", propertyID.String()), zap.Error(err))
		return nil, err
	}

	return &media, nil
}

func (r *PropertyMediaRepo) List(filters PropertyMediaFilters, limit, offset int) ([]models.PropertyMedia, int64, error) {
	var medias []models.PropertyMedia
	var total int64

	query := r.db.Model(&models.PropertyMedia{})

	// Apply filters
	if filters.PropertyID != nil {
		query = query.Where("property_id = ?", *filters.PropertyID)
	}
	if filters.MediaType != nil {
		query = query.Where("media_type = ?", *filters.MediaType)
	}
	if filters.IsMain != nil {
		query = query.Where("is_main = ?", *filters.IsMain)
	}
	if filters.IsPublic != nil {
		query = query.Where("is_public = ?", *filters.IsPublic)
	}
	if filters.CreatedBy != nil {
		query = query.Where("created_by = ?", *filters.CreatedBy)
	}

	if !filters.IncludeDeleted {
		query = query.Where("deleted_at IS NULL")
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		r.lg.Error("failed to count property media", zap.Error(err))
		return nil, 0, err
	}

	// Get records with preloading
	err := query.Preload("Property").
		Order("display_order ASC, created_at ASC").
		Limit(limit).Offset(offset).
		Find(&medias).Error

	if err != nil {
		r.lg.Error("failed to list property media", zap.Error(err))
		return nil, 0, err
	}

	return medias, total, nil
}

func (r *PropertyMediaRepo) Update(media *models.PropertyMedia) error {
	if err := r.db.Save(media).Error; err != nil {
		r.lg.Error("failed to update property media", zap.String("id", media.ID.String()), zap.Error(err))
		return err
	}
	r.lg.Info("property media updated successfully", zap.String("id", media.ID.String()))
	return nil
}

func (r *PropertyMediaRepo) Delete(id uuid.UUID) error {
	if err := r.db.Delete(&models.PropertyMedia{}, "id = ?", id).Error; err != nil {
		r.lg.Error("failed to delete property media", zap.String("id", id.String()), zap.Error(err))
		return err
	}
	r.lg.Info("property media deleted successfully", zap.String("id", id.String()))
	return nil
}

func (r *PropertyMediaRepo) SoftDelete(id uuid.UUID, deletedBy uuid.UUID) error {
	now := time.Now()
	err := r.db.Model(&models.PropertyMedia{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"deleted_at": &now,
			"updated_by": &deletedBy,
		}).Error

	if err != nil {
		r.lg.Error("failed to soft delete property media", zap.String("id", id.String()), zap.Error(err))
		return err
	}

	r.lg.Info("property media soft deleted successfully", zap.String("id", id.String()))
	return nil
}

func (r *PropertyMediaRepo) UpdateDisplayOrder(propertyID uuid.UUID, mediaOrder []uuid.UUID, updatedBy uuid.UUID) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	now := time.Now()

	for i, mediaID := range mediaOrder {
		err := tx.Model(&models.PropertyMedia{}).
			Where("id = ? AND property_id = ?", mediaID, propertyID).
			Updates(map[string]interface{}{
				"display_order": int32(i + 1), // 1-indexed display order
				"updated_at":    &now,
				"updated_by":    &updatedBy,
			}).Error

		if err != nil {
			tx.Rollback()
			r.lg.Error("failed to update media display order",
				zap.String("id", mediaID.String()),
				zap.Int("display_order", i+1),
				zap.Error(err))
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		r.lg.Error("failed to commit update display order transaction", zap.Error(err))
		return err
	}

	r.lg.Info("property media display order updated successfully",
		zap.String("property_id", propertyID.String()),
		zap.Int("count", len(mediaOrder)))
	return nil
}

func (r *PropertyMediaRepo) SetAsMain(id uuid.UUID, propertyID uuid.UUID, updatedBy uuid.UUID) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	now := time.Now()

	// First, unset any existing main media for this property
	err := tx.Model(&models.PropertyMedia{}).
		Where("property_id = ? AND is_main = true AND id != ?", propertyID, id).
		Updates(map[string]interface{}{
			"is_main":    false,
			"updated_at": &now,
			"updated_by": &updatedBy,
		}).Error

	if err != nil {
		tx.Rollback()
		r.lg.Error("failed to unset existing main media",
			zap.String("property_id", propertyID.String()),
			zap.Error(err))
		return err
	}

	// Then set the new main media
	err = tx.Model(&models.PropertyMedia{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_main":    true,
			"updated_at": &now,
			"updated_by": &updatedBy,
		}).Error

	if err != nil {
		tx.Rollback()
		r.lg.Error("failed to set new main media",
			zap.String("id", id.String()),
			zap.Error(err))
		return err
	}

	if err := tx.Commit().Error; err != nil {
		r.lg.Error("failed to commit set main media transaction", zap.Error(err))
		return err
	}

	r.lg.Info("property media set as main successfully",
		zap.String("id", id.String()),
		zap.String("property_id", propertyID.String()))
	return nil
}

func (r *PropertyMediaRepo) ReorderMedia(propertyID uuid.UUID, mediaOrders []struct {
	ID           uuid.UUID
	DisplayOrder int32
}, updatedBy uuid.UUID) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	now := time.Now()

	for _, order := range mediaOrders {
		err := tx.Model(&models.PropertyMedia{}).
			Where("id = ? AND property_id = ?", order.ID, propertyID).
			Updates(map[string]interface{}{
				"display_order": order.DisplayOrder,
				"updated_at":    &now,
				"updated_by":    &updatedBy,
			}).Error

		if err != nil {
			tx.Rollback()
			r.lg.Error("failed to update media display order",
				zap.String("id", order.ID.String()),
				zap.Int32("display_order", order.DisplayOrder),
				zap.Error(err))
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		r.lg.Error("failed to commit reorder media transaction", zap.Error(err))
		return err
	}

	r.lg.Info("property media reordered successfully",
		zap.String("property_id", propertyID.String()),
		zap.Int("count", len(mediaOrders)))
	return nil
}

func (r *PropertyMediaRepo) BulkCreateForProperty(propertyID uuid.UUID, media []models.PropertyMedia) error {
	if len(media) == 0 {
		return nil
	}

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Set property ID for all media items
	for i := range media {
		media[i].PropertyID = propertyID
	}

	// Create all media items
	if err := tx.Create(&media).Error; err != nil {
		tx.Rollback()
		r.lg.Error("failed to bulk create property media",
			zap.String("property_id", propertyID.String()),
			zap.Int("count", len(media)),
			zap.Error(err))
		return err
	}

	if err := tx.Commit().Error; err != nil {
		r.lg.Error("failed to commit bulk create media transaction", zap.Error(err))
		return err
	}

	r.lg.Info("property media bulk created successfully",
		zap.String("property_id", propertyID.String()),
		zap.Int("count", len(media)))
	return nil
}

func (r *PropertyMediaRepo) BulkDeleteForProperty(propertyID uuid.UUID) error {
	err := r.db.Delete(&models.PropertyMedia{}, "property_id = ?", propertyID).Error

	if err != nil {
		r.lg.Error("failed to bulk delete property media",
			zap.String("property_id", propertyID.String()),
			zap.Error(err))
		return err
	}

	r.lg.Info("property media bulk deleted successfully",
		zap.String("property_id", propertyID.String()))
	return nil
}
