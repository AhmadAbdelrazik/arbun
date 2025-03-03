package domain

import (
	"AhmadAbdelrazik/arbun/internal/pkg/validator"
	"maps"

	"github.com/Rhymond/go-money"
)

type Cart struct {
	Items []CartItem   `json:"items"`
	Price *money.Money `json:"price"`
}

type CartItem struct {
	ProductID   int64             `json:"product_id"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Amount      int               `json:"amount"`
	ItemPrice   *money.Money      `json:"item_price"`
	TotalPrice  *money.Money      `json:"total_price"`
	Properties  map[string]string `json:"properties,omitempty"`
	Images      []string          `json:"images,omitempty"`
}

func (c *CartItem) Populate(p Product, amount int) {
	c.ProductID = p.ID
	c.Name = p.Name
	c.Amount = amount
	c.ItemPrice = p.Price
	c.TotalPrice = p.Price.Multiply(int64(amount))
	c.Description = p.Description
	c.Properties = maps.Clone(p.Properties)
	c.Images = p.Images
}

func (c CartItem) Validate() *validator.Validator {
	v := validator.New()

	v.Check(c.Amount != 0, "quantity", "must not be 0")

	return v.Err()
}
