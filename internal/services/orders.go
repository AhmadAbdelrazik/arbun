package services

import (
	"AhmadAbdelrazik/arbun/internal/domain"
	"AhmadAbdelrazik/arbun/internal/models"
	"errors"
	"fmt"
	"time"
)

var (
	ErrOrderNotFound = errors.New("order not found")
)

type OrderService struct {
	models *models.Model
}

func newOrderService(models *models.Model) *OrderService {
	return &OrderService{
		models: models,
	}
}

func (s *OrderService) CreateOrder(customer domain.Customer, cartService *CartService) (domain.Order, error) {
	cart, err := cartService.GetCart(customer.ID)
	if err != nil {
		return domain.Order{}, fmt.Errorf("CreateOrder :%w", err)
	}

	o := domain.Order{
		CustomerID:  customer.ID,
		CreatedAt:   time.Now(),
		Cart:        cart,
		PaymentType: domain.PaymentCash,
		Address:     customer.Address,
		MobilePhone: customer.MobilePhone,
		Status:      domain.StatusDispatched,
	}

	order, err := s.models.Orders.Create(o)
	if err != nil {
		return domain.Order{}, fmt.Errorf("CreateOrder: %w", err)
	}

	//TODO: Implement Mailer Here

	return order, nil
}

func (s *OrderService) GetOrder(customer domain.Customer, orderID int64) (domain.Order, error) {
	order, err := s.models.Orders.Get(orderID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrOrderNotFound):
			return domain.Order{}, ErrOrderNotFound
		default:
			return domain.Order{}, fmt.Errorf("GetOrder: %w", err)
		}
	}

	if order.CustomerID != customer.ID {
		return domain.Order{}, ErrOrderNotFound
	}

	return order, nil
}

func (s *OrderService) GetAllUserOrders(customer domain.Customer) ([]domain.Order, error) {
	orders, err := s.models.Orders.GetAll(customer.ID)
	if err != nil {
		return nil, fmt.Errorf("GetOrder: %w", err)
	}

	return orders, nil
}
