-- Crear tabla users
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID UNIQUE NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    person_id UUID, -- Referencia externa a person-svc
    onboard_status VARCHAR(16) DEFAULT 'new' CHECK (onboard_status IN ('new', 'in_progress', 'done')),
    last_login TIMESTAMP,
    
    -- Auditoría
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    updated_by UUID,
    deleted_at TIMESTAMP
);

-- Índices para users
CREATE UNIQUE INDEX idx_users_account_id ON users(account_id);
CREATE INDEX idx_users_person_id ON users(person_id);
CREATE INDEX idx_users_onboard_status ON users(onboard_status);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
CREATE INDEX idx_users_last_login ON users(last_login);
