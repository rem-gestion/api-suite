# Organization Service - DTOs Documentation

## 📋 Overview

This directory contains Data Transfer Objects (DTOs) for the Organization Service. DTOs define the structure of data exchanged between the client and server, ensuring type safety and validation.

## 🗂️ File Structure

```
dto/
├── common_dto.go         # Common DTOs shared across all modules
├── organization_dto.go   # Organization-related DTOs
├── employee_dto.go       # Employee management DTOs
├── branch_dto.go         # Branch/Role management DTOs
├── invitation_dto.go     # Invitation system DTOs
└── owner_dto.go          # Ownership management DTOs
```

## 🔗 External Service Integration

The DTOs are designed to integrate with the following external services via gRPC:

### Auth-Identity Service
- **UserDetails**: Complete user information
- **UserBasicDetails**: Minimal user information
- **AccountDetails**: Account information with provider details

### Person Service  
- **PersonDetails**: Complete person information (Individual/Company)
- **PersonBasicDetails**: Minimal person information
- **IndividualDetails**: Individual-specific data (FirstName, LastName, DNI)
- **CompanyDetails**: Company-specific data (LegalName, CUIT, SocietyType)
- **ContactDetails**: Contact information (Email, Phone, WhatsApp)

### Address Service
- **AddressDetails**: Complete address information with geolocation

## 📝 DTO Categories

### 1. Common DTOs (`common_dto.go`)

#### Base Structures
- **BaseResponse**: Common response fields (ID, timestamps, audit fields)
- **PaginationRequest/Response**: Pagination handling
- **FilterRequest**: Common filtering parameters
- **StatusUpdateRequest**: Status change operations

#### Utility Structures
- **ErrorResponse**: Standardized error responses
- **SuccessResponse**: Success operation responses
- **BulkRequest/Response**: Bulk operations support
- **IncludeOptions**: Control nested data inclusion

### 2. Organization DTOs (`organization_dto.go`)

#### Request DTOs
- **CreateOrganizationRequest**: Create new organization
- **UpdateOrganizationRequest**: Update organization data
- **OrganizationFiltersRequest**: Search/filter organizations

#### Response DTOs
- **OrganizationResponse**: Complete organization data
- **OrganizationListResponse**: Paginated organization list
- **OrganizationBasicResponse**: Minimal organization data
- **OrganizationSettingResponse**: Settings key-value pairs
- **OrganizationStats**: Statistical information

### 3. Employee DTOs (`employee_dto.go`)

#### Request DTOs
- **AddEmployeeRequest**: Link user to organization
- **UpdateEmployeeRequest**: Update employee data
- **EmployeeFiltersRequest**: Search/filter employees
- **AssignRoleRequest**: Assign roles to employees

#### Response DTOs
- **EmployeeResponse**: Complete employee data
- **EmployeeListResponse**: Paginated employee list
- **EmployeeBasicResponse**: Minimal employee data
- **EmployeeRoleResponse**: Employee-role relationships
- **EmployeeStats**: Employee statistics

### 4. Branch & Role DTOs (`branch_dto.go`)

#### Branch DTOs
- **CreateBranchRequest**: Create organization branch
- **UpdateBranchRequest**: Update branch data
- **OrganizationBranchResponse**: Complete branch data
- **BranchStats**: Branch statistics

#### Role DTOs
- **CreateRoleRequest**: Create custom roles
- **UpdateRoleRequest**: Update role data
- **OrganizationRoleResponse**: Complete role data
- **RoleStats**: Role usage statistics

### 5. Invitation DTOs (`invitation_dto.go`)

#### Request DTOs
- **CreateInvitationRequest**: Send organization invitations
- **UpdateInvitationRequest**: Update invitation data
- **ResendInvitationRequest**: Resend invitations
- **AcceptInvitationRequest**: Accept invitation (public)
- **DeclineInvitationRequest**: Decline invitation (public)
- **ValidateInvitationRequest**: Validate invitation token (public)

#### Response DTOs
- **OrganizationInviteResponse**: Complete invitation data
- **ValidateInvitationResponse**: Invitation validation result
- **AcceptInvitationResponse**: Acceptance confirmation
- **OrganizationInviteLogResponse**: Invitation audit logs

### 6. Owner DTOs (`owner_dto.go`)

#### Request DTOs
- **AddOwnerRequest**: Add organization owner
- **UpdateOwnerRequest**: Update ownership data
- **TransferOwnershipRequest**: Transfer ownership between owners
- **RemoveOwnerRequest**: Remove owner with redistribution

#### Response DTOs
- **OrganizationOwnerResponse**: Complete owner data
- **OwnerListResponse**: Paginated owner list with summary
- **OwnershipSummary**: Ownership distribution analysis
- **TransferOwnershipResponse**: Ownership transfer result
- **OwnershipValidationResponse**: Ownership validation results

## 🔍 Validation Rules

### Common Validation Patterns

