package server

import (
	"golang.org/x/crypto/bcrypt"
)

// User represents an user.
// Each user is identified by an username and a bcrypt hashed password.
type User struct {
	Name           string
	HashedPassword []byte
}

// Users represents a collection of users.
type Users map[string]User

// NewUser creates a new User with the given username and an hashed password.
func NewUser(name, pass string) (User, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	return User{name, h}, nil
}
