package config

// PostgresConfig agrupa todo lo que necesita Postgres
type PostgresConfig struct {
	Host     string // REM_POSTGRES_HOST
	Port     int    // REM_POSTGRES_PORT
	User     string // REM_POSTGRES_USER
	Password string // REM_POSTGRES_PASSWORD
	DBName   string // REM_POSTGRES_DBNAME
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

// ServerConfig define puerto de escucha
type ServerConfig struct {
	Port int // REM_SERVER_PORT
}

// Config engloba todo, incluyendo el driver relacional
type Config struct {
	// DriverRelacional elige "postgres" o "mysql"
	DriverRelacional string `envconfig:"DB_DRIVER"`
	Postgres         PostgresConfig
	MySQL            MySQLConfig
	Mongo            MongoConfig
	Redis            RedisConfig
	Rabbit           RabbitConfig
	Logger           LoggerConfig
	Server           ServerConfig
	APIKey           string `envconfig:"API_KEY"` // REM_API_KEY: clave secreta que usa este servicio
}
