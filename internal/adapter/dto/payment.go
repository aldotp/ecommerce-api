package dto

type PaymentRequest struct {
	OrderID int `json:"order_id" binding:"required"`
}
