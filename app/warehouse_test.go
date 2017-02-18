package app

import (
	"database/sql"
	"os"
	"testing"
)

func cleanupDatabase(t *testing.T, database *sql.DB, dbPath string) {
	database.Close()
	os.Remove(dbPath)
}

func checkIfTableExists(t *testing.T, database *sql.DB, table string) {
	var name string

	query := `
		SELECT
			name
		FROM
			sqlite_master
		WHERE
			type='table'
		AND
			name=?
	`
	err := database.QueryRow(query, table).Scan(&name)

	if err != nil {
		t.Fatalf("%s", err)
	}

	if name != table {
		t.Fatalf("table %s does not exist in DB", table)
	}
}

func TestWarehouse(t *testing.T) {
	dbPath := "./test_database.sqlite"
	db := newDB(dbPath)

	t.Run("NewWarehouse_DoesNotPanic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf(`NewWarehouse panics`)
			}
		}()

		wh := NewWarehouse(db)
		checkNewItemCreating(t, wh, nil)
	})
	t.Run("NewWarehouse_CreatesTables", func(t *testing.T) {
		_ = NewWarehouse(db)
		checkIfTableExists(t, db, "warehouse")
		checkIfTableExists(t, db, "distributors")
	})
	t.Run("CreateStock", func(t *testing.T) {
		wh := NewWarehouse(db)

		item, _ := defaultExpirableStockItem(MEDICINE)
		wh.CreateStock(item)

		read, ok := wh.ReadStock(item.ID())

		if !ok {
			t.Fatalf(`new item not found in database`)
		}
		if compareStock(read, item) == false {
			t.Fatalf(`
				Read item is different from expected.
				Expected %+v, got %+v.`,
				item,
				read)
		}
	})
	t.Run("ReadStock", func(t *testing.T) {
		wh := NewWarehouse(db)

		item, _ := defaultExpirableStockItem(MEDICINE)
		wh.CreateStock(item)

		read, ok := wh.ReadStock(item.ID())
		checkNewItemCreatingWithOKStatus(t, read, ok)

		if !ok {
			t.Fatalf(`new item not found in database`)
		}
		if compareStock(read, item) == false {
			t.Fatalf(`
				Read item is different from expected.
				Expected %+v, got %+v.`,
				item,
				read)
		}

		fakeID := "I am a fake ID!"

		read, ok = wh.ReadStock(fakeID)
		checkNewItemCreatingWithOKStatus(t, read, ok)

		if ok || read != nil {
			t.Fatalf(`reads invalid item from database ???`)
		}
	})
	t.Run("UpdateStock", func(t *testing.T) {
		wh := NewWarehouse(db)

		item, _ := defaultExpirableStockItem(MEDICINE)
		wh.CreateStock(item)

		item.SetName("Aspirin")
		wh.UpdateStock(item)

		read, ok := wh.ReadStock(item.ID())
		if !ok || read == nil {
			t.Fatalf(`cannot read valid item from database`)
		}

		if compareStock(read, item) == false {
			t.Fatalf(`
				Read item is different from expected.
				Expected %+v, got %+v.`,
				item,
				read)
		}
	})
	t.Run("DeleteStock", func(t *testing.T) {
		wh := NewWarehouse(db)

		item, _ := defaultExpirableStockItem(MEDICINE)
		wh.CreateStock(item)

		wh.DeleteStock(item.ID())

		read, ok := wh.ReadStock(item.ID())
		if ok || read != nil {
			t.Fatalf(`reads deleted item from database ???`)
		}
	})

	cleanupDatabase(t, db, dbPath)
}

// TODO: add tests for "distributors" table
