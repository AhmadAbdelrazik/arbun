package repository

import "errors"

var (
	ErrDuplicateProduct = errors.New("duplicate product")
	ErrProductNotFound  = errors.New("product not found")
	ErrEditConflict     = errors.New("edit conflict")
)

type Product struct {
	ID              int64
	Name            string
	Description     string
	Properties      map[string]string
	AvailableAmount int
	Version         int64
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

func (m *ProductModel) InsertProduct(product Product) (int64, error) {
	for _, p := range m.products {
		if p.Name == product.Name {
			return 0, ErrDuplicateProduct
		}
	}

	product.ID = m.idCounter
	product.Version = 1
	m.products = append(m.products, product)

	m.idCounter++

	return product.ID, nil
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

func (m *ProductModel) UpdateProduct(product Product) error {
	for i, p := range m.products {
		if p.ID == product.ID {
			if p.Version != product.Version {
				return ErrEditConflict
			}
			m.products[i] = product
			m.products[i].Version++
			break
		}
	}

	return nil

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
