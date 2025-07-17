package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DomainType represents domain type enum
type DomainType string

const (
	DomainTypeCustom    DomainType = "custom"
	DomainTypeSubdomain DomainType = "subdomain"
	DomainTypeApex      DomainType = "apex"
)

// DomainStatus represents domain status enum
type DomainStatus string

const (
	DomainStatusPending   DomainStatus = "pending"
	DomainStatusActive    DomainStatus = "active"
	DomainStatusInactive  DomainStatus = "inactive"
	DomainStatusSuspended DomainStatus = "suspended"
	DomainStatusExpired   DomainStatus = "expired"
)

// DNSVerificationMethod represents DNS verification method enum
type DNSVerificationMethod string

const (
	DNSVerifyTXT   DNSVerificationMethod = "txt"
	DNSVerifyCNAME DNSVerificationMethod = "cname"
	DNSVerifyHTML  DNSVerificationMethod = "html"
)

// DNSRecordType represents DNS record type enum
type DNSRecordType string

const (
	DNSRecordA     DNSRecordType = "A"
	DNSRecordAAAA  DNSRecordType = "AAAA"
	DNSRecordCNAME DNSRecordType = "CNAME"
	DNSRecordMX    DNSRecordType = "MX"
	DNSRecordTXT   DNSRecordType = "TXT"
	DNSRecordNS    DNSRecordType = "NS"
)

// DomainVerificationType represents domain verification type enum
type DomainVerificationType string

const (
	DomainVerifyDNS  DomainVerificationType = "dns"
	DomainVerifyHTTP DomainVerificationType = "http"
	DomainVerifyFile DomainVerificationType = "file"
)

// VerificationStatus represents verification status enum
type VerificationStatus string

const (
	VerifyStatusPending VerificationStatus = "pending"
	VerifyStatusSuccess VerificationStatus = "success"
	VerifyStatusFailed  VerificationStatus = "failed"
)

// OrganizationDomain represents custom domains per organization
type OrganizationDomain struct {
	ID                    uuid.UUID              `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrganizationID        uuid.UUID              `json:"organization_id" gorm:"type:uuid;not null;index"`
	DomainName            string                 `json:"domain_name" gorm:"type:varchar(255);not null" validate:"required,fqdn"`
	Subdomain             *string                `json:"subdomain" gorm:"type:varchar(100)" validate:"omitempty,alphanum"`
	DomainType            DomainType             `json:"domain_type" gorm:"type:domain_type_enum;not null;default:'custom'" validate:"required"`
	Status                DomainStatus           `json:"status" gorm:"type:domain_status_enum;not null;default:'pending'" validate:"required"`
	IsPrimary             bool                   `json:"is_primary" gorm:"not null;default:false"`
	SSLEnabled            bool                   `json:"ssl_enabled" gorm:"not null;default:false"`
	SSLCertificate        *string                `json:"ssl_certificate" gorm:"type:text"`
	SSLPrivateKey         *string                `json:"ssl_private_key" gorm:"type:text"` // Encrypted in application
	SSLExpiresAt          *time.Time             `json:"ssl_expires_at" gorm:"type:timestamp with time zone"`
	DNSVerified           bool                   `json:"dns_verified" gorm:"not null;default:false"`
	DNSVerificationToken  string                 `json:"dns_verification_token" gorm:"type:varchar(100);default:gen_random_uuid()"`
	DNSVerificationMethod *DNSVerificationMethod `json:"dns_verification_method" gorm:"type:dns_verification_method_enum;default:'txt'"`
	VerificationAttempts  int                    `json:"verification_attempts" gorm:"not null;default:0" validate:"min=0,max=10"`
	LastVerificationAt    *time.Time             `json:"last_verification_at" gorm:"type:timestamp with time zone"`
	RedirectToPrimary     bool                   `json:"redirect_to_primary" gorm:"not null;default:false"`
	CustomHeaders         interface{}            `json:"custom_headers" gorm:"type:jsonb;default:'{}'"`

	// Audit fields
	CreatedAt time.Time       `json:"created_at" gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`
	CreatedBy uuid.UUID       `json:"created_by" gorm:"type:uuid;not null"` // FK to auth-identity-svc
	UpdatedAt time.Time       `json:"updated_at" gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`
	UpdatedBy *uuid.UUID      `json:"updated_by" gorm:"type:uuid"` // FK to auth-identity-svc
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	Organization     Organization                        `json:"organization,omitempty" gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
	DNSRecords       []OrganizationDomainDNS             `json:"dns_records,omitempty" gorm:"foreignKey:DomainID;constraint:OnDelete:CASCADE"`
	VerificationLogs []OrganizationDomainVerificationLog `json:"verification_logs,omitempty" gorm:"foreignKey:DomainID;constraint:OnDelete:CASCADE"`
}

