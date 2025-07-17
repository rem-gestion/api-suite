CREATE TABLE properties (
    id SERIAL PRIMARY KEY,
    type VARCHAR(32) NOT NULL,
    address_id INTEGER NOT NULL,
    owner_id INTEGER NOT NULL,
    title VARCHAR(128) NOT NULL,
    description TEXT,
    area NUMERIC(10,2),
    rooms INTEGER,
    bathrooms INTEGER,
    amenities TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
