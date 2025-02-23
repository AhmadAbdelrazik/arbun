package customer

import (
	"AhmadAbdelrazik/arbun/internal/domain/user"
	"AhmadAbdelrazik/arbun/internal/validator"
)

type Customer struct {
	ID       int64
	Email    string
	Password user.Password
	FullName string
	Version  int
}

func (c Customer) Validate() *validator.Validator {
	v := validator.New()
	v.Check(c.FullName != "", "full_name", "must not be empty")
	v.Check(len(c.FullName) <= 40, "full_name", "must not be more than 40")

	v.Check(c.Email != "", "email", "must not be empty")
	v.Check(v.Matches(c.Email, *validator.EmailRX), "email", "must be a valid email address")

	c.Password.Validate(v)

	return v.Err()
}
