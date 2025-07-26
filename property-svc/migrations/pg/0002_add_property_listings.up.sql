-- Migration: Add property_listings table
-- Description: Create property_listings table for managing property listings with status, pricing, and analytics

CREATE TABLE IF NOT EXISTS property_listings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    property_id UUID NOT NULL,
    operation_type VARCHAR(50) NOT NULL CHECK (operation_type IN ('sale', 'rent', 'lease', 'auction')),
    listing_status VARCHAR(50) NOT NULL DEFAULT 'draft' CHECK (listing_status IN ('draft', 'active', 'inactive', 'pending', 'expired', 'sold', 'rented', 'cancelled')),
    
    -- Pricing information
    list_price DECIMAL(15,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    price_per_sqm DECIMAL(10,2),
    price_history JSONB DEFAULT '[]',
    
    -- Listing details
    title VARCHAR(255) NOT NULL,
    description TEXT,
    highlights TEXT[],
    features TEXT[],
    
    -- Marketing settings
    highlighted BOOLEAN DEFAULT FALSE,
    publish_on_web BOOLEAN DEFAULT TRUE,
    publish_on_portals BOOLEAN DEFAULT FALSE,
    portal_urls JSONB DEFAULT '{}',
    
    -- Analytics
    views_count INTEGER DEFAULT 0,
    inquiries_count INTEGER DEFAULT 0,
    favorites_count INTEGER DEFAULT 0,
    last_viewed_at TIMESTAMP WITH TIME ZONE,
    
    -- Contact information
    contact_name VARCHAR(255),
    contact_phone VARCHAR(50),
    contact_email VARCHAR(255),
    contact_schedule TEXT,
    
    -- Dates
    available_from DATE,
    available_until DATE,
    last_price_update TIMESTAMP WITH TIME ZONE,
    
    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by UUID,
    updated_by UUID,
    deleted_by UUID,
    
    -- Foreign key constraints
    CONSTRAINT fk_property_listings_property_id 
        FOREIGN KEY (property_id) REFERENCES property(id) ON DELETE CASCADE
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_property_listings_property_id ON property_listings(property_id);
CREATE INDEX IF NOT EXISTS idx_property_listings_operation_type ON property_listings(operation_type);
CREATE INDEX IF NOT EXISTS idx_property_listings_status ON property_listings(listing_status);
CREATE INDEX IF NOT EXISTS idx_property_listings_price ON property_listings(list_price);
CREATE INDEX IF NOT EXISTS idx_property_listings_currency ON property_listings(currency);
CREATE INDEX IF NOT EXISTS idx_property_listings_highlighted ON property_listings(highlighted);
CREATE INDEX IF NOT EXISTS idx_property_listings_publish_web ON property_listings(publish_on_web);
CREATE INDEX IF NOT EXISTS idx_property_listings_created_at ON property_listings(created_at);
CREATE INDEX IF NOT EXISTS idx_property_listings_deleted_at ON property_listings(deleted_at) WHERE deleted_at IS NULL;

-- Composite indexes for common queries
CREATE INDEX IF NOT EXISTS idx_property_listings_active_listings 
    ON property_listings(operation_type, listing_status, deleted_at) 
    WHERE deleted_at IS NULL AND listing_status = 'active';

CREATE INDEX IF NOT EXISTS idx_property_listings_price_range 
    ON property_listings(operation_type, currency, list_price, deleted_at) 
    WHERE deleted_at IS NULL;

-- Update trigger for updated_at
CREATE OR REPLACE FUNCTION update_property_listings_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_property_listings_updated_at
    BEFORE UPDATE ON property_listings
    FOR EACH ROW
    EXECUTE FUNCTION update_property_listings_updated_at();

INSERT INTO property_listings (id, property_id, operation_type, listing_status, list_price, currency, price_per_sqm, title, description, highlighted, publish_on_web, publish_on_portals, portal_urls, views_count, inquiries_count, favorites_count, available_from, created_at, updated_at) VALUES
-- Departamento Palermo - Venta
('11111111-aaaa-1111-aaaa-111111111111'::uuid, '550e8400-e29b-41d4-a716-446655440001'::uuid, 'sale', 'active', 280000.00, 'USD', 3710.00, 'Departamento moderno en Palermo - VENTA', 'Moderno departamento de 2 ambientes en el corazón de Palermo. Luminoso, con balcón y excelente distribución. Edificio con amenities.', true, true, true, '{"zonaprop": true, "argenprop": true, "mercadolibre": true}', 245, 12, 8, '2025-07-01', NOW(), NOW()),

-- Casa San Isidro - Venta
('22222222-aaaa-2222-aaaa-222222222222'::uuid, '550e8400-e29b-41d4-a716-446655440002'::uuid, 'sale', 'active', 850000.00, 'USD', 3035.00, 'Casa familiar en San Isidro - VENTA', 'Hermosa casa familiar en zona residencial de San Isidro. 4 dormitorios, 3 baños, quincho con parrilla, jardín con pileta. Ideal para familias.', true, true, true, '{"zonaprop": true, "argenprop": true}', 156, 18, 5, '2025-07-01', NOW(), NOW()),

-- PH Villa Crespo - Alquiler
('33333333-bbbb-3333-bbbb-333333333333'::uuid, '550e8400-e29b-41d4-a716-446655440003'::uuid, 'rent', 'active', 1800.00, 'USD', 15.00, 'PH en Villa Crespo - ALQUILER', 'PH de 3 ambientes con terraza propia. Totalmente refaccionado, muy luminoso. Terraza de 30m2 con parrilla.', false, true, true, '{"zonaprop": true, "argenprop": true}', 189, 7, 12, '2025-07-05', NOW(), NOW()),

-- Oficina Puerto Madero - Alquiler
('44444444-cccc-4444-cccc-444444444444'::uuid, '550e8400-e29b-41d4-a716-446655440006'::uuid, 'rent', 'active', 8500.00, 'USD', 26.56, 'Oficina premium en Puerto Madero - ALQUILER', 'Oficina premium en edificio corporativo de Puerto Madero. Piso completo con vista al río. 4 cocheras incluidas.', true, true, false, '{"zonaprop": true}', 78, 9, 3, '2025-07-10', NOW(), NOW()),

-- Quinta Tigre - Alquiler temporal
('55555555-dddd-5555-dddd-555555555555'::uuid, '550e8400-e29b-41d4-a716-446655440004'::uuid, 'rent', 'active', 2500.00, 'USD', 5.56, 'Quinta en Tigre - ALQUILER TEMPORAL', 'Quinta familiar en Tigre con amplio parque y dock privado. Casa principal de 3 dormitorios, quincho independiente, pileta.', false, true, true, '{"airbnb": true, "booking": true}', 92, 15, 6, '2025-07-15', NOW(), NOW()),

-- Local comercial Microcentro - Alquiler
('66666666-eeee-6666-eeee-666666666666'::uuid, '550e8400-e29b-41d4-a716-446655440005'::uuid, 'rent', 'active', 2800.00, 'USD', 32.94, 'Local comercial en Microcentro - ALQUILER', 'Local comercial en pleno Microcentro, sobre avenida principal. Amplia vidriera, excelente ubicación para cualquier tipo de comercio.', false, true, true, '{"zonaprop": true, "argenprop": true}', 134, 22, 4, '2025-07-20', NOW(), NOW());
