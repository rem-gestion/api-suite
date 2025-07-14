# 🏢 Organization Service - Roadmap de Desarrollo 2025

## 📋 Estado Actual y Visión

El **Organization Service** es el microservicio central del sistema REM con una **base de datos completamente implementada** (8 migraciones evolutivas) y **documentación exhaustiva**. Este roadmap define el desarrollo de la capa de aplicación, APIs, integraciones y funcionalidades avanzadas.

### 🏆 **Logros Completados**
- ✅ **Sistema de Migraciones Completo**: 8 migraciones implementadas y documentadas
- ✅ **Arquitectura de Base de Datos**: 12 tablas, 62 funciones, 64 triggers, 97 índices
- ✅ **Documentación Técnica**: README principal y documentación detallada de migraciones
- ✅ **Patrones de Diseño**: Soft-delete, auditoría completa, particionado automático
- ✅ **Estrategias de Performance**: Índices condicionales, logs particionados

### 🎯 **Objetivos 2025**
- **🚀 Implementación Completa del Servicio**: APIs REST/gRPC, business logic
- **🔗 Integración con Microservicios**: auth-identity, address, person, notification
- **📊 Hub de Integraciones**: CRM, contabilidad, marketing, inmobiliarias
- **🌐 Dominios Personalizados**: Verificación DNS/SSL automática
- **⚡ Optimización y Escalabilidad**: Caching, performance, monitoreo
- **🛡️ Seguridad Avanzada**: Encriptación, validaciones, compliance

---

## 🗃️ Arquitectura y Entidades (Estado Implementado)

### 📊 **Sistema de Migraciones Completado**

| Migración | Estado | Componentes | Funcionalidad |
|-----------|--------|-------------|---------------|
| **0001** | ✅ | ENUMs base, funciones utilitarias | Fundamentos del sistema |
| **0002** | ✅ | organization, settings, owners | Núcleo organizacional |
| **0003** | ✅ | branches, roles, jerarquías | Estructura interna |
| **0004** | ✅ | employees, employee_roles | Gestión de personal |
| **0005** | ✅ | invitations, workflow automático | Sistema de invitaciones |
| **0006** | ✅ | integrations, logs particionados | Hub de integraciones |
| **0007** | ✅ | custom domains, DNS/SSL | Dominios personalizados |
| **0008** | ✅ | Optimizaciones, mantenimiento | Performance y escalabilidad |

### 🔧 **Características Técnicas Implementadas**
- **Soft Delete Universal**: Todas las tablas con protección automática
- **Auditoría Completa**: Timestamps y usuarios en todos los cambios
- **Foreign Keys Lógicas**: Independencia total entre microservicios
- **Particionado Automático**: Logs organizados por mes con limpieza automática
- **Índices Inteligentes**: Optimización condicional para registros activos
- **Validaciones Multi-Capa**: DB constraints + triggers + business logic

---

## 🚀 Roadmap de Desarrollo

### **FASE 1: Infraestructura del Servicio** ⏱️ (1-2 semanas)

> **🎯 Objetivo**: Implementar la estructura base del microservicio Go con todas las dependencias

#### 1.1 **Setup del Proyecto**
- [ ] **📁 Estructura de directorios**
```
organization-svc/
├── cmd/
│   ├── api/main.go              # Servidor HTTP + gRPC
│   ├── migrate/main.go          # Ejecutor de migraciones (ya existe)
│   └── worker/main.go           # Procesador de eventos asincrónicos
├── internal/
│   ├── api/
│   │   ├── rest/                # REST API handlers
│   │   │   ├── handlers/
│   │   │   │   ├── organization.go
│   │   │   │   ├── branch.go
│   │   │   │   ├── employee.go
│   │   │   │   ├── invite.go
│   │   │   │   ├── integration.go
│   │   │   │   └── domain.go
│   │   │   ├── middleware/
│   │   │   │   ├── auth.go
│   │   │   │   ├── rbac.go
│   │   │   │   ├── rate_limit.go
│   │   │   │   └── audit.go
│   │   │   └── router.go
│   │   └── grpc/                # gRPC server implementation
│   │       ├── server.go
│   │       └── organization_service.go
│   ├── domain/
│   │   ├── models/              # Domain models (GORM)
│   │   │   ├── organization.go
│   │   │   ├── branch.go
│   │   │   ├── employee.go
│   │   │   ├── invite.go
│   │   │   ├── integration.go
│   │   │   └── domain.go
│   │   ├── repositories/        # Repository interfaces
│   │   │   ├── organization.go
│   │   │   ├── employee.go
│   │   │   └── integration.go
│   │   └── services/            # Business logic
│   │       ├── organization.go
│   │       ├── employee.go
│   │       ├── invite.go
│   │       ├── integration.go
│   │       └── domain.go
│   ├── infrastructure/
│   │   ├── database/
│   │   │   ├── postgres.go
│   │   │   └── repositories/    # Repository implementations
│   │   ├── cache/
│   │   │   └── redis.go
│   │   ├── messaging/
│   │   │   └── rabbitmq.go
│   │   └── external/            # External service clients
│   │       ├── auth_client.go
│   │       ├── person_client.go
│   │       ├── address_client.go
│   │       └── notification_client.go
│   ├── config/
│   │   ├── config.go
│   │   └── dependencies.go      # DI container
│   └── utils/
│       ├── crypto.go            # Encriptación de credenciales
│       ├── tokens.go            # Generación de tokens seguros
│       └── validation.go        # Validaciones custom
├── pkg/                         # Packages públicos
│   ├── dto/                     # Data Transfer Objects
│   │   ├── organization.go
│   │   ├── employee.go
│   │   └── integration.go
│   └── events/                  # Event definitions
│       ├── organization_events.go
│       └── employee_events.go
├── deployments/
│   ├── docker/
│   │   ├── Dockerfile
│   │   └── docker-compose.yml
│   └── k8s/
│       ├── deployment.yaml
│       ├── service.yaml
│       └── configmap.yaml
└── tests/
    ├── integration/
    ├── unit/
    └── fixtures/
```

#### 1.2 **Configuración y Dependencies**
- [ ] **⚙️ Variables de entorno** (`.env`)
```bash
# Server
PORT=8083
GRPC_PORT=50053
ENVIRONMENT=development

# Database
DATABASE_URL=postgresql://user:pass@localhost:5432/organization_db
DATABASE_POOL_SIZE=20
DATABASE_TIMEOUT=30s

# Cache
REDIS_URL=redis://localhost:6379
REDIS_DB=0
CACHE_TTL=3600

# Message Queue
RABBITMQ_URL=amqp://user:pass@localhost:5672
QUEUE_PREFETCH_COUNT=10

# External Services
AUTH_SERVICE_URL=auth-identity-svc:8080
PERSON_SERVICE_URL=person-svc:8080
ADDRESS_SERVICE_URL=address-svc:8080
NOTIFICATION_SERVICE_URL=notification-svc:8080

# Security
JWT_SECRET=your-secret-key
ENCRYPTION_KEY=32-char-key-for-integrations

# Features
ENABLE_DOMAIN_VERIFICATION=true
ENABLE_INTEGRATION_HUB=true
MAX_INVITATIONS_PER_ORG=50

# Monitoring
JAEGER_ENDPOINT=http://jaeger:14268/api/traces
PROMETHEUS_PORT=9090
```

- [ ] **📦 Dependencies (go.mod)**
  - Echo v4 para REST API
  - gRPC para comunicación entre servicios
  - GORM v2 para ORM con PostgreSQL
  - Redis para caching
  - RabbitMQ para messaging
  - Prometheus para métricas
  - Jaeger para tracing

#### 1.3 **Docker & DevOps**
- [ ] **🐳 Dockerfile multi-stage**
- [ ] **📝 docker-compose.override.yml** para desarrollo
- [ ] **🔧 Air para hot reload** en desarrollo
- [ ] **📊 Health checks** y readiness probes

---

### **FASE 2: Implementación del Core del Servicio** ⏱️ (2-3 semanas)

> **🎯 Objetivo**: Implementar el CRUD completo para todas las entidades principales

#### 2.1 **Entidades Core (Semana 1)**

**📋 Prioridad 1: Organization Management**
- [ ] **Organization CRUD**
  - [ ] Modelo GORM con todas las relaciones
  - [ ] Repository con soft delete y filtros
  - [ ] Service con business logic y validaciones
  - [ ] REST handlers: POST, GET, PUT, DELETE
  - [ ] gRPC service implementation
  - [ ] Tests unitarios y de integración

- [ ] **Organization Settings**
  - [ ] CRUD de configuraciones key-value
  - [ ] Validación de schemas JSON para settings
  - [ ] Cache de configuraciones frecuentes
  - [ ] Versionado de cambios en settings

- [ ] **Organization Owners**
  - [ ] Gestión de propietarios múltiples
  - [ ] Validación de porcentajes (suma = 100%)
  - [ ] Integración con person-svc para detalles
  - [ ] Auditoría de cambios de propiedad

**📋 Prioridad 2: Estructura Organizacional**
- [ ] **Branches (Sucursales)**
  - [ ] CRUD con validación de sucursal principal
  - [ ] Integración con address-svc
  - [ ] Geolocalización y búsqueda por área
  - [ ] Protecciones contra eliminación indebida

- [ ] **Roles Internos**
  - [ ] Sistema de roles predefinidos y personalizados
  - [ ] Gestión de permisos como JSONB
  - [ ] Protección de roles críticos del sistema
  - [ ] Templates de roles por industria

#### 2.2 **Gestión de Empleados (Semana 2)**

- [ ] **Employee Management**
  - [ ] CRUD completo con estados (active, inactive, terminated)
  - [ ] Integración con auth-identity-svc para usuarios
  - [ ] Integración con person-svc para datos personales
  - [ ] Sistema de metadata flexible (JSONB)
  - [ ] Historial de cambios y auditoría

- [ ] **Employee Roles (Asignaciones Múltiples)**
  - [ ] Asignación/revocación de roles con fechas
  - [ ] Sistema de rol primario automático
  - [ ] Validaciones de jerarquía y permisos
  - [ ] Reportes de roles por empleado/organización

- [ ] **Business Logic Avanzada**
  - [ ] Validación de último administrador
  - [ ] Transiciones de estado automáticas
  - [ ] Notificaciones de cambios importantes
  - [ ] Métricas de empleados por organización

#### 2.3 **Sistema de Invitaciones (Semana 3)**

- [ ] **Invitation Workflow**
  - [ ] Creación de invitaciones con tokens seguros
  - [ ] Expiración automática y cleanup
  - [ ] Límites dinámicos por organización
  - [ ] Estados: pending, accepted, expired, cancelled

- [ ] **Automatización**
  - [ ] Integración con notification-svc para emails
  - [ ] Conversión automática a empleado al aceptar
  - [ ] Recordatorios automáticos antes de expiración
  - [ ] Logs detallados del workflow completo

---

### **FASE 3: Hub de Integraciones** ⏱️ (2-3 semanas)

> **🎯 Objetivo**: Implementar el sistema de integraciones externas con terceros

#### 3.1 **Infraestructura de Integraciones**
- [ ] **Catálogo de Tipos de Integración**
  - [ ] Sistema extensible de integration_types
  - [ ] Schemas JSON para configuración
  - [ ] Soporte para API keys, OAuth, webhooks
  - [ ] Categorización por industria/función

- [ ] **Gestión de Credenciales**
  - [ ] Encriptación AES-256 de credenciales sensibles
  - [ ] Rotación automática de tokens
  - [ ] Vault integration para mayor seguridad
  - [ ] Auditoría de acceso a credenciales

#### 3.2 **Integraciones Específicas (MVP)**
- [ ] **QuickBooks Online**
  - [ ] OAuth flow completo
  - [ ] Sincronización de cuentas y transacciones
  - [ ] Webhook para cambios en tiempo real
  - [ ] Rate limiting y error handling

