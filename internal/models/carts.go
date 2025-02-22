package models

import (
	"AhmadAbdelrazik/arbun/internal/validator"
	"fmt"
	"slices"
)

type CartItem struct {
	ProductID int64 `json:"product_id"`
	Amount    int   `json:"amount"`
}

func (c CartItem) Validate() *validator.Validator {
	v := validator.New()

	v.Check(c.Amount != 0, "quantity", "must not be 0")

	return v.Err()
}

type CartModel struct {
	carts     map[int64][]CartItem
	idCounter int64
}

func (m *CartModel) GetAll(customerID int64) ([]CartItem, error) {
	return m.carts[customerID], nil
}

func (m *CartModel) SetItem(customerID int64, item CartItem) error {
	cartItems, err := m.GetAll(customerID)
	if err != nil {
		return fmt.Errorf("InsertItem: %w", err)
	}

	for i, cartItem := range cartItems {
		if cartItem.ProductID == item.ProductID {
			cartItems[i].Amount = item.Amount
			return nil
		}
	}

	cartItems = append(cartItems, item)

	return nil
}

func (m *CartModel) DeleteItem(customerID, productID int64) error {
	for i, item := range m.carts[customerID] {
		if item.ProductID == productID {
			m.carts[customerID] = slices.Delete(m.carts[customerID], i, i+1)
			return nil
		}
	}
	return nil
}
