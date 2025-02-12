package services

import (
	"AhmadAbdelrazik/arbun/internal/repository"
	"fmt"
)

type ProductService struct {
	model repository.ProductModel
}

func (p ProductService) InsertProduct(name string, description string, properties map[string]string) (int64, error) {
	product := repository.Product{
		Name:        name,
		Description: description,
		Properties:  properties,
	}

	id, err := p.model.InsertProduct(product)
	if err != nil {
		return 0, fmt.Errorf("repository layer error: %w", err)
	}

	return id, nil
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

func (p ProductService) UpdateProduct(id int64, name string, description string, properties map[string]string) error {
	product := repository.Product{
		Name:        name,
		Description: description,
		Properties:  properties,
	}

	err := p.model.UpdateProduct(product)
	if err != nil {
		return fmt.Errorf("repository layer error: %w", err)
	}

	return nil
}

func (p ProductService) DeleteProduct(id int64) error {
	err := p.DeleteProduct(id)
	if err != nil {
		return fmt.Errorf("repository layer error: %w", err)
	}

	return nil
}
