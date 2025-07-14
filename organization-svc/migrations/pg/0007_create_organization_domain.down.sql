-- =============================================
-- Migration: 0007_create_organization_domain.down.sql
-- Description: Drop organization domain tables and related functions
-- Author: System Generated
-- Date: 2024
-- Version: 2.0 - Enhanced rollback
-- =============================================

-- Eliminar particiones de logs (las particiones se eliminan automáticamente con la tabla padre)
-- Pero listamos las conocidas por si hay que eliminarlas manualmente
DROP TABLE IF EXISTS organization_domain_verification_log_y2024m01;
DROP TABLE IF EXISTS organization_domain_verification_log_y2024m02;
DROP TABLE IF EXISTS organization_domain_verification_log_y2024m03;

-- Eliminar tablas en orden inverso para evitar problemas de FK
DROP TABLE IF EXISTS organization_domain_verification_log;
DROP TABLE IF EXISTS organization_domain_dns;
DROP TABLE IF EXISTS organization_domain;

-- Eliminar funciones
DROP FUNCTION IF EXISTS manage_domain_primary_redirect;
DROP FUNCTION IF EXISTS log_domain_verification;
DROP FUNCTION IF EXISTS create_domain_verification_log_partition;
DROP FUNCTION IF EXISTS sync_domain_dns_verification;
DROP FUNCTION IF EXISTS create_default_dns_records;
DROP FUNCTION IF EXISTS generate_dns_verification_token;

-- Nota: Los índices únicos se eliminan automáticamente con las tablas
