package app

import "sync"

type warehouse struct {}

func NewWarehouse() *warehouse {
	return &warehouse{}
}

func (wh *warehouse) Add(item Stock) string {
	item.Create()

	return item.Id()
}

func (wh *warehouse) Get(id string) (item Stock, ok bool) {
	item, ok = item.Read(id)

	return
}

// Removes the item with the given id from the warehouse.
func (wh *warehouse) Remove(id string) {
	stock := &defaultStock{id: id}
	stock.Delete()
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
