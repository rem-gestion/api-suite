# API Gateway para Desarrollo Local

## 🏗️ Arquitectura

Este proyecto usa nginx como API Gateway centralizado para todos los microservicios en desarrollo local.

```
┌─────────────────────────────────────────────────────────────┐
│                    API Gateway (nginx)                      │
│                  http://localhost:8081                      │
└─────────────────────┬───────────────────────────────────────┘
                      │
         ┌────────────┼────────────┐
         │            │            │
    ┌────▼────┐  ┌────▼────┐  ┌────▼────┐
    │  Auth   │  │ Person  │  │Address  │
    │ Service │  │ Service │  │Service  │
    │ :4002   │  │ :4001   │  │ :4000   │
    └─────────┘  └─────────┘  └─────────┘
```

## � Enrutamiento

### Rutas del API Gateway

| Gateway URL                   | Servicio Destino                 | Descripción                    |
|-------------------------------|-----------------------------------|--------------------------------|
| `/api/auth/*`                 | `auth-identity-svc:4002/`        | Autenticación                  |
| `/api/users/*`                | `auth-identity-svc:4002/users/`  | Gestión de usuarios           |
| `/api/persons/*`              | `person-svc:4001/`               | Gestión de personas           |
| `/api/contacts/*`             | `person-svc:4001/contacts/`      | Gestión de contactos          |
| `/api/addresses/*`            | `address-svc:4000/`              | Gestión de direcciones        |
| `/health`                     | API Gateway Health Check          | Estado del gateway            |
| `/api/health/{service}`       | Health check individual           | Estado de servicios           |
curl -H "X-Api-Key: supersecret-api-key-for-dev" http://localhost:8081/api/health/person  
curl -H "X-Api-Key: supersecret-api-key-for-dev" http://localhost:8081/api/health/address
```

### Operaciones CRUD
```bash
# Listar personas
curl -H "X-Api-Key: supersecret-api-key-for-dev" http://localhost:8081/api/persons/

# Crear persona
curl -X POST -H "X-Api-Key: supersecret-api-key-for-dev" \
  -H "Content-Type: application/json" \
  -d '{"nombre":"Juan","apellido":"Pérez"}' \
  http://localhost:8081/api/persons/

### Ejemplos de Uso

#### Autenticación
```bash
# Registro
curl -X POST http://localhost:8081/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"secret123"}'

# Login
curl -X POST http://localhost:8081/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"secret123"}'

# Perfil del usuario autenticado
curl -X GET http://localhost:8081/api/auth/me \
  -H "Authorization: Bearer <token>"
```

#### Gestión de Usuarios
```bash
# Listar usuarios (requiere auth)
curl -X GET http://localhost:8081/api/users/ \
  -H "Authorization: Bearer <token>"

# Buscar usuario por email
curl -X GET http://localhost:8081/api/users/email/test@example.com \
  -H "Authorization: Bearer <token>"
```

#### Gestión de Personas
```bash
# Crear persona
curl -X POST http://localhost:8081/api/persons/ \
  -H "Content-Type: application/json" \
  -d '{"nombre":"Juan","apellido":"Pérez"}'

# Listar personas
curl -X GET http://localhost:8081/api/persons/

# Obtener persona con dirección expandida
curl -X GET http://localhost:8081/api/persons/1/full

# Agregar contacto a persona
curl -X POST http://localhost:8081/api/persons/1/contacts \
  -H "Content-Type: application/json" \
  -d '{"tipo":"email","valor":"juan@example.com"}'
```

#### Gestión de Contactos
```bash
# Actualizar contacto
curl -X PUT http://localhost:8081/api/contacts/123 \
  -H "Content-Type: application/json" \
  -d '{"tipo":"telefono","valor":"555-1234"}'

# Eliminar contacto
curl -X DELETE http://localhost:8081/api/contacts/123

# Obtener contacto primario de persona
curl -X GET http://localhost:8081/api/contacts/primary/1
```

#### Gestión de Direcciones
```bash
# Crear dirección
curl -X POST http://localhost:8081/api/addresses/ \
  -H "Content-Type: application/json" \
  -d '{"calle":"Main St","numero":"123","ciudad":"Springfield"}'

# Obtener dirección
curl -X GET http://localhost:8081/api/addresses/1

# Actualizar dirección
curl -X PUT http://localhost:8081/api/addresses/1 \
  -H "Content-Type: application/json" \
  -d '{"calle":"Main Street","numero":"123","ciudad":"Springfield"}'
```

## 🚀 Uso del Entorno

### Iniciar Desarrollo
```bash
.\scripts\dev.bat
```

Este comando:
1. Levanta PostgreSQL y nginx en Docker
2. Ejecuta migraciones
3. Inicia los 3 microservicios con hot-reload (Air)

### Probar API Gateway
```bash
.\scripts\test-gateway.bat
```

### Limpiar Entorno
```bash
.\scripts\clean.bat
```

## 🛠️ Configuración

### Puertos de Servicios (Local)
- **API Gateway (nginx)**: `8081`
- **Auth Identity Service**: `4002`
- **Person Service**: `4001`
- **Address Service**: `4000`
- **PostgreSQL**: `5432`

### Variables de Entorno
Los servicios usan archivos `.env.development` individuales que heredan la configuración global de `services/.env.development`.

### Base de Datos
Todos los servicios comparten una única instancia de PostgreSQL en desarrollo:
- **Host**: `localhost`
- **Puerto**: `5432`
- **Base de datos**: `rem_development`
- **Usuario**: `lucho`
- **Password**: `supersecreta`

## 📝 Beneficios del API Gateway

1. **Centralización**: Una sola URL para todas las APIs
2. **Simplicidad**: Sin duplicación de prefijos en rutas
3. **CORS**: Configurado centralmente
4. **Desarrollo**: Fácil proxy a servicios locales
5. **Escalabilidad**: Preparado para load balancing futuro

## 🔧 Troubleshooting

### Error de Puerto Ocupado
Si el puerto 8081 está ocupado:
```bash
# Ver qué proceso usa el puerto
netstat -ano | findstr :8081

# Cambiar puerto en dev-full-compose.yml
# Editar: ports: - "8081:80"
```

### Servicio No Responde
Verificar que todos los servicios estén corriendo:
```bash
curl http://localhost:8081/api/health/auth
curl http://localhost:8081/api/health/person  
curl http://localhost:8081/api/health/address
```

### Logs de nginx
```bash
docker logs rem-api-gateway
```

### Logs de Base de Datos
```bash
docker logs rem-postgres-dev
```
