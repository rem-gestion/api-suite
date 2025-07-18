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

// BranchController maneja las peticiones HTTP para sucursales
type BranchController struct {
	branchService *services.BranchService
	logger        *zap.Logger
}

// NewBranchController crea una nueva instancia del controlador de sucursales
func NewBranchController(branchService *services.BranchService, logger *zap.Logger) *BranchController {
	return &BranchController{
		branchService: branchService,
		logger:        logger,
	}
}

// ===============================
// BRANCH CRUD OPERATIONS
// ===============================

// CreateBranch crea una nueva sucursal
// @Summary Crear sucursal
// @Description Crea una nueva sucursal para una organización
// @Tags branches
// @Accept json
// @Produce json
// @Param org_id path string true "Organization ID"
// @Param request body dto.CreateBranchRequest true "Datos de la sucursal"
// @Success 201 {object} dto.OrganizationBranchResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 422 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /organizations/{org_id}/branches [post]
func (c *BranchController) CreateBranch(ctx *gin.Context) {
	c.logger.Info("Creating branch", zap.String("method", "POST"), zap.String("path", ctx.Request.URL.Path))

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

	// 2. Parsear org_id de la URL
	orgIDStr := ctx.Param("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.logger.Warn("Invalid organization ID", zap.Error(err), zap.String("org_id", orgIDStr))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_org_id",
			Message: "Invalid organization identifier",
		})
		return
	}

	// 3. Parsear y validar request body
	var req dto.CreateBranchRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Warn("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	// 4. Obtener person_id del usuario
	creatorPersonID, err := uuid.Parse(userClaims.PersonID)
	if err != nil {
		c.logger.Error("Invalid person ID in claims", zap.Error(err), zap.String("person_id", userClaims.PersonID))
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "invalid_person",
			Message: "Invalid person identifier",
		})
		return
	}

	// 5. Crear sucursal usando el servicio
	branch, err := c.branchService.CreateBranch(ctx.Request.Context(), orgID, creatorPersonID, req)
	if err != nil {
		c.handleServiceError(ctx, err, "Failed to create branch")
		return
	}

	// 6. Mapear a response DTO
	response := c.mapBranchToResponse(branch)

	c.logger.Info("Branch created successfully",
		zap.String("branch_id", branch.ID.String()),
		zap.String("org_id", orgID.String()),
	)

	ctx.JSON(http.StatusCreated, response)
}

// ListBranches lista todas las sucursales de una organización
// @Summary Listar sucursales
// @Description Lista todas las sucursales de una organización con paginación
// @Tags branches
// @Accept json
// @Produce json
// @Param org_id path string true "Organization ID"
// @Param page query int false "Número de página (default: 1)"
// @Param per_page query int false "Items por página (default: 10)"
// @Success 200 {object} dto.BranchListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /organizations/{org_id}/branches [get]
func (c *BranchController) ListBranches(ctx *gin.Context) {
	c.logger.Info("Listing branches", zap.String("method", "GET"), zap.String("path", ctx.Request.URL.Path))

	// 1. Parsear org_id de la URL
	orgIDStr := ctx.Param("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.logger.Warn("Invalid organization ID", zap.Error(err), zap.String("org_id", orgIDStr))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_org_id",
			Message: "Invalid organization identifier",
		})
		return
	}

	// 2. Obtener información del usuario autenticado para permisos
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

	// 3. Parsear filtros de la query string
	filters, err := c.parseBranchFilters(ctx)
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

	// 4. Obtener sucursales usando el servicio
	branches, total, err := c.branchService.ListBranches(ctx.Request.Context(), requesterPersonID, *filters)
	if err != nil {
		c.handleServiceError(ctx, err, "Failed to list branches")
		return
	}

	// 5. Mapear a response DTOs
	branchResponses := make([]dto.OrganizationBranchResponse, len(branches))
	for i, branch := range branches {
		branchResponses[i] = c.mapBranchToResponse(&branch)
	}

	response := dto.BranchListResponse{
		Data: branchResponses,
		Pagination: dto.PaginationResponse{
			Page:       filters.Page,
			PerPage:    filters.PerPage,
			Total:      total,
			TotalPages: int((total + int64(filters.PerPage) - 1) / int64(filters.PerPage)),
			HasNext:    filters.Page < int((total+int64(filters.PerPage)-1)/int64(filters.PerPage)),
			HasPrev:    filters.Page > 1,
		},
	}

	ctx.JSON(http.StatusOK, response)
}

// GetBranch obtiene una sucursal específica
// @Summary Obtener sucursal
// @Description Obtiene los detalles de una sucursal específica
// @Tags branches
// @Accept json
// @Produce json
// @Param org_id path string true "Organization ID"
// @Param branch_id path string true "Branch ID"
// @Success 200 {object} dto.OrganizationBranchResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /organizations/{org_id}/branches/{branch_id} [get]
func (c *BranchController) GetBranch(ctx *gin.Context) {
	c.logger.Info("Getting branch", zap.String("method", "GET"), zap.String("path", ctx.Request.URL.Path))

	// 1. Parsear branch_id de la URL
	branchIDStr := ctx.Param("branch_id")
	branchID, err := uuid.Parse(branchIDStr)
	if err != nil {
		c.logger.Warn("Invalid branch ID", zap.Error(err), zap.String("branch_id", branchIDStr))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_branch_id",
			Message: "Invalid branch identifier",
		})
		return
	}

	// 2. Obtener información del usuario autenticado para permisos
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

	// 3. Obtener sucursal usando el servicio
	branch, err := c.branchService.GetBranch(ctx.Request.Context(), branchID, requesterPersonID)
	if err != nil {
		c.handleServiceError(ctx, err, "Failed to get branch")
		return
	}

	// 4. Mapear a response DTO
	response := c.mapBranchToResponse(branch)

	ctx.JSON(http.StatusOK, response)
}

