# 👥 Person Service (person-svc)

Microservicio centralizado para la gestión de personas (individuos y empresas) en **REM Gestión**. Maneja el registro, consulta y administración de entidades personales junto con sus datos de contacto, integrándose con el servicio de direcciones para proporcionar información completa.

---

## 🎯 Propósito

Este servicio fue diseñado como un microservicio independiente para:

- **Centralizar la gestión de personas** (individuos y empresas) evitando duplicación
- **Manejo unificado de contactos** (email, teléfono, WhatsApp) con soporte para contactos primarios
- **Integración con address-svc** vía gRPC para datos de ubicación completos
- **Soporte para operaciones masivas** (bulk operations) con manejo de errores parciales
- **Validación y normalización** de datos personales (DNI, CUIT, emails)
- **Escalabilidad independiente** según la demanda de gestión de personas

---

## 🏗️ Arquitectura

```
person-svc/
├── cmd/
│   ├── api/main.go           # Punto de entrada del servidor HTTP/gRPC
│   └── migrate/main.go       # Herramienta de migraciones de BD
├── src/
│   ├── controllers/          # Capa de presentación (HTTP handlers)
│   │   └── person_controller.go
│   ├── services/             # Lógica de negocio y validaciones
│   │   └── person_service.go
│   ├── repository/           # Acceso a datos (GORM)
│   │   └── person_repo.go
│   ├── models/              # Entidades de dominio (Person, Individual, Company, Contacto)
│   │   └── person_contact.go
│   ├── dto/                 # Data Transfer Objects (request/response)
│   │   └── person.go
│   ├── router/              # Configuración de rutas HTTP
│   │   └── router.go
│   └── grpc/                # Handlers para comunicación gRPC
│       └── handler.go
├── migrations/pg/           # Scripts de migración PostgreSQL
├── .env                     # Variables de entorno locales
├── .air.toml               # Configuración de hot-reload
└── go.mod                  # Dependencias del módulo
```

### Flujo de datos:
```
HTTP Request → Router → Controller → Service → Repository → Database
                ↓                      ↓
HTTP Response ← JSON ← DTO ← Business Logic ← Model ← PostgreSQL
                              ↓
                        gRPC Client → address-svc (para direcciones completas)
```

### Tipos de entidades:
```
Person (base)
├── Individual (persona física)
│   ├── first_name, last_name, dni
│   └── Contactos[]
└── Company (persona jurídica)
    ├── legal_name, cuit, society_type
    └── Contactos[]
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
REM_POSTGRES_DBNAME=remgestion
REM_POSTGRES_SSLMODE=disable

# Configuración del servicio
REM_API_KEY=PersonSvcSecretKey
REM_LOGGER_LEVEL=info
REM_SERVER_PORT=4000

# gRPC configuration (servidor)
REM_GRPC_HOST=0.0.0.0
REM_GRPC_PORT=50052

# gRPC address service (cliente para address-svc)
REM_ADDRESS_HOST=localhost
REM_ADDRESS_PORT=50051

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
go build -o person-svc ./cmd/api
./person-svc
```

---

## 📡 API Reference

**Base URL:** `http://localhost:4000`  
**Autenticación:** Header `X-Api-Key: PersonSvcSecretKey`

### **Personas (CRUD básico)**

#### **POST /persons** - Crear persona
Crea una nueva persona (individual o empresa) con opción de dirección y contactos.

