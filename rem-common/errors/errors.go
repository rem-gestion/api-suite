package errors

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