- [ ] **HubSpot CRM**
  - [ ] API key configuration
  - [ ] Sincronización de contactos y deals
  - [ ] Pipeline mapping personalizable
  - [ ] Métricas de performance de sync

- [ ] **Mailchimp Marketing**
  - [ ] Gestión de listas y campañas
  - [ ] Segmentación por sucursal/empleado
  - [ ] Analytics de campañas
  - [ ] Templates personalizables

#### 3.3 **Sistema de Eventos y Workers**
- [ ] **Event Processing**
  - [ ] Queue de eventos con RabbitMQ
  - [ ] Workers para sincronización asíncrona
  - [ ] Retry automático con backoff exponencial
  - [ ] Dead letter queue para errores persistentes

- [ ] **Monitoring y Observabilidad**
  - [ ] Métricas de integraciones por tipo/organización
  - [ ] Alertas por fallos repetidos
  - [ ] Dashboard de salud de integraciones
  - [ ] Logs particionados con cleanup automático

---

### **FASE 4: Dominios Personalizados** ⏱️ (1-2 semanas)

> **🎯 Objetivo**: Sistema completo de gestión de dominios con verificación DNS/SSL

#### 4.1 **Gestión de Dominios**
- [ ] **CRUD de Dominios**
  - [ ] Registro y validación de dominios
  - [ ] Verificación de propiedad via DNS TXT
  - [ ] Estados: pending, verified, failed, expired
  - [ ] Dominio primario único por organización

- [ ] **Verificación DNS Automática**
  - [ ] Generación automática de records requeridos
  - [ ] Worker para verificación periódica
  - [ ] Integración con proveedores DNS populares
  - [ ] Alertas por configuraciones incorrectas

#### 4.2 **SSL y Certificados**
- [ ] **Let's Encrypt Integration**
  - [ ] Generación automática de certificados SSL
  - [ ] Renovación automática antes de expiración
  - [ ] Monitoreo de salud de certificados
  - [ ] Fallback para certificados manuales

- [ ] **Routing y Redirecciones**
  - [ ] Sistema de redirecciones inteligentes
  - [ ] Soporte para subdominios
  - [ ] Cache de configuraciones DNS
  - [ ] Métricas de tráfico por dominio

---

### **FASE 5: APIs y Interfaces** ⏱️ (1-2 semanas)

> **🎯 Objetivo**: APIs REST y gRPC completas con documentación y versionado

#### 5.1 **REST API**
- [ ] **Endpoints Completos**
```http
# Organizations
POST   /api/v1/organizations
GET    /api/v1/organizations
GET    /api/v1/organizations/{id}
PUT    /api/v1/organizations/{id}
DELETE /api/v1/organizations/{id}

# Settings
GET    /api/v1/organizations/{id}/settings
PUT    /api/v1/organizations/{id}/settings

# Branches
POST   /api/v1/organizations/{orgId}/branches
GET    /api/v1/organizations/{orgId}/branches
GET    /api/v1/organizations/{orgId}/branches/{id}
PUT    /api/v1/organizations/{orgId}/branches/{id}
DELETE /api/v1/organizations/{orgId}/branches/{id}

# Employees
POST   /api/v1/organizations/{orgId}/employees
GET    /api/v1/organizations/{orgId}/employees
GET    /api/v1/organizations/{orgId}/employees/{id}
PUT    /api/v1/organizations/{orgId}/employees/{id}
DELETE /api/v1/organizations/{orgId}/employees/{id}

# Employee Roles
POST   /api/v1/employees/{id}/roles
DELETE /api/v1/employees/{id}/roles/{roleId}
PUT    /api/v1/employees/{id}/roles/{roleId}/primary

# Invitations
POST   /api/v1/organizations/{orgId}/invitations
GET    /api/v1/organizations/{orgId}/invitations
GET    /api/v1/invitations/{token}
POST   /api/v1/invitations/{token}/accept
DELETE /api/v1/invitations/{id}

# Integrations
GET    /api/v1/integration-types
POST   /api/v1/organizations/{orgId}/integrations
GET    /api/v1/organizations/{orgId}/integrations
PUT    /api/v1/integrations/{id}
DELETE /api/v1/integrations/{id}
POST   /api/v1/integrations/{id}/test
GET    /api/v1/integrations/{id}/logs

# Custom Domains
POST   /api/v1/organizations/{orgId}/domains
GET    /api/v1/organizations/{orgId}/domains
GET    /api/v1/domains/{id}/dns-records
POST   /api/v1/domains/{id}/verify
PUT    /api/v1/domains/{id}/primary
DELETE /api/v1/domains/{id}
```

#### 5.2 **gRPC Interface**
- [ ] **Service Definition (.proto)**
```protobuf
service OrganizationService {
  // Organizations
  rpc CreateOrganization(CreateOrganizationRequest) returns (OrganizationResponse);
  rpc GetOrganization(GetOrganizationRequest) returns (OrganizationResponse);
  rpc UpdateOrganization(UpdateOrganizationRequest) returns (OrganizationResponse);
  rpc DeleteOrganization(DeleteOrganizationRequest) returns (google.protobuf.Empty);
  rpc ListOrganizations(ListOrganizationsRequest) returns (ListOrganizationsResponse);
  
  // Employees
  rpc CreateEmployee(CreateEmployeeRequest) returns (EmployeeResponse);
  rpc GetEmployee(GetEmployeeRequest) returns (EmployeeResponse);
  rpc AssignEmployeeRole(AssignRoleRequest) returns (EmployeeRoleResponse);
  
  // Internal use
  rpc ValidateOrganizationExists(ValidateOrgRequest) returns (ValidateOrgResponse);
  rpc GetOrganizationSettings(GetSettingsRequest) returns (SettingsResponse);
  rpc GetEmployeesByOrganization(GetEmployeesRequest) returns (GetEmployeesResponse);
}
```

#### 5.3 **Documentación y Versionado**
- [ ] **OpenAPI/Swagger**
  - [ ] Documentación automática de REST API
  - [ ] Ejemplos de request/response
  - [ ] Códigos de error documentados
  - [ ] Interactive API explorer

- [ ] **API Versioning**
  - [ ] Versionado semántico (v1, v2)
  - [ ] Backward compatibility garantizada
  - [ ] Deprecation notices y migración
  - [ ] Feature flags para nuevas funcionalidades

---

### **FASE 6: Seguridad y Compliance** ⏱️ (1-2 semanas)

> **🎯 Objetivo**: Implementar seguridad robusta y compliance para entornos empresariales

#### 6.1 **Autenticación y Autorización**
- [ ] **JWT Integration**
  - [ ] Validación de tokens de auth-identity-svc
  - [ ] Refresh token handling
  - [ ] Role-based access control (RBAC)
  - [ ] Rate limiting por usuario/organización

- [ ] **Permission System**
  - [ ] Matriz detallada de permisos
  - [ ] Validación en cada endpoint
  - [ ] Audit trail de acciones sensibles
  - [ ] Emergency access procedures

#### 6.2 **Encriptación y Datos Sensibles**
- [ ] **Data Encryption**
  - [ ] Encriptación AES-256 para credenciales
  - [ ] Hashing seguro de tokens de invitación
  - [ ] PII encryption para datos personales
  - [ ] Key rotation automática

- [ ] **Compliance**
  - [ ] GDPR compliance para datos europeos
  - [ ] SOC 2 Type II requirements
  - [ ] Data retention policies
  - [ ] Right to be forgotten implementation

#### 6.3 **Auditoría y Monitoreo**
- [ ] **Audit Logs**
  - [ ] Logs completos de todas las acciones
  - [ ] Retención según políticas de compliance
  - [ ] Búsqueda y análisis de logs
  - [ ] Alertas por actividades sospechosas

- [ ] **Security Monitoring**
  - [ ] Detección de anomalías
  - [ ] Failed login attempts tracking
  - [ ] Privilege escalation detection
  - [ ] Integration con SIEM systems

---

### **FASE 7: Performance y Escalabilidad** ⏱️ (1-2 semanas)

> **🎯 Objetivo**: Optimizar performance para escalabilidad empresarial

#### 7.1 **Caching Strategy**
- [ ] **Redis Multi-Layer Cache**
  - [ ] Cache de organizaciones frecuentemente accedidas
  - [ ] Cache de configuraciones y roles
  - [ ] Cache de conteos y métricas
  - [ ] Cache invalidation strategies

- [ ] **Database Optimization**
  - [ ] Query optimization con EXPLAIN ANALYZE
  - [ ] Connection pooling optimizado
  - [ ] Read replicas para queries pesados
  - [ ] Partitioning adicional si es necesario

#### 7.2 **Horizontal Scaling**
- [ ] **Microservice Patterns**
  - [ ] Stateless service design
  - [ ] Load balancing strategies
  - [ ] Circuit breaker para external calls
  - [ ] Graceful degradation

- [ ] **Event-Driven Architecture**
  - [ ] Event sourcing para cambios críticos
  - [ ] CQRS para separation of concerns
  - [ ] Async processing optimization
  - [ ] Message deduplication

#### 7.3 **Monitoring y Observabilidad**
- [ ] **Metrics & Monitoring**
  - [ ] Prometheus metrics para business KPIs
  - [ ] Grafana dashboards
  - [ ] Distributed tracing con Jaeger
  - [ ] Custom alerts por thresholds

- [ ] **Performance Testing**
  - [ ] Load testing con k6 o Artillery
  - [ ] Stress testing scenarios
  - [ ] Performance regression tests
  - [ ] Capacity planning guidelines

---

### **FASE 8: Funcionalidades Avanzadas** ⏱️ (2-3 semanas)

> **🎯 Objetivo**: Funcionalidades innovadoras para diferenciación competitiva

#### 8.1 **Advanced Analytics**
- [ ] **Business Intelligence**
  - [ ] Dashboard de métricas organizacionales
  - [ ] Reportes de performance de empleados
  - [ ] Analytics de integraciones y uso
  - [ ] Forecasting con ML básico

- [ ] **Data Export/Import**
  - [ ] Bulk operations para migraciones
  - [ ] Export compliance (GDPR, etc.)
  - [ ] Integration con data warehouses
  - [ ] ETL pipelines para analytics

#### 8.2 **Workflow Automation**
- [ ] **Business Process Automation**
  - [ ] Workflows personalizables por organización
  - [ ] Triggers automáticos por eventos
  - [ ] Integration con tools externos (Zapier)
  - [ ] Visual workflow builder (futuro)

- [ ] **AI-Powered Features**
  - [ ] Smart role assignment suggestions
  - [ ] Anomaly detection en usage patterns
  - [ ] Predictive analytics para churn
  - [ ] NLP para categorización automática

#### 8.3 **Enterprise Features**
- [ ] **Multi-Tenancy Avanzado**
  - [ ] Tenant isolation strategies
  - [ ] Resource quotas por tenant
  - [ ] Custom branding por organización
  - [ ] White-label capabilities

- [ ] **Integration Marketplace**
  - [ ] Plugin system para integraciones custom
  - [ ] Marketplace de conectores
  - [ ] SDK para developers third-party
  - [ ] Revenue sharing model

---

## 📋 Plan de Releases

### **Release 1.0 - MVP Core** (6-8 semanas)
**🎯 Funcionalidades Mínimas Viables**
- ✅ CRUD completo de todas las entidades
- ✅ APIs REST documentadas
- ✅ Integraciones básicas (auth, person, address)
- ✅ Sistema de invitaciones funcional
- ✅ Seguridad básica y validaciones

### **Release 1.5 - Integration Hub** (2-3 semanas)
**🔗 Ecosistema de Integraciones**
- ✅ 3-5 integraciones principales implementadas
- ✅ Sistema de eventos robusto
- ✅ Monitoring y alertas de integraciones
- ✅ Documentación de APIs de integración

