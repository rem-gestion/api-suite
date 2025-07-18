package repository

import (
	"strings"

	"github.com/google/uuid"
	"github.com/rem-gestion/api-suite/property/src/models"
	rerrors "github.com/rem-gestion/rem-common/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

/* ───────────────── struct & ctor ────────────────── */

type PropertyRepo struct {
	db *gorm.DB
	lg *zap.Logger
}

func NewPropertyRepo(db *gorm.DB, lg *zap.Logger) *PropertyRepo {
	return &PropertyRepo{db: db, lg: lg}
}

/* ───────────────── Properties CRUD ─────────────────── */

func (r *PropertyRepo) Create(p *models.Property) (*models.Property, error) {
	err := r.withTx(func(tx *gorm.DB) error {
		return tx.Create(p).Error
	})
	return p, err
}

func (r *PropertyRepo) Get(id string) (*models.Property, error) {
	var p models.Property
	err := r.preloads(r.db).
		First(&p, "id = ?", id).Error
	return &p, err
}

func (r *PropertyRepo) List(search string, propertyTypeID *int, ownerPersonID *uuid.UUID, limit, offset int) ([]models.Property, int64, error) {
	var list []models.Property
	var total int64

	q := r.preloads(r.db.Model(&models.Property{})).
		Limit(limit).Offset(offset)

	if propertyTypeID != nil {
		q = q.Where("property_type_id = ?", *propertyTypeID)
	}

	if ownerPersonID != nil {
		q = q.Where("owner_person_id = ?", *ownerPersonID)
	}

	if search != "" {
		search = strings.ToLower(search)
		q = q.Where("LOWER(description) LIKE ? OR LOWER(internal_code) LIKE ?",
			"%"+search+"%", "%"+search+"%")
	}

	// Count total
	countQ := r.db.Model(&models.Property{})
	if propertyTypeID != nil {
		countQ = countQ.Where("property_type_id = ?", *propertyTypeID)
	}
	if ownerPersonID != nil {
		countQ = countQ.Where("owner_person_id = ?", *ownerPersonID)
	}
	if search != "" {
		search = strings.ToLower(search)
		countQ = countQ.Where("LOWER(description) LIKE ? OR LOWER(internal_code) LIKE ?",
			"%"+search+"%", "%"+search+"%")
	}
	countQ.Count(&total)

	err := q.Find(&list).Error
	return list, total, err
}

func (r *PropertyRepo) Update(id string, updates map[string]interface{}) (*models.Property, error) {
	var p models.Property
	err := r.withTx(func(tx *gorm.DB) error {
		if err := tx.Model(&p).Where("id = ?", id).Updates(updates).Error; err != nil {
			return err
		}
		return r.preloads(tx).First(&p, "id = ?", id).Error
	})
	return &p, err
}

func (r *PropertyRepo) Delete(id string) error {
	return r.withTx(func(tx *gorm.DB) error {
		return tx.Where("id = ?", id).Delete(&models.Property{}).Error
	})
}

func (r *PropertyRepo) GetByInternalCode(code string) (*models.Property, error) {
	var p models.Property
	err := r.preloads(r.db).
		Where("internal_code = ?", code).
		First(&p).Error
	return &p, err
}

/* ───────────────── Property Management CRUD ─────────────────── */

func (r *PropertyRepo) CreateManagement(pm *models.PropertyManagement) (*models.PropertyManagement, error) {
	err := r.withTx(func(tx *gorm.DB) error {
		return tx.Create(pm).Error
	})
	return pm, err
}

func (r *PropertyRepo) GetManagement(id string) (*models.PropertyManagement, error) {
	var pm models.PropertyManagement
	err := r.db.Preload("Property").First(&pm, "id = ?", id).Error
	return &pm, err
}

func (r *PropertyRepo) ListManagementByProperty(propertyID string) ([]models.PropertyManagement, error) {
	var list []models.PropertyManagement
	err := r.db.Where("property_id = ?", propertyID).
		Order("managed_since DESC").
		Find(&list).Error
	return list, err
}

func (r *PropertyRepo) UpdateManagement(id string, updates map[string]interface{}) (*models.PropertyManagement, error) {
	var pm models.PropertyManagement
	err := r.withTx(func(tx *gorm.DB) error {
		if err := tx.Model(&pm).Where("id = ?", id).Updates(updates).Error; err != nil {
			return err
		}
		return tx.Preload("Property").First(&pm, "id = ?", id).Error
	})
	return &pm, err
}

func (r *PropertyRepo) DeleteManagement(id string) error {
	return r.withTx(func(tx *gorm.DB) error {
		return tx.Where("id = ?", id).Delete(&models.PropertyManagement{}).Error
	})
}

/* ───────────────── Amenities CRUD ─────────────────── */

func (r *PropertyRepo) CreateAmenity(a *models.Amenity) (*models.Amenity, error) {
	err := r.withTx(func(tx *gorm.DB) error {
		return tx.Create(a).Error
	})
	return a, err
}

func (r *PropertyRepo) GetAmenity(id string) (*models.Amenity, error) {
	var a models.Amenity
	err := r.db.First(&a, "id = ?", id).Error
	return &a, err
}

func (r *PropertyRepo) ListAmenities(search string, category *string, limit, offset int) ([]models.Amenity, int64, error) {
	var list []models.Amenity
	var total int64

	q := r.db.Model(&models.Amenity{}).
		Limit(limit).Offset(offset)

	if category != nil && *category != "" {
		q = q.Where("category = ?", *category)
	}

	if search != "" {
		search = strings.ToLower(search)
		q = q.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?",
			"%"+search+"%", "%"+search+"%")
	}

	// Count total
	countQ := r.db.Model(&models.Amenity{})
	if category != nil && *category != "" {
		countQ = countQ.Where("category = ?", *category)
	}
	if search != "" {
		search = strings.ToLower(search)
		countQ = countQ.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?",
			"%"+search+"%", "%"+search+"%")
	}
	countQ.Count(&total)

	err := q.Find(&list).Error
	return list, total, err
}

