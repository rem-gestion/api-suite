# 🏢 Property Service - REM Platform

Microservicio completo para la gestión integral de propiedades inmobiliarias, tipos de propiedad, amenities, gestión de propiedades y tipos de manager dentro de la plataforma REM.

## 📋 Funcionalidades Completas

### 🏢 Gestión de Propiedades
- ✅ **CRUD completo** de propiedades inmobiliarias
- ✅ **Integración gRPC** con Address Service y Person Service
- ✅ **Validación automática** de owners con Person Service
- ✅ **Manejo inteligente** de direcciones con Address Service
- ✅ **Códigos internos únicos** autogenerados
- ✅ **Filtros avanzados**: por tipo, propietario, precio, características
- ✅ **Paginación optimizada** con metadatos
- ✅ **Búsqueda avanzada** por múltiples criterios
- ✅ **Detalles completos** con amenities expandidas

### 🏠 Gestión de Tipos de Propiedad
- ✅ **CRUD completo** de property types
- ✅ **Gestión de estados** activo/inactivo
- ✅ **Validación de uso** antes de eliminar
- ✅ **Filtros por estado** y búsqueda

### 🎯 Gestión de Amenities
- ✅ **CRUD completo** de amenities con categorías
- ✅ **Categorización avanzada**: seguridad, recreación, servicios, bienestar
- ✅ **Asociación flexible** entre propiedades y amenities
- ✅ **Filtros por categoría** y búsqueda por nombre
- ✅ **URLs de iconos** para interfaz visual

### 👔 Property Management
- ✅ **Gestión completa** de property managements
- ✅ **Asociación con manager types** y organizaciones
- ✅ **Períodos de gestión** con fechas y comisiones
- ✅ **Estados activos/inactivos** de gestiones
- ✅ **Validación de datos** y reglas de negocio

### 🎖️ Manager Types
- ✅ **Gestión de tipos** de administradores
- ✅ **Creación y listado** de manager types
- ✅ **Integración** con property managements

## 🌐 Integración con Servicios

### Address Service (gRPC)
- **Validación automática**: Verifica direcciones antes de crear propiedades
- **Creación inteligente**: Crea direcciones automáticamente si no existen  
- **Reutilización eficiente**: Usa direcciones existentes por ID
- **Fallback resiliente**: Mock services para desarrollo independiente

### Person Service (gRPC)
- **Validación de propietarios**: Verifica que el owner_person_id exista
- **Información expandida**: Obtiene datos del propietario para respuestas completas
- **Fallback resiliente**: Mock services para desarrollo independiente

## 🗄️ Modelo de Datos Completo

### Entidades Principales

| Entidad | Descripción | Campos Clave |
|---------|-------------|--------------|
| **Property** | Propiedades inmobiliarias | `owner_person_id`, `address_id`, `property_type_id`, `internal_code`, características físicas |
| **PropertyType** | Tipos de propiedad | `name`, `description`, `is_active` |
| **PropertyManagement** | Gestión de propiedades | `property_id`, `manager_type_id`, `organization_id`, `start_date`, `commission_percent` |
| **ManagerType** | Tipos de administradores | `name`, `description` |
| **Amenity** | Servicios/Comodidades | `name`, `category`, `icon_url` |
| **PropertyAmenity** | Relación N:N | `property_id`, `amenity_id` |

### Categorías de Amenities
- `seguridad` - Características de seguridad
- `recreacion` - Instalaciones recreativas
- `servicios` - Servicios y utilidades  
- `bienestar` - Instalaciones de bienestar

## 🚀 API Endpoints Completos

### 🏢 Properties Management
```http
GET    /api/properties                     # Listar con filtros avanzados
POST   /api/properties                     # Crear propiedad
GET    /api/properties/:id                 # Obtener propiedad específica
PUT    /api/properties/:id                 # Actualizar propiedad
DELETE /api/properties/:id                 # Eliminar propiedad
GET    /api/properties/:id/full-details    # Detalles completos con amenities
```

**Filtros avanzados disponibles:**
- `page`, `limit` - Paginación
- `owner_person_id` - Filtrar por propietario
- `property_type_id` - Filtrar por tipo de propiedad
- `min_price`, `max_price` - Rango de precios
- `bedrooms` - Número de habitaciones
- `bathrooms` - Número de baños
- `min_area`, `max_area` - Rango de área
- `search` - Búsqueda en descripción

### 🏠 Property Types
```http
GET    /api/property-types                 # Listar tipos de propiedad
POST   /api/property-types                 # Crear nuevo tipo
GET    /api/property-types/:id             # Obtener tipo específico
PUT    /api/property-types/:id             # Actualizar tipo
DELETE /api/property-types/:id             # Eliminar tipo
```

