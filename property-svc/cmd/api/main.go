package main

import (
	"github.com/rem-gestion/property-svc/src/controllers"
	"github.com/rem-gestion/property-svc/src/router"
	"github.com/rem-gestion/rem-common/config"
	"github.com/rem-gestion/rem-common/db"
)

func main() {
	cfg := config.Load()
	dbConn, err := db.NewPostgres(cfg.Postgres)
	if err != nil {
		panic(err)
	}
	controllers.SetDB(dbConn)
	r := router.SetupRouter()
	r.Run(":4003")
}
