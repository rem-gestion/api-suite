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
	"github.com/rem-gestion/rem-common/errors"
)

// EmployeeService maneja la lógica de negocio para empleados
type EmployeeService struct {
	repos   *repository.Repositories
	clients *clients.ClientManager
	cache   CacheService
	events  EventService
	logger  *zap.Logger
}

// NewEmployeeService crea una nueva instancia del servicio de empleados
func NewEmployeeService(
	repos *repository.Repositories,
	clients *clients.ClientManager,
	cache CacheService,
	events EventService,
	logger *zap.Logger,
) *EmployeeService {
	return &EmployeeService{
		repos:   repos,
		clients: clients,
		cache:   cache,
		events:  events,
		logger:  logger,
	}
}

// AddEmployee vincula un usuario a una organización como empleado
func (s *EmployeeService) AddEmployee(ctx context.Context, orgID string, req *dto.AddEmployeeRequest, createdBy uuid.UUID) (*models.Employee, error) {
	s.logger.Info("Adding employee to organization",
		zap.String("org_id", orgID),
		zap.String("user_id", req.UserID.String()),
		zap.String("created_by", createdBy.String()),
	)

	// Convertir orgID a UUID
	orgUUID, err := uuid.Parse(orgID)
	if err != nil {
		s.logger.Error("Invalid organization ID format", zap.Error(err))
		return nil, errors.NewValidationError("invalid organization ID format", err.Error())
	}

	// Verificar que la organización existe
	org, err := s.repos.Organization.GetByID(ctx, orgUUID)
	if err != nil {
		s.logger.Error("Organization not found", zap.Error(err))
		return nil, errors.NewNotFoundError("organization", orgID)
	}

	// Validar que el usuario creador existe
	_, err = s.clients.AuthIdentity.GetUserById(createdBy.String())
	if err != nil {
		s.logger.Error("Creator user validation failed", zap.Error(err))
		return nil, errors.NewValidationError("invalid creator user_id", err.Error())
	}

	// Validar que el usuario empleado existe
	_, err = s.clients.AuthIdentity.GetUserById(req.UserID.String())
	if err != nil {
		s.logger.Error("Employee user validation failed", zap.Error(err))
		return nil, errors.NewValidationError("invalid employee user_id", err.Error())
	}

	// Verificar que el usuario no es ya empleado de esta organización
	exists, err := s.repos.Employee.ExistsByUserAndOrg(ctx, req.UserID, orgUUID)
	if err != nil {
		s.logger.Error("Failed to check employee existence", zap.Error(err))
		return nil, errors.NewDatabaseError("failed to check employee existence", err)
	}
	if exists {
		return nil, errors.NewValidationError("employee already exists", "user is already an employee of this organization")
	}

	// Validar que la persona existe si se proporciona
	if req.PersonID != nil {
		_, err = s.clients.Person.GetPerson(req.PersonID.String())
		if err != nil {
			s.logger.Error("Person validation failed", zap.Error(err))
			return nil, errors.NewValidationError("invalid person_id", err.Error())
		}
	}

	// Validar que el rol primario existe si se proporciona
	if req.PrimaryRoleID != nil {
		role, err := s.repos.Role.GetByID(ctx, *req.PrimaryRoleID)
		if err != nil {
			s.logger.Error("Primary role validation failed", zap.Error(err))
			return nil, errors.NewValidationError("invalid primary_role_id", err.Error())
		}
		// Verificar que el rol pertenece a la organización
		if role.OrganizationID != orgUUID {
			return nil, errors.NewValidationError("role does not belong to organization", "primary role must belong to the same organization")
		}
	}

	// Crear empleado
	employee := &models.Employee{
		OrganizationID: orgUUID,
		UserID:         req.UserID,
		PersonID:       req.PersonID,
		PrimaryRoleID:  req.PrimaryRoleID,
		Status:         models.EmpStatusActive,
		HiredAt:        req.HiredAt,
		CreatedBy:      createdBy,
		UpdatedBy:      &createdBy,
	}

	// Si no se proporciona fecha de contratación, usar fecha actual
	if employee.HiredAt == nil {
		now := time.Now()
		employee.HiredAt = &now
	}

	// Guardar empleado en base de datos
	createdEmployee, err := s.repos.Employee.Create(ctx, employee)
	if err != nil {
		s.logger.Error("Failed to create employee", zap.Error(err))
		return nil, errors.NewDatabaseError("failed to create employee", err)
	}

	// Asignar roles adicionales si se proporcionan
	if len(req.RoleIDs) > 0 {
		for _, roleID := range req.RoleIDs {
			isPrimary := req.PrimaryRoleID != nil && roleID == *req.PrimaryRoleID

			_, err := s.repos.EmployeeRole.AssignRole(ctx, createdEmployee.ID, roleID, isPrimary, createdBy)
			if err != nil {
				s.logger.Warn("Failed to assign additional role",
					zap.String("employee_id", createdEmployee.ID.String()),
					zap.String("role_id", roleID.String()),
					zap.Error(err),
				)
			}
		}
	} else if req.PrimaryRoleID != nil {
		// Asignar solo el rol primario
		_, err := s.repos.EmployeeRole.AssignRole(ctx, createdEmployee.ID, *req.PrimaryRoleID, true, createdBy)
		if err != nil {
			s.logger.Warn("Failed to assign primary role",
				zap.String("employee_id", createdEmployee.ID.String()),
				zap.String("role_id", req.PrimaryRoleID.String()),
				zap.Error(err),
			)
		}
	}

	// Invalidar cache relacionado
	s.invalidateEmployeeCache(createdEmployee.ID.String())
	s.invalidateOrganizationEmployeesCache(orgID)

	// Publicar evento
	if err := s.events.PublishEmployeeAdded(org, createdEmployee); err != nil {
		s.logger.Warn("Failed to publish employee added event",
			zap.String("employee_id", createdEmployee.ID.String()),
			zap.Error(err),
		)
	}

	s.logger.Info("Employee added successfully",
		zap.String("employee_id", createdEmployee.ID.String()),
		zap.String("org_id", orgID),
	)

	return createdEmployee, nil
}