func (r *PropertyRepo) UpdateAmenity(id string, updates map[string]interface{}) (*models.Amenity, error) {
	var a models.Amenity
	err := r.withTx(func(tx *gorm.DB) error {
		if err := tx.Model(&a).Where("id = ?", id).Updates(updates).Error; err != nil {
			return err
		}
		return tx.First(&a, "id = ?", id).Error
	})
	return &a, err
}

func (r *PropertyRepo) DeleteAmenity(id string) error {
	return r.withTx(func(tx *gorm.DB) error {
		return tx.Where("id = ?", id).Delete(&models.Amenity{}).Error
	})
}

func (r *PropertyRepo) GetAmenityByName(name string) (*models.Amenity, error) {
	var a models.Amenity
	err := r.db.Where("name = ?", name).First(&a).Error
	return &a, err
}

/* ───────────────── Property Amenities Management ─────────────────── */

func (r *PropertyRepo) AddAmenityToProperty(propertyID, amenityID uuid.UUID) error {
	pa := &models.PropertyAmenity{
		PropertyID: propertyID,
		AmenityID:  amenityID,
	}
	return r.withTx(func(tx *gorm.DB) error {
		return tx.Create(pa).Error
	})
}

func (r *PropertyRepo) RemoveAmenityFromProperty(propertyID, amenityID uuid.UUID) error {
	return r.withTx(func(tx *gorm.DB) error {
		return tx.Where("property_id = ? AND amenity_id = ?", propertyID, amenityID).
			Delete(&models.PropertyAmenity{}).Error
	})
}

func (r *PropertyRepo) GetPropertyAmenities(propertyID uuid.UUID) ([]models.Amenity, error) {
	var amenities []models.Amenity
	err := r.db.Table("amenity").
		Joins("JOIN property_amenities ON amenity.id = property_amenities.amenity_id").
		Where("property_amenities.property_id = ?", propertyID).
		Find(&amenities).Error
	return amenities, err
}

/* ───────────────── Type validation methods ─────────────────── */

func (r *PropertyRepo) PropertyTypeExists(id int) (bool, error) {
	var count int64
	err := r.db.Model(&models.PropertyType{}).
		Where("id = ? AND is_active = true", id).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *PropertyRepo) ManagerTypeExists(id int) (bool, error) {
	var count int64
	err := r.db.Model(&models.ManagerType{}).
		Where("id = ? AND is_active = true", id).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *PropertyRepo) GetPropertyTypes() ([]models.PropertyType, error) {
	var types []models.PropertyType
	err := r.db.Where("is_active = true").
		Order("name").
		Find(&types).Error
	return types, err
}

func (r *PropertyRepo) GetManagerTypes() ([]models.ManagerType, error) {
	var types []models.ManagerType
	err := r.db.Where("is_active = true").
		Order("name").
		Find(&types).Error
	return types, err
}

/* ───────────────── helpers ─────────────────── */

func (r *PropertyRepo) preloads(db *gorm.DB) *gorm.DB {
	return db.
		Preload("PropertyType").
		Preload("PropertyManagement").
		Preload("PropertyManagement.ManagerType").
		Preload("PropertyAmenities").
		Preload("PropertyAmenities.Amenity")
}

func (r *PropertyRepo) withTx(fn func(*gorm.DB) error) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return &rerrors.InternalServerError{Msg: "failed to begin transaction: " + tx.Error.Error()}
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return &rerrors.InternalServerError{Msg: "transaction operation failed: " + err.Error()}
	}

	if err := tx.Commit().Error; err != nil {
		return &rerrors.InternalServerError{Msg: "failed to commit transaction: " + err.Error()}
	}

	return nil
}
