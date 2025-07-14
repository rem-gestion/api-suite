# Migración 0007: Sistema de Dominios Organizacionales

## 📋 Información General

| Atributo | Valor |
|----------|-------|
| **Archivo** | `0007_create_organization_domain.up.sql` |
| **Propósito** | Implementar sistema completo de dominios personalizados con SSL y verificación DNS |
| **Dependencias** | 0001 (ENUMs y funciones base), 0002 (organization) |
| **Reversible** | ✅ Sí |
| **Particionado** | ✅ Logs de verificación particionados mensualmente |

---

## 🎯 Objetivos

1. **Dominios Personalizados**: Gestión completa de dominios custom por organización
2. **Verificación DNS**: Sistema automático de verificación con registros DNS
3. **Certificados SSL**: Gestión de certificados con fechas de expiración
4. **Logs de Auditoría**: Trazabilidad completa de verificaciones
5. **Configuración Automática**: Registros DNS por defecto y tokens de verificación
6. **Soft Delete**: Retención de configuraciones históricas

---

## 🏗️ Arquitectura del Sistema

```mermaid
graph TB
    subgraph "Sistema de Dominios"
        OD[organization_domain]
        ODD[organization_domain_dns]
        ODVL[organization_domain_verification_log]
        
        subgraph "Funciones Core"
            GDVT[generate_dns_verification_token]
            CDDR[create_default_dns_records]
            SDDV[sync_domain_dns_verification]
            LDV[log_domain_verification]
        end
        
        subgraph "Verificación Automática"
            DNS_CHK[DNS Checker]
            SSL_CHK[SSL Checker]
            CRON[Cron Jobs]
        end
    end
    
    subgraph "Servicios Externos"
        DNS_PROV[DNS Provider]
        SSL_CA[Certificate Authority]
        CDN[CDN Service]
    end
    
    subgraph "Dependencias"
        ORG[organization]
    end
    
    OD --> ORG
    OD --> ODD
    OD --> ODVL
    
    GDVT --> OD
    CDDR --> ODD
    SDDV --> OD
    LDV --> ODVL
    
    DNS_CHK --> ODD
    SSL_CHK --> OD
    CRON --> DNS_CHK
    CRON --> SSL_CHK
    
    ODD -.-> DNS_PROV
    OD -.-> SSL_CA
    OD -.-> CDN
```

---

## 📊 Entidades Principales

### 1. **organization_domain** - Tabla Principal de Dominios

```sql
CREATE TABLE organization_domain (
    id                      UUID PRIMARY KEY DEFAULT generate_uuid(),
    organization_id         UUID NOT NULL,                    -- FK a organization
    domain_name             VARCHAR(255) NOT NULL,            -- Nombre del dominio
    subdomain               VARCHAR(100),                     -- Subdominio opcional
    domain_type             domain_type_enum NOT NULL DEFAULT 'custom',
    status                  domain_status_enum NOT NULL DEFAULT 'pending',
    is_primary              BOOLEAN NOT NULL DEFAULT false,   -- Dominio principal
    ssl_enabled             BOOLEAN NOT NULL DEFAULT false,   -- SSL habilitado
    ssl_certificate         TEXT,                             -- Certificado SSL
    ssl_private_key         TEXT,                             -- Clave privada SSL
    ssl_expires_at          TIMESTAMP WITH TIME ZONE,         -- Expiración SSL
    dns_verified            BOOLEAN NOT NULL DEFAULT false,   -- Verificación DNS
    dns_verification_token  VARCHAR(100) DEFAULT generate_dns_verification_token(),
    dns_verification_method dns_verification_method_enum DEFAULT 'txt',
    verification_attempts   INTEGER NOT NULL DEFAULT 0,        -- Intentos de verificación
    last_verification_at    TIMESTAMP WITH TIME ZONE,         -- Última verificación
    redirect_to_primary     BOOLEAN NOT NULL DEFAULT false,   -- Redirigir a primario
    custom_headers          JSONB DEFAULT '{}',               -- Headers personalizados
    -- Auditoría
    created_at              TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp_utc(),
    created_by              UUID NOT NULL,                    -- FK lógica a auth-identity-svc
    updated_at              TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp_utc(),
    updated_by              UUID,                             -- FK lógica a auth-identity-svc
    deleted_at              TIMESTAMP WITH TIME ZONE          -- Soft delete
);
```

#### **Tipos de Dominio**

| Tipo | Descripción | Uso |
|------|-------------|-----|
| `custom` | Dominio personalizado del cliente | empresa.com |
| `subdomain` | Subdominio de la plataforma | empresa.plataforma.com |
| `temporary` | Dominio temporal para testing | test-empresa.dev |