// GetEmployee obtiene un empleado por ID
func (s *EmployeeService) GetEmployee(ctx context.Context, empID string) (*models.Employee, error) {
	s.logger.Debug("Getting employee", zap.String("employee_id", empID))

	// Convertir string a UUID
	empUUID, err := uuid.Parse(empID)
	if err != nil {
		s.logger.Error("Invalid employee ID format", zap.Error(err))
		return nil, errors.NewValidationError("invalid employee ID format", err.Error())
	}

	// Intentar obtener del cache primero
	if cached, err := s.cache.GetEmployee(empID); err == nil && cached != nil {
		s.logger.Debug("Employee found in cache", zap.String("employee_id", empID))
		return cached, nil
	}

	// Obtener de base de datos con relaciones
	employee, err := s.repos.Employee.GetWithRoles(ctx, empUUID)
	if err != nil {
		s.logger.Error("Failed to get employee", zap.Error(err))
		return nil, errors.NewNotFoundError("employee", empID)
	}

	// Guardar en cache
	if err := s.cache.SetEmployee(employee, 15*time.Minute); err != nil {
		s.logger.Warn("Failed to cache employee",
			zap.String("employee_id", empID),
			zap.Error(err),
		)
	}

	return employee, nil
}

// ListEmployees lista empleados con filtros y paginación
func (s *EmployeeService) ListEmployees(ctx context.Context, filters *dto.EmployeeFiltersRequest) ([]*models.Employee, int64, error) {
	s.logger.Debug("Listing employees", zap.Any("filters", filters))

	// Validar filtros
	if err := s.validateEmployeeFilters(filters); err != nil {
		return nil, 0, err
	}

	employees, total, err := s.repos.Employee.List(ctx, *filters)
	if err != nil {
		s.logger.Error("Failed to list employees", zap.Error(err))
		return nil, 0, errors.NewDatabaseError("failed to list employees", err)
	}

	// Convertir []models.Employee a []*models.Employee
	result := make([]*models.Employee, len(employees))
	for i := range employees {
		result[i] = &employees[i]
	}

	s.logger.Debug("Employees listed successfully",
		zap.Int("count", len(result)),
		zap.Int64("total", total),
	)

	return result, total, nil
}

