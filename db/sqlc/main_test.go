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
var testDB *sql.DB

// TestMain Creates a test DB for test functions
func TestMain(m *testing.M) {
	var err error

	testDB, err = sql.Open(dsn.TestDBDriver, dsn.TestDBSource)
	if err != nil {
		log.Fatal("cannot connect to DB", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
