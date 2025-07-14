-- =============================================
-- Migration: 0002_create_organization_core.up.sql
-- Description: Create core organization tables with ENUM types and improved constraints
-- Author: System Generated  
-- Date: 2024
-- =============================================

-- =============================================================================
-- TABLA PRINCIPAL: organization
-- =============================================================================
CREATE TABLE organization (
    id                  UUID         PRIMARY KEY DEFAULT generate_uuid(),
    display_name        VARCHAR(120) NOT NULL,
    logo_url            TEXT,
    fiscal_address_id   UUID,        -- FK lógica externa → address-svc.address.id (no enforce DB)
    matricula           VARCHAR(32),
    status              organization_status_enum NOT NULL DEFAULT 'active',
    
    -- Auditoría completa
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp_utc(),
    created_by          UUID NOT NULL,  -- FK a auth-identity-svc users.id (external) - OBLIGATORIO
    updated_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp_utc(),
    updated_by          UUID,        -- FK a auth-identity-svc users.id (external)
    deleted_at          TIMESTAMP WITH TIME ZONE, -- Soft delete
    
    -- Constraints mejorados
    CONSTRAINT chk_organization_display_name_length CHECK (char_length(display_name) >= 2),
    CONSTRAINT chk_organization_matricula_format CHECK (matricula IS NULL OR matricula ~ '^[A-Z0-9\-]+$'),
    CONSTRAINT chk_organization_valid_created_by CHECK (is_valid_uuid(created_by::TEXT)),
    -- IMPORTANTE: Este constraint obliga a que status='deleted' cuando deleted_at IS NOT NULL
    -- Asegurar en la capa de servicio que TODAS las operaciones de soft-delete actualicen
    -- ambas columnas en el mismo UPDATE statement para evitar violaciones de constraint
    CONSTRAINT chk_organization_soft_delete_status CHECK (
        (deleted_at IS NULL AND status != 'deleted') OR 
        (deleted_at IS NOT NULL AND status = 'deleted')
    )
);

-- Índices con nombres explícitos y optimizados
CREATE INDEX ix_organization_status_active ON organization(status) WHERE deleted_at IS NULL;
CREATE INDEX ix_organization_display_name_search ON organization(display_name) WHERE deleted_at IS NULL;
CREATE INDEX ix_organization_created_by ON organization(created_by);
CREATE INDEX ix_organization_fiscal_address ON organization(fiscal_address_id) WHERE fiscal_address_id IS NOT NULL;
CREATE INDEX ix_organization_deleted_at ON organization(deleted_at) WHERE deleted_at IS NOT NULL;

-- =============================================================================
-- TABLA: organization_settings (KV store)
-- =============================================================================
CREATE TABLE organization_settings (
    organization_id     UUID         NOT NULL,
    setting_key         VARCHAR(64)  NOT NULL,
    setting_value       JSONB,
    
    -- Auditoría mejorada con trazabilidad
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp_utc(),
    updated_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp_utc(),
    updated_by          UUID,        -- FK lógica externa → auth-identity-svc.users.id (no enforce DB)
    
    -- PK compuesta con nombre explícito
    CONSTRAINT pk_organization_settings PRIMARY KEY (organization_id, setting_key),
    CONSTRAINT fk_organization_settings_organization 
        FOREIGN KEY (organization_id) REFERENCES organization(id) ON DELETE CASCADE,
    
    -- Constraints de validación
    CONSTRAINT chk_organization_settings_key_format 
        CHECK (setting_key ~ '^[a-z][a-z0-9_]*[a-z0-9]$'),
    CONSTRAINT chk_organization_settings_key_length 
        CHECK (char_length(setting_key) >= 2)
);

-- Índices con nombres explícitos
CREATE INDEX ix_organization_settings_updated_at ON organization_settings(updated_at);
CREATE INDEX ix_organization_settings_key ON organization_settings(setting_key);
CREATE INDEX ix_organization_settings_organization_id ON organization_settings(organization_id);

