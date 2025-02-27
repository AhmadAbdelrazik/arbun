package models

import (
	"AhmadAbdelrazik/arbun/internal/domain"
	"errors"
	"slices"
)

var (
	ErrDuplicateProduct          = errors.New("duplicate product")
	ErrProductNotFound           = errors.New("product not found")
	ErrEditConflict              = errors.New("edit conflict")
	ErrInsufficientProductAmount = errors.New("insufficient amount")
)

type ProductModel struct {
	products  []domain.Product
	idCounter int64
}

func newProductModel() *ProductModel {
	return &ProductModel{
		products:  make([]domain.Product, 0),
		idCounter: 1,
	}
}

func (m *ProductModel) InsertProduct(p domain.Product) (domain.Product, error) {
	for _, pp := range m.products {
		if pp.Name == p.Name && pp.Vendor == p.Vendor {
			return domain.Product{}, ErrDuplicateProduct
		}
	}

	p.ID = m.idCounter
	p.Version = 1
	m.products = append(m.products, p)

	m.idCounter++

	return p, nil
}

func (m *ProductModel) GetProductByID(id int64) (domain.Product, error) {
	for _, pp := range m.products {
		if pp.ID == id {
			return pp, nil
		}
	}

	return domain.Product{}, ErrProductNotFound
}

func (m *ProductModel) GetAllProducts() ([]domain.Product, error) {
	return m.products, nil
}

func (m *ProductModel) ChangeProductAmountBy(productID int64, amount int) error {
	product, err := m.GetProductByID(productID)
	if err != nil {
		return err
	}

	return m.ChangeProductAmountTo(productID, product.AvailableAmount+amount)
}

func (m *ProductModel) ChangeProductAmountTo(productID int64, amount int) error {
	if amount < 0 {
		return ErrInsufficientProductAmount
	}
	for i := range m.products {
		if m.products[i].ID == productID {
			m.products[i].AvailableAmount = amount
			return nil
		}
	}
	return ErrProductNotFound
}

func (m *ProductModel) UpdateProduct(p domain.Product) (domain.Product, error) {
	var result domain.Product

	for i, pp := range m.products {
		if pp.ID == p.ID {
			if pp.Version != p.Version {
				return domain.Product{}, ErrEditConflict
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
