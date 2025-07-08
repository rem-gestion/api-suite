# Auth Identity Service

Microservicio de autenticación e identidad para el sistema REM.

## Funcionalidades

### Autenticación
- ✅ Registro de usuarios con email/contraseña
- ✅ Login con validación de credenciales
- ✅ Refresh tokens para renovación de sesión
- ✅ Validación de tokens JWT
- ✅ Logout (invalidación del lado cliente)
- 🚧 Recuperación de contraseña (preparado para integración con servicio de comunicación)

### Gestión de Usuarios
- ✅ Perfiles de usuario con datos de persona
- ✅ Actualización de perfiles
- ✅ Estados de onboarding (new, in_progress, done)
- ✅ Integración con person-svc via gRPC
- ✅ Soft delete de usuarios
- ✅ Lista de usuarios con filtros (admin)

### Estados de Cuenta
- **pending**: Cuenta creada pero no activada
- **active**: Cuenta activa y funcional
- **suspended**: Cuenta suspendida por administrador

### Estados de Onboarding
- **new**: Usuario recién registrado
- **in_progress**: Usuario con datos parciales
- **done**: Usuario con perfil completo

## API Endpoints

### Públicos (sin autenticación)
```
POST /api/v1/auth/register          - Registro de usuario
POST /api/v1/auth/login             - Login
POST /api/v1/auth/refresh           - Renovar token
POST /api/v1/auth/forgot-password   - Solicitar reset de contraseña
POST /api/v1/auth/reset-password    - Resetear contraseña
GET  /api/v1/auth/validate          - Validar token
```

### Protegidos (requieren JWT)
```
GET  /api/v1/auth/me                - Obtener perfil actual
POST /api/v1/auth/logout            - Logout

PUT  /api/v1/users/profile          - Actualizar perfil
GET  /api/v1/users                  - Listar usuarios (admin)
GET  /api/v1/users/email/{email}    - Buscar por email (admin)
POST /api/v1/users/{id}/deactivate  - Desactivar usuario (admin)
```

## Estructura de Datos

### Account
```go
type Account struct {
    ID           string          // UUID
    Provider     AccountProvider // email | google
    Email        string          // único
    PasswordHash string
    Status       AccountStatus   // pending | active | suspended
    CreatedAt    time.Time
    UpdatedAt    *time.Time
    UpdatedBy    *string
    DeletedAt    *time.Time
}
```

### User
```go
type User struct {
    ID            string        // UUID
    AccountID     string        // FK a accounts
    PersonID      *string       // FK a person-svc
    OnboardStatus OnboardStatus // new | in_progress | done
    LastLogin     *time.Time
    CreatedAt     time.Time
    UpdatedAt     *time.Time
    UpdatedBy     *string
    DeletedAt     *time.Time
}
```

## Integración con Person Service

El servicio se integra con `person-svc` via gRPC para:
- Crear personas durante el registro con datos opcionales
- Obtener datos de persona para el perfil de usuario
- Manejar direcciones a través de person-svc que se comunica con address-svc

### Flujo de Registro con Persona
1. Usuario envía datos de registro + datos opcionales de persona
2. Se crea Account y User en auth-identity-svc
3. Si hay datos de persona, se llama a person-svc via gRPC
4. Person-svc maneja la creación de dirección via address-svc
5. Se vincula el PersonID al User
6. Se retornan los tokens JWT

## Configuración

### Variables de Entorno Requeridas
```bash
# Base de datos
REM_POSTGRES_HOST=localhost
REM_POSTGRES_PORT=5432
REM_POSTGRES_USER=postgres
REM_POSTGRES_PASSWORD=password
REM_POSTGRES_DBNAME=auth_identity

# JWT
REM_API_KEY=your-secret-jwt-key

# Servidor
REM_SERVER_PORT=8080

# Person Service (gRPC)
REM_ADDRESS_HOST=127.0.0.1
REM_ADDRESS_PORT=50051
```

## Comandos

### Ejecutar migraciones
```bash
# Aplicar migraciones
go run cmd/migrate/main.go -direction=up

# Revertir a versión específica
go run cmd/migrate/main.go -direction=down -version=0001
```

### Ejecutar servicio
```bash
go run cmd/api/main.go
```

## Desarrollo

### Estructura del Proyecto
```
auth-identity-svc/
├── cmd/
│   ├── api/main.go          # Punto de entrada HTTP
│   └── migrate/main.go      # Herramienta de migraciones
├── src/
│   ├── controllers/         # Controladores HTTP
│   ├── dto/                # Data Transfer Objects
│   ├── models/             # Modelos de base de datos
│   ├── repository/         # Capa de acceso a datos
│   ├── router/             # Configuración de rutas
│   └── services/           # Lógica de negocio
├── migrations/pg/          # Migraciones de PostgreSQL
└── go.mod
```

### Testing

Endpoints de ejemplo:

```bash
# Registro
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "person_data": {
      "type": "individual",
      "individual": {
        "first_name": "Juan",
        "last_name": "Pérez",
        "dni": "12345678"
      }
    }
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'

# Obtener perfil
curl -X GET http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Dependencias

- **Gin**: Framework web HTTP
- **GORM**: ORM para PostgreSQL
- **JWT**: Autenticación basada en tokens
- **gRPC**: Comunicación con person-svc
- **bcrypt**: Hash de contraseñas
- **rem-common**: Librería común del proyecto

## TODO

- [ ] Implementar reset de contraseña completo (requiere servicio de comunicación)
- [ ] Agregar middleware de autorización por roles
- [ ] Implementar blacklist de tokens para logout real
- [ ] Agregar rate limiting en endpoints de auth
- [ ] Implementar autenticación OAuth2 (Google)
- [ ] Agregar métricas y monitoring
- [ ] Tests unitarios y de integración
- [ ] Documentación OpenAPI/Swagger
