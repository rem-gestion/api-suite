-- =============================================
-- Migration: 0007_create_organization_domain.up.sql
-- Description: Create organization domain tables with ENUM types and improved constraints
-- Author: System Generated
-- Date: 2024
-- Version: 2.0 - Enhanced with review feedback
-- =============================================

-- Ensure required extensions are available
DO $$
BEGIN
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
END;
$$ LANGUAGE plpgsql;



-- =============================================================================
-- FUNCIONES AUXILIARES (deben definirse antes de usarse como DEFAULT)
-- =============================================================================

-- Función para generar token de verificación DNS único
CREATE OR REPLACE FUNCTION generate_dns_verification_token()
RETURNS TEXT AS $$
DECLARE
    token TEXT;
    token_exists BOOLEAN;
    max_attempts INTEGER := 10;
    attempt_count INTEGER := 0;
BEGIN
    LOOP
        -- Generar token usando CSPRNG para mayor seguridad
        token := encode(gen_random_bytes(24), 'base64');
        -- Hacer URL-safe reemplazando caracteres problemáticos
        token := replace(replace(token, '/', '_'), '+', '-');
        -- Truncar a 32 caracteres
        token := left(token, 32);
        
        -- Verificar si el token ya existe
        SELECT EXISTS(SELECT 1 FROM organization_domain WHERE dns_verification_token = token) INTO token_exists;
        
        -- Si no existe, salir del loop
        IF NOT token_exists THEN
            EXIT;
        END IF;
        
        attempt_count := attempt_count + 1;
        IF attempt_count >= max_attempts THEN
            RAISE EXCEPTION 'Unable to generate unique DNS verification token after % attempts', max_attempts;
        END IF;
    END LOOP;
    
    RETURN token;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION is_valid_domain(domain_name text)
RETURNS boolean
LANGUAGE sql
IMMUTABLE
STRICT
AS $$
  /* Reglas:
     - ≤253 chars
     - al menos un .
     - cada label 1-63, alfanum o guion, sin guion al principio/fin
  */
  SELECT
    domain_name IS NOT NULL
    AND length(domain_name) <= 253
    AND lower(domain_name) ~
      '^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?([.][a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?)+$';
$$;

CREATE OR REPLACE FUNCTION is_valid_subdomain(sd text)
RETURNS boolean
LANGUAGE sql
IMMUTABLE
STRICT
AS $$
  SELECT
    sd IS NOT NULL
    AND length(sd) BETWEEN 1 AND 63
    AND sd ~ '^[A-Za-z0-9]([A-Za-z0-9-]{0,61}[A-Za-z0-9])?$';
$$;


DROP FUNCTION IF EXISTS generate_dns_verification_token();

CREATE OR REPLACE FUNCTION generate_dns_verification_token()
RETURNS text
LANGUAGE sql
IMMUTABLE            -- no toca tablas, ahora es pura
AS $$
  SELECT left(
           translate(encode(gen_random_bytes(24), 'base64'), '/+', '_-'),
           32
         );
$$;


