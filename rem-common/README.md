# rem-common 📦  
Librería compartida para todos los microservicios de **REM Gestión** – centraliza
configuración, conexiones, logger, cache, broker, middlewares, utilidades
criptográficas, sistema de migraciones, **comunicación gRPC** y **definiciones de protobuf**.

---

## 🚀 Instalación

```bash
go get github.com/rem-gestion/rem-common@latest
```

Luego importa los paquetes que necesites:

```go
import (
    "github.com/rem-gestion/rem-common/config"
    "github.com/rem-gestion/rem-common/logger"
    "github.com/rem-gestion/rem-common/db"
    "github.com/rem-gestion/rem-common/middleware"
    "github.com/rem-gestion/rem-common/grpc"
    addresspb "github.com/rem-gestion/rem-common/protos/address/v1"
    personpb "github.com/rem-gestion/rem-common/protos/person/v1"
    // ...
)
```

---

## 📋 Contenido del módulo

| Carpeta / paquete          | Qué provee                                                                                           |
|----------------------------|-------------------------------------------------------------------------------------------------------|
| `config`                   | Carga variables `REM_*` (con soporte `.env`) en un struct único. Incluye configuración gRPC.        |
| `logger`                   | Instancia **Zap** con colores, caller y campo fijo `service`.                                         |
| `db`                       | Constructores `NewPostgres`, `NewMySQL`, `NewMongo`, `NewRedis`, `NewRabbit` + helper de migraciones. |
| `broker`                  | Helpers para RabbitMQ (DeclareExchange/Queue, Publish, Consume).                                      |
| `cache`                    | `NewRedis` con `Ping()`.                                                                              |
| `grpc`                     | **NUEVO**: Utilidades para servidor y cliente gRPC con interceptors de logging, health checks y keep-alive. |
| `protos`                   | **NUEVO**: Definiciones Protocol Buffers para address-svc y person-svc con código generado.           |
| `middleware`               | Gin middlewares: APIKeyAuth, UserAuth/JWT, AdminAuth, MembershipAuth, EmployeeAuth, GinLogger, RequestID, Recovery, ErrorHandler, RateLimit. |
| `errors`                   | Tipos de error tipificados para mapear a 400/401/403/404/409/422/429/500.                             |
| `utils`                    | Hash/Check de contraseñas (bcrypt) y JWT helper (Generate/Parse).                                     |
| `db/migrate.go`            | Helper de **golang-migrate** multi-motor (Postgres, MySQL, MongoDB).                                  |

---

## ⚙️ Variables de entorno principales

| Variable (`REM_…`)        | Descripción                               | Ejemplo                      |
|---------------------------|-------------------------------------------|------------------------------|
| `DB_DRIVER`               | `postgres`, `mysql` o `mongo`.            | `postgres`                   |
| `POSTGRES_HOST`           | Host de PostgreSQL                        | `localhost`                  |
| `POSTGRES_PORT`           | Puerto de PostgreSQL                      | `5432`                       |
| `POSTGRES_USER`           | Usuario de PostgreSQL                     | `admin`                      |
| `POSTGRES_PASSWORD`       | Contraseña de PostgreSQL                  | `secret`                     |
| `POSTGRES_DBNAME`         | Base de datos de PostgreSQL               | `rem_auth`                   |
| `POSTGRES_SSLMODE`        | Modo SSL de PostgreSQL                    | `disable`                    |
| `MYSQL_HOST`              | Host de MySQL                             | `localhost`                  |
| `MYSQL_PORT`              | Puerto de MySQL                           | `3306`                       |
| `MYSQL_USER`              | Usuario de MySQL                          | `root`                       |
| `MYSQL_PASSWORD`          | Contraseña de MySQL                       | `secret`                     |
| `MYSQL_DB`                | Base de datos de MySQL                    | `rem_auth`                   |
| `MYSQL_CHARSET`           | Charset de MySQL                          | `utf8mb4`                    |
| `MONGO_URI`               | URI completa de MongoDB                   | `mongodb://localhost:27017`  |
| `MONGO_HOST`              | Host de MongoDB                           | `localhost`                  |
| `MONGO_PORT`              | Puerto de MongoDB                         | `27017`                      |
| `MONGO_DB`                | Base de datos de MongoDB                  | `rem_auth`                   |
| `RABBIT_HOST`             | Host de RabbitMQ                          | `localhost`                  |
| `RABBIT_PORT`             | Puerto de RabbitMQ                        | `5672`                       |
| `RABBIT_USER`             | Usuario de RabbitMQ                       | `guest`                      |
| `RABBIT_PASSWORD`         | Contraseña de RabbitMQ                    | `guest`                      |
| `REDIS_ENABLED`           | `true/false` para activar cache           | `true`                       |
| `REDIS_ADDR`              | Dirección de Redis                        | `localhost:6379`             |
| `REDIS_PASSWORD`          | Contraseña de Redis                       | ` ` (vacío)                  |
| `REDIS_DB`                | Base de datos de Redis                    | `0`                          |
| `GRPC_HOST`               | **NUEVO**: Host del servidor gRPC         | `0.0.0.0`                    |
| `GRPC_PORT`               | **NUEVO**: Puerto del servidor gRPC       | `50051`                      |
| `ADDRESS_HOST`            | **NUEVO**: Host del address-svc           | `address-svc`                |
| `ADDRESS_PORT`            | **NUEVO**: Puerto del address-svc         | `50052`                      |
| `LOGGER_LEVEL`            | `debug`, `info`, `warn`, `error`          | `info`                       |
| `SERVER_PORT`             | Puerto del servidor HTTP                  | `8080`                       |
| `API_KEY`                 | Clave para `APIKeyAuth`                   | `your-secret-api-key`        |
| `JWT_SECRET`              | Secreto HS256 para tokens                 | `your-jwt-secret-key`        |

