-- 0001_initial.down.sql - Property Service Rollback
DROP TABLE IF EXISTS property_amenities CASCADE;
DROP TABLE IF EXISTS property_management CASCADE;
DROP TABLE IF EXISTS property CASCADE;
DROP TABLE IF EXISTS amenity CASCADE;
DROP TABLE IF EXISTS manager_type CASCADE;
DROP TABLE IF EXISTS property_type CASCADE;
DROP TYPE IF EXISTS amenity_category CASCADE;
