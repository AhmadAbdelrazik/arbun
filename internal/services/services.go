package services

import "AhmadAbdelrazik/arbun/internal/models"

type Services struct {
	Products *ProductService
	Users    *UserService
	Carts    *CartService
}

func New() *Services {
	models := models.NewModel()
	return &Services{
		Products: newProductService(models),
		Users:    newUserService(models),
		Carts:    newCartService(models),
	}
}
