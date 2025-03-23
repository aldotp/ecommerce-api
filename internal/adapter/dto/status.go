package dto

type UpdateOrderStatus struct {
	OrderID int    `json:"order_id"`
	Status  string `json:"status"`
}
