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
	ID          int64        `json:"id"`
	CreatedAt   time.Time    `json:"created_at"`
	Cart        Cart         `json:"cart"`
	PaymentType string       `json:"payment_type"`
	Status      string       `json:"status"`
	Customer    CustomerInfo `json:"customer"`
}
