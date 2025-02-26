package models

type Model struct {
	Products *ProductModel
	Tokens   *TokenModel
	Carts    *CartModel
	Orders   *OrderModel
	Users    *UserModel
}

func NewModel() *Model {
	return &Model{
		Products: newProductModel(),
		Tokens:   newTokenModel(),
		Carts:    newCartModel(),
		Orders:   newOrderModel(),
		Users:    newUserModel(),
	}
}
