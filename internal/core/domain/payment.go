package domain

import "time"

type Payment struct {
	ID            int       `json:"id"`
	OrderID       int       `json:"order_id"`
	PaymentMethod string    `json:"payment_method"`
	PaymentStatus string    `json:"payment_status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	ExpiredAt     time.Time `json:"expired_at"`
}
