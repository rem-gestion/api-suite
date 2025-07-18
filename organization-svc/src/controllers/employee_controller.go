package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/rem-gestion/api-suite/organization/src/dto"
	"github.com/rem-gestion/api-suite/organization/src/models"
	"github.com/rem-gestion/api-suite/organization/src/services"
	rcerrors "github.com/rem-gestion/rem-common/errors"
)

// EmployeeController maneja las peticiones HTTP para empleados
type EmployeeController struct {
	empService *services.EmployeeService
	logger     *zap.Logger
}

// NewEmployeeController crea una nueva instancia del controlador de empleados
func NewEmployeeController(empService *services.EmployeeService, logger *zap.Logger) *EmployeeController {
	return &EmployeeController{
		empService: empService,
		logger:     logger,
	}
}

// AddEmployee crea un nuevo empleado en la organización
// @Summary Crear empleado
// @Description Agrega un nuevo empleado a la organización
// @Tags employees
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param employee body dto.AddEmployeeRequest true "Employee data"
// @Success 201 {object} dto.EmployeeResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /organizations/{orgId}/employees [post]
func (c *EmployeeController) AddEmployee(ctx *gin.Context) {
	orgID := ctx.Param("orgId")
	if _, err := uuid.Parse(orgID); err != nil {
		c.logger.Warn("Invalid organization ID", zap.String("orgId", orgID), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "INVALID_ORG_ID",
			Message: "Invalid organization ID format",
		})
		return
	}

	var req dto.AddEmployeeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Warn("Invalid request body", zap.Error(err))

		// Check if it's a validation error
		if validationErr, ok := err.(*rcerrors.ValidationError); ok {
			ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "VALIDATION_FAILED",
				Message: "Request validation failed",
				Details: map[string]interface{}{
					"fields": validationErr.Fields,
				},
			})
			return
		}

		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "INVALID_REQUEST_BODY",
			Message: "Invalid request body format",
		})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("user_id")
	if !exists {
		c.logger.Error("User ID not found in context")
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "UNAUTHORIZED",
			Message: "User not authenticated",
		})
		return
	}

	employee, err := c.empService.AddEmployee(ctx, orgID, &req, uuid.MustParse(userID.(string)))
	if err != nil {
		c.logger.Error("Failed to add employee",
			zap.String("orgId", orgID),
			zap.String("userId", userID.(string)),
			zap.Error(err))

		// Handle different types of errors
		switch err.(type) {
		case *rcerrors.ValidationError:
			validationErr := err.(*rcerrors.ValidationError)
			ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "VALIDATION_FAILED",
				Message: validationErr.Error(),
				Details: map[string]interface{}{
					"fields": validationErr.Fields,
				},
			})
		case *rcerrors.NotFoundError:
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "NOT_FOUND",
				Message: err.Error(),
			})
		case *rcerrors.ForbiddenError:
			ctx.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "FORBIDDEN",
				Message: err.Error(),
			})
		default:
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "INTERNAL_ERROR",
				Message: "Failed to add employee",
			})
		}
		return
	}

	response := mapEmployeeToResponse(employee)
	ctx.JSON(http.StatusCreated, response)
}

// GetEmployee obtiene un empleado por ID
// @Summary Obtener empleado
// @Description Obtiene información de un empleado específico
// @Tags employees
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param empId path string true "Employee ID"
// @Success 200 {object} dto.EmployeeResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /organizations/{orgId}/employees/{empId} [get]
func (c *EmployeeController) GetEmployee(ctx *gin.Context) {
	orgID := ctx.Param("orgId")
	empID := ctx.Param("empId")

	if _, err := uuid.Parse(orgID); err != nil {
		c.logger.Warn("Invalid organization ID", zap.String("orgId", orgID), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "INVALID_ORG_ID",
			Message: "Invalid organization ID format",
		})
		return
	}

	if _, err := uuid.Parse(empID); err != nil {
		c.logger.Warn("Invalid employee ID", zap.String("empId", empID), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "INVALID_EMP_ID",
			Message: "Invalid employee ID format",
		})
		return
	}

	employee, err := c.empService.GetEmployee(ctx, empID)
	if err != nil {
		c.logger.Error("Failed to get employee", zap.String("empId", empID), zap.Error(err))

		switch err.(type) {
		case *rcerrors.NotFoundError:
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "EMPLOYEE_NOT_FOUND",
				Message: "Employee not found",
			})
		case *rcerrors.ForbiddenError:
			ctx.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "FORBIDDEN",
				Message: err.Error(),
			})
		default:
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "INTERNAL_ERROR",
				Message: "Failed to retrieve employee",
			})
		}
		return
	}

	response := mapEmployeeToResponse(employee)
	ctx.JSON(http.StatusOK, response)
}

