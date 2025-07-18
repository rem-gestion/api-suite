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

// branchRepository implements BranchRepository interface

type branchRepository struct {
	logger *zap.Logger
	db     *gorm.DB
}

// NewBranchRepository creates a new branch repository instance
func NewBranchRepository(db *gorm.DB, logger *zap.Logger) BranchRepository {
	return &branchRepository{db: db, logger: logger}
}

// ===============================
// BASIC CRUD OPERATIONS
// ===============================

func (r *branchRepository) Create(ctx context.Context, branch *models.OrganizationBranch) (*models.OrganizationBranch, error) {
	if err := r.db.WithContext(ctx).Create(branch).Error; err != nil {
		return nil, fmt.Errorf("failed to create branch: %w", err)
	}
	return branch, nil
}

func (r *branchRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.OrganizationBranch, error) {
	var branch models.OrganizationBranch
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&branch).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("branch not found")
		}
		return nil, fmt.Errorf("failed to get branch: %w", err)
	}
	return &branch, nil
}

func (r *branchRepository) Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.OrganizationBranch, error) {
	var branch models.OrganizationBranch

	// First check if branch exists
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&branch).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("branch not found")
		}
		return nil, fmt.Errorf("failed to get branch: %w", err)
	}

	// Add updated_at timestamp
	updates["updated_at"] = time.Now()

	// Perform update
	if err := r.db.WithContext(ctx).Model(&branch).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update branch: %w", err)
	}

	// Return updated branch
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&branch).Error; err != nil {
		return nil, fmt.Errorf("failed to get updated branch: %w", err)
	}

	return &branch, nil
}

func (r *branchRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Unscoped().Delete(&models.OrganizationBranch{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete branch: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("branch not found")
	}
	return nil
}

func (r *branchRepository) SoftDelete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	updates := map[string]interface{}{
		"deleted_at": time.Now(),
		"updated_by": deletedBy,
		"updated_at": time.Now(),
	}

	result := r.db.WithContext(ctx).Model(&models.OrganizationBranch{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to soft delete branch: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("branch not found")
	}
	return nil
}

// ===============================
// LISTING AND FILTERING
// ===============================

func (r *branchRepository) List(ctx context.Context, filters dto.BranchFiltersRequest) ([]models.OrganizationBranch, int64, error) {
	var branches []models.OrganizationBranch
	var total int64

	query := r.db.WithContext(ctx).Model(&models.OrganizationBranch{})

	// Apply filters
	query = r.applyBranchFilters(query, filters)

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count branches: %w", err)
	}

	// Apply pagination and sorting
	query = r.applyPaginationAndSorting(query, filters.PaginationRequest, filters.FilterRequest)

	// Execute query
	if err := query.Find(&branches).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list branches: %w", err)
	}

	return branches, total, nil
}

func (r *branchRepository) ListByOrganization(ctx context.Context, orgID uuid.UUID, filters dto.BranchFiltersRequest) ([]models.OrganizationBranch, int64, error) {
	var branches []models.OrganizationBranch
	var total int64

	query := r.db.WithContext(ctx).Model(&models.OrganizationBranch{}).Where("organization_id = ?", orgID)

	// Apply filters
	query = r.applyBranchFilters(query, filters)

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count branches by organization: %w", err)
	}

	// Apply pagination and sorting
	query = r.applyPaginationAndSorting(query, filters.PaginationRequest, filters.FilterRequest)

	// Execute query
	if err := query.Find(&branches).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list branches by organization: %w", err)
	}

	return branches, total, nil
}

// ===============================
// MAIN BRANCH MANAGEMENT
// ===============================

func (r *branchRepository) GetMainBranch(ctx context.Context, orgID uuid.UUID) (*models.OrganizationBranch, error) {
	var branch models.OrganizationBranch
	if err := r.db.WithContext(ctx).Where("organization_id = ? AND is_main = ?", orgID, true).First(&branch).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("main branch not found")
		}
		return nil, fmt.Errorf("failed to get main branch: %w", err)
	}
	return &branch, nil
}

func (r *branchRepository) SetAsMain(ctx context.Context, id uuid.UUID, updatedBy uuid.UUID) error {
	// First get the branch to find its organization
	var branch models.OrganizationBranch
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&branch).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("branch not found")
		}
		return fmt.Errorf("failed to get branch: %w", err)
	}

	// Use transaction to ensure atomicity
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Clear current main branch
		if err := tx.Model(&models.OrganizationBranch{}).
			Where("organization_id = ? AND is_main = ?", branch.OrganizationID, true).
			Updates(map[string]interface{}{
				"is_main":    false,
				"updated_by": updatedBy,
				"updated_at": time.Now(),
			}).Error; err != nil {
			return fmt.Errorf("failed to clear current main branch: %w", err)
		}

		// Set new main branch
		if err := tx.Model(&models.OrganizationBranch{}).
			Where("id = ?", id).
			Updates(map[string]interface{}{
				"is_main":    true,
				"updated_by": updatedBy,
				"updated_at": time.Now(),
			}).Error; err != nil {
			return fmt.Errorf("failed to set new main branch: %w", err)
		}

		return nil
	})
}