#### **Estados de Dominio**

```mermaid
stateDiagram-v2
    [*] --> pending: Crear dominio
    pending --> configuring: Configurar DNS
    configuring --> verifying: Iniciar verificación
    verifying --> active: DNS verificado + SSL OK
    verifying --> failed: Error en verificación
    failed --> verifying: Reintentar
    active --> suspended: Error crítico
    suspended --> active: Resolver problema
    active --> [*]: Soft delete
    suspended --> [*]: Soft delete
    failed --> [*]: Soft delete
```

#### **Reglas de Negocio**

| Regla | Descripción | Implementación |
|-------|-------------|----------------|
| **Dominio Único** | Un dominio solo puede estar en una organización | Constraint `uq_organization_domain_name_case_insensitive` |
| **Primario Único** | Solo un dominio primario por organización | Index `uq_organization_domain_one_primary_per_org` |
| **SSL Consistente** | Si SSL está habilitado, debe tener certificado | Constraint `chk_organization_domain_ssl_consistency` |
| **No Redirect en Primario** | El dominio primario no puede redirigir | Constraint `chk_organization_domain_redirect_logic` |

### 2. **organization_domain_dns** - Configuración DNS

```sql
CREATE TABLE organization_domain_dns (
    id                      UUID PRIMARY KEY DEFAULT generate_uuid(),
    domain_id               UUID NOT NULL,                     -- FK a organization_domain
    record_type             dns_record_type_enum NOT NULL,     -- A|AAAA|CNAME|TXT|MX
    record_name             VARCHAR(255) NOT NULL,             -- Nombre del registro
    record_value            TEXT NOT NULL,                     -- Valor del registro
    record_ttl              INTEGER NOT NULL DEFAULT 300,      -- TTL en segundos
    is_required             BOOLEAN NOT NULL DEFAULT true,     -- Requerido para verificación
    is_verified             BOOLEAN NOT NULL DEFAULT false,    -- Estado de verificación
    verification_attempts   INTEGER NOT NULL DEFAULT 0,        -- Intentos de verificación
    last_verification_at    TIMESTAMP WITH TIME ZONE,          -- Última verificación
    notes                   TEXT,                              -- Notas para el administrador
    -- Auditoría
    created_at              TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp_utc(),
    updated_at              TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp_utc()
);
```

#### **Tipos de Registros DNS**

| Tipo | Propósito | Ejemplo |
|------|-----------|---------|
| `A` | IPv4 del servidor | `empresa.com → 192.168.1.100` |
| `AAAA` | IPv6 del servidor | `empresa.com → 2001:db8::1` |
| `CNAME` | Alias a otro dominio | `www.empresa.com → empresa.com` |
| `TXT` | Verificación y metadatos | `_verification.empresa.com → token123` |
| `MX` | Servidor de email | `empresa.com → mail.empresa.com` |

#### **Registros Creados Automáticamente**

```sql
-- Al crear un dominio personalizado se generan automáticamente:

-- 1. Registro A principal (requerido)
INSERT INTO organization_domain_dns (domain_id, record_type, record_name, record_value, is_required)
VALUES (domain_id, 'A', 'empresa.com', '127.0.0.1', true);

-- 2. Registro CNAME para www (opcional)
INSERT INTO organization_domain_dns (domain_id, record_type, record_name, record_value, is_required)
VALUES (domain_id, 'CNAME', 'www.empresa.com', 'empresa.com', false);

-- 3. Registro TXT para verificación (requerido)
INSERT INTO organization_domain_dns (domain_id, record_type, record_name, record_value, is_required)
VALUES (domain_id, 'TXT', '_verification.empresa.com', 'generated-token', true);
```

### 3. **organization_domain_verification_log** - Auditoría de Verificación

```sql
CREATE TABLE organization_domain_verification_log (
    id                      UUID PRIMARY KEY DEFAULT generate_uuid(),
    domain_id               UUID NOT NULL,                        -- FK a organization_domain
    dns_record_id           UUID,                                 -- FK opcional a DNS record
    verification_type       domain_verification_type_enum NOT NULL, -- dns|ssl|http
    status                  verification_status_enum NOT NULL,    -- success|failed|pending
    details                 JSONB DEFAULT '{}',                  -- Detalles adicionales
    error_message           TEXT,                                -- Error si falla
    response_data           JSONB,                               -- Respuesta raw del servicio
    duration_ms             INTEGER,                             -- Duración en ms
    created_at              TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp_utc()
) PARTITION BY RANGE (created_at);
```

#### **Tipos de Verificación**

