package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/rem-gestion/rem-common/config"
	"go.uber.org/zap"
)

type AdaptiveMigrator struct {
	config    *config.Config
	db        *sql.DB
	logger    *zap.Logger
	tableName string
}

func NewAdaptiveMigrator(cfg *config.Config, db *sql.DB, logger *zap.Logger) *AdaptiveMigrator {
	// Crear nombre de tabla de migraciones adaptativo
	dbResolver := config.NewDatabaseConfigResolver(cfg.ServiceName)
	tableName := dbResolver.GetMigrationTableName()

	return &AdaptiveMigrator{
		config:    cfg,
		db:        db,
		logger:    logger.Named("migrator"),
		tableName: tableName,
	}
}

func (am *AdaptiveMigrator) RunMigrations(migrationPath string) error {
	am.logger.Info("starting migrations",
		zap.String("service", am.config.ServiceName),
		zap.String("environment", am.config.Environment.String()),
		zap.String("table", am.tableName))

	// Crear tabla de migraciones si no existe
	if err := am.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Obtener archivos de migración
	migrationFiles, err := am.getMigrationFiles(migrationPath)
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	if len(migrationFiles) == 0 {
		am.logger.Info("no migrations found", zap.String("path", migrationPath))
		return nil
	}

	// Aplicar migraciones
	applied := 0
	for _, file := range migrationFiles {
		if err := am.applyMigration(file); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", file, err)
		}
		applied++
	}

	am.logger.Info("migrations completed",
		zap.Int("applied", applied),
		zap.String("service", am.config.ServiceName))

	return nil
}

func (am *AdaptiveMigrator) createMigrationsTable() error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`, am.tableName)

	_, err := am.db.Exec(query)
	if err != nil {
		am.logger.Error("failed to create migrations table", zap.Error(err))
	}
	return err
}

func (am *AdaptiveMigrator) getMigrationFiles(migrationPath string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(migrationPath, "*.up.sql"))
	if err != nil {
		return nil, err
	}

	// Ordenar archivos por nombre (que debería incluir timestamp)
	sort.Strings(files)

	// Filtrar archivos ya aplicados
	var pendingFiles []string
	for _, file := range files {
		version := am.getVersionFromFilename(file)
		if !am.isMigrationApplied(version) {
			pendingFiles = append(pendingFiles, file)
		}
	}

	return pendingFiles, nil
}

func (am *AdaptiveMigrator) applyMigration(filePath string) error {
	version := am.getVersionFromFilename(filePath)

	am.logger.Info("applying migration",
		zap.String("version", version),
		zap.String("file", filepath.Base(filePath)))

	// Leer contenido del archivo
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Aplicar prefijos si está en desarrollo
	sql := string(content)
	if am.config.EnvironmentDetector.IsDevelopment() {
		sql = am.applyTablePrefixes(sql)
		am.logger.Debug("applied table prefixes for development environment")
	}

	// Ejecutar migración
	if _, err := am.db.Exec(sql); err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	// Marcar como aplicada
	if err := am.markMigrationAsApplied(version); err != nil {
		return fmt.Errorf("failed to mark migration as applied: %w", err)
	}

	am.logger.Info("migration applied successfully", zap.String("version", version))
	return nil
}

func (am *AdaptiveMigrator) applyTablePrefixes(sql string) string {
	if !am.config.EnvironmentDetector.IsDevelopment() {
		return sql
	}

	prefix := fmt.Sprintf("%s_", am.config.ServiceName)
	serviceTables := am.getServiceTables()

	am.logger.Debug("applying table prefixes",
		zap.String("prefix", prefix),
		zap.Strings("tables", serviceTables))

	// Aplicar prefijos usando regex más precisos
	result := sql
	for _, table := range serviceTables {
		// Usar expresiones regulares para ser más precisos
		oldPatterns := []string{
			// CREATE TABLE variantes
			fmt.Sprintf(`(?i)\bCREATE\s+TABLE\s+%s\b`, table),
			fmt.Sprintf(`(?i)\bCREATE\s+TABLE\s+IF\s+NOT\s+EXISTS\s+%s\b`, table),
			// CREATE INDEX ON table
			fmt.Sprintf(`(?i)\bON\s+%s\s*\(`, table),
			fmt.Sprintf(`(?i)\bON\s+%s\s*;`, table),
			fmt.Sprintf(`(?i)\bON\s+%s\s*$`, table),
			// INSERT, UPDATE, DELETE, ALTER
			fmt.Sprintf(`(?i)\bINSERT\s+INTO\s+%s\b`, table),
			fmt.Sprintf(`(?i)\bUPDATE\s+%s\s+`, table),
			fmt.Sprintf(`(?i)\bDELETE\s+FROM\s+%s\b`, table),
			fmt.Sprintf(`(?i)\bALTER\s+TABLE\s+%s\b`, table),
			fmt.Sprintf(`(?i)\bDROP\s+TABLE\s+%s\b`, table),
			// REFERENCES para foreign keys
			fmt.Sprintf(`(?i)\bREFERENCES\s+%s\b`, table),
		}

		newTable := prefix + table

		for _, pattern := range oldPatterns {
			result = am.regexReplace(result, pattern, table, newTable)
		}
	}

	am.logger.Debug("prefixes applied", zap.String("result_preview", result[:min(400, len(result))]))
	return result
}

func (am *AdaptiveMigrator) getServiceTables() []string {
	// Definir tablas por servicio
	tables := map[string][]string{
		"person":  {"person", "individual", "company", "contacto"},
		"address": {"addresses"},
		"auth":    {"accounts", "users"},
	}

	if serviceTables, exists := tables[am.config.ServiceName]; exists {
		return serviceTables
	}

	return []string{}
}

func (am *AdaptiveMigrator) getVersionFromFilename(filePath string) string {
	filename := filepath.Base(filePath)
	return strings.TrimSuffix(filename, ".up.sql")
}

func (am *AdaptiveMigrator) isMigrationApplied(version string) bool {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE version = $1", am.tableName)
	var count int
	err := am.db.QueryRow(query, version).Scan(&count)
	if err != nil {
		am.logger.Warn("failed to check migration status", zap.String("version", version), zap.Error(err))
		return false
	}
	return count > 0
}

func (am *AdaptiveMigrator) markMigrationAsApplied(version string) error {
	query := fmt.Sprintf("INSERT INTO %s (version, applied_at) VALUES ($1, $2)", am.tableName)
	_, err := am.db.Exec(query, version, time.Now())
	return err
}

func (am *AdaptiveMigrator) regexReplace(text, pattern, oldTable, newTable string) string {
	re, err := regexp.Compile(pattern)
	if err != nil {
		am.logger.Warn("invalid regex pattern", zap.String("pattern", pattern), zap.Error(err))
		return text
	}

	// Reemplazar todas las ocurrencias, pero solo el nombre de la tabla específica
	return re.ReplaceAllStringFunc(text, func(match string) string {
		// Reemplazar solo el nombre de la tabla en el match, manteniendo el resto
		return strings.ReplaceAll(match, oldTable, newTable)
	})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
