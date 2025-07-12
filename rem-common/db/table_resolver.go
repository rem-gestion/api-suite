package db

import (
	"fmt"
	"os"
)

// TableResolver maneja la resolución de nombres de tabla según el entorno
type TableResolver struct {
	ServiceName string
	Environment string
}

// NewTableResolver crea un nuevo resolvedor de tablas
func NewTableResolver(serviceName string) *TableResolver {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}

	return &TableResolver{
		ServiceName: serviceName,
		Environment: env,
	}
}

// GetTableName devuelve el nombre de tabla correcto según el entorno
func (tr *TableResolver) GetTableName(baseName string) string {
	if tr.Environment == "development" {
		// En desarrollo, usa prefijo del servicio
		return fmt.Sprintf("%s_%s", tr.ServiceName, baseName)
	}

	// En producción, usa el nombre base sin prefijo
	return baseName
}

// GetServicePrefix devuelve el prefijo del servicio para desarrollo
func (tr *TableResolver) GetServicePrefix() string {
	if tr.Environment == "development" {
		return tr.ServiceName + "_"
	}
	return ""
}

// IsDevEnvironment verifica si estamos en entorno de desarrollo
func (tr *TableResolver) IsDevEnvironment() bool {
	return tr.Environment == "development"
}

// TableNameFunc crea una función TableName() para modelos GORM
func (tr *TableResolver) TableNameFunc(baseName string) func() string {
	return func() string {
		return tr.GetTableName(baseName)
	}
}

// Global resolvers para cada servicio
var (
	PersonTableResolver  *TableResolver
	AddressTableResolver *TableResolver
	AuthTableResolver    *TableResolver
)

// InitializeResolvers inicializa los resolvers globales
func InitializeResolvers() {
	PersonTableResolver = NewTableResolver("person")
	AddressTableResolver = NewTableResolver("address")
	AuthTableResolver = NewTableResolver("auth")
}

// Helper functions para uso directo
func GetPersonTableName(baseName string) string {
	if PersonTableResolver == nil {
		PersonTableResolver = NewTableResolver("person")
	}
	return PersonTableResolver.GetTableName(baseName)
}

func GetAddressTableName(baseName string) string {
	if AddressTableResolver == nil {
		AddressTableResolver = NewTableResolver("address")
	}
	return AddressTableResolver.GetTableName(baseName)
}

func GetAuthTableName(baseName string) string {
	if AuthTableResolver == nil {
		AuthTableResolver = NewTableResolver("auth")
	}
	return AuthTableResolver.GetTableName(baseName)
}
