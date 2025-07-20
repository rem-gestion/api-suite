package dto

import (
	"time"

	"github.com/google/uuid"
)

// ===============================
// INVITATION REQUEST DTOs
// ===============================

// CreateInvitationRequest represents request to create a new invitation
type CreateInvitationRequest struct {
	InviteeEmail    string     `json:"invitee_email" validate:"required,email" example:"newemployee@example.com"`
	InviteePersonID *uuid.UUID `json:"invitee_person_id,omitempty" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174001"`
	RoleID          uuid.UUID  `json:"role_id" validate:"required,uuid" example:"123e4567-e89b-12d3-a456-426614174002"`
	BranchID        *uuid.UUID `json:"branch_id,omitempty" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174003"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty" example:"2024-02-01T00:00:00Z"`
	PersonalMessage *string    `json:"personal_message,omitempty" validate:"omitempty,max=500" example:"Welcome to our team! We're excited to have you join us."`
	MetadataRequest
}

// UpdateInvitationRequest represents request to update invitation data
type UpdateInvitationRequest struct {
	RoleID          *uuid.UUID `json:"role_id,omitempty" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174002"`
	BranchID        *uuid.UUID `json:"branch_id,omitempty" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174003"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty" example:"2024-02-01T00:00:00Z"`
	PersonalMessage *string    `json:"personal_message,omitempty" validate:"omitempty,max=500" example:"Welcome to our team! We're excited to have you join us."`
	MetadataRequest
}

// InvitationFiltersRequest represents filters for listing invitations
type InvitationFiltersRequest struct {
	FilterRequest
	OrganizationID  uuid.UUID  `json:"organization_id" form:"organization_id" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	InviterPersonID *uuid.UUID `json:"inviter_person_id" form:"inviter_person_id" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174001"`
	InviteeEmail    *string    `json:"invitee_email" form:"invitee_email" validate:"omitempty,email" example:"newemployee@example.com"`
	RoleID          *uuid.UUID `json:"role_id" form:"role_id" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174002"`
	BranchID        *uuid.UUID `json:"branch_id" form:"branch_id" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174003"`
	ExpiresAfter    *time.Time `json:"expires_after" form:"expires_after" example:"2024-01-01T00:00:00Z"`
	ExpiresBefore   *time.Time `json:"expires_before" form:"expires_before" example:"2024-12-31T23:59:59Z"`
	PaginationRequest
	IncludeOptions
}

// ResendInvitationRequest represents request to resend an invitation
type ResendInvitationRequest struct {
	NewExpiresAt    *time.Time `json:"new_expires_at,omitempty" example:"2024-02-15T00:00:00Z"`
	PersonalMessage *string    `json:"personal_message,omitempty" validate:"omitempty,max=500" example:"Friendly reminder to join our team!"`
}

// CancelInvitationRequest represents request to cancel an invitation
type CancelInvitationRequest struct {
	CancellationReason string `json:"cancellation_reason" validate:"required,max=500" example:"Position no longer available"`
}

// AcceptInvitationRequest represents request to accept an invitation (public endpoint)
type AcceptInvitationRequest struct {
	Token      string                   `json:"token" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	PersonData *CreatePersonDataRequest `json:"person_data,omitempty"`
	UserData   *CreateUserDataRequest   `json:"user_data,omitempty"`
}

// DeclineInvitationRequest represents request to decline an invitation (public endpoint)
type DeclineInvitationRequest struct {
	Token         string  `json:"token" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	DeclineReason *string `json:"decline_reason,omitempty" validate:"omitempty,max=500" example:"Not interested at this time"`
}

