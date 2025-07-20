package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	/* ───── rem-common ───── */
	"github.com/rem-gestion/rem-common/broker"
	"github.com/rem-gestion/rem-common/cache"
	"github.com/rem-gestion/rem-common/config"
	"github.com/rem-gestion/rem-common/db"
	rcgrpc "github.com/rem-gestion/rem-common/grpc"
	"github.com/rem-gestion/rem-common/logger"
	mw "github.com/rem-gestion/rem-common/middleware"

	/* ───── capas locales ───── */
	"github.com/rem-gestion/api-suite/organization/src/clients"
	"github.com/rem-gestion/api-suite/organization/src/controllers"
	"github.com/rem-gestion/api-suite/organization/src/repository"
	"github.com/rem-gestion/api-suite/organization/src/router"
	"github.com/rem-gestion/api-suite/organization/src/services"
)

func main() {
	/* ---------- carga de config & logger ---------- */
	cfg := config.Load()
	lg := logger.New(cfg.Logger, "organization-svc")

	// Mostrar información del entorno
	serverPort := cfg.GetServerPort()
	lg.Info("starting organization service",
		zap.String("service", cfg.ServiceName),
		zap.String("environment", cfg.Environment.String()),
		zap.String("database", cfg.Postgres.Database),
		zap.String("http_port", fmt.Sprintf("%d", serverPort)))

	/* ---------- Postgres ---------- */
	pg, err := db.NewPostgres(cfg.Postgres)
	if err != nil {
		lg.Fatal("postgres connect failed", zap.Error(err))
	}

	/* ---------- repositorios ---------- */
	repos := repository.NewRepositories(pg, lg)

	/* ---------- dial a external services ---------- */
	// Auth-identity service
	authTarget := cfg.GetServiceGRPCAddress("auth")
	authConn, err := rcgrpc.Dial(authTarget)
	if err != nil {
		lg.Warn("dial auth-identity-svc failed", zap.Error(err))
		lg.Warn("Auth-identity service will not be available")
	}
	if authConn != nil {
		defer authConn.Close()
	}

	// Person service
	personTarget := cfg.GetServiceGRPCAddress("person")
	personConn, err := rcgrpc.Dial(personTarget)
	if err != nil {
		lg.Warn("dial person-svc failed", zap.Error(err))
		lg.Warn("Person service will not be available")
	}
	if personConn != nil {
		defer personConn.Close()
	}

	// Address service
	addressTarget := cfg.GetServiceGRPCAddress("address")
	addressConn, err := rcgrpc.Dial(addressTarget)
	if err != nil {
		lg.Warn("dial address-svc failed", zap.Error(err))
		lg.Warn("Address service will not be available")
	}
	if addressConn != nil {
		defer addressConn.Close()
	}

	/* ---------- gRPC client manager ---------- */
	var clientManager *clients.ClientManager
	if authConn != nil && personConn != nil && addressConn != nil {
		clientManager = clients.NewClientManager(authConn, personConn, addressConn)
		lg.Info("gRPC clients initialized successfully")
	} else {
		lg.Warn("Some gRPC services unavailable - operating in degraded mode")
		// En modo degradado, crear client manager con conexiones nulas para evitar panics
		clientManager = &clients.ClientManager{
			AuthIdentity: nil,
			Person:       nil,
			Address:      nil,
		}
	}

	/* ---------- servicios de dominio ---------- */
	// Cache service - inicializar Redis si está disponible
	var cacheService services.CacheService
	redisClient, err := cache.NewRedis(cfg.Redis)
	if err != nil {
		lg.Warn("Redis connection failed - using no-op cache", zap.Error(err))
		cacheService = &services.NoOpCacheService{}
	} else {
		cacheService = services.NewCacheService(redisClient, lg.Named("cache"))
	}

	// Event service - inicializar RabbitMQ si está disponible
	var eventService services.EventService
	rabbitConn, err := broker.NewRabbit(cfg.Rabbit)
	if err != nil {
		lg.Warn("RabbitMQ connection failed - using no-op events", zap.Error(err))
		eventService = &services.NoOpEventService{}
	} else {
		eventService = services.NewEventService(rabbitConn, lg.Named("events"))
	}

	// Crear todos los servicios usando el agregador
	allServices := services.NewServices(repos, clientManager, cacheService, eventService, lg)

	/* ---------- controladores ---------- */
	// Crear todos los controladores usando el agregador
	allControllers := controllers.NewControllers(allServices, lg)

	/* ---------- HTTP ---------- */
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(
		mw.RequestID(),
		mw.GinLogger(lg),
		gin.Recovery(),
		mw.ErrorHandler(),
	)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "organization-svc",
		})
	})

	// Setup routes para el MVP completo
	router.SetupRoutes(r, allControllers, cfg.APIKey)

	/* ---------- HTTP Server ---------- */
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", serverPort),
		Handler: r,
	}

	go func() {
		lg.Info("HTTP server starting", zap.Int("port", serverPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			lg.Fatal("listen failed", zap.Error(err))
		}
	}()

	/* ---------- Graceful Shutdown ---------- */
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	lg.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		lg.Fatal("server forced to shutdown", zap.Error(err))
	}

	lg.Info("server exiting")
}
