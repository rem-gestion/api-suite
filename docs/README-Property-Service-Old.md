# 🏢 Property Service - Documentación Completa

El **Property Service** es un microservicio completo para la gestión de propiedades inmobiliarias, amenities y sus relaciones en la plataforma REM.

## 🚀 Estado Actual

✅ **COMPLETAMENTE IMPLEMENTADO Y FUNCIONAL**

- ✅ Migraciones de base de datos
- ✅ Modelos y DTOs completos
- ✅ Repository pattern implementado
- ✅ Servicios con validación externa
- ✅ Controladores REST completos
- ✅ Rutas configuradas y funcionando
- ✅ Integración con API Gateway
- ✅ Postman Collection actualizada
- ✅ Documentación del proyecto actualizada

## 🏗️ Arquitectura

### Tecnologías
- **Framework**: Go 1.23.0 con Gin
- **Base de Datos**: PostgreSQL con GORM
- **Puerto**: 4004
- **Comunicación**: gRPC con address-svc y person-svc
- **Autenticación**: API Key middleware

### Entidades Principales

| Entidad | Descripción | Campos Clave |
|---------|-------------|--------------|
| **Property** | Propiedad inmobiliaria | `owner_person_id`, `address_id`, `property_type`, `internal_code` |
| **PropertyAmenity** | Relación N:N | `property_id`, `amenity_id` |
| **Amenity** | Servicios/Comodidades | `name`, `category`, `icon` |

### Tipos de Propiedad
- `APARTMENT` (Apartamento)
- `HOUSE` (Casa)
- `COMMERCIAL_SPACE` (Local Comercial)
- `OFFICE` (Oficina)
- `LAND` (Terreno)
- `INDUSTRIAL_WAREHOUSE` (Galpón Industrial)

### Categorías de Amenities
- `Security` (Seguridad)
- `Recreation` (Recreación)
- `Services` (Servicios)
- `Transport` (Transporte)
- `Healthcare` (Salud)
- `Education` (Educación)

## 📡 API Endpoints

### 🏢 Propiedades

```bash
# Listar propiedades con paginación
GET /api/properties/?page=1&per_page=10

# Obtener propiedad específica
GET /api/properties/{id}

# Crear nueva propiedad
POST /api/properties/
{
  "owner_person_id": "uuid",
  "address_id": "uuid", 
  "property_type": "APARTMENT",
  "title": "Departamento 2 ambientes",
  "description": "Excelente ubicación"
}

# Actualizar propiedad
PUT /api/properties/{id}

# Eliminar propiedad
DELETE /api/properties/{id}
```

### 🎯 Amenities

```bash
# Listar amenities
GET /api/amenities/?page=1&per_page=10

# Crear amenity
POST /api/amenities/
{
  "name": "Pileta",
  "category": "Recreation",
  "icon": "🏊‍♂️"
}

# Obtener amenity específico
GET /api/amenities/{id}

# Actualizar amenity
PUT /api/amenities/{id}

# Eliminar amenity
DELETE /api/amenities/{id}
```

### 🔗 Relaciones Property-Amenity

```bash
# Asociar amenity a propiedad
POST /api/properties/manage/{property_id}/amenities/{amenity_id}

# Desasociar amenity de propiedad
DELETE /api/properties/manage/{property_id}/amenities/{amenity_id}
```

## 🔧 Configuración

### Variables de Entorno

```bash
# Servidor
REM_PROPERTY_HTTP_PORT=4004

# Base de datos (compartida)
REM_POSTGRES_HOST=localhost
REM_POSTGRES_PORT=5432
REM_POSTGRES_USER=user
REM_POSTGRES_PASSWORD=supersecreta
REM_POSTGRES_DBNAME=rem_development

# Seguridad
REM_API_KEY=supersecret-api-key-for-dev

# Servicios externos (gRPC)
REM_ADDRESS_HOST=localhost
REM_ADDRESS_PORT=50052
REM_PERSON_HOST=localhost
REM_PERSON_PORT=50051
```

### Servicios Mock

Durante el desarrollo, el servicio utiliza **servicios mock** para address-svc y person-svc:

