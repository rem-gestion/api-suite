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

// invitationLogRepository implements InvitationLogRepository interface
type invitationLogRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewInvitationLogRepository creates a new invitation log repository instance
func NewInvitationLogRepository(db *gorm.DB, logger *zap.Logger) InvitationLogRepository {
	return &invitationLogRepository{db: db, logger: logger}
}

// ===============================
// BASIC OPERATIONS
// ===============================

func (r *invitationLogRepository) Create(ctx context.Context, log *models.OrganizationInviteLog) (*models.OrganizationInviteLog, error) {
	if err := r.db.WithContext(ctx).Create(log).Error; err != nil {
		return nil, fmt.Errorf("failed to create invitation log: %w", err)
	}
	return log, nil
}

func (r *invitationLogRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.OrganizationInviteLog, error) {
	var log models.OrganizationInviteLog
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&log).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invitation log not found")
		}
		return nil, fmt.Errorf("failed to get invitation log: %w", err)
	}
	return &log, nil
}

// ===============================
// LISTING OPERATIONS
// ===============================

func (r *invitationLogRepository) ListByInvitation(ctx context.Context, inviteID uuid.UUID) ([]models.OrganizationInviteLog, error) {
	var logs []models.OrganizationInviteLog
	if err := r.db.WithContext(ctx).Where("invite_id = ?", inviteID).Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, fmt.Errorf("failed to list invitation logs by invitation: %w", err)
	}
	return logs, nil
}

func (r *invitationLogRepository) ListByOrganization(ctx context.Context, orgID uuid.UUID, limit int) ([]models.OrganizationInviteLog, error) {
	var logs []models.OrganizationInviteLog
	query := r.db.WithContext(ctx).
		Joins("JOIN organization_invite ON organization_invite_log.invite_id = organization_invite.id").
		Where("organization_invite.organization_id = ?", orgID).
		Order("organization_invite_log.created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&logs).Error; err != nil {
		return nil, fmt.Errorf("failed to list invitation logs by organization: %w", err)
	}
	return logs, nil
}

func (r *invitationLogRepository) ListByAction(ctx context.Context, action models.InvitationLogAction, limit int) ([]models.OrganizationInviteLog, error) {
	var logs []models.OrganizationInviteLog
	query := r.db.WithContext(ctx).Where("action = ?", action).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&logs).Error; err != nil {
		return nil, fmt.Errorf("failed to list invitation logs by action: %w", err)
	}
	return logs, nil
}

// ===============================
// CONVENIENCE METHODS FOR COMMON LOG ENTRIES
// ===============================

func (r *invitationLogRepository) LogInvitationCreated(ctx context.Context, inviteID uuid.UUID, actorPersonID *uuid.UUID, clientIP, userAgent *string) error {
	userType := "user"
	return r.createLogEntry(ctx, inviteID, models.InvLogCreated, actorPersonID, &userType, clientIP, userAgent, nil)
}

func (r *invitationLogRepository) LogInvitationAccepted(ctx context.Context, inviteID uuid.UUID, actorPersonID *uuid.UUID, clientIP, userAgent *string) error {
	userType := "user"
	return r.createLogEntry(ctx, inviteID, models.InvLogAccepted, actorPersonID, &userType, clientIP, userAgent, nil)
}

func (r *invitationLogRepository) LogInvitationRejected(ctx context.Context, inviteID uuid.UUID, actorPersonID *uuid.UUID, reason *string, clientIP, userAgent *string) error {
	userType := "user"
	notes := map[string]interface{}{}
	if reason != nil {
		notes["reason"] = *reason
	}
	return r.createLogEntry(ctx, inviteID, models.InvLogRejected, actorPersonID, &userType, clientIP, userAgent, notes)
}

func (r *invitationLogRepository) LogInvitationCancelled(ctx context.Context, inviteID uuid.UUID, actorPersonID *uuid.UUID, reason *string, clientIP, userAgent *string) error {
	userType := "user"
	notes := map[string]interface{}{}
	if reason != nil {
		notes["reason"] = *reason
	}
	return r.createLogEntry(ctx, inviteID, models.InvLogCancelled, actorPersonID, &userType, clientIP, userAgent, notes)
}

func (r *invitationLogRepository) LogInvitationExpired(ctx context.Context, inviteID uuid.UUID) error {
	systemType := "system"
	notes := map[string]interface{}{
		"message": "Automatically expired by system",
	}
	return r.createLogEntry(ctx, inviteID, models.InvLogExpired, nil, &systemType, nil, nil, notes)
}

func (r *invitationLogRepository) LogInvitationResent(ctx context.Context, inviteID uuid.UUID, actorPersonID *uuid.UUID, clientIP, userAgent *string) error {
	userType := "user"
	notes := map[string]interface{}{
		"message": "Invitation resent",
	}
	return r.createLogEntry(ctx, inviteID, models.InvLogUpdated, actorPersonID, &userType, clientIP, userAgent, notes)
}

// ===============================
// CLEANUP OPERATIONS
// ===============================

func (r *invitationLogRepository) DeleteOldLogs(ctx context.Context, olderThan int) (int64, error) {
	cutoffDate := time.Now().AddDate(0, 0, -olderThan)

	result := r.db.WithContext(ctx).Where("created_at < ?", cutoffDate).Delete(&models.OrganizationInviteLog{})
	if result.Error != nil {
		return 0, fmt.Errorf("failed to delete old invitation logs: %w", result.Error)
	}

	return result.RowsAffected, nil
}

// ===============================
// HELPER METHODS
// ===============================

func (r *invitationLogRepository) createLogEntry(
	ctx context.Context,
	inviteID uuid.UUID,
	action models.InvitationLogAction,
	actorPersonID *uuid.UUID,
	actorType *string,
	clientIP, userAgent *string,
	metadata map[string]interface{},
) error {
	// Convert metadata to notes if needed
	var notes *string
	if metadata != nil {
		if msg, exists := metadata["message"]; exists {
			if msgStr, ok := msg.(string); ok {
				notes = &msgStr
			}
		} else if reason, exists := metadata["reason"]; exists {
			if reasonStr, ok := reason.(string); ok {
				notes = &reasonStr
			}
		}
	}

	// Convert actorType pointer to string
	actorTypeStr := "user"
	if actorType != nil {
		actorTypeStr = *actorType
	}

	log := &models.OrganizationInviteLog{
		ID:            uuid.New(),
		InviteID:      inviteID,
		Action:        action,
		ActorPersonID: actorPersonID,
		ActorType:     actorTypeStr,
		ClientIP:      clientIP,
		UserAgent:     userAgent,
		Notes:         notes,
		ErrorDetails:  metadata,
		CreatedAt:     time.Now(),
	}

	if err := r.db.WithContext(ctx).Create(log).Error; err != nil {
		return fmt.Errorf("failed to create log entry for action %s: %w", action, err)
	}

	return nil
}
