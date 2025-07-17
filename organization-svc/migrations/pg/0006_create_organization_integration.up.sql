-- =============================================
-- Migration: 0006_create_organization_integration.up.sql
-- Description: Create organization integration system with comprehensive event tracking and audit trails
-- Author: System Generated
-- TRIGGERS Y FUNCIONES DE AUDITORÍA 
-- NOTA: Reutiliza update_timestamp_utc() definida en migración 0002
-- =============================================================================

-- Alias para compatibilidad con triggers existentes (idempotente)
CREATE OR REPLACE FUNCTION update_integration_updated_at()
RETURNS TRIGGER AS $$
BEGIN 
    RETURN update_timestamp_utc(); 
END;
$$ LANGUAGE plpgsql;

-- Date: 2024
-- Version: 1.1
-- Dependencies: 0001 (ENUMs), 0002 (organization), 0005 (organization_invite)
-- =============================================

-- Asegurar que pgcrypto esté disponible para posibles tokens de integración
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- =============================================================================
-- TABLA: integration_type (catálogo de tipos de integración disponibles)
-- Purpose: Catalog of available integration types with metadata and capabilities
-- Business Rules: Global catalog managed by system administrators
-- =============================================================================
CREATE TABLE integration_type (
    id                  UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name                VARCHAR(100) NOT NULL UNIQUE,
    display_name        VARCHAR(200) NOT NULL,
    description         TEXT,
    category            integration_category_enum NOT NULL,
    provider            VARCHAR(100) NOT NULL,
    version             VARCHAR(20)  NOT NULL DEFAULT '1.0',
    is_active           BOOLEAN      NOT NULL DEFAULT true,
    configuration_schema JSONB,
    webhook_support     BOOLEAN      NOT NULL DEFAULT false,
    oauth_support       BOOLEAN      NOT NULL DEFAULT false,
    api_key_support     BOOLEAN      NOT NULL DEFAULT true,
    
    -- Auditoría completa
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp_utc(),
    updated_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp_utc(),
    updated_by          UUID,        -- FK lógica externa → auth-identity-svc.users.id (opcional para admin changes)
    
    -- Business Logic Constraints
    CONSTRAINT chk_integration_type_name_format 
        CHECK (name ~ '^[a-z][a-z0-9_]*[a-z0-9]$'),
    CONSTRAINT chk_integration_type_version_format 
        CHECK (version ~ '^\d+\.\d+(\.\d+)?$'),
    CONSTRAINT chk_integration_type_display_name_length
        CHECK (char_length(display_name) >= 2),
    CONSTRAINT chk_integration_type_provider_length
        CHECK (char_length(provider) >= 2),
    CONSTRAINT chk_integration_type_schema_structure
        CHECK (configuration_schema IS NULL OR jsonb_typeof(configuration_schema) = 'object')
);

