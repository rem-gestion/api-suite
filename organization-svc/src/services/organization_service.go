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

// OrganizationService maneja la lógica de negocio para organizaciones
type OrganizationService struct {
	repos   *repository.Repositories
	clients *clients.ClientManager
	cache   CacheService
	events  EventService
	logger  *zap.Logger
}

// NewOrganizationService crea una nueva instancia del servicio
func NewOrganizationService(
	repos *repository.Repositories,
	clients *clients.ClientManager,
	cache CacheService,
	events EventService,
	logger *zap.Logger,
) *OrganizationService {
	return &OrganizationService{
		repos:   repos,
		clients: clients,
		cache:   cache,
		events:  events,
		logger:  logger,
	}
}

// ===============================
// CORE ORGANIZATION OPERATIONS
// ===============================

// CreateOrganization crea una nueva organización con validaciones y configuraciones por defecto
func (s *OrganizationService) CreateOrganization(ctx context.Context, req *dto.CreateOrganizationRequest, createdBy uuid.UUID) (*models.Organization, error) {
	s.logger.Info("Creating organization",
		zap.String("name", req.Name),
		zap.String("type", req.Type),
		zap.String("created_by", createdBy.String()),
	)

	// 1. Validar que el usuario creador existe
	_, err := s.clients.AuthIdentity.GetUserById(createdBy.String())

	fmt.Println("Creator user validation:", createdBy.String())
	if err != nil {
		s.logger.Error("Creator user validation failed", zap.Error(err))
		return nil, rcerrors.NewValidationError("created_by", "invalid creator user_id")
	}

	// 2. Validar dirección fiscal si se proporciona
	if req.FiscalAddressID != nil {
		_, err := s.clients.Address.GetAddress(req.FiscalAddressID.String())
		if err != nil {
			s.logger.Error("Fiscal address validation failed", zap.Error(err))
			return nil, rcerrors.NewValidationError("fiscal_address_id", "invalid fiscal address")
		}
	}

	// 3. Verificar que el slug no esté duplicado
	if req.Slug != "" {
		exists, err := s.repos.Organization.ExistsBySlug(ctx, req.Slug)
		if err != nil {
			s.logger.Error("Failed to check slug existence", zap.Error(err))
			return nil, rcerrors.NewInternalServerError("failed to validate slug")
		}
		if exists {
			return nil, rcerrors.NewConflictError("organization slug already exists")
		}
	}

	// 4. Crear organización
	org := &models.Organization{
		Name:            req.Name,
		DisplayName:     req.DisplayName,
		Slug:            &req.Slug,
		Description:     req.Description,
		Type:            &req.Type,
		LegalName:       req.LegalName,
		TaxID:           req.TaxID,
		Website:         req.Website,
		Phone:           req.Phone,
		Email:           req.Email,
		LogoURL:         req.LogoURL,
		TimezoneID:      &req.TimezoneID,
		FiscalAddressID: req.FiscalAddressID,
		Status:          models.OrgStatusPending,
		CreatedBy:       createdBy,
		UpdatedBy:       &createdBy,
	}

	// 5. Crear en base de datos
	createdOrg, err := s.repos.Organization.Create(ctx, org)
	if err != nil {
		s.logger.Error("Failed to create organization", zap.Error(err))
		return nil, rcerrors.NewInternalServerError("failed to create organization")
	}

	// 6. Crear configuraciones por defecto
	if err := s.createDefaultSettings(ctx, createdOrg.ID, createdBy); err != nil {
		s.logger.Warn("Failed to create default settings", zap.Error(err))
		// No es un error crítico, continúa
	}

	// 7. Crear branch principal por defecto
	if err := s.createMainBranch(ctx, createdOrg, createdBy); err != nil {
		s.logger.Warn("Failed to create main branch", zap.Error(err))
		// No es un error crítico para la creación de la organización
	}

	// 8. Crear rol de propietario por defecto
	if err := s.createDefaultOwnerRole(ctx, createdOrg.ID, createdBy); err != nil {
		s.logger.Warn("Failed to create default owner role", zap.Error(err))
	}

	// 9. Invalidar cache
	s.invalidateOrganizationCache(createdOrg.ID.String())

	// 10. Publicar evento
	if err := s.events.PublishOrganizationCreated(createdOrg); err != nil {
		s.logger.Warn("Failed to publish organization created event", zap.Error(err))
	}

	s.logger.Info("Organization created successfully",
		zap.String("org_id", createdOrg.ID.String()),
		zap.String("name", createdOrg.Name),
	)

	return createdOrg, nil
}

