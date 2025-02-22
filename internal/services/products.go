package services

import (
	"AhmadAbdelrazik/arbun/internal/domain/product"
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

type InsertProductParam struct {
	Name            string
	Description     string
	Vendor          string
	Properties      map[string]string
	Price           float32
	AvailableAmount int
}

func (s *ProductService) InsertProduct(param InsertProductParam) (product.Product, error) {
	p := product.Product{
		Name:            param.Name,
		Description:     param.Description,
		Properties:      param.Properties,
		Vendor:          param.Vendor,
		Price:           param.Price,
		AvailableAmount: param.AvailableAmount,
	}

	v := p.Validate()
	if v != nil {
		return product.Product{}, v
	}

	newProduct, err := s.models.Products.InsertProduct(p)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrDuplicateProduct):
			return product.Product{}, ErrDuplicateProduct
		default:
			return product.Product{}, fmt.Errorf("insert product: %w", err)
		}
	}

	return newProduct, nil
}

func (p *ProductService) GetProductByID(id int64) (product.Product, error) {
	produc, err := p.models.Products.GetProductByID(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrProductNotFound):
			return product.Product{}, ErrProductNotFound
		default:
			return product.Product{}, fmt.Errorf("get product by id: %w", err)
		}
	}
	return produc, nil
}

func (p *ProductService) GetAllProducts() ([]product.Product, error) {
	products, err := p.models.Products.GetAllProducts()
	if err != nil {
		return []product.Product{}, fmt.Errorf("get all products: %w", err)
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

func (u *UpdateProductParam) updateProduct(product product.Product) product.Product {
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

func (p *ProductService) UpdateProduct(param UpdateProductParam) (product.Product, error) {
	fetchedProduct, err := p.GetProductByID(param.ID)
	if err != nil {
		return product.Product{}, fmt.Errorf("update product: %w", err)
	}

	prod := param.updateProduct(fetchedProduct)

	v := prod.Validate()
	if v != nil {
		return product.Product{}, v
	}

	updatedProduct, err := p.models.Products.UpdateProduct(prod)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrEditConflict):
			return product.Product{}, ErrEditConflict
		default:
			return product.Product{}, fmt.Errorf("update product: %w", err)
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
