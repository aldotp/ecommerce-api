package port

import (
	"context"
	"time"

	"github.com/aldotp/ecommerce-go-api/internal/core/domain"
)

type PaymentRepository interface {
	Finds(ctx context.Context, filter map[string]interface{}) (response []domain.Payment, err error)
	FindOne(ctx context.Context, id int) (response *domain.Payment, err error)
	Store(ctx context.Context, data *domain.Payment) error
	Update(ctx context.Context, id int, updatedData *domain.Payment) error
	Delete(ctx context.Context, id int) error
	FindExpiredPayments(ctx context.Context, now time.Time) ([]*domain.Payment, error)
	FindByUserIDandOrderID(ctx context.Context, userID int, orderID int) (*domain.Payment, error)
}

type PaymentService interface {
	MakePayment(ctx context.Context, userID int, orderID int) error
}
