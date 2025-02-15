package services

type Services struct {
	Products *ProductService
	Admins   *AdminService
	Users    *UserService
}

func New() *Services {
	return &Services{
		Products: newProductService(),
		Admins:   newAdminService(),
		Users:    newUserService(),
	}
}
