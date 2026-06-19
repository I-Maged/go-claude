package login

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInsufficient  = errors.New("insufficient information")
	ErrNotFound      = errors.New("not found")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrAlreadyExists = errors.New("already exist")
)

type UserStore struct {
	users map[string]string
}

func NewUserStore() *UserStore {
	return &UserStore{users: make(map[string]string)}
}

func (u *UserStore) Create(name, password string) error {
	if strings.TrimSpace(name) == "" || strings.TrimSpace(password) == "" {
		return ErrInsufficient
	}

	if _, exists := u.users[name]; exists {
		return ErrAlreadyExists
	}

	u.users[name] = password

	return nil
}

func (u *UserStore) Login(name, password string) error {
	stored, exists := u.users[name]
	if !exists {
		return ErrNotFound
	}
	if stored != password {
		return ErrUnauthorized
	}
	return nil
}

func (u *UserStore) TestLogin() {
	u.Create("Maged", "secret123")

	err := u.Login("nobody", "secret123")
	if errors.Is(err, ErrNotFound) {
		fmt.Println("User does not exist")
	}

	err = u.Login("Maged", "WrongPassword")
	if errors.Is(err, ErrUnauthorized) {
		fmt.Println("Wrong Password")
	}

	err = u.Create("Maged", "newpass")
	if errors.Is(err, ErrAlreadyExists) {
		fmt.Println("User already registered")
	}

	err = u.Create("Maged", "  ")
	if errors.Is(err, ErrInsufficient) {
		// fmt.Println("name and password cannot be empty")
		fmt.Println(fmt.Errorf("Error: %w", ErrInsufficient))
	}
}
