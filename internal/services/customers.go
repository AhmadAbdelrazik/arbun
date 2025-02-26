package services

import (
	"AhmadAbdelrazik/arbun/internal/domain"
	"AhmadAbdelrazik/arbun/internal/models"
	"errors"
	"fmt"
	"time"
)

type CustomerService struct {
	models *models.Model
}

func newCustomerService(models *models.Model) *CustomerService {
	return &CustomerService{
		models: models,
	}
}

func (a *CustomerService) Signup(fullName, email, password string) (Token, error) {
	// 1. user provide credentials
	var newCustomer domain.Customer
	newCustomer.FullName = fullName
	newCustomer.Email = email
	newCustomer.Password.Set(password)

	v := newCustomer.Validate()
	if v != nil {
		return Token{}, v
	}

	// 2. check that email is not used
	admin, err := a.models.Customers.InsertCustomer(newCustomer)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrDuplicateCustomer):
			return Token{}, ErrEmailAlreadyTaken
		default:
			return Token{}, fmt.Errorf("admin signup: %w", err)
		}
	}

	return a.generateToken(admin.ID, models.ScopeAuth, 3*time.Hour)
}
func (a *CustomerService) Login(email, password string) (Token, error) {
	// 1. Fetch the provided email
	admin, err := a.models.Customers.GetCustomerByEmail(email)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrCustomerNotFound):
			return Token{}, ErrInvalidCredentials
		default:
			return Token{}, fmt.Errorf("admin login: %w", err)
		}
	}
	// 2. Check for password match
	match, err := admin.Password.Matches(password)
	if err != nil {
		return Token{}, fmt.Errorf("admin login: %w", err)
	}

	if !match {
		return Token{}, ErrInvalidCredentials
	}

	// 3. Return an Auth token
	return a.generateToken(admin.ID, models.ScopeAuth, 3*time.Hour)
}

func (a *CustomerService) Logout(token Token) error {
	admin, err := a.GetCustomerbyAuthToken(token.Plaintext)
	if err != nil {
		return fmt.Errorf("admin logout: %w", err)
	}

	err = a.models.Tokens.DeleteTokensByID(admin.ID)
	if err != nil {
		return err
	}

	return nil
}

func (a *CustomerService) generateToken(adminId int64, scope string, ttl time.Duration) (Token, error) {
	token, err := domain.Generate(adminId, scope, ttl)

	err = a.models.Tokens.InsertToken(token)
	if err != nil {
		return Token{}, fmt.Errorf("admin generate token: %w", err)
	}

	result := Token{
		Plaintext:  token.Plaintext,
		ExpiryTime: token.ExpiryTime,
	}
	return result, nil
}

func (a *CustomerService) GetCustomerbyAuthToken(tokenText string) (domain.Customer, error) {
	token, err := a.models.Tokens.GetToken(tokenText, models.ScopeAuth)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrTokenNotFound):
			return domain.Customer{}, ErrInvalidAuthToken
		default:
			return domain.Customer{}, fmt.Errorf("getCustomerByToken: %w", err)
		}
	}

	c, err := a.models.Customers.GetCustomerByID(token.UserID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrCustomerNotFound):
			return domain.Customer{}, ErrInvalidAuthToken
		default:
			return domain.Customer{}, fmt.Errorf("getCustomerByToken: %w", err)
		}
	}

	return c, nil
}
