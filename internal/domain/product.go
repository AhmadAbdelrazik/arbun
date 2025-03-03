package domain

import (
	"AhmadAbdelrazik/arbun/internal/pkg/validator"
	"encoding/json"

	"github.com/Rhymond/go-money"
)

type Product struct {
	ID              int64             `json:"id,omitempty"`
	Name            string            `json:"name,omitempty"`
	Description     string            `json:"description,omitempty"`
	Vendor          string            `json:"vendor,omitempty"`
	Properties      map[string]string `json:"properties,omitempty"`
	Price           *money.Money      `json:"price,omitempty"`
	AvailableAmount int               `json:"available_amount,omitempty"`
	Images          []string          `json:"images,omitempty"`
	Version         int               `json:"-"`
}

func (p Product) Validate() *validator.Validator {
	v := validator.New()

	v.Check(p.Name != "", "name", "can't be empty")
	v.Check(len(p.Name) < 100, "name", "must not be more than 100 bytes")

	v.Check(p.Description != "", "description", "can't be empty")
	v.Check(len(p.Description) < 2000, "description", "must not be more than 2000 bytes")

	v.Check(p.Price.IsPositive(), "price", "must be more than 0")

	v.Check(p.AvailableAmount > 0, "amount", "must be more than 0")

	return v.Err()
}

func (p Product) String() string {
	js, _ := json.MarshalIndent(p, "", "\t")
	return string(js)
}
