package app

type user struct {
	name     string
	password string
}

func NewUser(name, password string) (user, error) {
	return user{name: name, password: password}, nil
}
