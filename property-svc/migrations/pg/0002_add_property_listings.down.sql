-- Migration: Drop property_listings table
-- Description: Remove property_listings table and related indexes/triggers

-- Drop trigger
DROP TRIGGER IF EXISTS trigger_property_listings_updated_at ON property_listings;

-- Drop function
DROP FUNCTION IF EXISTS update_property_listings_updated_at();

-- Drop indexes
DROP INDEX IF EXISTS idx_property_listings_price_range;
DROP INDEX IF EXISTS idx_property_listings_active_listings;
DROP INDEX IF EXISTS idx_property_listings_deleted_at;
DROP INDEX IF EXISTS idx_property_listings_created_at;
DROP INDEX IF EXISTS idx_property_listings_publish_web;
DROP INDEX IF EXISTS idx_property_listings_highlighted;
DROP INDEX IF EXISTS idx_property_listings_currency;
DROP INDEX IF EXISTS idx_property_listings_price;
DROP INDEX IF EXISTS idx_property_listings_status;
DROP INDEX IF EXISTS idx_property_listings_operation_type;
DROP INDEX IF EXISTS idx_property_listings_property_id;

-- Drop table
DROP TABLE IF EXISTS property_listings;
