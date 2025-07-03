CREATE TABLE addresses (
  id CHAR(36) PRIMARY KEY,
  floor       VARCHAR(16),
  unit        VARCHAR(16),
  street      VARCHAR(120),
  number      INT,
  city        VARCHAR(64),
  state       VARCHAR(64),
  zip         VARCHAR(16),
  country     CHAR(2),
  created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at  TIMESTAMP NULL,
  updated_by  CHAR(36),
  deleted_at  TIMESTAMP NULL
);
