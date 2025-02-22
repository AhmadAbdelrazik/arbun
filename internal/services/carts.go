package services

import (
	"AhmadAbdelrazik/arbun/internal/domain/cart"
	"AhmadAbdelrazik/arbun/internal/models"
	"AhmadAbdelrazik/arbun/internal/validator"
	"fmt"
)

type CartService struct {
	models *models.Model
}

func newCartService(models *models.Model) *CartService {
	return &CartService{
		models: models,
	}
}

func (c *CartService) GetCart(customerID int64) (cart.Cart, error) {
	items, err := c.models.Carts.GetAll(customerID)
	if err != nil {
		return cart.Cart{}, fmt.Errorf("getCartItems: %w", err)
	}

	userCart := cart.Cart{
		Items: make([]cart.CartItem, 0, len(items)),
	}

	for _, item := range items {
		product, err := c.models.Products.GetProductByID(item.ProductID)
		if err != nil {
			return cart.Cart{}, fmt.Errorf("getCartItems: %w", err)
		}

		cartItem := cart.CartItem{}
		cartItem.Populate(product, item.Amount)

		userCart.Price += cartItem.TotalPrice
		userCart.Items = append(userCart.Items, cartItem)
	}

	return userCart, nil
}

type InputItem struct {
	ProductID int64 `json:"product_id"`
	Amount    int   `json:"amount"`
}

type AddItemsParam struct {
	CustomerID int64
	Items      []InputItem
}

// UpdateItems - will add non existent items and set the items by productAmount
// if item exists, then their amount will be productAmount
func (c *CartService) UpdateItems(input AddItemsParam) (cart.Cart, error) {
	// 1. check if items are legitimate
	items, err := c.checkItems(input.Items)
	if err != nil {
		return cart.Cart{}, fmt.Errorf("AddItems: %w", err)
	}

	// 2. Insert Items
	for _, item := range items {
		err := c.models.Carts.SetItem(input.CustomerID, item)
		if err != nil {
			return cart.Cart{}, fmt.Errorf("AddItems: %w", err)
		}
	}

	return c.GetCart(input.CustomerID)
}

func (c *CartService) DeleteItem(customerID, productID int64) (cart.Cart, error) {
	err := c.models.Carts.DeleteItem(customerID, productID)
	if err != nil {
		return cart.Cart{}, fmt.Errorf("DeleteItem: %w", err)
	}

	return c.GetCart(customerID)
}

func (c *CartService) checkItems(items []InputItem) ([]cart.CartItem, error) {
	if len(items) == 0 {
		return nil, fmt.Errorf("checkItems: items are empty")
	}

	result := make([]cart.CartItem, 0, len(items))

	for _, item := range items {
		product, err := c.models.Products.GetProductByID(item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("checkItems: %w", ErrProductNotFound)
		}

		cartItem := cart.CartItem{}

		if product.AvailableAmount < item.Amount {
			v := validator.New()
			v.AddError(
				fmt.Sprintf("product %v", item.ProductID),
				fmt.Sprintf("available %v only", product.AvailableAmount),
			)
			return nil, v.Err()
		}

		cartItem.Populate(product, item.Amount)

		result = append(result, cartItem)
	}

	return result, nil
}
