package app

import (
	"testing"
	"os"
	"database/sql"
)

func setupWarehouse(t *testing.T, dbPath string) (Warehouse, *sql.DB) {
	database := newDB(dbPath)

	wh := NewWarehouse(database)
	checkNewItemCreating(t, wh, nil)

	return wh, database
}
func cleanupWarehouse(t *testing.T, database *sql.DB, dbPath string) {
	database.Close()
	os.Remove(dbPath)
}

func TestNewWarehouse(t *testing.T) {
	dbPath := "./test_database.sqlite"

	_, db := setupWarehouse(t, dbPath)

	_ = `SELECT
			name
		FROM
			sqlite_master
		WHERE
			type='table'
		AND
			name='warehouse';`

	// test if "warehouse" and "distributors" tables exist

	cleanupWarehouse(t, db, dbPath)
}