package models

import (
	"AhmadAbdelrazik/arbun/internal/domain"

	"errors"
)

var (
	ErrCustomerNotFound  = errors.New("customer not found")
	ErrDuplicateCustomer = errors.New("duplicate customer")
)

type CustomerModel struct {
	customers []domain.Customer
	idCounter int64
}

func newCustomerModel() *CustomerModel {
	return &CustomerModel{
		customers: make([]domain.Customer, 0),
		idCounter: 1,
	}
}

func (m *CustomerModel) InsertCustomer(c domain.Customer) (domain.Customer, error) {
	for _, a := range m.customers {
		if a.Email == c.Email {
			return domain.Customer{}, ErrDuplicateCustomer
		}
	}

	c.ID = m.idCounter
	c.Version = 1
	m.customers = append(m.customers, c)
	m.idCounter++

	return c, nil
}

func (m *CustomerModel) GetCustomerByEmail(email string) (domain.Customer, error) {
	for _, a := range m.customers {
		if a.Email == email {
			return a, nil
		}
	}

	return domain.Customer{}, ErrCustomerNotFound
}

func (m *CustomerModel) GetCustomerByID(id int64) (domain.Customer, error) {
	for _, a := range m.customers {
		if a.ID == id {
			return a, nil
		}
	}

	return domain.Customer{}, ErrCustomerNotFound
}

func (m *CustomerModel) GetAllCustomers() ([]domain.Customer, error) {
	return m.customers, nil
}

func (m *CustomerModel) UpdateCustomer(c domain.Customer) (domain.Customer, error) {
	for i, cc := range m.customers {
		if cc.ID == c.ID {
			if cc.Version != c.Version {
				return domain.Customer{}, ErrEditConflict
			}

			c.Version = cc.Version + 1
			m.customers[i] = c
			return c, nil
		}
	}

	return domain.Customer{}, ErrCustomerNotFound
}

func (m *CustomerModel) DeleteCustomer(id int64) error {
	for i, a := range m.customers {
		if a.ID == id {
			m.customers = append(m.customers[:i], m.customers[i+1:]...)
			return nil
		}
	}
	return ErrCustomerNotFound
}
