CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- enum property_type
DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'property_type') THEN
    CREATE TYPE property_type AS ENUM (
      'APARTMENT',
      'HOUSE', 
      'COMMERCIAL_SPACE',
      'OFFICE',
      'LAND',
      'INDUSTRIAL_WAREHOUSE'
    );
  END IF;
END$$;

-- tabla property ────────────────
CREATE TABLE property (
  id                  uuid          PRIMARY KEY DEFAULT uuid_generate_v4(),
  owner_person_id     uuid          NOT NULL,  -- FK lógica (person-svc) → *sin* constraint hard
  address_id          uuid          UNIQUE NOT NULL,  -- FK lógica (address-svc) → *sin* constraint hard
  property_type       property_type NOT NULL,
  internal_code       varchar(50)   UNIQUE,
  year_built          int,
  bedrooms            int,
  bathrooms           decimal(3,1),
  total_area_sqm      decimal(10,2),
  covered_area_sqm    decimal(10,2),
  description         text,

  created_at          timestamp     NOT NULL DEFAULT now(),
  updated_at          timestamp,
  updated_by          uuid,
  deleted_at          timestamp
);

-- índices para performance
CREATE INDEX idx_property_owner_person_id ON property(owner_person_id);
CREATE INDEX idx_property_type ON property(property_type);
CREATE INDEX idx_property_internal_code ON property(internal_code) WHERE internal_code IS NOT NULL;

-- tabla property_management ────────────────
CREATE TABLE property_management (
  id                  uuid          PRIMARY KEY DEFAULT uuid_generate_v4(),
  property_id         uuid          NOT NULL REFERENCES property(id) ON DELETE CASCADE,
  organization_id     uuid          NOT NULL,  -- FK lógica (organization-svc) → *sin* constraint hard
  managed_since       timestamp     NOT NULL DEFAULT now(),
  managed_until       timestamp,
  is_active           boolean       NOT NULL DEFAULT true,
  commission_percent  decimal(5,2),  -- porcentaje de comisión
  notes               text,

  created_at          timestamp     NOT NULL DEFAULT now(),
  updated_at          timestamp,
  updated_by          uuid,
  deleted_at          timestamp
);

-- índices para property_management
CREATE INDEX idx_property_management_property_id ON property_management(property_id);
CREATE INDEX idx_property_management_organization_id ON property_management(organization_id);
CREATE INDEX idx_property_management_active ON property_management(is_active);

-- tabla amenity ────────────────
CREATE TABLE amenity (
  id          uuid         PRIMARY KEY DEFAULT uuid_generate_v4(),
  name        varchar(100) NOT NULL UNIQUE,
  description text,
  icon        varchar(50), -- para la UI
  category    varchar(50), -- ej: "security", "recreation", "services"

  created_at  timestamp    NOT NULL DEFAULT now(),
  updated_at  timestamp,
  deleted_at  timestamp
);

-- tabla property_amenities ────────────────
CREATE TABLE property_amenities (
  id          uuid      PRIMARY KEY DEFAULT uuid_generate_v4(),
  property_id uuid      NOT NULL REFERENCES property(id) ON DELETE CASCADE,
  amenity_id  uuid      NOT NULL REFERENCES amenity(id) ON DELETE CASCADE,
  
  created_at  timestamp NOT NULL DEFAULT now(),
  
  UNIQUE(property_id, amenity_id)
);

-- índices para property_amenities
CREATE INDEX idx_property_amenities_property_id ON property_amenities(property_id);
CREATE INDEX idx_property_amenities_amenity_id ON property_amenities(amenity_id);
