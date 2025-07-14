-- =============================================
-- Migration: 0006_create_organization_integration.down.sql
-- Description: Rollback organization integration system completely
-- Author: System Generated
-- Date: 2024
-- Version: 1.1
-- =============================================

-- Drop tables in dependency order (will cascade to triggers, indexes, constraints)
DROP TABLE IF EXISTS organization_integration_log;
DROP TABLE IF EXISTS organization_integration_event;
DROP TABLE IF EXISTS organization_integration;
DROP TABLE IF EXISTS integration_type;

-- Drop utility functions
DROP FUNCTION IF EXISTS get_pending_integration_events;
DROP FUNCTION IF EXISTS get_organization_integration_stats;
DROP FUNCTION IF EXISTS cleanup_integration_logs;
DROP FUNCTION IF EXISTS organization_integration_log_partition_trigger;
DROP FUNCTION IF EXISTS log_integration_status_change;
DROP FUNCTION IF EXISTS update_integration_updated_at;

-- Note: pgcrypto extension and ENUMs are managed by other migrations
-- Note: Partitions are automatically dropped with the main table

-- Rollback completed successfully
