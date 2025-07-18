package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/rem-gestion/api-suite/organization/src/dto"
	"github.com/rem-gestion/api-suite/organization/src/models"
)

// roleRepository implements RoleRepository interface
type roleRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewRoleRepository creates a new role repository instance
func NewRoleRepository(db *gorm.DB, logger *zap.Logger) RoleRepository {
	return &roleRepository{db: db, logger: logger}
}

// ===============================
// BASIC CRUD OPERATIONS
// ===============================

func (r *roleRepository) Create(ctx context.Context, role *models.OrganizationRole) (*models.OrganizationRole, error) {
	if err := r.db.WithContext(ctx).Create(role).Error; err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}
	return role, nil
}

func (r *roleRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.OrganizationRole, error) {
	var role models.OrganizationRole
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("role not found")
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}
	return &role, nil
}

func (r *roleRepository) Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.OrganizationRole, error) {
	var role models.OrganizationRole

	// First check if role exists
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("role not found")
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	// Add updated_at timestamp
	updates["updated_at"] = time.Now()

	// Perform update
	if err := r.db.WithContext(ctx).Model(&role).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	// Return updated role
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&role).Error; err != nil {
		return nil, fmt.Errorf("failed to get updated role: %w", err)
	}

	return &role, nil
}

func (r *roleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Unscoped().Delete(&models.OrganizationRole{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete role: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("role not found")
	}
	return nil
}

func (r *roleRepository) SoftDelete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	updates := map[string]interface{}{
		"deleted_at": time.Now(),
		"updated_by": deletedBy,
		"updated_at": time.Now(),
	}

	result := r.db.WithContext(ctx).Model(&models.OrganizationRole{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to soft delete role: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("role not found")
	}
	return nil
}

// ===============================
// LISTING AND FILTERING
// ===============================

func (r *roleRepository) List(ctx context.Context, filters dto.RoleFiltersRequest) ([]models.OrganizationRole, int64, error) {
	var roles []models.OrganizationRole
	var total int64

	query := r.db.WithContext(ctx).Model(&models.OrganizationRole{})

	// Apply filters
	query = r.applyRoleFilters(query, filters)

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count roles: %w", err)
	}

	// Apply pagination and sorting
	query = r.applyPaginationAndSorting(query, filters.PaginationRequest, filters.FilterRequest)

	// Execute query
	if err := query.Find(&roles).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list roles: %w", err)
	}

	return roles, total, nil
}

func (r *roleRepository) ListByOrganization(ctx context.Context, orgID uuid.UUID, filters dto.RoleFiltersRequest) ([]models.OrganizationRole, int64, error) {
	var roles []models.OrganizationRole
	var total int64

	query := r.db.WithContext(ctx).Model(&models.OrganizationRole{}).Where("organization_id = ?", orgID)

	// Apply filters
	query = r.applyRoleFilters(query, filters)

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count roles by organization: %w", err)
	}

	// Apply pagination and sorting
	query = r.applyPaginationAndSorting(query, filters.PaginationRequest, filters.FilterRequest)

	// Execute query
	if err := query.Find(&roles).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list roles by organization: %w", err)
	}

	return roles, total, nil
}

// ===============================
// DEFAULT ROLES MANAGEMENT
// ===============================

func (r *roleRepository) GetDefaultRoles(ctx context.Context, orgID uuid.UUID) ([]models.OrganizationRole, error) {
	var roles []models.OrganizationRole
	if err := r.db.WithContext(ctx).Where("organization_id = ? AND is_default = ?", orgID, true).Find(&roles).Error; err != nil {
		return nil, fmt.Errorf("failed to get default roles: %w", err)
	}
	return roles, nil
}

func (r *roleRepository) GetNonDefaultRoles(ctx context.Context, orgID uuid.UUID) ([]models.OrganizationRole, error) {
	var roles []models.OrganizationRole
	if err := r.db.WithContext(ctx).Where("organization_id = ? AND is_default = ?", orgID, false).Find(&roles).Error; err != nil {
		return nil, fmt.Errorf("failed to get non-default roles: %w", err)
	}
	return roles, nil
}

func (r *roleRepository) CreateDefaultRoles(ctx context.Context, orgID uuid.UUID, createdBy uuid.UUID) ([]models.OrganizationRole, error) {
	// Define default roles for real estate organizations
	defaultRoles := []models.OrganizationRole{
		{
			ID:             uuid.New(),
			OrganizationID: orgID,
			Name:           "Admin",
			Description:    stringPtr("Administrator with full access to organization management"),
			IsDefault:      true,
			CreatedAt:      time.Now(),
			CreatedBy:      createdBy,
			UpdatedAt:      time.Now(),
		},
		{
			ID:             uuid.New(),
			OrganizationID: orgID,
			Name:           "Manager",
			Description:    stringPtr("Manager with access to employee and property management"),
			IsDefault:      true,
			CreatedAt:      time.Now(),
			CreatedBy:      createdBy,
			UpdatedAt:      time.Now(),
		},
		{
			ID:             uuid.New(),
			OrganizationID: orgID,
			Name:           "Agent",
			Description:    stringPtr("Real estate agent with access to property listings and client management"),
			IsDefault:      true,
			CreatedAt:      time.Now(),
			CreatedBy:      createdBy,
			UpdatedAt:      time.Now(),
		},
		{
			ID:             uuid.New(),
			OrganizationID: orgID,
			Name:           "Assistant",
			Description:    stringPtr("Administrative assistant with limited access"),
			IsDefault:      true,
			CreatedAt:      time.Now(),
			CreatedBy:      createdBy,
			UpdatedAt:      time.Now(),
		},
	}

	// Create roles in a transaction
	var createdRoles []models.OrganizationRole
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, role := range defaultRoles {
			if err := tx.Create(&role).Error; err != nil {
				return fmt.Errorf("failed to create default role %s: %w", role.Name, err)
			}
			createdRoles = append(createdRoles, role)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdRoles, nil
}

// ===============================
// RELATIONSHIPS
// ===============================

func (r *roleRepository) GetWithEmployees(ctx context.Context, id uuid.UUID) (*models.OrganizationRole, error) {
	var role models.OrganizationRole
	if err := r.db.WithContext(ctx).Preload("EmployeeRoles.Employee").Where("id = ?", id).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("role not found")
		}
		return nil, fmt.Errorf("failed to get role with employees: %w", err)
	}
	return &role, nil
}

