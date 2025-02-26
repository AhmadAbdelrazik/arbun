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

func (s *OrderService) CreateOrder(customer int64) {
	// NOTE: Figure out how services communications should be
}