### **Release 2.0 - Enterprise Ready** (3-4 semanas)
**🏢 Características Empresariales**
- ✅ Dominios personalizados con SSL
- ✅ Advanced security y compliance
- ✅ Performance optimization completa
- ✅ Analytics y reporting básico

### **Release 2.5 - Advanced Features** (4-6 semanas)
**🚀 Diferenciación Competitiva**
- ✅ Workflow automation
- ✅ AI-powered features básicas
- ✅ Advanced analytics y BI
- ✅ Multi-tenancy avanzado

---

## 🔧 Consideraciones Técnicas

### **Arquitectura y Patrones**
- **Clean Architecture**: Separación clara de responsabilidades
- **Domain-Driven Design**: Modelado rico del dominio inmobiliario
- **Event Sourcing**: Para cambios críticos y auditoría
- **CQRS**: Separación de comandos y queries
- **Circuit Breaker**: Resilience ante fallos de servicios externos

### **Tecnologías Core**
- **Backend**: Go 1.21+ con Echo v4 framework
- **Database**: PostgreSQL 13+ con funcionalidades avanzadas
- **Cache**: Redis 6+ para performance
- **Messaging**: RabbitMQ para eventos asíncronos
- **Monitoring**: Prometheus + Grafana + Jaeger

### **DevOps y Deployment**
- **Containerization**: Docker multi-stage builds
- **Orchestration**: Kubernetes con Helm charts
- **CI/CD**: GitHub Actions con automated testing
- **Infrastructure**: Terraform para IaC
- **Security**: HashiCorp Vault para secrets

### **Testing Strategy**
- **Unit Tests**: 80%+ coverage mínimo
- **Integration Tests**: Database y external services
- **End-to-End Tests**: Critical user journeys
- **Performance Tests**: Load y stress testing
- **Security Tests**: OWASP compliance scanning

---

## 📊 Métricas y KPIs

### **Business Metrics**
- Organizaciones activas y crecimiento mensual
- Empleados por organización (promedio/distribución)
- Tasa de adopción de integraciones
- Tiempo de onboarding promedio
- Dominios personalizados activos

### **Technical Metrics**
- API response times (p95, p99)
- Error rates por endpoint
- Database query performance
- Cache hit ratios
- Event processing latency

### **User Experience Metrics**
- Invitation acceptance rate
- Time to first value
- Feature adoption rates
- Support ticket volume
- User satisfaction scores

---

## 🏷️ Versionado y Compatibility

### **Semantic Versioning**
- **MAJOR.MINOR.PATCH** (ej: 2.1.3)
- **MAJOR**: Breaking changes en API
- **MINOR**: New features backward compatible
- **PATCH**: Bug fixes y improvements

### **Backward Compatibility**
- Mantener compatibilidad por **2 versiones mayores**
- Deprecation notices con **6 meses** de anticipación
- Migration guides para breaking changes
- Feature flags para gradual rollouts

### **Database Migrations**
- **Never breaking changes** en migrations
- **Rollback capability** para todas las migrations
- **Blue-green deployments** para zero downtime
- **Data validation** post-migration

---

## 📞 Equipo y Responsabilidades

### **Roles Requeridos**
- **Backend Lead**: Arquitectura y code reviews
- **Full-Stack Developers (2-3)**: Feature implementation
- **DevOps Engineer**: Infrastructure y CI/CD
- **QA Engineer**: Testing strategy y automation
- **Product Owner**: Requirements y prioritization

### **Ceremoniás Ágiles**
- **Sprint Planning**: 2 weeks sprints
- **Daily Standups**: Progress y blockers
- **Sprint Reviews**: Demo y feedback
- **Retrospectives**: Continuous improvement

### **Definition of Done**
- [ ] Feature completamente implementada
- [ ] Tests unitarios y de integración passing
- [ ] Code review aprobado por 2+ developers
- [ ] Documentation actualizada
- [ ] Performance benchmarks validados
- [ ] Security review completado

---

*📅 **Roadmap actualizado**: Julio 2025*  
*🔄 **Próxima revisión**: Septiembre 2025*  
*📧 **Contacto**: backend@rem-system.com*

# Dependencies
GRPC_MOCK=true              # Usar mocks en desarrollo
SUBSCRIPTION_BILLING_ADDR=  # Vacío = mock
PROPERTY_SVC_ADDR=          # Vacío = mock
RBAC_SVC_ADDR=              # Vacío = mock

# Auth
JWT_SECRET=your-secret-key

# Crypto
ORGSVC_CRYPTO_KEY=32-byte-key-for-integrations

# Cache
REDIS_URL=redis://localhost:6379

# Features
SUB_CHECK=false             # Feature flag para validación de límites
```

#### 1.3 Server Setup
- [ ] **🚀 cmd/api/main.go** - Entry point
  - [ ] Echo v4 server en puerto 8083
  - [ ] Middlewares: logging, recovery, request_id, user_auth
  - [ ] Graceful shutdown
  - [ ] Health check endpoint
  - [ ] Integración con rem-common/config

#### 1.4 Database & Migrations
- [ ] **🗃️ cmd/migrate/main.go**
  - [ ] Reutilizar `rem-common/db` adaptive migrator
  - [ ] Soporte para rollback
  - [ ] Auto-migrate en desarrollo

#### 1.5 Docker Integration
- [ ] **🐳 docker-compose integration**
  - [ ] Puerto 8083 en `dev-full-compose.yml`
  - [ ] Health check: `GET /health`
  - [ ] Dependencies: postgres, redis
  - [ ] Volúmenes para hot-reload

---

### **FASE 2: Base de Datos y Migraciones** ⏱️ (2-3 días)

#### 2.1 Scripts de Migración
**📁 Ruta**: `organization-svc/migrations/pg/`

| Nº | Archivo | Contenido Principal |
|----|---------|-------------------|
| **0001** | `create_organization_core.up.sql` | `organization`, `organization_settings`, `organization_owner` |
| **0002** | `create_structure.up.sql` | `organization_branch`, `organization_role` |
| **0003** | `create_employee.up.sql` | `employees`, `employee_roles` |
| **0004** | `create_invite.up.sql` | `organization_invite` |
| **0005** | `create_integration_domain.up.sql` | `organization_integration`, `organization_domain` |
| **0006** | `seed_roles.up.sql` | INSERT roles predeterminados con función PL/pgSQL |
| **0007** | `add_audit_triggers.up.sql` | Triggers para updated_at automático |

#### 2.2 Ejemplo de Migración Principal
```sql
-- migrations/pg/0001_create_organization_core.up.sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Tabla principal de organizaciones
CREATE TABLE organization (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  display_name varchar(120) NOT NULL,
  logo_url text,
  fiscal_address_id uuid, -- Referencia a address-svc
  matricula varchar(32),
  status varchar(12) NOT NULL DEFAULT 'active'
    CHECK (status IN ('active', 'suspended', 'deleted')),
  
  -- Auditoría
  created_at timestamptz NOT NULL DEFAULT now(),
  created_by uuid, -- Referencia a auth-identity-svc users.id
  updated_at timestamptz,
  updated_by uuid,
  deleted_at timestamptz -- Soft delete
);

-- Configuraciones clave-valor por organización
CREATE TABLE organization_settings (
  organization_id uuid NOT NULL REFERENCES organization(id) ON DELETE CASCADE,
  key varchar(64) NOT NULL,
  value jsonb,
  updated_at timestamptz NOT NULL DEFAULT now(),
  updated_by uuid,
  PRIMARY KEY (organization_id, key)
);

-- Dueños de la organización
CREATE TABLE organization_owner (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  organization_id uuid NOT NULL REFERENCES organization(id) ON DELETE CASCADE,
  owner_type varchar(12) NOT NULL CHECK (owner_type IN ('individual', 'company')),
  owner_id uuid NOT NULL, -- Referencia a person-svc person.id
  created_at timestamptz NOT NULL DEFAULT now(),
  created_by uuid,
  deleted_at timestamptz,
  
  -- Constraint: un owner no puede estar duplicado en la misma org
  UNIQUE (owner_type, owner_id, organization_id)
);

