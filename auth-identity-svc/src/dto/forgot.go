package dto

// ForgotPasswordRequest representa la petición de recuperación de contraseña
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ForgotPasswordResponse representa la respuesta de recuperación de contraseña
type ForgotPasswordResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}