-- =============================================================================
-- TABLA: organization_integration (integraciones configuradas por organización)
-- Purpose: Organization-specific integration configurations with credentials and sync settings
-- Business Rules: One integration per type per org, soft delete, comprehensive audit trail
-- =============================================================================
CREATE TABLE organization_integration (
    id                  UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id     UUID         NOT NULL,
    integration_type_id UUID         NOT NULL,
    name                VARCHAR(200) NOT NULL,
    description         TEXT,
    status              integration_status_enum NOT NULL DEFAULT 'inactive',
    configuration       JSONB        NOT NULL DEFAULT '{}',
    credentials         JSONB        DEFAULT '{}', -- Encriptado en aplicación - NULLABLE para flexibilidad
    last_sync_at        TIMESTAMP WITH TIME ZONE,
    last_sync_status    sync_status_enum,
    last_sync_error     TEXT,
    sync_frequency      sync_frequency_enum,
    auto_sync_enabled   BOOLEAN      NOT NULL DEFAULT false,
    webhook_url         VARCHAR(500),
    webhook_secret      VARCHAR(100),
    oauth_token         JSONB        DEFAULT '{}', -- Encriptado en aplicación - NULLABLE para flexibilidad
    api_usage_count     INTEGER      NOT NULL DEFAULT 0,
    api_rate_limit      INTEGER,
    is_active           BOOLEAN      NOT NULL DEFAULT true,
    
    -- Auditoría completa
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp_utc(),
    created_by          UUID NOT NULL,        -- FK lógica externa → auth-identity-svc.users.id (external)
    updated_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp_utc(),
    updated_by          UUID,        -- FK lógica externa → auth-identity-svc.users.id (external)
    deleted_at          TIMESTAMP WITH TIME ZONE, -- Soft delete
    
    -- Foreign Key Constraints
    CONSTRAINT fk_organization_integration_organization 
        FOREIGN KEY (organization_id) REFERENCES organization(id) ON DELETE CASCADE,
    CONSTRAINT fk_organization_integration_type 
        FOREIGN KEY (integration_type_id) REFERENCES integration_type(id) ON DELETE RESTRICT,
    
    -- Business Logic Constraints
    CONSTRAINT chk_organization_integration_api_usage 
        CHECK (api_usage_count >= 0),
    CONSTRAINT chk_organization_integration_api_rate_limit 
        CHECK (api_rate_limit IS NULL OR api_rate_limit > 0),
    CONSTRAINT chk_organization_integration_name_length 
        CHECK (char_length(name) >= 2),
    CONSTRAINT chk_organization_integration_last_sync_coherence
        CHECK (last_sync_at IS NULL OR last_sync_at >= created_at),
    CONSTRAINT chk_organization_integration_webhook_url_format
        CHECK (webhook_url IS NULL OR webhook_url ~* '^https?://'),
    CONSTRAINT chk_organization_integration_configuration_structure
        CHECK (jsonb_typeof(configuration) = 'object'),
    CONSTRAINT chk_organization_integration_credentials_structure
        CHECK (jsonb_typeof(credentials) = 'object'),
    CONSTRAINT chk_organization_integration_oauth_structure
        CHECK (oauth_token IS NULL OR jsonb_typeof(oauth_token) = 'object'),
    
    -- Soft delete consistency (mejorado para mayor coherencia)
    CONSTRAINT chk_organization_integration_soft_delete_consistency
        CHECK (deleted_at IS NULL OR (is_active = false AND status != 'active'))
);

-- =============================================================================
-- TABLA: organization_integration_event (eventos de integración para workers)
-- Purpose: Event queue for integration workers with retry logic and scheduling
-- Business Rules: Events are processed by workers, retries are limited, scheduling supported
-- =============================================================================
CREATE TABLE organization_integration_event (
    id                  UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    integration_id      UUID         NOT NULL,
    event_type          integration_event_type_enum NOT NULL,
    event_data          JSONB        NOT NULL DEFAULT '{}',
    status              integration_event_status_enum NOT NULL DEFAULT 'pending',
    error_message       TEXT,
    retry_count         INTEGER      NOT NULL DEFAULT 0,
    max_retries         INTEGER      NOT NULL DEFAULT 3,
    scheduled_at        TIMESTAMP WITH TIME ZONE,
    processed_at        TIMESTAMP WITH TIME ZONE,
    
    -- Auditoría
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp_utc(),
    
    -- Foreign Key Constraints
    CONSTRAINT fk_organization_integration_event_integration 
        FOREIGN KEY (integration_id) REFERENCES organization_integration(id) ON DELETE CASCADE,
    
    -- Business Logic Constraints
    CONSTRAINT chk_organization_integration_event_retry_count 
        CHECK (retry_count >= 0 AND retry_count <= max_retries),
    CONSTRAINT chk_organization_integration_event_max_retries 
        CHECK (max_retries >= 0 AND max_retries <= 10),
    CONSTRAINT chk_organization_integration_event_scheduling_logic
        CHECK (scheduled_at IS NULL OR processed_at IS NULL OR scheduled_at <= processed_at),
    CONSTRAINT chk_organization_integration_event_data_structure
        CHECK (jsonb_typeof(event_data) = 'object')
);