*(Ver [migraciones.md](./migraciones.md) para guía completa de migraciones)*

---

## 🔧 Ejemplo rápido de uso en un microservicio

### Ejemplo básico (HTTP + gRPC)
```go
// services/person-svc/cmd/api/main.go
package main

import (
    "net"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    "google.golang.org/grpc"

    "github.com/rem-gestion/rem-common/config"
    "github.com/rem-gestion/rem-common/logger"
    "github.com/rem-gestion/rem-common/db"
    "github.com/rem-gestion/rem-common/cache"
    "github.com/rem-gestion/rem-common/grpc" as remgrpc
    mw "github.com/rem-gestion/rem-common/middleware"
    personpb "github.com/rem-gestion/rem-common/protos/person/v1"
)

func main() {
    cfg := config.Load()
    lg := logger.New(cfg.Logger, "person-svc")
    
    // Conexiones de BD y cache
    pg, _ := db.NewPostgres(cfg.Postgres)
    rdb, _ := cache.NewRedis(cfg.Redis)
    
    // Servidor HTTP
    r := gin.New()
    r.Use(
        mw.RequestID(),
        mw.APIKeyAuth("X-Api-Key", cfg.APIKey),
        mw.GinLogger(lg),
        mw.RecoveryWithZap(lg),
        mw.ErrorHandler(),
    )
    
    // Setup de rutas REST
    setupRoutes(r, pg, rdb, lg)
    
    // Servidor gRPC paralelo
    go startGRPCServer(cfg, lg, pg)
    
    // Iniciar servidor HTTP
    lg.Info("starting HTTP server", zap.Int("port", cfg.Server.Port))
    r.Run(fmt.Sprintf(":%d", cfg.Server.Port))
}

func startGRPCServer(cfg *config.Config, lg *zap.Logger, db *gorm.DB) {
    lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port))
    if err != nil {
        lg.Fatal("failed to listen gRPC", zap.Error(err))
    }
    
    // Crear servidor gRPC con interceptors
    s := remgrpc.NewServer(lg, grpc.KeepaliveParams{
        MaxConnectionIdle:     30 * time.Second,
        MaxConnectionAge:      5 * time.Minute,
        MaxConnectionAgeGrace: 5 * time.Second,
        Time:                  30 * time.Second,
        Timeout:              5 * time.Second,
    })
    
    // Registrar servicios
    personHandler := &PersonGRPCHandler{db: db, logger: lg}
    personpb.RegisterPersonServiceServer(s, personHandler)
    
    lg.Info("starting gRPC server", zap.String("addr", lis.Addr().String()))
    if err := s.Serve(lis); err != nil {
        lg.Fatal("failed to serve gRPC", zap.Error(err))
    }
}
```