| Tipo | Descripción | Frecuencia |
|------|-------------|------------|
| `dns` | Verificación de registros DNS | Cada 15 minutos |
| `ssl` | Verificación de certificado SSL | Diaria |
| `http` | Verificación de conectividad HTTP | Cada 5 minutos |

#### **Estados de Verificación**

| Estado | Descripción | Acción |
|--------|-------------|--------|
| `pending` | Verificación pendiente | Programar siguiente verificación |
| `success` | Verificación exitosa | Marcar como verificado |
| `failed` | Verificación falló | Incrementar contador de intentos |
| `timeout` | Timeout en verificación | Reintentar con backoff |

---

## 🔧 Funciones Especializadas

### 1. **generate_dns_verification_token()** - Generación de Tokens

```sql
CREATE OR REPLACE FUNCTION generate_dns_verification_token()
RETURNS TEXT
```

#### **Características**

- **Longitud**: 32 caracteres URL-safe
- **Entropía**: Usa `gen_random_bytes()` para máxima seguridad
- **Unicidad**: Verifica que no exista en la base
- **Formato**: Base64 con `/` → `_` y `+` → `-`

#### **Ejemplo de Uso**

```sql
-- Generar token automáticamente al crear dominio
INSERT INTO organization_domain (organization_id, domain_name, created_by)
VALUES ('org-uuid', 'empresa.com', 'user-uuid');
-- Token generado automáticamente: 'xK8nQ2mP_vR5hW9tL3bA7cF6dG1jS4uY'

-- Regenerar token manualmente
UPDATE organization_domain 
SET dns_verification_token = generate_dns_verification_token(),
    verification_attempts = 0
WHERE id = 'domain-uuid';
```

### 2. **create_default_dns_records()** - Configuración Automática

```sql
CREATE OR REPLACE FUNCTION create_default_dns_records()
RETURNS TRIGGER
```

#### **Registros Creados**

1. **Registro A Principal** (requerido): Apunta a IP del servidor
2. **Registro CNAME WWW** (opcional): Redirección de www al dominio principal
3. **Registro TXT Verificación** (requerido): Token para verificar propiedad

#### **Notas Incluidas**

```sql
-- Registro A
'Registro A principal - actualizar con IP real del servidor'

-- Registro CNAME
'Redirección opcional de www al dominio principal'

-- Registro TXT
'Token de verificación DNS - requerido para activar el dominio'
```

### 3. **sync_domain_dns_verification()** - Sincronización Automática

```sql
CREATE OR REPLACE FUNCTION sync_domain_dns_verification()
RETURNS TRIGGER
```

#### **Lógica de Sincronización**

1. **Contar Registros**: Total requeridos vs verificados
2. **Determinar Estado**: `dns_verified = (all required verified)`
3. **Actualizar Dominio**: Solo si el estado cambió
4. **Log Automático**: Registra el cambio en logs de verificación

#### **Ejemplo de Flujo**

```sql
-- Estado inicial: 3 registros requeridos, 0 verificados
SELECT dns_verified FROM organization_domain WHERE id = 'domain-uuid';
-- Resultado: false

-- Verificar registro A
UPDATE organization_domain_dns 
SET is_verified = true 
WHERE domain_id = 'domain-uuid' AND record_type = 'A';
-- dns_verified sigue siendo false (1/3 verificados)

-- Verificar registro TXT
UPDATE organization_domain_dns 
SET is_verified = true 
WHERE domain_id = 'domain-uuid' AND record_type = 'TXT';
-- dns_verified sigue siendo false (2/3 verificados)

-- Marcar CNAME como no requerido
UPDATE organization_domain_dns 
SET is_required = false 
WHERE domain_id = 'domain-uuid' AND record_type = 'CNAME';
-- dns_verified cambia a true (2/2 requeridos verificados)
```

### 4. **log_domain_verification()** - Logging Enriquecido

```sql
CREATE OR REPLACE FUNCTION log_domain_verification(
    p_domain_id UUID,
    p_dns_record_id UUID,
    p_verification_type domain_verification_type_enum,
    p_status verification_status_enum,
    p_details JSONB DEFAULT '{}',
    p_error_message TEXT DEFAULT NULL,
    p_response_data JSONB DEFAULT NULL,
    p_duration_ms INTEGER DEFAULT NULL
)
RETURNS UUID
```

#### **Contexto Enriquecido**

```json
{
  "logged_at": "2024-01-15T10:30:00Z",
  "session_id": "session-uuid-123",
  "user_id": "user-uuid-456",
  "verification_attempt": 3,
  "dns_server": "8.8.8.8",
  "response_time_ms": 150,
  "expected_value": "token123",
  "actual_value": "token123"
}
```

