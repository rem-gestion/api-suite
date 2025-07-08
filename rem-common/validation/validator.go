package validation

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidationError representa un error de validación
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", v.Field, v.Message)
}

// ValidationErrors es una colección de errores de validación
type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}
	var messages []string
	for _, err := range v {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

func (v ValidationErrors) HasErrors() bool {
	return len(v) > 0
}

// Validator es la interfaz base para validadores
type Validator interface {
	Validate(interface{}) ValidationErrors
}

// BaseValidator contiene validaciones comunes
type BaseValidator struct{}

// ValidateRequired valida que un campo no esté vacío
func (v *BaseValidator) ValidateRequired(value, fieldName string) *ValidationError {
	if strings.TrimSpace(value) == "" {
		return &ValidationError{
			Field:   fieldName,
			Message: "field is required",
			Code:    "REQUIRED",
		}
	}
	return nil
}

// ValidateEmail valida formato de email
func (v *BaseValidator) ValidateEmail(email, fieldName string) *ValidationError {
	if email == "" {
		return nil // Optional field
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return &ValidationError{
			Field:   fieldName,
			Message: "invalid email format",
			Code:    "INVALID_EMAIL",
		}
	}
	return nil
}

// ValidatePhone valida formato de teléfono
func (v *BaseValidator) ValidatePhone(phone, fieldName string) *ValidationError {
	if phone == "" {
		return nil // Optional field
	}

	// Permitir formatos: +54911234567, 011234567, etc.
	phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	cleanPhone := regexp.MustCompile(`[^\d+]`).ReplaceAllString(phone, "")

	if !phoneRegex.MatchString(cleanPhone) {
		return &ValidationError{
			Field:   fieldName,
			Message: "invalid phone format",
			Code:    "INVALID_PHONE",
		}
	}
	return nil
}

// ValidateDNI valida formato de DNI argentino
func (v *BaseValidator) ValidateDNI(dni, fieldName string) *ValidationError {
	if dni == "" {
		return nil // Optional field
	}

	// DNI debe tener exactamente 8 dígitos
	dniRegex := regexp.MustCompile(`^\d{8}$`)
	if !dniRegex.MatchString(dni) {
		return &ValidationError{
			Field:   fieldName,
			Message: "DNI must be exactly 8 digits",
			Code:    "INVALID_DNI",
		}
	}
	return nil
}

// ValidateCUIT valida formato de CUIT argentino
func (v *BaseValidator) ValidateCUIT(cuit, fieldName string) *ValidationError {
	if cuit == "" {
		return nil // Optional field
	}

	// CUIT debe tener exactamente 11 dígitos
	cuitRegex := regexp.MustCompile(`^\d{11}$`)
	if !cuitRegex.MatchString(cuit) {
		return &ValidationError{
			Field:   fieldName,
			Message: "CUIT must be exactly 11 digits",
			Code:    "INVALID_CUIT",
		}
	}
	return nil
}

// ValidateMaxLength valida longitud máxima
func (v *BaseValidator) ValidateMaxLength(value, fieldName string, maxLength int) *ValidationError {
	if len(value) > maxLength {
		return &ValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("field exceeds maximum length of %d characters", maxLength),
			Code:    "MAX_LENGTH_EXCEEDED",
		}
	}
	return nil
}

// ValidateEnum valida que el valor esté en una lista de valores permitidos
func (v *BaseValidator) ValidateEnum(value, fieldName string, allowedValues []string) *ValidationError {
	for _, allowed := range allowedValues {
		if value == allowed {
			return nil
		}
	}
	return &ValidationError{
		Field:   fieldName,
		Message: fmt.Sprintf("field must be one of: %s", strings.Join(allowedValues, ", ")),
		Code:    "INVALID_ENUM_VALUE",
	}
}
