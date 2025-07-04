# 🏠 Address Service (address-svc)

Microservicio centralizado para la gestión de direcciones en **REM Gestión**. Proporciona un punto único para crear, consultar y gestionar direcciones que son utilizadas por múltiples entidades del sistema (propiedades, usuarios, oficinas, etc.).

---

## 🎯 Propósito

Este servicio fue extraído como un microservicio independiente para:

- **Centralizar la gestión de direcciones** y evitar duplicación de código
- **Prevenir registros duplicados** mediante deduplicación inteligente
- **Reutilización entre servicios** (properties, users, organizations, etc.)
- **Mantener consistencia** en el formato y validación de direcciones
- **Escalabilidad independiente** según la demanda de geolocalización

---

## 🏗️ Arquitectura

```
address-svc/
├── cmd/
│   ├── api/main.go           # Punto de entrada del servidor HTTP
│   └── migrate/main.go       # Herramienta de migraciones
├── src/
│   ├── controllers/          # Capa de presentación (HTTP handlers)
│   ├── services/             # Lógica de negocio
│   ├── repository/           # Acceso a datos (GORM)
│   ├── models/              # Entidades de dominio
│   ├── dto/                 # Data Transfer Objects
│   └── router/              # Configuración de rutas
├── migrations/pg/           # Scripts de migración PostgreSQL
├── .env                     # Variables de entorno locales
├── .air.toml               # Configuración de hot-reload
└── go.mod                  # Dependencias del módulo
```

### Flujo de datos:
```
HTTP Request → Router → Controller → Service → Repository → Database
                ↓
HTTP Response ← JSON ← DTO ← Business Logic ← Model ← PostgreSQL
```

---

## 🚀 Instalación y Setup

### 1. Variables de entorno requeridas

```env
# Motor de base de datos
REM_DB_DRIVER=postgres

# Configuración PostgreSQL
REM_POSTGRES_HOST=localhost
REM_POSTGRES_PORT=5432
REM_POSTGRES_USER=usuario
REM_POSTGRES_PASSWORD=password
REM_POSTGRES_DB=rem_addresses
REM_POSTGRES_SSLMODE=disable

# Configuración del servicio
REM_API_KEY=AddressSvcSecretKey
REM_LOGGER_LEVEL=info
REM_SERVER_PORT=4000

# Redis/Rabbit (opcionales, valores por defecto)
REM_REDIS_ENABLED=false
REM_RABBIT_HOST=localhost
REM_RABBIT_PORT=5672
REM_RABBIT_USER=user
REM_RABBIT_PASSWORD=password
```

### 2. Preparar base de datos

```bash
# Ejecutar migraciones
go run ./cmd/migrate

# Rollback (cuidado!)
go run ./cmd/migrate down

# Rollback específico
go run ./cmd/migrate "steps -1"
```

### 3. Ejecutar el servicio

```bash
# Desarrollo con hot-reload
air

# Producción
go run ./cmd/api

# Build
go build -o address-svc ./cmd/api
./address-svc
```

---

## 📡 API Reference

**Base URL:** `http://localhost:4000`  
**Autenticación:** Header `X-Api-Key: AddressSvcSecretKey`

### **POST /addresses** - Crear dirección
Crea una nueva dirección con deduplicación automática.

```http
POST /addresses
X-Api-Key: AddressSvcSecretKey
Content-Type: application/json

{
  "floor": "2",
  "unit": "A",
  "street": "Av. Corrientes",
  "number": 1234,
  "city": "Buenos Aires",
  "state": "CABA",
  "zip": "C1043AAZ",
  "country": "AR"
}
```

**Respuesta exitosa (201):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "floor": "2",
  "unit": "A", 
  "street": "Av. Corrientes",
  "number": 1234,
  "city": "Buenos Aires",
  "state": "CABA",
  "zip": "C1043AAZ",
  "country": "AR",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

**Deduplicación:** Si una dirección idéntica ya existe, retorna la existente en lugar de crear duplicada.

