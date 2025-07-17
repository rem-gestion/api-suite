# Organization Service - Data Models

This document provides comprehensive documentation for all data models in the organization-svc, based on PostgreSQL migrations 0001-0008.

## Overview

The organization service implements a complete multi-tenant organization management system with:

- **Organization Management**: Core organization entities with settings and ownership
- **User & Employee Management**: Employee roles, permissions, and relationships
- **Invitation System**: Secure invitation workflow with audit trails
- **Integration System**: Third-party service integrations with event processing
- **Domain Management**: Custom domain configuration with DNS verification
- **Subscription Management**: Organization subscription and billing information

## Model Architecture

### Core Organization Models

#### Organization
The main organization entity representing a business or company.

**Key Features:**
- Multi-tenant support with unique slugs
- Status lifecycle (pending → active → inactive/suspended)
- Type classification (real_estate, property_management, etc.)
- External service references (address-svc, subscription-billing-svc)
- Soft delete support
- Complete audit trail

**External Dependencies:**
- `address_id` → address-svc.address.id
- `created_by/updated_by` → auth-identity-svc.users.id

#### OrganizationSetting
Flexible key-value configuration per organization using JSONB.

**Key Features:**
- JSONB value storage for complex configurations
- Editability control
- Audit trail for settings changes

#### OrganizationOwner
Defines ownership relationship between person and organization.

**Key Features:**
- Links to person-svc and auth-identity-svc
- Founder designation
- Ownership percentage tracking
- One owner per organization constraint

### Structure Models

#### OrganizationBranch
Physical locations/offices for an organization.

**Key Features:**
- Main branch designation (only one per organization)
- External address service integration
- Contact information (phone, email)
- Soft delete with business logic protection

#### OrganizationRole
Internal roles within organizations (Admin, Manager, Agent, Assistant).

**Key Features:**
- Default roles (non-deletable)
- Custom roles support
- Case-insensitive unique names per organization
- Soft delete protection for default roles

### Employee Management

#### Employee
Links users to organizations with employment details.

**Key Features:**
- Links to auth-identity-svc (user_id) and person-svc (person_id)
- Employment status lifecycle
- Hire/fire date tracking
- Primary role assignment
- Business logic: cannot terminate last admin

#### EmployeeRole
Many-to-many relationship between employees and roles.

**Key Features:**
- Primary role designation (one per employee)
- Soft delete for role history
- Automatic synchronization with Employee.primary_role_id

### Invitation System

#### OrganizationInvite
Secure invitation system for joining organizations.

**Key Features:**
- URL-safe token generation (32-128 chars)
- Expiration handling
- Status lifecycle (pending → accepted/rejected/expired/cancelled)
- Rate limiting (configurable per organization)
- Comprehensive metadata storage

#### OrganizationInviteLog
Complete audit trail for invitation lifecycle events.

**Key Features:**
- Partitioned by month for performance
- Actor tracking (user/system/cron)
- Security audit (IP, user agent)
- Error details for troubleshooting

### Integration System

#### IntegrationType
Global catalog of available third-party integrations.

**Key Features:**
- Configuration schema validation
- Support flags (webhook, OAuth, API key)
- Version management
- Provider categorization

#### OrganizationIntegration
Organization-specific integration configurations.

**Key Features:**
- Encrypted credentials storage (application-level)
- Sync frequency and status tracking
- API usage monitoring and rate limiting
- OAuth token management
- Soft delete support

#### OrganizationIntegrationEvent
Event queue for integration workers.

**Key Features:**
- Retry logic with configurable limits
- Scheduling support
- Status tracking (pending → processing → success/failed)
- Error message storage

### Domain Management

#### OrganizationDomain
Custom domain configuration per organization.

**Key Features:**
- SSL certificate management
- DNS verification with multiple methods
- Primary domain designation
- Redirect configuration
- Custom headers support

#### OrganizationDomainDNS
DNS record configuration and verification.

**Key Features:**
- Multiple record types (A, AAAA, CNAME, MX, TXT, NS)
- TTL configuration
- Verification tracking
- Required/optional designation

#### OrganizationDomainVerificationLog
Audit trail for domain verification attempts.

**Key Features:**
- Partitioned by month
- Multiple verification types (DNS, HTTP, file)
- Performance timing
- Error details and response data

### Subscription Management

#### OrganizationSubscription
Subscription and billing information.

**Key Features:**
- Multiple billing cycles (monthly, quarterly, yearly)
- Trial period support
- Feature and limits tracking
- Payment method integration
- Currency support

## Database Schema Features

### Advanced PostgreSQL Features

1. **Enum Types**: Extensive use of custom enum types for type safety
2. **JSONB**: Flexible schema for metadata, configuration, and settings
3. **UUID**: Primary keys for distributed system compatibility
4. **Partitioning**: Monthly partitions for log tables
5. **Triggers**: Business logic enforcement at database level
6. **Constraints**: Complex CHECK constraints for data integrity
7. **Indexes**: Optimized indexes for common query patterns

### Business Logic Protection

1. **Last Admin Protection**: Cannot delete/terminate the last admin
2. **Main Branch Protection**: Cannot delete main branch while others exist
3. **Default Role Protection**: Cannot delete default roles (Admin, Manager, etc.)
4. **Invitation Limits**: Configurable pending invitation limits
5. **Primary Constraints**: Only one primary per organization (branch, domain, role)

### Security Features

1. **Soft Delete**: Maintains data integrity while hiding records
2. **Audit Trails**: Complete change tracking with actor identification
3. **Token Security**: Cryptographically secure tokens for invitations
4. **Rate Limiting**: API usage tracking and limits
5. **Credential Encryption**: Application-level encryption for sensitive data

### External Service Integration

The models include logical foreign key references to external services:

- **auth-identity-svc**: User authentication and identity management
- **person-svc**: Personal information and contacts
- **address-svc**: Address and location data
- **subscription-billing-svc**: Billing and payment processing

These are implemented as UUID fields without database-level foreign key constraints to maintain service boundaries while preserving referential integrity at the application level.

## Usage Guidelines

### Model Initialization

```go
// Get all models for GORM migration
models := models.AllModels()

// Get only core models for minimal setup
coreModels := models.CoreModels()

// Get extended models for full features
extendedModels := models.ExtendedModels()
```

### External Service References

When working with external service references:

1. Always validate UUIDs before storing
2. Implement proper error handling for missing references
3. Use application-level integrity checks
4. Consider eventual consistency patterns

### Performance Considerations

1. **Partitioned Tables**: OrganizationInviteLog, OrganizationIntegrationEvent, OrganizationDomainVerificationLog
2. **Indexes**: Optimized for common query patterns
3. **Soft Deletes**: Use `WHERE deleted_at IS NULL` filters
4. **JSONB Queries**: Use GIN indexes for JSONB columns

### Security Best Practices

1. **Sensitive Data**: Encrypt credentials and tokens at application level
2. **Audit Trails**: Always populate created_by/updated_by fields
3. **Input Validation**: Use struct tags for validation
4. **Rate Limiting**: Implement invitation and API usage limits

This comprehensive model system provides a robust foundation for multi-tenant organization management with enterprise-grade features and security.
