CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- enum sexo
DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'sexo') THEN
    CREATE TYPE sexo AS ENUM ('masculino', 'femenino');
  END IF;
END$$;

-- tabla person  ────────────────
CREATE TABLE person (
  id          uuid         PRIMARY KEY DEFAULT uuid_generate_v4(),
  type        varchar(10)  NOT NULL  CHECK (type IN ('individual','company')),
  address_id  uuid,                    -- FK lógica (otro micro) → *sin* constraint hard
  avatar_url  text,
  sexo        sexo,

  created_at  timestamp    NOT NULL DEFAULT now(),
  updated_at  timestamp,
  updated_by  uuid,
  deleted_at  timestamp
);

-- subtipo individual ───────────
CREATE TABLE individual (
  person_id   uuid PRIMARY KEY REFERENCES person(id) ON DELETE CASCADE,
  first_name  varchar(64) NOT NULL,
  last_name   varchar(64) NOT NULL,
  dni         char(8) UNIQUE,

  created_at  timestamp NOT NULL DEFAULT now(),
  updated_at  timestamp,
  updated_by  uuid,
  deleted_at  timestamp
);

-- subtipo company ──────────────
CREATE TABLE company (
  person_id    uuid PRIMARY KEY REFERENCES person(id) ON DELETE CASCADE,
  legal_name   varchar(120) NOT NULL,
  cuit         char(11) UNIQUE,
  society_type varchar(12),

  created_at   timestamp NOT NULL DEFAULT now(),
  updated_at   timestamp,
  updated_by   uuid,
  deleted_at   timestamp
);

-- contactos ─────────────────────
CREATE TABLE contacto (
  id          uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  persona_id  uuid NOT NULL REFERENCES person(id) ON DELETE CASCADE,
  tipo        varchar(16) NOT NULL CHECK (tipo IN ('email','phone','whatsapp')),
  dato        varchar(128) NOT NULL,
  is_primary  boolean DEFAULT false,

  created_at  timestamp NOT NULL DEFAULT now(),
  created_by  uuid,
  updated_at  timestamp,
  updated_by  uuid,
  deleted_at  timestamp
);

CREATE INDEX idx_persona      ON contacto (persona_id);
CREATE INDEX idx_persona_tipo ON contacto (persona_id, tipo);
