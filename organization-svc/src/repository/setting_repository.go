package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/rem-gestion/api-suite/organization/src/models"
)

// organizationSettingRepository implements OrganizationSettingRepository interface
type organizationSettingRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewOrganizationSettingRepository creates a new organization setting repository instance
func NewOrganizationSettingRepository(db *gorm.DB, logger *zap.Logger) OrganizationSettingRepository {
	return &organizationSettingRepository{db: db, logger: logger}
}

// ===============================
// SETTINGS CRUD
// ===============================

func (r *organizationSettingRepository) Create(ctx context.Context, setting *models.OrganizationSetting) (*models.OrganizationSetting, error) {
	if err := r.db.WithContext(ctx).Create(setting).Error; err != nil {
		return nil, fmt.Errorf("failed to create organization setting: %w", err)
	}
	return setting, nil
}

func (r *organizationSettingRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.OrganizationSetting, error) {
	// Este método no es apropiado para settings con clave primaria compuesta
	// Se mantiene por compatibilidad con la interfaz, pero debería usar GetByKey
	return nil, fmt.Errorf("GetByID not supported for OrganizationSetting with composite primary key, use GetByKey instead")
}

func (r *organizationSettingRepository) GetByKey(ctx context.Context, orgID uuid.UUID, key string) (*models.OrganizationSetting, error) {
	var setting models.OrganizationSetting
	if err := r.db.WithContext(ctx).Where("organization_id = ? AND setting_key = ?", orgID, key).First(&setting).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("organization setting not found")
		}
		return nil, fmt.Errorf("failed to get organization setting by key: %w", err)
	}
	return &setting, nil
}

func (r *organizationSettingRepository) Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.OrganizationSetting, error) {
	// Este método no es apropiado para settings con clave primaria compuesta
	// Se mantiene por compatibilidad con la interfaz, pero debería usar UpdateByKey
	return nil, fmt.Errorf("Update not supported for OrganizationSetting with composite primary key, use UpdateByKey instead")
}

// UpdateByKey actualiza un setting por organización y clave
func (r *organizationSettingRepository) UpdateByKey(ctx context.Context, orgID uuid.UUID, key string, updates map[string]interface{}) (*models.OrganizationSetting, error) {
	var setting models.OrganizationSetting

	// First check if setting exists
	if err := r.db.WithContext(ctx).Where("organization_id = ? AND setting_key = ?", orgID, key).First(&setting).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("organization setting not found")
		}
		return nil, fmt.Errorf("failed to get organization setting: %w", err)
	}

	// Add updated_at timestamp
	updates["updated_at"] = time.Now()

	// Perform update
	if err := r.db.WithContext(ctx).Model(&setting).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update organization setting: %w", err)
	}

	// Return updated setting
	if err := r.db.WithContext(ctx).Where("organization_id = ? AND setting_key = ?", orgID, key).First(&setting).Error; err != nil {
		return nil, fmt.Errorf("failed to get updated organization setting: %w", err)
	}

	return &setting, nil
}

func (r *organizationSettingRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Este método no es apropiado para settings con clave primaria compuesta
	// Se mantiene por compatibilidad con la interfaz, pero debería usar DeleteByKey
	return fmt.Errorf("Delete not supported for OrganizationSetting with composite primary key, use DeleteByKey instead")
}

// ===============================
// BULK OPERATIONS
// ===============================

func (r *organizationSettingRepository) GetAllByOrganization(ctx context.Context, orgID uuid.UUID) ([]models.OrganizationSetting, error) {
	var settings []models.OrganizationSetting
	if err := r.db.WithContext(ctx).Where("organization_id = ?", orgID).Find(&settings).Error; err != nil {
		return nil, fmt.Errorf("failed to get organization settings: %w", err)
	}
	return settings, nil
}

