package services

import (
	"AhmadAbdelrazik/arbun/internal/repository"
	"errors"
	"fmt"
)

var (
	ErrEmailAlreadyTaken  = errors.New("email already taken")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAuthToken   = errors.New("invalid auth token")
)

type UserService struct {
	admins *AdminService
}

const (
	TypeAdmin    = "admin"
	TypeCustomer = "customer"
)

func newUserService() *UserService {
	return &UserService{
		admins: newAdminService(),
	}
}

func (a *UserService) Signup(fullName, email, password, userType string) (Token, error) {
	switch userType {
	case TypeAdmin:
		return a.admins.Signup(fullName, email, password)
	default:
		return Token{}, fmt.Errorf("invalid type")
	}
}

func (a *UserService) Login(email, password, userType string) (Token, error) {
	switch userType {
	case TypeAdmin:
		return a.admins.Login(email, password)
	default:
		return Token{}, fmt.Errorf("invalid type")
	}
}

func (a *UserService) Logout(token Token, userType string) error {
	switch userType {
	case TypeAdmin:
		return a.admins.Logout(token)
	default:
		return fmt.Errorf("invalid type")
	}
}

func (a *UserService) GetByToken(tokenText, scope, userType string) (repository.User, error) {
	switch userType {
	case TypeAdmin:
		return a.admins.GetAdminByToken(tokenText, scope)
	default:
		return nil, fmt.Errorf("invalid type")
	}
}
