package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

const (
	TestDBSource = "postgresql://root:password@localhost:5432/test_db?sslmode=disable"
)

// TestMain Creates a test DB for test functions
func TestMain(m *testing.M) {
	var err error

	testDB, err = sql.Open("postgres", TestDBSource)
	if err != nil {
		log.Fatal("cannot connect to DB", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