func (r *roleRepository) GetWithOrganization(ctx context.Context, id uuid.UUID) (*models.OrganizationRole, error) {
	var role models.OrganizationRole
	if err := r.db.WithContext(ctx).Preload("Organization").Where("id = ?", id).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("role not found")
		}
		return nil, fmt.Errorf("failed to get role with organization: %w", err)
	}
	return &role, nil
}

// ===============================
// BUSINESS LOGIC SUPPORT
// ===============================

func (r *roleRepository) ExistsByName(ctx context.Context, orgID uuid.UUID, name string, excludeID *uuid.UUID) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.OrganizationRole{}).
		Where("organization_id = ? AND LOWER(name) = LOWER(?)", orgID, name)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check role name existence: %w", err)
	}
	return count > 0, nil
}

func (r *roleRepository) CountByOrganization(ctx context.Context, orgID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.OrganizationRole{}).Where("organization_id = ?", orgID).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count roles by organization: %w", err)
	}
	return count, nil
}

func (r *roleRepository) GetRoleUsageStats(ctx context.Context, roleID uuid.UUID) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Count active employees with this role
	var activeEmployees int64
	if err := r.db.WithContext(ctx).Table("employee_roles").
		Joins("JOIN employees ON employee_roles.employee_id = employees.id").
		Where("employee_roles.role_id = ? AND employees.status = ?", roleID, "active").
		Count(&activeEmployees).Error; err != nil {
		return nil, fmt.Errorf("failed to count active employees: %w", err)
	}

	// Count total employees with this role
	var totalEmployees int64
	if err := r.db.WithContext(ctx).Model(&models.EmployeeRole{}).
		Where("role_id = ?", roleID).Count(&totalEmployees).Error; err != nil {
		return nil, fmt.Errorf("failed to count total employees: %w", err)
	}

	// Count primary role assignments
	var primaryAssignments int64
	if err := r.db.WithContext(ctx).Model(&models.EmployeeRole{}).
		Where("role_id = ? AND is_primary = ?", roleID, true).Count(&primaryAssignments).Error; err != nil {
		return nil, fmt.Errorf("failed to count primary assignments: %w", err)
	}

	stats["active_employees"] = activeEmployees
	stats["total_employees"] = totalEmployees
	stats["primary_assignments"] = primaryAssignments

	return stats, nil
}

func (r *roleRepository) CanDelete(ctx context.Context, id uuid.UUID) (bool, error) {
	// Get role details
	var role models.OrganizationRole
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&role).Error; err != nil {
		return false, fmt.Errorf("role not found")
	}

	// Cannot delete default roles
	if role.IsDefault {
		return false, nil
	}

	// Check if there are employees assigned to this role
	var employeeCount int64
	if err := r.db.WithContext(ctx).Model(&models.EmployeeRole{}).
		Where("role_id = ?", id).Count(&employeeCount).Error; err != nil {
		return false, fmt.Errorf("failed to count employees with role: %w", err)
	}

	if employeeCount > 0 {
		return false, nil // Cannot delete role with assigned employees
	}

	return true, nil
}

// ===============================
// HELPER METHODS
// ===============================

func (r *roleRepository) applyRoleFilters(query *gorm.DB, filters dto.RoleFiltersRequest) *gorm.DB {
	// Search filter
	if filters.Search != "" {
		searchTerm := "%" + filters.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", searchTerm, searchTerm)
	}

	// Organization filter
	if filters.OrganizationID != uuid.Nil {
		query = query.Where("organization_id = ?", filters.OrganizationID)
	}

	// Is default filter
	if filters.IsDefault != nil {
		query = query.Where("is_default = ?", *filters.IsDefault)
	}

	// Created date range
	if filters.CreatedAt != nil {
		query = query.Where("created_at >= ?", *filters.CreatedAt)
	}

	if filters.UpdatedAt != nil {
		query = query.Where("updated_at >= ?", *filters.UpdatedAt)
	}

	return query
}

func (r *roleRepository) applyPaginationAndSorting(query *gorm.DB, pagination dto.PaginationRequest, filter dto.FilterRequest) *gorm.DB {
	// Apply default sorting (default roles first, then by name)
	query = query.Order("is_default DESC, name ASC")

	// Apply pagination
	if pagination.Page > 0 && pagination.PerPage > 0 {
		offset := (pagination.Page - 1) * pagination.PerPage
		query = query.Offset(offset).Limit(pagination.PerPage)
	}

	return query
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
