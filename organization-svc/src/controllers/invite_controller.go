package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/rem-gestion/api-suite/organization/src/dto"
	"github.com/rem-gestion/api-suite/organization/src/services"
	rcerrors "github.com/rem-gestion/rem-common/errors"
	"github.com/rem-gestion/rem-common/middleware"
)

// InviteController maneja las peticiones HTTP para invitaciones
type InviteController struct {
	inviteService *services.InvitationService
	logger        *zap.Logger
}

// NewInviteController crea una nueva instancia del controlador de invitaciones
func NewInviteController(inviteService *services.InvitationService, logger *zap.Logger) *InviteController {
	return &InviteController{
		inviteService: inviteService,
		logger:        logger,
	}
}

// ===============================
// INVITATION OPERATIONS
// ===============================

// SendInvitation envía una nueva invitación a un usuario
// @Summary Enviar invitación
// @Description Envía una invitación para unirse a la organización
// @Tags invitations
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param request body dto.CreateInvitationRequest true "Datos de la invitación"
// @Success 201 {object} dto.OrganizationInviteResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 422 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /organizations/{orgId}/invitations [post]
func (c *InviteController) SendInvitation(ctx *gin.Context) {
	c.logger.Info("Sending invitation", zap.String("method", "POST"), zap.String("path", ctx.Request.URL.Path))

	// 1. Validar y parsear Organization ID
	orgIDStr := ctx.Param("orgId")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.logger.Warn("Invalid organization ID", zap.String("orgId", orgIDStr), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_org_id",
			Message: "Invalid organization ID format",
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

	inviterPersonID, err := uuid.Parse(userClaims.UserID)
	if err != nil {
		c.logger.Error("Invalid user ID in claims", zap.Error(err), zap.String("user_id", userClaims.UserID))
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "invalid_user",
			Message: "Invalid user identifier",
		})
		return
	}

	// 3. Parsear y validar request body
	var req dto.CreateInvitationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Warn("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	// 4. Llamar al servicio
	invitation, err := c.inviteService.SendInvitation(ctx.Request.Context(), orgID, inviterPersonID, req)
	if err != nil {
		c.logger.Error("Failed to send invitation",
			zap.String("orgId", orgID.String()),
			zap.String("inviterPersonId", inviterPersonID.String()),
			zap.String("inviteeEmail", req.InviteeEmail),
			zap.Error(err))

		switch err.(type) {
		case *rcerrors.ValidationError:
			validationErr := err.(*rcerrors.ValidationError)
			ctx.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
				Error:   "validation_failed",
				Message: validationErr.Msg,
				Details: map[string]interface{}{
					"fields": validationErr.Fields,
				},
			})
		case *rcerrors.NotFoundError:
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: err.Error(),
			})
		case *rcerrors.ForbiddenError:
			ctx.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "forbidden",
				Message: err.Error(),
			})
		case *rcerrors.ConflictError:
			ctx.JSON(http.StatusConflict, dto.ErrorResponse{
				Error:   "conflict",
				Message: err.Error(),
			})
		default:
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: "Failed to send invitation",
			})
		}
		return
	}

	// 5. Convertir a DTO de respuesta
	response := dto.OrganizationInviteResponse{
		BaseResponse: dto.BaseResponse{
			ID:        invitation.ID,
			CreatedAt: invitation.CreatedAt,
			UpdatedAt: invitation.UpdatedAt,
		},
		OrganizationID:  invitation.OrganizationID,
		InviterPersonID: invitation.InviterPersonID,
		InviteeEmail:    invitation.InviteeEmail,
		InviteePersonID: invitation.InviteePersonID,
		RoleID:          invitation.RoleID,
		BranchID:        invitation.BranchID,
		Status:          string(invitation.Status),
		Token:           invitation.Token,
		ExpiresAt:       invitation.ExpiresAt,
		AcceptedAt:      invitation.AcceptedAt,
		RejectedAt:      invitation.RejectedAt,
		CancelledAt:     invitation.CancelledAt,
	}

	c.logger.Info("Invitation sent successfully",
		zap.String("invitationId", invitation.ID.String()),
		zap.String("orgId", orgID.String()),
		zap.String("inviteeEmail", req.InviteeEmail))

	ctx.JSON(http.StatusCreated, response)
}

