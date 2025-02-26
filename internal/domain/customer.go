package domain

import (
	"AhmadAbdelrazik/arbun/internal/pkg/validator"
)

type Customer struct {
	User
}

func (c Customer) Validate() *validator.Validator {
	v := validator.New()
	c.User.Validate(v)
	return v.Err()
}
