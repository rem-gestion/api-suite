package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/rem-gestion/api-suite/property/src/models"
	rerrors "github.com/rem-gestion/rem-common/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PropertyListingRepo struct {
	db *gorm.DB
	lg *zap.Logger
}

func NewPropertyListingRepo(db *gorm.DB, lg *zap.Logger) PropertyListingRepository {
	return &PropertyListingRepo{db: db, lg: lg.Named("property-listing-repo")}
}

func (r *PropertyListingRepo) Create(listing *models.PropertyListing) (*models.PropertyListing, error) {
	if err := r.db.Create(listing).Error; err != nil {
		r.lg.Error("failed to create property listing", zap.Error(err))
		return nil, err
	}
	r.lg.Info("property listing created successfully", zap.String("id", listing.ID.String()))
	return listing, nil
}

func (r *PropertyListingRepo) GetByID(id uuid.UUID) (*models.PropertyListing, error) {
	var listing models.PropertyListing
	err := r.db.Preload("Property").
		First(&listing, "id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &rerrors.NotFoundError{Msg: "property listing not found"}
		}
		r.lg.Error("failed to get property listing by ID", zap.String("id", id.String()), zap.Error(err))
		return nil, err
	}

	return &listing, nil
}

func (r *PropertyListingRepo) GetByPropertyID(propertyID uuid.UUID) ([]models.PropertyListing, error) {
	var listings []models.PropertyListing
	err := r.db.Where("property_id = ?", propertyID).
		Order("created_at DESC").
		Find(&listings).Error

	if err != nil {
		r.lg.Error("failed to get property listings by property ID", zap.String("property_id", propertyID.String()), zap.Error(err))
		return nil, err
	}

	return listings, nil
}

func (r *PropertyListingRepo) GetActiveByPropertyID(propertyID uuid.UUID) (*models.PropertyListing, error) {
	var listing models.PropertyListing
	err := r.db.Where("property_id = ? AND listing_status = ?", propertyID, models.ListingStatusActive).
		First(&listing).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &rerrors.NotFoundError{Msg: "active property listing not found"}
		}
		r.lg.Error("failed to get active property listing", zap.String("property_id", propertyID.String()), zap.Error(err))
		return nil, err
	}

	return &listing, nil
}

func (r *PropertyListingRepo) GetActiveListingsByType(operationType models.OperationType, limit, offset int) ([]models.PropertyListing, int64, error) {
	var listings []models.PropertyListing
	var total int64

	query := r.db.Model(&models.PropertyListing{}).
		Where("operation_type = ? AND listing_status = ?", operationType, models.ListingStatusActive)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		r.lg.Error("failed to count active listings by type",
			zap.String("operation_type", string(operationType)),
			zap.Error(err))
		return nil, 0, err
	}

	// Get records with preloading
	err := query.Preload("Property").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&listings).Error

	if err != nil {
		r.lg.Error("failed to get active listings by type",
			zap.String("operation_type", string(operationType)),
			zap.Error(err))
		return nil, 0, err
	}

	return listings, total, nil
}

func (r *PropertyListingRepo) List(filters PropertyListingFilters, limit, offset int) ([]models.PropertyListing, int64, error) {
	var listings []models.PropertyListing
	var total int64

	query := r.db.Model(&models.PropertyListing{})

	// Apply filters
	if filters.PropertyID != nil {
		query = query.Where("property_id = ?", *filters.PropertyID)
	}
	if filters.OperationType != nil {
		query = query.Where("operation_type = ?", *filters.OperationType)
	}
	if filters.ListingStatus != nil {
		query = query.Where("listing_status = ?", *filters.ListingStatus)
	}
	if filters.MinPrice != nil {
		query = query.Where("price >= ?", *filters.MinPrice)
	}
	if filters.MaxPrice != nil {
		query = query.Where("price <= ?", *filters.MaxPrice)
	}
	if filters.Currency != nil {
		query = query.Where("currency = ?", *filters.Currency)
	}
	if filters.Highlighted != nil {
		query = query.Where("highlighted = ?", *filters.Highlighted)
	}
	if filters.PublishOnWeb != nil {
		query = query.Where("publish_on_web = ?", *filters.PublishOnWeb)
	}
	if filters.PublishOnPortals != nil {
		query = query.Where("publish_on_portals = ?", *filters.PublishOnPortals)
	}
	if filters.CreatedBy != nil {
		query = query.Where("created_by = ?", *filters.CreatedBy)
	}

	if !filters.IncludeDeleted {
		query = query.Where("deleted_at IS NULL")
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		r.lg.Error("failed to count property listings", zap.Error(err))
		return nil, 0, err
	}

	// Get records with preloading
	err := query.Preload("Property").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&listings).Error

	if err != nil {
		r.lg.Error("failed to list property listings", zap.Error(err))
		return nil, 0, err
	}

	return listings, total, nil
}

func (r *PropertyListingRepo) Update(listing *models.PropertyListing) error {
	if err := r.db.Save(listing).Error; err != nil {
		r.lg.Error("failed to update property listing", zap.String("id", listing.ID.String()), zap.Error(err))
		return err
	}
	r.lg.Info("property listing updated successfully", zap.String("id", listing.ID.String()))
	return nil
}