// GetInvitation obtiene una invitación específica
// @Summary Obtener invitación
// @Description Obtiene los detalles de una invitación específica
// @Tags invitations
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param invitationId path string true "Invitation ID"
// @Success 200 {object} dto.OrganizationInviteResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /organizations/{orgId}/invitations/{invitationId} [get]
func (c *InviteController) GetInvitation(ctx *gin.Context) {
	c.logger.Info("Getting invitation", zap.String("method", "GET"), zap.String("path", ctx.Request.URL.Path))

	// 1. Validar y parsear Invitation ID
	invitationIDStr := ctx.Param("invitationId")
	invitationID, err := uuid.Parse(invitationIDStr)
	if err != nil {
		c.logger.Warn("Invalid invitation ID", zap.String("invitationId", invitationIDStr), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_invitation_id",
			Message: "Invalid invitation ID format",
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

	requesterPersonID, err := uuid.Parse(userClaims.UserID)
	if err != nil {
		c.logger.Error("Invalid user ID in claims", zap.Error(err), zap.String("user_id", userClaims.UserID))
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "invalid_user",
			Message: "Invalid user identifier",
		})
		return
	}

	// 3. Llamar al servicio
	invitation, err := c.inviteService.GetInvitation(ctx.Request.Context(), invitationID, requesterPersonID)
	if err != nil {
		c.logger.Error("Failed to get invitation",
			zap.String("invitationId", invitationID.String()),
			zap.String("requesterPersonId", requesterPersonID.String()),
			zap.Error(err))

		switch err.(type) {
		case *rcerrors.NotFoundError:
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: err.Error(),
			})
		case *rcerrors.ForbiddenError:
			ctx.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "forbidden",
				Message: err.Error(),
			})
		default:
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: "Failed to get invitation",
			})
		}
		return
	}

	// 4. Convertir a DTO de respuesta
	response := dto.OrganizationInviteResponse{
		BaseResponse: dto.BaseResponse{
			ID:        invitation.ID,
			CreatedAt: invitation.CreatedAt,
			UpdatedAt: invitation.UpdatedAt,
		},
		OrganizationID:  invitation.OrganizationID,
		InviterPersonID: invitation.InviterPersonID,
		InviteeEmail:    invitation.InviteeEmail,
		InviteePersonID: invitation.InviteePersonID,
		RoleID:          invitation.RoleID,
		BranchID:        invitation.BranchID,
		Status:          string(invitation.Status),
		Token:           invitation.Token,
		ExpiresAt:       invitation.ExpiresAt,
		AcceptedAt:      invitation.AcceptedAt,
		RejectedAt:      invitation.RejectedAt,
		CancelledAt:     invitation.CancelledAt,
	}

	c.logger.Info("Invitation retrieved successfully",
		zap.String("invitationId", invitation.ID.String()))

	ctx.JSON(http.StatusOK, response)
}

