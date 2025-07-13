package dto

// RegisterRequest representa la petición de registro de usuario
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Provider string `json:"provider" binding:"omitempty,oneof=email google"`

	// Datos de persona (opcional en registro inicial)
	PersonData *PersonData `json:"person,omitempty"`
}

// PersonData representa los datos de persona para el registro
type PersonData struct {
	Type       string          `json:"type" binding:"required,oneof=individual company"`
	Individual *IndividualData `json:"individual,omitempty"`
	Company    *CompanyData    `json:"company,omitempty"`
	Address    *AddressData    `json:"address,omitempty"`
	Contacts   []ContactData   `json:"contacts,omitempty"`
	AvatarURL  string          `json:"avatar_url,omitempty"`
	Sexo       string          `json:"sexo,omitempty" binding:"omitempty,oneof=masculino femenino"`
}

// IndividualData representa datos de persona individual
type IndividualData struct {
	FirstName      string `json:"first_name" binding:"required"`
	LastName       string `json:"last_name" binding:"required"`
	DocumentType   string `json:"document_type,omitempty"`
	DocumentNumber string `json:"document_number,omitempty"`
}

// CompanyData representa datos de empresa
type CompanyData struct {
	LegalName   string `json:"legal_name" binding:"required"`
	CUIT        string `json:"cuit,omitempty"`
	SocietyType string `json:"society_type,omitempty"`
}

// AddressData representa datos de dirección
type AddressData struct {
	Floor      string `json:"floor,omitempty"`
	Unit       string `json:"unit,omitempty"`
	Street     string `json:"street" binding:"required"`
	Number     int    `json:"number,omitempty"` // Como string para flexibilidad
	City       string `json:"city" binding:"required"`
	State      string `json:"state" binding:"required"`
	Zip        string `json:"zip,omitempty"`         // Mantener zip para compatibilidad con protobuf
	PostalCode string `json:"postal_code,omitempty"` // Alias para zip
	Country    string `json:"country" binding:"required"`
}

// ContactData representa datos de contacto
type ContactData struct {
	Type      string `json:"type" binding:"required,oneof=email phone whatsapp"`
	Value     string `json:"value" binding:"required"`
	IsPrimary bool   `json:"is_primary,omitempty"`
}

// RegisterResponse representa la respuesta del registro
type RegisterResponse struct {
	UserID       string `json:"user_id"`
	Email        string `json:"email"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	PersonID     string `json:"person_id,omitempty"`
}