-- Índices de performance
CREATE INDEX idx_organization_status ON organization(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_organization_created_by ON organization(created_by);
CREATE INDEX idx_org_owner_lookup ON organization_owner(owner_type, owner_id) WHERE deleted_at IS NULL;
```

#### 2.3 Seed de Roles Predeterminados
```sql
-- migrations/pg/0006_seed_roles.up.sql
CREATE OR REPLACE FUNCTION seed_default_roles(org_id uuid)
RETURNS void AS $$
BEGIN
  INSERT INTO organization_role (organization_id, name, description, created_at)
  VALUES 
    (org_id, 'Admin', 'Administrador con acceso completo', now()),
    (org_id, 'Manager', 'Gerente con permisos de gestión', now()),
    (org_id, 'Agent', 'Agente con permisos básicos', now()),
    (org_id, 'Assistant', 'Asistente con permisos limitados', now())
  ON CONFLICT (organization_id, name) DO NOTHING;
END;
$$ LANGUAGE plpgsql;
```

---

### **FASE 3: Entidades Core** ⏱️ (5-7 días)

#### 3.1 Organization (Entidad Principal)

**📄 Archivos involucrados:**
- `src/models/organization.go` - GORM model
- `src/dto/organization.go` - DTOs para API
- `src/repository/organization_repo.go` - Data access
- `src/services/organization_service.go` - Business logic
- `src/controllers/organization_controller.go` - HTTP handlers
- `src/validations/organization.go` - Validation rules

**🔧 Especificaciones técnicas:**

- [ ] **Model** (`models/organization.go`)
```go
type Organization struct {
    ID               uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
    DisplayName      string     `gorm:"size:120;not null" validate:"required,max=120"`
    LogoURL          *string    `gorm:"type:text"`
    FiscalAddressID  *uuid.UUID `gorm:"type:uuid"`
    Matricula        *string    `gorm:"size:32"`
    Status           string     `gorm:"size:12;default:active" validate:"oneof=active suspended deleted"`
    
    // Auditoría
    CreatedAt time.Time  `gorm:"not null;default:now()"`
    CreatedBy *uuid.UUID `gorm:"type:uuid"`
    UpdatedAt *time.Time
    UpdatedBy *uuid.UUID `gorm:"type:uuid"`
    DeletedAt *time.Time `gorm:"index"`
    
    // Relaciones
    Settings    []OrganizationSetting `gorm:"foreignKey:OrganizationID"`
    Branches    []OrganizationBranch  `gorm:"foreignKey:OrganizationID"`
    Employees   []Employee            `gorm:"foreignKey:OrganizationID"`
    Roles       []OrganizationRole    `gorm:"foreignKey:OrganizationID"`
}
```

- [ ] **Repository** - CRUD básico con GORM
  - [ ] `Create(org *Organization) error`
  - [ ] `GetByID(id uuid.UUID) (*Organization, error)`
  - [ ] `List(filters ListFilters) ([]Organization, int64, error)`
  - [ ] `Update(id uuid.UUID, updates *Organization) error`
  - [ ] `SoftDelete(id uuid.UUID, deletedBy uuid.UUID) error`
  - [ ] `GetByUserID(userID uuid.UUID) ([]Organization, error)` // Para membresías

- [ ] **Service** - Lógica de negocio
  - [ ] Validar `display_name` único (soft-delete aware)
  - [ ] Auto-seed de roles predeterminados en creación
  - [ ] Validar `fiscal_address_id` existe (address-svc gRPC)
  - [ ] Enforcement de límites de plan (subscription-svc mock)

- [ ] **Controller** - HTTP Endpoints
  - [ ] `POST /organizations` - Crear inmobiliaria
    - [ ] Body: `CreateOrganizationRequest`
    - [ ] Auth: Usuario autenticado
    - [ ] Logic: Crear org + seed roles + asignar owner como Admin
  - [ ] `GET /organizations/:id` - Obtener por ID
    - [ ] Auth: Miembro de la organización
    - [ ] Include: settings, branches (opcional)
  - [ ] `PUT /organizations/:id` - Actualizar
    - [ ] Auth: Rol Admin en la organización
    - [ ] Campos actualizables: display_name, logo_url, matricula
  - [ ] `DELETE /organizations/:id` - Soft delete
    - [ ] Auth: Rol Admin + confirmación adicional
  - [ ] `GET /organizations` - Listar con filtros
    - [ ] Auth: Usuario autenticado
    - [ ] Filtros: status, search (display_name), created_after
    - [ ] Paginación: page, per_page

#### 3.2 Organization Settings (Sistema KV)

**📄 Archivos:**
- `src/models/organization_setting.go`
- `src/controllers/settings_controller.go`

**🔧 Implementación:**

- [ ] **Model**
```go
type OrganizationSetting struct {
    OrganizationID uuid.UUID   `gorm:"type:uuid;primary_key"`
    Key           string      `gorm:"size:64;primary_key"`
    Value         *string     `gorm:"type:jsonb"`
    UpdatedAt     time.Time   `gorm:"not null;default:now()"`
    UpdatedBy     *uuid.UUID  `gorm:"type:uuid"`
}
```

- [ ] **Endpoints**:
  - [ ] `GET /organizations/:id/settings` - Todas las configuraciones
    - [ ] Response: `map[string]interface{}`
    - [ ] Auth: Miembro de la organización
  - [ ] `PUT /organizations/:id/settings/:key` - Upsert setting
    - [ ] Body: `{"value": any}`
    - [ ] Auth: Rol Manager+ en la organización
  - [ ] `DELETE /organizations/:id/settings/:key` - Eliminar setting
    - [ ] Auth: Rol Admin en la organización

- [ ] **Configuraciones comunes**:
  - `timezone` - "America/Argentina/Buenos_Aires"
  - `currency` - "ARS"
  - `default_commission` - `{"percentage": 3.5}`
  - `branding` - `{"primary_color": "#FF6B35", "secondary_color": "#F5F5F5"}`

#### 3.3 Organization Owners

**📄 Archivos:**
- `src/models/organization_owner.go`
- `src/controllers/owner_controller.go`
- `src/grpc/client_person.go` (existente, mejorar)

**🔧 Integración con person-svc:**

- [ ] **Model**
```go
type OrganizationOwner struct {
    ID             uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
    OrganizationID uuid.UUID  `gorm:"type:uuid;not null"`
    OwnerType      string     `gorm:"size:12;not null" validate:"oneof=individual company"`
    OwnerID        uuid.UUID  `gorm:"type:uuid;not null"`
    CreatedAt      time.Time  `gorm:"not null;default:now()"`
    CreatedBy      *uuid.UUID `gorm:"type:uuid"`
    DeletedAt      *time.Time `gorm:"index"`
    
    // Constraint único en DB: (owner_type, owner_id, organization_id)
}
```

- [ ] **gRPC Client Enhancement**
  - [ ] `ValidatePersonExists(id uuid.UUID) (bool, error)`
  - [ ] `CreatePersonFromPayload(payload CreatePersonPayload) (*PersonResponse, error)`
  - [ ] `GetPersonDetails(id uuid.UUID) (*PersonResponse, error)`

- [ ] **Endpoints**:
  - [ ] `POST /organizations/:id/owners` - Agregar dueño
    - [ ] Body: `{"owner_type": "individual|company", "owner_id": "uuid"}` 
    - [ ] O: `{"owner_type": "company", "person_data": {...}}` (crear + link)
    - [ ] Auth: Rol Admin
    - [ ] Validar owner existe en person-svc
  - [ ] `GET /organizations/:id/owners` - Listar dueños con detalles
    - [ ] Include datos de person-svc: name, dni/cuit, contacts
    - [ ] Auth: Miembro de la organización
  - [ ] `DELETE /organizations/:id/owners/:owner_id` - Quitar dueño
    - [ ] Auth: Rol Admin
    - [ ] Soft delete (no eliminar la persona)

---

### **FASE 4: Estructura Organizacional** ⏱️ (4-6 días)

#### 4.1 Organization Branches (Sucursales)

**📄 Archivos:**
- `src/models/organization_branch.go`
- `src/services/branch_service.go`
- `src/controllers/branch_controller.go`
- `src/grpc/client_address.go` (existente)

**🔧 Especificaciones:**

- [ ] **Model**
```go
type OrganizationBranch struct {
    ID             uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
    OrganizationID uuid.UUID  `gorm:"type:uuid;not null"`
    DisplayName    string     `gorm:"size:120;not null"`
    AddressID      *uuid.UUID `gorm:"type:uuid"`
    Phone          *string    `gorm:"size:32"`
    Email          *string    `gorm:"size:160"`
    IsMain         bool       `gorm:"default:false"`
    
    // Auditoría
    CreatedAt time.Time  `gorm:"not null;default:now()"`
    CreatedBy *uuid.UUID
    UpdatedAt *time.Time
    UpdatedBy *uuid.UUID
    DeletedAt *time.Time `gorm:"index"`
}
```

- [ ] **Business Rules**
  - [ ] Solo UNA sucursal `is_main=true` por organización
  - [ ] Trigger/constraint en DB + validación en service
  - [ ] Al crear primera sucursal, auto-set `is_main=true`
  - [ ] No se puede eliminar la sucursal principal si hay otras

- [ ] **Address Integration**
  - [ ] Crear dirección automáticamente si se envía `address_payload`
  - [ ] Validar `address_id` existe si se envía ID
  - [ ] Incluir datos de dirección en responses (join via gRPC)

- [ ] **Endpoints**:
  - [ ] `POST /organizations/:id/branches` - Crear sucursal
  - [ ] `GET /organizations/:id/branches` - Listar sucursales
  - [ ] `PUT /organizations/:id/branches/:branch_id` - Actualizar
  - [ ] `PUT /organizations/:id/branches/:branch_id/set-main` - Cambiar principal
  - [ ] `DELETE /organizations/:id/branches/:branch_id` - Eliminar

#### 4.2 Organization Roles (Roles Internos)

**📄 Archivos:**
- `src/models/organization_role.go`
- `src/services/role_service.go`
- `src/controllers/role_controller.go`

**🔧 Implementación:**

- [ ] **Model**
```go
type OrganizationRole struct {
    ID             int        `gorm:"primary_key;auto_increment"`
    OrganizationID uuid.UUID  `gorm:"type:uuid;not null"`
    Name           string     `gorm:"size:64;not null"`
    Description    *string    `gorm:"type:text"`
    IsDefault      bool       `gorm:"default:false"` // Para roles seed
    
    CreatedAt time.Time  `gorm:"not null;default:now()"`
    CreatedBy *uuid.UUID
    UpdatedAt *time.Time  
    UpdatedBy *uuid.UUID
    DeletedAt *time.Time `gorm:"index"`
    
    // Constraint único: (organization_id, name) WHERE deleted_at IS NULL
}
```

- [ ] **Business Logic**
  - [ ] Roles predeterminados NO se pueden eliminar (`is_default=true`)
  - [ ] Validar que no hay empleados asignados antes de eliminar rol custom
  - [ ] `name` único por organización (considerando soft-delete)

- [ ] **Endpoints**:
  - [ ] `POST /organizations/:id/roles` - Crear rol personalizado
  - [ ] `GET /organizations/:id/roles` - Listar todos los roles
  - [ ] `PUT /organizations/:id/roles/:role_id` - Actualizar rol
  - [ ] `DELETE /organizations/:id/roles/:role_id` - Eliminar rol custom

---

### **FASE 5: Gestión de Empleados** ⏱️ (6-8 días)

#### 5.1 Employees

**📄 Archivos:**
- `src/models/employee.go`
- `src/models/employee_role.go`  
- `src/services/employee_service.go`
- `src/controllers/employee_controller.go`
- `src/grpc/client_auth.go` (existente)
- `src/grpc/client_subs.go` (mock)

**🔧 Models:**

```go
type Employee struct {
    ID             uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
    OrganizationID uuid.UUID  `gorm:"type:uuid;not null"`
    UserID         uuid.UUID  `gorm:"type:uuid;not null"` // auth-identity-svc
    PrimaryRoleID  *int       `gorm:"type:int"` // FK a organization_role
    HiredAt        *time.Time `gorm:"type:date"`
    FiredAt        *time.Time `gorm:"type:date"`
    
    CreatedAt time.Time  `gorm:"not null;default:now()"`
    CreatedBy *uuid.UUID
    UpdatedAt *time.Time
    UpdatedBy *uuid.UUID
    DeletedAt *time.Time `gorm:"index"`
    
    // Relations
    Roles []EmployeeRole `gorm:"foreignKey:EmployeeID"`
}

type EmployeeRole struct {
    EmployeeID uuid.UUID `gorm:"type:uuid;primary_key"`
    RoleID     int       `gorm:"type:int;primary_key"`
    Primary    bool      `gorm:"default:false"`
    CreatedAt  time.Time `gorm:"not null;default:now()"`
    CreatedBy  *uuid.UUID
}
```

**🔧 Business Rules:**
- [ ] **Constraint único**: `(user_id, organization_id)` - Un user solo puede estar una vez por org
- [ ] **Limits enforcement**: Validar `employees_limit` del plan (via subscription-svc mock)
- [ ] **Primary role único**: Solo un rol puede ser `primary=true` por employee
- [ ] **Estados**: active (hired_at set, fired_at null), fired (fired_at set)

**🔧 Integration Points:**

- [ ] **auth-identity-svc gRPC**
  - [ ] `GetUserById(user_id)` - Validar que user existe
  - [ ] `CreateUser(account_id, person_id?)` - Para invitaciones aceptadas
  - [ ] `DeactivateUser(user_id)` - Al despedir empleado

- [ ] **subscription-billing-svc Mock**
  - [ ] `GetOrgSubscription(org_id)` → mock responde employees_limit=50
  - [ ] Validar antes de agregar nuevo empleado
  - [ ] Feature flag `SUB_CHECK=false` para saltear validación

**🔧 Endpoints:**

- [ ] `POST /organizations/:id/employees` - Agregar empleado existente
  - [ ] Body: `{"user_id": "uuid", "role_id": 1, "hired_at": "2025-01-15"}`
  - [ ] Auth: Rol Manager+ en la organización
  - [ ] Validar user existe en auth-identity-svc
  - [ ] Validar límites del plan
  - [ ] Auto-asignar como primary role

- [ ] `GET /organizations/:id/employees` - Listar empleados
  - [ ] Include: user data (email), person data (name), roles
  - [ ] Filtros: status (active/fired), role_id, hired_after
  - [ ] Auth: Miembro de la organización

- [ ] `PUT /organizations/:id/employees/:employee_id` - Actualizar empleado
  - [ ] Campos: hired_at, primary_role_id
  - [ ] Auth: Rol Manager+

- [ ] `DELETE /organizations/:id/employees/:employee_id` - Despedir empleado
  - [ ] Set fired_at = now()
  - [ ] Llamar auth-identity-svc.DeactivateUser()
  - [ ] Auth: Rol Admin

#### 5.2 Employee Roles (Asignación N-a-N)

**🔧 Endpoints:**

- [ ] `POST /organizations/:id/employees/:employee_id/roles` - Asignar rol adicional
  - [ ] Body: `{"role_id": 2}`
  - [ ] Validar rol pertenece a la organización
  - [ ] Auth: Rol Admin

- [ ] `DELETE /organizations/:id/employees/:employee_id/roles/:role_id` - Quitar rol
  - [ ] No permitir quitar el primary role (usar PUT primary-role antes)
  - [ ] Auth: Rol Admin

- [ ] `PUT /organizations/:id/employees/:employee_id/primary-role` - Cambiar rol primario
  - [ ] Body: `{"role_id": 3}`
  - [ ] Validar empleado ya tiene ese rol asignado
  - [ ] Update primary=false en rol anterior, primary=true en nuevo
  - [ ] Auth: Rol Admin

#### 5.3 Middleware EmployeeAuth

**📄 Archivo:** `rem-common/middleware/employee_auth.go` (mejora del existente)

**🔧 Funcionalidad:**
- [ ] Extraer `org_id` de la URL (`:id` parameter)
- [ ] Validar que el usuario JWT es empleado de esa organización
- [ ] Setear en context: `X-Org-ID`, `X-Employee-ID`, `X-Primary-Role`
- [ ] Support para múltiples organizaciones por usuario

```go
// Context keys que se setean
type ContextKey string
const (
    OrgIDKey      ContextKey = "org_id"
    EmployeeIDKey ContextKey = "employee_id"
    PrimaryRoleKey ContextKey = "primary_role"
)
```

---

### **FASE 6: Sistema de Invitaciones** ⏱️ (4-5 días)

#### 6.1 Organization Invites

**📄 Archivos:**
- `src/models/organization_invite.go`
- `src/services/invite_service.go`
- `src/controllers/invite_controller.go`
- `src/utils/tokens.go`

**🔧 Model:**

```go
type OrganizationInvite struct {
    ID             uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
    OrganizationID uuid.UUID  `gorm:"type:uuid;not null"`
    RoleID         int        `gorm:"type:int;not null"` // FK a organization_role
    Email          string     `gorm:"size:160;not null"`
    Status         string     `gorm:"size:12;default:pending"`
    Token          string     `gorm:"size:64;unique;not null"`
    ExpiresAt      time.Time  `gorm:"not null"`
    AcceptedAt     *time.Time
    
    CreatedAt time.Time  `gorm:"not null;default:now()"`
    CreatedBy *uuid.UUID
    UpdatedAt *time.Time
    UpdatedBy *uuid.UUID
    DeletedAt *time.Time `gorm:"index"`
}
```

**🔧 Token Generation:**
```go
// utils/tokens.go
func GenerateInviteToken() string {
    // SHA-256 hash de random bytes
    // Formato: inv_[32 chars]
    // Expiry: 7 días por defecto
}
```

**🔧 Business Logic:**

- [ ] **Validaciones pre-creación**
  - [ ] Email no duplicado en invitaciones pending de la org
  - [ ] Email no corresponde a empleado actual de la org
  - [ ] Límite de empleados del plan no alcanzado

- [ ] **Estados**: `pending`, `accepted`, `declined`, `expired`
- [ ] **Auto-expiry**: Job/worker que marca como `expired` después de ExpiresAt
- [ ] **Email notifications**: Placeholder para integración futura

**🔧 Endpoints:**

- [ ] `POST /organizations/:id/invites` - Crear invitación
  - [ ] Body: `{"email": "user@example.com", "role_id": 2}`
  - [ ] Auth: Rol Manager+
  - [ ] Response: incluye `token` para testing

- [ ] `GET /organizations/:id/invites` - Listar invitaciones
  - [ ] Filtros: status, created_after
  - [ ] Auth: Rol Manager+

- [ ] `DELETE /organizations/:id/invites/:invite_id` - Cancelar invitación
  - [ ] Solo si status=pending
  - [ ] Auth: Rol Manager+

#### 6.2 Flujo de Aceptación (Endpoints Públicos)

**🔧 Public Endpoints:**

- [ ] `GET /invites/:token` - Ver detalles de invitación
  - [ ] Response: org name, role name, inviter name, expires_at
  - [ ] No auth required
  - [ ] Validar token válido y no expirado

- [ ] `PUT /invites/:token/accept` - Aceptar invitación
  - [ ] Body: `{"account_data": {...}}` (si no tiene cuenta)
  - [ ] O: autenticado con JWT existente
  - [ ] **Flujo**:
    1. Validar token + expiry
    2. Si no tiene account: crear via auth-identity-svc
    3. Si no tiene user: crear user linkado al account
    4. Crear employee record + asignar rol
    5. Marcar invite como accepted
  - [ ] Response: JWT token + organization data

- [ ] `PUT /invites/:token/decline` - Rechazar invitación
  - [ ] Marcar invite como declined
  - [ ] No auth required

#### 6.3 Auto-account Creation

**🔧 Integration con auth-identity-svc:**

```go
// Cuando se acepta invite sin account existente
func (s *InviteService) AcceptInviteNewUser(token string, accountData CreateAccountRequest) error {
    // 1. Validar invite
    invite := s.repo.GetInviteByToken(token)
    
    // 2. Crear account en auth-identity-svc
    account := s.authClient.CreateAccount(accountData)
    
    // 3. Crear user en auth-identity-svc
    user := s.authClient.CreateUser(account.ID, nil) // person_id null por ahora
    
    // 4. Crear employee
    employee := s.employeeService.CreateFromInvite(invite, user.ID)
    
    // 5. Update invite
    invite.Status = "accepted"
    invite.AcceptedAt = &now
}
```

---

### **FASE 7: Integraciones y Dominios** ⏱️ (3-4 días)

#### 7.1 Organization Integrations

**📄 Archivos:**
- `src/models/organization_integration.go`
- `src/services/integration_service.go`
- `src/controllers/integration_controller.go`
- `src/utils/encryption.go`

**🔧 Model:**

```go
type OrganizationIntegration struct {
    ID             uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
    OrganizationID uuid.UUID  `gorm:"type:uuid;not null"`
    Provider       string     `gorm:"size:32;not null"` // zonaprop, mercadolibre, whatsapp
    ExternalID     *string    `gorm:"size:128"`
    Config         *string    `gorm:"type:jsonb"` // Encrypted JSON
    Status         string     `gorm:"size:12;default:connected"`
    ConnectedAt    *time.Time
    
    CreatedAt time.Time  `gorm:"not null;default:now()"`
    CreatedBy *uuid.UUID
    UpdatedAt *time.Time
    UpdatedBy *uuid.UUID
    DeletedAt *time.Time `gorm:"index"`
    
    // Constraint único: (organization_id, provider)
}
```

**🔧 Encryption:**
```go
// utils/encryption.go
func EncryptConfig(plaintext string) (string, error) {
    // AES-256-GCM con key desde ORGSVC_CRYPTO_KEY
}

func DecryptConfig(ciphertext string) (string, error) {
    // Decrypt y return plain JSON
}
```

**🔧 Providers Support:**

- [ ] **ZonaProp**
  - [ ] Config: `{"api_key": "...", "account_id": "...", "sync_enabled": true}`
  - [ ] Validación: Test API call para verificar credenciales

- [ ] **WhatsApp Business**
  - [ ] Config: `{"phone_number": "+54911...", "access_token": "...", "webhook_verify_token": "..."}`

- [ ] **MercadoLibre**
  - [ ] Config: `{"client_id": "...", "client_secret": "...", "access_token": "..."}`

**🔧 Endpoints:**

- [ ] `POST /organizations/:id/integrations` - Crear integración
  - [ ] Body: `{"provider": "zonaprop", "config": {...}}`
  - [ ] Validar credenciales con provider (ping API)
  - [ ] Encriptar config antes de almacenar
  - [ ] Auth: Rol Admin

- [ ] `GET /organizations/:id/integrations` - Listar integraciones
  - [ ] Response: provider, status, connected_at (config omitido por seguridad)
  - [ ] Auth: Rol Manager+

- [ ] `PUT /organizations/:id/integrations/:integration_id` - Actualizar config
  - [ ] Re-validar credenciales
  - [ ] Auth: Rol Admin

- [ ] `PUT /organizations/:id/integrations/:integration_id/test` - Test conexión
  - [ ] Hacer ping al provider API
  - [ ] Response: connection status + details
  - [ ] Auth: Rol Manager+

- [ ] `DELETE /organizations/:id/integrations/:integration_id` - Eliminar
  - [ ] Confirmar que no hay sync jobs activos (placeholder)
  - [ ] Auth: Rol Admin

#### 7.2 Organization Domains

**📄 Archivos:**
- `src/models/organization_domain.go`
- `src/services/domain_service.go`
- `src/controllers/domain_controller.go`

**🔧 Model:**

```go
type OrganizationDomain struct {
    ID             uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
    OrganizationID uuid.UUID  `gorm:"type:uuid;not null"`
    Domain         string     `gorm:"size:255;unique;not null"`
    SSLEnabled     bool       `gorm:"default:false"`
    VerifiedAt     *time.Time
    
    CreatedAt time.Time  `gorm:"not null;default:now()"`
    CreatedBy *uuid.UUID
    UpdatedAt *time.Time
    UpdatedBy *uuid.UUID
    DeletedAt *time.Time `gorm:"index"`
}
```

**🔧 Domain Verification:**
- [ ] **DNS TXT Record**: `rem-verify=<token>` en el dominio
- [ ] **HTTP File**: `/.well-known/rem-verification.txt` con token
- [ ] **Meta Tag**: `<meta name="rem-verification" content="<token>">`

**🔧 Endpoints:**

- [ ] `POST /organizations/:id/domains` - Agregar dominio
  - [ ] Body: `{"domain": "example.com"}`
  - [ ] Validar formato de dominio
  - [ ] Generar verification token
  - [ ] Auth: Rol Admin

- [ ] `GET /organizations/:id/domains` - Listar dominios
  - [ ] Include verification status y método
  - [ ] Auth: Rol Manager+

- [ ] `PUT /organizations/:id/domains/:domain_id/verify` - Verificar dominio
  - [ ] Intentar los 3 métodos de verificación
  - [ ] Update verified_at si éxito
  - [ ] Response: verification status + details
  - [ ] Auth: Rol Admin

- [ ] `PUT /organizations/:id/domains/:domain_id/ssl` - Toggle SSL
  - [ ] Body: `{"ssl_enabled": true}`
  - [ ] Validar dominio verificado
  - [ ] Auth: Rol Admin

- [ ] `DELETE /organizations/:id/domains/:domain_id` - Eliminar dominio
  - [ ] Auth: Rol Admin

---

### **FASE 8: Subscripciones (Solo Lectura)** ⏱️ (2-3 días)

#### 8.1 Mock Subscription Client

**📄 Archivos:**
- `src/grpc/client_subs.go`
- `rem-common/testing/mocks/fake_subs_client.go`
- `src/grpc/interfaces.go`

**🔧 Interface:**

```go
// src/grpc/interfaces.go
type SubscriptionClient interface {
    GetOrgSubscription(ctx context.Context, orgID uuid.UUID) (*SubscriptionSummary, error)
    GetUsageLimits(ctx context.Context, orgID uuid.UUID) (*UsageLimits, error)
}

type SubscriptionSummary struct {
    OrganizationID   uuid.UUID `json:"organization_id"`
    PlanCode         string    `json:"plan_code"`
    Status           string    `json:"status"`
    ListingsLimit    int64     `json:"listings_limit"`
    ListingsUsed     int64     `json:"listings_used"`
    EmployeesLimit   int64     `json:"employees_limit"`
    EmployeesUsed    int64     `json:"employees_used"`
    StartsAt         time.Time `json:"starts_at"`
    RenewsAt         time.Time `json:"renews_at"`
}
```

**🔧 Mock Implementation:**

```go
// rem-common/testing/mocks/fake_subs_client.go
//go:build mock

type FakeSubsClient struct{}

func (f *FakeSubsClient) GetOrgSubscription(ctx context.Context, orgID uuid.UUID) (*SubscriptionSummary, error) {
    return &SubscriptionSummary{
        OrganizationID:   orgID,
        PlanCode:         "PRO",
        Status:           "active",
        ListingsLimit:    9999,
        ListingsUsed:     0,
        EmployeesLimit:   50,
        EmployeesUsed:    1, // Se cuenta dinámicamente
        StartsAt:         time.Now().AddDate(0, -1, 0),
        RenewsAt:         time.Now().AddDate(0, 11, 0),
    }, nil
}
```

**🔧 DI Container:**

```go
// src/config/dependencies.go
func NewSubscriptionClient(cfg *config.Config) grpc.SubscriptionClient {
    if cfg.GRPCMock || cfg.SubscriptionBillingAddr == "" {
        return &mocks.FakeSubsClient{}
    }
    return grpc.NewRealSubscriptionClient(cfg.SubscriptionBillingAddr)
}
```

#### 8.2 Cache Layer

**📄 Archivos:**
- `src/services/subscription_cache.go`
- Redis integration via `rem-common/cache`

**🔧 Cache Strategy:**
- [ ] **Key pattern**: `subs:{org_id}`
- [ ] **TTL**: 5 minutos
- [ ] **Invalidation**: Manual endpoint para testing
- [ ] **Fallback**: Si Redis falla, hit directo al client

**🔧 Implementation:**

```go
type SubscriptionCache struct {
    client SubscriptionClient
    redis  *redis.Client
}

func (c *SubscriptionCache) GetOrgSubscription(orgID uuid.UUID) (*SubscriptionSummary, error) {
    // 1. Try cache
    key := fmt.Sprintf("subs:%s", orgID)
    cached := c.redis.Get(key)
    if cached != nil {
        return unmarshal(cached)
    }
    
    // 2. Fetch from client
    summary := c.client.GetOrgSubscription(context.Background(), orgID)
    
    // 3. Cache for 5 min
    c.redis.Set(key, marshal(summary), 5*time.Minute)
    
    return summary
}
```

#### 8.3 Limits Enforcement

**📄 Archivos:**
- `src/middleware/limits_guard.go`
- Feature flag en config: `SUB_CHECK=true/false`

**🔧 Middleware:**

```go
func LimitsGuard(limits []string) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            if !config.SubCheckEnabled {
                return next(c) // Skip en desarrollo
            }
            
            orgID := c.Get("org_id").(uuid.UUID)
            summary := subscriptionService.GetOrgSubscription(orgID)
            
            for _, limit := range limits {
                if limit == "employees" && summary.EmployeesUsed >= summary.EmployeesLimit {
                    return echo.NewHTTPError(402, "Employee limit reached")
                }
                // Otros límites...
            }
            
            return next(c)
        }
    }
}