// ListInvitations lista todas las invitaciones con filtros
// @Summary Listar invitaciones
// @Description Lista todas las invitaciones de la organización con filtros opcionales
// @Tags invitations
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param page query int false "Número de página" default(1)
// @Param per_page query int false "Límite por página" default(10)
// @Param status query string false "Filtrar por estado" Enums(PENDING, ACCEPTED, DECLINED, EXPIRED, CANCELLED)
// @Param invitee_email query string false "Filtrar por email del invitado"
// @Param role_id query string false "Filtrar por ID de rol"
// @Param branch_id query string false "Filtrar por ID de sucursal"
// @Success 200 {object} dto.InvitationListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /organizations/{orgId}/invitations [get]
func (c *InviteController) ListInvitations(ctx *gin.Context) {
	c.logger.Info("Listing invitations", zap.String("method", "GET"), zap.String("path", ctx.Request.URL.Path))

	// 1. Validar y parsear Organization ID
	orgIDStr := ctx.Param("orgId")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.logger.Warn("Invalid organization ID", zap.String("orgId", orgIDStr), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_org_id",
			Message: "Invalid organization ID format",
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

	requesterPersonID, err := uuid.Parse(userClaims.UserID)
	if err != nil {
		c.logger.Error("Invalid user ID in claims", zap.Error(err), zap.String("user_id", userClaims.UserID))
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "invalid_user",
			Message: "Invalid user identifier",
		})
		return
	}

	// 3. Parsear filtros de query parameters
	filters := dto.InvitationFiltersRequest{
		OrganizationID: orgID,
	}

	// Pagination - defaults
	filters.Page = 1
	filters.PerPage = 10

	if pageStr := ctx.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			filters.Page = page
		}
	}
	if perPageStr := ctx.Query("per_page"); perPageStr != "" {
		if perPage, err := strconv.Atoi(perPageStr); err == nil && perPage > 0 {
			filters.PerPage = perPage
		}
	}

	// Status filter
	if status := ctx.Query("status"); status != "" {
		filters.Status = status
	}

	// Email filter
	if email := ctx.Query("invitee_email"); email != "" {
		filters.InviteeEmail = &email
	}

	// Role ID filter
	if roleIDStr := ctx.Query("role_id"); roleIDStr != "" {
		if roleID, err := uuid.Parse(roleIDStr); err == nil {
			filters.RoleID = &roleID
		}
	}

	// Branch ID filter
	if branchIDStr := ctx.Query("branch_id"); branchIDStr != "" {
		if branchID, err := uuid.Parse(branchIDStr); err == nil {
			filters.BranchID = &branchID
		}
	}

	// 4. Llamar al servicio
	invitations, total, err := c.inviteService.ListInvitations(ctx.Request.Context(), requesterPersonID, filters)
	if err != nil {
		c.logger.Error("Failed to list invitations",
			zap.String("orgId", orgID.String()),
			zap.String("requesterPersonId", requesterPersonID.String()),
			zap.Error(err))

		switch err.(type) {
		case *rcerrors.ForbiddenError:
			ctx.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "forbidden",
				Message: err.Error(),
			})
		default:
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: "Failed to list invitations",
			})
		}
		return
	}

	// 5. Convertir a DTOs de respuesta
	var invitationDTOs []dto.OrganizationInviteResponse
	for _, invitation := range invitations {
		invitationDTOs = append(invitationDTOs, dto.OrganizationInviteResponse{
			BaseResponse: dto.BaseResponse{
				ID:        invitation.ID,
				CreatedAt: invitation.CreatedAt,
				UpdatedAt: invitation.UpdatedAt,
			},
			OrganizationID:  invitation.OrganizationID,
			InviterPersonID: invitation.InviterPersonID,
			InviteeEmail:    invitation.InviteeEmail,
			InviteePersonID: invitation.InviteePersonID,
			RoleID:          invitation.RoleID,
			BranchID:        invitation.BranchID,
			Status:          string(invitation.Status),
			Token:           invitation.Token,
			ExpiresAt:       invitation.ExpiresAt,
			AcceptedAt:      invitation.AcceptedAt,
			RejectedAt:      invitation.RejectedAt,
			CancelledAt:     invitation.CancelledAt,
		})
	}

	// 6. Crear respuesta paginada
	totalPages := int((total + int64(filters.PerPage) - 1) / int64(filters.PerPage))
	response := dto.InvitationListResponse{
		Data: invitationDTOs,
		Pagination: dto.PaginationResponse{
			Page:       filters.Page,
			PerPage:    filters.PerPage,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    filters.Page < totalPages,
			HasPrev:    filters.Page > 1,
		},
	}

	c.logger.Info("Invitations listed successfully",
		zap.String("orgId", orgID.String()),
		zap.Int("count", len(invitationDTOs)),
		zap.Int64("total", total))

	ctx.JSON(http.StatusOK, response)
}

// ResendInvitation reenvía una invitación existente
// @Summary Reenviar invitación
// @Description Reenvía una invitación existente con nueva fecha de expiración
// @Tags invitations
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param invitationId path string true "Invitation ID"
// @Param request body dto.ResendInvitationRequest true "Datos para reenvío"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /organizations/{orgId}/invitations/{invitationId}/resend [post]
func (c *InviteController) ResendInvitation(ctx *gin.Context) {
	c.logger.Info("Resending invitation", zap.String("method", "POST"), zap.String("path", ctx.Request.URL.Path))

	// 1. Validar y parsear Invitation ID
	invitationIDStr := ctx.Param("invitationId")
	invitationID, err := uuid.Parse(invitationIDStr)
	if err != nil {
		c.logger.Warn("Invalid invitation ID", zap.String("invitationId", invitationIDStr), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_invitation_id",
			Message: "Invalid invitation ID format",
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

	requesterPersonID, err := uuid.Parse(userClaims.UserID)
	if err != nil {
		c.logger.Error("Invalid user ID in claims", zap.Error(err), zap.String("user_id", userClaims.UserID))
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "invalid_user",
			Message: "Invalid user identifier",
		})
		return
	}

	// 3. Parsear y validar request body
	var req dto.ResendInvitationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Warn("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	// 4. Llamar al servicio
	err = c.inviteService.ResendInvitation(ctx.Request.Context(), invitationID, requesterPersonID, req)
	if err != nil {
		c.logger.Error("Failed to resend invitation",
			zap.String("invitationId", invitationID.String()),
			zap.String("requesterPersonId", requesterPersonID.String()),
			zap.Error(err))

		switch err.(type) {
		case *rcerrors.NotFoundError:
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: err.Error(),
			})
		case *rcerrors.ForbiddenError:
			ctx.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "forbidden",
				Message: err.Error(),
			})
		case *rcerrors.ValidationError:
			validationErr := err.(*rcerrors.ValidationError)
			ctx.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
				Error:   "validation_failed",
				Message: validationErr.Msg,
				Details: map[string]interface{}{
					"fields": validationErr.Fields,
				},
			})
		default:
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: "Failed to resend invitation",
			})
		}
		return
	}

	c.logger.Info("Invitation resent successfully",
		zap.String("invitationId", invitationID.String()))

	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Invitation resent successfully",
	})
}

