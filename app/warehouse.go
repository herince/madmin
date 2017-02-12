package app

import "sync"

type warehouse struct {
	stock map[string]Stock

	sync.RWMutex
}

func NewWarehouse() *warehouse {
	return &warehouse{stock: make(map[string]Stock)}
}

func (wh *warehouse) Add(item Stock) string {
	wh.Lock()
	wh.stock[item.Id()] = item
	wh.Unlock()

	return item.Id()
}
func (wh *warehouse) Get(id string) (item Stock, ok bool) {
	wh.RLock()
	item, ok = wh.stock[id]
	wh.RUnlock()

	return
}
func (wh *warehouse) Remove(item Stock) (id string) {
	wh.Lock()
	id = item.Id()
	delete(wh.stock, id)
	wh.Unlock()

	return
}

/*
 * Returns a map with the items in the warehouse with ids as keys and stock items as their values.
 */
func (wh *warehouse) Stock() (stock map[string]Stock) {
	stock = make(map[string]Stock)
	wh.RLock()
	for key, value := range wh.stock {
		stock[key] = value
	}
	wh.RUnlock()

	return
}

func (wh *warehouse) Size() int {
	wh.RLock()
	defer wh.RUnlock()

	return len(wh.stock)
}
