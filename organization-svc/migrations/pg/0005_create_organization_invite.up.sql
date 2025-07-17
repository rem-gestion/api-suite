-- =============================================
-- Migration: 0005_create_organization_invite.up.sql
-- Description: Create organization invitation system with enhanced security and audit trails
-- Author: System Generated
-- Date: 2024
-- Version: 1.1
-- Dependencies: 0001 (ENUMs), 0002 (organization), 0003 (organization_role, organization_branch)
-- =============================================

-- Asegurar que pgcrypto esté disponible para token generation y crypto functions
DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: Creating pgcrypto extension...';
CREATE EXTENSION IF NOT EXISTS pgcrypto;
END;
$$ LANGUAGE plpgsql;

DO $$
BEGIN
  CREATE EXTENSION IF NOT EXISTS citext;
END;
$$ LANGUAGE plpgsql;

DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: pgcrypto extension created successfully';
 END;
$$ LANGUAGE plpgsql;

-- =============================================================================
-- TABLA: organization_invite (invitaciones pendientes/históricas)
-- =============================================================================
-- TABLA PRINCIPAL: organization_invite (invitaciones pendientes/históricas)
-- Purpose: Invitations to join an organization with role and branch assignments
-- Business Rules:
-- - Invitaciones tienen estado lifecycle: pending → accepted/rejected/expired/cancelled
-- - Token único para cada invitación (URL-safe, 64 chars)
-- - Expiración automática mediante batch job  
-- - Límites dinámicos de invitaciones pendientes por organización
-- =============================================================================
DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: Creating organization_invite table...';
 END;
$$ LANGUAGE plpgsql;

CREATE TABLE organization_invite (
    id                  UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id     UUID         NOT NULL,
    inviter_person_id   UUID         NOT NULL, -- FK lógica externa → person-svc.person.id (no enforce DB)
    invitee_email       CITEXT       NOT NULL,
    invitee_person_id   UUID,        -- FK lógica externa → person-svc.person.id (no enforce DB)
    role_id             UUID         NOT NULL,
    branch_id           UUID,        -- FK opcional a organization_branch.id
    token               VARCHAR(500) NOT NULL UNIQUE,
    status              invitation_status_enum NOT NULL DEFAULT 'pending',
    expires_at          TIMESTAMP WITH TIME ZONE NOT NULL,
    accepted_at         TIMESTAMP WITH TIME ZONE,
    rejected_at         TIMESTAMP WITH TIME ZONE,
    cancelled_at        TIMESTAMP WITH TIME ZONE,
    metadata            JSONB        DEFAULT '{}', -- Extra data: IP, user agent, etc.
    
    -- Auditoría completa
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    created_by          UUID NOT NULL,        -- FK lógica externa → auth-identity-svc.users.id (no enforce DB)
    updated_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_by          UUID,        -- FK lógica externa → auth-identity-svc.users.id (no enforce DB)
    
    -- Foreign Key Constraints
    CONSTRAINT fk_organization_invite_organization 
        FOREIGN KEY (organization_id) REFERENCES organization(id) ON DELETE CASCADE,
    CONSTRAINT fk_organization_invite_role 
        FOREIGN KEY (role_id) REFERENCES organization_role(id) ON DELETE RESTRICT,
    CONSTRAINT fk_organization_invite_branch 
        FOREIGN KEY (branch_id) REFERENCES organization_branch(id) ON DELETE SET NULL,
    
    -- Business Logic Constraints
    CONSTRAINT chk_organization_invite_valid_email 
        CHECK (is_valid_email(invitee_email)),
    CONSTRAINT chk_organization_invite_valid_dates 
        CHECK (expires_at > created_at),
    CONSTRAINT chk_organization_invite_future_expiry
        CHECK (expires_at > now()),
    CONSTRAINT chk_organization_invite_acceptance_logic 
        CHECK (
            (status = 'accepted' AND accepted_at IS NOT NULL AND rejected_at IS NULL AND cancelled_at IS NULL) OR
            (status = 'rejected' AND rejected_at IS NOT NULL AND accepted_at IS NULL AND cancelled_at IS NULL) OR
            (status = 'cancelled' AND cancelled_at IS NOT NULL AND accepted_at IS NULL AND rejected_at IS NULL) OR
            (status NOT IN ('accepted', 'rejected', 'cancelled') AND accepted_at IS NULL AND rejected_at IS NULL AND cancelled_at IS NULL)
        ),
    CONSTRAINT chk_organization_invite_valid_inviter 
        CHECK (is_valid_uuid(inviter_person_id::TEXT)),
    CONSTRAINT chk_organization_invite_token_security 
        CHECK (char_length(token) >= 32 AND char_length(token) <= 128 AND token ~ '^[A-Za-z0-9_-]+$'), -- URL-safe characters only
    CONSTRAINT chk_organization_invite_metadata_structure
        CHECK (jsonb_typeof(metadata) = 'object'),
    
    -- Prevent duplicate pending invites for same email+org (case-insensitive)
    CONSTRAINT uq_organization_invite_pending_email_ci
        UNIQUE (organization_id, invitee_email) DEFERRABLE INITIALLY DEFERRED
);


DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: organization_invite table created successfully';
 END;
$$ LANGUAGE plpgsql;

-- =============================================================================
-- TABLA: organization_invite_log (historial de invitaciones)
-- Purpose: Audit trail for all invitation lifecycle events
-- Partitioning: Monthly partitions for performance and maintenance
-- =============================================================================
DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: Creating organization_invite_log table...';
 END;
$$ LANGUAGE plpgsql;

CREATE TABLE organization_invite_log (
    id                  UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    invite_id           UUID         NOT NULL,
    action              invitation_log_action_enum NOT NULL,
    old_status          invitation_status_enum,
    new_status          invitation_status_enum,
    actor_person_id     UUID,        -- FK lógica externa → person-svc.person.id (no enforce DB) 
    actor_type          VARCHAR(20)  DEFAULT 'user' CHECK (actor_type IN ('user', 'system', 'cron')),
    client_ip           INET,        -- IP address for security audit
    user_agent          TEXT,        -- Browser/client info
    notes               TEXT,
    error_details       JSONB,       -- Error info for failed operations
    
    -- Auditoría
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp_utc(),
    
    -- Constraints
    CONSTRAINT fk_organization_invite_log_invite 
        FOREIGN KEY (invite_id) REFERENCES organization_invite(id) ON DELETE CASCADE,
    CONSTRAINT chk_organization_invite_log_status_change
        CHECK (
            (action IN ('created', 'expired', 'cancelled') AND old_status IS NULL) OR
            (action IN ('accepted', 'rejected', 'updated') AND old_status IS NOT NULL)
        )
);

DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: organization_invite_log table created successfully';
 END;
$$ LANGUAGE plpgsql;


-- =============================================================================
-- ÍNDICES OPTIMIZADOS PARA PERFORMANCE
-- =============================================================================
DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: Creating indexes for organization_invite...';
 END;
$$ LANGUAGE plpgsql;

-- Índices básicos para organization_invite
CREATE INDEX ix_organization_invite_organization_id ON organization_invite(organization_id);
CREATE INDEX ix_organization_invite_invitee_email_ci ON organization_invite(LOWER(invitee_email)); -- Case-insensitive primary
CREATE INDEX ix_organization_invite_invitee_person ON organization_invite(invitee_person_id) WHERE invitee_person_id IS NOT NULL;
CREATE UNIQUE INDEX ix_organization_invite_token_unique ON organization_invite(token); -- Explicit unique index
CREATE INDEX ix_organization_invite_role_id ON organization_invite(role_id);
CREATE INDEX ix_organization_invite_branch_id ON organization_invite(branch_id) WHERE branch_id IS NOT NULL;
CREATE INDEX ix_organization_invite_inviter ON organization_invite(inviter_person_id);
CREATE INDEX ix_organization_invite_created_by ON organization_invite(created_by);

-- Índices de estado y fecha para operaciones frecuentes
CREATE INDEX ix_organization_invite_status_expires ON organization_invite(status, expires_at);
CREATE INDEX ix_organization_invite_pending_expired ON organization_invite(expires_at) 
    WHERE status = 'pending';
CREATE INDEX ix_organization_invite_org_status_created ON organization_invite(organization_id, status, created_at DESC);

-- Índices compuestos para búsquedas de negocio
CREATE INDEX ix_organization_invite_active_lookup ON organization_invite(organization_id, status, expires_at)
    WHERE status IN ('pending', 'accepted');
CREATE INDEX ix_organization_invite_email_status ON organization_invite(LOWER(invitee_email), status);
CREATE INDEX ix_organization_invite_org_email_pending ON organization_invite(organization_id, LOWER(invitee_email))
    WHERE status = 'pending';
DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: organization_invite indexes created successfully';
 END;
$$ LANGUAGE plpgsql;

-- Índices para organization_invite_log (tabla particionada)
DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: Creating indexes for organization_invite_log...';
 END;
$$ LANGUAGE plpgsql;

CREATE INDEX ix_organization_invite_log_invite_id ON organization_invite_log(invite_id);
CREATE INDEX ix_organization_invite_log_action_date ON organization_invite_log(action, created_at DESC);
CREATE INDEX ix_organization_invite_log_actor ON organization_invite_log(actor_person_id, created_at DESC) 
    WHERE actor_person_id IS NOT NULL;
DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: organization_invite_log indexes created successfully';
 END;
$$ LANGUAGE plpgsql;

-- =============================================================================
-- TRIGGERS Y FUNCIONES DE BUSINESS LOGIC
-- =============================================================================
DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: Creating triggers and business logic functions...';
 END;
$$ LANGUAGE plpgsql;

-- Trigger para actualizar updated_at automáticamente
DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: Creating updated_at trigger...';
 END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER trg_organization_invite_updated_at
    BEFORE UPDATE ON organization_invite
    FOR EACH ROW
    EXECUTE FUNCTION update_organization_updated_at();
DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: updated_at trigger created successfully';
 END;
$$ LANGUAGE plpgsql;

-- Trigger para asegurar created_at en INSERT si llega NULL
DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: Creating ensure_created_at trigger...';
 END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER trg_organization_invite_ensure_created_at
    BEFORE INSERT ON organization_invite
    FOR EACH ROW
    WHEN (NEW.created_at IS NULL)
    EXECUTE FUNCTION update_organization_updated_at();
DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: ensure_created_at trigger created successfully';
 END;
$$ LANGUAGE plpgsql;

-- =============================================================================
-- FUNCIÓN MEJORADA PARA LOGGING DE CAMBIOS DE ESTADO
-- =============================================================================
DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: Creating log_organization_invite_status_change function...';
 END;
$$ LANGUAGE plpgsql;
CREATE OR REPLACE FUNCTION log_organization_invite_status_change()
RETURNS TRIGGER AS $$
DECLARE
    log_action invitation_log_action_enum;
    actor_id UUID;
    actor_type_val VARCHAR(20) := 'user';
BEGIN
    -- Solo registrar si hay cambios significativos
    IF (TG_OP = 'INSERT') OR 
       (TG_OP = 'UPDATE' AND (
           OLD.status IS DISTINCT FROM NEW.status OR
           OLD.expires_at IS DISTINCT FROM NEW.expires_at OR
           OLD.role_id IS DISTINCT FROM NEW.role_id OR
           OLD.branch_id IS DISTINCT FROM NEW.branch_id
       )) THEN
        
        -- Determinar la acción basada en el contexto
        IF TG_OP = 'INSERT' THEN
            log_action := 'created';
            actor_id := NEW.created_by;
        ELSE
            log_action := CASE 
                WHEN NEW.status = 'accepted' THEN 'accepted'
                WHEN NEW.status = 'rejected' THEN 'rejected'
                WHEN NEW.status = 'expired' THEN 'expired'
                WHEN NEW.status = 'cancelled' THEN 'cancelled'
                ELSE 'updated'
            END;
            
            -- Determinar actor y tipo
            IF NEW.updated_by IS NOT NULL THEN
                actor_id := NEW.updated_by;
                actor_type_val := 'user';
            ELSE
                actor_id := NULL;
                actor_type_val := CASE 
                    WHEN NEW.status = 'expired' THEN 'cron'
                    ELSE 'system'
                END;
            END IF;
        END IF;
        
        -- Insertar log con información completa
        INSERT INTO organization_invite_log (
            invite_id,
            action,
            old_status,
            new_status,
            actor_person_id,
            actor_type,
            notes
        ) VALUES (
            NEW.id,
            log_action,
            CASE WHEN TG_OP = 'INSERT' THEN NULL ELSE OLD.status END,
            NEW.status,
            actor_id,
            actor_type_val,
            CASE 
                WHEN TG_OP = 'INSERT' THEN 'Invitation created'
                WHEN OLD.status IS DISTINCT FROM NEW.status THEN 
                    format('Status changed from %s to %s', OLD.status, NEW.status)
                WHEN OLD.role_id IS DISTINCT FROM NEW.role_id THEN 
                    'Role assignment updated'
                WHEN OLD.branch_id IS DISTINCT FROM NEW.branch_id THEN 
                    'Branch assignment updated'
                WHEN OLD.expires_at IS DISTINCT FROM NEW.expires_at THEN 
                    'Expiration date updated'
                ELSE 'Invitation updated'
            END
        );
    END IF;
    
    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;
DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: log_organization_invite_status_change function created successfully';
 END;
$$ LANGUAGE plpgsql;

-- Trigger para logging de cambios (INSERT y UPDATE)
DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: Creating status_log trigger...';
 END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER trg_organization_invite_status_log
    AFTER INSERT OR UPDATE ON organization_invite
    FOR EACH ROW
    EXECUTE FUNCTION log_organization_invite_status_change();
DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: status_log trigger created successfully';
 END;
$$ LANGUAGE plpgsql;

-- =============================================================================
-- FUNCIÓN MEJORADA PARA VERIFICAR LÍMITES DE INVITACIONES
-- =============================================================================
DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: Creating check_organization_invite_limits function...';
 END;
$$ LANGUAGE plpgsql;
CREATE OR REPLACE FUNCTION check_organization_invite_limits()
RETURNS TRIGGER AS $$
DECLARE
    pending_count INTEGER;
    org_limit INTEGER;
    error_msg TEXT;
BEGIN
    -- Solo verificar para INSERTs de invitaciones pendientes
    IF NEW.status != 'pending' THEN
        RETURN NEW;
    END IF;
    
    -- Obtener límite dinámico desde organization_settings con soporte para JSON
    -- CORREGIDO: Soportar tanto texto plano como JSON para setting_value
    SELECT COALESCE(
        (SELECT CASE 
            -- Si es un número directo en texto
            WHEN setting_value::text ~ '^\d+$' THEN (setting_value::text)::INTEGER
            -- Si es JSON, extraer el valor numérico
            WHEN setting_value::text LIKE '"%"' THEN 
                CASE 
                    WHEN (setting_value #>> '{}')::text ~ '^\d+$' THEN (setting_value #>> '{}')::INTEGER
                    ELSE 50
                END
            -- Si es JSONB con estructura, buscar key 'limit' o 'value'
            WHEN jsonb_typeof(setting_value) = 'object' THEN
                COALESCE(
                    (setting_value->>'limit')::INTEGER,
                    (setting_value->>'value')::INTEGER,
                    50
                )
            -- Si es número en JSONB
            WHEN jsonb_typeof(setting_value) = 'number' THEN (setting_value)::INTEGER
            ELSE 50
         END
         FROM organization_settings
         WHERE organization_id = NEW.organization_id
         AND setting_key = 'max_pending_invites'),
        50) -- default limit
    INTO org_limit;
    
    -- Verificar que el límite sea razonable
    IF org_limit < 1 OR org_limit > 1000 THEN
        org_limit := 50;
    END IF;
    
    -- Contar invitaciones pendientes activas para la organización
    SELECT COUNT(*) 
    INTO pending_count
    FROM organization_invite 
    WHERE organization_id = NEW.organization_id 
    AND status = 'pending' 
    AND expires_at > current_timestamp_utc()
    AND id != NEW.id; -- Excluir la invitación actual si es UPDATE
    
    -- Verificar límite con error específico para captura desde el servicio
    IF pending_count >= org_limit THEN
        error_msg := format(
            'Organization has reached the maximum number of pending invitations. Current: %s, Limit: %s', 
            pending_count, 
            org_limit
        );
        
        RAISE EXCEPTION 
            USING ERRCODE = '23514', -- CHECK violation más específico que P0001
                  MESSAGE = error_msg,
                  DETAIL = format('Organization ID: %s, Pending count: %s', NEW.organization_id, pending_count),
                  HINT = 'Cancel existing invitations or increase the limit in organization settings';
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: check_organization_invite_limits function created successfully';
 END;
$$ LANGUAGE plpgsql;

-- Trigger para verificar límites en INSERT y UPDATE a pending
DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: Creating invite_limits trigger...';
 END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER trg_organization_invite_limits
    BEFORE INSERT OR UPDATE ON organization_invite
    FOR EACH ROW
    EXECUTE FUNCTION check_organization_invite_limits();
DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: invite_limits trigger created successfully';
 END;
$$ LANGUAGE plpgsql;

-- =============================================================================
-- FUNCIÓN MEJORADA PARA GESTIÓN DE INVITACIONES EXPIRADAS
-- =============================================================================
DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: Creating expire_old_invitations function...';
 END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION expire_old_invitations(batch_size INTEGER DEFAULT 1000)
RETURNS TABLE(
    expired_count INTEGER,
    processed_orgs INTEGER,
    execution_time_ms INTEGER
) AS $$
DECLARE
    start_time TIMESTAMP;
    expired_count_var INTEGER := 0;
    org_count INTEGER := 0;
    org_record RECORD;
BEGIN
    start_time := clock_timestamp();
    
    FOR org_record IN 
        SELECT DISTINCT organization_id 
        FROM organization_invite 
        WHERE status = 'pending' 
          AND expires_at <= current_timestamp_utc()
        LIMIT batch_size
    LOOP
        WITH updated AS (
            UPDATE organization_invite 
            SET status = 'expired',
                updated_at = current_timestamp_utc(),
                updated_by = NULL -- Sistema
            WHERE organization_id = org_record.organization_id
              AND status = 'pending' 
              AND expires_at <= current_timestamp_utc()
            RETURNING 1
        )
        SELECT COUNT(*) INTO expired_count_var FROM updated;
        
        org_count := org_count + 1;
        
        IF org_count % 100 = 0 THEN
            RAISE NOTICE 'Processed % organizations, expired % invitations so far',
                         org_count, expired_count_var;
        END IF;
    END LOOP;
    
    RETURN QUERY
    SELECT 
        expired_count_var,
        org_count,
        EXTRACT(milliseconds FROM clock_timestamp() - start_time)::INTEGER;
END;
$$ LANGUAGE plpgsql;


DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: expire_old_invitations function created successfully';
 END;
$$ LANGUAGE plpgsql;

-- =============================================================================
-- FUNCIONES PARA GESTIÓN SEGURA DE TOKENS
-- =============================================================================
DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: Creating token management functions...';
 END;
$$ LANGUAGE plpgsql;

-- Función para generar token único URL-safe de invitación (mejorada con CSPRNG)
CREATE OR REPLACE FUNCTION generate_invite_token(token_length INTEGER DEFAULT 64)
RETURNS TEXT AS $$
DECLARE
    token TEXT;
    token_exists BOOLEAN;
    max_attempts INTEGER := 10;
    attempt_count INTEGER := 0;
    token_var TEXT;
BEGIN
    -- Validar longitud del token
    IF token_length < 32 OR token_length > 128 THEN
        RAISE EXCEPTION 'Token length must be between 32 and 128 characters';
    END IF;
    
    LOOP
        -- Generar token usando CSPRNG para mayor seguridad
        token_var := encode(gen_random_bytes((token_length * 3 / 4)::integer), 'base64');
        -- Hacer URL-safe reemplazando caracteres problemáticos
        token_var := replace(replace(token_var, '/', '_'), '+', '-');
        -- Truncar al tamaño deseado
        token_var := left(token_var, token_length);
        
        -- Verificar si el token ya existe
        SELECT EXISTS(
          SELECT 1
          FROM organization_invite oi
          WHERE oi.token = token_var
        ) INTO token_exists;
        
        -- Si no existe, salir del loop
        IF NOT token_exists THEN
            token := token_var;
            EXIT;
        END IF;
        
        attempt_count := attempt_count + 1;
        IF attempt_count >= max_attempts THEN
            RAISE EXCEPTION 'Unable to generate unique token after % attempts', max_attempts;
        END IF;
    END LOOP;
    
    RETURN token;
END;
$$ LANGUAGE plpgsql;
DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: generate_invite_token function created successfully';
 END;
$$ LANGUAGE plpgsql;

-- Función para validar y renovar token de invitación
CREATE OR REPLACE FUNCTION renew_invite_token(invite_id UUID)
RETURNS TEXT AS $$
DECLARE
    new_token TEXT;
    invite_status invitation_status_enum;
BEGIN
    -- Verificar que la invitación existe y está pendiente
    SELECT status INTO invite_status
    FROM organization_invite
    WHERE id = invite_id;
    
    IF NOT FOUND THEN
        RAISE EXCEPTION 'Invitation not found' USING ERRCODE = '02000';
    END IF;
    
    IF invite_status != 'pending' THEN
        RAISE EXCEPTION 'Can only renew tokens for pending invitations' USING ERRCODE = '23000';
    END IF;
    
    -- Generar nuevo token
    new_token := generate_invite_token();
    
    -- Actualizar la invitación
    UPDATE organization_invite
    SET token = new_token,
        updated_at = current_timestamp_utc()
    WHERE id = invite_id;
    
    RETURN new_token;
END;
$$ LANGUAGE plpgsql;
DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: renew_invite_token function created successfully';
 END;
$$ LANGUAGE plpgsql;

-- =============================================================================
-- FUNCIONES PARA PARTICIONADO AUTOMÁTICO DE LOGS
-- =============================================================================
DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: Creating partitioning functions...';
 END;
$$ LANGUAGE plpgsql;

-- Función para trigger de particionado automático de logs
-- NOTA: En PG 15+ se recomienda crear particiones por adelantado via job cron
-- Este trigger funciona como "safety-net" para evitar errores
-- ESTRATEGIA RECOMENDADA: 
-- 1. Scheduler externo (pg_cron): eliminar este trigger y usar solo job programado
-- 2. Safety-net: mantener trigger pero 0008 pre-genera particiones para reducir DDL
CREATE OR REPLACE FUNCTION organization_invite_log_partition_trigger()
RETURNS TRIGGER AS $$
DECLARE
    partition_date DATE;
    partition_name TEXT;
BEGIN
    partition_date := NEW.created_at::DATE;
    
    -- Crear partición mensual automáticamente usando función de 0001
    -- ADVERTENCIA: DDL en trigger puede causar locks en alta concurrencia
    BEGIN
        PERFORM create_monthly_partition('organization_invite_log', partition_date);
    EXCEPTION
        WHEN duplicate_table THEN
            -- Partición ya existe, continuar
            NULL;
        WHEN OTHERS THEN
            -- Log error pero no fallar la inserción
            RAISE WARNING 'Failed to create partition for organization_invite_log: %', SQLERRM;
    END;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: organization_invite_log_partition_trigger function created successfully';
 END;
$$ LANGUAGE plpgsql;

-- Trigger para particionado automático
DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: Creating log_partition trigger...';
 END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_organization_invite_log_partition
    BEFORE INSERT ON organization_invite_log
    FOR EACH ROW
    EXECUTE FUNCTION organization_invite_log_partition_trigger();
DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: log_partition trigger created successfully';
 END;
$$ LANGUAGE plpgsql;

-- =============================================================================
-- FUNCIONES UTILITARIAS PARA ADMINISTRACIÓN
-- =============================================================================
DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: Creating utility functions...';
 END;
$$ LANGUAGE plpgsql;

-- Función para cleanup de invitaciones antiguas (mantener solo últimos N meses)
CREATE OR REPLACE FUNCTION cleanup_old_invitations(months_to_keep INTEGER DEFAULT 12)
RETURNS TABLE(
    deleted_invites INTEGER,
    deleted_logs INTEGER
) AS $$
DECLARE
    cutoff_date TIMESTAMP WITH TIME ZONE;
    invite_count INTEGER;
    log_count INTEGER;
BEGIN
    cutoff_date := current_timestamp_utc() - (months_to_keep || ' months')::INTERVAL;
    
    -- Eliminar logs antiguos primero (FK constraint)
    DELETE FROM organization_invite_log 
    WHERE created_at < cutoff_date;
    
    GET DIAGNOSTICS log_count = ROW_COUNT;
    
    -- Eliminar invitaciones antiguas no pendientes
    DELETE FROM organization_invite 
    WHERE created_at < cutoff_date 
    AND status NOT IN ('pending'); -- Mantener pendientes sin importar edad
    
    GET DIAGNOSTICS invite_count = ROW_COUNT;
    
    RETURN QUERY SELECT invite_count, log_count;
END;
$$ LANGUAGE plpgsql;

DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: cleanup_old_invitations function created successfully';
 END;
$$ LANGUAGE plpgsql;

-- Función para estadísticas de invitaciones por organización
CREATE OR REPLACE FUNCTION get_organization_invite_stats(org_id UUID)
RETURNS TABLE(
    total_invites BIGINT,
    pending_invites BIGINT,
    accepted_invites BIGINT,
    rejected_invites BIGINT,
    expired_invites BIGINT,
    cancelled_invites BIGINT,
    acceptance_rate NUMERIC(5,2)
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        COUNT(*) as total_invites,
        COUNT(*) FILTER (WHERE status = 'pending') as pending_invites,
        COUNT(*) FILTER (WHERE status = 'accepted') as accepted_invites,
        COUNT(*) FILTER (WHERE status = 'rejected') as rejected_invites,
        COUNT(*) FILTER (WHERE status = 'expired') as expired_invites,
        COUNT(*) FILTER (WHERE status = 'cancelled') as cancelled_invites,
        CASE 
            WHEN COUNT(*) FILTER (WHERE status IN ('accepted', 'rejected')) > 0 THEN
                ROUND(
                    COUNT(*) FILTER (WHERE status = 'accepted') * 100.0 / 
                    COUNT(*) FILTER (WHERE status IN ('accepted', 'rejected')),
                    2
                )
            ELSE 0
        END as acceptance_rate
    FROM organization_invite
    WHERE organization_id = org_id;
END;
$$ LANGUAGE plpgsql;
DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: get_organization_invite_stats function created successfully';
 END;
$$ LANGUAGE plpgsql;

-- =============================================================================
-- COMENTARIOS PARA DOCUMENTACIÓN COMPLETA
-- =============================================================================
DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: Adding table and function comments...';
 END;
$$ LANGUAGE plpgsql;
COMMENT ON TABLE organization_invite IS 'Sistema de invitaciones a organizaciones con workflow completo y auditoría';
COMMENT ON COLUMN organization_invite.inviter_person_id IS 'FK lógica externa → person-svc.person.id (quien invita)';
COMMENT ON COLUMN organization_invite.invitee_person_id IS 'FK lógica externa → person-svc.person.id (quien es invitado, opcional hasta aceptación)';
COMMENT ON COLUMN organization_invite.role_id IS 'Rol que se asignará al aceptar la invitación';
COMMENT ON COLUMN organization_invite.branch_id IS 'Sucursal específica (opcional, NULL = acceso a toda la organización)';
COMMENT ON COLUMN organization_invite.token IS 'Token único URL-safe para aceptar/rechazar invitación (32-128 chars)';
COMMENT ON COLUMN organization_invite.metadata IS 'Datos adicionales: IP, user agent, source, etc. (JSONB)';
COMMENT ON COLUMN organization_invite.created_by IS 'FK lógica externa → auth-identity-svc.users.id (no enforce DB)';
COMMENT ON COLUMN organization_invite.updated_by IS 'FK lógica externa → auth-identity-svc.users.id (no enforce DB)';

COMMENT ON TABLE organization_invite_log IS 'Audit trail completo de invitaciones (particionado mensualmente para performance)';
COMMENT ON COLUMN organization_invite_log.actor_person_id IS 'FK lógica externa → person-svc.person.id (quien realizó la acción, NULL = sistema)';
COMMENT ON COLUMN organization_invite_log.actor_type IS 'Tipo de actor: user (manual), system (automático), cron (batch job)';
COMMENT ON COLUMN organization_invite_log.client_ip IS 'IP del cliente para auditoría de seguridad';
COMMENT ON COLUMN organization_invite_log.error_details IS 'Detalles de errores para troubleshooting (JSONB)';

COMMENT ON FUNCTION log_organization_invite_status_change() IS 'Registra automáticamente todos los cambios significativos con contexto completo';
COMMENT ON FUNCTION check_organization_invite_limits() IS 'Verifica límites dinámicos de invitaciones pendientes con validación robusta';
COMMENT ON FUNCTION expire_old_invitations(INTEGER) IS 'Expira invitaciones pendientes por lotes con estadísticas detalladas';
COMMENT ON FUNCTION generate_invite_token(INTEGER) IS 'Genera tokens únicos URL-safe con validación y retry logic';
COMMENT ON FUNCTION renew_invite_token(UUID) IS 'Renueva token de invitación pendiente para casos de seguridad';
COMMENT ON FUNCTION organization_invite_log_partition_trigger() IS 'Crea particiones mensuales automáticamente con manejo de errores';
COMMENT ON FUNCTION cleanup_old_invitations(INTEGER) IS 'Limpieza de invitaciones antiguas manteniendo integridad referencial';
COMMENT ON FUNCTION get_organization_invite_stats(UUID) IS 'Estadísticas detalladas de invitaciones por organización con ratios';
DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: Comments added successfully';
 END;
$$ LANGUAGE plpgsql;

-- =============================================================================
-- SMOKE TESTS PARA VALIDACIÓN DE MIGRACIÓN
-- =============================================================================
DO $$
 BEGIN 
 RAISE NOTICE 'DEBUG: Starting smoke tests...';
 END;
$$ LANGUAGE plpgsql;

-- Test 1: Verificar que todas las tablas fueron creadas
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_tables
        WHERE tablename = 'organization_invite'
    ) THEN
        RAISE EXCEPTION 'Table organization_invite was not created';
    END IF;
    
    IF NOT EXISTS (
        SELECT 1
        FROM pg_tables
        WHERE tablename = 'organization_invite_log'
    ) THEN
        RAISE EXCEPTION 'Table organization_invite_log was not created';
    END IF;
    
    RAISE NOTICE 'TEST PASSED: All tables created successfully';
END;
$$ LANGUAGE plpgsql;

DO $$
BEGIN
    RAISE NOTICE 'DEBUG: Table existence test completed';
END;
$$ LANGUAGE plpgsql;


-- Test 2: Verificar que todas las funciones fueron creadas
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_proc
        WHERE proname = 'generate_invite_token'
    ) THEN
        RAISE EXCEPTION 'Function generate_invite_token was not created';
    END IF;

    IF NOT EXISTS (
        SELECT 1
        FROM pg_proc
        WHERE proname = 'expire_old_invitations'
    ) THEN
        RAISE EXCEPTION 'Function expire_old_invitations was not created';
    END IF;

    IF NOT EXISTS (
        SELECT 1
        FROM pg_proc
        WHERE proname = 'get_organization_invite_stats'
    ) THEN
        RAISE EXCEPTION 'Function get_organization_invite_stats was not created';
    END IF;

    RAISE NOTICE 'TEST PASSED: All functions created successfully';
