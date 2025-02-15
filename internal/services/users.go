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
	ErrInvalidUserType    = errors.New("invalid user type")
)

type UserService struct {
	admins *AdminService
}

const (
	TypeAdmin    = "admin"
	TypeCustomer = "customer"
)

func newUserService(admin *AdminService) *UserService {
	return &UserService{
		admins: admin,
	}
}

func (a *UserService) Signup(fullName, email, password, userType string) (Token, error) {
	switch userType {
	case TypeAdmin:
		return a.admins.Signup(fullName, email, password)
	default:
		return Token{}, fmt.Errorf("signup: %w", ErrInvalidUserType)
	}
}

func (a *UserService) Login(email, password, userType string) (Token, error) {
	switch userType {
	case TypeAdmin:
		return a.admins.Login(email, password)
	default:
		return Token{}, fmt.Errorf("login: %w", ErrInvalidUserType)
	}
}

func (a *UserService) Logout(token Token, userType string) error {
	switch userType {
	case TypeAdmin:
		return a.admins.Logout(token)
	default:
		return fmt.Errorf("logout: %w", ErrInvalidUserType)
	}
}

func (a *UserService) GetAuthToken(tokenText, userType string) (repository.User, error) {
	switch userType {
	case TypeAdmin:
		return a.admins.GetAdminbyAuthToken(tokenText)
	default:
		return nil, fmt.Errorf("getByToken: %w", ErrInvalidUserType)
	}
}
