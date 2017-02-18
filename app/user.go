package app

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"math/rand"
)

// User is an application user's interface
type User interface {
	ID() string

	Name() string
	SetName(string)

	Password() string
	SetPassword(string) error
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

	if len(name) == 0 {
		return nil, errors.New("cannot set empty string as name")
	}

	du := &defaultUser{id: id, name: name}
	err = du.SetPassword(password)
	if err != nil {
		return nil, err
	}

	return du, nil
}

func (du *defaultUser) ID() string {
	return du.id
}

func (du *defaultUser) Name() string {
	return du.name
}
func (du *defaultUser) SetName(name string) {
	du.name = name
}

func (du *defaultUser) Password() string {
	return du.password
}
func (du *defaultUser) SetPassword(password string) error {
	if len(password) == 0 {
		return errors.New("cannot set empty string as password")
	}
	saltBase := fmt.Sprintf("%d", rand.Intn(10000000))

	h := sha256.New()
	h.Write([]byte(saltBase))
	salt := h.Sum(nil)

	du.password = passwordHash(password, salt)
	du.salt = salt

	return nil
}

func (du *defaultUser) CheckPassword(password string) bool {
	newPasswordHash := passwordHash(password, du.salt)
	return du.password == newPasswordHash
}

func (du *defaultUser) Salt() []byte {
	return du.salt
}

func passwordHash(password string, salt []byte) string {
	saltedPassword := fmt.Sprintf("%s%s", salt, password)

	h := sha256.New()
	h.Write([]byte(saltedPassword))
	pwdHash := h.Sum(nil)

	return fmt.Sprintf("%s", pwdHash)
}