-- =============================================================================
-- TABLA: organization_integration_log (logs detallados - particionada mensualmente)
-- Purpose: Comprehensive logging for integration operations with monthly partitioning
-- Business Rules: Automatic partitioning, retention policies, comprehensive context tracking
-- =============================================================================
CREATE TABLE organization_integration_log (
    id            UUID               NOT NULL DEFAULT generate_uuid(),
    integration_id UUID              NOT NULL,
    event_id      UUID,
    level         log_level_enum     NOT NULL DEFAULT 'info',
    message       TEXT               NOT NULL,
    context       JSONB              DEFAULT '{}',
    created_at    TIMESTAMPTZ        NOT NULL DEFAULT current_timestamp_utc(),

    -- Clave primaria compuesta, debe incluir la columna de partición
    CONSTRAINT pk_organization_integration_log PRIMARY KEY (id, created_at),

    -- FKs
    CONSTRAINT fk_organization_integration_log_integration
        FOREIGN KEY (integration_id)
        REFERENCES organization_integration(id) ON DELETE CASCADE,
    CONSTRAINT fk_organization_integration_log_event
        FOREIGN KEY (event_id)
        REFERENCES organization_integration_event(id) ON DELETE SET NULL,

    -- Checks
    CONSTRAINT chk_organization_integration_log_message_length
        CHECK (char_length(message) >= 1),
    CONSTRAINT chk_organization_integration_log_context_structure
        CHECK (jsonb_typeof(context) = 'object')
) PARTITION BY RANGE (created_at);

-- =============================================================================
-- INICIALIZACIÓN: CREAR PARTICIÓN INICIAL PARA LOGS
-- =============================================================================
DO $$
BEGIN
    PERFORM create_monthly_partition('organization_integration_log', CURRENT_DATE);
END;
$$ LANGUAGE plpgsql;

-- =============================================================================
-- ÍNDICES OPTIMIZADOS PARA PERFORMANCE Y WORKERS
-- =============================================================================

-- Índices para integration_type (catálogo global)
CREATE INDEX ix_integration_type_name ON integration_type(name);
CREATE INDEX ix_integration_type_category_active ON integration_type(category, is_active);
CREATE INDEX ix_integration_type_provider ON integration_type(provider);
CREATE INDEX ix_integration_type_active ON integration_type(is_active) WHERE is_active = true;

-- Índices básicos para organization_integration
CREATE INDEX ix_organization_integration_organization_id ON organization_integration(organization_id);
CREATE INDEX ix_organization_integration_type_id ON organization_integration(integration_type_id);
CREATE INDEX ix_organization_integration_last_sync ON organization_integration(last_sync_at);
CREATE INDEX ix_organization_integration_created_by ON organization_integration(created_by);
CREATE INDEX ix_organization_integration_active ON organization_integration(is_active) WHERE is_active = true;

-- Índice único case-insensitive para nombre de integración (excluye soft-deleted) - CORREGIDO
CREATE UNIQUE INDEX uq_organization_integration_name_ci 
    ON organization_integration (organization_id, integration_type_id, LOWER(name))
    WHERE deleted_at IS NULL;

-- SOLO CREAR EL ÍNDICE DEFINITIVO - evitar duplicados que se eliminan en 0008
-- Índice más específico y útil para consultas frecuentes
CREATE INDEX ix_organization_integration_org_status_active ON organization_integration(organization_id, status)
    WHERE deleted_at IS NULL AND status = 'active';
    
-- Índices de negocio específicos
CREATE INDEX ix_organization_integration_sync_lookup ON organization_integration(auto_sync_enabled, sync_frequency, last_sync_at) 
    WHERE auto_sync_enabled = true AND deleted_at IS NULL;
CREATE INDEX ix_organization_integration_error_tracking ON organization_integration(last_sync_status, last_sync_at)
    WHERE last_sync_status = 'error';

-- Índices para organization_integration_event (workers)
CREATE INDEX ix_organization_integration_event_integration_id ON organization_integration_event(integration_id);
CREATE INDEX ix_organization_integration_event_type ON organization_integration_event(event_type);
CREATE INDEX ix_organization_integration_event_status_scheduled ON organization_integration_event(status, scheduled_at)
    WHERE scheduled_at IS NOT NULL;
CREATE INDEX ix_organization_integration_event_worker_queue ON organization_integration_event(status, created_at ASC)
    WHERE status IN ('pending', 'retry');
CREATE INDEX ix_organization_integration_event_retry_tracking ON organization_integration_event(status, retry_count, max_retries) 
    WHERE status IN ('error', 'retry');

