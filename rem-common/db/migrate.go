// rem-common/db/migrate.go
package db

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"github.com/rem-gestion/rem-common/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// RunMigrations aplica las migraciones encontradas en `migrationsPath`.
// driver: "postgres" | "mysql" | "mongo"
// direction: "up" | "down" | "steps N"
func RunMigrations(driver string, cfg config.Config, migrationsPath, direction string) error {
	switch driver {
	case "postgres":
		return runSQLMigrate(
			func(db *sql.DB) (database.Driver, error) {
				return postgres.WithInstance(db, &postgres.Config{})
			},
			buildPgDSN(cfg.Postgres),
			"postgres",
			migrationsPath,
			direction,
		)

	case "mysql":
		return runSQLMigrate(
			func(db *sql.DB) (database.Driver, error) {
				return mysql.WithInstance(db, &mysql.Config{})
			},
			buildMyDSN(cfg.MySQL),
			"mysql",
			migrationsPath,
			direction,
		)

	case "mongo":
		return runMongoMigrate(cfg.Mongo, migrationsPath, direction)

	default:
		return fmt.Errorf("driver %s no soportado", driver)
	}
}

/* ---------- helpers SQL ---------- */

func runSQLMigrate(
	newDriver func(*sql.DB) (database.Driver, error),
	dsn, dbName, migrationsPath, direction string,
) error {
	sqlDB, err := sql.Open(dbName, dsn)
	if err != nil {
		return err
	}
	driver, err := newDriver(sqlDB)
	if err != nil {
		return err
	}
	return applyMigrations(driver, dbName, migrationsPath, direction)
}

/* ---------- helper Mongo ---------- */

func runMongoMigrate(cfg config.MongoConfig, migrationsPath, direction string) error {
	uri := cfg.URI
	if uri == "" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%d",
			cfg.User, cfg.Password, cfg.Host, cfg.Port)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	driver, err := mongodb.WithInstance(client, &mongodb.Config{
		DatabaseName: cfg.DBName,
	})
	if err != nil {
		return err
	}
	return applyMigrations(driver, "mongodb", migrationsPath, direction)
}

/* ---------- motor común (up|down|steps) ---------- */

func applyMigrations(driver database.Driver, dbName, migrationsPath, direction string) error {
	abs, _ := filepath.Abs(migrationsPath)
	m, err := migrate.NewWithDatabaseInstance("file://"+abs, dbName, driver)
	if err != nil {
		return err
	}
	switch direction {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()
	default:
		var n int
		fmt.Sscanf(direction, "steps %d", &n)
		err = m.Steps(n)
	}
	if err == migrate.ErrNoChange {
		return nil
	}
	return err
}

/* ---------- DSN builders ---------- */

func buildPgDSN(c config.PostgresConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode)
}

func buildMyDSN(c config.MySQLConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.Charset)
}
