package domain

import (
	"AhmadAbdelrazik/arbun/internal/validator"
)

type Admin struct {
	User
}

func (a Admin) Validate() *validator.Validator {
	v := validator.New()
	a.User.Validate(v)
	return v.Err()
}