```http
POST /persons
X-Api-Key: PersonSvcSecretKey
Content-Type: application/json

# Ejemplo 1: Individual con dirección nueva
{
  "type": "individual",
  "first_name": "Juan",
  "last_name": "Pérez",
  "dni": "12345678",
  "sexo": "masculino",
  "address_payload": {
    "street": "Av. Corrientes",
    "number": 1234,
    "city": "Buenos Aires",
    "country": "AR"
  },
  "contacts": [
    {
      "tipo": "email",
      "dato": "juan.perez@email.com",
      "is_primary": true
    }
  ]
}

# Ejemplo 2: Empresa con address_id existente
{
  "type": "company",
  "legal_name": "Acme Corp S.A.",
  "cuit": "20123456789",
  "society_type": "SA",
  "address_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Respuesta exitosa (201):**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "type": "individual",
  "address_id": "550e8400-e29b-41d4-a716-446655440000",
  "avatar_url": null,
  "sexo": "masculino",
  "created_at": "2024-01-15T10:30:00Z",
  "individual": {
    "person_id": "123e4567-e89b-12d3-a456-426614174000",
    "first_name": "Juan",
    "last_name": "Pérez",
    "dni": "12345678"
  },
  "contactos": [
    {
      "id": "789e0123-e45f-67a8-9bcd-ef0123456789",
      "persona_id": "123e4567-e89b-12d3-a456-426614174000",
      "tipo": "email",
      "dato": "juan.perez@email.com",
      "is_primary": true
    }
  ]
}
```

#### **GET /persons/:id** - Obtener persona básica
```http
GET /persons/123e4567-e89b-12d3-a456-426614174000
X-Api-Key: PersonSvcSecretKey
```

#### **GET /persons/:id/full** - Obtener persona con dirección completa
Devuelve la persona con la dirección expandida desde address-svc vía gRPC.

```http
GET /persons/123e4567-e89b-12d3-a456-426614174000/full
X-Api-Key: PersonSvcSecretKey
```

**Respuesta (200):**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "type": "individual",
  "address_id": "550e8400-e29b-41d4-a716-446655440000",
  "sexo": "masculino",
  "individual": {
    "first_name": "Juan",
    "last_name": "Pérez",
    "dni": "12345678"
  },
  "contactos": [...],
  "address": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "street": "Av. Corrientes",
    "number": 1234,
    "city": "Buenos Aires",
    "country": "AR"
  }
}
```

#### **GET /persons** - Listar personas con filtros
```http
GET /persons?page=1&per_page=20&search=juan&type=individual&with_contacts=true
X-Api-Key: PersonSvcSecretKey
```

**Parámetros:**
- `page`: Número de página (default: 1)
- `per_page`: Elementos por página (default: 20, max: 100)
- `search`: Búsqueda en nombres/razón social
- `type`: Filtrar por tipo (`individual` | `company`)
- `with_contacts`: Incluir contactos (`true` | `false`)

#### **PUT /persons/:id** - Actualizar persona
```http
PUT /persons/123e4567-e89b-12d3-a456-426614174000
X-Api-Key: PersonSvcSecretKey
Content-Type: application/json

{
  "first_name": "Juan Carlos",
  "avatar_url": "https://example.com/avatar.jpg"
}
```

#### **DELETE /persons/:id** - Eliminar persona
```http
DELETE /persons/123e4567-e89b-12d3-a456-426614174000
X-Api-Key: PersonSvcSecretKey
```
**Respuesta (204):** Sin contenido

---

### **Operaciones masivas**

#### **POST /persons/bulk** - Creación masiva
Permite crear múltiples personas en una sola operación con manejo de errores parciales.

```http
POST /persons/bulk
X-Api-Key: PersonSvcSecretKey
Content-Type: application/json

[
  {
    "type": "individual",
    "first_name": "María",
    "last_name": "García",
    "dni": "87654321"
  },
  {
    "type": "company",
    "legal_name": "Tech Solutions SRL",
    "cuit": "20987654321"
  }
]
```

**Respuesta (201/207):**
- **201 Created**: Si todas las personas se crearon exitosamente
- **207 Multi-Status**: Si hubo errores parciales

```json
{
  "success": [
    {
      "index": 0,
      "person": {
        "id": "...",
        "type": "individual",
        "individual": {"first_name": "María", "last_name": "García"}
      }
    }
  ],
  "errors": [
    {
      "index": 1,
      "error": "CUIT ya existe"
    }
  ],
  "total": 2
}
```

---

### **Contactos**

#### **POST /persons/:id/contacts** - Agregar contacto
```http
POST /persons/123e4567-e89b-12d3-a456-426614174000/contacts
X-Api-Key: PersonSvcSecretKey
Content-Type: application/json

