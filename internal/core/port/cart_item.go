package port

import (
	"context"

	"github.com/aldotp/ecommerce-go-api/internal/core/domain"
)

type CartItemRepository interface {
	FindOne(ctx context.Context, id int) (*domain.CartItem, error)
	Store(ctx context.Context, data *domain.CartItem) error
	Update(ctx context.Context, id int, updatedData domain.CartItem) error
	DeleteByCartID(ctx context.Context, cartID int) error
	Finds(ctx context.Context, filter map[string]interface{}) ([]domain.CartItem, error)
	FindOneByFilters(ctx context.Context, filter map[string]interface{}) (*domain.CartItem, error)
	DeleteByProductID(ctx context.Context, product_id int) error
}
