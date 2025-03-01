package domain

import (
	"AhmadAbdelrazik/arbun/internal/pkg/validator"
	"time"
)

const (
	PaymentCash   PaymentMethod = "Cash"
	PaymentCredit               = "Credit"

	StatusDispatched OrderStatus = "dispatched"
	StatusCompleted              = "completed"
	StatusCanceled               = "canceled"
)

type Order struct {
	ID          int64         `json:"id"`
	CustomerID  int64         `json:"customer_id"`
	CreatedAt   time.Time     `json:"created_at"`
	Cart        Cart          `json:"cart"`
	PaymentType PaymentMethod `json:"payment_type"`
	Address     Address       `json:"address"`
	MobilePhone MobilePhone   `json:"mobile_phone"`
	Status      OrderStatus   `json:"status"`
}

func (o Order) Validate() *validator.Validator {
	v := validator.New()

	v.Add(o.Address.Validate())
	v.Add(o.MobilePhone.Validate())
	v.Add(o.Status.Validate())
	v.Add(o.PaymentType.Validate())

	return v.Err()
}

type PaymentMethod string

func (p PaymentMethod) Validate() *validator.Validator {
	v := validator.New()

	acceptedPayment := []PaymentMethod{
		PaymentCash,
		PaymentCredit,
	}

	v.Check(validator.In(p, acceptedPayment...), "payment", "invalid payment type")

	return v.Err()
}

type OrderStatus string

func (o OrderStatus) Validate() *validator.Validator {
	v := validator.New()

	acceptedStatuses := []OrderStatus{
		StatusDispatched,
		StatusCompleted,
		StatusCanceled,
	}

	v.Check(validator.In(o, acceptedStatuses...), "status", "invalid order status")

	return v.Err()
}

type CreditCard struct {
	Name       string `json:"name"`
	Number     int    `json:"number"`
	CCV        int    `json:"ccv"`
	ExpiryDate string `json:"expiry_date"`
}