// GetOrganization obtiene una organización por ID con validaciones de acceso
func (s *OrganizationService) GetOrganization(ctx context.Context, orgID uuid.UUID, requesterUserID uuid.UUID) (*models.Organization, error) {
	s.logger.Debug("Getting organization",
		zap.String("org_id", orgID.String()),
		zap.String("requester_user_id", requesterUserID.String()),
	)

	// 1. Verificar acceso del usuario a la organización
	hasAccess, err := s.validateUserAccess(ctx, orgID, requesterUserID)
	if err != nil {
		s.logger.Error("Failed to validate user access", zap.Error(err))
		return nil, rcerrors.NewInternalServerError("failed to validate access")
	}
	if !hasAccess {
		return nil, rcerrors.NewForbiddenError("user does not have access to this organization")
	}

	// 2. Obtener de cache primero
	if org := s.getOrganizationFromCache(orgID.String()); org != nil {
		return org, nil
	}

	// 3. Obtener de base de datos
	org, err := s.repos.Organization.GetByID(ctx, orgID)
	if err != nil {
		s.logger.Error("Failed to get organization", zap.Error(err))
		return nil, rcerrors.NewNotFoundError("organization", orgID.String())
	}

	// 4. Guardar en cache
	s.setOrganizationCache(orgID.String(), org)

	return org, nil
}

// ListOrganizations lista organizaciones según filtros y permisos del usuario
func (s *OrganizationService) ListOrganizations(ctx context.Context, requesterUserID uuid.UUID, filters dto.OrganizationFiltersRequest) ([]models.Organization, int64, error) {
	s.logger.Debug("Listing organizations",
		zap.String("requester_user_id", requesterUserID.String()),
		zap.Int("page", filters.Page),
		zap.Int("per_page", filters.PerPage),
	)

	// 1. Validar usuario existe
	_, err := s.clients.AuthIdentity.GetUserById(requesterUserID.String())
	if err != nil {
		s.logger.Error("User validation failed", zap.Error(err))
		return nil, 0, rcerrors.NewValidationError("requester_user_id", "invalid user")
	}

	// 2. Obtener organizaciones donde el usuario tiene acceso
	organizations, total, err := s.repos.Organization.List(ctx, filters)
	if err != nil {
		s.logger.Error("Failed to list organizations", zap.Error(err))
		return nil, 0, rcerrors.NewInternalServerError("failed to list organizations")
	}

	return organizations, total, nil
}

