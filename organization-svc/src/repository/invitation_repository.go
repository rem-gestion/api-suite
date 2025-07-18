package repository

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/rem-gestion/api-suite/organization/src/dto"
	"github.com/rem-gestion/api-suite/organization/src/models"
)

// invitationRepository implements InvitationRepository interface

type invitationRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewInvitationRepository creates a new invitation repository instance
func NewInvitationRepository(db *gorm.DB, logger *zap.Logger) InvitationRepository {
	return &invitationRepository{db: db, logger: logger}
}

// ===============================
// BASIC CRUD OPERATIONS
// ===============================

func (r *invitationRepository) Create(ctx context.Context, invitation *models.OrganizationInvite) (*models.OrganizationInvite, error) {
	// Generate token if not provided
	if invitation.Token == "" {
		token, err := r.generateSecureToken()
		if err != nil {
			return nil, fmt.Errorf("failed to generate invitation token: %w", err)
		}
		invitation.Token = token
	}

	if err := r.db.WithContext(ctx).Create(invitation).Error; err != nil {
		return nil, fmt.Errorf("failed to create invitation: %w", err)
	}
	return invitation, nil
}

func (r *invitationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.OrganizationInvite, error) {
	var invitation models.OrganizationInvite
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&invitation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invitation not found")
		}
		return nil, fmt.Errorf("failed to get invitation: %w", err)
	}
	return &invitation, nil
}

func (r *invitationRepository) GetByToken(ctx context.Context, token string) (*models.OrganizationInvite, error) {
	var invitation models.OrganizationInvite
	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&invitation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invitation not found")
		}
		return nil, fmt.Errorf("failed to get invitation by token: %w", err)
	}
	return &invitation, nil
}

func (r *invitationRepository) Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.OrganizationInvite, error) {
	var invitation models.OrganizationInvite

	// First check if invitation exists
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&invitation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invitation not found")
		}
		return nil, fmt.Errorf("failed to get invitation: %w", err)
	}

	// Add updated_at timestamp
	updates["updated_at"] = time.Now()

	// Perform update
	if err := r.db.WithContext(ctx).Model(&invitation).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update invitation: %w", err)
	}

	// Return updated invitation
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&invitation).Error; err != nil {
		return nil, fmt.Errorf("failed to get updated invitation: %w", err)
	}

	return &invitation, nil
}

func (r *invitationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Unscoped().Delete(&models.OrganizationInvite{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete invitation: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("invitation not found")
	}
	return nil
}

// ===============================
// LISTING AND FILTERING
// ===============================

func (r *invitationRepository) List(ctx context.Context, filters dto.InvitationFiltersRequest) ([]models.OrganizationInvite, int64, error) {
	var invitations []models.OrganizationInvite
	var total int64

	query := r.db.WithContext(ctx).Model(&models.OrganizationInvite{})

	// Apply filters
	query = r.applyInvitationFilters(query, filters)

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count invitations: %w", err)
	}

	// Apply pagination and sorting
	query = r.applyPaginationAndSorting(query, filters.PaginationRequest, filters.FilterRequest)

	// Execute query
	if err := query.Find(&invitations).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list invitations: %w", err)
	}

	return invitations, total, nil
}

func (r *invitationRepository) ListByOrganization(ctx context.Context, orgID uuid.UUID, filters dto.InvitationFiltersRequest) ([]models.OrganizationInvite, int64, error) {
	var invitations []models.OrganizationInvite
	var total int64

	query := r.db.WithContext(ctx).Model(&models.OrganizationInvite{}).Where("organization_id = ?", orgID)

	// Apply filters
	query = r.applyInvitationFilters(query, filters)

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count invitations by organization: %w", err)
	}

	// Apply pagination and sorting
	query = r.applyPaginationAndSorting(query, filters.PaginationRequest, filters.FilterRequest)

	// Execute query
	if err := query.Find(&invitations).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list invitations by organization: %w", err)
	}

	return invitations, total, nil
}

func (r *invitationRepository) ListByEmail(ctx context.Context, email string) ([]models.OrganizationInvite, error) {
	var invitations []models.OrganizationInvite
	if err := r.db.WithContext(ctx).Where("invitee_email = ?", email).Find(&invitations).Error; err != nil {
		return nil, fmt.Errorf("failed to list invitations by email: %w", err)
	}
	return invitations, nil
}

// ===============================
// STATUS MANAGEMENT
// ===============================

func (r *invitationRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.InvitationStatus, updatedBy *uuid.UUID) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	if updatedBy != nil {
		updates["updated_by"] = *updatedBy
	}

	result := r.db.WithContext(ctx).Model(&models.OrganizationInvite{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update invitation status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("invitation not found")
	}
	return nil
}

