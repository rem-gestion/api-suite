package services

import (
	"github.com/rem-gestion/property-svc/src/models"
	"gorm.io/gorm"
)

func CreateProperty(db *gorm.DB, property *models.Property) error {
	return db.Create(property).Error
}
