package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/streadway/amqp"
	"go.uber.org/zap"

	"github.com/rem-gestion/api-suite/organization/src/models"
)

// EventService define la interfaz para publicar eventos del dominio de organizaciones
type EventService interface {
	// Organization events
	PublishOrganizationCreated(org *models.Organization) error
	PublishOrganizationUpdated(org *models.Organization) error
	PublishOrganizationStatusChanged(org *models.Organization, newStatus models.OrganizationStatus) error
	PublishOrganizationDeleted(orgID string) error
	PublishOrganizationSettingsUpdated(orgID string, settings map[string]interface{}) error

	// Employee events
	PublishEmployeeAdded(org *models.Organization, emp *models.Employee) error
	PublishEmployeeUpdated(org *models.Organization, emp *models.Employee) error
	PublishEmployeeStatusChanged(org *models.Organization, emp *models.Employee, newStatus models.EmployeeStatus) error
	PublishEmployeeRemoved(orgID, empID string) error
	PublishEmployeeRoleAssigned(org *models.Organization, emp *models.Employee, role *models.OrganizationRole) error
	PublishEmployeeRoleRemoved(org *models.Organization, emp *models.Employee, roleID string) error

	// Branch events
	PublishBranchCreated(org *models.Organization, branch *models.OrganizationBranch) error
	PublishBranchUpdated(org *models.Organization, branch *models.OrganizationBranch) error
	PublishBranchDeleted(orgID, branchID string) error

	// Role events
	PublishRoleCreated(org *models.Organization, role *models.OrganizationRole) error
	PublishRoleUpdated(org *models.Organization, role *models.OrganizationRole) error
	PublishRoleDeleted(orgID, roleID string) error

	// Invitation events
	PublishInvitationSent(org *models.Organization, invite *models.OrganizationInvite) error
	PublishInvitationAccepted(org *models.Organization, invite *models.OrganizationInvite, emp *models.Employee) error
	PublishInvitationDeclined(org *models.Organization, invite *models.OrganizationInvite) error
	PublishInvitationCanceled(org *models.Organization, invite *models.OrganizationInvite) error

	// Owner events
	PublishOwnerAdded(org *models.Organization, owner *models.OrganizationOwner) error
	PublishOwnerRemoved(orgID, ownerID string) error
}

// RabbitMQEventService implementa EventService usando RabbitMQ
type RabbitMQEventService struct {
	channel  *amqp.Channel
	exchange string
	logger   *zap.Logger
}

// NewRabbitMQEventService crea una nueva instancia del servicio de eventos
func NewRabbitMQEventService(channel *amqp.Channel, exchange string, logger *zap.Logger) EventService {
	return &RabbitMQEventService{
		channel:  channel,
		exchange: exchange,
		logger:   logger,
	}
}

// NewEventService crea una nueva instancia de EventService
// Si channel es nil, retorna un NoOpEventService
func NewEventService(conn *amqp.Connection, logger *zap.Logger) EventService {
	if conn == nil {
		return &NoOpEventService{logger: logger}
	}

	channel, err := conn.Channel()
	if err != nil {
		logger.Error("Failed to create RabbitMQ channel", zap.Error(err))
		return &NoOpEventService{logger: logger}
	}

	// Declare exchange for organization events
	exchange := "organization.events"
	err = channel.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		logger.Error("Failed to declare exchange", zap.Error(err))
		return &NoOpEventService{logger: logger}
	}

	return NewRabbitMQEventService(channel, exchange, logger)
}

