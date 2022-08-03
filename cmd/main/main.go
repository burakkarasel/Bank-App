package main

import (
	"database/sql"
	"log"

	"github.com/burakkarasel/Bank-App/api"
	db "github.com/burakkarasel/Bank-App/db/sqlc"
	"github.com/burakkarasel/Bank-App/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".") // we pass the location of the file

	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)

	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