---

## 🚀 Performance y Optimización

### **Índices Especializados**

```sql
-- Búsquedas principales
CREATE INDEX ix_organization_domain_organization_id ON organization_domain(organization_id);
CREATE INDEX ix_organization_domain_active ON organization_domain(organization_id, status) 
    WHERE deleted_at IS NULL;

-- Verificación y polling
CREATE INDEX ix_organization_domain_dns_polling ON organization_domain_dns(is_verified, last_verification_at, verification_attempts)
    WHERE is_verified = false AND is_required = true;
CREATE INDEX ix_organization_domain_verification ON organization_domain(dns_verified, status);

-- SSL y expiración
CREATE INDEX ix_organization_domain_ssl_expiry ON organization_domain(ssl_expires_at) 
    WHERE ssl_expires_at IS NOT NULL;

-- Logs particionados
CREATE INDEX ix_organization_domain_verification_log_domain_id ON organization_domain_verification_log(domain_id);
CREATE INDEX ix_organization_domain_verification_log_status ON organization_domain_verification_log(status);
```

### **Constraints Únicos Inteligentes**

```sql
-- Solo un dominio primario por organización
CREATE UNIQUE INDEX uq_organization_domain_one_primary_per_org 
ON organization_domain(organization_id) 
WHERE is_primary = true AND deleted_at IS NULL;

-- Dominio único case-insensitive (excluye soft-deleted)
CREATE UNIQUE INDEX uq_organization_domain_name_case_insensitive
ON organization_domain(organization_id, LOWER(domain_name))
WHERE deleted_at IS NULL;

-- Registro DNS único por dominio
CREATE UNIQUE INDEX uq_organization_domain_dns_unique 
ON organization_domain_dns(domain_id, record_type, record_name);
```

---

## 📋 Ejemplos Prácticos

### **1. Configurar Dominio Personalizado Completo**

```sql
-- 1. Crear dominio personalizado
INSERT INTO organization_domain (
    organization_id,
    domain_name,
    domain_type,
    is_primary,
    ssl_enabled,
    created_by
) VALUES (
    'org-uuid-123',
    'empresa.com',
    'custom',
    true,
    true,
    'admin-user-uuid'
);

-- 2. Los registros DNS se crean automáticamente via trigger
-- Verificar registros creados
SELECT 
    record_type,
    record_name,
    record_value,
    is_required,
    notes
FROM organization_domain_dns 
WHERE domain_id = (SELECT id FROM organization_domain WHERE domain_name = 'empresa.com')
ORDER BY is_required DESC, record_type;

-- 3. Actualizar registro A con IP real
UPDATE organization_domain_dns 
SET record_value = '203.0.113.100'
WHERE domain_id = (SELECT id FROM organization_domain WHERE domain_name = 'empresa.com')
AND record_type = 'A';

-- 4. Configurar certificado SSL
UPDATE organization_domain 
SET ssl_certificate = '-----BEGIN CERTIFICATE-----...',
    ssl_private_key = 'encrypted:AES256:...',  -- Encriptado en aplicación
    ssl_expires_at = '2024-12-31 23:59:59+00'
WHERE domain_name = 'empresa.com';
```

### **2. Proceso de Verificación DNS**

