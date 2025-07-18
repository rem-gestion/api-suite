package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/rem-gestion/api-suite/organization/src/clients"
	"github.com/rem-gestion/api-suite/organization/src/dto"
	"github.com/rem-gestion/api-suite/organization/src/models"
	"github.com/rem-gestion/api-suite/organization/src/repository"
	rcerrors "github.com/rem-gestion/rem-common/errors"
)

// RoleService maneja la lógica de negocio para roles organizacionales
type RoleService struct {
	repos   *repository.Repositories
	clients *clients.ClientManager
	cache   CacheService
	events  EventService
	logger  *zap.Logger
}

// NewRoleService crea una nueva instancia del servicio de roles
func NewRoleService(
	repos *repository.Repositories,
	clients *clients.ClientManager,
	cache CacheService,
	events EventService,
	logger *zap.Logger,
) *RoleService {
	return &RoleService{
		repos:   repos,
		clients: clients,
		cache:   cache,
		events:  events,
		logger:  logger,
	}
}

// ===============================
// CORE ROLE OPERATIONS
// ===============================

// CreateRole crea un nuevo rol para una organización
func (s *RoleService) CreateRole(ctx context.Context, orgID uuid.UUID, creatorPersonID uuid.UUID, req dto.CreateRoleRequest) (*models.OrganizationRole, error) {
	s.logger.Info("Creating role",
		zap.String("org_id", orgID.String()),
		zap.String("creator_person_id", creatorPersonID.String()),
		zap.String("name", req.Name),
	)

	// 1. Validar que la organización existe y está activa
	org, err := s.repos.Organization.GetByID(ctx, orgID)
	if err != nil {
		s.logger.Error("Failed to get organization", zap.Error(err))
		return nil, rcerrors.NewNotFoundError("organization", orgID.String())
	}

	if org.Status != models.OrgStatusActive {
		return nil, rcerrors.NewValidationError("organization", "organization must be active to create roles")
	}

	// 2. Verificar permisos para crear roles
	hasAccess, err := s.validateOrganizationAccess(ctx, orgID, creatorPersonID)
	if err != nil {
		s.logger.Error("Failed to validate organization access", zap.Error(err))
		return nil, rcerrors.NewInternalServerError("failed to validate access")
	}
	if !hasAccess {
		return nil, rcerrors.NewForbiddenError("insufficient permissions to create roles")
	}

	// 3. Verificar que no existe otro rol con el mismo nombre en la organización
	exists, err := s.repos.Role.ExistsByName(ctx, orgID, req.Name, nil)
	if err != nil {
		s.logger.Error("Failed to check role name uniqueness", zap.Error(err))
		return nil, rcerrors.NewInternalServerError("failed to validate role name")
	}
	if exists {
		return nil, rcerrors.NewConflictError("role with this name already exists in the organization")
	}

	// 4. Si se marca como default, verificar que no haya otro rol default
	if req.IsDefault {
		defaultRoles, err := s.repos.Role.GetDefaultRoles(ctx, orgID)
		if err != nil {
			s.logger.Error("Failed to check default roles", zap.Error(err))
			return nil, rcerrors.NewInternalServerError("failed to validate default role")
		}
		if len(defaultRoles) > 0 {
			return nil, rcerrors.NewConflictError("organization already has a default role")
		}
	}

	// 5. Crear el rol
	role := &models.OrganizationRole{
		OrganizationID: orgID,
		Name:           req.Name,
		Description:    req.Description,
		IsDefault:      req.IsDefault,
		CreatedBy:      creatorPersonID,
	}

	// 6. Guardar en base de datos
	createdRole, err := s.repos.Role.Create(ctx, role)
	if err != nil {
		s.logger.Error("Failed to create role", zap.Error(err))
		return nil, rcerrors.NewInternalServerError("failed to create role")
	}

	// 7. Invalidar cache de roles de la organización
	if err := s.cache.InvalidateOrganizationCache(orgID.String()); err != nil {
		s.logger.Warn("Failed to invalidate organization cache", zap.Error(err))
	}

	// 8. Publicar evento de creación
	if err := s.publishRoleCreatedEvent(ctx, createdRole, org); err != nil {
		s.logger.Warn("Failed to publish role created event", zap.Error(err))
	}

	s.logger.Info("Role created successfully",
		zap.String("role_id", createdRole.ID.String()),
		zap.String("org_id", orgID.String()),
	)

	return createdRole, nil
}

