-- 0001_initial.up.sql - Property Service Initial Schema
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Enum for amenity categories
CREATE TYPE amenity_category AS ENUM ('general', 'services', 'environments', 'security', 'comfort');

-- Property type catalog table
CREATE TABLE property_type (
  id                SERIAL PRIMARY KEY,
  code              VARCHAR(32) UNIQUE NOT NULL,
  name              VARCHAR(100) NOT NULL,
  description       TEXT,
  is_active         BOOLEAN DEFAULT TRUE,
  created_at        TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at        TIMESTAMP
);

-- Manager type catalog table (for property management)
CREATE TABLE manager_type (
  id                SERIAL PRIMARY KEY,
  code              VARCHAR(32) UNIQUE NOT NULL,
  name              VARCHAR(100) NOT NULL,
  description       TEXT,
  is_active         BOOLEAN DEFAULT TRUE,
  created_at        TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at        TIMESTAMP
);

-- Amenity catalog table
CREATE TABLE amenity (
  id                SERIAL PRIMARY KEY,
  name              VARCHAR(100) UNIQUE NOT NULL,
  category          amenity_category,
  icon_url          TEXT
);

-- Main property table
CREATE TABLE property (
  id                    UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  owner_person_id       UUID NOT NULL,
  address_id            UUID UNIQUE NOT NULL,
  property_type_id      INTEGER NOT NULL REFERENCES property_type(id),
  internal_code         VARCHAR(50),
  year_built            INTEGER,
  bedrooms              INTEGER,
  bathrooms             DECIMAL(3,1),
  total_area_sqm        DECIMAL(10,2),
  covered_area_sqm      DECIMAL(10,2),
  description           TEXT,
  created_at            TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at            TIMESTAMP,
  updated_by            UUID,
  deleted_at            TIMESTAMP
);

-- Property management table
CREATE TABLE property_management (
  id                    UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  property_id           UUID NOT NULL REFERENCES property(id),
  manager_id            UUID NOT NULL,
  manager_type_id       INTEGER NOT NULL REFERENCES manager_type(id),
  start_date            DATE NOT NULL,
  end_date              DATE,
  commission_percentage DECIMAL(5,2),
  created_at            TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at            TIMESTAMP,
  updated_by            UUID,
  deleted_at            TIMESTAMP
);

-- Property amenities junction table
CREATE TABLE property_amenities (
  property_id           UUID NOT NULL REFERENCES property(id),
  amenity_id            JSONB,
  note                  TEXT,
  PRIMARY KEY (property_id, amenity_id)
);

-- Indexes for better performance
CREATE INDEX idx_property_owner ON property(owner_person_id);
CREATE INDEX idx_property_type ON property(property_type_id);
CREATE UNIQUE INDEX idx_property_internal_code ON property(internal_code) WHERE internal_code IS NOT NULL;

CREATE INDEX idx_property_management_property ON property_management(property_id);
CREATE INDEX idx_property_management_manager ON property_management(manager_id);
CREATE INDEX idx_property_management_type ON property_management(manager_type_id);
CREATE INDEX idx_property_management_dates ON property_management(property_id, start_date, end_date);

-- Insert initial property types
INSERT INTO property_type (code, name, description) VALUES
  ('APARTMENT', 'Departamento', 'Departamento en edificio'),
  ('HOUSE', 'Casa', 'Casa unifamiliar'),
  ('PH', 'PH', 'Propiedad horizontal'),
  ('VACATION_HOME', 'Quinta Vacacional', 'Propiedad para vacaciones'),
  ('COUNTRYSIDE', 'Campo', 'Propiedad rural'),
  ('TOMB', 'Bóveda, Nicho, Parcela', 'Propiedad funeraria'),
  ('LAND', 'Terreno', 'Terreno sin construcciones'),
  ('COMMERCIAL_OFFICE', 'Oficina Comercial', 'Espacio de oficina'),
  ('COMMERCIAL_STORE', 'Local Comercial', 'Local comercial'),
  ('BUILDING', 'Edificio', 'Edificio completo'),
  ('WAREHOUSE_STORAGE', 'Bodega / Galpón', 'Espacio de almacenamiento'),
  ('CLINIC', 'Consultorio', 'Espacio médico'),
  ('GARAGE', 'Cochera', 'Plaza de garaje'),
  ('STORAGE', 'Depósito', 'Depósito'),
  ('BUSINESS_BACKGROUND', 'Fondo de Comercio', 'Fondo de comercio'),
  ('HOTEL', 'Hotel', 'Establecimiento hotelero'),
  ('BOAT_SLIP', 'Cama Náutica', 'Plaza de amarre'),
  ('LAND_COMMERCIAL', 'Terreno Comercial', 'Terreno con uso comercial'),
  ('LAND_INDUSTRIAL', 'Terreno Industrial', 'Terreno con uso industrial'),
  ('WAREHOUSE_INDUSTRIAL', 'Galpón Industrial', 'Galpón para uso industrial'),
  ('INDUSTRIAL_BUILDING', 'Nave Industrial', 'Edificación industrial');