-- =============================================================================
-- TABLA: organization_owner (dueños de la inmobiliaria)
-- =============================================================================
CREATE TABLE organization_owner (
    id              UUID         PRIMARY KEY DEFAULT generate_uuid(),
    organization_id UUID         NOT NULL,
    owner_type      owner_type_enum NOT NULL,
    owner_id        UUID         NOT NULL, -- FK lógica externa → person-svc.person.id (no enforce DB)
    ownership_percentage DECIMAL(5,2) NOT NULL, -- Porcentaje de propiedad (0.00-100.00) - obligatorio incluso si soft-deleted
    
    -- Auditoría completa
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp_utc(),
    created_by      UUID NOT NULL,        -- FK a auth-identity-svc users.id (external)
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp_utc(),
    updated_by      UUID,
    deleted_at      TIMESTAMP WITH TIME ZONE, -- Soft delete
    
    -- Constraints mejorados
    CONSTRAINT fk_organization_owner_organization 
        FOREIGN KEY (organization_id) REFERENCES organization(id) ON DELETE CASCADE,
    CONSTRAINT chk_organization_owner_ownership_percentage 
        CHECK (ownership_percentage >= 0 AND ownership_percentage <= 100),
    CONSTRAINT chk_organization_owner_valid_owner_id 
        CHECK (is_valid_uuid(owner_id::TEXT)),
    
    -- Constraint único mejorado
    CONSTRAINT uq_organization_owner_unique 
        UNIQUE (organization_id, owner_type, owner_id) DEFERRABLE INITIALLY DEFERRED
);

