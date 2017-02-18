package app

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // needed for for working with sqlite3
)

func newDB(dbPath string) (db *sql.DB) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return
}
