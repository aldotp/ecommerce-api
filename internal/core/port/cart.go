package port

import (
	"context"

	"github.com/aldotp/ecommerce-go-api/internal/adapter/dto"
	"github.com/aldotp/ecommerce-go-api/internal/core/domain"
)

type CartRepository interface {
	FindOne(ctx context.Context, id int) (*domain.Cart, error)
	Store(ctx context.Context, data *domain.Cart) error
	Update(ctx context.Context, id int, updatedData domain.Cart) error
	DeleteByUserID(ctx context.Context, userID int) error
	FindByUserID(ctx context.Context, userID int) (*domain.Cart, error)
}

type CartService interface {
	GetCart(ctx context.Context, userID int) (dto.CartResponse, error)
	AddToCart(ctx context.Context, userID int, productID int, quantity int) error
	RemoveFromCart(ctx context.Context, userID int, productID int) error
	UpdateCart(ctx context.Context, userID int, request dto.UpdateCartRequest) error
}