```sql
-- Simular verificación DNS exitosa
-- (En producción esto lo haría un worker/cron job)

-- 1. Obtener registros pendientes de verificación
SELECT 
    dd.id,
    dd.record_type,
    dd.record_name,
    dd.record_value,
    dd.verification_attempts
FROM organization_domain_dns dd
JOIN organization_domain d ON dd.domain_id = d.id
WHERE dd.is_verified = false 
AND dd.is_required = true
AND dd.verification_attempts < 10
AND d.status IN ('configuring', 'verifying')
ORDER BY dd.last_verification_at ASC NULLS FIRST;

-- 2. Realizar verificación (simulada)
DO $$
DECLARE
    dns_record RECORD;
    verification_result BOOLEAN;
    error_msg TEXT;
    log_id UUID;
BEGIN
    FOR dns_record IN 
        SELECT * FROM organization_domain_dns 
        WHERE is_verified = false AND is_required = true
        LIMIT 5
    LOOP
        -- Simular verificación DNS
        verification_result := (random() > 0.3); -- 70% éxito
        
        IF verification_result THEN
            -- Verificación exitosa
            UPDATE organization_domain_dns 
            SET is_verified = true,
                last_verification_at = current_timestamp_utc()
            WHERE id = dns_record.id;
            
            -- Log exitoso
            SELECT log_domain_verification(
                dns_record.domain_id,
                dns_record.id,
                'dns'::domain_verification_type_enum,
                'success'::verification_status_enum,
                jsonb_build_object(
                    'record_type', dns_record.record_type,
                    'expected', dns_record.record_value,
                    'actual', dns_record.record_value,
                    'dns_server', '8.8.8.8'
                ),
                NULL,
                NULL,
                150 + (random() * 200)::integer
            ) INTO log_id;
        ELSE
            -- Verificación falló
            error_msg := 'DNS record not found or value mismatch';
            
            UPDATE organization_domain_dns 
            SET verification_attempts = verification_attempts + 1,
                last_verification_at = current_timestamp_utc()
            WHERE id = dns_record.id;
            
            -- Log con error
            SELECT log_domain_verification(
                dns_record.domain_id,
                dns_record.id,
                'dns'::domain_verification_type_enum,
                'failed'::verification_status_enum,
                jsonb_build_object(
                    'record_type', dns_record.record_type,
                    'expected', dns_record.record_value,
                    'attempt', dns_record.verification_attempts + 1
                ),
                error_msg,
                jsonb_build_object('dns_response', 'NXDOMAIN'),
                300 + (random() * 500)::integer
            ) INTO log_id;
        END IF;
        
        RAISE NOTICE 'Verified record %: % (success: %)', 
            dns_record.record_name, dns_record.record_type, verification_result;
    END LOOP;
END $$;

-- 3. Verificar estado de sincronización
SELECT 
    d.domain_name,
    d.dns_verified,
    COUNT(*) as total_records,
    COUNT(*) FILTER (WHERE dd.is_required) as required_records,
    COUNT(*) FILTER (WHERE dd.is_required AND dd.is_verified) as verified_required
FROM organization_domain d
LEFT JOIN organization_domain_dns dd ON d.id = dd.domain_id
WHERE d.domain_name = 'empresa.com'
GROUP BY d.id, d.domain_name, d.dns_verified;
```

### **3. Monitoreo de Certificados SSL**

```sql
-- Certificados próximos a expirar (próximos 30 días)
SELECT 
    d.domain_name,
    o.organization_name,
    d.ssl_expires_at,
    EXTRACT(EPOCH FROM d.ssl_expires_at - current_timestamp_utc())/86400 as days_until_expiry,
    d.status,
    d.created_by
FROM organization_domain d
JOIN organization o ON d.organization_id = o.id
WHERE d.ssl_enabled = true
AND d.ssl_expires_at IS NOT NULL
AND d.ssl_expires_at BETWEEN current_timestamp_utc() AND current_timestamp_utc() + INTERVAL '30 days'
AND d.deleted_at IS NULL
ORDER BY d.ssl_expires_at ASC;

-- Renovar certificado SSL
UPDATE organization_domain 
SET ssl_certificate = '-----BEGIN CERTIFICATE-----...',
    ssl_private_key = 'encrypted:AES256:new-key...',
    ssl_expires_at = current_timestamp_utc() + INTERVAL '1 year',
    updated_by = 'admin-user-uuid'
WHERE domain_name = 'empresa.com';

-- Log de renovación
SELECT log_domain_verification(
    (SELECT id FROM organization_domain WHERE domain_name = 'empresa.com'),
    NULL,
    'ssl'::domain_verification_type_enum,
    'success'::verification_status_enum,
    jsonb_build_object(
        'action', 'certificate_renewed',
        'new_expiry', current_timestamp_utc() + INTERVAL '1 year',
        'renewed_by', 'admin-user-uuid'
    )
);
```

### **4. Gestión de Dominios Múltiples**

