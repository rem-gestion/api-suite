# 📋 REM Platform - Quick Reference de Datos Mock

Esta es una referencia rápida de todos los IDs y datos disponibles en el environment de Postman, extraídos directamente del script de seed de desarrollo.

## 🔐 Credenciales de Acceso

| Variable | Valor | Descripción |
|----------|-------|-------------|
| `api_key` | `supersecret-api-key-for-dev` | API Key para todos los servicios |
| `test_password` | `password` | Password para todos los usuarios de prueba |

## 👤 Personas Individuales

### Juan Carlos Pérez (Usuario Principal)
- **Email login**: `juan.perez@example.com`
- **User ID**: `11111111-1111-1111-1111-111111111111`
- **Person ID**: `11111111-1111-1111-1111-111111111111`
- **Address ID**: `11111111-1111-1111-1111-111111111111` (Av. Corrientes 1234, Buenos Aires)
- **DNI**: `12345678`
- **Contactos**:
  - Email: `11111111-1111-1111-1111-111111111111` (juan.perez@personal.com)
  - Phone: `11111111-2222-1111-1111-111111111111` (+541123456789)
  - WhatsApp: `11111111-3333-1111-1111-111111111111` (+541123456789)

### María Elena González
- **Email login**: `maria.gonzalez@example.com`
- **User ID**: `22222222-2222-2222-2222-222222222222`
- **Person ID**: `22222222-2222-2222-2222-222222222222`
- **Address ID**: `22222222-2222-2222-2222-222222222222` (Av. Santa Fe 2567, Piso 5 A, Buenos Aires)
- **DNI**: `23456789`
- **Contactos**:
  - Email primario: `22222222-1111-2222-2222-222222222222` (maria.gonzalez@personal.com)
  - Phone: `22222222-2222-2222-2222-222222222222` (+541134567890)
  - Email trabajo: `22222222-3333-2222-2222-222222222222` (mgonzalez.work@gmail.com)

### Carlos Alberto Rodríguez (PENDING STATUS)
- **Email login**: `carlos.rodriguez@example.com`
- **User ID**: `33333333-3333-3333-3333-333333333333`
- **Person ID**: `33333333-3333-3333-3333-333333333333`
- **Address ID**: `33333333-3333-3333-3333-333333333333` (Belgrano 890, Mar del Plata)
- **DNI**: `34567890`
- **Status**: `pending` (usuario no ha completado onboarding)

### Ana Sofía Martínez
- **Email login**: `ana.martinez@example.com`
- **User ID**: `44444444-4444-4444-4444-444444444444`
- **Person ID**: `44444444-4444-4444-4444-444444444444`
- **Address ID**: `44444444-4444-4444-4444-444444444444` (Av. Pueyrredón 1456, Piso 12 B, Buenos Aires)
- **DNI**: `45678901`

### Laura Patricia Fernández (Google OAuth)
- **Email login**: `laura.fernandez@gmail.com`
- **User ID**: `55555555-5555-5555-5555-555555555555`
- **Person ID**: `55555555-5555-5555-5555-555555555555`
- **Address ID**: `55555555-5555-5555-5555-555555555555` (San Martín 345, Rosario)
- **DNI**: `56789012`
- **Provider**: `google` (autenticación OAuth)

## 🏢 Empresas

### ACME Corporation S.A.
- **Email login**: `admin@acmecorp.com`
- **User ID**: `66666666-6666-6666-6666-666666666666`
- **Person ID**: `66666666-6666-6666-6666-666666666666`
- **Address ID**: `88888888-8888-8888-8888-888888888888` (Av. Leandro N. Alem 456, Piso 15 OF 1501, Buenos Aires)
- **CUIT**: `30123456789`
- **Tipo**: `S.A.`
- **Contactos**:
  - Email principal: `66666666-1111-6666-6666-666666666666` (info@acmecorp.com)
  - Phone: `66666666-2222-6666-6666-666666666666` (+541156789013)
  - Email ventas: `66666666-3333-6666-6666-666666666666` (ventas@acmecorp.com)

### Tech Solutions SRL
- **Email login**: `contacto@techsolutions.com`
- **User ID**: `77777777-7777-7777-7777-777777777777`
- **Person ID**: `77777777-7777-7777-7777-777777777777`
- **Address ID**: `99999999-9999-9999-9999-999999999999` (Av. 9 de Julio 1020, Piso 8 OF 804, Buenos Aires)
- **CUIT**: `30234567890`
- **Tipo**: `S.R.L.`

