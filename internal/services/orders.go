package services

import (
	"AhmadAbdelrazik/arbun/internal/domain"
	"AhmadAbdelrazik/arbun/internal/models"
	"AhmadAbdelrazik/arbun/internal/stripe"
	"errors"
	"fmt"
	"time"
)

var (
	ErrOrderNotFound = errors.New("order not found")
)

type OrderService struct {
	models *models.Model
	carts  *CartService
	stripe *stripe.StripeService
}

func newOrderService(models *models.Model, cartService *CartService, stripeService *stripe.StripeService) *OrderService {
	return &OrderService{
		models: models,
		carts:  cartService,
		stripe: stripeService,
	}
}

func (s *OrderService) CreateCashOrder(customer domain.Customer) (domain.Order, error) {
	cart, err := s.carts.GetCart(customer.ID)
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

	return order, nil
}

func (s *OrderService) CreateCardOrder(customer domain.Customer) (string, error) {
	cart, err := s.carts.GetCart(customer.ID)
	if err != nil {
		return "", fmt.Errorf("CreateOrder :%w", err)
	}

	o := domain.Order{
		CustomerID:  customer.ID,
		CreatedAt:   time.Now(),
		Cart:        cart,
		PaymentType: domain.PaymentCard,
		Address:     customer.Address,
		MobilePhone: customer.MobilePhone,
		Status:      domain.StatusDispatched,
	}

	_, err = s.models.Orders.Create(o)
	if err != nil {
		return "", fmt.Errorf("CreateOrder: %w", err)
	}

	url, err := s.stripe.Checkout(o, customer)
	if err != nil {
		return "", nil
	}

	return url, nil
}

func (s *OrderService) returnItems(cart domain.Cart) {
	for _, item := range cart.Items {
		s.models.Products.ChangeProductAmountBy(item.ProductID, item.Amount)
	}
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

func (s *OrderService) ChangeOrderStatus(orderID int64, status domain.OrderStatus) error {
	order, err := s.models.Orders.Get(orderID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrOrderNotFound):
			return ErrOrderNotFound
		default:
			return fmt.Errorf("GetOrder: %w", err)
		}
	}

	if status == domain.StatusCanceled {
		s.returnItems(order.Cart)
	}

	err = s.models.Orders.Update(orderID, status)
	if err != nil {
		return fmt.Errorf("CancelOrder: %w", err)
	}

	return nil
}