{
  "tipo": "phone",
  "dato": "+54911234567",
  "is_primary": false
}
```

#### **GET /persons/:id/contacts** - Listar contactos de una persona
```http
GET /persons/123e4567-e89b-12d3-a456-426614174000/contacts
X-Api-Key: PersonSvcSecretKey
```

#### **GET /contacts/primary/:personId** - Obtener contacto primario
Devuelve solo el contacto marcado como primario de la persona.

```http
GET /contacts/primary/123e4567-e89b-12d3-a456-426614174000
X-Api-Key: PersonSvcSecretKey
```

#### **PUT /contacts/:contactID** - Actualizar contacto
```http
PUT /contacts/789e0123-e45f-67a8-9bcd-ef0123456789
X-Api-Key: PersonSvcSecretKey
Content-Type: application/json

{
  "dato": "nuevo.email@example.com",
  "is_primary": true
}
```

#### **DELETE /contacts/:contactID** - Eliminar contacto
```http
DELETE /contacts/789e0123-e45f-67a8-9bcd-ef0123456789
X-Api-Key: PersonSvcSecretKey
```

---

## 🛡️ Validaciones y reglas de negocio

### Tipos de persona:
- **`individual`**: Persona física (requiere `first_name`, `last_name`)
- **`company`**: Persona jurídica (requiere `legal_name`)

### Campos requeridos por tipo:

#### Individual:
- `type: "individual"`
- `first_name` (máx. 64 caracteres)
- `last_name` (máx. 64 caracteres)

#### Company:
- `type: "company"`  
- `legal_name` (máx. 120 caracteres)

### Campos opcionales:
- `dni` (8 caracteres, único) - Solo para individuales
- `cuit` (11 caracteres, único) - Solo para empresas
- `society_type` (máx. 12 caracteres) - Solo para empresas
- `sexo` (`masculino` | `femenino`) - Solo para individuales
- `avatar_url` - URL de imagen de perfil
- `address_id` - UUID de dirección existente
- `address_payload` - Datos para crear nueva dirección

### Contactos:
- **Tipos válidos**: `email`, `phone`, `whatsapp`
- **Campo `dato`**: Requerido, máx. 128 caracteres
- **Contacto primario**: Solo uno por persona (automáticamente establecido si es el primero)
- **Validación de email**: Formato válido para tipo `email`

### Reglas de dirección:
- **Escenario 1**: Proporcionar `address_id` existente
- **Escenario 2**: Proporcionar `address_payload` para crear nueva vía address-svc
- **Mutualidad**: No se puede enviar ambos al mismo tiempo

### Límites operacionales:
- **Bulk create**: Máximo 100 personas por operación
- **Paginación**: Máximo 100 elementos por página
- **Search**: Busca en nombres/razón social (mín. 2 caracteres)

---

## 🔍 Códigos de respuesta

| Código | Descripción | Cuándo ocurre |
|--------|-------------|---------------|
| **200** | OK | GET exitoso, PUT exitoso |
| **201** | Created | POST exitoso (persona o contacto creado) |
| **204** | No Content | DELETE exitoso |
| **207** | Multi-Status | Bulk create con errores parciales |
| **400** | Bad Request | JSON malformado, validación fallida, límites excedidos |
| **401** | Unauthorized | API Key faltante o inválida |
| **404** | Not Found | Persona o contacto no existe |
| **409** | Conflict | DNI/CUIT duplicado, violación de restricciones |
| **422** | Unprocessable Entity | Datos válidos pero lógica de negocio falló |
| **500** | Internal Server Error | Error de base de datos, gRPC o lógica interna |
| **503** | Service Unavailable | address-svc no disponible (para operaciones con dirección) |

---

## 🗄️ Esquema de base de datos

### Tabla principal `person`

```sql
CREATE TABLE person (
  id          uuid         PRIMARY KEY DEFAULT uuid_generate_v4(),
  type        varchar(10)  NOT NULL CHECK (type IN ('individual','company')),
  address_id  uuid,        -- FK lógica a address-svc (sin constraint)
  avatar_url  text,
  sexo        sexo,        -- enum: 'masculino', 'femenino'
  
  created_at  timestamp    NOT NULL DEFAULT now(),
  updated_at  timestamp,
  updated_by  uuid,
  deleted_at  timestamp
);
```

### Tabla `individual` (herencia)

```sql
CREATE TABLE individual (
  person_id   uuid PRIMARY KEY REFERENCES person(id) ON DELETE CASCADE,
  first_name  varchar(64) NOT NULL,
  last_name   varchar(64) NOT NULL,
  dni         char(8) UNIQUE,
  
  created_at  timestamp NOT NULL DEFAULT now(),
  updated_at  timestamp,
  updated_by  uuid,
  deleted_at  timestamp
);
```

### Tabla `company` (herencia)

```sql
CREATE TABLE company (
  person_id    uuid PRIMARY KEY REFERENCES person(id) ON DELETE CASCADE,
  legal_name   varchar(120) NOT NULL,
  cuit         char(11) UNIQUE,
  society_type varchar(12),
  
  created_at   timestamp NOT NULL DEFAULT now(),
  updated_at   timestamp,
  updated_by   uuid,
  deleted_at   timestamp
);
```

### Tabla `contacto`

```sql
CREATE TABLE contacto (
  id          uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  persona_id  uuid NOT NULL REFERENCES person(id) ON DELETE CASCADE,
  tipo        varchar(16) NOT NULL CHECK (tipo IN ('email','phone','whatsapp')),
  dato        varchar(128) NOT NULL,
  is_primary  boolean DEFAULT false,
  
  created_at  timestamp NOT NULL DEFAULT now(),
  created_by  uuid,
  updated_at  timestamp,
  updated_by  uuid,
  deleted_at  timestamp
);