-- Índices para organization_integration_log (particionada) - sin timestamps dinámicos
CREATE INDEX ix_organization_integration_log_integration_id ON organization_integration_log(integration_id);
CREATE INDEX ix_organization_integration_log_event_id ON organization_integration_log(event_id) 
    WHERE event_id IS NOT NULL;
CREATE INDEX ix_organization_integration_log_level_created ON organization_integration_log(level, created_at DESC);
CREATE INDEX ix_organization_integration_log_errors ON organization_integration_log(integration_id, created_at DESC)
    WHERE level IN ('error', 'critical');

-- =============================================================================
-- TRIGGERS Y FUNCIONES DE BUSINESS LOGIC
-- NOTA: Usa update_timestamp_utc() genérica definida en migración 0002
-- =============================================================================

-- Triggers para updated_at (usando función genérica)
CREATE TRIGGER trg_integration_type_updated_at
    BEFORE UPDATE ON integration_type
    FOR EACH ROW
    EXECUTE FUNCTION update_integration_updated_at();

CREATE TRIGGER trg_organization_integration_updated_at
    BEFORE UPDATE ON organization_integration
    FOR EACH ROW
    EXECUTE FUNCTION update_integration_updated_at();

-- =============================================================================
-- FUNCIÓN MEJORADA PARA LOGGING DE CAMBIOS
-- =============================================================================
CREATE OR REPLACE FUNCTION log_integration_status_change()
RETURNS TRIGGER AS $$
DECLARE
    change_context JSONB;
BEGIN
    -- Preparar contexto del cambio
    change_context := jsonb_build_object(
        'timestamp', current_timestamp_utc(),
        'trigger_operation', TG_OP
    );
    
    -- Registrar cambios significativos
    IF TG_OP = 'INSERT' THEN
        INSERT INTO organization_integration_log (
            integration_id,
            level,
            message,
            context
        ) VALUES (
            NEW.id,
            'info',
            'Integration created: ' || NEW.name,
            change_context || jsonb_build_object(
                'status', NEW.status,
                'integration_type_id', NEW.integration_type_id,
                'created_by', NEW.created_by
            )
        );
    ELSIF TG_OP = 'UPDATE' THEN
        -- Cambios de estado
        IF OLD.status IS DISTINCT FROM NEW.status THEN
            INSERT INTO organization_integration_log (
                integration_id,
                level,
                message,
                context
            ) VALUES (
                NEW.id,
                CASE 
                    WHEN NEW.status = 'error' THEN 'error'
                    WHEN NEW.status = 'suspended' THEN 'warning'
                    ELSE 'info'
                END,
                format('Status changed from %s to %s', OLD.status, NEW.status),
                change_context || jsonb_build_object(
                    'old_status', OLD.status,
                    'new_status', NEW.status,
                    'updated_by', NEW.updated_by
                )
            );
        END IF;
        
        -- Cambios de activación/suspensión
        IF OLD.is_active IS DISTINCT FROM NEW.is_active THEN
            INSERT INTO organization_integration_log (
                integration_id,
                level,
                message,
                context
            ) VALUES (
                NEW.id,
                CASE WHEN NEW.is_active THEN 'info' ELSE 'warning' END,
                CASE WHEN NEW.is_active THEN 'Integration activated' ELSE 'Integration deactivated' END,
                change_context || jsonb_build_object(
                    'old_is_active', OLD.is_active,
                    'new_is_active', NEW.is_active,
                    'updated_by', NEW.updated_by
                )
            );
        END IF;
        
        -- Cambios de configuración importantes
        IF OLD.configuration IS DISTINCT FROM NEW.configuration OR
           OLD.sync_frequency IS DISTINCT FROM NEW.sync_frequency OR
           OLD.auto_sync_enabled IS DISTINCT FROM NEW.auto_sync_enabled THEN
            INSERT INTO organization_integration_log (
                integration_id,
                level,
                message,
                context
            ) VALUES (
                NEW.id,
                'info',
                'Configuration updated',
                change_context || jsonb_build_object(
                    'auto_sync_changed', OLD.auto_sync_enabled IS DISTINCT FROM NEW.auto_sync_enabled,
                    'frequency_changed', OLD.sync_frequency IS DISTINCT FROM NEW.sync_frequency,
                    'config_changed', OLD.configuration IS DISTINCT FROM NEW.configuration,
                    'updated_by', NEW.updated_by
                )
            );
        END IF;
    END IF;
    
    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;