// Event topics/routing keys
const (
	// Organization events
	TopicOrganizationCreated         = "organization.created"
	TopicOrganizationUpdated         = "organization.updated"
	TopicOrganizationStatusChanged   = "organization.status_changed"
	TopicOrganizationDeleted         = "organization.deleted"
	TopicOrganizationSettingsUpdated = "organization.settings_updated"

	// Employee events
	TopicEmployeeAdded         = "organization.employee.added"
	TopicEmployeeUpdated       = "organization.employee.updated"
	TopicEmployeeStatusChanged = "organization.employee.status_changed"
	TopicEmployeeRemoved       = "organization.employee.removed"
	TopicEmployeeRoleAssigned  = "organization.employee.role_assigned"
	TopicEmployeeRoleRemoved   = "organization.employee.role_removed"

	// Branch events
	TopicBranchCreated = "organization.branch.created"
	TopicBranchUpdated = "organization.branch.updated"
	TopicBranchDeleted = "organization.branch.deleted"

	// Role events
	TopicRoleCreated = "organization.role.created"
	TopicRoleUpdated = "organization.role.updated"
	TopicRoleDeleted = "organization.role.deleted"

	// Invitation events
	TopicInvitationSent     = "organization.invitation.sent"
	TopicInvitationAccepted = "organization.invitation.accepted"
	TopicInvitationDeclined = "organization.invitation.declined"
	TopicInvitationCanceled = "organization.invitation.canceled"

	// Owner events
	TopicOwnerAdded   = "organization.owner.added"
	TopicOwnerRemoved = "organization.owner.removed"
)

// Event payload structures