// UpdateEmployee actualiza información de un empleado
func (s *EmployeeService) UpdateEmployee(ctx context.Context, empID string, req *dto.UpdateEmployeeRequest, updatedBy uuid.UUID) (*models.Employee, error) {
	s.logger.Info("Updating employee",
		zap.String("employee_id", empID),
		zap.String("updated_by", updatedBy.String()),
	)

	// Convertir string a UUID
	empUUID, err := uuid.Parse(empID)
	if err != nil {
		s.logger.Error("Invalid employee ID format", zap.Error(err))
		return nil, errors.NewValidationError("invalid employee ID format", err.Error())
	}

	// Verificar que el empleado existe
	employee, err := s.repos.Employee.GetByID(ctx, empUUID)
	if err != nil {
		s.logger.Error("Employee not found for update", zap.Error(err))
		return nil, errors.NewNotFoundError("employee", empID)
	}

	// Validar que el usuario que actualiza existe
	_, err = s.clients.AuthIdentity.GetUserById(updatedBy.String())
	if err != nil {
		s.logger.Error("Updater user validation failed", zap.Error(err))
		return nil, errors.NewValidationError("invalid updater user_id", err.Error())
	}

	// Validar que la nueva persona existe si se proporciona
	if req.PersonID != nil {
		_, err = s.clients.Person.GetPerson(req.PersonID.String())
		if err != nil {
			s.logger.Error("New person validation failed", zap.Error(err))
			return nil, errors.NewValidationError("invalid person_id", err.Error())
		}
	}

	// Validar que el nuevo rol primario existe si se proporciona
	if req.PrimaryRoleID != nil {
		role, err := s.repos.Role.GetByID(ctx, *req.PrimaryRoleID)
		if err != nil {
			s.logger.Error("New primary role validation failed", zap.Error(err))
			return nil, errors.NewValidationError("invalid primary_role_id", err.Error())
		}
		// Verificar que el rol pertenece a la organización
		if role.OrganizationID != employee.OrganizationID {
			return nil, errors.NewValidationError("role does not belong to organization", "primary role must belong to the same organization")
		}
	}

	// Preparar actualizaciones
	updates := make(map[string]interface{})

	if req.PersonID != nil {
		updates["person_id"] = req.PersonID
	}
	if req.PrimaryRoleID != nil {
		updates["primary_role_id"] = req.PrimaryRoleID
	}
	if req.HiredAt != nil {
		updates["hired_at"] = req.HiredAt
	}
	if req.FiredAt != nil {
		updates["fired_at"] = req.FiredAt
	}

	updates["updated_by"] = updatedBy
	updates["updated_at"] = time.Now()

	// Actualizar en base de datos
	updatedEmployee, err := s.repos.Employee.Update(ctx, empUUID, updates)
	if err != nil {
		s.logger.Error("Failed to update employee", zap.Error(err))
		return nil, errors.NewDatabaseError("failed to update employee", err)
	}

	// Si se cambió el rol primario, actualizar la relación empleado-rol
	if req.PrimaryRoleID != nil {
		if err := s.repos.EmployeeRole.SetPrimaryRole(ctx, empUUID, *req.PrimaryRoleID); err != nil {
			s.logger.Warn("Failed to update primary role relationship",
				zap.String("employee_id", empID),
				zap.String("role_id", req.PrimaryRoleID.String()),
				zap.Error(err),
			)
		}
	}

	// Invalidar cache
	s.invalidateEmployeeCache(empID)

	// Publicar evento
	if err := s.events.PublishEmployeeUpdated(&employee.Organization, updatedEmployee); err != nil {
		s.logger.Warn("Failed to publish employee updated event",
			zap.String("employee_id", empID),
			zap.Error(err),
		)
	}

	s.logger.Info("Employee updated successfully", zap.String("employee_id", empID))
	return updatedEmployee, nil
}

