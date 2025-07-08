package dto

// ResetPasswordRequest representa la petición de reset de contraseña
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ResetPasswordResponse representa la respuesta de reset de contraseña
type ResetPasswordResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}
