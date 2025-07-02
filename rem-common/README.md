# rem-common 📦  
Librería compartida para todos los microservicios de **REM Gestión** – centraliza
configuración, conexiones, logger, cache, broker, middlewares, utilidades
criptográficas y sistema de migraciones.

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
    // ...
)
```

---

## 📋 Contenido del módulo

| Carpeta / paquete          | Qué provee                                                                                           |
|----------------------------|-------------------------------------------------------------------------------------------------------|
| `config`                   | Carga variables `REM_*` (con soporte `.env`) en un struct único.                                      |
| `logger`                   | Instancia **Zap** con colores, caller y campo fijo `service`.                                         |
| `db`                       | Constructores `NewPostgres`, `NewMySQL`, `NewMongo`, `NewRedis`, `NewRabbit` + helper de migraciones. |
| `broker`                  | Helpers para RabbitMQ (DeclareExchange/Queue, Publish, Consume).                                      |
| `cache`                    | `NewRedis` con `Ping()`.                                                                              |
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
| `POSTGRES_DB`             | Base de datos de PostgreSQL               | `rem_auth`                   |
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
| `LOGGER_LEVEL`            | `debug`, `info`, `warn`, `error`          | `info`                       |
| `SERVER_PORT`             | Puerto del servidor                       | `8080`                       |
| `API_KEY`                 | Clave para `APIKeyAuth`                   | `your-secret-api-key`        |
| `JWT_SECRET`              | Secreto HS256 para tokens                 | `your-jwt-secret-key`        |

*(Ver [migraciones.md](./migraciones.md) para guía completa de migraciones)*

---

## 🔧 Ejemplo rápido de uso en un microservicio

```go
// services/auth-identity/config/builder/builder.go
package builder

import (
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"

    "github.com/rem-gestion/rem-common/config"
    "github.com/rem-gestion/rem-common/logger"
    "github.com/rem-gestion/rem-common/db"
    "github.com/rem-gestion/rem-common/cache"
    "github.com/rem-gestion/rem-common/broker"
    mw "github.com/rem-gestion/rem-common/middleware"
)

func Build() *gin.Engine {
    cfg := config.Load()

    lg  := logger.New(cfg.Logger, "auth-identity")
    pg, _ := db.NewPostgres(cfg.Postgres)
    rdb, _ := cache.NewRedis(cfg.Redis)
    rbConn, _ := broker.NewRabbit(cfg.Rabbit)
    _ = pg; _ = rdb; _ = rbConn     // inyectá donde corresponda

    r := gin.New()
    r.Use(
        mw.RequestID(),
        mw.APIKeyAuth("X-Api-Key", cfg.APIKey),
        mw.GinLogger(lg),
        mw.RecoveryWithZap(lg),
        mw.ErrorHandler(),
    )

    r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })
    return r
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

✅ **Extensible**: Agregás features (p.ej. tracing) en un solo lugar y todos los servicios lo heredan.

✅ **Listo para producción**: Migraciones versionadas, rate-limit, request-ID, JWT, Zap-JSON logging.

✅ **Multi-base de datos**: Soporte nativo para PostgreSQL, MySQL y MongoDB.

✅ **Microservicios**: Diseñado específicamente para arquitecturas distribuidas.

---

## 💡 Tip para desarrollo

Agregá `rem-common` en tu `go.work` para desarrollo en monorepo y tendrás autocompletado y rebuild instantáneo con Air:

```bash
# En la raíz del proyecto
go work init
go work use ./rem-common
go work use ./auth-identity
go work use ./profile
# ... otros servicios
```

---

## 📚 Documentación adicional

- [📝 Guía de migraciones](./migraciones.md) - Tutorial completo del sistema de migraciones
- [🔧 Variables de entorno](./migraciones.md#1-variables-de-entorno-necesarias) - Lista completa de configuraciones
