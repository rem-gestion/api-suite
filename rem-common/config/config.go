package config

import (
	"fmt"

	"github.com/rem-gestion/rem-common/environment"
)

// PostgresConfig agrupa todo lo que necesita Postgres
type PostgresConfig struct {
	Host     string // REM_POSTGRES_HOST
	Port     string // REM_POSTGRES_PORT (cambiado a string para flexibility)
	User     string // REM_POSTGRES_USER
	Password string // REM_POSTGRES_PASSWORD
	Database string // REM_POSTGRES_DB (renombrado para consistencia)
	SSLMode  string // REM_POSTGRES_SSLMODE
}

// MySQLConfig agrupa lo que necesita MySQL
type MySQLConfig struct {
	Host     string // REM_MYSQL_HOST
	Port     int    // REM_MYSQL_PORT
	User     string // REM_MYSQL_USER
	Password string // REM_MYSQL_PASSWORD
	DBName   string // REM_MYSQL_DB
	Charset  string // REM_MYSQL_CHARSET
}

// MongoConfig agrupa lo que necesita MongoDB
type MongoConfig struct {
	URI      string // REM_MONGO_URI (si queda vacío arma desde Host/Port/User/Password)
	Host     string // REM_MONGO_HOST
	Port     int    // REM_MONGO_PORT
	User     string // REM_MONGO_USER
	Password string // REM_MONGO_PASSWORD
	DBName   string // REM_MONGO_DB
	Timeout  int    // REM_MONGO_TIMEOUT en segundos
}

// RedisConfig para cache
type RedisConfig struct {
	Addr     string // REM_REDIS_ADDR
	Password string // REM_REDIS_PASSWORD
	DB       int    // REM_REDIS_DB
	Enabled  bool   // REM_REDIS_ENABLED
}

// RabbitConfig para RabbitMQ
type RabbitConfig struct {
	Host     string // REM_RABBIT_HOST
	Port     int    // REM_RABBIT_PORT
	User     string // REM_RABBIT_USER
	Password string // REM_RABBIT_PASSWORD
}

// LoggerConfig controla el nivel de logs
type LoggerConfig struct {
	Level string // REM_LOGGER_LEVEL (debug|info|warn|error)
}

// ServerConfig define puerto de escucha HTTP
type ServerConfig struct {
	Port int `envconfig:"REM_SERVER_PORT" default:"8080"`
}

// GRPCConfig define la configuración del servidor gRPC
// - Host: dirección IP o nombre de host donde escucha el servidor gRPC
// - Port: puerto en el que escucha el servidor gRPC
type GRPCConfig struct {
	Host string `envconfig:"REM_GRPC_HOST" default:"0.0.0.0"`
	Port int    `envconfig:"REM_GRPC_PORT" default:"50051"`
}

// Configuraciones gRPC específicas por servicio para desarrollo
type AuthGRPCConfig struct {
	Host string `envconfig:"REM_AUTH_GRPC_HOST" default:"0.0.0.0"`
	Port int    `envconfig:"REM_AUTH_GRPC_PORT" default:"50051"`
}

type AddressGRPCConfig struct {
	Host string `envconfig:"REM_ADDRESS_GRPC_HOST" default:"0.0.0.0"`
	Port int    `envconfig:"REM_ADDRESS_GRPC_PORT" default:"50052"`
}

type PersonGRPCConfig struct {
	Host string `envconfig:"REM_PERSON_GRPC_HOST" default:"0.0.0.0"`
	Port int    `envconfig:"REM_PERSON_GRPC_PORT" default:"50053"`
}

type OrganizationGRPCConfig struct {
	Host string `envconfig:"REM_ORGANIZATION_GRPC_HOST" default:"0.0.0.0"`
	Port int    `envconfig:"REM_ORGANIZATION_GRPC_PORT" default:"50054"`
}

type AddressServiceConfig struct {
	Host string `envconfig:"REM_ADDRESS_HOST" default:"127.0.0.1"`
	Port int    `envconfig:"REM_ADDRESS_PORT" default:"50051"`
}

type PersonServiceConfig struct {
	Host string `envconfig:"REM_PERSON_HOST" default:"127.0.0.1"`
	Port int    `envconfig:"REM_PERSON_PORT" default:"50052"`
}

type OrganizationServiceConfig struct {
	Host string `envconfig:"REM_ORGANIZATION_HOST" default:"127.0.0.1"`
	Port int    `envconfig:"REM_ORGANIZATION_PORT" default:"50054"`
}

type AuthServerConfig struct {
	Port int `envconfig:"REM_AUTH_HTTP_PORT" default:"4007"`
}

