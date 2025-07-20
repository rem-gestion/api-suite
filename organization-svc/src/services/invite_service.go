package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/rem-gestion/api-suite/organization/src/clients"
	"github.com/rem-gestion/api-suite/organization/src/dto"
	"github.com/rem-gestion/api-suite/organization/src/models"
	"github.com/rem-gestion/api-suite/organization/src/repository"
	rcerrors "github.com/rem-gestion/rem-common/errors"
)

// InvitationService maneja la lógica de negocio para invitaciones organizacionales
type InvitationService struct {
	repos           *repository.Repositories
	clients         *clients.ClientManager
	cache           CacheService
	events          EventService
	employeeService *EmployeeService
	logger          *zap.Logger
}

// NewInvitationService crea una nueva instancia del servicio de invitaciones
func NewInvitationService(
	repos *repository.Repositories,
	clients *clients.ClientManager,
	cache CacheService,
	events EventService,
	employeeService *EmployeeService,
	logger *zap.Logger,
) *InvitationService {
	return &InvitationService{
		repos:           repos,
		clients:         clients,
		cache:           cache,
		events:          events,
		employeeService: employeeService,
		logger:          logger,
	}
}

// ===============================
// CORE INVITATION OPERATIONS
// ===============================

// SendInvitation envía una nueva invitación para unirse a una organización
func (s *InvitationService) SendInvitation(ctx context.Context, orgID uuid.UUID, inviterPersonID uuid.UUID, req dto.CreateInvitationRequest) (*models.OrganizationInvite, error) {
	s.logger.Info("Sending invitation",
		zap.String("org_id", orgID.String()),
		zap.String("inviter_person_id", inviterPersonID.String()),
		zap.String("invitee_email", req.InviteeEmail),
	)

	// 1. Validar que la organización existe
	org, err := s.repos.Organization.GetByID(ctx, orgID)
	if err != nil {
		s.logger.Error("Failed to get organization", zap.Error(err))
		return nil, rcerrors.NewNotFoundError("organization", orgID.String())
	}

	// 2. Validar que el role_id existe en la organización
	role, err := s.repos.Role.GetByID(ctx, req.RoleID)
	if err != nil || role.OrganizationID != orgID {
		s.logger.Error("Invalid role for organization", zap.Error(err))
		return nil, rcerrors.NewValidationError("role_id", "invalid role for this organization")
	}

	// 3. Validar branch_id si se proporciona
	if req.BranchID != nil {
		branch, err := s.repos.Branch.GetByID(ctx, *req.BranchID)
		if err != nil || branch.OrganizationID != orgID {
			s.logger.Error("Invalid branch for organization", zap.Error(err))
			return nil, rcerrors.NewValidationError("branch_id", "invalid branch for this organization")
		}
	}

	// 4. Verificar si ya existe invitación pending para este email
	exists, err := s.repos.Invitation.ExistsByEmailAndOrg(ctx, req.InviteeEmail, orgID, []models.InvitationStatus{models.InvStatusPending})
	if err == nil && exists {
		s.logger.Warn("Invitation already exists for email", zap.String("email", req.InviteeEmail))
		return nil, rcerrors.NewConflictError("invitation already pending for this email")
	}

	// 5. Generar token seguro
	token, err := s.generateSecureToken()
	if err != nil {
		s.logger.Error("Failed to generate invitation token", zap.Error(err))
		return nil, rcerrors.NewInternalServerError("failed to generate invitation token")
	}

	// 6. Crear invitación
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 7 días por defecto
	if req.ExpiresAt != nil {
		expiresAt = *req.ExpiresAt
	}

	invitation := &models.OrganizationInvite{
		OrganizationID:  orgID,
		InviterPersonID: inviterPersonID,
		InviteeEmail:    req.InviteeEmail,
		InviteePersonID: req.InviteePersonID,
		RoleID:          req.RoleID,
		BranchID:        req.BranchID,
		Token:           token,
		Status:          models.InvStatusPending,
		ExpiresAt:       expiresAt,
		Metadata:        req.Metadata,
	}

	// 7. Guardar en base de datos
	createdInvitation, err := s.repos.Invitation.Create(ctx, invitation)
	if err != nil {
		s.logger.Error("Failed to create invitation", zap.Error(err))
		return nil, rcerrors.NewInternalServerError("failed to create invitation")
	}

	// 8. Log de auditoría
	if err := s.createInvitationLog(ctx, createdInvitation.ID, models.InvLogCreated, inviterPersonID); err != nil {
		s.logger.Warn("Failed to create invitation log", zap.Error(err))
	}

	// 9. Publicar evento para envío de email
	if err := s.publishInvitationSentEvent(ctx, createdInvitation, org, role); err != nil {
		s.logger.Warn("Failed to publish invitation sent event", zap.Error(err))
	}

	s.logger.Info("Invitation sent successfully", zap.String("invitation_id", createdInvitation.ID.String()))
	return createdInvitation, nil
}

