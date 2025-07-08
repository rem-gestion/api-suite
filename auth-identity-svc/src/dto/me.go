package dto

import "time"

// UserProfileResponse representa el perfil del usuario
type UserProfileResponse struct {
	UserID        string     `json:"user_id"`
	Email         string     `json:"email"`
	Status        string     `json:"status"`
	OnboardStatus string     `json:"onboard_status"`
	PersonID      string     `json:"person_id,omitempty"`
	LastLogin     *time.Time `json:"last_login,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`

	// Datos de la persona si están disponibles
	PersonData *PersonProfileData `json:"person_data,omitempty"`
}

// PersonProfileData representa datos de persona en el perfil
type PersonProfileData struct {
	ID         string          `json:"id"`
	Type       string          `json:"type"`
	AvatarURL  string          `json:"avatar_url,omitempty"`
	Sexo       string          `json:"sexo,omitempty"`
	Individual *IndividualData `json:"individual,omitempty"`
	Company    *CompanyData    `json:"company,omitempty"`
	Contacts   []ContactData   `json:"contacts,omitempty"`
}

// UpdateProfileRequest representa la petición para actualizar el perfil
type UpdateProfileRequest struct {
	PersonData *PersonData `json:"person_data,omitempty"`
}

// ListUsersResponse representa la respuesta de lista de usuarios
type ListUsersResponse struct {
	Users   []*UserProfileResponse `json:"users"`
	Total   int64                  `json:"total"`
	Page    int                    `json:"page"`
	PerPage int                    `json:"per_page"`
}
