-- Arreglar constraints para permitir soft delete
-- Eliminar el constraint UNIQUE actual en email
ALTER TABLE accounts DROP CONSTRAINT auth_accounts_email_key;

-- Crear un índice único condicional que solo aplique a registros no eliminados
CREATE UNIQUE INDEX accounts_email_unique_active ON accounts(email) WHERE deleted_at IS NULL;