// UpdateEmployeeStatus actualiza el estado de un empleado
func (s *EmployeeService) UpdateEmployeeStatus(ctx context.Context, empID string, status models.EmployeeStatus, updatedBy uuid.UUID) error {
	s.logger.Info("Updating employee status",
		zap.String("employee_id", empID),
		zap.String("status", string(status)),
		zap.String("updated_by", updatedBy.String()),
	)

	// Convertir string a UUID
	empUUID, err := uuid.Parse(empID)
	if err != nil {
		s.logger.Error("Invalid employee ID format", zap.Error(err))
		return errors.NewValidationError("invalid employee ID format", err.Error())
	}

	// Verificar que el empleado existe
	employee, err := s.repos.Employee.GetByID(ctx, empUUID)
	if err != nil {
		s.logger.Error("Employee not found for status update", zap.Error(err))
		return errors.NewNotFoundError("employee", empID)
	}

	// Validar que el usuario que actualiza existe
	_, err = s.clients.AuthIdentity.GetUserById(updatedBy.String())
	if err != nil {
		s.logger.Error("Updater user validation failed", zap.Error(err))
		return errors.NewValidationError("invalid updater user_id", err.Error())
	}

	// Validar transición de estado
	if err := s.validateStatusTransition(employee.Status, status); err != nil {
		return err
	}

	// Verificar regla de negocio: no se puede terminar al último admin
	if status == models.EmpStatusTerminated {
		if err := s.checkLastAdminRule(ctx, empUUID, employee.OrganizationID); err != nil {
			return err
		}
	}

	oldStatus := employee.Status

	// Actualizar estado
	updates := map[string]interface{}{
		"status":     status,
		"updated_by": updatedBy,
		"updated_at": time.Now(),
	}

	// Si se está terminando el empleado, establecer fecha de despido
	if status == models.EmpStatusTerminated {
		now := time.Now()
		updates["fired_at"] = &now
	}

	if err := s.repos.Employee.UpdateStatus(ctx, empUUID, status, updatedBy); err != nil {
		s.logger.Error("Failed to update employee status", zap.Error(err))
		return errors.NewDatabaseError("failed to update employee status", err)
	}

	// Invalidar cache
	s.invalidateEmployeeCache(empID)

	// Publicar evento
	if err := s.events.PublishEmployeeStatusChanged(&employee.Organization, employee, status); err != nil {
		s.logger.Warn("Failed to publish employee status changed event",
			zap.String("employee_id", empID),
			zap.Error(err),
		)
	}

	s.logger.Info("Employee status updated successfully",
		zap.String("employee_id", empID),
		zap.String("old_status", string(oldStatus)),
		zap.String("new_status", string(status)),
	)

	return nil
}

// AssignRole asigna un rol a un empleado
func (s *EmployeeService) AssignRole(ctx context.Context, empID string, req *dto.AssignRoleRequest, createdBy uuid.UUID) error {
	s.logger.Info("Assigning role to employee",
		zap.String("employee_id", empID),
		zap.String("role_id", req.RoleID.String()),
		zap.Bool("is_primary", req.IsPrimary),
		zap.String("created_by", createdBy.String()),
	)

	// Convertir string a UUID
	empUUID, err := uuid.Parse(empID)
	if err != nil {
		s.logger.Error("Invalid employee ID format", zap.Error(err))
		return errors.NewValidationError("invalid employee ID format", err.Error())
	}

	// Verificar que el empleado existe
	employee, err := s.repos.Employee.GetByID(ctx, empUUID)
	if err != nil {
		s.logger.Error("Employee not found for role assignment", zap.Error(err))
		return errors.NewNotFoundError("employee", empID)
	}

	// Verificar que el rol existe y pertenece a la organización
	role, err := s.repos.Role.GetByID(ctx, req.RoleID)
	if err != nil {
		s.logger.Error("Role not found for assignment", zap.Error(err))
		return errors.NewNotFoundError("role", req.RoleID.String())
	}

	if role.OrganizationID != employee.OrganizationID {
		return errors.NewValidationError("role does not belong to organization", "role must belong to the same organization as the employee")
	}

	// Verificar si ya tiene el rol
	hasRole, err := s.repos.EmployeeRole.HasRole(ctx, empUUID, req.RoleID)
	if err != nil {
		s.logger.Error("Failed to check existing role", zap.Error(err))
		return errors.NewDatabaseError("failed to check existing role", err)
	}
	if hasRole {
		return errors.NewValidationError("role already assigned", "employee already has this role")
	}

	// Asignar rol
	empRole, err := s.repos.EmployeeRole.AssignRole(ctx, empUUID, req.RoleID, req.IsPrimary, createdBy)
	if err != nil {
		s.logger.Error("Failed to assign role", zap.Error(err))
		return errors.NewDatabaseError("failed to assign role", err)
	}

	// Si es rol primario, actualizar el employee
	if req.IsPrimary {
		updates := map[string]interface{}{
			"primary_role_id": req.RoleID,
			"updated_by":      createdBy,
			"updated_at":      time.Now(),
		}
		if _, err := s.repos.Employee.Update(ctx, empUUID, updates); err != nil {
			s.logger.Warn("Failed to update employee primary role",
				zap.String("employee_id", empID),
				zap.Error(err),
			)
		}
	}

	// Invalidar cache
	s.invalidateEmployeeCache(empID)

	// Publicar evento
	if err := s.events.PublishEmployeeRoleAssigned(&employee.Organization, employee, role); err != nil {
		s.logger.Warn("Failed to publish employee role assigned event",
			zap.String("employee_id", empID),
			zap.String("role_id", req.RoleID.String()),
			zap.Error(err),
		)
	}

	s.logger.Info("Role assigned successfully",
		zap.String("employee_id", empID),
		zap.String("role_id", req.RoleID.String()),
		zap.String("emp_role_id", empRole.ID.String()),
	)

	return nil
}

