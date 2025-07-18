package main

import (
	"log"

	"github.com/rem-gestion/rem-common/config"
	"github.com/rem-gestion/rem-common/db"
	"github.com/rem-gestion/rem-common/logger"
)

func main() {
	// Carga las REM_* definidas en .env
	cfg := config.Load()
	lg := logger.New(cfg.Logger, "property-migrate")

	// Conectar a base de datos
	pgDB, err := db.NewPostgres(cfg.Postgres)
	if err != nil {
		log.Fatal("postgres connection failed:", err)
	}

	sqlDB, err := pgDB.DB()
	if err != nil {
		log.Fatal("get sql.DB failed:", err)
	}
	defer sqlDB.Close()

	// Usar migrador adaptativo que maneja prefijos en dev y tabla de migraciones separada
	migrator := db.NewAdaptiveMigrator(&cfg, sqlDB, lg)

	// Ejecuta las migraciones de la carpeta correspondiente al driver
	path := "./migrations/pg" // default Postgres
	if cfg.DriverRelacional == "mysql" {
		path = "./migrations/mysql"
	}

	log.Printf("service=%s environment=%s driver=%s", cfg.ServiceName, cfg.Environment, cfg.DriverRelacional)
	if err := migrator.RunMigrations(path); err != nil {
		log.Fatal(err)
	}
}
