package models

import (
	"AhmadAbdelrazik/arbun/internal/domain"
	"errors"
)

var (
	ErrOrderNotFound = errors.New("order not found")
)

type OrderModel struct {
	orders    []domain.Order
	idCounter int64
}

func newOrderModel() *OrderModel {
	return &OrderModel{
		orders:    make([]domain.Order, 0),
		idCounter: 1,
	}
}

func (m *OrderModel) Create(o domain.Order) error {
	o.ID = m.idCounter
	m.idCounter++
	m.orders = append(m.orders, o)
	return nil
}

func (m *OrderModel) Get(orderID int64) (domain.Order, error) {
	for _, o := range m.orders {
		if o.ID == orderID {
			return o, nil
		}
	}

	return domain.Order{}, ErrOrderNotFound
}

func (m *OrderModel) GetAll(customerID int64) ([]domain.Order, error) {
	// orders := make([]domain.Order, 0, 10)
	// for _, o := range orders {
	// 	if o.Customer.ID == customerID {
	// 		orders = append(orders, o)
	// 	}
	// }
	//
	// return orders, nil
	return nil, nil
}
