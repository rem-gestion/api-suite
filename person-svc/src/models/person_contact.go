package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ---- enums ----
type Sexo string

const (
	SexoMasculino Sexo = "masculino"
	SexoFemenino  Sexo = "femenino"
)

type PersonType string

const (
	PersonIndividual PersonType = "individual"
	PersonCompany    PersonType = "company"
)

// ---- core tables ----
type Person struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Type      PersonType `gorm:"size:10;not null"`
	AddressID *uuid.UUID `gorm:"type:uuid"`
	AvatarURL *string
	Sexo      *Sexo `gorm:"type:sexo"`

	CreatedAt time.Time
	UpdatedAt *time.Time
	UpdatedBy *uuid.UUID
	DeletedAt gorm.DeletedAt `gorm:"index"`
	// Relations
	Individual *Individual `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Company    *Company    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Contactos  []Contacto  `gorm:"foreignKey:PersonaID"`
}

func (Person) TableName() string { return "person" }

// -------- sub-tipos ----------
type Individual struct {
	PersonID  uuid.UUID `gorm:"type:uuid;primaryKey"`
	FirstName string    `gorm:"size:64;not null"`
	LastName  string    `gorm:"size:64;not null"`
	DNI       string    `gorm:"size:8;unique"`

	CreatedAt time.Time
	UpdatedAt *time.Time
	UpdatedBy *uuid.UUID
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Individual) TableName() string { return "individual" }

type Company struct {
	PersonID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	LegalName   string    `gorm:"size:120;not null"`
	CUIT        string    `gorm:"column:cuit;size:11;unique"`
	SocietyType string    `gorm:"size:12"`

	CreatedAt time.Time
	UpdatedAt *time.Time
	UpdatedBy *uuid.UUID
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Company) TableName() string { return "company" }

// ---------- contacto ----------
type Contacto struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	PersonaID uuid.UUID `gorm:"type:uuid;index"`
	Tipo      string    `gorm:"size:16;not null"` // email | phone | whatsapp
	Dato      string    `gorm:"size:128;not null"`
	IsPrimary bool      `gorm:"default:false"`
	CreatedAt time.Time
	CreatedBy *uuid.UUID
	UpdatedAt *time.Time
	UpdatedBy *uuid.UUID
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Contacto) TableName() string { return "contacto" }
