-- =============================================
-- Migration: 0004_create_employee.down.sql
-- Description: Drop employee and employee roles tables and related functions
-- Author: System Generated
-- Date: 2024
-- =============================================

-- Eliminar tablas en orden inverso para evitar problemas de FK
DROP TABLE IF EXISTS employee_roles;
DROP TABLE IF EXISTS employees;

-- Eliminar funciones (sintaxis PostgreSQL sin paréntesis)
DROP FUNCTION IF EXISTS prevent_delete_last_admin;
DROP FUNCTION IF EXISTS validate_last_admin_employee;
DROP FUNCTION IF EXISTS sync_employee_primary_role;
