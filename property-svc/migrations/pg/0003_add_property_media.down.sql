-- Migration: Drop property_media table
-- Description: Remove property_media table and related indexes/triggers/functions

-- Drop triggers
DROP TRIGGER IF EXISTS trigger_ensure_single_main_media ON property_media;
DROP TRIGGER IF EXISTS trigger_property_media_updated_at ON property_media;

-- Drop functions
DROP FUNCTION IF EXISTS ensure_single_main_media();
DROP FUNCTION IF EXISTS update_property_media_updated_at();

-- Drop indexes
DROP INDEX IF EXISTS idx_property_media_public_display;
DROP INDEX IF EXISTS idx_property_media_property_type_order;
DROP INDEX IF EXISTS idx_property_media_deleted_at;
DROP INDEX IF EXISTS idx_property_media_created_at;
DROP INDEX IF EXISTS idx_property_media_is_featured;
DROP INDEX IF EXISTS idx_property_media_is_public;
DROP INDEX IF EXISTS idx_property_media_is_main;
DROP INDEX IF EXISTS idx_property_media_display_order;
DROP INDEX IF EXISTS idx_property_media_type;
DROP INDEX IF EXISTS idx_property_media_property_id;

-- Drop table
DROP TABLE IF EXISTS property_media;
