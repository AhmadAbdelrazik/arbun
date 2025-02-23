package customer

import (
	"AhmadAbdelrazik/arbun/internal/domain/user"
	"AhmadAbdelrazik/arbun/internal/validator"
)

type Customer struct {
	user.User
}

func (c Customer) Validate() *validator.Validator {
	v := validator.New()
	c.User.Validate(v)
	return v.Err()
}
