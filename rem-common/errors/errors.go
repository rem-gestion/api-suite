package errors

import "fmt"

// ValidationError para errores de validación de campos
type ValidationError struct {
	Msg    string
	Fields map[string]string // campo -> mensaje
}

func (e *ValidationError) Error() string { return e.Msg }

// BadRequestError para 400
type BadRequestError struct{ Msg string }

func (e *BadRequestError) Error() string { return e.Msg }

// UnauthorizedError para 401
type UnauthorizedError struct{ Msg string }

func (e *UnauthorizedError) Error() string { return e.Msg }

// ForbiddenError para 403
type ForbiddenError struct{ Msg string }

func (e *ForbiddenError) Error() string { return e.Msg }

// NotFoundError para 404
type NotFoundError struct{ Msg string }

func (e *NotFoundError) Error() string { return e.Msg }

// ConflictError para 409
type ConflictError struct{ Msg string }

func (e *ConflictError) Error() string { return e.Msg }

// TooManyRequestsError para 429
type TooManyRequestsError struct{ Msg string }

func (e *TooManyRequestsError) Error() string { return e.Msg }

// InternalServerError para 500
type InternalServerError struct{ Msg string }

func (e *InternalServerError) Error() string { return e.Msg }

// DatabaseError para errores de base de datos
type DatabaseError struct{ Msg string }

func (e *DatabaseError) Error() string { return e.Msg }

// BusinessRuleError para errores de reglas de negocio
type BusinessRuleError struct{ Msg string }

func (e *BusinessRuleError) Error() string { return e.Msg }

// Funciones de creación de errores

// NewValidationError crea un nuevo error de validación
func NewValidationError(field, message string) error {
	return &ValidationError{
		Msg:    message,
		Fields: map[string]string{field: message},
	}
}

// NewBadRequestError crea un nuevo error de solicitud incorrecta
func NewBadRequestError(message string) error {
	return &BadRequestError{Msg: message}
}

// NewUnauthorizedError crea un nuevo error no autorizado
func NewUnauthorizedError(message string) error {
	return &UnauthorizedError{Msg: message}
}

// NewForbiddenError crea un nuevo error prohibido
func NewForbiddenError(message string) error {
	return &ForbiddenError{Msg: message}
}

// NewNotFoundError crea un nuevo error no encontrado
func NewNotFoundError(resource, id string) error {
	return &NotFoundError{Msg: fmt.Sprintf("%s with id '%s' not found", resource, id)}
}

// NewConflictError crea un nuevo error de conflicto
func NewConflictError(message string) error {
	return &ConflictError{Msg: message}
}

// NewTooManyRequestsError crea un nuevo error de demasiadas solicitudes
func NewTooManyRequestsError(message string) error {
	return &TooManyRequestsError{Msg: message}
}

// NewInternalServerError crea un nuevo error interno del servidor
func NewInternalServerError(message string) error {
	return &InternalServerError{Msg: message}
}

// NewDatabaseError crea un nuevo error de base de datos
func NewDatabaseError(message string, err error) error {
	if err != nil {
		return &DatabaseError{Msg: fmt.Sprintf("%s: %v", message, err)}
	}
	return &DatabaseError{Msg: message}
}

// NewBusinessRuleError crea un nuevo error de regla de negocio
func NewBusinessRuleError(rule, message string) error {
	return &BusinessRuleError{Msg: fmt.Sprintf("Business rule violation '%s': %s", rule, message)}
}
