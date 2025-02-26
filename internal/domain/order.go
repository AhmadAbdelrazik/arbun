package domain

import (
	"time"
)

const (
	PaymentCash  = "Cash"
	PaymentDebit = "Debit"
)

type CustomerInfo struct {
	ID              int64
	DeliveryAddress string
	MobilePhone     string
}

type Order struct {
	ID          int64
	Time        time.Time
	Cart        Cart
	PaymentType string
	Status      string
	Customer    CustomerInfo
}
