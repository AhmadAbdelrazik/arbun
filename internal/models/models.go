package models

type Model struct {
	Products  *ProductModel
	Admins    *AdminModel
	Customers *CustomerModel
	Tokens    *TokenModel
	Carts     *CartModel
	Orders    *OrderModel
	Users     *UserModel
}

func NewModel() *Model {
	return &Model{
		Products:  newProductModel(),
		Admins:    newAdminModel(),
		Customers: newCustomerModel(),
		Tokens:    newTokenModel(),
		Carts:     newCartModel(),
		Orders:    newOrderModel(),
		Users:     newUserModel(),
	}
}
