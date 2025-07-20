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

// OrganizationController maneja las peticiones HTTP para organizaciones
type OrganizationController struct {
	orgService *services.OrganizationService
	logger     *zap.Logger
}

// NewOrganizationController crea una nueva instancia del controlador de organizaciones
func NewOrganizationController(orgService *services.OrganizationService, logger *zap.Logger) *OrganizationController {
	return &OrganizationController{
		orgService: orgService,
		logger:     logger,
	}
}

// ===============================
// ORGANIZATION CRUD OPERATIONS
// ===============================

// CreateOrganization crea una nueva organización
// @Summary Crear organización
// @Description Crea una nueva organización con validaciones de negocio
// @Tags organizations
// @Accept json
// @Produce json
// @Param request body dto.CreateOrganizationRequest true "Datos de la organización"
// @Success 201 {object} dto.OrganizationResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 422 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /organizations [post]
func (c *OrganizationController) CreateOrganization(ctx *gin.Context) {
	c.logger.Info("Creating organization", zap.String("method", "POST"), zap.String("path", ctx.Request.URL.Path))

	// 1. Obtener información del usuario autenticado
	userClaims, exists := middleware.GetUserClaims(ctx)
	if !exists {
		c.logger.Error("User claims not found in context")
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "User authentication required",
		})
		return
	}

	createdBy, err := uuid.Parse(userClaims.UserID)
	if err != nil {
		c.logger.Error("Invalid user ID in claims", zap.Error(err), zap.String("user_id", userClaims.UserID))
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "invalid_user",
			Message: "Invalid user identifier",
		})
		return
	}

	// 2. Parsear y validar request body
	var req dto.CreateOrganizationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Warn("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	// 3. Crear organización usando el servicio
	org, err := c.orgService.CreateOrganization(ctx.Request.Context(), &req, createdBy)
	if err != nil {
		c.handleServiceError(ctx, err, "Failed to create organization")
		return
	}

	// 4. Mapear a response DTO
	response := c.mapOrganizationToResponse(org)

	c.logger.Info("Organization created successfully",
		zap.String("org_id", org.ID.String()),
		zap.String("created_by", createdBy.String()),
	)

	ctx.JSON(http.StatusCreated, response)
}

// GetOrganization obtiene una organización por ID
// @Summary Obtener organización
// @Description Obtiene los detalles de una organización por su ID
// @Tags organizations
// @Accept json
// @Produce json
// @Param id path string true "ID de la organización"
// @Param include query string false "Recursos relacionados a incluir (settings,stats)" example("settings,stats")
// @Success 200 {object} dto.OrganizationResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /organizations/{id} [get]
func (c *OrganizationController) GetOrganization(ctx *gin.Context) {
	c.logger.Debug("Getting organization", zap.String("method", "GET"), zap.String("path", ctx.Request.URL.Path))

	// 1. Obtener información del usuario autenticado
	userClaims, exists := middleware.GetUserClaims(ctx)
	if !exists {
		c.logger.Error("User claims not found in context")
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "User authentication required",
		})
		return
	}

	userID, err := uuid.Parse(userClaims.UserID)
	if err != nil {
		c.logger.Error("Invalid user ID in claims", zap.Error(err))
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "invalid_user",
			Message: "Invalid user identifier",
		})
		return
	}

	// 2. Validar y parsear ID de la organización
	orgIDStr := ctx.Param("org_id")
	if orgIDStr == "" {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "missing_org_id",
			Message: "Organization ID is required",
		})
		return
	}

	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_org_id",
			Message: "Invalid organization ID format",
		})
		return
	}

	// 3. Obtener organización del servicio
	org, err := c.orgService.GetOrganization(ctx.Request.Context(), orgID, userID)
	if err != nil {
		c.handleServiceError(ctx, err, "Failed to get organization")
		return
	}

	// 4. Mapear a response DTO
	response := c.mapOrganizationToResponse(org)

	// 4. Incluir datos opcionales basado en query parameters
	includeParams := ctx.Query("include")
	if includeParams != "" {
		if err := c.enrichOrganizationResponse(ctx, orgID.String(), includeParams, response); err != nil {
			c.logger.Warn("Failed to enrich organization response", zap.Error(err))
			// No falla la request, solo registra el warning
		}
	}

	ctx.JSON(http.StatusOK, response)
}

