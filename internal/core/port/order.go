package port

import (
	"context"

	"github.com/aldotp/ecommerce-go-api/internal/core/domain"
)

type OrderRepository interface {
	// Finds(ctx context.Context, filter map[string]interface{}) (response []domain.Order, err error)
	FindOne(ctx context.Context, id int, userID int) (*domain.Order, error)
	Store(ctx context.Context, data *domain.Order) error
	Update(ctx context.Context, id int, updatedData *domain.Order) error
	Delete(ctx context.Context, id int) error
	Finds(ctx context.Context, filter map[string]interface{}) ([]domain.Order, error)
}

type OrderService interface {
	UpdateStatusOrder(ctx context.Context, orderID int, status string) error
	ListOrders(ctx context.Context, userId int) ([]domain.Order, error)
	GetOrder(ctx context.Context, orderID int, userID int) (*domain.Order, error)
}
