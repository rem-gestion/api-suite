-- 0001_initial.up.sql - Property Service Initial Schema
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Enum for amenity categories
CREATE TYPE amenity_category AS ENUM ('general', 'services', 'environments', 'security', 'comfort');

CREATE TABLE property_type (
  id                SERIAL PRIMARY KEY,
  code              VARCHAR(32) UNIQUE NOT NULL,
  name              VARCHAR(100) NOT NULL,
  description       TEXT,
  category          VARCHAR(50),
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
  property_type_id      INTEGER NOT NULL REFERENCES type(id),
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

CREATE TABLE property_management (
  id                    UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  property_id           UUID NOT NULL REFERENCES property(id),
  manager_id            UUID NOT NULL,
  manager_type_id       INTEGER NOT NULL REFERENCES manager_type(id),
  organization_id       UUID,
  start_date            DATE NOT NULL,
  end_date              DATE,
  is_active             BOOLEAN DEFAULT TRUE,
  commission_percentage DECIMAL(5,2),
  fixed_fee             DECIMAL(10,2),
  exclusive_management  BOOLEAN DEFAULT FALSE,
  services_included     JSONB,
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

CREATE INDEX idx_management_property ON management(property_id);
CREATE INDEX idx_management_manager ON management(manager_id);
CREATE INDEX idx_management_type ON management(manager_type_id);
CREATE INDEX idx_management_dates ON management(property_id, start_date, end_date);

-- Insert initial property types
INSERT INTO type (code, name, description, category) VALUES
  ('APARTMENT', 'Departamento', 'Departamento en edificio', 'residential'),
  ('HOUSE', 'Casa', 'Casa unifamiliar', 'residential'),
  ('PH', 'PH', 'Propiedad horizontal', 'residential'),
  ('VACATION_HOME', 'Quinta Vacacional', 'Propiedad para vacaciones', 'residential'),
  ('COUNTRYSIDE', 'Campo', 'Propiedad rural', 'residential'),
  ('TOMB', 'Bóveda, Nicho, Parcela', 'Propiedad funeraria', 'residential'),
  ('LAND', 'Terreno', 'Terreno sin construcciones', 'residential'),
  ('COMMERCIAL_OFFICE', 'Oficina Comercial', 'Espacio de oficina', 'commercial'),
  ('COMMERCIAL_STORE', 'Local Comercial', 'Local comercial', 'commercial'),
  ('BUILDING', 'Edificio', 'Edificio completo', 'commercial'),
  ('WAREHOUSE_STORAGE', 'Bodega / Galpón', 'Espacio de almacenamiento', 'industrial'),
  ('CLINIC', 'Consultorio', 'Espacio médico', 'commercial'),
  ('GARAGE', 'Cochera', 'Plaza de garaje', 'residential'),
  ('STORAGE', 'Depósito', 'Depósito', 'industrial'),
  ('BUSINESS_BACKGROUND', 'Fondo de Comercio', 'Fondo de comercio', 'commercial'),
  ('HOTEL', 'Hotel', 'Establecimiento hotelero', 'commercial'),
  ('BOAT_SLIP', 'Cama Náutica', 'Plaza de amarre', 'residential'),
  ('LAND_COMMERCIAL', 'Terreno Comercial', 'Terreno con uso comercial', 'commercial'),
  ('LAND_INDUSTRIAL', 'Terreno Industrial', 'Terreno con uso industrial', 'industrial'),
  ('WAREHOUSE_INDUSTRIAL', 'Galpón Industrial', 'Galpón para uso industrial', 'industrial'),
  ('INDUSTRIAL_BUILDING', 'Nave Industrial', 'Edificación industrial', 'industrial');

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
  ('Cloacas', 'services'),
  ('APARTMENT', 'Departamento', 'Departamento en edificio', 'residential'),
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

-- Insert example properties
INSERT INTO property (
  id, 
  owner_person_id, 
  address_id, 
  property_type_id, 
  internal_code,
  year_built,
  bedrooms,
  bathrooms,
  total_area_sqm,
  covered_area_sqm,
  description
) VALUES
  -- Departamento moderno en Palermo
  (
    '550e8400-e29b-41d4-a716-446655440001',
    '550e8400-e29b-41d4-a716-446655440101', -- owner_person_id
    '550e8400-e29b-41d4-a716-446655440201', -- address_id
    1, -- APARTMENT
    'APT-PAL-001',
    2019,
    2,
    1.0,
    75.50,
    65.00,
    'Moderno departamento de 2 ambientes en el corazón de Palermo. Luminoso, con balcón y excelente distribución. Edificio con amenities: gimnasio, terraza, parrilla común. A metros del subte y centros comerciales.'
  ),
  
  -- Casa familiar en San Isidro
  (
    '550e8400-e29b-41d4-a716-446655440002',
    '550e8400-e29b-41d4-a716-446655440102',
    '550e8400-e29b-41d4-a716-446655440202',
    2, -- HOUSE
    'CSA-SI-002',
    2005,
    4,
    3.0,
    280.00,
    180.00,
    'Hermosa casa familiar en zona residencial de San Isidro. 4 dormitorios, 3 baños, living-comedor integrado, cocina moderna, quincho con parrilla, jardín con pileta. Garage para 2 autos. Ideal para familias.'
  ),
  
  -- PH en Villa Crespo
  (
    '550e8400-e29b-41d4-a716-446655440003',
    '550e8400-e29b-41d4-a716-446655440103',
    '550e8400-e29b-41d4-a716-446655440203',
    3, -- PH
    'PH-VC-003',
    1995,
    3,
    2.0,
    120.00,
    90.00,
    'PH de 3 ambientes con terraza propia. Totalmente refaccionado, muy luminoso. Terraza de 30m2 con parrilla. Ubicado en zona tranquila con fácil acceso a transporte público.'
  ),
  
  -- Quinta vacacional en Tigre
  (
    '550e8400-e29b-41d4-a716-446655440004',
    '550e8400-e29b-41d4-a716-446655440104',
    '550e8400-e29b-41d4-a716-446655440204',
    4, -- VACATION_HOME
    'QTA-TIG-004',
    2010,
    3,
    2.0,
    450.00,
    150.00,
    'Quinta familiar en Tigre con amplio parque y dock privado. Casa principal de 3 dormitorios, quincho independiente, pileta, muelle propio. Ideal para descanso y actividades náuticas.'
  ),
  
  -- Local comercial en Microcentro
  (
    '550e8400-e29b-41d4-a716-446655440005',
    '550e8400-e29b-41d4-a716-446655440105',
    '550e8400-e29b-41d4-a716-446655440205',
    9, -- COMMERCIAL_STORE
    'LOC-MC-005',
    1980,
    0,
    1.0,
    85.00,
    85.00,
    'Local comercial en pleno Microcentro, sobre avenida principal. Amplia vidriera, excelente ubicación para cualquier tipo de comercio. Alto tránsito peatonal. Baño y depósito.'
  ),
  
  -- Oficina corporativa en Puerto Madero
  (
    '550e8400-e29b-41d4-a716-446655440006',
    '550e8400-e29b-41d4-a716-446655440106',
    '550e8400-e29b-41d4-a716-446655440206',
    8, -- COMMERCIAL_OFFICE
    'OFC-PM-006',
    2015,
    0,
    2.0,
    320.00,
    320.00,
    'Oficina premium en edificio corporativo de Puerto Madero. Piso completo con vista al río. Incluye 8 oficinas privadas, sala de reuniones, recepción, kitchenette. 4 cocheras incluidas.'
  ),
  
  -- Galpón industrial en Avellaneda
  (
    '550e8400-e29b-41d4-a716-446655440007',
    '550e8400-e29b-41d4-a716-446655440107',
    '550e8400-e29b-41d4-a716-446655440207',
    11, -- WAREHOUSE_STORAGE
    'GAL-AVE-007',
    2000,
    0,
    1.0,
    1500.00,
    1200.00,
    'Galpón industrial con nave de 1200m2 cubiertos más 300m2 de playón descubierto. Altura libre 8 metros, portón de acceso para camiones, oficina administrativa, vestuarios. Ideal logística.'
  ),
  
  -- Terreno urbano en Nordelta
  (
    '550e8400-e29b-41d4-a716-446655440008',
    '550e8400-e29b-41d4-a716-446655440108',
    '550e8400-e29b-41d4-a716-446655440208',
    7, -- LAND
    'TER-NDL-008',
    NULL,
    0,
    0,
    800.00,
    0.00,
    'Lote en barrio privado Nordelta con vista al agua. 800m2, con todos los servicios disponibles. Apto para construcción de casa familiar. Seguridad 24hs, amenities del barrio incluidos.'
  ),
  
  -- Consultorio médico en Recoleta
  (
    '550e8400-e29b-41d4-a716-446655440009',
    '550e8400-e29b-41d4-a716-446655440109',
    '550e8400-e29b-41d4-a716-446655440209',
    12, -- CLINIC
    'CON-REC-009',
    1990,
    0,
    1.0,
    65.00,
    65.00,
    'Consultorio médico en edificio profesional de Recoleta. Totalmente equipado, sala de espera, consultorio principal, baño privado. Excelente ubicación cerca de hospitales y clínicas.'
  ),
  
  -- Cochera cubierta en Belgrano
  (
    '550e8400-e29b-41d4-a716-446655440010',
    '550e8400-e29b-41d4-a716-446655440110',
    '550e8400-e29b-41d4-a716-446655440210',
    13, -- GARAGE
    'COC-BEL-010',
    2008,
    0,
    0,
    25.00,
    25.00,
    'Cochera cubierta en edificio de Belgrano. Fácil acceso, ubicación central, portón automático. Ideal para resguardar vehículo en zona de alta demanda de estacionamiento.'
  );

-- Insert example property management relationships
INSERT INTO management (
  id,
  property_id,
  manager_id,
  manager_type_id,
  start_date,
  commission_percentage
) VALUES
  -- Gestión por persona individual
  (
    '650e8400-e29b-41d4-a716-446655440001',
    '550e8400-e29b-41d4-a716-446655440001', -- Departamento Palermo
    '650e8400-e29b-41d4-a716-446655440301', -- manager_id (persona)
    1, -- PERSON
    '2024-01-15',
    3.50
  ),
  
  -- Gestión por organización inmobiliaria
  (
    '650e8400-e29b-41d4-a716-446655440002',
    '550e8400-e29b-41d4-a716-446655440002', -- Casa San Isidro
    '650e8400-e29b-41d4-a716-446655440401', -- manager_id (organización)
    2, -- ORGANIZATION
    '2024-02-01',
    4.00
  ),
  
  -- Gestión corporativa para oficina
  (
    '650e8400-e29b-41d4-a716-446655440003',
    '550e8400-e29b-41d4-a716-446655440006', -- Oficina Puerto Madero
    '650e8400-e29b-41d4-a716-446655440402', -- manager_id (organización)
    2, -- ORGANIZATION
    '2024-01-01',
    5.00
  ),
  
  -- Gestión de galpón industrial
  (
    '650e8400-e29b-41d4-a716-446655440004',
    '550e8400-e29b-41d4-a716-446655440007', -- Galpón Avellaneda
    '650e8400-e29b-41d4-a716-446655440302', -- manager_id (persona especialista)
    1, -- PERSON
    '2024-03-01',
    4.50
  );

-- Insert example property amenities relationships
INSERT INTO property_amenities (property_id, amenity_id, note) VALUES
  -- Departamento Palermo (moderno, con amenities)
  ('550e8400-e29b-41d4-a716-446655440001', '6', 'Balcón con vista a la calle'),    -- Balcón
  ('550e8400-e29b-41d4-a716-446655440001', '13', 'Luminoso por orientación norte'), -- Luminoso
  ('550e8400-e29b-41d4-a716-446655440001', '31', 'Aire split en dormitorios'),     -- Aire Acondicionado
  ('550e8400-e29b-41d4-a716-446655440001', '38', 'Ascensor de alta velocidad'),    -- Ascensor
  ('550e8400-e29b-41d4-a716-446655440001', '74', 'Gimnasio completo en terraza'),  -- Gimnasio
  ('550e8400-e29b-41d4-a716-446655440001', '71', 'A estrenar, entrega inmediata'), -- A Estrenar
  
  -- Casa San Isidro (familiar, con jardín y pileta)
  ('550e8400-e29b-41d4-a716-446655440002', '9', 'Jardín con césped y plantas'),     -- Jardín
  ('550e8400-e29b-41d4-a716-446655440002', '12', 'Pileta climatizada 8x4 metros'),   -- Pileta
  ('550e8400-e29b-41d4-a716-446655440002', '10', 'Quincho con parrilla y horno'),    -- Quincho
  ('550e8400-e29b-41d4-a716-446655440002', '14', 'Cochera para 2 autos'),           -- Cochera
  ('550e8400-e29b-41d4-a716-446655440002', '16', 'Lavadero independiente'),         -- Lavadero
  ('550e8400-e29b-41d4-a716-446655440002', '32', 'Calefacción central por radiadores'), -- Calefacción Central
  ('550e8400-e29b-41d4-a716-446655440002', '74', 'En muy buen estado'),             -- Muy Bueno
  
  -- PH Villa Crespo (con terraza)
  ('550e8400-e29b-41d4-a716-446655440003', '7', 'Terraza de 30m2 con pérgola'),     -- Terraza
  ('550e8400-e29b-41d4-a716-446655440003', '11', 'Parrilla en terraza'),            -- Parrilla
  ('550e8400-e29b-41d4-a716-446655440003', '13', 'Muy luminoso, orientación norte'), -- Luminoso
  ('550e8400-e29b-41d4-a716-446655440003', '70', 'Totalmente reciclado'),           -- Reciclado
  
  -- Quinta Tigre (con muelle y parque)
  ('550e8400-e29b-41d4-a716-446655440004', '9', 'Parque de 300m2 con árboles'),     -- Jardín
  ('550e8400-e29b-41d4-a716-446655440004', '12', 'Pileta de natación'),             -- Pileta
  ('550e8400-e29b-41d4-a716-446655440004', '10', 'Quincho para 20 personas'),       -- Quincho
  ('550e8400-e29b-41d4-a716-446655440004', '18', 'Vista al río y canales'),         -- Vista al Mar (río)
  ('550e8400-e29b-41d4-a716-446655440004', '81', 'Deck de madera en muelle'),       -- Deck
  
  -- Local comercial Microcentro
  ('550e8400-e29b-41d4-a716-446655440005', '5', 'Sobre avenida principal'),         -- Frente
  ('550e8400-e29b-41d4-a716-446655440005', '13', 'Mucha luz natural'),              -- Luminoso
  ('550e8400-e29b-41d4-a716-446655440005', '1', 'Ideal para cualquier rubro'),      -- Apto Crédito
  
  -- Oficina Puerto Madero (corporativa)
  ('550e8400-e29b-41d4-a716-446655440006', '19', 'Vista panorámica al río'),        -- Vista Panorámica
  ('550e8400-e29b-41d4-a716-446655440006', '31', 'Aire acondicionado centralizado'), -- Aire Acondicionado
  ('550e8400-e29b-41d4-a716-446655440006', '38', 'Ascensores de alta velocidad'),   -- Ascensor
  ('550e8400-e29b-41d4-a716-446655440006', '57', 'Seguridad 24hs con recepcionista'), -- Seguridad 24hs
  ('550e8400-e29b-41d4-a716-446655440006', '14', '4 cocheras cubiertas incluidas'), -- Cochera
  
  -- Galpón Avellaneda (industrial)
  ('550e8400-e29b-41d4-a716-446655440007', '21', 'Entrada independiente para camiones'), -- Entrada Independiente
  ('550e8400-e29b-41d4-a716-446655440007', '39', 'Montacargas para carga pesada'),      -- Montacargas
  ('550e8400-e29b-41d4-a716-446655440007', '27', 'Electricidad trifásica'),             -- Electricidad
  ('550e8400-e29b-41d4-a716-446655440007', '40', 'Generador de emergencia'),            -- Generador Eléctrico
  
  -- Terreno Nordelta
  ('550e8400-e29b-41d4-a716-446655440008', '62', 'Barrio cerrado con amenities'),       -- Barrio Cerrado
  ('550e8400-e29b-41d4-a716-446655440008', '57', 'Seguridad 24hs del barrio'),          -- Seguridad 24hs
  ('550e8400-e29b-41d4-a716-446655440008', '19', 'Vista al canal'),                     -- Vista Panorámica
  
  -- Consultorio Recoleta
  ('550e8400-e29b-41d4-a716-446655440009', '2', 'Habilitado para actividad profesional'), -- Apto Profesional
  ('550e8400-e29b-41d4-a716-446655440009', '31', 'Aire acondicionado frío/calor'),        -- Aire Acondicionado
  ('550e8400-e29b-41d4-a716-446655440009', '38', 'Edificio con ascensor'),                -- Ascensor
  
  -- Cochera Belgrano
  ('550e8400-e29b-41d4-a716-446655440010', '56', 'Portón automático'),                    -- Portero Eléctrico
  ('550e8400-e29b-41d4-a716-446655440010', '21', 'Acceso independiente');                 -- Entrada Independiente
