package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/rem-gestion/rem-common/environment"
)

// LoadAdaptive carga la configuración de manera adaptativa según el entorno
func LoadAdaptive() Config {
	// Detectar el nombre del servicio automáticamente
	serviceName := detectServiceName()
	envDetector := environment.NewEnvironmentDetector()
	dbResolver := NewDatabaseConfigResolver(serviceName)

	// Cargar .env específico del entorno si existe
	loadEnvironmentFiles(envDetector.GetEnvironment())

	// Cargar configuración base
	var cfg Config
	if err := envconfig.Process("REM", &cfg); err != nil {
		log.Fatalf("rem-common/config: error cargando envs: %v", err)
	}

	// Sobrescribir configuración adaptativa
	cfg.Postgres = dbResolver.ResolvePostgresConfig()

	// Setear información del entorno
	cfg.ServiceName = serviceName
	cfg.Environment = envDetector.GetEnvironment()
	cfg.EnvironmentDetector = envDetector

	return cfg
}

// Load mantiene compatibilidad hacia atrás
func Load() Config {
	return LoadAdaptive()
}

func detectServiceName() string {
	// 1. Variable de entorno explícita
	if name := os.Getenv("REM_SERVICE_NAME"); name != "" {
		return name
	}

	// 2. Detectar por directorio actual
	wd, _ := os.Getwd()
	if strings.Contains(wd, "person-svc") {
		return "person"
	}
	if strings.Contains(wd, "address-svc") {
		return "address"
	}
	if strings.Contains(wd, "auth-identity-svc") {
		return "auth"
	}
	if strings.Contains(wd, "property-svc") {
		return "property"
	}
	if strings.Contains(wd, "organization-svc") {
		return "organization"
	}

	// 3. Detectar por ejecutable
	if exe, err := os.Executable(); err == nil {
		base := filepath.Base(exe)
		if strings.Contains(base, "person") {
			return "person"
		}
		if strings.Contains(base, "address") {
			return "address"
		}
		if strings.Contains(base, "auth") {
			return "auth"
		}
		if strings.Contains(base, "property") {
			return "property"
		}
		if strings.Contains(base, "organization") {
			return "organization"
		}
	}

	return "unknown"
}

func loadEnvironmentFiles(env environment.Environment) {
	// Lista de archivos a intentar cargar en orden de prioridad
	envFiles := []string{
		fmt.Sprintf(".env.%s.local", env),
		fmt.Sprintf(".env.%s", env),
		".env.local",
		".env",
	}

	for _, file := range envFiles {
		if err := godotenv.Load(file); err == nil {
			log.Printf("rem-common/config: loaded %s", file)
			break
		}
	}
}
