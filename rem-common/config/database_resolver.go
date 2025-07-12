package config

import (
	"fmt"
	"os"

	"github.com/rem-gestion/rem-common/environment"
)

type DatabaseConfigResolver struct {
	envDetector *environment.EnvironmentDetector
	serviceName string
}

func NewDatabaseConfigResolver(serviceName string) *DatabaseConfigResolver {
	return &DatabaseConfigResolver{
		envDetector: environment.NewEnvironmentDetector(),
		serviceName: serviceName,
	}
}

func (dcr *DatabaseConfigResolver) ResolvePostgresConfig() PostgresConfig {
	baseConfig := PostgresConfig{
		Host:     getEnvOrDefault("REM_POSTGRES_HOST", "localhost"),
		Port:     getEnvOrDefault("REM_POSTGRES_PORT", dcr.getDefaultPort()),
		User:     getEnvOrDefault("REM_POSTGRES_USER", "user"),
		Password: getEnvOrDefault("REM_POSTGRES_PASSWORD", "supersecreta"),
		SSLMode:  getEnvOrDefault("REM_POSTGRES_SSLMODE", "disable"),
	}

	// Resolver nombre de base de datos según entorno
	baseConfig.Database = dcr.resolveDatabaseName()

	return baseConfig
}

func (dcr *DatabaseConfigResolver) resolveDatabaseName() string {
	// Si está explícitamente configurado, usarlo
	if dbName := os.Getenv("REM_POSTGRES_DB"); dbName != "" {
		return dbName
	}

	// Resolución automática según entorno
	switch dcr.envDetector.GetEnvironment() {
	case environment.Development:
		return "rem_development" // Base de datos compartida
	case environment.Testing:
		return "rem_test"
	case environment.Production:
		return dcr.getProductionDatabaseName()
	default:
		return dcr.getProductionDatabaseName()
	}
}

func (dcr *DatabaseConfigResolver) getProductionDatabaseName() string {
	// Mapear nombres de servicio a nombres de DB existentes
	dbNames := map[string]string{
		"address": "remgestion_address",
		"person":  "remgestion_person",
		"auth":    "remgestion_auth",
	}

	if dbName, exists := dbNames[dcr.serviceName]; exists {
		return dbName
	}

	return fmt.Sprintf("remgestion_%s", dcr.serviceName)
}

func (dcr *DatabaseConfigResolver) getDefaultPort() string {
	// En desarrollo, todos usan el mismo puerto
	if dcr.envDetector.IsDevelopment() {
		return "5432"
	}

	// En producción, usar puertos específicos como están configurados
	ports := map[string]string{
		"address": "5432",
		"person":  "5433",
		"auth":    "5434",
	}

	if port, exists := ports[dcr.serviceName]; exists {
		return port
	}

	return "5432"
}

func (dcr *DatabaseConfigResolver) GetTablePrefix() string {
	if dcr.envDetector.IsDevelopment() {
		return fmt.Sprintf("%s_", dcr.serviceName)
	}
	return "" // Sin prefijo en producción
}

func (dcr *DatabaseConfigResolver) GetMigrationTableName() string {
	if dcr.envDetector.IsDevelopment() {
		return fmt.Sprintf("%s_schema_migrations", dcr.serviceName)
	}
	return "schema_migrations"
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
