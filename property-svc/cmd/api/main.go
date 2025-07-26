// cmd/api/main.go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	/* ───── rem-common ───── */
	"github.com/rem-gestion/rem-common/config"
	"github.com/rem-gestion/rem-common/db"
	"github.com/rem-gestion/rem-common/logger"
	mw "github.com/rem-gestion/rem-common/middleware"

	/* ───── capas locales ───── */
	"github.com/rem-gestion/api-suite/property/src/controllers"
	"github.com/rem-gestion/api-suite/property/src/grpc/clients"
	"github.com/rem-gestion/api-suite/property/src/grpc/server"
	"github.com/rem-gestion/api-suite/property/src/repository"
	"github.com/rem-gestion/api-suite/property/src/router"
	"github.com/rem-gestion/api-suite/property/src/services"
)

func main() {
	/* ---------- carga de config & logger ---------- */
	cfg := config.Load()
	lg := logger.New(cfg.Logger, "property-svc")

	// Get property-specific port from environment variable
	serverPort := 4004 // default
	if portStr := os.Getenv("REM_PROPERTY_HTTP_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			serverPort = port
		}
	}

	// Get gRPC port from environment variable
	grpcPort := 50054 // default gRPC port for property service (matches .env.development)
	if portStr := os.Getenv("REM_PROPERTY_GRPC_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			grpcPort = port
		}
	}

	// Mostrar información del entorno
	lg.Info("starting property service",
		zap.String("service", cfg.ServiceName),
		zap.String("environment", cfg.Environment.String()),
		zap.String("database", cfg.Postgres.Database),
		zap.String("http_port", fmt.Sprintf("%d", serverPort)),
		zap.String("grpc_port", fmt.Sprintf("%d", grpcPort)))

	/* ---------- Postgres ---------- */
	pg, err := db.NewPostgres(cfg.Postgres)
	if err != nil {
		lg.Fatal("postgres connect failed", zap.Error(err))
	}

	/* ---------- gRPC clients ---------- */
	clientConfig := clients.LoadClientConfigFromEnv(&cfg)
	clientManager, err := clients.NewClientManager(clientConfig, lg)
	if err != nil {
		lg.Warn("failed to initialize gRPC clients", zap.Error(err))
		// Continue without external validations
		clientManager = nil
	}

	/* ---------- repositories ---------- */
	propertyRepo := repository.NewPropertyRepo(pg, lg)
	propertyTypeRepo := repository.NewPropertyTypeRepo(pg, lg)
	managerTypeRepo := repository.NewManagerTypeRepo(pg, lg)
	amenityRepo := repository.NewAmenityRepo(pg, lg)
	propertyManagementRepo := repository.NewPropertyManagementRepo(pg, lg)
	propertyAmenityRepo := repository.NewPropertyAmenityRepo(pg, lg)
	propertyListingRepo := repository.NewPropertyListingRepo(pg, lg)
	propertyMediaRepo := repository.NewPropertyMediaRepo(pg, lg)
	propertyValuationRepo := repository.NewPropertyValuationRepo(pg, lg)

	/* ---------- services ---------- */
	propertyService := services.NewPropertyService(
		propertyRepo,
		propertyTypeRepo,
		propertyManagementRepo,
		propertyAmenityRepo,
		propertyListingRepo,
		propertyMediaRepo,
		propertyValuationRepo,
		clientManager,
		lg,
	)
	propertyTypeService := services.NewPropertyTypeService(propertyTypeRepo, lg)
	managerTypeService := services.NewManagerTypeService(managerTypeRepo, lg)
	amenityService := services.NewAmenityService(amenityRepo, lg)
	propertyManagementService := services.NewPropertyManagementService(propertyManagementRepo, clientManager, lg)
	propertyListingService := services.NewPropertyListingService(propertyListingRepo, propertyRepo, clientManager, lg)
	propertyMediaService := services.NewPropertyMediaService(propertyMediaRepo, propertyRepo, clientManager, lg)
	propertyValuationService := services.NewPropertyValuationService(propertyValuationRepo, propertyRepo, clientManager, lg)

	/* ---------- controllers ---------- */
	propertyController := controllers.NewPropertyController(propertyService)
	propertyTypeController := controllers.NewPropertyTypeController(propertyTypeService)
	managerTypeController := controllers.NewManagerTypeController(managerTypeService)
	amenityController := controllers.NewAmenityController(amenityService)
	propertyManagementController := controllers.NewPropertyManagementController(propertyManagementService)
	propertyListingController := controllers.NewPropertyListingController(propertyListingService)
	propertyMediaController := controllers.NewPropertyMediaController(propertyMediaService)
	propertyValuationController := controllers.NewPropertyValuationController(propertyValuationService)

	/* ---------- router ---------- */
	appRouter := router.New(
		propertyController,
		propertyTypeController,
		managerTypeController,
		amenityController,
		propertyManagementController,
		propertyListingController,
		propertyMediaController,
		propertyValuationController,
	)

	/* ---------- HTTP ---------- */
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Health endpoint sin autenticación
	r.GET("/health", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	// Aplicar middlewares para rutas protegidas
	r.Use(
		mw.RequestID(),
		mw.APIKeyAuth("X-Api-Key", cfg.APIKey),
		mw.GinLogger(lg),
		mw.RecoveryWithZap(lg),
		mw.ErrorHandler(),
	)
	appRouter.SetupRoutes(r) // /properties, /amenities...

	/* ---------- gRPC server ---------- */
	grpcServer := server.NewPropertyGRPCServer(propertyService, lg)

	// TODO: Combined gRPC server for new services - will be implemented after PostMan
	// combinedGRPCServer := server.NewCombinedGRPCServer(
	//     propertyListingService,
	//     propertyMediaService,
	//     propertyValuationService,
	//     lg,
	// )

	httpSrv := &http.Server{Addr: fmt.Sprintf(":%d", serverPort), Handler: r}

	/* ---------- lanzar servidores ---------- */
	// Start original gRPC server (property service)
	go func() {
		lg.Info("gRPC listening", zap.String("addr", fmt.Sprintf("0.0.0.0:%d", grpcPort)))
		if err := grpcServer.Start(grpcPort); err != nil {
			lg.Fatal("gRPC failed", zap.Error(err))
		}
	}()

	// TODO: Start new combined gRPC server (new services) on a different port
	// newGrpcPort := grpcPort + 1
	// go func() {
	//     lg.Info("Combined gRPC listening", zap.String("addr", fmt.Sprintf("0.0.0.0:%d", newGrpcPort)))
	//     if err := combinedGRPCServer.Start(newGrpcPort); err != nil {
	//         lg.Fatal("Combined gRPC failed", zap.Error(err))
	//     }
	// }()

	// Start HTTP server
	go func() {
		lg.Info("REST listening", zap.String("addr", httpSrv.Addr))
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			lg.Fatal("gin failed", zap.Error(err))
		}
	}()

	/* ---------- graceful-shutdown ---------- */
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	lg.Info("shutting down…")

	// Stop gRPC server
	grpcServer.Stop()
	// TODO: Stop combined gRPC server when implemented
	// combinedGRPCServer.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = httpSrv.Shutdown(ctx)

	// Close gRPC clients
	if clientManager != nil {
		_ = clientManager.Close()
	}

	if sqlDB, err := pg.DB(); err == nil {
		sqlDB.Close()
	}
}
