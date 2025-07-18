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

// OwnerController maneja las peticiones HTTP para propietarios de organizaciones
type OwnerController struct {
	ownerService *services.OwnerService
	logger       *zap.Logger
}

// NewOwnerController crea una nueva instancia del controlador de propietarios
func NewOwnerController(ownerService *services.OwnerService, logger *zap.Logger) *OwnerController {
	return &OwnerController{
		ownerService: ownerService,
		logger:       logger,
	}
}

// ===============================
// OWNER CRUD OPERATIONS
// ===============================

// AddOwner agrega un nuevo propietario a la organización
// @Summary Agregar propietario
// @Description Agrega un nuevo propietario a la organización
// @Tags owners
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param owner body dto.AddOwnerRequest true "Owner data"
// @Success 201 {object} dto.OwnerResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /organizations/{orgId}/owners [post]
func (c *OwnerController) AddOwner(ctx *gin.Context) {
	c.logger.Info("Adding owner", zap.String("method", "POST"), zap.String("path", ctx.Request.URL.Path))

	// 1. Validar organization ID
	orgIDStr := ctx.Param("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.logger.Warn("Invalid organization ID", zap.String("org_id", orgIDStr), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_org_id",
			Message: "Invalid organization identifier",
		})
		return
	}

	// 2. Obtener información del usuario autenticado
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

	// 3. Parsear y validar request body
	var req dto.AddOwnerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Warn("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	// 4. Crear propietario usando el servicio
	owner, err := c.ownerService.AddOwner(ctx.Request.Context(), orgID, createdBy, req)
	if err != nil {
		c.handleServiceError(ctx, err, "Failed to add owner")
		return
	}

	// 5. Mapear a response DTO
	response := c.mapOwnerToResponse(owner)

	c.logger.Info("Owner added successfully",
		zap.String("org_id", orgID.String()),
		zap.String("owner_id", owner.ID.String()),
		zap.String("created_by", createdBy.String()),
	)

	ctx.JSON(http.StatusCreated, response)
}

// ListOwners obtiene la lista de propietarios de una organización
// @Summary Listar propietarios
// @Description Obtiene la lista de propietarios de una organización con filtros opcionales
// @Tags owners
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param is_founder query bool false "Filter by founder status"
// @Param min_ownership query float64 false "Minimum ownership percentage"
// @Param max_ownership query float64 false "Maximum ownership percentage"
// @Success 200 {object} dto.OwnerListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /organizations/{orgId}/owners [get]
func (c *OwnerController) ListOwners(ctx *gin.Context) {
	c.logger.Info("Listing owners", zap.String("method", "GET"), zap.String("path", ctx.Request.URL.Path))

	// 1. Validar organization ID
	orgIDStr := ctx.Param("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.logger.Warn("Invalid organization ID", zap.String("org_id", orgIDStr), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_org_id",
			Message: "Invalid organization identifier",
		})
		return
	}

	// 2. Parsear filtros de la query string
	filters, err := c.parseOwnerFilters(ctx)
	if err != nil {
		c.logger.Warn("Invalid query filters", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_filters",
			Message: "Invalid query parameters: " + err.Error(),
		})
		return
	}

	// Asegurar que el orgID coincida con el parámetro de la URL
	filters.OrganizationID = orgID

	// 3. Obtener información del usuario autenticado para permisos
	userClaims, exists := middleware.GetUserClaims(ctx)
	if !exists {
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

	// 4. Obtener lista de propietarios
	owners, totalCount, err := c.ownerService.ListOwners(ctx.Request.Context(), requesterPersonID, *filters)
	if err != nil {
		c.handleServiceError(ctx, err, "Failed to list owners")
		return
	}

	// 5. Mapear a response DTOs
	ownerResponses := make([]dto.OrganizationOwnerResponse, len(owners))
	for i, owner := range owners {
		ownerResponses[i] = c.mapOwnerToResponse(&owner)
	}

	// 6. Crear respuesta paginada
	response := dto.OwnerListResponse{
		Data: ownerResponses,
		Pagination: dto.PaginationResponse{
			Page:       filters.Page,
			PerPage:    filters.PerPage,
			Total:      totalCount,
			TotalPages: int((totalCount + int64(filters.PerPage) - 1) / int64(filters.PerPage)),
			HasNext:    filters.Page < int((totalCount+int64(filters.PerPage)-1)/int64(filters.PerPage)),
			HasPrev:    filters.Page > 1,
		},
	}

	c.logger.Info("Owners listed successfully",
		zap.String("org_id", orgID.String()),
		zap.Int("count", len(owners)),
		zap.Int64("total", totalCount),
	)

	ctx.JSON(http.StatusOK, response)
}

