-- Migration: Add property_valuations table
-- Description: Create property_valuations table for managing property valuations and market history

CREATE TABLE IF NOT EXISTS property_valuations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    property_id UUID NOT NULL,
    
    -- Valuation information
    valuation_date DATE NOT NULL,
    valuation_type VARCHAR(100) NOT NULL, -- 'market', 'insurance', 'tax', 'mortgage', 'investment', 'liquidation'
    appraised_value DECIMAL(15,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    value_per_sqm DECIMAL(10,2),
    
    -- Appraiser information
    appraiser_name VARCHAR(255),
    appraiser_license VARCHAR(100),
    appraiser_organization VARCHAR(255),
    
    -- Valuation methodology
    valuation_method VARCHAR(100), -- 'comparative', 'income', 'cost', 'automated', 'hybrid'
    market_conditions TEXT,
    adjustments_applied TEXT,
    comparable_properties JSONB DEFAULT '[]',
    
    -- Supporting documentation
    valuation_report_url VARCHAR(500),
    photos_urls TEXT[],
    
    -- Validity and purpose
    valid_until DATE,
    purpose VARCHAR(255), -- 'sale', 'purchase', 'financing', 'insurance', 'tax_assessment', 'legal'
    
    -- Additional data
    notes TEXT,
    metadata JSONB DEFAULT '{}',
    
    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by UUID,
    updated_by UUID,
    deleted_by UUID,
    
    -- Foreign key constraints
    CONSTRAINT fk_property_valuations_property_id 
        FOREIGN KEY (property_id) REFERENCES property(id) ON DELETE CASCADE,
    
    -- Constraints
    CONSTRAINT chk_property_valuations_appraised_value CHECK (appraised_value > 0),
    CONSTRAINT chk_property_valuations_value_per_sqm CHECK (value_per_sqm IS NULL OR value_per_sqm > 0),
    CONSTRAINT chk_property_valuations_valid_until CHECK (valid_until IS NULL OR valid_until >= valuation_date)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_property_valuations_property_id ON property_valuations(property_id);
CREATE INDEX IF NOT EXISTS idx_property_valuations_valuation_date ON property_valuations(valuation_date);
CREATE INDEX IF NOT EXISTS idx_property_valuations_type ON property_valuations(valuation_type);
CREATE INDEX IF NOT EXISTS idx_property_valuations_appraised_value ON property_valuations(appraised_value);
CREATE INDEX IF NOT EXISTS idx_property_valuations_currency ON property_valuations(currency);
CREATE INDEX IF NOT EXISTS idx_property_valuations_appraiser ON property_valuations(appraiser_name);
CREATE INDEX IF NOT EXISTS idx_property_valuations_method ON property_valuations(valuation_method);
CREATE INDEX IF NOT EXISTS idx_property_valuations_valid_until ON property_valuations(valid_until);
CREATE INDEX IF NOT EXISTS idx_property_valuations_created_at ON property_valuations(created_at);
CREATE INDEX IF NOT EXISTS idx_property_valuations_deleted_at ON property_valuations(deleted_at) WHERE deleted_at IS NULL;

-- Composite indexes for common queries
CREATE INDEX IF NOT EXISTS idx_property_valuations_property_date_type 
    ON property_valuations(property_id, valuation_date DESC, valuation_type, deleted_at) 
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_property_valuations_market_history 
    ON property_valuations(property_id, valuation_date DESC, appraised_value, deleted_at) 
    WHERE deleted_at IS NULL AND valuation_type = 'market';

CREATE INDEX IF NOT EXISTS idx_property_valuations_latest 
    ON property_valuations(property_id, valuation_date DESC, deleted_at) 
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_property_valuations_valid_current 
    ON property_valuations(property_id, valuation_type, valid_until, deleted_at) 
    WHERE deleted_at IS NULL;

-- Update trigger for updated_at
CREATE OR REPLACE FUNCTION update_property_valuations_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_property_valuations_updated_at
    BEFORE UPDATE ON property_valuations
    FOR EACH ROW
    EXECUTE FUNCTION update_property_valuations_updated_at();


INSERT INTO property_valuations (id, property_id, valuation_type, appraised_value, currency, valuation_date, valid_until, valuation_method, appraiser_name, appraiser_license, notes, market_conditions, created_at, updated_at) VALUES
-- Casa en Palermo - Valuación de mercado
('550e8400-e29b-41d4-a716-446655441001'::uuid, '550e8400-e29b-41d4-a716-446655440001'::uuid, 'market', 450000.00, 'USD', '2025-06-15', '2025-12-15', 'Método comparativo de mercado', 'Ing. Patricia Morales', 'CMCABA 12345', 'Valuación realizada considerando propiedades similares en la zona de Palermo', 'Mercado inmobiliario estable con tendencia al alza en Palermo', NOW(), NOW()),

-- Casa en Palermo - Valuación para seguro
('550e8400-e29b-41d4-a716-446655441002'::uuid, '550e8400-e29b-41d4-a716-446655440001'::uuid, 'insurance', 380000.00, 'USD', '2025-06-20', '2026-06-20', 'Costo de reposición', 'Arq. Roberto Silva', 'CPAU 54321', 'Valuación para póliza de seguro contra incendio y otros riesgos', 'Costos de construcción en alza debido a inflación', NOW(), NOW()),

-- Departamento en Recoleta - Valuación de mercado
('550e8400-e29b-41d4-a716-446655441003'::uuid, '550e8400-e29b-41d4-a716-446655440002'::uuid, 'market', 280000.00, 'USD', '2025-07-01', '2026-01-01', 'Método comparativo de mercado', 'Lic. Carmen Fernández', 'CUCICBA 67890', 'Valuación basada en comparables de edificios similares en Recoleta', 'Recoleta mantiene precios altos por ubicación premium', NOW(), NOW()),

-- Oficina céntrica - Valuación comercial
('550e8400-e29b-41d4-a716-446655441004'::uuid, '550e8400-e29b-41d4-a716-446655440003'::uuid, 'commercial', 420000.00, 'USD', '2025-07-05', '2026-07-05', 'Capitalización de rentas', 'Mg. Fernando López', 'CMCABA 98765', 'Valuación comercial basada en renta potencial y comparables del área', 'Mercado de oficinas en recuperación post-pandemia', NOW(), NOW()),

-- Quinta en zona norte - Valuación rural/recreativa
('550e8400-e29b-41d4-a716-446655441005'::uuid, '550e8400-e29b-41d4-a716-446655440004'::uuid, 'market', 320000.00, 'USD', '2025-07-10', '2026-01-10', 'Método residual', 'Ing. Agr. Marcos Aguirre', 'CIAER 13579', 'Valuación considerando uso recreativo y potencial de desarrollo', 'Creciente demanda de propiedades recreativas post-pandemia', NOW(), NOW()),

-- Departamento en Córdoba - Valuación local
('550e8400-e29b-41d4-a716-446655441006'::uuid, '550e8400-e29b-41d4-a716-446655440005'::uuid, 'market', 115000.00, 'USD', '2025-07-15', '2026-01-15', 'Método comparativo', 'Arq. Silvia Herrera', 'CAPBA 24680', 'Valuación para el mercado cordobés considerando zona universitaria', 'Mercado estudiantil/profesional muy activo en la zona', NOW(), NOW()),

-- Valuaciones adicionales (históricas)
-- Casa en Palermo - Valuación anterior
('550e8400-e29b-41d4-a716-446655441007'::uuid, '550e8400-e29b-41d4-a716-446655440001'::uuid, 'market', 420000.00, 'USD', '2024-12-15', '2025-06-15', 'Método comparativo de mercado', 'Ing. Patricia Morales', 'CMCABA 12345', 'Valuación anterior para seguimiento de evolución del mercado', 'Mercado en crecimiento sostenido', NOW(), NOW());
