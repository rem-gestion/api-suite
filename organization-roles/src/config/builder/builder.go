package builder

import (
	"users-api/src/client"
	"users-api/src/config/db"
	"users-api/src/config/log"
	"users-api/src/config/redis"
	"users-api/src/controllers"
	"users-api/src/router"
	"users-api/src/services"

	"github.com/gin-gonic/gin"
	redisClient "github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AppBuilder struct {
	db             *gorm.DB
	redisClient    *redisClient.Client
	Logger         *zap.Logger
	userRepo       client.UserRepository
	userService    services.UserService
	authService    services.AuthService
	userController *controllers.UserController
	authController *controllers.AuthController
	router         *gin.Engine
}

func NewAppBuilder() *AppBuilder {
	return &AppBuilder{}
}

func BuildApp() *AppBuilder {
	return NewAppBuilder().
		BuildLogger().
		BuildDBConnection().
		BuildUserRepo().
		BuildUserService().
		BuildUserController().
		BuildRouter()
}

func (b *AppBuilder) BuildLogger() *AppBuilder {
	b.Logger = log.GetLogger()
	b.Logger.Info("[USERS-API] Logger inicializado")
	return b
}

func (b *AppBuilder) BuildDBConnection() *AppBuilder {
	var err error
	b.db, err = db.ConnectDB(b.Logger)
	if err != nil {
		b.Logger.Fatal("[USERS-API] Error al conectar a la base de datos", zap.Error(err))
	}
	b.redisClient = redis.ConnectRedis()
	if b.redisClient == nil {
		b.Logger.Warn("[USERS-API] Redis no está disponible. La aplicación funcionará sin caché.")
	}
	return b
}

func (b *AppBuilder) DisconnectDB() {
	if b.db != nil {
		sqlDB, err := b.db.DB()
		if err != nil {
			b.Logger.Error("[USERS-API] Error al obtener la conexión SQL", zap.Error(err))
		} else {
			if err := sqlDB.Close(); err != nil {
				b.Logger.Error("[USERS-API] Error al desconectar de la base de datos", zap.Error(err))
			} else {
				b.Logger.Info("[USERS-API] Conexión a la base de datos cerrada")
			}
		}
	}
	_ = b.Logger.Sync()
}

func (b *AppBuilder) BuildUserRepo() *AppBuilder {
	b.userRepo = client.NewUserRepository(b.db, b.Logger)
	b.Logger.Info("[USERS-API] Repositorio de usuarios inicializado")
	return b
}

func (b *AppBuilder) BuildUserService() *AppBuilder {
	b.userService = services.NewUserService(b.userRepo, b.redisClient, b.Logger)
	b.Logger.Info("[USERS-API] Servicio de usuarios inicializado")
	b.authService = services.NewAuthService(b.userRepo, b.redisClient, b.Logger)
	b.Logger.Info("[USERS-API] Servicio de autenticación inicializado")
	return b
}

func (b *AppBuilder) BuildUserController() *AppBuilder {
	b.userController = controllers.NewUserController(b.userService, b.Logger)
	b.Logger.Info("[USERS-API] Controlador de usuarios inicializado")
	b.authController = controllers.NewAuthController(b.authService, b.Logger)
	b.Logger.Info("[USERS-API] Controlador de autenticación inicializado")
	return b
}

func (b *AppBuilder) BuildRouter() *AppBuilder {
	b.router = gin.Default()
	router.SetupRoutes(b.router, b.userController, b.authController)
	b.Logger.Info("[USERS-API] Rutas configuradas")
	return b
}

func (b *AppBuilder) GetRouter() *gin.Engine {
	return b.router
}
