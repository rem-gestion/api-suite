package db

import (
	"sync"
	"users-api/src/config/envs"

	"go.uber.org/zap"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var once sync.Once
var dbInstance *gorm.DB

func ConnectDB(logger *zap.Logger) (*gorm.DB, error) {
	var err error
	POSTGRES_URI := envs.LoadEnvs(".env").Get("POSTGRES_URI")

	once.Do(func() {
		dbInstance, err = gorm.Open(postgres.Open(POSTGRES_URI), &gorm.Config{})
		if err != nil {
			logger.Fatal("[USERS-API] Error al conectar con PostgreSQL", zap.Error(err))
		}

		sqlDB, err := dbInstance.DB()
		if err != nil {
			logger.Fatal("[USERS-API] Error al obtener la conexión SQL", zap.Error(err))
		}
		_, err = sqlDB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
		if err != nil {
			logger.Warn("[USERS-API] Error al habilitar la extensión uuid-ossp", zap.Error(err))
		}

		// Si deseas realizar migraciones automáticas
		/*
			err = dbInstance.AutoMigrate(&models.Inscripto{})
			if err != nil {
				logger.Fatal("[INSCRIPTION-API] Error al realizar la migración automática", zap.Error(err))
			}
		*/

		logger.Info("[USERS-API] Conexión a PostgreSQL establecida")
	})

	return dbInstance, err
}
