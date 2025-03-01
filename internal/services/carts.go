package services

import (
	"AhmadAbdelrazik/arbun/internal/domain"
	"AhmadAbdelrazik/arbun/internal/models"
	"AhmadAbdelrazik/arbun/internal/pkg/validator"
	"fmt"

	"github.com/Rhymond/go-money"
)

type CartService struct {
	models *models.Model
}

func newCartService(models *models.Model) *CartService {
	return &CartService{
		models: models,
	}
}

func (c *CartService) GetCart(customerID int64) (domain.Cart, error) {
	items, err := c.models.Carts.GetAll(customerID)
	if err != nil {
		return domain.Cart{}, fmt.Errorf("getCartItems: %w", err)
	}

	userCart := domain.Cart{
		Items: make([]domain.CartItem, 0, len(items)),
		Price: money.New(0, money.EGP),
	}

	for _, item := range items {
		product, err := c.models.Products.GetProductByID(item.ProductID)
		if err != nil {
			return domain.Cart{}, fmt.Errorf("getCartItems: %w", err)
		}

		cartItem := domain.CartItem{}
		cartItem.Populate(product, item.Amount)

		userCart.Price, err = userCart.Price.Add(cartItem.TotalPrice)
		if err != nil {
			return domain.Cart{}, fmt.Errorf("getCartItems: %w", err)
		}
		userCart.Items = append(userCart.Items, cartItem)
	}

	return userCart, nil
}

// UpdateItems - will add non existent items and set the items by productAmount
// if item exists, then their amount will be productAmount
func (c *CartService) UpdateItems(customerID int64, items []domain.CartItem) (domain.Cart, error) {
	// 1. check if items exists, and amounts are available
	err := c.checkItems(items)
	if err != nil {
		return domain.Cart{}, fmt.Errorf("AddItems: %w", err)
	}

	// 2. Insert Items
	for _, item := range items {
		err := c.models.Carts.SetItem(customerID, item)
		if err != nil {
			return domain.Cart{}, fmt.Errorf("AddItems: %w", err)
		}
	}

	return c.GetCart(customerID)
}

func (c *CartService) DeleteItem(customerID, productID int64) (domain.Cart, error) {
	err := c.models.Carts.DeleteItem(customerID, productID)
	if err != nil {
		return domain.Cart{}, fmt.Errorf("DeleteItem: %w", err)
	}

	return c.GetCart(customerID)
}

func (c *CartService) checkItems(items []domain.CartItem) error {
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
