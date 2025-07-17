-- =============================================
-- Migration: 0004_create_employee.up.sql
-- Description: Create employee management tables with improved constraints and ENUM types
-- Author: System Generated
-- Date: 2024
-- =============================================

-- =============================================================================
-- TABLA: employees (vínculo user ↔ organization)
-- =============================================================================
CREATE TABLE employees (
    id               UUID         PRIMARY KEY DEFAULT generate_uuid(),
    organization_id  UUID         NOT NULL,
    user_id          UUID         NOT NULL, -- FK lógica externa → auth-identity-svc.users.id (no enforce DB)
    person_id        UUID,        -- FK lógica externa → person-svc.person.id (no enforce DB) - Agregado para evitar conflicto futuro
    primary_role_id  UUID,        -- FK → organization_role.id (ON DELETE SET NULL)
    status           employee_status_enum NOT NULL DEFAULT 'active',
    hired_at         DATE,
    fired_at         DATE,
    
    -- Auditoría completa
    created_at       TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp_utc(),
    created_by       UUID NOT NULL,        -- FK lógica externa → auth-identity-svc.users.id (no enforce DB)
    updated_at       TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp_utc(),
    updated_by       UUID,        -- FK lógica externa → auth-identity-svc.users.id (no enforce DB)
    deleted_at       TIMESTAMP WITH TIME ZONE, -- Soft delete
    
    -- Constraints mejorados
    CONSTRAINT fk_employees_organization 
        FOREIGN KEY (organization_id) REFERENCES organization(id) ON DELETE CASCADE,
    CONSTRAINT fk_employees_primary_role 
        FOREIGN KEY (primary_role_id) REFERENCES organization_role(id) ON DELETE SET NULL,
    CONSTRAINT chk_employees_valid_user_id 
        CHECK (is_valid_uuid(user_id::TEXT)),
    CONSTRAINT chk_employees_hire_fire_dates 
        CHECK (fired_at IS NULL OR hired_at IS NULL OR fired_at >= hired_at),
    CONSTRAINT chk_employees_hired_not_future 
        CHECK (hired_at IS NULL OR hired_at <= CURRENT_DATE),
    CONSTRAINT chk_employees_status_fire_consistency 
        CHECK ((status = 'terminated' AND fired_at IS NOT NULL) OR (status != 'terminated' AND fired_at IS NULL)),
    
    -- Constraint único mejorado: un usuario solo puede estar una vez por organización (considerando soft delete)
    CONSTRAINT uq_employees_user_organization 
        UNIQUE (user_id, organization_id) DEFERRABLE INITIALLY DEFERRED
);

