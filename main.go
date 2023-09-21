package main

import (
	"database/sql"
	"log"

	"github.com/corey888773/golang-course/api"
	db "github.com/corey888773/golang-course/db/sqlc"
	"github.com/corey888773/golang-course/util"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config")
	}

	connector, err := pq.NewConnector(config.DbSource)
	if err != nil {
		log.Fatal("cannot create connector", err)
	}
	connection := sql.OpenDB(connector)

	defer connection.Close()

	store := db.NewStore(connection)
	server := api.NewServer(store)

	err = server.Start(config.ServerAdress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
