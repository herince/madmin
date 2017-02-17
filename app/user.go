package app

type user struct {
	id       string
	name     string
	password string
}

func NewUser(name, password string) (user, error) {
	id, err := newUUID()
	if err != nil {
		return user{}, err
	}
	return user{id: id, name: name, password: password}, nil
}