### **GET /addresses/:id** - Obtener dirección
```http
GET /addresses/550e8400-e29b-41d4-a716-446655440000
X-Api-Key: AddressSvcSecretKey
```

**Respuesta (200):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "floor": "2",
  "unit": "A",
  "street": "Av. Corrientes",
  "number": 1234,
  "city": "Buenos Aires", 
  "state": "CABA",
  "zip": "C1043AAZ",
  "country": "AR",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

### **PUT /addresses/:id** - Actualizar dirección
⚠️ **Inmutabilidad:** Las direcciones son inmutables. Este endpoint siempre retorna `403 Forbidden`.

```http
PUT /addresses/550e8400-e29b-41d4-a716-446655440000
X-Api-Key: AddressSvcSecretKey
Content-Type: application/json

{
  "city": "Rosario"
}
```

**Respuesta (403):**
```json
{
  "error": "addresses are immutable; create a new record"
}
```

### **DELETE /addresses/:id** - Eliminar dirección
```http
DELETE /addresses/550e8400-e29b-41d4-a716-446655440000
X-Api-Key: AddressSvcSecretKey
```

**Respuesta (204):** Sin contenido

---

## 🛡️ Validaciones y reglas de negocio

### Campos requeridos:
- `street` (máx. 120 caracteres)
- `number` (entero)
- `city` (máx. 64 caracteres)
- `country` (código ISO 2 letras, ej: "AR", "US")

### Campos opcionales:
- `floor` (máx. 16 caracteres)
- `unit` (máx. 16 caracteres)
- `state` (máx. 64 caracteres)
- `zip` (máx. 16 caracteres)

### Deduplicación:
El sistema compara direcciones normalizando:
- Texto en minúsculas y sin espacios extra
- País en mayúsculas
- Número exacto

Si detecta duplicado durante CREATE, retorna la dirección existente.

### Inmutabilidad:
Las direcciones no se pueden modificar una vez creadas. Si necesitas cambios, debes crear una nueva dirección.

---

## 🔍 Códigos de respuesta

| Código | Descripción | Cuándo ocurre |
|--------|-------------|---------------|
| **200** | OK | GET exitoso |
| **201** | Created | POST exitoso (nueva o deduplicada) |
| **204** | No Content | DELETE exitoso |
| **400** | Bad Request | JSON malformado o validación fallida |
| **401** | Unauthorized | API Key faltante o inválida |
| **403** | Forbidden | Intento de UPDATE (inmutable) |
| **404** | Not Found | ID de dirección no existe |
| **500** | Internal Server Error | Error de base de datos o lógica |

---

## 🗄️ Esquema de base de datos

### Tabla `addresses`

```sql
CREATE TABLE addresses (
  id         CHAR(36) PRIMARY KEY,
  floor      VARCHAR(16),
  unit       VARCHAR(16), 
  street     VARCHAR(120),
  number     INT,
  city       VARCHAR(64),
  state      VARCHAR(64),
  zip        VARCHAR(16),
  country    CHAR(2),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL,
  updated_by CHAR(36),
  deleted_at TIMESTAMP NULL
);

-- Índice único para prevenir duplicados
CREATE UNIQUE INDEX addresses_unique_addr
    ON addresses (
        COALESCE(lower(trim(floor)),   ''),
        COALESCE(lower(trim(unit)),    ''),
        COALESCE(lower(trim(street)),  ''),
        number,
        COALESCE(lower(trim(city)),    ''),
        COALESCE(lower(trim(state)),   ''),
        COALESCE(trim(zip),            ''),
        upper(trim(country))
    );
```

---

## 🔧 Configuración y deployment

### Desarrollo local
```bash
# 1. Clonar y navegar
git clone <repo>
cd address-svc

# 2. Configurar .env
cp .env.example .env
# Editar variables según tu entorno

# 3. Instalar dependencias
go mod tidy

# 4. Ejecutar migraciones
go run ./cmd/migrate

# 5. Iniciar con hot-reload
air
```

### Docker
```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o address-svc ./cmd/api

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/address-svc .
EXPOSE 4000
CMD ["./address-svc"]
```