// UpdateBranch actualiza una sucursal existente
// @Summary Actualizar sucursal
// @Description Actualiza los datos de una sucursal existente
// @Tags branches
// @Accept json
// @Produce json
// @Param org_id path string true "Organization ID"
// @Param branch_id path string true "Branch ID"
// @Param request body dto.UpdateBranchRequest true "Datos actualizados de la sucursal"
// @Success 200 {object} dto.OrganizationBranchResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 422 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /organizations/{org_id}/branches/{branch_id} [put]
func (c *BranchController) UpdateBranch(ctx *gin.Context) {
	c.logger.Info("Updating branch", zap.String("method", "PUT"), zap.String("path", ctx.Request.URL.Path))

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

	// 2. Parsear branch_id de la URL
	branchIDStr := ctx.Param("branch_id")
	branchID, err := uuid.Parse(branchIDStr)
	if err != nil {
		c.logger.Warn("Invalid branch ID", zap.Error(err), zap.String("branch_id", branchIDStr))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_branch_id",
			Message: "Invalid branch identifier",
		})
		return
	}

	// 3. Parsear y validar request body
	var req dto.UpdateBranchRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Warn("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	// 4. Obtener person_id del usuario
	updaterPersonID, err := uuid.Parse(userClaims.PersonID)
	if err != nil {
		c.logger.Error("Invalid person ID in claims", zap.Error(err), zap.String("person_id", userClaims.PersonID))
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "invalid_person",
			Message: "Invalid person identifier",
		})
		return
	}

	// 5. Actualizar sucursal usando el servicio
	branch, err := c.branchService.UpdateBranch(ctx.Request.Context(), branchID, updaterPersonID, req)
	if err != nil {
		c.handleServiceError(ctx, err, "Failed to update branch")
		return
	}

	// 6. Mapear a response DTO
	response := c.mapBranchToResponse(branch)

	c.logger.Info("Branch updated successfully", zap.String("branch_id", branchID.String()))

	ctx.JSON(http.StatusOK, response)
}

// DeleteBranch elimina (soft delete) una sucursal
// @Summary Eliminar sucursal
// @Description Elimina (soft delete) una sucursal existente
// @Tags branches
// @Accept json
// @Produce json
// @Param org_id path string true "Organization ID"
// @Param branch_id path string true "Branch ID"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 422 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /organizations/{org_id}/branches/{branch_id} [delete]
func (c *BranchController) DeleteBranch(ctx *gin.Context) {
	c.logger.Info("Deleting branch", zap.String("method", "DELETE"), zap.String("path", ctx.Request.URL.Path))

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

	// 2. Parsear branch_id de la URL
	branchIDStr := ctx.Param("branch_id")
	branchID, err := uuid.Parse(branchIDStr)
	if err != nil {
		c.logger.Warn("Invalid branch ID", zap.Error(err), zap.String("branch_id", branchIDStr))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_branch_id",
			Message: "Invalid branch identifier",
		})
		return
	}

	// 3. Obtener person_id del usuario
	deleterPersonID, err := uuid.Parse(userClaims.PersonID)
	if err != nil {
		c.logger.Error("Invalid person ID in claims", zap.Error(err), zap.String("person_id", userClaims.PersonID))
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "invalid_person",
			Message: "Invalid person identifier",
		})
		return
	}

	// 4. Eliminar sucursal usando el servicio
	err = c.branchService.DeleteBranch(ctx.Request.Context(), branchID, deleterPersonID)
	if err != nil {
		c.handleServiceError(ctx, err, "Failed to delete branch")
		return
	}

	c.logger.Info("Branch deleted successfully", zap.String("branch_id", branchID.String()))

	ctx.Status(http.StatusNoContent)
}

// ===============================
// HELPER METHODS
// ===============================

// handleServiceError maneja errores del servicio y los convierte a respuestas HTTP apropiadas
func (c *BranchController) handleServiceError(ctx *gin.Context, err error, logMessage string) {
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

// mapBranchToResponse convierte un modelo OrganizationBranch a DTO response
func (c *BranchController) mapBranchToResponse(branch *models.OrganizationBranch) dto.OrganizationBranchResponse {
	return dto.OrganizationBranchResponse{
		BaseResponse: dto.BaseResponse{
			ID:        branch.ID,
			CreatedAt: branch.CreatedAt,
			UpdatedAt: branch.UpdatedAt,
			CreatedBy: branch.CreatedBy,
			UpdatedBy: branch.UpdatedBy,
		},
		OrganizationID: branch.OrganizationID,
		DisplayName:    branch.DisplayName,
		AddressID:      branch.AddressID,
		Phone:          branch.Phone,
		Email:          branch.Email,
		IsMain:         branch.IsMain,
	}
}

// parseBranchFilters parsea los filtros de la query string para branches
func (c *BranchController) parseBranchFilters(ctx *gin.Context) (*dto.BranchFiltersRequest, error) {
	filters := &dto.BranchFiltersRequest{
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

	// IsMain
	if isMainStr := ctx.Query("is_main"); isMainStr != "" {
		if isMain, err := strconv.ParseBool(isMainStr); err == nil {
			filters.IsMain = &isMain
		}
	}

	// HasAddress
	if hasAddressStr := ctx.Query("has_address"); hasAddressStr != "" {
		if hasAddress, err := strconv.ParseBool(hasAddressStr); err == nil {
			filters.HasAddress = &hasAddress
		}
	}

	// Search
	if search := ctx.Query("search"); search != "" {
		filters.Search = search
	}

	return filters, nil
}