-- Índices para performance
CREATE INDEX idx_persona      ON contacto (persona_id);
CREATE INDEX idx_persona_tipo ON contacto (persona_id, tipo);
```

### Relaciones:
- **1:1** - `person` ↔ `individual` | `company` (herencia por tabla)
- **1:N** - `person` → `contacto` (una persona tiene múltiples contactos)
- **1:1** - `person` → `address` (FK lógica, sin constraint por ser externo)

---

## 🔧 Configuración y deployment

### Desarrollo local
```bash
# 1. Clonar y navegar
git clone <repo>
cd person-svc

# 2. Configurar .env
cp .env.example .env
# Editar variables según tu entorno
# - Configurar PostgreSQL
# - Configurar address-svc host/port
# - Establecer API key única

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
RUN go mod tidy && go build -o person-svc ./cmd/api

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/person-svc .
EXPOSE 4000 50052
CMD ["./person-svc"]
```

### Variables de entorno en producción
```env
REM_DB_DRIVER=postgres
REM_POSTGRES_HOST=db-cluster.region.rds.amazonaws.com
REM_POSTGRES_PORT=5432
REM_POSTGRES_USER=rem_user
REM_POSTGRES_PASSWORD=secure_password_here
REM_POSTGRES_DBNAME=rem_production
REM_POSTGRES_SSLMODE=require

REM_API_KEY=ProductionPersonSvcApiKey123
REM_LOGGER_LEVEL=warn
REM_SERVER_PORT=4000

