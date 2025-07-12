# 📋 Postman Collections para REM API

Este directorio contiene las colecciones completas de Postman para probar todos los endpoints de la API REM a través del API Gateway.

## 📁 Archivos

- **`REM-API-Collection.postman_collection.json`** - Colección completa con todos los endpoints
- **`REM-Development.postman_environment.json`** - Environment con variables de desarrollo
- **`README-Postman.md`** - Este archivo con instrucciones

## 🚀 Instalación y Configuración

### 1. Importar en Postman

1. Abre Postman
2. Haz clic en **"Import"** en la parte superior izquierda
3. Arrastra y suelta los dos archivos JSON o usa **"Upload Files"**:
   - `REM-API-Collection.postman_collection.json`
   - `REM-Development.postman_environment.json`

### 2. Seleccionar Environment

1. En la esquina superior derecha de Postman, selecciona **"REM Development Environment"**
2. Verifica que las variables estén configuradas correctamente:
   - `base_url`: http://localhost:8081 (API Gateway)
   - `access_token`: (se llenará automáticamente al hacer login)
   - `refresh_token`: (se llenará automáticamente al hacer login)

### 3. Iniciar el entorno de desarrollo

Antes de usar las colecciones, asegúrate de que el entorno esté corriendo:

```bash
# En el directorio services/
.\scripts\dev.bat   # Windows
# o
./scripts/dev.sh    # Linux/Mac
```

## 🎯 Uso de las Colecciones

### Orden Recomendado de Pruebas

#### 1. **Health Checks** 🩺
Primero verifica que todos los servicios estén funcionando:
- API Gateway Health
- Auth Service Health
- Person Service Health
- Address Service Health

#### 2. **Authentication** 🔐
```
Login → (automáticamente guarda tokens) → Validate Token → Get My Profile
```

Los endpoints de login tienen scripts que automáticamente guardan los tokens en las variables de environment.

#### 3. **CRUD de Personas** 👥
```
List Persons → Create Person → Get Person → Update Person → Delete Person
```

#### 4. **Gestión de Contactos** 📱
```
List Contacts → Add Contact → Update Contact → Delete Contact
```

#### 5. **Gestión de Direcciones** 🏠
```
Create Address → Get Address → Update Address → Delete Address
```

### 🔑 Usuarios de Prueba (Datos Seed)

Los siguientes usuarios están disponibles en la base de datos de desarrollo:

```json
{
  "email": "juan.perez@example.com",
  "password": "password"
}

{
  "email": "maria.gonzalez@example.com", 
  "password": "password"
}

{
  "email": "admin@acmecorp.com",
  "password": "password"
}
```

## 📊 Estructura de la Colección

### 🔐 Authentication & Users
- **Register New User** - Registro de nuevos usuarios
- **Login** - Autenticación (guarda tokens automáticamente)
- **Refresh Token** - Renovación de tokens
- **Validate Token** - Validación de token actual
- **Get My Profile** - Información del perfil
- **Logout** - Cerrar sesión
- **Users Management** - Gestión de usuarios (admin)

### 👥 Persons & Companies
- **List Persons** - Listado con filtros y paginación
- **Get Person by ID** - Detalles de persona
- **Get Person with Address (Full)** - Persona con dirección expandida
- **Create Individual Person** - Crear persona individual
- **Create Company** - Crear empresa
- **Update Person** - Actualizar información
- **Delete Person** - Eliminar persona
- **Bulk Create Persons** - Creación masiva
- **Contacts Management** - Gestión de contactos

### 🏠 Addresses
- **Get Address by ID** - Obtener dirección
- **Create Address** - Crear nueva dirección
- **Update Address** - Actualizar dirección
- **Delete Address** - Eliminar dirección

### 🩺 Health Checks
- Health checks para todos los servicios a través del API Gateway

### 🔧 Direct Service Access (Debug)
- Acceso directo a servicios (para debugging, saltando el API Gateway)

## 🛠️ Variables de Environment

Las siguientes variables se configuran automáticamente:

| Variable | Descripción | Auto-Set |
|----------|-------------|----------|
| `base_url` | URL del API Gateway | ❌ |
| `access_token` | Token JWT de acceso | ✅ (login) |
| `refresh_token` | Token de renovación | ✅ (login) |
| `user_id` | ID del usuario actual | ✅ (login) |
| `person_id` | ID de persona para pruebas | ✅ (list/create) |
| `address_id` | ID de dirección para pruebas | ✅ (create) |
| `contact_id` | ID de contacto para pruebas | ✅ (list/create) |

## 🔄 Scripts Automáticos

Los siguientes endpoints tienen scripts que automáticamente guardan IDs para usar en otras peticiones:

- **Login** → guarda `access_token`, `refresh_token`, `user_id`, `person_id`
- **List Persons** → guarda `person_id` del primer resultado
- **Create Person** → guarda `person_id` de la nueva persona
- **Create Address** → guarda `address_id` de la nueva dirección
- **List Contacts** → guarda `contact_id` del primer contacto
- **Add Contact** → guarda `contact_id` del nuevo contacto

## 🚨 Troubleshooting

### Error 502 Bad Gateway
- Verifica que todos los servicios estén corriendo
- Usa los health checks para verificar el estado
- Revisa los logs en las terminales de los servicios

### Error 401 Unauthorized
- Haz login para obtener un token válido
- Verifica que `access_token` esté configurado en el environment
- Usa "Refresh Token" si el token expiró

### Error 404 Not Found
- Verifica que las URLs estén correctas
- Asegúrate de que el API Gateway esté corriendo en puerto 8081
- Verifica que los servicios estén configurados correctamente

### Variables no encontradas
- Asegúrate de tener seleccionado el environment "REM Development Environment"
- Ejecuta los endpoints que auto-configuran las variables (como Login o List Persons)

## 🎯 Endpoints Principales del API Gateway

| Ruta | Servicio | Descripción |
|------|----------|-------------|
| `/api/auth/*` | auth-identity-svc | Autenticación y autorización |
| `/api/users/*` | auth-identity-svc | Gestión de usuarios |
| `/api/persons/*` | person-svc | Gestión de personas |
| `/api/contacts/*` | person-svc | Gestión de contactos |
| `/api/addresses/*` | address-svc | Gestión de direcciones |
| `/health` | nginx | Health check del API Gateway |
| `/api/health/*` | services | Health checks individuales |

## 📝 Notas de Desarrollo

- Todos los endpoints usan el API Gateway como punto de entrada
- Los tokens JWT se manejan automáticamente
- Los scripts auto-guardan IDs importantes para facilitar las pruebas
- Se incluyen ejemplos de payloads realistas para todas las operaciones
- Los endpoints de debug permiten acceso directo a servicios para troubleshooting

---

¡Happy Testing! 🎉
