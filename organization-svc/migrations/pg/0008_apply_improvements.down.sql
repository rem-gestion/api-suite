-- =============================================
-- Migration: 0008_apply_improvements.down.sql
-- Description: Revert improvements applied in 0008
-- Author: System Generated
-- Date: 2024
-- Version: 2.0 - Enhanced rollback
-- =============================================

-- Eliminar vista
DROP VIEW IF EXISTS organization_subscription_details;

-- Eliminar funciones de mantenimiento
DROP FUNCTION IF EXISTS run_maintenance_tasks;
DROP FUNCTION IF EXISTS run_organization_maintenance;
DROP FUNCTION IF EXISTS validate_admin_exists_check;
DROP FUNCTION IF EXISTS trg_monthly_partition;
DROP FUNCTION IF EXISTS cleanup_old_invite_logs;
DROP FUNCTION IF EXISTS cleanup_old_integration_logs;
DROP FUNCTION IF EXISTS cleanup_old_domain_verification_logs;
DROP FUNCTION IF EXISTS ensure_default_integration_types;

-- Revertir constraint único corregido
ALTER TABLE employees DROP CONSTRAINT IF EXISTS uq_employees_organization_user;

-- Eliminar índices adicionales creados en 0008
DROP INDEX IF EXISTS ix_employees_org_status_active;
DROP INDEX IF EXISTS ix_organization_invite_org_status_expires;
DROP INDEX IF EXISTS ix_org_integration_org_status_active;
DROP INDEX IF EXISTS ix_organization_integration_org_active;
DROP INDEX IF EXISTS ix_organization_domain_org_primary;
DROP INDEX IF EXISTS ix_organization_subscription_organization_id;
DROP INDEX IF EXISTS ix_organization_subscription_status_active;

-- Nota: Triggers y FKs se deben revertir manualmente si es necesario
-- ya que el orden de eliminación puede afectar la integridad referencial

-- Nota: Las particiones creadas no se eliminan automáticamente
-- ya que pueden contener datos importantes. Se requiere 
-- eliminación manual si es necesario:
-- 
-- Para eliminar manualmente las particiones (CUIDADO - se pierden datos):
-- DROP TABLE IF EXISTS organization_invite_log_y*;
-- DROP TABLE IF EXISTS organization_integration_log_y*;
-- DROP TABLE IF EXISTS organization_domain_verification_log_y*;
