package server

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name           string
	HashedPassword []byte
}

type Users map[string]User

func NewUser(name, pass string) (User, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	return User{name, h}, nil
}
