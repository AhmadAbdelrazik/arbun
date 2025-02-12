package services

type Services struct {
	Products *ProductService
}

func New() *Services {
	return &Services{
		Products: newProductService(),
	}
}