type OrganizationEvent struct {
	EventID      string                 `json:"event_id"`
	EventType    string                 `json:"event_type"`
	Timestamp    time.Time              `json:"timestamp"`
	Source       string                 `json:"source"`
	Organization *models.Organization   `json:"organization"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

type EmployeeEvent struct {
	EventID      string                   `json:"event_id"`
	EventType    string                   `json:"event_type"`
	Timestamp    time.Time                `json:"timestamp"`
	Source       string                   `json:"source"`
	Organization *models.Organization     `json:"organization"`
	Employee     *models.Employee         `json:"employee"`
	Role         *models.OrganizationRole `json:"role,omitempty"`
	Metadata     map[string]interface{}   `json:"metadata,omitempty"`
}

type BranchEvent struct {
	EventID      string                     `json:"event_id"`
	EventType    string                     `json:"event_type"`
	Timestamp    time.Time                  `json:"timestamp"`
	Source       string                     `json:"source"`
	Organization *models.Organization       `json:"organization"`
	Branch       *models.OrganizationBranch `json:"branch"`
	Metadata     map[string]interface{}     `json:"metadata,omitempty"`
}

type RoleEvent struct {
	EventID      string                   `json:"event_id"`
	EventType    string                   `json:"event_type"`
	Timestamp    time.Time                `json:"timestamp"`
	Source       string                   `json:"source"`
	Organization *models.Organization     `json:"organization"`
	Role         *models.OrganizationRole `json:"role"`
	Metadata     map[string]interface{}   `json:"metadata,omitempty"`
}

type InvitationEvent struct {
	EventID      string                     `json:"event_id"`
	EventType    string                     `json:"event_type"`
	Timestamp    time.Time                  `json:"timestamp"`
	Source       string                     `json:"source"`
	Organization *models.Organization       `json:"organization"`
	Invitation   *models.OrganizationInvite `json:"invitation"`
	Employee     *models.Employee           `json:"employee,omitempty"`
	Metadata     map[string]interface{}     `json:"metadata,omitempty"`
}

type OwnerEvent struct {
	EventID      string                    `json:"event_id"`
	EventType    string                    `json:"event_type"`
	Timestamp    time.Time                 `json:"timestamp"`
	Source       string                    `json:"source"`
	Organization *models.Organization      `json:"organization"`
	Owner        *models.OrganizationOwner `json:"owner,omitempty"`
	OwnerID      string                    `json:"owner_id,omitempty"`
	Metadata     map[string]interface{}    `json:"metadata,omitempty"`
}

// Organization events implementation

func (e *RabbitMQEventService) PublishOrganizationCreated(org *models.Organization) error {
	event := OrganizationEvent{
		EventID:      generateEventID(),
		EventType:    TopicOrganizationCreated,
		Timestamp:    time.Now(),
		Source:       "organization-svc",
		Organization: org,
		Metadata: map[string]interface{}{
			"version": "1.0",
		},
	}

	return e.publishEvent(TopicOrganizationCreated, event)
}

func (e *RabbitMQEventService) PublishOrganizationUpdated(org *models.Organization) error {
	event := OrganizationEvent{
		EventID:      generateEventID(),
		EventType:    TopicOrganizationUpdated,
		Timestamp:    time.Now(),
		Source:       "organization-svc",
		Organization: org,
		Metadata: map[string]interface{}{
			"version": "1.0",
		},
	}

	return e.publishEvent(TopicOrganizationUpdated, event)
}

func (e *RabbitMQEventService) PublishOrganizationStatusChanged(org *models.Organization, newStatus models.OrganizationStatus) error {
	event := OrganizationEvent{
		EventID:      generateEventID(),
		EventType:    TopicOrganizationStatusChanged,
		Timestamp:    time.Now(),
		Source:       "organization-svc",
		Organization: org,
		Metadata: map[string]interface{}{
			"version":    "1.0",
			"new_status": string(newStatus),
			"old_status": string(org.Status),
		},
	}

	return e.publishEvent(TopicOrganizationStatusChanged, event)
}

func (e *RabbitMQEventService) PublishOrganizationDeleted(orgID string) error {
	event := OrganizationEvent{
		EventID:   generateEventID(),
		EventType: TopicOrganizationDeleted,
		Timestamp: time.Now(),
		Source:    "organization-svc",
		Metadata: map[string]interface{}{
			"version":         "1.0",
			"organization_id": orgID,
		},
	}

	return e.publishEvent(TopicOrganizationDeleted, event)
}

func (e *RabbitMQEventService) PublishOrganizationSettingsUpdated(orgID string, settings map[string]interface{}) error {
	event := OrganizationEvent{
		EventID:   generateEventID(),
		EventType: TopicOrganizationSettingsUpdated,
		Timestamp: time.Now(),
		Source:    "organization-svc",
		Metadata: map[string]interface{}{
			"version":         "1.0",
			"organization_id": orgID,
			"settings":        settings,
		},
	}

	return e.publishEvent(TopicOrganizationSettingsUpdated, event)
}

// Employee events implementation

func (e *RabbitMQEventService) PublishEmployeeAdded(org *models.Organization, emp *models.Employee) error {
	event := EmployeeEvent{
		EventID:      generateEventID(),
		EventType:    TopicEmployeeAdded,
		Timestamp:    time.Now(),
		Source:       "organization-svc",
		Organization: org,
		Employee:     emp,
		Metadata: map[string]interface{}{
			"version": "1.0",
		},
	}

	return e.publishEvent(TopicEmployeeAdded, event)
}

func (e *RabbitMQEventService) PublishEmployeeUpdated(org *models.Organization, emp *models.Employee) error {
	event := EmployeeEvent{
		EventID:      generateEventID(),
		EventType:    TopicEmployeeUpdated,
		Timestamp:    time.Now(),
		Source:       "organization-svc",
		Organization: org,
		Employee:     emp,
		Metadata: map[string]interface{}{
			"version": "1.0",
		},
	}

	return e.publishEvent(TopicEmployeeUpdated, event)
}

func (e *RabbitMQEventService) PublishEmployeeStatusChanged(org *models.Organization, emp *models.Employee, newStatus models.EmployeeStatus) error {
	event := EmployeeEvent{
		EventID:      generateEventID(),
		EventType:    TopicEmployeeStatusChanged,
		Timestamp:    time.Now(),
		Source:       "organization-svc",
		Organization: org,
		Employee:     emp,
		Metadata: map[string]interface{}{
			"version":    "1.0",
			"new_status": string(newStatus),
			"old_status": string(emp.Status),
		},
	}

	return e.publishEvent(TopicEmployeeStatusChanged, event)
}

func (e *RabbitMQEventService) PublishEmployeeRemoved(orgID, empID string) error {
	event := EmployeeEvent{
		EventID:   generateEventID(),
		EventType: TopicEmployeeRemoved,
		Timestamp: time.Now(),
		Source:    "organization-svc",
		Metadata: map[string]interface{}{
			"version":         "1.0",
			"organization_id": orgID,
			"employee_id":     empID,
		},
	}

	return e.publishEvent(TopicEmployeeRemoved, event)
}

func (e *RabbitMQEventService) PublishEmployeeRoleAssigned(org *models.Organization, emp *models.Employee, role *models.OrganizationRole) error {
	event := EmployeeEvent{
		EventID:      generateEventID(),
		EventType:    TopicEmployeeRoleAssigned,
		Timestamp:    time.Now(),
		Source:       "organization-svc",
		Organization: org,
		Employee:     emp,
		Role:         role,
		Metadata: map[string]interface{}{
			"version": "1.0",
		},
	}

	return e.publishEvent(TopicEmployeeRoleAssigned, event)
}

func (e *RabbitMQEventService) PublishEmployeeRoleRemoved(org *models.Organization, emp *models.Employee, roleID string) error {
	event := EmployeeEvent{
		EventID:      generateEventID(),
		EventType:    TopicEmployeeRoleRemoved,
		Timestamp:    time.Now(),
		Source:       "organization-svc",
		Organization: org,
		Employee:     emp,
		Metadata: map[string]interface{}{
			"version": "1.0",
			"role_id": roleID,
		},
	}

	return e.publishEvent(TopicEmployeeRoleRemoved, event)
}

// Branch events implementation

func (e *RabbitMQEventService) PublishBranchCreated(org *models.Organization, branch *models.OrganizationBranch) error {
	event := BranchEvent{
		EventID:      generateEventID(),
		EventType:    TopicBranchCreated,
		Timestamp:    time.Now(),
		Source:       "organization-svc",
		Organization: org,
		Branch:       branch,
		Metadata: map[string]interface{}{
			"version": "1.0",
		},
	}

	return e.publishEvent(TopicBranchCreated, event)
}

func (e *RabbitMQEventService) PublishBranchUpdated(org *models.Organization, branch *models.OrganizationBranch) error {
	event := BranchEvent{
		EventID:      generateEventID(),
		EventType:    TopicBranchUpdated,
		Timestamp:    time.Now(),
		Source:       "organization-svc",
		Organization: org,
		Branch:       branch,
		Metadata: map[string]interface{}{
			"version": "1.0",
		},
	}

	return e.publishEvent(TopicBranchUpdated, event)
}

func (e *RabbitMQEventService) PublishBranchDeleted(orgID, branchID string) error {
	event := BranchEvent{
		EventID:   generateEventID(),
		EventType: TopicBranchDeleted,
		Timestamp: time.Now(),
		Source:    "organization-svc",
		Metadata: map[string]interface{}{
			"version":         "1.0",
			"organization_id": orgID,
			"branch_id":       branchID,
		},
	}

	return e.publishEvent(TopicBranchDeleted, event)
}

// Role events implementation

func (e *RabbitMQEventService) PublishRoleCreated(org *models.Organization, role *models.OrganizationRole) error {
	event := RoleEvent{
		EventID:      generateEventID(),
		EventType:    TopicRoleCreated,
		Timestamp:    time.Now(),
		Source:       "organization-svc",
		Organization: org,
		Role:         role,
		Metadata: map[string]interface{}{
			"version": "1.0",
		},
	}

	return e.publishEvent(TopicRoleCreated, event)
}

func (e *RabbitMQEventService) PublishRoleUpdated(org *models.Organization, role *models.OrganizationRole) error {
	event := RoleEvent{
		EventID:      generateEventID(),
		EventType:    TopicRoleUpdated,
		Timestamp:    time.Now(),
		Source:       "organization-svc",
		Organization: org,
		Role:         role,
		Metadata: map[string]interface{}{
			"version": "1.0",
		},
	}

	return e.publishEvent(TopicRoleUpdated, event)
}

func (e *RabbitMQEventService) PublishRoleDeleted(orgID, roleID string) error {
	event := RoleEvent{
		EventID:   generateEventID(),
		EventType: TopicRoleDeleted,
		Timestamp: time.Now(),
		Source:    "organization-svc",
		Metadata: map[string]interface{}{
			"version":         "1.0",
			"organization_id": orgID,
			"role_id":         roleID,
		},
	}

	return e.publishEvent(TopicRoleDeleted, event)
}

// Invitation events implementation

func (e *RabbitMQEventService) PublishInvitationSent(org *models.Organization, invite *models.OrganizationInvite) error {
	event := InvitationEvent{
		EventID:      generateEventID(),
		EventType:    TopicInvitationSent,
		Timestamp:    time.Now(),
		Source:       "organization-svc",
		Organization: org,
		Invitation:   invite,
		Metadata: map[string]interface{}{
			"version": "1.0",
		},
	}

	return e.publishEvent(TopicInvitationSent, event)
}

func (e *RabbitMQEventService) PublishInvitationAccepted(org *models.Organization, invite *models.OrganizationInvite, emp *models.Employee) error {
	event := InvitationEvent{
		EventID:      generateEventID(),
		EventType:    TopicInvitationAccepted,
		Timestamp:    time.Now(),
		Source:       "organization-svc",
		Organization: org,
		Invitation:   invite,
		Employee:     emp,
		Metadata: map[string]interface{}{
			"version": "1.0",
		},
	}

	return e.publishEvent(TopicInvitationAccepted, event)
}

func (e *RabbitMQEventService) PublishInvitationDeclined(org *models.Organization, invite *models.OrganizationInvite) error {
	event := InvitationEvent{
		EventID:      generateEventID(),
		EventType:    TopicInvitationDeclined,
		Timestamp:    time.Now(),
		Source:       "organization-svc",
		Organization: org,
		Invitation:   invite,
		Metadata: map[string]interface{}{
			"version": "1.0",
		},
	}

	return e.publishEvent(TopicInvitationDeclined, event)
}

func (e *RabbitMQEventService) PublishInvitationCanceled(org *models.Organization, invite *models.OrganizationInvite) error {
	event := InvitationEvent{
		EventID:      generateEventID(),
		EventType:    TopicInvitationCanceled,
		Timestamp:    time.Now(),
		Source:       "organization-svc",
		Organization: org,
		Invitation:   invite,
		Metadata: map[string]interface{}{
			"version": "1.0",
		},
	}

	return e.publishEvent(TopicInvitationCanceled, event)
}

// Owner events implementation

func (e *RabbitMQEventService) PublishOwnerAdded(org *models.Organization, owner *models.OrganizationOwner) error {
	event := OwnerEvent{
		EventID:      generateEventID(),
		EventType:    TopicOwnerAdded,
		Timestamp:    time.Now(),
		Source:       "organization-svc",
		Organization: org,
		Owner:        owner,
		Metadata: map[string]interface{}{
			"version": "1.0",
		},
	}

	return e.publishEvent(TopicOwnerAdded, event)
}

func (e *RabbitMQEventService) PublishOwnerRemoved(orgID, ownerID string) error {
	event := OwnerEvent{
		EventID:   generateEventID(),
		EventType: TopicOwnerRemoved,
		Timestamp: time.Now(),
		Source:    "organization-svc",
		OwnerID:   ownerID,
		Metadata: map[string]interface{}{
			"version":         "1.0",
			"organization_id": orgID,
		},
	}

	return e.publishEvent(TopicOwnerRemoved, event)
}

// Helper methods

func (e *RabbitMQEventService) publishEvent(topic string, event interface{}) error {
	payload, err := json.Marshal(event)
	if err != nil {
		e.logger.Error("Failed to marshal event",
			zap.String("topic", topic),
			zap.Error(err))
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = e.channel.Publish(
		e.exchange, // exchange
		topic,      // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
			Timestamp:   time.Now(),
		})

	if err != nil {
		e.logger.Error("Failed to publish event",
			zap.String("topic", topic),
			zap.Error(err))
		return fmt.Errorf("failed to publish event: %w", err)
	}

	e.logger.Debug("Event published successfully",
		zap.String("topic", topic),
		zap.Int("payload_size", len(payload)))

	return nil
}

func generateEventID() string {
	return fmt.Sprintf("evt_%d", time.Now().UnixNano())
}

// NoOpEventService implements EventService but does nothing (for when RabbitMQ is not available)
type NoOpEventService struct {
	logger *zap.Logger
}

func (n *NoOpEventService) PublishOrganizationCreated(org *models.Organization) error { return nil }
func (n *NoOpEventService) PublishOrganizationUpdated(org *models.Organization) error { return nil }
func (n *NoOpEventService) PublishOrganizationStatusChanged(org *models.Organization, newStatus models.OrganizationStatus) error {
	return nil
}
func (n *NoOpEventService) PublishOrganizationDeleted(orgID string) error { return nil }
func (n *NoOpEventService) PublishOrganizationSettingsUpdated(orgID string, settings map[string]interface{}) error {
	return nil
}
func (n *NoOpEventService) PublishEmployeeAdded(org *models.Organization, emp *models.Employee) error {
	return nil
}
func (n *NoOpEventService) PublishEmployeeUpdated(org *models.Organization, emp *models.Employee) error {
	return nil
}
func (n *NoOpEventService) PublishEmployeeStatusChanged(org *models.Organization, emp *models.Employee, newStatus models.EmployeeStatus) error {
	return nil
}
func (n *NoOpEventService) PublishEmployeeRemoved(orgID, empID string) error { return nil }
func (n *NoOpEventService) PublishEmployeeRoleAssigned(org *models.Organization, emp *models.Employee, role *models.OrganizationRole) error {
	return nil
}
func (n *NoOpEventService) PublishEmployeeRoleRemoved(org *models.Organization, emp *models.Employee, roleID string) error {
	return nil
}
func (n *NoOpEventService) PublishBranchCreated(org *models.Organization, branch *models.OrganizationBranch) error {
	return nil
}
func (n *NoOpEventService) PublishBranchUpdated(org *models.Organization, branch *models.OrganizationBranch) error {
	return nil
}
func (n *NoOpEventService) PublishBranchDeleted(orgID, branchID string) error { return nil }
func (n *NoOpEventService) PublishRoleCreated(org *models.Organization, role *models.OrganizationRole) error {
	return nil
}
func (n *NoOpEventService) PublishRoleUpdated(org *models.Organization, role *models.OrganizationRole) error {
	return nil
}
func (n *NoOpEventService) PublishRoleDeleted(orgID, roleID string) error { return nil }
func (n *NoOpEventService) PublishInvitationSent(org *models.Organization, invite *models.OrganizationInvite) error {
	return nil
}
func (n *NoOpEventService) PublishInvitationAccepted(org *models.Organization, invite *models.OrganizationInvite, emp *models.Employee) error {
	return nil
}
func (n *NoOpEventService) PublishInvitationDeclined(org *models.Organization, invite *models.OrganizationInvite) error {
	return nil
}
func (n *NoOpEventService) PublishInvitationCanceled(org *models.Organization, invite *models.OrganizationInvite) error {
	return nil
}
func (n *NoOpEventService) PublishOwnerAdded(org *models.Organization, owner *models.OrganizationOwner) error {
	return nil
}
func (n *NoOpEventService) PublishOwnerRemoved(orgID, ownerID string) error { return nil }