// UpdateOrganization actualiza una organización existente
func (s *OrganizationService) UpdateOrganization(ctx context.Context, orgID uuid.UUID, req *dto.UpdateOrganizationRequest, updatedBy uuid.UUID) (*models.Organization, error) {
	s.logger.Info("Updating organization",
		zap.String("org_id", orgID.String()),
		zap.String("updated_by", updatedBy.String()),
	)

	// 1. Verificar que la organización existe
	_, err := s.repos.Organization.GetByID(ctx, orgID)
	if err != nil {
		return nil, rcerrors.NewNotFoundError("organization", orgID.String())
	}

	// 2. Validar permisos de edición
	hasPermission, err := s.validateEditPermissions(ctx, orgID, updatedBy)
	if err != nil {
		s.logger.Error("Failed to validate edit permissions", zap.Error(err))
		return nil, rcerrors.NewInternalServerError("failed to validate permissions")
	}
	if !hasPermission {
		return nil, rcerrors.NewForbiddenError("user does not have permission to edit this organization")
	}

	// 3. Validar dirección fiscal si se actualiza
	if req.FiscalAddressID != nil {
		_, err := s.clients.Address.GetAddress(req.FiscalAddressID.String())
		if err != nil {
			s.logger.Error("Fiscal address validation failed", zap.Error(err))
			return nil, rcerrors.NewValidationError("fiscal_address_id", "invalid fiscal address")
		}
	}

	// 4. Preparar updates
	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.DisplayName != nil {
		updates["display_name"] = *req.DisplayName
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.LegalName != nil {
		updates["legal_name"] = *req.LegalName
	}
	if req.TaxID != nil {
		updates["tax_id"] = *req.TaxID
	}
	if req.Website != nil {
		updates["website"] = *req.Website
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.LogoURL != nil {
		updates["logo_url"] = *req.LogoURL
	}
	if req.TimezoneID != nil {
		updates["timezone_id"] = *req.TimezoneID
	}
	if req.FiscalAddressID != nil {
		updates["fiscal_address_id"] = *req.FiscalAddressID
	}

	updates["updated_by"] = updatedBy
	updates["updated_at"] = time.Now()

	// 6. Actualizar en base de datos
	updatedOrg, err := s.repos.Organization.Update(ctx, orgID, updates)
	if err != nil {
		s.logger.Error("Failed to update organization", zap.Error(err))
		return nil, rcerrors.NewInternalServerError("failed to update organization")
	}

	// 7. Invalidar cache
	s.invalidateOrganizationCache(orgID.String())

	// 8. Publicar evento
	if err := s.events.PublishOrganizationUpdated(updatedOrg); err != nil {
		s.logger.Warn("Failed to publish organization updated event", zap.Error(err))
	}

	s.logger.Info("Organization updated successfully", zap.String("org_id", orgID.String()))
	return updatedOrg, nil
}

// UpdateOrganizationStatus actualiza el estado de una organización
func (s *OrganizationService) UpdateOrganizationStatus(ctx context.Context, orgID uuid.UUID, status string, updatedBy uuid.UUID) error {
	s.logger.Info("Updating organization status",
		zap.String("org_id", orgID.String()),
		zap.String("status", status),
		zap.String("updated_by", updatedBy.String()),
	)

	// 1. Validar que la organización existe
	org, err := s.repos.Organization.GetByID(ctx, orgID)
	if err != nil {
		return rcerrors.NewNotFoundError("organization", orgID.String())
	}

	// 2. Validar permisos (solo admins del sistema pueden cambiar status)
	isSystemAdmin, err := s.validateSystemAdminPermissions(ctx, updatedBy)
	if err != nil {
		s.logger.Error("Failed to validate system admin permissions", zap.Error(err))
		return rcerrors.NewInternalServerError("failed to validate permissions")
	}
	if !isSystemAdmin {
		return rcerrors.NewForbiddenError("only system administrators can change organization status")
	}

	// 3. Validar que el nuevo status es válido
	validStatuses := []string{
		string(models.OrgStatusPending),
		string(models.OrgStatusActive),
		string(models.OrgStatusInactive),
		string(models.OrgStatusSuspended),
	}
	isValidStatus := false
	for _, validStatus := range validStatuses {
		if status == validStatus {
			isValidStatus = true
			break
		}
	}
	if !isValidStatus {
		return rcerrors.NewValidationError("status", "invalid organization status")
	}

	// 4. Actualizar status
	updates := map[string]interface{}{
		"status":     models.OrganizationStatus(status),
		"updated_by": updatedBy,
		"updated_at": time.Now(),
	}

	_, err = s.repos.Organization.Update(ctx, orgID, updates)
	if err != nil {
		s.logger.Error("Failed to update organization status", zap.Error(err))
		return rcerrors.NewInternalServerError("failed to update organization status")
	}

	// 5. Invalidar cache
	s.invalidateOrganizationCache(orgID.String())

	// 6. Publicar evento
	org.Status = models.OrganizationStatus(status)
	if err := s.events.PublishOrganizationStatusChanged(org, models.OrganizationStatus(status)); err != nil {
		s.logger.Warn("Failed to publish organization status changed event", zap.Error(err))
	}

	s.logger.Info("Organization status updated successfully",
		zap.String("org_id", orgID.String()),
		zap.String("new_status", status),
	)
	return nil
}

// GetOrganizationSettings obtiene las configuraciones de una organización
func (s *OrganizationService) GetOrganizationSettings(ctx context.Context, orgID uuid.UUID, requesterUserID uuid.UUID) (map[string]interface{}, error) {
	s.logger.Debug("Getting organization settings",
		zap.String("org_id", orgID.String()),
		zap.String("requester_user_id", requesterUserID.String()),
	)

	// 1. Verificar acceso
	hasAccess, err := s.validateUserAccess(ctx, orgID, requesterUserID)
	if err != nil {
		s.logger.Error("Failed to validate user access", zap.Error(err))
		return nil, rcerrors.NewInternalServerError("failed to validate access")
	}
	if !hasAccess {
		return nil, rcerrors.NewForbiddenError("user does not have access to this organization")
	}

	// 2. Obtener configuraciones
	settings, err := s.repos.OrganizationSetting.GetAllByOrganization(ctx, orgID)
	if err != nil {
		s.logger.Error("Failed to get organization settings", zap.Error(err))
		return nil, rcerrors.NewInternalServerError("failed to get organization settings")
	}

	// 3. Convertir a mapa
	settingsMap := make(map[string]interface{})
	for _, setting := range settings {
		settingsMap[setting.SettingKey] = setting.SettingValue
	}

	return settingsMap, nil
}

// UpdateOrganizationSettings actualiza las configuraciones de una organización
func (s *OrganizationService) UpdateOrganizationSettings(ctx context.Context, orgID uuid.UUID, settings map[string]interface{}, updatedBy uuid.UUID) error {
	s.logger.Info("Updating organization settings",
		zap.String("org_id", orgID.String()),
		zap.String("updated_by", updatedBy.String()),
		zap.Int("settings_count", len(settings)),
	)

	// 1. Verificar permisos de edición
	hasPermission, err := s.validateEditPermissions(ctx, orgID, updatedBy)
	if err != nil {
		s.logger.Error("Failed to validate edit permissions", zap.Error(err))
		return rcerrors.NewInternalServerError("failed to validate permissions")
	}
	if !hasPermission {
		return rcerrors.NewForbiddenError("user does not have permission to edit organization settings")
	}

	// 2. Actualizar configuraciones
	if err := s.repos.OrganizationSetting.SetMultiple(ctx, orgID, settings, updatedBy); err != nil {
		s.logger.Error("Failed to update settings", zap.Error(err))
		return rcerrors.NewInternalServerError("failed to update organization settings")
	}

	// 3. Invalidar cache
	s.invalidateOrganizationSettingsCache(orgID.String())

	s.logger.Info("Organization settings updated successfully", zap.String("org_id", orgID.String()))
	return nil
}

// ===============================
// HELPER METHODS
// ===============================

// createDefaultSettings crea configuraciones por defecto para una nueva organización
func (s *OrganizationService) createDefaultSettings(ctx context.Context, orgID uuid.UUID, createdBy uuid.UUID) error {
	defaultSettings := map[string]interface{}{
		"invite_expiration_days":     7,
		"max_pending_invitations":    10,
		"allow_public_registration":  false,
		"require_email_verification": true,
		"default_employee_role":      "employee",
		"timezone":                   "UTC",
		"date_format":                "YYYY-MM-DD",
		"currency":                   "USD",
		"language":                   "en",
	}

	if err := s.repos.OrganizationSetting.SetMultiple(ctx, orgID, defaultSettings, createdBy); err != nil {
		s.logger.Error("Failed to create default settings", zap.Error(err))
		return err
	}

	return nil
}

// createDefaultOwnerRole crea el rol de propietario por defecto
func (s *OrganizationService) createDefaultOwnerRole(ctx context.Context, orgID uuid.UUID, createdBy uuid.UUID) error {
	ownerRole := &models.OrganizationRole{
		OrganizationID: orgID,
		Name:           "Owner",
		Description:    func() *string { desc := "Organization owner with full permissions"; return &desc }(),
		IsDefault:      true,
		CreatedBy:      createdBy,
	}

	_, err := s.repos.Role.Create(ctx, ownerRole)
	if err != nil {
		s.logger.Error("Failed to create default owner role", zap.Error(err))
		return err
	}

	return nil
}

// validateUserAccess verifica que un usuario tiene acceso a una organización
func (s *OrganizationService) validateUserAccess(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) (bool, error) {
	// Verificar si el usuario es empleado de la organización
	isEmployee, err := s.repos.Employee.ExistsByUserAndOrg(ctx, userID, orgID)
	if err != nil {
		return false, err
	}
	if isEmployee {
		return true, nil
	}

	// Verificar si el usuario es propietario de la organización
	isOwner, err := s.repos.Owner.ExistsByPersonAndOrg(ctx, userID, orgID)
	if err != nil {
		return false, err
	}
	if isOwner {
		return true, nil
	}

	return false, nil
}

// validateEditPermissions verifica que un usuario puede editar una organización
func (s *OrganizationService) validateEditPermissions(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) (bool, error) {
	// Solo propietarios pueden editar la organización
	isOwner, err := s.repos.Owner.ExistsByPersonAndOrg(ctx, userID, orgID)
	if err != nil {
		return false, err
	}

	return isOwner, nil
}

// validateSystemAdminPermissions verifica que un usuario es administrador del sistema
func (s *OrganizationService) validateSystemAdminPermissions(ctx context.Context, userID uuid.UUID) (bool, error) {
	// Verificar con el servicio de autenticación si el usuario es admin del sistema
	userDetails, err := s.clients.AuthIdentity.GetUserById(userID.String())
	if err != nil {
		return false, err
	}

	// Asumir que hay un campo que indica si es admin del sistema
	// Esta lógica dependerá de la implementación del servicio de auth
	// TODO: Implementar verificación de admin de sistema
	// Por ahora, solo verificamos si el usuario existe
	return userDetails != nil, nil
}

// getOrganizationFromCache obtiene una organización del cache
func (s *OrganizationService) getOrganizationFromCache(orgID string) *models.Organization {
	if org, err := s.cache.GetOrganization(orgID); err == nil && org != nil {
		return org
	}

	return nil
}

// setOrganizationCache guarda una organización en cache
func (s *OrganizationService) setOrganizationCache(orgID string, org *models.Organization) {
	ttl := 30 * time.Minute
	s.cache.SetOrganization(org, ttl)
}

// invalidateOrganizationCache invalida el cache de una organización
func (s *OrganizationService) invalidateOrganizationCache(orgID string) {
	s.cache.DeleteOrganization(orgID)
	s.cache.DeleteOrganizationSettings(orgID)
}

// invalidateOrganizationSettingsCache invalida el cache de configuraciones de una organización
func (s *OrganizationService) invalidateOrganizationSettingsCache(orgID string) {
	s.cache.DeleteOrganizationSettings(orgID)
}

// createMainBranch crea la sucursal principal por defecto para una nueva organización
func (s *OrganizationService) createMainBranch(ctx context.Context, org *models.Organization, createdBy uuid.UUID) error {
	s.logger.Info("Creating main branch for organization",
		zap.String("org_id", org.ID.String()),
		zap.String("org_name", org.DisplayName),
	)

	// Crear la branch principal con el mismo nombre que la organización
	mainBranch := &models.OrganizationBranch{
		OrganizationID: org.ID,
		DisplayName:    org.DisplayName + " - Oficina Principal", // Agregar sufijo para distinguir
		IsMain:         true,
		CreatedBy:      createdBy,
	}

	// Usar el mismo email y teléfono de la organización si están disponibles
	if org.Email != nil {
		mainBranch.Email = org.Email
	}
	if org.Phone != nil {
		mainBranch.Phone = org.Phone
	}

	// Crear la branch en la base de datos
	_, err := s.repos.Branch.Create(ctx, mainBranch)
	if err != nil {
		s.logger.Error("Failed to create main branch", zap.Error(err))
		return err
	}

	s.logger.Info("Main branch created successfully",
		zap.String("org_id", org.ID.String()),
		zap.String("branch_name", mainBranch.DisplayName),
	)

	return nil
}