// Uso en endpoints
e.POST("/organizations/:id/employees", employeeController.Create, LimitsGuard([]string{"employees"}))
```

#### 8.4 Read-Only Endpoints

**🔧 Endpoints:**

- [ ] `GET /organizations/:id/subscription` - Estado de suscripción
  - [ ] Response: plan_code, status, limits vs usage
  - [ ] Auth: Rol Manager+
  - [ ] Cache: 5 min TTL

- [ ] `GET /organizations/:id/subscription/limits` - Solo límites
  - [ ] Response: employees_limit, listings_limit, storage_limit
  - [ ] Auth: Cualquier miembro
  - [ ] Cache: 10 min TTL

- [ ] `GET /organizations/:id/subscription/usage` - Solo uso actual
  - [ ] Response: employees_used, listings_used (via property-svc mock)
  - [ ] Auth: Rol Manager+
  - [ ] Cache: 1 min TTL

- [ ] `DELETE /organizations/:id/subscription/cache` - Invalidar cache
  - [ ] Para testing y troubleshooting
  - [ ] Auth: Rol Admin

---

### **FASE 9: Observabilidad y Monitoreo** ⏱️ (2-3 días)

#### 9.1 Métricas (Prometheus)

**📄 Archivos:**
- `src/metrics/metrics.go`
- `src/controllers/metrics_controller.go`

**🔧 Métricas a exponer:**

```go
// src/metrics/metrics.go
var (
    OrganizationsTotal = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "organizations_total",
            Help: "Total number of organizations by status",
        },
        []string{"status"},
    )
    
    EmployeesTotal = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "employees_total", 
            Help: "Total number of employees by organization",
        },
        []string{"organization_id", "status"},
    )
    
    InvitesTotal = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "invites_total",
            Help: "Total number of invites by status",
        },
        []string{"status"},
    )
    
    HTTPRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint", "status_code"},
    )
)
```

**🔧 Endpoints:**
- [ ] `GET /metrics` - Prometheus scraping endpoint
- [ ] Middleware para medir latencia de requests
- [ ] Background worker para actualizar gauges cada 5 min

#### 9.2 Eventos de Negocio (RabbitMQ)

**📄 Archivos:**
- `src/events/publisher.go`
- `src/events/types.go`

**🔧 Event Types:**

```go
type EventType string

