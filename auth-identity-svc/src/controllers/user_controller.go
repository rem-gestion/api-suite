package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rem-gestion/api-suite/auth-identity/src/dto"
	"github.com/rem-gestion/api-suite/auth-identity/src/models"
	"github.com/rem-gestion/api-suite/auth-identity/src/repository"
	"github.com/rem-gestion/api-suite/auth-identity/src/services"
	rerrors "github.com/rem-gestion/rem-common/errors"
)

// UserController maneja los endpoints relacionados con usuarios
type UserController struct {
	userService *services.UserService
}

// NewUserController crea una nueva instancia del controlador
func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update current user's profile information
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateProfileRequest true "Profile update data"
// @Success 200 {object} dto.UserProfileResponse
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/profile [put]
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.Error(&rerrors.UnauthorizedError{Msg: "user not authenticated"})
		return
	}

	var req dto.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(&rerrors.BadRequestError{Msg: err.Error()})
		return
	}

	profile, err := c.userService.UpdateProfile(userID.(string), req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, profile)
}

// GetUserByEmail godoc
// @Summary Get user by email
// @Description Get user information by email (admin only)
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param email path string true "User email"
// @Success 200 {object} dto.UserProfileResponse
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/email/{email} [get]
func (c *UserController) GetUserByEmail(ctx *gin.Context) {
	email := ctx.Param("email")
	if email == "" {
		ctx.Error(&rerrors.BadRequestError{Msg: "email parameter is required"})
		return
	}

	profile, err := c.userService.GetUserByEmail(email)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, profile)
}

// ListUsers godoc
// @Summary List users
// @Description Get paginated list of users with optional filters (admin only)
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(20)
// @Param search query string false "Search in email"
// @Param onboard_status query string false "Filter by onboard status" Enums(new, in_progress, done)
// @Param account_status query string false "Filter by account status" Enums(pending, active, suspended)
// @Param has_person query bool false "Filter by person association"
// @Success 200 {object} dto.ListUsersResponse
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users [get]
func (c *UserController) ListUsers(ctx *gin.Context) {
	// Parsear parámetros de paginación
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "20"))

	// Construir filtros
	filter := &repository.UserFilter{}

	// Filtro de búsqueda
	if search := ctx.Query("search"); search != "" {
		filter.Search = search
	}

	// Filtro de estado de onboarding
	if onboardStatusStr := ctx.Query("onboard_status"); onboardStatusStr != "" {
		onboardStatus := models.OnboardStatus(onboardStatusStr)
		filter.OnboardStatus = &onboardStatus
	}

	// Filtro de estado de cuenta
	if accountStatusStr := ctx.Query("account_status"); accountStatusStr != "" {
		accountStatus := models.AccountStatus(accountStatusStr)
		filter.AccountStatus = &accountStatus
	}

	// Filtro de asociación con persona
	if hasPersonStr := ctx.Query("has_person"); hasPersonStr != "" {
		if hasPersonStr == "true" {
			hasPerson := true
			filter.HasPerson = &hasPerson
		} else if hasPersonStr == "false" {
			hasPerson := false
			filter.HasPerson = &hasPerson
		}
	}

	response, err := c.userService.ListUsers(page, perPage, filter)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// DeactivateUser godoc
// @Summary Deactivate user
// @Description Deactivate a user account (admin only)
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{} "User deactivated successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/{id}/deactivate [post]
func (c *UserController) DeactivateUser(ctx *gin.Context) {
	userID := ctx.Param("id")
	if userID == "" {
		ctx.Error(&rerrors.BadRequestError{Msg: "user ID parameter is required"})
		return
	}

	// Obtener ID del usuario que realiza la acción (para auditoría)
	updatedBy, _ := ctx.Get("userID")
	updatedByStr := ""
	if updatedBy != nil {
		updatedByStr = updatedBy.(string)
	}

	err := c.userService.DeactivateUser(userID, updatedByStr)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User deactivated successfully",
		"success": true,
	})
}
