package user

import (
	"AhmadAbdelrazik/arbun/internal/validator"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type IUser interface {
}

type Password struct {
	plaintext *string
	hash      []byte
}

func (p *Password) Validate(v *validator.Validator) {
	if p.plaintext != nil {
		password := *p.plaintext
		v.Check(password != "", "password", "must not be empty")
		v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
		v.Check(len(password) <= 72, "password", "must be more than 72 bytes long")
	}

	if p.hash == nil {
		panic("missing password hash for user")
	}
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