### Ejemplo como cliente gRPC
```go
// Conectar a address-svc desde person-svc
package service

import (
    "context"
    "fmt"
    
    "github.com/rem-gestion/rem-common/grpc"
    addresspb "github.com/rem-gestion/rem-common/protos/address/v1"
)

type PersonService struct {
    addressClient addresspb.AddressServiceClient
}

func NewPersonService(cfg *config.Config) *PersonService {
    // Conectar a address-svc vía gRPC
    conn, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.Address.Host, cfg.Address.Port))
    if err != nil {
        panic(err)
    }
    
    return &PersonService{
        addressClient: addresspb.NewAddressServiceClient(conn),
    }
}

func (s *PersonService) GetFullPerson(personID string) (*FullPersonResponse, error) {
    // ... obtener person de BD local ...
    
    // Si tiene address_id, obtener dirección vía gRPC
    if person.AddressID != nil {
        addressResp, err := s.addressClient.Get(context.Background(), &addresspb.GetAddressRequest{
            Id: person.AddressID.String(),
        })
        if err != nil {
            return nil, fmt.Errorf("failed to get address: %w", err)
        }
        
        // Combinar datos
        return &FullPersonResponse{
            Person: person,
            Address: convertAddress(addressResp.Address),
        }, nil
    }
    
    return &FullPersonResponse{Person: person}, nil
}
```

---

## 🗄️ Migraciones de base de datos

### Estructura de archivos
Guardá tus scripts en `migrations/`:
```
migrations/
├─ 0001_initial.up.sql      # crea tablas o colecciones
├─ 0001_initial.down.sql    # revierte lo anterior
├─ 0002_add_index.up.sql
└─ 0002_add_index.down.sql
```

### Configuración
Crea `cmd/migrate/main.go`:

```go
package main

import (
    "log"
    "os"

    "github.com/rem-gestion/rem-common/config"
    "github.com/rem-gestion/rem-common/db"
)

func main() {
    cfg := config.Load()
    dir := "up"
    if len(os.Args) > 1 { 
        dir = os.Args[1] 
    }
    if err := db.RunMigrations(cfg.DriverRelacional, cfg, "./migrations", dir); err != nil {
        log.Fatal(err)
    }
}
```

### Comandos de uso
```bash
# Aplica todas las migraciones pendientes
go run ./cmd/migrate            # equivale a "up"

# Revierte todas (¡cuidado!)
go run ./cmd/migrate down       

# Retrocede exactamente una migración
go run ./cmd/migrate "steps -1" 

# Avanza dos migraciones
go run ./cmd/migrate "steps 2"
```

Soporta **Postgres**, **MySQL** y **MongoDB**.

---

## 🛡️ Middlewares disponibles

| Middleware              | Función                                          | Uso                                    |
|-------------------------|--------------------------------------------------|----------------------------------------|
| `APIKeyAuth`            | Valida clave API en header                       | `mw.APIKeyAuth("X-Api-Key", apiKey)`   |
| `UserAuth`              | Valida JWT de usuario                            | `mw.UserAuth(jwtSecret)`               |
| `AdminAuth`             | Valida JWT con rol admin                         | `mw.AdminAuth(jwtSecret)`              |
| `EmployeeAuth`          | Valida JWT con rol empleado                      | `mw.EmployeeAuth(jwtSecret)`           |
| `MembershipAuth`        | Valida JWT con membresía                         | `mw.MembershipAuth(jwtSecret)`         |
| `GinLogger`             | Logging estructurado con Zap                    | `mw.GinLogger(zapLogger)`              |
| `RequestID`             | Genera ID único por request                      | `mw.RequestID()`                       |
| `RecoveryWithZap`       | Recovery con logging estructurado               | `mw.RecoveryWithZap(zapLogger)`        |
| `ErrorHandler`          | Mapea errores tipificados a códigos HTTP        | `mw.ErrorHandler()`                    |
| `RateLimit`             | Limitador de velocidad por IP                   | `mw.RateLimit(requestsPerSecond)`      |

---

## 🎯 Tipos de error disponibles

| Tipo de Error            | Código HTTP | Uso                                      |
|--------------------------|-------------|------------------------------------------|
| `ValidationError`        | 422         | Errores de validación de campos         |
| `BadRequestError`        | 400         | Solicitud malformada                     |
| `UnauthorizedError`      | 401         | Sin autenticación                        |
| `ForbiddenError`         | 403         | Sin permisos                             |
| `NotFoundError`          | 404         | Recurso no encontrado                    |
| `ConflictError`          | 409         | Conflicto (ej: email duplicado)          |
| `TooManyRequestsError`   | 429         | Rate limit excedido                      |
| `InternalServerError`    | 500         | Error interno del servidor               |