type AddressServerConfig struct {
	Port int `envconfig:"REM_ADDRESS_HTTP_PORT" default:"4011"`
}

type PersonServerConfig struct {
	Port int `envconfig:"REM_PERSON_HTTP_PORT" default:"4019"`
}

type OrganizationServerConfig struct {
	Port int `envconfig:"REM_ORGANIZATION_HTTP_PORT" default:"4020"`
}

type Config struct {
	// Información del entorno y servicio
	ServiceName         string                           // Detectado automáticamente
	Environment         environment.Environment          // Entorno actual
	EnvironmentDetector *environment.EnvironmentDetector // Detector para funciones helper

	DriverRelacional string `envconfig:"REM_DB_DRIVER"`
	Postgres         PostgresConfig
	MySQL            MySQLConfig
	Mongo            MongoConfig
	Redis            RedisConfig
	Rabbit           RabbitConfig

	GRPC             GRPCConfig             // Configuración gRPC genérica (para compatibilidad)
	AuthGRPC         AuthGRPCConfig         // Configuración gRPC específica para Auth
	AddressGRPC      AddressGRPCConfig      // Configuración gRPC específica para Address
	PersonGRPC       PersonGRPCConfig       // Configuración gRPC específica para Person
	OrganizationGRPC OrganizationGRPCConfig // Configuración gRPC específica para Organization

	Server             ServerConfig             // Configuración HTTP genérica (para compatibilidad)
	AuthServer         AuthServerConfig         // Configuración HTTP específica para Auth
	AddressServer      AddressServerConfig      // Configuración HTTP específica para Address
	PersonServer       PersonServerConfig       // Configuración HTTP específica para Person
	OrganizationServer OrganizationServerConfig // Configuración HTTP específica para Organization

	Address      AddressServiceConfig
	Person       PersonServiceConfig
	Organization OrganizationServiceConfig

	Logger LoggerConfig
	APIKey string `envconfig:"REM_API_KEY"` // REM_API_KEY: clave secreta que usa este servicio
}

// GetGRPCConfig retorna la configuración gRPC específica para el servicio actual
func (c *Config) GetGRPCConfig() (string, int) {
	switch c.ServiceName {
	case "auth":
		return c.AuthGRPC.Host, c.AuthGRPC.Port
	case "address":
		return c.AddressGRPC.Host, c.AddressGRPC.Port
	case "person":
		return c.PersonGRPC.Host, c.PersonGRPC.Port
	case "organization":
		return c.OrganizationGRPC.Host, c.OrganizationGRPC.Port
	default:
		// Fallback a configuración genérica
		return c.GRPC.Host, c.GRPC.Port
	}
}

// GetGRPCAddress retorna la dirección gRPC completa para el servicio actual
func (c *Config) GetGRPCAddress() string {
	host, port := c.GetGRPCConfig()
	return fmt.Sprintf("%s:%d", host, port)
}

// GetServerPort retorna el puerto HTTP específico para el servicio actual
func (c *Config) GetServerPort() int {
	switch c.ServiceName {
	case "auth":
		return c.AuthServer.Port
	case "address":
		return c.AddressServer.Port
	case "person":
		return c.PersonServer.Port
	case "organization":
		return c.OrganizationServer.Port
	default:
		// Fallback a configuración genérica
		return c.Server.Port
	}
}

// GetServiceGRPCAddress retorna la dirección gRPC de otro servicio específico
func (c *Config) GetServiceGRPCAddress(serviceName string) string {
	switch serviceName {
	case "auth":
		return fmt.Sprintf("%s:%d", c.AuthGRPC.Host, c.AuthGRPC.Port)
	case "address":
		return fmt.Sprintf("%s:%d", c.AddressGRPC.Host, c.AddressGRPC.Port)
	case "person":
		return fmt.Sprintf("%s:%d", c.PersonGRPC.Host, c.PersonGRPC.Port)
	case "organization":
		return fmt.Sprintf("%s:%d", c.OrganizationGRPC.Host, c.OrganizationGRPC.Port)
	default:
		// Fallback para Address y Person legacy configs
		if serviceName == "address-svc" {
			return fmt.Sprintf("%s:%d", c.Address.Host, c.Address.Port)
		}
		if serviceName == "person-svc" {
			return fmt.Sprintf("%s:%d", c.Person.Host, c.Person.Port)
		}
		if serviceName == "organization-svc" {
			return fmt.Sprintf("%s:%d", c.Organization.Host, c.Organization.Port)
		}
		return fmt.Sprintf("%s:%d", c.GRPC.Host, c.GRPC.Port)
	}
}