-- Índices con nombres explícitos y mejorados
CREATE INDEX ix_organization_owner_organization_active ON organization_owner(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX ix_organization_owner_lookup ON organization_owner(owner_type, owner_id) WHERE deleted_at IS NULL;
CREATE INDEX ix_organization_owner_created_by ON organization_owner(created_by);
CREATE INDEX ix_organization_owner_ownership_percentage ON organization_owner(ownership_percentage) WHERE ownership_percentage IS NOT NULL;

-- =============================================================================
-- TRIGGERS Y FUNCIONES
-- =============================================================================

-- Función genérica para actualizar updated_at (consolidada para reutilización)
CREATE OR REPLACE FUNCTION update_timestamp_utc()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = current_timestamp_utc();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Alias para compatibilidad con triggers existentes (idempotente)
CREATE OR REPLACE FUNCTION update_organization_updated_at()
RETURNS TRIGGER AS $$
BEGIN 
    RETURN update_timestamp_utc(); 
END;
$$ LANGUAGE plpgsql;

-- Triggers para updated_at (usando función genérica)
    CREATE TRIGGER trg_organization_updated_at
        BEFORE UPDATE ON organization
        FOR EACH ROW
        EXECUTE FUNCTION update_timestamp_utc();

    CREATE TRIGGER trg_organization_settings_updated_at
        BEFORE UPDATE ON organization_settings
        FOR EACH ROW
        EXECUTE FUNCTION update_timestamp_utc();

    CREATE TRIGGER trg_organization_owner_updated_at
        BEFORE UPDATE ON organization_owner
        FOR EACH ROW
        EXECUTE FUNCTION update_timestamp_utc();

-- Función para prevenir DELETE directo (forzar soft delete)
-- NOTA: Esta función solo protege la tabla 'organization'
-- IMPORTANTE: Confirmar que la capa de servicio nunca hace DELETE en tablas hijas
-- (organization_branch, employees, etc.) o crear triggers análogos en migraciones posteriores
CREATE OR REPLACE FUNCTION prevent_hard_delete_organization()
RETURNS TRIGGER AS $$
BEGIN
    -- Solo aplicar en operaciones DELETE del esquema público
    IF TG_OP = 'DELETE' AND TG_TABLE_SCHEMA = 'public' THEN
        RAISE EXCEPTION 'Direct DELETE not allowed on %. Use soft delete by setting deleted_at timestamp.', TG_TABLE_NAME;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Trigger para prevenir DELETE directo en production
CREATE TRIGGER trg_organization_prevent_hard_delete
    BEFORE DELETE ON organization
    FOR EACH ROW
    EXECUTE FUNCTION prevent_hard_delete_organization();

-- Función para validar porcentajes de propiedad
CREATE OR REPLACE FUNCTION validate_ownership_percentages()
RETURNS TRIGGER AS $$
DECLARE
    total_percentage DECIMAL(5,2);
BEGIN
    -- Calcular total actual excluyendo el registro que se está modificando
    SELECT COALESCE(SUM(ownership_percentage), 0)
    INTO total_percentage
    FROM organization_owner
    WHERE organization_id = NEW.organization_id
    AND deleted_at IS NULL
    AND id != COALESCE(NEW.id, '00000000-0000-0000-0000-000000000000'::UUID);
    
    -- Verificar que el total no exceda 100%
    IF (total_percentage + NEW.ownership_percentage) > 100 THEN
        RAISE EXCEPTION 'Total ownership percentage cannot exceed 100%%. Current total: %%, trying to add: %%', 
            total_percentage, NEW.ownership_percentage;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger para validar porcentajes
CREATE TRIGGER trg_organization_owner_validate_percentages
    BEFORE INSERT OR UPDATE ON organization_owner
    FOR EACH ROW
    EXECUTE FUNCTION validate_ownership_percentages();

-- =============================================================================
-- COMENTARIOS PARA DOCUMENTACIÓN
-- =============================================================================
COMMENT ON TABLE organization IS 'Inmobiliarias - entidad principal del dominio';
COMMENT ON COLUMN organization.fiscal_address_id IS 'FK lógica externa → address-svc.address.id (no enforce DB)';
COMMENT ON COLUMN organization.created_by IS 'FK lógica externa → auth-identity-svc.users.id (no enforce DB)';
COMMENT ON COLUMN organization.updated_by IS 'FK lógica externa → auth-identity-svc.users.id (no enforce DB)';

COMMENT ON TABLE organization_settings IS 'Configuraciones clave-valor por organización';
COMMENT ON COLUMN organization_settings.setting_key IS 'Clave de configuración en formato snake_case';
COMMENT ON COLUMN organization_settings.setting_value IS 'Valor JSON de la configuración (usar ::jsonb para tipos correctos)';

COMMENT ON TABLE organization_owner IS 'Propietarios de las inmobiliarias con porcentajes de propiedad';
COMMENT ON COLUMN organization_owner.owner_id IS 'FK lógica externa → person-svc.person.id (no enforce DB)';
COMMENT ON COLUMN organization_owner.ownership_percentage IS 'Porcentaje de propiedad obligatorio (0-100%)';
COMMENT ON COLUMN organization_owner.created_by IS 'FK lógica externa → auth-identity-svc.users.id (no enforce DB)';
COMMENT ON COLUMN organization_owner.updated_by IS 'FK lógica externa → auth-identity-svc.users.id (no enforce DB)';

-- =============================================================================
-- DATOS INICIALES DE CONFIGURACIÓN
-- =============================================================================

-- Función para crear configuraciones por defecto
CREATE OR REPLACE FUNCTION create_default_organization_settings(org_id UUID)
RETURNS VOID AS $$
BEGIN
    INSERT INTO organization_settings (organization_id, setting_key, setting_value) VALUES
    (org_id, 'enable_notifications', true::jsonb),
    (org_id, 'timezone', '"UTC"'::jsonb),
    (org_id, 'language', '"es"'::jsonb),
    (org_id, 'currency', '"USD"'::jsonb),
    (org_id, 'date_format', '"dd/mm/yyyy"'::jsonb),
    (org_id, 'enable_public_listings', true::jsonb),
    (org_id, 'max_images_per_property', 10::jsonb),
    (org_id, 'enable_agent_commissions', true::jsonb),
    (org_id, 'default_commission_percentage', 3.0::jsonb)
    ON CONFLICT (organization_id, setting_key) DO NOTHING;
END;
$$ LANGUAGE plpgsql;

-- Trigger para crear configuraciones por defecto al crear organización
CREATE OR REPLACE FUNCTION create_organization_defaults()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM create_default_organization_settings(NEW.id);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_organization_create_defaults
    AFTER INSERT ON organization
    FOR EACH ROW
    EXECUTE FUNCTION create_organization_defaults();

COMMENT ON FUNCTION create_default_organization_settings(UUID) IS 'Crea configuraciones por defecto para una nueva organización';

-- =============================================================================
-- SMOKE TEST PARA VALIDACIÓN EN CI
-- =============================================================================
/*
Simple smoke test para verificar que la migración se aplicó correctamente:

-- Verificar que las tablas existen y los ENUMs están disponibles
SELECT 1 FROM organization_status_enum LIMIT 1;
SELECT 1 FROM owner_type_enum LIMIT 1;

-- Verificar que las funciones están disponibles
SELECT generate_uuid() IS NOT NULL;
SELECT current_timestamp_utc() IS NOT NULL;

-- Verificar constraints básicos
INSERT INTO organization (display_name, created_by) 
VALUES ('Test Org', generate_uuid()) 
RETURNING id;
*/
