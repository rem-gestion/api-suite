package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/rem-gestion/rem-common/config"
	"github.com/rem-gestion/rem-common/db"
	"github.com/rem-gestion/rem-common/logger"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	// Carga las REM_* definidas en .env
	cfg := config.Load()
	lg := logger.New(cfg.Logger, "migrate")

	lg.Info("starting address service migrations",
		zap.String("service", cfg.ServiceName),
		zap.String("environment", cfg.Environment.String()),
		zap.String("database", cfg.Postgres.Database))

	// up | down | "steps -1"  (por defecto: up)
	dir := "up"
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	// Construir string de conexión
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password,
		cfg.Postgres.Database, cfg.Postgres.SSLMode)

	// Conectar a la base de datos
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	// Verificar conexión
	if err := dbConn.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Usar migrador adaptativo si está en desarrollo
	if cfg.EnvironmentDetector.IsDevelopment() {
		lg.Info("using adaptive migrator for development environment")
		migrator := db.NewAdaptiveMigrator(&cfg, dbConn, lg)
		if err := migrator.RunMigrations("./migrations/pg"); err != nil {
			log.Fatalf("Adaptive migration failed: %v", err)
		}
	} else {
		// Usar migrador tradicional para producción
		lg.Info("using traditional migrator for production environment")
		path := "./migrations/pg"
		if err := db.RunMigrations(cfg.DriverRelacional, cfg, path, dir); err != nil {
			log.Fatalf("Traditional migration failed: %v", err)
		}
	}

	lg.Info("migrations completed successfully")
}