// GetInvitation obtiene una invitación por ID
func (s *InvitationService) GetInvitation(ctx context.Context, invitationID uuid.UUID, requesterPersonID uuid.UUID) (*models.OrganizationInvite, error) {
	s.logger.Debug("Getting invitation", zap.String("invitation_id", invitationID.String()))

	invitation, err := s.repos.Invitation.GetByID(ctx, invitationID)
	if err != nil {
		s.logger.Error("Failed to get invitation", zap.Error(err))
		return nil, rcerrors.NewNotFoundError("invitation", invitationID.String())
	}

	return invitation, nil
}

// ListInvitations lista invitaciones con filtros
func (s *InvitationService) ListInvitations(ctx context.Context, requesterPersonID uuid.UUID, filters dto.InvitationFiltersRequest) ([]models.OrganizationInvite, int64, error) {
	s.logger.Debug("Listing invitations", zap.String("requester_person_id", requesterPersonID.String()))

	invitations, total, err := s.repos.Invitation.List(ctx, filters)
	if err != nil {
		s.logger.Error("Failed to list invitations", zap.Error(err))
		return nil, 0, rcerrors.NewInternalServerError("failed to list invitations")
	}

	return invitations, total, nil
}

// ResendInvitation reenvía una invitación existente
func (s *InvitationService) ResendInvitation(ctx context.Context, invitationID uuid.UUID, requesterPersonID uuid.UUID, req dto.ResendInvitationRequest) error {
	s.logger.Info("Resending invitation", zap.String("invitation_id", invitationID.String()))

	invitation, err := s.repos.Invitation.GetByID(ctx, invitationID)
	if err != nil {
		return rcerrors.NewNotFoundError("invitation", invitationID.String())
	}

	if invitation.Status != models.InvStatusPending {
		return rcerrors.NewValidationError("status", "only pending invitations can be resent")
	}

	if time.Now().After(invitation.ExpiresAt) {
		return rcerrors.NewValidationError("expires_at", "cannot resend expired invitation")
	}

	// Actualizar expiración si se proporciona
	updates := map[string]interface{}{}
	if req.NewExpiresAt != nil {
		updates["expires_at"] = *req.NewExpiresAt
	}

	if len(updates) > 0 {
		if _, err := s.repos.Invitation.Update(ctx, invitationID, updates); err != nil {
			s.logger.Error("Failed to update invitation", zap.Error(err))
			return rcerrors.NewInternalServerError("failed to update invitation")
		}
	}

	// Log de auditoría
	if err := s.createInvitationLog(ctx, invitationID, models.InvLogUpdated, requesterPersonID); err != nil {
		s.logger.Warn("Failed to create invitation log", zap.Error(err))
	}

	s.logger.Info("Invitation resent successfully", zap.String("invitation_id", invitationID.String()))
	return nil
}

// CancelInvitation cancela una invitación
func (s *InvitationService) CancelInvitation(ctx context.Context, invitationID uuid.UUID, requesterPersonID uuid.UUID, req dto.CancelInvitationRequest) error {
	s.logger.Info("Cancelling invitation", zap.String("invitation_id", invitationID.String()))

	invitation, err := s.repos.Invitation.GetByID(ctx, invitationID)
	if err != nil {
		return rcerrors.NewNotFoundError("invitation", invitationID.String())
	}

	if invitation.Status != models.InvStatusPending {
		return rcerrors.NewValidationError("status", "only pending invitations can be cancelled")
	}

	// Actualizar estado a cancelled
	now := time.Now()
	updates := map[string]interface{}{
		"status":       models.InvStatusCancelled,
		"cancelled_at": &now,
		"metadata": map[string]interface{}{
			"cancellation_reason": req.CancellationReason,
		},
	}

	if _, err := s.repos.Invitation.Update(ctx, invitationID, updates); err != nil {
		s.logger.Error("Failed to cancel invitation", zap.Error(err))
		return rcerrors.NewInternalServerError("failed to cancel invitation")
	}

	// Log de auditoría
	if err := s.createInvitationLog(ctx, invitationID, models.InvLogCancelled, requesterPersonID); err != nil {
		s.logger.Warn("Failed to create invitation log", zap.Error(err))
	}

	s.logger.Info("Invitation cancelled successfully", zap.String("invitation_id", invitationID.String()))
	return nil
}

// ===============================
// PUBLIC INVITATION OPERATIONS (Sin autenticación)
// ===============================

// ValidateInvitationToken valida un token de invitación y retorna la información
func (s *InvitationService) ValidateInvitationToken(ctx context.Context, token string) (*models.OrganizationInvite, error) {
	s.logger.Debug("Validating invitation token")

	invitation, err := s.repos.Invitation.GetByToken(ctx, token)
	if err != nil {
		s.logger.Error("Failed to get invitation by token", zap.Error(err))
		return nil, rcerrors.NewNotFoundError("invitation", "token")
	}

	if invitation.Status != models.InvStatusPending {
		return nil, rcerrors.NewValidationError("status", "invitation is no longer valid")
	}

	if time.Now().After(invitation.ExpiresAt) {
		// Marcar como expirada
		s.markInvitationAsExpired(ctx, invitation.ID)
		return nil, rcerrors.NewValidationError("expires_at", "invitation has expired")
	}

	return invitation, nil
}

