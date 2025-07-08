-- Crear extensión para UUID si no existe
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Crear tabla accounts
CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    provider VARCHAR(16) NOT NULL CHECK (provider IN ('email', 'google')),
    email VARCHAR(160) UNIQUE NOT NULL,
    password_hash VARCHAR(255),
    status VARCHAR(16) DEFAULT 'pending' CHECK (status IN ('pending', 'active', 'suspended')),
    
    -- Auditoría
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    updated_by UUID,
    deleted_at TIMESTAMP
);

-- Índices para accounts
CREATE INDEX idx_accounts_email ON accounts(email);
CREATE INDEX idx_accounts_status ON accounts(status);
CREATE INDEX idx_accounts_deleted_at ON accounts(deleted_at);
CREATE INDEX idx_accounts_provider ON accounts(provider);
