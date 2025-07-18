-- Migration to convert property_type and manager_type from enums to tables
-- This migration creates the new tables and migrates existing data

-- First, we need to handle the existing enum carefully
-- Create temporary columns for the migration

-- Add new columns to property table ────────────────
ALTER TABLE property ADD COLUMN property_type_id int;

-- Create property_type table ────────────────
CREATE TABLE property_type (
  id          int          PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  code        varchar(50)  UNIQUE NOT NULL,
  name        varchar(100) NOT NULL,
  description text,
  is_active   boolean      NOT NULL DEFAULT true,
  created_at  timestamp    NOT NULL DEFAULT now(),
  updated_at  timestamp
);

-- Insert property type data
INSERT INTO property_type (code, name, description) VALUES
('APARTMENT', 'Departamento', 'Unidad habitacional en edificio de departamentos'),
('HOUSE', 'Casa', 'Vivienda unifamiliar independiente'),
('COMMERCIAL_SPACE', 'Local Comercial', 'Espacio destinado a actividad comercial'),
('OFFICE', 'Oficina', 'Espacio destinado a actividades de oficina'),
('LAND', 'Terreno', 'Lote de tierra sin construcciones'),
('INDUSTRIAL_WAREHOUSE', 'Galpón Industrial', 'Nave industrial para uso manufacturero o depósito');

-- Migrate existing property_type enum data to property_type_id
UPDATE property SET property_type_id = pt.id
FROM property_type pt
WHERE property.property_type::text = pt.code;

-- Make property_type_id NOT NULL after migration
ALTER TABLE property ALTER COLUMN property_type_id SET NOT NULL;

-- Add foreign key constraint
ALTER TABLE property ADD CONSTRAINT fk_property_property_type 
FOREIGN KEY (property_type_id) REFERENCES property_type(id);

-- Drop old property_type column
ALTER TABLE property DROP COLUMN property_type;

-- Add index for property_type_id
CREATE INDEX idx_property_property_type_id ON property(property_type_id);

-- Drop old index
DROP INDEX IF EXISTS idx_property_type;

-- Create manager_type table ────────────────
CREATE TABLE manager_type (
  id          int          PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  code        varchar(50)  UNIQUE NOT NULL,
  name        varchar(100) NOT NULL,
  description text,
  is_active   boolean      NOT NULL DEFAULT true,
  created_at  timestamp    NOT NULL DEFAULT now(),
  updated_at  timestamp
);

-- Insert manager type data
INSERT INTO manager_type (code, name, description) VALUES
('PERSON', 'Persona', 'Persona física responsable de la gestión'),
('ORGANIZATION', 'Organización', 'Empresa u organización responsable de la gestión');

-- Update property_management table structure ────────────────
-- Rename organization_id to manager_id for consistency
ALTER TABLE property_management RENAME COLUMN organization_id TO manager_id;

-- Add manager_type_id column
ALTER TABLE property_management ADD COLUMN manager_type_id int NOT NULL DEFAULT 2; -- Default to ORGANIZATION for existing records

-- Add foreign key constraint for manager_type
ALTER TABLE property_management ADD CONSTRAINT fk_property_management_manager_type 
FOREIGN KEY (manager_type_id) REFERENCES manager_type(id);

-- Rename columns for better semantics
ALTER TABLE property_management RENAME COLUMN managed_since TO start_date;
ALTER TABLE property_management RENAME COLUMN managed_until TO end_date;
ALTER TABLE property_management RENAME COLUMN commission_percent TO commission_percentage;

-- Remove is_active column (will be determined by end_date)
ALTER TABLE property_management DROP COLUMN is_active;

-- Drop old index and create new ones
DROP INDEX IF EXISTS idx_property_management_organization_id;
DROP INDEX IF EXISTS idx_property_management_active;

CREATE INDEX idx_property_management_manager_id ON property_management(manager_id);
CREATE INDEX idx_property_management_manager_type_id ON property_management(manager_type_id);
CREATE INDEX idx_property_management_dates ON property_management(property_id, start_date, end_date);

-- Drop the old property_type enum (this will work now since the column is dropped)
DROP TYPE IF EXISTS property_type;
