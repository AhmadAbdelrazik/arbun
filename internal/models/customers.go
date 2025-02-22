package models

import (
	"AhmadAbdelrazik/arbun/internal/validator"
	"errors"
)

type Customer struct {
	ID       int64
	Email    string
	Password Password
	FullName string
	Version  int
}

func (c Customer) Validate() *validator.Validator {
	v := validator.New()
	v.Check(c.FullName != "", "full_name", "must not be empty")
	v.Check(len(c.FullName) <= 40, "full_name", "must not be more than 40")

	v.Check(c.Email != "", "email", "must not be empty")
	v.Check(v.Matches(c.Email, *validator.EmailRX), "email", "must be a valid email address")

	if c.Password.plaintext != nil {
		password := *c.Password.plaintext
		v.Check(password != "", "password", "must not be empty")
		v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
		v.Check(len(password) <= 72, "password", "must be more than 72 bytes long")
	}

	if c.Password.hash == nil {
		panic("missing password hash for user")
	}

	return v.Err()
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
