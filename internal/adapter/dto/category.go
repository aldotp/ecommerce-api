package dto

import (
	"github.com/aldotp/ecommerce-go-api/internal/core/domain"
)

type CategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

type GetCategoryRequest struct {
	ID int `uri:"id" binding:"required"`
}

type ListCategoryRequest struct {
	Page     int    `form:"page"`
	PageSize uint64 `form:"pageSize"`
}

func NewCategoryResponse(category *domain.Category) CategoryResponse {
	return CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}
}

func NewCategoryResponses(products []domain.Category) []CategoryResponse {
	var res []CategoryResponse
	for _, product := range products {
		res = append(res, NewCategoryResponse(&product))
	}
	return res
}

type CategoryResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ParamCategoryRequest struct {
	ID int `uri:"id" binding:"required"`
}
