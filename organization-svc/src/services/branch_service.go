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

// BranchService maneja la lógica de negocio para sucursales organizacionales
type BranchService struct {
	repos   *repository.Repositories
	clients *clients.ClientManager
	cache   CacheService
	events  EventService
	logger  *zap.Logger
}

// NewBranchService crea una nueva instancia del servicio de sucursales
func NewBranchService(
	repos *repository.Repositories,
	clients *clients.ClientManager,
	cache CacheService,
	events EventService,
	logger *zap.Logger,
) *BranchService {
	return &BranchService{
		repos:   repos,
		clients: clients,
		cache:   cache,
		events:  events,
		logger:  logger,
	}
}

// ===============================
// CORE BRANCH OPERATIONS
// ===============================

// CreateBranch crea una nueva sucursal para una organización
func (s *BranchService) CreateBranch(ctx context.Context, orgID uuid.UUID, creatorPersonID uuid.UUID, req dto.CreateBranchRequest) (*models.OrganizationBranch, error) {
	s.logger.Info("Creating branch",
		zap.String("org_id", orgID.String()),
		zap.String("creator_person_id", creatorPersonID.String()),
		zap.String("display_name", req.DisplayName),
	)

	// 1. Validar que la organización existe y está activa
	org, err := s.repos.Organization.GetByID(ctx, orgID)
	if err != nil {
		s.logger.Error("Failed to get organization", zap.Error(err))
		return nil, rcerrors.NewNotFoundError("organization", orgID.String())
	}

	if org.Status != models.OrgStatusActive {
		return nil, rcerrors.NewValidationError("organization", "organization must be active to create branches")
	}

	// 2. Validar address_id via gRPC si se proporciona
	if req.AddressID != nil {
		if err := s.validateAddressViaGRPC(ctx, *req.AddressID); err != nil {
			s.logger.Error("Address validation failed", zap.Error(err))
			return nil, rcerrors.NewValidationError("address_id", "invalid address")
		}
	}

	// 3. Si se marca como principal, verificar que no haya otra principal
	if req.IsMain {
		hasMain, err := s.repos.Branch.HasMainBranch(ctx, orgID)
		if err != nil {
			s.logger.Error("Failed to check main branch existence", zap.Error(err))
			return nil, rcerrors.NewInternalServerError("failed to validate main branch")
		}
		if hasMain {
			return nil, rcerrors.NewConflictError("organization already has a main branch")
		}
	}

	// 4. Verificar que no existe otra sucursal con el mismo nombre en la organización
	exists, err := s.repos.Branch.ExistsByName(ctx, orgID, req.DisplayName, nil)
	if err != nil {
		s.logger.Error("Failed to check branch name uniqueness", zap.Error(err))
		return nil, rcerrors.NewInternalServerError("failed to validate branch name")
	}
	if exists {
		return nil, rcerrors.NewConflictError("branch with this name already exists in the organization")
	}

	// 5. Crear la sucursal
	branch := &models.OrganizationBranch{
		OrganizationID: orgID,
		DisplayName:    req.DisplayName,
		AddressID:      req.AddressID,
		Phone:          req.Phone,
		Email:          req.Email,
		IsMain:         req.IsMain,
		CreatedBy:      creatorPersonID,
	}

	// 6. Guardar en base de datos
	createdBranch, err := s.repos.Branch.Create(ctx, branch)
	if err != nil {
		s.logger.Error("Failed to create branch", zap.Error(err))
		return nil, rcerrors.NewInternalServerError("failed to create branch")
	}

	// 7. Invalidar cache de branches de la organización
	if err := s.cache.InvalidateOrganizationCache(orgID.String()); err != nil {
		s.logger.Warn("Failed to invalidate organization cache", zap.Error(err))
	}

	// 8. Publicar evento de creación
	if err := s.publishBranchCreatedEvent(ctx, createdBranch, org); err != nil {
		s.logger.Warn("Failed to publish branch created event", zap.Error(err))
	}

	s.logger.Info("Branch created successfully",
		zap.String("branch_id", createdBranch.ID.String()),
		zap.String("org_id", orgID.String()),
	)

	return createdBranch, nil
}

