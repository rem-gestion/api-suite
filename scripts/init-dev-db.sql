-- Script de inicialización para base de datos de desarrollo
-- Este script se ejecuta automáticamente cuando se crea el contenedor

-- Crear extensiones necesarias
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- El usuario principal 'user' ya existe por las variables de entorno de postgres
-- Solo necesitamos otorgar permisos adicionales y configurar la base de datos

-- Otorgar permisos necesarios al usuario principal
GRANT ALL PRIVILEGES ON DATABASE rem_development TO "user";
GRANT ALL PRIVILEGES ON SCHEMA public TO "user";

-- Configurar permisos por defecto para objetos futuros
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO "user";
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO "user";
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON FUNCTIONS TO "user";

-- Log de inicialización
DO $$
BEGIN
    RAISE NOTICE 'REM Development Database initialized successfully';
    RAISE NOTICE 'Database: rem_development';
    RAISE NOTICE 'User: user (with full privileges)';
    RAISE NOTICE 'Extensions: uuid-ossp';
END
$$;