// GetRole obtiene un rol por ID
func (s *RoleService) GetRole(ctx context.Context, roleID uuid.UUID, requesterPersonID uuid.UUID) (*models.OrganizationRole, error) {
	s.logger.Debug("Getting role",
		zap.String("role_id", roleID.String()),
		zap.String("requester_person_id", requesterPersonID.String()),
	)

	// 1. Obtener el rol
	role, err := s.repos.Role.GetByID(ctx, roleID)
	if err != nil {
		s.logger.Error("Failed to get role", zap.Error(err))
		return nil, rcerrors.NewNotFoundError("role", roleID.String())
	}

	// 2. Verificar que el usuario tiene acceso a la organización
	hasAccess, err := s.validateOrganizationAccess(ctx, role.OrganizationID, requesterPersonID)
	if err != nil {
		s.logger.Error("Failed to validate organization access", zap.Error(err))
		return nil, rcerrors.NewInternalServerError("failed to validate access")
	}
	if !hasAccess {
		return nil, rcerrors.NewForbiddenError("insufficient permissions to access this role")
	}

	return role, nil
}

// ListRoles lista roles de una organización con filtros
func (s *RoleService) ListRoles(ctx context.Context, requesterPersonID uuid.UUID, filters dto.RoleFiltersRequest) ([]models.OrganizationRole, int64, error) {
	s.logger.Debug("Listing roles",
		zap.String("requester_person_id", requesterPersonID.String()),
		zap.String("org_id", filters.OrganizationID.String()),
	)

	// 1. Verificar acceso a la organización
	hasAccess, err := s.validateOrganizationAccess(ctx, filters.OrganizationID, requesterPersonID)
	if err != nil {
		s.logger.Error("Failed to validate organization access", zap.Error(err))
		return nil, 0, rcerrors.NewInternalServerError("failed to validate access")
	}
	if !hasAccess {
		return nil, 0, rcerrors.NewForbiddenError("insufficient permissions to access organization roles")
	}

	// 2. Listar roles
	roles, total, err := s.repos.Role.List(ctx, filters)
	if err != nil {
		s.logger.Error("Failed to list roles", zap.Error(err))
		return nil, 0, rcerrors.NewInternalServerError("failed to list roles")
	}

	return roles, total, nil
}

// UpdateRole actualiza un rol existente
func (s *RoleService) UpdateRole(ctx context.Context, roleID uuid.UUID, updaterPersonID uuid.UUID, req dto.UpdateRoleRequest) (*models.OrganizationRole, error) {
	s.logger.Info("Updating role",
		zap.String("role_id", roleID.String()),
		zap.String("updater_person_id", updaterPersonID.String()),
	)

	// 1. Obtener el rol actual
	role, err := s.repos.Role.GetByID(ctx, roleID)
	if err != nil {
		s.logger.Error("Failed to get role", zap.Error(err))
		return nil, rcerrors.NewNotFoundError("role", roleID.String())
	}

	// 2. Verificar permisos de edición
	hasAccess, err := s.validateOrganizationAccess(ctx, role.OrganizationID, updaterPersonID)
	if err != nil {
		s.logger.Error("Failed to validate organization access", zap.Error(err))
		return nil, rcerrors.NewInternalServerError("failed to validate access")
	}
	if !hasAccess {
		return nil, rcerrors.NewForbiddenError("insufficient permissions to update this role")
	}

	// 3. Validar unicidad de nombre si se está cambiando
	if req.Name != nil && *req.Name != role.Name {
		excludeID := roleID
		exists, err := s.repos.Role.ExistsByName(ctx, role.OrganizationID, *req.Name, &excludeID)
		if err != nil {
			s.logger.Error("Failed to check role name uniqueness", zap.Error(err))
			return nil, rcerrors.NewInternalServerError("failed to validate role name")
		}
		if exists {
			return nil, rcerrors.NewConflictError("role with this name already exists in the organization")
		}
	}

	// 4. Validar cambio de rol default - Por ahora removido ya que UpdateRoleRequest no tiene IsDefault
	// if req.IsDefault != nil && *req.IsDefault && !role.IsDefault {
	//     defaultRoles, err := s.repos.Role.GetDefaultRoles(ctx, role.OrganizationID)
	//     if err != nil {
	//         s.logger.Error("Failed to check default roles", zap.Error(err))
	//         return nil, rcerrors.NewInternalServerError("failed to validate default role")
	//     }
	//     if len(defaultRoles) > 0 {
	//         return nil, rcerrors.NewConflictError("organization already has a default role")
	//     }
	// }

	// 5. Preparar actualizaciones
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	// IsDefault no está disponible en UpdateRoleRequest por ahora
	// if req.IsDefault != nil {
	//     updates["is_default"] = *req.IsDefault
	// }

	// Siempre actualizar updated_at y updated_by
	updates["updated_at"] = time.Now()
	updates["updated_by"] = updaterPersonID

	// 6. Aplicar actualizaciones
	updatedRole, err := s.repos.Role.Update(ctx, roleID, updates)
	if err != nil {
		s.logger.Error("Failed to update role", zap.Error(err))
		return nil, rcerrors.NewInternalServerError("failed to update role")
	}

	// 7. Invalidar cache
	if err := s.cache.InvalidateOrganizationCache(role.OrganizationID.String()); err != nil {
		s.logger.Warn("Failed to invalidate organization cache", zap.Error(err))
	}

	// 8. Publicar evento de actualización
	if err := s.publishRoleUpdatedEvent(ctx, updatedRole); err != nil {
		s.logger.Warn("Failed to publish role updated event", zap.Error(err))
	}

	s.logger.Info("Role updated successfully", zap.String("role_id", roleID.String()))
	return updatedRole, nil
}

