package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/rem-gestion/api-suite/organization/src/clients"
	"github.com/rem-gestion/api-suite/organization/src/dto"
	"github.com/rem-gestion/api-suite/organization/src/models"
	"github.com/rem-gestion/api-suite/organization/src/repository"
	rcerrors "github.com/rem-gestion/rem-common/errors"
)

// OwnerService maneja la lógica de negocio para propietarios organizacionales
type OwnerService struct {
	repos   *repository.Repositories
	clients *clients.ClientManager
	cache   CacheService
	events  EventService
	logger  *zap.Logger
}

// NewOwnerService crea una nueva instancia del servicio de propietarios
func NewOwnerService(
	repos *repository.Repositories,
	clients *clients.ClientManager,
	cache CacheService,
	events EventService,
	logger *zap.Logger,
) *OwnerService {
	return &OwnerService{
		repos:   repos,
		clients: clients,
		cache:   cache,
		events:  events,
		logger:  logger,
	}
}

// ===============================
// CORE OWNER OPERATIONS
// ===============================

// AddOwner añade un nuevo propietario a una organización
func (s *OwnerService) AddOwner(ctx context.Context, orgID uuid.UUID, adderPersonID uuid.UUID, req dto.AddOwnerRequest) (*models.OrganizationOwner, error) {
	s.logger.Info("Adding owner",
		zap.String("org_id", orgID.String()),
		zap.String("adder_person_id", adderPersonID.String()),
		zap.String("person_id", req.PersonID.String()),
	)

	// 1. Validar que la organización existe y está activa
	org, err := s.repos.Organization.GetByID(ctx, orgID)
	if err != nil {
		s.logger.Error("Failed to get organization", zap.Error(err))
		return nil, rcerrors.NewNotFoundError("organization", orgID.String())
	}

	if org.Status != models.OrgStatusActive {
		return nil, rcerrors.NewValidationError("organization", "organization must be active to add owners")
	}

	// 2. Verificar permisos para agregar propietarios (solo propietarios existentes pueden agregar)
	hasPermission, err := s.validateOwnershipPermissions(ctx, orgID, adderPersonID)
	if err != nil {
		s.logger.Error("Failed to validate ownership permissions", zap.Error(err))
		return nil, rcerrors.NewInternalServerError("failed to validate permissions")
	}
	if !hasPermission {
		return nil, rcerrors.NewForbiddenError("only existing owners can add new owners")
	}

	// 3. Validar que la persona existe via gRPC
	if err := s.validatePersonViaGRPC(ctx, req.PersonID); err != nil {
		s.logger.Error("Person validation failed", zap.Error(err))
		return nil, rcerrors.NewValidationError("person_id", "invalid person")
	}

	// 4. Verificar que la persona no es ya propietaria
	exists, err := s.repos.Owner.ExistsByPersonAndOrg(ctx, req.PersonID, orgID)
	if err != nil {
		s.logger.Error("Failed to check existing ownership", zap.Error(err))
		return nil, rcerrors.NewInternalServerError("failed to validate existing ownership")
	}
	if exists {
		return nil, rcerrors.NewConflictError("person is already an owner of this organization")
	}

	// 5. Validar porcentaje de propiedad
	if req.OwnershipPct != nil && (*req.OwnershipPct <= 0 || *req.OwnershipPct > 100) {
		return nil, rcerrors.NewValidationError("ownership_pct", "ownership percentage must be between 0.01 and 100")
	}

	// 6. Verificar que el porcentaje total no excede el 100%
	currentTotal, err := s.repos.Owner.GetTotalOwnership(ctx, orgID)
	if err != nil {
		s.logger.Error("Failed to get total ownership", zap.Error(err))
		return nil, rcerrors.NewInternalServerError("failed to validate ownership percentages")
	}

	var ownershipPct float64
	if req.OwnershipPct != nil {
		ownershipPct = *req.OwnershipPct
		if currentTotal+ownershipPct > 100 {
			return nil, rcerrors.NewValidationError("ownership_pct",
				fmt.Sprintf("adding this percentage would exceed 100%%. Current total: %.2f%%", currentTotal))
		}
	}

	// 7. Crear el propietario
	owner := &models.OrganizationOwner{
		OrganizationID: orgID,
		PersonID:       req.PersonID,
		UserID:         req.UserID,
		IsFounder:      req.IsFounder,
		OwnershipPct:   req.OwnershipPct,
		CreatedBy:      adderPersonID,
	}

	// 8. Guardar en base de datos
	createdOwner, err := s.repos.Owner.Create(ctx, owner)
	if err != nil {
		s.logger.Error("Failed to create owner", zap.Error(err))
		return nil, rcerrors.NewInternalServerError("failed to create owner")
	}

	// 9. Actualizar el estatus de majority owner si es necesario
	if err := s.updateMajorityOwnerStatus(ctx, orgID); err != nil {
		s.logger.Warn("Failed to update majority owner status", zap.Error(err))
	}

	// 10. Invalidar cache
	if err := s.cache.InvalidateOrganizationCache(orgID.String()); err != nil {
		s.logger.Warn("Failed to invalidate organization cache", zap.Error(err))
	}

	// 11. Publicar evento
	if err := s.publishOwnerAddedEvent(ctx, createdOwner, org); err != nil {
		s.logger.Warn("Failed to publish owner added event", zap.Error(err))
	}

	s.logger.Info("Owner added successfully",
		zap.String("owner_id", createdOwner.ID.String()),
		zap.String("org_id", orgID.String()),
	)

	return createdOwner, nil
}