# gRPC Configuration
REM_GRPC_HOST=0.0.0.0
REM_GRPC_PORT=50052
REM_ADDRESS_HOST=address-svc.internal
REM_ADDRESS_PORT=50051
```

---

## 🧪 Testing

### Healthcheck
```bash
curl -H "X-Api-Key: PersonSvcSecretKey" \
     http://localhost:4000/health
```

### Crear persona individual
```bash
curl -X POST \
  -H "X-Api-Key: PersonSvcSecretKey" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "individual",
    "first_name": "Juan",
    "last_name": "Pérez",
    "dni": "12345678",
    "address_payload": {
      "street": "Av. 9 de Julio",
      "number": 1000,
      "city": "Buenos Aires",
      "country": "AR"
    }
  }' \
  http://localhost:4000/persons
```

### Crear empresa
```bash
curl -X POST \
  -H "X-Api-Key: PersonSvcSecretKey" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "company",
    "legal_name": "Tech Solutions SRL",
    "cuit": "20123456789",
    "society_type": "SRL"
  }' \
  http://localhost:4000/persons
```

### Obtener persona con dirección completa
```bash
curl -H "X-Api-Key: PersonSvcSecretKey" \
     http://localhost:4000/persons/{id}/full
```

### Agregar contacto
```bash
curl -X POST \
  -H "X-Api-Key: PersonSvcSecretKey" \
  -H "Content-Type: application/json" \
  -d '{
    "tipo": "email",
    "dato": "contacto@empresa.com",
    "is_primary": true
  }' \
  http://localhost:4000/persons/{id}/contacts
```

### Creación masiva
```bash
curl -X POST \
  -H "X-Api-Key: PersonSvcSecretKey" \
  -H "Content-Type: application/json" \
  -d '[
    {
      "type": "individual",
      "first_name": "María",
      "last_name": "García"
    },
    {
      "type": "company", 
      "legal_name": "Acme Corp"
    }
  ]' \
  http://localhost:4000/persons/bulk
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
  "path": "/persons",
  "status": 201,
  "latency_ms": 145.2,
  "user_agent": "curl/7.68.0",
  "svc": "[PERSON-SVC]"
}

{
  "level": "info", 
  "ts": "2024-01-15T10:30:05Z",
  "msg": "grpc call to address-svc",
  "method": "GetAddress",
  "address_id": "550e8400-e29b-41d4-a716-446655440000",
  "latency_ms": 23.1,
  "success": true
}
```

### Métricas importantes
- **Latencia promedio** por endpoint (person CRUD, contacts, bulk)
- **Rate de personas individuales vs empresas** creadas
- **Errores 4xx/5xx** por minuto y endpoint
- **Throughput** de requests por segundo
- **gRPC calls a address-svc**: latencia, éxito/fallo
- **Bulk operations**: personas procesadas vs errores
- **Contactos primarios**: distribución por tipo (email/phone/whatsapp)

### Alertas sugeridas
- **Latencia > 500ms** en endpoints básicos (GET, POST individual)
- **Error rate > 5%** en 5 minutos
- **address-svc unreachable** por más de 30 segundos
- **Bulk errors > 20%** en operaciones masivas

---

## 🤝 Integración con otros servicios

### Con address-svc (gRPC):
```go
// person-svc como cliente gRPC
func (s *PersonService) GetFullPerson(personID string) (*dto.FullPersonResponse, error) {
    // 1. Obtener persona de BD local
    person, err := s.repo.GetByID(personID)
    if err != nil {
        return nil, err
    }
    
    // 2. Si tiene address_id, consultar address-svc vía gRPC
    if person.AddressID != nil {
        addressResp, err := s.addressClient.GetAddress(*person.AddressID)
        if err != nil {
            // Log error pero no fallar - devolver persona sin dirección
            s.logger.Warn("failed to get address", zap.Error(err))
        }
    }
    
    return &dto.FullPersonResponse{
        Person: person,
        Address: addressResponse,
    }, nil
}
```

### Desde property-service:
```go
type Property struct {
    ID       string `json:"id"`
    OwnerID  string `json:"owner_id"`  // FK a person-svc (individual/company)
    TenantID string `json:"tenant_id"` // FK a person-svc (individual)
    // ...otros campos
}

