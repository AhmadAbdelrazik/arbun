package domain

import (
	"AhmadAbdelrazik/arbun/internal/pkg/validator"

	"github.com/Rhymond/go-money"
)

type Cart struct {
	Items []CartItem   `json:"items"`
	Price *money.Money `json:"price"`
}

type CartItem struct {
	ProductID  int64        `json:"product_id"`
	Name       string       `json:"name"`
	Amount     int          `json:"amount"`
	ItemPrice  *money.Money `json:"item_price"`
	TotalPrice *money.Money `json:"total_price"`
}

func (c *CartItem) Populate(p Product, amount int) {
	c.ProductID = p.ID
	c.Name = p.Name
	c.Amount = amount
	c.ItemPrice = p.Price
	c.TotalPrice = p.Price.Multiply(int64(amount))
}

func (c CartItem) Validate() *validator.Validator {
	v := validator.New()

	v.Check(c.Amount != 0, "quantity", "must not be 0")

	return v.Err()
}