// ListEmployees obtiene la lista de empleados de una organización
// @Summary Listar empleados
// @Description Obtiene la lista de empleados de una organización con filtros opcionales
// @Tags employees
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param status query string false "Employee status filter"
// @Param role_id query string false "Role ID filter"
// @Param search query string false "Search term"
// @Success 200 {object} dto.EmployeeListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /organizations/{orgId}/employees [get]
func (c *EmployeeController) ListEmployees(ctx *gin.Context) {
	orgID := ctx.Param("orgId")
	if _, err := uuid.Parse(orgID); err != nil {
		c.logger.Warn("Invalid organization ID", zap.String("orgId", orgID), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "INVALID_ORG_ID",
			Message: "Invalid organization ID format",
		})
		return
	}

	var filters dto.EmployeeFiltersRequest
	if err := ctx.ShouldBindQuery(&filters); err != nil {
		c.logger.Warn("Invalid query parameters", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "INVALID_QUERY_PARAMS",
			Message: "Invalid query parameters",
		})
		return
	}

	// Set default pagination
	if filters.Page <= 0 {
		filters.Page = 1
	}
	if filters.PerPage <= 0 {
		filters.PerPage = 10
	}

	employees, total, err := c.empService.ListEmployees(ctx, &filters)
	if err != nil {
		c.logger.Error("Failed to list employees", zap.String("orgId", orgID), zap.Error(err))

		switch err.(type) {
		case *rcerrors.NotFoundError:
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "ORGANIZATION_NOT_FOUND",
				Message: "Organization not found",
			})
		case *rcerrors.ForbiddenError:
			ctx.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "FORBIDDEN",
				Message: err.Error(),
			})
		default:
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "INTERNAL_ERROR",
				Message: "Failed to retrieve employees",
			})
		}
		return
	}

	// Map employees to responses
	employeeResponses := make([]dto.EmployeeResponse, len(employees))
	for i, emp := range employees {
		employeeResponses[i] = *mapEmployeeToResponse(emp)
	}

	response := dto.EmployeeListResponse{
		Data: employeeResponses,
		Pagination: dto.PaginationResponse{
			Page:       filters.Page,
			PerPage:    filters.PerPage,
			Total:      total,
			TotalPages: int((total + int64(filters.PerPage) - 1) / int64(filters.PerPage)),
		},
	}

	ctx.JSON(http.StatusOK, response)
}

