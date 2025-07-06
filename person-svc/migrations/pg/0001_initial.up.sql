CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE sexo AS ENUM ('masculino', 'femenino');

CREATE TABLE person (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  type varchar(10) NOT NULL,              -- individual | company
  address_id uuid REFERENCES address(id),
  avatar_url text,
  sexo sexo,
  created_at timestamp DEFAULT now() NOT NULL,
  updated_at timestamp,
  updated_by uuid,
  deleted_at timestamp
);

CREATE TABLE individual (
  person_id uuid PRIMARY KEY REFERENCES person(id),
  first_name varchar(64) NOT NULL,
  last_name varchar(64)  NOT NULL,
  dni char(8) UNIQUE,
  created_at timestamp DEFAULT now() NOT NULL,
  updated_at timestamp,
  updated_by uuid,
  deleted_at timestamp
);

CREATE TABLE company (
  person_id uuid PRIMARY KEY REFERENCES person(id),
  legal_name varchar(120) NOT NULL,
  cuit char(11) UNIQUE,
  society_type varchar(12),
  created_at timestamp DEFAULT now() NOT NULL,
  updated_at timestamp,
  updated_by uuid,
  deleted_at timestamp
);

CREATE TABLE contacto (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  persona_id uuid NOT NULL REFERENCES person(id),
  tipo varchar(16) NOT NULL,
  dato varchar(128) NOT NULL,
  is_primary boolean DEFAULT false,
  created_at timestamp DEFAULT now() NOT NULL,
  created_by uuid,
  updated_at timestamp,
  updated_by uuid,
  deleted_at timestamp
);

CREATE INDEX idx_persona      ON contacto (persona_id);
CREATE INDEX idx_persona_tipo ON contacto (persona_id, tipo);
