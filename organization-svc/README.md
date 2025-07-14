# Organization Service

## 📋 Tabla de Contenidos

1. [Descripción General](#descripción-general)
2. [Arquitectura del Sistema](#arquitectura-del-sistema)
3. [Funcionalidades Principales](#funcionalidades-principales)
4. [Sistema de Migraciones](#sistema-de-migraciones)
5. [Estructura de Base de Datos](#estructura-de-base-de-datos)
6. [API y Endpoints](#api-y-endpoints)
7. [Integración con Otros Servicios](#integración-con-otros-servicios)
8. [Configuración y Deployment](#configuración-y-deployment)
9. [Monitoreo y Métricas](#monitoreo-y-métricas)
10. [Seguridad y Validaciones](#seguridad-y-validaciones)
11. [Ejemplos de Uso](#ejemplos-de-uso)
12. [Desarrollo y Contribución](#desarrollo-y-contribución)
13. [Troubleshooting](#troubleshooting)

---

## 📖 Descripción General

El **Organization Service** es el microservicio central del sistema de gestión inmobiliaria REM (Real Estate Management), diseñado para manejar la estructura organizacional completa de empresas inmobiliarias. Proporciona una **arquitectura evolutiva** que se construye en 8 migraciones progresivas, desde fundamentos básicos hasta características avanzadas de escalabilidad y mantenimiento.

### 🎯 Objetivos del Servicio

- **Gestión Completa de Organizaciones**: Desde la creación básica hasta estructuras complejas
- **Escalabilidad Empresarial**: Soporte para múltiples sucursales, empleados y roles
- **Integraciones Externas**: Conectividad con servicios de terceros (CRM, contabilidad, marketing)
- **Dominios Personalizados**: Gestión de dominios propios con verificación DNS automática
- **Seguridad y Compliance**: Auditoría completa y soft-delete obligatorio
- **Performance Optimizada**: Índices inteligentes y particionado proactivo

### 🏗️ Características Principales

- **🏢 Núcleo Organizacional**: Gestión de empresas con configuraciones flexibles
- **👥 Estructura Jerárquica**: Sucursales, roles y empleados con asignaciones múltiples
- **� Sistema de Invitaciones**: Workflow completo con tokens seguros y expiración automática
- **🔗 Hub de Integraciones**: Catálogo extensible de conectores externos
- **🌐 Dominios Personalizados**: Verificación DNS/SSL automática y gestión de certificados
- **📊 Auditoría y Logs**: Tracking completo con logs particionados por performance
- **🔒 Seguridad Avanzada**: Validaciones multi-capa y protecciones de datos críticos

## 🏗️ Arquitectura del Sistema

### 🎯 Visión General

El Organization Service implementa una **arquitectura de microservicio** con las siguientes capas:

```mermaid
graph TB
    subgraph "API Layer"
        REST[REST API]
        GRPC[gRPC Interface]
        EVENTS[Event Handlers]
    end
    
    subgraph "Business Logic"
        ORG_SVC[Organization Service]
        EMP_SVC[Employee Service]
        INVITE_SVC[Invitation Service]
        INTEGRATION_SVC[Integration Service]
        DOMAIN_SVC[Domain Service]
    end
    
    subgraph "Data Layer"
        POSTGRES[(PostgreSQL)]
        REDIS[(Redis Cache)]
        RABBIT[RabbitMQ]
    end
    
    subgraph "External Services"
        AUTH[auth-identity-svc]
        ADDRESS[address-svc]
        PERSON[person-svc]
        NOTIFICATION[notification-svc]
    end
    
    REST --> ORG_SVC
    GRPC --> ORG_SVC
    EVENTS --> ORG_SVC
    
    ORG_SVC --> EMP_SVC
    ORG_SVC --> INVITE_SVC
    ORG_SVC --> INTEGRATION_SVC
    ORG_SVC --> DOMAIN_SVC
    
    ORG_SVC --> POSTGRES
    ORG_SVC --> REDIS
    ORG_SVC --> RABBIT
    
    ORG_SVC --> AUTH
    ORG_SVC --> ADDRESS
    ORG_SVC --> PERSON
    ORG_SVC --> NOTIFICATION
    
    style REST fill:#e3f2fd,stroke:#1565c0
    style POSTGRES fill:#e8f5e8,stroke:#2e7d32
    style REDIS fill:#ffebee,stroke:#c62828
    style RABBIT fill:#fff3e0,stroke:#ef6c00
```

### 🔄 Flujo de Datos Principal

```mermaid
sequenceDiagram
    participant Client
    participant API
    participant OrgService
    participant Database
    participant EventBus
    participant ExternalServices
    
    Client->>API: Create Organization
    API->>OrgService: ProcessOrganizationCreation
    OrgService->>Database: Insert organization
    Database->>OrgService: Return organization_id
    OrgService->>Database: Create default settings
    OrgService->>Database: Create main branch
    OrgService->>Database: Create default roles
    OrgService->>EventBus: Publish OrganizationCreated
    EventBus->>ExternalServices: Notify dependent services
    OrgService->>API: Return organization details
    API->>Client: 201 Created
```

---

## 🚀 Funcionalidades Principales

### 1. 🏢 Gestión de Organizaciones

#### **Núcleo Organizacional**
- **Creación y configuración** de empresas inmobiliarias
- **Settings flexibles** almacenados como JSONB
- **Propietarios múltiples** con porcentajes de participación
- **Soft delete obligatorio** con protecciones automáticas

#### **Configuraciones Avanzadas**
- Configuraciones por defecto automáticas al crear organización
- Validación de porcentajes de propiedad (máximo 100%)
- Protección contra eliminación directa (hard delete)
- Auditoría completa de cambios

### 2. 🏘️ Estructura Organizacional

#### **Sucursales (Branches)**
- **Sucursal principal única** por organización
- **Sucursales múltiples** con información de contacto
- **Integración con address-svc** para geolocalización
- **Protecciones** contra eliminación de sucursal principal

#### **Sistema de Roles**
- **Roles por defecto**: Admin, Manager, Agent, Assistant
- **Roles personalizados** según necesidades específicas
- **Protección de roles críticos** (no eliminables)
- **Jerarquía flexible** de permisos

### 3. 👥 Gestión de Empleados

#### **Empleados y Asignaciones**
- **Empleados asignados** a sucursales específicas
- **Roles múltiples** por empleado con uno primario
- **Estados de empleado** (activo, terminado, suspendido)
- **Metadata flexible** para información adicional

#### **Sistema de Roles Múltiples**
- Un empleado puede tener **múltiples roles**
- **Rol primario** obligatorio para jerarquía
- **Fechas de vigencia** y revocación de roles
- **Constraints automáticos** para integridad

### 4. 📧 Sistema de Invitaciones

#### **Workflow Completo**
- **Tokens seguros** generados automáticamente
- **Expiración configurable** por organización
- **Límites dinámicos** de invitaciones pendientes
- **Estados de seguimiento** (pending, accepted, expired, cancelled)

#### **Automatización**
- **Notificaciones automáticas** via notification-svc
- **Cleanup automático** de invitaciones expiradas
- **Conversión automática** a empleados al aceptar
- **Logs detallados** de todo el flujo

### 5. 🔗 Hub de Integraciones

#### **Catálogo Extensible**
- **Tipos predefinidos**: CRM, Contabilidad, Marketing, ERP
- **Configuración flexible** con schemas JSON
- **Credenciales seguras** encriptadas en JSONB
- **Soporte para APIs** y webhooks

#### **Integraciones Disponibles**
- **QuickBooks Online**: Sincronización contable bidireccional
- **HubSpot CRM**: Gestión de clientes, leads y pipeline de ventas
- **Mailchimp**: Campañas de marketing y automatización de emails
- **Zapier**: Automatizaciones personalizadas y conectores adicionales
- **Stripe/PayPal**: Procesamiento de pagos (próximamente)
- **DocuSign**: Firma electrónica de contratos (próximamente)

#### **Sistema de Eventos**
- **Queue de eventos** para workers externos con retry automático
- **Logs particionados** por rendimiento y mantenimiento automático
- **Retry automático** con backoff exponencial y límites configurables
- **Monitoreo de salud** de integraciones con alertas proactivas
- **Event sourcing** para auditoría completa de acciones críticas

### 6. 🌐 Dominios Personalizados

#### **Gestión de Dominios**
- **Verificación DNS automática** con records TXT
- **Configuración SSL** automática
- **Dominio primario único** por organización
- **Redirecciones inteligentes** entre dominios

#### **Verificación Automática**
- **Records DNS requeridos** creados automáticamente
- **Verificación periódica** de configuración
- **Estado de SSL** y certificados
- **Logs detallados** de verificación

---

## 📊 Sistema de Migraciones

El servicio utiliza un **sistema evolutivo de 8 migraciones** que construye progresivamente la funcionalidad completa. Para documentación técnica detallada de cada migración, consultar:

### 📚 **[Documentación Completa de Migraciones](migrations/README.md)**

*Esta documentación contiene todos los detalles técnicos, diagramas avanzados, funciones específicas, y guías de implementación para cada migración individual.*

### 🏗️ Resumen de Fases

| Fase | Migraciones | Descripción | Componentes |
|------|------------|-------------|-------------|
| **Fundamentos** | 0001-0003 | Base del sistema | ENUMs, funciones, tablas core, estructura |
| **RRHH** | 0004-0005 | Gestión de personal | Empleados, roles, invitaciones |
| **Conectividad** | 0006-0007 | Servicios externos | Integraciones, dominios, DNS |
| **Optimización** | 0008 | Performance y escalabilidad | Particionado, vistas, mantenimiento |

### 📈 Estadísticas del Sistema

- **12 tablas principales** con auditoría completa y soft-delete obligatorio
- **62 funciones especializadas** para business logic y validaciones
- **64 triggers automáticos** para validaciones, auditoría e integridad
- **97 índices optimizados** para performance (incluyendo índices condicionales)
- **9 ENUMs** para tipos de datos consistentes y extensibles
- **35 constraints** para protección de datos críticos y reglas de negocio

## 🗃️ Estructura de Base de Datos

### 📊 Diagrama Entidad-Relación Global

```mermaid
erDiagram
    %% Núcleo Organizacional (0001-0003)
    organization {
        UUID id PK
        VARCHAR display_name
        UUID fiscal_address_id FK
        organization_status_enum status
        TIMESTAMPTZ created_at
        UUID created_by FK
        TIMESTAMPTZ updated_at
        UUID updated_by FK
        TIMESTAMPTZ deleted_at
    }
    
    organization_settings {
        UUID organization_id PK,FK
        VARCHAR setting_key PK
        JSONB setting_value
        TIMESTAMPTZ updated_at
        UUID updated_by FK
    }
    
    organization_owner {
        UUID id PK
        UUID organization_id FK
        owner_type_enum owner_type
        UUID owner_id FK
        DECIMAL ownership_percentage
        TIMESTAMPTZ deleted_at
    }
    
    organization_branch {
        UUID id PK
        UUID organization_id FK
        VARCHAR display_name
        UUID address_id FK
        VARCHAR phone
        VARCHAR email
        BOOLEAN is_main
        TIMESTAMPTZ deleted_at
    }
    
    organization_role {
        UUID id PK
        UUID organization_id FK
        VARCHAR name
        TEXT description
        BOOLEAN is_default
        TIMESTAMPTZ deleted_at
    }
    
    %% Gestión de Empleados (0004-0005)
    organization_employee {
        UUID id PK
        UUID organization_id FK
        UUID person_id FK
        UUID branch_id FK
        employee_status_enum status
        DATE hire_date
        DECIMAL salary
        TIMESTAMPTZ deleted_at
    }
    
    organization_employee_role {
        UUID id PK
        UUID employee_id FK
        UUID role_id FK
        BOOLEAN is_primary
        DATE effective_date
        DATE revoked_date
    }
    
    organization_invitation {
        UUID id PK
        UUID organization_id FK
        VARCHAR email
        invitation_status_enum status
        VARCHAR token
        TIMESTAMPTZ expires_at
        TIMESTAMPTZ deleted_at
    }
    
    %% Integraciones y Dominios (0006-0007)
    organization_integration {
        UUID id PK
        UUID organization_id FK
        integration_type_enum integration_type
        VARCHAR integration_name
        JSONB credentials
        JSONB configuration
        TIMESTAMPTZ deleted_at
    }
    
    organization_integration_log {
        UUID id PK
        UUID integration_id FK
        integration_log_level_enum log_level
        JSONB log_data
        TIMESTAMPTZ created_at
    }
    
    organization_custom_domain {
        UUID id PK
        UUID organization_id FK
        VARCHAR domain_name
        domain_status_enum status
        BOOLEAN ssl_enabled
        BOOLEAN is_primary
        TIMESTAMPTZ deleted_at
    }
    
    %% Relaciones principales
    organization ||--o{ organization_settings : "configura"
    organization ||--o{ organization_owner : "posee"
    organization ||--o{ organization_branch : "tiene"
    organization ||--o{ organization_role : "define"
    organization ||--o{ organization_employee : "emplea"
    organization ||--o{ organization_invitation : "invita"
    organization ||--o{ organization_integration : "integra"
    organization ||--o{ organization_custom_domain : "personaliza"
    
    organization_branch ||--o{ organization_employee : "asigna"
    organization_role ||--o{ organization_employee_role : "asigna"
    organization_employee ||--o{ organization_employee_role : "tiene"
    organization_integration ||--o{ organization_integration_log : "registra"
```

### 🏗️ Principales Entidades y Relaciones

#### **Núcleo Organizacional**
- **organization**: Entidad principal con soft-delete obligatorio
- **organization_settings**: Configuraciones KV con soporte JSONB
- **organization_owner**: Propietarios con validación de porcentajes
- **organization_branch**: Sucursales con protección de sucursal principal
- **organization_role**: Roles con protección de roles por defecto

#### **Gestión de Personal** 
- **organization_employee**: Empleados con asignación a sucursales
- **organization_employee_role**: Roles múltiples con fechas de vigencia
- **organization_invitation**: Sistema de invitaciones con tokens seguros

#### **Integraciones y Personalización**
- **organization_integration**: Configuraciones de APIs externas
- **organization_integration_log**: Logs particionados por rendimiento
- **organization_custom_domain**: Dominios personalizados con verificación DNS

### 🔧 Características Técnicas

#### **Patrones de Diseño**
- **Soft Delete Universal**: Todas las tablas principales con protección automática contra hard deletes
- **Auditoría Completa**: created_at, created_by, updated_at, updated_by con triggers automáticos
- **Foreign Keys Lógicas**: Sin constraints DB para independencia total entre microservicios
- **Validaciones Multi-Capa**: DB constraints + triggers + application logic para máxima robustez

#### **Optimizaciones de Performance**
- **Índices Condicionales**: Solo registros activos (deleted_at IS NULL) para máxima eficiencia
- **Particionado Mensual**: Tablas de logs por fecha con creación automática
- **JSONB Indexing**: Índices GIN para consultas JSON eficientes en configuraciones
- **Unique Constraints Inteligentes**: Evitan duplicados críticos (sucursal principal, roles por defecto)

---

## 🔌 API y Endpoints

### 🎯 REST API Principal

#### **Organizations**
```http
POST   /api/v1/organizations                    # Crear organización
GET    /api/v1/organizations                    # Listar organizaciones
GET    /api/v1/organizations/{id}               # Obtener organización
PUT    /api/v1/organizations/{id}               # Actualizar organización
DELETE /api/v1/organizations/{id}               # Soft delete organización

GET    /api/v1/organizations/{id}/settings      # Obtener configuraciones
PUT    /api/v1/organizations/{id}/settings      # Actualizar configuraciones
```

#### **Branches (Sucursales)**
```http
POST   /api/v1/organizations/{orgId}/branches          # Crear sucursal
GET    /api/v1/organizations/{orgId}/branches          # Listar sucursales
GET    /api/v1/organizations/{orgId}/branches/{id}     # Obtener sucursal
PUT    /api/v1/organizations/{orgId}/branches/{id}     # Actualizar sucursal
DELETE /api/v1/organizations/{orgId}/branches/{id}     # Eliminar sucursal
PUT    /api/v1/organizations/{orgId}/branches/{id}/main # Marcar como principal
```

#### **Employees (Empleados)**
```http
POST   /api/v1/organizations/{orgId}/employees         # Crear empleado
GET    /api/v1/organizations/{orgId}/employees         # Listar empleados
GET    /api/v1/organizations/{orgId}/employees/{id}    # Obtener empleado
PUT    /api/v1/organizations/{orgId}/employees/{id}    # Actualizar empleado
DELETE /api/v1/organizations/{orgId}/employees/{id}    # Soft delete empleado

POST   /api/v1/employees/{id}/roles                    # Asignar rol
DELETE /api/v1/employees/{id}/roles/{roleId}           # Revocar rol
PUT    /api/v1/employees/{id}/roles/{roleId}/primary   # Marcar rol primario
```

#### **Invitations (Invitaciones)**
```http
POST   /api/v1/organizations/{orgId}/invitations       # Crear invitación
GET    /api/v1/organizations/{orgId}/invitations       # Listar invitaciones
GET    /api/v1/invitations/{token}                     # Obtener por token
POST   /api/v1/invitations/{token}/accept              # Aceptar invitación
DELETE /api/v1/invitations/{id}                        # Cancelar invitación
```

#### **Integrations (Integraciones)**
```http
GET    /api/v1/integration-types                       # Catálogo de tipos
POST   /api/v1/organizations/{orgId}/integrations      # Crear integración  
GET    /api/v1/organizations/{orgId}/integrations      # Listar integraciones
PUT    /api/v1/integrations/{id}                       # Actualizar configuración
DELETE /api/v1/integrations/{id}                       # Eliminar integración
POST   /api/v1/integrations/{id}/test                  # Probar conexión
GET    /api/v1/integrations/{id}/logs                  # Obtener logs
```

#### **Custom Domains (Dominios)**
```http
POST   /api/v1/organizations/{orgId}/domains           # Agregar dominio
GET    /api/v1/organizations/{orgId}/domains           # Listar dominios
GET    /api/v1/domains/{id}/dns-records                # Ver records DNS requeridos
POST   /api/v1/domains/{id}/verify                     # Verificar configuración DNS
PUT    /api/v1/domains/{id}/primary                    # Marcar como primario
DELETE /api/v1/domains/{id}                            # Eliminar dominio
```

### 🔄 gRPC Interface

```protobuf
service OrganizationService {
  // Organizations
  rpc CreateOrganization(CreateOrganizationRequest) returns (OrganizationResponse);
  rpc GetOrganization(GetOrganizationRequest) returns (OrganizationResponse);
  rpc UpdateOrganization(UpdateOrganizationRequest) returns (OrganizationResponse);
  rpc DeleteOrganization(DeleteOrganizationRequest) returns (Empty);
  
  // Employees
  rpc CreateEmployee(CreateEmployeeRequest) returns (EmployeeResponse);
  rpc GetEmployee(GetEmployeeRequest) returns (EmployeeResponse);
  rpc AssignEmployeeRole(AssignRoleRequest) returns (EmployeeRoleResponse);
  
  // Internal use for other services
  rpc ValidateOrganizationExists(ValidateOrgRequest) returns (ValidateOrgResponse);
  rpc GetOrganizationSettings(GetSettingsRequest) returns (SettingsResponse);
}
```

### 📡 Event System

#### **Published Events**
```json
{
  "event_type": "OrganizationCreated",
  "organization_id": "uuid",
  "timestamp": "2025-07-14T10:00:00Z",
  "data": {
    "display_name": "Inmobiliaria Los Pinos",
    "created_by": "user-uuid",
    "default_settings_created": true
  }
}
```

```json
{
  "event_type": "EmployeeInvited", 
  "organization_id": "uuid",
  "invitation_id": "uuid",
  "data": {
    "email": "nuevo@empleado.com",
    "role_id": "uuid",
    "expires_at": "2025-07-21T10:00:00Z"
  }
}
```

#### **Consumed Events**
- `UserCreated` (from auth-identity-svc): Para crear empleados automáticamente
- `AddressUpdated` (from address-svc): Actualizar direcciones de sucursales
- `PersonUpdated` (from person-svc): Sincronizar datos de empleados

---

## 🔗 Integración con Otros Servicios

### 🏛️ Arquitectura de Microservicios

```mermaid
graph TB
    subgraph "Organization Service"
        ORG_API[Organization API]
        ORG_CORE[Core Logic]
        ORG_DB[(PostgreSQL)]
    end
    
    subgraph "Identity & Auth"
        AUTH_SVC[auth-identity-svc]
        JWT[JWT Validation]
        RBAC[RBAC Engine]
    end
    
    subgraph "Supporting Services"
        ADDRESS_SVC[address-svc]
        PERSON_SVC[person-svc]
        NOTIFICATION_SVC[notification-svc]
    end
    
    subgraph "Infrastructure"
        GATEWAY[API Gateway]
        EVENT_BUS[Event Bus - RabbitMQ]
        CACHE[Redis Cache]
    end
    
    GATEWAY --> ORG_API
    ORG_API --> JWT
    ORG_CORE --> AUTH_SVC
    ORG_CORE --> ADDRESS_SVC
    ORG_CORE --> PERSON_SVC
    ORG_CORE --> NOTIFICATION_SVC
    ORG_CORE --> EVENT_BUS
    ORG_CORE --> CACHE
    ORG_CORE --> ORG_DB
    
    style ORG_API fill:#e3f2fd,stroke:#1565c0
    style ORG_CORE fill:#e8f5e8,stroke:#2e7d32
    style ORG_DB fill:#fff3e0,stroke:#ef6c00
```

### 🔄 Service Dependencies

#### **auth-identity-svc**
- **Validación de usuarios**: created_by, updated_by fields
- **Autorización**: Verificar permisos para operaciones
- **JWT tokens**: Validación de autenticación

#### **address-svc**
- **Direcciones de sucursales**: fiscal_address_id en organizations
- **Geolocalización**: address_id en organization_branch
- **Validación**: Verificar existencia de direcciones

#### **person-svc**
- **Datos de empleados**: person_id en organization_employee
- **Propietarios**: owner_id en organization_owner
- **Contactos**: Información personal para invitaciones

#### **notification-svc**
- **Emails de invitación**: Workflow de invitaciones
- **Notificaciones**: Cambios importantes en organización
- **Alerts**: Integraciones fallidas, dominios no verificados

### 📨 Event-Driven Communication

#### **Outbound Events**
```yaml
OrganizationCreated:
  triggers:
    - notification-svc: Enviar email de bienvenida
    - property-svc: Crear namespace para propiedades
    - billing-svc: Crear cuenta de facturación

EmployeeAdded:
  triggers:
    - auth-identity-svc: Crear credenciales si no existen
    - notification-svc: Email de incorporación

IntegrationActivated:
  triggers:
    - integration-workers: Iniciar sincronización
    - monitoring-svc: Crear métricas
```

#### **Inbound Events**
```yaml
UserRegistered:
  action: Verificar si hay invitación pendiente
  
AddressValidated:
  action: Actualizar estado de verificación de sucursal
  
PersonUpdated:
  action: Sincronizar cambios en empleados
```

---

## ⚙️ Configuración y Deployment

### 🔧 Variables de Entorno

```bash
# Database
DATABASE_URL=postgresql://user:pass@localhost:5432/organization_db
DATABASE_POOL_SIZE=20
DATABASE_TIMEOUT=30s

# Redis Cache
REDIS_URL=redis://localhost:6379
REDIS_DB=0
CACHE_TTL=3600

# Message Queue
RABBITMQ_URL=amqp://user:pass@localhost:5672
QUEUE_PREFETCH_COUNT=10

# External Services
AUTH_SERVICE_URL=http://auth-identity-svc:8080
ADDRESS_SERVICE_URL=http://address-svc:8080
PERSON_SERVICE_URL=http://person-svc:8080
NOTIFICATION_SERVICE_URL=http://notification-svc:8080

# Service Configuration
SERVER_PORT=8080
GRPC_PORT=9090
LOG_LEVEL=info
ENVIRONMENT=production

# Feature Flags
ENABLE_DOMAIN_VERIFICATION=true
ENABLE_INTEGRATION_HUB=true
MAX_INVITATIONS_PER_ORG=50

# Security
JWT_SECRET=your-secret-key
ENCRYPTION_KEY=32-char-encryption-key
CORS_ORIGINS=https://app.rem-system.com
```

### 🐳 Docker Configuration

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o organization-service ./cmd/api

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/organization-service .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080 9090
CMD ["./organization-service"]
```

### 🚀 Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: organization-service
  labels:
    app: organization-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: organization-service
  template:
    metadata:
      labels:
        app: organization-service
    spec:
      containers:
      - name: organization-service
        image: rem/organization-service:latest
        ports:
        - containerPort: 8080
          name: http
        - containerPort: 9090
          name: grpc
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: organization-secrets
              key: database-url
        - name: REDIS_URL
          valueFrom:
            configMapKeyRef:
              name: organization-config
              key: redis-url
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

### 📊 Database Migrations

```bash
# Ejecutar migraciones en orden
migrate -path ./migrations/pg -database $DATABASE_URL up

# Verificar estado
migrate -path ./migrations/pg -database $DATABASE_URL version

# Rollback si es necesario
migrate -path ./migrations/pg -database $DATABASE_URL down 1
```

---

## 📈 Monitoreo y Métricas

### 🔍 Health Checks

```http
GET /health      # Health check básico
GET /ready       # Readiness check (DB, cache, deps)
GET /metrics     # Prometheus metrics
```

### 📊 Métricas Principales

#### **Business Metrics**
```prometheus
# Organizaciones
organization_total{status="active"}
organization_created_total
organization_deleted_total

# Empleados
employee_total{organization_id, status}
employee_created_total
employee_terminated_total

# Invitaciones
invitation_sent_total
invitation_accepted_total
invitation_expired_total

# Integraciones
integration_active_total{type}
integration_sync_success_total{integration_id}
integration_sync_error_total{integration_id}
```

#### **Technical Metrics**
```prometheus
# Performance
http_request_duration_seconds{method, endpoint}
grpc_request_duration_seconds{method}
database_query_duration_seconds{query_type}

# Resources
memory_usage_bytes
cpu_usage_percent
database_connections_active
cache_hit_ratio

# Events
events_published_total{event_type}
events_consumed_total{event_type}
events_failed_total{event_type}
```

### 🚨 Alerting Rules

```yaml
- alert: HighErrorRate
  expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.1
  for: 5m
  annotations:
    summary: "High error rate in organization service"

- alert: DatabaseConnectionsHigh
  expr: database_connections_active > 15
  for: 2m
  annotations:
    summary: "Database connection pool exhaustion"

- alert: IntegrationSyncFailing
  expr: increase(integration_sync_error_total[1h]) > 10
  for: 0m
  annotations:
    summary: "Integration sync failures detected"
```

### 📝 Structured Logging

```json
{
  "timestamp": "2025-07-14T10:00:00Z",
  "level": "info",
  "service": "organization-service",
  "trace_id": "abc123",
  "organization_id": "uuid",
  "user_id": "uuid",
  "operation": "create_employee",
  "duration_ms": 150,
  "message": "Employee created successfully"
}
```
## 🔒 Seguridad y Validaciones

### 🛡️ Principios de Seguridad

#### **1. Soft Delete Obligatorio**
```sql
-- Patrón universal aplicado en todas las tablas principales
deleted_at TIMESTAMP WITH TIME ZONE

-- Con constraints que aseguran consistencia
CONSTRAINT chk_soft_delete_status CHECK (
    (deleted_at IS NULL AND status != 'deleted') OR 
    (deleted_at IS NOT NULL AND status = 'deleted')
)
```

#### **2. Auditoría Completa**
```sql
-- Campos estándar en todas las tablas
created_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp_utc(),
created_by    UUID NOT NULL,
updated_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp_utc(),
updated_by    UUID
```

#### **3. Foreign Keys Lógicas**
```sql
-- Evitar dependencias circulares entre microservicios
fiscal_address_id   UUID,        -- FK lógica → address-svc
created_by          UUID NOT NULL, -- FK lógica → auth-identity-svc
owner_id            UUID NOT NULL  -- FK lógica → person-svc
```

### 🔐 Validaciones Multi-Capa

#### **Base de Datos**
- **Constraints**: Validaciones de integridad críticas
- **Triggers**: Business rules automatizadas
- **Functions**: Validaciones complejas en PL/pgSQL

#### **Aplicación**
- **Input validation**: Sanitización de entrada
- **Business logic**: Reglas de negocio específicas
- **Authorization**: Verificación de permisos

#### **API Gateway**
- **Rate limiting**: Por organización y usuario
- **JWT validation**: Tokens válidos y no expirados
- **CORS**: Configuración de orígenes permitidos

### 🔑 Gestión de Credenciales

#### **Encriptación de Integraciones**
```go
// Credenciales encriptadas antes de almacenar
type IntegrationCredentials struct {
    APIKey       string `json:"api_key,omitempty"`
    ClientID     string `json:"client_id,omitempty"`
    ClientSecret string `json:"client_secret,omitempty"`
    AccessToken  string `json:"access_token,omitempty"`
}

// Almacenado como JSONB encriptado
{
  "api_key": "encrypted:AES256:base64encodedcontent",
  "client_secret": "encrypted:AES256:anothersecret"
}
```

#### **Tokens de Invitación**
```sql
-- Generación segura de tokens
CREATE OR REPLACE FUNCTION generate_invite_token(token_length INTEGER DEFAULT 64)
RETURNS VARCHAR AS $$
DECLARE
    chars TEXT := 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_';
    result TEXT := '';
    i INTEGER;
BEGIN
    FOR i IN 1..token_length LOOP
        result := result || substr(chars, floor(random() * length(chars) + 1)::integer, 1);
    END LOOP;
    RETURN result;
END;
$$ LANGUAGE plpgsql;
```

### 🚫 Protecciones Automáticas

#### **Prevención de Hard Delete**
```sql
CREATE OR REPLACE FUNCTION prevent_hard_delete_organization()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'DELETE' AND TG_TABLE_SCHEMA = 'public' THEN
        RAISE EXCEPTION 'Direct DELETE not allowed on %. Use soft delete by setting deleted_at timestamp.', TG_TABLE_NAME;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;
```

#### **Validación de Roles Críticos**
```sql
CREATE OR REPLACE FUNCTION prevent_delete_default_roles()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.is_default = true THEN
        RAISE EXCEPTION 'Cannot delete default role: %', OLD.name;
    END IF;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;
```

#### **Protección de Sucursal Principal**
```sql
CREATE OR REPLACE FUNCTION prevent_delete_only_main_branch()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.is_main = true THEN
        IF (SELECT COUNT(*) FROM organization_branch 
            WHERE organization_id = OLD.organization_id 
            AND id != OLD.id 
            AND deleted_at IS NULL) > 0 THEN
            RAISE EXCEPTION 'Cannot delete the main branch while other branches exist.';
        END IF;
    END IF;
    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;
```

---

## 💡 Ejemplos de Uso

### 📝 Casos de Uso Principales

#### **1. Crear Nueva Organización Completa**

```go
// 1. Crear organización
orgRequest := &CreateOrganizationRequest{
    DisplayName: "Inmobiliaria Los Pinos",
    LogoUrl:     "https://example.com/logo.png",
    Matricula:   "INM-2024-001",
    CreatedBy:   userID,
}

org, err := orgService.CreateOrganization(ctx, orgRequest)
if err != nil {
    return err
}

// El sistema automáticamente:
// - Crea configuraciones por defecto
// - Crea sucursal principal
// - Crea roles básicos (Admin, Manager, Agent, Assistant)
// - Publica evento OrganizationCreated
```

#### **2. Gestión de Empleados y Roles**

```go
// Crear empleado
employee, err := orgService.CreateEmployee(ctx, &CreateEmployeeRequest{
    OrganizationId: orgID,
    PersonId:       personID,
    BranchId:       branchID,
    HireDate:       time.Now(),
    CreatedBy:      userID,
})

// Asignar múltiples roles
roles := []string{adminRoleID, agentRoleID}
for i, roleID := range roles {
    _, err := orgService.AssignEmployeeRole(ctx, &AssignRoleRequest{
        EmployeeId: employee.Id,
        RoleId:     roleID,
        IsPrimary:  i == 0, // Primer rol es primario
        CreatedBy:  userID,
    })
}
```

#### **3. Sistema de Invitaciones**

```go
// Crear invitación
invitation, err := orgService.CreateInvitation(ctx, &CreateInvitationRequest{
    OrganizationId: orgID,
    Email:          "nuevo@empleado.com",
    RoleId:         agentRoleID,
    BranchId:       branchID,
    ExpiresInDays:  7,
    CreatedBy:      userID,
})

// El sistema automáticamente:
// - Genera token seguro único
// - Verifica límites de invitaciones
// - Programa envío de email
// - Logs la acción
```

#### **4. Configurar Integración Externa**

```go
// Crear integración con HubSpot
integration, err := orgService.CreateIntegration(ctx, &CreateIntegrationRequest{
    OrganizationId: orgID,
    IntegrationType: "hubspot_crm",
    Configuration: map[string]interface{}{
        "api_key":   "encrypted-api-key",
        "portal_id": "12345",
    },
    CreatedBy: userID,
})

// Activar sincronización
err = orgService.ActivateIntegration(ctx, integration.Id)
```

#### **5. Gestión de Dominios Personalizados**

```go
// Agregar dominio personalizado
domain, err := orgService.CreateCustomDomain(ctx, &CreateDomainRequest{
    OrganizationId: orgID,
    DomainName:     "lospinos.com",
    CreatedBy:      userID,
})

// El sistema automáticamente:
// - Crea records DNS requeridos
// - Programa verificación periódica
// - Configura SSL cuando sea verificado
```

### 🔄 Workflows Complejos

#### **Workflow: Onboarding de Nueva Organización**

```go
func OnboardNewOrganization(ctx context.Context, req *OnboardingRequest) error {
    // 1. Crear organización
    org, err := createOrganization(ctx, req)
    if err != nil {
        return err
    }
    
    // 2. Configurar sucursal principal con dirección
    branch, err := createMainBranch(ctx, org.Id, req.Address)
    if err != nil {
        return err
    }
    
    // 3. Crear primer empleado admin
    employee, err := createFirstEmployee(ctx, org.Id, branch.Id, req.OwnerPersonId)
    if err != nil {
        return err
    }
    
    // 4. Asignar rol de administrador
    err = assignAdminRole(ctx, employee.Id, org.Id)
    if err != nil {
        return err
    }
    
    // 5. Configurar integraciones básicas si se especifican
    if len(req.Integrations) > 0 {
        err = setupInitialIntegrations(ctx, org.Id, req.Integrations)
        if err != nil {
            log.Warn("Failed to setup some integrations", "error", err)
        }
    }
    
    // 6. Configurar dominio personalizado si se proporciona
    if req.CustomDomain != "" {
        err = setupCustomDomain(ctx, org.Id, req.CustomDomain)
        if err != nil {
            log.Warn("Failed to setup custom domain", "error", err)
        }
    }
    
    return nil
}
```

#### **Workflow: Invitación y Onboarding de Empleado**

```go
func InviteAndOnboardEmployee(ctx context.Context, req *InviteEmployeeRequest) error {
    // 1. Validar límites de invitaciones
    err := validateInvitationLimits(ctx, req.OrganizationId)
    if err != nil {
        return err
    }
    
    // 2. Crear invitación
    invitation, err := createInvitation(ctx, req)
    if err != nil {
        return err
    }
    
    // 3. Enviar email de invitación (async)
    err = publishInvitationEvent(ctx, invitation)
    if err != nil {
        log.Error("Failed to publish invitation event", "error", err)
    }
    
    // 4. Configurar recordatorio automático
    err = scheduleInvitationReminder(ctx, invitation.Id)
    if err != nil {
        log.Warn("Failed to schedule reminder", "error", err)
    }
    
    return nil
}

// Callback cuando se acepta invitación
func OnInvitationAccepted(ctx context.Context, token string, personId string) error {
    // 1. Validar token y obtener invitación
    invitation, err := getInvitationByToken(ctx, token)
    if err != nil {
        return err
    }
    
    // 2. Crear empleado automáticamente
    employee, err := createEmployeeFromInvitation(ctx, invitation, personId)
    if err != nil {
        return err
    }
    
    // 3. Asignar rol especificado en la invitación
    err = assignRoleFromInvitation(ctx, employee.Id, invitation)
    if err != nil {
        return err
    }
    
    // 4. Marcar invitación como aceptada
    err = markInvitationAccepted(ctx, invitation.Id, personId)
    if err != nil {
        return err
    }
    
    // 5. Notificar al organizador
    err = notifyInvitationAccepted(ctx, invitation)
    if err != nil {
        log.Warn("Failed to notify invitation accepted", "error", err)
    }
    
    return nil
}
```

### 📊 Consultas de Análisis

#### **Reportes de Organización**
```sql
-- Resumen ejecutivo de organización
SELECT 
    o.display_name,
    o.status,
    COUNT(DISTINCT ob.id) as total_branches,
    COUNT(DISTINCT oe.id) as total_employees,
    COUNT(DISTINCT oe.id) FILTER (WHERE oe.status = 'active') as active_employees,
    COUNT(DISTINCT oi.id) as pending_invitations,
    COUNT(DISTINCT oint.id) as active_integrations,
    COUNT(DISTINCT od.id) as custom_domains
FROM organization o
LEFT JOIN organization_branch ob ON o.id = ob.organization_id AND ob.deleted_at IS NULL
LEFT JOIN organization_employee oe ON o.id = oe.organization_id AND oe.deleted_at IS NULL
LEFT JOIN organization_invitation oi ON o.id = oi.organization_id AND oi.status = 'pending'
LEFT JOIN organization_integration oint ON o.id = oint.organization_id AND oint.deleted_at IS NULL
LEFT JOIN organization_custom_domain od ON o.id = od.organization_id AND od.deleted_at IS NULL
WHERE o.deleted_at IS NULL
GROUP BY o.id, o.display_name, o.status;
```

#### **Análisis de Performance de Integraciones**
```sql
-- Top integraciones por uso
SELECT 
    it.display_name,
    COUNT(oi.id) as organizations_using,
    SUM(oi.api_usage_count) as total_api_calls,
    AVG(oi.api_usage_count) as avg_api_calls_per_org
FROM integration_type it
JOIN organization_integration oi ON it.id = oi.integration_type_id
WHERE oi.deleted_at IS NULL
AND oi.status = 'active'
GROUP BY it.id, it.display_name
ORDER BY total_api_calls DESC;
```

#### **Métricas de Adopción**
```sql
-- Adopción de características por organización
WITH org_features AS (
    SELECT 
        o.id,
        o.display_name,
        COUNT(DISTINCT ob.id) > 1 as has_multiple_branches,
        COUNT(DISTINCT oe.id) > 5 as has_team,
        COUNT(DISTINCT oint.id) > 0 as uses_integrations,
        COUNT(DISTINCT od.id) > 0 as has_custom_domain,
        o.created_at
    FROM organization o
    LEFT JOIN organization_branch ob ON o.id = ob.organization_id AND ob.deleted_at IS NULL
    LEFT JOIN organization_employee oe ON o.id = oe.organization_id AND oe.deleted_at IS NULL
    LEFT JOIN organization_integration oint ON o.id = oint.organization_id AND oint.deleted_at IS NULL
    LEFT JOIN organization_custom_domain od ON o.id = od.organization_id AND od.deleted_at IS NULL
    WHERE o.deleted_at IS NULL
    GROUP BY o.id, o.display_name, o.created_at
)
SELECT 
    COUNT(*) as total_organizations,
    SUM(CASE WHEN has_multiple_branches THEN 1 ELSE 0 END) as with_multiple_branches,
    SUM(CASE WHEN has_team THEN 1 ELSE 0 END) as with_team,
    SUM(CASE WHEN uses_integrations THEN 1 ELSE 0 END) as using_integrations,
    SUM(CASE WHEN has_custom_domain THEN 1 ELSE 0 END) as with_custom_domain,
    ROUND(AVG(EXTRACT(days FROM NOW() - created_at))) as avg_age_days
FROM org_features;
```
        uuid created_by FK
        timestamptz updated_at
        uuid updated_by FK
        timestamptz deleted_at
    }
    
    employee_roles {
        uuid id PK
        uuid employee_id FK
        uuid role_id FK
        boolean is_primary
        date assigned_date
        date revoked_date
        timestamptz created_at
        uuid created_by FK
    }
    
    organization_invite {
        uuid id PK
        uuid organization_id FK
        uuid inviter_person_id FK
        varchar invitee_email
        uuid invitee_person_id FK
        uuid role_id FK
        uuid branch_id FK
        varchar token
        invitation_status_enum status
        timestamptz expires_at
        timestamptz accepted_at
        timestamptz rejected_at
        timestamptz cancelled_at
        jsonb metadata
        timestamptz created_at
        uuid created_by FK
        timestamptz updated_at
        uuid updated_by FK
    }
    
    integration_type {
        uuid id PK
        varchar name
        varchar display_name
        text description
        integration_category_enum category
        varchar provider
        varchar version
        jsonb configuration_schema
        boolean webhook_support
        boolean oauth_support
        boolean api_key_support
        timestamptz created_at
        timestamptz updated_at
        uuid updated_by FK
    }
    
    organization_integration {
        uuid id PK
        uuid organization_id FK
        uuid integration_type_id FK
        varchar name
        text description
        integration_status_enum status
        jsonb configuration
        jsonb credentials
        timestamptz last_sync_at
        sync_status_enum last_sync_status
        text last_sync_error
        sync_frequency_enum sync_frequency
        boolean auto_sync_enabled
        varchar webhook_url
        varchar webhook_secret
        jsonb oauth_token
        integer api_usage_count
        integer api_rate_limit
        boolean is_active
        timestamptz created_at
        uuid created_by FK
        timestamptz updated_at
        uuid updated_by FK
        timestamptz deleted_at
    }
    
    organization_domain {
        uuid id PK
        uuid organization_id FK
        varchar domain_name
        domain_type_enum domain_type
        domain_status_enum status
        boolean is_primary
        boolean redirect_to_primary
        varchar dns_verification_token
        boolean dns_verified
        timestamptz verified_at
        timestamptz expires_at
        jsonb ssl_config
        jsonb metadata
        timestamptz created_at
        uuid created_by FK
        timestamptz updated_at
        uuid updated_by FK
        timestamptz deleted_at
    }
```

---

## 📚 Migraciones Detalladas

### 0001_create_enums_and_types.up.sql

**Propósito**: Establece los tipos de datos fundamentales y funciones utilitarias del sistema.

#### ENUMs Creados

```sql
-- Estados de organización
CREATE TYPE organization_status_enum AS ENUM (
    'active',        -- Organización activa y operativa
    'suspended',     -- Temporalmente suspendida
    'deleted',       -- Eliminada (soft delete)
    'pending_activation' -- Esperando activación inicial
);

-- Tipos de propietario
CREATE TYPE owner_type_enum AS ENUM (
    'individual',    -- Persona física
    'company'        -- Persona jurídica
);

-- Estados de empleado
CREATE TYPE employee_status_enum AS ENUM (
    'active',        -- Empleado activo
    'inactive',      -- Inactivo (licencia, vacaciones extendidas)
    'suspended',     -- Suspendido disciplinariamente
    'terminated'     -- Relación laboral terminada
);

-- Estados de invitación
CREATE TYPE invitation_status_enum AS ENUM (
    'pending',       -- Esperando respuesta
    'accepted',      -- Aceptada y procesada
    'rejected',      -- Rechazada por el invitado
    'expired',       -- Expiró por tiempo
    'cancelled'      -- Cancelada por el organizador
);
```

#### Funciones Utilitarias Principales

```sql
-- Generación de UUIDs seguros
CREATE OR REPLACE FUNCTION generate_uuid()
RETURNS UUID AS $$
BEGIN
    RETURN gen_random_uuid();
END;
$$ LANGUAGE plpgsql IMMUTABLE SET search_path = public;
```

**¿Por qué?**: Centraliza la generación de UUIDs y permite cambiar la implementación sin tocar todas las tablas.

```sql
-- Timestamps UTC consistentes
CREATE OR REPLACE FUNCTION current_timestamp_utc()
RETURNS TIMESTAMP WITH TIME ZONE AS $$
BEGIN
    RETURN CURRENT_TIMESTAMP AT TIME ZONE 'UTC';
END;
$$ LANGUAGE plpgsql STABLE SET search_path = public;
```

**¿Por qué?**: Garantiza que todos los timestamps se almacenen en UTC, evitando problemas de zona horaria.

### 0002_create_organization_core.up.sql

**Propósito**: Crea las entidades centrales del sistema organizacional.

#### Tabla `organization`

La tabla principal que representa una empresa inmobiliaria.

```sql
CREATE TABLE organization (
    id                  UUID         PRIMARY KEY DEFAULT generate_uuid(),
    display_name        VARCHAR(120) NOT NULL,
    logo_url            TEXT,
    fiscal_address_id   UUID,        -- FK lógica → address-svc
    matricula           VARCHAR(32),
    status              organization_status_enum NOT NULL DEFAULT 'active',
    
    -- Auditoría completa
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp_utc(),
    created_by          UUID NOT NULL,  -- FK → auth-identity-svc
    updated_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp_utc(),
    updated_by          UUID,
    deleted_at          TIMESTAMP WITH TIME ZONE -- Soft delete
);
```

**Características clave**:
- **Soft Delete**: Nunca se eliminan registros físicamente
- **Auditoría completa**: Quién y cuándo creó/modificó
- **FKs lógicas**: Referencias a otros microservicios sin constraints DB

#### Constraints de Negocio

```sql
-- Consistency en soft delete
CONSTRAINT chk_organization_soft_delete_status CHECK (
    (deleted_at IS NULL AND status != 'deleted') OR 
    (deleted_at IS NOT NULL AND status = 'deleted')
)
```

**¿Por qué?**: Evita estados inconsistentes donde una organización esté marcada como eliminada pero no tenga `deleted_at`.

### 0003_create_structure.up.sql

**Propósito**: Define la estructura jerárquica interna de las organizaciones.

#### Tabla `organization_branch` (Sucursales)

```sql
CREATE TABLE organization_branch (
    id              UUID         PRIMARY KEY DEFAULT generate_uuid(),
    organization_id UUID         NOT NULL,
    name            VARCHAR(200) NOT NULL,
    description     TEXT,
    address_id      UUID,        -- FK lógica → address-svc
    phone           VARCHAR(20),
    email           VARCHAR(255),
    is_main         BOOLEAN      NOT NULL DEFAULT false,
    is_active       BOOLEAN      NOT NULL DEFAULT true
);
```

**Reglas de Negocio**:
- Solo puede haber **una sucursal principal** por organización
- Cada organización debe tener **al menos una sucursal**
- Las sucursales inactivas no pueden recibir nuevos empleados

#### Tabla `organization_role` (Roles)

```sql
CREATE TABLE organization_role (
    id              UUID         PRIMARY KEY DEFAULT generate_uuid(),
    organization_id UUID         NOT NULL,
    name            VARCHAR(100) NOT NULL,
    description     TEXT,
    permissions     JSONB        NOT NULL DEFAULT '{}',
    is_default      BOOLEAN      NOT NULL DEFAULT false,
    is_admin        BOOLEAN      NOT NULL DEFAULT false,
    is_active       BOOLEAN      NOT NULL DEFAULT true
);
```

**Sistema de Permisos**: Los permisos se almacenan como JSONB:

```json
{
  "properties": {
    "read": true,
    "write": true,
    "delete": false
  },
  "reports": {
    "financial": true,
    "operational": false
  },
  "admin": {
    "user_management": true,
    "system_config": false
  }
}
```

### 0004_create_employee.up.sql

**Propósito**: Gestiona empleados y sus asignaciones de roles.

#### Tabla `employees`

```sql
CREATE TABLE employees (
    id                  UUID         PRIMARY KEY DEFAULT generate_uuid(),
    organization_id     UUID         NOT NULL,
    person_id           UUID         NOT NULL,  -- FK lógica → person-svc
    branch_id           UUID         NOT NULL,  -- FK → organization_branch
    primary_role_id     UUID,                   -- FK → organization_role
    employee_code       VARCHAR(50),
    hire_date           DATE         NOT NULL DEFAULT CURRENT_DATE,
    termination_date    DATE,
    status              employee_status_enum NOT NULL DEFAULT 'active',
    metadata            JSONB        DEFAULT '{}'
);
```

#### Sistema de Roles Múltiples

Los empleados pueden tener múltiples roles, pero uno debe ser primario:

```sql
CREATE TABLE employee_roles (
    id              UUID         PRIMARY KEY DEFAULT generate_uuid(),
    employee_id     UUID         NOT NULL,
    role_id         UUID         NOT NULL,
    is_primary      BOOLEAN      NOT NULL DEFAULT false,
    assigned_date   DATE         NOT NULL DEFAULT CURRENT_DATE,
    revoked_date    DATE
);
```

#### Función de Sincronización

```sql
CREATE OR REPLACE FUNCTION sync_employee_primary_role()
RETURNS TRIGGER AS $$
BEGIN
    -- Prevenir recursión infinita
    IF pg_trigger_depth() > 1 THEN
        RETURN COALESCE(NEW, OLD);
    END IF;
    
    -- Si se marca un rol como primario
    IF NEW.is_primary = true THEN
        -- Actualizar la tabla employees
        UPDATE employees 
        SET primary_role_id = NEW.role_id,
            updated_at = current_timestamp_utc()
        WHERE id = NEW.employee_id;
        
        -- Desmarcar otros roles como primarios
        UPDATE employee_roles 
        SET is_primary = false 
        WHERE employee_id = NEW.employee_id 
        AND id != NEW.id;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
```

**¿Por qué esta complejidad?**: Mantener la consistencia entre `employees.primary_role_id` y `employee_roles.is_primary` automáticamente.

### 0005_create_organization_invite.up.sql

**Propósito**: Sistema completo de invitaciones con tokens seguros y límites dinámicos.

#### Tabla `organization_invite`

```sql
CREATE TABLE organization_invite (
    id                  UUID         PRIMARY KEY DEFAULT generate_uuid(),
    organization_id     UUID         NOT NULL,
    inviter_person_id   UUID         NOT NULL,  -- Quien invita
    invitee_email       VARCHAR(255) NOT NULL,  -- Email del invitado
    invitee_person_id   UUID,                   -- Se llena al aceptar
    role_id             UUID         NOT NULL,  -- Rol asignado
    branch_id           UUID,                   -- Sucursal opcional
    token               VARCHAR(500) NOT NULL UNIQUE,
    status              invitation_status_enum NOT NULL DEFAULT 'pending',
    expires_at          TIMESTAMP WITH TIME ZONE NOT NULL,
    accepted_at         TIMESTAMP WITH TIME ZONE,
    rejected_at         TIMESTAMP WITH TIME ZONE,
    cancelled_at        TIMESTAMP WITH TIME ZONE,
    metadata            JSONB        DEFAULT '{}'
);
```

#### Generación de Tokens Seguros

```sql
CREATE OR REPLACE FUNCTION generate_invite_token(token_length INTEGER DEFAULT 64)
RETURNS VARCHAR AS $$
DECLARE
    chars TEXT := 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_';
    result TEXT := '';
    i INTEGER;
BEGIN
    FOR i IN 1..token_length LOOP
        result := result || substr(chars, floor(random() * length(chars) + 1)::INTEGER, 1);
    END LOOP;
    
    RETURN result;
END;
$$ LANGUAGE plpgsql;
```

#### Control de Límites Dinámicos

```sql
CREATE OR REPLACE FUNCTION check_organization_invite_limits()
RETURNS TRIGGER AS $$
DECLARE
    pending_count INTEGER;
    org_limit INTEGER;
BEGIN
    -- Obtener límite desde configuración
    SELECT COALESCE(
        (setting_value->>'value')::INTEGER,
        50  -- Default
    ) INTO org_limit
    FROM organization_settings
    WHERE organization_id = NEW.organization_id
    AND setting_key = 'max_pending_invites';
    
    -- Contar invitaciones pendientes
    SELECT COUNT(*) INTO pending_count
    FROM organization_invite 
    WHERE organization_id = NEW.organization_id 
    AND status = 'pending' 
    AND expires_at > current_timestamp_utc();
    
    -- Verificar límite
    IF pending_count >= org_limit THEN
        RAISE EXCEPTION 
            USING ERRCODE = '23514',
                  MESSAGE = format('Organization has reached max pending invitations: %s', org_limit);
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
```

### 0006_create_organization_integration.up.sql

**Propósito**: Hub central de integraciones con sistemas externos.

#### Catálogo de Integraciones

```sql
CREATE TABLE integration_type (
    id                      UUID         PRIMARY KEY DEFAULT generate_uuid(),
    name                    VARCHAR(100) NOT NULL UNIQUE,
    display_name           VARCHAR(200) NOT NULL,
    description            TEXT,
    category               integration_category_enum NOT NULL,
    provider               VARCHAR(100) NOT NULL,
    version                VARCHAR(20)  NOT NULL DEFAULT '1.0',
    configuration_schema   JSONB,
    webhook_support        BOOLEAN      NOT NULL DEFAULT false,
    oauth_support          BOOLEAN      NOT NULL DEFAULT false,
    api_key_support        BOOLEAN      NOT NULL DEFAULT false
);
```

**Tipos de Integración Preconfigurados**:

```sql
INSERT INTO integration_type (name, display_name, category, provider, configuration_schema) VALUES
('quickbooks_online', 'QuickBooks Online', 'accounting', 'Intuit', 
 '{"required": ["client_id", "client_secret"], "optional": ["sandbox_mode"]}'::jsonb),
('hubspot_crm', 'HubSpot CRM', 'crm', 'HubSpot',
 '{"required": ["api_key"], "optional": ["portal_id"]}'::jsonb),
('mailchimp', 'Mailchimp', 'marketing', 'Mailchimp',
 '{"required": ["api_key"], "optional": ["datacenter", "list_id"]}'::jsonb);
```

#### Configuración por Organización

```sql
CREATE TABLE organization_integration (
    id                  UUID         PRIMARY KEY DEFAULT generate_uuid(),
    organization_id     UUID         NOT NULL,
    integration_type_id UUID         NOT NULL,
    name                VARCHAR(200) NOT NULL,
    status              integration_status_enum NOT NULL DEFAULT 'inactive',
    configuration       JSONB        NOT NULL DEFAULT '{}',
    credentials         JSONB        DEFAULT '{}',    -- Encriptado en app
    last_sync_at        TIMESTAMP WITH TIME ZONE,
    last_sync_status    sync_status_enum,
    sync_frequency      sync_frequency_enum,
    auto_sync_enabled   BOOLEAN      NOT NULL DEFAULT false,
    api_usage_count     INTEGER      NOT NULL DEFAULT 0
);
```

#### Sistema de Eventos para Workers

```sql
CREATE TABLE organization_integration_event (
    id              UUID         PRIMARY KEY DEFAULT generate_uuid(),
    integration_id  UUID         NOT NULL,
    event_type      integration_event_type_enum NOT NULL,
    event_data      JSONB        NOT NULL DEFAULT '{}',
    status          integration_event_status_enum NOT NULL DEFAULT 'pending',
    retry_count     INTEGER      NOT NULL DEFAULT 0,
    max_retries     INTEGER      NOT NULL DEFAULT 3,
    scheduled_at    TIMESTAMP WITH TIME ZONE,
    processed_at    TIMESTAMP WITH TIME ZONE
);
```

**Flujo de Trabajo**:
1. Se crea evento `sync_contacts` con status `pending`
2. Worker lo toma y marca como `processing`
3. Si falla, incrementa `retry_count` y reintenta
4. Al éxito marca como `completed`

### 0007_create_organization_domain.up.sql

**Propósito**: Gestión de dominios personalizados con verificación DNS automática.

#### Tabla Principal de Dominios

```sql
CREATE TABLE organization_domain (
    id                      UUID         PRIMARY KEY DEFAULT generate_uuid(),
    organization_id         UUID         NOT NULL,
    domain_name             VARCHAR(255) NOT NULL,
    domain_type             domain_type_enum NOT NULL DEFAULT 'subdomain',
    status                  domain_status_enum NOT NULL DEFAULT 'pending',
    is_primary              BOOLEAN      NOT NULL DEFAULT false,
    redirect_to_primary     BOOLEAN      NOT NULL DEFAULT false,
    dns_verification_token  VARCHAR(255) NOT NULL DEFAULT generate_dns_verification_token(),
    dns_verified            BOOLEAN      NOT NULL DEFAULT false,
    verified_at             TIMESTAMP WITH TIME ZONE,
    expires_at              TIMESTAMP WITH TIME ZONE,
    ssl_config              JSONB        DEFAULT '{}',
    metadata                JSONB        DEFAULT '{}'
);
```

#### Constraint Inteligente

```sql
-- Un dominio no puede redirigir a sí mismo si es primario
CONSTRAINT chk_organization_domain_redirect_logic 
    CHECK (NOT (is_primary = true AND redirect_to_primary = true))
```

#### Verificación DNS Automática

```sql
CREATE OR REPLACE FUNCTION sync_domain_dns_verification()
RETURNS TRIGGER AS $$
DECLARE
    required_records INTEGER;
    verified_records INTEGER;
BEGIN
    -- Contar registros DNS requeridos y verificados
    SELECT 
        COUNT(*) FILTER (WHERE is_required = true),
        COUNT(*) FILTER (WHERE is_required = true AND is_verified = true)
    INTO required_records, verified_records
    FROM organization_domain_dns 
    WHERE domain_id = NEW.domain_id;
    
    -- Actualizar estado de verificación
    UPDATE organization_domain 
    SET 
        dns_verified = (required_records > 0 AND verified_records = required_records),
        verified_at = CASE 
            WHEN verified_records = required_records AND verified_records > 0 
            THEN current_timestamp_utc() 
            ELSE NULL 
        END
    WHERE id = NEW.domain_id;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
```

#### Registros DNS Automáticos

```sql
CREATE OR REPLACE FUNCTION create_default_dns_records()
RETURNS TRIGGER AS $$
BEGIN
    -- Crear registros DNS por defecto al agregar un dominio
    INSERT INTO organization_domain_dns (domain_id, record_type, name, value, is_required, created_by) VALUES
    (NEW.id, 'TXT', '_rem-verify', NEW.dns_verification_token, true, NEW.created_by),
    (NEW.id, 'CNAME', 'www', NEW.domain_name, false, NEW.created_by),
    (NEW.id, 'A', '@', '0.0.0.0', true, NEW.created_by);
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
```

### 0008_apply_improvements.up.sql

**Propósito**: Optimizaciones finales, particionado automático y funciones de mantenimiento.

#### Sistema de Particionado Automático

```sql
CREATE OR REPLACE FUNCTION trg_monthly_partition()
RETURNS TRIGGER AS $$
BEGIN
    -- ADVERTENCIA: DDL en trigger puede causar dead-locks en alta concurrencia
    PERFORM create_monthly_partition(TG_TABLE_NAME, NEW.created_at::date);
    RETURN NEW;
EXCEPTION
    WHEN duplicate_table THEN
        RETURN NEW;
    WHEN OTHERS THEN
        RAISE WARNING 'Failed to create partition for %: %', TG_TABLE_NAME, SQLERRM;
        RETURN NEW;
END;
$$ LANGUAGE plpgsql;
```

**¿Por qué particiones?**: Las tablas de logs pueden crecer rápidamente. El particionado mensual:
- Mejora el rendimiento de consultas por rango de fechas
- Facilita el archivado/eliminación de datos antiguos
- Permite mantenimiento en caliente

#### Vista de Suscripciones

```sql
CREATE OR REPLACE VIEW organization_subscription_details AS
SELECT 
    o.id as organization_id,
    o.display_name,
    o.status as organization_status,
    o.created_at as organization_created_at,
    -- Campos básicos hasta que billing-svc defina schema completo
    NULL::TEXT as subscription_status,
    NULL::DATE as current_period_end,
    NULL::JSONB as plan_features
FROM organization o
WHERE o.deleted_at IS NULL;
```

---

## 🔧 Funciones Utilitarias

### Funciones de Validación

#### `is_valid_email(email TEXT)`
```sql
CREATE OR REPLACE FUNCTION is_valid_email(email TEXT)
RETURNS BOOLEAN AS $$
BEGIN
    RETURN email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$';
END;
$$ LANGUAGE plpgsql IMMUTABLE;
```

#### `is_valid_domain(domain_name TEXT)`
```sql
CREATE OR REPLACE FUNCTION is_valid_domain(domain_name TEXT)
RETURNS BOOLEAN AS $$
BEGIN
    -- Verificar que no sea NULL o vacío
    IF domain_name IS NULL OR domain_name = '' THEN
        RETURN FALSE;
    END IF;
    
    -- Verificar longitud (max 253 caracteres para FQDN)
    IF char_length(domain_name) > 253 THEN
        RETURN FALSE;
    END IF;
    
    -- Regex con anclas estrictas para prevenir bypass CRLF
    RETURN domain_name ~ E'^\\A([a-z0-9]([a-z0-9\\-]{0,61}[a-z0-9])?\\\.)+[a-z0-9]{2,}\\z$';
END;
$$ LANGUAGE plpgsql IMMUTABLE;
```

### Funciones de Mantenimiento

#### `cleanup_old_partitions(table_name TEXT, months_to_keep INTEGER)`
```sql
CREATE OR REPLACE FUNCTION cleanup_old_partitions(
    table_name TEXT,
    months_to_keep INTEGER DEFAULT 12
)
RETURNS INTEGER AS $$
DECLARE
    cutoff_date DATE;
    partition_name TEXT;
    dropped_count INTEGER := 0;
BEGIN
    cutoff_date := CURRENT_DATE - INTERVAL '1 month' * months_to_keep;
    
    FOR partition_name IN
        SELECT schemaname||'.'||tablename
        FROM pg_tables
        WHERE tablename LIKE table_name || '_y%m%'
        AND tablename < table_name || '_' || to_char(cutoff_date, 'YYYY"m"MM')
    LOOP
        EXECUTE 'DROP TABLE IF EXISTS ' || partition_name;
        dropped_count := dropped_count + 1;
        RAISE NOTICE 'Dropped old partition: %', partition_name;
    END LOOP;
    
    RETURN dropped_count;
END;
$$ LANGUAGE plpgsql;
```

---

## ⚡ Triggers y Automatizaciones

### Triggers de Auditoría

Todas las tablas principales tienen triggers que actualizan automáticamente `updated_at`:

```sql
-- Función genérica reutilizable
CREATE OR REPLACE FUNCTION update_timestamp_utc()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = current_timestamp_utc();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Aplicado a todas las tablas
CREATE TRIGGER trg_organization_updated_at
    BEFORE UPDATE ON organization
    FOR EACH ROW
    EXECUTE FUNCTION update_timestamp_utc();
```

### Triggers de Business Logic

#### Prevención de Hard Deletes
```sql
CREATE OR REPLACE FUNCTION prevent_hard_delete_organization()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'DELETE' THEN
        RAISE EXCEPTION 'Direct DELETE not allowed. Use soft delete by setting deleted_at timestamp.';
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;
```

#### Validación de Último Admin
```sql
CREATE OR REPLACE FUNCTION validate_last_admin_employee()
RETURNS TRIGGER AS $$
DECLARE
    admin_count INTEGER;
BEGIN
    -- Contar administradores activos restantes
    SELECT COUNT(*) INTO admin_count
    FROM employees e
    JOIN employee_roles er ON e.id = er.employee_id
    JOIN organization_role r ON er.role_id = r.id
    WHERE e.organization_id = OLD.organization_id
    AND e.status = 'active'
    AND e.deleted_at IS NULL
    AND r.is_admin = true
    AND er.revoked_date IS NULL
    AND e.id != OLD.id;  -- Excluir el que se está eliminando/desactivando
    
    IF admin_count = 0 THEN
        RAISE EXCEPTION 'Cannot remove the last active administrator from organization';
    END IF;
    
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;
```

---

## 📊 Estrategias de Particionado

### Particionado Mensual por Rango

Las tablas de logs se particionan mensualmente para optimizar rendimiento:

```sql
-- Tabla padre
CREATE TABLE organization_invite_log (
    id              UUID         NOT NULL DEFAULT generate_uuid(),
    invite_id       UUID         NOT NULL,
    action          invitation_log_action_enum NOT NULL,
    old_status      invitation_status_enum,
    new_status      invitation_status_enum,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp_utc()
) PARTITION BY RANGE (created_at);

-- Particiones automáticas
CREATE TABLE organization_invite_log_y2024m01 
PARTITION OF organization_invite_log
FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');

CREATE TABLE organization_invite_log_y2024m02 
PARTITION OF organization_invite_log
FOR VALUES FROM ('2024-02-01') TO ('2024-03-01');
```

### Creación Automática de Particiones

```sql
CREATE OR REPLACE FUNCTION create_monthly_partition(
    table_name TEXT,
    partition_date DATE DEFAULT CURRENT_DATE
)
RETURNS VOID AS $$
DECLARE
    partition_name TEXT;
    start_date DATE;
    end_date DATE;
BEGIN
    start_date := date_trunc('month', partition_date);
    end_date := start_date + INTERVAL '1 month';
    partition_name := table_name || '_y' || to_char(start_date, 'YYYY') || 'm' || to_char(start_date, 'MM');
    
    EXECUTE format($f$
        CREATE TABLE IF NOT EXISTS %I 
        PARTITION OF %I
        FOR VALUES FROM (%L) TO (%L)
    $f$, partition_name, table_name, start_date, end_date);
    
    RAISE NOTICE 'Created partition % for table %', partition_name, table_name;
END;
$$ LANGUAGE plpgsql;
```

---

## 🔒 Seguridad y Validaciones

### Soft Delete Consistency

Todos los soft deletes deben mantener consistencia entre `deleted_at` y `status`:

```sql
CONSTRAINT chk_organization_soft_delete_status CHECK (
    (deleted_at IS NULL AND status != 'deleted') OR 
    (deleted_at IS NOT NULL AND status = 'deleted')
)
```

### Validaciones de Integridad Referencial

Aunque usamos FKs lógicas a otros microservicios, mantenemos validaciones básicas:

```sql
-- UUIDs válidos
CONSTRAINT chk_organization_valid_created_by 
    CHECK (is_valid_uuid(created_by::TEXT))

-- Emails válidos
CONSTRAINT chk_organization_invite_email_format
    CHECK (is_valid_email(invitee_email))

-- Dominios válidos
CONSTRAINT chk_organization_domain_name_format
    CHECK (is_valid_domain(domain_name))
```

### Encriptación de Credenciales

Las credenciales de integración se almacenan como JSONB pero deben encriptarse a nivel de aplicación:

```json
{
  "api_key": "encrypted:AES256:base64encodedcontent",
  "secret": "encrypted:AES256:anothersecret",
  "oauth_token": {
    "access_token": "encrypted:AES256:token",
    "refresh_token": "encrypted:AES256:refresh"
  }
}
```

---

## 💡 Ejemplos de Uso

### Crear Nueva Organización Completa

```sql
-- 1. Crear organización
INSERT INTO organization (display_name, created_by) 
VALUES ('Inmobiliaria Los Pinos', 'user-uuid-123')
RETURNING id;

-- 2. El trigger create_organization_defaults automáticamente:
--    - Crea configuración por defecto
--    - Crea sucursal principal
--    - Crea roles básicos (admin, agent, viewer)

-- 3. Agregar primer empleado/admin
INSERT INTO employees (organization_id, person_id, branch_id, created_by)
VALUES ('org-uuid', 'person-uuid', 'branch-uuid', 'user-uuid-123');

-- 4. Asignar rol de administrador
INSERT INTO employee_roles (employee_id, role_id, is_primary, created_by)
SELECT 'employee-uuid', r.id, true, 'user-uuid-123'
FROM organization_role r 
WHERE r.organization_id = 'org-uuid' AND r.is_admin = true;
```

### Configurar Integración con CRM

```sql
-- 1. Crear integración con HubSpot
INSERT INTO organization_integration (
    organization_id, 
    integration_type_id,
    name,
    configuration,
    credentials,
    created_by
) VALUES (
    'org-uuid',
    (SELECT id FROM integration_type WHERE name = 'hubspot_crm'),
    'HubSpot CRM Principal',
    '{"portal_id": "12345", "sync_contacts": true, "sync_deals": true}'::jsonb,
    '{"api_key": "encrypted:AES256:hubspotapikey"}'::jsonb,
    'user-uuid-123'
);

-- 2. Activar sincronización automática
UPDATE organization_integration 
SET 
    status = 'active',
    auto_sync_enabled = true,
    sync_frequency = 'hourly'
WHERE id = 'integration-uuid';

-- 3. Crear evento de sincronización inicial
INSERT INTO organization_integration_event (
    integration_id,
    event_type,
    event_data
) VALUES (
    'integration-uuid',
    'full_sync',
    '{"entities": ["contacts", "deals"], "since": null}'::jsonb
);
```

### Invitar Nuevo Usuario

```sql
-- 1. Crear invitación
INSERT INTO organization_invite (
    organization_id,
    inviter_person_id,
    invitee_email,
    role_id,
    branch_id,
    token,
    expires_at,
    created_by
) VALUES (
    'org-uuid',
    'inviter-person-uuid',
    'nuevo@ejemplo.com',
    'role-uuid',
    'branch-uuid',
    generate_invite_token(),
    current_timestamp_utc() + INTERVAL '7 days',
    'user-uuid-123'
);

-- 2. El sistema automáticamente:
--    - Verifica límites de invitaciones pendientes
--    - Genera token seguro único
--    - Programa email de invitación (via evento)
--    - Logs la acción

-- 3. Al aceptar la invitación:
UPDATE organization_invite 
SET 
    status = 'accepted',
    accepted_at = current_timestamp_utc(),
    invitee_person_id = 'new-person-uuid'
WHERE token = 'invitation-token';

-- 4. Crear empleado automáticamente
INSERT INTO employees (organization_id, person_id, branch_id, created_by)
SELECT organization_id, invitee_person_id, branch_id, invitee_person_id
FROM organization_invite 
WHERE id = 'invite-uuid';
```

### Configurar Dominio Personalizado

```sql
-- 1. Agregar dominio
INSERT INTO organization_domain (
    organization_id,
    domain_name,
    domain_type,
    created_by
) VALUES (
    'org-uuid',
    'lospinos.com',
    'custom',
    'user-uuid-123'
);

-- 2. El trigger create_default_dns_records automáticamente crea:
--    TXT _rem-verify.lospinos.com = "verification-token"
--    CNAME www.lospinos.com = lospinos.com
--    A lospinos.com = 0.0.0.0 (placeholder)

-- 3. Cliente configura DNS en su proveedor
-- 4. Verificar DNS
SELECT sync_domain_dns_verification() FROM organization_domain_dns 
WHERE domain_id = 'domain-uuid';

-- 5. Marcar como primario si verificado
UPDATE organization_domain 
SET is_primary = true
WHERE id = 'domain-uuid' AND dns_verified = true;
```

---

## 🎯 Mejores Prácticas

### Desarrollo

1. **Siempre usar transacciones** para operaciones multi-tabla
2. **Validar FKs lógicas** en la capa de aplicación antes de insertar
3. **Encriptar credenciales** antes de almacenar en JSONB
4. **Usar UUIDs** para todas las PKs para facilitar replicación/sharding futuro

### Rendimiento

1. **Índices apropiados**: Todos los WHERE comunes están indexados
2. **Particionado**: Logs se particionan mensualmente
3. **Soft deletes**: Usar `WHERE deleted_at IS NULL` en índices parciales
4. **JSONB**: Usar `jsonb_typeof()` en constraints para validar estructura

### Operaciones

1. **Backups**: Especial atención a credenciales encriptadas
2. **Monitoreo**: Triggers de particionado pueden indicar alta carga
3. **Mantenimiento**: Ejecutar `cleanup_old_partitions()` mensualmente
4. **Rollbacks**: Los `.down.sql` están probados y documentados

### Seguridad

1. **Nunca hard delete**: Solo soft delete en producción
2. **Auditoría completa**: Todos los cambios tienen `created_by`/`updated_by`
3. **Validación de entrada**: Functions de validación para emails/dominios
4. **Principio de menor privilegio**: Roles de DB específicos por función

---

## 🔍 Troubleshooting

### Problemas Comunes

#### "Function does not exist" al aplicar migraciones

**Causa**: Orden incorrecto de migraciones o rollback parcial.

**Solución**:
```sql
-- Verificar funciones faltantes
SELECT 'update_timestamp_utc'::regproc;

-- Re-ejecutar migración 0002 si falta
\i 0002_create_organization_core.up.sql
```

#### Partición no encontrada en INSERT

**Causa**: Partition trigger falló o no se pre-crearon particiones.

**Solución**:
```sql
-- Crear partición manualmente
SELECT create_monthly_partition('organization_invite_log', CURRENT_DATE);

-- Verificar triggers
SELECT tgname FROM pg_trigger WHERE tgrelid = 'organization_invite_log'::regclass;
```

#### Constraint violation en soft delete

**Causa**: Trying to set `deleted_at` without updating `status`.

**Solución**:
```sql
-- ❌ Incorrecto
UPDATE organization SET deleted_at = NOW() WHERE id = 'org-uuid';

-- ✅ Correcto
UPDATE organization 
SET deleted_at = NOW(), status = 'deleted' 
WHERE id = 'org-uuid';
```

### 🔍 Queries de Diagnóstico

#### **Estado General del Sistema**
```sql
-- Resumen de organizaciones por estado
SELECT 
    status,
    COUNT(*) as count,
    COUNT(*) * 100.0 / SUM(COUNT(*)) OVER() as percentage
FROM organization 
WHERE deleted_at IS NULL
GROUP BY status;

-- Distribución de empleados por organización
SELECT 
    o.display_name,
    COUNT(oe.id) as total_employees,
    COUNT(CASE WHEN oe.status = 'active' THEN 1 END) as active_employees,
    COUNT(DISTINCT ob.id) as total_branches,
    COUNT(DISTINCT oint.id) as active_integrations
FROM organization o
LEFT JOIN organization_employee oe ON o.id = oe.organization_id AND oe.deleted_at IS NULL
LEFT JOIN organization_branch ob ON o.id = ob.organization_id AND ob.deleted_at IS NULL
LEFT JOIN organization_integration oint ON o.id = oint.organization_id 
    AND oint.deleted_at IS NULL AND oint.status = 'active'
WHERE o.deleted_at IS NULL
GROUP BY o.id, o.display_name
ORDER BY total_employees DESC;
```

#### **Performance y Utilización**
```sql
-- Tamaño de tablas principales
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size,
    pg_total_relation_size(schemaname||'.'||tablename) as size_bytes
FROM pg_tables 
WHERE tablename LIKE 'organization%'
ORDER BY size_bytes DESC;

-- Efectividad de índices
SELECT 
    schemaname,
    tablename,
    indexname,
    pg_size_pretty(pg_relation_size(indexname::regclass)) as index_size,
    idx_scan as times_used,
    idx_tup_read as tuples_read,
    idx_tup_fetch as tuples_fetched,
    ROUND(
        CASE WHEN idx_scan > 0 
        THEN idx_tup_fetch::numeric / idx_scan 
        ELSE 0 END, 2
    ) as avg_tuples_per_scan
FROM pg_stat_user_indexes 
WHERE schemaname = 'public' 
AND tablename LIKE 'organization%'
ORDER BY times_used DESC;

-- Estado de particiones de logs
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size,
    (SELECT COUNT(*) FROM information_schema.tables t2 
     WHERE t2.table_name LIKE tablename || '_y%m%') as partition_count
FROM pg_tables 
WHERE tablename LIKE '%_log' AND tablename LIKE 'organization%'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```
```sql
-- Organizaciones con estado inconsistente
SELECT id, display_name, status, deleted_at
FROM organization
WHERE (deleted_at IS NULL AND status = 'deleted') 
   OR (deleted_at IS NOT NULL AND status != 'deleted');

-- Empleados sin rol primario
SELECT e.id, e.employee_code, e.primary_role_id
FROM employees e
LEFT JOIN employee_roles er ON e.primary_role_id = er.role_id 
    AND er.employee_id = e.id 
    AND er.is_primary = true
WHERE e.deleted_at IS NULL 
AND er.id IS NULL;

-- Invitaciones expiradas sin limpiar
SELECT COUNT(*)
FROM organization_invite
WHERE status = 'pending' 
AND expires_at < current_timestamp_utc();

-- Organizaciones sin sucursal principal
SELECT o.id, o.display_name
FROM organization o
LEFT JOIN organization_branch ob ON o.id = ob.organization_id 
    AND ob.is_main = true AND ob.deleted_at IS NULL
WHERE o.deleted_at IS NULL AND ob.id IS NULL;
```

#### **Estadísticas de Uso**
```sql
-- Top organizaciones por integraciones
SELECT o.display_name, COUNT(oi.id) as integration_count
FROM organization o
LEFT JOIN organization_integration oi ON o.id = oi.organization_id
WHERE o.deleted_at IS NULL
GROUP BY o.id, o.display_name
ORDER BY integration_count DESC;

-- Uso de particiones
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size
FROM pg_tables 
WHERE tablename LIKE '%_log_y%'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

---

## 📈 Roadmap y Extensibilidad

### 🎯 Próximas Funcionalidades

1. **Sistema de Facturación Avanzado**: Expansión completa de `organization_subscription_details`
2. **Multi-tenancy Horizontal**: Sharding por organización para escalabilidad masiva
3. **API Gateway Integration**: Rate limiting inteligente por organización y usuario
4. **Advanced Analytics**: Métricas de uso por feature con ML predictions
5. **Workflow Automation**: Sistema de flujos de trabajo personalizables
6. **Document Management**: Integración con gestión documental y firma electrónica

### 🔧 Extensión de Integraciones

El sistema está diseñado para agregar nuevas integraciones de manera sencilla:

```sql
-- Nuevo tipo de integración
INSERT INTO integration_type (name, display_name, category, provider, configuration_schema) 
VALUES (
    'stripe_payments',
    'Stripe Payments', 
    'financial',
    'Stripe',
    '{"required": ["api_key", "webhook_secret"], "optional": ["test_mode", "currency"]}'::jsonb
);

-- Nueva categoría de integración
ALTER TYPE integration_category_enum ADD VALUE 'real_estate_specific';
ALTER TYPE integration_category_enum ADD VALUE 'legal_compliance';
```

### ⚡ Scaling Considerations

1. **Read Replicas**: Para queries analíticos pesados y reportes
2. **Connection Pooling**: PgBouncer configurado para máximo throughput
3. **Archival Strategy**: Mover particiones antiguas a storage frío automáticamente
4. **Event Sourcing**: Para auditoría completa de cambios críticos
5. **Caching Strategy**: Redis multi-layer para datos frecuentemente accedidos
6. **Microservice Decomposition**: Separación de dominios específicos cuando sea necesario

---

### 📋 Checklist para PRs

- [ ] Tests pasan en CI/CD
- [ ] Documentación actualizada  
- [ ] Scripts de migración probados (up/down)
- [ ] Performance no degradado
- [ ] Security review completado
- [ ] Backward compatibility mantenida

---


### 🔄 Compatibilidad

- **PostgreSQL**: 13+ (requerido para funcionalidades avanzadas)
- **Go**: 1.21+ (para el servicio)
- **API Versioning**: Mantiene compatibilidad hacia atrás desde v1.0

---
