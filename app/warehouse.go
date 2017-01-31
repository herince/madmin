package app

type Warehouse struct {
	Stock map[string]Stock
}

func (w *Warehouse) Add(item Stock) {
	w.Stock[item.Id()] = item
}
func (w *Warehouse) Get(id string) (Stock, bool) {
	item, ok := w.Stock[id]
	return item, ok
}
func (w *Warehouse) Remove(item *Stock) error {
	return nil
}
