// cmd/api/main.go
package main

import (
	"context"
	"fmt"
	"net"
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
	rcgrpc "github.com/rem-gestion/rem-common/grpc"
	"github.com/rem-gestion/rem-common/logger"
	mw "github.com/rem-gestion/rem-common/middleware"

	/* ───── capas locales ───── */
	controller "github.com/rem-gestion/api-suite/property/src/controllers"
	"github.com/rem-gestion/api-suite/property/src/repository"
	"github.com/rem-gestion/api-suite/property/src/router"
	"github.com/rem-gestion/api-suite/property/src/services"
)

// waitForService waits for a service to be available at the given address
func waitForService(address string, serviceName string, lg *zap.Logger, maxWait time.Duration) error {
	lg.Info("waiting for service to be available",
		zap.String("service", serviceName),
		zap.String("address", address),
		zap.Duration("max_wait", maxWait))

	timeout := time.After(maxWait)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout waiting for %s at %s after %v", serviceName, address, maxWait)
		case <-ticker.C:
			conn, err := net.DialTimeout("tcp", address, 2*time.Second)
			if err == nil {
				conn.Close()
				lg.Info("service is available", zap.String("service", serviceName))
				return nil
			}
			lg.Debug("service not yet available",
				zap.String("service", serviceName),
				zap.Error(err))
		}
	}
}

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

	repo := repository.NewPropertyRepo(pg, lg.Named("repo"))

	/* ---------- wait for dependencies ---------- */
	// Wait for address service to be available
	addrTarget := cfg.GetServiceGRPCAddress("address")
	if err := waitForService(addrTarget, "address-svc", lg, 30*time.Second); err != nil {
		lg.Fatal("address service not available", zap.Error(err))
	}

	/* ---------- dial a address-svc ---------- */
	lg.Info("attempting to connect to address service", zap.String("target", addrTarget))

	addrConn, err := rcgrpc.Dial(addrTarget) // helper con timeout & keep-alive
	if err != nil {
		lg.Fatal("dial address-svc failed", zap.String("target", addrTarget), zap.Error(err))
	}
	defer addrConn.Close()

	lg.Info("successfully connected to address service")

	/* ---------- servicio de dominio ---------- */
	svc := services.New(repo, addrConn, lg)
	ctrl := controller.New(svc)

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
	router.Setup(r, ctrl) // /properties, /amenities...

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
