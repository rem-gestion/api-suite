-- =============================================
-- Migration: 0003_create_structure.down.sql
-- Description: Drop organization structure tables and related functions
-- Author: System Generated
-- Date: 2024
-- =============================================

-- Eliminar tablas en orden inverso
DROP TABLE IF EXISTS organization_role;
DROP TABLE IF EXISTS organization_branch;

-- Eliminar funciones (sintaxis PostgreSQL sin paréntesis)
DROP FUNCTION IF EXISTS create_default_main_branch;
DROP FUNCTION IF EXISTS create_default_organization_roles;
DROP FUNCTION IF EXISTS prevent_delete_default_roles;
DROP FUNCTION IF EXISTS prevent_delete_only_main_branch;
