package models

type Model struct {
	Products IProductModel
	Tokens   ITokenModel
	Carts    ICartModel
	Orders   IOrderModel
	Users    IUserModel
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
