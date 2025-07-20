package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/rem-gestion/api-suite/organization/src/dto"
	"github.com/rem-gestion/api-suite/organization/src/models"
	"github.com/rem-gestion/api-suite/organization/src/services"
	rcerrors "github.com/rem-gestion/rem-common/errors"
	"github.com/rem-gestion/rem-common/middleware"
)

// RoleController maneja las peticiones HTTP para roles
type RoleController struct {
	roleService *services.RoleService
	logger      *zap.Logger
}

// NewRoleController crea una nueva instancia del controlador de roles
func NewRoleController(roleService *services.RoleService, logger *zap.Logger) *RoleController {
	return &RoleController{
		roleService: roleService,
		logger:      logger,
	}
}

// ===============================
// ROLE CRUD OPERATIONS
// ===============================

// CreateRole crea un nuevo rol para la organización
// @Summary Crear rol
// @Description Crea un nuevo rol personalizado para la organización
// @Tags roles
// @Accept json
// @Produce json
// @Param org_id path string true "Organization ID"
// @Param request body dto.CreateRoleRequest true "Datos del rol"
// @Success 201 {object} dto.OrganizationRoleResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 422 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /organizations/{org_id}/roles [post]
func (c *RoleController) CreateRole(ctx *gin.Context) {
	// 1. Parsear y validar org_id
	orgID, err := uuid.Parse(ctx.Param("org_id"))
	if err != nil {
		c.logger.Warn("Invalid organization ID", zap.String("org_id", ctx.Param("org_id")), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_org_id",
			Message: "Invalid organization ID format",
		})
		return
	}

	// 2. Parsear request body
	var req dto.CreateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Warn("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request format",
		})
		return
	}

	// 3. Obtener información del usuario autenticado
	userClaims, _ := middleware.GetUserClaims(ctx)
	if userClaims == nil {
		c.logger.Error("User claims not found in context")
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "User authentication required",
		})
		return
	}

	creatorPersonID, err := uuid.Parse(userClaims.UserID)
	if err != nil {
		c.logger.Error("Invalid user ID in claims", zap.Error(err), zap.String("user_id", userClaims.UserID))
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "invalid_user",
			Message: "Invalid user identifier",
		})
		return
	}

	// 4. Crear rol usando el servicio
	role, err := c.roleService.CreateRole(ctx.Request.Context(), orgID, creatorPersonID, req)
	if err != nil {
		c.handleServiceError(ctx, err, "Failed to create role")
		return
	}

	// 5. Mapear a response DTO
	response := c.mapRoleToResponse(role)

	c.logger.Info("Role created successfully",
		zap.String("role_id", role.ID.String()),
		zap.String("org_id", orgID.String()),
	)

	ctx.JSON(http.StatusCreated, response)
}

// GetRole obtiene un rol específico
// @Summary Obtener rol
// @Description Obtiene información de un rol específico
// @Tags roles
// @Accept json
// @Produce json
// @Param org_id path string true "Organization ID"
// @Param role_id path string true "Role ID"
// @Success 200 {object} dto.OrganizationRoleResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /organizations/{org_id}/roles/{role_id} [get]
func (c *RoleController) GetRole(ctx *gin.Context) {
	// 1. Parsear y validar role_id
	roleID, err := uuid.Parse(ctx.Param("role_id"))
	if err != nil {
		c.logger.Warn("Invalid role ID", zap.String("role_id", ctx.Param("role_id")), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_role_id",
			Message: "Invalid role ID format",
		})
		return
	}

	// 2. Obtener información del usuario autenticado
	userClaims, _ := middleware.GetUserClaims(ctx)
	if userClaims == nil {
		c.logger.Error("User claims not found in context")
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "User authentication required",
		})
		return
	}

	requesterPersonID, err := uuid.Parse(userClaims.UserID)
	if err != nil {
		c.logger.Error("Invalid user ID in claims", zap.Error(err), zap.String("user_id", userClaims.UserID))
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "invalid_user",
			Message: "Invalid user identifier",
		})
		return
	}

	// 3. Obtener rol usando el servicio
	role, err := c.roleService.GetRole(ctx.Request.Context(), roleID, requesterPersonID)
	if err != nil {
		c.handleServiceError(ctx, err, "Failed to get role")
		return
	}

	// 4. Mapear a response DTO
	response := c.mapRoleToResponse(role)

	ctx.JSON(http.StatusOK, response)
}

