package domain

import (
	"AhmadAbdelrazik/arbun/internal/pkg/validator"
	"encoding/json"
)

type Product struct {
	ID              int64             `json:"id"`
	Name            string            `json:"nme"`
	Description     string            `json:"description"`
	Vendor          string            `json:"vendor"`
	Properties      map[string]string `json:"properties"`
	Price           float32           `json:"price"`
	AvailableAmount int               `json:"available_amount"`
	Version         int               `json:"-"`
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
