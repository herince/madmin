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

func (wh *warehouse) Add(item Stock) string {
	wh.createStock(item)

	return item.Id()
}

func (wh *warehouse) Get(id string) (item Stock, ok bool) {
	item, ok = wh.readStock(id)
	return
}

// Removes the item with the given id from the warehouse.
func (wh *warehouse) Remove(id string) {
	wh.deleteStock(id)
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
	query := `
	SELECT COUNT(*) FROM warehouse;
	`

	rows, err := wh.database.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if err := rows.Scan(&size); err != nil {
		panic(err)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	return
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
func (wh *warehouse) createStock(s Stock) {}

// read from DB
func (wh *warehouse) readStock(id string) (Stock, bool) {
	return nil, false
}

// update in DB
func (wh *warehouse) updateStock(s Stock) {}

// remove from DB
func (wh *warehouse) deleteStock(id string) {}

// Database methods for distributors
// insert in DB
func (wh *warehouse) createDistributor(d distributor) {}

// read from DB
func (wh *warehouse) readDistributor(id string) (distributor, bool) {
	return distributor{}, false
}

// update in DB
func (wh *warehouse) updateDistributor(d distributor) {}

// remove from DB
func (wh *warehouse) deleteDistributor(id string) {}
