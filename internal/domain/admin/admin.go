package admin

import (
	"AhmadAbdelrazik/arbun/internal/domain/user"
	"AhmadAbdelrazik/arbun/internal/validator"
)

type Admin struct {
	ID       int64
	Email    string
	Password user.Password
	FullName string
	Version  int
}

func (a Admin) Validate() *validator.Validator {
	v := validator.New()

	v.Check(a.FullName != "", "full_name", "must not be empty")
	v.Check(len(a.FullName) <= 40, "full_name", "must not be more than 40")

	v.Check(a.Email != "", "email", "must not be empty")
	v.Check(v.Matches(a.Email, *validator.EmailRX), "email", "must be a valid email address")

	a.Password.Validate(v)

	return v.Err()
}
