package app

// User is an application user's interface
type User interface {
	ID() string

	Name() string
	SetName(string)

	Password() string
	SetPassword(string)
}

type defaultUser struct {
	id       string
	name     string
	password string
}

// NewUser creates a new default user with a valid UUID, name and password
func NewUser(name, password string) (User, error) {
	id, err := newUUID()
	if err != nil {
		return nil, err
	}

	return &defaultUser{id: id, name: name, password: password}, nil
}

func (du defaultUser) ID() (id string) {
	return
}

func (du defaultUser) Name() (name string) {
	return
}
func (du defaultUser) SetName(name string) {}

func (du defaultUser) Password() (password string) {
	return
}
func (du defaultUser) SetPassword(password string) {}
