package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rem-gestion/api-suite/organization/src/dto"
	"github.com/rem-gestion/api-suite/organization/src/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ownerRepository implements OwnerRepository interface
type ownerRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewOwnerRepository creates a new owner repository instance
func NewOwnerRepository(db *gorm.DB, logger *zap.Logger) OwnerRepository {
	return &ownerRepository{db: db, logger: logger}
}

// ===============================
// BASIC CRUD OPERATIONS
// ===============================

func (r *ownerRepository) Create(ctx context.Context, owner *models.OrganizationOwner) (*models.OrganizationOwner, error) {
	if err := r.db.WithContext(ctx).Create(owner).Error; err != nil {
		return nil, fmt.Errorf("failed to create owner: %w", err)
	}
	return owner, nil
}

func (r *ownerRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.OrganizationOwner, error) {
	var owner models.OrganizationOwner
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&owner).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("owner not found")
		}
		return nil, fmt.Errorf("failed to get owner: %w", err)
	}
	return &owner, nil
}

func (r *ownerRepository) GetByOrganization(ctx context.Context, orgID uuid.UUID) (*models.OrganizationOwner, error) {
	var owner models.OrganizationOwner
	if err := r.db.WithContext(ctx).Where("organization_id = ?", orgID).First(&owner).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("owner not found")
		}
		return nil, fmt.Errorf("failed to get owner by organization: %w", err)
	}
	return &owner, nil
}

func (r *ownerRepository) GetByPersonAndOrg(ctx context.Context, personID, orgID uuid.UUID) (*models.OrganizationOwner, error) {
	var owner models.OrganizationOwner
	if err := r.db.WithContext(ctx).Where("person_id = ? AND organization_id = ?", personID, orgID).First(&owner).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("owner not found")
		}
		return nil, fmt.Errorf("failed to get owner by person and organization: %w", err)
	}
	return &owner, nil
}

func (r *ownerRepository) Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.OrganizationOwner, error) {
	var owner models.OrganizationOwner

	// First check if owner exists
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&owner).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("owner not found")
		}
		return nil, fmt.Errorf("failed to get owner: %w", err)
	}

	// Add updated_at timestamp
	updates["updated_at"] = time.Now()

	// Perform update
	if err := r.db.WithContext(ctx).Model(&owner).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update owner: %w", err)
	}

	// Return updated owner
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&owner).Error; err != nil {
		return nil, fmt.Errorf("failed to get updated owner: %w", err)
	}

	return &owner, nil
}

func (r *ownerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Unscoped().Delete(&models.OrganizationOwner{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete owner: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("owner not found")
	}
	return nil
}

// ===============================
// LISTING AND FILTERING
// ===============================

func (r *ownerRepository) List(ctx context.Context, filters dto.OwnerFiltersRequest) ([]models.OrganizationOwner, int64, error) {
	var owners []models.OrganizationOwner
	var total int64

	query := r.db.WithContext(ctx).Model(&models.OrganizationOwner{})

	// Apply filters
	query = r.applyOwnerFilters(query, filters)

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count owners: %w", err)
	}

	// Apply pagination and sorting
	query = r.applyPaginationAndSorting(query, filters.PaginationRequest, filters.FilterRequest)

	// Execute query
	if err := query.Find(&owners).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list owners: %w", err)
	}

	return owners, total, nil
}

func (r *ownerRepository) ListByOrganization(ctx context.Context, orgID uuid.UUID, filters dto.OwnerFiltersRequest) ([]models.OrganizationOwner, int64, error) {
	var owners []models.OrganizationOwner
	var total int64

	query := r.db.WithContext(ctx).Model(&models.OrganizationOwner{}).Where("organization_id = ?", orgID)

	// Apply filters
	query = r.applyOwnerFilters(query, filters)

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count owners by organization: %w", err)
	}

	// Apply pagination and sorting
	query = r.applyPaginationAndSorting(query, filters.PaginationRequest, filters.FilterRequest)

	// Execute query
	if err := query.Find(&owners).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list owners by organization: %w", err)
	}

	return owners, total, nil
}

// ===============================
// OWNERSHIP MANAGEMENT
// ===============================

func (r *ownerRepository) GetFounders(ctx context.Context, orgID uuid.UUID) ([]models.OrganizationOwner, error) {
	var owners []models.OrganizationOwner
	if err := r.db.WithContext(ctx).Where("organization_id = ? AND is_founder = ?", orgID, true).Find(&owners).Error; err != nil {
		return nil, fmt.Errorf("failed to get founders: %w", err)
	}
	return owners, nil
}

