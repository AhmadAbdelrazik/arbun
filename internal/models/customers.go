package models

import (
	"AhmadAbdelrazik/arbun/internal/domain/customer"

	"errors"
)

var (
	ErrCustomerNotFound  = errors.New("customer not found")
	ErrDuplicateCustomer = errors.New("duplicate customer")
)

type CustomerModel struct {
	customers []customer.Customer
	idCounter int64
}

func newCustomerModel() *CustomerModel {
	return &CustomerModel{
		customers: make([]customer.Customer, 0),
		idCounter: 1,
	}
}

func (m *CustomerModel) InsertCustomer(c customer.Customer) (customer.Customer, error) {
	for _, a := range m.customers {
		if a.Email == c.Email {
			return customer.Customer{}, ErrDuplicateCustomer
		}
	}

	c.ID = m.idCounter
	c.Version = 1
	m.customers = append(m.customers, c)
	m.idCounter++

	return c, nil
}

func (m *CustomerModel) GetCustomerByEmail(email string) (customer.Customer, error) {
	for _, a := range m.customers {
		if a.Email == email {
			return a, nil
		}
	}

	return customer.Customer{}, ErrCustomerNotFound
}

func (m *CustomerModel) GetCustomerByID(id int64) (customer.Customer, error) {
	for _, a := range m.customers {
		if a.ID == id {
			return a, nil
		}
	}

	return customer.Customer{}, ErrCustomerNotFound
}

func (m *CustomerModel) GetAllCustomers() ([]customer.Customer, error) {
	return m.customers, nil
}

func (m *CustomerModel) UpdateCustomer(c customer.Customer) (customer.Customer, error) {
	for i, cc := range m.customers {
		if cc.ID == c.ID {
			if cc.Version != c.Version {
				return customer.Customer{}, ErrEditConflict
			}

			c.Version = cc.Version + 1
			m.customers[i] = c
			return c, nil
		}
	}

	return customer.Customer{}, ErrCustomerNotFound
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
