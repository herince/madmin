package app

type distributor struct {
	id   string
	name string
}

func NewDistributor(name string) (*distributor, error) {
	id, err := newUUID()
	if err != nil {
		return nil, err
	}

	return &distributor{id: id, name: name}, nil
}

func (d distributor) Name() string {
	return d.name
}
func (d distributor) SetName(name string) {
	d.name = name
}

/*
 * more TBD
 */