// ListRoles lista todos los roles de la organización
// @Summary Listar roles
// @Description Lista todos los roles de la organización con filtros opcionales
// @Tags roles
// @Accept json
// @Produce json
// @Param org_id path string true "Organization ID"
// @Param page query int false "Número de página" default(1)
// @Param per_page query int false "Elementos por página" default(10)
// @Param is_default query bool false "Filtrar por roles default"
// @Param search query string false "Búsqueda por nombre"
// @Success 200 {object} dto.RoleListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /organizations/{org_id}/roles [get]
func (c *RoleController) ListRoles(ctx *gin.Context) {
	// 1. Parsear y validar org_id
	orgID, err := uuid.Parse(ctx.Param("org_id"))
	if err != nil {
		c.logger.Warn("Invalid organization ID", zap.String("org_id", ctx.Param("org_id")), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_org_id",
			Message: "Invalid organization ID format",
		})
		return
	}

	// 2. Parsear parámetros de paginación
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	// 3. Parsear filtros
	filters := dto.RoleFiltersRequest{
		FilterRequest: dto.FilterRequest{
			Search: ctx.Query("search"),
		},
		PaginationRequest: dto.PaginationRequest{
			Page:    page,
			PerPage: perPage,
		},
		OrganizationID: orgID,
	}

	// Parsear is_default si está presente
	if isDefaultStr := ctx.Query("is_default"); isDefaultStr != "" {
		if isDefaultValue, err := strconv.ParseBool(isDefaultStr); err == nil {
			filters.IsDefault = &isDefaultValue
		}
	}

	// 4. Obtener información del usuario autenticado
	userClaims, _ := middleware.GetUserClaims(ctx)
	if userClaims == nil {
		c.logger.Error("User claims not found in context")
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "User authentication required",
		})
		return
	}

	requesterPersonID, err := uuid.Parse(userClaims.UserID)
	if err != nil {
		c.logger.Error("Invalid user ID in claims", zap.Error(err), zap.String("user_id", userClaims.UserID))
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "invalid_user",
			Message: "Invalid user identifier",
		})
		return
	}

	// 5. Listar roles usando el servicio
	roles, total, err := c.roleService.ListRoles(ctx.Request.Context(), requesterPersonID, filters)
	if err != nil {
		c.handleServiceError(ctx, err, "Failed to list roles")
		return
	}

	// 6. Mapear a response DTOs
	roleResponses := make([]dto.OrganizationRoleResponse, len(roles))
	for i, role := range roles {
		roleResponses[i] = *c.mapRoleToResponse(&role)
	}

	// 7. Crear respuesta con metadatos de paginación
	response := dto.RoleListResponse{
		Data: roleResponses,
		Pagination: dto.PaginationResponse{
			Page:       page,
			PerPage:    perPage,
			Total:      total,
			TotalPages: int((total + int64(perPage) - 1) / int64(perPage)),
			HasNext:    page < int((total+int64(perPage)-1)/int64(perPage)),
			HasPrev:    page > 1,
		},
	}

	ctx.JSON(http.StatusOK, response)
}

// UpdateRole actualiza un rol existente
// @Summary Actualizar rol
// @Description Actualiza un rol existente de la organización
// @Tags roles
// @Accept json
// @Produce json
// @Param org_id path string true "Organization ID"
// @Param role_id path string true "Role ID"
// @Param request body dto.UpdateRoleRequest true "Datos actualizados del rol"
// @Success 200 {object} dto.RoleResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 422 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /organizations/{org_id}/roles/{role_id} [put]
func (c *RoleController) UpdateRole(ctx *gin.Context) {
	// 1. Parsear y validar role_id
	roleID, err := uuid.Parse(ctx.Param("role_id"))
	if err != nil {
		c.logger.Warn("Invalid role ID", zap.String("role_id", ctx.Param("role_id")), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_role_id",
			Message: "Invalid role ID format",
		})
		return
	}

	// 2. Parsear request body
	var req dto.UpdateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Warn("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request format",
		})
		return
	}

	// 3. Obtener información del usuario autenticado
	userClaims, _ := middleware.GetUserClaims(ctx)
	if userClaims == nil {
		c.logger.Error("User claims not found in context")
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "User authentication required",
		})
		return
	}

	updaterPersonID, err := uuid.Parse(userClaims.UserID)
	if err != nil {
		c.logger.Error("Invalid user ID in claims", zap.Error(err), zap.String("user_id", userClaims.UserID))
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "invalid_user",
			Message: "Invalid user identifier",
		})
		return
	}

	// 4. Actualizar rol usando el servicio
	role, err := c.roleService.UpdateRole(ctx.Request.Context(), roleID, updaterPersonID, req)
	if err != nil {
		c.handleServiceError(ctx, err, "Failed to update role")
		return
	}

	// 5. Mapear a response DTO
	response := c.mapRoleToResponse(role)

	c.logger.Info("Role updated successfully", zap.String("role_id", roleID.String()))

	ctx.JSON(http.StatusOK, response)
}

