package repository

import (
	"github.com/google/uuid"
	"github.com/rem-gestion/api-suite/property/src/models"
	rerrors "github.com/rem-gestion/rem-common/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PropertyRepo struct {
	db *gorm.DB
	lg *zap.Logger
}

func NewPropertyRepo(db *gorm.DB, lg *zap.Logger) PropertyRepository {
	return &PropertyRepo{db: db, lg: lg.Named("property-repo")}
}

func (r *PropertyRepo) Create(property *models.Property) (*models.Property, error) {
	if err := r.db.Create(property).Error; err != nil {
		r.lg.Error("failed to create property", zap.Error(err))
		return nil, err
	}
	r.lg.Info("property created successfully", zap.String("id", property.ID.String()))
	return property, nil
}

func (r *PropertyRepo) GetByID(id uuid.UUID) (*models.Property, error) {
	var property models.Property
	err := r.db.Preload("PropertyType").
		Preload("PropertyManagements").
		Preload("PropertyAmenities").
		First(&property, "id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &rerrors.NotFoundError{Msg: "property not found"}
		}
		r.lg.Error("failed to get property by ID", zap.String("id", id.String()), zap.Error(err))
		return nil, err
	}

	return &property, nil
}

func (r *PropertyRepo) GetByInternalCode(code string) (*models.Property, error) {
	var property models.Property
	err := r.db.Preload("PropertyType").
		Preload("PropertyManagements").
		Preload("PropertyAmenities").
		First(&property, "internal_code = ?", code).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &rerrors.NotFoundError{Msg: "property not found"}
		}
		r.lg.Error("failed to get property by internal code", zap.String("code", code), zap.Error(err))
		return nil, err
	}

	return &property, nil
}

func (r *PropertyRepo) List(filters PropertyFilters, limit, offset int) ([]models.Property, int64, error) {
	var properties []models.Property
	var total int64

	query := r.db.Model(&models.Property{})

	// Apply filters
	if filters.OwnerPersonID != nil {
		query = query.Where("owner_person_id = ?", *filters.OwnerPersonID)
	}
	if filters.PropertyTypeID != nil {
		query = query.Where("property_type_id = ?", *filters.PropertyTypeID)
	}
	if filters.InternalCode != nil {
		query = query.Where("internal_code ILIKE ?", "%"+*filters.InternalCode+"%")
	}
	if filters.YearBuilt != nil {
		query = query.Where("year_built = ?", *filters.YearBuilt)
	}
	if filters.Bedrooms != nil {
		query = query.Where("bedrooms = ?", *filters.Bedrooms)
	}
	if filters.Bathrooms != nil {
		query = query.Where("bathrooms = ?", *filters.Bathrooms)
	}
	if filters.MinTotalArea != nil {
		query = query.Where("total_area_sqm >= ?", *filters.MinTotalArea)
	}
	if filters.MaxTotalArea != nil {
		query = query.Where("total_area_sqm <= ?", *filters.MaxTotalArea)
	}

	if !filters.IncludeDeleted {
		query = query.Where("deleted_at IS NULL")
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		r.lg.Error("failed to count properties", zap.Error(err))
		return nil, 0, err
	}

	// Get records with preloading
	err := query.Preload("PropertyType").
		Preload("PropertyManagements").
		Preload("PropertyAmenities").
		Limit(limit).Offset(offset).
		Find(&properties).Error

	if err != nil {
		r.lg.Error("failed to list properties", zap.Error(err))
		return nil, 0, err
	}

	return properties, total, nil
}

func (r *PropertyRepo) Update(property *models.Property) error {
	if err := r.db.Save(property).Error; err != nil {
		r.lg.Error("failed to update property", zap.String("id", property.ID.String()), zap.Error(err))
		return err
	}
	r.lg.Info("property updated successfully", zap.String("id", property.ID.String()))
	return nil
}

func (r *PropertyRepo) SoftDelete(id uuid.UUID, deletedBy uuid.UUID) error {
	err := r.db.Model(&models.Property{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"deleted_at": gorm.DeletedAt{},
			"updated_by": deletedBy,
		}).Error

	if err != nil {
		r.lg.Error("failed to soft delete property", zap.String("id", id.String()), zap.Error(err))
		return err
	}

	r.lg.Info("property soft deleted successfully", zap.String("id", id.String()))
	return nil
}

func (r *PropertyRepo) Delete(id uuid.UUID) error {
	if err := r.db.Unscoped().Delete(&models.Property{}, "id = ?", id).Error; err != nil {
		r.lg.Error("failed to delete property", zap.String("id", id.String()), zap.Error(err))
		return err
	}
	r.lg.Info("property deleted permanently", zap.String("id", id.String()))
	return nil
}

func (r *PropertyRepo) Search(query string, filters PropertyFilters, limit, offset int) ([]models.Property, int64, error) {
	var properties []models.Property
	var total int64

	dbQuery := r.db.Model(&models.Property{})

	// Full-text search on description and internal_code
	searchQuery := "%" + query + "%"
	dbQuery = dbQuery.Where("description ILIKE ? OR internal_code ILIKE ?", searchQuery, searchQuery)

	// Apply additional filters
	if filters.OwnerPersonID != nil {
		dbQuery = dbQuery.Where("owner_person_id = ?", *filters.OwnerPersonID)
	}
	if filters.PropertyTypeID != nil {
		dbQuery = dbQuery.Where("property_type_id = ?", *filters.PropertyTypeID)
	}
	if !filters.IncludeDeleted {
		dbQuery = dbQuery.Where("deleted_at IS NULL")
	}

	// Count total
	if err := dbQuery.Count(&total).Error; err != nil {
		r.lg.Error("failed to count search results", zap.String("query", query), zap.Error(err))
		return nil, 0, err
	}

	// Get records
	err := dbQuery.Preload("PropertyType").
		Preload("PropertyManagements").
		Preload("PropertyAmenities").
		Limit(limit).Offset(offset).
		Find(&properties).Error

	if err != nil {
		r.lg.Error("failed to search properties", zap.String("query", query), zap.Error(err))
		return nil, 0, err
	}

	return properties, total, nil
}

func (r *PropertyRepo) BulkCreate(properties []models.Property) ([]models.Property, []BulkError, error) {
	var created []models.Property
	var errors []BulkError

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for i, property := range properties {
		if err := tx.Create(&property).Error; err != nil {
			errors = append(errors, BulkError{
				Index: i,
				Error: err.Error(),
			})
			continue
		}
		created = append(created, property)
	}

	if err := tx.Commit().Error; err != nil {
		r.lg.Error("failed to commit bulk create transaction", zap.Error(err))
		return nil, nil, err
	}

	r.lg.Info("bulk create completed",
		zap.Int("total", len(properties)),
		zap.Int("created", len(created)),
		zap.Int("errors", len(errors)))

	return created, errors, nil
}
