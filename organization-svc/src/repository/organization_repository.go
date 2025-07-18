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

// organizationRepository implements OrganizationRepository interface
type organizationRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewOrganizationRepository creates a new organization repository instance
func NewOrganizationRepository(db *gorm.DB, logger *zap.Logger) OrganizationRepository {
	return &organizationRepository{db: db, logger: logger}
}

// ===============================
// BASIC CRUD OPERATIONS
// ===============================

func (r *organizationRepository) Create(ctx context.Context, org *models.Organization) (*models.Organization, error) {
	if err := r.db.WithContext(ctx).Create(org).Error; err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}
	return org, nil
}

func (r *organizationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	var org models.Organization
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&org).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("organization not found")
		}
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}
	return &org, nil
}

func (r *organizationRepository) GetBySlug(ctx context.Context, slug string) (*models.Organization, error) {
	var org models.Organization
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&org).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("organization not found")
		}
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}
	return &org, nil
}

func (r *organizationRepository) Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.Organization, error) {
	var org models.Organization

	// First check if organization exists
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&org).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("organization not found")
		}
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}

	// Add updated_at timestamp
	updates["updated_at"] = time.Now()

	// Perform update
	if err := r.db.WithContext(ctx).Model(&org).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update organization: %w", err)
	}

	// Return updated organization
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&org).Error; err != nil {
		return nil, fmt.Errorf("failed to get updated organization: %w", err)
	}

	return &org, nil
}

func (r *organizationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Unscoped().Delete(&models.Organization{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete organization: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("organization not found")
	}
	return nil
}

func (r *organizationRepository) SoftDelete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	updates := map[string]interface{}{
		"deleted_at": time.Now(),
		"deleted_by": deletedBy,
	}

	result := r.db.WithContext(ctx).Model(&models.Organization{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to soft delete organization: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("organization not found")
	}
	return nil
}

// ===============================
// LISTING AND FILTERING
// ===============================

func (r *organizationRepository) List(ctx context.Context, filters dto.OrganizationFiltersRequest) ([]models.Organization, int64, error) {
	var orgs []models.Organization
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Organization{})

	// Apply filters
	query = r.applyOrganizationFilters(query, filters)

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count organizations: %w", err)
	}

	// Apply pagination and sorting
	query = r.applyPaginationAndSorting(query, filters.PaginationRequest, filters.FilterRequest)

	// Execute query
	if err := query.Find(&orgs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list organizations: %w", err)
	}

	return orgs, total, nil
}

func (r *organizationRepository) ListByOwner(ctx context.Context, userID uuid.UUID, filters dto.OrganizationFiltersRequest) ([]models.Organization, int64, error) {
	var orgs []models.Organization
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Organization{}).
		Joins("JOIN organization_owners ON organizations.id = organization_owners.organization_id").
		Where("organization_owners.person_id = ?", userID)

	// Apply filters
	query = r.applyOrganizationFilters(query, filters)

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count organizations by owner: %w", err)
	}

	// Apply pagination and sorting
	query = r.applyPaginationAndSorting(query, filters.PaginationRequest, filters.FilterRequest)

	// Execute query
	if err := query.Find(&orgs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list organizations by owner: %w", err)
	}

	return orgs, total, nil
}

func (r *organizationRepository) ListByEmployee(ctx context.Context, userID uuid.UUID, filters dto.OrganizationFiltersRequest) ([]models.Organization, int64, error) {
	var orgs []models.Organization
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Organization{}).
		Joins("JOIN employees ON organizations.id = employees.organization_id").
		Where("employees.user_id = ? AND employees.status = ?", userID, "active")

	// Apply filters
	query = r.applyOrganizationFilters(query, filters)

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count organizations by employee: %w", err)
	}

	// Apply pagination and sorting
	query = r.applyPaginationAndSorting(query, filters.PaginationRequest, filters.FilterRequest)

	// Execute query
	if err := query.Find(&orgs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list organizations by employee: %w", err)
	}

	return orgs, total, nil
}

// ===============================
// STATUS MANAGEMENT
// ===============================

func (r *organizationRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.OrganizationStatus, updatedBy uuid.UUID) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_by": updatedBy,
		"updated_at": time.Now(),
	}

	result := r.db.WithContext(ctx).Model(&models.Organization{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update organization status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("organization not found")
	}
	return nil
}

func (r *organizationRepository) GetByStatus(ctx context.Context, status models.OrganizationStatus, limit int) ([]models.Organization, error) {
	var orgs []models.Organization
	query := r.db.WithContext(ctx).Where("status = ?", status)

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&orgs).Error; err != nil {
		return nil, fmt.Errorf("failed to get organizations by status: %w", err)
	}

	return orgs, nil
}

// ===============================
// VERIFICATION
// ===============================

