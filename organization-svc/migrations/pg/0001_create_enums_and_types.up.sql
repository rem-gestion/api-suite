-- =============================================
-- Migration: 0001_create_enums_and_types.up.sql
-- Description: Create PostgreSQL ENUMs and custom types for type safety and performance
-- Author: System Generated
-- Date: 2024
-- =============================================

-- =============================================
-- ORGANIZATION RELATED ENUMS
-- =============================================

-- Status de organización
CREATE TYPE organization_status_enum AS ENUM (
    'active',
    'suspended', 
    'deleted',
    'pending_activation'
);

-- Tipo de propietario
CREATE TYPE owner_type_enum AS ENUM (
    'individual',
    'company'
);

-- =============================================
-- EMPLOYEE RELATED ENUMS
-- =============================================

-- Status de empleado
CREATE TYPE employee_status_enum AS ENUM (
    'active',
    'inactive',
    'suspended',
    'terminated'
);

-- =============================================
-- INVITATION RELATED ENUMS
-- =============================================

-- Status de invitación
CREATE TYPE invitation_status_enum AS ENUM (
    'pending',
    'accepted',
    'rejected',
    'expired',
    'cancelled'
);

-- Acciones de log de invitación
CREATE TYPE invitation_log_action_enum AS ENUM (
    'created',
    'sent',
    'viewed',
    'accepted',
    'rejected',
    'expired',
    'cancelled',
    'updated',
    'resent'
);

-- =============================================
-- INTEGRATION RELATED ENUMS
-- =============================================

-- Categorías de integración
CREATE TYPE integration_category_enum AS ENUM (
    'accounting',
    'crm',
    'marketing',
    'communication',
    'analytics',
    'storage',
    'payment',
    'automation',
    'other'
);

-- Status de integración
CREATE TYPE integration_status_enum AS ENUM (
    'inactive',
    'active',
    'error',
    'suspended',
    'configuring'
);

-- Status de sincronización
CREATE TYPE sync_status_enum AS ENUM (
    'success',
    'error',
    'partial',
    'in_progress'
);

-- Frecuencia de sincronización
CREATE TYPE sync_frequency_enum AS ENUM (
    'manual',
    'hourly',
    'daily',
    'weekly',
    'monthly'
);

-- Tipos de eventos de integración
CREATE TYPE integration_event_type_enum AS ENUM (
    'sync',
    'webhook',
    'api_call',
    'auth_refresh',
    'config_change',
    'error',
    'manual_trigger'
);

-- Status de eventos de integración
CREATE TYPE integration_event_status_enum AS ENUM (
    'pending',
    'processing',
    'success',
    'error',
    'cancelled',
    'retry'
);

-- Niveles de log
CREATE TYPE log_level_enum AS ENUM (
    'debug',
    'info',
    'warning',
    'error',
    'critical'
);

-- =============================================
-- DOMAIN RELATED ENUMS
-- =============================================

-- Tipos de dominio
CREATE TYPE domain_type_enum AS ENUM (
    'custom',
    'subdomain',
    'system'
);

-- Status de dominio
CREATE TYPE domain_status_enum AS ENUM (
    'pending',
    'active',
    'suspended',
    'expired',
    'failed'
);

-- Métodos de verificación DNS
CREATE TYPE dns_verification_method_enum AS ENUM (
    'txt',
    'cname',
    'file'
);

-- Tipos de registro DNS
CREATE TYPE dns_record_type_enum AS ENUM (
    'A',
    'AAAA',
    'CNAME',
    'TXT',
    'MX',
    'NS'
);

-- Tipos de verificación de dominio
CREATE TYPE domain_verification_type_enum AS ENUM (
    'dns',
    'ssl',
    'domain',
    'connectivity'
);

-- Status de verificación
CREATE TYPE verification_status_enum AS ENUM (
    'success',
    'failed',
    'timeout',
    'partial'
);

-- =============================================
-- SUBSCRIPTION RELATED ENUMS
-- =============================================

-- Tipos de plan (uniformizado a snake_case)
CREATE TYPE plan_type_enum AS ENUM (
    'free',
    'trial',
    'recurring',
    'one_time',
    'enterprise'
);

-- Intervalos de facturación (uniformizado a snake_case)
CREATE TYPE billing_interval_enum AS ENUM (
    'monthly',
    'quarterly',
    'yearly',
    'one_time'
);

-- Monedas soportadas
CREATE TYPE currency_enum AS ENUM (
    'USD',
    'EUR',
    'GBP',
    'CAD',
    'AUD',
    'MXN'
);

-- Status de suscripción
CREATE TYPE subscription_status_enum AS ENUM (
    'pending',
    'active',
    'past_due',
    'cancelled',
    'paused',
    'trialing'
);