// Al crear propiedad, primero validar que la persona existe
personResp, err := personServiceClient.GetPerson(ownerID)
if err != nil {
    return fmt.Errorf("owner not found: %w", err)
}
```

### Desde contract-service:
```go
type Contract struct {
    ID         string `json:"id"`
    LandlordID string `json:"landlord_id"` // FK a person-svc
    TenantID   string `json:"tenant_id"`   // FK a person-svc
    
    // Obtener contactos para notificaciones
    LandlordContacts []Contact `json:"landlord_contacts"`
    TenantContacts   []Contact `json:"tenant_contacts"`
}

// Obtener contactos primarios para envío de contratos
landlordPrimary, err := personServiceClient.GetPrimaryContact(contract.LandlordID)
tenantPrimary, err := personServiceClient.GetPrimaryContact(contract.TenantID)
```

### Desde notification-service:
```go
// Enviar notificación usando contacto primario
func SendNotification(personID string, message string) error {
    contact, err := personServiceClient.GetPrimaryContact(personID)
    if err != nil {
        return err
    }
    
    switch contact.Tipo {
    case "email":
        return emailService.Send(contact.Dato, message)
    case "whatsapp":
        return whatsappService.Send(contact.Dato, message) 
    case "phone":
        return smsService.Send(contact.Dato, message)
    }
}
```

---

## 📡 Comunicación gRPC

### Como servidor gRPC (puerto 50052):
person-svc expone servicios gRPC para otros microservicios:

```protobuf
// person/v1/person.proto
service PersonService {
  rpc GetPerson(GetPersonRequest) returns (GetPersonResponse);
  rpc GetPersonWithAddress(GetPersonRequest) returns (GetPersonWithAddressResponse);
  rpc ValidatePersons(ValidatePersonsRequest) returns (ValidatePersonsResponse);
  rpc GetPrimaryContact(GetPersonRequest) returns (GetContactResponse);
}
```

### Como cliente gRPC (a address-svc):
person-svc consume address-svc para obtener direcciones completas:

```go
// Configuración del cliente gRPC a address-svc
addressClient := addresspb.NewAddressServiceClient(conn)

// Uso en GetFullPerson
addressResp, err := addressClient.GetAddress(context.Background(), &addresspb.GetAddressRequest{
    Id: person.AddressID.String(),
})
```

---

## 🚨 Troubleshooting

### Error: "connection refused"
```bash
# Verificar que PostgreSQL esté corriendo
pg_isready -h localhost -p 5432

# Verificar variables de entorno
echo $REM_POSTGRES_HOST
echo $REM_POSTGRES_DBNAME
```

### Error: "api key inválida"
```bash
# Verificar header en requests
curl -H "X-Api-Key: PersonSvcSecretKey" ...

# Verificar variable de entorno
echo $REM_API_KEY
```

### Error: "migration failed"
```bash
# Verificar permisos de usuario en BD
GRANT ALL PRIVILEGES ON DATABASE remgestion TO usuario;

# Forzar rollback y re-aplicar
go run ./cmd/migrate down
go run ./cmd/migrate up
```

### Error: "DNI/CUIT ya existe"
```bash
# Es esperado por validación de unicidad
# Opción 1: Usar el registro existente
curl -H "X-Api-Key: PersonSvcSecretKey" \
     http://localhost:4000/persons?search=12345678

# Opción 2: Actualizar el registro existente
curl -X PUT \
  -H "X-Api-Key: PersonSvcSecretKey" \
  -d '{"first_name": "Nuevo Nombre"}' \
  http://localhost:4000/persons/{existing-id}