const (
    OrgCreated           EventType = "org.created"
    OrgUpdated           EventType = "org.updated"
    OrgDeleted           EventType = "org.deleted"
    EmployeeHired        EventType = "employee.hired"
    EmployeeFired        EventType = "employee.fired"
    EmployeeRoleChanged  EventType = "employee.role_changed"
    InviteSent           EventType = "invite.sent"
    InviteAccepted       EventType = "invite.accepted"
    IntegrationConnected EventType = "integration.connected"
    DomainVerified       EventType = "domain.verified"
)

type OrganizationEvent struct {
    EventID       uuid.UUID   `json:"event_id"`
    EventType     EventType   `json:"event_type"`
    OrganizationID uuid.UUID  `json:"organization_id"`
    UserID        *uuid.UUID  `json:"user_id,omitempty"`
    Timestamp     time.Time   `json:"timestamp"`
    Payload       interface{} `json:"payload"`
    Source        string      `json:"source"` // "organization-svc"
}
```

**🔧 Publisher Integration:**
- [ ] Usar `rem-common/broker` RabbitMQ
- [ ] Exchange: `org.events`
- [ ] Routing keys: `org.created`, `employee.hired`, etc.
- [ ] Async publish (no bloquear requests)
- [ ] Retry logic con backoff

#### 9.3 Structured Logging

**🔧 Log Standards:**
- [ ] Usar `rem-common/logger` (Zap)
- [ ] Structured JSON en producción
- [ ] Context propagation (request_id, user_id, org_id)
- [ ] Log levels apropiados

```go
// Ejemplo de logging estructurado
logger.Info("Organization created",
    zap.String("org_id", org.ID.String()),
    zap.String("created_by", userID.String()),
    zap.String("display_name", org.DisplayName),
    zap.Int64("employee_limit", subscription.EmployeesLimit),
)
```

---

### **FASE 10: Testing y Documentación** ⏱️ (3-4 días)

#### 10.1 Testing Strategy

**📄 Estructura de tests:**
```
organization-svc/
├── tests/
│   ├── unit/
│   │   ├── services/
│   │   ├── repository/
│   │   └── utils/
│   ├── integration/
│   │   ├── controllers/
│   │   └── grpc/
│   └── e2e/
│       └── scenarios/
```

**🔧 Unit Tests:**
- [ ] **Services**: Mock repositories, test business logic
- [ ] **Repositories**: Test DB interactions con base in-memory
- [ ] **Utils**: Token generation, encryption, validations
- [ ] **Coverage target**: 80%+

**🔧 Integration Tests:**
- [ ] **Controllers**: Test endpoints con database real
- [ ] **gRPC Clients**: Test mocks vs interfaces
- [ ] **Middleware**: Auth, limits, logging

**🔧 E2E Tests:**
- [ ] **Complete flows**: Crear org → invitar → aceptar → asignar roles
- [ ] **Error scenarios**: Límites excedidos, tokens expirados
- [ ] **Multi-org scenarios**: Usuario en múltiples organizaciones

#### 10.2 Colección Postman Completa

**📄 Archivos:**
- `Organization-Service.postman_collection.json`
- `Organization-Service.postman_environment.json`

**🔧 Collections Structure:**
```
📁 Organization Service
├── 📁 Health & Auth
│   ├── Health Check
│   └── Login (inherit from main collection)
├── 📁 Organizations
│   ├── Create Organization
│   ├── Get Organization by ID
│   ├── List Organizations
│   ├── Update Organization  
│   └── Delete Organization
├── 📁 Settings
│   ├── Get All Settings
│   ├── Update Setting
│   └── Delete Setting
├── 📁 Owners
│   ├── Add Owner (existing person)
│   ├── Add Owner (create person)
│   ├── List Owners
│   └── Remove Owner
├── 📁 Branches
│   ├── Create Branch
│   ├── List Branches
│   ├── Update Branch
│   ├── Set Main Branch
│   └── Delete Branch
├── 📁 Roles
│   ├── Create Custom Role
│   ├── List Roles
│   ├── Update Role
│   └── Delete Role
├── 📁 Employees
│   ├── Add Employee
│   ├── List Employees
│   ├── Update Employee
│   ├── Fire Employee
│   ├── Assign Role
│   ├── Remove Role
│   └── Change Primary Role
├── 📁 Invitations
│   ├── Create Invite
│   ├── List Invites
│   ├── Cancel Invite
│   ├── View Invite (public)
│   ├── Accept Invite (existing user)
│   ├── Accept Invite (new user)
│   └── Decline Invite
├── 📁 Integrations
│   ├── Connect ZonaProp
│   ├── Connect WhatsApp
│   ├── List Integrations
│   ├── Test Integration
│   └── Disconnect Integration
├── 📁 Domains
│   ├── Add Domain
│   ├── List Domains
│   ├── Verify Domain
│   ├── Enable SSL
│   └── Remove Domain
├── 📁 Subscriptions (Read-Only)
│   ├── Get Subscription
│   ├── Get Limits
│   ├── Get Usage
│   └── Clear Cache
└── 📁 Error Scenarios
    ├── Duplicate Organization Name
    ├── Invalid Token
    ├── Exceeded Limits
    └── Permission Denied
```

**🔧 Auto-Variables:**
- [ ] `{{org_id}}` - Auto-set al crear organización
- [ ] `{{employee_id}}` - Auto-set al agregar empleado
- [ ] `{{invite_token}}` - Auto-set al crear invitación
- [ ] `{{branch_id}}` - Auto-set al crear sucursal
- [ ] `{{role_id}}` - Auto-set al crear rol

#### 10.3 Documentación Completa

**📄 README.md del servicio:**
- [ ] **Overview** y arquitectura
- [ ] **Quick start** con docker-compose
- [ ] **API endpoints** documentados
- [ ] **gRPC integrations** con ejemplos
- [ ] **Environment variables** completas
- [ ] **Testing** instructions
- [ ] **Troubleshooting** common issues

**📄 API Documentation:**
- [ ] **OpenAPI 3.0 spec** generado desde controllers
- [ ] **Request/Response examples** para cada endpoint
- [ ] **Error codes** y meanings
- [ ] **Authentication** requirements

**📄 Architecture Diagrams:**
- [ ] **Service dependencies** (gRPC calls)
- [ ] **Database schema** con relaciones
- [ ] **Event flow** para business events
- [ ] **Authentication/Authorization** flow

#### 10.4 Scripts de Desarrollo

**📄 Integration en project scripts:**

- [ ] **docker-compose.yml** update
```yaml
# En dev-full-compose.yml
organization-svc:
  build: ./organization-svc
  ports:
    - "8083:8083"
    - "50053:50053"
  environment:
    - DB_DSN=postgres://user:pass@postgres:5432/rem_db
    - GRPC_MOCK=true
    - REDIS_URL=redis://redis:6379
  depends_on:
    - postgres
    - redis
  volumes:
    - ./organization-svc:/app
```

- [ ] **dev.bat** update - Include organization-svc hot-reload
- [ ] **seed-db.bat** update - Seed organizations test data
- [ ] **test-gateway.bat** update - Include org endpoints testing

**📄 Seed Data:**
```sql
-- En seed-dev-db.sql
INSERT INTO organization (id, display_name, status, created_at) VALUES
('00000000-0000-0000-0000-000000000001', 'Inmobiliaria Central', 'active', now()),
('00000000-0000-0000-0000-000000000002', 'PropHouse Real Estate', 'active', now()),
('00000000-0000-0000-0000-000000000003', 'Casa & Terreno SA', 'suspended', now());