```go
// UUID validation
validate:"required,uuid"

// Email validation
validate:"required,email"

// String length validation
validate:"required,min=2,max=100"

// Enum validation
validate:"required,oneof=active inactive suspended"

// Percentage validation
validate:"omitempty,min=0,max=100"

// URL validation
validate:"omitempty,url"

// Phone validation
validate:"omitempty,phone"
```

### Field-Specific Validations

#### Organization Fields
- **DisplayName**: 2-300 characters
- **Slug**: 2-100 characters, slug format
- **Type**: Must be valid organization type enum
- **Email**: Valid email format
- **Website**: Valid URL format
- **TimezoneID**: Valid timezone identifier

#### Employee Fields
- **Status**: active, inactive, suspended, terminated
- **UserID**: Valid UUID, must exist in auth-identity-svc
- **PersonID**: Valid UUID, must exist in person-svc
- **RoleIDs**: Array of valid UUIDs

#### Invitation Fields
- **InviteeEmail**: Valid email format
- **ExpiresAt**: Future timestamp
- **Token**: Required for public endpoints
- **Status**: pending, accepted, rejected, expired, cancelled

## 📊 Include Options

Control nested data inclusion using `IncludeOptions`:

```go
type IncludeOptions struct {
    IncludeDetails    bool // Include detailed nested objects
    IncludeStats      bool // Include statistical data
    IncludeHistorical bool // Include historical/audit data
    IncludeAddress    bool // Include address details via gRPC
    IncludePerson     bool // Include person details via gRPC
    IncludeRoles      bool // Include role information
}
```

## 🔄 gRPC Integration Examples

### User Validation in Employee Service
```go
// Validate user exists before creating employee
userDetails, err := grpcClient.GetUserByID(request.UserID)
if err != nil {
    return nil, errors.New("invalid user_id")
}

// Populate UserDetails in response
response.UserDetails = &dto.UserDetails{
    ID:            userDetails.Id,
    AccountID:     userDetails.AccountId,
    OnboardStatus: userDetails.OnboardStatus.String(),
    // ... other fields
}
```

### Person Data Integration
```go
// Fetch person details for employee
if request.IncludePerson {
    personDetails, err := grpcClient.GetPersonByID(employee.PersonID)
    if err == nil {
        response.PersonDetails = &dto.PersonDetails{
            ID:   personDetails.Id,
            Type: personDetails.Type.String(),
            // ... map other fields
        }
    }
}
```

### Address Validation
```go
// Validate fiscal address for organization
if request.FiscalAddressID != nil {
    addressDetails, err := grpcClient.GetAddress(*request.FiscalAddressID)
    if err != nil {
        return nil, errors.New("invalid fiscal_address_id")
    }
}
```

## 📈 Statistics DTOs

Each main entity includes statistics DTOs for analytics:

- **OrganizationStats**: Employee count, activity metrics
- **EmployeeStats**: Role count, performance metrics
- **BranchStats**: Employee count, property metrics
- **RoleStats**: Usage statistics, assignment metrics
- **OwnerStats**: Ownership duration, dividend information
- **InvitationStats**: Response rates, activity metrics

## 🔧 Usage Examples

### Creating an Organization
```go
request := &dto.CreateOrganizationRequest{
    Name:            "Tech Solutions Inc",
    DisplayName:     "Tech Solutions - Real Estate",
    Type:            "real_estate",
    Email:           "contact@tech-solutions.com",
    TimezoneID:      "America/New_York",
    FiscalAddressID: &addressUUID,
}
```

### Adding an Employee
```go
request := &dto.AddEmployeeRequest{
    UserID:        userUUID,
    PersonID:      &personUUID,
    RoleIDs:       []uuid.UUID{roleUUID},
    HiredAt:       &time.Now(),
    Status:        "active",
}
```

### Filtering Organizations
```go
filters := &dto.OrganizationFiltersRequest{
    FilterRequest: dto.FilterRequest{
        Status: "active",
        Search: "tech",
    },
    Type:       "real_estate",
    IsVerified: &verified,
    PaginationRequest: dto.PaginationRequest{
        Page:    1,
        PerPage: 20,
    },
    IncludeOptions: dto.IncludeOptions{
        IncludeStats:   true,
        IncludeAddress: true,
    },
}
```

## ⚡ Performance Considerations

1. **Lazy Loading**: Use `IncludeOptions` to control data fetching
2. **Pagination**: Always use pagination for list endpoints
3. **gRPC Caching**: Cache frequently accessed external data
4. **Selective Fields**: Only include necessary fields in responses
5. **Bulk Operations**: Use bulk DTOs for multiple operations

## 🛡️ Security Considerations

1. **Token Exposure**: Invitation tokens only included in admin responses
2. **PII Protection**: Personal data only included when explicitly requested
3. **Ownership Validation**: Ensure users can only access their organization data
4. **Role-Based Access**: Different response data based on user roles
5. **Audit Trails**: All modification requests include audit fields

## 📚 Next Steps

After implementing these DTOs:

1. Create Repository interfaces and implementations
2. Implement Service layer with business logic
3. Configure gRPC clients for external services
4. Implement Controllers using these DTOs
5. Add comprehensive validation middleware
6. Implement caching strategies
7. Add event publishing for state changes