### Innova Tech Argentina S.A.
- **Email login**: `ventas@innovatech.com.ar`
- **User ID**: `88888888-8888-8888-8888-888888888888`
- **Person ID**: `88888888-8888-8888-8888-888888888888`
- **Address ID**: `aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa` (Rivadavia 567, Mendoza)
- **CUIT**: `30345678901`
- **Tipo**: `S.A.`

## 🏠 Direcciones Adicionales

### Córdoba
- **Dean Funes**: `66666666-6666-6666-6666-666666666666` (Dean Funes 234, Piso 3 C, Córdoba)
- **Av. Colón**: `77777777-7777-7777-7777-777777777777` (Av. Colón 789, Córdoba)

## 📞 Tipos de Contacto Disponibles

| Tipo | Descripción | Ejemplo |
|------|-------------|---------|
| `email` | Correo electrónico | juan.perez@personal.com |
| `phone` | Teléfono | +541123456789 |
| `whatsapp` | WhatsApp | +541123456789 |

## 🔄 Variables de Postman Configuradas

### URLs Base
- `base_url`: `http://localhost:8081` (API Gateway)
- `auth_service_url`: `http://localhost:4002`
- `person_service_url`: `http://localhost:4001`
- `address_service_url`: `http://localhost:4000`

### Datos de Prueba Principales
- `test_email`: `juan.perez@example.com`
- `test_password`: `password`
- `test_dni`: `12345678`
- `test_phone`: `+541123456789`
- `test_personal_email`: `juan.perez@personal.com`

### IDs Principales (Auto-actualizables)
- `user_id`: Se actualiza con login/me
- `person_id`: Se actualiza con requests de personas
- `address_id`: Se actualiza con requests de direcciones
- `contact_id`: Se actualiza con requests de contactos
- `access_token`: Se actualiza con login
- `refresh_token`: Se actualiza con login

### IDs Específicos por Usuario
- `maria_*`: Variables para María González
- `carlos_*`: Variables para Carlos Rodríguez
- `ana_*`: Variables para Ana Martínez
- `laura_*`: Variables para Laura Fernández
- `acme_*`: Variables para ACME Corporation
- `tech_*`: Variables para Tech Solutions
- `innovatech_*`: Variables para Innova Tech

## 🎯 Casos de Uso Recomendados

### 1. Testing de Autenticación
```
1. Login con juan.perez@example.com / password
2. Verificar auto-guardado de access_token
3. Usar endpoint /api/auth/me para obtener info del usuario
```

### 2. Testing de Personas
```
1. Listar todas las personas (auto-guarda primer person_id)
2. Obtener persona específica con {{person_id}}
3. Crear nueva persona con datos mock
```

### 3. Testing de Direcciones
```
1. Listar todas las direcciones (auto-guarda primer address_id)
2. Obtener dirección específica con {{address_id}}
3. Crear nueva dirección y asociarla a persona
```

### 4. Testing de Empresas
```
1. Usar {{acme_person_id}} para obtener ACME Corporation
2. Verificar datos de empresa vs individual
3. Testing de contactos empresariales
```

### 5. Testing de Estados
```
1. Usar {{carlos_person_id}} para testing de usuario PENDING
2. Verificar manejo de estados de onboarding
3. Testing de activación de cuentas
```

## 🚨 Notas Importantes

- **Todos los passwords** de usuarios de prueba son `password`
- **Los UUIDs** están correlacionados entre tablas para facilitar testing
- **Carlos Rodríguez** tiene status `pending` para testing de flujos incompletos
- **Laura Fernández** usa Google OAuth (sin password hash)
- **API Key** es obligatoria para endpoints de personas, direcciones y contactos
- **JWT Token** es obligatorio para endpoints autenticados (/me, /logout, etc.)

## 🔧 Scripts Automáticos de Postman

### Pre-request Script
- Auto-inyecta `X-Api-Key` cuando es necesario
- Auto-inyecta `Authorization: Bearer` para endpoints autenticados
- Reemplaza variables en request bodies
- Logs informativos en consola

### Test Script
- Auto-guarda `access_token` y `refresh_token` desde login
- Auto-guarda IDs de respuestas (person_id, address_id, contact_id)
- Auto-guarda datos de usuario desde `/me`
- Logs detallados de errores y éxitos
- Manejo inteligente de respuestas de listas vs objetos individuales

Esta referencia te permite usar Postman de forma inmediata sin necesidad de buscar o recordar IDs específicos. Todos los datos están precargados y listos para uso.