// CancelInvitation cancela una invitación existente
// @Summary Cancelar invitación
// @Description Cancela una invitación existente
// @Tags invitations
// @Accept json
// @Produce json
// @Param orgId path string true "Organization ID"
// @Param invitationId path string true "Invitation ID"
// @Param request body dto.CancelInvitationRequest true "Datos de cancelación"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /organizations/{orgId}/invitations/{invitationId}/cancel [post]
func (c *InviteController) CancelInvitation(ctx *gin.Context) {
	c.logger.Info("Canceling invitation", zap.String("method", "POST"), zap.String("path", ctx.Request.URL.Path))

	// 1. Validar y parsear Invitation ID
	invitationIDStr := ctx.Param("invitationId")
	invitationID, err := uuid.Parse(invitationIDStr)
	if err != nil {
		c.logger.Warn("Invalid invitation ID", zap.String("invitationId", invitationIDStr), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_invitation_id",
			Message: "Invalid invitation ID format",
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

	requesterPersonID, err := uuid.Parse(userClaims.UserID)
	if err != nil {
		c.logger.Error("Invalid user ID in claims", zap.Error(err), zap.String("user_id", userClaims.UserID))
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "invalid_user",
			Message: "Invalid user identifier",
		})
		return
	}

	// 3. Parsear y validar request body
	var req dto.CancelInvitationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Warn("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	// 4. Llamar al servicio
	err = c.inviteService.CancelInvitation(ctx.Request.Context(), invitationID, requesterPersonID, req)
	if err != nil {
		c.logger.Error("Failed to cancel invitation",
			zap.String("invitationId", invitationID.String()),
			zap.String("requesterPersonId", requesterPersonID.String()),
			zap.Error(err))

		switch err.(type) {
		case *rcerrors.NotFoundError:
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: err.Error(),
			})
		case *rcerrors.ForbiddenError:
			ctx.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "forbidden",
				Message: err.Error(),
			})
		case *rcerrors.ValidationError:
			validationErr := err.(*rcerrors.ValidationError)
			ctx.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
				Error:   "validation_failed",
				Message: validationErr.Msg,
				Details: map[string]interface{}{
					"fields": validationErr.Fields,
				},
			})
		default:
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: "Failed to cancel invitation",
			})
		}
		return
	}

	c.logger.Info("Invitation canceled successfully",
		zap.String("invitationId", invitationID.String()))

	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Invitation canceled successfully",
	})
}

// ===============================
// PUBLIC INVITATION ENDPOINTS
// ===============================

// ValidateInvitationToken valida un token de invitación (endpoint público)
// @Summary Validar token de invitación
// @Description Valida un token de invitación y retorna la información de la invitación
// @Tags public-invitations
// @Accept json
// @Produce json
// @Param request body dto.ValidateInvitationRequest true "Token de invitación"
// @Success 200 {object} dto.ValidateInvitationResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 410 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /public/invitations/validate [post]
func (c *InviteController) ValidateInvitationToken(ctx *gin.Context) {
	c.logger.Info("Validating invitation token", zap.String("method", "POST"), zap.String("path", ctx.Request.URL.Path))

	// 1. Parsear y validar request body
	var req dto.ValidateInvitationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Warn("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	// 2. Llamar al servicio
	invitation, err := c.inviteService.ValidateInvitationToken(ctx.Request.Context(), req.Token)
	if err != nil {
		c.logger.Error("Failed to validate invitation token",
			zap.String("token", req.Token[:10]+"..."), // Log only first 10 chars for security
			zap.Error(err))

		switch err.(type) {
		case *rcerrors.NotFoundError:
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: "Invitation not found or invalid token",
			})
		case *rcerrors.ValidationError:
			ctx.JSON(http.StatusGone, dto.ErrorResponse{
				Error:   "expired",
				Message: err.Error(),
			})
		default:
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: "Failed to validate invitation token",
			})
		}
		return
	}

	// 3. Convertir a DTO de respuesta (usar ValidateInvitationResponse para endpoint público)
	response := dto.ValidateInvitationResponse{
		IsValid:      true,
		InvitationID: &invitation.ID,
		// Organization, Role, Branch se cargarían con includes si fuera necesario
	}

	c.logger.Info("Invitation token validated successfully",
		zap.String("invitationId", invitation.ID.String()))

	ctx.JSON(http.StatusOK, response)
}