**Filtros disponibles:**
- `page`, `limit` - Paginación
- `is_active` - Filtrar por estado activo/inactivo

### 👔 Property Management
```http
GET    /api/property-managements           # Listar gestiones
POST   /api/property-managements           # Crear nueva gestión
GET    /api/property-managements/:id       # Obtener gestión específica
PUT    /api/property-managements/:id       # Actualizar gestión
DELETE /api/property-managements/:id       # Eliminar gestión
```

**Filtros disponibles:**
- `page`, `limit` - Paginación
- `property_id` - Filtrar por propiedad
- `manager_type_id` - Filtrar por tipo de manager
- `organization_id` - Filtrar por organización
- `is_active` - Filtrar por estado activo/inactivo

### 🎖️ Manager Types
```http
GET    /api/manager-types                  # Listar tipos de manager
POST   /api/manager-types                  # Crear nuevo tipo de manager
```

### 🎯 Amenities
```http
GET    /api/amenities                      # Listar amenities con filtros
POST   /api/amenities                      # Crear amenity
GET    /api/amenities/:id                  # Obtener amenity específica
PUT    /api/amenities/:id                  # Actualizar amenity
DELETE /api/amenities/:id                  # Eliminar amenity
```

**Filtros disponibles:**
- `page`, `limit` - Paginación
- `category` - Filtrar por categoría (seguridad, recreacion, servicios, bienestar)

### 🔗 Property-Amenity Relations
```http
GET    /api/properties/:id/amenities       # Amenities de una propiedad
POST   /api/properties/:id/amenities       # Agregar amenity a propiedad  
DELETE /api/properties/:id/amenities/:aid  # Quitar amenity de propiedad
```

### 🩺 Health Check
```http
GET    /health                             # Estado del servicio
```

## 📤 Ejemplos de Request/Response

### Crear Propiedad Completa
```json
POST /api/properties
{
  "owner_person_id": "11111111-1111-1111-1111-111111111111",
  "address_id": "22222222-2222-2222-2222-222222222222",
  "property_type_id": 1,
  "internal_code": "DEPT-2024-001",
  "year_built": 2023,
  "bedrooms": 3,
  "bathrooms": 2.5,
  "total_area_sqm": 95.75,
  "covered_area_sqm": 85.50,
  "description": "Departamento moderno de 3 ambientes con vista al parque"
}
```

### Respuesta con Detalles Completos
```json
GET /api/properties/:id/full-details
{
  "id": "prop-uuid-123",
  "owner_person_id": "11111111-1111-1111-1111-111111111111",
  "address_id": "22222222-2222-2222-2222-222222222222",
  "property_type_id": 1,
  "internal_code": "DEPT-2024-001",
  "year_built": 2023,
  "bedrooms": 3,
  "bathrooms": 2.5,
  "total_area_sqm": 95.75,
  "covered_area_sqm": 85.50,
  "description": "Departamento moderno de 3 ambientes con vista al parque",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z",
  "amenities": [
    {
      "id": 1,
      "name": "Piscina",
      "category": "recreacion",
      "icon_url": "https://example.com/icons/pool.svg"
    },
    {
      "id": 2,
      "name": "Seguridad 24hs",
      "category": "seguridad", 
      "icon_url": "https://example.com/icons/security.svg"
    }
  ]
}
```

### Crear Amenity
```json
POST /api/amenities
{
  "name": "Gimnasio Premium",
  "category": "bienestar",
  "icon_url": "https://example.com/icons/gym.svg"
}
```

### Crear Property Management
```json
POST /api/property-managements
{
  "property_id": "prop-uuid-123",
  "manager_type_id": 1,
  "organization_id": "org-uuid-456",
  "start_date": "2024-01-01T00:00:00Z",
  "end_date": "2024-12-31T23:59:59Z",
  "commission_percent": 5.5,
  "is_active": true
}
```

## 🔧 Configuración

### Variables de Entorno (.env.development)
```bash
# Service Configuration
REM_SERVICE_NAME=property
REM_ENVIRONMENT=development
REM_PROPERTY_HTTP_PORT=4004
REM_PROPERTY_GRPC_PORT=50054

# Database Configuration
REM_POSTGRES_HOST=localhost
REM_POSTGRES_PORT=5432
REM_POSTGRES_USER=rem_user
REM_POSTGRES_PASSWORD=rem_password
REM_POSTGRES_DATABASE=rem_development

# External Services
REM_ADDRESS_GRPC_HOST=0.0.0.0
REM_ADDRESS_GRPC_PORT=50052
REM_PERSON_GRPC_HOST=0.0.0.0
REM_PERSON_GRPC_PORT=50051

# Logging
REM_LOG_LEVEL=debug
```