// DeleteRole elimina (soft delete) un rol
// @Summary Eliminar rol
// @Description Elimina (soft delete) un rol personalizado de la organización
// @Tags roles
// @Accept json
// @Produce json
// @Param org_id path string true "Organization ID"
// @Param role_id path string true "Role ID"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 422 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /organizations/{org_id}/roles/{role_id} [delete]
func (c *RoleController) DeleteRole(ctx *gin.Context) {
	// 1. Parsear y validar role_id
	roleID, err := uuid.Parse(ctx.Param("role_id"))
	if err != nil {
		c.logger.Warn("Invalid role ID", zap.String("role_id", ctx.Param("role_id")), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_role_id",
			Message: "Invalid role ID format",
		})
		return
	}

	// 2. Obtener información del usuario autenticado
	userClaims, _ := middleware.GetUserClaims(ctx)
	if userClaims == nil {
		c.logger.Error("User claims not found in context")
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "User authentication required",
		})
		return
	}

	deleterPersonID, err := uuid.Parse(userClaims.UserID)
	if err != nil {
		c.logger.Error("Invalid user ID in claims", zap.Error(err), zap.String("user_id", userClaims.UserID))
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "invalid_user",
			Message: "Invalid user identifier",
		})
		return
	}

	// 3. Eliminar rol usando el servicio
	err = c.roleService.DeleteRole(ctx.Request.Context(), roleID, deleterPersonID)
	if err != nil {
		c.handleServiceError(ctx, err, "Failed to delete role")
		return
	}

	c.logger.Info("Role deleted successfully", zap.String("role_id", roleID.String()))

	ctx.JSON(http.StatusNoContent, nil)
}

// ===============================
// HELPER METHODS
// ===============================

// handleServiceError maneja errores del servicio y los convierte a respuestas HTTP apropiadas
func (c *RoleController) handleServiceError(ctx *gin.Context, err error, logMessage string) {
	c.logger.Error(logMessage, zap.Error(err))

	switch e := err.(type) {
	case *rcerrors.NotFoundError:
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "not_found",
			Message: e.Error(),
		})
	case *rcerrors.ValidationError:
		ctx.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
			Error:   "validation_error",
			Message: e.Error(),
		})
	case *rcerrors.ConflictError:
		ctx.JSON(http.StatusConflict, dto.ErrorResponse{
			Error:   "conflict",
			Message: e.Error(),
		})
	case *rcerrors.ForbiddenError:
		ctx.JSON(http.StatusForbidden, dto.ErrorResponse{
			Error:   "forbidden",
			Message: e.Error(),
		})
	default:
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "An internal error occurred",
		})
	}
}

// mapRoleToResponse convierte un modelo OrganizationRole a DTO response
func (c *RoleController) mapRoleToResponse(role *models.OrganizationRole) *dto.OrganizationRoleResponse {
	return &dto.OrganizationRoleResponse{
		BaseResponse: dto.BaseResponse{
			ID:        role.ID,
			CreatedAt: role.CreatedAt,
			UpdatedAt: role.UpdatedAt,
			CreatedBy: role.CreatedBy,
			UpdatedBy: role.UpdatedBy,
		},
		OrganizationID: role.OrganizationID,
		Name:           role.Name,
		Description:    role.Description,
		IsDefault:      role.IsDefault,
	}
}
