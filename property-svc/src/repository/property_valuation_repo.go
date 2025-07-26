package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/rem-gestion/api-suite/property/src/models"
	rerrors "github.com/rem-gestion/rem-common/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PropertyValuationRepo struct {
	db *gorm.DB
	lg *zap.Logger
}

func NewPropertyValuationRepo(db *gorm.DB, lg *zap.Logger) PropertyValuationRepository {
	return &PropertyValuationRepo{db: db, lg: lg.Named("property-valuation-repo")}
}

func (r *PropertyValuationRepo) Create(valuation *models.PropertyValuation) (*models.PropertyValuation, error) {
	if err := r.db.Create(valuation).Error; err != nil {
		r.lg.Error("failed to create property valuation", zap.Error(err))
		return nil, err
	}
	r.lg.Info("property valuation created successfully", zap.String("id", valuation.ID.String()))
	return valuation, nil
}

func (r *PropertyValuationRepo) GetByID(id uuid.UUID) (*models.PropertyValuation, error) {
	var valuation models.PropertyValuation
	err := r.db.Preload("Property").
		First(&valuation, "id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &rerrors.NotFoundError{Msg: "property valuation not found"}
		}
		r.lg.Error("failed to get property valuation by ID", zap.String("id", id.String()), zap.Error(err))
		return nil, err
	}

	return &valuation, nil
}

func (r *PropertyValuationRepo) GetByPropertyID(propertyID uuid.UUID) ([]models.PropertyValuation, error) {
	var valuations []models.PropertyValuation
	err := r.db.Where("property_id = ?", propertyID).
		Order("valuation_date DESC").
		Find(&valuations).Error

	if err != nil {
		r.lg.Error("failed to get property valuations by property ID", zap.String("property_id", propertyID.String()), zap.Error(err))
		return nil, err
	}

	return valuations, nil
}

func (r *PropertyValuationRepo) GetLatestByPropertyID(propertyID uuid.UUID) (*models.PropertyValuation, error) {
	var valuation models.PropertyValuation
	err := r.db.Where("property_id = ?", propertyID).
		Order("valuation_date DESC").
		First(&valuation).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &rerrors.NotFoundError{Msg: "property valuation not found"}
		}
		r.lg.Error("failed to get latest property valuation", zap.String("property_id", propertyID.String()), zap.Error(err))
		return nil, err
	}

	return &valuation, nil
}

func (r *PropertyValuationRepo) GetByPropertyIDAndType(propertyID uuid.UUID, valuationType string) ([]models.PropertyValuation, error) {
	var valuations []models.PropertyValuation
	err := r.db.Where("property_id = ? AND valuation_type = ?", propertyID, valuationType).
		Order("valuation_date DESC").
		Find(&valuations).Error

	if err != nil {
		r.lg.Error("failed to get property valuations by property ID and type",
			zap.String("property_id", propertyID.String()),
			zap.String("valuation_type", valuationType),
			zap.Error(err))
		return nil, err
	}

	return valuations, nil
}

func (r *PropertyValuationRepo) List(filters PropertyValuationFilters, limit, offset int) ([]models.PropertyValuation, int64, error) {
	var valuations []models.PropertyValuation
	var total int64

	query := r.db.Model(&models.PropertyValuation{})

	// Apply filters
	if filters.PropertyID != nil {
		query = query.Where("property_id = ?", *filters.PropertyID)
	}
	if filters.ValuationType != nil {
		query = query.Where("valuation_type = ?", *filters.ValuationType)
	}
	if filters.ValuationMethod != nil {
		query = query.Where("valuation_method = ?", *filters.ValuationMethod)
	}
	if filters.MinValue != nil {
		query = query.Where("appraised_value >= ?", *filters.MinValue)
	}
	if filters.MaxValue != nil {
		query = query.Where("appraised_value <= ?", *filters.MaxValue)
	}
	if filters.Currency != nil {
		query = query.Where("currency = ?", *filters.Currency)
	}
	if filters.ValidOnly {
		query = query.Where("valid_until IS NULL OR valid_until >= ?", time.Now())
	}
	if filters.CreatedBy != nil {
		query = query.Where("created_by = ?", *filters.CreatedBy)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		r.lg.Error("failed to count property valuations", zap.Error(err))
		return nil, 0, err
	}

	// Get records with preloading
	err := query.Preload("Property").
		Order("valuation_date DESC").
		Limit(limit).Offset(offset).
		Find(&valuations).Error

	if err != nil {
		r.lg.Error("failed to list property valuations", zap.Error(err))
		return nil, 0, err
	}

	return valuations, total, nil
}

