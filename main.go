package main

import (
	"database/sql"
	"fmt"
	"github.com/micaelapucciariello/simplebank/api"
	db "github.com/micaelapucciariello/simplebank/db/sqlc"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	//cfg, err := utils.LoadConfig(".")
	//if err != nil {
	//	log.Fatal("cannot get config: ", err)
	//}

	driverName := "postgres"
	sourceName := "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress := "0.0.0.0:8080"

	conn, err := sql.Open(driverName, sourceName)
	if err != nil {
		log.Fatal(fmt.Sprintf("cannot connect to db: %s", err))
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal(fmt.Sprintf("cannot start server: %s", err))
	}
}
