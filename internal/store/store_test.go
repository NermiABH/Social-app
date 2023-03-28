package store_test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
)

var (
	databaseURL string
	TestDB      = func(t *testing.T, databaseURL string) (*sql.DB, func(...string)) {
		t.Helper()
		db, err := sql.Open("postgres", databaseURL)
		if err != nil {
			log.Fatalln(err)
		}
		if err = db.Ping(); err != nil {
			log.Fatalln(err)
		}
		return db, func(tables ...string) {
			if len(tables) > 0 {
				db.Exec(fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ",")))
			}
			db.Close()
		}
	}
)

func TestMain(m *testing.M) {
	databaseURL = os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "user=postgres host=localhost password=pusinu48 dbname=socialapptest sslmode=disable"
	}
	os.Exit(m.Run())
}
