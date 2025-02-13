package services

type Services struct {
	Products *ProductService
	Admins   *AdminService
}

func New() *Services {
	return &Services{
		Products: newProductService(),
		Admins:   newAdminService(),
	}
}
