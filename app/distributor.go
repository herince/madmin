package app

// Distributor is a simple type interface for data about a stock distributor
type Distributor interface {
	ID() string

	Name() string
	SetName(string)
}

type defaultDistributor struct {
	id   string
	name string
}

// NewDistributor creates a new distributor with a UUID and a given name
func NewDistributor(name string) (Distributor, error) {
	id, err := newUUID()
	if err != nil {
		return nil, err
	}

	return &defaultDistributor{id: id, name: name}, nil
}

func (d defaultDistributor) ID() string {
	return d.id
}

func (d defaultDistributor) Name() string {
	return d.name
}
func (d defaultDistributor) SetName(name string) {
	d.name = name
}

/*
 * more TBD
 */