-- Insert initial manager types
INSERT INTO manager_type (code, name, description) VALUES
  ('PERSON', 'Persona', 'Gestor individual'),
  ('ORGANIZATION', 'Organización', 'Gestor organizacional');

-- Insert initial amenities
INSERT INTO amenity (name, category) VALUES
  -- General
  ('Apto Crédito', 'general'),
  ('Apto Profesional', 'general'),
  ('Luminoso', 'general'),
  ('Contra Frente', 'general'),
  ('Frente', 'general'),
  ('Balcón', 'general'),
  ('Terraza', 'general'),
  ('Patio', 'general'),
  ('Jardín', 'general'),
  ('Quincho', 'general'),
  ('Parrilla', 'general'),
  ('Pileta', 'general'),
  ('Solarium', 'general'),
  ('Cochera', 'general'),
  ('Baulera', 'general'),
  ('Lavadero', 'general'),
  ('Dependencia de Servicio', 'general'),
  ('Vista al Mar', 'general'),
  ('Vista Panorámica', 'general'),
  ('Esquina', 'general'),
  ('Entrada Independiente', 'general'),
  ('Doble Altura', 'general'),
  
  -- Services
  ('Agua Corriente', 'services'),
  ('Gas Natural', 'services'),
  ('Gas Envasado', 'services'),
  ('Electricidad', 'services'),
  ('Cloacas', 'services'),
  ('Internet', 'services'),
  ('Cable', 'services'),
  ('Teléfono', 'services'),
  ('Aire Acondicionado', 'services'),
  ('Calefacción Central', 'services'),
  ('Calefacción Individual', 'services'),
  ('Agua Caliente Central', 'services'),
  ('Agua Caliente Individual', 'services'),
  ('Ascensor', 'services'),
  ('Montacargas', 'services'),
  ('Generador Eléctrico', 'services'),
  ('Pozo de Agua', 'services'),
  ('Tanque de Agua', 'services'),
  ('Bomba de Agua', 'services'),
  
  -- Environments
  ('Living', 'environments'),
  ('Comedor', 'environments'),
  ('Living Comedor', 'environments'),
  ('Cocina', 'environments'),
  ('Cocina Integrada', 'environments'),
  ('Office', 'environments'),
  ('Toilette', 'environments'),
  ('Hall de Distribución', 'environments'),
  ('Vestidor', 'environments'),
  ('Suite', 'environments'),
  ('Escritorio', 'environments'),
  ('Playroom', 'environments'),
  ('Sala de Estar', 'environments'),
  ('Biblioteca', 'environments'),
  ('Estar íntimo', 'environments'),
  ('Family Room', 'environments'),
  ('Sala de TV', 'environments'),
  ('Sala de Juegos', 'environments'),
  ('Habitación de Servicio', 'environments'),
  ('Baño de Servicio', 'environments'),
  
  -- Security
  ('Portero', 'security'),
  ('Portero Eléctrico', 'security'),
  ('Seguridad 24hs', 'security'),
  ('Cámaras de Seguridad', 'security'),
  ('Alarma', 'security'),
  ('Rejas', 'security'),
  ('Barrio Cerrado', 'security'),
  ('Country', 'security'),
  ('Vigilancia', 'security'),
  ('Control de Acceso', 'security'),
  ('Cerco Eléctrico', 'security'),
  ('Puerta Blindada', 'security'),
  ('Circuito Cerrado TV', 'security'),
  
  -- Comfort
  ('Amoblado', 'comfort'),
  ('Semi Amoblado', 'comfort'),
  ('A Estrenar', 'comfort'),
  ('Reciclado', 'comfort'),
  ('En Construcción', 'comfort'),
  ('Muy Bueno', 'comfort'),
  ('Bueno', 'comfort'),
  ('A Refaccionar', 'comfort'),
  ('SUM', 'comfort'),
  ('Gimnasio', 'comfort'),
  ('Spa', 'comfort'),
  ('Sauna', 'comfort'),
  ('Cine', 'comfort'),
  ('Cancha de Tenis', 'comfort'),
  ('Cancha de Fútbol', 'comfort'),
  ('Asador', 'comfort'),
  ('Deck', 'comfort'),
  ('Jacuzzi', 'comfort'),
  ('Hidromasaje', 'comfort');
