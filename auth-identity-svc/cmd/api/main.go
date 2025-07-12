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
	"go.uber.org/zap"
	"google.golang.org/grpc/keepalive"

	/* ───── rem-common ───── */
	"github.com/rem-gestion/rem-common/config"
	"github.com/rem-gestion/rem-common/db"
	rcgrpc "github.com/rem-gestion/rem-common/grpc"
	"github.com/rem-gestion/rem-common/logger"
	mw "github.com/rem-gestion/rem-common/middleware"

	/* ───── protos ───── */
	authpb "github.com/rem-gestion/rem-common/protos/auth-identity/v1"

	/* ───── capas locales ───── */
	"github.com/rem-gestion/api-suite/auth-identity/src/controllers"
	grpcHandler "github.com/rem-gestion/api-suite/auth-identity/src/grpc"
	"github.com/rem-gestion/api-suite/auth-identity/src/repository"
	"github.com/rem-gestion/api-suite/auth-identity/src/router"
	"github.com/rem-gestion/api-suite/auth-identity/src/services"
)

func main() {
	/* ---------- carga de config & logger ---------- */
	cfg := config.Load()
	lg := logger.New(cfg.Logger, "auth-identity-svc")

	// Mostrar información del entorno
	serverPort := cfg.GetServerPort()
	_, grpcPort := cfg.GetGRPCConfig()
	lg.Info("starting auth-identity service",
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

	/* ---------- repositorios ---------- */
	accountRepo := repository.NewAccountRepository(pg)
	userRepo := repository.NewUserRepository(pg)

	/* ---------- dial a person-svc ---------- */
	personTarget := cfg.GetServiceGRPCAddress("person")
	personConn, err := rcgrpc.Dial(personTarget) // helper con timeout & keep-alive
	if err != nil {
		lg.Warn("dial person-svc failed", zap.Error(err))
		lg.Warn("Person service will not be available")
	}
	if personConn != nil {
		defer personConn.Close()
	}

	/* ---------- cliente person-svc ---------- */
	var personClient *services.PersonServiceClient
	if personConn != nil {
		personClient = services.NewPersonServiceClientWithConn(personConn)
	}

	/* ---------- servicios de dominio ---------- */
	authService := services.NewAuthService(
		accountRepo,
		userRepo,
		personClient,
		cfg.APIKey, // Usar APIKey como JWT secret por ahora
	)

	userService := services.NewUserService(
		userRepo,
		accountRepo,
		personClient,
	)

	/* ---------- controladores ---------- */
	authController := controllers.NewAuthController(authService)
	userController := controllers.NewUserController(userService)

	/* ---------- HTTP ---------- */
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(
		mw.RequestID(),
		mw.GinLogger(lg),
		mw.RecoveryWithZap(lg),
		mw.ErrorHandler(),
	)

	// Configurar rutas
	router.SetupRoutes(r, authController, userController, cfg.APIKey)
	r.GET("/health", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	httpSrv := &http.Server{
		Addr:         fmt.Sprintf(":%d", serverPort),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	/* ---------- gRPC ---------- */
	kp := keepalive.ServerParameters{Time: 2 * time.Hour, Timeout: 20 * time.Second}
	grpcSrv := rcgrpc.NewServer(lg, kp)
	grpcAddr := cfg.GetGRPCAddress()

	authpb.RegisterAuthIdentityServiceServer(grpcSrv, grpcHandler.New(authService, userService))

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		lg.Fatal("grpc listen failed", zap.Error(err))
	}

	/* ---------- lanzamos ambos servidores ---------- */
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

	/* ---------- graceful-shutdown ---------- */
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	lg.Info("shutting down…")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = httpSrv.Shutdown(ctx)
	grpcSrv.GracefulStop()

	if sqlDB, err := pg.DB(); err == nil {
		sqlDB.Close()
	}
}
