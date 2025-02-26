package services

import (
	"AhmadAbdelrazik/arbun/internal/domain"
	"AhmadAbdelrazik/arbun/internal/models"
	"errors"
	"fmt"
	"time"
)

type AdminService struct {
	models *models.Model
}

func newAdminService(models *models.Model) *AdminService {
	return &AdminService{
		models: models,
	}
}

func (a *AdminService) Signup(fullName, email, password string) (Token, error) {
	// 1. user provide credentials
	var newAdmin domain.Admin
	newAdmin.FullName = fullName
	newAdmin.Email = email
	newAdmin.Password.Set(password)

	v := newAdmin.Validate()
	if v != nil {
		return Token{}, v
	}

	// 2. check that email is not used
	admin, err := a.models.Admins.InsertAdmin(newAdmin)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrDuplicateAdmin):
			return Token{}, ErrEmailAlreadyTaken
		default:
			return Token{}, fmt.Errorf("admin signup: %w", err)
		}
	}

	return a.generateToken(admin.ID, models.ScopeAuth, 3*time.Hour)
}
func (a *AdminService) Login(email, password string) (Token, error) {
	// 1. Fetch the provided email
	admin, err := a.models.Admins.GetAdminByEmail(email)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrAdminNotFound):
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

func (a *AdminService) Logout(token Token) error {
	admin, err := a.GetAdminbyAuthToken(token.Plaintext)
	if err != nil {
		return fmt.Errorf("admin logout: %w", err)
	}

	err = a.models.Tokens.DeleteTokensByID(admin.ID)
	if err != nil {
		return err
	}

	return nil
}

func (a *AdminService) generateToken(adminId int64, scope string, ttl time.Duration) (Token, error) {
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

func (s *AdminService) GetAdminbyAuthToken(tokenText string) (domain.Admin, error) {
	token, err := s.models.Tokens.GetToken(tokenText, models.ScopeAuth)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrTokenNotFound):
			return domain.Admin{}, ErrInvalidAuthToken
		default:
			return domain.Admin{}, fmt.Errorf("getAdminByToken: %w", err)
		}
	}

	a, err := s.models.Admins.GetAdminByID(token.UserID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrAdminNotFound):
			return domain.Admin{}, ErrInvalidAuthToken
		default:
			return domain.Admin{}, fmt.Errorf("getAdminByToken: %w", err)
		}
	}

	return a, nil
}
