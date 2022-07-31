package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/burakkarasel/Bank-App/db/dsn"
	_ "github.com/lib/pq"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open(dsn.DBDriver, dsn.DBSource)
	if err != nil {
		log.Fatal("cannot connect to DB", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
