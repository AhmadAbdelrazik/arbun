package services

import (
	"AhmadAbdelrazik/arbun/internal/domain/user"
	"AhmadAbdelrazik/arbun/internal/models"
	"errors"
	"fmt"
	"time"
)

var (
	ErrEmailAlreadyTaken  = errors.New("email already taken")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAuthToken   = errors.New("invalid auth token")
	ErrInvalidUserType    = errors.New("invalid user type")
)

type UserService struct {
	admins    *AdminService
	customers *CustomerService
}

type Token struct {
	Plaintext  string
	ExpiryTime time.Time
}

const (
	TypeAdmin    = "admin"
	TypeCustomer = "customer"
)

func newUserService(models *models.Model) *UserService {
	return &UserService{
		admins:    newAdminService(models),
		customers: newCustomerService(models),
	}
}

func (a *UserService) Signup(fullName, email, password, userType string) (Token, error) {
	switch userType {
	case TypeAdmin:
		return a.admins.Signup(fullName, email, password)
	case TypeCustomer:
		return a.customers.Signup(fullName, email, password)
	default:
		return Token{}, fmt.Errorf("signup: %w", ErrInvalidUserType)
	}
}

func (a *UserService) Login(email, password, userType string) (Token, error) {
	switch userType {
	case TypeAdmin:
		return a.admins.Login(email, password)
	case TypeCustomer:
		return a.customers.Login(email, password)
	default:
		return Token{}, fmt.Errorf("login: %w", ErrInvalidUserType)
	}
}

func (a *UserService) Logout(token Token, userType string) error {
	switch userType {
	case TypeAdmin:
		return a.admins.Logout(token)
	case TypeCustomer:
		return a.customers.Logout(token)
	default:
		return fmt.Errorf("logout: %w", ErrInvalidUserType)
	}
}

func (a *UserService) GetAuthToken(tokenText, userType string) (user.IUser, error) {
	switch userType {
	case TypeAdmin:
		return a.admins.GetAdminbyAuthToken(tokenText)
	case TypeCustomer:
		return a.customers.GetCustomerbyAuthToken(tokenText)
	default:
		return nil, fmt.Errorf("getByToken: %w", ErrInvalidUserType)
	}
}
