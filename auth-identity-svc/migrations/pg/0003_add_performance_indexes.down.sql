-- Eliminar índices adicionales
DROP INDEX IF EXISTS idx_accounts_email_status;
DROP INDEX IF EXISTS idx_users_account_person;
DROP INDEX IF EXISTS idx_accounts_created_at;
DROP INDEX IF EXISTS idx_users_created_at;
