package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/burakkarasel/Bank-App/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

// TestMain Creates a test DB for test functions
func TestMain(m *testing.M) {
	var err error

	config, err := util.LoadConfig("../..")

	if err != nil {
		log.Fatal("cannot get env variables:", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to DB", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
