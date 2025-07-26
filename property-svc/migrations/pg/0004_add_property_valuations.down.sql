-- Migration: Drop property_valuations table
-- Description: Remove property_valuations table and related indexes/triggers/functions

-- Drop trigger
DROP TRIGGER IF EXISTS trigger_property_valuations_updated_at ON property_valuations;

-- Drop function
DROP FUNCTION IF EXISTS update_property_valuations_updated_at();

-- Drop indexes
DROP INDEX IF EXISTS idx_property_valuations_valid_current;
DROP INDEX IF EXISTS idx_property_valuations_latest;
DROP INDEX IF EXISTS idx_property_valuations_market_history;
DROP INDEX IF EXISTS idx_property_valuations_property_date_type;
DROP INDEX IF EXISTS idx_property_valuations_deleted_at;
DROP INDEX IF EXISTS idx_property_valuations_created_at;
DROP INDEX IF EXISTS idx_property_valuations_valid_until;
DROP INDEX IF EXISTS idx_property_valuations_method;
DROP INDEX IF EXISTS idx_property_valuations_appraiser;
DROP INDEX IF EXISTS idx_property_valuations_currency;
DROP INDEX IF EXISTS idx_property_valuations_appraised_value;
DROP INDEX IF EXISTS idx_property_valuations_type;
DROP INDEX IF EXISTS idx_property_valuations_valuation_date;
DROP INDEX IF EXISTS idx_property_valuations_property_id;

-- Drop table
DROP TABLE IF EXISTS property_valuations;