// GetOwner obtiene un propietario específico
// @Summary Obtener propietario
// @Description Obtiene los detalles de un propietario específico
// @Tags owners
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param ownerId path string true "Owner ID"
// @Success 200 {object} dto.OwnerResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /organizations/{orgId}/owners/{ownerId} [get]
func (c *OwnerController) GetOwner(ctx *gin.Context) {
	c.logger.Info("Getting owner", zap.String("method", "GET"), zap.String("path", ctx.Request.URL.Path))

	// 1. Validar organization ID
	orgIDStr := ctx.Param("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.logger.Warn("Invalid organization ID", zap.String("org_id", orgIDStr), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_org_id",
			Message: "Invalid organization identifier",
		})
		return
	}

	// 2. Validar owner ID
	ownerIDStr := ctx.Param("id")
	ownerID, err := uuid.Parse(ownerIDStr)
	if err != nil {
		c.logger.Warn("Invalid owner ID", zap.String("owner_id", ownerIDStr), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_owner_id",
			Message: "Invalid owner identifier",
		})
		return
	}

	// 3. Obtener información del usuario autenticado para permisos
	userClaims, exists := middleware.GetUserClaims(ctx)
	if !exists {
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

	// 4. Obtener propietario
	owner, err := c.ownerService.GetOwner(ctx.Request.Context(), ownerID, requesterPersonID)
	if err != nil {
		c.handleServiceError(ctx, err, "Failed to get owner")
		return
	}

	// 5. Verificar que el propietario pertenece a la organización
	if owner.OrganizationID != orgID {
		c.logger.Warn("Owner does not belong to organization",
			zap.String("owner_id", ownerID.String()),
			zap.String("org_id", orgID.String()),
			zap.String("owner_org_id", owner.OrganizationID.String()),
		)
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "owner_not_found",
			Message: "Owner not found in this organization",
		})
		return
	}

	// 6. Mapear a response DTO
	response := c.mapOwnerToResponse(owner)

	c.logger.Info("Owner retrieved successfully",
		zap.String("org_id", orgID.String()),
		zap.String("owner_id", ownerID.String()),
	)

	ctx.JSON(http.StatusOK, response)
}

// RemoveOwner elimina un propietario de la organización
// @Summary Eliminar propietario
// @Description Elimina un propietario de la organización
// @Tags owners
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param ownerId path string true "Owner ID"
// @Param request body dto.RemoveOwnerRequest true "Remove owner data"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /organizations/{orgId}/owners/{ownerId} [delete]
func (c *OwnerController) RemoveOwner(ctx *gin.Context) {
	c.logger.Info("Removing owner", zap.String("method", "DELETE"), zap.String("path", ctx.Request.URL.Path))

	// 1. Validar organization ID
	orgIDStr := ctx.Param("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.logger.Warn("Invalid organization ID", zap.String("org_id", orgIDStr), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_org_id",
			Message: "Invalid organization identifier",
		})
		return
	}

	// 2. Validar owner ID
	ownerIDStr := ctx.Param("id")
	ownerID, err := uuid.Parse(ownerIDStr)
	if err != nil {
		c.logger.Warn("Invalid owner ID", zap.String("owner_id", ownerIDStr), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_owner_id",
			Message: "Invalid owner identifier",
		})
		return
	}

	// 3. Obtener información del usuario autenticado
	userClaims, exists := middleware.GetUserClaims(ctx)
	if !exists {
		c.logger.Error("User claims not found in context")
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "User authentication required",
		})
		return
	}

	removedBy, err := uuid.Parse(userClaims.UserID)
	if err != nil {
		c.logger.Error("Invalid user ID in claims", zap.Error(err), zap.String("user_id", userClaims.UserID))
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "invalid_user",
			Message: "Invalid user identifier",
		})
		return
	}

	// 4. Parsear request body
	var req dto.RemoveOwnerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Warn("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	// 5. Eliminar propietario usando el servicio
	err = c.ownerService.RemoveOwner(ctx.Request.Context(), ownerID, removedBy, req)
	if err != nil {
		c.handleServiceError(ctx, err, "Failed to remove owner")
		return
	}

	c.logger.Info("Owner removed successfully",
		zap.String("org_id", orgID.String()),
		zap.String("owner_id", ownerID.String()),
		zap.String("removed_by", removedBy.String()),
	)

	ctx.Status(http.StatusNoContent)
}