// RemoveRole remueve un rol de un empleado
func (s *EmployeeService) RemoveRole(ctx context.Context, empID, roleID string) error {
	s.logger.Info("Removing role from employee",
		zap.String("employee_id", empID),
		zap.String("role_id", roleID),
	)

	// Convertir strings a UUIDs
	empUUID, err := uuid.Parse(empID)
	if err != nil {
		s.logger.Error("Invalid employee ID format", zap.Error(err))
		return errors.NewValidationError("invalid employee ID format", err.Error())
	}

	roleUUID, err := uuid.Parse(roleID)
	if err != nil {
		s.logger.Error("Invalid role ID format", zap.Error(err))
		return errors.NewValidationError("invalid role ID format", err.Error())
	}

	// Verificar que el empleado existe
	employee, err := s.repos.Employee.GetByID(ctx, empUUID)
	if err != nil {
		s.logger.Error("Employee not found for role removal", zap.Error(err))
		return errors.NewNotFoundError("employee", empID)
	}

	// Verificar que tiene el rol
	hasRole, err := s.repos.EmployeeRole.HasRole(ctx, empUUID, roleUUID)
	if err != nil {
		s.logger.Error("Failed to check existing role", zap.Error(err))
		return errors.NewDatabaseError("failed to check existing role", err)
	}
	if !hasRole {
		return errors.NewValidationError("role not assigned", "employee does not have this role")
	}

	// Verificar si es el rol primario
	if employee.PrimaryRoleID != nil && *employee.PrimaryRoleID == roleUUID {
		// Limpiar rol primario del empleado
		updates := map[string]interface{}{
			"primary_role_id": nil,
			"updated_at":      time.Now(),
		}
		if _, err := s.repos.Employee.Update(ctx, empUUID, updates); err != nil {
			s.logger.Warn("Failed to clear employee primary role",
				zap.String("employee_id", empID),
				zap.Error(err),
			)
		}
	}

	// Remover rol
	if err := s.repos.EmployeeRole.RemoveRole(ctx, empUUID, roleUUID); err != nil {
		s.logger.Error("Failed to remove role", zap.Error(err))
		return errors.NewDatabaseError("failed to remove role", err)
	}

	// Invalidar cache
	s.invalidateEmployeeCache(empID)

	// Publicar evento
	if err := s.events.PublishEmployeeRoleRemoved(&employee.Organization, employee, roleID); err != nil {
		s.logger.Warn("Failed to publish employee role removed event",
			zap.String("employee_id", empID),
			zap.String("role_id", roleID),
			zap.Error(err),
		)
	}

	s.logger.Info("Role removed successfully",
		zap.String("employee_id", empID),
		zap.String("role_id", roleID),
	)

	return nil
}

// GetEmployeeRoles obtiene todos los roles de un empleado
func (s *EmployeeService) GetEmployeeRoles(ctx context.Context, empID string) ([]*models.EmployeeRole, error) {
	s.logger.Debug("Getting employee roles", zap.String("employee_id", empID))

	// Convertir string a UUID
	empUUID, err := uuid.Parse(empID)
	if err != nil {
		s.logger.Error("Invalid employee ID format", zap.Error(err))
		return nil, errors.NewValidationError("invalid employee ID format", err.Error())
	}

	// Verificar que el empleado existe
	_, err = s.repos.Employee.GetByID(ctx, empUUID)
	if err != nil {
		s.logger.Error("Employee not found", zap.Error(err))
		return nil, errors.NewNotFoundError("employee", empID)
	}

	// Obtener roles del empleado
	empRoles, err := s.repos.EmployeeRole.GetByEmployee(ctx, empUUID)
	if err != nil {
		s.logger.Error("Failed to get employee roles", zap.Error(err))
		return nil, errors.NewDatabaseError("failed to get employee roles", err)
	}

	// Convertir []models.EmployeeRole a []*models.EmployeeRole
	result := make([]*models.EmployeeRole, len(empRoles))
	for i := range empRoles {
		result[i] = &empRoles[i]
	}

	s.logger.Debug("Employee roles retrieved successfully",
		zap.String("employee_id", empID),
		zap.Int("roles_count", len(result)),
	)

	return result, nil
}

