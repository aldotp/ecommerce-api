package dto

type AddCartRequest struct {
	UserID    int `json:"user_id"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type CartResponse struct {
	TotalItems    int                `json:"total_items"`
	TotalProducts int                `json:"total_products"`
	TotalPrice    float64            `json:"total_price"`
	Items         []CartItemResponse `json:"items"`
}

type RemoveCartRequest struct {
	ProductID int `json:"product_id"`
}

type UpdateCartRequest struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}
