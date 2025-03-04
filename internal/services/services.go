package services

import (
	"AhmadAbdelrazik/arbun/internal/models"
	"AhmadAbdelrazik/arbun/internal/stripe"
)

type Services struct {
	Products *ProductService
	Users    *UserService
	Carts    *CartService
	Orders   *OrderService
}

func New() *Services {
	models := models.NewModel()
	cartService := newCartService(models)
	// TODO: Provide method of passing secret key
	// TODO: Provide the stripe api success and canceled URLs
	stripeService := stripe.New("", "", "")
	return &Services{
		Products: newProductService(models),
		Users:    newUserService(models),
		Carts:    cartService,
		Orders:   newOrderService(models, cartService, stripeService),
	}
}
