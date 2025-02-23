package order

import "time"

type Order struct {
	Time        time.Time
	PaymentType string
}
