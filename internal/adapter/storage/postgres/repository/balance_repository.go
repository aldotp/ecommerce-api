package repository

import (
	"context"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/aldotp/ecommerce-go-api/internal/adapter/storage/postgres"
	"github.com/aldotp/ecommerce-go-api/internal/core/domain"
	"github.com/aldotp/ecommerce-go-api/pkg/consts"

	"github.com/jackc/pgx/v5"
)

type BalanceRepository struct {
	db        *postgres.DB
	TableName string
}

func NewBalanceRepository(db *postgres.DB) *BalanceRepository {
	return &BalanceRepository{
		db:        db,
		TableName: "balances",
	}
}

func (br *BalanceRepository) Withdraw(ctx context.Context, userID uint64, amount float64) error {

	tx, err := br.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := sq.Select("balance").
		From(br.TableName).
		Where(sq.Eq{"user_id": userID}).
		Suffix("FOR UPDATE").PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	var balance float64
	err = tx.QueryRow(ctx, sql, args...).Scan(&balance)
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.New("balance record not found")
		}
		return err
	}

	if balance < amount {
		return consts.ErrInsufficientBalance
	}

	newBalance := balance - amount
	updateQuery := sq.Update(br.TableName).
		Set("balance", newBalance).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"user_id": userID}).PlaceholderFormat(sq.Dollar)

	sql, args, err = updateQuery.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (br *BalanceRepository) Deposit(ctx context.Context, userID uint64, amount float64) error {
	tx, err := br.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := sq.Select("balance").
		From(br.TableName).
		Where(sq.Eq{"user_id": userID}).
		Suffix("FOR UPDATE").PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	var balance float64
	err = tx.QueryRow(ctx, sql, args...).Scan(&balance)
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.New("balance record not found")
		}
		return err
	}

	// Update balance
	newBalance := balance + amount
	updateQuery := sq.Update(br.TableName).
		Set("balance", newBalance).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"user_id": userID}).PlaceholderFormat(sq.Dollar)

	sql, args, err = updateQuery.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (br *BalanceRepository) GetBalance(ctx context.Context, userID uint64) (float64, error) {
	var balance float64
	query := br.db.QueryBuilder.Select("balance").
		From(br.TableName).
		Where(sq.Eq{"user_id": userID})

	sql, args, err := query.ToSql()
	if err != nil {
		return 0, err
	}

	err = br.db.QueryRow(ctx, sql, args...).Scan(&balance)
	if err != nil {
		return 0, err
	}

	return balance, nil
}
func (br *BalanceRepository) Transfer(ctx context.Context, fromUserID, toUserID uint64, amount float64) (*domain.Balance, *domain.Balance, error) {
	// Ensure sender and receiver are different
	if fromUserID == toUserID {
		return nil, nil, errors.New("cannot transfer to the same account")
	}

	// Start database transaction
	tx, err := br.db.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback(ctx)

	// Fetch sender's balance with row lock
	var senderBalance float64
	senderQuery := sq.Select("balance").
		From(br.TableName).
		Where(sq.Eq{"user_id": fromUserID}).
		Suffix("FOR UPDATE").PlaceholderFormat(sq.Dollar)

	sql, args, err := senderQuery.ToSql()
	if err != nil {
		return nil, nil, err
	}

	err = tx.QueryRow(ctx, sql, args...).Scan(&senderBalance)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil, errors.New("sender balance record not found")
		}
		return nil, nil, err
	}

	// Check if sender has enough balance
	if senderBalance < amount {
		return nil, nil, consts.ErrInsufficientBalance
	}

	// Fetch receiver's balance with row lock
	var receiverBalance float64
	receiverQuery := sq.Select("balance").
		From(br.TableName).
		Where(sq.Eq{"user_id": toUserID}).
		Suffix("FOR UPDATE").PlaceholderFormat(sq.Dollar)

	sql, args, err = receiverQuery.ToSql()
	if err != nil {
		return nil, nil, err
	}

	err = tx.QueryRow(ctx, sql, args...).Scan(&receiverBalance)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil, errors.New("receiver balance record not found")
		}
		return nil, nil, err
	}

	// Update sender's balance
	updateSenderQuery := sq.Update(br.TableName).
		Set("balance", senderBalance-amount).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"user_id": fromUserID}).PlaceholderFormat(sq.Dollar)

	sql, args, err = updateSenderQuery.ToSql()
	if err != nil {
		return nil, nil, err
	}
	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return nil, nil, err
	}

	// Update receiver's balance
	updateReceiverQuery := sq.Update(br.TableName).
		Set("balance", receiverBalance+amount).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"user_id": toUserID}).PlaceholderFormat(sq.Dollar)

	sql, args, err = updateReceiverQuery.ToSql()
	if err != nil {
		return nil, nil, err
	}
	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return nil, nil, err
	}

	// Fetch updated balances
	var updatedSender, updatedReceiver domain.Balance

	err = tx.QueryRow(ctx, "SELECT id, user_id, balance, created_at, updated_at FROM balances WHERE user_id = $1", fromUserID).
		Scan(&updatedSender.ID, &updatedSender.UserID, &updatedSender.Balance, &updatedSender.CreatedAt, &updatedSender.UpdatedAt)
	if err != nil {
		return nil, nil, err
	}

	err = tx.QueryRow(ctx, "SELECT id, user_id, balance, created_at, updated_at FROM balances WHERE user_id = $1", toUserID).
		Scan(&updatedReceiver.ID, &updatedReceiver.UserID, &updatedReceiver.Balance, &updatedReceiver.CreatedAt, &updatedReceiver.UpdatedAt)
	if err != nil {
		return nil, nil, err
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, nil, err
	}

	return &updatedSender, &updatedReceiver, nil
}

func (br *BalanceRepository) Store(ctx context.Context, data *domain.Balance) error {
	query := br.db.QueryBuilder.Insert(br.TableName).
		Columns("user_id", "balance", "created_at", "updated_at").
		Values(data.UserID, data.Balance, data.CreatedAt, data.UpdatedAt).
		Suffix("RETURNING id, user_id, balance, created_at, updated_at")

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	err = br.db.QueryRow(ctx, sql, args...).Scan(
		&data.ID,
		&data.UserID,
		&data.Balance,
		&data.CreatedAt,
		&data.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
