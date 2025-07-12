package environment

import (
	"os"
	"path/filepath"
	"strings"
)

type Environment string

const (
	Development Environment = "development"
	Production  Environment = "production"
	Testing     Environment = "testing"
)

type EnvironmentDetector struct {
	env Environment
}

func NewEnvironmentDetector() *EnvironmentDetector {
	return &EnvironmentDetector{
		env: detectEnvironment(),
	}
}

func (ed *EnvironmentDetector) GetEnvironment() Environment {
	return ed.env
}

func (ed *EnvironmentDetector) IsDevelopment() bool {
	return ed.env == Development
}

func (ed *EnvironmentDetector) IsProduction() bool {
	return ed.env == Production
}

func (ed *EnvironmentDetector) IsTesting() bool {
	return ed.env == Testing
}

func detectEnvironment() Environment {
	// 1. Verificar variable de entorno explícita
	if env := os.Getenv("REM_ENVIRONMENT"); env != "" {
		switch strings.ToLower(env) {
		case "production", "prod":
			return Production
		case "development", "dev":
			return Development
		case "testing", "test":
			return Testing
		}
	}

	// 2. Detectar por herramienta de ejecución
	if isRunningWithAir() {
		return Development
	}

	// 3. Detectar por binario compilado
	if isCompiledBinary() {
		return Production
	}

	// 4. Detectar por directorio de trabajo
	if isInDevelopmentStructure() {
		return Development
	}

	// Default: production por seguridad
	return Production
}

func isRunningWithAir() bool {
	// Air típicamente setea esta variable
	if os.Getenv("AIR_TMP_DIR") != "" {
		return true
	}

	// Verificar si el proceso padre es air
	args := strings.Join(os.Args, " ")
	return strings.Contains(args, "tmp") ||
		strings.Contains(args, "__debug_bin") ||
		strings.Contains(args, "go_build")
}

func isCompiledBinary() bool {
	exe, err := os.Executable()
	if err != nil {
		return false
	}

	// Si el ejecutable no está en un directorio temporal, probablemente es producción
	return !strings.Contains(exe, "tmp") &&
		!strings.Contains(exe, "debug") &&
		!strings.Contains(exe, "go-build") &&
		!strings.Contains(exe, "go_build")
}

func isInDevelopmentStructure() bool {
	wd, err := os.Getwd()
	if err != nil {
		return false
	}

	// Buscar indicadores de estructura de desarrollo
	indicators := []string{
		"go.mod",
		".air.toml",
		"src/",
		"cmd/",
		".git/",
	}

	for _, indicator := range indicators {
		if _, err := os.Stat(filepath.Join(wd, indicator)); err == nil {
			return true
		}
	}

	return false
}

// String convierte el environment a string para logging
func (e Environment) String() string {
	return string(e)
}