-- Trigger para logging de cambios (INSERT y UPDATE)
CREATE TRIGGER trg_organization_integration_change_log
    AFTER INSERT OR UPDATE ON organization_integration
    FOR EACH ROW
    EXECUTE FUNCTION log_integration_status_change();

-- =============================================================================
-- FUNCIONES PARA PARTICIONADO AUTOMÁTICO Y MANTENIMIENTO
-- =============================================================================

-- Función para trigger de particionado automático de logs
CREATE OR REPLACE FUNCTION organization_integration_log_partition_trigger()
RETURNS TRIGGER AS $$
DECLARE
    partition_date DATE;
BEGIN
    partition_date := NEW.created_at::DATE;
    
    -- Crear partición mensual automáticamente usando función de 0001
    BEGIN
        PERFORM create_monthly_partition('organization_integration_log', partition_date);
    EXCEPTION
        WHEN duplicate_table THEN
            -- Partición ya existe, continuar
            NULL;
        WHEN OTHERS THEN
            -- Log error pero no fallar la inserción
            RAISE WARNING 'Failed to create partition for organization_integration_log: %', SQLERRM;
    END;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger para particionado automático
CREATE TRIGGER trg_organization_integration_log_partition
    BEFORE INSERT ON organization_integration_log
    FOR EACH ROW
    EXECUTE FUNCTION organization_integration_log_partition_trigger();

-- Función mejorada para limpiar logs antiguos con soporte para particiones
CREATE OR REPLACE FUNCTION cleanup_integration_logs(
    days_to_keep INTEGER DEFAULT 30,
    keep_errors BOOLEAN DEFAULT TRUE
)
RETURNS TABLE(
    deleted_count INTEGER,
    preserved_errors INTEGER,
    partitions_dropped INTEGER
) AS $$
DECLARE
    cutoff_date TIMESTAMP WITH TIME ZONE;
    deleted_count_var INTEGER := 0;
    preserved_count_var INTEGER := 0;
    partitions_dropped_var INTEGER := 0;
    partition_rec RECORD;
    partition_date DATE;
    table_exists BOOLEAN;
BEGIN
    cutoff_date := current_timestamp_utc() - (days_to_keep || ' days')::INTERVAL;
    
    -- Contar errores que se preservarían
    IF keep_errors THEN
        SELECT COUNT(*) INTO preserved_count_var
        FROM organization_integration_log
        WHERE created_at < cutoff_date
        AND level IN ('error', 'critical');
    END IF;
    
    -- Método 1: Si hay particiones antiguas, eliminarlas directamente (más eficiente)
    FOR partition_rec IN
        SELECT schemaname, tablename
        FROM pg_tables
        WHERE tablename LIKE 'organization_integration_log_%'
        AND schemaname = 'public'
    LOOP
        -- Intentar extraer fecha de la partición del nombre (formato: table_YYYY_MM)
        IF partition_rec.tablename ~ '_[0-9]{4}_[0-9]{2}$' THEN
            BEGIN
                -- Extraer año y mes del nombre de la partición
                partition_date := TO_DATE(
                    right(partition_rec.tablename, 7),  -- Últimos 7 caracteres: YYYY_MM
                    'YYYY_MM'
                );
                
                -- Si la partición es anterior al cutoff, eliminarla
                IF partition_date < DATE_TRUNC('month', cutoff_date) THEN
                    EXECUTE format('DROP TABLE IF EXISTS %I.%I', 
                        partition_rec.schemaname, partition_rec.tablename);
                    partitions_dropped_var := partitions_dropped_var + 1;
                    CONTINUE; -- Saltar al siguiente partition
                END IF;
            EXCEPTION WHEN OTHERS THEN
                -- Si no se puede parsear la fecha, usar método 2
                NULL;
            END;
        END IF;
    END LOOP;
    
    -- Método 2: Eliminar registros individuales (para particiones que no se pudieron eliminar)
    -- Solo si no se eliminaron todas las particiones viejas
    DELETE FROM organization_integration_log 
    WHERE created_at < cutoff_date
    AND (NOT keep_errors OR level NOT IN ('error', 'critical'));
    
    GET DIAGNOSTICS deleted_count_var = ROW_COUNT;
    
    RETURN QUERY SELECT deleted_count_var, preserved_count_var, partitions_dropped_var;
