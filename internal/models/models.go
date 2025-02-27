package models

type Model struct {
	Products *ProductModel
	Tokens   *TokenModel
	Carts    *CartModel
	Orders   *OrderModel
	Users    *UserModel
}

func NewModel() *Model {
	product := newProductModel()
	return &Model{
		Products: product,
		Tokens:   newTokenModel(),
		Carts:    newCartModel(),
		Orders:   newOrderModel(product),
		Users:    newUserModel(),
	}
}