// AcceptInvitation acepta una invitación (endpoint público)
// @Summary Aceptar invitación
// @Description Acepta una invitación y crea el usuario/empleado correspondiente
// @Tags public-invitations
// @Accept json
// @Produce json
// @Param request body dto.AcceptInvitationRequest true "Datos de aceptación"
// @Success 200 {object} dto.AcceptInvitationResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 410 {object} dto.ErrorResponse
// @Failure 422 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /public/invitations/accept [post]
func (c *InviteController) AcceptInvitation(ctx *gin.Context) {
	c.logger.Info("Accepting invitation", zap.String("method", "POST"), zap.String("path", ctx.Request.URL.Path))

	// 1. Parsear y validar request body
	var req dto.AcceptInvitationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Warn("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	// 2. Llamar al servicio
	employee, err := c.inviteService.AcceptInvitation(ctx.Request.Context(), req.Token, req)
	if err != nil {
		c.logger.Error("Failed to accept invitation",
			zap.String("token", req.Token[:10]+"..."), // Log only first 10 chars for security
			zap.Error(err))

		switch err.(type) {
		case *rcerrors.NotFoundError:
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: "Invitation not found or invalid token",
			})
		case *rcerrors.ValidationError:
			validationErr := err.(*rcerrors.ValidationError)
			// Could be expired invitation or validation error
			if validationErr.Msg == "invitation has expired" {
				ctx.JSON(http.StatusGone, dto.ErrorResponse{
					Error:   "expired",
					Message: validationErr.Msg,
				})
			} else {
				ctx.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
					Error:   "validation_failed",
					Message: validationErr.Msg,
					Details: map[string]interface{}{
						"fields": validationErr.Fields,
					},
				})
			}
		case *rcerrors.ConflictError:
			ctx.JSON(http.StatusConflict, dto.ErrorResponse{
				Error:   "conflict",
				Message: err.Error(),
			})
		default:
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: "Failed to accept invitation",
			})
		}
		return
	}

	// 3. Convertir a DTO de respuesta (usar AcceptInvitationResponse para endpoint público)
	response := dto.AcceptInvitationResponse{
		// Se necesitará verificar qué campos tiene AcceptInvitationResponse
		// Employee: employee,
		// ... otros campos según el DTO
	}

	c.logger.Info("Invitation accepted successfully",
		zap.String("employeeId", employee.ID.String()),
		zap.String("organizationId", employee.OrganizationID.String()))

	ctx.JSON(http.StatusOK, response)
}

// DeclineInvitation rechaza una invitación (endpoint público)
// @Summary Rechazar invitación
// @Description Rechaza una invitación
// @Tags public-invitations
// @Accept json
// @Produce json
// @Param request body dto.DeclineInvitationRequest true "Datos de rechazo"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 410 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /public/invitations/decline [post]
func (c *InviteController) DeclineInvitation(ctx *gin.Context) {
	c.logger.Info("Declining invitation", zap.String("method", "POST"), zap.String("path", ctx.Request.URL.Path))

	// 1. Parsear y validar request body
	var req dto.DeclineInvitationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Warn("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	// 2. Llamar al servicio
	err := c.inviteService.DeclineInvitation(ctx.Request.Context(), req.Token, req)
	if err != nil {
		c.logger.Error("Failed to decline invitation",
			zap.String("token", req.Token[:10]+"..."), // Log only first 10 chars for security
			zap.Error(err))

		switch err.(type) {
		case *rcerrors.NotFoundError:
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: "Invitation not found or invalid token",
			})
		case *rcerrors.ValidationError:
			validationErr := err.(*rcerrors.ValidationError)
			// Could be expired invitation
			if validationErr.Msg == "invitation has expired" {
				ctx.JSON(http.StatusGone, dto.ErrorResponse{
					Error:   "expired",
					Message: validationErr.Msg,
				})
			} else {
				ctx.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
					Error:   "validation_failed",
					Message: validationErr.Msg,
				})
			}
		default:
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: "Failed to decline invitation",
			})
		}
		return
	}

	c.logger.Info("Invitation declined successfully")

	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Invitation declined successfully",
	})
}