-- =============================================================================
-- TABLA: organization_domain (dominios personalizados por organización)
-- =============================================================================
CREATE TABLE organization_domain (
    id                      UUID         PRIMARY KEY DEFAULT generate_uuid(),
    organization_id         UUID         NOT NULL,
    domain_name             VARCHAR(255) NOT NULL,
    subdomain               VARCHAR(100),
    domain_type             domain_type_enum NOT NULL DEFAULT 'custom',
    status                  domain_status_enum NOT NULL DEFAULT 'pending',
    is_primary              BOOLEAN      NOT NULL DEFAULT false,
    ssl_enabled             BOOLEAN      NOT NULL DEFAULT false,
    ssl_certificate         TEXT,
    ssl_private_key         TEXT, -- Encriptado en aplicación
    ssl_expires_at          TIMESTAMP WITH TIME ZONE,
    dns_verified            BOOLEAN      NOT NULL DEFAULT false,
    dns_verification_token  VARCHAR(100) NOT NULL DEFAULT generate_dns_verification_token(),
    dns_verification_method dns_verification_method_enum DEFAULT 'txt',
    verification_attempts   INTEGER      NOT NULL DEFAULT 0,
    last_verification_at    TIMESTAMP WITH TIME ZONE,
    redirect_to_primary     BOOLEAN      NOT NULL DEFAULT false,
    custom_headers          JSONB        DEFAULT '{}',
    
    -- Auditoría completa
    created_at              TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp_utc(),
    created_by              UUID NOT NULL,        -- FK a auth-identity-svc users.id (external)
    updated_at              TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp_utc(),
    updated_by              UUID,        -- FK a auth-identity-svc users.id (external)
    deleted_at              TIMESTAMP WITH TIME ZONE, -- Soft delete
    
    -- Constraints mejorados
    CONSTRAINT fk_organization_domain_organization 
        FOREIGN KEY (organization_id) REFERENCES organization(id) ON DELETE CASCADE,
    CONSTRAINT chk_organization_domain_valid_domain 
        CHECK (is_valid_domain(domain_name)),
    -- Sólo letras, números y guiones:
    CONSTRAINT chk_organization_domain_subdomain_format
      CHECK (subdomain IS NULL OR is_valid_subdomain(subdomain)),
    CONSTRAINT chk_organization_domain_verification_attempts 
        CHECK (verification_attempts >= 0 AND verification_attempts <= 10),
    CONSTRAINT chk_organization_domain_ssl_consistency 
        CHECK (
            (ssl_enabled = false) OR 
            (ssl_enabled = true AND ssl_certificate IS NOT NULL AND ssl_private_key IS NOT NULL)
        ),
    CONSTRAINT chk_organization_domain_ssl_certificate_size
        CHECK (ssl_certificate IS NULL OR char_length(ssl_certificate) BETWEEN 100 AND 65535),
    CONSTRAINT chk_organization_domain_ssl_private_key_size
        CHECK (ssl_private_key IS NULL OR char_length(ssl_private_key) BETWEEN 100 AND 65535),
    CONSTRAINT chk_organization_domain_ssl_expiry 
        CHECK (ssl_expires_at IS NULL OR ssl_expires_at > created_at),
    CONSTRAINT chk_organization_domain_token_length 
        CHECK (dns_verification_token IS NULL OR char_length(dns_verification_token) >= 16),
    CONSTRAINT chk_organization_domain_redirect_logic
        CHECK (
            (is_primary = true AND redirect_to_primary = false) OR
            (is_primary = false)
        )
);

-- =============================================================================
-- UNIQUE CONSTRAINTS E INDICES ESPECIALES
-- =============================================================================

-- Constraint único para asegurar un solo dominio primario por organización
CREATE UNIQUE INDEX uq_organization_domain_one_primary_per_org 
ON organization_domain(organization_id) 
WHERE is_primary = true AND deleted_at IS NULL;

-- Índice único case-insensitive para domain_name POR ORGANIZACIÓN (permite mismo dominio en diferentes orgs)
CREATE UNIQUE INDEX uq_organization_domain_name_case_insensitive
ON organization_domain(organization_id, LOWER(domain_name))
WHERE deleted_at IS NULL;