func (r *invitationRepository) AcceptInvitation(ctx context.Context, id uuid.UUID, acceptedBy *uuid.UUID) error {
	updates := map[string]interface{}{
		"status":      models.InvStatusAccepted,
		"accepted_at": time.Now(),
		"updated_at":  time.Now(),
	}

	if acceptedBy != nil {
		updates["accepted_by"] = *acceptedBy
		updates["updated_by"] = *acceptedBy
	}

	result := r.db.WithContext(ctx).Model(&models.OrganizationInvite{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to accept invitation: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("invitation not found")
	}
	return nil
}

func (r *invitationRepository) RejectInvitation(ctx context.Context, id uuid.UUID, rejectedBy *uuid.UUID, reason *string) error {
	updates := map[string]interface{}{
		"status":      models.InvStatusRejected,
		"rejected_at": time.Now(),
		"updated_at":  time.Now(),
	}

	if rejectedBy != nil {
		updates["rejected_by"] = *rejectedBy
		updates["updated_by"] = *rejectedBy
	}

	if reason != nil {
		updates["rejection_reason"] = *reason
	}

	result := r.db.WithContext(ctx).Model(&models.OrganizationInvite{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to reject invitation: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("invitation not found")
	}
	return nil
}

func (r *invitationRepository) CancelInvitation(ctx context.Context, id uuid.UUID, cancelledBy uuid.UUID, reason string) error {
	updates := map[string]interface{}{
		"status":              models.InvStatusCancelled,
		"cancelled_at":        time.Now(),
		"cancellation_reason": reason,
		"updated_by":          cancelledBy,
		"updated_at":          time.Now(),
	}

	result := r.db.WithContext(ctx).Model(&models.OrganizationInvite{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to cancel invitation: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("invitation not found")
	}
	return nil
}

func (r *invitationRepository) ExpireInvitation(ctx context.Context, id uuid.UUID) error {
	updates := map[string]interface{}{
		"status":     models.InvStatusExpired,
		"updated_at": time.Now(),
	}

	result := r.db.WithContext(ctx).Model(&models.OrganizationInvite{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to expire invitation: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("invitation not found")
	}
	return nil
}

// ===============================
// TOKEN MANAGEMENT
// ===============================

func (r *invitationRepository) RegenerateToken(ctx context.Context, id uuid.UUID) (string, error) {
	newToken, err := r.generateSecureToken()
	if err != nil {
		return "", fmt.Errorf("failed to generate new token: %w", err)
	}

	updates := map[string]interface{}{
		"token":      newToken,
		"updated_at": time.Now(),
	}

	result := r.db.WithContext(ctx).Model(&models.OrganizationInvite{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return "", fmt.Errorf("failed to regenerate token: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return "", fmt.Errorf("invitation not found")
	}

	return newToken, nil
}

func (r *invitationRepository) IsTokenValid(ctx context.Context, token string) (bool, error) {
	var invitation models.OrganizationInvite
	err := r.db.WithContext(ctx).Where("token = ? AND status = ? AND expires_at > ?",
		token, models.InvStatusPending, time.Now()).First(&invitation).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, fmt.Errorf("failed to validate token: %w", err)
	}

	return true, nil
}

// ===============================
// EXPIRATION MANAGEMENT
// ===============================

func (r *invitationRepository) GetExpiredInvitations(ctx context.Context, limit int) ([]models.OrganizationInvite, error) {
	var invitations []models.OrganizationInvite
	query := r.db.WithContext(ctx).Where("status = ? AND expires_at < ?", models.InvStatusPending, time.Now())

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&invitations).Error; err != nil {
		return nil, fmt.Errorf("failed to get expired invitations: %w", err)
	}

	return invitations, nil
}

func (r *invitationRepository) MarkExpiredInvitations(ctx context.Context) (int64, error) {
	result := r.db.WithContext(ctx).Model(&models.OrganizationInvite{}).
		Where("status = ? AND expires_at < ?", models.InvStatusPending, time.Now()).
		Updates(map[string]interface{}{
			"status":     models.InvStatusExpired,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return 0, fmt.Errorf("failed to mark expired invitations: %w", result.Error)
	}

	return result.RowsAffected, nil
}

// ===============================
// RELATIONSHIPS
// ===============================

func (r *invitationRepository) GetWithOrganization(ctx context.Context, id uuid.UUID) (*models.OrganizationInvite, error) {
	var invitation models.OrganizationInvite
	if err := r.db.WithContext(ctx).Preload("Organization").Where("id = ?", id).First(&invitation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invitation not found")
		}
		return nil, fmt.Errorf("failed to get invitation with organization: %w", err)
	}
	return &invitation, nil
}

func (r *invitationRepository) GetWithRole(ctx context.Context, id uuid.UUID) (*models.OrganizationInvite, error) {
	var invitation models.OrganizationInvite
	if err := r.db.WithContext(ctx).Preload("Role").Where("id = ?", id).First(&invitation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invitation not found")
		}
		return nil, fmt.Errorf("failed to get invitation with role: %w", err)
	}
	return &invitation, nil
}

func (r *invitationRepository) GetWithBranch(ctx context.Context, id uuid.UUID) (*models.OrganizationInvite, error) {
	var invitation models.OrganizationInvite
	if err := r.db.WithContext(ctx).Preload("Branch").Where("id = ?", id).First(&invitation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invitation not found")
		}
		return nil, fmt.Errorf("failed to get invitation with branch: %w", err)
	}
	return &invitation, nil
}

func (r *invitationRepository) GetWithLogs(ctx context.Context, id uuid.UUID) (*models.OrganizationInvite, error) {
	var invitation models.OrganizationInvite
	if err := r.db.WithContext(ctx).Preload("Logs").Where("id = ?", id).First(&invitation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invitation not found")
		}
		return nil, fmt.Errorf("failed to get invitation with logs: %w", err)
	}
	return &invitation, nil
}

// ===============================
// BUSINESS LOGIC SUPPORT
// ===============================

func (r *invitationRepository) ExistsByEmailAndOrg(ctx context.Context, email string, orgID uuid.UUID, excludeStatuses []models.InvitationStatus) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.OrganizationInvite{}).
		Where("invitee_email = ? AND organization_id = ?", email, orgID)

	if len(excludeStatuses) > 0 {
		query = query.Where("status NOT IN ?", excludeStatuses)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check invitation existence: %w", err)
	}
	return count > 0, nil
}

