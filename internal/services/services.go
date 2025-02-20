package services

import "AhmadAbdelrazik/arbun/internal/repository"

type Services struct {
	Products *ProductService
	Users    *UserService
	Carts    *CartService
}

func New() *Services {
	models := repository.NewModel()
	return &Services{
		Products: newProductService(models),
		Users:    newUserService(models),
		Carts:    newCartService(models),
	}
}
