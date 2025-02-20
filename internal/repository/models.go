package repository

type Model struct {
	Products  *ProductModel
	Admins    *AdminModel
	Customers *CustomerModel
	Tokens    *TokenModel
	Carts     *CartModel
}

func NewModel() *Model {
	return &Model{
		Products:  NewProductModel(),
		Admins:    NewAdminModel(),
		Customers: NewCustomerModel(),
		Tokens:    NewTokenModel(),
	}
}
