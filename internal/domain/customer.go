package domain

import (
	"AhmadAbdelrazik/arbun/internal/pkg/validator"
)

type Address struct {
	Governorate    string `json:"governorate"`
	City           string `json:"city"`
	Street         string `json:"street"`
	AdditionalInfo string `json:"additional_info"`
}

type Customer struct {
	User
	Address     Address `json:"address"`
	MobilePhone string  `json:"mobile_phone"`
}

func (c Customer) Validate() *validator.Validator {
	v := validator.New()
	c.User.Validate()
	return v.Err()
}
