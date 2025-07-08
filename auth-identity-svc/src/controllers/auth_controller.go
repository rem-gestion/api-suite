package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rem-gestion/api-suite/auth-identity/src/dto"
	"github.com/rem-gestion/api-suite/auth-identity/src/services"
	rerrors "github.com/rem-gestion/rem-common/errors"
)

// AuthController maneja los endpoints de autenticación
type AuthController struct {
	authService *services.AuthService
}

// NewAuthController crea una nueva instancia del controlador
func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email and password, optionally with person data
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration data"
// @Success 201 {object} dto.RegisterResponse
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 409 {object} map[string]interface{} "Email already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(&rerrors.BadRequestError{Msg: err.Error()})
		return
	}

	response, err := c.authService.Register(req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

// Login godoc
// @Summary User login
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Invalid credentials"
// @Failure 403 {object} map[string]interface{} "Account suspended or pending"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(&rerrors.BadRequestError{Msg: err.Error()})
		return
	}

	response, err := c.authService.Login(req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get a new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} dto.RefreshTokenResponse
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Invalid refresh token"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/refresh [post]
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(&rerrors.BadRequestError{Msg: err.Error()})
		return
	}

	response, err := c.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// ForgotPassword godoc
// @Summary Request password reset
// @Description Send password reset email to user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.ForgotPasswordRequest true "Email for password reset"
// @Success 200 {object} dto.ForgotPasswordResponse
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/forgot-password [post]
func (c *AuthController) ForgotPassword(ctx *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(&rerrors.BadRequestError{Msg: err.Error()})
		return
	}

	response, err := c.authService.ForgotPassword(req.Email)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// ResetPassword godoc
// @Summary Reset password
// @Description Reset password using reset token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.ResetPasswordRequest true "Reset password data"
// @Success 200 {object} dto.ResetPasswordResponse
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Invalid reset token"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/reset-password [post]
func (c *AuthController) ResetPassword(ctx *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(&rerrors.BadRequestError{Msg: err.Error()})
		return
	}

	response, err := c.authService.ResetPassword(req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get current user's profile information
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.UserProfileResponse
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/me [get]
func (c *AuthController) GetProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.Error(&rerrors.UnauthorizedError{Msg: "user not authenticated"})
		return
	}

	profile, err := c.authService.GetUserProfile(userID.(string))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, profile)
}

// ValidateToken godoc
// @Summary Validate access token
// @Description Validate if the provided access token is valid
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Token is valid"
// @Failure 401 {object} map[string]interface{} "Invalid token"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/validate [get]
func (c *AuthController) ValidateToken(ctx *gin.Context) {
	// Obtener token del header Authorization
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.Error(&rerrors.UnauthorizedError{Msg: "authorization header missing"})
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		ctx.Error(&rerrors.UnauthorizedError{Msg: "invalid authorization header format"})
		return
	}

	claims, err := c.authService.ValidateToken(parts[1])
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"valid":      true,
		"user_id":    claims.UserID,
		"role":       claims.Role,
		"expires_at": claims.ExpiresAt.Unix(),
	})
}

// Logout godoc
// @Summary User logout
// @Description Logout user (client should discard tokens)
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Logout successful"
// @Router /auth/logout [post]
func (c *AuthController) Logout(ctx *gin.Context) {
	// En una implementación completa, aquí se podría:
	// 1. Invalidar el refresh token en Redis/BD
	// 2. Agregar el access token a una blacklist
	// 3. Registrar el evento de logout

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Logout successful. Please discard your tokens.",
		"success": true,
	})
}
