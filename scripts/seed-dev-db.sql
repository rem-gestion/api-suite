-- ══════════════════════════════════════════════════════════════════
-- Script de Población de Base de Datos para Desarrollo
-- ══════════════════════════════════════════════════════════════════
-- Este script popula la base de datos con datos de prueba para desarrollo
-- Se ejecuta después de las migraciones y antes de iniciar los servicios

-- Nota: Usamos UUIDs específicos para poder referenciarlos entre tablas
-- En producción, estos datos no deberían existir

BEGIN;

-- ═══════════════════════════════════════════════════════════════════════════════
-- 1. ADDRESSES - Crear direcciones primero ya que person las referencia
-- ═══════════════════════════════════════════════════════════════════════════════

INSERT INTO address_addresses (id, floor, unit, street, number, city, state, zip, country, created_at) VALUES
-- Direcciones en Buenos Aires (address.id es character(36))
('11111111-1111-1111-1111-111111111111', NULL, NULL, 'Av. Corrientes', 1234, 'Buenos Aires', 'CABA', '1043', 'AR', NOW()),
('22222222-2222-2222-2222-222222222222', '5', 'A', 'Av. Santa Fe', 2567, 'Buenos Aires', 'CABA', '1123', 'AR', NOW()),
('33333333-3333-3333-3333-333333333333', NULL, NULL, 'Belgrano', 890, 'Mar del Plata', 'Buenos Aires', '7600', 'AR', NOW()),
('44444444-4444-4444-4444-444444444444', '12', 'B', 'Av. Pueyrredón', 1456, 'Buenos Aires', 'CABA', '1118', 'AR', NOW()),
('55555555-5555-5555-5555-555555555555', NULL, NULL, 'San Martín', 345, 'Rosario', 'Santa Fe', '2000', 'AR', NOW()),
-- Direcciones en Córdoba
('66666666-6666-6666-6666-666666666666', '3', 'C', 'Dean Funes', 234, 'Córdoba', 'Córdoba', '5000', 'AR', NOW()),
('77777777-7777-7777-7777-777777777777', NULL, NULL, 'Av. Colón', 789, 'Córdoba', 'Córdoba', '5000', 'AR', NOW()),
-- Oficinas/Empresas
('88888888-8888-8888-8888-888888888888', '15', 'OF 1501', 'Av. Leandro N. Alem', 456, 'Buenos Aires', 'CABA', '1001', 'AR', NOW()),
('99999999-9999-9999-9999-999999999999', '8', 'OF 804', 'Av. 9 de Julio', 1020, 'Buenos Aires', 'CABA', '1063', 'AR', NOW()),
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', NULL, NULL, 'Rivadavia', 567, 'Mendoza', 'Mendoza', '5500', 'AR', NOW());

-- ═══════════════════════════════════════════════════════════════════════════════
-- 2. ACCOUNTS - Crear cuentas de autenticación
-- ═══════════════════════════════════════════════════════════════════════════════

