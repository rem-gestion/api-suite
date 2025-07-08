package models

import (
	"time"

	"gorm.io/gorm"
)

// OnboardStatus define los estados de onboarding del usuario
type OnboardStatus string

const (
	OnboardNew        OnboardStatus = "new"
	OnboardInProgress OnboardStatus = "in_progress"
	OnboardDone       OnboardStatus = "done"
)

// User representa un usuario del sistema
type User struct {
	ID            string        `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	AccountID     string        `gorm:"type:uuid;uniqueIndex;not null" json:"account_id"`
	PersonID      *string       `gorm:"type:uuid" json:"person_id,omitempty"` // Referencia a person-svc
	OnboardStatus OnboardStatus `gorm:"type:varchar(16);default:'new'" json:"onboard_status"`
	LastLogin     *time.Time    `json:"last_login,omitempty"`

	// Auditoría
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	UpdatedBy *string    `gorm:"type:uuid" json:"updated_by,omitempty"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// Relación con Account
	Account *Account `gorm:"foreignKey:AccountID" json:"account,omitempty"`
}

// TableName especifica el nombre de la tabla
func (User) TableName() string {
	return "users"
}

// BeforeCreate hook de GORM para generar UUID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		// GORM maneja la generación de UUID con default:uuid_generate_v4()
	}
	return nil
}

// IsOnboardingComplete verifica si el usuario completó el onboarding
func (u *User) IsOnboardingComplete() bool {
	return u.OnboardStatus == OnboardDone
}

// IsOnboardingInProgress verifica si el usuario está en proceso de onboarding
func (u *User) IsOnboardingInProgress() bool {
	return u.OnboardStatus == OnboardInProgress
}

// IsNewUser verifica si es un usuario nuevo
func (u *User) IsNewUser() bool {
	return u.OnboardStatus == OnboardNew
}

// HasPerson verifica si el usuario tiene una persona asociada
func (u *User) HasPerson() bool {
	return u.PersonID != nil && *u.PersonID != ""
}

// UpdateLastLogin actualiza la fecha de último login
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLogin = &now
}
