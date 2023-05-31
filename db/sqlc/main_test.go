package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	driverName = "postgres"
	sourceName = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open(driverName, sourceName)
	if err != nil {
		log.Fatal(fmt.Sprintf("cannot connect to db: %s", err))
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
