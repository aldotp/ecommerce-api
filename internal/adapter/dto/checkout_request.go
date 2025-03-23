package dto

type CheckoutRequest struct {
	PaymentMethod string `json:"payment_method"`
}

type CheckoutResponse struct {
	OrderID       int    `json:"order_id"`
	PaymentMethod string `json:"payment_method"`
	Total         int    `json:"total"`
}
