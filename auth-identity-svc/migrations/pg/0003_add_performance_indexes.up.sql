-- Agregar índices adicionales para optimización
CREATE INDEX idx_accounts_email_status ON accounts(email, status) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_account_person ON users(account_id, person_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_accounts_created_at ON accounts(created_at);
CREATE INDEX idx_users_created_at ON users(created_at);
