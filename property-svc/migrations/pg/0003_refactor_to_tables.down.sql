-- Rollback migration 0003: Convert tables back to enums

-- Add back the old property_type enum
CREATE TYPE property_type AS ENUM (
  'APARTMENT',
  'HOUSE', 
  'COMMERCIAL_SPACE',
  'OFFICE',
  'LAND',
  'INDUSTRIAL_WAREHOUSE'
);

-- Add back property_type column to property table
ALTER TABLE property ADD COLUMN property_type property_type;

-- Migrate data back from property_type_id to property_type enum
UPDATE property SET property_type = pt.code::property_type
FROM property_type pt
WHERE property.property_type_id = pt.id;

-- Make property_type NOT NULL
ALTER TABLE property ALTER COLUMN property_type SET NOT NULL;

-- Drop the new property_type_id column and constraints
ALTER TABLE property DROP CONSTRAINT IF EXISTS fk_property_property_type;
DROP INDEX IF EXISTS idx_property_property_type_id;
ALTER TABLE property DROP COLUMN property_type_id;

-- Add back the old index
CREATE INDEX idx_property_type ON property(property_type);

-- Revert property_management table structure
-- Rename columns back
ALTER TABLE property_management RENAME COLUMN start_date TO managed_since;
ALTER TABLE property_management RENAME COLUMN end_date TO managed_until;
ALTER TABLE property_management RENAME COLUMN commission_percentage TO commission_percent;

-- Add back is_active column
ALTER TABLE property_management ADD COLUMN is_active boolean NOT NULL DEFAULT true;

-- Update is_active based on managed_until
UPDATE property_management SET is_active = (managed_until IS NULL);

-- Rename manager_id back to organization_id
ALTER TABLE property_management RENAME COLUMN manager_id TO organization_id;

-- Drop manager_type constraints and column
ALTER TABLE property_management DROP CONSTRAINT IF EXISTS fk_property_management_manager_type;
DROP INDEX IF EXISTS idx_property_management_manager_type_id;
ALTER TABLE property_management DROP COLUMN manager_type_id;

-- Recreate old indexes
DROP INDEX IF EXISTS idx_property_management_manager_id;
DROP INDEX IF EXISTS idx_property_management_dates;

CREATE INDEX idx_property_management_organization_id ON property_management(organization_id);
CREATE INDEX idx_property_management_active ON property_management(is_active);

-- Drop the new tables
DROP TABLE IF EXISTS manager_type;
DROP TABLE IF EXISTS property_type;
