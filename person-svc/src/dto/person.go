package dto

import "github.com/google/uuid"

/* ---------- PERSON ---------- */
type CreatePersonDTO struct {
	Type      string     `json:"type" validate:"required,oneof=individual company"`
	AddressID *uuid.UUID `json:"address_id"`
	AvatarURL *string    `json:"avatar_url,omitempty"`
	Sexo      *string    `json:"sexo,omitempty" validate:"omitempty,oneof=masculino femenino"`

	// If individual
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	DNI       *string `json:"dni,omitempty"`

	// If company
	LegalName   *string `json:"legal_name,omitempty"`
	CUIT        *string `json:"cuit,omitempty"`
	SocietyType *string `json:"society_type,omitempty"`
}

type AddressPayloadDTO struct {
	Floor   *string `json:"floor,omitempty"`
	Unit    *string `json:"unit,omitempty"`
	Street  string  `json:"street"  validate:"required"`
	Number  int     `json:"number"  validate:"required"`
	City    string  `json:"city"    validate:"required"`
	State   string  `json:"state,omitempty"`
	Zip     string  `json:"zip,omitempty"`
	Country string  `json:"country" validate:"required,len=2"`
}

type UpdatePersonDTO struct {
	AvatarURL *string `json:"avatar_url,omitempty"`
	Sexo      *string `json:"sexo,omitempty" validate:"omitempty,oneof=masculino femenino"`

	// Sub-type edits
	FirstName   *string `json:"first_name,omitempty"`
	LastName    *string `json:"last_name,omitempty"`
	DNI         *string `json:"dni,omitempty"`
	LegalName   *string `json:"legal_name,omitempty"`
	CUIT        *string `json:"cuit,omitempty"`
	SocietyType *string `json:"society_type,omitempty"`
}

/* ---------- CONTACTO ---------- */
type CreateContactoDTO struct {
	PersonaID uuid.UUID `json:"persona_id" validate:"required"`
	Tipo      string    `json:"tipo"       validate:"required,oneof=email phone whatsapp"`
	Dato      string    `json:"dato"       validate:"required"`
	IsPrimary bool      `json:"is_primary"`
}

type UpdateContactoDTO struct {
	Dato      *string `json:"dato,omitempty"`
	IsPrimary *bool   `json:"is_primary,omitempty"`
}