// ListOrganizations lista organizaciones con filtros y paginación
// @Summary Listar organizaciones
// @Description Lista organizaciones con filtros, paginación y ordenamiento
// @Tags organizations
// @Accept json
// @Produce json
// @Param page query int false "Número de página" default(1)
// @Param per_page query int false "Elementos por página" default(10)
// @Param status query string false "Filtrar por estado" Enums(active, inactive, suspended)
// @Param type query string false "Filtrar por tipo" Enums(company, nonprofit, government)
// @Param search query string false "Búsqueda por nombre"
// @Param sort_by query string false "Campo de ordenamiento" Enums(name, created_at, updated_at)
// @Param sort_order query string false "Dirección de ordenamiento" Enums(asc, desc)
// @Param include query string false "Recursos relacionados a incluir"
// @Success 200 {object} dto.OrganizationListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /organizations [get]
func (c *OrganizationController) ListOrganizations(ctx *gin.Context) {
	c.logger.Debug("Listing organizations", zap.String("method", "GET"))

	// 1. Obtener información del usuario autenticado
	userClaims, exists := middleware.GetUserClaims(ctx)
	if !exists {
		c.logger.Error("User claims not found in context")
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "User authentication required",
		})
		return
	}

	userID, err := uuid.Parse(userClaims.UserID)
	if err != nil {
		c.logger.Error("Invalid user ID in claims", zap.Error(err))
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "invalid_user",
			Message: "Invalid user identifier",
		})
		return
	}

	// 2. Parsear y validar filtros de la query string
	filters, err := c.parseOrganizationFilters(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_filters",
			Message: err.Error(),
		})
		return
	}

	// 3. Listar organizaciones usando el servicio
	organizations, total, err := c.orgService.ListOrganizations(ctx.Request.Context(), userID, *filters)
	if err != nil {
		c.handleServiceError(ctx, err, "Failed to list organizations")
		return
	}

	// 4. Mapear a response DTOs
	responses := make([]dto.OrganizationResponse, len(organizations))
	for i, org := range organizations {
		responses[i] = *c.mapOrganizationToResponse(&org)
	}

	// 4. Crear respuesta paginada
	pagination := dto.PaginationResponse{
		Page:       filters.Page,
		PerPage:    filters.PerPage,
		Total:      total,
		TotalPages: int((total + int64(filters.PerPage) - 1) / int64(filters.PerPage)),
	}

	response := dto.OrganizationListResponse{
		Data:       responses,
		Pagination: pagination,
	}

	ctx.JSON(http.StatusOK, response)
}

// UpdateOrganization actualiza una organización existente
// @Summary Actualizar organización
// @Description Actualiza los datos de una organización existente
// @Tags organizations
// @Accept json
// @Produce json
// @Param id path string true "ID de la organización"
// @Param request body dto.UpdateOrganizationRequest true "Datos a actualizar"
// @Success 200 {object} dto.OrganizationResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 422 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /organizations/{id} [put]
func (c *OrganizationController) UpdateOrganization(ctx *gin.Context) {
	c.logger.Info("Updating organization", zap.String("method", "PUT"))

	// 1. Obtener información del usuario autenticado
	userClaims, exists := middleware.GetUserClaims(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "User authentication required",
		})
		return
	}

	updatedBy, err := uuid.Parse(userClaims.UserID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "invalid_user",
			Message: "Invalid user identifier",
		})
		return
	}

	// 2. Validar y parsear ID de la organización
	orgIDStr := ctx.Param("org_id")
	if orgIDStr == "" {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "missing_org_id",
			Message: "Organization ID is required",
		})
		return
	}

	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_org_id",
			Message: "Invalid organization ID format",
		})
		return
	}

	// 3. Parsear y validar request body
	var req dto.UpdateOrganizationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Warn("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	// 4. Actualizar organización usando el servicio
	org, err := c.orgService.UpdateOrganization(ctx.Request.Context(), orgID, &req, updatedBy)
	if err != nil {
		c.handleServiceError(ctx, err, "Failed to update organization")
		return
	}

	// 5. Mapear a response DTO
	response := c.mapOrganizationToResponse(org)

	c.logger.Info("Organization updated successfully",
		zap.String("org_id", orgID.String()),
		zap.String("updated_by", updatedBy.String()),
	)

	ctx.JSON(http.StatusOK, response)
}