### Ejemplo de uso:
```go
import "github.com/rem-gestion/rem-common/errors"

func GetUser(id string) error {
    if id == "" {
        return &errors.BadRequestError{Msg: "ID is required"}
    }
    
    user, err := userRepo.FindByID(id)
    if err != nil {
        return &errors.NotFoundError{Msg: "User not found"}
    }
    
    return nil
}
```

---

## 🔐 Utilidades criptográficas

### Hashing de contraseñas (bcrypt)
```go
import "github.com/rem-gestion/rem-common/utils"

// Hash de contraseña
hash, err := utils.HashPassword("mi-password-secreto")

// Verificación de contraseña
err := utils.CheckPassword(hash, "mi-password-secreto")
if err != nil {
    // Contraseña incorrecta
}
```

### JWT Tokens
```go
import (
    "time"
    "github.com/rem-gestion/rem-common/utils"
)

// Generar token
token, err := utils.GenerateToken(
    "user-123",           // UserID
    "admin",              // Role (opcional)
    "jwt-secret-key",     // Secret
    time.Hour * 24,       // Duración
)

// Parsear token
claims, err := utils.ParseToken(token, "jwt-secret-key")
if err == nil {
    userID := claims.UserID
    role := claims.Role
}
```

---

## 🎯 Ventajas de usar rem-common

✅ **DRY**: Cero copy-paste de conexiones y middlewares entre servicios.

✅ **Consistente**: Mismo logger, manejo de errores y seguridad en toda la plataforma.

✅ **gRPC Ready**: Servidores y clientes gRPC preconfigurados con logging y health checks.

✅ **Protocol Buffers centralizados**: Definiciones compartidas para address y person services.

✅ **Extensible**: Agregás features (p.ej. tracing) en un solo lugar y todos los servicios lo heredan.

✅ **Listo para producción**: Migraciones versionadas, rate-limit, request-ID, JWT, Zap-JSON logging.

✅ **Multi-base de datos**: Soporte nativo para PostgreSQL, MySQL y MongoDB.

✅ **Microservicios**: Diseñado específicamente para arquitecturas distribuidas con comunicación HTTP y gRPC.

✅ **Versionado de APIs**: Soporte para versionado de protobuf (v1, v2, etc.).

---

## 💡 Tip para desarrollo

Agregá `rem-common` en tu `go.work` para desarrollo en monorepo y tendrás autocompletado y rebuild instantáneo con Air:

```bash
# En la raíz del proyecto
go work init
go work use ./rem-common
go work use ./auth-identity-svc
go work use ./person-svc
go work use ./address-svc
# ... otros servicios
```

### 🔄 Regenerar protobuf (si modificas .proto)

```bash
# Instalar dependencias
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Regenerar address service
protoc --go_out=. --go-grpc_out=. protos/address/v1/address.proto

# Regenerar person service  
protoc --go_out=. --go-grpc_out=. protos/person/v1/person.proto

# O usar script de automatización (si existe)
./scripts/generate-protos.sh
```

### 🧪 Testing gRPC

```bash
# Instalar grpcurl para testing manual
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Listar servicios disponibles
grpcurl -plaintext localhost:50051 list

# Llamar método específico
grpcurl -plaintext -d '{"id":"123"}' localhost:50051 address.v1.AddressService/Get

# Health check
grpcurl -plaintext localhost:50051 grpc.health.v1.Health/Check
```

---

## 📚 Documentación adicional

