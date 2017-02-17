package app

import (
	"database/sql"
	// 	"github.com/shopspring/decimal"
)

// Warehouse is a warehouse interface.
// A warehouse must manage two datasets -
// one with the existing stock items and one with the stock items' distributors.
type Warehouse interface {
	CreateStock(Stock) string
	ReadStock(string) (Stock, bool)
	UpdateStock(Stock)
	DeleteStock(string)

	CreateDistributor(Distributor)
	ReadDistributor(string) (Distributor, bool)
	UpdateDistributor(Distributor)
	DeleteDistributor(string)

	// Stock() returns a map with the ids of the current stock items in the DB,
	// mapped to the corresponding stock items
	Stock() map[string]Stock

	// Size() returns number of unique stock items in DB
	// TODO: should return number of all stock items in DB
	Size() int
}

type dafaultWarehouse struct {
	database *sql.DB
}

// NewWarehouse creates a warehouse that holds the stock items'
// and distriubutors' data in two separate sqlite3 tables inside the db
// that is passed as an argument.
func NewWarehouse(db *sql.DB) Warehouse {
	wh := &dafaultWarehouse{database: db}

	wh.initStockTable()
	wh.initDistributorsTable()

	return wh
}

func (wh *dafaultWarehouse) initStockTable() {
	stockTable := `
	CREATE TABLE IF NOT EXISTS warehouse(
		id BLOB NOT NULL PRIMARY KEY,
		type TEXT NOT NULL,
		name TEXT,
		min_quantity TEXT,
		expiration_date INTEGER,
		distributor_id BLOB,
		FOREIGN KEY (distributor_id) REFERENCES distributors (Id)
	);
	`

	_, err := wh.database.Exec(stockTable)
	if err != nil {
		panic(err)
	}
}

func (wh *dafaultWarehouse) initDistributorsTable() {
	distributorsTable := `
	CREATE TABLE IF NOT EXISTS distributors(
		id BLOB NOT NULL PRIMARY KEY,
		name TEXT
	);
	`

	_, err := wh.database.Exec(distributorsTable)
	if err != nil {
		panic(err)
	}
}

// Database CRUD methods for stock items
// insert in DB
func (wh *dafaultWarehouse) CreateStock(item Stock) string {
	transaction, err := wh.database.Begin()
	if err != nil {
		panic(err)
	}

	stmt, err := transaction.Prepare(`
		INSERT INTO warehouse (id, type, name, min_quantity, expiration_date, distributor_id) VALUES(?, ?, ?, ?, ?, ?)
	`)
	defer stmt.Close()
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(
		item.ID(),
		item.Type(),
		item.Name(),
		item.MinQuantity(),
		item.ExpirationDate().UnixNano(),
		item.DistributorID())
	transaction.Commit()

	return item.ID()
}

// read from DB
func (wh *dafaultWarehouse) ReadStock(id string) (item Stock, ok bool) {
	stmt, err := wh.database.Prepare(`
	SELECT
		type, name, min_quantity, expiration_date, distributor_id FROM warehouse
	WHERE
		id = ?
	`)
	defer stmt.Close()
	if err != nil {
		panic(err)
	}

	var (
		dto           NewStockDTO
		distributorID string
	)
	err = stmt.QueryRow(id).Scan(&dto.Type, &dto.Name, &dto.MinQuantity, &dto.ExpirationDate, &distributorID)
	switch {
	case err == sql.ErrNoRows:
		return nil, false
	case err != nil:
		panic(err)
	}

	return
}

// update in DB
func (wh *dafaultWarehouse) UpdateStock(s Stock) {

}

// remove from DB
func (wh *dafaultWarehouse) DeleteStock(id string) {}

// Database CRUD methods for distributors
// insert in DB
func (wh *dafaultWarehouse) CreateDistributor(d Distributor) {}

// read from DB
func (wh *dafaultWarehouse) ReadDistributor(id string) (Distributor, bool) {
	return defaultDistributor{}, false
}

// update in DB
func (wh *dafaultWarehouse) UpdateDistributor(d Distributor) {}

// remove from DB
func (wh *dafaultWarehouse) DeleteDistributor(id string) {}

// Returns a map with the items in the warehouse with ids as keys and stock items as their values.
func (wh *dafaultWarehouse) Stock() (stock map[string]Stock) {
	stock = make(map[string]Stock)

	// return stock items in db as a Go map structure
	// something like...
	// 	for key, value := range wh.stock {
	// 		stock[key] = value
	// 	}

	return
}

func (wh *dafaultWarehouse) Size() (size int) {
	// 	return number of stock items in db
	// 	query := `
	// 	SELECT COUNT(*) FROM warehouse;
	// 	`
	//
	// 	rows, err := wh.database.Query(query)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	defer rows.Close()
	//
	// 	if err := rows.Scan(&size); err != nil {
	// 		panic(err)
	// 	}
	// 	if err := rows.Err(); err != nil {
	// 		panic(err)
	// 	}
	return
}
