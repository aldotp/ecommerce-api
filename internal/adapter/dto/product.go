package dto

import (
	"time"

	"github.com/aldotp/ecommerce-go-api/internal/core/domain"
)

type ProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	Stock       int     `json:"stock" binding:"required"`
	CategoryID  int     `json:"category_id" binding:"required"`
}

type GetProductRequest struct {
	ID int `uri:"id" binding:"required"`
}

type ListProductRequest struct {
	Page     int    `form:"page"`
	PageSize uint64 `form:"pageSize"`
	Search   string `form:"search"`
}

func NewProductResponse(user *domain.Product) ProductResponse {
	return ProductResponse{
		ID:        user.ID,
		Name:      user.Name,
		Price:     user.Price,
		Stock:     user.Stock,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func NewProductsResponse(products []domain.Product) []ProductResponse {
	var res []ProductResponse
	for _, product := range products {
		res = append(res, NewProductResponse(&product))
	}
	return res
}

type ProductResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ParamProductRequest struct {
	ID int `uri:"id" binding:"required"`
}
