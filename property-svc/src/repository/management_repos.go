package repository

import (
	"github.com/google/uuid"
	"github.com/rem-gestion/api-suite/property/src/models"
	rerrors "github.com/rem-gestion/rem-common/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// PropertyManagementRepo implementation
type PropertyManagementRepo struct {
	db *gorm.DB
	lg *zap.Logger
}

func NewPropertyManagementRepo(db *gorm.DB, lg *zap.Logger) PropertyManagementRepository {
	return &PropertyManagementRepo{db: db, lg: lg.Named("property-management-repo")}
}

func (r *PropertyManagementRepo) Create(management *models.PropertyManagement) (*models.PropertyManagement, error) {
	if err := r.db.Create(management).Error; err != nil {
		r.lg.Error("failed to create property management", zap.Error(err))
		return nil, err
	}
	r.lg.Info("property management created successfully", zap.String("id", management.ID.String()))
	return management, nil
}

func (r *PropertyManagementRepo) GetByID(id uuid.UUID) (*models.PropertyManagement, error) {
	var management models.PropertyManagement
	err := r.db.Preload("Property").
		Preload("ManagerType").
		First(&management, "id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &rerrors.NotFoundError{Msg: "property management not found"}
		}
		r.lg.Error("failed to get property management by ID", zap.String("id", id.String()), zap.Error(err))
		return nil, err
	}

	return &management, nil
}

func (r *PropertyManagementRepo) GetByPropertyID(propertyID uuid.UUID, active bool) ([]models.PropertyManagement, error) {
	var managements []models.PropertyManagement

	query := r.db.Preload("Property").
		Preload("ManagerType").
		Where("property_id = ?", propertyID)

	if active {
		query = query.Where("end_date IS NULL OR end_date > NOW()")
	}

	err := query.Find(&managements).Error
	if err != nil {
		r.lg.Error("failed to get property managements by property ID",
			zap.String("property_id", propertyID.String()), zap.Error(err))
		return nil, err
	}

	return managements, nil
}

func (r *PropertyManagementRepo) GetByManagerID(managerID uuid.UUID, active bool) ([]models.PropertyManagement, error) {
	var managements []models.PropertyManagement

	query := r.db.Preload("Property").
		Preload("ManagerType").
		Where("manager_id = ?", managerID)

	if active {
		query = query.Where("end_date IS NULL OR end_date > NOW()")
	}

	err := query.Find(&managements).Error
	if err != nil {
		r.lg.Error("failed to get property managements by manager ID",
			zap.String("manager_id", managerID.String()), zap.Error(err))
		return nil, err
	}

	return managements, nil
}

func (r *PropertyManagementRepo) List(filters PropertyManagementFilters, limit, offset int) ([]models.PropertyManagement, int64, error) {
	var managements []models.PropertyManagement
	var total int64

	query := r.db.Model(&models.PropertyManagement{})

	// Apply filters
	if filters.PropertyID != nil {
		query = query.Where("property_id = ?", *filters.PropertyID)
	}
	if filters.ManagerID != nil {
		query = query.Where("manager_id = ?", *filters.ManagerID)
	}
	if filters.ManagerTypeID != nil {
		query = query.Where("manager_type_id = ?", *filters.ManagerTypeID)
	}
	if filters.ActiveOnly {
		query = query.Where("end_date IS NULL OR end_date > NOW()")
	}
	if !filters.IncludeDeleted {
		query = query.Where("deleted_at IS NULL")
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		r.lg.Error("failed to count property managements", zap.Error(err))
		return nil, 0, err
	}

	// Get records with preloading
	err := query.Preload("Property").
		Preload("ManagerType").
		Limit(limit).Offset(offset).
		Find(&managements).Error

	if err != nil {
		r.lg.Error("failed to list property managements", zap.Error(err))
		return nil, 0, err
	}

	return managements, total, nil
}

func (r *PropertyManagementRepo) Update(management *models.PropertyManagement) error {
	if err := r.db.Save(management).Error; err != nil {
		r.lg.Error("failed to update property management",
			zap.String("id", management.ID.String()), zap.Error(err))
		return err
	}
	r.lg.Info("property management updated successfully", zap.String("id", management.ID.String()))
	return nil
}

