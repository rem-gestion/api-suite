# Repository Layer - Organization Service

## 📋 Overview

Esta implementación del **Repository Layer** para el servicio de organizaciones sigue el patrón Repository con interfaces bien definidas, proporcionando una abstracción limpia para operaciones de datos utilizando GORM y PostgreSQL.

## 🏗️ Estructura Implementada

### ✅ Repositorios Completados

1. **Organization Repository** (`organization_repository.go`)
   - CRUD completo para organizaciones
   - Búsqueda y filtrado avanzado
   - Gestión de estado y verificación
   - Relaciones con entidades asociadas
   - Soporte para paginación

2. **Organization Setting Repository** (`setting_repository.go`)
   - Gestión de configuraciones clave-valor
   - Operaciones bulk para múltiples settings
   - Configuraciones por defecto
   - Reset a valores predeterminados

3. **Employee Repository** (`employee_repository.go`)
   - CRUD para empleados
   - Filtrado por organización, rol y sucursal
   - Gestión de estados de empleados
   - Terminación de empleados con auditoría

4. **Employee Role Repository** (`employee_role_repository.go`)
   - Asignación y remoción de roles
   - Gestión de roles primarios
   - Consultas por empleado o rol
   - Validaciones de negocio

5. **Branch Repository** (`branch_repository.go`)
   - CRUD completo para sucursales
   - Gestión de sucursal principal
   - Filtrado por organización
   - Validaciones de negocio (no eliminar única sucursal)

6. **Role Repository** (`role_repository.go`)
   - CRUD para roles organizacionales
   - Creación automática de roles por defecto
   - Gestión de roles personalizados
   - Estadísticas de uso de roles

7. **Owner Repository** (`owner_repository.go`)
   - Gestión de propietarios y porcentajes
   - Validación de porcentajes de propiedad
   - Transferencia de propiedad
   - Identificación de fundadores y mayoristas

8. **Invitation Repository** (`invitation_repository.go`)
   - Sistema completo de invitaciones
   - Gestión de tokens seguros
   - Estados de invitación (pending, accepted, rejected, etc.)
   - Expiración automática de invitaciones

9. **Invitation Log Repository** (`invitation_log_repository.go`)
   - Auditoría completa de invitaciones
   - Métodos de conveniencia para logging
   - Limpieza automática de logs antiguos
   - Seguimiento de actores y acciones

## 🔧 Características Técnicas

### ✅ Implementadas

- **GORM Integration**: Uso completo de GORM para operaciones ORM
- **Error Handling**: Manejo consistente de errores con mensajes descriptivos
- **Soft Delete**: Soporte para eliminación suave con auditoría
- **Transactions**: Operaciones transaccionales donde sea necesario
- **Context Support**: Soporte completo para contexto de Go
- **UUID Primary Keys**: Uso de UUIDs como claves primarias
- **Audit Fields**: Campos de auditoría (created_at, updated_at, created_by, etc.)
- **Pagination**: Soporte para paginación en listados
- **Filtering**: Sistema de filtros flexible y extensible
- **Relationships**: Carga de relaciones con Preload
- **Business Logic**: Métodos de soporte para lógica de negocio

### 🎯 Interfaces

Todas las interfaces están definidas en `interfaces.go` con:
- Separación clara de responsabilidades
- Métodos bien documentados
- Soporte para todas las operaciones CRUD
- Métodos especializados para cada entidad
- Validaciones de negocio integradas

### 📊 Patrones Implementados

1. **Repository Pattern**: Abstracción de la capa de datos
2. **Dependency Injection**: Interfaces para facilitar testing
3. **Factory Pattern**: Constructores para cada repository
4. **Manager Pattern**: Agregador de repositories (`repository.go`)

## 🛠️ Uso y Configuración

### Inicialización

```go
import (
    "github.com/rem-gestion/api-suite/organization/src/repository"
    "gorm.io/gorm"
)

// Crear una instancia del manager de repositories
repoManager := repository.NewRepositoryManager(db)

// Obtener repositories específicos
orgRepo := repoManager.GetOrganizationRepository()
empRepo := repoManager.GetEmployeeRepository()
```

### Ejemplo de Uso

```go
// Crear una organización
org := &models.Organization{
    Name: "Tech Company",
    Slug: "tech-company",
    Type: models.OrgTypeRealEstate,
    Status: models.OrgStatusActive,
}

createdOrg, err := orgRepo.Create(ctx, org)
if err != nil {
    return fmt.Errorf("failed to create organization: %w", err)
}

// Buscar organizaciones con filtros
filters := dto.OrganizationFiltersRequest{
    FilterRequest: dto.FilterRequest{
        Search: "Tech",
        Status: "active",
    },
    Type: "real_estate",
    PaginationRequest: dto.PaginationRequest{
        Page: 1,
        PerPage: 10,
    },
}

orgs, total, err := orgRepo.List(ctx, filters)
```

## 🔒 Separación de Responsabilidades

### ✅ Incluidas en este Servicio
- **Organizations**: Entidad principal
- **Employees**: Empleados vinculados a organizaciones
- **Branches**: Sucursales organizacionales
- **Roles**: Roles internos de la organización
- **Owners**: Propietarios y porcentajes de participación
- **Invitations**: Sistema de invitaciones a empleados
- **Settings**: Configuraciones organizacionales

### ❌ Excluidas (Servicios Externos)
- **Subscriptions**: Manejado por `subscription-billing-svc`
- **Person Data**: Manejado por `person-svc` 
- **Address Data**: Manejado por `address-svc`
- **User Authentication**: Manejado por `auth-identity-svc`

## 📈 Próximos Pasos

### Task 1.2 - Completar Repositories Restantes

1. **Branch Repository**
   - Implementar CRUD para sucursales
   - Gestión de sucursal principal
   - Validaciones de negocio

2. **Role Repository** 
   - Roles organizacionales personalizados
   - Roles por defecto del sistema
   - Gestión de permisos

3. **Owner Repository**
   - Gestión de propietarios
   - Cálculo de porcentajes
   - Validaciones de propiedad

4. **Invitation System**
   - Repository de invitaciones
   - Audit log completo
   - Gestión de tokens y expiración

### Task 1.3 - Service Layer
Implementar la capa de servicios que utilize estos repositories.

### Task 1.4 - Controller Layer  
Implementar los controladores REST que expongan la funcionalidad.

## 🚀 Estado Actual

**✅ Task 1.1**: DTOs completamente implementados
**🔄 Task 1.2**: Repository Layer 50% completado
- ✅ Interfaces definidas (100%)
- ✅ Core repositories implementados (4/9 - 44%)
- 🚧 Repositories restantes pendientes (5/9 - 56%)

La base sólida está establecida y los patterns están implementados. Los repositories restantes seguirán el mismo patrón establecido.
