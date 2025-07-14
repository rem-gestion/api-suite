-- =============================================
-- Migration: 0003_create_structure.up.sql
-- Description: Create organization structure tables (branches and roles)
-- Author: System Generated
-- Date: 2024
-- =============================================

-- =============================================================================
-- TABLA: organization_branch (sucursales/oficinas)
-- =============================================================================
CREATE TABLE organization_branch (
    id              UUID         PRIMARY KEY DEFAULT generate_uuid(),
    organization_id UUID         NOT NULL,
    display_name    VARCHAR(120) NOT NULL,
    address_id      UUID,        -- FK lógica externa → address-svc.address.id (no enforce DB)
    phone           VARCHAR(32),
    email           VARCHAR(160),
    is_main         BOOLEAN      NOT NULL DEFAULT false,
    
    -- Auditoría completa
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp_utc(),
    created_by      UUID NOT NULL,        -- FK a auth-identity-svc users.id (external)
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp_utc(),
    updated_by      UUID,        -- FK a auth-identity-svc users.id (external)
    deleted_at      TIMESTAMP WITH TIME ZONE, -- Soft delete
    
    -- Constraints mejorados
    CONSTRAINT fk_organization_branch_organization 
        FOREIGN KEY (organization_id) REFERENCES organization(id) ON DELETE CASCADE,
    CONSTRAINT chk_organization_branch_display_name_length 
        CHECK (char_length(display_name) >= 2),
    CONSTRAINT chk_organization_branch_valid_email 
        CHECK (email IS NULL OR is_valid_email(email)),
    CONSTRAINT chk_organization_branch_phone_format 
        CHECK (phone IS NULL OR phone ~ '^[\+]?[0-9\s\-\(\)\.]{7,20}$'),
    CONSTRAINT chk_organization_branch_main_not_deleted 
        CHECK (deleted_at IS NULL OR is_main = false) -- Evita soft-delete de sucursal principal sin reemplazo
);

