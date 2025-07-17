package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// OrganizationStatus represents the status enum for organizations
type OrganizationStatus string

const (
	OrgStatusActive    OrganizationStatus = "active"
	OrgStatusInactive  OrganizationStatus = "inactive"
	OrgStatusSuspended OrganizationStatus = "suspended"
	OrgStatusPending   OrganizationStatus = "pending"
)

// OrganizationType represents the type enum for organizations
type OrganizationType string

const (
	OrgTypeRealEstate     OrganizationType = "real_estate"
	OrgTypePropertyMgmt   OrganizationType = "property_management"
	OrgTypeConstruction   OrganizationType = "construction"
	OrgTypeArchitecture   OrganizationType = "architecture"
	OrgTypeAppraisal      OrganizationType = "appraisal"
	OrgTypeMortgageBroker OrganizationType = "mortgage_broker"
	OrgTypeInsurance      OrganizationType = "insurance"
	OrgTypeLegal          OrganizationType = "legal"
	OrgTypeConsulting     OrganizationType = "consulting"
	OrgTypeInvestment     OrganizationType = "investment"
)

// Organization represents the main organization entity
type Organization struct {
	ID          uuid.UUID          `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string             `json:"name" gorm:"type:varchar(200);not null" validate:"required,min=2,max=200"`
	DisplayName string             `json:"display_name" gorm:"type:varchar(300);not null" validate:"required,min=2,max=300"`
	Slug        string             `json:"slug" gorm:"type:varchar(100);not null;uniqueIndex" validate:"required,min=2,max=100,slug"`
	Description *string            `json:"description" gorm:"type:text"`
	Status      OrganizationStatus `json:"status" gorm:"type:organization_status_enum;not null;default:'pending'" validate:"required,oneof=active inactive suspended pending"`
	Type        OrganizationType   `json:"type" gorm:"type:organization_type_enum;not null" validate:"required"`
	LegalName   *string            `json:"legal_name" gorm:"type:varchar(300)"`
	TaxID       *string            `json:"tax_id" gorm:"type:varchar(50)"`
	Website     *string            `json:"website" gorm:"type:varchar(255)" validate:"omitempty,url"`
	Phone       *string            `json:"phone" gorm:"type:varchar(32)" validate:"omitempty,phone"`
	Email       *string            `json:"email" gorm:"type:varchar(160)" validate:"omitempty,email"`
	LogoURL     *string            `json:"logo_url" gorm:"type:varchar(500)" validate:"omitempty,url"`
	TimezoneID  string             `json:"timezone_id" gorm:"type:varchar(50);not null;default:'UTC'" validate:"required,timezone"`
	AddressID   *uuid.UUID         `json:"address_id" gorm:"type:uuid"` // FK to address-svc
	IsVerified  bool               `json:"is_verified" gorm:"not null;default:false"`
	VerifiedAt  *time.Time         `json:"verified_at" gorm:"type:timestamp with time zone"`

	// Subscription fields
	SubscriptionID *uuid.UUID `json:"subscription_id" gorm:"type:uuid"`
	PlanID         *uuid.UUID `json:"plan_id" gorm:"type:uuid"`

	// Audit fields
	CreatedAt time.Time       `json:"created_at" gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`
	CreatedBy uuid.UUID       `json:"created_by" gorm:"type:uuid;not null"` // FK to auth-identity-svc
	UpdatedAt time.Time       `json:"updated_at" gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`
	UpdatedBy *uuid.UUID      `json:"updated_by" gorm:"type:uuid"` // FK to auth-identity-svc
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	Settings     []OrganizationSetting     `json:"settings,omitempty" gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
	Owner        *OrganizationOwner        `json:"owner,omitempty" gorm:"foreignKey:OrganizationID"`
	Branches     []OrganizationBranch      `json:"branches,omitempty" gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
	Roles        []OrganizationRole        `json:"roles,omitempty" gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
	Employees    []Employee                `json:"employees,omitempty" gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
	Invitations  []OrganizationInvite      `json:"invitations,omitempty" gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
	Integrations []OrganizationIntegration `json:"integrations,omitempty" gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
	Domains      []OrganizationDomain      `json:"domains,omitempty" gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
	Subscription *OrganizationSubscription `json:"subscription,omitempty" gorm:"foreignKey:OrganizationID"`
}

