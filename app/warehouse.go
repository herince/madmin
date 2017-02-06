package app

type Warehouse struct {
	Stock map[string]Stock
}

func (w *Warehouse) Add(item Stock) {
	w.Stock[item.Id()] = item
}
func (w *Warehouse) Remove(item *Stock) error {
	return nil
}
