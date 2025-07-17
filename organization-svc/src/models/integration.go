package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// IntegrationCategory represents integration category enum
type IntegrationCategory string

const (
	IntegCategoryAccounting    IntegrationCategory = "accounting"
	IntegCategoryCRM           IntegrationCategory = "crm"
	IntegCategoryMarketing     IntegrationCategory = "marketing"
	IntegCategoryCommunication IntegrationCategory = "communication"
	IntegCategoryAnalytics     IntegrationCategory = "analytics"
	IntegCategoryAutomation    IntegrationCategory = "automation"
	IntegCategoryPayment       IntegrationCategory = "payment"
	IntegCategoryDocument      IntegrationCategory = "document"
	IntegCategoryProductivity  IntegrationCategory = "productivity"
	IntegCategorySecurity      IntegrationCategory = "security"
)

// IntegrationStatus represents integration status enum
type IntegrationStatus string

const (
	IntegStatusInactive  IntegrationStatus = "inactive"
	IntegStatusActive    IntegrationStatus = "active"
	IntegStatusError     IntegrationStatus = "error"
	IntegStatusSuspended IntegrationStatus = "suspended"
)

// SyncStatus represents sync status enum
type SyncStatus string

const (
	SyncStatusSuccess SyncStatus = "success"
	SyncStatusError   SyncStatus = "error"
	SyncStatusPartial SyncStatus = "partial"
)

// SyncFrequency represents sync frequency enum
type SyncFrequency string

const (
	SyncFreqRealtime SyncFrequency = "realtime"
	SyncFreqHourly   SyncFrequency = "hourly"
	SyncFreqDaily    SyncFrequency = "daily"
	SyncFreqWeekly   SyncFrequency = "weekly"
	SyncFreqMonthly  SyncFrequency = "monthly"
	SyncFreqManual   SyncFrequency = "manual"
)

