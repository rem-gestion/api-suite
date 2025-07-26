-- 0001_initial.down.sql - Property Service Rollback

DROP TYPE IF EXISTS amenity_category CASCADE;
DROP TABLE IF EXISTS property_amenity CASCADE;
DROP TABLE IF EXISTS property_listings CASCADE;
DROP TABLE IF EXISTS property_manager_type CASCADE;
DROP TABLE IF EXISTS property_media CASCADE;
DROP TABLE IF EXISTS property_property CASCADE;
DROP TABLE IF EXISTS property_property_amenities CASCADE;
DROP TABLE IF EXISTS property_property_management CASCADE;
DROP TABLE IF EXISTS property_property_type CASCADE;
DROP TABLE IF EXISTS property_valuations CASCADE;
