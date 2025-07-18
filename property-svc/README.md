# 🏠 Property Service - REM Platform

Microservicio para la gestión de propiedades inmobiliarias, amenities y management de propiedades dentro de la plataforma REM.

## 📋 Funcionalidades

### 🏢 Gestión de Propiedades
- ✅ **CRUD completo** de propiedades inmobiliarias
- ✅ **Tipos de propiedad**: Departamento, Casa, Local comercial, Oficina, Terreno, Nave industrial
- ✅ **Filtros avanzados**: por tipo, propietario, búsqueda por texto
- ✅ **Paginación** optimizada
- ✅ **Códigos internos** únicos para organización

### 🎯 Gestión de Amenities
- ✅ **CRUD completo** de amenities (servicios/comodidades)
- ✅ **Categorización** por tipo (seguridad, recreación, servicios)
- ✅ **Asociación** flexible entre propiedades y amenities
- ✅ **Búsqueda** por nombre y categoría

### 🏢 Property Management
- ✅ **Gestión de administración** de propiedades por organizaciones
- ✅ **Períodos de gestión** con fechas de inicio y fin
- ✅ **Comisiones** por gestión
- ✅ **Estados activos/inactivos**

## 🗄️ Modelo de Datos

### Principales Entidades

| Entidad | Descripción | Campos Clave |
|---------|-------------|--------------|
| **Property** | Propiedad inmobiliaria | `owner_person_id`, `address_id`, `property_type`, `internal_code` |
| **PropertyManagement** | Gestión de propiedades | `property_id`, `organization_id`, `commission_percent` |
| **Amenity** | Servicios/Comodidades | `name`, `category`, `icon` |
| **PropertyAmenity** | Relación N:N | `property_id`, `amenity_id` |

### Tipos de Propiedad
- `APARTMENT` - Departamento
- `HOUSE` - Casa
- `COMMERCIAL_SPACE` - Local comercial
- `OFFICE` - Oficina
- `LAND` - Terreno
- `INDUSTRIAL_WAREHOUSE` - Nave industrial

## 🚀 API Endpoints

### Properties
```http
POST   /properties              # Crear propiedad
GET    /properties              # Listar propiedades (con filtros)
GET    /properties/:id          # Obtener propiedad específica
PUT    /properties/:id          # Actualizar propiedad
DELETE /properties/:id          # Eliminar propiedad
```

### Amenities
```http
POST   /amenities               # Crear amenity
GET    /amenities               # Listar amenities
GET    /amenities/:id           # Obtener amenity específico
PUT    /amenities/:id           # Actualizar amenity
DELETE /amenities/:id           # Eliminar amenity
```

### Property-Amenity Management
```http
POST   /properties/:property_id/amenities/:amenity_id    # Agregar amenity a propiedad
DELETE /properties/:property_id/amenities/:amenity_id    # Quitar amenity de propiedad
```

### Health Check
```http
GET    /health                  # Estado del servicio
```

## 🔧 Configuración

### Variables de Entorno (.env.development)
```bash
REM_SERVICE_NAME=property
REM_ENVIRONMENT=development
REM_PROPERTY_HTTP_PORT=4004
REM_GRPC_PORT=50054
```

### Puerto Asignado
- **HTTP**: `4004` (desarrollo)
- **gRPC**: `50054` (futuro uso)

## 📊 Testing

### Datos Dummy/Mock
El servicio incluye IDs mock para testing:

```json
{
  "mock_owner_person_ids": [
    "11111111-1111-1111-1111-111111111111",
    "22222222-2222-2222-2222-222222222222"
  ],
  "mock_address_ids": [
    "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
    "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
  ],
  "mock_organization_ids": [
    "99999999-9999-9999-9999-999999999999"
  ]
}
```

### Ejemplo de Payload - Crear Propiedad
```json
{
  "owner_person_id": "11111111-1111-1111-1111-111111111111",
  "address_id": "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
  "property_type": "APARTMENT",
  "internal_code": "DEPT-001",
  "year_built": 2020,
  "bedrooms": 2,
  "bathrooms": 1.5,
  "total_area_sqm": 85.50,
  "covered_area_sqm": 75.00,
  "description": "Departamento moderno en zona céntrica",
  "management": {
    "organization_id": "99999999-9999-9999-9999-999999999999",
    "commission_percent": 3.5
  }
}
```

## 🏗️ Arquitectura Futura

### Comunicación gRPC
- **Address Service**: Validación y obtención de direcciones
- **Person Service**: Validación de propietarios
- **Organization Service**: Validación de organizaciones gestoras

### Proto Files (Preparación)
```proto
service PropertyService {
  rpc GetProperty(GetPropertyRequest) returns (PropertyResponse);
  rpc ListProperties(ListPropertiesRequest) returns (ListPropertiesResponse);
  rpc CreateProperty(CreatePropertyRequest) returns (PropertyResponse);
}
```

## 📝 Comandos de Desarrollo

```bash
# Iniciar con hot-reload
air

# Ejecutar migraciones
go run cmd/migrate/main.go

# Build manual
go build -o tmp/main.exe cmd/api/main.go
```

## 🔄 Estado del Desarrollo

✅ **COMPLETADO**:
- Migraciones de base de datos
- Modelos y DTOs
- Repository con CRUD completo
- Service con lógica de negocio
- Controllers con validaciones
- Routers configurados
- Configuración de entorno

🚧 **PENDIENTE**:
- Integración gRPC con otros servicios
- Validación real de IDs externos
- Tests unitarios
- Documentación OpenAPI/Swagger
