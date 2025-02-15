package services

type Services struct {
	Products *ProductService
	Admins   *AdminService
	Users    *UserService
}

func New() *Services {
	adminService := newAdminService()
	return &Services{
		Products: newProductService(),
		Admins:   adminService,
		Users:    newUserService(adminService),
	}
}
