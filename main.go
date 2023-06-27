package main

import (
	"database/sql"
	"fmt"
	"github.com/micaelapucciariello/simplebank/api/client"
	db "github.com/micaelapucciariello/simplebank/db/sqlc"
	"github.com/micaelapucciariello/simplebank/utils"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := utils.LoadConfig("")
	if err != nil {
		log.Fatal("cannot get config: ", err)
	}
	conn, err := sql.Open(cfg.DriverName, cfg.SourceName)
	if err != nil {
		log.Fatal(fmt.Sprintf("cannot connect to db: %s", err))
	}

	store := db.NewStore(conn)
	server := client.NewServer(store)

	err = server.Start(cfg.ServerAddress)
	if err != nil {
		log.Fatal(fmt.Sprintf("cannot start server: %s", err))
	}
}
