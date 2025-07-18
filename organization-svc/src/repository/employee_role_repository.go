package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rem-gestion/api-suite/organization/src/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// employeeRoleRepository implements EmployeeRoleRepository interface
type employeeRoleRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewEmployeeRoleRepository creates a new employee role repository instance
func NewEmployeeRoleRepository(db *gorm.DB, logger *zap.Logger) EmployeeRoleRepository {
	return &employeeRoleRepository{db: db, logger: logger}
}

// ===============================
// BASIC CRUD OPERATIONS
// ===============================

func (r *employeeRoleRepository) Create(ctx context.Context, empRole *models.EmployeeRole) (*models.EmployeeRole, error) {
	if err := r.db.WithContext(ctx).Create(empRole).Error; err != nil {
		return nil, fmt.Errorf("failed to create employee role: %w", err)
	}
	return empRole, nil
}

func (r *employeeRoleRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.EmployeeRole, error) {
	var empRole models.EmployeeRole
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&empRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("employee role not found")
		}
		return nil, fmt.Errorf("failed to get employee role: %w", err)
	}
	return &empRole, nil
}

func (r *employeeRoleRepository) GetByEmployeeAndRole(ctx context.Context, employeeID, roleID uuid.UUID) (*models.EmployeeRole, error) {
	var empRole models.EmployeeRole
	if err := r.db.WithContext(ctx).Where("employee_id = ? AND role_id = ?", employeeID, roleID).First(&empRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("employee role not found")
		}
		return nil, fmt.Errorf("failed to get employee role: %w", err)
	}
	return &empRole, nil
}

func (r *employeeRoleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.EmployeeRole{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete employee role: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("employee role not found")
	}
	return nil
}

func (r *employeeRoleRepository) SoftDelete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	updates := map[string]interface{}{
		"deleted_at": time.Now(),
		"deleted_by": deletedBy,
	}

	result := r.db.WithContext(ctx).Model(&models.EmployeeRole{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to soft delete employee role: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("employee role not found")
	}
	return nil
}

// ===============================
// ROLE ASSIGNMENTS
// ===============================

func (r *employeeRoleRepository) AssignRole(ctx context.Context, employeeID, roleID uuid.UUID, isPrimary bool, createdBy uuid.UUID) (*models.EmployeeRole, error) {
	// If this is a primary role, first clear any existing primary role for this employee
	if isPrimary {
		if err := r.ClearPrimaryRole(ctx, employeeID); err != nil {
			return nil, fmt.Errorf("failed to clear existing primary role: %w", err)
		}
	}

	empRole := &models.EmployeeRole{
		ID:         uuid.New(),
		EmployeeID: employeeID,
		RoleID:     roleID,
		IsPrimary:  isPrimary,
		CreatedBy:  createdBy,
		CreatedAt:  time.Now(),
	}

	if err := r.db.WithContext(ctx).Create(empRole).Error; err != nil {
		return nil, fmt.Errorf("failed to assign role to employee: %w", err)
	}

	return empRole, nil
}

func (r *employeeRoleRepository) RemoveRole(ctx context.Context, employeeID, roleID uuid.UUID) error {
	result := r.db.WithContext(ctx).Where("employee_id = ? AND role_id = ?", employeeID, roleID).Delete(&models.EmployeeRole{})
	if result.Error != nil {
		return fmt.Errorf("failed to remove role from employee: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("employee role assignment not found")
	}
	return nil
}

func (r *employeeRoleRepository) SetPrimaryRole(ctx context.Context, employeeID, roleID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// First clear any existing primary role
		if err := tx.Model(&models.EmployeeRole{}).
			Where("employee_id = ? AND is_primary = ?", employeeID, true).
			Update("is_primary", false).Error; err != nil {
			return fmt.Errorf("failed to clear existing primary role: %w", err)
		}

		// Set the new primary role
		result := tx.Model(&models.EmployeeRole{}).
			Where("employee_id = ? AND role_id = ?", employeeID, roleID).
			Update("is_primary", true)

		if result.Error != nil {
			return fmt.Errorf("failed to set primary role: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("employee role assignment not found")
		}
		return nil
	})
}

func (r *employeeRoleRepository) ClearPrimaryRole(ctx context.Context, employeeID uuid.UUID) error {
	result := r.db.WithContext(ctx).Model(&models.EmployeeRole{}).
		Where("employee_id = ? AND is_primary = ?", employeeID, true).
		Update("is_primary", false)

	if result.Error != nil {
		return fmt.Errorf("failed to clear primary role: %w", result.Error)
	}
	return nil
}

// ===============================
// LISTING OPERATIONS
// ===============================

func (r *employeeRoleRepository) GetByEmployee(ctx context.Context, employeeID uuid.UUID) ([]models.EmployeeRole, error) {
	var empRoles []models.EmployeeRole
	if err := r.db.WithContext(ctx).Preload("Role").Where("employee_id = ?", employeeID).Find(&empRoles).Error; err != nil {
		return nil, fmt.Errorf("failed to get employee roles: %w", err)
	}
	return empRoles, nil
}

func (r *employeeRoleRepository) GetByRole(ctx context.Context, roleID uuid.UUID) ([]models.EmployeeRole, error) {
	var empRoles []models.EmployeeRole
	if err := r.db.WithContext(ctx).Preload("Employee").Where("role_id = ?", roleID).Find(&empRoles).Error; err != nil {
		return nil, fmt.Errorf("failed to get role assignments: %w", err)
	}
	return empRoles, nil
}

func (r *employeeRoleRepository) GetPrimaryRoleByEmployee(ctx context.Context, employeeID uuid.UUID) (*models.EmployeeRole, error) {
	var empRole models.EmployeeRole
	if err := r.db.WithContext(ctx).Preload("Role").Where("employee_id = ? AND is_primary = ?", employeeID, true).First(&empRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("primary role not found for employee")
		}
		return nil, fmt.Errorf("failed to get primary role: %w", err)
	}
	return &empRole, nil
}

// ===============================
// BUSINESS LOGIC SUPPORT
// ===============================

func (r *employeeRoleRepository) HasRole(ctx context.Context, employeeID, roleID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.EmployeeRole{}).
		Where("employee_id = ? AND role_id = ?", employeeID, roleID).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check role assignment: %w", err)
	}
	return count > 0, nil
}

func (r *employeeRoleRepository) CountByRole(ctx context.Context, roleID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.EmployeeRole{}).Where("role_id = ?", roleID).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count employees by role: %w", err)
	}
	return count, nil
}

func (r *employeeRoleRepository) GetEmployeeRoleCount(ctx context.Context, employeeID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.EmployeeRole{}).Where("employee_id = ?", employeeID).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count roles for employee: %w", err)
	}
	return count, nil
}

func (r *employeeRoleRepository) RemoveAllRoles(ctx context.Context, employeeID uuid.UUID) error {
	result := r.db.WithContext(ctx).Where("employee_id = ?", employeeID).Delete(&models.EmployeeRole{})
	if result.Error != nil {
		return fmt.Errorf("failed to remove all roles from employee: %w", result.Error)
	}
	return nil
}