-- Seed roles para cada organización
SELECT seed_default_roles('00000000-0000-0000-0000-000000000001');
SELECT seed_default_roles('00000000-0000-0000-0000-000000000002');
SELECT seed_default_roles('00000000-0000-0000-0000-000000000003');
```

---

## 📊 Definición Detallada de .proto

### subscription-billing/v1/subscription.proto

```protobuf
syntax = "proto3";

package subscription.v1;
option go_package = "github.com/rem-gestion/rem-common/protos/subscription/v1;subscriptionpb";

import "google/protobuf/timestamp.proto";

service SubscriptionService {
  rpc GetOrgSubscription(GetOrgSubscriptionRequest) returns (SubscriptionSummary);
  rpc GetUsageLimits(GetUsageLimitsRequest) returns (UsageLimits);
}

message GetOrgSubscriptionRequest {
  string organization_id = 1;
}

message GetUsageLimitsRequest {
  string organization_id = 1;
}

message SubscriptionSummary {
  string organization_id = 1;
  string plan_code = 2;           // "PRO", "TEAM", "ENTERPRISE"
  string status = 3;              // "active", "trial", "canceled"
  int64 listings_limit = 4;
  int64 listings_used = 5;
  int64 employees_limit = 6;
  int64 employees_used = 7;
  google.protobuf.Timestamp starts_at = 8;
  google.protobuf.Timestamp renews_at = 9;
  
  // Reserved para expansiones futuras
  reserved 10, 11, 12;
}

message UsageLimits {
  string organization_id = 1;
  int64 listings_limit = 2;
  int64 employees_limit = 3;
  int64 storage_limit_mb = 4;
  int64 api_calls_limit = 5;
  
  // Reserved para nuevos límites
  reserved 6, 7, 8;
}
```

### property/v1/property_counter.proto

```protobuf
syntax = "proto3";

package property.v1;
option go_package = "github.com/rem-gestion/rem-common/protos/property/v1;propertypb";

service PropertyService {
  rpc CountProperties(CountPropertiesRequest) returns (CountPropertiesResponse);
}

message CountPropertiesRequest {
  string organization_id = 1;
  string branch_id = 2;        // Optional - filter by branch
  string status_filter = 3;    // Optional - "active", "sold", "rented"
}

message CountPropertiesResponse {
  int64 total = 1;
  repeated BranchCount by_branch = 2;
  repeated StatusCount by_status = 3;
  
  // Reserved para métricas adicionales
  reserved 4, 5, 6;
}

message BranchCount {
  string branch_id = 1;
  string branch_name = 2;
  int64 count = 3;
}

message StatusCount {
  string status = 1;
  int64 count = 2;
}
```

### role-rbac/v1/rbac.proto

```protobuf
syntax = "proto3";

package rbac.v1;
option go_package = "github.com/rem-gestion/rem-common/protos/rbac/v1;rbacpb";

service RbacService {
  rpc ValidatePermission(ValidatePermissionRequest) returns (ValidatePermissionResponse);
  rpc GetUserPermissions(GetUserPermissionsRequest) returns (GetUserPermissionsResponse);
}

message ValidatePermissionRequest {
  string user_id = 1;
  string organization_id = 2;
  string permission = 3;      // "org.create", "employee.invite", etc.
  string resource_id = 4;     // Optional - specific resource
}

message ValidatePermissionResponse {
  bool allowed = 1;
  string reason = 2;          // Optional - explanation if denied
}

message GetUserPermissionsRequest {
  string user_id = 1;
  string organization_id = 2;
}

message GetUserPermissionsResponse {
  repeated string permissions = 1;
  repeated string roles = 2;
  
  // Reserved para context adicional
  reserved 3, 4, 5;
}

enum Permission {
  PERMISSION_UNSPECIFIED = 0;
  ORG_CREATE = 1;
  ORG_UPDATE = 2;
  ORG_DELETE = 3;
  EMPLOYEE_INVITE = 4;
  EMPLOYEE_FIRE = 5;
  EMPLOYEE_ROLE_ASSIGN = 6;
  BRANCH_CREATE = 7;
  INTEGRATION_MANAGE = 8;
  DOMAIN_MANAGE = 9;
  
