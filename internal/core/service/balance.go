package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aldotp/ecommerce-go-api/internal/adapter/dto"
	"github.com/aldotp/ecommerce-go-api/internal/core/port"
	"github.com/aldotp/ecommerce-go-api/pkg/consts"
)

type BalanceService struct {
	repo  port.BalanceRepository
	redis port.CacheInterface
}

func NewBalanceService(repo port.BalanceRepository, redis port.CacheInterface) *BalanceService {
	return &BalanceService{repo: repo, redis: redis}
}

func (bs *BalanceService) Withdraw(ctx context.Context, userID uint64, amount float64) error {
	lockKey := fmt.Sprintf("balance_lock:%d", userID)
	lockTTL := 5 * time.Second

	acquired, err := bs.redis.AcquireLock(ctx, lockKey, lockTTL)
	if err != nil {
		return err
	}
	if !acquired {
		return errors.New("another transaction is in progress, try again later")
	}
	defer bs.redis.ReleaseLock(ctx, lockKey)

	return bs.repo.Withdraw(ctx, userID, amount)
}

func (bs *BalanceService) Deposit(ctx context.Context, userID uint64, amount float64) (*dto.DepositResponse, error) {
	lockKey := fmt.Sprintf("balance_lock:%d", userID)
	lockTTL := 5 * time.Second

	acquired, err := bs.redis.AcquireLock(ctx, lockKey, lockTTL)
	if err != nil {
		return nil, err
	}
	if !acquired {
		return nil, errors.New("another transaction is in progress, try again later")
	}
	defer bs.redis.ReleaseLock(ctx, lockKey)

	err = bs.repo.Deposit(ctx, userID, amount)
	if err != nil {
		return nil, err
	}

	amount, err = bs.repo.GetBalance(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &dto.DepositResponse{Balance: amount}, nil
}
func (bs *BalanceService) Transfer(ctx context.Context, fromUserID, toUserID uint64, amount float64) (*dto.TransferResponse, error) {

	// Validate input parameters
	if amount <= 0 {
		return nil, errors.New("transfer amount must be positive")
	}

	if fromUserID == toUserID {
		return nil, consts.ErrCannotSendBalanceSameAccount
	}

	// Acquire locks in a consistent order to prevent deadlocks
	// Always lock the lower ID first to prevent potential deadlock scenarios
	firstLockID, secondLockID := fromUserID, toUserID
	if toUserID < fromUserID {
		firstLockID, secondLockID = toUserID, fromUserID
	}

	firstLockKey := fmt.Sprintf("balance_lock:%d", firstLockID)
	secondLockKey := fmt.Sprintf("balance_lock:%d", secondLockID)
	lockTTL := 5 * time.Second

	// Acquire first lock
	firstAcquired, err := bs.redis.AcquireLock(ctx, firstLockKey, lockTTL)
	if err != nil || !firstAcquired {
		return nil, errors.New("unable to lock first account, try again later")
	}
	defer bs.redis.ReleaseLock(ctx, firstLockKey)

	// Acquire second lock
	secondAcquired, err := bs.redis.AcquireLock(ctx, secondLockKey, lockTTL)
	if err != nil || !secondAcquired {
		return nil, errors.New("unable to lock second account, try again later")
	}
	defer bs.redis.ReleaseLock(ctx, secondLockKey)

	// Now perform the actual transfer
	fromBalance, toBalance, err := bs.repo.Transfer(ctx, fromUserID, toUserID, amount)
	if err != nil {
		return nil, err
	}

	response := &dto.TransferResponse{
		From: struct {
			UserID  uint64  `json:"user_id"`
			Balance float64 `json:"balance"`
		}{UserID: fromUserID, Balance: fromBalance.Balance},

		To: struct {
			UserID  uint64  `json:"user_id"`
			Balance float64 `json:"balance"`
		}{UserID: toUserID, Balance: toBalance.Balance},
	}

	return response, nil
}

func (bs *BalanceService) CheckBalance(ctx context.Context, userID uint64) (dto.BalanceResponse, error) {
	balance, err := bs.repo.GetBalance(ctx, userID)
	if err != nil {
		return dto.BalanceResponse{}, err
	}

	return dto.BalanceResponse{Balance: balance}, nil
}
