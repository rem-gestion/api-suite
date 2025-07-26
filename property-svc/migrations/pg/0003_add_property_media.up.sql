-- Migration: Add property_media table
-- Description: Create property_media table for managing property images, videos, and documents

CREATE TABLE IF NOT EXISTS property_media (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    property_id UUID NOT NULL,
    
    -- Media information
    media_type VARCHAR(50) NOT NULL CHECK (media_type IN ('image', 'video', 'document', 'virtual_tour', 'floor_plan', 'blueprint')),
    file_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_url VARCHAR(500),
    file_size BIGINT,
    mime_type VARCHAR(100),
    
    -- Media metadata
    title VARCHAR(255),
    description TEXT,
    alt_text VARCHAR(255),
    
    -- Display settings
    display_order INTEGER DEFAULT 0,
    is_main BOOLEAN DEFAULT FALSE,
    is_public BOOLEAN DEFAULT TRUE,
    is_featured BOOLEAN DEFAULT FALSE,
    
    -- Image/Video specific metadata
    width INTEGER,
    height INTEGER,
    duration INTEGER, -- for videos, in seconds
    quality VARCHAR(50), -- 'low', 'medium', 'high', 'ultra'
    
    -- Additional metadata
    metadata JSONB DEFAULT '{}',
    tags TEXT[],
    
    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by UUID,
    updated_by UUID,
    deleted_by UUID,
    
    -- Foreign key constraints
    CONSTRAINT fk_property_media_property_id 
        FOREIGN KEY (property_id) REFERENCES property(id) ON DELETE CASCADE,
    
    -- Constraints
    CONSTRAINT chk_property_media_display_order CHECK (display_order >= 0),
    CONSTRAINT chk_property_media_file_size CHECK (file_size IS NULL OR file_size > 0),
    CONSTRAINT chk_property_media_dimensions CHECK (
        (width IS NULL AND height IS NULL) OR 
        (width > 0 AND height > 0)
    )
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_property_media_property_id ON property_media(property_id);
CREATE INDEX IF NOT EXISTS idx_property_media_type ON property_media(media_type);
CREATE INDEX IF NOT EXISTS idx_property_media_display_order ON property_media(property_id, display_order);
CREATE INDEX IF NOT EXISTS idx_property_media_is_main ON property_media(property_id, is_main) WHERE is_main = TRUE;
CREATE INDEX IF NOT EXISTS idx_property_media_is_public ON property_media(is_public);
CREATE INDEX IF NOT EXISTS idx_property_media_is_featured ON property_media(is_featured) WHERE is_featured = TRUE;
CREATE INDEX IF NOT EXISTS idx_property_media_created_at ON property_media(created_at);
CREATE INDEX IF NOT EXISTS idx_property_media_deleted_at ON property_media(deleted_at) WHERE deleted_at IS NULL;

-- Composite indexes for common queries
CREATE INDEX IF NOT EXISTS idx_property_media_property_type_order 
    ON property_media(property_id, media_type, display_order, deleted_at) 
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_property_media_public_display 
    ON property_media(property_id, is_public, display_order, deleted_at) 
    WHERE deleted_at IS NULL AND is_public = TRUE;

-- Update trigger for updated_at
CREATE OR REPLACE FUNCTION update_property_media_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_property_media_updated_at
    BEFORE UPDATE ON property_media
    FOR EACH ROW
    EXECUTE FUNCTION update_property_media_updated_at();

-- Trigger to ensure only one main image per property
CREATE OR REPLACE FUNCTION ensure_single_main_media()
RETURNS TRIGGER AS $$
BEGIN
    -- If setting this media as main, unset all other main media for this property
    IF NEW.is_main = TRUE THEN
        UPDATE property_media 
        SET is_main = FALSE, updated_at = NOW()
        WHERE property_id = NEW.property_id 
          AND id != NEW.id 
          AND is_main = TRUE
          AND deleted_at IS NULL;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_ensure_single_main_media
    BEFORE INSERT OR UPDATE ON property_media
    FOR EACH ROW
    WHEN (NEW.is_main = TRUE)
    EXECUTE FUNCTION ensure_single_main_media();


INSERT INTO property_media (id, property_id, media_type, file_name, file_path, file_url, title, description, display_order, is_main, file_size, mime_type, width, height, alt_text, created_at, updated_at) VALUES
-- Casa en Palermo - Medios (using property ID: 550e8400-e29b-41d4-a716-446655440001)
('aaaaaaaa-1111-aaaa-1111-aaaaaaaaaaaa'::uuid, '550e8400-e29b-41d4-a716-446655440001'::uuid, 'image', 'casa-palermo-frente.jpg', '/uploads/properties/casa-palermo-frente.jpg', 'https://cdn.example.com/props/casa-palermo-frente.jpg', 'Frente de la casa', 'Vista frontal de la propiedad desde la calle', 1, true, 2621440, 'image/jpeg', 1920, 1080, 'Frente de casa en Palermo', NOW(), NOW()),
('aaaaaaaa-2222-aaaa-2222-aaaaaaaaaaaa'::uuid, '550e8400-e29b-41d4-a716-446655440001'::uuid, 'image', 'casa-palermo-living.jpg', '/uploads/properties/casa-palermo-living.jpg', 'https://cdn.example.com/props/casa-palermo-living.jpg', 'Living comedor', 'Amplio living comedor con grandes ventanales', 2, false, 2202009, 'image/jpeg', 1920, 1080, 'Living de casa en Palermo', NOW(), NOW()),
('aaaaaaaa-3333-aaaa-3333-aaaaaaaaaaaa'::uuid, '550e8400-e29b-41d4-a716-446655440001'::uuid, 'image', 'casa-palermo-jardin.jpg', '/uploads/properties/casa-palermo-jardin.jpg', 'https://cdn.example.com/props/casa-palermo-jardin.jpg', 'Jardín privado', 'Hermoso jardín con césped y plantas', 3, false, 2936012, 'image/jpeg', 1920, 1080, 'Jardín de casa en Palermo', NOW(), NOW()),
('aaaaaaaa-4444-aaaa-4444-aaaaaaaaaaaa'::uuid, '550e8400-e29b-41d4-a716-446655440001'::uuid, 'floor_plan', 'casa-palermo-planos.pdf', '/uploads/properties/casa-palermo-planos.pdf', 'https://cdn.example.com/props/casa-palermo-planos.pdf', 'Planos de la casa', 'Planos arquitectónicos completos', 4, false, 1258291, 'application/pdf', 2100, 2970, 'Planos de casa en Palermo', NOW(), NOW()),

-- Departamento en Recoleta - Medios (using property ID: 550e8400-e29b-41d4-a716-446655440002)
('bbbbbbbb-1111-bbbb-1111-bbbbbbbbbbbb'::uuid, '550e8400-e29b-41d4-a716-446655440002'::uuid, 'image', 'depto-recoleta-living.jpg', '/uploads/properties/depto-recoleta-living.jpg', 'https://cdn.example.com/props/depto-recoleta-living.jpg', 'Living con vista', 'Living comedor con vista panorámica', 1, true, 1992294, 'image/jpeg', 1920, 1080, 'Living de departamento en Recoleta', NOW(), NOW()),
('bbbbbbbb-2222-bbbb-2222-bbbbbbbbbbbb'::uuid, '550e8400-e29b-41d4-a716-446655440002'::uuid, 'image', 'depto-recoleta-dormitorio.jpg', '/uploads/properties/depto-recoleta-dormitorio.jpg', 'https://cdn.example.com/props/depto-recoleta-dormitorio.jpg', 'Dormitorio principal', 'Dormitorio principal con placard empotrado', 2, false, 1782579, 'image/jpeg', 1920, 1080, 'Dormitorio de departamento en Recoleta', NOW(), NOW()),
('bbbbbbbb-3333-bbbb-3333-bbbbbbbbbbbb'::uuid, '550e8400-e29b-41d4-a716-446655440002'::uuid, 'video', 'depto-recoleta-tour.mp4', '/uploads/properties/depto-recoleta-tour.mp4', 'https://cdn.example.com/props/depto-recoleta-tour.mp4', 'Video tour del departamento', 'Recorrido completo por todo el departamento', 3, false, 27048550, 'video/mp4', 1920, 1080, 'Video tour de departamento en Recoleta', NOW(), NOW()),

-- Oficina céntrica - Medios (using property ID: 550e8400-e29b-41d4-a716-446655440003)
('cccccccc-1111-cccc-1111-cccccccccccc'::uuid, '550e8400-e29b-41d4-a716-446655440003'::uuid, 'image', 'oficina-micro-recepcion.jpg', '/uploads/properties/oficina-micro-recepcion.jpg', 'https://cdn.example.com/props/oficina-micro-recepcion.jpg', 'Área de recepción', 'Recepción moderna con mobiliario corporativo', 1, true, 2306867, 'image/jpeg', 1920, 1080, 'Recepción de oficina en Microcentro', NOW(), NOW()),
('cccccccc-2222-cccc-2222-cccccccccccc'::uuid, '550e8400-e29b-41d4-a716-446655440003'::uuid, 'image', 'oficina-micro-openspace.jpg', '/uploads/properties/oficina-micro-openspace.jpg', 'https://cdn.example.com/props/oficina-micro-openspace.jpg', 'Open space de trabajo', 'Amplio espacio de trabajo colaborativo', 2, false, 2097152, 'image/jpeg', 1920, 1080, 'Open space de oficina en Microcentro', NOW(), NOW()),

-- Quinta en zona norte - Medios (using property ID: 550e8400-e29b-41d4-a716-446655440004)
('dddddddd-1111-dddd-1111-dddddddddddd'::uuid, '550e8400-e29b-41d4-a716-446655440004'::uuid, 'image', 'quinta-norte-pileta.jpg', '/uploads/properties/quinta-norte-pileta.jpg', 'https://cdn.example.com/props/quinta-norte-pileta.jpg', 'Piscina y quincho', 'Área de piscina con quincho equipado', 1, true, 3250585, 'image/jpeg', 1920, 1080, 'Piscina de quinta en zona norte', NOW(), NOW()),
('dddddddd-2222-dddd-2222-dddddddddddd'::uuid, '550e8400-e29b-41d4-a716-446655440004'::uuid, 'image', 'quinta-norte-parque.jpg', '/uploads/properties/quinta-norte-parque.jpg', 'https://cdn.example.com/props/quinta-norte-parque.jpg', 'Parque y jardines', 'Extenso parque con árboles frutales', 2, false, 3565158, 'image/jpeg', 1920, 1080, 'Parque de quinta en zona norte', NOW(), NOW()),
('dddddddd-3333-dddd-3333-dddddddddddd'::uuid, '550e8400-e29b-41d4-a716-446655440004'::uuid, 'virtual_tour', 'quinta-norte-360.html', '/uploads/properties/quinta-norte-360.html', 'https://cdn.example.com/props/quinta-norte-360.html', 'Tour 360° de la quinta', 'Recorrido virtual interactivo por toda la propiedad', 3, false, 15942451, 'text/html', 2048, 1024, 'Tour 360° de quinta en zona norte', NOW(), NOW()),

-- Departamento en Córdoba - Medios (using property ID: 550e8400-e29b-41d4-a716-446655440005)
('eeeeeeee-1111-eeee-1111-eeeeeeeeeeee'::uuid, '550e8400-e29b-41d4-a716-446655440005'::uuid, 'image', 'depto-cordoba-ambiente.jpg', '/uploads/properties/depto-cordoba-ambiente.jpg', 'https://cdn.example.com/props/depto-cordoba-ambiente.jpg', 'Ambiente principal', 'Living comedor integrado con cocina', 1, true, 1887436, 'image/jpeg', 1920, 1080, 'Ambiente de departamento en Córdoba', NOW(), NOW()),
('eeeeeeee-2222-eeee-2222-eeeeeeeeeeee'::uuid, '550e8400-e29b-41d4-a716-446655440005'::uuid, 'image', 'depto-cordoba-balcon.jpg', '/uploads/properties/depto-cordoba-balcon.jpg', 'https://cdn.example.com/props/depto-cordoba-balcon.jpg', 'Balcón con vista', 'Balcón con vista a la ciudad', 2, false, 1677721, 'image/jpeg', 1920, 1080, 'Balcón de departamento en Córdoba', NOW(), NOW());