-- Índices con nombres explícitos y optimizados (eliminando duplicados desde 0004)
CREATE INDEX ix_employees_organization_active ON employees(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX ix_employees_user_active ON employees(user_id) WHERE deleted_at IS NULL;
CREATE INDEX ix_employees_person_active ON employees(person_id) WHERE person_id IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX ix_employees_primary_role ON employees(primary_role_id) WHERE primary_role_id IS NOT NULL;
CREATE INDEX ix_employees_created_by ON employees(created_by);
-- NOTA: Solo mantener el constraint único básico, 0008 se encarga de la optimización
-- No crear índices únicos adicionales que se eliminarán en 0008

-- =============================================================================
-- TABLA: employee_roles (N-a-N entre employees y roles)
-- =============================================================================
CREATE TABLE employee_roles (
    id               UUID         PRIMARY KEY DEFAULT generate_uuid(),
    employee_id      UUID         NOT NULL,
    role_id          UUID         NOT NULL,
    is_primary       BOOLEAN      NOT NULL DEFAULT false,
    
    -- Auditoría mejorada
    created_at       TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp_utc(),
    created_by       UUID NOT NULL,        -- FK lógica externa → auth-identity-svc.users.id (no enforce DB)
    deleted_at       TIMESTAMP WITH TIME ZONE, -- Soft delete para roles
    
    -- Constraints mejorados
    CONSTRAINT fk_employee_roles_employee 
        FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
    CONSTRAINT fk_employee_roles_role 
        FOREIGN KEY (role_id) REFERENCES organization_role(id) ON DELETE CASCADE,
    
    -- Constraint único mejorado
    CONSTRAINT uq_employee_roles_unique 
        UNIQUE (employee_id, role_id) DEFERRABLE INITIALLY DEFERRED
);

-- Índices con nombres explícitos y mejorados
CREATE INDEX ix_employee_roles_employee_active ON employee_roles(employee_id) WHERE deleted_at IS NULL;
CREATE INDEX ix_employee_roles_role_active ON employee_roles(role_id) WHERE deleted_at IS NULL;
CREATE INDEX ix_employee_roles_primary ON employee_roles(employee_id, is_primary) WHERE deleted_at IS NULL AND is_primary = true;
CREATE INDEX ix_employee_roles_created_by ON employee_roles(created_by);

-- Constraint único para asegurar solo un rol primario por empleado
CREATE UNIQUE INDEX uq_employee_roles_one_primary_per_employee 
ON employee_roles(employee_id) 
WHERE is_primary = true AND deleted_at IS NULL;

-- =============================================================================
-- TRIGGERS Y FUNCIONES DE BUSINESS LOGIC
-- =============================================================================

-- Trigger para actualizar updated_at automáticamente
CREATE TRIGGER trg_employees_updated_at
    BEFORE UPDATE ON employees
    FOR EACH ROW
    EXECUTE FUNCTION update_organization_updated_at();

-- Función para sincronizar primary_role_id en employees con employee_roles
-- CORREGIDO: Prevenir recursión usando pg_trigger_depth() en lugar de session_replication_role
CREATE OR REPLACE FUNCTION sync_employee_primary_role()
RETURNS TRIGGER AS $$
DECLARE
    current_depth INTEGER;
BEGIN
    -- Prevenir recursión infinita usando pg_trigger_depth (disponible desde PG 9.0+)
    current_depth := pg_trigger_depth();
    IF current_depth > 1 THEN
        RETURN COALESCE(NEW, OLD);
    END IF;
    
    IF TG_OP = 'INSERT' OR TG_OP = 'UPDATE' THEN
        -- Solo procesar si realmente cambió is_primary
        IF TG_OP = 'UPDATE' AND OLD.is_primary = NEW.is_primary THEN
            RETURN NEW;
        END IF;
        
        -- Si se marca un rol como primario, actualizar employees.primary_role_id
        IF NEW.is_primary = true THEN
            UPDATE employees 
            SET primary_role_id = NEW.role_id,
                updated_at = current_timestamp_utc(),
                updated_by = COALESCE(NEW.created_by, NEW.created_by) -- Usar created_by ya que employee_roles solo tiene created_by
            WHERE id = NEW.employee_id
            AND (primary_role_id IS DISTINCT FROM NEW.role_id); -- Solo actualizar si es diferente
            
            -- Desmarcar otros roles como primarios para este empleado
            UPDATE employee_roles 
            SET is_primary = false
            WHERE employee_id = NEW.employee_id 
            AND role_id != NEW.role_id 
            AND is_primary = true
            AND deleted_at IS NULL;
        END IF;
        RETURN NEW;
    END IF;
    
    IF TG_OP = 'DELETE' THEN
        -- Si se elimina el rol primario, limpiar employees.primary_role_id
        IF OLD.is_primary = true THEN
            UPDATE employees 
            SET primary_role_id = NULL,
                updated_at = current_timestamp_utc()
            WHERE id = OLD.employee_id
            AND primary_role_id = OLD.role_id; -- Solo limpiar si coincide
        END IF;
        RETURN OLD;
    END IF;
    
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_employee_roles_sync_primary ON employee_roles;

CREATE TRIGGER trg_employee_roles_sync_primary
    AFTER INSERT OR UPDATE OF is_primary OR DELETE ON employee_roles
    FOR EACH ROW
    EXECUTE FUNCTION sync_employee_primary_role();

-- Función para validar que no se puede despedir el último admin
CREATE OR REPLACE FUNCTION validate_last_admin_employee()
RETURNS TRIGGER AS $$
DECLARE
    admin_role_id UUID;
    admin_count INTEGER;
BEGIN
    -- Solo validar si se está cambiando el status a 'terminated' o seteando fired_at
    IF (NEW.status = 'terminated' AND OLD.status != 'terminated') OR 
       (NEW.fired_at IS NOT NULL AND OLD.fired_at IS NULL) THEN
        
        -- Obtener el ID del rol Admin para esta organización (case-insensitive)
        SELECT id INTO admin_role_id 
        FROM organization_role 
        WHERE organization_id = NEW.organization_id 
        AND LOWER(name) = 'admin' 
        AND deleted_at IS NULL;
        
        -- Verificar si este empleado tiene rol Admin
        IF EXISTS (SELECT 1 FROM employee_roles 
                  WHERE employee_id = NEW.id 
                  AND role_id = admin_role_id
                  AND deleted_at IS NULL) THEN
            
            -- Contar otros admins activos en la organización
            SELECT COUNT(*) INTO admin_count
            FROM employees e
            JOIN employee_roles er ON e.id = er.employee_id
            WHERE e.organization_id = NEW.organization_id
            AND e.id != NEW.id  -- Excluir el empleado actual
            AND e.status = 'active'
            AND e.deleted_at IS NULL
            AND er.role_id = admin_role_id
            AND er.deleted_at IS NULL;
            
            -- Si no hay otros admins, no permitir despedir
            IF admin_count = 0 THEN
                RAISE EXCEPTION 'Cannot terminate the last Admin of the organization. Assign Admin role to another employee first.';
            END IF;
        END IF;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_employees_validate_last_admin
    BEFORE UPDATE ON employees
    FOR EACH ROW
    EXECUTE FUNCTION validate_last_admin_employee();

-- Función para prevenir hard-delete del último admin
CREATE OR REPLACE FUNCTION prevent_delete_last_admin()
RETURNS TRIGGER AS $$
DECLARE
    admin_role_id UUID;
    admin_count INTEGER;
BEGIN
    -- Obtener el ID del rol Admin para esta organización (case-insensitive)
    SELECT id INTO admin_role_id 
    FROM organization_role 
    WHERE organization_id = OLD.organization_id 
    AND LOWER(name) = 'admin' 
    AND deleted_at IS NULL;
    
    -- Verificar si este empleado tiene rol Admin
    IF EXISTS (SELECT 1 FROM employee_roles 
              WHERE employee_id = OLD.id 
              AND role_id = admin_role_id
              AND deleted_at IS NULL) THEN
        
        -- Contar otros admins activos en la organización
        SELECT COUNT(*) INTO admin_count
        FROM employees e
        JOIN employee_roles er ON e.id = er.employee_id
        WHERE e.organization_id = OLD.organization_id
        AND e.id != OLD.id  -- Excluir el empleado actual
        AND e.status = 'active'
        AND e.deleted_at IS NULL
        AND er.role_id = admin_role_id
        AND er.deleted_at IS NULL;
        
        -- Si no hay otros admins, no permitir eliminación
        IF admin_count = 0 THEN
            RAISE EXCEPTION 'Cannot hard-delete the last Admin of the organization. Assign Admin role to another employee first.';
        END IF;
    END IF;
    
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- Trigger para prevenir hard-delete del último admin
CREATE TRIGGER trg_employees_prevent_delete_last_admin
    BEFORE DELETE ON employees
    FOR EACH ROW
    EXECUTE FUNCTION prevent_delete_last_admin();

-- =============================================================================
-- COMENTARIOS PARA DOCUMENTACIÓN
-- =============================================================================
COMMENT ON TABLE employees IS 'Empleados - vínculo entre usuarios y organizaciones';
COMMENT ON COLUMN employees.user_id IS 'FK lógica externa → auth-identity-svc.users.id (no enforce DB)';
COMMENT ON COLUMN employees.person_id IS 'FK lógica externa → person-svc.person.id (no enforce DB) - Agregado para consistency futura';
COMMENT ON COLUMN employees.primary_role_id IS 'Rol principal del empleado (ON DELETE SET NULL)';
COMMENT ON COLUMN employees.status IS 'Estado del empleado (active, inactive, suspended, terminated)';
COMMENT ON COLUMN employees.fired_at IS 'Fecha de despido - NULL significa empleado no despedido';
COMMENT ON COLUMN employees.hired_at IS 'Fecha de contratación - no puede ser futura';
COMMENT ON COLUMN employees.created_by IS 'FK lógica externa → auth-identity-svc.users.id (no enforce DB)';
COMMENT ON COLUMN employees.updated_by IS 'FK lógica externa → auth-identity-svc.users.id (no enforce DB)';

COMMENT ON TABLE employee_roles IS 'Roles asignados a empleados (N-a-N)';
COMMENT ON COLUMN employee_roles.is_primary IS 'Solo uno puede ser true por empleado';
COMMENT ON COLUMN employee_roles.created_by IS 'FK lógica externa → auth-identity-svc.users.id (no enforce DB)';

COMMENT ON FUNCTION sync_employee_primary_role() IS 'Sincroniza el campo primary_role_id en employees con la tabla employee_roles';
COMMENT ON FUNCTION validate_last_admin_employee() IS 'Previene despedir al último administrador de una organización (soft-delete)';
COMMENT ON FUNCTION prevent_delete_last_admin() IS 'Previene hard-delete del último administrador de una organización';

-- =============================================================================
-- SMOKE TEST PARA VALIDACIÓN EN CI
-- =============================================================================
/*
Simple smoke test para verificar que la migración 0004 se aplicó correctamente:

-- Verificar que las tablas existen
SELECT 1 FROM employees LIMIT 0;
SELECT 1 FROM employee_roles LIMIT 0;

-- Verificar que las FK están correctamente configuradas
SELECT conname, contype FROM pg_constraint 
WHERE conrelid IN (
    'employees'::regclass, 
    'employee_roles'::regclass
);

-- Verificar índices únicos 
SELECT indexname FROM pg_indexes 
WHERE tablename IN ('employees', 'employee_roles') 
AND indexname LIKE 'uq_%';

-- Test de constraint de fecha no futura (debe fallar)
-- INSERT INTO employees (organization_id, user_id, hired_at, created_by) 
-- VALUES (generate_uuid(), generate_uuid(), CURRENT_DATE + 1, generate_uuid());
*/
