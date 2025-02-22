package product

import (
	"AhmadAbdelrazik/arbun/internal/validator"
	"encoding/json"
)

type Product struct {
	ID              int64
	Name            string
	Description     string
	Vendor          string
	Properties      map[string]string
	Price           float32
	AvailableAmount int
	Version         int
}

func (p Product) Validate() *validator.Validator {
	v := validator.New()

	v.Check(p.Name != "", "name", "can't be empty")
	v.Check(len(p.Name) < 100, "name", "must not be more than 100 bytes")

	v.Check(p.Description != "", "description", "can't be empty")
	v.Check(len(p.Description) < 2000, "description", "must not be more than 2000 bytes")

	v.Check(p.Price != 0, "price", "can't be zero")

	v.Check(p.AvailableAmount > 0, "amount", "must be more than 0")

	return v.Err()
}

func (p Product) String() string {
	js, _ := json.MarshalIndent(p, "", "\t")
	return string(js)
}
