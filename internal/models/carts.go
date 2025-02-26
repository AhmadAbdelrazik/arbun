package models

import (
	"AhmadAbdelrazik/arbun/internal/domain"
	"slices"
)

type CartModel struct {
	carts     map[int64][]domain.CartItem
	idCounter int64
}

func newCartModel() *CartModel {
	return &CartModel{
		carts:     make(map[int64][]domain.CartItem),
		idCounter: 1,
	}
}

func (m *CartModel) GetAll(customerID int64) ([]domain.CartItem, error) {
	if _, ok := m.carts[customerID]; !ok {
		m.carts[customerID] = make([]domain.CartItem, 0)
	}

	return m.carts[customerID], nil
}

func (m *CartModel) SetItem(customerID int64, item domain.CartItem) error {

	for i, cartItem := range m.carts[customerID] {
		if cartItem.ProductID == item.ProductID {
			m.carts[customerID][i].Amount = item.Amount
			return nil
		}
	}

	m.carts[customerID] = append(m.carts[customerID], item)

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