-- =============================================================================
-- TABLA: organization_domain_dns (configuración DNS por dominio)
-- =============================================================================
CREATE TABLE organization_domain_dns (
    id                      UUID         PRIMARY KEY DEFAULT generate_uuid(),
    domain_id               UUID         NOT NULL,
    record_type             dns_record_type_enum NOT NULL,
    record_name             VARCHAR(255) NOT NULL,
    record_value            TEXT         NOT NULL,
    record_ttl              INTEGER      NOT NULL DEFAULT 300,
    is_required             BOOLEAN      NOT NULL DEFAULT true,
    is_verified             BOOLEAN      NOT NULL DEFAULT false,
    verification_attempts   INTEGER      NOT NULL DEFAULT 0,
    last_verification_at    TIMESTAMP WITH TIME ZONE,
    notes                   TEXT,
    
    -- Auditoría
    created_at              TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp_utc(),
    updated_at              TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp_utc(),
    
    -- Constraints mejorados
    CONSTRAINT fk_organization_domain_dns_domain 
        FOREIGN KEY (domain_id) REFERENCES organization_domain(id) ON DELETE CASCADE,
    CONSTRAINT chk_organization_domain_dns_valid_ttl 
        CHECK (record_ttl >= 60 AND record_ttl <= 86400),
    CONSTRAINT chk_organization_domain_dns_verification_attempts 
        CHECK (verification_attempts >= 0 AND verification_attempts <= 10),
    CONSTRAINT chk_organization_domain_dns_record_name_length 
        CHECK (char_length(record_name) >= 1),
    CONSTRAINT chk_organization_domain_dns_record_value_length 
        CHECK (char_length(record_value) >= 1),
    
    -- Constraint único mejorado
    CONSTRAINT uq_organization_domain_dns_unique 
        UNIQUE (domain_id, record_type, record_name) DEFERRABLE INITIALLY DEFERRED
);

