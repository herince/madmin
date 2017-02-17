package app

import "database/sql"

type warehouse struct {
	database *sql.DB
}

func NewWarehouse(db *sql.DB) *warehouse {
	wh := &warehouse{database: db}

	wh.initStockTable()
	wh.initDistributorsTable()

	return wh
}

func (wh *warehouse) initStockTable() {
	stock_table := `
	CREATE TABLE IF NOT EXISTS warehouse(
		id BLOB NOT NULL PRIMARY KEY,
		type TEXT NOT NULL,
		name TEXT,
		min_quantity REAL,
		expiration_date TEXT,
		distributor_id BLOB,
		FOREIGN KEY (distributor_id) REFERENCES distributors (Id)
	);
	`

	_, err := wh.database.Exec(stock_table)
	if err != nil {
		panic(err)
	}
}

func (wh *warehouse) initDistributorsTable() {
	distributors_table := `
	CREATE TABLE IF NOT EXISTS distributors(
		id BLOB NOT NULL PRIMARY KEY,
		name TEXT
	);
	`

	_, err := wh.database.Exec(distributors_table)
	if err != nil {
		panic(err)
	}
}

// Database methods for stock items
// insert in DB
func (wh *warehouse) createItem(item Stock) {}

// read from DB
func (wh *warehouse) readItem(id string) (Stock, bool) {
	return nil, false
}

// update in DB
func (wh *warehouse) updateItem(item Stock) {}

// remove from DB
func (wh *warehouse) deleteItem(id string) {}

func (wh *warehouse) Add(item Stock) string {
	wh.createItem(item)

	return item.Id()
}

func (wh *warehouse) Get(id string) (item Stock, ok bool) {
	item, ok = wh.readItem(id)
	return
}

// Removes the item with the given id from the warehouse.
func (wh *warehouse) Remove(id string) {
	wh.deleteItem(id)
}

// Returns a map with the items in the warehouse with ids as keys and stock items as their values.
func (wh *warehouse) Stock() (stock map[string]Stock) {
	stock = make(map[string]Stock)

	// return stock items in db as a Go map structure
	// something like...
	// 	for key, value := range wh.stock {
	// 		stock[key] = value
	// 	}

	return
}

func (wh *warehouse) Size() (size int) {
	// 	return number of stock items in db

	return
}