// GetBranch obtiene una sucursal por ID
func (s *BranchService) GetBranch(ctx context.Context, branchID uuid.UUID, requesterPersonID uuid.UUID) (*models.OrganizationBranch, error) {
	s.logger.Debug("Getting branch",
		zap.String("branch_id", branchID.String()),
		zap.String("requester_person_id", requesterPersonID.String()),
	)

	// 1. Obtener la sucursal
	branch, err := s.repos.Branch.GetByID(ctx, branchID)
	if err != nil {
		s.logger.Error("Failed to get branch", zap.Error(err))
		return nil, rcerrors.NewNotFoundError("branch", branchID.String())
	}

	// 2. Verificar que el usuario tiene acceso a la organización
	hasAccess, err := s.validateOrganizationAccess(ctx, branch.OrganizationID, requesterPersonID)
	if err != nil {
		s.logger.Error("Failed to validate organization access", zap.Error(err))
		return nil, rcerrors.NewInternalServerError("failed to validate access")
	}
	if !hasAccess {
		return nil, rcerrors.NewForbiddenError("insufficient permissions to access this branch")
	}

	return branch, nil
}

// ListBranches lista sucursales con filtros
func (s *BranchService) ListBranches(ctx context.Context, requesterPersonID uuid.UUID, filters dto.BranchFiltersRequest) ([]models.OrganizationBranch, int64, error) {
	s.logger.Debug("Listing branches",
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
		return nil, 0, rcerrors.NewForbiddenError("insufficient permissions to access organization branches")
	}

	// 2. Listar sucursales
	branches, total, err := s.repos.Branch.List(ctx, filters)
	if err != nil {
		s.logger.Error("Failed to list branches", zap.Error(err))
		return nil, 0, rcerrors.NewInternalServerError("failed to list branches")
	}

	return branches, total, nil
}

