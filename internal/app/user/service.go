package user

import (
	"errors"
)

var (
	ErrEmailAlreadyTaken  = errors.New("email already taken")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAuthToken   = errors.New("invalid auth token")
	ErrInvalidUserType    = errors.New("invalid user type")
)

type UserService struct {
}

type Token struct {
}

func (s *UserService) Signup(user User) (Token, error) {

	return Token{}, nil
}
func (s *UserService) Login(username, password string) (Token, error) {
	return Token{}, nil
}
