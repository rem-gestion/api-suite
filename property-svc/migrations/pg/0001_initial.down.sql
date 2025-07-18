-- Eliminar tablas en orden inverso (por las FK)
DROP TABLE IF EXISTS property_amenities;
DROP TABLE IF EXISTS amenity;
DROP TABLE IF EXISTS property_management;
DROP TABLE IF EXISTS property;

-- Eliminar tipos personalizados
DROP TYPE IF EXISTS property_type;