// UpdateEmployee actualiza la información de un empleado
// @Summary Actualizar empleado
// @Description Actualiza la información de un empleado
// @Tags employees
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param empId path string true "Employee ID"
// @Param employee body dto.UpdateEmployeeRequest true "Employee update data"
// @Success 200 {object} dto.EmployeeResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /organizations/{orgId}/employees/{empId} [put]
func (c *EmployeeController) UpdateEmployee(ctx *gin.Context) {
	orgID := ctx.Param("orgId")
	empID := ctx.Param("empId")

	if _, err := uuid.Parse(orgID); err != nil {
		c.logger.Warn("Invalid organization ID", zap.String("orgId", orgID), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "INVALID_ORG_ID",
			Message: "Invalid organization ID format",
		})
		return
	}

	if _, err := uuid.Parse(empID); err != nil {
		c.logger.Warn("Invalid employee ID", zap.String("empId", empID), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "INVALID_EMP_ID",
			Message: "Invalid employee ID format",
		})
		return
	}

	var req dto.UpdateEmployeeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Warn("Invalid request body", zap.Error(err))

		if validationErr, ok := err.(*rcerrors.ValidationError); ok {
			ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "VALIDATION_FAILED",
				Message: "Request validation failed",
				Details: map[string]interface{}{
					"fields": validationErr.Fields,
				},
			})
			return
		}

		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "INVALID_REQUEST_BODY",
			Message: "Invalid request body format",
		})
		return
	}

	// Get user ID from context for updatedBy
	userID, exists := ctx.Get("user_id")
	if !exists {
		c.logger.Error("User ID not found in context")
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "UNAUTHORIZED",
			Message: "User not authenticated",
		})
		return
	}

	employee, err := c.empService.UpdateEmployee(ctx, empID, &req, uuid.MustParse(userID.(string)))
	if err != nil {
		c.logger.Error("Failed to update employee", zap.String("empId", empID), zap.Error(err))

		switch err.(type) {
		case *rcerrors.ValidationError:
			validationErr := err.(*rcerrors.ValidationError)
			ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "VALIDATION_FAILED",
				Message: validationErr.Error(),
				Details: map[string]interface{}{
					"fields": validationErr.Fields,
				},
			})
		case *rcerrors.NotFoundError:
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "EMPLOYEE_NOT_FOUND",
				Message: "Employee not found",
			})
		case *rcerrors.ForbiddenError:
			ctx.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "FORBIDDEN",
				Message: err.Error(),
			})
		default:
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "INTERNAL_ERROR",
				Message: "Failed to update employee",
			})
		}
		return
	}

	response := mapEmployeeToResponse(employee)
	ctx.JSON(http.StatusOK, response)
}

// UpdateEmployeeStatus actualiza el estado de un empleado
// @Summary Actualizar estado del empleado
// @Description Actualiza el estado de un empleado (activo/inactivo)
// @Tags employees
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param empId path string true "Employee ID"
// @Param status body dto.UpdateStatusRequest true "Status update data"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /organizations/{orgId}/employees/{empId}/status [patch]
func (c *EmployeeController) UpdateEmployeeStatus(ctx *gin.Context) {
	orgID := ctx.Param("orgId")
	empID := ctx.Param("empId")

	if _, err := uuid.Parse(orgID); err != nil {
		c.logger.Warn("Invalid organization ID", zap.String("orgId", orgID), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "INVALID_ORG_ID",
			Message: "Invalid organization ID format",
		})
		return
	}

	if _, err := uuid.Parse(empID); err != nil {
		c.logger.Warn("Invalid employee ID", zap.String("empId", empID), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "INVALID_EMP_ID",
			Message: "Invalid employee ID format",
		})
		return
	}

	var req dto.StatusUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Warn("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "INVALID_REQUEST_BODY",
			Message: "Invalid request body format",
		})
		return
	}

	// Get user ID from context for updatedBy
	userID, exists := ctx.Get("user_id")
	if !exists {
		c.logger.Error("User ID not found in context")
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "UNAUTHORIZED",
			Message: "User not authenticated",
		})
		return
	}

	// Convert string status to EmployeeStatus enum
	var status models.EmployeeStatus
	switch req.Status {
	case "active":
		status = models.EmpStatusActive
	case "inactive":
		status = models.EmpStatusInactive
	case "suspended":
		status = models.EmpStatusSuspended
	case "terminated":
		status = models.EmpStatusTerminated
	default:
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "INVALID_STATUS",
			Message: "Invalid employee status",
		})
		return
	}

	err := c.empService.UpdateEmployeeStatus(ctx, empID, status, uuid.MustParse(userID.(string)))
	if err != nil {
		c.logger.Error("Failed to update employee status", zap.String("empId", empID), zap.Error(err))

		switch err.(type) {
		case *rcerrors.ValidationError:
			validationErr := err.(*rcerrors.ValidationError)
			ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "VALIDATION_FAILED",
				Message: validationErr.Error(),
			})
		case *rcerrors.NotFoundError:
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "EMPLOYEE_NOT_FOUND",
				Message: "Employee not found",
			})
		case *rcerrors.ForbiddenError:
			ctx.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "FORBIDDEN",
				Message: err.Error(),
			})
		default:
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "INTERNAL_ERROR",
				Message: "Failed to update employee status",
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Employee status updated successfully",
	})
}

