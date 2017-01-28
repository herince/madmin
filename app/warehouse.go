package app

type Warehouse struct {
	Stock map[string]Stock
}

func (w *Warehouse) AddItem(item *Stock) error {
	return nil
}
func (w *Warehouse) GetItem(id string) (Stock, error) {
	return w.Stock[id], nil
}
func (w *Warehouse) RemoveItem(item *Stock) error {
	return nil
}
