package app

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
)

// User is an application user's interface
type User interface {
	ID() string

	Name() string
	SetName(string)

	Password() string
	SetPassword(string)
	CheckPassword(string) bool

	Salt() []byte
}

type defaultUser struct {
	id       string
	name     string
	password string
	salt     []byte
}

// NewUser creates a new default user with a valid UUID, name and password
func NewUser(name, password string) (User, error) {
	id, err := newUUID()
	if err != nil {
		return nil, err
	}

	du := &defaultUser{id: id, name: name}
	du.SetPassword(password)

	return du, nil
}

func (du defaultUser) ID() string {
	return du.id
}

func (du defaultUser) Name() string {
	return du.name
}
func (du defaultUser) SetName(name string) {
	du.name = name
}

func (du defaultUser) Password() string {
	return du.password
}
func (du defaultUser) SetPassword(password string) {
	saltBase := string(rand.Intn(10000000))[12:]

	h := sha256.New()
	h.Write([]byte(saltBase))
	salt := h.Sum(nil)

	du.password = passwordHash(password, salt)
	du.salt = salt
}

func (du defaultUser) CheckPassword(password string) bool {
	newPasswordHash := passwordHash(password, du.salt)
	return du.password == newPasswordHash
}

func (du defaultUser) Salt () []byte {
	return du.salt
}

func passwordHash(password string, salt []byte) string {
	saltedPassword := fmt.Sprintf("%s%s", salt, password)

	h := sha256.New()
	h.Write([]byte(saltedPassword))
	passwordHash := h.Sum(nil)

	return fmt.Sprintf("%s", passwordHash)
}
