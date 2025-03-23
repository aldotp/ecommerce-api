package port

import (
	"context"

	"github.com/aldotp/ecommerce-go-api/internal/adapter/dto"
	"github.com/aldotp/ecommerce-go-api/internal/core/domain"
)

type OrderItemRepository interface {
	Finds(ctx context.Context, filter map[string]interface{}) (response []domain.OrderItem, err error)
	FindOne(ctx context.Context, id int) (response *domain.OrderItem, err error)
	Store(ctx context.Context, data *domain.OrderItem) error
	Update(ctx context.Context, id int, updatedData domain.OrderItem) error
	Delete(ctx context.Context, id int) error
}

type OrderItemService interface {
	FindOne(ctx context.Context, productID int) (response *domain.OrderItem, err error)
	Store(ctx context.Context, data dto.ProductRequest) error
	Finds(ctx context.Context, param dto.ListProductRequest) (response []domain.OrderItem, err error)
	Update(ctx context.Context, id int, data dto.ProductRequest) error
	Delete(ctx context.Context, id int) error
}
