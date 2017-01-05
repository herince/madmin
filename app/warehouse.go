package madmin

type Warehouse struct {
	stock []Stock
}

func (w *Warehouse) Add(item *Stock) error {}
func (w *Warehouse) Remove(item *Stock) error {}