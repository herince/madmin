package app

import (
	"database/sql"

	// 	"github.com/shopspring/decimal"
	_ "github.com/mattn/go-sqlite3"
)

func newDB() (db *sql.DB) {
	dbPath := "./database/database.sqlite"

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return
}

func truncateDB() error {
	return nil
}

func dropDB() error {
	return nil
}