// GetEmployeeRoles obtiene los roles de un empleado
// @Summary Obtener roles del empleado
// @Description Obtiene la lista de roles asignados a un empleado
// @Tags employees
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param empId path string true "Employee ID"
// @Success 200 {object} dto.EmployeeRolesResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /organizations/{orgId}/employees/{empId}/roles [get]
func (c *EmployeeController) GetEmployeeRoles(ctx *gin.Context) {
	orgID := ctx.Param("orgId")
	empID := ctx.Param("empId")

	if _, err := uuid.Parse(orgID); err != nil {
		c.logger.Warn("Invalid organization ID", zap.String("orgId", orgID), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "INVALID_ORG_ID",
			Message: "Invalid organization ID format",
		})
		return
	}

	if _, err := uuid.Parse(empID); err != nil {
		c.logger.Warn("Invalid employee ID", zap.String("empId", empID), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "INVALID_EMP_ID",
			Message: "Invalid employee ID format",
		})
		return
	}

	roles, err := c.empService.GetEmployeeRoles(ctx, empID)
	if err != nil {
		c.logger.Error("Failed to get employee roles", zap.String("empId", empID), zap.Error(err))

		switch err.(type) {
		case *rcerrors.NotFoundError:
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "EMPLOYEE_NOT_FOUND",
				Message: "Employee not found",
			})
		case *rcerrors.ForbiddenError:
			ctx.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "FORBIDDEN",
				Message: err.Error(),
			})
		default:
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "INTERNAL_ERROR",
				Message: "Failed to retrieve employee roles",
			})
		}
		return
	}

	// Map roles to responses
	roleResponses := make([]dto.EmployeeRoleResponse, len(roles))
	for i, role := range roles {
		roleResponses[i] = dto.EmployeeRoleResponse{
			BaseResponse: dto.BaseResponse{
				ID:        role.ID,
				CreatedAt: role.CreatedAt,
				UpdatedAt: time.Now(), // EmployeeRole doesn't have UpdatedAt, use current time
				CreatedBy: role.CreatedBy,
			},
			EmployeeID: role.EmployeeID,
			RoleID:     role.RoleID,
			IsPrimary:  role.IsPrimary,
			Role: &dto.OrganizationRoleResponse{
				BaseResponse: dto.BaseResponse{
					ID:        role.Role.ID,
					CreatedAt: role.Role.CreatedAt,
					UpdatedAt: role.Role.UpdatedAt,
					CreatedBy: role.Role.CreatedBy,
					UpdatedBy: role.Role.UpdatedBy,
				},
				OrganizationID: role.Role.OrganizationID,
				Name:           role.Role.Name,
				Description:    role.Role.Description,
				IsDefault:      role.Role.IsDefault,
			},
		}
	}

	response := dto.EmployeeRolesResponse{
		EmployeeID: uuid.MustParse(empID),
		Roles:      roleResponses,
	}

	ctx.JSON(http.StatusOK, response)
}

// AssignRole asigna un rol a un empleado
// @Summary Asignar rol a empleado
// @Description Asigna un rol específico a un empleado
// @Tags employees
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param empId path string true "Employee ID"
// @Param roleId path string true "Role ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /organizations/{orgId}/employees/{empId}/roles/{roleId} [post]
func (c *EmployeeController) AssignRole(ctx *gin.Context) {
	orgID := ctx.Param("orgId")
	empID := ctx.Param("empId")
	roleID := ctx.Param("roleId")

	if _, err := uuid.Parse(orgID); err != nil {
		c.logger.Warn("Invalid organization ID", zap.String("orgId", orgID), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "INVALID_ORG_ID",
			Message: "Invalid organization ID format",
		})
		return
	}

	if _, err := uuid.Parse(empID); err != nil {
		c.logger.Warn("Invalid employee ID", zap.String("empId", empID), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "INVALID_EMP_ID",
			Message: "Invalid employee ID format",
		})
		return
	}

	if _, err := uuid.Parse(roleID); err != nil {
		c.logger.Warn("Invalid role ID", zap.String("roleId", roleID), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "INVALID_ROLE_ID",
			Message: "Invalid role ID format",
		})
		return
	}

	// Get user ID from context for createdBy
	userID, exists := ctx.Get("user_id")
	if !exists {
		c.logger.Error("User ID not found in context")
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "UNAUTHORIZED",
			Message: "User not authenticated",
		})
		return
	}

	// Create assign role request
	assignReq := &dto.AssignRoleRequest{
		RoleID:    uuid.MustParse(roleID),
		IsPrimary: false, // Default to false, can be made configurable
	}

	err := c.empService.AssignRole(ctx, empID, assignReq, uuid.MustParse(userID.(string)))
	if err != nil {
		c.logger.Error("Failed to assign role",
			zap.String("empId", empID),
			zap.String("roleId", roleID),
			zap.Error(err))

		switch err.(type) {
		case *rcerrors.ValidationError:
			validationErr := err.(*rcerrors.ValidationError)
			ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "VALIDATION_FAILED",
				Message: validationErr.Error(),
			})
		case *rcerrors.NotFoundError:
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "NOT_FOUND",
				Message: err.Error(),
			})
		case *rcerrors.ConflictError:
			ctx.JSON(http.StatusConflict, dto.ErrorResponse{
				Error:   "ROLE_ALREADY_ASSIGNED",
				Message: "Role already assigned to employee",
			})
		case *rcerrors.ForbiddenError:
			ctx.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "FORBIDDEN",
				Message: err.Error(),
			})
		default:
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "INTERNAL_ERROR",
				Message: "Failed to assign role",
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Role assigned successfully",
	})
}

