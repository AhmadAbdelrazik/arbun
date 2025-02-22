package models

import "AhmadAbdelrazik/arbun/internal/validator"

type User interface {
	Validate() *validator.Validator
}
