-- =============================================
-- Migration: 0001_create_enums_and_types.down.sql
-- Description: Drop PostgreSQL ENUMs and custom types
-- Author: System Generated
-- Date: 2024
-- 
-- IMPORTANTE: Este script debe ejecutarse DESPUÉS de eliminar todas las tablas
-- que dependan de estos ENUMs. Si se ejecuta 'migrate down 1' sin haber 
-- eliminado las tablas primero, los DROP TYPE fallarán por dependencias.
-- Secuencia correcta: eliminar tablas → eliminar funciones → eliminar tipos
-- =============================================

-- Eliminar funciones utilitarias
DROP FUNCTION IF EXISTS organization_integration_log_insert_trigger;
DROP FUNCTION IF EXISTS cleanup_old_partitions;
DROP FUNCTION IF EXISTS create_monthly_partition;
DROP FUNCTION IF EXISTS is_valid_domain;
DROP FUNCTION IF EXISTS is_valid_uuid;
DROP FUNCTION IF EXISTS is_valid_email;
DROP FUNCTION IF EXISTS current_timestamp_utc;
DROP FUNCTION IF EXISTS generate_uuid;

-- Eliminar ENUMs de suscripción
DROP TYPE IF EXISTS aggregation_type_enum;
DROP TYPE IF EXISTS usage_metric_enum;
DROP TYPE IF EXISTS invoice_status_enum;
DROP TYPE IF EXISTS subscription_status_enum;
DROP TYPE IF EXISTS currency_enum;
DROP TYPE IF EXISTS billing_interval_enum;
DROP TYPE IF EXISTS plan_type_enum;

-- Eliminar ENUMs de dominio
DROP TYPE IF EXISTS verification_status_enum;
DROP TYPE IF EXISTS domain_verification_type_enum;
DROP TYPE IF EXISTS dns_record_type_enum;
DROP TYPE IF EXISTS dns_verification_method_enum;
DROP TYPE IF EXISTS domain_status_enum;
DROP TYPE IF EXISTS domain_type_enum;

-- Eliminar ENUMs de integración
DROP TYPE IF EXISTS log_level_enum;
DROP TYPE IF EXISTS integration_event_status_enum;
DROP TYPE IF EXISTS integration_event_type_enum;
DROP TYPE IF EXISTS sync_frequency_enum;
DROP TYPE IF EXISTS sync_status_enum;
DROP TYPE IF EXISTS integration_status_enum;
DROP TYPE IF EXISTS integration_category_enum;

-- Eliminar ENUMs de invitación
DROP TYPE IF EXISTS invitation_log_action_enum;
DROP TYPE IF EXISTS invitation_status_enum;

-- Eliminar ENUMs de empleado
DROP TYPE IF EXISTS employee_status_enum;

-- Eliminar ENUMs de organización
DROP TYPE IF EXISTS owner_type_enum;
DROP TYPE IF EXISTS organization_status_enum;

-- Nota: No eliminamos la extensión uuid-ossp ya que puede ser utilizada por otros servicios

-- Comentario de finalización
-- ENUMs and utility types removed successfully