// RemoveRole remueve un rol de un empleado
// @Summary Remover rol de empleado
// @Description Remueve un rol específico de un empleado
// @Tags employees
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param empId path string true "Employee ID"
// @Param roleId path string true "Role ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /organizations/{orgId}/employees/{empId}/roles/{roleId} [delete]
func (c *EmployeeController) RemoveRole(ctx *gin.Context) {
	orgID := ctx.Param("orgId")
	empID := ctx.Param("empId")
	roleID := ctx.Param("roleId")

	if _, err := uuid.Parse(orgID); err != nil {
		c.logger.Warn("Invalid organization ID", zap.String("orgId", orgID), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "INVALID_ORG_ID",
			Message: "Invalid organization ID format",
		})
		return
	}

	if _, err := uuid.Parse(empID); err != nil {
		c.logger.Warn("Invalid employee ID", zap.String("empId", empID), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "INVALID_EMP_ID",
			Message: "Invalid employee ID format",
		})
		return
	}

	if _, err := uuid.Parse(roleID); err != nil {
		c.logger.Warn("Invalid role ID", zap.String("roleId", roleID), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "INVALID_ROLE_ID",
			Message: "Invalid role ID format",
		})
		return
	}

	err := c.empService.RemoveRole(ctx, empID, roleID)
	if err != nil {
		c.logger.Error("Failed to remove role",
			zap.String("empId", empID),
			zap.String("roleId", roleID),
			zap.Error(err))

		switch err.(type) {
		case *rcerrors.ValidationError:
			validationErr := err.(*rcerrors.ValidationError)
			ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "VALIDATION_FAILED",
				Message: validationErr.Error(),
			})
		case *rcerrors.NotFoundError:
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "NOT_FOUND",
				Message: err.Error(),
			})
		case *rcerrors.ForbiddenError:
			ctx.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "FORBIDDEN",
				Message: err.Error(),
			})
		default:
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "INTERNAL_ERROR",
				Message: "Failed to remove role",
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Role removed successfully",
	})
}

// mapEmployeeToResponse convierte un modelo Employee a EmployeeResponse
func mapEmployeeToResponse(emp *models.Employee) *dto.EmployeeResponse {
	if emp == nil {
		return nil
	}

	return &dto.EmployeeResponse{
		BaseResponse: dto.BaseResponse{
			ID:        emp.ID,
			CreatedAt: emp.CreatedAt,
			UpdatedAt: emp.UpdatedAt,
			CreatedBy: emp.CreatedBy,
			UpdatedBy: emp.UpdatedBy,
		},
		OrganizationID: emp.OrganizationID,
		UserID:         emp.UserID,
		PersonID:       emp.PersonID,
		PrimaryRoleID:  emp.PrimaryRoleID,
		Status:         string(emp.Status),
		HiredAt:        emp.HiredAt,
		FiredAt:        emp.FiredAt,
	}
}
