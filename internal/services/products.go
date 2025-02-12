package services

import (
	"AhmadAbdelrazik/arbun/internal/repository"
	"fmt"
)

type ProductService struct {
	model *repository.ProductModel
}

func newProductService() *ProductService {
	return &ProductService{
		model: repository.NewProductModel(),
	}
}

type InsertProductParam struct {
	Name            string
	Description     string
	Vendor          string
	Properties      map[string]string
	AvailableAmount int
}

func (p ProductService) InsertProduct(param InsertProductParam) (repository.Product, error) {
	product := repository.Product{
		Name:            param.Name,
		Description:     param.Description,
		Properties:      param.Properties,
		Vendor:          param.Vendor,
		AvailableAmount: param.AvailableAmount,
	}

	newProduct, err := p.model.InsertProduct(product)
	if err != nil {
		return repository.Product{}, fmt.Errorf("repository layer error: %w", err)
	}

	return newProduct, nil
}

func (p ProductService) GetProductByID(id int64) (repository.Product, error) {
	product, err := p.model.GetProductByID(id)
	if err != nil {
		return repository.Product{}, fmt.Errorf("repository layer error: %w", err)
	}
	return product, nil
}

func (p ProductService) GetAllProducts() ([]repository.Product, error) {
	products, err := p.model.GetAllProducts()
	if err != nil {
		return []repository.Product{}, fmt.Errorf("repository layer error: %w", err)
	}
	return products, nil
}

type UpdateProductParam struct {
	ID              int64
	Name            *string
	Description     *string
	Vendor          *string
	Properties      map[string]string
	AvailableAmount *int
}

func (u *UpdateProductParam) UpdateProduct(product repository.Product) repository.Product {
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
	if u.Properties != nil {
		result.Properties = u.Properties
	}

	return result
}

func (p ProductService) UpdateProduct(param UpdateProductParam) (repository.Product, error) {
	fetchedProduct, err := p.GetProductByID(param.ID)
	if err != nil {
		return repository.Product{}, fmt.Errorf("repository layer error: %w", err)
	}

	product := param.UpdateProduct(fetchedProduct)

	updatedProduct, err := p.model.UpdateProduct(product)
	if err != nil {
		return repository.Product{}, fmt.Errorf("repository layer error: %w", err)
	}

	return updatedProduct, nil
}

func (p ProductService) UpdateProductAmount(id int64, newAmount int) (repository.Product, error) {
	product, err := p.GetProductByID(id)
	if err != nil {
		return repository.Product{}, fmt.Errorf("repository layer error: %w", err)
	}

	product.AvailableAmount = newAmount

	updatedProduct, err := p.model.UpdateProduct(product)
	if err != nil {
		return repository.Product{}, fmt.Errorf("repository layer error: %w", err)
	}

	return updatedProduct, nil
}

func (p ProductService) DeleteProduct(id int64) error {
	err := p.DeleteProduct(id)
	if err != nil {
		return fmt.Errorf("repository layer error: %w", err)
	}

	return nil
}
