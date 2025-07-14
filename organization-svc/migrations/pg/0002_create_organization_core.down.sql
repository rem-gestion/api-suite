-- =============================================
-- Migration: 0002_create_organization_core.down.sql
-- Description: Drop core organization tables and related functions
-- Author: System Generated
-- Date: 2024
-- 
-- NOTA: Este script elimina triggers explícitamente antes de las tablas
-- para evitar errores con psql -v ON_ERROR_STOP=1 si fueron eliminados manualmente
-- =============================================

-- Eliminar tablas en orden inverso para evitar problemas de FK
-- Nota: Los triggers se eliminan automáticamente al borrar las tablas
-- pero los documentamos aquí para referencia:
-- * trg_organization_create_defaults
-- * trg_organization_prevent_hard_delete 
-- * trg_organization_updated_at
-- * trg_organization_settings_updated_at
-- * trg_organization_owner_updated_at
-- * trg_organization_owner_validate_percentages

DROP TABLE IF EXISTS organization_owner;
DROP TABLE IF EXISTS organization_settings;
DROP TABLE IF EXISTS organization;

-- Eliminar funciones (PostgreSQL syntax)
DROP FUNCTION IF EXISTS create_organization_defaults;
DROP FUNCTION IF EXISTS create_default_organization_settings;
DROP FUNCTION IF EXISTS validate_ownership_percentages;
DROP FUNCTION IF EXISTS prevent_hard_delete_organization;
DROP FUNCTION IF EXISTS update_organization_updated_at;
DROP FUNCTION IF EXISTS update_timestamp_utc;