INSERT INTO auth_accounts (id, provider, email, password_hash, status, created_at) VALUES
-- Usuarios individuales activos (auth_accounts.id es uuid)
('11111111-1111-1111-1111-111111111111'::uuid, 'email', 'juan.perez@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'active', NOW()),
('22222222-2222-2222-2222-222222222222'::uuid, 'email', 'maria.gonzalez@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'active', NOW()),
('33333333-3333-3333-3333-333333333333'::uuid, 'email', 'carlos.rodriguez@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'pending', NOW()),
('44444444-4444-4444-4444-444444444444'::uuid, 'email', 'ana.martinez@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'active', NOW()),
('55555555-5555-5555-5555-555555555555'::uuid, 'google', 'laura.fernandez@gmail.com', NULL, 'active', NOW()),
-- Cuentas para empresas
('66666666-6666-6666-6666-666666666666'::uuid, 'email', 'admin@acmecorp.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'active', NOW()),
('77777777-7777-7777-7777-777777777777'::uuid, 'email', 'contacto@techsolutions.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'active', NOW()),
('88888888-8888-8888-8888-888888888888'::uuid, 'email', 'ventas@innovatech.com.ar', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'active', NOW());

-- ═══════════════════════════════════════════════════════════════════════════════
-- 3. PERSON - Crear personas (individuales y empresas)
-- ═══════════════════════════════════════════════════════════════════════════════

INSERT INTO person_person (id, type, address_id, avatar_url, sexo, created_at) VALUES
-- Personas individuales (id es uuid, address_id debe coincidir con address.id que es character(36))
-- Convertimos las IDs de address de string a uuid para que coincidan con el tipo de la columna
('11111111-1111-1111-1111-111111111111'::uuid, 'individual', '11111111-1111-1111-1111-111111111111'::uuid, 'https://avatar.example.com/juan.jpg', 'masculino', NOW()),
('22222222-2222-2222-2222-222222222222'::uuid, 'individual', '22222222-2222-2222-2222-222222222222'::uuid, 'https://avatar.example.com/maria.jpg', 'femenino', NOW()),
('33333333-3333-3333-3333-333333333333'::uuid, 'individual', '33333333-3333-3333-3333-333333333333'::uuid, NULL, 'masculino', NOW()),
('44444444-4444-4444-4444-444444444444'::uuid, 'individual', '44444444-4444-4444-4444-444444444444'::uuid, 'https://avatar.example.com/ana.jpg', 'femenino', NOW()),
('55555555-5555-5555-5555-555555555555'::uuid, 'individual', '55555555-5555-5555-5555-555555555555'::uuid, NULL, 'femenino', NOW()),
-- Empresas
('66666666-6666-6666-6666-666666666666'::uuid, 'company', '88888888-8888-8888-8888-888888888888'::uuid, NULL, NULL, NOW()),
('77777777-7777-7777-7777-777777777777'::uuid, 'company', '99999999-9999-9999-9999-999999999999'::uuid, NULL, NULL, NOW()),
('88888888-8888-8888-8888-888888888888'::uuid, 'company', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'::uuid, NULL, NULL, NOW());

-- ═══════════════════════════════════════════════════════════════════════════════
-- 4. INDIVIDUAL - Detalles de personas individuales
-- ═══════════════════════════════════════════════════════════════════════════════

INSERT INTO person_individual (person_id, first_name, last_name, dni, created_at) VALUES
('11111111-1111-1111-1111-111111111111'::uuid, 'Juan Carlos', 'Pérez', '12345678', NOW()),
('22222222-2222-2222-2222-222222222222'::uuid, 'María Elena', 'González', '23456789', NOW()),
('33333333-3333-3333-3333-333333333333'::uuid, 'Carlos Alberto', 'Rodríguez', '34567890', NOW()),
('44444444-4444-4444-4444-444444444444'::uuid, 'Ana Sofía', 'Martínez', '45678901', NOW()),
('55555555-5555-5555-5555-555555555555'::uuid, 'Laura Patricia', 'Fernández', '56789012', NOW());

-- ═══════════════════════════════════════════════════════════════════════════════
-- 5. COMPANY - Detalles de empresas
-- ═══════════════════════════════════════════════════════════════════════════════

INSERT INTO person_company (person_id, legal_name, cuit, society_type, created_at) VALUES
('66666666-6666-6666-6666-666666666666'::uuid, 'ACME Corporation S.A.', '30123456789', 'S.A.', NOW()),
('77777777-7777-7777-7777-777777777777'::uuid, 'Tech Solutions SRL', '30234567890', 'S.R.L.', NOW()),
('88888888-8888-8888-8888-888888888888'::uuid, 'Innova Tech Argentina S.A.', '30345678901', 'S.A.', NOW());

-- ═══════════════════════════════════════════════════════════════════════════════
-- 6. USERS - Conectar accounts con persons
-- ═══════════════════════════════════════════════════════════════════════════════

INSERT INTO auth_users (id, account_id, person_id, onboard_status, last_login, created_at) VALUES
-- Usuarios individuales (todos los IDs son uuid)
('11111111-1111-1111-1111-111111111111'::uuid, '11111111-1111-1111-1111-111111111111'::uuid, '11111111-1111-1111-1111-111111111111'::uuid, 'done', '2025-07-12 10:30:00', NOW()),
('22222222-2222-2222-2222-222222222222'::uuid, '22222222-2222-2222-2222-222222222222'::uuid, '22222222-2222-2222-2222-222222222222'::uuid, 'done', '2025-07-12 09:15:00', NOW()),
('33333333-3333-3333-3333-333333333333'::uuid, '33333333-3333-3333-3333-333333333333'::uuid, '33333333-3333-3333-3333-333333333333'::uuid, 'new', NULL, NOW()),
('44444444-4444-4444-4444-444444444444'::uuid, '44444444-4444-4444-4444-444444444444'::uuid, '44444444-4444-4444-4444-444444444444'::uuid, 'in_progress', '2025-07-11 16:45:00', NOW()),
('55555555-5555-5555-5555-555555555555'::uuid, '55555555-5555-5555-5555-555555555555'::uuid, '55555555-5555-5555-5555-555555555555'::uuid, 'done', '2025-07-12 08:20:00', NOW()),
-- Usuarios de empresas
('66666666-6666-6666-6666-666666666666'::uuid, '66666666-6666-6666-6666-666666666666'::uuid, '66666666-6666-6666-6666-666666666666'::uuid, 'done', '2025-07-12 11:00:00', NOW()),
('77777777-7777-7777-7777-777777777777'::uuid, '77777777-7777-7777-7777-777777777777'::uuid, '77777777-7777-7777-7777-777777777777'::uuid, 'done', '2025-07-12 07:30:00', NOW()),
('88888888-8888-8888-8888-888888888888'::uuid, '88888888-8888-8888-8888-888888888888'::uuid, '88888888-8888-8888-8888-888888888888'::uuid, 'in_progress', '2025-07-11 14:15:00', NOW());

-- ═══════════════════════════════════════════════════════════════════════════════
-- 7. CONTACTO - Información de contacto para personas
-- ═══════════════════════════════════════════════════════════════════════════════

INSERT INTO person_contacto (id, persona_id, tipo, dato, is_primary, created_at) VALUES
-- Contactos para Juan Carlos Pérez (todos los IDs son uuid)
('11111111-1111-1111-1111-111111111111'::uuid, '11111111-1111-1111-1111-111111111111'::uuid, 'email', 'juan.perez@personal.com', true, NOW()),
('11111111-2222-1111-1111-111111111111'::uuid, '11111111-1111-1111-1111-111111111111'::uuid, 'phone', '+541123456789', true, NOW()),
('11111111-3333-1111-1111-111111111111'::uuid, '11111111-1111-1111-1111-111111111111'::uuid, 'whatsapp', '+541123456789', false, NOW()),

-- Contactos para María Elena González
('22222222-1111-2222-2222-222222222222'::uuid, '22222222-2222-2222-2222-222222222222'::uuid, 'email', 'maria.gonzalez@personal.com', true, NOW()),
('22222222-2222-2222-2222-222222222222'::uuid, '22222222-2222-2222-2222-222222222222'::uuid, 'phone', '+541134567890', true, NOW()),
('22222222-3333-2222-2222-222222222222'::uuid, '22222222-2222-2222-2222-222222222222'::uuid, 'email', 'mgonzalez.work@gmail.com', false, NOW()),

-- Contactos para Carlos Alberto Rodríguez
('33333333-1111-3333-3333-333333333333'::uuid, '33333333-3333-3333-3333-333333333333'::uuid, 'email', 'carlos.rodriguez@email.com', true, NOW()),
('33333333-2222-3333-3333-333333333333'::uuid, '33333333-3333-3333-3333-333333333333'::uuid, 'phone', '+542234567891', true, NOW()),

-- Contactos para Ana Sofía Martínez
('44444444-1111-4444-4444-444444444444'::uuid, '44444444-4444-4444-4444-444444444444'::uuid, 'email', 'ana.martinez@personal.com', true, NOW()),
('44444444-2222-4444-4444-444444444444'::uuid, '44444444-4444-4444-4444-444444444444'::uuid, 'phone', '+541145678902', true, NOW()),
('44444444-3333-4444-4444-444444444444'::uuid, '44444444-4444-4444-4444-444444444444'::uuid, 'whatsapp', '+541145678902', false, NOW()),

-- Contactos para Laura Patricia Fernández
('55555555-1111-5555-5555-555555555555'::uuid, '55555555-5555-5555-5555-555555555555'::uuid, 'email', 'laura.fernandez@personal.com', true, NOW()),
('55555555-2222-5555-5555-555555555555'::uuid, '55555555-5555-5555-5555-555555555555'::uuid, 'phone', '+543414567893', true, NOW()),

-- Contactos para empresas
-- ACME Corporation
('66666666-1111-6666-6666-666666666666'::uuid, '66666666-6666-6666-6666-666666666666'::uuid, 'email', 'info@acmecorp.com', true, NOW()),
('66666666-2222-6666-6666-666666666666'::uuid, '66666666-6666-6666-6666-666666666666'::uuid, 'phone', '+541156789013', true, NOW()),
('66666666-3333-6666-6666-666666666666'::uuid, '66666666-6666-6666-6666-666666666666'::uuid, 'email', 'ventas@acmecorp.com', false, NOW()),

-- Tech Solutions SRL
('77777777-1111-7777-7777-777777777777'::uuid, '77777777-7777-7777-7777-777777777777'::uuid, 'email', 'info@techsolutions.com', true, NOW()),
('77777777-2222-7777-7777-777777777777'::uuid, '77777777-7777-7777-7777-777777777777'::uuid, 'phone', '+541167890124', true, NOW()),
('77777777-3333-7777-7777-777777777777'::uuid, '77777777-7777-7777-7777-777777777777'::uuid, 'email', 'soporte@techsolutions.com', false, NOW()),

-- Innova Tech Argentina
('88888888-1111-8888-8888-888888888888'::uuid, '88888888-8888-8888-8888-888888888888'::uuid, 'email', 'contacto@innovatech.com.ar', true, NOW()),
('88888888-2222-8888-8888-888888888888'::uuid, '88888888-8888-8888-8888-888888888888'::uuid, 'phone', '+542614567894', true, NOW()),
('88888888-3333-8888-8888-888888888888'::uuid, '88888888-8888-8888-8888-888888888888'::uuid, 'whatsapp', '+542614567894', false, NOW());

COMMIT;

-- ═══════════════════════════════════════════════════════════════════════════════
-- Resumen de datos insertados
-- ═══════════════════════════════════════════════════════════════════════════════

DO $$
BEGIN
    RAISE NOTICE '═══════════════════════════════════════════════════════════════════';
    RAISE NOTICE 'BASE DE DATOS POBLADA EXITOSAMENTE PARA DESARROLLO';
    RAISE NOTICE '═══════════════════════════════════════════════════════════════════';
    RAISE NOTICE 'Addresses insertadas: %', (SELECT COUNT(*) FROM address_addresses);
    RAISE NOTICE 'Accounts insertadas: %', (SELECT COUNT(*) FROM auth_accounts);
    RAISE NOTICE 'Persons insertadas: %', (SELECT COUNT(*) FROM person_person);
    RAISE NOTICE '  - Individual: %', (SELECT COUNT(*) FROM person_individual);
    RAISE NOTICE '  - Company: %', (SELECT COUNT(*) FROM person_company);
    RAISE NOTICE 'Users insertados: %', (SELECT COUNT(*) FROM auth_users);
    RAISE NOTICE 'Contactos insertados: %', (SELECT COUNT(*) FROM person_contacto);
    RAISE NOTICE '═══════════════════════════════════════════════════════════════════';
    RAISE NOTICE 'USUARIOS DE PRUEBA DISPONIBLES:';
    RAISE NOTICE '• juan.perez@example.com (password: password)';
    RAISE NOTICE '• maria.gonzalez@example.com (password: password)';
    RAISE NOTICE '• carlos.rodriguez@example.com (password: password) - PENDING';
    RAISE NOTICE '• ana.martinez@example.com (password: password)';
    RAISE NOTICE '• laura.fernandez@gmail.com (Google OAuth)';
    RAISE NOTICE '• admin@acmecorp.com (password: password) - EMPRESA';
    RAISE NOTICE '• contacto@techsolutions.com (password: password) - EMPRESA';
    RAISE NOTICE '• ventas@innovatech.com.ar (password: password) - EMPRESA';
    RAISE NOTICE '═══════════════════════════════════════════════════════════════════';
END
$$;
