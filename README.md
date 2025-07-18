# 🏢 REM - Real Estate Management Platform

**REM** es una plataforma completa de gestión inmobiliaria construida con arquitectura de microservicios, diseñada para proporcionar una solución escalable, robusta y moderna para la gestión integral de propiedades, personas, direcciones, amenities y operaciones inmobiliarias.

---

## 📖 Tabla de Contenidos

- [🎯 Descripción General](#-descripción-general)
- [🏗️ Arquitectura](#️-arquitectura)
- [🚀 Inicio Rápido](#-inicio-rápido)
- [📋 Prerrequisitos](#-prerrequisitos)
- [⚡ Scripts de Desarrollo](#-scripts-de-desarrollo)
- [🔄 Migraciones y Datos](#-migraciones-y-datos)
- [🌐 API Gateway](#-api-gateway)
- [📡 Testing con Postman](#-testing-con-postman)
- [🧩 Microservicios](#-microservicios)
- [📚 Rem-Common](#-rem-common)
- [📁 Estructura del Proyecto](#-estructura-del-proyecto)
- [🔧 Configuración](#-configuración)
- [🚧 Desarrollo Futuro](#-desarrollo-futuro)
- [🤝 Contribución](#-contribución)

---

## 🎯 Descripción General

REM es una plataforma de gestión inmobiliaria que centraliza:

- **👥 Gestión de Personas**: Individuos y empresas (físicas y jurídicas)
- **🏠 Gestión de Direcciones**: Sistema centralizado de ubicaciones
- **🏢 Gestión de Propiedades**: Inmuebles con amenities y características
- **🔐 Autenticación e Identidad**: Sistema seguro de usuarios y permisos
- **🌐 API Gateway**: Punto único de entrada para todos los servicios
- **📊 Base de Datos Compartida**: Arquitectura optimizada para desarrollo

### Características Principales

- ✅ **Arquitectura de Microservicios** con comunicación gRPC
- ✅ **API Gateway centralizado** con nginx
- ✅ **Base de datos compartida** para desarrollo optimizado
- ✅ **Resiliencia automática** con fallback a memoria
- ✅ **Scripts automatizados** para desarrollo y deployment
- ✅ **Colecciones Postman completas** para testing
- ✅ **Población automática** de datos de prueba (seed)
- ✅ **Hot-reload** en desarrollo con Air
- ✅ **Migraciones automáticas** de base de datos

---

## 🏗️ Arquitectura

### Vista General

```
┌─────────────────────────────────────────────────────────────┐
│                    API Gateway (nginx)                      │
│                  http://localhost:8081                      │
└───────────────────────────┬─────────────────────────────────┘
                            │HTTP
       ┌────────────────────┼────────────────────┐
       │                    │                    │
  ┌────▼────┐          ┌────▼────┐          ┌────▼────┐
  │  Auth   │   gRPC   │ Person  │   gRPC   │Address  │
  │ Service │ <──────> │ Service │ <──────> │Service  │
  │ :4002   │          │ :4001   │          │ :4000   │
  └─────────┘          └─────────┘          └─────────┘
       │                    │                    │
       │               ┌────▼────┐               │
       │               │Property │               │
       │               │ Service │ <─────────────┤
       │               │ :4004   │ gRPC          │
       │               └─────────┘               │
       │                    │                    │
       └────────────────────┼────────────────────┘
                            │
                      ┌─────▼─────┐
                      │PostgreSQL │
                      │    DB     │
                      │rem_develop│
                      └───────────┘
```

### Tecnologías

- **Backend**: Go con Gin Framework
- **Base de Datos**: PostgreSQL con GORM
- **API Gateway**: nginx
- **Comunicación**: HTTP REST + gRPC
- **Autenticación**: JWT
- **Containerización**: Docker & Docker Compose
- **Hot Reload**: Air
- **Migraciones**: golang-migrate

---

## 🚀 Inicio Rápido

### 1. **Clonar y Configurar**

```bash
git clone <repository-url>
cd rem-backend/services
```

### 2. **Iniciar Entorno de Desarrollo**

```bash
# Windows
.\scripts\dev.bat

# Linux/Mac  
./scripts/dev.sh
```

El script automáticamente:
- ✅ Verifica Docker
- ✅ Levanta infraestructura (DB + API Gateway)
- ✅ Ejecuta migraciones en todos los servicios
- ✅ **Pregunta si poblar la base con datos de prueba**
- ✅ Instala dependencias (Air)
- ✅ Abre terminales con hot-reload para cada servicio

### 3. **Verificar Funcionamiento**

```bash
# Health check del API Gateway
curl http://localhost:8081/health

# Health check de servicios
curl http://localhost:8081/api/health/auth
curl http://localhost:8081/api/health/person  
curl http://localhost:8081/api/health/address
curl http://localhost:8081/api/health/property
```

### 4. **URLs Disponibles**

- **🌐 API Gateway**: http://localhost:8081
- **🔐 Auth Service**: http://localhost:4002
- **👥 Person Service**: http://localhost:4001  
- **🏠 Address Service**: http://localhost:4000
- **🏢 Property Service**: http://localhost:4004

---

## 📋 Prerrequisitos

### Software Requerido

- **Docker & Docker Compose** - Para infraestructura
- **Go 1.21+** - Para compilar servicios
- **Git** - Para versionado
- **Make** (opcional) - Para comandos adicionales

### Puertos Utilizados

| Puerto | Servicio | Descripción |
|--------|----------|-------------|
| `8081` | nginx | API Gateway |
| `4000` | address-svc | Servicio de Direcciones |
| `4001` | person-svc | Servicio de Personas |
| `4002` | auth-identity-svc | Servicio de Autenticación |
| `4004` | property-svc | Servicio de Propiedades |
| `5432` | PostgreSQL | Base de Datos |
| `6379` | Redis | Cache (opcional) |
| `5672` | RabbitMQ | Message Broker (opcional) |

### Verificación de Prerrequisitos

```bash
# Verificar Docker
docker --version
docker-compose --version

# Verificar Go
go version

# Verificar puertos disponibles
netstat -an | findstr ":8081"  # Windows
lsof -i :8081                  # Linux/Mac
```

---

## ⚡ Scripts de Desarrollo

### Scripts Principales

| Script | Propósito | Descripción |
|--------|-----------|-------------|
| `dev.bat/sh` | **🚀 Desarrollo** | Inicia entorno completo con población opcional |
| `clean.bat/sh` | **🧹 Limpieza** | Detiene servicios y limpia archivos temporales |
| `seed-db.bat/sh` | **📊 Datos** | Puebla la base con datos de prueba |
| `test-gateway.bat` | **🧪 Testing** | Prueba endpoints del API Gateway |

### Uso de Scripts

#### **Desarrollo Normal**
```bash
# Inicio completo
.\scripts\dev.bat

# Limpieza después del trabajo
.\scripts\clean.bat
```

#### **Solo Poblado de Datos**
```bash
# Si ya tienes los servicios corriendo
.\scripts\seed-db.bat
```

#### **Testing Rápido**
```bash
# Prueba automática de endpoints
.\scripts\test-gateway.bat
```

### Características de los Scripts

- ✅ **Verificación automática** de prerrequisitos
- ✅ **Espera inteligente** para que servicios estén listos
- ✅ **Manejo de errores** con mensajes claros
- ✅ **Población interactiva** de datos de prueba
- ✅ **Hot-reload automático** con Air
- ✅ **Múltiples terminales** para desarrollo paralelo

---

## 🔄 Migraciones y Datos

### Sistema de Migraciones

Cada servicio tiene su propio sistema de migraciones:

```bash
# Ejecutar migraciones manualmente
cd auth-identity-svc && go run ./cmd/migrate
cd person-svc && go run ./cmd/migrate  
cd address-svc && go run ./cmd/migrate
cd property-svc && go run ./cmd/migrate
```

### Estructura de Migraciones

```
service-name/migrations/pg/
├── 0001_initial.up.sql      # Crear tablas
├── 0001_initial.down.sql    # Rollback
├── 0002_indexes.up.sql      # Agregar índices
└── 0002_indexes.down.sql    # Remover índices
```

### Datos de Prueba (Seed)

El sistema incluye datos de prueba listos para usar:

#### **👥 Usuarios de Prueba**
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

#### **🏢 Datos Incluidos**
- **Accounts & Users** - 8 usuarios con diferentes estados
- **Persons** - 5 individuos + 3 empresas
- **Addresses** - 10 direcciones en diferentes ciudades argentinas
- **Contacts** - Emails, teléfonos y contactos web
- **Relaciones** - Datos interconectados y consistentes

#### **📊 Comandos de Población**
```bash
# Población automática (incluida en dev.bat)
.\scripts\seed-db.bat

# Limpieza de datos + población
.\scripts\clean.bat
.\scripts\seed-db.bat
```

---

## 🌐 API Gateway

### Configuración

El API Gateway usa **nginx** para enrutar peticiones a los microservicios:

```nginx
# nginx.conf
location /api/auth/ {
    proxy_pass http://auth_backend/;
}
location /api/persons/ {
    proxy_pass http://person_backend/;
}
location /api/addresses/ {
    proxy_pass http://address_backend/;
}
location /api/properties/ {
    proxy_pass http://property_backend/;
}
location /api/amenities/ {
    proxy_pass http://property_backend/;
}
```

### Rutas Disponibles

| Gateway URL | Servicio Destino | Descripción |
|-------------|------------------|-------------|
| `/api/auth/*` | auth-identity-svc | Autenticación y registro |
| `/api/users/*` | auth-identity-svc | Gestión de usuarios |
| `/api/persons/*` | person-svc | Gestión de personas |
| `/api/contacts/*` | person-svc | Gestión de contactos |
| `/api/addresses/*` | address-svc | Gestión de direcciones |
| `/api/properties/*` | property-svc | Gestión de propiedades |
| `/api/amenities/*` | property-svc | Gestión de amenities |
| `/health` | nginx | Estado del gateway |
| `/api/health/*` | servicios | Health checks individuales |

### Ejemplos de Uso

```bash
# A través del API Gateway (RECOMENDADO)
curl http://localhost:8081/api/persons/
curl http://localhost:8081/api/auth/login

# Acceso directo a servicios (solo para debug)
curl http://localhost:4001/
curl http://localhost:4002/login
```

### Beneficios del API Gateway

- 🎯 **Punto único de entrada** para clientes
- 🔒 **Centralización de seguridad** y autenticación
- 📊 **Logging unificado** de todas las peticiones
- ⚡ **Load balancing** futuro entre instancias
- 🌐 **CORS** y headers manejados centralmente

---

## 📡 Testing con Postman

### ✨ Archivos Completamente Configurados

- **`REM-API-Collection.postman_collection.json`** - Colección completa con scripts inteligentes
- **`REM-Development.postman_environment.json`** - **60+ variables** con datos mock del seed

### 🚀 Importar en Postman

1. **Abrir Postman**
2. **Import** → Arrastrar ambos archivos JSON
3. **Seleccionar environment** "REM Development Environment"
4. **¡Listo para testing inmediato con datos mock!**

### 🎯 Características Automáticas Avanzadas

#### ✅ **Auto-inyección de Headers**
- **X-Api-Key** se agrega automáticamente a endpoints que lo necesitan
- **Authorization Bearer** se agrega automáticamente a endpoints autenticados
- **Logs informativos** en consola confirman qué headers se agregaron

#### ✅ **Auto-guardado Inteligente de IDs**
- **access_token** y **refresh_token** desde login
- **person_id**, **address_id**, **contact_id** desde respuestas
- **user_id** desde endpoint /me
- **Detección automática** de listas vs objetos individuales

#### ✅ **Auto-reemplazo de Variables**
- Request bodies usan variables como `{{person_id}}`, `{{test_email}}`
- **Reemplazo automático** en requests POST y PUT
- **Datos mock precargados** extraídos del script de seed

### 📋 Datos Mock Precargados (Listos para Usar)

#### 👤 **Personas Individuales**
- **Juan Pérez**: `person_id`, `address_id`, `contact_id` 
- **María González**: `maria_person_id`, `maria_address_id`
- **Carlos Rodríguez**: `carlos_person_id` (estado PENDING)
- **Ana Martínez**: `ana_person_id`, `ana_address_id`

#### 🏢 **Empresas**  
- **ACME Corporation**: `acme_person_id`, `acme_cuit`
- **Tech Solutions**: `tech_person_id`, `tech_email`
- **Innova Tech**: `innovatech_person_id`

#### � **Credenciales**
- **api_key**: `supersecret-api-key-for-dev`
- **test_email**: `juan.perez@example.com`
- **test_password**: `password` (todos los usuarios)

### 🎯 Workflow de Testing Inmediato

1. **Ejecutar "Login"** → `{{access_token}}` se guarda automáticamente
2. **Ejecutar "List Persons"** → `{{person_id}}` se guarda automáticamente  
3. **Ejecutar "Get Person by ID"** → Usa `{{person_id}}` automáticamente
4. **Cambiar a datos específicos** → Usar `{{maria_person_id}}`, `{{acme_person_id}}`, etc.

### 📚 Documentación Adicional

- **[📖 Guía Completa de Postman](./docs/README-Postman.md)** - Instrucciones detalladas
- **[📋 Referencia Rápida](./docs/Postman-Quick-Reference.md)** - Todos los IDs y datos mock
- **[📊 Resumen de Funcionalidades](./docs/Postman-Final-Summary.md)** - Lista completa de mejoras
- **[💡 Ejemplos de Payloads](./docs/Postman-Examples.md)** - Casos de uso específicos

### 🔧 Scripts Disponibles

```bash
# Ver funcionalidades avanzadas de Postman
.\scripts\test-postman-features.bat

# Probar endpoints del API Gateway  
.\scripts\test-gateway.bat
```

### 🌟 Estructura de la Colección

#### 🔐 **Authentication & Users**
- Register, Login, Refresh Token (con auto-save de tokens)
- Profile management (/me endpoint)
- User administration

#### 👥 **Persons & Companies**  
- CRUD completo con auto-save de IDs
- Individuos vs empresas (datos específicos)
- Gestión de contactos múltiples
- Operaciones bulk

#### 🏠 **Addresses**
- CRUD con datos de direcciones reales de Argentina
- Validación de campos

#### 🏢 **Properties & Real Estate**
- CRUD completo de propiedades inmobiliarias
- Gestión de amenities/comodidades por categorías
- Relaciones propiedades-amenities
- Filtros por tipo de propiedad y características

#### 🩺 **Health Checks**
- Monitoreo de API Gateway y servicios individuales

#### 🔧 **Direct Service Access**
- Endpoints directos para debugging (bypassing gateway)

---

## 🧩 Microservicios

### 🔐 Auth Identity Service

**Puerto**: 4002 | **[📖 Documentación](./auth-identity-svc/README.md)**

**Responsabilidades**:
- Registro y autenticación de usuarios
- Gestión de sesiones con JWT
- Perfiles de usuario y onboarding
- Integración con person-svc via gRPC

**Endpoints principales**:
```bash
POST /register    # Registro con datos opcionales de persona
POST /login       # Autenticación
GET  /me          # Perfil actual
PUT  /profile     # Actualizar perfil
GET  /users       # Listar usuarios (admin)
```

### 👥 Person Service

**Puerto**: 4001 | **[📖 Documentación](./person-svc/README.md)**

**Responsabilidades**:
- Gestión de personas (individuos y empresas)
- Manejo de contactos (email, teléfono, web)
- Integración con address-svc para direcciones
- Operaciones bulk y búsquedas avanzadas

**Tipos de entidad**:
- **Individual**: first_name, last_name, dni, sexo
- **Company**: legal_name, cuit, society_type

**Endpoints principales**:
```bash
GET  /persons          # Listar con filtros
POST /persons          # Crear (individual/empresa)
GET  /persons/:id/full # Con dirección expandida
POST /persons/bulk     # Creación masiva
POST /persons/:id/contacts # Agregar contacto
```

### 🏠 Address Service

**Puerto**: 4000 | **[📖 Documentación](./address-svc/README.md)**

**Responsabilidades**:
- Gestión centralizada de direcciones
- Deduplicación automática de direcciones
- Resiliencia con fallback a memoria
- Inmutabilidad de direcciones

**Características especiales**:
- 🛡️ **Resiliencia automática**: Funciona sin DB
- 🔄 **Reconexión inteligente**: Recovery automático
- 🚫 **Inmutabilidad**: No se pueden modificar direcciones
- ✨ **Deduplicación**: Evita duplicados automáticamente

**Endpoints principales**:
```bash
POST /addresses     # Crear (con deduplicación)
GET  /addresses/:id # Obtener dirección
PUT  /addresses/:id # ❌ No permitido (inmutable)
GET  /health/detailed # Estado del servicio
```

### 🏢 Property Service

**Puerto**: 4004 | **[📖 Documentación](./property-svc/README.md)**

**Responsabilidades**:
- Gestión completa de propiedades inmobiliarias
- Manejo de amenities/comodidades con categorías
- Relaciones propiedades-amenities
- Validación con address-svc y person-svc

**Tipos de propiedad**:
- APARTMENT, HOUSE, COMMERCIAL_SPACE, OFFICE, LAND, INDUSTRIAL_WAREHOUSE

**Categorías de amenities**:
- Security, Recreation, Services, Transport, Healthcare, Education

**Endpoints principales**:
```bash
GET  /properties              # Listar propiedades
POST /properties              # Crear propiedad
GET  /properties/:id          # Obtener propiedad específica
PUT  /properties/:id          # Actualizar propiedad
DELETE /properties/:id        # Eliminar propiedad
GET  /amenities               # Listar amenities
POST /amenities               # Crear amenity
POST /manage/:property_id/amenities/:amenity_id  # Asociar amenity
DELETE /manage/:property_id/amenities/:amenity_id # Desasociar amenity
```

**Características especiales**:
- 🏗️ **Códigos internos únicos**: Generación automática de códigos de propiedad
- 🔗 **Integración externa**: Validación con address-svc y person-svc via gRPC
- 🎯 **Mock services**: Servicios mock para desarrollo independiente
- ✨ **Gestión de relaciones**: Sistema flexible de amenities por propiedad

---

## 📚 Rem-Common

**[📖 Documentación Completa](./rem-common/README.md)**

Librería compartida que centraliza funcionalidades comunes:

### Módulos Incluidos

| Módulo | Funcionalidad |
|--------|---------------|
| **config** | Carga de variables de entorno y configuración |
| **logger** | Logging estructurado con Zap |
| **db** | Conexiones a PostgreSQL, MySQL, MongoDB |
| **grpc** | Servidor y cliente gRPC con interceptors |
| **protos** | Definiciones Protocol Buffers |
| **middleware** | Gin middlewares (auth, logging, recovery) |
| **errors** | Tipos de error tipificados |
| **utils** | Utilidades (bcrypt, JWT, validaciones) |
| **broker** | Conexión a RabbitMQ |
| **cache** | Conexión a Redis |

### Beneficios

- ✅ **DRY**: No repetir código entre servicios
- ✅ **Consistencia**: Mismo comportamiento en todos lados  
- ✅ **Mantenibilidad**: Cambios centralizados
- ✅ **Escalabilidad**: Facilita agregar nuevos servicios

### Ejemplo de Uso

```go
import (
    "github.com/rem-gestion/rem-common/config"
    "github.com/rem-gestion/rem-common/logger"
    "github.com/rem-gestion/rem-common/db"
)

func main() {
    cfg := config.Load()
    lg := logger.NewZap(cfg.Logger.Level)
    database := db.NewPostgres(cfg.Postgres, lg)
    // ...usar en el servicio
}
```

---

## 📁 Estructura del Proyecto

```
services/
├── 📋 README.md                    # 👈 Este archivo
├── ⚙️  .env.development             # Variables de desarrollo
├── 🐳 dev-full-compose.yml         # Docker Compose completo
├── 🌐 nginx.conf                   # Configuración del API Gateway
├── 📦 go.work                      # Go workspace
│
├── 🔐 auth-identity-svc/           # Servicio de Autenticación
│   ├── 📋 README.md
│   ├── cmd/api/main.go
│   ├── cmd/migrate/main.go
│   ├── src/                        # Código fuente
│   └── migrations/pg/              # Migraciones SQL
│
├── 👥 person-svc/                  # Servicio de Personas
│   ├── 📋 README.md
│   ├── cmd/api/main.go
│   ├── cmd/migrate/main.go
│   ├── src/                        # Código fuente
│   └── migrations/pg/              # Migraciones SQL
│
├── 🏠 address-svc/                 # Servicio de Direcciones
│   ├── 📋 README.md
│   ├── cmd/api/main.go
│   ├── cmd/migrate/main.go
│   ├── src/                        # Código fuente
│   └── migrations/pg/              # Migraciones SQL
│
├── 📚 rem-common/                  # Librería Compartida
│   ├── 📋 README.md
│   ├── config/                     # Configuración
│   ├── logger/                     # Logging
│   ├── db/                         # Conexiones DB
│   ├── grpc/                       # Utilidades gRPC
│   ├── protos/                     # Protocol Buffers
│   ├── middleware/                 # Gin middlewares
│   ├── errors/                     # Tipos de error
│   └── utils/                      # Utilidades
│
├── ⚡ scripts/                     # Scripts de Automatización
│   ├── dev.bat / dev.sh           # Inicio de desarrollo
│   ├── clean.bat / clean.sh       # Limpieza
│   ├── seed-db.bat / seed-db.sh   # Población de datos
│   ├── seed-dev-db.sql            # Datos de prueba
│   └── test-gateway.bat           # Testing automático
│
├── 📡 REM-API-Collection.postman_collection.json     # Colección Postman
├── 📡 REM-Development.postman_environment.json   # Environment Postman
│
├── 📚 docs/                        # Documentación adicional
│   ├── README-API-Gateway.md
│   ├── README-Environment-Config.md
│   ├── README-Postman.md
│   └── Postman-Examples.md
│
└── 🚧 [servicios-futuros]/         # Preparados para desarrollo
    ├── organization-svc/
    ├── role-rbac-svc/
    ├── subscription-billing-svc/
    └── terms-compliance-svc/
```

---

## 🔧 Configuración

### Variables de Entorno

#### **🌍 Globales (services/.env.development)**
```env
REM_ENVIRONMENT=development
REM_POSTGRES_DB=rem_development   # BD compartida
```

#### **🔐 Por Servicio**
Cada servicio tiene su propio `.env.development`:
```env
# Ejemplo: auth-identity-svc/.env.development
REM_POSTGRES_HOST=localhost
REM_POSTGRES_PORT=5432
REM_POSTGRES_USER=user
REM_POSTGRES_PASSWORD=password
REM_POSTGRES_DBNAME=rem_development
REM_SERVER_PORT=4002
REM_API_KEY=supersecret-api-key-for-dev
```

#### **🚀 Producción**
```bash
# Crear archivos de producción desde plantillas
cp auth-identity-svc/.env.production.example auth-identity-svc/.env.production
cp person-svc/.env.production.example person-svc/.env.production
cp address-svc/.env.production.example address-svc/.env.production

# Editar con credenciales reales
vim auth-identity-svc/.env.production
```

### Docker Compose

#### **Desarrollo** - `dev-full-compose.yml`
- PostgreSQL compartida
- nginx API Gateway  
- Redis y RabbitMQ opcionales
- Volúmenes persistentes
- Health checks automáticos

#### **Infraestructura mínima** - `postgres-compose.yml`
- Solo PostgreSQL para testing manual

### Configuración de Go Workspace

El proyecto usa **Go Workspaces** para manejar múltiples módulos:

```go
// go.work
go 1.21

use (
    ./address-svc
    ./auth-identity-svc
    ./person-svc
    ./rem-common
)
```

Beneficios:
- ✅ **Desarrollo local** sin go.mod replace
- ✅ **Hot-reload** de cambios en rem-common
- ✅ **IDE** reconoce todos los módulos
- ✅ **Testing** integrado entre servicios

---

## 🚧 Desarrollo Futuro

REM está diseñado como una plataforma modular que se expandirá gradualmente. Los servicios actuales (auth, persons, addresses) forman la **base sólida** sobre la cual se construirán las siguientes funcionalidades:

### 🏢 **Organization Service** *(Próximo)*
Gestión de empresas, inmobiliarias y organizaciones
- Estructura organizacional jerárquica
- Membresías y afiliaciones
- Configuraciones por organización
- Integración con persons para empleados

### 🏠 **Property Service** *(Core Business)*
Gestión de propiedades inmobiliarias
- Tipos de propiedad (casa, departamento, terreno, comercial)
- Características y amenities
- Historial y valuaciones
- Integración con addresses y persons (propietarios)

### 📋 **Contract Service** *(Operaciones)*
Gestión de contratos inmobiliarios
- Contratos de venta, alquiler, administración
- Estados y workflows de contratos  
- Comisiones y términos comerciales
- Integración con properties y persons

### 🏛️ **Role & RBAC Service** *(Seguridad)*
Sistema avanzado de roles y permisos
- Roles organizacionales (admin, agente, vendedor)
- Permisos granulares por recurso
- Jerarquías de autorización
- Audit logs de permisos

### 💰 **Subscription & Billing Service** *(Monetización)*
Gestión de planes y facturación
- Planes de suscripción por organización
- Facturación automática y manual
- Métricas de uso y límites
- Integración con gateways de pago

### 📊 **Accounting Service** *(Contabilidad)*
Sistema contable integrado
- Movimientos contables automáticos
- Comisiones y liquidaciones
- Reportes financieros
- Integración con contratos y propiedades

### 🔗 **Integration Service** *(Conectividad)*
Hub de integraciones externas
- APIs de portales inmobiliarios
- Servicios de geolocalización
- Sistemas de comunicación
- Webhooks y notificaciones

### 📈 **Analytics Service** *(Business Intelligence)*
Análisis de datos y reportes
- KPIs inmobiliarios
- Dashboards ejecutivos
- Reportes automáticos
- Data warehouse

### 👥 **CRM Service** *(Gestión Comercial)*
Gestión de relaciones con clientes
- Leads y oportunidades
- Seguimiento comercial
- Automatización de marketing
- Integración con persons y contracts

### 💬 **Communication Service** *(Mensajería)*
Centro de comunicaciones
- Email transaccional
- SMS y WhatsApp
- Notificaciones push
- Templates y automatización

### 📜 **Terms & Compliance Service** *(Legal)*
Gestión legal y normativa
- Términos y condiciones
- Políticas de privacidad
- Compliance regulatorio
- Audit trails

### 💰 **Finance Service** *(Finanzas Avanzadas)*
Operaciones financieras complejas
- Créditos hipotecarios
- Inversiones inmobiliarias
- Análisis de rentabilidad
- Simuladores financieros

### 🤖 **Assistant Bot Service** *(IA)*
Asistente inteligente
- Chatbot con NLP
- Automatización de consultas
- Recomendaciones inteligentes
- Integración con todos los servicios

### 🛠️ **Operations Service** *(Gestión Operativa)*
Gestión de operaciones diarias
- Workflows de procesos
- Task management
- Calendarios y citas
- Gestión documental

---

### 🎯 **Roadmap de Desarrollo**

#### **Fase 1 - Fundación** *(Completada)*
- ✅ Auth Identity Service
- ✅ Person Service  
- ✅ Address Service
- ✅ rem-common library
- ✅ API Gateway
- ✅ Scripts de desarrollo
- ✅ Colecciones Postman

#### **Fase 2 - Core Business** *(Q1 2025)*
- 🏢 Organization Service
- 🏠 Property Service (MVP)
- 📋 Contract Service (básico)
- 🏛️ Role & RBAC Service

#### **Fase 3 - Operaciones** *(Q2 2025)*
- 💰 Subscription & Billing Service
- 📊 Accounting Service (básico)
- 👥 CRM Service (MVP)
- 📈 Analytics Service (básico)

#### **Fase 4 - Integración** *(Q3 2025)*
- 🔗 Integration Service
- 💬 Communication Service
- 📜 Terms & Compliance Service
- 🛠️ Operations Service

#### **Fase 5 - Avanzado** *(Q4 2025)*
- 💰 Finance Service
- 🤖 Assistant Bot Service
- 📈 Analytics Service (avanzado)
- 📊 Accounting Service (completo)

---

### 🔮 **Visión Técnica Futura**

#### **Escalabilidad**
- **Kubernetes** deployment para producción
- **Event-driven architecture** con Apache Kafka
- **CQRS** para separar lecturas de escrituras
- **Database sharding** por organización

#### **Observabilidad**
- **Distributed tracing** con Jaeger
- **Metrics** con Prometheus + Grafana
- **Centralized logging** con ELK Stack
- **APM** con New Relic o Datadog

#### **Seguridad Avanzada**
- **OAuth2/OIDC** con proveedores externos
- **Multi-factor authentication** (MFA)
- **API rate limiting** avanzado
- **Encryption at rest** para datos sensibles

#### **Performance**
- **Redis clustering** para cache distribuido
- **CDN** para assets estáticos
- **Database read replicas** para consultas
- **GraphQL** BFF para clientes móviles

#### **DevOps**
- **CI/CD pipelines** con GitHub Actions
- **Infrastructure as Code** con Terraform
- **Automated testing** en múltiples niveles
- **Blue-green deployments** sin downtime

---

## 🤝 Contribución

### Principios de Desarrollo

1. **🏗️ Arquitectura Limpia**: Separación clara de responsabilidades
2. **📋 API First**: Documentación antes de implementación  
3. **🧪 Testing**: Unit tests + integration tests obligatorios
4. **📚 Documentación**: README actualizado con cada cambio
5. **🔒 Seguridad**: Security by design en cada feature
6. **⚡ Performance**: Optimización desde el diseño
7. **🌐 Escalabilidad**: Diseño para millones de usuarios

### Flujo de Desarrollo

1. **Fork** del repositorio
2. **Feature branch** desde `develop`
3. **Desarrollo** siguiendo estándares
4. **Tests** unitarios y de integración
5. **Pull Request** con descripción detallada
6. **Code Review** por el equipo
7. **Merge** después de aprobación

### Estándares de Código

- **Go**: Seguir `gofmt`, `golint`, `go vet`
- **SQL**: Naming conventions consistentes
- **API**: RESTful design + OpenAPI docs
- **Git**: Conventional commits
- **Documentación**: Markdown con ejemplos

### Testing

```bash
# Tests unitarios por servicio
cd auth-identity-svc && go test ./...
cd person-svc && go test ./...
cd address-svc && go test ./...

# Tests de integración
./scripts/test-integration.sh

# Tests de performance
./scripts/test-performance.sh
```

### Herramientas de Desarrollo

- **Air** - Hot reload en desarrollo
- **golangci-lint** - Linting avanzado  
- **Postman** - Testing de APIs
- **Docker** - Containerización consistente
- **Makefile** - Comandos de desarrollo

---

## 📞 Soporte y Contacto

### Documentación

- **📋 Este README** - Información general
- **📁 docs/** - Documentación específica
- **📖 Servicios/** - README de cada microservicio
- **🔧 rem-common/** - Documentación de librería compartida

### Enlaces Útiles

- [API Gateway](./docs/README-API-Gateway.md)
- [Variables de Entorno](./docs/README-Environment-Config.md)  
- [Postman Collections](./docs/README-Postman.md)
- [Ejemplos de Payloads](./docs/Postman-Examples.md)

### Issues y Bugs

Reportar issues en GitHub con:
- 🐛 **Bug reports** con steps to reproduce
- 💡 **Feature requests** con justificación
- 📖 **Documentation** mejoras y correcciones

---

## 📄 Licencia

Este proyecto es **privado** y está protegido por derechos de autor. El acceso está limitado al equipo de desarrollo autorizado.

---

## 🎉 Estado del Proyecto

```
🟢 READY FOR DEVELOPMENT
├── ✅ Microservicios base funcionando
├── ✅ API Gateway configurado  
├── ✅ Base de datos con seed data
├── ✅ Scripts de desarrollo automatizados
├── ✅ Colecciones Postman completas
├── ✅ Documentación completa
└── 🚀 Listo para agregar Property Service
```

**¡El proyecto está completamente funcional y listo para el desarrollo del core business!** 🎯

---

*Última actualización: 12 de Julio 2025 | Versión: 1.0.0-alpha*
