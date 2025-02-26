package domain

import (
	"AhmadAbdelrazik/arbun/internal/pkg/validator"
)

type Cart struct {
	Items []CartItem `json:"items"`
	Price float32    `json:"price"`
}

type CartItem struct {
	ProductID  int64   `json:"product_id"`
	Name       string  `json:"name"`
	Amount     int     `json:"amount"`
	ItemPrice  float32 `json:"item_price"`
	TotalPrice float32 `json:"total_price"`
}

func (c *CartItem) Populate(p Product, amount int) {
	c.ProductID = p.ID
	c.Name = p.Name
	c.Amount = amount
	c.ItemPrice = p.Price
	c.TotalPrice = c.ItemPrice * float32(c.Amount)
}

func (c CartItem) Validate() *validator.Validator {
	v := validator.New()

	v.Check(c.Amount != 0, "quantity", "must not be 0")

	return v.Err()
}