// GetOwner obtiene un propietario por ID
func (s *OwnerService) GetOwner(ctx context.Context, ownerID uuid.UUID, requesterPersonID uuid.UUID) (*models.OrganizationOwner, error) {
	s.logger.Debug("Getting owner",
		zap.String("owner_id", ownerID.String()),
		zap.String("requester_person_id", requesterPersonID.String()),
	)

	// 1. Obtener el propietario
	owner, err := s.repos.Owner.GetByID(ctx, ownerID)
	if err != nil {
		s.logger.Error("Failed to get owner", zap.Error(err))
		return nil, rcerrors.NewNotFoundError("owner", ownerID.String())
	}

	// 2. Verificar permisos de acceso
	hasAccess, err := s.validateAccessPermissions(ctx, owner.OrganizationID, requesterPersonID)
	if err != nil {
		s.logger.Error("Failed to validate access permissions", zap.Error(err))
		return nil, rcerrors.NewInternalServerError("failed to validate access")
	}
	if !hasAccess {
		return nil, rcerrors.NewForbiddenError("insufficient permissions to access owner information")
	}

	return owner, nil
}

// ListOwners lista propietarios de una organización
func (s *OwnerService) ListOwners(ctx context.Context, requesterPersonID uuid.UUID, filters dto.OwnerFiltersRequest) ([]models.OrganizationOwner, int64, error) {
	s.logger.Debug("Listing owners",
		zap.String("requester_person_id", requesterPersonID.String()),
		zap.String("org_id", filters.OrganizationID.String()),
	)

	// 1. Verificar permisos de acceso
	hasAccess, err := s.validateAccessPermissions(ctx, filters.OrganizationID, requesterPersonID)
	if err != nil {
		s.logger.Error("Failed to validate access permissions", zap.Error(err))
		return nil, 0, rcerrors.NewInternalServerError("failed to validate access")
	}
	if !hasAccess {
		return nil, 0, rcerrors.NewForbiddenError("insufficient permissions to access organization owners")
	}

	// 2. Listar propietarios
	owners, total, err := s.repos.Owner.List(ctx, filters)
	if err != nil {
		s.logger.Error("Failed to list owners", zap.Error(err))
		return nil, 0, rcerrors.NewInternalServerError("failed to list owners")
	}

	return owners, total, nil
}