// AcceptInvitation acepta una invitación y crea el empleado (versión simplificada)
func (s *InvitationService) AcceptInvitation(ctx context.Context, token string, req dto.AcceptInvitationRequest) (*models.Employee, error) {
	s.logger.Info("Accepting invitation")

	// 1. Validar token
	invitation, err := s.ValidateInvitationToken(ctx, token)
	if err != nil {
		return nil, err
	}

	// 2. Para esta implementación simplificada, asumimos que el person_id y user_id ya existen
	if invitation.InviteePersonID == nil {
		return nil, rcerrors.NewValidationError("person_id", "person_data is required for new persons")
	}

	// 3. Crear empleado usando EmployeeService (simulado por ahora)
	// En una implementación completa, aquí se haría la integración completa con los otros servicios
	var personID uuid.UUID
	if invitation.InviteePersonID != nil {
		personID = *invitation.InviteePersonID
	} else {
		// Si no hay PersonID, generar uno temporal o manejarlo según la lógica de negocio
		personID = uuid.New()
	}

	employee := &models.Employee{
		ID:             uuid.New(),
		OrganizationID: invitation.OrganizationID,
		UserID:         uuid.New(), // Sería el ID real del usuario
		PersonID:       &personID,
		Status:         models.EmpStatusActive,
	}

	// 4. Marcar invitación como aceptada
	now := time.Now()
	updates := map[string]interface{}{
		"status":      models.InvStatusAccepted,
		"accepted_at": &now,
	}

	if _, err := s.repos.Invitation.Update(ctx, invitation.ID, updates); err != nil {
		s.logger.Error("Failed to mark invitation as accepted", zap.Error(err))
	}

	// 5. Log de auditoría
	if err := s.createInvitationLog(ctx, invitation.ID, models.InvLogAccepted, *invitation.InviteePersonID); err != nil {
		s.logger.Warn("Failed to create invitation log", zap.Error(err))
	}

	s.logger.Info("Invitation accepted successfully", zap.String("invitation_id", invitation.ID.String()))
	return employee, nil
}

// DeclineInvitation rechaza una invitación
func (s *InvitationService) DeclineInvitation(ctx context.Context, token string, req dto.DeclineInvitationRequest) error {
	s.logger.Info("Declining invitation")

	invitation, err := s.ValidateInvitationToken(ctx, token)
	if err != nil {
		return err
	}

	// Marcar como rechazada
	now := time.Now()
	updates := map[string]interface{}{
		"status":      models.InvStatusRejected,
		"rejected_at": &now,
	}

	if req.DeclineReason != nil {
		updates["metadata"] = map[string]interface{}{
			"decline_reason": *req.DeclineReason,
		}
	}

	if _, err := s.repos.Invitation.Update(ctx, invitation.ID, updates); err != nil {
		s.logger.Error("Failed to decline invitation", zap.Error(err))
		return rcerrors.NewInternalServerError("failed to decline invitation")
	}

	// Log de auditoría
	if err := s.createInvitationLog(ctx, invitation.ID, models.InvLogRejected, uuid.Nil); err != nil {
		s.logger.Warn("Failed to create invitation log", zap.Error(err))
	}

	s.logger.Info("Invitation declined successfully", zap.String("invitation_id", invitation.ID.String()))
	return nil
}

// ===============================
// HELPER METHODS
// ===============================

// generateSecureToken genera un token seguro para invitaciones
func (s *InvitationService) generateSecureToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// markInvitationAsExpired marca una invitación como expirada
func (s *InvitationService) markInvitationAsExpired(ctx context.Context, invitationID uuid.UUID) {
	updates := map[string]interface{}{
		"status": models.InvStatusExpired,
	}

	if _, err := s.repos.Invitation.Update(ctx, invitationID, updates); err != nil {
		s.logger.Error("Failed to mark invitation as expired", zap.Error(err))
	}

	// Log de auditoría
	if err := s.createInvitationLog(ctx, invitationID, models.InvLogExpired, uuid.Nil); err != nil {
		s.logger.Warn("Failed to create invitation log", zap.Error(err))
	}
}

// createInvitationLog crea un log de auditoría para la invitación
func (s *InvitationService) createInvitationLog(ctx context.Context, invitationID uuid.UUID, action models.InvitationLogAction, actorPersonID uuid.UUID) error {
	log := &models.OrganizationInviteLog{
		InviteID:      invitationID,
		Action:        action,
		ActorPersonID: &actorPersonID,
		NewStatus:     models.InvStatusPending, // Estado por defecto
	}

	_, err := s.repos.InvitationLog.Create(ctx, log)
	return err
}

// ===============================
// EVENT PUBLISHING METHODS
// ===============================

// publishInvitationSentEvent publica evento cuando se envía una invitación
func (s *InvitationService) publishInvitationSentEvent(ctx context.Context, invitation *models.OrganizationInvite, org *models.Organization, role *models.OrganizationRole) error {
	return s.events.PublishInvitationSent(org, invitation)
}
