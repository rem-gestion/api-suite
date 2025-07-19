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

	// Mostrar información del entorno
	lg.Info("starting property service",
		zap.String("service", cfg.ServiceName),
		zap.String("environment", cfg.Environment.String()),
		zap.String("database", cfg.Postgres.Database),
		zap.String("http_port", fmt.Sprintf("%d", serverPort)))

	/* ---------- Postgres ---------- */
	pg, err := db.NewPostgres(cfg.Postgres)
	if err != nil {
		lg.Fatal("postgres connect failed", zap.Error(err))
	}

	/* ---------- repositories ---------- */
	propertyRepo := repository.NewPropertyRepo(pg, lg)
	propertyTypeRepo := repository.NewPropertyTypeRepo(pg, lg)
	managerTypeRepo := repository.NewManagerTypeRepo(pg, lg)
	amenityRepo := repository.NewAmenityRepo(pg, lg)
	propertyManagementRepo := repository.NewPropertyManagementRepo(pg, lg)
	propertyAmenityRepo := repository.NewPropertyAmenityRepo(pg, lg)

	/* ---------- services ---------- */
	propertyService := services.NewPropertyService(
		propertyRepo,
		propertyTypeRepo,
		propertyManagementRepo,
		propertyAmenityRepo,
		lg,
	)
	propertyTypeService := services.NewPropertyTypeService(propertyTypeRepo, lg)
	managerTypeService := services.NewManagerTypeService(managerTypeRepo, lg)
	amenityService := services.NewAmenityService(amenityRepo, lg)
	propertyManagementService := services.NewPropertyManagementService(propertyManagementRepo, lg)

	/* ---------- controllers ---------- */
	propertyController := controllers.NewPropertyController(propertyService)
	propertyTypeController := controllers.NewPropertyTypeController(propertyTypeService)
	managerTypeController := controllers.NewManagerTypeController(managerTypeService)
	amenityController := controllers.NewAmenityController(amenityService)
	propertyManagementController := controllers.NewPropertyManagementController(propertyManagementService)

	/* ---------- router ---------- */
	appRouter := router.New(propertyController, propertyTypeController, managerTypeController, amenityController, propertyManagementController)

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

	httpSrv := &http.Server{Addr: fmt.Sprintf(":%d", serverPort), Handler: r}

	/* ---------- lanzar servidor HTTP ---------- */
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = httpSrv.Shutdown(ctx)

	if sqlDB, err := pg.DB(); err == nil {
		sqlDB.Close()
	}
}
