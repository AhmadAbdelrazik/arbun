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
	Stripe   *stripe.StripeService
}

func New() *Services {
	model := models.NewModel()
	cartService := newCartService(model)
	// TODO: Provide method of passing secret key
	// TODO: Provide the stripe api success and canceled URLs
	stripeService := stripe.New("", "", "", "", model)
	return &Services{
		Products: newProductService(model),
		Users:    newUserService(model),
		Carts:    cartService,
		Orders:   newOrderService(model, cartService, stripeService),
		Stripe:   stripeService,
	}
}
