package services

import (
	"AhmadAbdelrazik/arbun/internal/domain"
	"AhmadAbdelrazik/arbun/internal/models"
	"errors"
	"fmt"
)

var (
	ErrDuplicateProduct = errors.New("product already exists")
	ErrProductNotFound  = errors.New("product not found")
	ErrEditConflict     = errors.New("edit conflict")
)

type ProductService struct {
	models *models.Model
}

func newProductService(model *models.Model) *ProductService {
	return &ProductService{
		models: model,
	}
}

func (s *ProductService) InsertProduct(p domain.Product) (domain.Product, error) {
	v := p.Validate()
	if v != nil {
		return domain.Product{}, v
	}

	newProduct, err := s.models.Products.InsertProduct(p)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrDuplicateProduct):
			return domain.Product{}, ErrDuplicateProduct
		default:
			return domain.Product{}, fmt.Errorf("insert product: %w", err)
		}
	}

	return newProduct, nil
}

func (p *ProductService) GetProductByID(id int64) (domain.Product, error) {
	produc, err := p.models.Products.GetProductByID(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrProductNotFound):
			return domain.Product{}, ErrProductNotFound
		default:
			return domain.Product{}, fmt.Errorf("get product by id: %w", err)
		}
	}
	return produc, nil
}

func (p *ProductService) GetAllProducts() ([]domain.Product, error) {
	products, err := p.models.Products.GetAllProducts()
	if err != nil {
		return []domain.Product{}, fmt.Errorf("get all products: %w", err)
	}
	return products, nil
}

func updateProductFields(newProduct, oldProduct domain.Product) domain.Product {
	result := oldProduct

	if newProduct.Name != "" {
		result.Name = newProduct.Name
	}
	if newProduct.Description != "" {
		result.Description = newProduct.Description
	}
	if newProduct.Vendor != "" {
		result.Vendor = newProduct.Vendor
	}
	if newProduct.Properties != nil {
		result.Properties = newProduct.Properties
	}
	if newProduct.Price != nil {
		result.Price = newProduct.Price
	}
	if newProduct.AvailableAmount != 0 {
		result.AvailableAmount = newProduct.AvailableAmount
	}

	return result
}

func (p *ProductService) UpdateProduct(update domain.Product) (domain.Product, error) {
	oldProduct, err := p.GetProductByID(update.ID)
	if err != nil {
		return domain.Product{}, fmt.Errorf("update product: %w", err)
	}

	updatedProduct := updateProductFields(update, oldProduct)

	result, err := p.models.Products.UpdateProduct(updatedProduct)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrEditConflict):
			return domain.Product{}, ErrEditConflict
		default:
			return domain.Product{}, fmt.Errorf("update product: %w", err)
		}
	}

	return result, nil
}

func (p *ProductService) DeleteProduct(id int64) error {
	err := p.models.Products.DeleteProduct(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrProductNotFound):
			return ErrProductNotFound
		default:
			return fmt.Errorf("get product by id: %w", err)
		}
	}

	return nil
}