```

### Error: "address-svc unreachable"
```bash
# Verificar conectividad a address-svc
telnet localhost 50051

# Verificar variables de entorno
echo $REM_ADDRESS_HOST
echo $REM_ADDRESS_PORT

# Revisar logs de address-svc
docker logs address-svc

# Fallback: person-svc funcionará sin direcciones expandidas
# Solo las operaciones GetFull fallarán
```

### Error: "person not found" en otros servicios
```bash
# Verificar que person-svc esté corriendo
curl -H "X-Api-Key: PersonSvcSecretKey" \
     http://localhost:4000/persons/{id}

# Verificar gRPC server
grpcurl -plaintext localhost:50052 list

# Verificar logs de person-svc
tail -f logs/person-svc.log
```

### Performance issues
```bash
# Verificar índices de BD
SELECT schemaname, tablename, indexname 
FROM pg_indexes 
WHERE schemaname = 'public';

# Analizar queries lentas
SELECT query, mean_time, calls 
FROM pg_stat_statements 
ORDER BY mean_time DESC 
LIMIT 10;

# Verificar conexiones de BD
SELECT count(*) FROM pg_stat_activity 
WHERE datname = 'remgestion';
```

---

## 📚 Dependencias clave

### Core dependencies:
- **rem-common** - Configuración, logging, middleware, DB connection, gRPC utilities
- **gin-gonic/gin** - Framework HTTP para API REST
- **gorm.io/gorm** - ORM para PostgreSQL con soporte para soft deletes
- **google/uuid** - Generación y validación de UUIDs
- **golang-migrate/migrate** - Sistema de migraciones de base de datos

### gRPC & Communication:
- **google.golang.org/grpc** - Framework gRPC para comunicación entre servicios
- **google.golang.org/protobuf** - Protocol Buffers para serialización
- **address-svc** - Servicio externo para gestión de direcciones (via gRPC)

### Database & Validation:
- **github.com/lib/pq** - Driver PostgreSQL
- **github.com/go-playground/validator/v10** - Validación de structs y DTOs
- **gorm.io/driver/postgres** - Driver GORM para PostgreSQL

### Development & Monitoring:
- **github.com/cosmtrek/air** - Hot reload para desarrollo
- **go.uber.org/zap** - Logger estructurado de alto rendimiento
- **github.com/gin-contrib/cors** - Middleware CORS para desarrollo

---

## 🚀 Estado del servicio

### ✅ Implementado y funcional:
- [x] CRUD completo de personas (individual/company)
- [x] Gestión de contactos con soporte para contactos primarios  
- [x] Integración gRPC con address-svc para direcciones completas
- [x] Operaciones masivas (bulk create) con manejo de errores parciales
- [x] Validaciones de negocio (DNI único, CUIT único, emails válidos)
- [x] API REST con autenticación via API key
- [x] Servidor gRPC para otros microservicios
- [x] Migraciones de base de datos
- [x] Logging estructurado
- [x] Soft deletes y auditoría básica

### 🔄 Próximas mejoras sugeridas:
- [ ] Tests unitarios y de integración automatizados
- [ ] Validación avanzada de DNI/CUIT (algoritmo de verificación)
- [ ] Cache con Redis para consultas frecuentes
- [ ] Rate limiting por API key
- [ ] Métricas de Prometheus/monitoring avanzado
- [ ] Documentación OpenAPI/Swagger
- [ ] Webhooks para notificar cambios a servicios consumidores
- [ ] Bulk update y bulk delete operations
- [ ] Versionado de API (v1, v2)
- [ ] Compresión gRPC y HTTP

---

**Repositorio:** [rem-backend/services/person-svc](.)  
**Documentación:** Este README.md + [rem-common docs](../rem-common/README.md)  
**Proto definitions:** [rem-common/protos/person/](../rem-common/protos/person/)

**Versión:** v1.0.0  
**Última actualización:** Diciembre 2024

