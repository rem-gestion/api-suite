package repository

import (
	"github.com/rem-gestion/api-suite/property/src/models"
	rerrors "github.com/rem-gestion/rem-common/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PropertyTypeRepo struct {
	db *gorm.DB
	lg *zap.Logger
}

func NewPropertyTypeRepo(db *gorm.DB, lg *zap.Logger) PropertyTypeRepository {
	return &PropertyTypeRepo{db: db, lg: lg.Named("property-type-repo")}
}

func (r *PropertyTypeRepo) Create(propertyType *models.PropertyType) (*models.PropertyType, error) {
	if err := r.db.Create(propertyType).Error; err != nil {
		r.lg.Error("failed to create property type", zap.Error(err))
		return nil, err
	}
	r.lg.Info("property type created successfully", zap.String("code", propertyType.Code))
	return propertyType, nil
}

func (r *PropertyTypeRepo) GetByID(id int32) (*models.PropertyType, error) {
	var propertyType models.PropertyType
	err := r.db.First(&propertyType, "id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &rerrors.NotFoundError{Msg: "property type not found"}
		}
		r.lg.Error("failed to get property type by ID", zap.Int32("id", id), zap.Error(err))
		return nil, err
	}

	return &propertyType, nil
}

func (r *PropertyTypeRepo) GetByCode(code string) (*models.PropertyType, error) {
	var propertyType models.PropertyType
	err := r.db.First(&propertyType, "code = ?", code).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &rerrors.NotFoundError{Msg: "property type not found"}
		}
		r.lg.Error("failed to get property type by code", zap.String("code", code), zap.Error(err))
		return nil, err
	}

	return &propertyType, nil
}

func (r *PropertyTypeRepo) List(category *string, isActive *bool, limit, offset int) ([]models.PropertyType, int64, error) {
	var propertyTypes []models.PropertyType
	var total int64

	query := r.db.Model(&models.PropertyType{})

	if category != nil {
		query = query.Where("category = ?", *category)
	}

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		r.lg.Error("failed to count property types", zap.Error(err))
		return nil, 0, err
	}

	// Get records
	err := query.Limit(limit).Offset(offset).Find(&propertyTypes).Error
	if err != nil {
		r.lg.Error("failed to list property types", zap.Error(err))
		return nil, 0, err
	}

	return propertyTypes, total, nil
}

func (r *PropertyTypeRepo) Update(propertyType *models.PropertyType) error {
	if err := r.db.Save(propertyType).Error; err != nil {
		r.lg.Error("failed to update property type", zap.Int32("id", propertyType.ID), zap.Error(err))
		return err
	}
	r.lg.Info("property type updated successfully", zap.Int32("id", propertyType.ID))
	return nil
}

func (r *PropertyTypeRepo) Delete(id int32) error {
	if err := r.db.Delete(&models.PropertyType{}, "id = ?", id).Error; err != nil {
		r.lg.Error("failed to delete property type", zap.Int32("id", id), zap.Error(err))
		return err
	}
	r.lg.Info("property type deleted successfully", zap.Int32("id", id))
	return nil
}

// ManagerTypeRepo implementation
type ManagerTypeRepo struct {
	db *gorm.DB
	lg *zap.Logger
}

func NewManagerTypeRepo(db *gorm.DB, lg *zap.Logger) ManagerTypeRepository {
	return &ManagerTypeRepo{db: db, lg: lg.Named("manager-type-repo")}
}

func (r *ManagerTypeRepo) Create(managerType *models.ManagerType) (*models.ManagerType, error) {
	if err := r.db.Create(managerType).Error; err != nil {
		r.lg.Error("failed to create manager type", zap.Error(err))
		return nil, err
	}
	r.lg.Info("manager type created successfully", zap.String("code", managerType.Code))
	return managerType, nil
}

func (r *ManagerTypeRepo) GetByID(id int32) (*models.ManagerType, error) {
	var managerType models.ManagerType
	err := r.db.First(&managerType, "id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &rerrors.NotFoundError{Msg: "manager type not found"}
		}
		r.lg.Error("failed to get manager type by ID", zap.Int32("id", id), zap.Error(err))
		return nil, err
	}

	return &managerType, nil
}

