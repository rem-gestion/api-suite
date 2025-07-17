# 🏢 Organization Service - Estrategia de Desarrollo MVP

## 📋 Análisis Exhaustivo del Sistema Actual

### 🎯 Estado Actual de los Servicios

#### ✅ **Servicios Completados (Base Sólida)**
- **auth-identity-svc**: Usuarios, autenticación JWT, onboarding
- **person-svc**: Gestión de personas (individuos/empresas) con contactos
- **address-svc**: Direcciones centralizadas con resiliencia automática
- **rem-common**: Librería compartida con utilidades comunes

#### 🔄 **Servicios en Desarrollo**
- **organization-svc**: Gestión de inmobiliarias (este documento)
- **subscription-billing-svc**: Estructura básica creada

#### 📋 **Pendientes para MVP Completo**
- **property-svc**: Gestión de propiedades inmobiliarias

---

## 🏗️ Análisis del Esquema de Base de Datos

### 📊 **Entidades Principales Identificadas**

#### **Núcleo Organizacional (Esencial MVP)**
1. **`organization`** - Entidad principal de la inmobiliaria
2. **`organization_settings`** - Configuraciones KV flexibles
3. **`organization_owner`** - Propietarios con porcentajes
4. **`organization_branch`** - Sucursales/oficinas
5. **`organization_role`** - Roles internos (Admin, Manager, Agent, Assistant)

#### **Gestión de Personal (Esencial MVP)**
6. **`employees`** - Vínculo usuario ↔ organización
7. **`employee_roles`** - Roles múltiples por empleado

#### **Sistema de Invitaciones (MVP - Funcionalidad Crítica)**
8. **`organization_invite`** - Invitaciones con tokens seguros

#### **Funcionalidades Avanzadas (Post-MVP)**
9. **`organization_integration`** - Hub de integraciones externas
10. **`organization_domain`** - Dominios personalizados
11. **`organization_subscription`** - Planes y facturación (preparado)

---

## 🎯 Estrategia de Desarrollo MVP

### **Fase 1: MVP Core (2-3 semanas)**

#### **Endpoints Esenciales (28 de 74 endpoints)**

##### 🏢 **Organization Management**
```http
# CRUD Básico de Organizaciones
POST   /organizations                     # Crear inmobiliaria
GET    /organizations                     # Listar mis organizaciones
GET    /organizations/{orgId}             # Obtener organización específica
PUT    /organizations/{orgId}             # Actualizar organización
PATCH  /organizations/{orgId}/status      # Cambiar estado

# Settings Básicos
GET    /organizations/{orgId}/settings    # Obtener configuraciones
PUT    /organizations/{orgId}/settings    # Actualizar configuraciones
```

##### 👥 **Employee Management**
```http
# Gestión de Empleados
POST   /organizations/{orgId}/employees           # Vincular empleado
GET    /organizations/{orgId}/employees           # Listar empleados
GET    /organizations/{orgId}/employees/{empId}   # Detalle empleado
PUT    /organizations/{orgId}/employees/{empId}   # Actualizar empleado
PATCH  /organizations/{orgId}/employees/{empId}/status # Cambiar estado empleado

# Gestión de Roles de Empleados
GET    /organizations/{orgId}/employees/{empId}/roles        # Roles del empleado
POST   /organizations/{orgId}/employees/{empId}/roles/{roleId} # Asignar rol
DELETE /organizations/{orgId}/employees/{empId}/roles/{roleId} # Remover rol
```

##### 🏘️ **Branch Management**
```http
# Gestión de Sucursales
POST   /organizations/{orgId}/branches           # Crear sucursal
GET    /organizations/{orgId}/branches           # Listar sucursales
GET    /organizations/{orgId}/branches/{branchId} # Detalle sucursal
PUT    /organizations/{orgId}/branches/{branchId} # Actualizar sucursal
DELETE /organizations/{orgId}/branches/{branchId} # Eliminar sucursal
```

##### 🔐 **Role Management**
```http
# Gestión de Roles
POST   /organizations/{orgId}/roles      # Crear rol personalizado
GET    /organizations/{orgId}/roles      # Listar roles
GET    /organizations/{orgId}/roles/{roleId} # Detalle rol
PUT    /organizations/{orgId}/roles/{roleId} # Actualizar rol
DELETE /organizations/{orgId}/roles/{roleId} # Eliminar rol
```

##### 📧 **Invitation System**
```http
# Sistema de Invitaciones
POST   /organizations/{orgId}/invites           # Enviar invitación
GET    /organizations/{orgId}/invites           # Listar invitaciones
GET    /organizations/{orgId}/invites/{invId}   # Detalle invitación
POST   /organizations/{orgId}/invites/{invId}/resend # Reenviar
POST   /organizations/{orgId}/invites/{invId}/cancel # Cancelar

# Endpoints Públicos (sin auth)
POST   /invites/accept    # Aceptar invitación con token
GET    /invites/validate  # Validar token de invitación
POST   /invites/decline   # Rechazar invitación
```

