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

type UpdateProductParam struct {
	ID              int64
	Name            *string
	Description     *string
	Vendor          *string
	Price           *float32
	Properties      map[string]string
	AvailableAmount *int
}

func (u *UpdateProductParam) updateProduct(product domain.Product) domain.Product {
	result := product
	if u.Name != nil {
		result.Name = *u.Name
	}
	if u.Description != nil {
		result.Description = *u.Description
	}
	if u.Vendor != nil {
		result.Vendor = *u.Vendor
	}
	if u.AvailableAmount != nil {
		result.AvailableAmount = *u.AvailableAmount
	}
	if u.Price != nil {
		result.Price = *u.Price
	}
	if u.Properties != nil {
		result.Properties = u.Properties
	}

	return result
}

func (p *ProductService) UpdateProduct(param UpdateProductParam) (domain.Product, error) {
	fetchedProduct, err := p.GetProductByID(param.ID)
	if err != nil {
		return domain.Product{}, fmt.Errorf("update product: %w", err)
	}

	prod := param.updateProduct(fetchedProduct)

	v := prod.Validate()
	if v != nil {
		return domain.Product{}, v
	}

	updatedProduct, err := p.models.Products.UpdateProduct(prod)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrEditConflict):
			return domain.Product{}, ErrEditConflict
		default:
			return domain.Product{}, fmt.Errorf("update product: %w", err)
		}
	}

	return updatedProduct, nil
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
