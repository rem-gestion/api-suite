package models

import (
	"time"

	"github.com/google/uuid"
)

// InvitationStatus represents the status enum for invitations
type InvitationStatus string

const (
	InvStatusPending   InvitationStatus = "pending"
	InvStatusAccepted  InvitationStatus = "accepted"
	InvStatusRejected  InvitationStatus = "rejected"
	InvStatusExpired   InvitationStatus = "expired"
	InvStatusCancelled InvitationStatus = "cancelled"
)

// InvitationLogAction represents the log action enum for invitation events
type InvitationLogAction string

const (
	InvLogCreated   InvitationLogAction = "created"
	InvLogAccepted  InvitationLogAction = "accepted"
	InvLogRejected  InvitationLogAction = "rejected"
	InvLogExpired   InvitationLogAction = "expired"
	InvLogCancelled InvitationLogAction = "cancelled"
	InvLogUpdated   InvitationLogAction = "updated"
)

// OrganizationInvite represents invitations to join an organization
type OrganizationInvite struct {
	ID              uuid.UUID        `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrganizationID  uuid.UUID        `json:"organization_id" gorm:"type:uuid;not null;index"`
	InviterPersonID uuid.UUID        `json:"inviter_person_id" gorm:"type:uuid;not null"` // FK to person-svc
	InviteeEmail    string           `json:"invitee_email" gorm:"type:varchar(255);not null" validate:"required,email"`
	InviteePersonID *uuid.UUID       `json:"invitee_person_id" gorm:"type:uuid"` // FK to person-svc
	RoleID          uuid.UUID        `json:"role_id" gorm:"type:uuid;not null"`
	BranchID        *uuid.UUID       `json:"branch_id" gorm:"type:uuid"` // FK to organization_branch
	Token           string           `json:"token" gorm:"type:varchar(500);not null;uniqueIndex" validate:"required"`
	Status          InvitationStatus `json:"status" gorm:"type:invitation_status_enum;not null;default:'pending'" validate:"required"`
	ExpiresAt       time.Time        `json:"expires_at" gorm:"type:timestamp with time zone;not null"`
	AcceptedAt      *time.Time       `json:"accepted_at" gorm:"type:timestamp with time zone"`
	RejectedAt      *time.Time       `json:"rejected_at" gorm:"type:timestamp with time zone"`
	CancelledAt     *time.Time       `json:"cancelled_at" gorm:"type:timestamp with time zone"`
	Metadata        interface{}      `json:"metadata" gorm:"type:jsonb;default:'{}'"`

	// Audit fields
	CreatedAt time.Time  `json:"created_at" gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`
	CreatedBy uuid.UUID  `json:"created_by" gorm:"type:uuid;not null"` // FK to auth-identity-svc
	UpdatedAt time.Time  `json:"updated_at" gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`
	UpdatedBy *uuid.UUID `json:"updated_by" gorm:"type:uuid"` // FK to auth-identity-svc

	// Relationships
	Organization Organization            `json:"organization,omitempty" gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
	Role         OrganizationRole        `json:"role,omitempty" gorm:"foreignKey:RoleID;constraint:OnDelete:RESTRICT"`
	Branch       *OrganizationBranch     `json:"branch,omitempty" gorm:"foreignKey:BranchID;constraint:OnDelete:SET NULL"`
	Logs         []OrganizationInviteLog `json:"logs,omitempty" gorm:"foreignKey:InviteID;constraint:OnDelete:CASCADE"`
}

// OrganizationInviteLog represents audit trail for invitation lifecycle events
type OrganizationInviteLog struct {
	ID            uuid.UUID           `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	InviteID      uuid.UUID           `json:"invite_id" gorm:"type:uuid;not null;index"`
	Action        InvitationLogAction `json:"action" gorm:"type:invitation_log_action_enum;not null"`
	OldStatus     *InvitationStatus   `json:"old_status" gorm:"type:invitation_status_enum"`
	NewStatus     InvitationStatus    `json:"new_status" gorm:"type:invitation_status_enum;not null"`
	ActorPersonID *uuid.UUID          `json:"actor_person_id" gorm:"type:uuid"` // FK to person-svc
	ActorType     string              `json:"actor_type" gorm:"type:varchar(20);default:'user'" validate:"oneof=user system cron"`
	ClientIP      *string             `json:"client_ip" gorm:"type:inet"`
	UserAgent     *string             `json:"user_agent" gorm:"type:text"`
	Notes         *string             `json:"notes" gorm:"type:text"`
	ErrorDetails  interface{}         `json:"error_details" gorm:"type:jsonb"`

	// Audit field
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`

	// Relationships
	Invite OrganizationInvite `json:"invite,omitempty" gorm:"foreignKey:InviteID;constraint:OnDelete:CASCADE"`
}

// TableName specifies the table name for OrganizationInvite
func (OrganizationInvite) TableName() string {
	return "organization_invite"
}

// TableName specifies the table name for OrganizationInviteLog
func (OrganizationInviteLog) TableName() string {
	return "organization_invite_log"
}