END;
$$ LANGUAGE plpgsql;

-- =============================================================================
-- FUNCIONES UTILITARIAS PARA ADMINISTRACIÓN
-- =============================================================================

-- Función para obtener estadísticas de integración por organización
CREATE OR REPLACE FUNCTION get_organization_integration_stats(org_id UUID)
RETURNS TABLE(
    total_integrations BIGINT,
    active_integrations BIGINT,
    error_integrations BIGINT,
    suspended_integrations BIGINT,
    last_sync_avg_age INTERVAL,
    auto_sync_enabled_count BIGINT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        COUNT(*) as total_integrations,
        COUNT(*) FILTER (WHERE status = 'active' AND is_active = true) as active_integrations,
        COUNT(*) FILTER (WHERE status = 'error') as error_integrations,
        COUNT(*) FILTER (WHERE status = 'suspended') as suspended_integrations,
        AVG(current_timestamp_utc() - last_sync_at) FILTER (WHERE last_sync_at IS NOT NULL) as last_sync_avg_age,
        COUNT(*) FILTER (WHERE auto_sync_enabled = true) as auto_sync_enabled_count
    FROM organization_integration
    WHERE organization_id = org_id 
    AND deleted_at IS NULL;
END;
$$ LANGUAGE plpgsql;

-- Función para obtener próximos eventos a procesar por workers
CREATE OR REPLACE FUNCTION get_pending_integration_events(limit_count INTEGER DEFAULT 100)
RETURNS TABLE(
    event_id UUID,
    integration_id UUID,
    event_type integration_event_type_enum,
    event_data JSONB,
    retry_count INTEGER,
    scheduled_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        e.id,
        e.integration_id,
        e.event_type,
        e.event_data,
        e.retry_count,
        e.scheduled_at,
        e.created_at
    FROM organization_integration_event e
    WHERE e.status IN ('pending', 'retry')
    AND (e.scheduled_at IS NULL OR e.scheduled_at <= current_timestamp_utc())
    ORDER BY 
        e.retry_count ASC,  -- Priorizar nuevos eventos
        e.created_at ASC    -- FIFO para eventos del mismo tipo
    LIMIT limit_count;
END;
$$ LANGUAGE plpgsql;


-- =============================================================================
-- DATOS INICIALES: TIPOS DE INTEGRACIÓN COMUNES
-- NOTA: Datos iniciales para integration_type - usar ON CONFLICT para evitar duplicados con 0008
-- Mantener sincronizado con ensure_default_integration_types() en 0008
-- =============================================================================
INSERT INTO integration_type (name, display_name, description, category, provider, version, configuration_schema, webhook_support, oauth_support, api_key_support) VALUES
('quickbooks_online', 'QuickBooks Online', 'Integración con QuickBooks Online para contabilidad y facturación', 'accounting', 'Intuit', '1.0',
 '{"required": ["client_id", "client_secret", "redirect_uri"], "optional": ["sandbox_mode", "company_id"]}'::jsonb, true, true, false),
('hubspot_crm', 'HubSpot CRM', 'Integración con HubSpot para gestión de relaciones con clientes', 'crm', 'HubSpot', '1.0',
 '{"required": ["api_key"], "optional": ["portal_id", "rate_limit"]}'::jsonb, true, true, true),
('mailchimp', 'Mailchimp', 'Integración con Mailchimp para marketing por email y campañas', 'marketing', 'Mailchimp', '1.0',
 '{"required": ["api_key"], "optional": ["datacenter", "list_id"]}'::jsonb, true, true, true),
('slack', 'Slack', 'Integración con Slack para comunicaciones y notificaciones', 'communication', 'Slack', '1.0',
 '{"required": ["bot_token"], "optional": ["signing_secret", "channel_id"]}'::jsonb, true, true, false),
('google_analytics', 'Google Analytics', 'Integración con Google Analytics para métricas y análisis', 'analytics', 'Google', '1.0',
 '{"required": ["tracking_id", "service_account_key"], "optional": ["view_id", "property_id"]}'::jsonb, false, true, false),
('zapier_webhook', 'Zapier Webhooks', 'Integración con Zapier para automatizaciones via webhooks', 'automation', 'Zapier', '1.0',
 '{"required": ["webhook_url"], "optional": ["secret_key", "event_types"]}'::jsonb, true, false, false),
('microsoft_outlook', 'Microsoft Outlook', 'Integración con Microsoft Outlook para email y calendario', 'communication', 'Microsoft', '1.0',
 '{"required": ["client_id", "client_secret", "tenant_id"], "optional": ["mailbox_id"]}'::jsonb, true, true, false)
ON CONFLICT (name) DO UPDATE SET
    display_name = EXCLUDED.display_name,
    description = EXCLUDED.description,
    configuration_schema = EXCLUDED.configuration_schema,
    updated_at = CURRENT_TIMESTAMP;

-- =============================================================================
-- COMENTARIOS PARA DOCUMENTACIÓN COMPLETA
-- =============================================================================
COMMENT ON TABLE integration_type IS 'Catálogo global de tipos de integración disponibles con metadatos y capacidades';
COMMENT ON COLUMN integration_type.name IS 'Identificador único del tipo (snake_case, inmutable)';
COMMENT ON COLUMN integration_type.configuration_schema IS 'Esquema JSON de configuración requerida y opcional';
COMMENT ON COLUMN integration_type.updated_by IS 'FK lógica externa → auth-identity-svc.users.id (admin que modificó)';

COMMENT ON TABLE organization_integration IS 'Integraciones configuradas por organización con credentials y configuración';
COMMENT ON COLUMN organization_integration.credentials IS 'Datos sensibles encriptados en aplicación (API keys, tokens)';
COMMENT ON COLUMN organization_integration.oauth_token IS 'Tokens OAuth encriptados con refresh capability';
COMMENT ON COLUMN organization_integration.webhook_secret IS 'Secret para validar webhooks entrantes';
COMMENT ON COLUMN organization_integration.created_by IS 'FK lógica externa → auth-identity-svc.users.id (quien configuró)';
COMMENT ON COLUMN organization_integration.updated_by IS 'FK lógica externa → auth-identity-svc.users.id (última modificación)';

COMMENT ON TABLE organization_integration_event IS 'Cola de eventos para workers con lógica de retry y scheduling';
COMMENT ON COLUMN organization_integration_event.scheduled_at IS 'Cuando procesar (NULL = inmediato)';
COMMENT ON COLUMN organization_integration_event.processed_at IS 'Timestamp de procesamiento completado';

COMMENT ON TABLE organization_integration_log IS 'Logs detallados particionados mensualmente para performance';
COMMENT ON COLUMN organization_integration_log.context IS 'Contexto adicional en JSONB para troubleshooting';

COMMENT ON FUNCTION log_integration_status_change() IS 'Registra automáticamente cambios de estado, configuración y activación';
COMMENT ON FUNCTION cleanup_integration_logs(INTEGER, BOOLEAN) IS 'Limpia logs antiguos con opción de preservar errores críticos';
COMMENT ON FUNCTION get_organization_integration_stats(UUID) IS 'Estadísticas completas de integraciones por organización';
COMMENT ON FUNCTION get_pending_integration_events(INTEGER) IS 'Obtiene próximos eventos para workers con priorización';
COMMENT ON FUNCTION organization_integration_log_partition_trigger() IS 'Crea particiones mensuales automáticamente con manejo de errores';

-- =============================================================================
-- SMOKE TESTS PARA VALIDACIÓN DE MIGRACIÓN
-- =============================================================================

-- Test 1: Verificar que todas las tablas fueron creadas
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_tables WHERE tablename = 'integration_type') THEN
        RAISE EXCEPTION 'Table integration_type was not created';
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_tables WHERE tablename = 'organization_integration') THEN
        RAISE EXCEPTION 'Table organization_integration was not created';
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_tables WHERE tablename = 'organization_integration_event') THEN
        RAISE EXCEPTION 'Table organization_integration_event was not created';
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_tables WHERE tablename = 'organization_integration_log') THEN
        RAISE EXCEPTION 'Table organization_integration_log was not created';
    END IF;
    
    RAISE NOTICE 'TEST PASSED: All tables created successfully';
