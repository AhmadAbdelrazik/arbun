package admin

import (
	"AhmadAbdelrazik/arbun/internal/validator"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Admin struct {
	ID       int64
	Email    string
	Password Password
	FullName string
	Version  int
}

func (a Admin) Validate() *validator.Validator {
	v := validator.New()

	v.Check(a.FullName != "", "full_name", "must not be empty")
	v.Check(len(a.FullName) <= 40, "full_name", "must not be more than 40")

	v.Check(a.Email != "", "email", "must not be empty")
	v.Check(v.Matches(a.Email, *validator.EmailRX), "email", "must be a valid email address")

	if a.Password.plaintext != nil {
		password := *a.Password.plaintext
		v.Check(password != "", "password", "must not be empty")
		v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
		v.Check(len(password) <= 72, "password", "must be more than 72 bytes long")
	}

	if a.Password.hash == nil {
		panic("missing password hash for user")
	}

	return v.Err()
}

type Password struct {
	plaintext *string
	hash      []byte
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
