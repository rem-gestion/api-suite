package main

import (
	"log"
	"os"

	"github.com/rem-gestion/rem-common/config"
	"github.com/rem-gestion/rem-common/db"
)

func main() {
	// Carga las REM_* definidas en .env
	cfg := config.Load()

	// up | down | "steps -1"  (por defecto: up)
	dir := "up"
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	// Ejecuta las migraciones de la carpeta correspondiente al driver
	path := "./migrations/pg" // default Postgres
	if cfg.DriverRelacional == "mysql" {
		path = "./migrations/mysql"
	}
	log.Printf("driver=%s", cfg.DriverRelacional)
	if err := db.RunMigrations(cfg.DriverRelacional, cfg, path, dir); err != nil {
		log.Fatal(err)
	}
}
