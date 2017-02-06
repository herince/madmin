package app

type Warehouse struct {
	Stock map[string]Stock
}

func (w *Warehouse) Add(item Stock) string{
	w.Stock[item.Id()] = item
	return item.Id()
}
func (w *Warehouse) Remove(item Stock) (string, error) {
	return item.Id(), nil
}