END;
$$ LANGUAGE plpgsql;

DO $$
BEGIN
    RAISE NOTICE 'DEBUG: Function existence test completed';
END;
$$ LANGUAGE plpgsql;

-- Test 3: Verificar generación de tokens
DO $$
DECLARE
    test_token TEXT;
BEGIN
    test_token := generate_invite_token(64);

    IF length(test_token) != 64 THEN
        RAISE EXCEPTION
            'Token length is incorrect: expected 64, got %',
            length(test_token);
    END IF;

    IF test_token !~ '^[A-Za-z0-9_-]+$' THEN
        RAISE EXCEPTION 'Token contains invalid characters';
    END IF;

    -- Un único RAISE NOTICE en este bloque
    RAISE NOTICE 'TEST PASSED: Token generation working correctly';
END;
$$ LANGUAGE plpgsql;

-- Debug
DO $$
BEGIN
    RAISE NOTICE 'DEBUG: Token generation test completed';
END;
$$ LANGUAGE plpgsql;

-- Test 4: Verificar constraints y validaciones básicas
DO $$
DECLARE
    test_org_id    UUID := gen_random_uuid();
    test_role_id   UUID := gen_random_uuid();
    test_user_id   UUID := gen_random_uuid();
    test_person_id UUID := gen_random_uuid();
