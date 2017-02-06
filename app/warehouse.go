package app

type Warehouse struct {
	Stock map[string]Stock
}

func (wh *Warehouse) Add(item Stock) string {
	wh.Stock[item.Id()] = item
	return item.Id()
}
func (wh *Warehouse) Remove(item Stock) (string) {
	id := item.Id()
	delete(wh.Stock, id)
	return id
}
