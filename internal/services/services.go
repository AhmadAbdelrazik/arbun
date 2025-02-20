package services

type Services struct {
	Products *ProductService
	Users    *UserService
}

func New() *Services {
	return &Services{
		Products: newProductService(),
		Users:    newUserService(),
	}
}