func (r *organizationRepository) MarkAsVerified(ctx context.Context, id uuid.UUID, verifiedBy uuid.UUID) error {
	updates := map[string]interface{}{
		"is_verified": true,
		"verified_at": time.Now(),
		"verified_by": verifiedBy,
		"updated_at":  time.Now(),
	}

	result := r.db.WithContext(ctx).Model(&models.Organization{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to mark organization as verified: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("organization not found")
	}
	return nil
}

func (r *organizationRepository) GetUnverifiedOrganizations(ctx context.Context, limit int) ([]models.Organization, error) {
	var orgs []models.Organization
	query := r.db.WithContext(ctx).Where("is_verified = ?", false)

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&orgs).Error; err != nil {
		return nil, fmt.Errorf("failed to get unverified organizations: %w", err)
	}

	return orgs, nil
}

// ===============================
// RELATIONSHIPS
// ===============================

func (r *organizationRepository) GetWithSettings(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	var org models.Organization
	if err := r.db.WithContext(ctx).Preload("Settings").Where("id = ?", id).First(&org).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("organization not found")
		}
		return nil, fmt.Errorf("failed to get organization with settings: %w", err)
	}
	return &org, nil
}

func (r *organizationRepository) GetWithOwner(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	var org models.Organization
	if err := r.db.WithContext(ctx).Preload("Owners").Where("id = ?", id).First(&org).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("organization not found")
		}
		return nil, fmt.Errorf("failed to get organization with owner: %w", err)
	}
	return &org, nil
}

func (r *organizationRepository) GetWithEmployees(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	var org models.Organization
	if err := r.db.WithContext(ctx).Preload("Employees").Where("id = ?", id).First(&org).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("organization not found")
		}
		return nil, fmt.Errorf("failed to get organization with employees: %w", err)
	}
	return &org, nil
}

func (r *organizationRepository) GetWithBranches(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	var org models.Organization
	if err := r.db.WithContext(ctx).Preload("Branches").Where("id = ?", id).First(&org).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("organization not found")
		}
		return nil, fmt.Errorf("failed to get organization with branches: %w", err)
	}
	return &org, nil
}

func (r *organizationRepository) GetWithRoles(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	var org models.Organization
	if err := r.db.WithContext(ctx).Preload("Roles").Where("id = ?", id).First(&org).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("organization not found")
		}
		return nil, fmt.Errorf("failed to get organization with roles: %w", err)
	}
	return &org, nil
}

// ===============================
// BUSINESS LOGIC SUPPORT
// ===============================

func (r *organizationRepository) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Organization{}).Where("slug = ?", slug).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check organization slug existence: %w", err)
	}
	return count > 0, nil
}

func (r *organizationRepository) ExistsByName(ctx context.Context, name string, excludeID *uuid.UUID) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.Organization{}).Where("LOWER(name) = LOWER(?)", name)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check organization name existence: %w", err)
	}
	return count > 0, nil
}

func (r *organizationRepository) CountByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.Organization{}).
		Joins("LEFT JOIN organization_owners ON organizations.id = organization_owners.organization_id").
		Joins("LEFT JOIN employees ON organizations.id = employees.organization_id").
		Where("organization_owners.person_id = ? OR employees.user_id = ?", userID, userID)

	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count organizations by user: %w", err)
	}
	return count, nil
}

func (r *organizationRepository) GetActiveCount(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Organization{}).Where("status = ?", models.OrgStatusActive).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count active organizations: %w", err)
	}
	return count, nil
}

// ===============================
// HELPER METHODS
// ===============================

func (r *organizationRepository) applyOrganizationFilters(query *gorm.DB, filters dto.OrganizationFiltersRequest) *gorm.DB {
	// Search filter
	if filters.Search != "" {
		searchTerm := "%" + filters.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ? OR slug ILIKE ?", searchTerm, searchTerm, searchTerm)
	}

	// Status filter
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}

	// Type filter
	if filters.Type != "" {
		query = query.Where("type = ?", filters.Type)
	}

	// Verification filter
	if filters.IsVerified != nil {
		query = query.Where("is_verified = ?", *filters.IsVerified)
	}

	// Created date range (from FilterRequest in DTO structure)
	if filters.CreatedAt != nil {
		query = query.Where("created_at >= ?", *filters.CreatedAt)
	}

	if filters.UpdatedAt != nil {
		query = query.Where("updated_at >= ?", *filters.UpdatedAt)
	}

	return query
}

func (r *organizationRepository) applyPaginationAndSorting(query *gorm.DB, pagination dto.PaginationRequest, filter dto.FilterRequest) *gorm.DB {
	// Apply default sorting
	query = query.Order("created_at DESC")

	// Apply pagination
	if pagination.Page > 0 && pagination.PerPage > 0 {
		offset := (pagination.Page - 1) * pagination.PerPage
		query = query.Offset(offset).Limit(pagination.PerPage)
	}

	return query
}