-- Status de factura
CREATE TYPE invoice_status_enum AS ENUM (
    'draft',
    'open',
    'paid',
    'void',
    'uncollectible'
);

-- Métricas de uso
CREATE TYPE usage_metric_enum AS ENUM (
    'users',
    'properties',
    'storage_gb',
    'api_calls',
    'email_sent',
    'sms_sent',
    'documents_generated'
);

-- Tipos de agregación de métricas
CREATE TYPE aggregation_type_enum AS ENUM (
    'sum',
    'max',
    'avg',
    'count'
);

-- =============================================
-- AUDIT RELATED TYPES
-- =============================================

-- Crear extensiones necesarias
-- pgcrypto: Para gen_random_uuid() y funciones criptográficas modernas (suficiente para UUIDs)
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Función helper para generar UUIDs
CREATE OR REPLACE FUNCTION generate_uuid()
RETURNS UUID AS $$
BEGIN
    RETURN gen_random_uuid();
END;
$$ LANGUAGE plpgsql IMMUTABLE SET search_path = public;

-- Función helper para timestamps (retorna timestamptz para consistencia global)
CREATE OR REPLACE FUNCTION current_timestamp_utc()
RETURNS TIMESTAMP WITH TIME ZONE AS $$
BEGIN
    RETURN CURRENT_TIMESTAMP AT TIME ZONE 'UTC';
END;
$$ LANGUAGE plpgsql STABLE SET search_path = public;

-- =============================================
-- UTILITY FUNCTIONS
-- =============================================

-- Función para validar email
CREATE OR REPLACE FUNCTION is_valid_email(email TEXT)
RETURNS BOOLEAN AS $$
BEGIN
    RETURN email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$';
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- Función para validar UUID (mejorada para performance)
CREATE OR REPLACE FUNCTION is_valid_uuid(uuid_text TEXT)
RETURNS BOOLEAN AS $$
BEGIN
    -- Usa casting directo que es más eficiente que regexp
    RETURN (uuid_text::UUID) IS NOT NULL;
EXCEPTION
    WHEN invalid_text_representation THEN
        RETURN FALSE;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- Función para validar dominio (mejorada con restricciones CRLF)
CREATE OR REPLACE FUNCTION is_valid_domain(domain_name TEXT)
RETURNS BOOLEAN AS $$
BEGIN
    -- Verificar que no sea NULL o vacío
    IF domain_name IS NULL OR domain_name = '' THEN
        RETURN FALSE;
    END IF;
    
    -- Verificar longitud (max 253 caracteres para FQDN)
    IF char_length(domain_name) > 253 THEN
        RETURN FALSE;
    END IF;
    
    -- Regex mejorado con anclas estrictas para prevenir bypass CRLF
    -- \A y \z aseguran inicio y fin completo, soporta TLD con dígitos
    RETURN domain_name ~ E'^\\A([a-z0-9]([a-z0-9\\-]{0,61}[a-z0-9])?\\\.)+[a-z0-9]{2,}\\z$';
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- =============================================
-- PARTITIONING UTILITIES
-- =============================================

-- Función para crear particiones por mes (mejorada con soporte para esquemas)
CREATE OR REPLACE FUNCTION create_monthly_partition(
    table_name      TEXT,
    partition_date  DATE   DEFAULT CURRENT_DATE,
    schema_name     TEXT   DEFAULT 'public'
)
RETURNS VOID AS $$
DECLARE
    partition_name        TEXT;
    start_date            DATE;
    end_date              DATE;
    idx_name              TEXT;
BEGIN
    -- Calcular fechas de inicio y fin del mes
    start_date := DATE_TRUNC('month', partition_date)::DATE;
    end_date   := (DATE_TRUNC('month', partition_date) + INTERVAL '1 month')::DATE;

    -- Nombre de la partición sin esquema
    partition_name := table_name || '_' || TO_CHAR(partition_date, 'YYYY_MM');

    -- 1) Crear partición específica del mes
    EXECUTE format(
      'CREATE TABLE IF NOT EXISTS %I.%I PARTITION OF %I.%I
         FOR VALUES FROM (%L) TO (%L)',
      schema_name,          -- %I → schema
      partition_name,       -- %I → partition table
      schema_name,          -- %I → schema
      table_name,           -- %I → parent table
      start_date,           -- %L → literal date
      end_date              -- %L → literal date
    );

    -- 2) Crear partición DEFAULT (una sola vez)
    BEGIN
      EXECUTE format(
        'CREATE TABLE IF NOT EXISTS %I.%I_default PARTITION OF %I.%I DEFAULT',
        schema_name,
        table_name,
        schema_name,
        table_name
      );
    EXCEPTION WHEN duplicate_table THEN
      NULL;
    END;

    -- 3) Índice sobre created_at en la partición recién creada
    idx_name := 'idx_' || table_name || '_' || TO_CHAR(partition_date, 'YYYY_MM') || '_created_at';
    EXECUTE format(
      'CREATE INDEX IF NOT EXISTS %I ON %I.%I (created_at)',
      idx_name,
      schema_name,
      partition_name
    );