func (r *PropertyValuationRepo) Update(valuation *models.PropertyValuation) error {
	if err := r.db.Save(valuation).Error; err != nil {
		r.lg.Error("failed to update property valuation", zap.String("id", valuation.ID.String()), zap.Error(err))
		return err
	}
	r.lg.Info("property valuation updated successfully", zap.String("id", valuation.ID.String()))
	return nil
}

func (r *PropertyValuationRepo) Delete(id uuid.UUID) error {
	if err := r.db.Delete(&models.PropertyValuation{}, "id = ?", id).Error; err != nil {
		r.lg.Error("failed to delete property valuation", zap.String("id", id.String()), zap.Error(err))
		return err
	}
	r.lg.Info("property valuation deleted successfully", zap.String("id", id.String()))
	return nil
}

func (r *PropertyValuationRepo) SoftDelete(id uuid.UUID, deletedBy uuid.UUID) error {
	now := time.Now()
	err := r.db.Model(&models.PropertyValuation{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"deleted_at": &now,
			"updated_by": &deletedBy,
		}).Error

	if err != nil {
		r.lg.Error("failed to soft delete property valuation", zap.String("id", id.String()), zap.Error(err))
		return err
	}

	r.lg.Info("property valuation soft deleted successfully", zap.String("id", id.String()))
	return nil
}

func (r *PropertyValuationRepo) BulkCreateForProperty(propertyID uuid.UUID, valuations []models.PropertyValuation) error {
	if len(valuations) == 0 {
		return nil
	}

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Set property ID for all valuation items
	for i := range valuations {
		valuations[i].PropertyID = propertyID
	}

	// Create all valuation items
	if err := tx.Create(&valuations).Error; err != nil {
		tx.Rollback()
		r.lg.Error("failed to bulk create property valuations",
			zap.String("property_id", propertyID.String()),
			zap.Int("count", len(valuations)),
			zap.Error(err))
		return err
	}

	if err := tx.Commit().Error; err != nil {
		r.lg.Error("failed to commit bulk create valuations transaction", zap.Error(err))
		return err
	}

	r.lg.Info("property valuations bulk created successfully",
		zap.String("property_id", propertyID.String()),
		zap.Int("count", len(valuations)))
	return nil
}

func (r *PropertyValuationRepo) BulkDeleteForProperty(propertyID uuid.UUID) error {
	err := r.db.Delete(&models.PropertyValuation{}, "property_id = ?", propertyID).Error

	if err != nil {
		r.lg.Error("failed to bulk delete property valuations",
			zap.String("property_id", propertyID.String()),
			zap.Error(err))
		return err
	}

	r.lg.Info("property valuations bulk deleted successfully",
		zap.String("property_id", propertyID.String()))
	return nil
}

func (r *PropertyValuationRepo) GetValidValuations(propertyID uuid.UUID, asOfDate time.Time) ([]models.PropertyValuation, error) {
	var valuations []models.PropertyValuation
	err := r.db.Where("property_id = ? AND (valid_until IS NULL OR valid_until >= ?)", propertyID, asOfDate).
		Order("valuation_date DESC").
		Find(&valuations).Error

	if err != nil {
		r.lg.Error("failed to get valid property valuations",
			zap.String("property_id", propertyID.String()),
			zap.Time("as_of_date", asOfDate),
			zap.Error(err))
		return nil, err
	}

	return valuations, nil
}

func (r *PropertyValuationRepo) GetMarketValueHistory(propertyID uuid.UUID, dateFrom, dateTo time.Time) ([]models.PropertyValuation, error) {
	var valuations []models.PropertyValuation
	err := r.db.Where("property_id = ? AND valuation_type = ? AND valuation_date BETWEEN ? AND ?",
		propertyID, "market", dateFrom, dateTo).
		Order("valuation_date ASC").
		Find(&valuations).Error

	if err != nil {
		r.lg.Error("failed to get market value history",
			zap.String("property_id", propertyID.String()),
			zap.Time("date_from", dateFrom),
			zap.Time("date_to", dateTo),
			zap.Error(err))
		return nil, err
	}

	return valuations, nil
}
