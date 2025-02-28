package domain

import (
	"AhmadAbdelrazik/arbun/internal/pkg/validator"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const (
	TypeAdmin    = "admin"
	TypeCustomer = "customer"
)

type IUser interface {
}

type User struct {
	ID       int64    `json:"id"`
	Email    string   `json:"email"`
	Password Password `json:"-"`
	Name     string   `json:"name"`
	Type     string   `json:"-"`
	Version  int      `json:"version"`
}

func (u *User) Validate() *validator.Validator {
	v := validator.New()
	v.Check(u.Name != "", "full_name", "must not be empty")
	v.Check(len(u.Name) <= 40, "full_name", "must not be more than 40")

	v.Check(u.Email != "", "email", "must not be empty")
	v.Check(validator.Matches(u.Email, *validator.EmailRX), "email", "must be a valid email address")

	v.Add(u.Password.Validate())

	return v.Err()
}

type Password struct {
	plaintext *string
	hash      []byte
}

func (p *Password) Validate() *validator.Validator {
	v := validator.New()
	if p.plaintext != nil {
		password := *p.plaintext
		v.Check(password != "", "password", "must not be empty")
		v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
		v.Check(len(password) <= 72, "password", "must be more than 72 bytes long")
	}

	if p.hash == nil {
		panic("missing password hash for user")
	}

	return v.Err()
}

func (p *Password) Set(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.hash = hash
	p.plaintext = &password
	return nil
}

func (p *Password) Matches(plaintext string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintext))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