  // Reserved para nuevos permisos
  reserved 10 to 50;
}
```

---

## 📊 Modelo de Datos Detallado

### Organization
```json
{
  "id": "uuid",
  "display_name": "string",
  "logo_url": "string",
  "fiscal_address_id": "uuid  // service: address -> address.id",
  "matricula": "string",
  "status": "string  // active|suspended|deleted",
  "created_at": "timestamp",
  "created_by": "uuid  // service: auth_identity -> users.id",
  "updated_at": "timestamp",
  "updated_by": "uuid  // auth_identity",
  "deleted_at": "timestamp"
}
```

### OrganizationBranch
```json
{
  "id": "uuid",
  "organization_id": "uuid  // local FK",
  "display_name": "string",
  "address_id": "uuid  // address-svc",
  "phone": "string",
  "email": "string",
  "is_main": "bool",
  "created_at": "timestamp",
  "created_by": "uuid  // auth_identity",
  "updated_at": "timestamp",
  "updated_by": "uuid",
  "deleted_at": "timestamp"
}
```

### OrganizationSettings
```json
{
  "organization_id": "uuid",
  "key": "string",
  "value": "json",
  "updated_at": "timestamp",
  "updated_by": "uuid  // auth_identity"
}
```

### OrganizationOwner
```json
{
  "id": "uuid",
  "organization_id": "uuid",
  "owner_type": "string  // individual|company",
  "owner_id": "uuid  // service: person -> person.id",
  "created_at": "timestamp",
  "created_by": "uuid  // auth_identity",
  "deleted_at": "timestamp"
}
```

### OrganizationRole
```json
{
  "id": "int",
  "organization_id": "uuid",
  "name": "string",
  "description": "string",
  "created_at": "timestamp",
  "created_by": "uuid"
}
```

### Employee
```json
{
  "id": "uuid",
  "organization_id": "uuid",
  "user_id": "uuid  // auth_identity -> users.id",
  "primary_role_id": "int  // local FK -> organization_role.id",
  "hired_at": "date",
  "fired_at": "date",
  "created_at": "timestamp",
  "created_by": "uuid"
}
```

### EmployeeRole
```json
{
  "employee_id": "uuid  // local FK -> employees.id",
  "role_id": "int  // local FK -> organization_role.id",
  "primary": "bool",
  "created_at": "timestamp",
  "created_by": "uuid"
}
```

### OrganizationInvite
```json
{
  "id": "uuid",
  "organization_id": "uuid",
  "role_id": "int  // local FK -> organization_role.id",
  "email": "string",
  "status": "string  // pending|accepted|declined|expired",
  "token": "string",
  "expires_at": "timestamp",
  "accepted_at": "timestamp",
  "created_at": "timestamp",
  "created_by": "uuid"
}
```

### OrganizationIntegration
```json
{
  "id": "uuid",
  "organization_id": "uuid",
  "provider": "string",
  "external_id": "string",
  "config": "json",
  "status": "string  // connected|error|revoked",
  "connected_at": "timestamp",
  "created_at": "timestamp",
  "created_by": "uuid"
}
```

### OrganizationDomain
```json
{
  "id": "uuid",
  "organization_id": "uuid",
  "domain": "string",
  "ssl_enabled": "bool",
  "verified_at": "timestamp",
  "created_at": "timestamp",
  "created_by": "uuid"
}
```

### OrganizationSubscription (Read-Only)
```json
{
  "id": "uuid",
  "organization_id": "uuid",
  "plan_id": "int  // service: subscription_billing -> plan.id",
  "status": "string  // active|trial|canceled",
  "starts_at": "date",
  "renews_at": "date",
  "current_usage": "json",
  "limit_snapshot": "json",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

---

## 🔄 Flujos de Negocio Críticos Detallados

### 1. Creación de Nueva Inmobiliaria
```
┌─────────────────────────────────────────────────────────────────┐
│ POST /organizations                                             │
│ Body: {display_name, fiscal_address_id, owner_data, ...}        │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 1. Validar datos de entrada                                    │
│    • display_name único (soft-delete aware)                    │
│    • fiscal_address_id válido (address-svc gRPC)              │
│    • owner_data válido si se provee                           │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 2. Verificar límites de suscripción (mock)                     │
│    • subscription-svc.GetOrgSubscription(user.current_org?)    │
│    • Validar no excede organizations_limit                     │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 3. Crear/validar owner si se provee                           │
│    • Si owner_id: validar existe (person-svc gRPC)            │
│    • Si owner_data: crear persona (person-svc gRPC)           │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 4. Transacción DB                                             │
│    • INSERT organization                                       │
│    • CALL seed_default_roles(org_id)                          │
│    • INSERT organization_owner (si aplica)                     │
│    • INSERT employee (creator como Admin)                      │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 5. Crear sucursal principal (opcional)                        │
│    • Si se provee branch_data                                 │
│    • Crear address (address-svc) + branch con is_main=true    │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 6. Eventos y respuesta                                        │
│    • Publish "org.created" event (RabbitMQ)                   │
│    • Update metrics (organizations_total)                      │
│    • Return organization data + roles + owner info            │
└─────────────────────────────────────────────────────────────────┘
```

### 2. Flujo Completo de Invitación → Aceptación → Empleado
```
┌─────────────────────────────────────────────────────────────────┐
│ POST /organizations/:id/invites                                │
│ Body: {email, role_id}                                         │
│ Auth: Rol Manager+ en la organización                          │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 1. Validaciones pre-creación                                   │
│    • Email no duplicado en invites pending de la org          │
│    • Email no pertenece a empleado actual de la org           │
│    • role_id válido y pertenece a la organización             │
│    • Límite employees_limit no alcanzado (subs-svc mock)      │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 2. Generar y almacenar invitación                             │
│    • token = SHA-256(random_bytes) con prefijo "inv_"         │
│    • expires_at = now() + 7 days                              │
│    • status = "pending"                                        │
│    • INSERT organization_invite                                │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 3. Notificación (placeholder)                                 │
│    • TODO: Send email con link de aceptación                  │
│    • Publish "invite.sent" event                              │
│    • Return invite data (incluye token para testing)          │
└─────────────────────────────────────────────────────────────────┘

         ... Tiempo pasa, usuario recibe link ...

┌─────────────────────────────────────────────────────────────────┐
│ GET /invites/:token                                            │
│ (Público - sin auth)                                           │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 1. Validar token y mostrar detalles                           │
│    • Token existe y no expirado                               │
│    • status = "pending"                                        │
│    • Return: org_name, role_name, inviter_name, expires_at    │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ PUT /invites/:token/accept                                     │
│ Body: {account_data} o JWT existente                           │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 2. Procesar aceptación                                        │
│    ┌─ Si no tiene account ─────────────────────────────────┐   │
│    │ • CreateAccount(email, password) → auth-identity-svc │   │
│    │ • CreateUser(account_id) → auth-identity-svc         │   │
│    └─────────────────────────────────────────────────────┘   │
│    ┌─ Si ya tiene account ────────────────────────────────┐   │
│    │ • Extraer user_id del JWT                           │   │
│    │ • Validar user existe                               │   │
│    └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 3. Crear empleado y actualizar invitación                     │
│    • INSERT employee (user_id, org_id, primary_role_id)       │
│    • INSERT employee_roles (employee_id, role_id, primary=t)  │
│    • UPDATE invite SET status='accepted', accepted_at=now()   │
│    • Publish "invite.accepted" + "employee.hired" events      │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 4. Respuesta exitosa                                          │
│    • Generate JWT con org context                             │
│    • Return: token, user_data, organization_data, role_data   │
└─────────────────────────────────────────────────────────────────┘
```

### 3. Integración con Proveedor Externo (ZonaProp)
```
┌─────────────────────────────────────────────────────────────────┐
│ POST /organizations/:id/integrations                           │
│ Body: {provider: "zonaprop", config: {api_key, account_id}}    │
│ Auth: Rol Admin en la organización                             │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 1. Validar credenciales con ZonaProp                          │
│    • HTTP GET a ZonaProp API con api_key                      │
│    • Verificar response 200 y account_id match                │
│    • Si falla: return error "Invalid credentials"             │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 2. Encriptar y almacenar configuración                        │
│    • config_encrypted = AES-256-GCM(config_json)             │
│    • INSERT organization_integration                           │
│    • status = "connected", connected_at = now()               │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 3. Sincronización inicial (placeholder)                       │
│    • TODO: Trigger initial sync job                           │
│    • TODO: Import existing properties from ZonaProp          │
│    • Publish "integration.connected" event                    │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 4. Respuesta exitosa                                          │
│    • Return integration data (sin config por seguridad)       │
│    • Include connection status y última sync                  │
└─────────────────────────────────────────────────────────────────┘
```

### 4. Enforcement de Límites de Suscripción
```
┌─────────────────────────────────────────────────────────────────┐
│ Cualquier endpoint que consume límites                         │
│ Ej: POST /organizations/:id/employees                          │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ Middleware: LimitsGuard(["employees"])                         │
│ (Solo si SUB_CHECK=true)                                       │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 1. Obtener suscripción (con cache)                            │
│    • Cache hit: return cached data                            │
│    • Cache miss: subscription-svc.GetOrgSubscription()        │
│    • Si mock: return unlimited plan                           │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 2. Calcular uso actual                                        │
│    • employees_used = COUNT(*) FROM employees WHERE ...       │
│    • listings_used = property-svc.CountProperties() (mock=0)  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 3. Validar límites                                            │
│    • IF employees_used >= employees_limit THEN                │
│    •   RETURN 402 "Employee limit reached"                    │
│    • IF listings_used >= listings_limit THEN                 │
│    •   RETURN 402 "Listing limit reached"                     │
└─────────────────────────────────────────────────────────────────┘
                              │ (Límites OK)
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 4. Continuar con request normal                               │
│    • next(c) - Ejecutar handler del endpoint                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🎯 Criterios de Aceptación Detallados por Fase

### ✅ Fase 0 - Contratos
- [ ] **Archivos .proto**: 3 servicios con métodos mínimos definidos
- [ ] **Generación de código**: `buf generate` ejecuta sin errores  
- [ ] **Mock clients**: Implementados en `rem-common/testing/mocks/`
- [ ] **Build tags**: Mocks vs real clients controlados por env
- [ ] **Permissions matrix**: CSV con 15+ permisos documentados

### ✅ Fase 1 - Infraestructura
- [ ] **Servicio HTTP**: Responde en puerto 8083 con Echo v4
- [ ] **Health check**: `GET /health` retorna status ok
- [ ] **Conexión DB**: PostgreSQL conectado y migraciones automáticas
- [ ] **Docker integration**: Service inicia en docker-compose
- [ ] **Hot-reload**: Air configurado y funcional
- [ ] **Environment**: Variables de entorno cargadas correctamente

### ✅ Fase 2 - Datos y Migraciones
- [ ] **Scripts SQL**: 7 migraciones up/down ejecutan correctamente
- [ ] **Constraints**: Unique, FK, check constraints funcionan
- [ ] **Soft delete**: deleted_at implementado en todas las tablas
- [ ] **Audit triggers**: updated_at se actualiza automáticamente
- [ ] **Indexes**: Performance indexes creados para queries comunes
- [ ] **Seed data**: Función seed_default_roles() funciona

### ✅ Fase 3 - Entidades Core
- [ ] **Organization CRUD**: 5 endpoints funcionan con validaciones
- [ ] **Settings KV**: Sistema clave-valor con upsert/delete
- [ ] **Owners management**: Integración person-svc + CRUD completo
- [ ] **Validaciones**: Business rules aplicadas (nombre único, etc.)
- [ ] **Error handling**: Responses consistentes con códigos HTTP apropiados
- [ ] **Postman tests**: Endpoints básicos testeados y documentados

### ✅ Fase 4 - Estructura Organizacional  
- [ ] **Branches**: CRUD + constraint sucursal principal única
- [ ] **Roles**: Sistema roles personalizable + predeterminados
- [ ] **Address integration**: Crear/validar direcciones vía address-svc
- [ ] **Business rules**: No eliminar sucursal principal, roles default protegidos
- [ ] **Relaciones**: FK relationships funcionan correctamente

### ✅ Fase 5 - Empleados
- [ ] **Employee CRUD**: Gestión completa con estados (activo/despedido)
- [ ] **Role assignment**: N-a-N roles con primary único por empleado
- [ ] **Auth integration**: Validación usuarios vía auth-identity-svc
- [ ] **Limits enforcement**: Mock subscription limits aplicados
- [ ] **Employee middleware**: Context org_id + employee_id seteable
- [ ] **Constraint validation**: user_id único por organización

### ✅ Fase 6 - Invitaciones
- [ ] **Token system**: Generación SHA-256 + expiración funcional
- [ ] **Public endpoints**: Ver/aceptar/rechazar sin autenticación
- [ ] **Account creation**: Flujo completo new user vía auth-identity-svc
- [ ] **State management**: Estados pending/accepted/declined/expired
- [ ] **Integration flow**: Invite → Accept → Employee creation seamless
- [ ] **Email placeholder**: Structure lista para email integration

### ✅ Fase 7 - Integraciones y Dominios
- [ ] **Provider support**: ZonaProp, WhatsApp, MercadoLibre configurables
- [ ] **Config encryption**: AES-256-GCM protege credenciales sensibles
- [ ] **Connection testing**: Validation pings a provider APIs
- [ ] **Domain verification**: 3 métodos (DNS, HTTP, Meta) implementados
- [ ] **SSL management**: SSL toggle post-verificación
- [ ] **Unique constraints**: Un provider por org, domain único global

### ✅ Fase 8 - Subscriptions (RO)
- [ ] **Mock client**: subscription-svc mock retorna data realista
- [ ] **Cache layer**: Redis TTL 5min + fallback funcional
- [ ] **Limits middleware**: 402 errors cuando se exceden límites
- [ ] **Feature flags**: SUB_CHECK=false disables validation
- [ ] **RO endpoints**: 4 endpoints subscription data + cache clear
- [ ] **Performance**: Cache hit/miss metrics visible

### ✅ Fase 9 - Observabilidad
- [ ] **Prometheus metrics**: /metrics endpoint con 4+ metrics
- [ ] **RabbitMQ events**: 8+ event types published a org.events exchange
- [ ] **Structured logging**: Zap JSON logs con context propagation
- [ ] **Background workers**: Metrics update cada 5min
- [ ] **Request tracing**: HTTP duration histograms por endpoint

### ✅ Fase 10 - Testing y Docs
- [ ] **Unit tests**: 80%+ coverage en services y utils
- [ ] **Integration tests**: Controllers con real DB
- [ ] **E2E scenarios**: 3+ complete workflows testados
- [ ] **Postman collection**: 50+ requests con auto-variables
- [ ] **Documentation**: README + API docs + architecture diagrams
- [ ] **Seed scripts**: Test data generation para development

---

## 🚀 Comandos de Desarrollo Actualizados

### Estructura de Scripts Final

```bash
# Development workflow
.\scripts\dev.bat                    # Inicia todo (incluye organization-svc)
.\scripts\dev-org.bat               # Solo organization-svc + dependencies
.\scripts\clean.bat                 # Limpia todo
.\scripts\seed-db.bat               # Pobla todas las tablas + organizations

# Testing específico
.\scripts\test-organization.bat     # Test solo endpoints de organization
.\scripts\test-postman-org.bat      # Run Postman collection específica

# Database específico
.\scripts\migrate-org.bat           # Solo migraciones de organization-svc
.\scripts\rollback-org.bat          # Rollback una migración
```

### URLs y Puertos Final

| Servicio | Puerto HTTP | Puerto gRPC | Health Check |
|----------|-------------|-------------|--------------|
| **API Gateway** | 8081 | - | http://localhost:8081/health |
| **auth-identity-svc** | 4002 | 50052 | http://localhost:4002/health |
| **person-svc** | 4001 | 50051 | http://localhost:4001/health |
| **address-svc** | 4000 | 50050 | http://localhost:4000/health |
| **organization-svc** | 8083 | 50053 | http://localhost:8083/health |

### Comandos de Testing Específicos

```bash
# Endpoints via Gateway (producción-like)
curl http://localhost:8081/api/organizations/health

# Endpoint directo (debug only)
curl http://localhost:8083/health

# Postman collection específica
newman run Organization-Service.postman_collection.json \
  -e Organization-Service.postman_environment.json

# Test de integración específico
cd organization-svc && go test ./tests/integration/...

# Test de mocks
GRPC_MOCK=true go test ./src/services/...
```

---

## 📈 Métricas de Éxito Detalladas

### 📊 **Cobertura y Calidad**
- **⏱️ Tiempo total estimado**: 6-8 semanas (vs 4-6 original por mayor detalle)
- **🔄 Cobertura de endpoints**: 40+ endpoints (vs 20+ estimado original)
- **🧪 Testing cobertura**: 80%+ unit tests + integration tests completos
- **📚 Documentación**: README + OpenAPI spec + diagramas + Postman
- **🔗 Integraciones**: 6 servicios integrados (4 reales + 2 mocks)
- **⚡ Performance**: < 200ms endpoints básicos, < 500ms con gRPC calls

### 📈 **Funcionalidades Implementadas**
- **✅ 11 entidades CRUD** completas con business rules
- **✅ 1 entidad read-only** (subscriptions) con cache
- **✅ 3 servicios mock** con interfaces futuras
- **✅ Sistema completo de auth** multi-organización  
- **✅ Event-driven architecture** con RabbitMQ
- **✅ Observabilidad production-ready** con metrics y logging

### 🎯 **Criterios de "Production Ready"**
- [ ] **Docker production builds** < 100MB
- [ ] **Zero-downtime migrations** strategy documented
- [ ] **Error handling** comprehensive con circuit breakers
- [ ] **Rate limiting** per organization implemented
- [ ] **API versioning** strategy defined (v1, v2 paths)
- [ ] **Security audit** completed (encryption, tokens, permissions)

---

## 🎯 Próximos Pasos Detallados

### **Semana 1-2: Foundation**
1. **✅ Revisar arquitectura** existente en otros servicios
2. **📝 Definir .proto files** para mocks (subscription, property, rbac)  
3. **🏗️ Setup infraestructura**: Docker, migrations, basic HTTP server
4. **🗃️ Implementar core entities**: Organization, Settings, Owners

### **Semana 3-4: Structure & People**
5. **🏢 Branches y Roles**: Estructura organizacional completa
6. **👥 Employee management**: CRUD + role assignment + auth integration
7. **📧 Invitation system**: Token generation → acceptance → employee creation

### **Semana 5-6: Integrations & Subscriptions**
8. **🔌 External integrations**: ZonaProp, WhatsApp setup con encryption
9. **🌐 Domain management**: Verification + SSL management
10. **💰 Subscription limits**: Mock enforcement + cache layer

### **Semana 7-8: Polish & Production**
11. **📊 Observability**: Metrics, events, structured logging
12. **🧪 Testing completeness**: Unit + integration + E2E scenarios
13. **📚 Documentation**: README + API docs + Postman collections
14. **🚀 Production preparation**: Docker builds + deployment strategy

### **Post-MVP: Real Integrations**
- **🔄 Replace mocks** con servicios reales cuando estén disponibles
- **📧 Email notifications** para invitaciones
- **🔄 Background jobs** para sync con integraciones externas
- **📊 Advanced analytics** y reporting
- **🛡️ Advanced RBAC** con fine-grained permissions

¡Listo para comenzar el desarrollo con una hoja de ruta completamente detallada! 🚀

**Cada fase tiene deliverables claros, criterios de aceptación específicos y está diseñada para ser completamente funcional al final.**
