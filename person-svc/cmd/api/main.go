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

	"github.com/rem-gestion/rem-common/db"
	"go.uber.org/zap"
	"google.golang.org/grpc/keepalive"

	controller "github.com/rem-gestion/api-suite/person/src/controllers"
	grpcHandler "github.com/rem-gestion/api-suite/person/src/grpc"
	"github.com/rem-gestion/api-suite/person/src/repository"
	"github.com/rem-gestion/api-suite/person/src/router"
	"github.com/rem-gestion/api-suite/person/src/services"

	pb "github.com/rem-gestion/api-suite/person/internal/pb"
)

func main() {
	cfg := config.Load()
	lg := logger.New(cfg.Logger, "person-svc")

	/* DB ------------------------------------------------- */
	pg, err := db.NewPostgres(cfg.Postgres)
	if err != nil {
		lg.Fatal("postgres connect failed", zap.Error(err))
	}
	repo := repository.New(pg, lg.Named("repo"))
	svc := services.New(repo, lg)
	ctrl := controller.New(svc)

	/* REST ---------------------------------------------- */
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(
		mw.RequestID(),
		mw.APIKeyAuth("X-Api-Key", cfg.APIKey),
		mw.GinLogger(lg),
		mw.RecoveryWithZap(lg),
		mw.ErrorHandler(),
	)
	router.Setup(r, ctrl)
	r.GET("/health", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	httpSrv := &http.Server{Addr: ":4000", Handler: r}

	/* gRPC ---------------------------------------------- */
	kp := keepalive.ServerParameters{Time: 2 * time.Hour, Timeout: 20 * time.Second}
	grpcSrv := rcgrpc.NewServer(lg, kp)
	grpcAddr := fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port)

	pb.RegisterPersonServiceServer(grpcSrv, grpcHandler.New(svc))

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

	if sqlDB, err := pg.DB(); err == nil {
		sqlDB.Close()
	}
}