```sql
-- Configurar dominio secundario con redirección
INSERT INTO organization_domain (
    organization_id,
    domain_name,
    domain_type,
    is_primary,
    redirect_to_primary,
    ssl_enabled,
    created_by
) VALUES (
    'org-uuid-123',
    'empresa.net',
    'custom',
    false,
    true,  -- Redirigir al dominio primario
    true,
    'admin-user-uuid'
);

-- Listar todos los dominios de una organización
SELECT 
    d.domain_name,
    d.domain_type,
    d.status,
    d.is_primary,
    d.redirect_to_primary,
    d.dns_verified,
    d.ssl_enabled,
    CASE 
        WHEN d.ssl_expires_at IS NULL THEN 'No SSL'
        WHEN d.ssl_expires_at < current_timestamp_utc() THEN 'SSL Expired'
        WHEN d.ssl_expires_at < current_timestamp_utc() + INTERVAL '30 days' THEN 'SSL Expiring Soon'
        ELSE 'SSL Valid'
    END as ssl_status,
    COUNT(dd.id) as dns_records_count,
    COUNT(dd.id) FILTER (WHERE dd.is_verified) as verified_records_count
FROM organization_domain d
LEFT JOIN organization_domain_dns dd ON d.id = dd.domain_id
WHERE d.organization_id = 'org-uuid-123'
AND d.deleted_at IS NULL
GROUP BY d.id, d.domain_name, d.domain_type, d.status, d.is_primary, 
         d.redirect_to_primary, d.dns_verified, d.ssl_enabled, d.ssl_expires_at
ORDER BY d.is_primary DESC, d.domain_name ASC;

-- Cambiar dominio primario
BEGIN;
    -- Quitar primario del actual
    UPDATE organization_domain 
    SET is_primary = false,
        updated_by = 'admin-user-uuid'
    WHERE organization_id = 'org-uuid-123' 
    AND is_primary = true;
    
    -- Establecer nuevo primario
    UPDATE organization_domain 
    SET is_primary = true,
        redirect_to_primary = false,  -- Se maneja automáticamente por trigger
        updated_by = 'admin-user-uuid'
    WHERE organization_id = 'org-uuid-123' 
    AND domain_name = 'empresa.net';
COMMIT;
```

### **5. Troubleshooting y Diagnóstico**

```sql
-- Dominios con problemas de verificación
SELECT 
    d.domain_name,
    o.organization_name,
    d.status,
    d.verification_attempts,
    d.last_verification_at,
    age(current_timestamp_utc(), d.last_verification_at) as time_since_last_attempt,
    COUNT(dd.id) FILTER (WHERE dd.is_required AND NOT dd.is_verified) as failed_required_records
FROM organization_domain d
JOIN organization o ON d.organization_id = o.id
LEFT JOIN organization_domain_dns dd ON d.id = dd.domain_id
WHERE d.deleted_at IS NULL
AND (
    d.verification_attempts > 5 OR
    d.last_verification_at < current_timestamp_utc() - INTERVAL '6 hours' OR
    d.dns_verified = false
)
GROUP BY d.id, d.domain_name, o.organization_name, d.status, 
         d.verification_attempts, d.last_verification_at
HAVING COUNT(dd.id) FILTER (WHERE dd.is_required AND NOT dd.is_verified) > 0
ORDER BY d.verification_attempts DESC, d.last_verification_at ASC;

-- Historial de verificación de un dominio específico
SELECT 
    dvl.created_at,
    dvl.verification_type,
    dvl.status,
    dvl.duration_ms,
    dvl.error_message,
    dvl.details,
    dd.record_type,
    dd.record_name
FROM organization_domain_verification_log dvl
LEFT JOIN organization_domain_dns dd ON dvl.dns_record_id = dd.id
WHERE dvl.domain_id = (SELECT id FROM organization_domain WHERE domain_name = 'empresa.com')
ORDER BY dvl.created_at DESC
LIMIT 50;

-- Estadísticas de éxito de verificación por tipo
SELECT 
    verification_type,
    COUNT(*) as total_attempts,
    COUNT(*) FILTER (WHERE status = 'success') as successful,
    COUNT(*) FILTER (WHERE status = 'failed') as failed,
    ROUND(COUNT(*) FILTER (WHERE status = 'success') * 100.0 / COUNT(*), 2) as success_rate,
    AVG(duration_ms) as avg_duration_ms
FROM organization_domain_verification_log
WHERE created_at > current_timestamp_utc() - INTERVAL '7 days'
GROUP BY verification_type
ORDER BY success_rate DESC;
```

---

## 🔍 Consultas de Diagnóstico

### **1. Health Check General**

```sql
-- Estado general del sistema de dominios
SELECT 
    'total_domains' as metric,
    COUNT(*) as value,
    'Active custom domains in system' as description
FROM organization_domain 
WHERE deleted_at IS NULL AND domain_type = 'custom'

UNION ALL

SELECT 
    'verified_domains',
    COUNT(*),
    'Domains with complete DNS verification'
FROM organization_domain 
WHERE deleted_at IS NULL AND dns_verified = true

UNION ALL

SELECT 
    'ssl_enabled_domains',
    COUNT(*),
    'Domains with SSL certificates'
FROM organization_domain 
WHERE deleted_at IS NULL AND ssl_enabled = true

UNION ALL

SELECT 
    'primary_domains',
    COUNT(*),
    'Organizations with primary domain set'
FROM organization_domain 
WHERE deleted_at IS NULL AND is_primary = true

UNION ALL

SELECT 
    'expiring_ssl_soon',
    COUNT(*),
    'SSL certificates expiring in next 30 days'
FROM organization_domain 
WHERE deleted_at IS NULL 
AND ssl_expires_at BETWEEN current_timestamp_utc() AND current_timestamp_utc() + INTERVAL '30 days';
```