##### 👤 **Owner Management**
```http
# Gestión de Propietarios
GET    /organizations/{orgId}/owners           # Listar propietarios
POST   /organizations/{orgId}/owners           # Añadir propietario
GET    /organizations/{orgId}/owners/{ownerId} # Detalle propietario
DELETE /organizations/{orgId}/owners/{ownerId} # Remover propietario
```

---

### **Fase 2: Post-MVP (Futuro - Integraciones)**

#### **Endpoints de Integraciones (26 endpoints)**
- Todo el sistema de `organization_integration`
- Logs de integración y eventos
- Sincronizaciones con CRM, contabilidad, etc.

#### **Endpoints de Dominios (14 endpoints)**
- Sistema completo de `organization_domain`
- Verificación DNS/SSL
- Dominios personalizados

#### **Endpoints de Mantenimiento (6 endpoints)**
- Logs de auditoría
- Mantenimiento programado
- Métricas avanzadas

---

## 🔄 Comunicación Entre Servicios

### **Necesidades Identificadas de Comunicación gRPC**

#### **Con auth-identity-svc**
```go
// Validar que user_id existe al crear empleado
service AuthService {
    rpc ValidateUser(ValidateUserRequest) returns (ValidateUserResponse);
    rpc GetUserDetails(GetUserRequest) returns (UserDetailsResponse);
}
```

#### **Con person-svc**
```go
// Validar y obtener datos de personas (propietarios, empleados)
service PersonService {
    rpc ValidatePerson(ValidatePersonRequest) returns (ValidatePersonResponse);
    rpc GetPersonDetails(GetPersonRequest) returns (PersonDetailsResponse);
    rpc GetPersonsByIds(GetPersonsByIdsRequest) returns (GetPersonsByIdsResponse);
}
```

#### **Con address-svc**
```go
// Validar direcciones fiscales y de sucursales
service AddressService {
    rpc ValidateAddress(ValidateAddressRequest) returns (ValidateAddressResponse);
    rpc GetAddressDetails(GetAddressRequest) returns (AddressDetailsResponse);
}
```

---

## 🚀 Arquitectura de Controllers

### **Estructura de Controllers MVP**

```go
organization-svc/src/controllers/
├── organization_controller.go      // CRUD básico de organizaciones
├── employee_controller.go          // Gestión de empleados
├── branch_controller.go            // Gestión de sucursales
├── role_controller.go              // Gestión de roles
├── invitation_controller.go        // Sistema de invitaciones
├── owner_controller.go             // Gestión de propietarios
└── health_controller.go            // Health checks
```

### **Servicios de Negocio**

```go
organization-svc/src/services/
├── organization_service.go         // Lógica de organizaciones
├── employee_service.go             // Lógica de empleados + gRPC calls
├── invitation_service.go           // Lógica de invitaciones + email queue
├── grpc_client_service.go          // Cliente gRPC unificado
└── validation_service.go           // Validaciones complejas
```

---

## 📦 Cache y Colas (Redis + RabbitMQ)

### **🔴 Redis Cache - Uso Estratégico**

#### **Casos de Uso Identificados**
```yaml
Cache Strategy:
  organization_settings:
    key: "org:settings:{org_id}"
    ttl: 1800  # 30 minutos
    reason: "Configuraciones accedidas frecuentemente"
    
  organization_roles:
    key: "org:roles:{org_id}"
    ttl: 3600  # 1 hora
    reason: "Roles usados en cada request de autorización"
    
  employee_roles:
    key: "emp:roles:{employee_id}"
    ttl: 900   # 15 minutos
    reason: "Validación de permisos en cada endpoint"
    
  organization_basic:
    key: "org:basic:{org_id}"
    ttl: 1800  # 30 minutos
    fields: "id, display_name, status, is_active"
    reason: "Datos básicos para validaciones"
```

#### **Implementación en rem-common**
```go
// Ya disponible en rem-common/cache/redis.go
cache := remcommon.NewRedisCache(config.Redis)

// Patrón de uso en services
func (s *OrganizationService) GetSettings(orgID string) (*models.Settings, error) {
    // 1. Verificar cache
    if cached, err := s.cache.Get(fmt.Sprintf("org:settings:%s", orgID)); err == nil {
        return parseSettings(cached), nil
    }
    
    // 2. Consultar DB
    settings, err := s.repo.GetSettings(orgID)
    if err != nil {
        return nil, err
    }
    
    // 3. Guardar en cache
    s.cache.Set(fmt.Sprintf("org:settings:%s", orgID), settings, 30*time.Minute)
    return settings, nil
}
```

### **🐰 RabbitMQ - Event-Driven Architecture**

#### **Eventos Identificados**
```yaml
Organization Events:
  organization.created:
    payload: { org_id, created_by, display_name }
    consumers: 
      - subscription-billing-svc  # Crear cuenta de facturación
      - notification-svc          # Email de bienvenida
      
  employee.added:
    payload: { org_id, employee_id, user_id, roles[] }
    consumers:
      - auth-identity-svc         # Actualizar permisos
      - notification-svc          # Email de incorporación
      
  invitation.sent:
    payload: { org_id, invite_id, email, role_id }
    consumers:
      - notification-svc          # Enviar email de invitación
      
  invitation.accepted:
    payload: { org_id, invite_id, new_employee_id }
    consumers:
      - auth-identity-svc         # Crear/vincular usuario
      - notification-svc          # Confirmación de aceptación
```

