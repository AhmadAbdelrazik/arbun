package admin

import (
	"AhmadAbdelrazik/arbun/internal/domain/user"
	"AhmadAbdelrazik/arbun/internal/validator"
)

type Admin struct {
	user.User
}

func (a Admin) Validate() *validator.Validator {
	v := validator.New()
	a.User.Validate(v)
	return v.Err()
}
