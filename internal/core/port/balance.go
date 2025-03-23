package port

import (
	"context"

	"github.com/aldotp/ecommerce-go-api/internal/adapter/dto"
	"github.com/aldotp/ecommerce-go-api/internal/core/domain"
)

type BalanceRepository interface {
	GetBalance(ctx context.Context, userID uint64) (float64, error)
	Deposit(ctx context.Context, userID uint64, amount float64) error
	Withdraw(ctx context.Context, userID uint64, amount float64) error
	// Transfer(ctx context.Context, fromUserID, toUserID uint64, amount float64) error
	Transfer(ctx context.Context, fromUserID, toUserID uint64, amount float64) (*domain.Balance, *domain.Balance, error)
	Store(ctx context.Context, data *domain.Balance) error
}

type BalanceService interface {
	Withdraw(ctx context.Context, userID uint64, amount float64) error
	// Deposit(ctx context.Context, userID uint64, amount float64) error
	Deposit(ctx context.Context, userID uint64, amount float64) (*dto.DepositResponse, error)
	// Transfer(ctx context.Context, fromUserID, toUserID uint64, amount float64) error
	Transfer(ctx context.Context, fromUserID, toUserID uint64, amount float64) (*dto.TransferResponse, error)
	CheckBalance(ctx context.Context, userID uint64) (dto.BalanceResponse, error)
}
