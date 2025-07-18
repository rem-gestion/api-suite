-- Insertar amenities básicas para testing
INSERT INTO amenity (id, name, description, icon, category, created_at) VALUES
('11111111-1111-1111-1111-111111111111', 'Piscina', 'Piscina climatizada', 'pool', 'recreation', now()),
('22222222-2222-2222-2222-222222222222', 'Gimnasio', 'Gimnasio con equipamiento completo', 'gym', 'recreation', now()),
('33333333-3333-3333-3333-333333333333', 'Portero 24hs', 'Servicio de portería las 24 horas', 'security', 'security', now()),
('44444444-4444-4444-4444-444444444444', 'Cochera', 'Cochera cubierta', 'garage', 'services', now()),
('55555555-5555-5555-5555-555555555555', 'Balcón', 'Balcón con vista', 'balcony', 'amenity', now()),
('66666666-6666-6666-6666-666666666666', 'Terraza', 'Terraza privada', 'terrace', 'amenity', now()),
('77777777-7777-7777-7777-777777777777', 'Aire Acondicionado', 'Aire acondicionado central', 'ac', 'services', now()),
('88888888-8888-8888-8888-888888888888', 'Laundry', 'Servicio de lavandería', 'laundry', 'services', now()),
('99999999-9999-9999-9999-999999999999', 'Jardín', 'Jardín compartido', 'garden', 'recreation', now()),
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'Parrilla', 'Zona de parrillas', 'grill', 'recreation', now());

-- Datos dummy para propiedades (usando IDs mock)
INSERT INTO property (id, owner_person_id, address_id, property_type, internal_code, year_built, bedrooms, bathrooms, total_area_sqm, covered_area_sqm, description, created_at) VALUES
('aaaa1111-aaaa-1111-aaaa-111111111111', '11111111-1111-1111-1111-111111111111', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'APARTMENT', 'DEPT-001', 2020, 2, 1.5, 85.50, 75.00, 'Departamento moderno en zona céntrica con excelente iluminación', now()),
('bbbb2222-bbbb-2222-bbbb-222222222222', '22222222-2222-2222-2222-222222222222', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'HOUSE', 'CASA-001', 2018, 3, 2.0, 150.00, 120.00, 'Casa familiar con jardín en barrio tranquilo', now()),
('cccc3333-cccc-3333-cccc-333333333333', '11111111-1111-1111-1111-111111111111', 'cccccccc-cccc-cccc-cccc-cccccccccccc', 'COMMERCIAL_SPACE', 'LOCAL-001', 2019, 0, 1.0, 80.00, 80.00, 'Local comercial sobre avenida principal con gran exposición', now()),
('dddd4444-dddd-4444-dddd-444444444444', '33333333-3333-3333-3333-333333333333', 'dddddddd-dddd-dddd-dddd-dddddddddddd', 'OFFICE', 'OFICINA-001', 2021, 0, 2.0, 95.00, 95.00, 'Oficina corporativa en torre moderna con vista panorámica', now()),
('eeee5555-eeee-5555-eeee-555555555555', '44444444-4444-4444-4444-444444444444', 'eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', 'LAND', 'TERRENO-001', NULL, NULL, NULL, 500.00, 0.00, 'Terreno para desarrollo en zona en crecimiento', now());

-- Datos dummy para property management (usando organization IDs mock)
INSERT INTO property_management (id, property_id, organization_id, managed_since, is_active, commission_percent, notes, created_at) VALUES
('aaaa1111-bbbb-2222-cccc-333333333333', 'aaaa1111-aaaa-1111-aaaa-111111111111', '99999999-9999-9999-9999-999999999999', '2024-01-01', true, 3.5, 'Gestión integral del departamento céntrico', now()),
('bbbb2222-cccc-3333-dddd-444444444444', 'bbbb2222-bbbb-2222-bbbb-222222222222', '99999999-9999-9999-9999-999999999999', '2024-02-15', true, 4.0, 'Administración de casa familiar', now()),
('cccc3333-dddd-4444-eeee-555555555555', 'cccc3333-cccc-3333-cccc-333333333333', '88888888-8888-8888-8888-888888888888', '2024-03-01', true, 5.0, 'Gestión comercial especializada', now());

-- Relacionar propiedades con amenities
INSERT INTO property_amenities (id, property_id, amenity_id, created_at) VALUES
-- Departamento moderno con piscina, gym, portero, cochera
('1111aaaa-1111-aaaa-1111-aaaaaaaaaaaa', 'aaaa1111-aaaa-1111-aaaa-111111111111', '11111111-1111-1111-1111-111111111111', now()),
('2222bbbb-2222-bbbb-2222-bbbbbbbbbbbb', 'aaaa1111-aaaa-1111-aaaa-111111111111', '22222222-2222-2222-2222-222222222222', now()),
('3333cccc-3333-cccc-3333-cccccccccccc', 'aaaa1111-aaaa-1111-aaaa-111111111111', '33333333-3333-3333-3333-333333333333', now()),
('4444dddd-4444-dddd-4444-dddddddddddd', 'aaaa1111-aaaa-1111-aaaa-111111111111', '44444444-4444-4444-4444-444444444444', now()),
('5555eeee-5555-eeee-5555-eeeeeeeeeeee', 'aaaa1111-aaaa-1111-aaaa-111111111111', '55555555-5555-5555-5555-555555555555', now()),
('6666ffff-6666-ffff-6666-ffffffffffff', 'aaaa1111-aaaa-1111-aaaa-111111111111', '77777777-7777-7777-7777-777777777777', now()),

-- Casa familiar con jardín, parrilla, cochera
('7777aaaa-7777-aaaa-7777-aaaaaaaaaaaa', 'bbbb2222-bbbb-2222-bbbb-222222222222', '99999999-9999-9999-9999-999999999999', now()),
('8888bbbb-8888-bbbb-8888-bbbbbbbbbbbb', 'bbbb2222-bbbb-2222-bbbb-222222222222', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', now()),
('9999cccc-9999-cccc-9999-cccccccccccc', 'bbbb2222-bbbb-2222-bbbb-222222222222', '44444444-4444-4444-4444-444444444444', now()),
('aaaadddd-aaaa-dddd-aaaa-dddddddddddd', 'bbbb2222-bbbb-2222-bbbb-222222222222', '66666666-6666-6666-6666-666666666666', now()),

-- Local comercial básico
('bbbbeeee-bbbb-eeee-bbbb-eeeeeeeeeeee', 'cccc3333-cccc-3333-cccc-333333333333', '77777777-7777-7777-7777-777777777777', now()),

-- Oficina corporativa con aires, cochera
('ccccffff-cccc-ffff-cccc-ffffffffffff', 'dddd4444-dddd-4444-dddd-444444444444', '77777777-7777-7777-7777-777777777777', now()),
('dddd1111-dddd-1111-dddd-111111111111', 'dddd4444-dddd-4444-dddd-444444444444', '44444444-4444-4444-4444-444444444444', now());
