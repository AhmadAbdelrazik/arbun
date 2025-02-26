package services

import (
	"AhmadAbdelrazik/arbun/internal/domain"
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
	models    *models.Model
}

type Token struct {
	Plaintext  string
	ExpiryTime time.Time
}

func newUserService(models *models.Model) *UserService {
	return &UserService{
		admins:    newAdminService(models),
		customers: newCustomerService(models),
		models:    models,
	}
}

func (s *UserService) Signup(fullName, email, password, userType string) (Token, error) {
	// 1. user provide credentials
	user := domain.User{
		Name:  fullName,
		Email: email,
		Type:  userType,
	}
	user.Password.Set(password)

	v := user.Validate()
	if v != nil {
		return Token{}, v
	}

	// 2. check that email is not used
	newUser, err := s.models.Users.InsertUser(user)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrDuplicateUser):
			return Token{}, ErrEmailAlreadyTaken
		default:
			return Token{}, fmt.Errorf("user signup: %w", err)
		}
	}

	return s.generateToken(newUser.ID, newUser.Type, domain.ScopeAuth, 3*time.Hour)
}

func (s *UserService) Login(email, password string) (Token, error) {
	// 1. Fetch the provided email
	user, err := s.models.Users.GetUserByEmail(email)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrUserNotFound):
			return Token{}, ErrInvalidCredentials
		default:
			return Token{}, fmt.Errorf("user login: %w", err)
		}
	}
	// 2. Check for password match
	match, err := user.Password.Matches(password)
	if err != nil {
		return Token{}, fmt.Errorf("user login: %w", err)
	}

	if !match {
		return Token{}, ErrInvalidCredentials
	}

	// 3. Return an Auth token
	return s.generateToken(user.ID, user.Type, domain.ScopeAuth, 3*time.Hour)
}

func (s *UserService) Logout(token Token, userType string) error {
	switch userType {
	case domain.TypeAdmin:
		return s.admins.Logout(token)
	case domain.TypeCustomer:
		return s.customers.Logout(token)
	default:
		return fmt.Errorf("logout: %w", ErrInvalidUserType)
	}
}

func (s *UserService) GetUserByToken(tokenText string) (domain.User, error) {
	token, err := s.models.Tokens.GetToken(tokenText, domain.ScopeAuth)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrTokenNotFound):
			return domain.User{}, ErrInvalidAuthToken
		default:
			return domain.User{}, fmt.Errorf("getUserByToken: %w", err)
		}
	}

	a, err := s.models.Users.GetUserByID(token.UserID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrUserNotFound):
			return domain.User{}, ErrInvalidAuthToken
		default:
			return domain.User{}, fmt.Errorf("getUserByToken: %w", err)
		}
	}

	return a, nil
}

func (a *UserService) generateToken(userID int64, userType string, scope string, ttl time.Duration) (Token, error) {
	token, err := domain.NewToken(userID, userType, scope, ttl)

	err = a.models.Tokens.InsertToken(token)
	if err != nil {
		return Token{}, fmt.Errorf("userService generateToken: %w", err)
	}

	result := Token{
		Plaintext:  token.Plaintext,
		ExpiryTime: token.ExpiryTime,
	}
	return result, nil
}
