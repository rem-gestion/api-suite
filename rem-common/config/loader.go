package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Load carga variables de entorno desde un archivo .env (si existe)
// y luego procesa todas las vars con prefijo REM_ en el struct Config.
func Load() Config {
	// Intentar cargar .env (no es error si no existe)
	if err := godotenv.Load(); err != nil {
		log.Printf("rem-common/config: .env no encontrado o no se pudo cargar: %v", err)
	}

	var cfg Config
	if err := envconfig.Process("REM", &cfg); err != nil {
		log.Fatalf("rem-common/config: error cargando envs: %v", err)
	}
	return cfg
}