// validateEmployeeFilters valida los filtros de listado de empleados
func (s *EmployeeService) validateEmployeeFilters(filters *dto.EmployeeFiltersRequest) error {
	if filters.Page < 1 {
		return errors.NewValidationError("invalid page", "page must be greater than 0")
	}
	if filters.PerPage < 1 || filters.PerPage > 100 {
		return errors.NewValidationError("invalid page_size", "page_size must be between 1 and 100")
	}
	return nil
}

// validateStatusTransition valida que la transición de estado sea válida
func (s *EmployeeService) validateStatusTransition(oldStatus, newStatus models.EmployeeStatus) error {
	// Definir transiciones válidas
	validTransitions := map[models.EmployeeStatus][]models.EmployeeStatus{
		models.EmpStatusActive: {
			models.EmpStatusSuspended,
			models.EmpStatusInactive,
			models.EmpStatusTerminated,
		},
		models.EmpStatusSuspended: {
			models.EmpStatusActive,
			models.EmpStatusTerminated,
		},
		models.EmpStatusInactive: {
			models.EmpStatusActive,
			models.EmpStatusTerminated,
		},
		models.EmpStatusTerminated: {
			// No hay transiciones desde terminado (es estado final)
		},
	}

	allowedStatuses, exists := validTransitions[oldStatus]
	if !exists {
		return errors.NewValidationError("invalid_status_transition",
			fmt.Sprintf("no transitions allowed from status %s", oldStatus))
	}

	for _, allowed := range allowedStatuses {
		if allowed == newStatus {
			return nil
		}
	}

	return errors.NewValidationError("invalid_status_transition",
		fmt.Sprintf("transition from %s to %s is not allowed", oldStatus, newStatus))
}

// checkLastAdminRule verifica que no se esté terminando al último admin
func (s *EmployeeService) checkLastAdminRule(ctx context.Context, empID, orgID uuid.UUID) error {
	// Obtener roles de admin de la organización
	adminRoleFilters := &dto.RoleFiltersRequest{
		FilterRequest: dto.FilterRequest{
			Search: "admin",
		},
		OrganizationID: orgID,
		PaginationRequest: dto.PaginationRequest{
			Page:    1,
			PerPage: 50,
		},
	}
	adminRoles, _, err := s.repos.Role.ListByOrganization(ctx, orgID, *adminRoleFilters)
	if err != nil {
		s.logger.Error("Failed to get admin roles", zap.Error(err))
		return errors.NewDatabaseError("failed to check admin roles", err)
	}

	if len(adminRoles) == 0 {
		// No hay roles de admin definidos, permitir la operación
		return nil
	}

	// Contar empleados activos con roles de admin (excluyendo el empleado actual)
	activeAdmins := 0
	for _, role := range adminRoles {
		count, err := s.repos.EmployeeRole.CountByRole(ctx, role.ID)
		if err != nil {
			continue
		}

		// Si es el empleado actual y tiene este rol, no contarlo
		hasRole, err := s.repos.EmployeeRole.HasRole(ctx, empID, role.ID)
		if err == nil && hasRole {
			count--
		}

		activeAdmins += int(count)
	}

	if activeAdmins <= 0 {
		return errors.NewBusinessRuleError("cannot terminate last admin", "organization must have at least one active administrator")
	}

	return nil
}

// invalidateEmployeeCache limpia el cache relacionado con un empleado
func (s *EmployeeService) invalidateEmployeeCache(empID string) {
	if err := s.cache.DeleteEmployee(empID); err != nil {
		s.logger.Warn("Failed to invalidate employee cache",
			zap.String("employee_id", empID),
			zap.Error(err),
		)
	}
}

// invalidateOrganizationEmployeesCache limpia el cache de empleados de una organización
func (s *EmployeeService) invalidateOrganizationEmployeesCache(orgID string) {
	if err := s.cache.InvalidateOrganizationCache(orgID); err != nil {
		s.logger.Warn("Failed to invalidate organization employees cache",
			zap.String("org_id", orgID),
			zap.Error(err),
		)
	}
}