// ValidateInvitationRequest represents request to validate an invitation token (public endpoint)
type ValidateInvitationRequest struct {
	Token string `json:"token" form:"token" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// CreatePersonDataRequest represents person data for invitation acceptance
type CreatePersonDataRequest struct {
	Type       string                   `json:"type" validate:"required,oneof=INDIVIDUAL COMPANY" example:"INDIVIDUAL"`
	AvatarURL  *string                  `json:"avatar_url,omitempty" validate:"omitempty,url" example:"https://example.com/avatar.jpg"`
	Sexo       *string                  `json:"sexo,omitempty" validate:"omitempty,oneof=MASCULINO FEMENINO" example:"MASCULINO"`
	Individual *CreateIndividualRequest `json:"individual,omitempty"`
	Company    *CreateCompanyRequest    `json:"company,omitempty"`
	Address    *CreateAddressRequest    `json:"address,omitempty"`
	Contacts   []CreateContactRequest   `json:"contacts,omitempty"`
}

// CreateUserDataRequest represents user data for invitation acceptance
type CreateUserDataRequest struct {
	OnboardStatus *string `json:"onboard_status,omitempty" validate:"omitempty,oneof=NEW IN_PROGRESS DONE" example:"NEW"`
}

// CreateIndividualRequest represents individual person data
type CreateIndividualRequest struct {
	FirstName string  `json:"first_name" validate:"required,min=1,max=100" example:"John"`
	LastName  string  `json:"last_name" validate:"required,min=1,max=100" example:"Doe"`
	DNI       *string `json:"dni,omitempty" validate:"omitempty,max=20" example:"12345678"`
}

// CreateCompanyRequest represents company person data
type CreateCompanyRequest struct {
	LegalName   string  `json:"legal_name" validate:"required,min=1,max=300" example:"Acme Corporation"`
	CUIT        *string `json:"cuit,omitempty" validate:"omitempty,max=20" example:"20-12345678-9"`
	SocietyType *string `json:"society_type,omitempty" validate:"omitempty,max=10" example:"SA"`
}

// CreateAddressRequest represents address data
type CreateAddressRequest struct {
	Floor   *string `json:"floor,omitempty" example:"2"`
	Unit    *string `json:"unit,omitempty" example:"A"`
	Street  string  `json:"street" validate:"required,min=1,max=200" example:"Main Street"`
	Number  int32   `json:"number" validate:"required,min=1" example:"123"`
	City    string  `json:"city" validate:"required,min=1,max=100" example:"New York"`
	State   string  `json:"state" validate:"required,min=1,max=100" example:"NY"`
	Zip     string  `json:"zip" validate:"required,min=1,max=20" example:"10001"`
	Country string  `json:"country" validate:"required,len=2" example:"US"`
}

// CreateContactRequest represents contact data
type CreateContactRequest struct {
	Tipo      string `json:"tipo" validate:"required,oneof=EMAIL PHONE WHATSAPP" example:"EMAIL"`
	Dato      string `json:"dato" validate:"required,min=1,max=255" example:"contact@example.com"`
	IsPrimary bool   `json:"is_primary" example:"true"`
}

// ===============================
// INVITATION RESPONSE DTOs
// ===============================

// OrganizationInviteResponse represents invitation data in responses
type OrganizationInviteResponse struct {
	BaseResponse
	OrganizationID  uuid.UUID  `json:"organization_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	InviterPersonID uuid.UUID  `json:"inviter_person_id" example:"123e4567-e89b-12d3-a456-426614174001"`
	InviteeEmail    string     `json:"invitee_email" example:"newemployee@example.com"`
	InviteePersonID *uuid.UUID `json:"invitee_person_id,omitempty" example:"123e4567-e89b-12d3-a456-426614174002"`
	RoleID          uuid.UUID  `json:"role_id" example:"123e4567-e89b-12d3-a456-426614174003"`
	BranchID        *uuid.UUID `json:"branch_id,omitempty" example:"123e4567-e89b-12d3-a456-426614174004"`
	Token           string     `json:"token,omitempty" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // Only included for admin views
	Status          string     `json:"status" example:"pending"`
	ExpiresAt       time.Time  `json:"expires_at" example:"2024-02-01T00:00:00Z"`
	AcceptedAt      *time.Time `json:"accepted_at,omitempty" example:"2024-01-20T10:30:00Z"`
	RejectedAt      *time.Time `json:"rejected_at,omitempty" example:"2024-01-20T10:30:00Z"`
	CancelledAt     *time.Time `json:"cancelled_at,omitempty" example:"2024-01-20T10:30:00Z"`
	PersonalMessage *string    `json:"personal_message,omitempty" example:"Welcome to our team!"`

	// Embedded optional data
	Organization *OrganizationBasicResponse      `json:"organization,omitempty"`
	Role         *OrganizationRoleResponse       `json:"role,omitempty"`
	Branch       *BranchBasicResponse            `json:"branch,omitempty"`
	Logs         []OrganizationInviteLogResponse `json:"logs,omitempty"`

	// External service data (populated via gRPC)
	InviterPerson *PersonBasicDetails `json:"inviter_person,omitempty"`
	InviteePerson *PersonBasicDetails `json:"invitee_person,omitempty"`

	// Statistics (included when include_stats=true)
	Stats *InvitationStats `json:"stats,omitempty"`
}

// InvitationListResponse represents paginated invitation list
type InvitationListResponse struct {
	Data       []OrganizationInviteResponse `json:"data"`
	Pagination PaginationResponse           `json:"pagination"`
}

// InvitationBasicResponse represents minimal invitation data
type InvitationBasicResponse struct {
	ID             uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	OrganizationID uuid.UUID `json:"organization_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	InviteeEmail   string    `json:"invitee_email" example:"newemployee@example.com"`
	Status         string    `json:"status" example:"pending"`
	ExpiresAt      time.Time `json:"expires_at" example:"2024-02-01T00:00:00Z"`
	CreatedAt      time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

// ValidateInvitationResponse represents invitation validation response (public endpoint)
type ValidateInvitationResponse struct {
	IsValid         bool                       `json:"is_valid" example:"true"`
	InvitationID    *uuid.UUID                 `json:"invitation_id,omitempty" example:"123e4567-e89b-12d3-a456-426614174000"`
	Organization    *OrganizationBasicResponse `json:"organization,omitempty"`
	Role            *RoleBasicResponse         `json:"role,omitempty"`
	Branch          *BranchBasicResponse       `json:"branch,omitempty"`
	InviterPerson   *PersonBasicDetails        `json:"inviter_person,omitempty"`
	ExpiresAt       *time.Time                 `json:"expires_at,omitempty" example:"2024-02-01T00:00:00Z"`
	PersonalMessage *string                    `json:"personal_message,omitempty" example:"Welcome to our team!"`
	TimeRemaining   *string                    `json:"time_remaining,omitempty" example:"5 days"`
	Error           *string                    `json:"error,omitempty" example:"Invitation has expired"`
}

// AcceptInvitationResponse represents invitation acceptance response
type AcceptInvitationResponse struct {
	Success   bool              `json:"success" example:"true"`
	Employee  *EmployeeResponse `json:"employee,omitempty"`
	Message   string            `json:"message" example:"Invitation accepted successfully"`
	NextSteps []string          `json:"next_steps,omitempty" example:"[\"Complete your profile\", \"Download the mobile app\"]"`
}

// ===============================
// INVITATION LOG DTOs
// ===============================

// OrganizationInviteLogResponse represents invitation log data
type OrganizationInviteLogResponse struct {
	ID            uuid.UUID   `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	InviteID      uuid.UUID   `json:"invite_id" example:"123e4567-e89b-12d3-a456-426614174001"`
	Action        string      `json:"action" example:"created"`
	OldStatus     *string     `json:"old_status,omitempty" example:"pending"`
	NewStatus     string      `json:"new_status" example:"accepted"`
	ActorPersonID *uuid.UUID  `json:"actor_person_id,omitempty" example:"123e4567-e89b-12d3-a456-426614174002"`
	ActorType     string      `json:"actor_type" example:"user"`
	ClientIP      *string     `json:"client_ip,omitempty" example:"192.168.1.1"`
	UserAgent     *string     `json:"user_agent,omitempty" example:"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"`
	Notes         *string     `json:"notes,omitempty" example:"User accepted invitation via email link"`
	ErrorDetails  interface{} `json:"error_details,omitempty"`
	CreatedAt     time.Time   `json:"created_at" example:"2024-01-01T00:00:00Z"`

	// External service data (populated via gRPC)
	ActorPerson *PersonBasicDetails `json:"actor_person,omitempty"`
}

// InvitationLogListResponse represents paginated invitation log list
type InvitationLogListResponse struct {
	Data       []OrganizationInviteLogResponse `json:"data"`
	Pagination PaginationResponse              `json:"pagination"`
}

// ===============================
// INVITATION STATISTICS DTOs
// ===============================

// InvitationStats represents invitation statistics
type InvitationStats struct {
	TotalResends    int       `json:"total_resends" example:"2"`
	DaysUntilExpiry int       `json:"days_until_expiry" example:"10"`
	LastActivityAt  time.Time `json:"last_activity_at" example:"2024-01-15T10:30:00Z"`
	AcceptanceRate  *float64  `json:"acceptance_rate,omitempty" example:"85.5"` // Organization-wide acceptance rate
}
