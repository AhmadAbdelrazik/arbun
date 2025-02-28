package models

import (
	"AhmadAbdelrazik/arbun/internal/domain"
	"errors"
	"fmt"
)

var (
	ErrOrderNotFound = errors.New("order not found")
)

type OrderModel struct {
	products  *ProductModel
	orders    []domain.Order
	idCounter int64
}

func newOrderModel(products *ProductModel) *OrderModel {
	return &OrderModel{
		orders:    make([]domain.Order, 0),
		idCounter: 1,
		products:  products,
	}
}

func (m *OrderModel) Create(order domain.Order) (domain.Order, error) {
	committedChanges := make(map[int64]int)

	for _, item := range order.Cart.Items {
		err := m.products.ChangeProductAmountBy(item.ProductID, -item.Amount)
		if err != nil {
			for productID, amount := range committedChanges {
				err := m.products.ChangeProductAmountBy(productID, amount)
				if err != nil {
					panic(err)
				}
			}

			return domain.Order{}, fmt.Errorf(
				"product %s with id %d: %w",
				item.Name,
				item.ProductID,
				ErrInsufficientProductAmount,
			)
		}
		committedChanges[item.ProductID] = item.Amount
	}

	order.ID = m.idCounter
	m.idCounter++
	m.orders = append(m.orders, order)

	return order, nil
}

func (m *OrderModel) Update(orderID int64, status domain.OrderStatus) error {
	for i, order := range m.orders {
		if order.ID == orderID {
			m.orders[i].Status = status
			return nil
		}
	}

	return ErrOrderNotFound
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
	orders := make([]domain.Order, 0, 10)
	for _, o := range orders {
		if o.CustomerID == customerID {
			orders = append(orders, o)
		}
	}

	return orders, nil
}