// OrganizationSetting represents organization-specific settings
type OrganizationSetting struct {
	ID             uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrganizationID uuid.UUID   `json:"organization_id" gorm:"type:uuid;not null;index"`
	SettingKey     string      `json:"setting_key" gorm:"type:varchar(100);not null" validate:"required,min=1,max=100"`
	SettingValue   interface{} `json:"setting_value" gorm:"type:jsonb;not null"`
	Description    *string     `json:"description" gorm:"type:text"`
	IsEditable     bool        `json:"is_editable" gorm:"not null;default:true"`

	// Audit fields
	CreatedAt time.Time  `json:"created_at" gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`
	CreatedBy uuid.UUID  `json:"created_by" gorm:"type:uuid;not null"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`
	UpdatedBy *uuid.UUID `json:"updated_by" gorm:"type:uuid"`

	// Relationships
	Organization Organization `json:"organization,omitempty" gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
}

// OrganizationOwner represents the owner relationship between person and organization
type OrganizationOwner struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrganizationID uuid.UUID `json:"organization_id" gorm:"type:uuid;not null;uniqueIndex"`
	PersonID       uuid.UUID `json:"person_id" gorm:"type:uuid;not null;index"` // FK to person-svc
	UserID         uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`   // FK to auth-identity-svc
	IsFounder      bool      `json:"is_founder" gorm:"not null;default:false"`
	OwnershipPct   *float64  `json:"ownership_pct" gorm:"type:decimal(5,2)" validate:"omitempty,min=0,max=100"`

	// Audit fields
	CreatedAt time.Time  `json:"created_at" gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`
	CreatedBy uuid.UUID  `json:"created_by" gorm:"type:uuid;not null"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`
	UpdatedBy *uuid.UUID `json:"updated_by" gorm:"type:uuid"`

	// Relationships
	Organization Organization `json:"organization,omitempty" gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
}

// OrganizationSubscription represents subscription information for organizations
type OrganizationSubscription struct {
	ID              uuid.UUID          `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrganizationID  uuid.UUID          `json:"organization_id" gorm:"type:uuid;not null;uniqueIndex"`
	PlanName        string             `json:"plan_name" gorm:"type:varchar(100);not null" validate:"required"`
	Status          SubscriptionStatus `json:"status" gorm:"type:subscription_status_enum;not null;default:'trial'" validate:"required"`
	StartDate       time.Time          `json:"start_date" gorm:"type:date;not null"`
	EndDate         *time.Time         `json:"end_date" gorm:"type:date"`
	TrialEndDate    *time.Time         `json:"trial_end_date" gorm:"type:date"`
	BillingCycle    BillingCycle       `json:"billing_cycle" gorm:"type:billing_cycle_enum;not null;default:'monthly'" validate:"required"`
	Amount          float64            `json:"amount" gorm:"type:decimal(10,2);not null" validate:"min=0"`
	Currency        string             `json:"currency" gorm:"type:varchar(3);not null;default:'USD'" validate:"required,len=3,uppercase"`
	Features        pq.StringArray     `json:"features" gorm:"type:text[]"`
	Limits          interface{}        `json:"limits" gorm:"type:jsonb"`
	PaymentMethodID *string            `json:"payment_method_id" gorm:"type:varchar(255)"`

	// Audit fields
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`

	// Relationships
	Organization Organization `json:"organization,omitempty" gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
}

// SubscriptionStatus represents subscription status enum
type SubscriptionStatus string

const (
	SubStatusTrial     SubscriptionStatus = "trial"
	SubStatusActive    SubscriptionStatus = "active"
	SubStatusSuspended SubscriptionStatus = "suspended"
	SubStatusCanceled  SubscriptionStatus = "canceled"
	SubStatusExpired   SubscriptionStatus = "expired"
)

// BillingCycle represents billing cycle enum
type BillingCycle string

const (
	BillingMonthly   BillingCycle = "monthly"
	BillingQuarterly BillingCycle = "quarterly"
	BillingYearly    BillingCycle = "yearly"
)

// TableName specifies the table name for Organization
func (Organization) TableName() string {
	return "organization"
}

// TableName specifies the table name for OrganizationSetting
func (OrganizationSetting) TableName() string {
	return "organization_settings"
}

// TableName specifies the table name for OrganizationOwner
func (OrganizationOwner) TableName() string {
	return "organization_owner"
}

// TableName specifies the table name for OrganizationSubscription
func (OrganizationSubscription) TableName() string {
	return "organization_subscription"
}
