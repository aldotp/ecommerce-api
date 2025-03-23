package port

import (
	"context"

	"github.com/aldotp/ecommerce-go-api/internal/adapter/dto"
	"github.com/aldotp/ecommerce-go-api/internal/core/domain"
)

type ProductRepository interface {
	Finds(ctx context.Context, filter map[string]interface{}) (response []domain.Product, err error)
	FindOne(ctx context.Context, id int) (response *domain.Product, err error)
	Store(ctx context.Context, data *domain.Product) error
	Update(ctx context.Context, id int, updatedData domain.Product) error
	Delete(ctx context.Context, id int) error
	UpdateStock(ctx context.Context, id, newStock int) error
}

type ProductService interface {
	FindOne(ctx context.Context, productID int) (response *domain.Product, err error)
	Store(ctx context.Context, data dto.ProductRequest) error
	Finds(ctx context.Context, param dto.ListProductRequest) (response []domain.Product, err error)
	Update(ctx context.Context, id int, data dto.ProductRequest) error
	Delete(ctx context.Context, id int) error
}