// UpdateOrganizationStatus actualiza el estado de una organización
// @Summary Actualizar estado de organización
// @Description Actualiza únicamente el estado de una organización
// @Tags organizations
// @Accept json
// @Produce json
// @Param id path string true "ID de la organización"
// @Param request body dto.UpdateOrganizationStatusRequest true "Nuevo estado"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 422 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /organizations/{id}/status [patch]
func (c *OrganizationController) UpdateOrganizationStatus(ctx *gin.Context) {
	c.logger.Info("Updating organization status", zap.String("method", "PATCH"))

	// 1. Obtener información del usuario autenticado
	userClaims, exists := middleware.GetUserClaims(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "User authentication required",
		})
		return
	}

	updatedBy, err := uuid.Parse(userClaims.UserID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "invalid_user",
			Message: "Invalid user identifier",
		})
		return
	}

	// 2. Validar y parsear ID de la organización
	orgIDStr := ctx.Param("org_id")
	if orgIDStr == "" {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "missing_org_id",
			Message: "Organization ID is required",
		})
		return
	}

	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_org_id",
			Message: "Invalid organization ID format",
		})
		return
	}

	// 3. Parsear y validar request body
	var req dto.UpdateOrganizationStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Warn("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	// 4. Actualizar estado usando el servicio
	err = c.orgService.UpdateOrganizationStatus(ctx.Request.Context(), orgID, string(req.Status), updatedBy)
	if err != nil {
		c.handleServiceError(ctx, err, "Failed to update organization status")
		return
	}

	c.logger.Info("Organization status updated successfully",
		zap.String("org_id", orgID.String()),
		zap.String("new_status", string(req.Status)),
		zap.String("updated_by", userClaims.UserID),
	)

	ctx.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Organization status updated successfully",
	})
}

// ===============================
// ORGANIZATION SETTINGS OPERATIONS
// ===============================

// GetOrganizationSettings obtiene las configuraciones de una organización
// @Summary Obtener configuraciones de organización
// @Description Obtiene todas las configuraciones de una organización
// @Tags organizations
// @Accept json
// @Produce json
// @Param id path string true "ID de la organización"
// @Success 200 {object} dto.OrganizationSettingsResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /organizations/{id}/settings [get]
func (c *OrganizationController) GetOrganizationSettings(ctx *gin.Context) {
	c.logger.Debug("Getting organization settings")

	// 1. Validar y parsear ID de la organización
	orgID := ctx.Param("org_id")
	if orgID == "" {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "missing_org_id",
			Message: "Organization ID is required",
		})
		return
	}

	// 1. Parsear el ID de la organización
	orgUUID, err := uuid.Parse(orgID)
	if err != nil {
		c.logger.Error("Invalid organization ID format", zap.String("org_id", orgID), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid organization ID format",
			Code:    "INVALID_ID_FORMAT",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	// 2. Obtener configuraciones del servicio
	settings, err := c.orgService.GetOrganizationSettings(ctx.Request.Context(), orgUUID, uuid.Nil) // TODO: Agregar validación de usuario
	if err != nil {
		c.handleServiceError(ctx, err, "Failed to get organization settings")
		return
	}

	// 3. Crear respuesta
	response := dto.OrganizationSettingsResponse{
		OrganizationID: orgUUID,
		Settings:       settings,
	}

	ctx.JSON(http.StatusOK, response)
}

// UpdateOrganizationSettings actualiza las configuraciones de una organización
// @Summary Actualizar configuraciones de organización
// @Description Actualiza las configuraciones de una organización
// @Tags organizations
// @Accept json
// @Produce json
// @Param id path string true "ID de la organización"
// @Param request body dto.UpdateOrganizationSettingsRequest true "Configuraciones a actualizar"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 422 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /organizations/{id}/settings [put]
func (c *OrganizationController) UpdateOrganizationSettings(ctx *gin.Context) {
	c.logger.Info("Updating organization settings")

	// 1. Obtener información del usuario autenticado
	userClaims, exists := middleware.GetUserClaims(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "User authentication required",
		})
		return
	}

	// 2. Validar y parsear ID de la organización
	orgIDStr := ctx.Param("org_id")
	if orgIDStr == "" {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "missing_org_id",
			Message: "Organization ID is required",
		})
		return
	}

	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_org_id",
			Message: "Invalid organization ID format",
		})
		return
	}

	// 3. Parsear y validar request body
	var req dto.UpdateOrganizationSettingsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Warn("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	// 4. Actualizar configuraciones usando el servicio
	updatedBy, _ := uuid.Parse(userClaims.UserID) // Error ya validado arriba
	err = c.orgService.UpdateOrganizationSettings(ctx.Request.Context(), orgID, req.Settings, updatedBy)
	if err != nil {
		c.handleServiceError(ctx, err, "Failed to update organization settings")
		return
	}

	c.logger.Info("Organization settings updated successfully",
		zap.String("org_id", orgID.String()),
		zap.String("updated_by", userClaims.UserID),
	)

	ctx.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Organization settings updated successfully",
	})
}