func (r *ManagerTypeRepo) GetByCode(code string) (*models.ManagerType, error) {
	var managerType models.ManagerType
	err := r.db.First(&managerType, "code = ?", code).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &rerrors.NotFoundError{Msg: "manager type not found"}
		}
		r.lg.Error("failed to get manager type by code", zap.String("code", code), zap.Error(err))
		return nil, err
	}

	return &managerType, nil
}

func (r *ManagerTypeRepo) List(isActive *bool, limit, offset int) ([]models.ManagerType, int64, error) {
	var managerTypes []models.ManagerType
	var total int64

	query := r.db.Model(&models.ManagerType{})

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		r.lg.Error("failed to count manager types", zap.Error(err))
		return nil, 0, err
	}

	// Get records
	err := query.Limit(limit).Offset(offset).Find(&managerTypes).Error
	if err != nil {
		r.lg.Error("failed to list manager types", zap.Error(err))
		return nil, 0, err
	}

	return managerTypes, total, nil
}

func (r *ManagerTypeRepo) Update(managerType *models.ManagerType) error {
	if err := r.db.Save(managerType).Error; err != nil {
		r.lg.Error("failed to update manager type", zap.Int32("id", managerType.ID), zap.Error(err))
		return err
	}
	r.lg.Info("manager type updated successfully", zap.Int32("id", managerType.ID))
	return nil
}

func (r *ManagerTypeRepo) Delete(id int32) error {
	if err := r.db.Delete(&models.ManagerType{}, "id = ?", id).Error; err != nil {
		r.lg.Error("failed to delete manager type", zap.Int32("id", id), zap.Error(err))
		return err
	}
	r.lg.Info("manager type deleted successfully", zap.Int32("id", id))
	return nil
}

// AmenityRepo implementation
type AmenityRepo struct {
	db *gorm.DB
	lg *zap.Logger
}

func NewAmenityRepo(db *gorm.DB, lg *zap.Logger) AmenityRepository {
	return &AmenityRepo{db: db, lg: lg.Named("amenity-repo")}
}

func (r *AmenityRepo) Create(amenity *models.Amenity) (*models.Amenity, error) {
	if err := r.db.Create(amenity).Error; err != nil {
		r.lg.Error("failed to create amenity", zap.Error(err))
		return nil, err
	}
	r.lg.Info("amenity created successfully", zap.String("name", amenity.Name))
	return amenity, nil
}

func (r *AmenityRepo) GetByID(id int32) (*models.Amenity, error) {
	var amenity models.Amenity
	err := r.db.First(&amenity, "id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &rerrors.NotFoundError{Msg: "amenity not found"}
		}
		r.lg.Error("failed to get amenity by ID", zap.Int32("id", id), zap.Error(err))
		return nil, err
	}

	return &amenity, nil
}

func (r *AmenityRepo) GetByName(name string) (*models.Amenity, error) {
	var amenity models.Amenity
	err := r.db.First(&amenity, "name = ?", name).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &rerrors.NotFoundError{Msg: "amenity not found"}
		}
		r.lg.Error("failed to get amenity by name", zap.String("name", name), zap.Error(err))
		return nil, err
	}

	return &amenity, nil
}

func (r *AmenityRepo) List(category *models.AmenityCategory, isActive *bool, limit, offset int) ([]models.Amenity, int64, error) {
	var amenities []models.Amenity
	var total int64

	query := r.db.Model(&models.Amenity{})

	if category != nil {
		query = query.Where("category = ?", *category)
	}

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		r.lg.Error("failed to count amenities", zap.Error(err))
		return nil, 0, err
	}

	// Get records
	err := query.Limit(limit).Offset(offset).Find(&amenities).Error
	if err != nil {
		r.lg.Error("failed to list amenities", zap.Error(err))
		return nil, 0, err
	}

	return amenities, total, nil
}

func (r *AmenityRepo) Update(amenity *models.Amenity) error {
	if err := r.db.Save(amenity).Error; err != nil {
		r.lg.Error("failed to update amenity", zap.Int32("id", amenity.ID), zap.Error(err))
		return err
	}
	r.lg.Info("amenity updated successfully", zap.Int32("id", amenity.ID))
	return nil
}

func (r *AmenityRepo) Delete(id int32) error {
	if err := r.db.Delete(&models.Amenity{}, "id = ?", id).Error; err != nil {
		r.lg.Error("failed to delete amenity", zap.Int32("id", id), zap.Error(err))
		return err
	}
	r.lg.Info("amenity deleted successfully", zap.Int32("id", id))
	return nil
}