### Variables de entorno en producción
```env
REM_DB_DRIVER=postgres
REM_POSTGRES_HOST=db-cluster.region.rds.amazonaws.com
REM_POSTGRES_PORT=5432
REM_POSTGRES_USER=rem_user
REM_POSTGRES_PASSWORD=secure_password_here
REM_POSTGRES_DB=rem_production
REM_POSTGRES_SSLMODE=require
REM_API_KEY=ProductionSecretApiKey123
REM_LOGGER_LEVEL=warn
REM_SERVER_PORT=4000
```

---

## 🧪 Testing

### Healthcheck
```bash
curl -H "X-Api-Key: AddressSvcSecretKey" \
     http://localhost:4000/addresses/health
```

### Crear dirección de prueba
```bash
curl -X POST \
  -H "X-Api-Key: AddressSvcSecretKey" \
  -H "Content-Type: application/json" \
  -d '{
    "street": "Av. 9 de Julio",
    "number": 1000,
    "city": "Buenos Aires",
    "country": "AR"
  }' \
  http://localhost:4000/addresses
```

### Obtener dirección
```bash
curl -H "X-Api-Key: AddressSvcSecretKey" \
     http://localhost:4000/addresses/{id}
```

---

## 📊 Monitoring y logs

### Logs estructurados (JSON)
```json
{
  "level": "info",
  "ts": "2024-01-15T10:30:00Z",
  "msg": "http request",
  "method": "POST",
  "path": "/addresses",
  "status": 201,
  "latency_ms": 45.2,
  "svc": "[ADDRESS-SVC]"
}
```

### Métricas importantes
- **Latencia promedio** por endpoint
- **Rate de deduplicación** (direcciones existentes vs nuevas)
- **Errores 4xx/5xx** por minuto
- **Throughput** de requests por segundo

---

## 🤝 Integración con otros servicios

### Desde property-service:
```go
type Property struct {
    ID        string `json:"id"`
    AddressID string `json:"address_id"` // FK a address-svc
    // ...otros campos
}

// Al crear propiedad, primero crear/obtener address
addressResp := createAddress(addressData)
property.AddressID = addressResp.ID
```

### Desde user-service:
```go
type User struct {
    ID              string `json:"id"`
    HomeAddressID   string `json:"home_address_id,omitempty"`
    OfficeAddressID string `json:"office_address_id,omitempty"`
    // ...otros campos
}
```

---

## 🚨 Troubleshooting

### Error: "connection refused"
```bash
# Verificar que PostgreSQL esté corriendo
pg_isready -h localhost -p 5432

# Verificar variables de entorno
echo $REM_POSTGRES_HOST
```

### Error: "api key inválida"
```bash
# Verificar header en requests
curl -H "X-Api-Key: AddressSvcSecretKey" ...
```

### Error: "migration failed"
```bash
# Verificar permisos de usuario en DB
GRANT ALL PRIVILEGES ON DATABASE rem_addresses TO usuario;

# Forzar rollback y re-aplicar
go run ./cmd/migrate down
go run ./cmd/migrate up
```

### Error: "addresses are immutable"
```bash
# Es esperado. Para cambiar una dirección:
# 1. Crear nueva dirección
curl -X POST .../addresses -d '{...nueva dirección...}'

# 2. Actualizar referencia en el servicio padre
curl -X PUT .../properties/123 -d '{"address_id": "nueva-direccion-id"}'

# 3. Opcionalmente eliminar dirección antigua si no se usa
curl -X DELETE .../addresses/direccion-vieja-id
```

---

## 📚 Dependencias clave

- **rem-common** - Configuración, logging, middleware, DB
- **gin-gonic/gin** - Framework HTTP
- **gorm.io/gorm** - ORM para PostgreSQL  
- **google/uuid** - Generación de UUIDs
- **golang-migrate** - Sistema de migraciones

---

**Repositorio:** [rem-backend/services/address-svc](.)  
**Documentación:** Este README.md + [rem-common docs](../rem-common/README.md)

