package madmin

type Warehouse struct {
	stock map[string]Stock
}

func (w *Warehouse) Add(item *Stock) error {
	return nil
}
func (w *Warehouse) Remove(item *Stock) error {
	return nil
}