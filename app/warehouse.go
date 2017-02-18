package app

import (
	"database/sql"

	// used for decimal operations
	_ "github.com/shopspring/decimal"
)

// Warehouse is a warehouse interface.
// A warehouse must manage two datasets -
// one with the existing stock items and one with the stock items' distributors.
type Warehouse interface {
	CreateStock(Stock)
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
		quantity NUMERIC NOT NULL,
		min_quantity NUMERIC,
		expiration_date DATETIME,
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
	CREATE TABLE IF NOT EXISTS
		distributors (
			id BLOB NOT NULL PRIMARY KEY,
			name TEXT);
	`
	_, err := wh.database.Exec(distributorsTable)
	if err != nil {
		panic(err)
	}
}

// Database CRUD methods for stock items
// insert in DB
func (wh *dafaultWarehouse) CreateStock(item Stock) {
	stmt, err := wh.database.Prepare(`
		INSERT INTO
			warehouse (
				id,
				type,
				name,
				quantity,
				min_quantity,
				expiration_date,
				distributor_id)
		VALUES(?, ?, ?, ?, ?, ?, ?)
	`)
	defer stmt.Close()
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(
		item.ID(),
		item.Type(),
		item.Name(),
		item.Quantity().String(),
		item.MinQuantity().String(),
		item.ExpirationDate(),
		item.DistributorID())
	if err != nil {
		panic(err)
	}
}

// read from DB
func (wh *dafaultWarehouse) ReadStock(id string) (item Stock, ok bool) {
	stmt, err := wh.database.Prepare(`
	SELECT
		type,
		name,
		quantity,
		min_quantity,
		expiration_date,
		distributor_id
	FROM
		warehouse
	WHERE
		id = ?
	`)
	defer stmt.Close()
	if err != nil {
		panic(err)
	}

	var (
		stockItem = defaultStock{id: id}
		sType     int8
	)
	err = stmt.QueryRow(id).Scan(
		&sType,
		&stockItem.name,
		&stockItem.quantity,
		&stockItem.minQuantity,
		&stockItem.expirationDate,
		&stockItem.distributorID)
	switch {
	case err == sql.ErrNoRows:
		return nil, false
	case err != nil:
		panic(err)
	}

	switch stockType(sType) {
	case MEDICINE:
		return &medicine{stockItem}, true
	case FEED:
		return &feed{stockItem}, true
	case ACCESSORY:
		return &accessory{stockItem}, true
	default:
		panic("invalid stock type in DB record")
	}
}

// update in DB
func (wh *dafaultWarehouse) UpdateStock(item Stock) {
	stmt, err := wh.database.Prepare(`
	UPDATE
		warehouse
	SET
		type = ?,
		name = ?,
		quantity = ?,
		min_quantity = ?,
		expiration_date = ?,
		distributor_id = ?
	WHERE
		id = ?
	`)
	defer stmt.Close()
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(
		item.Type(),
		item.Name(),
		item.Quantity().String(),
		item.MinQuantity().String(),
		item.ExpirationDate(),
		item.DistributorID(),
		item.ID())
	if err != nil {
		panic(err)
	}
}

// remove from DB
func (wh *dafaultWarehouse) DeleteStock(id string) {
	stmt, err := wh.database.Prepare(`
		DELETE FROM
			warehouse
		WHERE
			id = ?
	`)
	defer stmt.Close()
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		panic(err)
	}
}

// Database CRUD methods for distributors
// insert in DB
func (wh *dafaultWarehouse) CreateDistributor(d Distributor) {
	stmt, err := wh.database.Prepare(`
		INSERT INTO
			distributors (
				id,
				name)
		VALUES (?, ?)
	`)
	defer stmt.Close()
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(
		d.ID(),
		d.Name())
	if err != nil {
		panic(err)
	}
}

// read from DB
func (wh *dafaultWarehouse) ReadDistributor(id string) (Distributor, bool) {
	stmt, err := wh.database.Prepare(`
	SELECT
		name
	FROM
		distributors
	WHERE
		id = ?
	`)
	defer stmt.Close()
	if err != nil {
		panic(err)
	}

	var d = defaultDistributor{id: id}

	err = stmt.QueryRow(id).Scan(&d.name)
	switch {
	case err == sql.ErrNoRows:
		return nil, false
	case err != nil:
		panic(err)
	}

	return d, true
}

// update in DB
func (wh *dafaultWarehouse) UpdateDistributor(d Distributor) {
	stmt, err := wh.database.Prepare(`
	UPDATE
		distributors
	SET
		name = ?,
	WHERE
		id = ?
	`)
	defer stmt.Close()
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(
		d.Name(),
		d.ID())
	if err != nil {
		panic(err)
	}
}

// remove from DB
func (wh *dafaultWarehouse) DeleteDistributor(id string) {
	stmt, err := wh.database.Prepare(`
		DELETE FROM
			distributors
		WHERE
			id = ?
	`)
	defer stmt.Close()
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		panic(err)
	}
}

// Returns a map with the items in the warehouse with ids as keys and stock items as their values.
func (wh *dafaultWarehouse) Stock() (stock map[string]Stock) {
	stock = make(map[string]Stock)

	query := `
		SELECT
			id,
			type,
			name,
			quantity,
			min_quantity,
			expiration_date,
			distributor_id
		FROM
			warehouse
	`

	rows, err := wh.database.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			stockItem = defaultStock{}
			sType     int8
		)
		err = rows.Scan(
			&stockItem.id,
			&sType,
			&stockItem.name,
			&stockItem.quantity,
			&stockItem.minQuantity,
			&stockItem.expirationDate,
			&stockItem.distributorID)

		if err != nil {
			panic(err)
		}

		switch stockType(sType) {
		case MEDICINE:
			stock[stockItem.ID()] = &medicine{stockItem}
		case FEED:
			stock[stockItem.ID()] = &feed{stockItem}
		case ACCESSORY:
			stock[stockItem.ID()] = &accessory{stockItem}
		default:
			panic("invalid stock type in DB record")
		}
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return
}

// Size returns the number of rows in the warehouse table in the DB
func (wh *dafaultWarehouse) Size() (size int) {
	stmt, err := wh.database.Prepare("SELECT COUNT(*) FROM warehouse;")
	if err != nil {
		panic(err)
	}

	err = stmt.QueryRow().Scan(&size)
	switch {
	case err == sql.ErrNoRows:
		return 0
	case err != nil:
		panic(err)
	}

	return
}