### **2. Performance y Recursos**

```sql
-- Tamaño de tablas de dominios
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size,
    CASE 
        WHEN tablename LIKE '%verification_log%' THEN 'Log Partition'
        WHEN tablename LIKE '%_dns' THEN 'DNS Records'
        ELSE 'Core Table'
    END as table_type
FROM pg_tables 
WHERE tablename LIKE 'organization_domain%'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

-- Análisis de particiones de logs
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size
FROM pg_tables 
WHERE tablename LIKE 'organization_domain_verification_log_%'
ORDER BY tablename DESC;

-- Uso de índices más importantes
SELECT 
    indexname,
    idx_tup_read,
    idx_tup_fetch,
    idx_scan,
    pg_size_pretty(pg_relation_size(indexrelname)) as index_size
FROM pg_stat_user_indexes 
WHERE relname LIKE 'organization_domain%'
ORDER BY idx_scan DESC
LIMIT 20;
```

### **3. Patrones de Uso y Tendencias**

```sql
-- Organizaciones más activas en configuración de dominios
SELECT 
    o.organization_name,
    COUNT(d.id) as total_domains,
    COUNT(d.id) FILTER (WHERE d.is_primary) as primary_domains,
    COUNT(d.id) FILTER (WHERE d.dns_verified) as verified_domains,
    COUNT(d.id) FILTER (WHERE d.ssl_enabled) as ssl_domains,
    MAX(d.created_at) as last_domain_added
FROM organization o
LEFT JOIN organization_domain d ON o.id = d.organization_id AND d.deleted_at IS NULL
GROUP BY o.id, o.organization_name
HAVING COUNT(d.id) > 0
ORDER BY total_domains DESC, last_domain_added DESC
LIMIT 20;

-- Análisis de tiempo de verificación por tipo de registro
SELECT 
    dd.record_type,
    COUNT(*) as total_records,
    COUNT(*) FILTER (WHERE dd.is_verified) as verified_count,
    AVG(dd.verification_attempts) as avg_attempts,
    AVG(EXTRACT(EPOCH FROM dd.last_verification_at - dd.created_at)/3600) FILTER (WHERE dd.is_verified) as avg_hours_to_verify
FROM organization_domain_dns dd
WHERE dd.is_required = true
GROUP BY dd.record_type
ORDER BY verified_count DESC;

-- Patrones de fallo en verificación
SELECT 
    DATE_TRUNC('day', dvl.created_at) as verification_date,
    dvl.verification_type,
    COUNT(*) as total_attempts,
    COUNT(*) FILTER (WHERE dvl.status = 'success') as successful,
    COUNT(*) FILTER (WHERE dvl.status = 'failed') as failed,
    ROUND(AVG(dvl.duration_ms)) as avg_duration_ms
FROM organization_domain_verification_log dvl
WHERE dvl.created_at > current_timestamp_utc() - INTERVAL '30 days'
GROUP BY DATE_TRUNC('day', dvl.created_at), dvl.verification_type
ORDER BY verification_date DESC, dvl.verification_type;
```

---

## ⚠️ Consideraciones de Seguridad

### **1. Protección de Datos Sensibles**

```sql
-- ❌ NUNCA almacenar claves SSL en texto plano
UPDATE organization_domain 
SET ssl_private_key = '-----BEGIN PRIVATE KEY-----...'  -- ¡INSEGURO!
WHERE id = 'domain-uuid';

-- ✅ Encriptar en la aplicación antes de guardar
UPDATE organization_domain 
SET ssl_private_key = 'encrypted:AES256:iv:encrypted_data'
WHERE id = 'domain-uuid';
```

### **2. Validación de Dominios**

```sql
-- Verificar que los dominios son válidos y seguros
SELECT 
    domain_name,
    CASE 
        WHEN domain_name ~ '^[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9]?\.[a-zA-Z]{2,}$' THEN 'Valid'
        WHEN domain_name LIKE '%.local' OR domain_name LIKE '%.test' THEN 'Test domain'
        WHEN domain_name LIKE '%localhost%' THEN 'Localhost - insecure'
        ELSE 'Invalid format'
    END as validation_status
FROM organization_domain
WHERE deleted_at IS NULL;
```

### **3. Verificación de Propiedad**

