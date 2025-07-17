package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rem-gestion/rem-common/config"
	"github.com/rem-gestion/rem-common/db"
)

func main() {
	cfg := config.Load()
	dbConn, err := db.NewPostgres(cfg.Postgres)
	if err != nil {
		panic(err)
	}

	direction := "up"
	if len(os.Args) > 1 && (os.Args[1] == "down" || os.Args[1] == "up") {
		direction = os.Args[1]
	}

	var sqlFile string
	if direction == "down" {
		sqlFile = "property-svc/migrations/pg/0001_initial.down.sql"
	} else {
		sqlFile = "property-svc/migrations/pg/0001_initial.up.sql"
	}

	sqlBytes, err := ioutil.ReadFile(sqlFile)
	if err != nil {
		fmt.Printf("Error leyendo migración %s: %v\n", direction, err)
		os.Exit(1)
	}

	err = dbConn.Exec(string(sqlBytes)).Error
	if err != nil {
		fmt.Printf("Error ejecutando migración %s: %v\n", direction, err)
		os.Exit(1)
	}

	fmt.Printf("Migración %s ejecutada correctamente.\n", direction)
}
