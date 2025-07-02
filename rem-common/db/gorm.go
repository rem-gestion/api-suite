package db

import (
	"fmt"
	"time"

	"github.com/rem-gestion/rem-common/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewGormDB abre la conexión a Postgres con GORM
func NewGormDB(cfg config.PostgresConfig) (*gorm.DB, error) {
	// Armamos el DSN con el config
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode,
	)

	// Configuramos el logger de GORM con nivel (podrías mapear el cfg.Logger.Level acá)
	gormLogger := logger.Default.LogMode(logger.Silent)

	// Abrimos la conexión
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("rem-common/db: no pude conectar con GORM: %w", err)
	}

	// Ajustes de pool de conexión:
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("rem-common/db: error obteniendo sql.DB de GORM: %w", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)

	return db, nil
}
