package services

import "AhmadAbdelrazik/arbun/internal/models"

type OrderService struct {
	models *models.Model
}

func newOrderService(models *models.Model) *OrderService {
	return &OrderService{
		models: models,
	}
}

func (s *OrderService) CreateOrder(customerID int64, cart *CartService) {

}
