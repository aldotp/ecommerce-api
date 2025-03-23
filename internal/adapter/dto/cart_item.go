package dto

type CartItemResponse struct {
	Name      string  `json:"name"`
	ProductID int     `json:"product_id"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
}