// RemoveOwner elimina un propietario de la organización
func (s *OwnerService) RemoveOwner(ctx context.Context, ownerID uuid.UUID, removerPersonID uuid.UUID, req dto.RemoveOwnerRequest) error {
	s.logger.Info("Removing owner",
		zap.String("owner_id", ownerID.String()),
		zap.String("remover_person_id", removerPersonID.String()),
	)

	// 1. Obtener el propietario
	owner, err := s.repos.Owner.GetByID(ctx, ownerID)
	if err != nil {
		s.logger.Error("Failed to get owner", zap.Error(err))
		return rcerrors.NewNotFoundError("owner", ownerID.String())
	}

	// 2. Verificar permisos para remover propietarios
	hasPermission, err := s.validateOwnershipPermissions(ctx, owner.OrganizationID, removerPersonID)
	if err != nil {
		s.logger.Error("Failed to validate ownership permissions", zap.Error(err))
		return rcerrors.NewInternalServerError("failed to validate permissions")
	}
	if !hasPermission {
		return rcerrors.NewForbiddenError("only existing owners can remove owners")
	}

	// 3. Verificar que no se está removiendo a sí mismo si es el único propietario
	if owner.PersonID == removerPersonID {
		totalOwners, err := s.repos.Owner.CountByOrganization(ctx, owner.OrganizationID)
		if err != nil {
			s.logger.Error("Failed to count owners", zap.Error(err))
			return rcerrors.NewInternalServerError("failed to validate owner count")
		}
		if totalOwners == 1 {
			return rcerrors.NewValidationError("owner", "cannot remove the last owner of the organization")
		}
	}

	// 4. Verificar reglas de negocio para la remoción
	canRemove, err := s.repos.Owner.CanRemoveOwner(ctx, ownerID)
	if err != nil {
		s.logger.Error("Failed to check if owner can be removed", zap.Error(err))
		return rcerrors.NewInternalServerError("failed to validate owner removal")
	}
	if !canRemove {
		return rcerrors.NewValidationError("owner", "owner cannot be removed due to business rules")
	}

	// 5. Eliminar el propietario
	if err := s.repos.Owner.Delete(ctx, ownerID); err != nil {
		s.logger.Error("Failed to remove owner", zap.Error(err))
		return rcerrors.NewInternalServerError("failed to remove owner")
	}

	// 6. Actualizar el estatus de majority owner si es necesario
	if err := s.updateMajorityOwnerStatus(ctx, owner.OrganizationID); err != nil {
		s.logger.Warn("Failed to update majority owner status", zap.Error(err))
	}

	// 7. Invalidar cache
	if err := s.cache.InvalidateOrganizationCache(owner.OrganizationID.String()); err != nil {
		s.logger.Warn("Failed to invalidate organization cache", zap.Error(err))
	}

	// 8. Publicar evento
	if err := s.publishOwnerRemovedEvent(ctx, owner); err != nil {
		s.logger.Warn("Failed to publish owner removed event", zap.Error(err))
	}

	s.logger.Info("Owner removed successfully", zap.String("owner_id", ownerID.String()))
	return nil
}

// UpdateOwnership actualiza la información de propiedad de un propietario
func (s *OwnerService) UpdateOwnership(ctx context.Context, ownerID uuid.UUID, updaterPersonID uuid.UUID, req dto.UpdateOwnerRequest) (*models.OrganizationOwner, error) {
	s.logger.Info("Updating ownership",
		zap.String("owner_id", ownerID.String()),
		zap.String("updater_person_id", updaterPersonID.String()),
	)

	// 1. Obtener el propietario actual
	owner, err := s.repos.Owner.GetByID(ctx, ownerID)
	if err != nil {
		s.logger.Error("Failed to get owner", zap.Error(err))
		return nil, rcerrors.NewNotFoundError("owner", ownerID.String())
	}

	// 2. Verificar permisos
	hasPermission, err := s.validateOwnershipPermissions(ctx, owner.OrganizationID, updaterPersonID)
	if err != nil {
		s.logger.Error("Failed to validate ownership permissions", zap.Error(err))
		return nil, rcerrors.NewInternalServerError("failed to validate permissions")
	}
	if !hasPermission {
		return nil, rcerrors.NewForbiddenError("only existing owners can update ownership information")
	}

	// 3. Validar nuevo porcentaje si se proporciona
	if req.OwnershipPct != nil {
		if *req.OwnershipPct <= 0 || *req.OwnershipPct > 100 {
			return nil, rcerrors.NewValidationError("ownership_pct", "ownership percentage must be between 0.01 and 100")
		}

		// 4. Verificar que el cambio no excede el 100%
		currentTotal, err := s.repos.Owner.GetTotalOwnership(ctx, owner.OrganizationID)
		if err != nil {
			s.logger.Error("Failed to get total ownership", zap.Error(err))
			return nil, rcerrors.NewInternalServerError("failed to validate ownership percentages")
		}

		var currentOwnership float64
		if owner.OwnershipPct != nil {
			currentOwnership = *owner.OwnershipPct
		}

		newTotal := currentTotal - currentOwnership + *req.OwnershipPct
		if newTotal > 100 {
			return nil, rcerrors.NewValidationError("ownership_pct",
				fmt.Sprintf("new percentage would result in total of %.2f%%, which exceeds 100%%", newTotal))
		}
	}

	// 5. Preparar las actualizaciones
	updates := make(map[string]interface{})
	updates["updated_by"] = updaterPersonID

	if req.IsFounder != nil {
		updates["is_founder"] = *req.IsFounder
	}
	if req.OwnershipPct != nil {
		updates["ownership_pct"] = *req.OwnershipPct
	}

	// 6. Actualizar el propietario
	updatedOwner, err := s.repos.Owner.Update(ctx, ownerID, updates)
	if err != nil {
		s.logger.Error("Failed to update owner", zap.Error(err))
		return nil, rcerrors.NewInternalServerError("failed to update owner")
	}

	// 7. Invalidar cache
	if err := s.cache.InvalidateOrganizationCache(owner.OrganizationID.String()); err != nil {
		s.logger.Warn("Failed to invalidate organization cache", zap.Error(err))
	}

	s.logger.Info("Ownership updated successfully", zap.String("owner_id", ownerID.String()))
	return updatedOwner, nil
}

