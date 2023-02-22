package database

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/abc_valera/flugo/config"
	_ "github.com/lib/pq"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	config, err := config.LoadConfig("./..")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	conn, err := sql.Open(config.DatabaseDriver, config.DatabaseUrl)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	err = conn.Ping()
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