// ===============================
// HELPER METHODS
// ===============================

// handleServiceError maneja errores del servicio y los convierte a respuestas HTTP apropiadas
func (c *OwnerController) handleServiceError(ctx *gin.Context, err error, logMessage string) {
	c.logger.Error(logMessage, zap.Error(err))

	switch e := err.(type) {
	case *rcerrors.ValidationError:
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: e.Error(),
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

// mapOwnerToResponse convierte un modelo Owner a DTO response
func (c *OwnerController) mapOwnerToResponse(owner *models.OrganizationOwner) dto.OrganizationOwnerResponse {
	return dto.OrganizationOwnerResponse{
		BaseResponse: dto.BaseResponse{
			ID:        owner.ID,
			CreatedAt: owner.CreatedAt,
			UpdatedAt: owner.UpdatedAt,
			CreatedBy: owner.CreatedBy,
			UpdatedBy: owner.UpdatedBy,
		},
		OrganizationID: owner.OrganizationID,
		PersonID:       owner.PersonID,
		UserID:         owner.UserID,
		IsFounder:      owner.IsFounder,
		OwnershipPct:   owner.OwnershipPct,
	}
}

// parseOwnerFilters parsea los filtros de la query string
func (c *OwnerController) parseOwnerFilters(ctx *gin.Context) (*dto.OwnerFiltersRequest, error) {
	filters := &dto.OwnerFiltersRequest{
		PaginationRequest: dto.PaginationRequest{
			Page:    1,
			PerPage: 10,
		},
	}

	// Page
	if pageStr := ctx.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			filters.Page = page
		}
	}

	// PerPage
	if perPageStr := ctx.Query("per_page"); perPageStr != "" {
		if perPage, err := strconv.Atoi(perPageStr); err == nil && perPage > 0 && perPage <= 100 {
			filters.PerPage = perPage
		}
	}

	// IsFounder
	if isFounderStr := ctx.Query("is_founder"); isFounderStr != "" {
		if isFounder, err := strconv.ParseBool(isFounderStr); err == nil {
			filters.IsFounder = &isFounder
		}
	}

	// MinOwnership
	if minOwnershipStr := ctx.Query("min_ownership"); minOwnershipStr != "" {
		if minOwnership, err := strconv.ParseFloat(minOwnershipStr, 64); err == nil {
			filters.MinOwnership = &minOwnership
		}
	}

	// MaxOwnership
	if maxOwnershipStr := ctx.Query("max_ownership"); maxOwnershipStr != "" {
		if maxOwnership, err := strconv.ParseFloat(maxOwnershipStr, 64); err == nil {
			filters.MaxOwnership = &maxOwnership
		}
	}

	// PersonID
	if personIDStr := ctx.Query("person_id"); personIDStr != "" {
		if personID, err := uuid.Parse(personIDStr); err == nil {
			filters.PersonID = &personID
		}
	}

	// UserID
	if userIDStr := ctx.Query("user_id"); userIDStr != "" {
		if userID, err := uuid.Parse(userIDStr); err == nil {
			filters.UserID = &userID
		}
	}

	return filters, nil
}