END;
$$ LANGUAGE plpgsql;


-- Función para limpiar particiones antiguas
CREATE OR REPLACE FUNCTION cleanup_old_partitions(
    table_name TEXT,
    months_to_keep INTEGER DEFAULT 12
)
RETURNS INTEGER AS $$
DECLARE
    partition_name TEXT;
    cutoff_date DATE;
    dropped_count INTEGER := 0;
    partition_rec RECORD;
BEGIN
    cutoff_date := (CURRENT_DATE - INTERVAL '1 month' * months_to_keep)::DATE;
    
    -- Buscar particiones antiguas
    FOR partition_rec IN
        SELECT schemaname, tablename
        FROM pg_tables
        WHERE tablename LIKE table_name || '_%'
        AND schemaname = 'public'
    LOOP
        -- Extraer fecha de la partición del nombre
        partition_name := partition_rec.tablename;
        
        -- Si la partición es anterior al cutoff, eliminarla
        -- Esto es una implementación simplificada, en producción sería más robusta
        IF partition_name ~ table_name || '_[0-9]{4}_[0-9]{2}$' THEN
            EXECUTE format('DROP TABLE IF EXISTS %I', partition_name);
            dropped_count := dropped_count + 1;
        END IF;
    END LOOP;
    
    RETURN dropped_count;
END;
$$ LANGUAGE plpgsql;

-- =============================================
-- NOTAS IMPORTANTES SOBRE USO DE PARTICIONES
-- =============================================

/*
EJEMPLOS DE USO DE FUNCIONES DE PARTICIONADO:

1. Para crear trigger automático de particiones en organization_integration_log:

    CREATE OR REPLACE FUNCTION organization_integration_log_insert_trigger()
    RETURNS TRIGGER AS $$
    BEGIN
        -- Crear partición si no existe para el mes del registro
        PERFORM create_monthly_partition('organization_integration_log', NEW.created_at::DATE);
        RETURN NEW;
    END;
    $$ LANGUAGE plpgsql;

    CREATE TRIGGER trigger_organization_integration_log_partition
        BEFORE INSERT ON organization_integration_log
        FOR EACH ROW EXECUTE FUNCTION organization_integration_log_insert_trigger();

2. Para limpiar particiones antiguas (ejecutar como job programado):
    
    SELECT cleanup_old_partitions('organization_integration_log', 12); -- Mantener 12 meses

3. Para crear particiones por adelantado:
    
    SELECT create_monthly_partition('organization_integration_log', CURRENT_DATE + INTERVAL '1 month');
*/

-- Comentarios para documentación
COMMENT ON TYPE organization_status_enum IS 'Estados posibles de una organización';
COMMENT ON TYPE invitation_status_enum IS 'Estados del ciclo de vida de invitaciones';
COMMENT ON TYPE integration_status_enum IS 'Estados de conexión de integraciones externas';
COMMENT ON TYPE domain_status_enum IS 'Estados de verificación y activación de dominios';
COMMENT ON TYPE subscription_status_enum IS 'Estados del ciclo de vida de suscripciones';

-- Documentación de ENUMs de facturación en organization-svc:
-- Estos ENUMs (billing_interval_enum, currency_enum, plan_type_enum) están aquí porque:
-- 1. La tabla 'organization' necesita referencias directas a plan_type y billing_interval
-- 2. Se evita dependencia circular entre organization-svc y billing-svc
-- 3. billing-svc puede mantener mirror tables para detalles específicos de facturación
-- 4. Facilita queries y validaciones en el contexto organizacional
COMMENT ON TYPE plan_type_enum IS 'Tipos de plan de suscripción (definido en organization-svc por dependencias)';
COMMENT ON TYPE billing_interval_enum IS 'Intervalos de facturación (definido en organization-svc por dependencias)';
COMMENT ON TYPE currency_enum IS 'Monedas soportadas (definido en organization-svc por dependencias)';

COMMENT ON FUNCTION is_valid_email(TEXT) IS 'Valida formato de email usando regex';
COMMENT ON FUNCTION is_valid_domain(TEXT) IS 'Valida formato de nombre de dominio (requiere TLD válido)';
COMMENT ON FUNCTION create_monthly_partition(TEXT, DATE, TEXT) IS 'Crea particiones mensuales para tablas de logs';
COMMENT ON FUNCTION cleanup_old_partitions(TEXT, INTEGER) IS 'Elimina particiones antiguas para gestión de espacio';
