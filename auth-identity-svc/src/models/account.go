package models

import (
	"time"

	"gorm.io/gorm"
)

// AccountProvider define los tipos de proveedores de autenticación
type AccountProvider string

const (
	ProviderEmail  AccountProvider = "email"
	ProviderGoogle AccountProvider = "google"
)

// AccountStatus define los estados de una cuenta
type AccountStatus string

const (
	StatusPending   AccountStatus = "pending"
	StatusActive    AccountStatus = "active"
	StatusSuspended AccountStatus = "suspended"
)

// Account representa una cuenta de autenticación
type Account struct {
	ID           string          `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Provider     AccountProvider `gorm:"type:varchar(16);not null" json:"provider"`
	Email        string          `gorm:"type:varchar(160);uniqueIndex;not null" json:"email"`
	PasswordHash string          `gorm:"type:varchar(255)" json:"-"` // No exponer en JSON
	Status       AccountStatus   `gorm:"type:varchar(16);default:'pending'" json:"status"`

	// Auditoría
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	UpdatedBy *string    `gorm:"type:uuid" json:"updated_by,omitempty"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// Relación 1:1 con User
	User *User `gorm:"foreignKey:AccountID" json:"user,omitempty"`
}

// TableName especifica el nombre de la tabla
func (Account) TableName() string {
	return "auth_accounts"
}

// BeforeCreate hook de GORM para generar UUID
func (a *Account) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		// GORM maneja la generación de UUID con default:uuid_generate_v4()
	}
	return nil
}

// IsActive verifica si la cuenta está activa
func (a *Account) IsActive() bool {
	return a.Status == StatusActive
}

// IsPending verifica si la cuenta está pendiente
func (a *Account) IsPending() bool {
	return a.Status == StatusPending
}

// IsSuspended verifica si la cuenta está suspendida
func (a *Account) IsSuspended() bool {
	return a.Status == StatusSuspended
}