func (r *ownerRepository) GetMajorityOwner(ctx context.Context, orgID uuid.UUID) (*models.OrganizationOwner, error) {
	var owner models.OrganizationOwner
	if err := r.db.WithContext(ctx).Where("organization_id = ? AND ownership_pct > ?", orgID, 50.0).First(&owner).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("majority owner not found")
		}
		return nil, fmt.Errorf("failed to get majority owner: %w", err)
	}
	return &owner, nil
}

func (r *ownerRepository) UpdateOwnershipPercentage(ctx context.Context, id uuid.UUID, percentage float64, updatedBy uuid.UUID) error {
	updates := map[string]interface{}{
		"ownership_pct": percentage,
		"updated_by":    updatedBy,
		"updated_at":    time.Now(),
	}

	result := r.db.WithContext(ctx).Model(&models.OrganizationOwner{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update ownership percentage: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("owner not found")
	}
	return nil
}

func (r *ownerRepository) TransferOwnership(ctx context.Context, fromOwnerID, toOwnerID uuid.UUID, percentage float64, updatedBy uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get current ownership percentages
		var fromOwner, toOwner models.OrganizationOwner

		if err := tx.Where("id = ?", fromOwnerID).First(&fromOwner).Error; err != nil {
			return fmt.Errorf("from owner not found: %w", err)
		}

		if err := tx.Where("id = ?", toOwnerID).First(&toOwner).Error; err != nil {
			return fmt.Errorf("to owner not found: %w", err)
		}

		// Validate same organization
		if fromOwner.OrganizationID != toOwner.OrganizationID {
			return fmt.Errorf("owners must belong to the same organization")
		}

		// Calculate new percentages
		fromPercentage := 0.0
		if fromOwner.OwnershipPct != nil {
			fromPercentage = *fromOwner.OwnershipPct
		}

		toPercentage := 0.0
		if toOwner.OwnershipPct != nil {
			toPercentage = *toOwner.OwnershipPct
		}

		if fromPercentage < percentage {
			return fmt.Errorf("insufficient ownership percentage to transfer")
		}

		newFromPercentage := fromPercentage - percentage
		newToPercentage := toPercentage + percentage

		// Update from owner
		if err := tx.Model(&fromOwner).Updates(map[string]interface{}{
			"ownership_pct": newFromPercentage,
			"updated_by":    updatedBy,
			"updated_at":    time.Now(),
		}).Error; err != nil {
			return fmt.Errorf("failed to update from owner: %w", err)
		}

		// Update to owner
		if err := tx.Model(&toOwner).Updates(map[string]interface{}{
			"ownership_pct": newToPercentage,
			"updated_by":    updatedBy,
			"updated_at":    time.Now(),
		}).Error; err != nil {
			return fmt.Errorf("failed to update to owner: %w", err)
		}

		return nil
	})
}

// ===============================
// VALIDATION AND STATISTICS
// ===============================

func (r *ownerRepository) GetTotalOwnership(ctx context.Context, orgID uuid.UUID) (float64, error) {
	var result struct {
		Total *float64 `json:"total"`
	}

	if err := r.db.WithContext(ctx).Model(&models.OrganizationOwner{}).
		Select("SUM(ownership_pct) as total").
		Where("organization_id = ?", orgID).
		Scan(&result).Error; err != nil {
		return 0, fmt.Errorf("failed to calculate total ownership: %w", err)
	}

	if result.Total == nil {
		return 0, nil
	}

	return *result.Total, nil
}

func (r *ownerRepository) ValidateOwnershipPercentages(ctx context.Context, orgID uuid.UUID) (bool, error) {
	total, err := r.GetTotalOwnership(ctx, orgID)
	if err != nil {
		return false, err
	}

	// Allow small floating point precision errors
	return total >= 99.99 && total <= 100.01, nil
}