func (r *PropertyManagementRepo) SoftDelete(id uuid.UUID, deletedBy uuid.UUID) error {
	err := r.db.Model(&models.PropertyManagement{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"deleted_at": gorm.DeletedAt{},
			"updated_by": deletedBy,
		}).Error

	if err != nil {
		r.lg.Error("failed to soft delete property management", zap.String("id", id.String()), zap.Error(err))
		return err
	}

	r.lg.Info("property management soft deleted successfully", zap.String("id", id.String()))
	return nil
}

func (r *PropertyManagementRepo) Delete(id uuid.UUID) error {
	if err := r.db.Unscoped().Delete(&models.PropertyManagement{}, "id = ?", id).Error; err != nil {
		r.lg.Error("failed to delete property management", zap.String("id", id.String()), zap.Error(err))
		return err
	}
	r.lg.Info("property management deleted permanently", zap.String("id", id.String()))
	return nil
}

// PropertyAmenityRepo implementation
type PropertyAmenityRepo struct {
	db *gorm.DB
	lg *zap.Logger
}

func NewPropertyAmenityRepo(db *gorm.DB, lg *zap.Logger) PropertyAmenityRepository {
	return &PropertyAmenityRepo{db: db, lg: lg.Named("property-amenity-repo")}
}

func (r *PropertyAmenityRepo) Create(amenity *models.PropertyAmenity) (*models.PropertyAmenity, error) {
	if err := r.db.Create(amenity).Error; err != nil {
		r.lg.Error("failed to create property amenity", zap.Error(err))
		return nil, err
	}
	r.lg.Info("property amenity created successfully",
		zap.String("property_id", amenity.PropertyID.String()),
		zap.Int32("amenity_id", amenity.AmenityID))
	return amenity, nil
}

func (r *PropertyAmenityRepo) GetByPropertyAndAmenity(propertyID uuid.UUID, amenityID int32) (*models.PropertyAmenity, error) {
	var amenity models.PropertyAmenity
	err := r.db.Preload("Property").
		First(&amenity, "property_id = ? AND amenity_id = ?", propertyID, amenityID).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &rerrors.NotFoundError{Msg: "property amenity not found"}
		}
		r.lg.Error("failed to get property amenity",
			zap.String("property_id", propertyID.String()),
			zap.Int32("amenity_id", amenityID),
			zap.Error(err))
		return nil, err
	}

	return &amenity, nil
}

func (r *PropertyAmenityRepo) GetByPropertyID(propertyID uuid.UUID) ([]models.PropertyAmenity, error) {
	var amenities []models.PropertyAmenity

	err := r.db.Preload("Property").
		Where("property_id = ?", propertyID).
		Find(&amenities).Error

	if err != nil {
		r.lg.Error("failed to get property amenities by property ID",
			zap.String("property_id", propertyID.String()), zap.Error(err))
		return nil, err
	}

	return amenities, nil
}

func (r *PropertyAmenityRepo) Update(amenity *models.PropertyAmenity) error {
	if err := r.db.Save(amenity).Error; err != nil {
		r.lg.Error("failed to update property amenity",
			zap.String("property_id", amenity.PropertyID.String()),
			zap.Int32("amenity_id", amenity.AmenityID),
			zap.Error(err))
		return err
	}
	r.lg.Info("property amenity updated successfully",
		zap.String("property_id", amenity.PropertyID.String()),
		zap.Int32("amenity_id", amenity.AmenityID))
	return nil
}

func (r *PropertyAmenityRepo) Delete(propertyID uuid.UUID, amenityID int32) error {
	err := r.db.Where("property_id = ? AND amenity_id = ?", propertyID, amenityID).
		Delete(&models.PropertyAmenity{}).Error

	if err != nil {
		r.lg.Error("failed to delete property amenity",
			zap.String("property_id", propertyID.String()),
			zap.Int32("amenity_id", amenityID),
			zap.Error(err))
		return err
	}
	r.lg.Info("property amenity deleted successfully",
		zap.String("property_id", propertyID.String()),
		zap.Int32("amenity_id", amenityID))
	return nil
}

