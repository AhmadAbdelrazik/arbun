package repository

import (
	"encoding/json"
	"errors"
)

var (
	ErrDuplicateProduct = errors.New("duplicate product")
	ErrProductNotFound  = errors.New("product not found")
	ErrEditConflict     = errors.New("edit conflict")
)

type Product struct {
	ID              int64
	Name            string
	Description     string
	Vendor          string
	Properties      map[string]string
	AvailableAmount int
	Version         int
}

func (p Product) String() string {
	js, _ := json.MarshalIndent(p, "", "\t")
	return string(js)
}

type ProductModel struct {
	products  []Product
	idCounter int64
}

func NewProductModel() *ProductModel {
	return &ProductModel{
		products:  make([]Product, 0),
		idCounter: 1,
	}
}

func (m *ProductModel) InsertProduct(product Product) (Product, error) {
	for _, p := range m.products {
		if p.Name == product.Name && p.Vendor == product.Vendor {
			return Product{}, ErrDuplicateProduct
		}
	}

	product.ID = m.idCounter
	product.Version = 1
	m.products = append(m.products, product)

	m.idCounter++

	return product, nil
}

func (m *ProductModel) GetProductByID(id int64) (Product, error) {
	for _, p := range m.products {
		if p.ID == id {
			return p, nil
		}
	}

	return Product{}, ErrProductNotFound
}

func (m *ProductModel) GetAllProducts() ([]Product, error) {
	return m.products, nil
}

func (m *ProductModel) UpdateProduct(product Product) (Product, error) {
	var result Product

	for i, p := range m.products {
		if p.ID == product.ID {
			if p.Version != product.Version {
				return Product{}, ErrEditConflict
			}
			m.products[i] = product
			m.products[i].Version++

			result = m.products[i]
			break
		}
	}

	return result, nil

}

func (m *ProductModel) DeleteProduct(id int64) error {
	for i, p := range m.products {
		if p.ID == id {
			products := m.products[:i]
			products = append(products, m.products[i+1:]...)
			m.products = products
			return nil
		}
	}
	return ErrProductNotFound
}
