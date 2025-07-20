package controllers

import (
	"go.uber.org/zap"
)

// IntegrationController maneja las peticiones HTTP para integraciones
type IntegrationController struct {
	logger *zap.Logger
}

// NewIntegrationController crea una nueva instancia del controlador de integraciones
func NewIntegrationController(logger *zap.Logger) *IntegrationController {
	return &IntegrationController{
		logger: logger,
	}
}

// TODO: Implementar métodos del controlador de integraciones
// CreateIntegration, GetIntegration, ListIntegrations, UpdateIntegration, etc.