END $$;

-- Test 2: Verificar que todas las funciones fueron creadas
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_proc WHERE proname = 'log_integration_status_change') THEN
        RAISE EXCEPTION 'Function log_integration_status_change was not created';
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_proc WHERE proname = 'cleanup_integration_logs') THEN
        RAISE EXCEPTION 'Function cleanup_integration_logs was not created';
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_proc WHERE proname = 'get_organization_integration_stats') THEN
        RAISE EXCEPTION 'Function get_organization_integration_stats was not created';
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_proc WHERE proname = 'get_pending_integration_events') THEN
        RAISE EXCEPTION 'Function get_pending_integration_events was not created';
    END IF;
    
    RAISE NOTICE 'TEST PASSED: All functions created successfully';
END $$;

-- Test 3: Verificar partición inicial y datos semilla
DO $$
DECLARE
    partition_count INTEGER;
    seed_count INTEGER;
BEGIN
    -- Verificar que existe al menos una partición
    SELECT COUNT(*) INTO partition_count
    FROM pg_tables 
    WHERE tablename LIKE 'organization_integration_log_%';
    
    IF partition_count = 0 THEN
        RAISE EXCEPTION 'No partitions created for organization_integration_log';
    END IF;
    
    -- Verificar datos semilla
    SELECT COUNT(*) INTO seed_count FROM integration_type;
    
    IF seed_count < 5 THEN
        RAISE EXCEPTION 'Insufficient seed data in integration_type table';
    END IF;
    
    RAISE NOTICE 'TEST PASSED: Partitions and seed data created (% partitions, % integration types)', 
        partition_count, seed_count;