// ===============================
// VALIDATION HELPER METHODS
// ===============================

// validatePersonViaGRPC valida que una persona existe usando el servicio de personas
func (s *OwnerService) validatePersonViaGRPC(ctx context.Context, personID uuid.UUID) error {
	s.logger.Debug("Validating person via gRPC", zap.String("person_id", personID.String()))

	// Para esta implementación simplificada, simulamos la validación gRPC
	// En una implementación completa, se usaría el cliente real de Person
	// exists, err := s.clients.Person.ValidatePerson(ctx, personID.String())
	// if err != nil {
	//     s.logger.Error("gRPC person validation failed", zap.Error(err))
	//     return fmt.Errorf("failed to validate person: %w", err)
	// }
	// if !exists {
	//     return fmt.Errorf("person not found")
	// }

	// Simulación: asumimos que la persona es válida por ahora
	s.logger.Info("Person validation skipped (mock implementation)", zap.String("person_id", personID.String()))
	return nil
}

// validateOwnershipPermissions verifica que un usuario tiene permisos de propietario
func (s *OwnerService) validateOwnershipPermissions(ctx context.Context, orgID uuid.UUID, personID uuid.UUID) (bool, error) {
	s.logger.Debug("Validating ownership permissions",
		zap.String("org_id", orgID.String()),
		zap.String("person_id", personID.String()),
	)

	// Verificar si es propietario existente de la organización
	isOwner, err := s.repos.Owner.ExistsByPersonAndOrg(ctx, personID, orgID)
	if err != nil {
		return false, fmt.Errorf("failed to check ownership: %w", err)
	}

	return isOwner, nil
}

// validateAccessPermissions verifica que un usuario tiene acceso a la información de propietarios
func (s *OwnerService) validateAccessPermissions(ctx context.Context, orgID uuid.UUID, personID uuid.UUID) (bool, error) {
	s.logger.Debug("Validating access permissions",
		zap.String("org_id", orgID.String()),
		zap.String("person_id", personID.String()),
	)

	// Para acceso a información de propietarios, verificamos si es propietario o empleado
	// Verificar si es propietario
	isOwner, err := s.repos.Owner.ExistsByPersonAndOrg(ctx, personID, orgID)
	if err != nil {
		return false, fmt.Errorf("failed to check ownership: %w", err)
	}
	if isOwner {
		return true, nil
	}

	// TODO: Verificar si es empleado con permisos administrativos
	// Para simplificar, por ahora solo los propietarios pueden ver la información

	return false, nil
}

// updateMajorityOwnerStatus actualiza el estatus de propietario mayoritario
func (s *OwnerService) updateMajorityOwnerStatus(ctx context.Context, orgID uuid.UUID) error {
	s.logger.Debug("Updating majority owner status", zap.String("org_id", orgID.String()))

	// Esta operación debe ser implementada en el repositorio para ser atómica
	// Por ahora solo registramos que debería hacerse
	s.logger.Info("Majority owner status update placeholder", zap.String("org_id", orgID.String()))

	// TODO: Implementar lógica real de actualización de majority owner
	// 1. Obtener todos los propietarios usando GetMajorityOwner del repositorio
	// 2. Actualizar los registros según los porcentajes de propiedad

	return nil
}

// ===============================
// EVENT PUBLISHING METHODS
// ===============================

// publishOwnerAddedEvent publica evento cuando se añade un propietario
func (s *OwnerService) publishOwnerAddedEvent(ctx context.Context, owner *models.OrganizationOwner, org *models.Organization) error {
	return s.events.PublishOwnerAdded(org, owner)
}

// publishOwnerRemovedEvent publica evento cuando se remueve un propietario
func (s *OwnerService) publishOwnerRemovedEvent(ctx context.Context, owner *models.OrganizationOwner) error {
	return s.events.PublishOwnerRemoved(owner.OrganizationID.String(), owner.ID.String())
}