func (r *ownerRepository) GetOwnershipSummary(ctx context.Context, orgID uuid.UUID) (*dto.OwnershipSummary, error) {
	summary := &dto.OwnershipSummary{
		LastUpdated: time.Now(),
	}

	// Get total ownership
	total, err := r.GetTotalOwnership(ctx, orgID)
	if err != nil {
		return nil, err
	}
	summary.TotalOwnership = total
	summary.AllocatedOwnership = total
	summary.UnallocatedOwnership = 100.0 - total

	// Count owners
	var ownerCount int64
	if err := r.db.WithContext(ctx).Model(&models.OrganizationOwner{}).
		Where("organization_id = ?", orgID).Count(&ownerCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count owners: %w", err)
	}
	summary.TotalOwners = int(ownerCount)

	// Count founders
	var founderCount int64
	if err := r.db.WithContext(ctx).Model(&models.OrganizationOwner{}).
		Where("organization_id = ? AND is_founder = ?", orgID, true).Count(&founderCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count founders: %w", err)
	}
	summary.TotalFounders = int(founderCount)

	// Get majority owner
	if majorityOwner, err := r.GetMajorityOwner(ctx, orgID); err == nil {
		summary.MajorityOwner = &dto.OwnerBasicResponse{
			ID:             majorityOwner.ID,
			OrganizationID: majorityOwner.OrganizationID,
			PersonID:       majorityOwner.PersonID,
			UserID:         majorityOwner.UserID,
			IsFounder:      majorityOwner.IsFounder,
			OwnershipPct:   majorityOwner.OwnershipPct,
		}
	}

	// Create ownership breakdown (this is a simplified version)
	summary.OwnershipBreakdown = []dto.OwnershipBreakdown{
		{
			Range:      "0-25%",
			Count:      0,
			Percentage: 0,
		},
		{
			Range:      "25-50%",
			Count:      0,
			Percentage: 0,
		},
		{
			Range:      "50-75%",
			Count:      0,
			Percentage: 0,
		},
		{
			Range:      "75-100%",
			Count:      0,
			Percentage: 0,
		},
	}

	return summary, nil
}

// ===============================
// RELATIONSHIPS
// ===============================

func (r *ownerRepository) GetWithOrganization(ctx context.Context, id uuid.UUID) (*models.OrganizationOwner, error) {
	var owner models.OrganizationOwner
	if err := r.db.WithContext(ctx).Preload("Organization").Where("id = ?", id).First(&owner).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("owner not found")
		}
		return nil, fmt.Errorf("failed to get owner with organization: %w", err)
	}
	return &owner, nil
}

// ===============================
// BUSINESS LOGIC SUPPORT
// ===============================

func (r *ownerRepository) ExistsByPersonAndOrg(ctx context.Context, personID, orgID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.OrganizationOwner{}).
		Where("person_id = ? AND organization_id = ?", personID, orgID).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check owner existence: %w", err)
	}
	return count > 0, nil
}

func (r *ownerRepository) CountByOrganization(ctx context.Context, orgID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.OrganizationOwner{}).Where("organization_id = ?", orgID).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count owners by organization: %w", err)
	}
	return count, nil
}

func (r *ownerRepository) IsUserOwnerOfOrg(ctx context.Context, userID, orgID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.OrganizationOwner{}).
		Where("user_id = ? AND organization_id = ?", userID, orgID).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check user ownership: %w", err)
	}
	return count > 0, nil
}

func (r *ownerRepository) CanRemoveOwner(ctx context.Context, id uuid.UUID) (bool, error) {
	// Get owner details
	var owner models.OrganizationOwner
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&owner).Error; err != nil {
		return false, fmt.Errorf("owner not found")
	}

	// Check if it's the only owner in the organization
	var ownerCount int64
	if err := r.db.WithContext(ctx).Model(&models.OrganizationOwner{}).
		Where("organization_id = ?", owner.OrganizationID).Count(&ownerCount).Error; err != nil {
		return false, fmt.Errorf("failed to count owners: %w", err)
	}

	if ownerCount <= 1 {
		return false, nil // Cannot remove the last owner
	}

	return true, nil
}

// ===============================
// HELPER METHODS
// ===============================

func (r *ownerRepository) applyOwnerFilters(query *gorm.DB, filters dto.OwnerFiltersRequest) *gorm.DB {
	// Organization filter
	if filters.OrganizationID != uuid.Nil {
		query = query.Where("organization_id = ?", filters.OrganizationID)
	}

	// Is founder filter
	if filters.IsFounder != nil {
		query = query.Where("is_founder = ?", *filters.IsFounder)
	}

	// Ownership percentage range
	if filters.MinOwnership != nil {
		query = query.Where("ownership_pct >= ?", *filters.MinOwnership)
	}

	if filters.MaxOwnership != nil {
		query = query.Where("ownership_pct <= ?", *filters.MaxOwnership)
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

func (r *ownerRepository) applyPaginationAndSorting(query *gorm.DB, pagination dto.PaginationRequest, filter dto.FilterRequest) *gorm.DB {
	// Apply default sorting (founder first, then by ownership percentage descending)
	query = query.Order("is_founder DESC, ownership_pct DESC")

	// Apply pagination
	if pagination.Page > 0 && pagination.PerPage > 0 {
		offset := (pagination.Page - 1) * pagination.PerPage
		query = query.Offset(offset).Limit(pagination.PerPage)
	}

	return query
}
