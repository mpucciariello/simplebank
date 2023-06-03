package main

import (
	"database/sql"
	"fmt"
	api "github.com/micaelapucciariello/simplebank/api/server"
	db "github.com/micaelapucciariello/simplebank/db/sqlc"
	"log"

	_ "github.com/lib/pq"
)

const (
	driverName    = "postgres"
	sourceName    = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(driverName, sourceName)
	if err != nil {
		log.Fatal(fmt.Sprintf("cannot connect to db: %s", err))
	}

	store := db.NewStore(conn)
	server := api.New(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal(fmt.Sprintf("cannot start server: %s", err))
	}
}
