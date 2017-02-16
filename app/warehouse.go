package app

import "database/sql"

type warehouse struct {
	database *sql.DB
}

func NewWarehouse(database *sql.DB) *warehouse {
	return &warehouse{database: database}
}

// Database methods for stock items
// insert in DB
func CreateItem(item Stock) {}

// read from DB
func ReadItem(id string) (Stock, bool) {
	return nil, false
}

// update in DB
func UpdateItem(item Stock) {}

// remove from DB
func DeleteItem(id string) {}

func (wh *warehouse) Add(item Stock) string {
	item.Create(wh.database)

	return item.Id()
}

func (wh *warehouse) Get(id string) (item Stock, ok bool) {
	item, ok = item.Read(wh.database, id)

	return
}

// Removes the item with the given id from the warehouse.
func (wh *warehouse) Remove(id string) {
	stock := &defaultStock{id: id}
	stock.Delete(wh.database)
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
