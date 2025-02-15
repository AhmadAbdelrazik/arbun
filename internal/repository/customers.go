package repository

import (
	"errors"
)

type Customer struct {
	ID       int64
	Email    string
	Password Password
	FullName string
	Version  int
}

var (
	ErrCustomerNotFound  = errors.New("customer not found")
	ErrDuplicateCustomer = errors.New("duplicate customer")
)

type CustomerModel struct {
	customers []Customer
	idCounter int64
}

func NewCustomerModel() *CustomerModel {
	return &CustomerModel{
		customers: make([]Customer, 0),
		idCounter: 1,
	}
}

func (m *CustomerModel) InsertCustomer(customer Customer) (Customer, error) {
	for _, a := range m.customers {
		if a.Email == customer.Email {
			return Customer{}, ErrDuplicateCustomer
		}
	}

	customer.ID = m.idCounter
	customer.Version = 1
	m.customers = append(m.customers, customer)
	m.idCounter++

	return customer, nil
}

func (m *CustomerModel) GetCustomerByEmail(email string) (Customer, error) {
	for _, a := range m.customers {
		if a.Email == email {
			return a, nil
		}
	}

	return Customer{}, ErrCustomerNotFound
}

func (m *CustomerModel) GetCustomerByID(id int64) (Customer, error) {
	for _, a := range m.customers {
		if a.ID == id {
			return a, nil
		}
	}

	return Customer{}, ErrCustomerNotFound
}

func (m *CustomerModel) GetAllCustomers() ([]Customer, error) {
	return m.customers, nil
}

func (m *CustomerModel) UpdateCustomer(customer Customer) (Customer, error) {
	for i, a := range m.customers {
		if a.ID == customer.ID {
			if a.Version != customer.Version {
				return Customer{}, ErrEditConflict
			}

			customer.Version = a.Version + 1
			m.customers[i] = customer
			return customer, nil
		}
	}

	return Customer{}, ErrCustomerNotFound
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