// DeleteRole elimina un rol (soft delete)
func (s *RoleService) DeleteRole(ctx context.Context, roleID uuid.UUID, deleterPersonID uuid.UUID) error {
	s.logger.Info("Deleting role",
		zap.String("role_id", roleID.String()),
		zap.String("deleter_person_id", deleterPersonID.String()),
	)

	// 1. Obtener el rol
	role, err := s.repos.Role.GetByID(ctx, roleID)
	if err != nil {
		s.logger.Error("Failed to get role", zap.Error(err))
		return rcerrors.NewNotFoundError("role", roleID.String())
	}

	// 2. Verificar permisos
	hasAccess, err := s.validateOrganizationAccess(ctx, role.OrganizationID, deleterPersonID)
	if err != nil {
		s.logger.Error("Failed to validate organization access", zap.Error(err))
		return rcerrors.NewInternalServerError("failed to validate access")
	}
	if !hasAccess {
		return rcerrors.NewForbiddenError("insufficient permissions to delete this role")
	}

	// 3. Verificar que no es el rol por defecto
	if role.IsDefault {
		return rcerrors.NewValidationError("is_default", "cannot delete the default role")
	}

	// 4. Verificar que no hay empleados activos con este rol
	canDelete, err := s.repos.Role.CanDelete(ctx, roleID)
	if err != nil {
		s.logger.Error("Failed to check if role can be deleted", zap.Error(err))
		return rcerrors.NewInternalServerError("failed to validate role deletion")
	}
	if !canDelete {
		return rcerrors.NewValidationError("employees", "cannot delete role assigned to active employees")
	}

	// 5. Realizar soft delete
	if err := s.repos.Role.SoftDelete(ctx, roleID, deleterPersonID); err != nil {
		s.logger.Error("Failed to delete role", zap.Error(err))
		return rcerrors.NewInternalServerError("failed to delete role")
	}

	// 6. Invalidar cache
	if err := s.cache.InvalidateOrganizationCache(role.OrganizationID.String()); err != nil {
		s.logger.Warn("Failed to invalidate organization cache", zap.Error(err))
	}

	// 7. Publicar evento de eliminación
	if err := s.publishRoleDeletedEvent(ctx, role); err != nil {
		s.logger.Warn("Failed to publish role deleted event", zap.Error(err))
	}

	s.logger.Info("Role deleted successfully", zap.String("role_id", roleID.String()))
	return nil
}

// ===============================
// VALIDATION HELPER METHODS
// ===============================

// validateOrganizationAccess verifica que un usuario tiene acceso a una organización
func (s *RoleService) validateOrganizationAccess(ctx context.Context, orgID uuid.UUID, personID uuid.UUID) (bool, error) {
	s.logger.Debug("Validating organization access",
		zap.String("org_id", orgID.String()),
		zap.String("person_id", personID.String()),
	)

	// Para esta implementación simplificada, verificamos que la organización existe
	_, err := s.repos.Organization.GetByID(ctx, orgID)
	if err != nil {
		return false, fmt.Errorf("organization not found: %w", err)
	}

	// TODO: Implementar verificación real de permisos via OwnerRepository y EmployeeRepository
	// isOwner, err := s.repos.Owner.IsOwner(ctx, orgID, personID)
	// isEmployee, err := s.repos.Employee.IsActiveEmployee(ctx, orgID, personID)

	// Por ahora retornamos true para permitir el desarrollo
	s.logger.Info("Access validation simplified (mock implementation)",
		zap.String("org_id", orgID.String()),
		zap.String("person_id", personID.String()),
	)
	return true, nil
}

// ===============================
// EVENT PUBLISHING METHODS
// ===============================

// publishRoleCreatedEvent publica evento cuando se crea un rol
func (s *RoleService) publishRoleCreatedEvent(ctx context.Context, role *models.OrganizationRole, org *models.Organization) error {
	return s.events.PublishRoleCreated(org, role)
}

// publishRoleUpdatedEvent publica evento cuando se actualiza un rol
func (s *RoleService) publishRoleUpdatedEvent(ctx context.Context, role *models.OrganizationRole) error {
	// Para obtener la organización
	org, err := s.repos.Organization.GetByID(ctx, role.OrganizationID)
	if err != nil {
		s.logger.Error("Failed to get organization for event", zap.Error(err))
		return err
	}
	return s.events.PublishRoleUpdated(org, role)
}

// publishRoleDeletedEvent publica evento cuando se elimina un rol
func (s *RoleService) publishRoleDeletedEvent(ctx context.Context, role *models.OrganizationRole) error {
	return s.events.PublishRoleDeleted(role.OrganizationID.String(), role.ID.String())
}