```go
// Mock data para development
var mockOwnerPersonIDs = []string{
    "11111111-1111-1111-1111-111111111111",
    "22222222-2222-2222-2222-222222222222",
    // ...
}

var mockAddressIDs = []string{
    "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
    "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
    // ...
}
```

## 🚀 Uso en Desarrollo

### 1. Iniciar el entorno completo

```bash
# Windows
.\scripts\dev.bat

# Linux/Mac
./scripts/dev.sh
```

Esto automáticamente:
- ✅ Ejecuta migraciones del property-svc
- ✅ Inicia el servicio en puerto 4004
- ✅ Configura el API Gateway con las rutas

### 2. Verificar funcionamiento

```bash
# Health check
curl http://localhost:4004/health

# A través del API Gateway (RECOMENDADO)
curl -H "X-Api-Key: supersecret-api-key-for-dev" \
     http://localhost:8081/api/properties/

# Listar amenities
curl -H "X-Api-Key: supersecret-api-key-for-dev" \
     http://localhost:8081/api/amenities/
```

### 3. Testing con Postman

La colección `REM-API-Collection.postman_collection.json` incluye:

- **🏢 Properties & Real Estate** section completa
- Variables de entorno preconfiguradas
- Ejemplos de payloads listos para usar
- Auto-save de property_id y amenity_id

## 🔄 Integración con otros servicios

### Address Service (gRPC)
```bash
# Validación de direcciones
- GetAddress(address_id) -> AddressResponse
- Fallback a mock si no está disponible
```

### Person Service (gRPC)
```bash
# Validación de propietarios
- GetPerson(person_id) -> PersonResponse  
- Fallback a mock si no está disponible
```

### API Gateway
```nginx
# nginx.conf - Ya configurado
location /api/properties {
    proxy_pass http://property_backend/;
}

location /api/amenities {
    proxy_pass http://property_backend/amenities/;
}
```

## 📊 Base de Datos

### Migraciones Aplicadas

1. **001_create_properties_table.up.sql** - Tabla principal de propiedades
2. **002_create_amenities_table.up.sql** - Tabla de amenities/comodidades
3. **003_create_property_amenities_table.up.sql** - Tabla de relaciones N:N

### Estructura de Tablas

```sql
-- Propiedades
properties (
    id UUID PRIMARY KEY,
    owner_person_id UUID NOT NULL,
    address_id UUID NOT NULL,
    property_type property_type_enum NOT NULL,
    internal_code VARCHAR(20) UNIQUE NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- Amenities  
amenities (
    id UUID PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    category amenity_category_enum NOT NULL,
    icon VARCHAR(10),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- Relaciones
property_amenities (
    id UUID PRIMARY KEY,
    property_id UUID REFERENCES properties(id),
    amenity_id UUID REFERENCES amenities(id),
    created_at TIMESTAMP,
    UNIQUE(property_id, amenity_id)
);
```

## 🎯 Próximos Pasos

Una vez que los otros servicios (address-svc y person-svc) estén corriendo:

1. **Actualizar configuración** para usar servicios reales en lugar de mocks
2. **Implementar validaciones gRPC** completas
3. **Agregar filtros avanzados** por tipo de propiedad, amenities, etc.
4. **Implementar búsqueda geoespacial** con address-svc
5. **Agregar imágenes** y media para propiedades

## 📋 Resumen de Archivos

### Nuevos archivos creados:
- `property-svc/` - Microservicio completo
- `property-svc/migrations/pg/` - Migraciones de BD
- `property-svc/src/` - Código fuente (models, DTOs, controllers, etc.)
- `property-svc/cmd/` - Entry points (api, migrate)

### Archivos actualizados:
- `REM-API-Collection.postman_collection.json` - Endpoints agregados
- `REM-Development.postman_environment.json` - Variables agregadas
- `README.md` - Documentación actualizada
- `nginx.conf` - Rutas ya configuradas
- `scripts/dev.bat` y `scripts/dev.sh` - Scripts actualizados

---

**✅ El Property Service está completamente implementado y listo para uso en desarrollo.**
