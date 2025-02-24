package order

import (
	"AhmadAbdelrazik/arbun/internal/domain/cart"
	"time"
)

const (
	PaymentCash  = "Cash"
	PaymentDebit = "Debit"
)

type Customer struct {
	ID              int64
	DeliveryAddress string
	MobilePhone     string
}

type Order struct {
	ID          int64
	Time        time.Time
	Cart        cart.Cart
	PaymentType string
	Status      string
	Customer    Customer
}
