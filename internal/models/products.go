package models

import (
	"AhmadAbdelrazik/arbun/internal/domain/product"
	"errors"
	"slices"
)

var (
	ErrDuplicateProduct = errors.New("duplicate product")
	ErrProductNotFound  = errors.New("product not found")
	ErrEditConflict     = errors.New("edit conflict")
)

type ProductModel struct {
	products  []product.Product
	idCounter int64
}

func NewProductModel() *ProductModel {
	return &ProductModel{
		products:  make([]product.Product, 0),
		idCounter: 1,
	}
}

func (m *ProductModel) InsertProduct(p product.Product) (product.Product, error) {
	for _, pp := range m.products {
		if pp.Name == p.Name && pp.Vendor == p.Vendor {
			return product.Product{}, ErrDuplicateProduct
		}
	}

	p.ID = m.idCounter
	p.Version = 1
	m.products = append(m.products, p)

	m.idCounter++

	return p, nil
}

func (m *ProductModel) GetProductByID(id int64) (product.Product, error) {
	for _, pp := range m.products {
		if pp.ID == id {
			return pp, nil
		}
	}

	return product.Product{}, ErrProductNotFound
}

func (m *ProductModel) GetAllProducts() ([]product.Product, error) {
	return m.products, nil
}

func (m *ProductModel) UpdateProduct(p product.Product) (product.Product, error) {
	var result product.Product

	for i, pp := range m.products {
		if pp.ID == p.ID {
			if pp.Version != p.Version {
				return product.Product{}, ErrEditConflict
			}
			m.products[i] = p
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
			m.products = slices.Delete(m.products, i, i+1)
			return nil
		}
	}

	return ErrProductNotFound
}
