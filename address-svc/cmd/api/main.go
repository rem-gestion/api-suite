package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/rem-gestion/rem-common/config"
	rcgrpc "github.com/rem-gestion/rem-common/grpc"
	"github.com/rem-gestion/rem-common/logger"
	mw "github.com/rem-gestion/rem-common/middleware"

	"go.uber.org/zap"
	"google.golang.org/grpc/keepalive"

	controller "github.com/rem-gestion/api-suite/address/src/controllers"
	grpcHandler "github.com/rem-gestion/api-suite/address/src/grpc"
	"github.com/rem-gestion/api-suite/address/src/repository"
	"github.com/rem-gestion/api-suite/address/src/router"
	"github.com/rem-gestion/api-suite/address/src/services"

	pb "github.com/rem-gestion/rem-common/protos/address/v1"
)

func main() {
	cfg := config.Load()
	lg := logger.New(cfg.Logger, "address-svc")

	// Mostrar información del entorno
	serverPort := cfg.GetServerPort()
	_, grpcPort := cfg.GetGRPCConfig()
	lg.Info("starting address service",
		zap.String("service", cfg.ServiceName),
		zap.String("environment", cfg.Environment.String()),
		zap.String("database", cfg.Postgres.Database),
		zap.String("http_port", fmt.Sprintf("%d", serverPort)),
		zap.String("grpc_port", fmt.Sprintf("%d", grpcPort)))

	/* DB ------------------------------------------------- */
	lg.Info("initializing adaptive repository with database resilience")

	// Usar repositorio adaptativo mejorado que maneja conexión internamente
	repo := repository.NewAdaptive(&cfg.Postgres, lg.Named("adaptive-repo"))

	// Dar tiempo al repositorio para intentar la conexión inicial
	time.Sleep(1 * time.Second)

	if repo.IsConnected() {
		lg.Info("database connection established - operating in database mode")
	} else {
		lg.Warn("database unavailable - operating in memory fallback mode with automatic retry every 15s")
	}

	svc := services.New(repo, lg)
	ctrl := controller.New(svc)

	/* REST ---------------------------------------------- */
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Health endpoints sin autenticación
	r.GET("/health", func(c *gin.Context) { c.String(http.StatusOK, "ok") })
	r.GET("/health/detailed", func(c *gin.Context) {
		stats := svc.GetRepositoryStats()
		c.JSON(http.StatusOK, gin.H{
			"status":                "ok",
			"repository":            stats,
			"using_memory_fallback": svc.IsUsingMemoryFallback(),
		})
	})

	// Aplicar middlewares para rutas protegidas
	r.Use(
		mw.RequestID(),
		mw.APIKeyAuth("X-Api-Key", cfg.APIKey),
		mw.GinLogger(lg),
		mw.RecoveryWithZap(lg),
		mw.ErrorHandler(),
	)
	router.Setup(r, ctrl)
	r.POST("/admin/sync-memory", func(c *gin.Context) {
		if err := svc.ForceMemorySync(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "sync triggered"})
	})
	r.POST("/admin/force-reconnect", func(c *gin.Context) {
		svc.ForceDatabaseReconnection()
		c.JSON(http.StatusOK, gin.H{"message": "reconnection attempt triggered"})
	})

	httpSrv := &http.Server{Addr: fmt.Sprintf(":%d", serverPort), Handler: r}

	/* gRPC ---------------------------------------------- */
	kp := keepalive.ServerParameters{Time: 2 * time.Hour, Timeout: 20 * time.Second}
	grpcSrv := rcgrpc.NewServer(lg, kp)
	grpcAddr := cfg.GetGRPCAddress()

	pb.RegisterAddressServiceServer(grpcSrv, grpcHandler.New(svc))

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		lg.Fatal("grpc listen failed", zap.Error(err))
	}

	/* arrancamos ambos en paralelo ---------------------- */
	go func() {
		lg.Info("REST listening", zap.String("addr", httpSrv.Addr))
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			lg.Fatal("gin failed", zap.Error(err))
		}
	}()
	go func() {
		lg.Info("gRPC listening", zap.String("addr", grpcAddr))
		if err := grpcSrv.Serve(lis); err != nil {
			lg.Fatal("grpc failed", zap.Error(err))
		}
	}()

	/* graceful-shutdown --------------------------------- */
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	lg.Info("shutting down…")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = httpSrv.Shutdown(ctx)
	grpcSrv.GracefulStop()

	// Cerrar repositorio adaptativo (detiene retry loop y cierra conexiones)
	repo.Close()
}