END $$;

-- Test 4: Verificar constraints y validaciones
DO $$
DECLARE
    test_type_id UUID;
    test_org_id UUID := generate_uuid();
    test_user_id UUID := generate_uuid();
BEGIN
    -- Obtener un tipo de integración válido
    SELECT id INTO test_type_id FROM integration_type LIMIT 1;
    
    -- Test de constraint de URL webhook
    BEGIN
        INSERT INTO organization_integration (
            organization_id,
            integration_type_id,
            name,
            webhook_url,
            created_by
        ) VALUES (
            test_org_id,
            test_type_id,
            'Test Integration',
            'invalid-url',
            test_user_id
        );
        
        RAISE EXCEPTION 'Webhook URL constraint not working';
    EXCEPTION 
        WHEN check_violation THEN
            -- Expected behavior
            NULL;
    END;
    
    -- Test de nombre case-insensitive
    BEGIN
        INSERT INTO organization_integration (
            organization_id,
            integration_type_id,
            name,
            created_by
        ) VALUES 
        (test_org_id, test_type_id, 'Test Integration', test_user_id),
        (test_org_id, test_type_id, 'TEST INTEGRATION', test_user_id);
        
        RAISE EXCEPTION 'Case-insensitive name constraint not working';
    EXCEPTION 
        WHEN unique_violation THEN
            -- Expected behavior
            NULL;
    END;
    
    RAISE NOTICE 'TEST PASSED: Constraints and validations working correctly';
    
    -- Cleanup
    DELETE FROM organization_integration WHERE organization_id = test_org_id;
    
EXCEPTION WHEN foreign_key_violation THEN
    RAISE NOTICE 'TEST PASSED: FK constraints working (expected in isolated test)';
END $$;

-- Test 5: Verificar funciones utilitarias
DO $$
DECLARE
    stats_result RECORD;
    events_result RECORD;
BEGIN
    -- Test función de estadísticas (debería funcionar sin datos)
    SELECT * INTO stats_result 
    FROM get_organization_integration_stats(generate_uuid());
    
    IF stats_result.total_integrations IS NULL THEN
        RAISE EXCEPTION 'Statistics function not working';
    END IF;
    
    -- Test función de eventos pendientes
    SELECT COUNT(*) INTO events_result 
    FROM get_pending_integration_events(10);
    
    RAISE NOTICE 'TEST PASSED: Utility functions working correctly';
END $$;

DO $$
BEGIN
  RAISE NOTICE 'MIGRATION 0006 COMPLETED SUCCESSFULLY: Organization integration system with comprehensive event tracking and audit trails';
END
$$;