func (r *invitationRepository) CountPendingByOrganization(ctx context.Context, orgID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.OrganizationInvite{}).
		Where("organization_id = ? AND status = ?", orgID, models.InvStatusPending).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count pending invitations: %w", err)
	}
	return count, nil
}

func (r *invitationRepository) GetInvitationStats(ctx context.Context, orgID uuid.UUID) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Count by status
	var statusCounts []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}

	if err := r.db.WithContext(ctx).Model(&models.OrganizationInvite{}).
		Select("status, COUNT(*) as count").
		Where("organization_id = ?", orgID).
		Group("status").
		Scan(&statusCounts).Error; err != nil {
		return nil, fmt.Errorf("failed to get invitation stats: %w", err)
	}

	for _, sc := range statusCounts {
		stats[sc.Status] = sc.Count
	}

	// Total invitations
	var total int64
	if err := r.db.WithContext(ctx).Model(&models.OrganizationInvite{}).
		Where("organization_id = ?", orgID).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count total invitations: %w", err)
	}
	stats["total"] = total

	return stats, nil
}

// ===============================
// HELPER METHODS
// ===============================

func (r *invitationRepository) generateSecureToken() (string, error) {
	// Generate 32 random bytes (256 bits)
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Encode to URL-safe base64
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (r *invitationRepository) applyInvitationFilters(query *gorm.DB, filters dto.InvitationFiltersRequest) *gorm.DB {
	// Search filter
	if filters.Search != "" {
		searchTerm := "%" + filters.Search + "%"
		query = query.Where("invitee_email ILIKE ?", searchTerm)
	}

	// Organization filter
	if filters.OrganizationID != uuid.Nil {
		query = query.Where("organization_id = ?", filters.OrganizationID)
	}

	// Status filter
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}

	// Role filter
	if filters.RoleID != nil {
		query = query.Where("role_id = ?", *filters.RoleID)
	}

	// Branch filter
	if filters.BranchID != nil {
		query = query.Where("branch_id = ?", *filters.BranchID)
	}

	// Email filter
	if filters.InviteeEmail != nil {
		query = query.Where("invitee_email = ?", *filters.InviteeEmail)
	}

	// Expiration filters
	if filters.ExpiresAfter != nil {
		query = query.Where("expires_at >= ?", *filters.ExpiresAfter)
	}

	if filters.ExpiresBefore != nil {
		query = query.Where("expires_at <= ?", *filters.ExpiresBefore)
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

func (r *invitationRepository) applyPaginationAndSorting(query *gorm.DB, pagination dto.PaginationRequest, filter dto.FilterRequest) *gorm.DB {
	// Apply default sorting (pending first, then by created date desc)
	query = query.Order("CASE WHEN status = 'pending' THEN 0 ELSE 1 END, created_at DESC")

	// Apply pagination
	if pagination.Page > 0 && pagination.PerPage > 0 {
		offset := (pagination.Page - 1) * pagination.PerPage
		query = query.Offset(offset).Limit(pagination.PerPage)
	}

	return query
}