#### **Implementación con rem-common**
```go
// Ya disponible en rem-common/broker/rabbit.go
broker := remcommon.NewRabbitBroker(config.RabbitMQ)

// Publisher en organization service
func (s *OrganizationService) CreateOrganization(data *dto.CreateOrgRequest) error {
    org, err := s.repo.Create(data)
    if err != nil {
        return err
    }
    
    // Publicar evento
    event := events.OrganizationCreated{
        OrgID:       org.ID,
        CreatedBy:   data.CreatedBy,
        DisplayName: org.DisplayName,
        Timestamp:   time.Now(),
    }
    
    return s.broker.Publish("organization.created", event)
}
```

---

## 📋 Plan de Desarrollo Detallado

### **Sprint 1 (Semana 1): Fundamentos**
```yaml
Día 1-2: Setup & Infrastructure
  - Configurar estructura de proyecto
  - Setup gRPC clients en rem-common
  - Configurar Redis y RabbitMQ connections
  
Día 3-5: Core Models & Repository
  - Implementar modelos base (Organization, Employee, Branch, Role)
  - Repository layer con GORM
  - Pruebas unitarias de modelos
```

### **Sprint 2 (Semana 2): Core Business Logic**
```yaml
Día 1-3: Organization & Employee Services
  - organization_service.go con gRPC calls
  - employee_service.go con validaciones
  - Cache layer para settings y roles
  
Día 4-5: Controllers & Routes
  - Controllers básicos con validaciones
  - Middleware de autorización
  - Pruebas de integración
```

### **Sprint 3 (Semana 3): Invitation System**
```yaml
Día 1-3: Invitation Workflow
  - invitation_service.go completo
  - Token generation y validación
  - Email queue integration
  
Día 4-5: Testing & Polish
  - Pruebas end-to-end
  - Documentación de APIs
  - Performance testing
```

---

## 🎯 Criterios de Éxito MVP

### **Funcionalidades Críticas**
- ✅ Crear inmobiliaria con sucursal principal
- ✅ Invitar empleados con roles específicos
- ✅ Gestionar estructura organizacional básica
- ✅ Comunicación confiable con otros servicios
- ✅ Cache eficiente para operaciones frecuentes

### **Métricas de Performance**
```yaml
Targets:
  - Tiempo de respuesta < 200ms (90% requests)
  - Cache hit rate > 80% para settings/roles
  - Email delivery rate > 95% para invitaciones
  - Uptime > 99.5%
```

### **Preparación para Subscription-Billing**
- Vista `organization_subscription_details` lista
- Eventos de facturación configurados
- Campos preparados para límites de plan

### **Preparación para Property Service**
- Estructura de empleados lista para asignación de propiedades
- Roles definidos para gestión de propiedades
- Sucursales listas para asociar propiedades

---

## 🚀 Próximos Pasos Inmediatos

### **1. Configuración Inicial (2-3 días)**
```bash
# Estructura básica de proyecto
organization-svc/
├── cmd/api/main.go
├── src/
│   ├── controllers/
│   ├── services/
│   ├── models/
│   ├── dto/
│   └── grpc/
└── go.mod
```

### **2. Desarrollo Core (Semanas 1-2)**
- Implementar 28 endpoints MVP identificados
- Configurar comunicación gRPC
- Integrar cache y eventos

### **3. Testing y Refinamiento (Semana 3)**
- Pruebas completas del sistema
- Optimización de performance
- Preparación para property-svc

---

## 💡 Consideraciones Técnicas Clave

### **Consistencia de Datos**
- Transacciones para operaciones críticas (crear org + sucursal + admin)
- Validaciones referidas via gRPC antes de commits
- Rollback automático en caso de fallos

### **Seguridad**
- Autorización granular por organización
- Validación de propietarios/administradores
- Tokens seguros para invitaciones

### **Escalabilidad**
- Cache strategy desde el inicio
- Event-driven para integraciones futuras
- Preparado para sharding por organización

---

## 🎯 Conclusión

El MVP de Organization Service será la **piedra angular** del sistema inmobiliario. Con los 28 endpoints identificados, proporcionará:

1. **Base sólida** para property-svc y subscription-billing-svc
2. **Gestión completa** de inmobiliarias y empleados
3. **Sistema robusto** de invitaciones y onboarding
4. **Arquitectura escalable** para funcionalidades futuras

**Tiempo estimado**: 3 semanas para MVP completo
**Prioridad**: **Crítica** - Bloqueante para property-svc
**Complejidad**: **Media-Alta** - Requiere coordinación entre múltiples servicios

---

*Este README será actualizado durante el desarrollo con detalles específicos de implementación y decisiones arquitectónicas.*
