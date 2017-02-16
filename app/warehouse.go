package app

import "sync"

type warehouse struct {
	stock map[string]Stock

	sync.RWMutex
}

func NewWarehouse() *warehouse {
	return &warehouse{}
}

func (wh *warehouse) Add(item Stock) string {
	wh.Lock()
	item.Create()
	wh.Unlock()

	return item.Id()
}

func (wh *warehouse) Get(id string) (item Stock, ok bool) {
	wh.RLock()
	item, ok = item.Read(id)
	wh.RUnlock()

	return
}

// Removes the item with the given id from the warehouse.
func (wh *warehouse) Remove(id string) {
	wh.Lock()
	stock := &defaultStock{id: id}
	stock.Delete()
	wh.Unlock()
}

// Returns a map with the items in the warehouse with ids as keys and stock items as their values.
func (wh *warehouse) Stock() (stock map[string]Stock) {
	stock = make(map[string]Stock)
	wh.RLock()
	// 	for key, value := range wh.stock {
	// 		stock[key] = value
	// 	}
	wh.RUnlock()

	return
}

func (wh *warehouse) Size() (size int) {
	wh.RLock()
	// 	size = len(wh.stock)
	wh.RUnlock()

	return
}
