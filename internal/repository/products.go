package repository

import "errors"

var (
	ErrDuplicateProduct = errors.New("duplicate product")
	ErrProductNotFound  = errors.New("product not found")
	ErrEditConflict     = errors.New("edit conflict")
)

type Product struct {
	ID          int64
	Name        string
	Description string
	Properties  map[string]string
	Version     int64
}

type ProductModel struct {
	Products  []Product
	IDCounter int64
}

func (m *ProductModel) InsertProduct(product Product) (int64, error) {
	for _, p := range m.Products {
		if p.Name == product.Name {
			return 0, ErrDuplicateProduct
		}
	}

	product.ID = m.IDCounter
	product.Version = 1
	m.Products = append(m.Products, product)

	m.IDCounter++

	return product.ID, nil
}

func (m *ProductModel) GetProductByID(id int64) (Product, error) {
	for _, p := range m.Products {
		if p.ID == id {
			return p, nil
		}
	}

	return Product{}, ErrProductNotFound
}

func (m *ProductModel) GetAllProducts() ([]Product, error) {
	return m.Products, nil
}

func (m *ProductModel) UpdateProduct(product Product) error {
	for i, p := range m.Products {
		if p.ID == product.ID {
			if p.Version != product.Version {
				return ErrEditConflict
			}
			m.Products[i] = product
			m.Products[i].Version++
			break
		}
	}

	return nil

}

func (m *ProductModel) DeleteProduct(id int64) error {
	for i, p := range m.Products {
		if p.ID == id {
			products := m.Products[:i]
			products = append(products, m.Products[i+1:]...)
			m.Products = products
			return nil
		}
	}
	return ErrProductNotFound
}
