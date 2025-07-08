// cmd/api/main.go
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
	personpb "github.com/rem-gestion/rem-common/protos/person/v1"

	/* ───── capas locales ───── */
	controller "github.com/rem-gestion/api-suite/person/src/controllers"
	grpcHandler "github.com/rem-gestion/api-suite/person/src/grpc"
	"github.com/rem-gestion/api-suite/person/src/repository"
	"github.com/rem-gestion/api-suite/person/src/router"
	"github.com/rem-gestion/api-suite/person/src/services"
)

func main() {
	/* ---------- carga de config & logger ---------- */
	cfg := config.Load()
	lg := logger.New(cfg.Logger, "person-svc")

	/* ---------- Postgres ---------- */
	pg, err := db.NewPostgres(cfg.Postgres)
	if err != nil {
		lg.Fatal("postgres connect failed", zap.Error(err))
	}

	repo := repository.NewPersonRepo(pg, lg.Named("repo"))

	/* ---------- dial a address-svc ---------- */
	addrTarget := fmt.Sprintf("%s:%d", cfg.Address.Host, cfg.Address.Port)
	addrConn, err := rcgrpc.Dial(addrTarget) // helper con timeout & keep-alive
	if err != nil {
		lg.Fatal("dial address-svc failed", zap.Error(err))
	}
	defer addrConn.Close()

	/* ---------- servicio de dominio ---------- */
	svc := services.New(repo, addrConn, lg)
	ctrl := controller.New(svc)

	/* ---------- HTTP ---------- */
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(
		mw.RequestID(),
		mw.APIKeyAuth("X-Api-Key", cfg.APIKey),
		mw.GinLogger(lg),
		mw.RecoveryWithZap(lg),
		mw.ErrorHandler(),
	)
	router.Setup(r, ctrl) // /persons, /contacts…
	r.GET("/health", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	httpSrv := &http.Server{Addr: ":4001", Handler: r}

	/* ---------- gRPC ---------- */
	kp := keepalive.ServerParameters{Time: 2 * time.Hour, Timeout: 20 * time.Second}
	grpcSrv := rcgrpc.NewServer(lg, kp)
	grpcAddr := fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port)

	personpb.RegisterPersonServiceServer(grpcSrv, grpcHandler.New(svc))

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