// ===============================
// HELPER METHODS
// ===============================

// handleServiceError maneja errores del servicio y los convierte a respuestas HTTP apropiadas
func (c *OrganizationController) handleServiceError(ctx *gin.Context, err error, logMessage string) {
	c.logger.Error(logMessage, zap.Error(err))

	switch e := err.(type) {
	case *rcerrors.ValidationError:
		ctx.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
			Error:   "validation_error",
			Message: e.Error(),
			Details: map[string]interface{}{
				"fields": e.Fields,
			},
		})
	case *rcerrors.NotFoundError:
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "not_found",
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

// mapOrganizationToResponse convierte un modelo Organization a DTO response
func (c *OrganizationController) mapOrganizationToResponse(org *models.Organization) *dto.OrganizationResponse {
	response := &dto.OrganizationResponse{
		BaseResponse: dto.BaseResponse{
			ID:        org.ID,
			CreatedAt: org.CreatedAt,
			UpdatedAt: org.UpdatedAt,
		},
		Name:        org.Name,
		DisplayName: org.DisplayName,
		Status:      string(org.Status),
	}

	// Campos opcionales
	if org.Slug != nil {
		response.Slug = *org.Slug
	}
	if org.Description != nil {
		response.Description = org.Description
	}
	if org.Type != nil {
		response.Type = *org.Type
	}
	if org.LegalName != nil {
		response.LegalName = org.LegalName
	}
	if org.TaxID != nil {
		response.TaxID = org.TaxID
	}
	if org.Website != nil {
		response.Website = org.Website
	}
	if org.Phone != nil {
		response.Phone = org.Phone
	}
	if org.Email != nil {
		response.Email = org.Email
	}
	if org.LogoURL != nil {
		response.LogoURL = org.LogoURL
	}
	if org.TimezoneID != nil {
		response.TimezoneID = *org.TimezoneID
	}
	if org.FiscalAddressID != nil {
		response.FiscalAddressID = org.FiscalAddressID
	}

	return response
}

// parseOrganizationFilters parsea los filtros de la query string
func (c *OrganizationController) parseOrganizationFilters(ctx *gin.Context) (*dto.OrganizationFiltersRequest, error) {
	filters := &dto.OrganizationFiltersRequest{
		FilterRequest: dto.FilterRequest{
			Search: ctx.Query("search"),
		},
		PaginationRequest: dto.PaginationRequest{
			Page:    1,
			PerPage: 10,
		},
	}

	// Parsear paginación
	if pageStr := ctx.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			filters.Page = page
		}
	}

	if perPageStr := ctx.Query("per_page"); perPageStr != "" {
		if perPage, err := strconv.Atoi(perPageStr); err == nil && perPage > 0 && perPage <= 100 {
			filters.PerPage = perPage
		}
	}

	// Parsear filtros específicos
	if status := ctx.Query("status"); status != "" {
		filters.Status = status
	}

	if orgType := ctx.Query("type"); orgType != "" {
		filters.Type = orgType
	}

	return filters, nil
}

// enrichOrganizationResponse enriquece la respuesta con datos adicionales basado en el parámetro include
func (c *OrganizationController) enrichOrganizationResponse(ctx *gin.Context, orgID string, includeParams string, response *dto.OrganizationResponse) error {
	// Para esta implementación simplificada, solo registramos que debería implementarse
	c.logger.Debug("Enriching organization response",
		zap.String("org_id", orgID),
		zap.String("include", includeParams),
	)

	// TODO: Implementar enrichment basado en parámetros include
	// Ejemplos: settings, stats, branches, employees_count, etc.

	return nil
}