```sql
-- Asegurar que los tokens de verificación son únicos y seguros
SELECT 
    domain_name,
    dns_verification_token,
    char_length(dns_verification_token) as token_length,
    verification_attempts,
    CASE 
        WHEN verification_attempts > 5 THEN 'Suspicious activity'
        WHEN char_length(dns_verification_token) < 16 THEN 'Weak token'
        ELSE 'OK'
    END as security_status
FROM organization_domain
WHERE deleted_at IS NULL
AND dns_verification_token IS NOT NULL;
```

---

## 🔄 Workflows de Integración

### **1. Configuración Completa de Dominio**

```mermaid
sequenceDiagram
    participant A as Admin
    participant S as organization-svc
    participant DNS as DNS Provider
    participant CA as Certificate Authority
    participant DB as PostgreSQL
    
    A->>S: POST /domains {domain: "empresa.com"}
    S->>DB: INSERT organization_domain
    DB->>DB: Trigger: create_default_dns_records()
    DB-->>S: Dominio + registros DNS creados
    S-->>A: Configuración inicial completa
    
    Note over A: Admin configura DNS en proveedor
    A->>DNS: Configurar registros A, TXT, CNAME
    
    Note over S: Cron job verifica DNS cada 15 min
    S->>DNS: Verificar registros DNS
    DNS-->>S: Registros encontrados
    S->>DB: UPDATE dns records verified
    DB->>DB: Trigger: sync_domain_dns_verification()
    DB->>DB: UPDATE domain.dns_verified = true
    
    A->>S: POST /domains/{id}/ssl {certificate, key}
    S->>CA: Validar certificado
    CA-->>S: Certificado válido
    S->>DB: UPDATE SSL info
    DB-->>S: SSL configurado
    S-->>A: Dominio totalmente configurado
```

### **2. Verificación Automática Continua**

```mermaid
sequenceDiagram
    participant C as Cron Job
    participant S as organization-svc
    participant DNS as DNS Checker
    participant SSL as SSL Checker
    participant DB as PostgreSQL
    
    Note over C: Cada 15 minutos
    C->>S: Verificar dominios pendientes
    S->>DB: SELECT domains needing verification
    DB-->>S: Lista de dominios
    
    loop Para cada dominio
        S->>DNS: Verificar registros DNS
        DNS-->>S: Resultado verificación
        S->>DB: UPDATE verification status
        
        alt DNS verificado
            S->>SSL: Verificar certificado SSL
            SSL-->>S: Estado SSL
            S->>DB: Log verification results
        end
    end
    
    S->>DB: Cleanup old verification logs
    DB-->>S: Logs limpiados
```

---

## 📚 Referencias y Recursos

### **Documentación Relacionada**

- [README_0001.md](./README_0001.md) - ENUMs y funciones base
- [README_0002.md](./README_0002.md) - Organización core

### **Extensiones PostgreSQL Utilizadas**

- `pgcrypto`: Generación de tokens de verificación seguros

### **Patrones Implementados**

- **Auto-Configuration**: Registros DNS creados automáticamente
- **State Synchronization**: dns_verified sincronizado con registros DNS
- **Audit Trail**: Log completo de verificaciones con contexto
- **Partitioning**: Logs particionados por fecha para performance
- **Soft Delete**: Retención de configuraciones históricas

### **Herramientas de Verificación DNS**

```bash
# Verificar registro A
dig A empresa.com @8.8.8.8

# Verificar registro TXT de verificación
dig TXT _verification.empresa.com @8.8.8.8

# Verificar certificado SSL
openssl s_client -connect empresa.com:443 -servername empresa.com

# Verificar fecha de expiración SSL
echo | openssl s_client -connect empresa.com:443 2>/dev/null | openssl x509 -noout -dates
```

### **Métricas de Monitoreo**

```sql
-- Dashboard de dominios para monitoreo
SELECT 'domains_total' as metric, COUNT(*) as value
FROM organization_domain WHERE deleted_at IS NULL
UNION ALL
SELECT 'domains_verified', COUNT(*)
FROM organization_domain WHERE deleted_at IS NULL AND dns_verified = true
UNION ALL
SELECT 'ssl_expiring_30d', COUNT(*)
FROM organization_domain WHERE deleted_at IS NULL 
AND ssl_expires_at BETWEEN CURRENT_TIMESTAMP AND CURRENT_TIMESTAMP + INTERVAL '30 days'
UNION ALL
SELECT 'verification_success_rate_24h', 
       ROUND(COUNT(*) FILTER (WHERE status = 'success') * 100.0 / COUNT(*), 2)
FROM organization_domain_verification_log 
WHERE created_at > CURRENT_TIMESTAMP - INTERVAL '24 hours';
```

---

**Migración completada exitosamente** ✅  
*Sistema de dominios implementado con verificación DNS automática, gestión SSL y auditoría completa*
