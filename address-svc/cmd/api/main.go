package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rem-gestion/rem-common/config"
	"github.com/rem-gestion/rem-common/db"
	"github.com/rem-gestion/rem-common/logger"
	mw "github.com/rem-gestion/rem-common/middleware"
	"go.uber.org/zap"

	controller "github.com/rem-gestion/api-suite/address/src/controllers"
	"github.com/rem-gestion/api-suite/address/src/repository"
	"github.com/rem-gestion/api-suite/address/src/router"
	"github.com/rem-gestion/api-suite/address/src/services"
)

func main() {
	cfg := config.Load()
	lg := logger.New(cfg.Logger, "address-svc")
	pg, err := db.NewPostgres(cfg.Postgres)

	if err != nil {
		lg.Fatal("fatal: no pude conectar a Postgres", zap.Error(err))
	}

	repo := repository.New(pg)
	ctrl := controller.New(services.New(repo))

	r := gin.New()
	r.Use(
		mw.RequestID(),
		mw.APIKeyAuth("X-Api-Key", cfg.APIKey),
		mw.GinLogger(lg),
		mw.RecoveryWithZap(lg),
		mw.ErrorHandler(),
	)

	router.Setup(r, ctrl)

	r.Run(":4000")
}