func (r *branchRepository) ClearMainStatus(ctx context.Context, orgID uuid.UUID) error {
	result := r.db.WithContext(ctx).Model(&models.OrganizationBranch{}).
		Where("organization_id = ? AND is_main = ?", orgID, true).
		Updates(map[string]interface{}{
			"is_main":    false,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("failed to clear main status: %w", result.Error)
	}

	return nil
}

// ===============================
// RELATIONSHIPS
// ===============================

func (r *branchRepository) GetWithOrganization(ctx context.Context, id uuid.UUID) (*models.OrganizationBranch, error) {
	var branch models.OrganizationBranch
	if err := r.db.WithContext(ctx).Preload("Organization").Where("id = ?", id).First(&branch).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("branch not found")
		}
		return nil, fmt.Errorf("failed to get branch with organization: %w", err)
	}
	return &branch, nil
}

// ===============================
// BUSINESS LOGIC SUPPORT
// ===============================

func (r *branchRepository) CountByOrganization(ctx context.Context, orgID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.OrganizationBranch{}).Where("organization_id = ?", orgID).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count branches by organization: %w", err)
	}
	return count, nil
}

func (r *branchRepository) ExistsByName(ctx context.Context, orgID uuid.UUID, name string, excludeID *uuid.UUID) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.OrganizationBranch{}).
		Where("organization_id = ? AND LOWER(display_name) = LOWER(?)", orgID, name)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check branch name existence: %w", err)
	}
	return count > 0, nil
}

func (r *branchRepository) HasMainBranch(ctx context.Context, orgID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.OrganizationBranch{}).
		Where("organization_id = ? AND is_main = ?", orgID, true).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check main branch existence: %w", err)
	}
	return count > 0, nil
}

func (r *branchRepository) CanDelete(ctx context.Context, id uuid.UUID) (bool, error) {
	// Get branch details
	var branch models.OrganizationBranch
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&branch).Error; err != nil {
		return false, fmt.Errorf("branch not found")
	}

	// Check if it's the only branch in the organization
	var branchCount int64
	if err := r.db.WithContext(ctx).Model(&models.OrganizationBranch{}).
		Where("organization_id = ?", branch.OrganizationID).Count(&branchCount).Error; err != nil {
		return false, fmt.Errorf("failed to count branches: %w", err)
	}

	if branchCount <= 1 {
		return false, nil // Cannot delete the last branch
	}

	// Check if there are active employees assigned to this branch
	var employeeCount int64
	if err := r.db.WithContext(ctx).Model(&models.Employee{}).
		Where("organization_id = ? AND status IN ?", branch.OrganizationID, []string{"active", "inactive"}).
		Count(&employeeCount).Error; err != nil {
		return false, fmt.Errorf("failed to count employees in branch: %w", err)
	}

	if employeeCount > 0 {
		return false, nil // Cannot delete branch with active employees
	}

	return true, nil
}

// ===============================
// HELPER METHODS
// ===============================

func (r *branchRepository) applyBranchFilters(query *gorm.DB, filters dto.BranchFiltersRequest) *gorm.DB {
	// Search filter
	if filters.Search != "" {
		searchTerm := "%" + filters.Search + "%"
		query = query.Where("display_name ILIKE ? OR phone ILIKE ? OR email ILIKE ?", searchTerm, searchTerm, searchTerm)
	}

	// Organization filter
	if filters.OrganizationID != uuid.Nil {
		query = query.Where("organization_id = ?", filters.OrganizationID)
	}

	// Main branch filter
	if filters.IsMain != nil {
		query = query.Where("is_main = ?", *filters.IsMain)
	}

	// Has address filter
	if filters.HasAddress != nil {
		if *filters.HasAddress {
			query = query.Where("address_id IS NOT NULL")
		} else {
			query = query.Where("address_id IS NULL")
		}
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

func (r *branchRepository) applyPaginationAndSorting(query *gorm.DB, pagination dto.PaginationRequest, filter dto.FilterRequest) *gorm.DB {
	// Apply default sorting (main branch first, then by name)
	query = query.Order("is_main DESC, display_name ASC")

	// Apply pagination
	if pagination.Page > 0 && pagination.PerPage > 0 {
		offset := (pagination.Page - 1) * pagination.PerPage
		query = query.Offset(offset).Limit(pagination.PerPage)
	}

	return query
}