- [📝 Guía de migraciones](./migraciones.md) - Tutorial completo del sistema de migraciones
- [🔧 Variables de entorno](./migraciones.md#1-variables-de-entorno-necesarias) - Lista completa de configuraciones
- [📡 Protocol Buffers](./protos/) - Definiciones de servicios address y person
- [🏗️ gRPC Utilities](./grpc/) - Utilidades para servidor y cliente gRPC

---

## 🚀 Roadmap de mejoras

### ✅ Implementado:
- [x] Configuración centralizada con soporte .env
- [x] Logger estructurado con Zap
- [x] Conexiones multi-base de datos (PostgreSQL, MySQL, MongoDB)
- [x] Cache con Redis
- [x] Broker con RabbitMQ
- [x] Middlewares completos para Gin
- [x] Sistema de migraciones
- [x] Utilidades JWT y bcrypt
- [x] Manejo de errores tipificados
- [x] **Servidor gRPC con interceptors**
- [x] **Cliente gRPC optimizado**
- [x] **Protocol Buffers para address y person services**
- [x] **Health checks para gRPC**
- [x] **Reflection para debugging**

### 🔄 Próximas mejoras:
- [ ] **TLS/SSL para gRPC en producción**
- [ ] **OpenTelemetry tracing integrado**
- [ ] **Circuit breaker para clientes gRPC**
- [ ] **Rate limiting para gRPC**
- [ ] **Métricas Prometheus automáticas**
- [ ] **Configuración avanzada de keep-alive**
- [ ] **Retry policies configurable**
- [ ] **Load balancing para clientes gRPC**
- [ ] **Compresión gRPC automática**
- [ ] **Validación automática de protobuf**

---

**Repositorio:** [rem-backend/services/rem-common](.)  
**Protobuf definitions:** [./protos/](./protos/)  
**Versión:** v2.0.0 (con soporte gRPC)  
**Última actualización:** Diciembre 2024

---

## 📡 Comunicación gRPC

### 🔧 Utilidades para servidor gRPC (`grpc` package)

rem-common provee utilidades para crear servidores gRPC con funcionalidades avanzadas:

```go
import remgrpc "github.com/rem-gestion/rem-common/grpc"

// Crear servidor con interceptors de logging y health checks
server := remgrpc.NewServer(logger, keepaliveParams)

// Incluye automáticamente:
// ✅ Logging interceptor (unary y stream)
// ✅ Health service registrado
// ✅ Reflection habilitado (útil para grpcurl)
// ✅ Recovery interceptor para evitar panics
```

### 🔗 Cliente gRPC (`grpc.Dial`)

```go
import remgrpc "github.com/rem-gestion/rem-common/grpc"

// Conectar con configuración optimizada
conn, err := remgrpc.Dial("localhost:50051")

// Incluye automáticamente:
// ✅ Keep-alive configurado (30s time, 5s timeout)
// ✅ Timeout de conexión (5s)
// ✅ Credenciales inseguras (solo para desarrollo)
// ⚠️  TODO: Soporte TLS para producción
```

### 📋 Protocol Buffers definidos

#### Address Service (`protos/address/v1/`)
```protobuf
service AddressService {
  rpc Create (CreateAddressRequest) returns (CreateAddressResponse);
  rpc Get    (GetAddressRequest)    returns (GetAddressResponse);
}

message Address {
  string id = 1;
  string floor = 2;
  string unit = 3;
  string street = 4;
  int32 number = 5;
  string city = 6;
  string state = 7;
  string zip = 8;
  string country = 9;
}
```

#### Person Service (`protos/person/v1/`)
```protobuf
service PersonService {
  // Personas
  rpc CreatePerson (CreatePersonRequest) returns (PersonResponse);
  rpc GetPerson    (GetPersonRequest)    returns (PersonResponse);
  rpc ListPersons  (ListPersonsRequest)  returns (ListPersonsResponse);
  rpc UpdatePerson (UpdatePersonRequest) returns (PersonResponse);
  rpc DeletePerson (DeletePersonRequest) returns (google.protobuf.Empty);
  
  // Contactos
  rpc AddContact    (AddContactRequest)    returns (ContactResponse);
  rpc UpdateContact (UpdateContactRequest) returns (ContactResponse);
  rpc DeleteContact (DeleteContactRequest) returns (google.protobuf.Empty);
  rpc ListContacts  (ListContactsRequest)  returns (ListContactsResponse);
}
```

### 📦 Importar protobuf generados

```go
// Address service
import addresspb "github.com/rem-gestion/rem-common/protos/address/v1"

// Person service  
import personpb "github.com/rem-gestion/rem-common/protos/person/v1"

// Uso
client := addresspb.NewAddressServiceClient(conn)
response, err := client.Get(ctx, &addresspb.GetAddressRequest{Id: "123"})
```

### ⚙️ Variables de entorno gRPC

```env
# Servidor gRPC (donde escucha tu servicio)
REM_GRPC_HOST=0.0.0.0
REM_GRPC_PORT=50051

# Cliente gRPC (donde está address-svc)
REM_ADDRESS_HOST=localhost
REM_ADDRESS_PORT=50052
```
