package port

import (
	"context"
)

type CheckoutService interface {
	Checkout(ctx context.Context, userID int) error
}
