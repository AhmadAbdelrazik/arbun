package services

import (
	"AhmadAbdelrazik/arbun/internal/domain/product"
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

type CartItem struct {
	ProductID  int64
	Name       string
	Amount     int
	ItemPrice  float32
	TotalPrice float32
}

func (c *CartItem) Populate(p product.Product, amount int) {
	c.ProductID = p.ID
	c.Name = p.Name
	c.Amount = amount
	c.ItemPrice = p.Price
	c.TotalPrice = c.ItemPrice * float32(c.Amount)
}

type Cart struct {
	Items []CartItem
	Price float32
}

func (c *CartService) GetCart(customerID int64) (Cart, error) {
	items, err := c.models.Carts.GetAll(customerID)
	if err != nil {
		return Cart{}, fmt.Errorf("getCartItems: %w", err)
	}

	cartItems := make([]CartItem, 0, len(items))
	cart := Cart{}

	for _, item := range items {
		product, err := c.models.Products.GetProductByID(item.ProductID)
		if err != nil {
			return Cart{}, fmt.Errorf("getCartItems: %w", err)
		}

		cartItem := CartItem{}
		cartItem.Populate(product, item.Amount)

		cart.Price += cartItem.TotalPrice

		cartItems = append(cartItems, cartItem)
	}

	return cart, nil
}

type AddItemsParam struct {
	CustomerID int64
	Items      []models.CartItem
}

// UpdateItems - will add non existent items and set the items by productAmount
// if item exists, then their amount will be productAmount
func (c *CartService) UpdateItems(input AddItemsParam) (Cart, error) {
	// 1. check if items are legitimate
	err := c.checkItems(input.Items)
	if err != nil {
		return Cart{}, fmt.Errorf("AddItems: %w", err)
	}

	// 2. Insert Items
	for _, item := range input.Items {
		err := c.models.Carts.SetItem(input.CustomerID, item)
		if err != nil {
			return Cart{}, fmt.Errorf("AddItems: %w", err)
		}
	}

	return c.GetCart(input.CustomerID)
}

func (c *CartService) DeleteItem(customerID, productID int64) (Cart, error) {
	err := c.models.Carts.DeleteItem(customerID, productID)
	if err != nil {
		return Cart{}, fmt.Errorf("DeleteItem: %w", err)
	}

	return Cart{}, nil
}

func (c *CartService) checkItems(items []models.CartItem) error {
	if len(items) == 0 {
		return fmt.Errorf("checkItems: items are empty")
	}

	for _, item := range items {
		product, err := c.models.Products.GetProductByID(item.ProductID)
		if err != nil {
			return fmt.Errorf("checkItems: %w", ErrProductNotFound)
		}

		if product.AvailableAmount < item.Amount {
			v := validator.New()
			v.AddError(
				fmt.Sprintf("product %v", item.ProductID),
				fmt.Sprintf("available %v only", product.AvailableAmount),
			)
			return v.Err()
		}
	}

	return nil
}