// UpdateBranch actualiza una sucursal existente
func (s *BranchService) UpdateBranch(ctx context.Context, branchID uuid.UUID, updaterPersonID uuid.UUID, req dto.UpdateBranchRequest) (*models.OrganizationBranch, error) {
	s.logger.Info("Updating branch",
		zap.String("branch_id", branchID.String()),
		zap.String("updater_person_id", updaterPersonID.String()),
	)

	// 1. Obtener la sucursal actual
	branch, err := s.repos.Branch.GetByID(ctx, branchID)
	if err != nil {
		s.logger.Error("Failed to get branch", zap.Error(err))
		return nil, rcerrors.NewNotFoundError("branch", branchID.String())
	}

	// 2. Verificar permisos de edición
	hasAccess, err := s.validateOrganizationAccess(ctx, branch.OrganizationID, updaterPersonID)
	if err != nil {
		s.logger.Error("Failed to validate organization access", zap.Error(err))
		return nil, rcerrors.NewInternalServerError("failed to validate access")
	}
	if !hasAccess {
		return nil, rcerrors.NewForbiddenError("insufficient permissions to update this branch")
	}

	// 3. Validar address_id via gRPC si se proporciona
	if req.AddressID != nil {
		if err := s.validateAddressViaGRPC(ctx, *req.AddressID); err != nil {
			s.logger.Error("Address validation failed", zap.Error(err))
			return nil, rcerrors.NewValidationError("address_id", "invalid address")
		}
	}

	// 4. Validar unicidad de nombre si se está cambiando
	if req.DisplayName != nil && *req.DisplayName != branch.DisplayName {
		excludeID := branchID
		exists, err := s.repos.Branch.ExistsByName(ctx, branch.OrganizationID, *req.DisplayName, &excludeID)
		if err != nil {
			s.logger.Error("Failed to check branch name uniqueness", zap.Error(err))
			return nil, rcerrors.NewInternalServerError("failed to validate branch name")
		}
		if exists {
			return nil, rcerrors.NewConflictError("branch with this name already exists in the organization")
		}
	}

	// 5. Preparar actualizaciones
	updates := make(map[string]interface{})
	if req.DisplayName != nil {
		updates["display_name"] = *req.DisplayName
	}
	if req.AddressID != nil {
		updates["address_id"] = *req.AddressID
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.Metadata != nil {
		updates["metadata"] = req.Metadata
	}

	// Siempre actualizar updated_at
	updates["updated_at"] = time.Now()

	// 6. Aplicar actualizaciones
	updatedBranch, err := s.repos.Branch.Update(ctx, branchID, updates)
	if err != nil {
		s.logger.Error("Failed to update branch", zap.Error(err))
		return nil, rcerrors.NewInternalServerError("failed to update branch")
	}

	// 7. Invalidar cache
	if err := s.cache.InvalidateOrganizationCache(branch.OrganizationID.String()); err != nil {
		s.logger.Warn("Failed to invalidate organization cache", zap.Error(err))
	}

	// 8. Publicar evento de actualización
	if err := s.publishBranchUpdatedEvent(ctx, updatedBranch); err != nil {
		s.logger.Warn("Failed to publish branch updated event", zap.Error(err))
	}

	s.logger.Info("Branch updated successfully", zap.String("branch_id", branchID.String()))
	return updatedBranch, nil
}

// DeleteBranch elimina una sucursal (soft delete)
func (s *BranchService) DeleteBranch(ctx context.Context, branchID uuid.UUID, deleterPersonID uuid.UUID) error {
	s.logger.Info("Deleting branch",
		zap.String("branch_id", branchID.String()),
		zap.String("deleter_person_id", deleterPersonID.String()),
	)

	// 1. Obtener la sucursal
	branch, err := s.repos.Branch.GetByID(ctx, branchID)
	if err != nil {
		s.logger.Error("Failed to get branch", zap.Error(err))
		return rcerrors.NewNotFoundError("branch", branchID.String())
	}

	// 2. Verificar permisos
	hasAccess, err := s.validateOrganizationAccess(ctx, branch.OrganizationID, deleterPersonID)
	if err != nil {
		s.logger.Error("Failed to validate organization access", zap.Error(err))
		return rcerrors.NewInternalServerError("failed to validate access")
	}
	if !hasAccess {
		return rcerrors.NewForbiddenError("insufficient permissions to delete this branch")
	}

	// 3. Verificar que no es la sucursal principal
	if branch.IsMain {
		return rcerrors.NewValidationError("is_main", "cannot delete the main branch")
	}

	// 4. Verificar que no hay empleados activos asignados a esta sucursal
	// Para esta implementación simplificada, asumimos que el repositorio tiene el método CanDelete
	canDelete, err := s.repos.Branch.CanDelete(ctx, branchID)
	if err != nil {
		s.logger.Error("Failed to check if branch can be deleted", zap.Error(err))
		return rcerrors.NewInternalServerError("failed to validate branch deletion")
	}
	if !canDelete {
		return rcerrors.NewValidationError("employees", "cannot delete branch with active employees or other dependencies")
	}

	// 5. Realizar soft delete
	if err := s.repos.Branch.SoftDelete(ctx, branchID, deleterPersonID); err != nil {
		s.logger.Error("Failed to delete branch", zap.Error(err))
		return rcerrors.NewInternalServerError("failed to delete branch")
	}

	// 6. Invalidar cache
	if err := s.cache.InvalidateOrganizationCache(branch.OrganizationID.String()); err != nil {
		s.logger.Warn("Failed to invalidate organization cache", zap.Error(err))
	}

	// 7. Publicar evento de eliminación
	if err := s.publishBranchDeletedEvent(ctx, branch); err != nil {
		s.logger.Warn("Failed to publish branch deleted event", zap.Error(err))
	}

	s.logger.Info("Branch deleted successfully", zap.String("branch_id", branchID.String()))
	return nil
}

// SetMainBranch establece una sucursal como principal
func (s *BranchService) SetMainBranch(ctx context.Context, branchID uuid.UUID, setterPersonID uuid.UUID) error {
	s.logger.Info("Setting main branch",
		zap.String("branch_id", branchID.String()),
		zap.String("setter_person_id", setterPersonID.String()),
	)

	// 1. Obtener la sucursal
	branch, err := s.repos.Branch.GetByID(ctx, branchID)
	if err != nil {
		s.logger.Error("Failed to get branch", zap.Error(err))
		return rcerrors.NewNotFoundError("branch", branchID.String())
	}

	// 2. Verificar permisos
	hasAccess, err := s.validateOrganizationAccess(ctx, branch.OrganizationID, setterPersonID)
	if err != nil {
		s.logger.Error("Failed to validate organization access", zap.Error(err))
		return rcerrors.NewInternalServerError("failed to validate access")
	}
	if !hasAccess {
		return rcerrors.NewForbiddenError("insufficient permissions to set main branch")
	}

	// 3. Verificar que no es ya la principal
	if branch.IsMain {
		return rcerrors.NewValidationError("is_main", "branch is already the main branch")
	}

	// 4. Transacción para cambiar sucursal principal
	if err := s.repos.Branch.SetAsMain(ctx, branchID, branch.OrganizationID); err != nil {
		s.logger.Error("Failed to set main branch", zap.Error(err))
		return rcerrors.NewInternalServerError("failed to set main branch")
	}

	// 5. Invalidar cache
	if err := s.cache.InvalidateOrganizationCache(branch.OrganizationID.String()); err != nil {
		s.logger.Warn("Failed to invalidate organization cache", zap.Error(err))
	}

	// 6. Publicar evento
	if err := s.publishMainBranchChangedEvent(ctx, branch); err != nil {
		s.logger.Warn("Failed to publish main branch changed event", zap.Error(err))
	}

	s.logger.Info("Main branch set successfully", zap.String("branch_id", branchID.String()))
	return nil
}

// ===============================
// VALIDATION HELPER METHODS
// ===============================

// validateAddressViaGRPC valida que una dirección existe usando el servicio de direcciones
func (s *BranchService) validateAddressViaGRPC(ctx context.Context, addressID uuid.UUID) error {
	s.logger.Debug("Validating address via gRPC", zap.String("address_id", addressID.String()))

	// Para esta implementación simplificada, simulamos la validación gRPC
	// En una implementación completa, se usaría el cliente real de Address
	// exists, err := s.clients.Address.ValidateAddress(ctx, addressID.String())
	// if err != nil {
	//     s.logger.Error("gRPC address validation failed", zap.Error(err))
	//     return fmt.Errorf("failed to validate address: %w", err)
	// }
	// if !exists {
	//     return fmt.Errorf("address not found")
	// }

	// Simulación: asumimos que la dirección es válida por ahora
	s.logger.Info("Address validation skipped (mock implementation)", zap.String("address_id", addressID.String()))
	return nil
}

// validateOrganizationAccess verifica que un usuario tiene acceso a una organización
func (s *BranchService) validateOrganizationAccess(ctx context.Context, orgID uuid.UUID, personID uuid.UUID) (bool, error) {
	s.logger.Debug("Validating organization access",
		zap.String("org_id", orgID.String()),
		zap.String("person_id", personID.String()),
	)

	// Para esta implementación simplificada, verificamos mediante employee repository
	// En una implementación completa, se verificaría ownership y employee status

	// Verificar si es empleado activo de la organización
	// Para simplificar, asumimos que el usuario tiene acceso si la organización existe
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

// publishBranchCreatedEvent publica evento cuando se crea una sucursal
func (s *BranchService) publishBranchCreatedEvent(ctx context.Context, branch *models.OrganizationBranch, org *models.Organization) error {
	return s.events.PublishBranchCreated(org, branch)
}

// publishBranchUpdatedEvent publica evento cuando se actualiza una sucursal
func (s *BranchService) publishBranchUpdatedEvent(ctx context.Context, branch *models.OrganizationBranch) error {
	// Para obtener la organización
	org, err := s.repos.Organization.GetByID(ctx, branch.OrganizationID)
	if err != nil {
		s.logger.Error("Failed to get organization for event", zap.Error(err))
		return err
	}
	return s.events.PublishBranchUpdated(org, branch)
}

// publishBranchDeletedEvent publica evento cuando se elimina una sucursal
func (s *BranchService) publishBranchDeletedEvent(ctx context.Context, branch *models.OrganizationBranch) error {
	return s.events.PublishBranchDeleted(branch.OrganizationID.String(), branch.ID.String())
}

// publishMainBranchChangedEvent publica evento cuando cambia la sucursal principal
func (s *BranchService) publishMainBranchChangedEvent(ctx context.Context, branch *models.OrganizationBranch) error {
	// Para el evento de cambio de sucursal principal, podemos usar el evento de actualización
	org, err := s.repos.Organization.GetByID(ctx, branch.OrganizationID)
	if err != nil {
		s.logger.Error("Failed to get organization for event", zap.Error(err))
		return err
	}
	return s.events.PublishBranchUpdated(org, branch)
}
