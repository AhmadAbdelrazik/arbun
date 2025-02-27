package domain

import (
	"time"
)

const (
	PaymentCash  = "Cash"
	PaymentDebit = "Debit"

	StatusDispatched = "dispatched"
	StatusCompleted  = "completed"
)

type Order struct {
	ID          int64       `json:"id"`
	CustomerID  int64       `json:"customer_id"`
	CreatedAt   time.Time   `json:"created_at"`
	Cart        Cart        `json:"cart"`
	PaymentType string      `json:"payment_type"`
	Address     Address     `json:"address"`
	MobilePhone MobilePhone `json:"mobile_phone"`
	Status      string      `json:"status"`
}