// IntegrationType represents the catalog of available integration types
type IntegrationType struct {
	ID                  uuid.UUID           `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name                string              `json:"name" gorm:"type:varchar(100);not null;uniqueIndex" validate:"required,min=2,max=100"`
	DisplayName         string              `json:"display_name" gorm:"type:varchar(200);not null" validate:"required,min=2,max=200"`
	Description         *string             `json:"description" gorm:"type:text"`
	Category            IntegrationCategory `json:"category" gorm:"type:integration_category_enum;not null" validate:"required"`
	Provider            string              `json:"provider" gorm:"type:varchar(100);not null" validate:"required,min=2,max=100"`
	Version             string              `json:"version" gorm:"type:varchar(20);not null;default:'1.0'" validate:"required"`
	IsActive            bool                `json:"is_active" gorm:"not null;default:true"`
	ConfigurationSchema interface{}         `json:"configuration_schema" gorm:"type:jsonb"`
	WebhookSupport      bool                `json:"webhook_support" gorm:"not null;default:false"`
	OAuthSupport        bool                `json:"oauth_support" gorm:"not null;default:false"`
	APIKeySupport       bool                `json:"api_key_support" gorm:"not null;default:true"`

	// Audit fields
	CreatedAt time.Time  `json:"created_at" gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`
	UpdatedBy *uuid.UUID `json:"updated_by" gorm:"type:uuid"` // FK to auth-identity-svc

	// Relationships
	OrganizationIntegrations []OrganizationIntegration `json:"organization_integrations,omitempty" gorm:"foreignKey:IntegrationTypeID;constraint:OnDelete:RESTRICT"`
}

// OrganizationIntegration represents organization-specific integration configurations
type OrganizationIntegration struct {
	ID                uuid.UUID         `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrganizationID    uuid.UUID         `json:"organization_id" gorm:"type:uuid;not null;index"`
	IntegrationTypeID uuid.UUID         `json:"integration_type_id" gorm:"type:uuid;not null;index"`
	Name              string            `json:"name" gorm:"type:varchar(200);not null" validate:"required,min=2,max=200"`
	Description       *string           `json:"description" gorm:"type:text"`
	Status            IntegrationStatus `json:"status" gorm:"type:integration_status_enum;not null;default:'inactive'" validate:"required"`
	Configuration     interface{}       `json:"configuration" gorm:"type:jsonb;not null;default:'{}'"`
	Credentials       interface{}       `json:"credentials" gorm:"type:jsonb;default:'{}'"`
	LastSyncAt        *time.Time        `json:"last_sync_at" gorm:"type:timestamp with time zone"`
	LastSyncStatus    *SyncStatus       `json:"last_sync_status" gorm:"type:sync_status_enum"`
	LastSyncError     *string           `json:"last_sync_error" gorm:"type:text"`
	SyncFrequency     *SyncFrequency    `json:"sync_frequency" gorm:"type:sync_frequency_enum"`
	AutoSyncEnabled   bool              `json:"auto_sync_enabled" gorm:"not null;default:false"`
	WebhookURL        *string           `json:"webhook_url" gorm:"type:varchar(500)" validate:"omitempty,url"`
	WebhookSecret     *string           `json:"webhook_secret" gorm:"type:varchar(100)"`
	OAuthToken        interface{}       `json:"oauth_token" gorm:"type:jsonb;default:'{}'"`
	APIUsageCount     int               `json:"api_usage_count" gorm:"not null;default:0" validate:"min=0"`
	APIRateLimit      *int              `json:"api_rate_limit" gorm:"" validate:"omitempty,min=1"`
	IsActive          bool              `json:"is_active" gorm:"not null;default:true"`

	// Audit fields
	CreatedAt time.Time       `json:"created_at" gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`
	CreatedBy uuid.UUID       `json:"created_by" gorm:"type:uuid;not null"` // FK to auth-identity-svc
	UpdatedAt time.Time       `json:"updated_at" gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`
	UpdatedBy *uuid.UUID      `json:"updated_by" gorm:"type:uuid"` // FK to auth-identity-svc
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	Organization    Organization                   `json:"organization,omitempty" gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
	IntegrationType IntegrationType                `json:"integration_type,omitempty" gorm:"foreignKey:IntegrationTypeID;constraint:OnDelete:RESTRICT"`
	Events          []OrganizationIntegrationEvent `json:"events,omitempty" gorm:"foreignKey:IntegrationID;constraint:OnDelete:CASCADE"`
}

// IntegrationEventType represents integration event type enum
type IntegrationEventType string

const (
	IntegEventSync    IntegrationEventType = "sync"
	IntegEventWebhook IntegrationEventType = "webhook"
	IntegEventAuth    IntegrationEventType = "auth"
	IntegEventConfig  IntegrationEventType = "config"
	IntegEventTest    IntegrationEventType = "test"
)

// IntegrationEventStatus represents integration event status enum
type IntegrationEventStatus string

const (
	IntegEventStatusPending    IntegrationEventStatus = "pending"
	IntegEventStatusProcessing IntegrationEventStatus = "processing"
	IntegEventStatusSuccess    IntegrationEventStatus = "success"
	IntegEventStatusFailed     IntegrationEventStatus = "failed"
	IntegEventStatusRetrying   IntegrationEventStatus = "retrying"
)

// OrganizationIntegrationEvent represents events for integration workers
type OrganizationIntegrationEvent struct {
	ID            uuid.UUID              `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	IntegrationID uuid.UUID              `json:"integration_id" gorm:"type:uuid;not null;index"`
	EventType     IntegrationEventType   `json:"event_type" gorm:"type:integration_event_type_enum;not null"`
	EventData     interface{}            `json:"event_data" gorm:"type:jsonb;not null;default:'{}'"`
	Status        IntegrationEventStatus `json:"status" gorm:"type:integration_event_status_enum;not null;default:'pending'"`
	ErrorMessage  *string                `json:"error_message" gorm:"type:text"`
	RetryCount    int                    `json:"retry_count" gorm:"not null;default:0" validate:"min=0"`
	MaxRetries    int                    `json:"max_retries" gorm:"not null;default:3" validate:"min=0"`
	ScheduledAt   *time.Time             `json:"scheduled_at" gorm:"type:timestamp with time zone"`
	ProcessedAt   *time.Time             `json:"processed_at" gorm:"type:timestamp with time zone"`

	// Audit field
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`

	// Relationships
	Integration OrganizationIntegration `json:"integration,omitempty" gorm:"foreignKey:IntegrationID;constraint:OnDelete:CASCADE"`
}

// TableName specifies the table name for IntegrationType
func (IntegrationType) TableName() string {
	return "integration_type"
}

// TableName specifies the table name for OrganizationIntegration
func (OrganizationIntegration) TableName() string {
	return "organization_integration"
}

// TableName specifies the table name for OrganizationIntegrationEvent
func (OrganizationIntegrationEvent) TableName() string {
	return "organization_integration_event"
}
