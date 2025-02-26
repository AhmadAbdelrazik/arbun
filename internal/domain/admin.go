package domain

import (
	"AhmadAbdelrazik/arbun/internal/pkg/validator"
)

type Admin struct {
	User
}

func (a Admin) Validate() *validator.Validator {
	v := validator.New()
	a.User.Validate()
	return v.Err()
}