// OrganizationDomainDNS represents DNS configuration per domain
type OrganizationDomainDNS struct {
	ID                   uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	DomainID             uuid.UUID     `json:"domain_id" gorm:"type:uuid;not null;index"`
	RecordType           DNSRecordType `json:"record_type" gorm:"type:dns_record_type_enum;not null"`
	RecordName           string        `json:"record_name" gorm:"type:varchar(255);not null" validate:"required"`
	RecordValue          string        `json:"record_value" gorm:"type:text;not null" validate:"required"`
	RecordTTL            int           `json:"record_ttl" gorm:"not null;default:300" validate:"min=60,max=86400"`
	IsRequired           bool          `json:"is_required" gorm:"not null;default:true"`
	IsVerified           bool          `json:"is_verified" gorm:"not null;default:false"`
	VerificationAttempts int           `json:"verification_attempts" gorm:"not null;default:0" validate:"min=0,max=10"`
	LastVerificationAt   *time.Time    `json:"last_verification_at" gorm:"type:timestamp with time zone"`
	Notes                *string       `json:"notes" gorm:"type:text"`

	// Audit fields
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`

	// Relationships
	Domain           OrganizationDomain                  `json:"domain,omitempty" gorm:"foreignKey:DomainID;constraint:OnDelete:CASCADE"`
	VerificationLogs []OrganizationDomainVerificationLog `json:"verification_logs,omitempty" gorm:"foreignKey:DNSRecordID;constraint:OnDelete:SET NULL"`
}

// OrganizationDomainVerificationLog represents verification logs for domains
type OrganizationDomainVerificationLog struct {
	ID               uuid.UUID              `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	DomainID         uuid.UUID              `json:"domain_id" gorm:"type:uuid;not null;index"`
	DNSRecordID      *uuid.UUID             `json:"dns_record_id" gorm:"type:uuid"` // FK to organization_domain_dns
	VerificationType DomainVerificationType `json:"verification_type" gorm:"type:domain_verification_type_enum;not null"`
	Status           VerificationStatus     `json:"status" gorm:"type:verification_status_enum;not null"`
	Details          interface{}            `json:"details" gorm:"type:jsonb;default:'{}'"`
	ErrorMessage     *string                `json:"error_message" gorm:"type:text"`
	ResponseData     interface{}            `json:"response_data" gorm:"type:jsonb"`
	DurationMs       *int                   `json:"duration_ms" gorm:"" validate:"omitempty,min=0"`

	// Audit field
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`

	// Relationships
	Domain    OrganizationDomain     `json:"domain,omitempty" gorm:"foreignKey:DomainID;constraint:OnDelete:CASCADE"`
	DNSRecord *OrganizationDomainDNS `json:"dns_record,omitempty" gorm:"foreignKey:DNSRecordID;constraint:OnDelete:SET NULL"`
}

// TableName specifies the table name for OrganizationDomain
func (OrganizationDomain) TableName() string {
	return "organization_domain"
}

// TableName specifies the table name for OrganizationDomainDNS
func (OrganizationDomainDNS) TableName() string {
	return "organization_domain_dns"
}

// TableName specifies the table name for OrganizationDomainVerificationLog
func (OrganizationDomainVerificationLog) TableName() string {
	return "organization_domain_verification_log"
}