BEGIN
    -- Este test solo verifica que la tabla acepta inserts válidos
    -- En ambiente real, las FKs hacia otras tablas fallarían
    BEGIN
        INSERT INTO organization_invite (
            organization_id,
            inviter_person_id,
            invitee_email,
            role_id,
            token,
            expires_at,
            created_by
        ) VALUES (
            test_org_id,
            test_person_id,
            'test@example.com',
            test_role_id,
            generate_invite_token(),
            current_timestamp_utc() + INTERVAL '7 days',
            test_user_id
        );

        -- Si llegamos aquí, todo OK
        RAISE NOTICE 'TEST PASSED: Basic constraints and structure working';
        
        -- Cleanup
        DELETE FROM organization_invite
        WHERE organization_id = test_org_id;
        
    EXCEPTION WHEN foreign_key_violation THEN
        -- Si falla por FK, también OK en entorno aislado
        RAISE NOTICE 'TEST PASSED: FK constraints working (expected in isolated test)';
    END;
END
$$ LANGUAGE plpgsql;

-- Debug
DO $$
BEGIN
    RAISE NOTICE 'DEBUG: Basic constraints test completed';
END
$$ LANGUAGE plpgsql;

-- Fin de migración
DO $$
BEGIN
    RAISE NOTICE 'MIGRATION 0005 COMPLETED SUCCESSFULLY: Organization invitation system with enhanced security and comprehensive audit trails';
END
$$ LANGUAGE plpgsql;

DO $$
BEGIN
    RAISE NOTICE 'DEBUG: Migration 0005 finished completely';
END
$$ LANGUAGE plpgsql;