### Puertos Asignados
- **HTTP**: `4004` (desarrollo)
- **gRPC**: `50054` (comunicación entre servicios)

### Dependencias de Servicios
- **Address Service**: Puerto `50052` (gRPC) - Para manejo de direcciones
- **Person Service**: Puerto `50051` (gRPC) - Para validación de propietarios
- **Database**: PostgreSQL `rem_development` - Base de datos compartida

## 🚀 Ejecutar el Servicio

### Desarrollo
```bash
# Desde el directorio property-svc
cd property-svc

# Instalar dependencias
go mod tidy

# Ejecutar migraciones
go run ./cmd/migrate

# Iniciar con hot-reload (requiere Air)
air

# O ejecutar directamente
go run ./cmd/api/main.go
```

### Construcción
```bash
# Compilar
go build -o api.exe ./cmd/api/main.go

# Ejecutar
./api.exe
```

## 🧪 Testing

### Health Check
```bash
curl http://localhost:4004/health
```

### Postman Collection
El servicio está completamente documentado en la colección de Postman incluida:
- `REM-API-Collection.postman_collection.json`
- `REM-Development.postman_environment.json`

La colección incluye:
- ✅ **Todos los endpoints** con datos pre-llenados
- ✅ **Scripts automáticos** para guardar IDs
- ✅ **Variables de entorno** configuradas
- ✅ **Ejemplos realistas** listos para ejecutar

### Test de Integración
```bash
# Desde el directorio raíz api-suite
.\scripts\test-postman-features.bat
```

## 🏗️ Arquitectura del Servicio

### Estructura de Directorios
```
property-svc/
├── cmd/
│   ├── api/main.go              # Entry point HTTP server
│   └── migrate/main.go          # Database migrations
├── migrations/pg/               # SQL migrations
├── src/
│   ├── controllers/             # HTTP request handlers
│   ├── services/               # Business logic
│   ├── repositories/           # Data access layer
│   ├── models/                 # Data models
│   ├── grpc/                   # gRPC client integrations
│   └── utils/                  # Utilities and helpers
├── tmp/                        # Air hot-reload artifacts
├── tests/                      # Test files
├── go.mod                      # Go module
├── go.sum                      # Go dependencies
├── .air.toml                   # Air configuration
└── README.md                   # This file
```

### Patrones Implementados
- ✅ **Repository Pattern**: Abstracción de acceso a datos
- ✅ **Service Layer**: Lógica de negocio centralizada
- ✅ **Controller Pattern**: Manejo de requests HTTP
- ✅ **Dependency Injection**: Inyección de dependencias
- ✅ **Circuit Breaker**: Para servicios externos
- ✅ **Graceful Shutdown**: Cierre ordenado del servicio

## 📊 Base de Datos

### Tablas Principales
- `properties` - Propiedades inmobiliarias
- `property_types` - Tipos de propiedad
- `property_managements` - Gestión de propiedades
- `manager_types` - Tipos de administradores
- `amenities` - Servicios y comodidades
- `property_amenities` - Relación N:N propiedades-amenities

### Migraciones
Las migraciones se ejecutan automáticamente al iniciar el servicio o pueden ejecutarse manualmente:
```bash
go run ./cmd/migrate
```

## 🔄 Estado del Desarrollo

### ✅ Completado
- [x] CRUD completo de Properties
- [x] CRUD completo de Property Types
- [x] CRUD completo de Amenities
- [x] CRUD completo de Property Managements
- [x] CRUD de Manager Types
- [x] Relaciones Property-Amenity
- [x] Integración gRPC con Address Service
- [x] Integración gRPC con Person Service
- [x] Filtros avanzados y paginación
- [x] Validaciones y manejo de errores
- [x] Documentación completa
- [x] Colección Postman actualizada
- [x] Health checks
- [x] Migraciones de base de datos

### 🚧 Pendiente (Futuras Iteraciones)
- [ ] Implementación gRPC server
- [ ] Tests unitarios e integración
- [ ] Métricas y monitoring
- [ ] Cache con Redis
- [ ] Eventos con RabbitMQ
- [ ] Subida de imágenes
- [ ] Geolocalización avanzada

## 🤝 Contribución

Para contribuir al desarrollo:

1. Revisa los endpoints en la colección Postman
2. Verifica que todos los tests pasen
3. Actualiza la documentación si es necesario
4. Asegúrate de que las migraciones estén correctas

## 📞 Soporte

Para dudas o problemas:
- Revisa la colección Postman para ejemplos de uso
- Verifica los logs del servicio en desarrollo
- Consulta la documentación de servicios dependientes