func (r *PropertyListingRepo) Delete(id uuid.UUID) error {
	if err := r.db.Delete(&models.PropertyListing{}, "id = ?", id).Error; err != nil {
		r.lg.Error("failed to delete property listing", zap.String("id", id.String()), zap.Error(err))
		return err
	}
	r.lg.Info("property listing deleted successfully", zap.String("id", id.String()))
	return nil
}

func (r *PropertyListingRepo) SoftDelete(id uuid.UUID, deletedBy uuid.UUID) error {
	now := time.Now()
	err := r.db.Model(&models.PropertyListing{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"deleted_at": &now,
			"updated_by": &deletedBy,
		}).Error

	if err != nil {
		r.lg.Error("failed to soft delete property listing", zap.String("id", id.String()), zap.Error(err))
		return err
	}

	r.lg.Info("property listing soft deleted successfully", zap.String("id", id.String()))
	return nil
}

func (r *PropertyListingRepo) UpdateStatus(id uuid.UUID, status models.ListingStatus, updatedBy uuid.UUID) error {
	now := time.Now()
	err := r.db.Model(&models.PropertyListing{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"listing_status": status,
			"updated_at":     &now,
			"updated_by":     &updatedBy,
		}).Error

	if err != nil {
		r.lg.Error("failed to update property listing status",
			zap.String("id", id.String()),
			zap.String("status", string(status)),
			zap.Error(err))
		return err
	}

	r.lg.Info("property listing status updated successfully",
		zap.String("id", id.String()),
		zap.String("status", string(status)))
	return nil
}

func (r *PropertyListingRepo) IncrementViews(id uuid.UUID) error {
	err := r.db.Model(&models.PropertyListing{}).
		Where("id = ?", id).
		Update("views_count", gorm.Expr("views_count + 1")).Error

	if err != nil {
		r.lg.Error("failed to increment property listing views", zap.String("id", id.String()), zap.Error(err))
		return err
	}

	return nil
}

func (r *PropertyListingRepo) IncrementInquiries(id uuid.UUID) error {
	err := r.db.Model(&models.PropertyListing{}).
		Where("id = ?", id).
		Update("inquiries_count", gorm.Expr("inquiries_count + 1")).Error

	if err != nil {
		r.lg.Error("failed to increment property listing inquiries", zap.String("id", id.String()), zap.Error(err))
		return err
	}

	return nil
}

func (r *PropertyListingRepo) IncrementFavorites(id uuid.UUID) error {
	err := r.db.Model(&models.PropertyListing{}).
		Where("id = ?", id).
		Update("favorites_count", gorm.Expr("favorites_count + 1")).Error

	if err != nil {
		r.lg.Error("failed to increment property listing favorites", zap.String("id", id.String()), zap.Error(err))
		return err
	}

	return nil
}

func (r *PropertyListingRepo) DecrementFavorites(id uuid.UUID) error {
	err := r.db.Model(&models.PropertyListing{}).
		Where("id = ? AND favorites_count > 0", id).
		Update("favorites_count", gorm.Expr("favorites_count - 1")).Error

	if err != nil {
		r.lg.Error("failed to decrement property listing favorites", zap.String("id", id.String()), zap.Error(err))
		return err
	}

	return nil
}

func (r *PropertyListingRepo) Search(query string, filters PropertyListingFilters, limit, offset int) ([]models.PropertyListing, int64, error) {
	var listings []models.PropertyListing
	var total int64

	dbQuery := r.db.Model(&models.PropertyListing{})

	// Full-text search on title and description
	if query != "" {
		searchQuery := "%" + query + "%"
		dbQuery = dbQuery.Where("title ILIKE ? OR description ILIKE ?", searchQuery, searchQuery)
	}

	// Apply additional filters (same as List method)
	if filters.PropertyID != nil {
		dbQuery = dbQuery.Where("property_id = ?", *filters.PropertyID)
	}
	if filters.OperationType != nil {
		dbQuery = dbQuery.Where("operation_type = ?", *filters.OperationType)
	}
	if filters.ListingStatus != nil {
		dbQuery = dbQuery.Where("listing_status = ?", *filters.ListingStatus)
	}
	if filters.MinPrice != nil {
		dbQuery = dbQuery.Where("price >= ?", *filters.MinPrice)
	}
	if filters.MaxPrice != nil {
		dbQuery = dbQuery.Where("price <= ?", *filters.MaxPrice)
	}
	if filters.Currency != nil {
		dbQuery = dbQuery.Where("currency = ?", *filters.Currency)
	}
	if filters.Highlighted != nil {
		dbQuery = dbQuery.Where("highlighted = ?", *filters.Highlighted)
	}
	if filters.PublishOnWeb != nil {
		dbQuery = dbQuery.Where("publish_on_web = ?", *filters.PublishOnWeb)
	}
	if filters.PublishOnPortals != nil {
		dbQuery = dbQuery.Where("publish_on_portals = ?", *filters.PublishOnPortals)
	}

	// Count total
	if err := dbQuery.Count(&total).Error; err != nil {
		r.lg.Error("failed to count search results", zap.String("query", query), zap.Error(err))
		return nil, 0, err
	}

	// Get records
	err := dbQuery.Preload("Property").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&listings).Error

	if err != nil {
		r.lg.Error("failed to search property listings", zap.String("query", query), zap.Error(err))
		return nil, 0, err
	}

	return listings, total, nil
}
