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

func (a Address) Validate() *validator.Validator {
	v := validator.New()

	v.Check(a.Governorate != "", "governorate", "must not be empty")
	v.Check(a.City != "", "city", "must not be empty")
	v.Check(a.Street != "", "street", "must not be empty")
	v.Check(len(a.AdditionalInfo) < 1000, "additional info", "must be less than 1000 bytes")

	return v.Err()
}

type MobilePhone string

func (p MobilePhone) Validate() *validator.Validator {
	v := validator.New()
	v.Check(validator.Matches(string(p), *validator.EgyPhoneNumbersRX), "mobile_phone", "invalid mobile phone")
	return v.Err()
}

type Customer struct {
	User
	Address     Address     `json:"address"`
	MobilePhone MobilePhone `json:"mobile_phone"`
	StripeID    string      `json:"-"`
}

func (c Customer) Validate() *validator.Validator {
	v := validator.New()

	v.Add(c.User.Validate())
	v.Add(c.Address.Validate())
	v.Add(c.MobilePhone.Validate())

	return v.Err()
}
