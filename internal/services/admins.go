package services

import (
	"AhmadAbdelrazik/arbun/internal/repository"
	"errors"
	"fmt"
	"time"
)

type AdminService struct {
	model  *repository.AdminModel
	tokens *repository.TokenModel
}

var (
	ErrEmailAlreadyTaken  = errors.New("email already taken")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAuthToken   = errors.New("invalid auth token")
)

func newAdminService() *AdminService {
	return &AdminService{
		model:  repository.NewAdminModel(),
		tokens: repository.NewTokenModel(),
	}
}

type Token struct {
	Plaintext  string
	ExpiryTime time.Time
}

func (a *AdminService) GenerateToken(adminId int64, scope string, ttl time.Duration) (Token, error) {
	token, err := repository.GenerateToken(adminId, scope, ttl)

	err = a.tokens.InsertToken(token)
	if err != nil {
		return Token{}, fmt.Errorf("admin generate token: %w", err)
	}

	result := Token{
		Plaintext:  token.Plaintext,
		ExpiryTime: token.ExpiryTime,
	}
	return result, nil
}

func (a *AdminService) Signup(fullName, email, password string) (Token, error) {
	// 1. user provide credentials
	// TODO: Implement Regex Validation
	newAdmin := repository.Admin{
		FullName: fullName,
		Email:    email,
	}
	newAdmin.Password.Set(password)

	// 2. check that email is not used
	admin, err := a.model.InsertAdmin(newAdmin)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrDuplicateAdmin):
			return Token{}, ErrEmailAlreadyTaken
		default:
			return Token{}, fmt.Errorf("admin signup: %w", err)
		}
	}

	return a.GenerateToken(admin.ID, repository.ScopeAuth, 3*time.Hour)
}
func (a *AdminService) Login(email, password string) (Token, error) {
	// 1. Fetch the provided email
	admin, err := a.model.GetAdminByEmail(email)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrAdminNotFound):
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
	return a.GenerateToken(admin.ID, repository.ScopeAuth, 3*time.Hour)
}

func (a *AdminService) Logout(token Token) error {
	admin, err := a.getAdminByToken(token.Plaintext, repository.ScopeAuth)
	if err != nil {
		return fmt.Errorf("admin logout: %w", err)
	}

	err = a.tokens.DeleteTokensByID(admin.ID)
	if err != nil {
		return err
	}

	return nil
}

func (a *AdminService) getAdminByToken(tokenText, scope string) (repository.Admin, error) {
	token, err := a.tokens.GetToken(tokenText, scope)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrTokenNotFound):
			return repository.Admin{}, ErrInvalidAuthToken
		default:
			return repository.Admin{}, fmt.Errorf("getAdminByToken: %w", err)
		}
	}

	admin, err := a.model.GetAdminByID(token.AdminID)
	if err != nil {
		return repository.Admin{}, fmt.Errorf("getAdminByToken: %w", err)
	}

	return admin, nil
}