func (r *PropertyAmenityRepo) BulkCreateForProperty(propertyID uuid.UUID, amenities []models.PropertyAmenity) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, amenity := range amenities {
		amenity.PropertyID = propertyID // Ensure property ID is set
		if err := tx.Create(&amenity).Error; err != nil {
			tx.Rollback()
			r.lg.Error("failed to bulk create property amenities",
				zap.String("property_id", propertyID.String()), zap.Error(err))
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		r.lg.Error("failed to commit bulk create amenities transaction",
			zap.String("property_id", propertyID.String()), zap.Error(err))
		return err
	}

	r.lg.Info("bulk created property amenities successfully",
		zap.String("property_id", propertyID.String()),
		zap.Int("count", len(amenities)))
	return nil
}

func (r *PropertyAmenityRepo) BulkDeleteForProperty(propertyID uuid.UUID) error {
	err := r.db.Where("property_id = ?", propertyID).Delete(&models.PropertyAmenity{}).Error
	if err != nil {
		r.lg.Error("failed to bulk delete property amenities",
			zap.String("property_id", propertyID.String()), zap.Error(err))
		return err
	}

	r.lg.Info("bulk deleted property amenities successfully",
		zap.String("property_id", propertyID.String()))
	return nil
}

// ListByPropertyID gets all amenities for a property with full amenity details
func (r *PropertyAmenityRepo) ListByPropertyID(propertyID uuid.UUID) ([]*models.Amenity, error) {
	var amenities []*models.Amenity

	// Get amenities through proper JOIN with property_property_amenities
	err := r.db.Table("property_amenity as a").
		Select("a.id, a.name, a.category, a.icon_url").
		Joins("INNER JOIN property_property_amenities as pa ON CAST(pa.amenity_id AS INTEGER) = a.id").
		Where("pa.property_id = ?", propertyID).
		Scan(&amenities).Error

	if err != nil {
		r.lg.Error("failed to list amenities for property",
			zap.String("property_id", propertyID.String()), zap.Error(err))
		return nil, err
	}

	r.lg.Debug("amenities listed successfully",
		zap.String("property_id", propertyID.String()),
		zap.Int("count", len(amenities)))
	return amenities, nil
}

// AddAmenityToProperty adds an amenity to a property
func (r *PropertyAmenityRepo) AddAmenityToProperty(propertyID uuid.UUID, amenityID int32, note string) error {
	// First check if the relation already exists
	var existing models.PropertyAmenity
	err := r.db.Where("property_id = ? AND amenity_id = ?", propertyID, amenityID).
		First(&existing).Error

	if err == nil {
		// Relation already exists
		r.lg.Warn("amenity already associated with property",
			zap.String("property_id", propertyID.String()),
			zap.Int32("amenity_id", amenityID))
		return &rerrors.ConflictError{Msg: "amenity already associated with this property"}
	}

	// Create the relation
	propertyAmenity := &models.PropertyAmenity{
		PropertyID: propertyID,
		AmenityID:  amenityID,
		Note:       &note,
	}

	err = r.db.Create(propertyAmenity).Error
	if err != nil {
		r.lg.Error("failed to add amenity to property",
			zap.String("property_id", propertyID.String()),
			zap.Int32("amenity_id", amenityID),
			zap.Error(err))
		return err
	}

	r.lg.Info("amenity added to property successfully",
		zap.String("property_id", propertyID.String()),
		zap.Int32("amenity_id", amenityID))
	return nil
}

// RemoveAmenityFromProperty removes an amenity from a property
func (r *PropertyAmenityRepo) RemoveAmenityFromProperty(propertyID uuid.UUID, amenityID int32) error {
	result := r.db.Where("property_id = ? AND amenity_id = ?", propertyID, amenityID).
		Delete(&models.PropertyAmenity{})

	if result.Error != nil {
		r.lg.Error("failed to remove amenity from property",
			zap.String("property_id", propertyID.String()),
			zap.Int32("amenity_id", amenityID),
			zap.Error(result.Error))
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.lg.Warn("attempt to remove non-existent amenity from property",
			zap.String("property_id", propertyID.String()),
			zap.Int32("amenity_id", amenityID),
			zap.String("reason", "no rows affected - amenity was not associated with this property"),
			zap.Int64("rows_affected", result.RowsAffected))
		return &rerrors.NotFoundError{Msg: "amenity not associated with this property"}
	}

	r.lg.Info("amenity removed from property successfully",
		zap.String("property_id", propertyID.String()),
		zap.Int32("amenity_id", amenityID),
		zap.Int64("rows_affected", result.RowsAffected))
	return nil
}