func (r *organizationSettingRepository) SetMultiple(ctx context.Context, orgID uuid.UUID, settings map[string]interface{}, createdBy uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for key, value := range settings {
			var setting models.OrganizationSetting

			// Serializar el valor a JSON para asegurar compatibilidad con JSONB
			jsonValue, err := json.Marshal(value)
			if err != nil {
				return fmt.Errorf("failed to marshal setting value for %s: %w", key, err)
			}

			// Check if setting exists
			err = tx.Where("organization_id = ? AND setting_key = ?", orgID, key).First(&setting).Error

			if err == gorm.ErrRecordNotFound {
				// Create new setting
				setting = models.OrganizationSetting{
					OrganizationID: orgID,
					SettingKey:     key,
					SettingValue:   string(jsonValue), // Usar el JSON como string
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
					UpdatedBy:      &createdBy,
				}
				if err := tx.Create(&setting).Error; err != nil {
					return fmt.Errorf("failed to create setting %s: %w", key, err)
				}
			} else if err != nil {
				return fmt.Errorf("failed to check setting %s: %w", key, err)
			} else {
				// Update existing setting
				updates := map[string]interface{}{
					"setting_value": string(jsonValue), // Usar el JSON como string
					"updated_by":    createdBy,
					"updated_at":    time.Now(),
				}
				if err := tx.Model(&setting).Updates(updates).Error; err != nil {
					return fmt.Errorf("failed to update setting %s: %w", key, err)
				}
			}
		}
		return nil
	})
}

func (r *organizationSettingRepository) DeleteByKey(ctx context.Context, orgID uuid.UUID, key string) error {
	result := r.db.WithContext(ctx).Where("organization_id = ? AND setting_key = ?", orgID, key).Delete(&models.OrganizationSetting{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete organization setting: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("organization setting not found")
	}
	return nil
}

// ===============================
// CONFIGURATION MANAGEMENT
// ===============================

func (r *organizationSettingRepository) GetEditableSettings(ctx context.Context, orgID uuid.UUID) ([]models.OrganizationSetting, error) {
	var settings []models.OrganizationSetting
	if err := r.db.WithContext(ctx).Where("organization_id = ? AND is_editable = ?", orgID, true).Find(&settings).Error; err != nil {
		return nil, fmt.Errorf("failed to get editable organization settings: %w", err)
	}
	return settings, nil
}

func (r *organizationSettingRepository) GetDefaultSettings(ctx context.Context) (map[string]interface{}, error) {
	// Define default settings for organizations
	defaults := map[string]interface{}{
		"timezone":                   "UTC",
		"date_format":                "YYYY-MM-DD",
		"time_format":                "24h",
		"currency":                   "USD",
		"language":                   "en",
		"allow_employee_invites":     true,
		"require_email_verification": true,
		"max_employee_count":         100,
		"enable_2fa":                 false,
		"session_timeout_minutes":    60,
		"password_policy_enabled":    true,
		"backup_enabled":             true,
		"audit_log_enabled":          true,
		"notification_email":         true,
		"notification_sms":           false,
		"notification_push":          true,
	}
	return defaults, nil
}

func (r *organizationSettingRepository) ResetToDefaults(ctx context.Context, orgID uuid.UUID, keys []string, updatedBy uuid.UUID) error {
	defaults, err := r.GetDefaultSettings(ctx)
	if err != nil {
		return fmt.Errorf("failed to get default settings: %w", err)
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, key := range keys {
			if defaultValue, exists := defaults[key]; exists {
				updates := map[string]interface{}{
					"setting_value": fmt.Sprintf("%v", defaultValue),
					"updated_by":    updatedBy,
					"updated_at":    time.Now(),
				}

				result := tx.Model(&models.OrganizationSetting{}).
					Where("organization_id = ? AND setting_key = ?", orgID, key).
					Updates(updates)

				if result.Error != nil {
					return fmt.Errorf("failed to reset setting %s: %w", key, result.Error)
				}
			}
		}
		return nil
	})
}

// ===============================
// BUSINESS LOGIC SUPPORT
// ===============================

func (r *organizationSettingRepository) SettingExists(ctx context.Context, orgID uuid.UUID, key string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.OrganizationSetting{}).
		Where("organization_id = ? AND setting_key = ?", orgID, key).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check setting existence: %w", err)
	}
	return count > 0, nil
}

func (r *organizationSettingRepository) GetSettingsAsMap(ctx context.Context, orgID uuid.UUID) (map[string]interface{}, error) {
	var settings []models.OrganizationSetting
	if err := r.db.WithContext(ctx).Where("organization_id = ?", orgID).Find(&settings).Error; err != nil {
		return nil, fmt.Errorf("failed to get organization settings: %w", err)
	}

	result := make(map[string]interface{})
	for _, setting := range settings {
		result[setting.SettingKey] = setting.SettingValue
	}

	return result, nil
}