-- Índices con nombres explícitos y mejorados
CREATE INDEX ix_organization_branch_organization_active ON organization_branch(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX ix_organization_branch_is_main ON organization_branch(organization_id, is_main) WHERE deleted_at IS NULL AND is_main = true;
CREATE INDEX ix_organization_branch_address ON organization_branch(address_id) WHERE address_id IS NOT NULL;
CREATE INDEX ix_organization_branch_email ON organization_branch(email) WHERE email IS NOT NULL;
CREATE INDEX ix_organization_branch_created_by ON organization_branch(created_by);

-- =============================================================================
-- TABLA: organization_role (roles internos por organización)
-- =============================================================================
CREATE TABLE organization_role (
    id              UUID         PRIMARY KEY DEFAULT generate_uuid(),
    organization_id UUID         NOT NULL,
    name            VARCHAR(64)  NOT NULL,
    description     TEXT,
    is_default      BOOLEAN      NOT NULL DEFAULT false, -- Para roles predeterminados
    
    -- Auditoría completa
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp_utc(),
    created_by      UUID NOT NULL,        -- FK a auth-identity-svc users.id (external)
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp_utc(),
    updated_by      UUID,        -- FK a auth-identity-svc users.id (external)
    deleted_at      TIMESTAMP WITH TIME ZONE, -- Soft delete
    
    -- Constraints mejorados
    CONSTRAINT fk_organization_role_organization 
        FOREIGN KEY (organization_id) REFERENCES organization(id) ON DELETE CASCADE,
    CONSTRAINT chk_organization_role_name_length 
        CHECK (char_length(name) >= 2 AND char_length(name) <= 64),
    CONSTRAINT chk_organization_role_name_format 
        CHECK (name ~ '^[a-zA-Z][a-zA-Z0-9\s\-_]*[a-zA-Z0-9]$')
);

-- Índices con nombres explícitos y mejorados
CREATE INDEX ix_organization_role_organization_active ON organization_role(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX ix_organization_role_name ON organization_role(organization_id, name) WHERE deleted_at IS NULL;
CREATE INDEX ix_organization_role_is_default ON organization_role(organization_id, is_default) WHERE deleted_at IS NULL AND is_default = true;
CREATE INDEX ix_organization_role_created_by ON organization_role(created_by);

-- =============================================================================
-- CONSTRAINTS ÚNICOS CON NOMBRES EXPLÍCITOS (OPTIMIZADOS)
-- =============================================================================

-- Constraint para asegurar una sola sucursal principal por organización
CREATE UNIQUE INDEX uq_organization_branch_one_main_per_org 
ON organization_branch(organization_id) 
WHERE is_main = true AND deleted_at IS NULL;

-- Constraint para asegurar nombre de rol único por organización (case-insensitive only)
CREATE UNIQUE INDEX uq_organization_role_name_case_insensitive 
ON organization_role(organization_id, LOWER(name)) 
WHERE deleted_at IS NULL;

-- =============================================================================
-- TRIGGERS Y FUNCIONES DE BUSINESS LOGIC
-- =============================================================================

-- Trigger para actualizar updated_at automáticamente
CREATE TRIGGER trg_organization_branch_updated_at
    BEFORE UPDATE ON organization_branch
    FOR EACH ROW
    EXECUTE FUNCTION update_organization_updated_at();

CREATE TRIGGER trg_organization_role_updated_at
    BEFORE UPDATE ON organization_role
    FOR EACH ROW
    EXECUTE FUNCTION update_organization_updated_at();

-- Función para prevenir eliminar la única sucursal principal
-- CORREGIDO: Ahora incluye cambios de is_main para evitar bypass del check
CREATE OR REPLACE FUNCTION prevent_delete_only_main_branch()
RETURNS TRIGGER AS $$
BEGIN
    -- Aplicar tanto para soft-delete como hard-delete, y cambio de is_main
    IF (TG_OP = 'DELETE' OR TG_OP = 'UPDATE') THEN
        -- Para UPDATE, verificar si se está quitando is_main o soft-deleting
        IF TG_OP = 'UPDATE' THEN
            -- Si no hay cambios relevantes, permitir
            IF (OLD.is_main = NEW.is_main) AND (OLD.deleted_at = NEW.deleted_at OR (OLD.deleted_at IS NULL AND NEW.deleted_at IS NULL)) THEN
                RETURN NEW;
            END IF;
            -- Si se está cambiando is_main de true a false o soft-deleting una main branch
            IF NOT ((OLD.is_main = true AND NEW.is_main = false) OR (OLD.is_main = true AND NEW.deleted_at IS NOT NULL AND OLD.deleted_at IS NULL)) THEN
                RETURN NEW;
            END IF;
        END IF;
        
        -- Verificar si hay otras sucursales en la organización (para DELETE o UPDATE de main branch)
        IF OLD.is_main = true THEN
            IF (SELECT COUNT(*) FROM organization_branch 
                WHERE organization_id = OLD.organization_id 
                AND id != OLD.id 
                AND deleted_at IS NULL) > 0 THEN
                IF TG_OP = 'DELETE' THEN
                    RAISE EXCEPTION 'Cannot hard-delete the main branch while other branches exist. Set another branch as main first.';
                ELSE
                    RAISE EXCEPTION 'Cannot modify main branch while other branches exist. Set another branch as main first.';
                END IF;
            END IF;
        END IF;
    END IF;
    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_organization_branch_prevent_delete_only_main
    BEFORE UPDATE OF deleted_at, is_main ON organization_branch
    FOR EACH ROW
    EXECUTE FUNCTION prevent_delete_only_main_branch();

-- Trigger para prevenir hard-delete de sucursal principal también
CREATE TRIGGER trg_organization_branch_prevent_hard_delete_main
    BEFORE DELETE ON organization_branch
    FOR EACH ROW
    EXECUTE FUNCTION prevent_delete_only_main_branch();

-- Función para prevenir eliminar roles predeterminados
CREATE OR REPLACE FUNCTION prevent_delete_default_roles()
RETURNS TRIGGER AS $$
BEGIN
    -- Aplicar tanto para soft-delete como hard-delete
    IF OLD.is_default = true THEN
        IF TG_OP = 'DELETE' THEN
            RAISE EXCEPTION 'Cannot hard-delete default role: %', OLD.name;
        ELSE
            RAISE EXCEPTION 'Cannot soft-delete default role: %', OLD.name;
        END IF;
    END IF;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_organization_role_prevent_delete_default
    BEFORE UPDATE OF deleted_at ON organization_role
    FOR EACH ROW
    WHEN (NEW.deleted_at IS NOT NULL AND OLD.deleted_at IS NULL)
    EXECUTE FUNCTION prevent_delete_default_roles();

-- Trigger para prevenir hard-delete de roles default también
CREATE TRIGGER trg_organization_role_prevent_hard_delete_default
    BEFORE DELETE ON organization_role
    FOR EACH ROW
    EXECUTE FUNCTION prevent_delete_default_roles();

-- NOTA: El índice único uq_organization_branch_one_main_per_org es suficiente
-- para garantizar una sola sucursal principal por organización.
-- Eliminar trigger duplicado manage_main_branch_change para evitar DML extra
-- y confiar en el índice parcial único que ya previene duplicados de manera más eficiente.

-- CREATE OR REPLACE FUNCTION manage_main_branch_change() -- ELIMINADO: redundante con índice único
-- CREATE TRIGGER trg_organization_branch_manage_main_change -- ELIMINADO: redundante con índice único

-- =============================================================================
-- FUNCIONES PARA DATOS INICIALES
-- =============================================================================

-- Función para crear roles por defecto
CREATE OR REPLACE FUNCTION create_default_organization_roles(org_id UUID, creator_id UUID)
RETURNS VOID AS $$
BEGIN
    INSERT INTO organization_role (organization_id, name, description, is_default, created_by) VALUES
    (org_id, 'Admin', 'Full administrative access to all organization features', true, creator_id),
    (org_id, 'Manager', 'Management access to most organization features', true, creator_id),
    (org_id, 'Agent', 'Real estate agent with property and client management access', true, creator_id),
    (org_id, 'Assistant', 'Limited access for administrative support tasks', true, creator_id)
    ON CONFLICT (organization_id, name) DO NOTHING;
END;
$$ LANGUAGE plpgsql;

-- Función para crear sucursal principal por defecto
CREATE OR REPLACE FUNCTION create_default_main_branch(org_id UUID, org_name VARCHAR, creator_id UUID)
RETURNS VOID AS $$
BEGIN
    INSERT INTO organization_branch (organization_id, display_name, is_main, created_by) VALUES
    (org_id, org_name || ' - Oficina Principal', true, creator_id)
    ON CONFLICT DO NOTHING;
END;
$$ LANGUAGE plpgsql;

-- =============================================================================
-- COMENTARIOS PARA DOCUMENTACIÓN
-- =============================================================================
COMMENT ON TABLE organization_branch IS 'Sucursales u oficinas de la inmobiliaria';
COMMENT ON COLUMN organization_branch.address_id IS 'FK lógica externa → address-svc.address.id (no enforce DB)';
COMMENT ON COLUMN organization_branch.is_main IS 'Solo una sucursal puede ser principal por organización';
COMMENT ON COLUMN organization_branch.created_by IS 'FK lógica externa → auth-identity-svc.users.id (no enforce DB)';
COMMENT ON COLUMN organization_branch.updated_by IS 'FK lógica externa → auth-identity-svc.users.id (no enforce DB)';

COMMENT ON TABLE organization_role IS 'Roles internos personalizables por organización';
COMMENT ON COLUMN organization_role.is_default IS 'Roles predeterminados (Admin, Manager, Agent, Assistant) no eliminables';
COMMENT ON COLUMN organization_role.name IS 'Nombre único por organización (case-insensitive)';
COMMENT ON COLUMN organization_role.created_by IS 'FK lógica externa → auth-identity-svc.users.id (no enforce DB)';
COMMENT ON COLUMN organization_role.updated_by IS 'FK lógica externa → auth-identity-svc.users.id (no enforce DB)';

COMMENT ON FUNCTION create_default_organization_roles(UUID, UUID) IS 'Crea roles por defecto para una nueva organización';
COMMENT ON FUNCTION create_default_main_branch(UUID, VARCHAR, UUID) IS 'Crea sucursal principal por defecto para una nueva organización';

-- =============================================================================
-- SMOKE TEST PARA VALIDACIÓN EN CI
-- =============================================================================
/*
Simple smoke test para verificar que la migración 0003 se aplicó correctamente:

-- Verificar que las tablas existen
SELECT 1 FROM organization_branch LIMIT 0;
SELECT 1 FROM organization_role LIMIT 0;

-- Verificar que los constraints funcionan
-- (Esto fallará correctamente si hay problemas)
INSERT INTO organization_branch (organization_id, display_name, is_main, created_by) 
VALUES (generate_uuid(), 'Test Branch', true, generate_uuid());

-- Verificar funciones de seeds
SELECT create_default_organization_roles(generate_uuid(), generate_uuid());

-- Verificar índices únicos
SELECT indexname FROM pg_indexes 
WHERE tablename IN ('organization_branch', 'organization_role') 
AND indexname LIKE 'uq_%';
*/