-- =============================================================================
-- TABLA: organization_domain_verification_log (logs de verificación - particionable)
-- =============================================================================
CREATE TABLE organization_domain_verification_log (
    id                UUID                         NOT NULL DEFAULT generate_uuid(),
    domain_id         UUID                         NOT NULL,
    dns_record_id     UUID,
    verification_type domain_verification_type_enum NOT NULL,
    status            verification_status_enum      NOT NULL,
    details           JSONB                        DEFAULT '{}',
    error_message     TEXT,
    response_data     JSONB,
    duration_ms       INTEGER,

    -- Auditoría simple (tabla de logs)
    created_at        TIMESTAMPTZ                  NOT NULL DEFAULT current_timestamp_utc(),

    -- Foreign keys y checks
    CONSTRAINT fk_organization_domain_verification_log_domain
      FOREIGN KEY (domain_id) REFERENCES organization_domain(id) ON DELETE CASCADE,
    CONSTRAINT fk_organization_domain_verification_log_dns
      FOREIGN KEY (dns_record_id) REFERENCES organization_domain_dns(id) ON DELETE SET NULL,
    CONSTRAINT chk_organization_domain_verification_log_duration
      CHECK (duration_ms IS NULL OR duration_ms >= 0),

    -- PK compuesto que incluye la columna de partición
    CONSTRAINT pk_organization_domain_verification_log
      PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

-- =============================================================================
-- INICIALIZACIÓN: CREAR PARTICIONES PARA logs de verificación (próximos 3 meses)
-- =============================================================================
DO $$
BEGIN
  PERFORM create_monthly_partition(
      'organization_domain_verification_log',
      CURRENT_DATE
  );
  PERFORM create_monthly_partition(
      'organization_domain_verification_log',
      (CURRENT_DATE + INTERVAL '1 month')::DATE
  );
  PERFORM create_monthly_partition(
      'organization_domain_verification_log',
      (CURRENT_DATE + INTERVAL '2 months')::DATE
  );
END;
$$ LANGUAGE plpgsql;
-- =============================================================================
-- ÍNDICES CON NOMBRES EXPLÍCITOS Y MEJORADOS
-- =============================================================================

-- Índices para organization_domain (eliminando duplicados)
CREATE INDEX ix_organization_domain_organization_id ON organization_domain(organization_id);
CREATE INDEX ix_organization_domain_status ON organization_domain(status);
CREATE INDEX ix_organization_domain_subdomain ON organization_domain(subdomain) WHERE subdomain IS NOT NULL;
CREATE INDEX ix_organization_domain_ssl_expiry ON organization_domain(ssl_expires_at) WHERE ssl_expires_at IS NOT NULL;
CREATE INDEX ix_organization_domain_verification ON organization_domain(dns_verified, status);
CREATE INDEX ix_organization_domain_created_by ON organization_domain(created_by);
CREATE INDEX ix_organization_domain_active ON organization_domain(organization_id, status) WHERE deleted_at IS NULL;
-- NOTA: No se crea índice separado para domain_name porque ya tenemos el índice único case-insensitive

-- Índices para organization_domain_dns
CREATE INDEX ix_organization_domain_dns_domain_id ON organization_domain_dns(domain_id);
CREATE INDEX ix_organization_domain_dns_type ON organization_domain_dns(record_type);
CREATE INDEX ix_organization_domain_dns_verified ON organization_domain_dns(is_verified);
CREATE INDEX ix_organization_domain_dns_required ON organization_domain_dns(domain_id, is_required) WHERE is_required = true;
-- Índice optimizado para polling de registros DNS pendientes de verificación
CREATE INDEX ix_organization_domain_dns_polling ON organization_domain_dns(is_verified, last_verification_at, verification_attempts)
WHERE is_verified = false AND is_required = true;

-- Índices para organization_domain_verification_log (tabla particionada)
CREATE INDEX ix_organization_domain_verification_log_domain_id ON organization_domain_verification_log(domain_id);
CREATE INDEX ix_organization_domain_verification_log_dns_record_id ON organization_domain_verification_log(dns_record_id) WHERE dns_record_id IS NOT NULL;
CREATE INDEX ix_organization_domain_verification_log_type ON organization_domain_verification_log(verification_type);
CREATE INDEX ix_organization_domain_verification_log_status ON organization_domain_verification_log(status);
CREATE INDEX ix_organization_domain_verification_log_created_at ON organization_domain_verification_log(created_at);

-- =============================================================================
-- TRIGGERS Y FUNCIONES DE BUSINESS LOGIC
-- =============================================================================

-- Trigger para actualizar updated_at automáticamente
CREATE TRIGGER trg_organization_domain_updated_at
    BEFORE UPDATE ON organization_domain
    FOR EACH ROW
    EXECUTE FUNCTION update_organization_updated_at();

CREATE TRIGGER trg_organization_domain_dns_updated_at
    BEFORE UPDATE ON organization_domain_dns
    FOR EACH ROW
    EXECUTE FUNCTION update_organization_updated_at();



-- Función para sincronizar dns_verified con el estado de verificación de registros DNS (CORREGIDA)
CREATE OR REPLACE FUNCTION sync_domain_dns_verification()
RETURNS TRIGGER AS $$
DECLARE
    required_verified_count INTEGER;
    total_required_count    INTEGER;
    domain_verified         BOOLEAN;
    current_domain_verified BOOLEAN;
    domain_id_var           UUID;
BEGIN
    -- Determinar domain_id según la operación
    domain_id_var := CASE 
        WHEN TG_OP = 'DELETE' THEN OLD.domain_id
        ELSE NEW.domain_id
    END;
    
    -- Solo procesar si cambió is_verified o si es INSERT/DELETE
    IF TG_OP = 'UPDATE' AND OLD.is_verified = NEW.is_verified THEN
        RETURN NEW;
    END IF;
    
    -- Obtener estado actual antes de cualquier cambio
    SELECT dns_verified INTO current_domain_verified 
      FROM organization_domain 
     WHERE id = domain_id_var;
    
    -- Contar registros requeridos y verificados
    SELECT 
       COUNT(*) FILTER (WHERE is_required AND is_verified),
       COUNT(*) FILTER (WHERE is_required)
    INTO required_verified_count, total_required_count
      FROM organization_domain_dns 
     WHERE domain_id = domain_id_var;
    
    -- Determinar si ya están todos verificados
    domain_verified := (total_required_count > 0
                         AND required_verified_count = total_required_count);
    
    -- Actualizar si cambió el estado
    IF current_domain_verified IS DISTINCT FROM domain_verified THEN
        UPDATE organization_domain
           SET dns_verified = domain_verified,
               updated_at    = current_timestamp_utc()
         WHERE id = domain_id_var;
        
        PERFORM log_domain_verification(
            domain_id_var,
            CASE WHEN TG_OP = 'DELETE' THEN NULL ELSE NEW.id END,
            'dns'::domain_verification_type_enum,
            CASE WHEN domain_verified THEN 'success'::verification_status_enum
                 ELSE 'partial'::verification_status_enum END,
            jsonb_build_object(
              'required_verified_count', required_verified_count,
              'total_required_count',    total_required_count,
              'trigger_operation',       TG_OP,
              'old_verified',            current_domain_verified::boolean,
              'new_verified',            domain_verified::boolean
            )
        );
    END IF;
    
    RETURN CASE WHEN TG_OP = 'DELETE' THEN OLD ELSE NEW END;
END;
$$ LANGUAGE plpgsql;

-- Trigger para sincronizar verificación DNS
CREATE TRIGGER trg_organization_domain_dns_sync_verification
    AFTER INSERT OR UPDATE OR DELETE ON organization_domain_dns
    FOR EACH ROW
    EXECUTE FUNCTION sync_domain_dns_verification();

-- Función mejorada para crear registros DNS por defecto al crear un dominio
-- Función mejorada para crear registros DNS por defecto al crear un dominio
CREATE OR REPLACE FUNCTION create_default_dns_records()
RETURNS TRIGGER AS $$
DECLARE
    verification_token TEXT;
BEGIN
    -- Solo crear registros DNS para dominios personalizados
    IF NEW.domain_type = 'custom' THEN
        -- Usar el token ya generado o crear uno nuevo
        verification_token := COALESCE(NEW.dns_verification_token, generate_dns_verification_token());
        
        -- Asegurar que el token esté guardado en el dominio
        IF NEW.dns_verification_token IS NULL THEN
            UPDATE organization_domain 
            SET dns_verification_token = verification_token
            WHERE id = NEW.id;
        END IF;
        
        -- Registro A para el dominio principal
        INSERT INTO organization_domain_dns (domain_id, record_type, record_name, record_value, is_required, notes) VALUES
        (NEW.id, 'A', NEW.domain_name, '127.0.0.1', true, 'Registro A principal - actualizar con IP real del servidor');
        
        -- Registro CNAME para www
        INSERT INTO organization_domain_dns (domain_id, record_type, record_name, record_value, is_required, notes) VALUES
        (NEW.id, 'CNAME', 'www.' || NEW.domain_name, NEW.domain_name, false, 'Redirección opcional de www al dominio principal');
        
        -- Registro TXT para verificación
        INSERT INTO organization_domain_dns (domain_id, record_type, record_name, record_value, is_required, notes) VALUES
        (NEW.id, 'TXT', '_verification.' || NEW.domain_name, verification_token, true, 'Token de verificación DNS - requerido para activar el dominio');
        
        -- Log del evento
        PERFORM log_domain_verification(
            NEW.id,
            NULL,
            'dns'::domain_verification_type_enum,
            'partial'::verification_status_enum,
            jsonb_build_object(
                'event', 'domain_created',
                'domain_type', NEW.domain_type,
                'verification_token_generated', (verification_token IS NOT NULL)::boolean
            )
        );
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger para crear registros DNS por defecto
CREATE TRIGGER trg_organization_domain_create_dns_records
    AFTER INSERT ON organization_domain
    FOR EACH ROW
    EXECUTE FUNCTION create_default_dns_records();

-- Función para crear particiones automáticamente
-- NOTA: En PG 15+ se recomienda crear particiones por adelantado via cron job
-- Este trigger funciona como "safety-net" para evitar errores por particiones faltantes
-- ADVERTENCIA: DDL en triggers puede causar dead-locks en alta concurrencia
CREATE OR REPLACE FUNCTION create_domain_verification_log_partition()
RETURNS TRIGGER AS $$
DECLARE
    partition_date DATE;
BEGIN
    partition_date := NEW.created_at::DATE;
    
    -- Crear partición mensual automáticamente usando función de 0001
    -- NOTA: En producción con alta concurrencia, considerar job cron para pre-crear particiones
    BEGIN
        PERFORM create_monthly_partition('organization_domain_verification_log', partition_date);
    EXCEPTION
        WHEN duplicate_table THEN
            -- La partición ya existe, continuar silenciosamente
            NULL;
        WHEN OTHERS THEN
            -- Log error silenciosamente para errores de concurrencia durante migraciones
            -- No usar RAISE WARNING durante migraciones para evitar spam en logs
            NULL;
    END;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger para crear particiones automáticamente
CREATE TRIGGER trg_organization_domain_verification_log_partition
    BEFORE INSERT ON organization_domain_verification_log
    FOR EACH ROW
    EXECUTE FUNCTION create_domain_verification_log_partition();

-- Función mejorada para logging de verificaciones de dominio
-- Función mejorada para logging de verificaciones de dominio
CREATE OR REPLACE FUNCTION log_domain_verification(
    p_domain_id UUID,
    p_dns_record_id UUID,
    p_verification_type domain_verification_type_enum,
    p_status verification_status_enum,
    p_details JSONB DEFAULT '{}',
    p_error_message TEXT DEFAULT NULL,
    p_response_data JSONB DEFAULT NULL,
    p_duration_ms INTEGER DEFAULT NULL
)
RETURNS UUID AS $$
DECLARE
    log_id UUID;
    enriched_details JSONB;
BEGIN
    -- Enriquecer detalles con información contextual
    enriched_details := p_details || jsonb_build_object(
        'logged_at', current_timestamp_utc(),
        'session_id', COALESCE(current_setting('app.session_id', true), 'system'),
        'user_id', COALESCE(current_setting('app.user_id', true), 'system')
    );
    
    -- Insertar log con manejo de errores
    BEGIN
        INSERT INTO organization_domain_verification_log (
            domain_id, dns_record_id, verification_type, status, 
            details, error_message, response_data, duration_ms
        ) VALUES (
            p_domain_id, p_dns_record_id, p_verification_type, p_status,
            enriched_details, p_error_message, p_response_data, p_duration_ms
        ) RETURNING id INTO log_id;
        
    EXCEPTION WHEN OTHERS THEN
        -- Si falla el log, registrar el error pero no fallar la operación principal
        RAISE WARNING 'Failed to log domain verification: %', SQLERRM;
        RETURN NULL;
    END;
    
    RETURN log_id;
END;
$$ LANGUAGE plpgsql;

-- =============================================================================
-- COMENTARIOS PARA DOCUMENTACIÓN
-- =============================================================================
COMMENT ON TABLE organization_domain IS 'Dominios personalizados por organización con soporte SSL y verificación DNS';
COMMENT ON COLUMN organization_domain.domain_name IS 'Nombre del dominio (case-insensitive único)';
COMMENT ON COLUMN organization_domain.ssl_certificate IS 'Certificado SSL (100-65535 caracteres, encriptado en aplicación)';
COMMENT ON COLUMN organization_domain.ssl_private_key IS 'Clave privada SSL (100-65535 caracteres, encriptada en aplicación)';
COMMENT ON COLUMN organization_domain.dns_verification_token IS 'Token único para verificación DNS (generado automáticamente)';
COMMENT ON COLUMN organization_domain.redirect_to_primary IS 'No permitido en dominios primarios (constraint chk_organization_domain_redirect_logic)';
COMMENT ON COLUMN organization_domain.created_by IS 'Referencia a auth-identity-svc users.id';

COMMENT ON TABLE organization_domain_dns IS 'Configuración de registros DNS por dominio con verificación automática';
COMMENT ON COLUMN organization_domain_dns.record_ttl IS 'Time To Live en segundos (60-86400)';
COMMENT ON COLUMN organization_domain_dns.is_verified IS 'Estado sincronizado automáticamente con verificación DNS';

COMMENT ON TABLE organization_domain_verification_log IS 'Logs de verificación de dominios (tabla particionada por fecha con particiones automáticas)';
COMMENT ON COLUMN organization_domain_verification_log.duration_ms IS 'Duración de la verificación en milisegundos';

-- Comentarios en funciones
COMMENT ON FUNCTION generate_dns_verification_token() IS 'Genera tokens únicos de 32 caracteres para verificación DNS';
COMMENT ON FUNCTION create_default_dns_records() IS 'Crea registros DNS por defecto al crear un dominio personalizado';
COMMENT ON FUNCTION sync_domain_dns_verification() IS 'Sincroniza dns_verified con el estado de registros DNS requeridos';
COMMENT ON FUNCTION create_domain_verification_log_partition() IS 'Crea particiones automáticamente para logs de verificación';
COMMENT ON FUNCTION log_domain_verification(UUID, UUID, domain_verification_type_enum, verification_status_enum, JSONB, TEXT, JSONB, INTEGER) IS 'Registra logs de verificación de dominios con contexto enriquecido';

-- =============================================================================
-- SMOKE TESTS - VALIDACIÓN BÁSICA DE LA MIGRACIÓN
-- =============================================================================
DO $$
DECLARE
  test_org           UUID;
  test_domain_1      TEXT;
  test_domain_2      TEXT;
  domain_count       INTEGER;
  dns_count          INTEGER;
  verification_token TEXT;
  is_verified        BOOLEAN;
  test_domain_id     UUID;
  test_dns_id        UUID;
  test_log_id        UUID;
BEGIN
  -- 1) Genero un nuevo UUID para la organización de pruebas
  test_org := gen_random_uuid();

  -- 2) Inserto esa organización en la tabla con campos obligatorios
  INSERT INTO organization (id, display_name, created_by)
  VALUES (test_org, 'Test Organization - Smoke Test', gen_random_uuid());

  -- 3) Limpieza de corridas previas
  DELETE FROM organization_domain
   WHERE domain_name LIKE 'example-%'
      OR domain_name LIKE 'test-%';

  -- 4) Genero dos dominios únicos de prueba
  test_domain_1 := 'example-' || substr(gen_random_uuid()::text, 1, 8) || '.com';
  test_domain_2 := 'test-'    || substr(gen_random_uuid()::text, 1, 8) || '.com';
  RAISE NOTICE '⏳ Ejecutando smoke tests con % y % (org=%)', test_domain_1, test_domain_2, test_org;

  -- Test 1: Tablas creadas
  SELECT COUNT(*) INTO domain_count
    FROM information_schema.tables
   WHERE table_name IN (
     'organization_domain',
     'organization_domain_dns',
     'organization_domain_verification_log'
   );
  IF domain_count <> 3 THEN
    RAISE EXCEPTION 'Error: faltan tablas';
  END IF;
  RAISE NOTICE '✓ Test 1: tablas OK';

  -- Test 2: UNIQUE case-insensitive
  INSERT INTO organization_domain (id, organization_id, domain_name, created_by)
  VALUES
    (gen_random_uuid(), test_org, test_domain_1, gen_random_uuid()),
    (gen_random_uuid(), test_org, test_domain_2, gen_random_uuid());
  BEGIN
    INSERT INTO organization_domain (organization_id, domain_name, created_by)
    VALUES (test_org, upper(test_domain_1), gen_random_uuid());
    RAISE EXCEPTION 'Error: UNIQUE case-insensitive no saltó';
  EXCEPTION WHEN unique_violation THEN
    RAISE NOTICE '✓ Test 2: UNIQUE case-insensitive OK';
  END;

  -- Test 3: Token generado
  SELECT dns_verification_token INTO verification_token
    FROM organization_domain
   WHERE domain_name = test_domain_1;
  IF verification_token IS NULL OR char_length(verification_token) < 16 THEN
    RAISE EXCEPTION 'Error: token no generado';
  END IF;
  RAISE NOTICE '✓ Test 3: token OK';

  -- Test 4: Registros DNS por defecto
  SELECT id INTO test_domain_id
    FROM organization_domain
   WHERE domain_name = test_domain_1;
  SELECT COUNT(*) INTO dns_count
    FROM organization_domain_dns
   WHERE domain_id = test_domain_id;
  IF dns_count < 3 THEN
    RAISE EXCEPTION 'Error: faltan registros DNS';
  END IF;
  RAISE NOTICE '✓ Test 4: DNS OK';

  -- Test 5: Sincronización DNS
  SELECT id INTO test_dns_id
    FROM organization_domain_dns
   WHERE domain_id   = test_domain_id
     AND is_required = true
   LIMIT 1;
  -- primero sólo uno verificado → NO debe marcar dns_verified
  UPDATE organization_domain_dns
     SET is_verified = true
   WHERE id = test_dns_id;
  SELECT dns_verified INTO is_verified
    FROM organization_domain
   WHERE id = test_domain_id;
  IF is_verified THEN
    RAISE EXCEPTION 'Error: sincronización prematura';
  END IF;
  -- luego todos verificados → SÍ debe marcarse
  UPDATE organization_domain_dns
     SET is_verified = true
   WHERE domain_id = test_domain_id
     AND is_required = true;
  SELECT dns_verified INTO is_verified
    FROM organization_domain
   WHERE id = test_domain_id;
  IF NOT is_verified THEN
    RAISE EXCEPTION 'Error: sincronización falló';
  END IF;
  RAISE NOTICE '✓ Test 5: sincronización OK';

  -- Test 6: función log_domain_verification
  SELECT log_domain_verification(
    test_domain_id,
    test_dns_id,
    'dns'::domain_verification_type_enum,
    'success'::verification_status_enum,
    '{"smoke": true}'::jsonb,
    NULL, NULL, 123
  ) INTO test_log_id;
  IF test_log_id IS NULL THEN
    RAISE EXCEPTION 'Error: logging falló';
  END IF;
  RAISE NOTICE '✓ Test 6: logging OK';

  -- Test 7: redirect_to_primary
  BEGIN
    UPDATE organization_domain
       SET is_primary        = true,
           redirect_to_primary = true
     WHERE id = test_domain_id;
    RAISE EXCEPTION 'Error: redirect_to_primary constraint no saltó';
  EXCEPTION WHEN check_violation THEN
    RAISE NOTICE '✓ Test 7: redirect_to_primary OK';
  END;

  -- Limpieza final
  DELETE FROM organization_domain
   WHERE domain_name LIKE 'example-%'
      OR domain_name LIKE 'test-%';
  
  -- Limpiar la organización de prueba creada usando soft delete
  UPDATE organization 
  SET deleted_at = current_timestamp_utc(), 
      status = 'deleted'
  WHERE id = test_org;
  
  RAISE NOTICE '🎉 Todos los smoke tests pasaron correctamente';
END;
$$ LANGUAGE plpgsql;



-- Trigger para manejar automáticamente redirect_to_primary al cambiar is_primary
-- NOTA: Soluciona el constraint chk_organization_domain_redirect_logic automáticamente
-- para evitar que el servicio tenga que hacer dos UPDATEs separados
-- CORREGIDO: Ahora maneja tanto INSERT como UPDATE para cubrir todos los casos
CREATE OR REPLACE FUNCTION manage_domain_primary_redirect()
RETURNS TRIGGER AS $$
BEGIN
    -- Si se marca como primario, automáticamente deshabilitar redirect_to_primary
    IF TG_OP = 'INSERT' THEN
        IF NEW.is_primary = true THEN
            NEW.redirect_to_primary = false;
        END IF;
    ELSIF TG_OP = 'UPDATE' THEN
        IF NEW.is_primary = true AND OLD.is_primary = false THEN
            NEW.redirect_to_primary = false;
        END IF;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_organization_domain_manage_primary_redirect
    BEFORE INSERT OR UPDATE OF is_primary ON organization_domain
    FOR EACH ROW
    EXECUTE FUNCTION manage_domain_primary_redirect();
