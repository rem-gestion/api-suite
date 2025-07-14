-- =============================================
-- Migration: 0005_create_organization_invite.down.sql
-- Description: Rollback organization invitation system completely
-- Author: System Generated
-- Date: 2024
-- Version: 1.1
-- =============================================

-- Drop tables (will cascade to triggers and indexes automatically)
DROP TABLE IF EXISTS organization_invite_log;
DROP TABLE IF EXISTS organization_invite;

-- Drop utility functions
DROP FUNCTION IF EXISTS cleanup_old_invitations;
DROP FUNCTION IF EXISTS get_organization_invite_stats;
DROP FUNCTION IF EXISTS organization_invite_log_partition_trigger;
DROP FUNCTION IF EXISTS renew_invite_token;
DROP FUNCTION IF EXISTS generate_invite_token;
DROP FUNCTION IF EXISTS expire_old_invitations;
DROP FUNCTION IF EXISTS check_organization_invite_limits;
DROP FUNCTION IF EXISTS log_organization_invite_status_change;

-- Note: pgcrypto extension and ENUMs are managed by other migrations

-- Rollback completed successfully
