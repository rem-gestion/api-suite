-- Revertir los cambios de constraint para soft delete
DROP INDEX IF EXISTS accounts_email_unique_active;

-- Recrear el constraint UNIQUE original
ALTER TABLE accounts ADD CONSTRAINT auth_accounts_email_key UNIQUE (email);
